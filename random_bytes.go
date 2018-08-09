package util

import (
	"crypto/rand"
	"io"
)

// RandomBytesE random by rand.Reader
// it will use cpu more
func RandomBytesE(n int) ([]byte, error) {
	buff := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, buff)
	if err != nil {
		return buff, err
	}
	return buff, nil
}

func RandomBytes(n int) []byte {
	buff, _ := RandomBytesE(n)
	return buff
}
