package qwikscript

import (
	"encoding/binary"
	"fmt"
	"sort"
	"unsafe"
)

type ByteOrder interface {
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	String() string
}

func TestBigEndian() {
	ID := "00abt72ndh27dj4fg592nhd2kfie67h3"
	ID = "!!!!"
	idBytes := []byte(ID)

	fmt.Println("Uint32...")
	buf := make([]byte, 4)
	fmt.Println("buf: ", buf)

	y := binary.BigEndian.Uint32(idBytes)
	binary.BigEndian.PutUint32(buf, y)
	fmt.Println("Big Endian buf: ", buf)

	/* z := binary.BigEndian.Uint32(buf)
	fmt.Println("size z: ", unsafe.Sizeof(z))
	str := strconv.Itoa(int(z))
	fmt.Println("str: ", str)

	buf := make([]byte, 8)
	fmt.Println("buf: ", buf)

	y := binary.BigEndian.Uint64(idBytes)
	binary.BigEndian.PutUint64(buf, y)
	fmt.Println("Big Endian buf: ", buf)

	z := binary.BigEndian.Uint32(buf)
	fmt.Println("size z: ", unsafe.Sizeof(z))
	str := strconv.Itoa(int(z))
	fmt.Println("str: ", str) */

	/* fmt.Println("Uint64...")

	buf = make([]byte, 8)
	fmt.Println("buf: ", buf)

	v := binary.BigEndian.Uint64(idBytes)
	binary.BigEndian.PutUint64(buf, v)
	fmt.Println("Big Endian buf: ", buf)

	x := binary.BigEndian.Uint64(buf)
	str = strconv.Itoa(int(x))
	fmt.Println("str: ", str) */
}

func TestBigEndianSplit() {
	ID := "00abt72ndh27dj4fg592nhd2kfie67h3"
	idBytes := []byte(ID)
	fmt.Println(idBytes)
	fmt.Println("idBytes: ", len(idBytes))

	fmt.Println("Uint64...")
	buf := make([]byte, 8)
	fmt.Println("buf: ", buf)

	y := binary.BigEndian.Uint64(idBytes[:8])
	binary.BigEndian.PutUint64(buf, y)
	fmt.Println("Big Endian buf: ", buf)
	fmt.Println("  size: ", unsafe.Sizeof(buf))
	x := binary.BigEndian.Uint64(buf)
	fmt.Println("big endian: ", x)
	fmt.Println("  size: ", unsafe.Sizeof(x))

	buf = make([]byte, 8)
	y = binary.BigEndian.Uint64(idBytes[8:16])
	binary.BigEndian.PutUint64(buf, y)
	fmt.Println("Big Endian buf: ", buf)
	fmt.Println("  size: ", unsafe.Sizeof(buf))
	x = binary.BigEndian.Uint64(buf)
	fmt.Println("big endian: ", x)
	fmt.Println("  size: ", unsafe.Sizeof(x))

	buf = make([]byte, 8)
	y = binary.BigEndian.Uint64(idBytes[16:24])
	binary.BigEndian.PutUint64(buf, y)
	fmt.Println("Big Endian buf: ", buf)
	fmt.Println("  size: ", unsafe.Sizeof(buf))
	x = binary.BigEndian.Uint64(buf)
	fmt.Println("big endian: ", x)
	fmt.Println("  size: ", unsafe.Sizeof(x))

	buf = make([]byte, 8)
	y = binary.BigEndian.Uint64(idBytes[24:])
	binary.BigEndian.PutUint64(buf, y)
	fmt.Println("Big Endian buf: ", buf)
	fmt.Println("  size: ", unsafe.Sizeof(buf))
	x = binary.BigEndian.Uint64(buf)
	fmt.Println("big endian: ", x)
	fmt.Println("  size: ", unsafe.Sizeof(x))
}

func EndianComp() {
	ID0 := "00abt72ndh27dj4fg592nhd2kfie67h3"
	ID1 := "C5abt7he5147dj4g592hd103fie6453"
	ID2 := "30abt7he5147rt4fg592hd103fie6453"
	ID3 := "C0000456"
	ID4 := "01hyat7h451b7rt4fg5d2hda43fia5j53"
	ids := []int{}

	idBytes := []byte(ID0)
	buf := make([]byte, 8)
	v := binary.BigEndian.Uint64(idBytes)
	binary.BigEndian.PutUint64(buf, v)
	a := binary.BigEndian.Uint64(buf)
	fmt.Println("0: ", ID0)
	fmt.Println("a: ", a)
	ids = append(ids, int(a))

	idBytes = []byte(ID1)
	buf = make([]byte, 8)
	v = binary.BigEndian.Uint64(idBytes)
	binary.BigEndian.PutUint64(buf, v)
	b := binary.BigEndian.Uint64(buf)
	fmt.Println("1: ", ID1)
	fmt.Println("b: ", b)
	ids = append(ids, int(b))

	idBytes = []byte(ID2)
	buf = make([]byte, 8)
	v = binary.BigEndian.Uint64(idBytes)
	binary.BigEndian.PutUint64(buf, v)
	c := binary.BigEndian.Uint64(buf)
	fmt.Println("2: ", ID2)
	fmt.Println("c: ", c)
	ids = append(ids, int(c))

	idBytes = []byte(ID3)
	buf = make([]byte, 8)
	v = binary.BigEndian.Uint64(idBytes)
	binary.BigEndian.PutUint64(buf, v)
	d := binary.BigEndian.Uint64(buf)
	fmt.Println("3: ", ID3)
	fmt.Println("d: ", d)
	ids = append(ids, int(d))

	idBytes = []byte(ID4)
	buf = make([]byte, 8)
	v = binary.BigEndian.Uint64(idBytes)
	binary.BigEndian.PutUint64(buf, v)
	e := binary.BigEndian.Uint64(buf)
	fmt.Println("3: ", ID4)
	fmt.Println("e: ", e)
	ids = append(ids, int(e))

	sort.Ints(ids)
	fmt.Println("ids ", ids)
}

func ShardSizeMax() {
	ids := []string{}
	for i := 0; i < 100; i++ {
		ID := "00abt72ndh27dj4fg592nhd2kfie67h3"
		ids = append(ids, ID)
	}
	fmt.Println("len IDs: ", len(ids))
	fmt.Println("IDs shard size - no encoding: ", unsafe.Sizeof(ids))

	var ends []uint32
	for i := 0; i < 100; i++ {
		ID := "00abt72ndh27dj4fg592nhd2kfie67h3"
		buf := make([]byte, 4)
		v := binary.BigEndian.Uint32([]byte(ID))
		binary.BigEndian.PutUint32(buf, v)
		x := binary.BigEndian.Uint32(buf)
		ends = append(ends, x)
	}
	fmt.Println("len IDs: ", len(ends))
	fmt.Println("IDs shard size - bigEndian: ", unsafe.Sizeof(ends))
}
