package admin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/elections/source/cache"
	"github.com/elections/source/databuilder"
	"github.com/elections/source/donations"
	"github.com/elections/source/parse"
	"github.com/elections/source/persist"
	"github.com/elections/source/ui"
)

/* DB CREATE OPERATIONS */

// ProcessNewRecords processes the FEC bulk data files for the given year.
// All directories must have the following files:
//   input/[year]/cmte/cn.txt
//   input/[year]/cand/cm.txt
//   input/[year]/pac/webk.txt
//   input/[year]/ctx/itoth.txt
//   input/[year]/indiv/itcont.txt
//   input/[year]/exp/oppexp.txt
func ProcessNewRecords() error {
	fmt.Println("******************************************")
	fmt.Println()

	// get input folder path
	input, err := persist.GetPath(true)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}
	if input == "" {
		input, err = getPath(true)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("ProcessNewRecords failed: %v", err)
		}
	}

	// derive year from filepath
	year := ui.GetYear()
	fmt.Println("Chosen year: ", year)

	// get output path
	output, err := persist.GetPath(false)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}
	if output == "" {
		output, err = getPath(false)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("ProcessNewRecords failed: %v", err)
		}
	}

	persist.OUTPUT_PATH = output

	fmt.Println("filepaths set - continue with data processing?")
	yes := ui.Ask4confirm()
	if !yes {
		fmt.Println("Returning to menu...")
		return nil
	}

	fmt.Println("PROCESS NEW RECORDS BEGIN - YEAR: ", year)

	// get year from Command Line input
	// input file paths - placeholders
	root := filepath.Join(input, year)
	candPath := filepath.Join(root, "cand", "cn.txt")
	cmtePath := filepath.Join(root, "cmte", "cm.txt")
	icPath := filepath.Join(root, "indiv", "itcont.txt")
	ccPath := filepath.Join(root, "ctx", "itoth.txt")
	disbPath := filepath.Join(root, "exp", "oppexp.txt")
	cmteFinPath := filepath.Join(root, "pac", "webk.txt")

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
		c, err := cache.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println("processCmteContributions failed: ", err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}

		// break if finished
		if len(c) == 0 {
			fmt.Println("cmteContTest DONE - no cache")
			return nil
		}

		// update cached objects for each transaction
		err = databuilder.TransactionUpdate(year, txQueue, c)
		if err != nil {
			fmt.Println("processCmteContributions failed: ", err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}

		// persist items in cache
		ser := cache.SerializeCache(c)
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
		c, err := cache.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println("processIndvContributions faield: ", err)
			return err
		}

		// break if finsished
		if len(c) == 0 {
			fmt.Println("processIndvContributions DONE - no cache")
			return nil
		}

		// update object data for each transaction
		err = databuilder.TransactionUpdate(year, txQueue, c)
		if err != nil {
			fmt.Println("processIndvContributions faield: ", err)
			return err
		}

		// persist objects in cache
		ser := cache.SerializeCache(c)
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
		c, err := cache.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println("processDisbursements failed: ", err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}

		// break if finished
		if len(c) == 0 {
			fmt.Println("disbTest DONE - no cache")
			return nil
		}

		// update object data
		err = databuilder.TransactionUpdate(year, txQueue, c)
		if err != nil {
			fmt.Println("processDisbursements failed: ", err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}

		// persist objects in cache
		ser := cache.SerializeCache(c)
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

// pathExists checks to see if given file path is valid
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func getPath(input bool) (string, error) {
	var name string
	var msg string
	if input {
		name = "input"
		msg = "Enter input folder filepath (schema: /root/input): "
	} else {
		name = "output"
		msg = "Enter database & search index output folder (schema: /root/output): "
	}
	// get folder path
	path, err := persist.GetPath(input)
	fmt.Println("Found path: ", path)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("getPath failed: %v", err)
	}
	if path != "" { // if path set previously
		fmt.Printf("Set new %s path? Current path will be overwritten:\n\t%s\n", name, path)
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Printf("Continuing with current %s path: %s\n", name, path)
		} else { // confirm overwrite
			fmt.Println(">>> Are you sure you want to overwrite the existing path?")
			yes := ui.Ask4confirm()
			if !yes {
				fmt.Printf("Continuing with current %s path: %s\n", name, path)
			} else { // set new path
				fmt.Println(msg)
				path = ui.GetPathFromUser()
				err := persist.LogPath(path, input)
				fmt.Printf("New path: %s saved\n", path)
				if err != nil {
					fmt.Println(err)
					return "", fmt.Errorf("getPath failed: %v", err)
				}
			}
		}
		return path, nil
	}

	// path is not set previously
	fmt.Println(msg)
	path = ui.GetPathFromUser()
	err = persist.LogPath(path, input)
	fmt.Printf("New path: %s saved\n", path)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("getPath failed: %v", err)
	}
	return path, nil
}

/* func getYear() (string, error) {
	yearStr := "" // default return value

	for {
		// get flag and check validity of file paths in loop
		yearFlag := flag.Int("year", 0, "'year' flag defines which election cycle's dataset to process")
		flag.Parse()
		year := *yearFlag
		if year == 0 {
			fmt.Println("'year' flag must be set to valid year")
			continue
		}

		// convert int value to string
		yearStr = strconv.Itoa(year)

		indvCont := "../raw_data/cmte_cont/" + yearStr + "/itcont.txt"
		cmteCont := "../raw_data/cmte_cont/" + yearStr + "/itoth.txt"
		disb := "../raw_data/cmte_cont/" + yearStr + "/oppenxp.txt"
		cand := "../raw_data/cmte_cont/" + yearStr + "/cn.txt"
		cmte := "../raw_data/cmte_cont/" + yearStr + "/cm.txt"
		paths := []string{indvCont, cmteCont, disb, cand, cmte}

		// check validity of each path
		valid := true
		for _, path := range paths {
			exists, err := pathExists(path)
			if err != nil {
				fmt.Println("getYear failed: ", err)
				return "", fmt.Errorf("getYear failed: %v", err)
			}
			if exists != true {
				valid = false
			}
		}

		if valid == false {
			fmt.Printf("Invalid file path! Year %v not found in one or more parent directories!\n", year)
			continue
		} else {
			break
		}
	}
	return yearStr, nil
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
*/
