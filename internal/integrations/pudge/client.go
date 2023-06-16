package pudgedb

import (
	"github.com/fastogt/pudge"
	"github.com/powerman/structlog"
)

// 0 - file first;  2 - with empty file - memory without persist.
type Mode int

const (
	Db    Mode = 0
	Cache Mode = 2
)

func New(storeMode Mode, name string, log *structlog.Logger) (*pudge.Db, error) {

	pathDb := "./db/" + name
	cfg := &pudge.Config{
		StoreMode: int(storeMode)}

	db, err := pudge.Open(pathDb, cfg)
	if err != nil {
		return nil, log.Err(err)
	}

	return db, nil
}
