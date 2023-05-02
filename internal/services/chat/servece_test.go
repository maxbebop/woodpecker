package chatservice_test

import (
	"context"
	"testing"
	chatservice "woodpecker/internal/services/chat"
	"woodpecker/mocks"
)

func TestStartChat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//log := structlog.New()
	inMsgChannel := make(chan chatservice.Message)
	chatBot := &mocks.ChatBot{}
	chatBot.On("GetMessagesLoop", ctx, inMsgChannel, nil)
	chatBot.On("SendMessage", chatservice.OutMessage{
		Message: chatservice.Message{},
		Type:    chatservice.Common})
	close(inMsgChannel)
}
