package util

import "testing"

func TestArrayReverseInt64(t *testing.T) {
	a := []int64{1, 2}
	ArrayReverseInt64(a)
	t.Log(a)
}

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

func TestArrayJoinUint32(t *testing.T) {
	t.Log(ArrayJoinUint32([]uint32{1, 3, 2}, ","))
}

func TestArrayContainsString(t *testing.T) {
	t.Log(ArrayContainsString([]string{"hello", "world"}, "hello"))
	t.Log(ArrayContainsString([]string{"hello", "world"}, "nihao"))
}

func TestArrayIntersectString(t *testing.T) {
	t.Log(ArrayIntersectString([]string{"hello", "world"}, []string{"hello", "nihao"}))
}

func TestArrayUnionString(t *testing.T) {
	t.Log(ArrayUnionString([]string{"hello", "world"}, []string{"hello", "nihao"}))
}

func TestArraySymDiffString(t *testing.T) {
	t.Log(ArraySymDiffString([]string{"1", "2", "3", "4", "5"}, []string{"4", "5", "6", "7", "8"}))
}
