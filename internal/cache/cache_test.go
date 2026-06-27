//go:build unit

package cache

import (
	"testing"
)

func TestCache_Functionality(t *testing.T) {
	c := NewCache[string](3)

	valA := "A"
	valB := "B"
	valC := "C"
	valD := "D"
	valANew := "A_Updated"

	t.Run("Add and Get Items", func(t *testing.T) {
		c.Add("keyA", &valA)
		c.Add("keyB", &valB)

		if got, found := c.Get("keyA"); !found || *got != valA {
			t.Errorf("Expected to find keyA with value %s, got found=%v, value=%v", valA, found, got)
		}

		if got, found := c.Get("keyB"); !found || *got != valB {
			t.Errorf("Expected to find keyB with value %s, got found=%v, value=%v", valB, found, got)
		}

		if _, found := c.Get("keyC"); found {
			t.Errorf("Expected keyC to not exist yet")
		}

		if c.Size() != 2 {
			t.Errorf("Expected cache size to be 2, got %d", c.Size())
		}
	})

	t.Run("Update Existing Key", func(t *testing.T) {
		c.Add("keyA", &valANew)

		got, found := c.Get("keyA")
		if !found || *got != valANew {
			t.Errorf("Expected updated value %s for keyA, got found=%v, value=%v", valANew, found, got)
		}

		// Size should remain 2 because it was an in-place update
		if c.Size() != 2 {
			t.Errorf("Expected cache size to still be 2 after update, got %d", c.Size())
		}
	})

	t.Run("Eviction Strategy", func(t *testing.T) {
		c.Add("keyC", &valC)
		c.Add("keyD", &valD)

		if _, found := c.Get("keyA"); found {
			t.Errorf("Expected keyA to be evicted, but it was found")
		}

		keys := []string{"keyB", "keyC", "keyD"}
		for _, k := range keys {
			if _, found := c.Get(k); !found {
				t.Errorf("Expected %s to still be in the cache", k)
			}
		}

		if c.Size() != 3 {
			t.Errorf("Expected cache size to be capped at 3, got %d", c.Size())
		}
	})
}
