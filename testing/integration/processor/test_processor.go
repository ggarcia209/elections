package main

import (
	"fmt"
	"os"

	"github.com/elections/databuilder"
	"github.com/elections/donations"
	"github.com/elections/parse"
	"github.com/elections/persist"
)

func main() {
	// input file paths
	candPath := "../../test_data/test_cands.txt"
	cmtePath := "../../test_data/test_cmtes.txt"
	icPath := "../../test_data/test_ics.txt"
	ccPath := "../../test_data/test_ccs.txt"
	disbPath := "../../test_data/test_disbs.txt"

	year := "2018" // derived from flag inputs

	// initialize database and TopOverallData objects
	persist.Init()
	topOverallList := donations.InitTopOverallDataObjs(50)
	err := persist.CacheTopOverall(year, topOverallList)
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
	fmt.Println("MAIN DONE")
}

func indvContTest(year, filepath string) error {
	// defer wg.Done()

	// open file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset("indv")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// get current donorID value; 1 if none
	donorID, err := persist.GetDonorID()
	if err != nil {
		fmt.Println(err)
		return err
	}
	cacheStart := true

	// parse file
	for {
		// parse 25 records per iteration
		ICQueue, DQueue, offset, err := parse.Parse25IndvCont(year, file, start, int32(donorID))
		if err != nil {
			fmt.Println(err)
			return err
		}

		// update & store donorID after ojbects returned
		donorID += len(DQueue)
		err = persist.StoreDonorID(donorID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save objects to disk
		if start > 0 {
			cacheStart = false
		}
		err = persist.CacheAndPersistIndvDonor(year, DQueue, cacheStart)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// update Indvidual Donors and receiving committees
		err = databuilder.IndvContUpdate(year, ICQueue)
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
		cacheStart = false

		// break if at end of file
		if len(ICQueue) < 25 {
			break
		}
	}

	fmt.Println("indvContTest DONE")
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

	// get current recID value; 1 if none
	recID, err := persist.GetRecID()
	if err != nil {
		fmt.Println(err)
		return err
	}
	cacheStart := true

	// parse file
	for {
		// parse 25 records per iteration
		DQueue, RQueue, offset, err := parse.Parse25Disbursements(year, file, start, int32(recID))
		if err != nil {
			fmt.Println(err)
			return err
		}

		// update & store recID after ojbects returned
		recID += len(DQueue)
		err = persist.StoreRecID(recID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save objects to disk
		if start > 0 {
			cacheStart = false
		}
		err = persist.CacheAndPersistDisbRecipient(year, RQueue, cacheStart)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// update Indvidual Donors and receiving committees
		err = databuilder.CmteDisbUpdate(year, DQueue)
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
		cacheStart = false

		// break if at end of file
		if len(DQueue) < 25 {
			break
		}
	}

	fmt.Println("disbTest DONE")
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
		objQueue, offset, err := parse.Parse25CmteCont(file, start)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// update sending and receiving committee values
		err = databuilder.CmteContUpdate(year, objQueue)
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
		if len(objQueue) < 25 {
			break
		}
	}

	fmt.Println("cmteContTest DONE")
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
	cacheStart := true

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, offset, err := parse.Parse25Candidate(file, start)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save objects to disk
		if start > 0 {
			cacheStart = false
		}
		err = persist.InitialCacheCand(year, objQueue, cacheStart)
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
		cacheStart = false

		// break if at end of file
		if len(objQueue) < 25 {
			break
		}
	}

	fmt.Println("candTest DONE")
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
	cacheStart := true

	// parse file
	for {
		// parse 25 records per iteration
		objQueue, offset, err := parse.Parse25Committee(file, start)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// save objects to disk
		if start > 0 {
			cacheStart = false
		}
		err = persist.InitialCacheCmte(year, objQueue, cacheStart)
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
		cacheStart = false

		// break if at end of file
		if len(objQueue) < 25 {
			break
		}
	}

	fmt.Println("cmteTest DONE")
	return nil
}
