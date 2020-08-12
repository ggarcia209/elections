package main

import (
	"fmt"
	"os"

	"github.com/elections/testing/unit_tests/indexing"
)

func main() {
	// indexing.TestComponents()
	err := indexing.TestGetResults()
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	fmt.Println("main done")
}
