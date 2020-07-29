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
	testDB.TestDeriveTxTypes()

	/* TopOverall units */
	// part 6 - test compareTopOverall
	// 7/22/20 - SUCCESS
	/* err := testDB.TestCompareTopOverall()
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// part 7 - test updateAndSave
	err := testDB.TestUpdateAndSave()
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// part 8  - test party switch cases (SUCCESS)
	err := testDB.TestPartySwitchCases()
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// part 9 - test type assertion switch cases and nested if/esle statements (SUCCESS)
	err := testDB.TestTopOverallInternalLogic()
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// part 10 - test updateTopOverall (SUCCESS)
	err := testDB.TestUpdateTopOverall()
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	} */

	// part 11 - test TransactionUpdate
	err := testDB.TestTransactionUpdate()
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
}
