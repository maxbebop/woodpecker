package userstatemanager

import (
	"fmt"
	"strings"
	models "woodpecker/internal/models"
)

type WaitTmsTokenState struct {
	userStateManager *UserStateManager
}

const (
	testTokenMask  = "token:"
	testSuccessMsg = "you have successfully registered!"
)

func (i *WaitTmsTokenState) compute(env models.Environment, handler StateHandler) error {
	if len(env.Msg) == 0 {
		return fmt.Errorf("tms token is empty")
	}

	env.User.TMSToken = getTMSToken(env.Msg)
	i.userStateManager.environment = env

	if err := i.userStateManager.userStorage.Set(string(env.User.MessengerToken), env.User); err != nil {
		return err //nolint:wrapcheck // intentional
	}

	handler.SendMessageByState(env.User, env.User.MessengerToken, testSuccessMsg, i.userStateManager.log)

	return i.userStateManager.setState(i.userStateManager.waitTask)
}

func getTMSToken(text string) models.TMSToken {
	return models.TMSToken(strings.ReplaceAll(text, testTokenMask, ""))
}
