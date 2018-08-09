package util

// ArrayReverseInt64 int64 array reverse
func ArrayReverseInt64(a []int64) []int64 {
	l := len(a)
	b := make([]int64, l)
	for i := l; i > 0; i-- {
		b[l-i] = a[i-1]
	}
	return b
}
