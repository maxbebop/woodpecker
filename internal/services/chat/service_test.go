package chatservice_test

import (
	context "context"
	"testing"

	config "woodpecker/internal/configs"

	chatservice "woodpecker/internal/services/chat"
	slackservice "woodpecker/internal/services/slack"

	structlog "github.com/powerman/structlog"
	"github.com/stretchr/testify/require"
)

// this is not just a comment but the way to instruct `go generate` command
// you can run `go generate ./...` in the project root to process all the files across the project
//go:generate mockery --all --testonly --outpkg chatservice_test --output .

func TestProcessMsg(t *testing.T) {
	mockChatBot := NewChatBot(t)
	outMsg := chatservice.OutMessage{Message: chatservice.Message{User: "test", Channel: "test", Error: nil}, Type: chatservice.Common, Pretext: "", Error: nil}
	mockChatBot.On("SendMessage", outMsg).Return(nil)
	err := mockChatBot.SendMessage(outMsg)

	require.NoError(t, err, "chatservice SendMessage")
}

func TestStartChat(t *testing.T) {
	t.Parallel()

	log := structlog.New()
	cfg := config.New("../../../slack.config.yml")
	chatBot := slackservice.New(cfg, log)

	require.NotNil(t, chatBot, "chatBot")

	mockChatBot := NewChatBot(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chatChannel := make(chan chatservice.Message)

	mockChatBot.On("GetMessagesLoop", ctx, chatChannel, log).Return()
	mockChatBot.GetMessagesLoop(ctx, chatChannel, log)

	mockChatBot.On("Run").Return(nil)
	mockChatBotErr := mockChatBot.Run()
	ctx.Done()

	require.NoError(t, mockChatBotErr, "chatservice Run")
}
