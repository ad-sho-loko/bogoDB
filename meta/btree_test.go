package meta

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

type Int64 int64

func (i Int64) Less(than Item) bool{
	l, ok := than.(Int64)
	if !ok{ return false }
	return i < l
}

func TestNoSplit(t *testing.T){
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

	// test balance
	assert.Equal(t, btree.Top.Items[0], Int64(2))
	assert.Equal(t, btree.Top.Children[0].Items[0], Int64(1))
	assert.Equal(t, btree.Top.Children[1].Items[0], Int64(3))
}

func TestSplitChild(t *testing.T){
	btree := NewBTree()
	btree.Insert(Int64(1))
	btree.Insert(Int64(2))
	btree.Insert(Int64(3))
	btree.Insert(Int64(4))
	btree.Insert(Int64(5))

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

	// test balance
	assert.Equal(t, btree.Top.Items[0], Int64(2))
	assert.Equal(t, btree.Top.Items[1], Int64(4))
	assert.Equal(t, btree.Top.Children[0].Items[0], Int64(1))
	assert.Equal(t, btree.Top.Children[1].Items[0], Int64(3))
	assert.Equal(t, btree.Top.Children[2].Items[0], Int64(5))
}

func TestBlanced(t *testing.T){
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

	// test balance
	assert.Equal(t, btree.Top.Items[0], Int64(4))
	assert.Equal(t, btree.Top.Children[0].Items[0], Int64(2))
	assert.Equal(t, btree.Top.Children[1].Items[0], Int64(6))
	assert.Equal(t, btree.Top.Children[0].Children[0].Items[0], Int64(1))
	assert.Equal(t, btree.Top.Children[0].Children[1].Items[0], Int64(3))
	assert.Equal(t, btree.Top.Children[1].Children[0].Items[0], Int64(5))
	assert.Equal(t, btree.Top.Children[1].Children[1].Items[0], Int64(7))
}

func TestBlancedReversed(t *testing.T){
	btree := NewBTree()
	btree.Insert(Int64(7))
	btree.Insert(Int64(6))
	btree.Insert(Int64(5))
	btree.Insert(Int64(4))
	btree.Insert(Int64(3))
	btree.Insert(Int64(2))
	btree.Insert(Int64(1))

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

	// test balance
	assert.Equal(t, btree.Top.Items[0], Int64(4))
	assert.Equal(t, btree.Top.Children[0].Items[0], Int64(2))
	assert.Equal(t, btree.Top.Children[1].Items[0], Int64(6))
	assert.Equal(t, btree.Top.Children[0].Children[0].Items[0], Int64(1))
	assert.Equal(t, btree.Top.Children[0].Children[1].Items[0], Int64(3))
	assert.Equal(t, btree.Top.Children[1].Children[0].Items[0], Int64(5))
	assert.Equal(t, btree.Top.Children[1].Children[1].Items[0], Int64(7))
}

func TestGet(t *testing.T){
	btree := NewBTree()
	btree.Insert(Int64(1))
	btree.Insert(Int64(2))
	btree.Insert(Int64(3))
	btree.Insert(Int64(4))
	btree.Insert(Int64(5))
	btree.Insert(Int64(6))
	btree.Insert(Int64(7))

	item := btree.Get(Int64(1))
	i := item.(Int64)
	assert.Equal(t, i, Int64(1))

	item = btree.Get(Int64(7))
	i = item.(Int64)
	assert.Equal(t, i, Int64(7))
}

func TestRandom(t *testing.T){
	btree := NewBTree()

	for i:=0; i<10000; i++{
		v := rand.Intn(1000)
		btree.Insert(Int64(v))
	}

	assert.Equal(t, btree.Len(), 10000)
}

func TestEmpty(t *testing.T){
	btree := NewBTree()
	found, _ := btree.Find(Int64(1))
	assert.False(t, found)
}
