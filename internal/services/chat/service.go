package chatservice

import (
	"context"
	"fmt"
	"strings"

	models "woodpecker/internal/models/user"
	storage "woodpecker/internal/storage/pudge"
	usertaskmanager "woodpecker/internal/userTaskManager"

	"github.com/powerman/structlog"
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

type (
	ChatService interface {
		StartChat(chatBot ChatBot, log *structlog.Logger) error
	}

	chatService struct {
		utmStorage  storage.Client[usertaskmanager.UserTaskManager]
		userStorage storage.Client[models.User]
	}
)

type ChatBot interface {
	GetMessagesLoop(ctx context.Context, inMsgChannel chan Message, log *structlog.Logger)
	SendMessage(msg OutMessage) error
	Run() error
}

func New(utmStorage storage.Client[usertaskmanager.UserTaskManager], userStorage storage.Client[models.User]) ChatService {
	c := &chatService{
		utmStorage:  utmStorage,
		userStorage: userStorage,
	}

	return c
}
func (s *chatService) StartChat(chatBot ChatBot, log *structlog.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chatChannel := make(chan Message)

	s.utmStorage.DebugAllValues()
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

	isUserNew := s.isNewUser(msg.User, log)
	if isUserNew {
		s.saveNewUser(msg.User, log)
	}
	user, err := s.userStorage.Get(msg.User)
	if err != nil {
		log.Err("get user from db", "err", err) //nolint:errcheck // intentional
	}

	//outMsg := OutMessage{Message: msg, Type: Common, Pretext: "", Error: nil}
	outMsg := s.createOutMessageByStatus(user, msg, log)

	if !outMsg.Empty {
		if err := chatBot.SendMessage(outMsg); err != nil {
			log.Err("send message", "err", err) //nolint:errcheck // intentional
		}
	}

}

func (s *chatService) initState(userChatToken string, chatBot ChatBot, log *structlog.Logger) error {

	var userTm usertaskmanager.UserTaskManager
	if s.isNewUser(userChatToken, log) {
		user := models.User{ChatToken: userChatToken}
		userTm = usertaskmanager.New()
		if err := userTm.AddUser(user); err != nil {
			return err
		}
	} else {
		user, err := s.userStorage.Get(userChatToken)
		if err != nil {
			return err
		}
		userTm = usertaskmanager.NewByUser(user)
	}

	if err := s.utmStorage.Set(userChatToken, userTm); err != nil {
		return err
	}

	return nil
}

func (s *chatService) createOutMessageByStatus(user models.User, msg Message, log *structlog.Logger) OutMessage {
	if !s.hasTMSToken(user) && !s.isMsgHasTMSToken(msg) {
		text := s.createRegistrationMsg(log)
		msg.Text = text
		return OutMessage{Message: msg, Type: Common, Pretext: "", Error: nil}
	} else if !s.hasTMSToken(user) && s.isMsgHasTMSToken(msg) {
		user.TMSToken = strings.ReplaceAll(msg.Text, "token:", "")
		s.saveUser(user, log)
		msg.Text = "you have successfully registered!"
		return OutMessage{Message: msg, Type: Common, Pretext: "", Error: nil}
	}

	msg.Text += fmt.Sprintf(" <db: %v>", user)
	return OutMessage{Message: msg, Type: Common, Pretext: "", Error: nil}
}

func (s *chatService) isMsgHasTMSToken(msg Message) bool {
	return strings.Contains(msg.Text, "token:")
}
func (s *chatService) createRegistrationMsg(log *structlog.Logger) string {
	return "Hello! Send me your tsm token as string token:you_token"
}

func (s *chatService) hasTMSToken(user models.User) bool {
	return len(user.TMSToken) > 0
}
func (s *chatService) isNewUser(user string, log *structlog.Logger) bool {
	flag, err := s.userStorage.Has(user)
	if err != nil {
		log.Err(err)
	}

	return !flag
}

func (s *chatService) saveNewUser(userChatToken string, log *structlog.Logger) {
	if err := s.userStorage.Set(userChatToken, models.User{ChatToken: userChatToken}); err != nil {
		log.Err("save new user error", "user", userChatToken, "error", err)
	}
}

func (s *chatService) saveUser(user models.User, log *structlog.Logger) {
	if err := s.userStorage.Set(user.ChatToken, user); err != nil {
		log.Err("save user error", "user", user, "error", err)
	}
}
