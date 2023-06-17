package pudgedb

import (
	"github.com/fastogt/pudge"
	"github.com/powerman/structlog"
)

// 0 - file first;  2 - with empty file - memory without persist.
type Mode int

const (
	DB    Mode = 0
	Cache Mode = 2
)

func New(storeMode Mode, name string, log *structlog.Logger) (*pudge.Db, error) {
	pathDB := "./db/" + name
	//nolint:default struct params check // intentional
	cfg := &pudge.Config{
		StoreMode: int(storeMode),
	}

	db, err := pudge.Open(pathDB, cfg)
	if err != nil {
		return nil, log.Err(err)
	}

	return db, nil
}
