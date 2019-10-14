package meta

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Int64 int64

func (i Int64) Less(than item) bool{
	l, ok := than.(Int64)
	if !ok{ return false }
	return i < l
}

func TestNoSplit(t *testing.T){
	t.Skip()
	btree := NewBTree()

	btree.Insert(Int64(1))
	btree.Insert(Int64(2))

	found, _ := btree.Find(Int64(1))
	assert.True(t, found)

	found, _ = btree.Find(Int64(2))
	assert.True(t, found)

	found, _ = btree.Find(Int64(3))
	assert.False(t, found)
}

func TestSplitParent(t *testing.T){
	t.Skip()
	btree := NewBTree()

	btree.Insert(Int64(1))
	btree.Insert(Int64(2))
	btree.Insert(Int64(3))

	found, _ := btree.Find(Int64(1))
	assert.True(t, found)

	found, _ = btree.Find(Int64(2))
	assert.True(t, found)

	found, _ = btree.Find(Int64(3))
	assert.True(t, found)
}

func TestSplitChild(t *testing.T){
	btree := NewBTree()
	btree.Insert(Int64(1))
	btree.Insert(Int64(2))
	btree.Insert(Int64(3))
	btree.Insert(Int64(4))
	btree.Insert(Int64(5))
	btree.Insert(Int64(6))
	btree.Insert(Int64(7))

	found, _ := btree.Find(Int64(1))
	assert.True(t, found)

	found, _ = btree.Find(Int64(2))
	assert.True(t, found)

	found, _ = btree.Find(Int64(3))
	assert.True(t, found)

	found, _ = btree.Find(Int64(4))
	assert.True(t, found)

	found, _ = btree.Find(Int64(5))
	assert.True(t, found)

	found, _ = btree.Find(Int64(6))
	assert.True(t, found)

	found, _ = btree.Find(Int64(7))
	assert.True(t, found)
}

func TestEmpty(t *testing.T){
	btree := NewBTree()
	found, _ := btree.Find(Int64(1))
	assert.False(t, found)
}

