package meta

import "sort"

const(
	maxDegree = 3
)

type BTree struct {
	top *node
	length int
}

type node struct {
	items items
	children []*node
}

type items []Item

// Item must be comarable for b-tree implementation.
type Item interface {
	Less(than Item) bool
}

func (i *items) find(item Item) (bool, int){
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

func (i *items) insertAt(index int, item Item){
	*i = append(*i, nil)
	if index < len(*i) {
		copy((*i)[index+1:], (*i)[index:])
	}
	(*i)[index] = item
}

func NewBTree() *BTree{
	return &BTree{
		top:nil,
		length:0,
	}
}

func (n *node) deleteChildAt(index int){
	first := n.children[:index]
	second := n.children[index+1:]
	n.children = append(first, second...)
}

func (n *node) splitChild(index int, item Item){
	_, innerIndex := n.children[index].items.find(item)
	n.children[index].items.insertAt(innerIndex, item)

	leftItem := n.children[index].items[maxDegree/2-1]
	midItem := n.children[index].items[maxDegree/2]
	rightItem := n.children[index].items[maxDegree/2+1]

	n.deleteChildAt(index)

	_, midIndex := n.items.find(midItem)
	n.items.insertAt(midIndex, midItem)

	left := new(node)
	left.insert(leftItem)

	right := new(node)
	right.insert(rightItem)

	n.children = append(n.children, left)
	n.children = append(n.children, right)
	sort.Slice(n.children, func(i, j int) bool{
		return n.children[i].items[0].Less(n.children[j].items[0])
	})
}

func (n *node) splitMe(){
	left := new(node)
	left.items.insertAt(0, n.items[maxDegree/2-1])

	right := new(node)
	right.items.insertAt(0, n.items[maxDegree/2+1])

	mid := n.items[maxDegree/2]
	n.items = append([]Item{}, mid)

	if len(n.children) == maxDegree + 1{
		var nodes []*node

		left.children = append(left.children, n.children[0])
		left.children = append(left.children, n.children[1])
		nodes = append(nodes, left)

		right.children = append(right.children, n.children[2])
		right.children = append(right.children, n.children[3])
		nodes = append(nodes, right)
		n.children = nodes
	}else{
		n.children = append(n.children, left)
		n.children = append(n.children, right)
	}
}

func (n *node) insert(item Item){
	found, index := n.items.find(item)
	if found{
		return
	}

	if len(n.children) == 0{
		n.items.insertAt(index, item)

		if len(n.items) == maxDegree{
			n.splitMe()
		}

		return
	}

	if len(n.children[index].items) == maxDegree - 1{
		n.splitChild(index, item)

		if len(n.items) == maxDegree{
			n.splitMe()
		}

		return
	}

	n.children[index].insert(item)
}

func (n *node) get(key Item) Item{
	found, i := n.items.find(key)
	if found {
		return n.items[i]
	} else if len(n.children) > 0 {
		return n.children[i].get(key)
	}
	return nil
}

func (n *node) find(item Item) (bool, int){
	found, index := n.items.find(item)
	if found{
		return found, index
	}

	if len(n.children) == 0{
		return false, -1
	}

	return n.children[index].find(item)
}

func (b *BTree) Insert(item Item){
	b.length++

	if b.top == nil{
		b.top = new(node)
		b.top.items.insertAt(0, item)
		return
	}

	b.top.insert(item)
}

func (b *BTree) Find(item Item) (bool, int){
	if b.top == nil{
		return false, -1
	}

	return b.top.find(item)
}

func (b *BTree) Get(key Item) Item {
	if b.top == nil{
		return nil
	}

	return b.top.get(key)
}

func (b *BTree) Len() int{
	return b.length
}

func SerializeBTree(tree *BTree) ([]byte, error){
	return nil, nil
}

func DeserializeBTree(bytes []byte) (*BTree, error){
	return nil, nil
}