package util

import (
	"bytes"
	"testing"
)

func TestGbkToUtf8(t *testing.T) {
	tt := []byte("你好，世界")
	b, _ := Utf8ToGbk(tt)
	b, _ = GbkToUtf8(b)
	if !bytes.Equal(b, tt) {
		t.Error("TestGbkToUtf8 error", string(b), string(tt))
	}
}