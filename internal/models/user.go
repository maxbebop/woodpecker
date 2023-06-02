package models

type User struct {
	Id        int
	ChatToken string
	Email     string
	Name      string
	TMSToken  string
}

func (u User) HasTMSToken() bool {
	return len(u.TMSToken) > 0
}
