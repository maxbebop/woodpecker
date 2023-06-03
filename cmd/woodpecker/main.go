package main

import (
	config "woodpecker/internal/configs"
	pudgedb "woodpecker/internal/integrations/pudge"
	models "woodpecker/internal/models"
	storage "woodpecker/internal/storage/pudge"

	chatservice "woodpecker/internal/services/chat"
	slackservice "woodpecker/internal/services/slack"

	"github.com/powerman/structlog"
)

func main() {
	log := structlog.New()

	log.Info("start task management bot - woodpecker")

	cfg := config.New("slack.config.yml")
	usersdbClient, err := storage.New[models.User](pudgedb.Db, "users", log) //pudgedb.NewClient(pudgedb.DB, log)
	if err != nil {
		log.PrintErr(err)
		panic(err)
	}

	chatBot := slackservice.New(cfg, log)
	chatService := chatservice.New(usersdbClient)

	if err := chatService.StartChat(chatBot, log); err != nil {
		log.PrintErr(err)
	}
}
