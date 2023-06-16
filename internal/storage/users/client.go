package users

import (
	"woodpecker/internal/models"

	pudgeclient "woodpecker/internal/storage/pudge"

	"github.com/powerman/structlog"
)

func New(log *structlog.Logger) (*pudgeclient.Client[models.User], error) {
	return pudgeclient.New[models.User](log)
}
