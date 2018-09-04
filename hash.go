package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
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
	return hex.EncodeToString(b), nil
}

// Sha256 golang sha256
func Sha256(p []byte) ([]byte, error) {
	hash := sha256.New()
	if _, err := hash.Write(p); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// HmacSha1 call golang hmac
func HmacSha1(p, key []byte) ([]byte, error) {
	hash := hmac.New(sha1.New, key)
	if _, err := hash.Write(p); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// HmacSha256 hmac sha256
func HmacSha256(p, key []byte) ([]byte, error) {
	hash := hmac.New(sha256.New, key)
	if _, err := hash.Write(p); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}
