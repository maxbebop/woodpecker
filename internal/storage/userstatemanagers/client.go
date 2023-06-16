package userstatemanagers

import (
	cacheclient "woodpecker/internal/storage/cache"
	"woodpecker/internal/userstatemanager"

	"github.com/powerman/structlog"
)

/*
type Client interface {
	Has(key string) bool
	Get(key string) (*userstatemanager.UserStateManager, bool)
	Set(key string, value *userstatemanager.UserStateManager) error
	GetAllItems() ([]*userstatemanager.UserStateManager, error)
	//todo: for development only. this should be removed
	DebugAllValues()
}
*/

func New(log *structlog.Logger) (*cacheclient.Client[userstatemanager.UserStateManager], error) {
	return cacheclient.New[userstatemanager.UserStateManager](log)
}
