package main

import (
	"context"

	"github.com/powerman/structlog"
	"github.com/sourcegraph/conc"

	"woodpecker/internal/configs"
	"woodpecker/pkg/slack"
)

func main() {
	log := structlog.New()

	log.Info("start task managment bot - woodpecker")

	config := config.New("config.yml")

	client := slack.New(
		config.Slack.OAuthToken,
		config.Slack.AppToken,
		config.Slack.AppUserId,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg conc.WaitGroup

	chatChannel := make(chan slack.Message)

	wg.Go(func() { client.GetMessagesLoop(ctx, chatChannel, log) })

	wg.Go(func() {
		if err := client.Run(); err != nil {
			log.Err("client run", "err", err)
		}
	})

	wg.Go(func() {
		for message := range chatChannel {
			if message.Error != nil {
				log.Fatal(message.Error)
			}

			log.Debug("msg", "from", message.User, "text", message.Text)

			if err := client.SendMessage(message); err != nil {
				log.Err("send message", "err", err)
			}
		}
	})

	wg.Wait()
}
