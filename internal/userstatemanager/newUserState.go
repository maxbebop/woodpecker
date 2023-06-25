package userstatemanager

import (
	"fmt"
	"woodpecker/internal/models"
)

type NewUserState struct {
	userStateManager *UserStateManager
}

func (i *NewUserState) compute(env models.Environment, handler StateHandler) error {
	i.userStateManager.environment = env
	if !i.userStateManager.userStorage.Has(string(env.User.MessengerToken)) {
		err := i.userStateManager.userStorage.Set(string(env.User.MessengerToken), env.User)
		if err != nil {
			return fmt.Errorf("compute state: failed save new user %w", err)
		}
	}

	if env.User.TMSToken == "" {
		if err := i.userStateManager.setState(i.userStateManager.noTMSToken); err != nil {
			return fmt.Errorf("compute state: failed set noTMSToken state %w", err)
		}

		if err := i.userStateManager.currentState.compute(env, handler); err != nil {
			return fmt.Errorf("failed compute noTMSToken state: %v; %w", env, err)
		}

		return nil
	}

	if err := i.userStateManager.setState(i.userStateManager.waitTask); err != nil {
		return fmt.Errorf("failed compute state: set waitTask state %w", err)
	}

	return nil
}
