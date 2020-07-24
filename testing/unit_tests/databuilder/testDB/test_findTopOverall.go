package testDB

import (
	"fmt"

	"github.com/elections/donations"
)

func TestCompareTopOverall() error {
	od1 := &donations.TopOverallData{
		Category:  "top_indv",
		Amts:      map[string]float32{"indv00": 100, "indv01": 200, "indv02": 50, "indv04": 150, "indv05": 250},
		Threshold: []*donations.Entry{},
		SizeLimit: 5,
	}
	od2 := &donations.TopOverallData{
		Category:  "top_indv",
		Amts:      map[string]float32{"indv00": 100, "indv01": 200, "indv02": 50},
		Threshold: []*donations.Entry{},
		SizeLimit: 5,
	}
	e := &donations.Entry{ID: "indv03", Total: 75}

	err := compareTopOverall(e, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("0: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	err = compareTopOverall(e, od2)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("0: od2")
	fmt.Println(od2)
	fmt.Println()

	e2 := &donations.Entry{ID: "indv06", Total: 300}
	err = compareTopOverall(e2, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("1: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e3 := &donations.Entry{ID: "indv07", Total: 120}
	err = compareTopOverall(e3, od1) // call removes lowest value but fails to update threshold
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("2: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e4 := &donations.Entry{ID: "indv08", Total: 135}
	err = compareTopOverall(e4, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("3: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e5 := &donations.Entry{ID: "indv09", Total: 175}
	err = compareTopOverall(e5, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("4: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e6 := &donations.Entry{ID: "indv10", Total: 210}
	err = compareTopOverall(e6, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("5: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e7 := &donations.Entry{ID: "indv11", Total: 500}
	err = compareTopOverall(e7, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("6: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	return nil
}

func printODThreshold(th []*donations.Entry) {
	for i, e := range th {
		fmt.Printf("%d: ID: %s\tTotal: %v\n", i, e.ID, e.Total)
	}
}

// compare obj to top overall threshold
func compareTopOverall(e *donations.Entry, od *donations.TopOverallData) error {
	// add to Amts map if len(Amts) < Size Limit
	if len(od.Amts) < od.SizeLimit {
		od.Amts[e.ID] = e.Total
		return nil
	}

	// if len(Amts) == SizeLimit
	// set/reset least threshold list
	var least Entries
	var err error
	if len(od.Threshold) == 0 {
		es := sortTopX(od.Amts)
		least, err = setThresholdLeast3(es)
		if err != nil {
			fmt.Println("compareTopOverall failed: ", err)
			return fmt.Errorf("compareTopOverall failed: %v", err)
		}
	} else {
		for _, entry := range od.Threshold {
			least = append(least, entry)
		}
	}

	// compare sen cmte's total received value to threshold
	threshold := least[len(least)-1].Total // last/smallest obj in least
	if e.Total > threshold {
		new := newEntry(e.ID, e.Total)
		delID, newEntries := reSortLeast(new, least)
		least = newEntries
		delete(od.Amts, delID)
		od.Amts[e.ID] = e.Total
	} else {
		newTh := []*donations.Entry{}
		for _, e := range least {
			newTh = append(newTh, e)
		}
		od.Threshold = append(od.Threshold[:0], newTh...)
		return nil
	}

	// update threshold
	newTh := []*donations.Entry{}
	for _, e := range least {
		newTh = append(newTh, e)
	}
	od.Threshold = append(od.Threshold[:0], newTh...)

	return nil
}
