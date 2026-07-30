package gin

type Context struct {
	Keys map[string]any
}

func (c *Context) Copy() *Context {
	cp := &Context{
		Keys: make(map[string]any, len(c.Keys)),
	}
	for k, v := range c.Keys {
		cp.Keys[k] = v
	}
	return cp
}
