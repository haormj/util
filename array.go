package util

import (
	"fmt"
	"sort"
	"strings"
)

// ArrayReverseInt64 int64 array reverse
func ArrayReverseInt64(a []int64) []int64 {
	l := len(a)
	b := make([]int64, l)
	for i := l; i > 0; i-- {
		b[l-i] = a[i-1]
	}
	return b
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
