// TEST CacheAndPersistIndvDonor SUCCESSFUL
// TEST GetIndvDonor - RETRIEVAL FROM DB SUCCESSFUL

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/elections/donations"
	"github.com/elections/parse"
	"github.com/elections/persist"
)

const indv = "../parse/txt_samples/test_indv.txt"
const disb = "../parse/txt_samples/test_disb.txt"

func main() {
	persist.Init()
	indvContTest()
	disbTest()

	testDonor, err := persist.GetIndvDonor("ID00000024")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	testRec, err := persist.GetDisbRecipient("ID00000002")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printIndv(testDonor)
	printRec(testRec)

}

func indvContTest() {
	// defer wg.Done()

	// open file
	file, err := os.Open(indv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset("indv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get current donorID value; 1 if none
	donorID, err := persist.GetDonorID()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// parse file
	for {
		// parse 25 records per iteration
		ICQueue, DQueue, offset, err := parse.Parse25IndvCont(file, start, int32(donorID))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// update & store donorID after ojbects returned
		donorID += len(DQueue)
		err = persist.StoreDonorID(donorID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// save objects to disk
		err = persist.CacheAndPersistIndvDonor(DQueue)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// save offset value after objects persisted
		err = persist.LogOffset("indv", offset)
		if err != nil {
			fmt.Println("indvContTest failed: ", err)
			os.Exit(1)
		}
		start = offset

		// print each Donor object
		/* for _, indv := range DQueue {
			printIndv(indv)
		} */

		// break if at end of file
		if len(ICQueue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
}

func printIndv(indv *donations.Individual) {
	fmt.Printf("ID: %s, Name: %s, City: %s, State: %s, Zip: %s, Occupation: %s, Employer: %s, Donations: %v, TotalDonations: %d, TotalDonated: %d AvgDonation: %d\n",
		indv.ID, indv.Name, indv.City, indv.State, indv.Zip, indv.Occupation, indv.Employer, indv.Donations, indv.TotalDonations, indv.TotalDonated, indv.AvgDonation)
}

func disbTest() {
	// defer wg.Done()

	// open db
	file, err := os.Open(disb)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// get starting offset value; 0 if none
	start, err := persist.GetOffset("disb_rec")
	fmt.Println("starting offset: ", start)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get current RecID value; 1 if none
	recID, err := persist.GetRecID()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// parse file
	for {
		// parse 25 records for each iteration
		DQueue, RQueue, offset, err := parse.Parse25Disbursements(file, start, int32(recID))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// update and store RecID after objs returned
		recID += len(RQueue)
		err = persist.StoreRecID(recID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// save DisbRecipient objs to disk
		err = persist.CacheAndPersistDisbRecipient(RQueue)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// save offset value after objects persisted
		err = persist.LogOffset("disb_rec", offset)
		if err != nil {
			fmt.Println("disbTest failed: ", err)
			os.Exit(1)
		}
		fmt.Println("offset logged: ", offset)
		start = offset

		for _, obj := range RQueue {
			printRec(obj)
		}

		// break if at end of file
		if len(DQueue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
}

func printRec(disb *donations.DisbRecipient) {
	fmt.Printf("ID: %s, Name: %s, City: %s, State: %s, Zip: %s, Disbursements: %v, Total Disbursements: %d, TotalReceived: %v, AvgReceived: %v\n",
		disb.ID, disb.Name, disb.City, disb.State, disb.Zip, disb.Disbursements, disb.TotalDisbursements, disb.TotalReceived, disb.AvgReceived)
}

func deriveID(s string) int32 {
	conv := strings.Split(s, "")
	var ID []string
	for i, n := range conv {
		if i > 1 && n != "0" {
			ID = conv[i:]
			break
		}
	}
	donorIDint, _ := strconv.Atoi(strings.Join(ID, ""))
	donorID := int32(donorIDint)

	return donorID
}
