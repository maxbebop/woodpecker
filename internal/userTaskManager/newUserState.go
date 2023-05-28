package usertaskmanager

import (
	"fmt"
	models "woodpecker/internal/models/user"
	storage "woodpecker/internal/storage/pudge"
)

type NewUserState struct {
	userTaskManager *UserTaskManager
	storaeClient    storage.Client[models.User]
}

func (i *NewUserState) AddUser(user models.User) error {
	if err := i.storaeClient.Set(user.ChatToken, user); err != nil {
		return err
	}

	i.userTaskManager.setState(i.userTaskManager.noTSMToken)

	return nil
}

func (i *NewUserState) RequestTsmToken(user models.User) error {
	return fmt.Errorf("user not registered")
}

func (i *NewUserState) SaveTsmToken(user models.User) error {
	return fmt.Errorf("user not registered")
}
