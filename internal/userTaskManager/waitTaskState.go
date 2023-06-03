package usertaskmanager

import models "woodpecker/internal/models"

type WaitTaskState struct {
	userTaskManager *UserTaskManager
}

func (i *WaitTaskState) compute(env models.Environment, handler StateHandler) error {
	handler.SendMessageByState(i.userTaskManager.environment.User, i.userTaskManager.environment.User.ChatToken, "you don't have any task!", i.userTaskManager.log)
	return nil
}
