package pudgedb

import (
	"github.com/fastogt/pudge"
	"github.com/powerman/structlog"
)

type (
	Client interface {
		UserClient
	}

	client struct {
		log *structlog.Logger
		db  *pudge.Db
	}
)

func NewClient(log *structlog.Logger) (Client, error) {
	cfg := &pudge.Config{
		SyncInterval: 0} //  0 - file first
	db, err := pudge.Open("./db/users", cfg)
	if err != nil {
		log.Err(err.Error())
		return nil, err
	}

	c := &client{
		log: log,
		db:  db,
	}

	return c, nil
}
