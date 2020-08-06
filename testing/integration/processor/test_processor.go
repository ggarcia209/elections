package main

/* TEST NOTES */
// Refactor dontations.CmteTxData object - add maps & split Org / Cmte contributors
// Refactored persist.GetObject & persist.encodeToProto to return pointer to donations obj type, not pointer to interface
// PutObject fails in transactions stage if Other cmte in transaction is not saved to database prior
//   - implement createCmte function? - verify if missing committees exist in complete bulk data records
// Individual derived from individual contributions file returned as Organziation objects
//  -

import (
	"fmt"
	"os"

	"github.com/elections/databuilder"
	"github.com/elections/donations"
	"github.com/elections/parse"
	"github.com/elections/persist"
)

func main() {
	year := "2018" // derived from flag inputs

	// runProcessorTest(year)

	err := viewTopOverall(year)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	fmt.Println("MAIN DONE")
}

func runProcessorTest(year string) {
	// input file paths
	candPath := "../../test_data/test_cands.txt"
	cmtePath := "../../test_data/test_cmtes.txt"
	icPath := "../../test_data/test_ics.txt"
	ccPath := "../../test_data/test_ccs.txt"
	disbPath := "../../test_data/test_disbs.txt"

	// initialize database and TopOverallData objects
	persist.Init(year)
	topOverallList := donations.InitTopOverallDataObjs(100)
	err := persist.StoreObjects(year, topOverallList)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// process candidates and committee objects first
	err = candTest(year, candPath)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	err = cmteTest(year, cmtePath)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// process transactions
	err = cmteContTest(year, ccPath)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	err = indvContTest(year, icPath)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	err = disbTest(year, disbPath)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	fmt.Println("PROCESSOR TEST DONE")
}

func indvContTest(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("indvContTest faield: ", err)
		os.Exit(1)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset("indv")
	if err != nil {
		fmt.Println("indvContTest faield: ", err)
		return err
	}

	// parse file
	for {
		// parse records
		txQueue, offset, err := parse.ScanContributions(year, file, start)
		if err != nil {
			fmt.Println("indvContTest faield: ", err)
			return err
		}

		// create cache from record IDs
		cache, err := databuilder.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println("indvContTest faield: ", err)
			return err
		}

		// break if finsished
		if len(cache) == 0 {
			fmt.Println("indvContTest DONE - no cache")
			return nil
		}

		// update object data for each transaction
		err = databuilder.TransactionUpdate(year, txQueue, cache)
		if err != nil {
			fmt.Println("indvContTest faield: ", err)
			return err
		}

		// persist objects in cache
		ser := databuilder.SerializeCache(cache)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save offset value after objects persisted
		err = persist.LogOffset("indv", offset)
		if err != nil {
			fmt.Println(err)
			return err
		}
		start = offset

		// break if at end of file
		if len(txQueue) < 1000 {
			break
		}
	}

	fmt.Println("Individual Contributions -  DONE")

	return nil
}

func disbTest(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset("disb")
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		// parse records
		txQueue, offset, err := parse.ScanDisbursements(year, file, start)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// create cache from
		cache, err := databuilder.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if len(cache) == 0 {
			fmt.Println("disbTest DONE - no cache")
			return nil
		}

		// update Indvidual Donors and receiving committees
		err = databuilder.TransactionUpdate(year, txQueue, cache)
		if err != nil {
			fmt.Println(err)
			return err
		}

		ser := databuilder.SerializeCache(cache)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save offset value after objects persisted
		err = persist.LogOffset("disb", offset)
		if err != nil {
			fmt.Println(err)
			return err
		}
		start = offset

		// break if at end of file
		if len(txQueue) < 1000 {
			break
		}
	}

	fmt.Println("Disbursements -  DONE")

	return nil
}

func cmteContTest(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset("cmte_cont")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// parse file
	for {
		// parse 25 records per iteration
		txQueue, offset, err := parse.ScanContributions(year, file, start)
		if err != nil {
			fmt.Println("cmteContTest failed: ", err)
			return err
		}

		cache, err := databuilder.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if len(cache) == 0 {
			fmt.Println("cmteContTest DONE - no cache")
			return nil
		}

		// update Indvidual Donors and receiving committees
		err = databuilder.TransactionUpdate(year, txQueue, cache)
		if err != nil {
			fmt.Println(err)
			return err
		}

		ser := databuilder.SerializeCache(cache)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save offset value after objects persisted
		err = persist.LogOffset("cmte_cont", offset)
		if err != nil {
			fmt.Println(err)
			return err
		}
		start = offset

		// break if at end of file
		if len(txQueue) < 1000 {
			break
		}
	}

	fmt.Println("Committee Contributions -  DONE")

	return nil
}

func candTest(year, filePath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset("cand")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, offset, err := parse.ScanCandidates(file, start)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save offset value after objects persisted
		err = persist.LogOffset("cand", offset)
		if err != nil {
			fmt.Println(err)
			return err
		}
		start = offset

		// break if at end of file
		if len(objQueue) < 100 {
			break
		}
	}

	fmt.Println("Candidates - DONE")

	return nil
}

func cmteTest(year, filePath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset("cmte")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, txDataQueue, offset, err := parse.ScanCommittees(file, start)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = persist.StoreObjects(year, txDataQueue)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save offset value after objects persisted
		err = persist.LogOffset("cmte", offset)
		if err != nil {
			fmt.Println(err)
			return err
		}
		start = offset

		// break if at end of file
		if len(objQueue) < 100 {
			break
		}
	}

	fmt.Println("Committees - DONE")
	return nil
}

func viewTopOverall(year string) error {
	testCmte, err := persist.GetObject(year, "cmte_tx_data", "C00343871") // C00401224
	if err != nil {
		fmt.Println(err)
		return err
	}

	odIndv, err := persist.GetObject(year, "top_overall", "indv")
	if err != nil {
		fmt.Println(err)
		return err
	}

	odCmte, err := persist.GetObject(year, "top_overall", "cmte_recs_all")
	if err != nil {
		fmt.Println(err)
		return err
	}

	testIndv, err := persist.GetObject(year, "individuals", "01489ace1a99994034ef8df6455bb6de")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("*** TEST COMMITTEES ***")
	fmt.Println(testCmte)
	fmt.Println()
	fmt.Println("*** TEST INDIVIDUAL ***")
	fmt.Println(testIndv)
	fmt.Println()

	fmt.Println("*** TOP COMMITTEES ***")
	fmt.Println(odCmte)
	fmt.Println()
	fmt.Println("*** TOP INDIVIDUALS ***")
	fmt.Println(odIndv)
	fmt.Println()

	return nil
}
