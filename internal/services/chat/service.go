package chatservice

import (
	"context"
	"fmt"

	models "woodpecker/internal/models"
	storage "woodpecker/internal/storage/pudge"
	usertaskmanager "woodpecker/internal/userTaskManager"

	"github.com/powerman/structlog"
)

type (
	ChatService interface {
		StartChat(chatBot ChatBot, log *structlog.Logger) error
	}

	chatService struct {
		chatBot     ChatBot
		userStorage storage.Client[models.User]
		urmCache    map[string]*usertaskmanager.UserTaskManager
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

func New(userStorage storage.Client[models.User]) ChatService {
	c := &chatService{
		userStorage: userStorage,
		chatBot:     nil,
		urmCache:    make(map[string]*usertaskmanager.UserTaskManager),
	}

	return c
}
func (s *chatService) StartChat(chatBot ChatBot, log *structlog.Logger) error {
	s.chatBot = chatBot
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.initStateManagrersCache(log)
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

	user, err := s.getUser(msg.Channel, msg.User, log)
	if err != nil {
		log.Err("get user from db", "err", err) //nolint:errcheck // intentional
		return
	}

	stateManager := s.getStameManager(user.ChatToken, log)
	env := models.Environment{
		User:         user,
		ChatChanelId: msg.Channel,
		Msg:          msg.Text,
	}

	if err := stateManager.Compute(env, s); err != nil {
		log.Err("state compute", "err", err) //nolint:errcheck // intentional
		return
	}

	s.setStameManager(user.ChatToken, stateManager, log)
	s.userStorage.DebugAllValues()
}

func (s *chatService) SendMessageByState(user models.User, chatChannel string, msg string, log *structlog.Logger) {
	baseMsg := Message{
		User:    user.ChatToken,
		Channel: chatChannel,
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
		stateManager := s.getStameManager(user.ChatToken, log)
		env := models.Environment{
			User:         user,
			ChatChanelId: user.ChatToken,
		}

		if err := stateManager.Compute(env, s); err != nil {
			return log.Err("state compute", "err", err) //nolint:errcheck // intentional
		}

		s.setStameManager(user.ChatToken, stateManager, log)
	}
	s.userStorage.DebugAllValues()
	return nil
}

func (s *chatService) getUser(chatChanelId string, userChatToken string, log *structlog.Logger) (models.User, error) {
	if !s.userStorage.Has(chatChanelId) {
		return models.User{ChatToken: userChatToken}, nil
	}

	return s.userStorage.Get(chatChanelId)
}

func (s *chatService) getStameManager(userChatToken string, log *structlog.Logger) *usertaskmanager.UserTaskManager {
	utm := s.urmCache[userChatToken]
	if utm == nil {
		utm = usertaskmanager.New(s.userStorage, log)
	}

	return utm
}

func (s *chatService) setStameManager(userChatToken string, stameManager *usertaskmanager.UserTaskManager, log *structlog.Logger) {
	s.urmCache[userChatToken] = stameManager
}
