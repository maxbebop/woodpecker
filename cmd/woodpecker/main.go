package main

import (
	config "woodpecker/configs"
	chatservice "woodpecker/internal/services/chat"
	slackservice "woodpecker/internal/services/slack"

	"github.com/powerman/structlog"
)

func main() {
	log := structlog.New()

	log.Info("start task management bot - woodpecker")

	cfg := config.New("slack.config.yml")

	chatBot := slackservice.New(cfg, log)
	err := chatservice.StartChat(chatBot, log)

	if err != nil {
		log.PrintErr(err)
	}
}
