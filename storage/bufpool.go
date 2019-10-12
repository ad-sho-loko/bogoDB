package storage

import "github.com/ad-sho-loko/bogodb/meta"

type bufferPool struct {
	lru *meta.Lru
}

type pageDescriptor struct {
	dirty bool
	ref uint64
	page *Page
}

func newBufferPool() *bufferPool{
	return &bufferPool{
		lru: meta.NewLru(1000),
	}
}

func (b *bufferPool) toPid(tid uint64) uint64{
	return tid / TupleNumber
}

func (b *bufferPool) pinDirty(pg *pageDescriptor){
	pg.ref++
}

func (b *bufferPool) unpinDirty(pg *pageDescriptor){
	pg.ref--
}

func (b *bufferPool) newPage(){
}

func (b *bufferPool) readPage(tableName string, tid uint64) (*Page, error){
	pgid := b.toPid(tid)
	pd := b.lru.Get(pgid).(*pageDescriptor)
	return pd.page, nil
}

func (b *bufferPool) putPage(pgid uint64, p *Page) (bool, *Page){
	pd := &pageDescriptor{
		page:p,
		ref:0,
		dirty:true,
	}

	victim := b.lru.Insert(pgid, pd).(*pageDescriptor)
	if victim == nil{
		return false, nil
	}

	return victim.dirty, victim.page
}
