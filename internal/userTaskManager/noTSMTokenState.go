package usertaskmanager

import (
	"fmt"
	models "woodpecker/internal/models/user"
)

type NoTSMTokenState struct {
	userTaskManager *UserTaskManager
	//storaeClient    storage.Client[models.User]
}

func (i *NoTSMTokenState) AddUser(user models.User) error {
	return fmt.Errorf("user already added")
}

func (i *NoTSMTokenState) RequestTsmToken(user models.User) error {
	i.userTaskManager.setState(i.userTaskManager.waitTSMToken)

	return nil
}

func (i *NoTSMTokenState) SaveTsmToken(user models.User) error {
	return fmt.Errorf("user don't have tsm token")
}
