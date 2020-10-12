// Package admin contains operations for running the local admin console service.
// Only the functions in this package are exposed to the admin service; lower
// level source packages remain encapsulated.
// This file contains the operations for building the datasets from the
// raw .txt bulk data files.
// NOTE: logic is not UX optimized and may contain unresolved errors.
package admin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/elections/source/cache"
	"github.com/elections/source/databuilder"
	"github.com/elections/source/parse"
	"github.com/elections/source/persist"
	"github.com/elections/source/ui"
)

/* DB CREATE OPERATIONS */

// ProcessData contains options for processing raw data and creating secondary datasets
// All input directories must have the following files:
//   input/[year]/cmte/cn.txt - candidate master
//   input/[year]/cand/cm.txt - committee master
//   input/[year]/pac/webk.txt - PAC summary
//   input/[year]/ctx/itoth.txt - any tx between committees
//   input/[year]/indiv/itcont.txt - individiual contributions
//   input/[year]/exp/oppexp.txt - operating expenses
func ProcessData() error {
	fmt.Println("***** PROCESS DATA *****")
	opts := []string{"Process Raw Data", "Create Secondary Datasets", "Return"}
	menu := ui.CreateMenu("process-data-main", opts)

	for {
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("ProcessData failed: %v", err)
		}
		switch {
		case menu.OptionsMap[ch] == "Process Raw Data":
			err := processNewRecords()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ProcessData failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Create Secondary Datasets":
			err := createSecondaryDatasets()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ProcessData failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Return":
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

// processNewRecords processes the FEC bulk data files for the given year.
func processNewRecords() error {
	fmt.Println("******************************************")
	fmt.Println()

	// get input folder path
	input, err := persist.GetPath(true)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processNewRecords failed: %v", err)
	}
	if input != "" {
		fmt.Println("Set new input path?")
		yes := ui.Ask4confirm()
		if yes {
			input, err = getPath(true)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("processNewRecords failed: %v", err)
			}
		}
	} else {
		input, err = getPath(true)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processNewRecords failed: %v", err)
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
	if output != "" {
		fmt.Println("Set new output path?")
		yes := ui.Ask4confirm()
		if yes {
			output, err = getPath(false)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("processNewRecords failed: %v", err)
			}
		}
	} else {
		output, err = getPath(false)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processNewRecords failed: %v", err)
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
	candFinPath := filepath.Join(root, "cmpn", "webl.txt")
	cmtePath := filepath.Join(root, "cmte", "cm.txt")
	icPath := filepath.Join(root, "indiv", "itcont.txt")
	ccPath := filepath.Join(root, "ctx", "itoth.txt")
	disbPath := filepath.Join(root, "exp", "oppexp.txt")
	cmteFinPath := filepath.Join(root, "pac", "webk.txt")

	// initialize database and TopOverallData objects
	persist.Init(year)

	// process candidates and committee objects first
	err = processCandidates(year, candPath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	err = processCommittees(year, cmtePath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	if year >= "1996" { // no data prior to 1996
		err = processCmpnFinancials(year, candFinPath)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("ProcessNewRecords failed: %v", err)
		}
		err = processCmteFinancials(year, cmteFinPath)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("ProcessNewRecords failed: %v", err)
		}
	}

	// process transactions
	err = processCmteContributions(year, ccPath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	err = processIndvContributions(year, icPath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}

	if year >= "2004" { // no data prior to 2004
		err = processDisbursements(year, disbPath)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("ProcessNewRecords failed: %v", err)
		}
	}

	fmt.Println("PROCESS NEW RECORDS COMPLETE - YEAR: ", year)
	fmt.Println("******************************************")
	fmt.Println()

	return nil
}

func processCandidates(year, filePath string) error {
	// defer wg.Done()
	j := 0

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	fs := fi.Size()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cand")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	fmt.Println("got offset cand: ", start)

	// parse file
	for {
		// parse 10000 records per iteration
		objQueue, offset, err := parse.ScanCandidates(file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCandidates failed: %v", err)
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCandidates failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cand", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCandidates failed: %v", err)
		}
		start = offset
		j += len(objQueue)

		// break if at end of file
		if start >= fs {
			break
		}
	}
	fmt.Println("ending offset cands: ", start)
	fmt.Println("Candidate records scanned: ", j)
	fmt.Println("Candidates - DONE")

	return nil
}

func processCommittees(year, filePath string) error {
	// defer wg.Done()
	j, k := 0, 0

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCommittees failed: %v", err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	fs := fi.Size()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmte")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCommittees failed: %v", err)
	}
	fmt.Println("got offset cmte: ", start)

	// parse file
	for {
		// parse 10000 records per iteration
		objQueue, txDataQueue, offset, err := parse.ScanCommittees(file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCommittees failed: %v", err)
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCommittees failed: %v", err)
		}
		err = persist.StoreObjects(year, txDataQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCommittees failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmte", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCommittees failed: %v", err)
		}
		start = offset

		j += len(objQueue)
		k += len(txDataQueue)

		// break if at end of file
		if start >= fs {
			break
		}
	}

	fmt.Println("ending offset cmte: ", start)
	fmt.Println("Committee records scanned: ", j)
	fmt.Println("CmteTxData objects created: ", k)
	fmt.Println("Committees - DONE")
	return nil
}

func processCmpnFinancials(year, filePath string) error {
	// defer wg.Done()
	j := 0

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandFinancials failed: %v", err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	fs := fi.Size()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmpn_fin")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandFinancials failed: %v", err)
	}
	fmt.Println("got offset cmpn_fin: ", start)

	// parse file
	for {
		// parse 10000 records per iteration
		objQueue, offset, err := parse.ScanCmpnFin(file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCandFinancials failed: %v", err)
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCandFinancials failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmpn_fin", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCandFinancials failed: %v", err)
		}
		start = offset
		j += len(objQueue)

		// break if at end of file
		if start >= fs {
			break
		}
	}

	fmt.Println("ending offset cmpn_fin: ", start)
	fmt.Println("CmpnFinancials records scanned: ", j)
	fmt.Println("Candidates - DONE")

	return nil
}

func processCmteFinancials(year, filePath string) error {
	// defer wg.Done()
	j := 0

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCmteFinancials failed: %v", err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	fs := fi.Size()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmte_fin")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCmteFinancials failed: %v", err)
	}
	fmt.Println("got offset cmte_fin: ", start)

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, offset, err := parse.ScanCmteFin(file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCmteFinancials failed: %v", err)
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCmteFinancials failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmte_fin", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCmteFinancials failed: %v", err)
		}
		start = offset
		j += len(objQueue)

		// break if at end of file
		if start >= fs {
			break
		}
	}

	fmt.Println("ending offset cmte_fins: ", start)
	fmt.Println("CmteFinancials records scanned: ", j)
	fmt.Println("CmteFinancials - DONE")

	return nil
}

func processCmteContributions(year, filepath string) error {
	// defer wg.Done()
	j := 0

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCmteContributions failed: %v", err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	fs := fi.Size()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmte_cont")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCmteContributions failed: %v", err)
	}
	fmt.Println("got offset cmte cont ", start)

	// parse file
	for {
		// parse records
		txQueue, offset, err := parse.ScanContributions(year, file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}

		// create cache from record IDs
		c, err := cache.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}

		// persist items in cache
		ser := cache.SerializeCache(c)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmte_cont", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processCmteContributions failed: %v", err)
		}
		start = offset

		// break if at end of file
		if start >= fs {
			j += len(txQueue)
			break
		}
		j += 100000
	}

	fmt.Println("ending offset cmte_cont: ", start)
	fmt.Println("Committee Contribution records scanned: ", j)
	fmt.Println("Committee Contributions -  DONE")

	return nil
}

func processIndvContributions(year, filepath string) error {
	fmt.Println("starting Individual contributions...")
	i := 0
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processIndvContributions faield: %v", err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	fs := fi.Size()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "indv")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processIndvContributions faield: %v", err)
	}
	fmt.Println("got offset indv: ", start)

	// parse file
	for {
		// parse records
		txQueue, offset, err := parse.ScanContributions(year, file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processIndvContributions faield: %v", err)
		}

		// create cache from record IDs
		c, err := cache.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processIndvContributions faield: %v", err)
		}

		// break if finsished
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processIndvContributions faield: %v", err)
		}

		// update object data for each transaction
		err = databuilder.TransactionUpdate(year, txQueue, c)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processIndvContributions faield: %v", err)
		}

		// persist objects in cache
		ser := cache.SerializeCache(c)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processIndvContributions faield: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "indv", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processIndvContributions faield: %v", err)
		}
		start = offset

		// break if at end of file
		if start >= fs {
			i += len(txQueue)
			break
		}
		i += 100000
	}

	fmt.Println("ending offset indv: ", start)
	fmt.Println("Individual Contribution records scanned: ", i)
	fmt.Println("Individual Contributions -  DONE")

	return nil
}

func processDisbursements(year, filepath string) error {
	// defer wg.Done()
	i := 0
	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processDisbursements failed: %v", err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processCandidates failed: %v", err)
	}
	fs := fi.Size()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "disb")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("processDisbursements failed: %v", err)
	}
	fmt.Println("got offset disbs: ", start)

	for {
		// parse records
		txQueue, offset, err := parse.ScanDisbursements(year, file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}

		// create cache from
		c, err := cache.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}

		// persist objects in cache
		ser := cache.SerializeCache(c)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "disb", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("processDisbursements failed: %v", err)
		}
		start = offset
		i += 100000

		// break if at end of file
		if start >= fs {
			i += len(txQueue)
			break
		}
	}

	fmt.Println("ending offset disbs: ", start)
	fmt.Println("Disbursements records scanned: ", i)
	fmt.Println("Disbursements -  DONE")

	return nil
}

// pathExists checks to see if given file path is valid
// move to util
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
			fmt.Println()
		} else { // confirm overwrite
			fmt.Println(">>> Are you sure you want to overwrite the existing path?")
			yes := ui.Ask4confirm()
			if !yes {
				fmt.Printf("Continuing with current %s path: %s\n", name, path)
				fmt.Println()
			} else { // set new path
				fmt.Println(msg)
				path = ui.GetPathFromUser()
				err := persist.LogPath(path, input)
				fmt.Printf("New path: %s saved\n", path)
				fmt.Println()
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
	fmt.Println()
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("getPath failed: %v", err)
	}
	return path, nil
}

/* DEPRECATED

func getYear() (string, error) {
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

*/
