package util

import (
	"container/list"
	"fmt"
	"strings"
)

// Tree definition
type Tree struct {
	Root *Node
}

// Depth the largest depth of the node
func (t *Tree) Depth() int {
	var d int
	t.BFS(func(n *Node) {
		dep := n.Depth()
		if dep > d {
			d = dep
		}
	})
	return d
}

// Height is the length of the longest path to a leaf
func (t *Tree) Height() int {
	return t.Root.Height()
}

// Degree the largest degree of the node
func (t *Tree) Degree() int {
	var d int
	t.BFS(func(n *Node) {
		deg := n.Degree()
		if deg > d {
			d = deg
		}
	})
	return d
}

// Level 1 + the number of edges between the node and the root
func (t *Tree) Level() int {
	return t.Depth() + 1
}

// DFSPreOrder Depth-First Search (DFS) is an algorithm for traversing or searching tree data structure.
// One starts at the root and explores as far as possible along each branch before backtracking.
// parent comes before children; overall root first
func (t *Tree) DFSPreOrder(fn func(*Node)) {
	t.Root.DFSPreOrder(fn)
}

// DFSPostOrder  Depth-First Search (DFS) is an algorithm for traversing or searching tree data structure.
// One starts at the root and explores as far as possible along each branch before backtracking.
// parent comes after children; overall root last
func (t *Tree) DFSPostOrder(fn func(*Node)) {
	t.Root.DFSPostOrder(fn)
}

// BFS Breadth-First Search (BFS) is an algorithm for traversing or searching tree data structure.
// It starts at the tree root and explores the neighbor nodes first, before moving to the next level neighbors.
func (t *Tree) BFS(fn func(*Node)) {
	t.Root.BFS(fn)
}

// Search node by key
func (t *Tree) Search(key string, compare func(interface{}, interface{}) bool,
	values ...string) (*Node, bool) {
	return t.Root.Search(key, compare, values...)
}

func (t *Tree) Insert(key string, compare func(interface{}, interface{}) bool,
	values ...string) *Node {
	return t.Root.Insert(key, compare, values...)
}

func (t *Tree) String() string {
	return t.Root.String()
}

// Node definition
type Node struct {
	Parent   *Node
	Data     map[string]interface{}
	Children []*Node
	Format   func(*Node) string
}

func NewNode() *Node {
	return &Node{
		Parent:   nil,
		Data:     make(map[string]interface{}),
		Children: make([]*Node, 0),
		Format: func(n *Node) string {
			return fmt.Sprintf("%+v", n.Data)
		},
	}
}

// Depth the number of edges between the node and the root.
func (n *Node) Depth() int {
	if n.Parent == nil {
		return 0
	}
	return n.Parent.Depth() + 1
}

// Height the number of edges on the longest path between that node and a leaf.
func (n *Node) Height() int {
	if len(n.Children) == 0 {
		return 0
	}
	var max int
	for _, c := range n.Children {
		h := c.Height()
		if max < h {
			max = h
		}
	}
	return max + 1
}

// Degree the number of children. A leaf is necessarily degree zero.
func (n *Node) Degree() int {
	return len(n.Children)
}

// Level 1 + the number of edges between the Node and the root
func (n *Node) Level() int {
	return n.Depth() + 1
}

// DFSPreOrder Depth-First Search (DFS) is an algorithm for traversing or searching tree data structure.
// One starts at the root and explores as far as possible along each branch before backtracking.
// parent comes before children; overall root first
func (n *Node) DFSPreOrder(fn func(*Node)) {
	fn(n)
	for _, c := range n.Children {
		c.DFSPreOrder(fn)
	}
}

// DFSPostOrder  Depth-First Search (DFS) is an algorithm for traversing or searching tree data structure.
// One starts at the root and explores as far as possible along each branch before backtracking.
// parent comes after children; overall root last
func (n *Node) DFSPostOrder(fn func(*Node)) {
	for i := len(n.Children); i > 0; i-- {
		n.Children[i-1].DFSPostOrder(fn)
	}
	fn(n)
}

// BFS Breadth-First Search (BFS) is an algorithm for traversing or searching tree data structure.
// It starts at the tree root and explores the neighbor nodes first, before moving to the next level neighbors.
func (n *Node) BFS(fn func(*Node)) {
	l := list.New()
	l.PushBack(n)
	for e := l.Front(); e != nil; e = e.Next() {
		ele := e.Value.(*Node)
		fn(ele)
		for _, c := range ele.Children {
			l.PushBack(c)
		}
	}
}

// AddChild -
func (n *Node) AddChild(c *Node) {
	c.Parent = n
	for _, cc := range c.Children {
		cc.Parent = c
	}
	if c.Format == nil {
		c.Format = n.Format
	}
	n.Children = append(n.Children, c)
}

func (n *Node) Put(key string, value interface{}) {
	n.Data[key] = value
}

func (n *Node) Get(key string) (interface{}, bool) {
	value, ok := n.Data[key]
	return value, ok
}

func (n *Node) Search(key string, compare func(interface{}, interface{}) bool,
	values ...string) (*Node, bool) {
	if len(values) == 0 {
		return nil, false
	}
	value := values[0]
	values = values[1:]
	for _, child := range n.Children {
		v, ok := child.Get(key)
		if !ok {
			continue
		}
		if compare(v, value) {
			if len(values) == 0 {
				return child, true
			} else {
				return child.Search(key, compare, values...)
			}
		}
	}
	return nil, false
}

func (n *Node) Insert(key string, compare func(interface{}, interface{}) bool,
	values ...string) *Node {
	if len(values) == 0 {
		return n
	}
	value := values[0]
	values = values[1:]
	var flag bool
	var child *Node
	for _, child = range n.Children {
		v, ok := child.Get(key)
		if !ok {
			continue
		}
		if compare(v, value) {
			flag = true
			break
		}
	}
	// node don't exist
	var node *Node
	if !flag {
		node = NewNode()
		node.Put(key, value)
		n.AddChild(node)
	} else {
		node = child
	}
	if len(values) == 0 {
		return node
	} else {
		return node.Insert(key, compare, values...)
	}
}

func (n *Node) stringWithFormat(fn func(*Node) string) string {
	str := strings.Repeat("    ", n.Depth()) + fn(n)
	for _, child := range n.Children {
		str += "\n" + child.stringWithFormat(fn)
	}
	return str
}

// String
func (n *Node) String() string {
	return n.stringWithFormat(n.Format)
}
