// SUCCESS
package main

import (
	"fmt"
	"os"

	"github.com/elections/parse"
)

func main() {
	key1, key2 := "abc", "xyz"
	offset1, offset2 := 50, 100

	parse.CreateDB()
	if err := parse.LogOffset(key1, offset1); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := parse.LogOffset(key2, offset2); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	val1, err := parse.GetOffset(key1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	val2, err := parse.GetOffset(key2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	val3 := val1 + val2

	fmt.Println("val1: ", val1)
	fmt.Println("val2: ", val2)
	fmt.Println("val3: ", val3)

}
