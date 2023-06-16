package models

type UserMessengerToken string
type TMSToken string
type User struct {
	ID             int
	MessengerToken UserMessengerToken
	Email          string
	Name           string
	TMSToken       TMSToken
}

func (u User) HasTMSToken() bool {
	return len(u.TMSToken) > 0
}
