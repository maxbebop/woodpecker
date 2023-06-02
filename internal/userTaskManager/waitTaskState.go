package usertaskmanager

import models "woodpecker/internal/models"

type WaitTaskState struct {
	userTaskManager UserTaskManager
	//dbClient        pudgedb.Client
}

func (i *WaitTaskState) Compute(env models.Environment, process StateProcess) error {
	process.SendMessage(i.userTaskManager.environment.User, i.userTaskManager.environment.User.ChatToken, "you don't have any task!", i.userTaskManager.log)
	return nil
}

/*
func (i *WaitTaskState) AddUser(userChatToken string) error {
	return fmt.Errorf("user already added")
}

func (i *WaitTaskState) RequestTsmToken() error {
	return fmt.Errorf("tsm token already saved")
}

func (i *WaitTaskState) SaveTsmToken(tsmToken string) error {
	return fmt.Errorf("tsm token already saved")
}
*/
