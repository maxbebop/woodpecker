package models

type UserMessengerToken string
type TMSToken string
type User struct {
	ID             int                `exhaustruct:"optional"`
	MessengerToken UserMessengerToken `exhaustruct:"optional"`
	Email          string             `exhaustruct:"optional"`
	Name           string             `exhaustruct:"optional"`
	TMSToken       TMSToken           `exhaustruct:"optional"`
}

func (u User) HasTMSToken() bool {
	return len(u.TMSToken) > 0
}
