package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// BytesToUint32 byte to uint32
func BytesToUint32(src []byte) (dst uint32, err error) {
	if len(src) != 4 {
		err = fmt.Errorf("len(src) != 4")
		return
	}
	buff := bytes.NewBuffer(src)
	if err = binary.Read(buff, binary.BigEndian, &dst); err != nil {
		return
	}
	return
}

// BytesCombine combine bytes
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// BytesToHex bytes to hexadecimal
func BytesToHex(src []byte) string {
	return fmt.Sprintf("%x", src)
}
