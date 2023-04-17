package main

import (
	"context"

	"github.com/powerman/structlog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	"github.com/maxbebop/woodpecker/internal/configs"
	"github.com/maxbebop/woodpecker/pkg/slackclient"
	"github.com/maxbebop/woodpecker/pkg/socketmodewrap"
)

func main() {
	log := structlog.New()

	log.Info("start task managment bot - woodpecker")

	cfg := config.New("config.yml")

	slackClient := slack.New(
		cfg.Slack.OAuthToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(cfg.Slack.AppToken),
	)

	socketClient := socketmode.New(
		slackClient,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(socketmodewrap.Log{Logger: log}),
	)

	client := slackclient.New(slackClient, socketmodewrap.New(socketClient), cfg.Slack.AppUserId)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chatChannel := make(chan slackclient.Message)

	go client.GetMessagesLoop(ctx, chatChannel, log)
	go processMsgLoop(ctx, client, chatChannel, log)

	if err := client.Run(); err != nil {
		log.Fatal("client run", "err", err)
	}
}

func processMsgLoop(
	ctx context.Context,
	slackClient *slackclient.Client,
	in <-chan slackclient.Message,
	log *structlog.Logger,
) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down processing loop")
			return
		case msg := <-in:
			processMsg(slackClient, msg, log)
		}
	}

}

func processMsg(slackClient *slackclient.Client, msg slackclient.Message, log *structlog.Logger) {
	if msg.Error != nil {
		log.Fatal(msg.Error)
	}

	log.Debug("msg", "from", msg.User, "text", msg.Text)

	if err := slackClient.SendMessage(msg); err != nil {
		log.Err("send message", "err", err)
	}
}
