package cache

import (
	"container/list"
	"log"
	"sync"
)

type cacheItem[V any] struct {
	key   string
	value *V
}

// cache uses a FIFO list to keep the ordering by age
// and a map to keep the values for quick access/update/delete
type Cache[V any] struct {
	maxSize int
	mu      sync.RWMutex
	items   map[string]*list.Element
	queue   *list.List
}

func NewCache[V any](maxSize int) *Cache[V] {
	if maxSize <= 0 {
		log.Fatal("can't setup a cache with this size")
	}
	return &Cache[V]{
		maxSize: maxSize,
		items:   make(map[string]*list.Element),
		queue:   list.New(),
	}
}

func (c *Cache[V]) Get(key string) (*V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	element, exists := c.items[key]
	if !exists {
		return nil, false
	}

	return element.Value.(*cacheItem[V]).value, true
}

func (c *Cache[V]) Add(key string, value *V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.items[key]; exists {
		element.Value.(*cacheItem[V]).value = value
		return
	}

	// delete only if it's full
	if c.queue.Len() >= c.maxSize {
		c.evictOldest()
	}

	item := &cacheItem[V]{key: key, value: value}
	// push at the back
	element := c.queue.PushBack(item)
	c.items[key] = element
}

func (c *Cache[V]) evictOldest() {
	// remove at the front
	oldestElement := c.queue.Front()
	if oldestElement == nil {
		return
	}

	c.queue.Remove(oldestElement)

	item := oldestElement.Value.(*cacheItem[V])
	delete(c.items, item.key)
}

func (c *Cache[V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.queue.Len()
}
