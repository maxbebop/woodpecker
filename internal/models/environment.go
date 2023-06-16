package models

type ChatChanelID string
type Environment struct {
	User         User
	ChatChanelID ChatChanelID
	Msg          string
}
