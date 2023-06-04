package usertaskmanager

import models "woodpecker/internal/models"

type NoTSMTokenState struct {
	userTaskManager *UserTaskManager
}

const testNoTSMTokenStateMsg = "Hello! Send me your tsm token as string token:you_token"

func (i *NoTSMTokenState) compute(env models.Environment, handler StateHandler) error {
	handler.SendMessageByState(i.userTaskManager.environment.User, i.userTaskManager.environment.ChatChanelId, testNoTSMTokenStateMsg, i.userTaskManager.log)
	return i.userTaskManager.setState(i.userTaskManager.waitTMSToken)
}
