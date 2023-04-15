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
	for message := range client.GetMessages() {
		log.Printf("msg ->  %v\n", message)
		if message.Error != nil {
			log.Fatal(message.Error)
		}
		log.Printf("msg -> user: %v; text: %v\n", message.User, message.Text)
	}
}
