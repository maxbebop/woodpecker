package cacheclient

import (
	"sync"

	"github.com/powerman/structlog"
)

type client[T any] struct {
	log *structlog.Logger
	mu  *sync.RWMutex
	db  map[string]*T
}

func New[T any](log *structlog.Logger) (*client[T], error) {

	db := make(map[string]*T)
	c := &client[T]{
		db:  db,
		log: log,
		mu:  &sync.RWMutex{},
	}

	return c, nil
}

func (c *client[T]) Has(key string) bool {
	c.mu.RLock()
	has := false
	_, ok := c.db[key]
	has = ok

	c.mu.RUnlock()

	return has
}

func (c *client[T]) Get(key string) (*T, bool) {
	c.mu.RLock()

	val, ok := c.db[key]

	c.mu.RUnlock()

	return val, ok
}

func (c *client[T]) Set(key string, value *T) error {
	c.mu.Lock()
	c.db[key] = value
	c.mu.Unlock()

	return nil
}

func (c *client[T]) GetAllItems() ([]*T, error) {
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

func (c *client[T]) DebugAllValues() {
	c.mu.RLock()
	c.log.Debug("All key value --")
	for key := range c.db {
		val, ok := c.db[key]
		c.log.Debug("key: %v;, hasKey: %v; val: %v", string(key), ok, val)
	}
	c.log.Debug("-- -- -- --")
	c.mu.RUnlock()
}