package cacheclient

import (
	"sync"

	"github.com/powerman/structlog"
)

type Client[T any] struct {
	log *structlog.Logger
	mu  *sync.RWMutex
	db  map[string]*T
}

func New[T any](log *structlog.Logger) (*Client[T], error) {
	db := make(map[string]*T)
	c := &Client[T]{
		db:  db,
		log: log,
		mu:  &sync.RWMutex{},
	}

	return c, nil
}

func (c *Client[T]) Has(key string) bool {
	c.mu.RLock()
	has := false
	_, ok := c.db[key]
	has = ok

	c.mu.RUnlock()

	return has
}

func (c *Client[T]) Get(key string) (*T, bool) {
	c.mu.RLock()

	val, ok := c.db[key]

	c.mu.RUnlock()

	return val, ok
}

func (c *Client[T]) Set(key string, value *T) error {
	c.mu.Lock()
	c.db[key] = value
	c.mu.Unlock()

	return nil
}

func (c *Client[T]) GetAllItems() ([]*T, error) {
	c.mu.RLock()
	result := []*T{}

	for key := range c.db {
		if val, ok := c.db[key]; ok {
			result = append(result, val)
		}
	}

	c.mu.RUnlock()

	return result, nil
}

func (c *Client[T]) DebugAllValues() {
	c.mu.RLock()
	c.log.Debug("All key value --")

	for key := range c.db {
		val, ok := c.db[key]
		c.log.Debug("key: %v;, hasKey: %v; val: %v", key, ok, val)
	}

	c.log.Debug("-- -- -- --")
	c.mu.RUnlock()
}
