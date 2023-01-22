package main

import (
	"fmt"
	"strconv"
	"time"

	cache "github.com/cothromachd/in-memory-cache"
)

func main() {
	c := cache.New()
	for i := 0; i < 100; i++ {
		go c.Set(strconv.Itoa(i), i, time.Second*5)
	}
	time.Sleep(time.Second * 2)
	id, _ := c.Get("0")
	id1, _ := c.Get("1")
	fmt.Println(id, id1)
	time.Sleep(time.Second * 4)
	id, _ = c.Get("0")
	id1, _ = c.Get("1")
	fmt.Println(id, id1)
}
