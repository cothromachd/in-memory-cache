package cache

import "fmt"

type cache struct {
	data map[string]interface{}
}

func New() *cache {
	m := make(map[string]interface{})
	cache := &cache{data: m}
	return cache
}

func (c *cache) Set(key string, value interface{}) {
	c.data[key] = value
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
