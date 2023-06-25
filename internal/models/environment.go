package models

type ChatChanelID string
type Environment struct {
	User         User
	ChatChanelID ChatChanelID `exhaustruct:"optional"`
	Msg          string       `exhaustruct:"optional"`
}
