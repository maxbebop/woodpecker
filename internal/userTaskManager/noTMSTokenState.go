package usertaskmanager

import models "woodpecker/internal/models"

type NoTSMTokenState struct {
	userTaskManager *UserTaskManager
}

func (i *NoTSMTokenState) compute(env models.Environment, handler StateHandler) error {
	msg := "Hello! Send me your tsm token as string token:you_token"
	handler.SendMessageByState(i.userTaskManager.environment.User, i.userTaskManager.environment.ChatChanelId, msg, i.userTaskManager.log)
	return i.userTaskManager.setState(i.userTaskManager.waitTMSToken)
}
