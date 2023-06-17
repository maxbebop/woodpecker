package main

import (
	config "woodpecker/internal/configs"
	"woodpecker/internal/storage/users"
	"woodpecker/internal/storage/userstatemanagers"

	chatservice "woodpecker/internal/services/chat"
	slackservice "woodpecker/internal/services/slack"

	"github.com/powerman/structlog"
)

func main() {
	log := structlog.New()

	log.Info("start task management bot - woodpecker")

	cfg := config.New("slack.config.yml")
	usersdbClient, err := users.New(log)

	if err != nil {
		log.PrintErr(err)
		panic(err)
	}

	utmCachClient, err := userstatemanagers.New(log)

	if err != nil {
		log.PrintErr(err)
		panic(err)
	}

	chatBot := slackservice.New(cfg, log)
	chatService := chatservice.New(usersdbClient, utmCachClient)

	if err := chatService.StartChat(chatBot, log); err != nil {
		log.PrintErr(err)
	}
}
