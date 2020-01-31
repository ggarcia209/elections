package main

import (
	"fmt"
	"os"

	"github.com/elections/donations"
	"github.com/elections/parse"
	"github.com/elections/persist"
)

const indv = "../parse/txt_samples/test_indv.txt"
const disb = "../parse/txt_samples/test_disb.txt"

func main() {
	persist.CreateDB()
	persist.CreateLookupBucket()
	persist.CreateIDBucket()
	offset, err := persist.GetOffset("indv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(offset)

	// indvContTest()
	// disbTest()

	/* testDonor, err := persist.GetIndvDonor("ID00000001")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	 testRec, err := persist.GetDisbRecipient("ID00000001")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} */

	// printIndv(testDonor)
	// printRec(testRec)

}

func indvContTest() {
	// defer wg.Done()

	file, err := os.Open(indv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// simulated start and donorID values
	start := 0
	donorID, err := persist.GetDonorID()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for {
		ICQueue, DQueue, offset, err := parse.Parse25IndvCont(file, start, int32(donorID))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		start = offset
		donorID += len(DQueue)

		// test
		err = persist.StoreDonorID(donorID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		seen := make(map[string]bool)
		err = persist.CacheAndPersistIndvDonor(DQueue, seen)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// end test

		/* fmt.Println("***** indvContTest *****")
		fmt.Println("***** Results *****")
		fmt.Println("len ICQueue: ", len(ICQueue))
		fmt.Println("len DQueue: ", len(DQueue))
		// fmt.Println("0: ", queue[0])
		// fmt.Println("-1: ", queue[len(queue)-1])
		for _, donor := range DQueue {
			printIndv(donor)
		}
		fmt.Println() */

		for _, indv := range DQueue {
			printIndv(indv)
		}

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
