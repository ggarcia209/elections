package testDB

import (
	"fmt"

	"github.com/elections/donations"
)

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
		least, err = setThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateOverall failed: ", err)
			return fmt.Errorf("updateOverall failed: %v", err)
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
		delID := reSortLeast(new, &least)
		delete(od.Amts, delID)
		od.Amts[e.ID] = e.Total
	}

	return nil
}
