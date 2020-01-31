package main

import (
	"fmt"
	"os"

	"github.com/elections/donations"
	"github.com/elections/parse"
	"github.com/elections/persist"
)

const indv = "txt_samples/test_indv.txt"

func main() {
	parse.CreateDB()
	persist.CreateLookupBucket()
	indvContTest()

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
	donorID := 1
	for {
		ICQueue, DQueue, offset, err := parse.Parse25IndvCont(file, start, int32(donorID))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		start = offset
		donorID += len(DQueue)

		fmt.Println("***** indvContTest *****")
		fmt.Println("***** Results *****")
		fmt.Println("len ICQueue: ", len(ICQueue))
		fmt.Println("len DQueue: ", len(DQueue))
		// fmt.Println("0: ", queue[0])
		// fmt.Println("-1: ", queue[len(queue)-1])
		for _, donor := range DQueue {
			printIndv(donor)
		}
		fmt.Println()

		if len(ICQueue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
}

func printIndv(indv *donations.Individual) {
	fmt.Printf("ID: %s, Name: %s, City: %s, State: %s, Zip: %s, Occupation: %s, Employer: %s, Donations: %v, TotalDonations: %d, TotalDonated: %d\n",
		indv.ID, indv.Name, indv.City, indv.State, indv.Zip, indv.Occupation, indv.Employer, indv.Donations, indv.TotalDonations, indv.TotalDonated)
}
