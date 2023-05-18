package main

import (
	config "woodpecker/internal/configs"
	pudgedb "woodpecker/internal/integrations/pudge_db"
	chatservice "woodpecker/internal/services/chat"
	slackservice "woodpecker/internal/services/slack"

	"github.com/powerman/structlog"
)

func main() {
	log := structlog.New()

	log.Info("start task management bot - woodpecker")

	cfg := config.New("slack.config.yml")
	dbClient, err := pudgedb.NewClient(log)
	if err != nil {
		log.PrintErr(err)
		panic(err)
	}
	chatBot := slackservice.New(cfg, log)
	chatService := chatservice.New(dbClient)

	if err := chatService.StartChat(chatBot, log); err != nil {
		log.PrintErr(err)
	}
}
