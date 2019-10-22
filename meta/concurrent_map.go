package meta

import "sync"

type ConcurrentMap struct {
	mu sync.RWMutex
	m map[interface{}]interface{}
}

func (c *ConcurrentMap) Put(key, value interface{}){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

func (c *ConcurrentMap) Get(key interface{}) interface{}{
	c.mu.RLock()
	defer c.mu.Unlock()
	return c.m[key]
}