package cache

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type cache struct {
	data *sync.Map
}

func New() *cache {
	m := new(sync.Map)
	cache := &cache{data: m}
	return cache
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) {
	c.data.Store(key, value)
	time.AfterFunc(ttl, func() {
		if _, ok := c.data.Load(key); ok {
			c.data.Delete(key)
		} else {
			fmt.Printf("Auto-delete error: there is no such key '%s' in storage\n", key)
		}
	})
}

func (c *cache) Get(key string) (interface{}, error) {
	if val, ok := c.data.Load(key); ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("Get error: there is no such key '%s' in storage", key)
	}
}

func (c *cache) Delete(key string) error {
	if _, ok := c.data.Load(key); ok {
		c.data.Delete(key)
		return nil
	} else {
		return fmt.Errorf("Delete error: there is no such key '%s' in storage", key)
	}
}

func (c *cache) Store() error {
	file, err := os.OpenFile("./bag.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	c.data.Range(func(key, value any) bool {
		_, err = file.Write([]byte(fmt.Sprintf("%v: %v\n", key, value)))
		return err == nil
	})
	
	if err != nil {
		return err
	}

	return nil
}