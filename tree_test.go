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
	tree := Tree{
		Root: &Node{
			Data: map[string]interface{}{
				"ID": "A",
			},
		},
	}
	tree.Root.AddChild(&Node{
		Data: map[string]interface{}{
			"ID": "B",
		},
		Children: []*Node{
			{
				Data: map[string]interface{}{
					"ID": "E",
				},
			},
		},
	})
	tree.Root.AddChild(&Node{
		Data: map[string]interface{}{
			"ID": "C",
		},
	})
	tree.Root.AddChild(&Node{
		Data: map[string]interface{}{
			"ID": "D",
		},
	})
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
}
