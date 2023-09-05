package cache

import (
	"sync"
)

// Cache represents a simple cache structure
type Cache struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

// New creates a new Cache instance
func New() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Set sets a key-value pair in the cache
func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

// Get retrieves a value based on the key from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	val, ok := c.data[key]
	return val, ok
}
