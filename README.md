# In memory cache Go package

An in-memory cache is a data storage layer that sits between applications and databases to deliver responses with high speeds. 

### Installation
```sh
go get -u "github.com/cothromachd/in-memory-cache"
```

### Example
```sh
package main

import (
	"fmt"
	"github.com/cothromachd/in-memory-cache"
)

func main() {
	c := cache.New()

	c.Set("userId", 12)

	id, err := c.Get("userId")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(id)
	}
	// Output: 12

	err = c.Delete("userId")
	if err != nil {
		fmt.Println(err)
	}

	id, err = c.Get("userId")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(id)
	}
	// Output: "Get error: there is no such key 'userId' in storage"

	err = c.Delete("userId")
	if err != nil {
		fmt.Println(err)
	}
	// Output: "Delete error: there is no such key 'userId' in storage"
}
```