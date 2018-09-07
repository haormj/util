package util

import (
	"fmt"
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
