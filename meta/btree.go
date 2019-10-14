package meta


const(
	maxDegree = 3
)

// A b-tree similar to google/btree
// https://www.cs.usfca.edu/~galles/visualization/BTree.html
type BTree struct {
	top *node
	length int
}

type node struct {
	items items
	children []*node
}

type items []item

// item must be comarable for b-tree implementation.
type item interface {
	Less(than item) bool
}
func (i *items) find(item item) (bool, int){
	for index, itm := range *i{
		if !itm.Less(item){
			if !item.Less(itm){
				return true, index
			}
			return false, index
		}
	}
	return false, len(*i)
}

func (i *items) insertAt(index int, item item){
	*i = append(*i, nil)
	(*i)[index] = item
}

func NewBTree() *BTree{
	return &BTree{
		top:nil,
		length:0,
	}
}

func (n *node) splitChild(i int){
	mid := new(node)
	mid.insert(n.children[i].items[maxDegree/2])

	// mid.items = append(mid.children[i], n.children[i].items[maxDegree/2-1])
	// mid.items = append(mid.children[i], n.children[i].items[maxDegree/2+1])

	// n.items[i] = mid
}

func (n *node) splitMe(){
	left := new(node)
	left.items.insertAt(0, n.items[maxDegree/2-1])

	right := new(node)
	right.items.insertAt(0, n.items[maxDegree/2+1])

	mid := n.items[maxDegree/2]
	n.items = append([]item{}, mid)

	n.children = append(n.children, left)
	n.children = append(n.children, right)
}

func (n *node) insert(item item){
	found, index := n.items.find(item)
	if found{
		return
	}

	if len(n.items) == maxDegree - 1{
		n.items.insertAt(index, item)
		n.splitMe()
		return
	}

	if len(n.children) == 0{
		n.items.insertAt(index, item)
		return
	}

	if len(n.children[index].items) == maxDegree - 1{
		n.splitChild(index)
		return
	}

	n.children[index].insert(item)
}

func (n *node) find(item item) (bool, int){
	found, index := n.items.find(item)
	if found{
		return found, index
	}

	if len(n.children) == 0{
		return false, -1
	}

	return n.children[index].find(item)
}

func (b *BTree) Insert(item item){
	if b.top == nil{
		b.top = new(node)
		b.top.items.insertAt(0, item)
		return
	}

	b.top.insert(item)
	b.length++
}

func (b *BTree) Find(item item) (bool, int){
	if b.top == nil{
		return false, -1
	}

	return b.top.find(item)
}

func (b *BTree) Len() int{
	return b.length
}

// func MarshallBTree(b *BTree) []byte{
// }
