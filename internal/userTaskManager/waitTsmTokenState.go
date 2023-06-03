package usertaskmanager

import (
	"errors"
	"strings"
	models "woodpecker/internal/models"
)

type WaitTmsTokenState struct {
	userTaskManager *UserTaskManager
}

func (i *WaitTmsTokenState) compute(env models.Environment, handler StateHandler) error {
	if len(env.Msg) == 0 {
		return errors.New("tms token is empty")
	}
	env.User.TMSToken = strings.ReplaceAll(env.Msg, "token:", "")
	i.userTaskManager.environment = env

	if err := i.userTaskManager.userStorage.Set(env.User.ChatToken, env.User); err != nil {
		return err
	}

	i.userTaskManager.userStorage.Set(env.ChatChanelId, env.User)
	handler.SendMessageByState(env.User, env.ChatChanelId, "you have successfully registered!", i.userTaskManager.log)

	return i.userTaskManager.setState(i.userTaskManager.waitTask)
}
