package main

import (
	"fmt"
	"sync"
)

// Context with thread-safe Copy.
// The bug: (*Context).Copy() shallow copies the Keys map, causing data races.
// Fix: Deep copy the Keys map in Copy().
type Context struct {
	mu   sync.RWMutex
	Keys map[string]interface{}
}

func (c *Context) Set(key string, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = val
}

func (c *Context) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Keys[key]
}

func (c *Context) Copy() *Context {
	c.mu.RLock()
	defer c.mu.RUnlock()
	copied := &Context{Keys: make(map[string]interface{}, len(c.Keys))}
	for k, v := range c.Keys {
		copied.Keys[k] = v
	}
	return copied
}

func main() {
	ctx := &Context{}
	ctx.Set("user", "alice")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); ctx.Set("key", "val") }()
		go func() { defer wg.Done(); _ = ctx.Copy() }()
	}
	wg.Wait()
	fmt.Println("No data race - Context.Copy() deep copies Keys map")
}
