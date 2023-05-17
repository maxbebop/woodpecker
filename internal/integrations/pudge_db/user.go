package pudgedb

import (
	models "woodpecker/internal/models/models"
)

type UserClient interface {
	Get(chatToken string) (models.User, error)
	Set(user models.User) error
}

func (c *client) Get(chatToken string) (models.User, error) {
	return nil, nil
}

func (c *client) Set(user models.User) error {
	return nil
}
