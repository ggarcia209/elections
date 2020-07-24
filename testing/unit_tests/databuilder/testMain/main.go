package main

import (
	"fmt"
	"os"

	"github.com/elections/testing/unit_tests/databuilder/testDB"
)

func main() {
	// Part 1 - compare units (SUCCESS)
	// 7/22/20 - successful test after refactor
	/* err := testDB.TestCompareUnits()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// part 2 - compare (SUCCESS)
	// 7/22/20 - successful test after refactor
	err := testDB.TestCompare()
	if err != nil {
		fmt.Println("main failed: ")
		os.Exit(1)
	}

	// part 3 - checkThreshold (SUCCESS)
	/* err := testDB.TestCheckThreshold()
	if err != nil {
		fmt.Println("main failed: ")
		os.Exit(1)
	}

	// part 4 - cmteCompGen (SUCCESS)
	testDB.TestCmteCompGen()

	// part 5 - test updateTop internal logic (SUCCESS)
	// testDB.TestUpdateTopInternalLogic()
	testDB.TestTxUpdateInternalLogic()
	testDB.TestDeriveTxTypes() */

	// part 6 - test updateTopOverall
	// 7/22/20 - SUCCESS
	err := testDB.TestCompareTopOverall()
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

}
