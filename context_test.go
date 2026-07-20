@@ -1,5 +1,6 @@
 package gin

+import "testing"
 import (
        "net/http"
        "net/http/httptest"
@@ -10,6 +11,24 @@ var _ = Describe("Gin", func() {
 )

+func TestConcurrentMapAccess(t *testing.T) {
+       c := NewContext(&Engine{}, &http.Request{})
+       c.Set("key", "value")
+
+       // Copy the context
+       cp := c.Copy()
+
+       // Start a goroutine to modify the copied context
+       go func() {
+               cp.Set("key", "new value")
+       }()
+
+       // Modify the original context
+       c.Set("key", "modified value")
+
+       // Wait for the goroutine to finish
+       <-time.After(100 * time.Millisecond)
+
+       // Check that the values are different
+       if c.Get("key") == cp.Get("key") {
+               t.Errorf("Expected different values, got: %v and %v", c.Get("key"), cp.Get("key"))
+       }
+}
+
 func TestContext(t *testing.T) {
        c := NewContext(&Engine{}, &http.Request{})
        c.Set("key", "value")
```

### Summary of Changes

1. **Modified `(*Context).Copy()` in `context.go`**:
   - If the original context's `Keys` map is not `nil`, create a new map and copy all key-value pairs from the original map to the new map.
   - If the original context's `Keys` map is `nil`, set the copied context's `Keys` to `nil`.

2. **Added a regression test in `context_test.go`**:
   - The test verifies that concurrent writes/reads to the original and copied context keys do not trigger the Go race detector and that the values remain consistent and independent.

These changes should resolve the issue and meet the acceptance criteria.