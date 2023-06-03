package usertaskmanager

import (
	models "woodpecker/internal/models"
	storage "woodpecker/internal/storage/pudge"

	"github.com/powerman/structlog"
)

/*
type state interface {
	compute(env models.Environment, process StateProcess) error
}*/

type StateHandler interface {
	SendMessageByState(user models.User, chatChannel string, msg string, log *structlog.Logger)
}

type (
	state interface {
		compute(env models.Environment, handler StateHandler) error
	}

	UserTaskManager struct {
		newUser      state
		noTMSToken   state
		waitTMSToken state
		waitTask     state

		currentState state

		environment models.Environment

		userStorage storage.Client[models.User]
		log         *structlog.Logger
	}
)

/*
type UserTaskManager struct {
	newUser      state
	noTMSToken   state
	waitTMSToken state
	waitTask     state

	currentState state

	environment models.Environment

	userStorage storage.Client[models.User]
	log         *structlog.Logger
} */

func New(userStorage storage.Client[models.User], log *structlog.Logger) *UserTaskManager {
	utm := &UserTaskManager{userStorage: userStorage, log: log}
	utm.initStates()
	utm.setState(utm.newUser)
	return utm
}

func (utm *UserTaskManager) initStates() {
	newUser := &NewUserState{
		userTaskManager: utm,
	}

	noTMSToken := &NoTSMTokenState{
		userTaskManager: utm,
	}

	waitTMSToken := &WaitTmsTokenState{
		userTaskManager: utm,
	}

	waitTask := &WaitTaskState{
		userTaskManager: utm,
	}

	utm.newUser = newUser
	utm.noTMSToken = noTMSToken
	utm.waitTMSToken = waitTMSToken
	utm.waitTask = waitTask
}

func (utm *UserTaskManager) Compute(env models.Environment, handler StateHandler) error {
	return utm.currentState.compute(env, handler)
}

func (utm *UserTaskManager) setState(s state) error {
	utm.currentState = s

	return nil
}
