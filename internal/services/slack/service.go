package slack

import (
	"log"
	"strings"
	config "woodpecker/configs"
	"woodpecker/pkg/slack"
)

type Service struct {
	client *slack.Client
}

type Message struct {
	User    string
	Channel string
	Text    string
	Error   error
}

type OutMessage struct {
	Message Message
	Pretext string
	Type    MessageType
	Error   error
}

type MessageType int

const (
	Common MessageType = iota
	Warning
	Attention
)

func New(slackConfig *config.Config) *Service {

	client := slack.New(slackConfig.Slack.OAuthToken, slackConfig.Slack.AppToken, slackConfig.Slack.AppUserId)
	return &Service{client: client}
}

func (service *Service) GetMessages(inMsgChannel chan Message) {

	chatChannel := make(chan slack.Message)
	go service.client.GetMessages(chatChannel)
	for msg := range chatChannel {
		if msg.Error != nil {
			log.Fatal(msg.Error)
		}
		inMsgChannel <- Message{User: msg.User, Channel: string(msg.Channel), Text: msg.Text, Error: msg.Error}
	}
}

func (service *Service) SendMessage(msg OutMessage) {
	if service.client != nil {
		slackOutMsg := msg.createSlackOutMessage()
		service.client.SendMessage(slackOutMsg)
	}
}

func (msg OutMessage) createSlackOutMessage() slack.OutMessage {
	slackOutMsg := slack.OutMessage{User: msg.Message.User, Channel: slack.ChannelID(msg.Message.Channel), Text: msg.Message.Text, Pretext: msg.Pretext}
	slackOutMsg.Color = msg.Type.getMsgColor()

	if strings.TrimSpace(slackOutMsg.Pretext) == "" {

		slackOutMsg.Pretext = msg.Type.getMsgPretext()
	}

	return slackOutMsg
}

func (msgType MessageType) getMsgPretext() string {

	switch msgType {
	case Warning:
		return "Warning"
	case Attention:
		return "Attention"
	}

	return ""
}

func (msgType MessageType) getMsgColor() string {

	switch msgType {
	case Warning:
		return "#ff9900"
	case Attention:
		return "#ff471a"
	}

	return "#4af030"
}
