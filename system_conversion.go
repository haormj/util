package util

// 2 binary
// 8 octal
// 10 decimal
// 16 hexadecimal

// DecimalToAny decimal to any
// ignore positive and negative
func DecimalToAny(i int64, base int) []int64 {
	if i == 0 {
		return []int64{0}
	}
	if i < 0 {
		i = -i
	}
	var dividend, divisor, quotient, remainder int64
	dividend = i
	divisor = int64(base)
	a := make([]int64, 0)
	for dividend > 0 {
		quotient = dividend / divisor
		remainder = dividend % divisor
		dividend = quotient
		a = append(a, remainder)
	}
	// reverse a
	ArrayReverseInt64(a)
	return a
}

// AnyToDecimal any to decimal
func AnyToDecimal(b []int64, base int) int64 {
	a := make([]int64, len(b))
	copy(a, b)
	ArrayReverseInt64(a)
	var i, power int64
	for j, v := range a {
		if j == 0 {
			power = 1
		} else {
			power *= int64(base)
		}
		i += v * power
	}
	return i
}
