package usertaskmanager

import (
	"errors"
	models "woodpecker/internal/models"
	storage "woodpecker/internal/storage/pudge"
)

type WaitTmsTokenState struct {
	userTaskManager UserTaskManager
	userStorage     storage.Client[models.User]
}

func (i *WaitTmsTokenState) Compute(env models.Environment, process StateProcess) error {
	if len(env.Msg) == 0 {
		return errors.New("tms token is empty")
	}
	i.userTaskManager.environment = env

	if err := i.userStorage.Set(env.User.ChatToken, env.User); err != nil {
		return err
	}

	return i.userTaskManager.setState(i.userTaskManager.waitTask)
}

/*
func (i *WaitTsmTokenState) AddUser(userChatToken string) error {
	return fmt.Errorf("user already added")
}

func (i *WaitTsmTokenState) RequestTsmToken() error {
	return fmt.Errorf("tsm token already saved")
}

func (i *WaitTsmTokenState) SaveTsmToken(tsmToken string) error {
	i.userTaskManager.user.TMSToken = tsmToken
	if err := i.userStorage.Set(tsmToken, i.userTaskManager.user); err != nil {
		return err
	}

	if err := i.userTaskManager.setState(i.userTaskManager.waitTask); err != nil {
		return err
	}

	return i.userStorage.Set(i.userTaskManager.user.ChatToken, i.userTaskManager.user)
}
*/
