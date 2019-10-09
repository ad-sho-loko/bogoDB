package meta

// Inspired by google/btree
// https://www.cs.usfca.edu/~galles/visualization/BTree.html
const(
	maxDegree = 3
)

type BTree struct {
	top items
}

type items []Item

func (i *items) find(item Item) bool{
	for _, i := range *i{
		if i == item{
			return true
		}
	}
	return false
}

func (i *items) insertAt(index int, item Item){
	*i = append(*i, nil)
	if len(*i) >= maxDegree{
		// should rotate
	}
	(*i)[index] = item
}

type Item interface {
	Less(than Item)
}

func NewBTree() *BTree{
	return &BTree{
	}
}

func (b *BTree) Insert(item Item){
	// index, found := b.Find(item)
	// b.top.insertAt(index, item)
}

func (b *BTree) Find(item Item) (int, bool){
	return 0, true
}

func (b *BTree) Delete(item Item){
}