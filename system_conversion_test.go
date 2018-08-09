package util

import (
	"testing"
)

func TestDecimalToAny(t *testing.T) {
	t.Log(DecimalToAny(1709, 90))
}

func TestAnyToDecimal(t *testing.T) {
	t.Log(AnyToDecimal([]int64{18, 89}, 90))
}
