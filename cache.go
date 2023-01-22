package cache

import (
	"fmt"
	"time"
)

type cache struct {
	data map[string]interface{}
}

func New() *cache {
	m := make(map[string]interface{})
	cache := &cache{data: m}
	return cache
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) {
	c.data[key] = value
	time.AfterFunc(ttl, func() {
		if _, ok := c.data[key]; ok {
			delete(c.data, key)
		} else {
			fmt.Printf("Auto-delete error: there is no such key '%s' in storage\n", key)
		}
	})
}

func (c *cache) Get(key string) (interface{}, error) {
	if _, ok := c.data[key]; ok {
		return c.data[key], nil
	} else {
		return nil, fmt.Errorf("Get error: there is no such key '%s' in storage", key)
	}
}

func (c *cache) Delete(key string) error {
	if _, ok := c.data[key]; ok {
		delete(c.data, key)
		return nil
	} else {
		return fmt.Errorf("Delete error: there is no such key '%s' in storage", key)
	}
}
