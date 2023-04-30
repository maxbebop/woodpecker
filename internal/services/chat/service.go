package chat

import (
	"woodpecker/internal/services/slack"
)

type service struct {
	chatBot ChatBot
}

type ChatBot interface {
	GetMessages(inMsgChannel chan slack.Message)
	SendMessage(msg slack.OutMessage)
}

func StartChat(chatBot ChatBot) error {
	service := new(chatBot)

	inMsgChannel := make(chan slack.Message)
	go service.chatBot.GetMessages(inMsgChannel)
	for inMsg := range inMsgChannel {
		if inMsg.Error != nil {
			return inMsg.Error
		}

		outMsg := slack.OutMessage{Message: inMsg}
		outMsg.Type = slack.Common

		service.chatBot.SendMessage(outMsg)
	}

	return nil
}

func new(chatBot ChatBot) *service {
	return &service{chatBot: chatBot}
}
