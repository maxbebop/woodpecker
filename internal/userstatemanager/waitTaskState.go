package userstatemanager

import models "woodpecker/internal/models"

type WaitTaskState struct {
	userStateManager *UserStateManager
}

const testWaitTaskStateMsg = "you don't have any task!"

func (i *WaitTaskState) compute(env models.Environment, handler StateHandler) error {
	handler.SendMessageByState(
		env.User,
		env.User.MessengerToken,
		testWaitTaskStateMsg,
		i.userStateManager.log,
	)
	return nil
}
