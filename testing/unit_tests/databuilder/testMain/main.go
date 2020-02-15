package main

import (
	"github.com/elections/testing/unit_tests/databuilder/testDB"
)

func main() {
	// Part 1 - compare units (SUCCESS)
	/* err := testDB.TestCompareUnits()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// part 2 - compare (SUCCESS)
	err := testDB.TestCompare()
	if err != nil {
		fmt.Println("main failed: ")
		os.Exit(1)
	}

	// part 3 - checkThreshold (SUCCESS)
	err := testDB.TestCheckThreshold()
	if err != nil {
		fmt.Println("main failed: ")
		os.Exit(1)
	}

	// part 4 - cmteCompGen (SUCCESS)
	testDB.TestCmteCompGen() */

	// part 5 - test updateTop internal logic (SUCCESS)
	testDB.TestUpdateTopInternalLogic()

}
