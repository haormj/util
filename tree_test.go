package util

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	/*
		A
		|-B-E
		|-C
		|-D
	*/

	b := NewNode()
	e := NewNode()
	e.Put("ID", "E")
	f := NewNode()
	f.Put("ID", "F")
	b.AddChild(e)
	b.AddChild(f)
	b.Put("ID", "B")

	c := NewNode()
	c.Put("ID", "C")

	d := NewNode()
	d.Put("ID", "D")

	root := NewNode()
	root.Put("ID", "A")
	root.AddChild(b)
	root.AddChild(c)
	root.AddChild(d)
	root.Format = func(n *Node) string {
		return n.Data["ID"].(string)
	}
	tree := Tree{
		Root: root,
	}
	fmt.Println(tree)

	fmt.Println("BFS")
	tree.BFS(func(n *Node) {
		fmt.Print(n.Data, " ")
	})
	fmt.Println("\npre order")
	tree.DFSPreOrder(func(n *Node) {
		fmt.Print(n.Data, " ")
	})
	fmt.Println("\npost order")
	tree.DFSPostOrder(func(n *Node) {
		fmt.Print(n.Data, " ")
	})
	fmt.Println("\ndepth: ", tree.Depth())
	fmt.Println("degree: ", tree.Degree())
	fmt.Println("height: ", tree.Height())
	fmt.Println("level: ", tree.Level())

	n, ok := tree.Search("ID", func(i, j interface{}) bool {
		return i == j
	}, "B", "E")
	fmt.Println(n, ok)

	fmt.Println(tree)

	n = tree.Insert("ID", func(i, j interface{}) bool {
		return i == j
	}, "B", "E", "G")
	fmt.Println("insert")
	fmt.Println(n)

	fmt.Println(tree)
}
