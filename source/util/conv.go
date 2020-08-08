package util

import (
	"encoding/binary"
)

// Itob encodes single int to byte slice
func Itob(i int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}

// Btoi decodes byte slice representing single uint16 to single int
func Btoi(b []byte) int {
	if b == nil {
		return 0
	}
	return int(binary.BigEndian.Uint32(b))
}
