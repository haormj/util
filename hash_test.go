package util

import (
	"testing"
)

func TestMd5(t *testing.T) {
	v, err := Md5([]byte("hao"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(v))
}

func TestHmacSha256(t *testing.T) {
	b, err := HmacSha256([]byte("hello world"), []byte("password"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(b))
}

func TestSha256(t *testing.T) {
	b, err := Sha256([]byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(b))
}
