package userstatemanagers

import (
	cacheclient "woodpecker/internal/storage/cache"
	"woodpecker/internal/userstatemanager"

	"github.com/powerman/structlog"
)

func New(log *structlog.Logger) (*cacheclient.Client[userstatemanager.UserStateManager], error) {
	return cacheclient.New[userstatemanager.UserStateManager](log)
}
