package util

import "testing"

func TestArrayJoinInt(t *testing.T) {
	t.Log(ArrayJoinInt([]int{1, 2, 3}, ","))
}

func TestArrayDistinctInt(t *testing.T) {
	t.Log(ArrayDistinctInt([]int{1, 2, 3, 1, 2}))
}

func TestArraySortAscInt(t *testing.T) {
	a := []int{3, 1, 2}
	ArraySortAscInt(a)
	t.Log(a)
}

func TestArraySortDescInt(t *testing.T) {
	a := []int{3, 1, 2}
	ArraySortDescInt(a)
	t.Log(a)
}
