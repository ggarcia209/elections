package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

func NewHash(input string) string {
	sum := md5.Sum([]byte(input))
	// pass := base64.StdEncoding.EncodeToString(sum[:])
	pass := hex.EncodeToString(sum[:])

	return pass
}

func main() {
	in1 := "ANSARY, HUSHANG HON.|HOUSTON|TX|770025014"
	in2 := "KHOURI, LAURA|IRVINE|CA|92614564743636336"
	in3 := "MCPHEE, JEFF|OAKDALE|CA|95361"
	in4 := "ROBERTS, KELLY|NEWPORT BEACH|CA|92660"

	start := time.Now()
	hash1 := NewHash(in1)
	elapsed := time.Since(start)
	fmt.Println("time elapsed: ", elapsed)

	start = time.Now()
	hash2 := NewHash(in2)
	elapsed = time.Since(start)
	fmt.Println("time elapsed: ", elapsed)

	start = time.Now()
	hash3 := NewHash(in3)
	elapsed = time.Since(start)
	fmt.Println("time elapsed: ", elapsed)

	start = time.Now()
	hash4 := NewHash(in4)
	elapsed = time.Since(start)
	fmt.Println("time elapsed: ", elapsed)

	fmt.Println(string(hash1))
	fmt.Println()
	fmt.Println(string(hash2))
	fmt.Println()
	fmt.Println(string(hash3))
	fmt.Println(string(hash4))
}
