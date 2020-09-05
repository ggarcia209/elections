package util

import (
	"encoding/binary"
)

// Itob encodes single int to byte slice
func Itob(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

// Btoi decodes byte slice representing single uint16 to single int
func Btoi(b []byte) int64 {
	if b == nil {
		return 0
	}
	return int64(binary.BigEndian.Uint64(b))
}
