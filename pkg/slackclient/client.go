package slackclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackClient interface {
	GetUserInfo(user string) (*slack.User, error)
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

type SocketmodeClient interface {
	Ack(req socketmode.Request, payload ...interface{})
	Run() error
	EventsIn() <-chan socketmode.Event
}

type Client struct {
	api          SlackClient
	socketClient SocketmodeClient
	botId        string
	log          Logger
}

type Message struct {
	User    string
	Channel ChannelID
	Text    string
	Error   error
}

type ChannelID string

func New(api SlackClient, sockClient SocketmodeClient, appUserId string, log Logger) *Client {
	return &Client{api: api, socketClient: sockClient, botId: appUserId, log: log}
}

type Logger interface {
	Printf(s string, params ...any)
	Err(msg interface{}, keyvals ...interface{}) error
}

func (c *Client) GetMessagesLoop(ctx context.Context, res chan<- Message) {
	for {
		select {
		case <-ctx.Done():
			c.log.Printf("Shutting down socketmode listener")
			return
		case event := <-c.socketClient.EventsIn():
			c.log.Printf("event: %v\n", event)

			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					c.log.Printf("Could not type cast the event to the EventsAPI: %T: %+v\n", event, event)
					continue
				}

				c.socketClient.Ack(*event.Request)
				c.handleBotEvent(eventsAPI, res)
			}
		}
	}
}

func (c *Client) Run() error {
	return c.socketClient.Run() //nolint:wrapcheck // intentional
}

func (c *Client) handleBotEvent(event slackevents.EventsAPIEvent, chatChannel chan<- Message) {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent

		switch evnt := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			handleBotEventMessage(evnt, c.api, c.botId, chatChannel)
		default:
			c.log.Err("unsupported inner event type", "type", evnt) //nolint:errcheck // intentional
		}
	default:
		c.log.Printf("unsupported event type", "type", event.Type)
	}
}

func handleBotEventMessage(
	event *slackevents.MessageEvent,
	_ SlackClient,
	botId string,
	chatChannel chan<- Message,
) {
	if botId != event.User {
		text := strings.ToLower(event.Text)
		chatChannel <- Message{User: event.User, Channel: ChannelID(event.Channel), Text: text}
	}
}

func (c *Client) SendMessage(message Message) error {
	user, err := c.api.GetUserInfo(message.User)
	if err != nil {
		c.log.Err("failed to post message: %v\n", err)
	}

	attachment := slack.Attachment{}
	attachment.Text = fmt.Sprintf("Hello %s! you msg: %s", user.Name, message.Text)
	attachment.Pretext = "Test answer"
	attachment.Color = "#4af030"

	_, _, errPostMsg := c.api.PostMessage(string(message.Channel), slack.MsgOptionAttachments(attachment))

	return errPostMsg //nolint:wrapcheck // intentional
}
