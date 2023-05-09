package main

import (
	"woodpecker/internal/configs"
	chatservice "woodpecker/internal/services/chat"
	slackservice "woodpecker/internal/services/slack"

	"github.com/powerman/structlog"
)

func main() {
	log := structlog.New()

	log.Info("start task management bot - woodpecker")

	cfg := config.New("slack.config.yml")

	chatBot := slackservice.New(cfg, log)

	if err := chatservice.StartChat(chatBot, log); err != nil {
		log.PrintErr(err)
	}
}
