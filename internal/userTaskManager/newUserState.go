package usertaskmanager

import (
	models "woodpecker/internal/models"
	storage "woodpecker/internal/storage/pudge"
)

type NewUserState struct {
	userTaskManager UserTaskManager
	userStorage     storage.Client[models.User]
}

func (i *NewUserState) Compute(env models.Environment, process StateProcess) error {
	i.userTaskManager.environment = env
	if !i.userStorage.Has(env.User.ChatToken) {
		if err := i.userStorage.Set(env.User.ChatToken, env.User); err != nil {
			return err
		}
	}

	return i.userTaskManager.setState(i.userTaskManager.noTMSToken)
}

/*
func (i *NewUserState) addUser(userChatToken string) error {
	i.userTaskManager.user = models.User{ChatToken: userChatToken}
	if err := i.userStorage.Set(userChatToken, i.userTaskManager.user); err != nil {
		return err
	}

	return i.userTaskManager.setState(i.userTaskManager.noTSMToken)
}

func (i *NewUserState) RequestTsmToken() error {
	return fmt.Errorf("user not registered")
}

func (i *NewUserState) SaveTsmToken(tsmTocken string) error {
	return fmt.Errorf("user not registered")
}
*/
