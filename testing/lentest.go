package main

import "fmt"

func main() {
	fmt.Println(len("ID00123456"))
	fmt.Println(len([]byte("ID00123456")))

	m := make(map[string]int16)
	m["ID00123456"] = 2000
	m["ID00534111"] = 100

	fmt.Println(len(m))
}
