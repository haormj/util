package util

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GbkToUtf8 convert gbk to utf-8
func GbkToUtf8(in []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(in), simplifiedchinese.GBK.NewDecoder())
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Utf8ToGbk convert utf-8 to gbk
func Utf8ToGbk(in []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(in), simplifiedchinese.GB18030.NewEncoder())
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return b, nil
}