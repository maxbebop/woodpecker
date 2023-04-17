package slack

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackClient interface {
}

type SockClient interface {
	Ack(req socketmode.Request, payload ...interface{})
}

type Client struct {
	api          *slack.Client
	socketClient *socketmode.Client
	botId        string
}

type Message struct {
	User    string
	Channel ChannelID
	Text    string
	Error   error
}

type ChannelID string

func New(oauthToken string, appToken string, appUserId string) *Client {
	api := slack.New(oauthToken, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	socketClient := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	return &Client{api: api, socketClient: socketClient, botId: appUserId}
}

type Logger interface {
	Printf(s string, params ...any)
	Err(msg interface{}, keyvals ...interface{}) error
}

func (c *Client) GetMessagesLoop(ctx context.Context, res chan Message, log Logger) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down socketmode listener")
			return
		case event := <-c.socketClient.Events:
			log.Printf("event: %v\n", event)

			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Printf("Could not type cast the event to the EventsAPI: %T: %+v\n", event, event)
					continue
				}

				c.socketClient.Ack(*event.Request)
				c.handleBotEvent(eventsAPI, res, log)
			}
		}
	}
}

func (c *Client) Run() error {
	return c.socketClient.Run()
}

func (c *Client) handleBotEvent(event slackevents.EventsAPIEvent, chatChannel chan Message, log Logger) {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent

		switch evnt := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			handleBotEventMessage(evnt, c.api, c.botId, chatChannel)
		default:
			log.Err("unsupported inner event type", "type", evnt)
		}
	default:
		log.Printf("unsupported event type", "type", event.Type)
	}
}

func handleBotEventMessage(
	event *slackevents.MessageEvent,
	_ *slack.Client,
	botId string,
	chatChannel chan Message,
) {
	if botId != event.User {
		text := strings.ToLower(event.Text)
		chatChannel <- Message{User: event.User, Channel: ChannelID(event.Channel), Text: text}
	}
}

func (c *Client) SendMessage(message Message) error {
	user, err := c.api.GetUserInfo(message.User)
	if err != nil {
		log.Printf("failed to post message: %v\n", err)
	}

	attachment := slack.Attachment{}
	attachment.Text = fmt.Sprintf("Hello %s! you msg: %s", user.Name, message.Text)
	attachment.Pretext = "Test answer"
	attachment.Color = "#4af030"

	_, _, errPostMsg := c.api.PostMessage(string(message.Channel), slack.MsgOptionAttachments(attachment))

	return errPostMsg
}
