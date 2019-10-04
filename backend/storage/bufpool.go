package storage

import "bogoDB/backend/meta"

type bufferPool struct {
	lru *meta.Lru
}

func newBufferPool() *bufferPool{
	return &bufferPool{
		lru: meta.NewLru(1000),
	}
}

func (b *bufferPool) pinDirty(){}
func (b *bufferPool) unpinDirty(){}
func (b *bufferPool) newPage(){}

func (b *bufferPool) fetchPage(pgid uint64) (*Page, error){
	return b.lru.Get(pgid).(*Page), nil
}