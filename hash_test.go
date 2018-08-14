package util

import "testing"

func TestMd5(t *testing.T) {
	v, err := Md5([]byte("hao"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(v))
}
