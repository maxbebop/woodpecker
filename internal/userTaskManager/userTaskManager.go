package usertaskmanager

import (
	models "woodpecker/internal/models/user"
)

type State interface {
	AddUser(user models.User) error
	RequestTsmToken(user models.User) error
	SaveTsmToken(user models.User) error
}

type UserTaskManager struct {
	newUser      State
	noTSMToken   State
	waitTSMToken State
	waitTask     State

	currentState State

	user models.User
}

func New() UserTaskManager {
	utm := UserTaskManager{}
	utm.initStates()
	utm.setState(utm.newUser)
	return utm
}

func NewByUser(user models.User) UserTaskManager {
	utm := UserTaskManager{user: user}
	utm.initStates()
	if user.HasTMSToken() {
		utm.setState(utm.waitTask)
	} else {
		utm.setState(utm.noTSMToken)
	}

	return utm
}

func (utm UserTaskManager) initStates() {
	newUser := &NewUserState{
		userTaskManager: &utm,
	}

	noTSMToken := &NoTSMTokenState{
		userTaskManager: &utm,
	}

	waitTSMToken := &WaitTsmTokenState{
		userTaskManager: &utm,
	}

	waitTask := &WaitTaskState{
		userTaskManager: &utm,
	}

	utm.newUser = newUser
	utm.noTSMToken = noTSMToken
	utm.waitTSMToken = waitTSMToken
	utm.waitTask = waitTask
}

func (utm *UserTaskManager) AddUser(user models.User) error {
	return utm.currentState.AddUser(user)
}

func (utm *UserTaskManager) RequestTsmToken(user models.User) error {
	return utm.currentState.RequestTsmToken(user)
}

func (utm *UserTaskManager) SaveTsmToken(user models.User) error {
	return utm.currentState.SaveTsmToken(user)
}

func (utm *UserTaskManager) setState(s State) {
	utm.currentState = s
}
