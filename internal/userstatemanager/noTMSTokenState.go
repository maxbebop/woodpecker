package userstatemanager

import "woodpecker/internal/models"

type NoTSMTokenState struct {
	userStateManager *UserStateManager
}

const testNoTSMTokenStateMsg = "Hello! Send me your tsm token as string token:you_token"

func (i *NoTSMTokenState) compute(env models.Environment, handler StateHandler) error {
	handler.SendMessageByState(
		env.User,
		env.User.MessengerToken,
		testNoTSMTokenStateMsg,
		i.userStateManager.log,
	)
	return i.userStateManager.setState(i.userStateManager.waitTMSToken)
}
