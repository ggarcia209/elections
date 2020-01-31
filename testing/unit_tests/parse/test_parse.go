// CONCURRENCY TEST SUCCESSFUL
// UPDATE - SUCCESSFUL IMPLEMENTATION OF build_indv.go
package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/elections/donations"
	"github.com/elections/parse"
)

const (
	indv     = "txt_samples/test_indv.txt"
	cmteCont = "txt_samples/test_cmte_cont.txt"
	cn       = "txt_samples/test_cn.txt"
	ccl      = "txt_samples/test_ccl.txt"
	cm       = "txt_samples/test_cm.txt"
	disb     = "txt_samples/test_disb.txt"
	pac      = "txt_samples/test_pac.txt"
)

var wg sync.WaitGroup

func main() {

	parse.CreateDB()

	wg.Add(1)
	go indvContTest()
	wg.Add(1)
	go cmteContTest()
	wg.Add(1)
	go candTest()
	wg.Add(1)
	go cmteLinkTest()
	wg.Add(1)
	go cmteTest()
	// cmteFinTest(file, start) // FAILED
	wg.Add(1)
	go disbTest()

	wg.Wait()
	fmt.Println("all goroutines finished")

}

/* func indvContTest() {
	defer wg.Done()

	file, err := os.Open(indv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	for {
		queue, err := parse.Parse25IndvCont(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("***** indvContTest *****")
		fmt.Println("***** Results *****")
		fmt.Println("len queue: ", len(queue))
		fmt.Println("0: ", queue[0])
		fmt.Println("-1: ", queue[len(queue)-1])
		fmt.Println()

		if len(queue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
} */

func cmteContTest() {
	defer wg.Done()

	file, err := os.Open(cmteCont)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	start := 0
	for {
		queue, offset, err := parse.Parse25CmteCont(file, start)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		start = offset

		fmt.Println("***** cmteContTest *****")
		fmt.Println("***** Results *****")
		fmt.Println("len queue: ", len(queue))
		fmt.Println("0: ", queue[0])
		fmt.Println("-1: ", queue[len(queue)-1])
		fmt.Println()

		if len(queue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
}

func candTest() {
	defer wg.Done()

	file, err := os.Open(cn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	start := 0
	for {
		queue, offset, err := parse.Parse25Candidate(file, start)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		start = offset

		fmt.Println("***** candTest *****")
		fmt.Println("***** Results *****")
		fmt.Println("len queue: ", len(queue))
		fmt.Println("0: ", queue[0])
		fmt.Println("-1: ", queue[len(queue)-1])
		fmt.Println()

		if len(queue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
}

func cmteLinkTest() {
	defer wg.Done()

	file, err := os.Open(ccl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	start := 0
	for {
		queue, offset, err := parse.Parse25CmteLink(file, start)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		start = offset

		fmt.Println("***** cmteLinkTest *****")
		fmt.Println("***** Results *****")
		fmt.Println("len queue: ", len(queue))
		fmt.Println("0: ", queue[0])
		fmt.Println("-1: ", queue[len(queue)-1])
		fmt.Println()

		if len(queue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
}

func cmteTest() {
	defer wg.Done()

	file, err := os.Open(cm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	start := 0
	for {
		queue, offset, err := parse.Parse25Committee(file, start)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		start = offset

		fmt.Println("***** cmteTest *****")
		fmt.Println("***** Results *****")
		fmt.Println("len queue: ", len(queue))
		fmt.Println("0: ", queue[0])
		fmt.Println("-1: ", queue[len(queue)-1])
		fmt.Println()

		if len(queue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
}

// FAILED / CURRENTLY UNUSED
func cmteFinTest(file *os.File, start int64) {
	for {
		queue, offset, err := parse.Parse25CmteFin(file, start)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		start += offset

		fmt.Println("***** Results *****")
		fmt.Println("len queue: ", len(queue))
		fmt.Println("0: ", queue[0])
		fmt.Println("-1: ", queue[len(queue)-1])
		fmt.Println()

		if len(queue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
}

func indvContTest() {
	defer wg.Done()

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

func disbTest() {
	defer wg.Done()

	file, err := os.Open(disb)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// simulated start and recID values
	start := 0
	recID := 1
	for {
		DQueue, RQueue, offset, err := parse.Parse25Disbursements(file, start, int32(recID))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		start = offset
		recID += len(RQueue)

		fmt.Println("***** disbTest *****")
		fmt.Println("***** Results *****")
		fmt.Println("len DQueue: ", len(DQueue))
		fmt.Println("len RQueue: ", len(RQueue))
		// fmt.Println("0: ", queue[0])
		// fmt.Println("-1: ", queue[len(queue)-1])
		for _, rec := range RQueue {
			printRec(rec)
		}
		fmt.Println()

		if len(DQueue) < 25 {
			break
		}
	}

	fmt.Println("DONE")
}

func printRec(disb *donations.DisbRecipient) {
	fmt.Printf("ID: %s, Name: %s, City: %s, State: %s, Zip: %s, Disbursements: %v, TotalReceived: %d, AvgReceived: %d\n",
		disb.ID, disb.Name, disb.City, disb.State, disb.Zip, disb.Disbursements, disb.TotalReceived, disb.AvgReceived)
}
