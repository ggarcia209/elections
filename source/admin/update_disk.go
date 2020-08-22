package admin

import (
	"fmt"
	"os"

	"github.com/elections/source/cache"
	"github.com/elections/source/databuilder"
	"github.com/elections/source/donations"
	"github.com/elections/source/parse"
	"github.com/elections/source/persist"
)

/* IN PROGRESS */
// determine if new records are appended to .txt files

// UpdateRecordsOnDisk updates the on disk database with all files contained
// in the /update folder for the given year.
func UpdateRecordsOnDisk(year string) error {
	//path, err := persist.GetPath()
	// open update subdir and read files
	dirName := "../../raw_data/" + year + "/update"
	fInfo, err := getFiles(dirName)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("UpdateRecordsOnDisk failed: %v", err)
	}

	for _, f := range fInfo {
		fmt.Println("name: ", f.Name())
	}
	return nil
}

func getFiles(dirName string) ([]os.FileInfo, error) {
	f, err := os.Open(dirName)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("getFiles failed: %v", err)
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("getFiles failed: %v", err)
	}

	for _, file := range files {
		fmt.Println("/update/", file.Name())
	}

	return files, nil
}

// idempotent - will not overwrite existing objects
func updateCandidates(year, fileName string) error {
	// open file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCandidates failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cand - update")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCandidates failed: %v", err)
	}

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, offset, err := parse.ScanCandidates(file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCandidates failed: %v", err)
		}

		// filter results - persist nil objects only to maintain data integrity
		lookupIDs := []string{}
		objMap := make(map[string]interface{})
		for _, obj := range objQueue {
			lookupIDs = append(lookupIDs, obj.(*donations.Individual).ID)
			objMap[obj.(*donations.Individual).ID] = obj
		}

		_, nilIDs, err := persist.BatchGetByID(year, "individuals", lookupIDs)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCandidates failed: %v", err)
		}

		newObjs := []interface{}{}
		for _, id := range nilIDs {
			newObjs = append(newObjs, objMap[id])
		}

		// save objects to disk
		err = persist.StoreObjects(year, newObjs)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCandidates failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cand - update", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCandidates failed: %v", err)
		}
		start = offset

		// break if at end of file
		if len(objQueue) < 100 {
			break
		}
	}

	// reset offset value at EOF
	err = persist.LogOffset(year, "cand - update", 0)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCandidates failed: %v", err)
	}

	fmt.Println("Candidates - UPDATE - DONE")

	return nil
}

// NOT idempotent - will overwrite CmteTxData objs
func updateCommittees(year, fileName string) error {
	// open file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCommittees failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmte - update")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCommittees failed: %v", err)
	}

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, txDataQueue, offset, err := parse.ScanCommittees(file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCommittees failed: %v", err)
		}

		// filter results - persist nil objects only to maintain data integrity
		lookupIDs := []string{}
		objMap := make(map[string]interface{})
		for _, obj := range objQueue {
			lookupIDs = append(lookupIDs, obj.(*donations.Committee).ID)
			objMap[obj.(*donations.Committee).ID] = obj
		}

		_, nilIDs, err := persist.BatchGetByID(year, "committees", lookupIDs)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCommittees failed: %v", err)
		}

		newObjs := []interface{}{}
		for _, id := range nilIDs {
			newObjs = append(newObjs, objMap[id])
		}

		// repeat for txDataQueue
		txDataMap := make(map[string]interface{})
		for _, obj := range txDataQueue {
			txDataMap[obj.(*donations.Committee).ID] = obj
		}

		// repeat BatchGetByID for redundancy
		_, nilIDs, err = persist.BatchGetByID(year, "cmte_tx_data", lookupIDs)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCommittees failed: %v", err)
		}

		newTxData := []interface{}{}
		for _, id := range nilIDs {
			newTxData = append(newTxData, txDataMap[id])
		}

		// save objects to disk
		err = persist.StoreObjects(year, newObjs)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCommittees failed: %v", err)
		}

		err = persist.StoreObjects(year, newTxData)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCommittees failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmte - update", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCommittees failed: %v", err)
		}
		start = offset

		// break if at end of file
		if len(objQueue) < 100 {
			break
		}
	}

	// reset offset value at EOF
	err = persist.LogOffset(year, "cmte - update", 0)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCommittees failed: %v", err)
	}

	fmt.Println("Committees - UPDATE - DONE")

	return nil
}

// idempotent - data will be overwritten with identical data
// special case - do not filter results - update all records
func updateCmteFinancials(year, filePath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCmteFinancials failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmte_fin - update")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCmteFinancials failed: %v", err)
	}

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, offset, err := parse.ScanCmteFin(file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCmteFinancials failed: %v", err)
		}

		// save objects to disk
		err = persist.StoreObjects(year, objQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCmteFinancials failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmte_fin - update", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCmteFinancials failed: %v", err)
		}
		start = offset

		// break if at end of file
		if len(objQueue) < 100 {
			break
		}
	}

	// reset offset value at EOF
	err = persist.LogOffset(year, "cmte_fin - update", 0)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCmteFinancials failed: %v", err)
	}

	fmt.Println("Committee Financials - UPDATE - DONE")

	return nil
}

// NOT idempotent - susequent updates from same file will corrupt data if called after EOF
func updateCmteContributions(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCmteContributions failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "cmte_cont - update")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCmteContributions failed: %v", err)
	}

	// parse file
	for {
		// parse records
		txQueue, offset, err := parse.ScanContributions(year, file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCmteContributions failed: %v", err)
		}

		// create cache from record IDs
		c, err := cache.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCmteContributions failed: %v", err)
		}

		// break if finished
		if len(c) == 0 {
			fmt.Println("updateCmteContributions DONE - no cache")
			return nil
		}

		// update cached objects for each transaction
		err = databuilder.TransactionUpdate(year, txQueue, c)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCmteContributions failed: %v", err)
		}

		// persist items in cache
		ser := cache.SerializeCache(c)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCmteContributions failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "cmte_cont - update", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateCmteContributions failed: %v", err)
		}
		start = offset

		// break if at end of file
		if len(txQueue) < 1000 {
			break
		}
	}

	// reset offset value if EOF
	err = persist.LogOffset(year, "cmte_cont - update", 0)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateCmteContributions failed: %v", err)
	}

	fmt.Println("Committee Contributions - UPDATE - DONE")

	return nil
}

// NOT idempotent - susequent updates from same file will corrupt data if called after EOF
func updateIndvContributions(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateIndvContributions failed %v: ", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "indv - update")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateIndvContributions failed %v: ", err)
	}

	// parse file
	for {
		// parse records
		txQueue, offset, err := parse.ScanContributions(year, file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateIndvContributions failed %v: ", err)
		}

		// create cache from record IDs
		c, err := cache.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateIndvContributions failed %v: ", err)
		}

		// break if finsished
		if len(c) == 0 {
			fmt.Println("processIndvContributions DONE - no cache")
			return nil
		}

		// update object data for each transaction
		err = databuilder.TransactionUpdate(year, txQueue, c)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateIndvContributions failed %v: ", err)
		}

		// persist objects in cache
		ser := cache.SerializeCache(c)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateIndvContributions failed %v: ", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "indv - update", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateIndvContributions failed %v: ", err)
		}
		start = offset

		// break if at end of file
		if len(txQueue) < 1000 {
			break
		}
	}

	fmt.Println("Individual Contributions - UPDATE -  DONE")

	// reset offset value at EOF
	err = persist.LogOffset(year, "indv - update", 0)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateIndvContributions failed %v: ", err)
	}

	return nil
}

// NOT idempotent - susequent updates from same file will corrupt data if called after EOF
func updateDisbursements(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateDisbursements failed: %v", err)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset(year, "disb - update")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateDisbursements failed: %v", err)
	}

	for {
		// parse records
		txQueue, offset, err := parse.ScanDisbursements(year, file, start)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateDisbursements failed: %v", err)
		}

		// create cache from
		c, err := cache.CreateCache(year, txQueue)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateDisbursements failed: %v", err)
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
			return fmt.Errorf("updateDisbursements failed: %v", err)
		}

		// persist objects in cache
		ser := cache.SerializeCache(c)
		err = persist.StoreObjects(year, ser)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateDisbursements failed: %v", err)
		}

		// save offset value after objects persisted
		err = persist.LogOffset(year, "disb - update", offset)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("updateDisbursements failed: %v", err)
		}
		start = offset

		// break if at end of file
		if len(txQueue) < 1000 {
			break
		}
	}

	// reset offset value at EOF
	err = persist.LogOffset(year, "disb - update", 0)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("updateDisbursements failed: %v", err)
	}

	fmt.Println("Disbursements - UPDATE -  DONE")

	return nil
}
