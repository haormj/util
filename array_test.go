package util

import "testing"

func TestArrayJoinInt(t *testing.T) {
	t.Log(ArrayJoinInt([]int{1, 2, 3}, ","))
}

func TestArrayDistinctInt(t *testing.T) {
	t.Log(ArrayDistinctInt([]int{1, 2, 3, 1, 2}))
}
