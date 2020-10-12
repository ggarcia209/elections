// Package databuilder conatins operations for updating datasets in memory.
// This package is primarily used by the admin service to create the
// primary datasets from the raw input, followed by the secondary
// datasets.
// This file contains operations for ranking objects and calculating
// the yearly totals for given categories.
package databuilder

import (
	"fmt"

	"github.com/elections/source/donations"
)

// CompareTopOverall compares an object's total to the smalles value in the map.
// Entry is added to the map and smalles enty is removed from map if
// new total > least.
func CompareTopOverall(ID string, total float32, od *donations.TopOverallData) error {
	if total == 0 {
		return nil
	}

	if od.Amts == nil {
		od.Amts = make(map[string]float32)
	}
	// add to Amts map if len(Amts) < Size Limit
	if len(od.Amts) < od.SizeLimit {
		od.Amts[ID] = total
		return nil
	}

	// check threshold when updating existing entry
	if od.Amts[ID] != 0 {
		od.Amts[ID] = total
		th, err := checkODThreshold(ID, od.Amts, od.Threshold)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("CompareTopOverall failed: %v", err)
		}
		od.Threshold = th
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
			fmt.Println(err)
			return fmt.Errorf("CompareTopOverall failed: %v", err)
		}
	} else {
		for _, entry := range od.Threshold {
			least = append(least, entry)
		}
	}

	// compare sen cmte's total received value to threshold
	threshold := least[len(least)-1].Total // last/smallest obj in least
	if total > threshold {
		new := newEntry(ID, total)
		delID, newEntries := reSortLeast(new, least)
		least = newEntries
		delete(od.Amts, delID)
		od.Amts[ID] = total
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

// UpdateYearlyTotal updates the total of a YearlyTotal object.
func UpdateYearlyTotal(amt float32, yt *donations.YearlyTotal) {
	yt.Total += amt
}

// Check to see if previous total of entry is in threshold range when updating existing entry.
func checkODThreshold(newID string, m map[string]float32, th []*donations.Entry) ([]*donations.Entry, error) {
	inRange := false
	check := map[string]bool{newID: true}
	for _, e := range th {
		if check[e.ID] == true {
			inRange = true
		}
	}
	if inRange {
		es := sortTopX(m)
		newRange, err := setThresholdLeast10(es)
		if err != nil {
			fmt.Println(err)
			return []*donations.Entry{}, fmt.Errorf("checkODThreshold failed: %v", err)
		}
		// update object's threshold list
		newTh := []*donations.Entry{}
		for _, entry := range newRange {
			newTh = append(newTh, entry)
		}
		return newTh, nil
	}
	return th, nil
}
