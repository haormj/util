package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ByteToUint32 byte to uint32
func ByteToUint32(src []byte) (dst uint32, err error) {
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
