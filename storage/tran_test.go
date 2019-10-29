package storage

import (
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
)

func TestTxidAtomicity(t *testing.T) {
	var wg sync.WaitGroup
	var exists [100001]bool
	manager := &TransactionManager{}

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			id := manager.newTxid()

			if exists[id] {
				log.Fatal("txid duplicated")
			}

			exists[id] = true
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, uint64(100000), manager.currentTxid)
}
