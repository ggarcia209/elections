package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/elections/databuilder"
	"github.com/elections/parse"
	"github.com/elections/persist"
)

var year string

func main() {

	year, err := getYear()
	if err != nil {
		fmt.Println("getYear failed: ", err)
		os.Exit(1)
	}

	// file path variables
	indvCont := "../raw_data/cmte_cont/" + year + "/itcont.txt"
	cmteCont := "../raw_data/cmte_cont/" + year + "/itoth.txt"
	disb := "../raw_data/cmte_cont/" + year + "/oppenxp.txt"
	cand := "../raw_data/cmte_cont/" + year + "/cn.txt"
	cmte := "../raw_data/cmte_cont/" + year + "/cm.txt"

	fmt.Println(indvCont, cmteCont, disb, cand, cmte)

	// create Candidate objects
	err = getCandidates(year, cand)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// create Committe objects
	err = getCommittees(year, cmte)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// process CmteContributions objects
	err = getCmteContributions(year, cmteCont)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// process IndvContributions objects
	err = getIndvContributions(year, indvCont)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// process Disbursement objects
	err = getDisbursements(year, disb)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

}

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

func getIndvContributions(year, path string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("getIndvContributions failed: ", err)
		return fmt.Errorf("getIndvContributions failed: %v", err)
	}
	defer file.Close()
	fmt.Printf("file open: %s\n", path)

	// get starting offset value; 0 if none
	offsetKey := "indv_cont"
	start, err := persist.GetOffset(offsetKey)
	if err != nil {
		fmt.Println("getIndvContributions failed: ", err)
		return fmt.Errorf("getIndvContributions failed: %v", err)
	}
	fmt.Printf("offset retreived: %v\n", start)

	// parse file
	begin := true
	for {
		// parse 25 records per iteration
		ICQueue, DQueue, offset, err := parse.Parse25IndvCont(year, file, start)
		if err != nil {
			if err != nil {
				fmt.Println("getIndvContributions failed: ", err)
				return fmt.Errorf("getIndvContributions failed: %v", err)
			}
		}

		// save objects to disk
		err = persist.CacheAndPersistIndvDonor(year, DQueue, begin)
		if err != nil {
			fmt.Println("getIndvContributions failed: ", err)
			return fmt.Errorf("getIndvContributions failed: %v", err)
		}
		fmt.Printf("year %s Individual objects persisted: %v\n", year, len(DQueue))

		// update IndvDonor objects
		err = databuilder.IndvContUpdate(year, ICQueue)
		if err != nil {
			fmt.Println("getIndvContributions failed: ", err)
			return fmt.Errorf("getIndvContributions failed: %v", err)
		}
		fmt.Printf("year %s Individual objects updated: %v\n", year, len(ICQueue))

		// save offset value after objects persisted
		err = persist.LogOffset(offsetKey, offset)
		if err != nil {
			fmt.Println("getIndvContributions failed: ", err)
			return fmt.Errorf("getIndvContributions failed: %v", err)
		}
		start = offset
		fmt.Printf("offset stored: %v\n", offset)

		// break if at end of file
		if len(ICQueue) < 25 {
			break
		}
		begin = false
	}

	fmt.Println("INDVIDUAL CONTRIBUTIONS DONE")
	return nil
}

func getCmteContributions(year, path string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("getCmteContributions failed: ", err)
		return fmt.Errorf("getCmteContributions failed: %v", err)
	}
	defer file.Close()
	fmt.Printf("file open: %s\n", path)

	// get starting offset value; 0 if none
	offsetKey := "cmte_cont"
	start, err := persist.GetOffset(offsetKey)
	if err != nil {
		fmt.Println("getCmteContributions failed: ", err)
		return fmt.Errorf("getCmteContributions failed: %v", err)
	}
	fmt.Printf("offset retreived: %v\n", start)

	// parse file
	for {
		// parse 25 records per iteration
		ObjQueue, offset, err := parse.Parse25CmteCont(file, start)
		if err != nil {
			fmt.Println("getCmteContributions failed: ", err)
			return fmt.Errorf("getCmteContributions failed: %v", err)
		}

		// update Committee objects
		err = databuilder.CmteContUpdate(year, ObjQueue)
		if err != nil {
			fmt.Println("getCmteContributions failed: ", err)
			return fmt.Errorf("getCmteContributions failed: %v", err)
		}
		fmt.Printf("year %s Committee objects updated: %v\n", year, len(ObjQueue))

		// save offset value after objects persisted
		err = persist.LogOffset(offsetKey, offset)
		if err != nil {
			fmt.Println("getCmteContributions failed: ", err)
			return fmt.Errorf("getCmteContributions failed: %v", err)
		}
		start = offset
		fmt.Printf("offset stored: %v\n", offset)

		// break if at end of file
		if len(ObjQueue) < 25 {
			break
		}
	}

	fmt.Println("COMMITTEE CONTRIBUTIONS DONE")
	return nil
}

func getCandidates(year, path string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("getCandidates failed: ", err)
		return fmt.Errorf("getCandidates failed: %v", err)
	}
	defer file.Close()
	fmt.Printf("file open: %s\n", path)

	// get starting offset value; 0 if none
	offsetKey := "cand"
	start, err := persist.GetOffset(offsetKey)
	if err != nil {
		fmt.Println("getCandidates failed: ", err)
		return fmt.Errorf("getCandidates failed: %v", err)
	}
	fmt.Printf("offset retreived: %v\n", start)

	// parse file
	begin := true
	for {
		// parse 25 records per iteration
		ObjQueue, offset, err := parse.Parse25Candidate(file, start)
		if err != nil {
			fmt.Println("getCandidates failed: ", err)
			return fmt.Errorf("getCandidates failed: %v", err)
		}

		// save objects to disk
		err = persist.InitialCacheCand(year, ObjQueue, begin)
		if err != nil {
			fmt.Println("getCandidates failed: ", err)
			return fmt.Errorf("getCandidates failed: %v", err)
		}
		fmt.Printf("year %s Candidate objects persisted: %v\n", year, len(ObjQueue))

		// save offset value after objects persisted
		err = persist.LogOffset(offsetKey, offset)
		if err != nil {
			fmt.Println("getCandidates failed: ", err)
			return fmt.Errorf("getCandidates failed: %v", err)
		}
		start = offset
		fmt.Printf("offset stored: %v\n", offset)

		// break if at end of file
		if len(ObjQueue) < 25 {
			break
		}
		begin = false
	}

	fmt.Println("CANDIDATES DONE")
	return nil
}

func getCommittees(year, path string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("getCommittees failed: ", err)
		return fmt.Errorf("getCommittees failed: %v", err)
	}
	defer file.Close()
	fmt.Printf("file open: %s\n", path)

	// get starting offset value; 0 if none
	offsetKey := "cmte"
	start, err := persist.GetOffset(offsetKey)
	if err != nil {
		fmt.Println("getCommittees failed: ", err)
		return fmt.Errorf("getCommittees failed: %v", err)
	}
	fmt.Printf("offset retreived: %v\n", start)

	// parse file
	begin := true
	for {
		// parse 25 records per iteration
		ObjQueue, offset, err := parse.Parse25Committee(file, start)
		if err != nil {
			fmt.Println("getCommittees failed: ", err)
			return fmt.Errorf("getCommittees failed: %v", err)
		}

		// save objects to disk
		err = persist.InitialCacheCmte(year, ObjQueue, begin)
		if err != nil {
			fmt.Println("getCommittees failed: ", err)
			return fmt.Errorf("getCommittees failed: %v", err)
		}
		fmt.Printf("year %s Committee objects persisted: %v\n", year, len(ObjQueue))

		// save offset value after objects persisted
		err = persist.LogOffset(offsetKey, offset)
		if err != nil {
			fmt.Println("getCommittees failed: ", err)
			return fmt.Errorf("getCommittees failed: %v", err)
		}
		start = offset
		fmt.Printf("offset stored: %v\n", offset)

		// break if at end of file
		if len(ObjQueue) < 25 {
			break
		}
		begin = false
	}

	fmt.Println("COMMITTEES DONE")
	return nil
}

func getDisbursements(year, path string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("getDisbursements failed: ", err)
		return fmt.Errorf("getDisbursements failed: %v", err)
	}
	defer file.Close()
	fmt.Printf("file open: %s\n", path)

	// get starting offset value; 0 if none
	offsetKey := "disb"
	start, err := persist.GetOffset(offsetKey)
	if err != nil {
		fmt.Println("getDisbursements failed: ", err)
		return fmt.Errorf("getDisbursements failed: %v", err)
	}
	fmt.Printf("offset retreived: %v\n", start)

	// parse file
	begin := true
	for {
		// parse 25 records per iteration
		TxQueue, ObjQueue, offset, err := parse.Parse25Disbursements(year, file, start)
		if err != nil {
			if err != nil {
				fmt.Println("getDisbursements failed: ", err)
				return fmt.Errorf("getDisbursements failed: %v", err)
			}
		}

		// save objects to disk
		err = persist.CacheAndPersistDisbRecipient(year, ObjQueue, begin)
		if err != nil {
			fmt.Println("getDisbursements failed: ", err)
			return fmt.Errorf("getDisbursements failed: %v", err)
		}
		fmt.Printf("year %s DisbRecipient objects persisted: %v\n", year, len(ObjQueue))

		// update IndvDonor objects
		err = databuilder.CmteDisbUpdate(year, TxQueue)
		if err != nil {
			fmt.Println("getDisbursements failed: ", err)
			return fmt.Errorf("getDisbursements failed: %v", err)
		}
		fmt.Printf("year %s DisbRecipients objects updated: %v\n", year, len(ObjQueue))

		// save offset value after objects persisted
		err = persist.LogOffset(offsetKey, offset)
		if err != nil {
			fmt.Println("getDisbursements failed: ", err)
			return fmt.Errorf("getDisbursements failed: %v", err)
		}
		start = offset
		fmt.Printf("offset stored: %v\n", offset)

		// break if at end of file
		if len(TxQueue) < 25 {
			break
		}
		begin = false
	}

	fmt.Println("DISBURSEMENTS DONE")
	return nil
}
