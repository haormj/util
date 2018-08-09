package aes

import (
	"bytes"
	"testing"
)

func TestCBC(t *testing.T) {
	key := []byte("0123456789abcdef")
	plainText := []byte("AES CBC PKCS7")
	cipherText, err := CBCEncrypt(plainText, key)
	if err != nil {
		t.Fatal(err)
	}
	pt, err := CBCDecrypt(cipherText, key)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(pt, plainText) {
		t.Errorf("expected:%v,got:%v", plainText, pt)
	}
}
