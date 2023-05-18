package pudgedb

import (
	"fmt"
	models "woodpecker/internal/models/user"
)

type UserClient interface {
	Has(chatToken string) (bool, error)
	Get(chatToken string) (models.User, error)
	Set(user models.User) error
	DebugAllValues()
}

func (c *client) Has(chatToken string) (bool, error) {
	return c.db.Has(chatToken)
}
func (c *client) Get(chatToken string) (models.User, error) {
	var user models.User
	if err := c.db.Get(chatToken, &user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (c *client) Set(user models.User) error {
	return c.db.Set(user.ChatToken, &user)
}

func (c *client) DebugAllValues() {
	fmt.Println("All key value --")
	keys, _ := c.db.Keys(nil, 0, 0, true)
	for _, key := range keys {
		var u models.User
		c.db.Get(key, &u)
		fmt.Println(u)
	}
	fmt.Println("-- -- -- --")
}
