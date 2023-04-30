package main

import (
	"fmt"
	"log"
	"strings"
	config "woodpecker/configs"
	"woodpecker/internal/services/slack"
)

func main() {
	fmt.Println("start task managment bot - woodpicker")
	config := config.New("slack.config.yml")

	slackService := slack.New(config)
	inMsgChannel := make(chan slack.Message)
	go slackService.GetMessages(inMsgChannel)
	for inMsg := range inMsgChannel {
		if inMsg.Error != nil {
			log.Fatal(inMsg.Error)
		}

		outMsg := slack.OutMessage{Message: inMsg}
		//outMsg.Pretext = "Test answer"

		outMsg.Type = slack.Common

		if strings.Contains(inMsg.Text, "ok") {
			outMsg.Type = slack.Common
		} else if strings.Contains(inMsg.Text, "attention") {
			outMsg.Type = slack.Attention
		} else if strings.Contains(inMsg.Text, "warning") {
			outMsg.Type = slack.Warning
		}
		slackService.SendMessage(outMsg)
	}
}
