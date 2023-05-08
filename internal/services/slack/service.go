package slackservice

import (
	"context"
	"errors"
	"fmt"
	"strings"
	config "woodpecker/configs"
	chatservice "woodpecker/internal/services/chat"
	slackclient "woodpecker/pkg/slack"
	socketmodewrap "woodpecker/pkg/socketmodewrap"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	"github.com/powerman/structlog"
)

type Service struct {
	client *slackclient.Client
}

var isNilErr = errors.New("is nil")

func New(cfg *config.Config, log *structlog.Logger) *Service {
	slackClient := slack.New(
		cfg.Slack.OAuthToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(cfg.Slack.AppToken),
	)

	socketClient := socketmode.New(
		slackClient,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(socketmodewrap.Log{Logger: log}),
	)

	client := slackclient.New(slackClient, socketmodewrap.New(socketClient), log)
	service := Service{client: client}

	return &service
}

func (service *Service) Run() error {
	err := service.client.Run()

	return fmt.Errorf("run %w", err)
}

func (service *Service) GetMessagesLoop(
	ctx context.Context,
	inChatMsgChannel chan chatservice.Message,
	log *structlog.Logger) {
	slackChatChannel := make(chan slackclient.Message)

	go service.client.GetMessagesLoop(ctx, slackChatChannel)
	go processMsgLoop(ctx, slackChatChannel, inChatMsgChannel, log)
}

func processMsgLoop(
	ctx context.Context,
	slackChatChannel chan slackclient.Message,
	inChatMsgChannel chan chatservice.Message,
	log *structlog.Logger,
) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down processing loop")
			return
		case msg := <-slackChatChannel:
			inChatMsgChannel <- chatservice.Message{
				User:    msg.User,
				Channel: string(msg.Channel),
				Text:    msg.Text,
				Error:   msg.Error}
		}
	}
}

func (service *Service) SendMessage(msg chatservice.OutMessage) error {
	if service.client != nil {
		slackOutMsg := createSlackOutMessage(msg)
		errPostMsg := service.client.SendMessage(slackOutMsg)

		return fmt.Errorf("send message. %w", errPostMsg)
	}

	return fmt.Errorf("service.client %w", isNilErr)
}

func createSlackOutMessage(msg chatservice.OutMessage) slackclient.OutMessage {
	slackOutMsg := slackclient.OutMessage{
		User:    msg.Message.User,
		Channel: slackclient.ChannelID(msg.Message.Channel),
		Text:    msg.Message.Text,
		Pretext: msg.Pretext,
		Color:   getMsgColor(msg.Type),
		Error:   nil,
	}

	if strings.TrimSpace(slackOutMsg.Pretext) == "" {
		slackOutMsg.Pretext = getMsgPretext(msg.Type)
	}

	return slackOutMsg
}

func getMsgPretext(msgType chatservice.MessageType) string {
	switch msgType {
	case chatservice.Warning:
		return "Warning"
	case chatservice.Attention:
		return "Attention"
	}

	return ""
}

func getMsgColor(msgType chatservice.MessageType) string {
	switch msgType {
	case chatservice.Warning:
		return "#ff9900"
	case chatservice.Attention:
		return "#ff471a"
	}

	return "#4af030"
}
