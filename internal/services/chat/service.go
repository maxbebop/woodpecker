package chatservice

import (
	"context"
	"fmt"

	"woodpecker/internal/models"
	"woodpecker/internal/storage/users"
	"woodpecker/internal/storage/userstatemanagers"
	"woodpecker/internal/userstatemanager"

	"github.com/powerman/structlog"
)

type (
	ChatService interface {
		StartChat(chatBot ChatBot, log *structlog.Logger) error
	}

	chatService struct {
		chatBot     ChatBot
		userStorage users.Client
		utmStorage  userstatemanagers.Client
	}
)

type Message struct {
	User    string
	Channel string
	Text    string
	Error   error
}

type OutMessage struct {
	Message Message
	Pretext string
	Type    MessageType
	Error   error
	Empty   bool
}

type MessageType int

const (
	Common MessageType = iota
	Warning
	Attention
)

type ChatBot interface {
	GetMessagesLoop(ctx context.Context, inMsgChannel chan Message, log *structlog.Logger)
	SendMessage(msg OutMessage) error
	Run() error
}

func New(userStorage users.Client, utmStorage userstatemanagers.Client) ChatService {
	c := &chatService{
		userStorage: userStorage,
		chatBot:     nil,
		utmStorage:  utmStorage,
	}

	return c
}
func (s *chatService) StartChat(chatBot ChatBot, log *structlog.Logger) error {
	s.chatBot = chatBot
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := s.initStateManagrersCache(log); err != nil {
		return err
	}
	chatChannel := make(chan Message)

	go chatBot.GetMessagesLoop(ctx, chatChannel, log)
	go s.processMsgLoop(ctx, chatBot, chatChannel, log)

	if err := chatBot.Run(); err != nil {
		return fmt.Errorf("%w", log.Err("client run", "err", err))
	}

	return nil
}

func (s *chatService) processMsgLoop(
	ctx context.Context,
	chatBot ChatBot,
	in <-chan Message,
	log *structlog.Logger,
) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down processing loop")

			return
		case msg := <-in:
			log.Debug("processMsgLoop", "from", msg.User, "text", msg.Text)
			s.processMsg(chatBot, msg, log)
		}
	}
}

func (s *chatService) processMsg(chatBot ChatBot, msg Message, log *structlog.Logger) {
	if msg.Error != nil {
		log.Fatal(msg.Error)
	}

	log.Debug("msg", "from", msg.User, "text", msg.Text)

	chatChanelId := models.ChatChanelId(msg.Channel)
	user, ok := s.getUser(models.UserMessengerToken(msg.User))
	if !ok {
		log.Err("get user from db") //nolint:errcheck // intentional
		return
	}

	stateManager := s.getStameManager(user.MessengerToken, log)
	env := models.Environment{
		User:         user,
		ChatChanelId: chatChanelId,
		Msg:          msg.Text,
	}

	if err := stateManager.Compute(env, s); err != nil {
		log.Err("state compute", "err", err) //nolint:errcheck // intentional
		return
	}

	s.setStameManager(user.MessengerToken, stateManager, log)
}

func (s *chatService) SendMessageByState(user models.User, messengerToken models.UserMessengerToken, msg string, log *structlog.Logger) {
	baseMsg := Message{
		User:    string(user.MessengerToken),
		Channel: string(user.MessengerToken),
		Text:    msg,
		Error:   nil,
	}
	outMsg := OutMessage{Message: baseMsg, Type: Common, Pretext: "", Error: nil}

	if !outMsg.Empty {
		if err := s.chatBot.SendMessage(outMsg); err != nil {
			log.Err("send message", "err", err) //nolint:errcheck // intentional
		}
	}
}

func (s *chatService) initStateManagrersCache(log *structlog.Logger) error {
	users, err := s.userStorage.GetAllItems()
	if err != nil {
		return err
	}

	for i := range users {
		user := users[i]
		stateManager := s.getStameManager(user.MessengerToken, log)
		env := models.Environment{
			User: user,
		}

		if err := stateManager.Compute(env, s); err != nil {
			return log.Err("state compute", "err", err) //nolint:errcheck // intentional
		}

		if err := s.setStameManager(user.MessengerToken, stateManager, log); err != nil {
			return log.Err("save state manager", "err", err) //nolint:errcheck // intentional
		}
	}
	s.userStorage.DebugAllValues()
	s.utmStorage.DebugAllValues()
	return nil
}

func (s *chatService) hasUser(messengerToken models.UserMessengerToken) bool {
	return s.userStorage.Has(string(messengerToken))
}

func (s *chatService) createUser(messengerToken models.UserMessengerToken) models.User {
	return models.User{MessengerToken: messengerToken}
}

func (s *chatService) getUser(messangerToken models.UserMessengerToken) (models.User, bool) {

	if !s.hasUser(messangerToken) {
		return s.createUser(messangerToken), true
	}
	return s.userStorage.Get(string(messangerToken))
}

func (s *chatService) getStameManager(userToken models.UserMessengerToken, log *structlog.Logger) *userstatemanager.UserStateManager {
	utm, ok := s.utmStorage.Get(string(userToken))
	if !ok {
		utm = userstatemanager.New(s.userStorage, log)
	}

	return utm
}

func (s *chatService) setStameManager(userToken models.UserMessengerToken, stameManager *userstatemanager.UserStateManager, log *structlog.Logger) error {
	log.Debug("setStameManager: %v\n", stameManager.GetCode())
	return s.utmStorage.Set(string(userToken), stameManager)
}
