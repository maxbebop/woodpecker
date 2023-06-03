package usertaskmanager

import (
	models "woodpecker/internal/models"
)

type NewUserState struct {
	userTaskManager *UserTaskManager
}

func (i *NewUserState) compute(env models.Environment, handler StateHandler) error {
	i.userTaskManager.environment = env
	if !i.userTaskManager.userStorage.Has(env.ChatChanelId) {
		if err := i.userTaskManager.userStorage.Set(env.ChatChanelId, env.User); err != nil {
			return err
		}
	}

	if env.User.TMSToken == "" {
		if err := i.userTaskManager.setState(i.userTaskManager.noTMSToken); err != nil {
			return err
		}
		return i.userTaskManager.currentState.compute(env, handler)
	}

	return i.userTaskManager.setState(i.userTaskManager.waitTask)
}
