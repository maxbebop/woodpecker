package slack

import (
	"context"
	"errors"
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

func (client *Client) Test() chan Message {
	res := make(chan Message)
	for i := 0; i < 5; i++ {
		go func(i int) {
			res <- Message{User: string(rune(i)), Channel: "test", Text: "test test"}
		}(i)
	}

	return res
}
func (client *Client) GetMessages() chan Message {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	chatChannel := make(chan Message)
	go func(ctx context.Context, client *Client) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-client.socketClient.Events:

				switch event.Type {

				case socketmode.EventTypeEventsAPI:
					eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPI: %v\n", event)
						continue
					}
					client.socketClient.Ack(*event.Request)
					log.Printf("EventsAPI: %v\n", eventsAPI)
					chatChannel <- Message{User: "test user", Channel: "test Channel", Text: "test Text"}
					/*
						err := HandleBotEvent(eventsAPI, api, botId, chatChannel)
						if err != nil {
							log.Println(err)
						} */
					//if msg.isAvailable() {
					//	res <- msg
					//}
				}
			}
		}
	}(ctx, client)

	client.socketClient.Run()

	return chatChannel
}

func HandleBotEvent(event slackevents.EventsAPIEvent, client *Client, chatChannel chan Message) error {

	switch event.Type {

	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent

		switch evnt := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			err := HandleBotEventMessage(evnt, client.api, client.botId, chatChannel)
			if err != nil {
				return err
			}
			//case *slackevents.AppMentionEvent:
			//	printAppMentionEventInfo(evnt)
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

func HandleBotEventMessage(event *slackevents.MessageEvent, api *slack.Client, botId string, chatChannel chan Message) error {

	if botId == event.User {
		return nil
	}

	text := strings.ToLower(event.Text)

	chatChannel <- Message{User: event.User, Channel: ChannelID(event.Channel), Text: text}

	return nil
}
