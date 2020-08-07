package admin

import (
	"fmt"
	"os"

	"github.com/elections/databuilder"
	"github.com/elections/donations"
	"github.com/elections/parse"
	"github.com/elections/persist"
)

// ProcessNewRecords processes the FEC bulk data files for the given year.
// All directories must have the following files:
//   /[year]/cands.txt
//   /[year]/cmtex.txt
//   /[year]/cmte_fin.txt
//   /[year]/ccs.txt
//   /[year]/ics.txt
//   /[year]/disbs.txt
func ProcessNewRecords(year string) error {
	fmt.Println("******************************************")
	fmt.Println("PROCESS NEW RECORDS BEGIN - YEAR: ", year)
	fmt.Println()

	// get year from Command Line input
	// input file paths - placeholders
	candPath := "../../" + year + "/cands.txt"
	cmtePath := "../../" + year + "/cmtes.txt"
	icPath := "../../" + year + "/ics.txt"
	ccPath := "../../" + year + "/ccs.txt"
	disbPath := "../../" + year + "/disbs.txt"
	cmteFinPath := "../../" + year + "/cmte_fin.txt"

	// initialize database and TopOverallData objects
	//   REFACTOR FOR IDEMPOTENCY
	persist.Init(year)

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cand")
	if err != nil {
		fmt.Println("processCandidates failed: ", err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}

	if start == 0 {
		// initialize and persist TopOverallData objects
		topOverallList := donations.InitTopOverallDataObjs(100)
		err := persist.StoreObjects(year, topOverallList)
		if err != nil {
			fmt.Println("ProcessNewRecords failed: ", err)
			return fmt.Errorf("ProcessNewRecords failed: %v", err)
		}
	}

	// process candidates and committee objects first
	err = processCandidates(year, candPath)
	if err != nil {
		fmt.Println("ProcessNewRecords failed: ", err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	err = processCommittees(year, cmtePath)
	if err != nil {
		fmt.Println("ProcessNewRecords failed: ", err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	err = processCmteFinancials(year, cmteFinPath)
	if err != nil {
		fmt.Println("ProcessNewRecords failed: ", err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	// process transactions
	err = processCmteContributions(year, ccPath)
	if err != nil {
		fmt.Println("ProcessNewRecords failed: ", err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	err = processIndvContributions(year, icPath)
	if err != nil {
		fmt.Println("ProcessNewRecords failed: ", err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	err = processDisbursements(year, disbPath)
	if err != nil {
		fmt.Println("ProcessNewRecords failed: ", err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	fmt.Println("PROCESS NEW RECORDS COMPLETE - YEAR: ", year)
	fmt.Println("******************************************")
	fmt.Println()

	return nil
}

func processCandidates(year, filePath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("processCandidates failed: ", err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cand")
	if err != nil {
		fmt.Println("processCandidates failed: ", err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, offset, err := parse.ScanCandidates(file, start)
		if err != nil {
			fmt.Println("processCandidates failed: ", err)
			return fmt.Errorf("processCandidates failed: %v", err)
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println("processCandidates failed: ", err)
			return fmt.Errorf("processCandidates failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cand", offset)
		if err != nil {
			fmt.Println("processCandidates failed: ", err)
			return fmt.Errorf("processCandidates failed: %v", err)
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

func processCommittees(year, filePath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("processCommittees failed: ", err)
		return fmt.Errorf("processCommittees failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmte")
	if err != nil {
		fmt.Println("processCommittees failed: ", err)
		return fmt.Errorf("processCommittees failed: %v", err)
	}

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, txDataQueue, offset, err := parse.ScanCommittees(file, start)
		if err != nil {
			fmt.Println("processCommittees failed: ", err)
			return fmt.Errorf("processCommittees failed: %v", err)
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println("processCommittees failed: ", err)
			return fmt.Errorf("processCommittees failed: %v", err)
		}
		err = persist.StoreObjects(year, txDataQueue)
		if err != nil {
			fmt.Println("processCommittees failed: ", err)
			return fmt.Errorf("processCommittees failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmte", offset)
		if err != nil {
			fmt.Println("processCommittees failed: ", err)
			return fmt.Errorf("processCommittees failed: %v", err)
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

func processCmteFinancials(year, filePath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("processCmteFinancials failed: ", err)
		return fmt.Errorf("processCmteFinancials failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmte_fin")
	if err != nil {
		fmt.Println("processCmteFinancials failed: ", err)
		return fmt.Errorf("processCmteFinancials failed: %v", err)
	}

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, offset, err := parse.ScanCmteFin(file, start)
		if err != nil {
			fmt.Println("processCmteFinancials failed: ", err)
			return fmt.Errorf("processCmteFinancials failed: %v", err)
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println("processCmteFinancials failed: ", err)
			return fmt.Errorf("processCmteFinancials failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmte_fin", offset)
		if err != nil {
			fmt.Println("processCmteFinancials failed: ", err)
			return fmt.Errorf("processCmteFinancials failed: %v", err)
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

func processCmteContributions(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("processCmteContributions failed: ", err)
		return fmt.Errorf("processCmteContributions failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmte_cont")
	if err != nil {
		fmt.Println("processCmteContributions failed: ", err)
		return fmt.Errorf("processCmteContributions failed: %v", err)
	}

	// parse file
	for {
		// parse records
		txQueue, offset, err := parse.ScanContributions(year, file, start)
		if err != nil {
			fmt.Println("processCmteContributions failed: ", err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}

		// create cache from record IDs
		cache, err := databuilder.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println("processCmteContributions failed: ", err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}

		// break if finished
		if len(cache) == 0 {
			fmt.Println("cmteContTest DONE - no cache")
			return nil
		}

		// update cached objects for each transaction
		err = databuilder.TransactionUpdate(year, txQueue, cache)
		if err != nil {
			fmt.Println("processCmteContributions failed: ", err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}

		// persist items in cache
		ser := databuilder.SerializeCache(cache)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println("processCmteContributions failed: ", err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmte_cont", offset)
		if err != nil {
			fmt.Println("processCmteContributions failed: ", err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
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

func processIndvContributions(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("processIndvContributions faield: ", err)
		os.Exit(1)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "indv")
	if err != nil {
		fmt.Println("processIndvContributions faield: ", err)
		return err
	}

	// parse file
	for {
		// parse records
		txQueue, offset, err := parse.ScanContributions(year, file, start)
		if err != nil {
			fmt.Println("processIndvContributions faield: ", err)
			return err
		}

		// create cache from record IDs
		cache, err := databuilder.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println("processIndvContributions faield: ", err)
			return err
		}

		// break if finsished
		if len(cache) == 0 {
			fmt.Println("processIndvContributions DONE - no cache")
			return nil
		}

		// update object data for each transaction
		err = databuilder.TransactionUpdate(year, txQueue, cache)
		if err != nil {
			fmt.Println("processIndvContributions faield: ", err)
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
		err = persist.LogOffset(year, "indv", offset)
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

func processDisbursements(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("processDisbursements failed: ", err)
		return fmt.Errorf("processDisbursements failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "disb")
	if err != nil {
		fmt.Println("processDisbursements failed: ", err)
		return fmt.Errorf("processDisbursements failed: %v", err)
	}

	for {
		// parse records
		txQueue, offset, err := parse.ScanDisbursements(year, file, start)
		if err != nil {
			fmt.Println("processDisbursements failed: ", err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}

		// create cache from
		cache, err := databuilder.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println("processDisbursements failed: ", err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}

		// break if finished
		if len(cache) == 0 {
			fmt.Println("disbTest DONE - no cache")
			return nil
		}

		// update object data
		err = databuilder.TransactionUpdate(year, txQueue, cache)
		if err != nil {
			fmt.Println("processDisbursements failed: ", err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}

		// persist objects in cache
		ser := databuilder.SerializeCache(cache)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println("processDisbursements failed: ", err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "disb", offset)
		if err != nil {
			fmt.Println("processDisbursements failed: ", err)
			return fmt.Errorf("processDisbursements failed: %v", err)
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
