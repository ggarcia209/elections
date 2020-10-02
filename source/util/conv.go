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

// PutBigEndianStr is used to encode string values to big endian bytes
func PutBigEndianStr(str string) []byte {
	buf := make([]byte, 8)
	v := binary.BigEndian.Uint64([]byte(str))
	binary.BigEndian.PutUint64(buf, v)
	return buf
}

// PutBigEndianUint64 is used to encode uint64 values to big endian bytes
func PutBigEndianUint64(ui uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, ui)
	return buf
}

// GetBigEndian is used to convert bytes encoded in Big Endian format to uint64
func GetBigEndian(buf []byte) uint64 {
	x := binary.BigEndian.Uint64(buf)
	return x
}
