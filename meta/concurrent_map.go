package meta

import "sync"

type ConcurrentMap struct {
	mu sync.RWMutex
	m  map[interface{}]interface{}
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		m: make(map[interface{}]interface{}),
	}
}

func (c *ConcurrentMap) Put(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

func (c *ConcurrentMap) Get(key interface{}) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.Unlock()
	v, found := c.m[key]
	return v, found
}
