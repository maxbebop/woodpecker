package userstatemanager_test

import (
	"fmt"
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
	code := usm.GetCode()
	require.Equal(t, "NewUserState", code, "state newUser")
}

func testComputeWaitTmsTokenState(t *testing.T, usm *userstatemanager.UserStateManager) {
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
	stateHandler := NewStateHandler(t)
	env := createMokeEnvironment()
	err := usm.Compute(env, stateHandler)

	require.Error(t, err, "tms token is empty")

}
func testComputeWaitTaskState(t *testing.T, usm *userstatemanager.UserStateManager) {
	var log *structlog.Logger
	stateHandler := NewStateHandler(t)
	env := createMokeEnvironment()
	env.Msg = "token:user_task_managment_system_token"
	env.User.TMSToken = "user_task_managment_system_token"
	msg := "you have successfully registered!"
	stateHandler.On("SendMessageByState", env.User, env.User.MessengerToken, msg, log).Once().Return(nil)
	err := usm.Compute(env, stateHandler)

	require.NoError(t, err, "WaitTaskState Compute")
	require.Equal(t, "WaitTaskState", usm.GetCode(), "state WaitTmsToken")
}
func createMokeEnvironment() models.Environment {
	return models.Environment{
		User: createMokeUser(),
	}
}

func createMokeUser() models.User {
	return models.User{
		MessengerToken: "user_messager_token",
	}
}

func (ut *usersStorage) Has(key string) bool {
	_, ok := ut.db[key]

	return ok
}

func (ut *usersStorage) Get(key string) (models.User, bool) {
	val, ok := ut.db[key]

	return val, ok
}

func (ut *usersStorage) Set(key string, value models.User) error {
	ut.db[key] = value

	return nil
}

func (ut *usersStorage) GetAllItems() ([]models.User, error) {
	result := []models.User{}
	for key := range ut.db {
		if val, ok := ut.db[key]; ok {
			result = append(result, val)
		}
	}

	return result, nil
}

func (ut *usersStorage) DebugAllValues() {
	fmt.Println(ut)
}
