package util

import "testing"

func TestArrayJoinInt(t *testing.T) {
	t.Log(ArrayJoinInt([]int{1, 2, 3}, ","))
}
