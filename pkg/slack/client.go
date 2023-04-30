package slack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

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

type OutMessage struct {
	User    string
	Channel ChannelID
	Text    string
	Pretext string
	Color   string
	Error   error
}

type ChannelID string

func New(oauthToken string, appToken string, appUserId string) *Client {
	log.Printf("oauthToken: %v, appToken: %v, appUserId: %v\n", oauthToken, appToken, appUserId)
	api := slack.New(oauthToken, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	socketClient := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	return &Client{api: api, socketClient: socketClient, botId: appUserId}
}

func (client *Client) GetMessages(res chan Message) {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func(ctx context.Context, client *Client) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-client.socketClient.Events:
				log.Printf("event: %v\n", event)
				switch event.Type {

				case socketmode.EventTypeEventsAPI:
					eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPI: %v\n", event)
						continue
					}
					client.socketClient.Ack(*event.Request)
					HandleBotEvent(eventsAPI, client, res)
				}
			}
		}
	}(ctx, client)

	client.socketClient.Run()
}

func HandleBotEvent(event slackevents.EventsAPIEvent, client *Client, chatChannel chan Message) {

	switch event.Type {

	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent

		switch evnt := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			HandleBotEventMessage(evnt, client.api, client.botId, chatChannel)
		}
	default:
		log.Println(errors.New("unsupported event type"))
	}
}

func HandleBotEventMessage(event *slackevents.MessageEvent, api *slack.Client, botId string, chatChannel chan Message) {

	if botId != event.User {
		text := strings.ToLower(event.Text)
		chatChannel <- Message{User: event.User, Channel: ChannelID(event.Channel), Text: text}
	}
}

func (client *Client) SendMessage(message OutMessage) {
	log.Printf("message out: %v\n", message)
	err := message.validate(client)
	if err != nil {
		log.Printf("failed to send message: %v\n", err)
	}

	user, err := client.api.GetUserInfo(message.User)
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
	}

	attachment := slack.Attachment{}
	attachment.Text = fmt.Sprintf("%s -> %s", user.Name, message.Text)
	attachment.Pretext = message.Pretext
	attachment.Color = message.Color
	attachment.AuthorID = client.botId

	_, _, errPostMsg := client.api.PostMessage(string(message.Channel), slack.MsgOptionAttachments(attachment))
	if errPostMsg != nil {
		log.Printf("failed to post message: %v\n", errPostMsg)
	}
}

func (message OutMessage) validate(client *Client) error {
	if message.User == client.botId {
		return errors.New("failed to post message: sending a message to yourself")
	}

	return nil
}
