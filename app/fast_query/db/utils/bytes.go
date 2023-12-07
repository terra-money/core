package utils

import (
	"bytes"
	"encoding/binary"
)

func ConcatBytes(items ...[]byte) []byte {
	buf := new(bytes.Buffer)
	for _, item := range items {
		buf.Write(item)
	}
	return buf.Bytes()
}

func UintToBigEndian(n uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)

	return buf
}

func BigEndianToUint(n []byte) uint64 {
	return binary.BigEndian.Uint64(n)
}
