package main

import (
	config "woodpecker/internal/configs"
	pudgedb "woodpecker/internal/integrations/pudge"
	models "woodpecker/internal/models/user"
	storage "woodpecker/internal/storage/pudge"
	userTaskManager "woodpecker/internal/userTaskManager"

	//pudgedb "woodpecker/internal/integrations/pudgedb"
	chatservice "woodpecker/internal/services/chat"
	slackservice "woodpecker/internal/services/slack"

	//storage "woodpecker/internal/storage"

	"github.com/powerman/structlog"
	//"golang.org/x/mod/sumdb/storage"
)

func main() {
	log := structlog.New()

	log.Info("start task management bot - woodpecker")

	cfg := config.New("slack.config.yml")
	dbClient, err := storage.New[models.User](pudgedb.Db, "users", log) //pudgedb.NewClient(pudgedb.DB, log)
	if err != nil {
		log.PrintErr(err)
		panic(err)
	}
	cacheClient, err := storage.New[userTaskManager.UserTaskManager](pudgedb.Db, "users", log) //udgedb.NewClient(pudgedb.Cache, log)
	if err != nil {
		log.PrintErr(err)
		panic(err)
	}

	chatBot := slackservice.New(cfg, log)
	chatService := chatservice.New(cacheClient, dbClient)

	if err := chatService.StartChat(chatBot, log); err != nil {
		log.PrintErr(err)
	}
}
