package models

type ChatChanelId string
type Environment struct {
	User         User
	ChatChanelId ChatChanelId
	Msg          string
}
