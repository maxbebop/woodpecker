package usertaskmanager

import (
	"errors"
	"strings"
	models "woodpecker/internal/models"
)

type WaitTmsTokenState struct {
	userTaskManager *UserTaskManager
}

const (
	testTokenMask  = "token:"
	testSuccessMsg = "you have successfully registered!"
)

func (i *WaitTmsTokenState) compute(env models.Environment, handler StateHandler) error {
	if len(env.Msg) == 0 {
		return errors.New("tms token is empty")
	}
	env.User.TMSToken = getTMSToken(env.Msg)
	i.userTaskManager.environment = env

	if err := i.userTaskManager.userStorage.Set(env.User.ChatToken, env.User); err != nil {
		return err
	}

	i.userTaskManager.userStorage.Set(env.ChatChanelId, env.User)
	handler.SendMessageByState(env.User, env.ChatChanelId, testSuccessMsg, i.userTaskManager.log)

	return i.userTaskManager.setState(i.userTaskManager.waitTask)
}

func getTMSToken(text string) string {
	return strings.ReplaceAll(text, testTokenMask, "")
}
