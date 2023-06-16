package userstatemanager

import (
	"reflect"
	models "woodpecker/internal/models"

	"github.com/powerman/structlog"
)

type StateHandler interface {
	SendMessageByState(user models.User, messengerToken models.UserMessengerToken, msg string, log *structlog.Logger)
}

type UsersClient interface {
	Has(key string) bool
	Set(key string, value models.User) error
}

type (
	state interface {
		compute(env models.Environment, handler StateHandler) error
	}

	UserStateManager struct {
		newUser      state
		noTMSToken   state
		waitTMSToken state
		waitTask     state

		currentState state

		environment models.Environment

		userStorage UsersClient
		log         *structlog.Logger
	}
)

func New(userStorage UsersClient, log *structlog.Logger) *UserStateManager {
	usm := &UserStateManager{userStorage: userStorage, log: log}
	usm.initStates()
	_ = usm.setState(usm.newUser)

	return usm
}

func (usm *UserStateManager) initStates() {
	newUser := &NewUserState{
		userStateManager: usm,
	}

	noTMSToken := &NoTSMTokenState{
		userStateManager: usm,
	}

	waitTMSToken := &WaitTmsTokenState{
		userStateManager: usm,
	}

	waitTask := &WaitTaskState{
		userStateManager: usm,
	}

	usm.newUser = newUser
	usm.noTMSToken = noTMSToken
	usm.waitTMSToken = waitTMSToken
	usm.waitTask = waitTask
}

func (utm *UserStateManager) GetCode() string {
	t := reflect.TypeOf(utm.currentState).Elem()
	return t.Name()
}

func (utm *UserStateManager) Compute(env models.Environment, handler StateHandler) error {
	return utm.currentState.compute(env, handler)
}

func (utm *UserStateManager) setState(s state) error {
	utm.currentState = s

	return nil
}
