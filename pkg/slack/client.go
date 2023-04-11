package slack

import (
	"log"

	"github.com/slack-go/slack"
)

func Test() {

	OAUTH_TOKEN := "xoxb-2222867095812-5106208214065-ZnFlXa1EQWBAWXNsFBRhN9Gb"
	//CHANNEL_ID := "D052K8CAENS"  //it's me and woodpecker
	//CHANNEL_ID := "C052REFPXGB" //it.s test channel
	api := slack.New(OAUTH_TOKEN, slack.OptionDebug(true))
	attachment := slack.Attachment{
		Pretext: "task managment bot woodpecker messag",
		Text:    "Hello from woodpecker",
	}

	responce, err := api.AuthTest()
	if err != nil {
		log.Fatalf("AuthTest error: %s\n", err)
	}
	log.Printf("AuthTest responce %s \n", responce)

	user, err := api.GetUserByEmail("maxim.kubarskiy@microavia.com")
	if err != nil {
		log.Fatalf("GetUserByEmail error: %s\n", err)
	}
	log.Printf("GetUserByEmail user %s - %s \n", user.ID, user.Name)
	channelId, timestamp, err := api.PostMessage(
		string(user.ID),
		//slack.MsgOptionText("This is the main message", false),
		slack.MsgOptionAttachments(attachment),
		//slack.MsgOptionAsUser(true),
	)

	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	log.Printf("Message successfully sent to Channel %s at %s\n", channelId, timestamp)
}
