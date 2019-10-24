package meta

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestInsertAndGet(t *testing.T) {
	lru := NewLru(1)
	lru.Insert(10, 100)
	v := lru.Get(10)
	assert.Equal(t, v, 100)
}

func TestEvicted(t *testing.T) {
	lru := NewLru(1)
	lru.Insert(10, 100)
	lru.Insert(11, 110)

	recent := lru.Get(11)
	assert.Equal(t, recent, 110)

	old := lru.Get(10)
	assert.Nil(t, old)
}

func TestConcurrency(t *testing.T) {
	lru := NewLru(1000)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			lru.Insert(i, i)
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 1000, lru.Len())
}
