package pudgedb

import (
	"fmt"

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

	cfg := &pudge.Config{StoreMode: int(storeMode)} //nolint:exhaustruct

	db, err := pudge.Open(pathDB, cfg)
	if err != nil {
		return nil, fmt.Errorf("faild opend db %v; %w", name, log.Err("faild opend db", err))
	}

	return db, nil
}
