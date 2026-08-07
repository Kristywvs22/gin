// Author: RawNuke
// Copyright (c) 2026 RawNuke. All rights reserved.

package gin

import "sync"

type Context struct {
	writermem    responseWriter
	Writer       ResponseWriter
	Request      interface{}
	Params       interface{}
	index        int
	handlers     interface{}
	params       interface{}
	skippedNodes interface{}
	Keys         map[string]interface{}
	mu           sync.RWMutex
}

type ResponseWriter interface{}
type responseWriter struct{ Size, Status int }

func NewContext() *Context {
	return &Context{}
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
	if c.Keys == nil {
		return nil
	}
	return c.Keys[key]
}

func (c *Context) Copy() *Context {
	c.mu.RLock()
	keys := c.Keys
	c.mu.RUnlock()

	cp := &Context{
		index:        abortIndex,
		handlers:     nil,
		params:       nil,
		skippedNodes: nil,
	}

	if keys != nil {
		cp.Keys = make(map[string]interface{}, len(keys))
		for k, v := range keys {
			cp.Keys[k] = v
		}
	}

	return cp
}

const abortIndex = -1
