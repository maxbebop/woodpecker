package userstatemanager

import (
	"reflect"
	models "woodpecker/internal/models"
	"woodpecker/internal/storage/users"

	"github.com/powerman/structlog"
)

type StateHandler interface {
	SendMessageByState(user models.User, messengerToken models.UserMessengerToken, msg string, log *structlog.Logger)
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

		userStorage users.Client
		log         *structlog.Logger
	}
)

func New(userStorage users.Client, log *structlog.Logger) *UserStateManager {
	usm := &UserStateManager{userStorage: userStorage, log: log}
	usm.initStates()
	usm.setState(usm.newUser)
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
