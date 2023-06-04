package usertaskmanager

import models "woodpecker/internal/models"

type WaitTaskState struct {
	userTaskManager *UserTaskManager
}

const testWaitTaskStateMsg = "you don't have any task!"

func (i *WaitTaskState) compute(env models.Environment, handler StateHandler) error {
	handler.SendMessageByState(i.userTaskManager.environment.User, i.userTaskManager.environment.User.ChatToken, testWaitTaskStateMsg, i.userTaskManager.log)
	return nil
}
