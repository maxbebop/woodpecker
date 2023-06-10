package userstatemanager

import models "woodpecker/internal/models"

type WaitTaskState struct {
	userStateManager *UserStateManager
}

// todo: test msg
const testWaitTaskStateMsg = "you don't have any task!"

func (i *WaitTaskState) compute(env models.Environment, handler StateHandler) error {
	handler.SendMessageByState(i.userStateManager.environment.User, i.userStateManager.environment.User.MessengerToken, testWaitTaskStateMsg, i.userStateManager.log)
	return nil
}
