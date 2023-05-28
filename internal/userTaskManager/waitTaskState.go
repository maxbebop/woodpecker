package usertaskmanager

import (
	"fmt"
	models "woodpecker/internal/models/user"
)

type WaitTaskState struct {
	userTaskManager *UserTaskManager
	//dbClient        pudgedb.Client
}

func (i *WaitTaskState) AddUser(user models.User) error {
	return fmt.Errorf("user already added")
}

func (i *WaitTaskState) RequestTsmToken(user models.User) error {
	return fmt.Errorf("tsm token already saved")
}

func (i *WaitTaskState) SaveTsmToken(user models.User) error {
	return fmt.Errorf("tsm token already saved")
}
