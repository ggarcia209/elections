package main

import (
	"fmt"
	"os"

	"github.com/elections/testing/unit_tests/databuilder/testDB"
)

func main() {
	// Part 1 - compare units
	/* err := testDB.TestCompareUnits()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} */

	// part 2 - compare
	err := testDB.TestCompare()
	if err != nil {
		fmt.Println("main failed: ")
		os.Exit(1)
	}
}
