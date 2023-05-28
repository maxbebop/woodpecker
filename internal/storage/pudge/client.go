package storage

import (
	"fmt"
	pudgedb "woodpecker/internal/integrations/pudge"

	"github.com/fastogt/pudge"
	"github.com/powerman/structlog"
)

type Client[T any] interface {
	Has(key string) (bool, error)
	Get(key string) (T, error)
	Set(key string, value T) error
	DebugAllValues()
}

type client[T any] struct {
	log *structlog.Logger
	db  *pudge.Db
}

func New[T any](storeMode pudgedb.Mode, name string, log *structlog.Logger) (Client[T], error) {

	db, err := pudgedb.New(storeMode, name, log)
	if err != nil {
		return nil, err
	}
	c := &client[T]{
		db:  db,
		log: log,
	}

	return c, nil
}

func (c *client[T]) Has(key string) (bool, error) {
	return c.db.Has(key)
}

func (c *client[T]) Get(key string) (T, error) {
	var val T
	if err := c.db.Get(key, &val); err != nil {
		return val, err
	}
	return val, nil
}

func (c *client[T]) Set(key string, value T) error {
	return c.db.Set(key, value)
}

func (c *client[T]) DebugAllValues() {
	fmt.Println("All key value --")
	keys, _ := c.db.Keys(nil, 0, 0, true)
	for _, key := range keys {
		var u T
		c.db.Get(key, &u)
		fmt.Println(u)
	}
	fmt.Println("-- -- -- --")
}
