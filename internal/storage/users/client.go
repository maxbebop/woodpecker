package users

import (
	"woodpecker/internal/models"

	pudgeclient "woodpecker/internal/storage/pudge"

	"github.com/powerman/structlog"
)

type Client interface {
	Has(key string) bool
	Get(key string) (models.User, bool)
	Set(key string, value models.User) error
	GetAllItems() ([]models.User, error)
	//todo: for development only. this should be removed
	DebugAllValues()
}

func New(log *structlog.Logger) (Client, error) {
	return pudgeclient.New[models.User](log)
}
