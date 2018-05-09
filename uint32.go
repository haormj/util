package util

import (
	"bytes"
	"encoding/binary"
)

// Uint32ToByte uint32 to byte
func Uint32ToByte(src uint32) (dst []byte, err error) {
	buff := bytes.NewBuffer(nil)
	if err = binary.Write(buff, binary.BigEndian, src); err != nil {
		return
	}
	dst = buff.Bytes()
	return
}
