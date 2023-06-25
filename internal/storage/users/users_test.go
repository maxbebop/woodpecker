package users_test

import (
	"errors"
	"testing"
	models "woodpecker/internal/models"

	"github.com/stretchr/testify/require"
)

const userToken = "test_user_token"

var errSetValueByKey = errors.New("set value by key")

func TestHas(t *testing.T) {
	t.Parallel()

	usersClient := NewClient(t)

	usersClient.On("Has", userToken).Return(false)
	has := usersClient.Has(userToken)

	require.Equal(t, has, false, "has user token  (key not found)")
}

func TestGet(t *testing.T) {
	t.Parallel()

	usersClient := NewClient(t)
	usersClient.On("Get", userToken).Return(models.User{}, false)
	user, ok := usersClient.Get(userToken)
	require.Equal(t, user, models.User{}, "get user token (key not found) - value")
	require.Equal(t, ok, false, "get user token (key not found) - flag")
}

func TestSet(t *testing.T) {
	t.Parallel()

	testUser := models.User{
		MessengerToken: userToken,
	}
	usersClient := NewClient(t)
	usersClient.On("Set", userToken, testUser).Return(errSetValueByKey)
	err := usersClient.Set(userToken, testUser)
	require.EqualError(t, err, "set value by key")
}
