package util

import "testing"

func TestGetUUIDV1(t *testing.T) {
	t.Log(GetUUIDV1())
}

func TestGetUUIDV1WithoutLine(t *testing.T) {
	t.Log(GetUUIDV1WithoutLine())
}
