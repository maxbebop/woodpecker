package userstatemanager_test

import (
	"testing"
	"woodpecker/internal/models"
	"woodpecker/internal/userstatemanager"

	structlog "github.com/powerman/structlog"
	"github.com/stretchr/testify/require"
)

type usersStorage struct {
	db map[string]models.User
}

func TestCompute(t *testing.T) {
	t.Parallel()

	testUsersStorage := &usersStorage{
		db: make(map[string]models.User),
	}
	usm := userstatemanager.New(testUsersStorage, nil)
	testInitState(t, usm)

	testComputeWaitTmsTokenState(t, usm)
	testComputeWaitTaskStateError(t, usm)
	testComputeWaitTaskState(t, usm)
}

func testInitState(t *testing.T, usm *userstatemanager.UserStateManager) {
	t.Helper()

	code := usm.GetCode()
	require.Equal(t, "NewUserState", code, "state newUser")
}

func testComputeWaitTmsTokenState(t *testing.T, usm *userstatemanager.UserStateManager) {
	t.Helper()

	var log *structlog.Logger

	stateHandler := NewStateHandler(t)
	env := createMokeEnvironment()
	msg := "Hello! Send me your tsm token as string token:you_token"
	stateHandler.On("SendMessageByState", env.User, env.User.MessengerToken, msg, log).Once().Return(nil)
	err := usm.Compute(env, stateHandler)

	require.NoError(t, err, "NewUserState Compute")
	require.Equal(t, "WaitTmsTokenState", usm.GetCode(), "state WaitTmsToken")
}

func testComputeWaitTaskStateError(t *testing.T, usm *userstatemanager.UserStateManager) {
	t.Helper()

	stateHandler := NewStateHandler(t)
	env := createMokeEnvironment()
	err := usm.Compute(env, stateHandler)

	require.Error(t, err, "tms token is empty")
}
func testComputeWaitTaskState(t *testing.T, usm *userstatemanager.UserStateManager) {
	t.Helper()

	var log *structlog.Logger

	stateHandler := NewStateHandler(t)
	env := createMokeEnvironment()
	env.Msg = "token:user_task_management_system_token"
	env.User.TMSToken = "user_task_management_system_token"
	msg := "you have successfully registered!"
	stateHandler.On("SendMessageByState", env.User, env.User.MessengerToken, msg, log).Once().Return(nil)
	err := usm.Compute(env, stateHandler)

	require.NoError(t, err, "WaitTaskState Compute")
	require.Equal(t, "WaitTaskState", usm.GetCode(), "state WaitTmsToken")
}
func createMokeEnvironment() models.Environment {
	return models.Environment{
		User:         createMokeUser(),
		ChatChanelID: "",
		Msg:          "",
	}
}

func createMokeUser() models.User {
	return models.User{
		MessengerToken: "user_messager_token",
		ID:             0,
		Email:          "",
		Name:           "",
		TMSToken:       "",
	}
}

func (us *usersStorage) Has(key string) bool {
	_, ok := us.db[key]

	return ok
}

func (us *usersStorage) Get(key string) (models.User, bool) {
	val, ok := us.db[key]

	return val, ok
}

func (us *usersStorage) Set(key string, value models.User) error {
	us.db[key] = value

	return nil
}

func (us *usersStorage) GetAllItems() ([]models.User, error) {
	result := []models.User{}

	for key := range us.db {
		if val, ok := us.db[key]; ok {
			result = append(result, val)
		}
	}

	return result, nil
}

func (us *usersStorage) DebugAllValues() {

}
