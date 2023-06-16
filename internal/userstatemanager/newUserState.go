package userstatemanager

import (
	models "woodpecker/internal/models"
)

type NewUserState struct {
	userStateManager *UserStateManager
}

func (i *NewUserState) compute(env models.Environment, handler StateHandler) error {
	i.userStateManager.environment = env
	if !i.userStateManager.userStorage.Has(string(env.User.MessengerToken)) {
		if err := i.userStateManager.userStorage.Set(string(env.User.MessengerToken), env.User); err != nil {
			return err
		}
	}

	if env.User.TMSToken == "" {
		if err := i.userStateManager.setState(i.userStateManager.noTMSToken); err != nil {
			return err
		}

		return i.userStateManager.currentState.compute(env, handler)
	}

	return i.userStateManager.setState(i.userStateManager.waitTask)
}
