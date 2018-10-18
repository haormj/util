package util

import (
	"fmt"
	"sort"
	"strings"
)

// ArrayReverseInt64 int64 array reverse
func ArrayReverseInt64(a []int64) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// ArrayJoinInt join int array
func ArrayJoinInt(a []int, sep string) string {
	str := ""
	for _, i := range a {
		str += fmt.Sprintf("%d%s", i, sep)
	}
	str = strings.Trim(str, sep)
	return str
}

// ArrayJoinUint32 join uint32 array
func ArrayJoinUint32(a []uint32, sep string) string {
	str := ""
	for _, i := range a {
		str += fmt.Sprintf("%d%s", i, sep)
	}
	str = strings.Trim(str, sep)
	return str
}

// ArrayDistinctInt array distinct
func ArrayDistinctInt(src []int) []int {
	dst := make([]int, 0)
	for _, s := range src {
		flag := false
		for _, d := range dst {
			if s == d {
				flag = true
				break
			}
		}
		if !flag {
			dst = append(dst, s)
		}
	}
	return dst
}

// ArraySortAscInt sort int array asc
func ArraySortAscInt(src []int) {
	sort.Slice(src, func(i int, j int) bool {
		if src[i] < src[j] {
			return true
		} else {
			return false
		}
	})
}

// ArraySortDescInt sort int array desc
func ArraySortDescInt(src []int) {
	sort.Slice(src, func(i int, j int) bool {
		if src[i] > src[j] {
			return true
		} else {
			return false
		}
	})
}
