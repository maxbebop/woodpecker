package usertaskmanager

import (
	"fmt"
	models "woodpecker/internal/models/user"
)

type WaitTsmTokenState struct {
	userTaskManager *UserTaskManager
	//dbClient        pudgedb.Client
}

func (i *WaitTsmTokenState) AddUser(user models.User) error {
	return fmt.Errorf("user already added")
}

func (i *WaitTsmTokenState) RequestTsmToken(user models.User) error {
	return fmt.Errorf("tsm token already saved")
}

func (i *WaitTsmTokenState) SaveTsmToken(user models.User) error {
	i.userTaskManager.setState(i.userTaskManager.waitTask)

	return nil
}
