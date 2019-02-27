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
			Entries: []Entry{
				{
					Key:   "ID",
					Value: "A",
				},
			},
		},
	}
	tree.Root.AddChild(&Node{
		Entries: []Entry{
			{
				Key:   "ID",
				Value: "B",
			},
		},
		Children: []*Node{
			{
				Entries: []Entry{
					{
						Key:   "ID",
						Value: "E",
					},
				},
			},
		},
	})
	tree.Root.AddChild(&Node{
		Entries: []Entry{
			{
				Key:   "ID",
				Value: "C",
			},
		},
	})
	tree.Root.AddChild(&Node{
		Entries: []Entry{
			{
				Key:   "ID",
				Value: "D",
			},
		},
	})
	fmt.Println("BFS")
	tree.BFS(func(n *Node) {
		fmt.Print(n.Entries, " ")
	})
	fmt.Println("\npre order")
	tree.DFSPreOrder(func(n *Node) {
		fmt.Print(n.Entries, " ")
	})
	fmt.Println("\npost order")
	tree.DFSPostOrder(func(n *Node) {
		fmt.Print(n.Entries, " ")
	})
	fmt.Println("\ndepth: ", tree.Depth())
	fmt.Println("degree: ", tree.Degree())
	fmt.Println("height: ", tree.Height())
	fmt.Println("level: ", tree.Level())
}
