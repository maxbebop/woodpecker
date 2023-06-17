package pudgeclient

import (
	"fmt"
	pudgedb "woodpecker/internal/integrations/pudge"

	"github.com/fastogt/pudge"
	"github.com/powerman/structlog"
)

type Client[T any] struct {
	log *structlog.Logger
	db  *pudge.Db
}

func New[T any](log *structlog.Logger) (*Client[T], error) {
	name := fmt.Sprintf("%T", *new(T))
	db, err := pudgedb.New(pudgedb.DB, name, log)

	if err != nil {
		return nil, err
	}

	c := &Client[T]{
		db:  db,
		log: log,
	}

	return c, nil
}

func (c *Client[T]) Has(key string) bool {
	has, err := c.db.Has(key)

	if err != nil {
		return false
	}

	return has
}

func (c *Client[T]) Get(key string) (T, bool) {
	var val T

	if err := c.db.Get(key, &val); err != nil {
		return val, false
	}

	return val, true
}

func (c *Client[T]) Set(key string, value T) error {
	return c.db.Set(key, value)
}

func (c *Client[T]) GetAllItems() ([]T, error) {
	result := []T{}
	keys, err := c.db.Keys(nil, 0, 0, true)

	if err != nil {
		return result, err
	}

	for _, key := range keys {
		var u T
		if err := c.db.Get(key, &u); err == nil {
			result = append(result, u)
		}
	}

	return result, nil
}

func (c *Client[T]) DebugAllValues() {
	c.log.Debug("All key value --")
	keys, _ := c.db.Keys(nil, 0, 0, true)

	for _, key := range keys {
		var u T
		err := c.db.Get(key, &u)
		c.log.Debug("key: %v; val: %v; err: %v", string(key), u, err)
	}

	c.log.Debug("-- -- -- --")
}
