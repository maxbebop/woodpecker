package main

import (
	"fmt"
	"log"
	config "woodpecker/configs"

	"woodpecker/pkg/slack"
)

func main() {
	fmt.Println("start task managment bot - woodpicker")
	config := config.New("config.yml")
	client := slack.New(config.Slack.OAuthToken, config.Slack.AppToken, config.Slack.AppUserId)
	chatChannel := make(chan slack.Message)
	go client.GetMessages(chatChannel)
	for message := range chatChannel {
		if message.Error != nil {
			log.Fatal(message.Error)
		}
		log.Printf("msg from user: %v; text: %v\n", message.User, message.Text)
		client.SendMessage(message)
	}
}
