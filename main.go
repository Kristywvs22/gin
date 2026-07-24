package main

import (
	"fmt"
	"sync"
)

type Context struct {
	mu   sync.RWMutex
	Keys map[string]interface{}
}

func (c *Context) Copy() *Context {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Shallow copy the Keys map to prevent data races
	cp := &Context{
		Keys: make(map[string]interface{}, len(c.Keys)),
	}
	for k, v := range c.Keys {
		cp.Keys[k] = v
	}
	return cp
}

func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

func (c *Context) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.Keys[key]
	return val, ok
}

func main() {
	ctx := &Context{}
	ctx.Set("user", "alice")

	ctxCopy := ctx.Copy()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if val, ok := ctxCopy.Get("user"); ok {
			fmt.Printf("Background: user=%v\n", val)
		}
	}()

	ctx.Set("user", "bob")
	wg.Wait()
	fmt.Println("No data race. Keys map is shallow copied.")
}
