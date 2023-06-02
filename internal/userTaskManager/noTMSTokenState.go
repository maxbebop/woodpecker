package usertaskmanager

import models "woodpecker/internal/models"

type NoTSMTokenState struct {
	userTaskManager UserTaskManager
}

func (i *NoTSMTokenState) Compute(env models.Environment, process StateProcess) error {
	msg := "Hello! Send me your tsm token as string token:you_token"
	process.SendMessage(i.userTaskManager.environment.User, i.userTaskManager.environment.User.ChatToken, msg, i.userTaskManager.log)
	return i.userTaskManager.setState(i.userTaskManager.waitTMSToken)
}

/*
func (i *NoTSMTokenState) AddUser(userChatToken string) error {
	return fmt.Errorf("user already added")
}

func (i *NoTSMTokenState) RequestTsmToken() error {
	return i.userTaskManager.setState(i.userTaskManager.waitTSMToken)
}

func (i *NoTSMTokenState) SaveTsmToken(tsmToken string) error {
	return fmt.Errorf("user don't have tsm token")
}
*/
