package util

import (
	"testing"
)

func TestKB(t *testing.T) {
	if KB != 1024 {
		t.Fail()
	}
}
