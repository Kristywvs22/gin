@@ -123,7 +123,10 @@ func (c *Context) Copy() *Context {
        cp.writermem = c.writermem
        cp.Request = c.Request
        cp.Params = c.Params
-       cp.Keys = c.Keys
+       if c.Keys != nil {
+               cp.Keys = make(map[string]any)
+               for k, v := range c.Keys {
+                       cp.Keys[k] = v
+               }
+       } else {
+               cp.Keys = nil
+       }
        return &cp
 }
```

### Regression Test

To ensure that the fix works as expected, we need to add a regression test. Here is the test code:

### Patch for `context_test.go`

```go
diff --git a/context_test.go b/context_test.go
--- a/context_test.go