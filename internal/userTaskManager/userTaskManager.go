package usertaskmanager

import (
	models "woodpecker/internal/models"

	"github.com/powerman/structlog"
)

type State interface {
	Compute(env models.Environment, process StateProcess) error
	/*addUser(userChatToken string) error
	requestTsmToken() error
	saveTsmToken(tsmTocken string) error */
}

type StateProcess interface {
	SendMessage(user models.User, chatChannel string, msg string, log *structlog.Logger)
}

/*
	type EnvironmentAction interface {
		Process(env Environment) error
	}
*/

type UserTaskManager struct {
	newUser      State
	noTMSToken   State
	waitTMSToken State
	waitTask     State

	currentState State

	environment models.Environment
	log         *structlog.Logger
}

func New(log *structlog.Logger) UserTaskManager {
	utm := UserTaskManager{log: log}
	utm.initStates()
	utm.setState(utm.newUser)
	return utm
}

/*
func NewByUser(user models.User) UserTaskManager {
	utm := UserTaskManager{user: user}
	utm.initStates()
	if user.HasTMSToken() {
		utm.setState(utm.waitTask)
	} else {
		utm.setState(utm.noTSMToken)
	}

	return utm
} */

func (utm UserTaskManager) initStates() {
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

func (utm UserTaskManager) Compute(env models.Environment, process StateProcess) error {
	return utm.currentState.Compute(env, process)
}

/*
	func (utm UserTaskManager) addUser(userChatToken string) error {
		return utm.currentState.addUser(userChatToken)
	}

	func (utm UserTaskManager) requestTsmToken() error {
		return utm.currentState.requestTsmToken()
	}

	func (utm UserTaskManager) saveTsmToken(tsmTocken string) error {
		return utm.currentState.saveTsmToken(tsmTocken)
	}
*/
func (utm UserTaskManager) setState(s State) error {
	utm.currentState = s

	return nil
}
