// Author: RawNuke
// Copyright (c) 2026 RawNuke. All rights reserved.

package gin

import (
	"fmt"
	"sync"
	"testing"
)

func TestContextCopyKeysRace(t *testing.T) {
	c := NewContext()
	c.Set("foo", "bar")

	cCopy := c.Copy()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			c.Set(fmt.Sprintf("key-%d", i), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			cCopy.Set(fmt.Sprintf("copy-key-%d", i), i)
		}
	}()

	wg.Wait()
}

func TestContextCopyKeysIndependence(t *testing.T) {
	c := NewContext()
	c.Set("foo", "bar")

	cCopy := c.Copy()

	if cCopy.Get("foo") != "bar" {
		t.Errorf("copied context missing key: want 'bar', got %v", cCopy.Get("foo"))
	}

	cCopy.Set("foo", "modified")
	if c.Get("foo") != "bar" {
		t.Errorf("original context modified by copy: want 'bar', got %v", c.Get("foo"))
	}

	c.Set("foo", "original-modified")
	if cCopy.Get("foo") != "modified" {
		t.Errorf("copied context modified by original: want 'modified', got %v", cCopy.Get("foo"))
	}
}

func TestContextCopyNilKeys(t *testing.T) {
	c := NewContext()

	cCopy := c.Copy()

	if cCopy.Keys != nil {
		t.Errorf("copied context Keys should be nil when original Keys is nil, got %v", cCopy.Keys)
	}

	cCopy.Set("foo", "bar")
	if c.Keys != nil {
		t.Errorf("original context Keys should remain nil after copy modification")
	}
}

func TestContextCopyEmptyKeys(t *testing.T) {
	c := NewContext()
	c.Set("foo", "bar")
	c.Set("foo", "")
	c.mu.Lock()
	c.Keys = make(map[string]interface{})
	c.mu.Unlock()

	cCopy := c.Copy()

	if len(cCopy.Keys) != 0 {
		t.Errorf("copied context Keys should be empty, got %d entries", len(cCopy.Keys))
	}
}
