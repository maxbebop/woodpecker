package userstatemanager

import models "woodpecker/internal/models"

type NoTSMTokenState struct {
	userStateManager *UserStateManager
}

// todo: test msg
const testNoTSMTokenStateMsg = "Hello! Send me your tsm token as string token:you_token"

func (i *NoTSMTokenState) compute(env models.Environment, handler StateHandler) error {
	handler.SendMessageByState(i.userStateManager.environment.User, i.userStateManager.environment.User.MessengerToken, testNoTSMTokenStateMsg, i.userStateManager.log)
	return i.userStateManager.setState(i.userStateManager.waitTMSToken)
}
