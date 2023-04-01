package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Cache struct {
	data *sync.Map
}

func New() *Cache {
	m := new(sync.Map)
	cache := &Cache{data: m}
	return cache
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.data.Store(key, value)
	time.AfterFunc(ttl, func() {
		if _, ok := c.data.Load(key); ok {
			c.data.Delete(key)
		} else {
			fmt.Printf("Auto-delete error: there is no such key '%s' in storage\n", key)
		}
	})
}

func (c *Cache) Get(key string) (interface{}, error) {
	if val, ok := c.data.Load(key); ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("Get() error: there is no such key '%s' in storage", key)
	}
}

func (c *Cache) Delete(key string) error {
	if _, ok := c.data.Load(key); ok {
		c.data.Delete(key)
		return nil
	} else {
		return fmt.Errorf("Delete() error: there is no such key '%s' in storage", key)
	}
}

func (c *Cache) Store() error {
	file, err := os.OpenFile("./bag.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("Store() error: %v", err)
	}

	defer file.Close()

	m := make(map[string]interface{})

	c.data.Range(func(key, value interface{}) bool {
		m[key.(string)] = value
		return true
	})

	err = json.NewEncoder(file).Encode(m)
	if err != nil {
		return fmt.Errorf("Store() error: %v", err)
	}

	return nil
}

func (c *Cache) Load(ttl time.Duration) error {
	file, err := os.OpenFile("./bag.txt", os.O_RDWR, 0755)
	if err != nil {
		return fmt.Errorf("Load() error: %v", err)
	}

	m := make(map[string]interface{})
	json.NewDecoder(file).Decode(&m)

	for key, value := range m {
		c.Set(key, value, ttl)
	}

	return nil
}
