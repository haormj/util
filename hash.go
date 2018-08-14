package util

import (
	"crypto/md5"
)

// Md5 call golang md5
func Md5(p []byte) ([]byte, error) {
	hash := md5.New()
	if _, err := hash.Write(p); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// Md5ToString result to hex
func Md5ToString(p []byte) (string, error) {
	b, err := Md5(p)
	if err != nil {
		return "", err
	}
	return BytesToHex(b), nil
}
