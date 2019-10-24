package meta

import (
	"encoding/json"
	"sort"
)

const (
	maxDegree = 3
)

type BTree struct {
	Top    *node `json:"top"`
	Length int   `json:"length"`
}

type node struct {
	Items    items   `json:"items"`
	Children []*node `json:"children"`
}

type items []Item
type IntItem int32

// Item must be comarable for b-tree implementation.
type Item interface {
	Less(than Item) bool
}

func (i IntItem) Less(than Item) bool {
	v, ok := than.(IntItem)
	if !ok {
		return false
	}
	return i < v
}

func (i *items) find(item Item) (bool, int) {
	for index, itm := range *i {
		if !itm.Less(item) {
			if !item.Less(itm) {
				return true, index
			}
			return false, index
		}
	}
	return false, len(*i)
}

func (i *items) insertAt(index int, item Item) {
	*i = append(*i, nil)
	if index < len(*i) {
		copy((*i)[index+1:], (*i)[index:])
	}
	(*i)[index] = item
}

func NewBTree() *BTree {
	return &BTree{
		Top:    nil,
		Length: 0,
	}
}

func (n *node) deleteChildAt(index int) {
	first := n.Children[:index]
	second := n.Children[index+1:]
	n.Children = append(first, second...)
}

func (n *node) splitChild(index int, item Item) {
	_, innerIndex := n.Children[index].Items.find(item)
	n.Children[index].Items.insertAt(innerIndex, item)

	leftItem := n.Children[index].Items[maxDegree/2-1]
	midItem := n.Children[index].Items[maxDegree/2]
	rightItem := n.Children[index].Items[maxDegree/2+1]

	n.deleteChildAt(index)

	_, midIndex := n.Items.find(midItem)
	n.Items.insertAt(midIndex, midItem)

	left := new(node)
	left.insert(leftItem)

	right := new(node)
	right.insert(rightItem)

	n.Children = append(n.Children, left)
	n.Children = append(n.Children, right)
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].Items[0].Less(n.Children[j].Items[0])
	})
}

func (n *node) splitMe() {
	left := new(node)
	left.Items.insertAt(0, n.Items[maxDegree/2-1])

	right := new(node)
	right.Items.insertAt(0, n.Items[maxDegree/2+1])

	mid := n.Items[maxDegree/2]
	n.Items = append([]Item{}, mid)

	if len(n.Children) == maxDegree+1 {
		var nodes []*node

		left.Children = append(left.Children, n.Children[0])
		left.Children = append(left.Children, n.Children[1])
		nodes = append(nodes, left)

		right.Children = append(right.Children, n.Children[2])
		right.Children = append(right.Children, n.Children[3])
		nodes = append(nodes, right)
		n.Children = nodes
	} else {
		n.Children = append(n.Children, left)
		n.Children = append(n.Children, right)
	}
}

func (n *node) insert(item Item) {
	found, index := n.Items.find(item)
	if found {
		return
	}

	if len(n.Children) == 0 {
		n.Items.insertAt(index, item)

		if len(n.Items) == maxDegree {
			n.splitMe()
		}

		return
	}

	if len(n.Children[index].Items) == maxDegree-1 {
		n.splitChild(index, item)

		if len(n.Items) == maxDegree {
			n.splitMe()
		}

		return
	}

	n.Children[index].insert(item)
}

func (n *node) get(key Item) Item {
	found, i := n.Items.find(key)
	if found {
		return n.Items[i]
	} else if len(n.Children) > 0 {
		return n.Children[i].get(key)
	}
	return nil
}

func (n *node) find(item Item) (bool, int) {
	found, index := n.Items.find(item)
	if found {
		return found, index
	}

	if len(n.Children) == 0 {
		return false, -1
	}

	return n.Children[index].find(item)
}

func (b *BTree) Insert(item Item) {
	b.Length++

	if b.Top == nil {
		b.Top = new(node)
		b.Top.Items.insertAt(0, item)
		return
	}

	b.Top.insert(item)
}

func (b *BTree) Find(item Item) (bool, int) {
	if b.Top == nil {
		return false, -1
	}

	return b.Top.find(item)
}

func (b *BTree) Get(key Item) Item {
	if b.Top == nil {
		return nil
	}

	return b.Top.get(key)
}

func (b *BTree) Len() int {
	return b.Length
}

func SerializeBTree(tree *BTree) ([]byte, error) {
	return json.Marshal(tree)
}

func DeserializeBTree(bytes []byte) (*BTree, error) {
	var tree BTree
	err := json.Unmarshal(bytes, &tree)
	return &tree, err
}

func (i items) MarshalJSON() ([]byte, error) {
	var intItems []IntItem
	for _, item := range i {
		intItems = append(intItems, item.(IntItem))
	}
	return json.Marshal(intItems)
}

func (i *items) UnmarshalJSON(b []byte) error {
	var intItems []IntItem
	err := json.Unmarshal([]byte(b), &intItems)

	for index, item := range intItems {
		i.insertAt(index, item)
	}

	return err
}
