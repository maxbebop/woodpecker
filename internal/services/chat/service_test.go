package chatservice_test

import (
	"context"
	"testing"

	chatservice "woodpecker/internal/services/chat"
)

// this is not just a comment but the way to instruct `go generate` command
// you can run `go generate ./...` in the project root to process all the files across the project
//go:generate mockery --all --testonly --outpkg chatservice_test --output .

func TestStartChat(t *testing.T) {
	// ToDo: fix the test to cover at leas something :)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// log := structlog.New()
	inMsgChannel := make(chan chatservice.Message)
	chatBot := &ChatBot{}
	chatBot.On("GetMessagesLoop", ctx, inMsgChannel, nil)
	chatBot.On("SendMessage", chatservice.OutMessage{
		Message: chatservice.Message{
			User:    "",
			Channel: "",
			Text:    "",
			Error:   nil,
		},
		Type:    chatservice.Common,
		Pretext: "",
		Error:   nil})
	close(inMsgChannel)
}
