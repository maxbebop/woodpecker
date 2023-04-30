package main

import (
	config "woodpecker/configs"
	"woodpecker/internal/services/chat"
	"woodpecker/internal/services/slack"
)

func main() {
	config := config.New("slack.config.yml")

	slackService := slack.New(config)
	chat.StartChat(slackService)
}
