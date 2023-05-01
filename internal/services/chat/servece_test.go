package chatservice_test

import (
	"context"
	"testing"
	chatservice "woodpecker/internal/services/chat"
	"woodpecker/mocks"

	"github.com/powerman/structlog"
)

func TestStartChat(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := structlog.New()
	inMsgChannel := make(chan chatservice.Message)
	chatBot := &mocks.ChatBot{}
	chatBot.On("GetMessagesLoop", ctx, inMsgChannel, log)
	outMsg := chatservice.OutMessage{Message: chatservice.Message{}}
	outMsg.Type = chatservice.Common
	chatBot.On("SendMessage", outMsg)
	close(inMsgChannel)
}
