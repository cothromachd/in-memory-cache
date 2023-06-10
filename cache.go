package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cacher interface {
	Cap() int
	Clear()
	Add(key string, value interface{})
	AddWithTTL(key string, value interface{}, ttl time.Duration)
	Get(key string) (value interface{}, ok bool)
	Remove(key string) bool
}

type Node struct {
	Data       interface{}
	KeyPtr     *list.Element
	Expiration time.Time
}

type Cache struct {
	mx       sync.Mutex
	Queue    *list.List
	Items    map[string]*Node
	Capacity int
	ttl      time.Duration
}

func New(capacity int, ttl time.Duration) Cache {
	return Cache{
		mx:       sync.Mutex{},
		Queue:    list.New(),
		Items:    make(map[string]*Node),
		Capacity: capacity,
		ttl:      ttl}
}

func (c *Cache) Add(key string, value interface{}) {
	c.mx.Lock()

	for key, val := range c.Items {
		if time.Now().After(val.Expiration) {
			c.Queue.Remove(val.KeyPtr)
			delete(c.Items, key)
		}
	}

	if item, ok := c.Items[key]; !ok {
		if c.Capacity == len(c.Items) {
			back := c.Queue.Back()
			c.Queue.Remove(back)
			delete(c.Items, back.Value.(string))
		}

		c.Items[key] = &Node{Data: value, KeyPtr: c.Queue.PushFront(key), Expiration: time.Now().Add(c.ttl)}
	} else {
		item.Data = value
		item.Expiration = time.Now().Add(c.ttl)
		c.Items[key] = item
		c.Queue.MoveToFront(item.KeyPtr)
	}

	c.mx.Unlock()
}

func (c *Cache) AddWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mx.Lock()

	for key, val := range c.Items {
		if time.Now().After(val.Expiration) {
			c.Queue.Remove(val.KeyPtr)
			delete(c.Items, key)
		}
	}

	if item, ok := c.Items[key]; !ok {
		if c.Capacity == len(c.Items) {
			back := c.Queue.Back()
			c.Queue.Remove(back)
			delete(c.Items, back.Value.(string))
		}

		c.Items[key] = &Node{Data: value, KeyPtr: c.Queue.PushFront(key), Expiration: time.Now().Add(ttl)}
	} else {
		item.Data = value
		item.Expiration = time.Now().Add(ttl)
		c.Items[key] = item
		c.Queue.MoveToFront(item.KeyPtr)
	}

	c.mx.Unlock()
}

func (c *Cache) Get(key string) (value interface{}, ok bool) {
	c.mx.Lock()

	defer c.mx.Unlock()

	if item, ok := c.Items[key]; ok {
		if !time.Now().After(item.Expiration) {
			c.Queue.MoveToFront(item.KeyPtr)
			value = item.Data
			return value, true
		} else {
			c.Queue.Remove(item.KeyPtr)
			delete(c.Items, key)
		}
	}
	return nil, false
}

func (c *Cache) Remove(key string) bool {
	c.mx.Lock()

	defer c.mx.Unlock()

	elem, ok := c.Items[key]
	if !ok {
		return false
	}

	c.Queue.Remove(elem.KeyPtr)
	delete(c.Items, key)
	return true
}

func (c *Cache) Clear() int {
	var cnt int
	var ok bool

	c.mx.Lock()
	for key := range c.Items {
		c.mx.Unlock()
		ok = c.Remove(key)
		if ok {
			cnt++
		} else {
			c.mx.Lock()
			return -1
		}

		c.mx.Lock()
	}

	defer c.mx.Unlock()
	return cnt
}

func (c *Cache) Cap() int {
	return c.Capacity
}
