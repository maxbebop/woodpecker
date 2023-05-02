package chatservice

import (
	"context"

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

func StartChat(chatBot ChatBot, log *structlog.Logger) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chatChannel := make(chan Message)

	go chatBot.GetMessagesLoop(ctx, chatChannel, log)
	go processMsgLoop(ctx, chatBot, chatChannel, log)

	if err := chatBot.Run(); err != nil {
		return log.Err("client run", "err", err)
	}

	return nil
}

func processMsgLoop(
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
			processMsg(chatBot, msg, log)
		}
	}
}

func processMsg(chatBot ChatBot, msg Message, log *structlog.Logger) {
	if msg.Error != nil {
		log.Fatal(msg.Error)
	}

	log.Debug("msg", "from", msg.User, "text", msg.Text)

	outMsg := OutMessage{Message: msg}
	outMsg.Type = Common

	if err := chatBot.SendMessage(outMsg); err != nil {
		log.Err("send message", "err", err) //nolint:errcheck // intentional
	}
}
