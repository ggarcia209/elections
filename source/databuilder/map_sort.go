package databuilder

import (
	"fmt"
	"sort"

	"github.com/elections/source/donations"
)

// Entries is a list of entries to be sorted.
type Entries []*donations.Entry

func (s Entries) Len() int           { return len(s) }
func (s Entries) Less(i, j int) bool { return s[i].Total > s[j].Total }
func (s Entries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// PopLeast pops the smalles value from the list of least values
func (s *Entries) popLeast() *donations.Entry {
	a := *s
	if len(a) == 0 {
		return &donations.Entry{}
	}
	del := a[len(a)-1]
	*s = a[:len(a)-1]
	return del
}

// reSortLeast re-sorts the least 5 or 10 values when a new value breaks the threshold (least[len(least)-1].Total)
// and returns the ID of the key to be deleted and the new sorted list of least values
func reSortLeast(new *donations.Entry, es Entries) (string, Entries) {
	copy := es
	// if new.Total >= largest value in threshold list
	if new.Total >= copy[0].Total {
		// pop smallest value and get it's ID to delete from records
		delID := copy.popLeast().ID
		return delID, copy
	}
	// value falls between threshold range:
	// add new value to copy of threshold list (# of items remains the same)
	//   len + 1 (append) - 1 (popLeast)
	copy = append(copy, new)
	// reSort with new value included
	sort.Sort(copy)
	// remove smallest item by value from list and return ID
	delID := copy.popLeast().ID

	return delID, copy
}

// sortTopX sorts the Top x Donors/Recipients maps from greatest -> smallest (decreasing order)
func sortTopX(m map[string]float32) Entries {
	var es Entries
	for k, v := range m {
		es = append(es, &donations.Entry{ID: k, Total: v})
	}
	sort.Sort(es)

	return es
}

// setThresholdLeast10 sets a threshold with the smallest 10 values in the Top x
func setThresholdLeast10(es Entries) (Entries, error) {
	if len(es) < 10 {
		for i, e := range es {
			fmt.Printf("%d) ID: %s\tTotal: %v\n", i, e.ID, e.Total)
		}
		return nil, fmt.Errorf("setThresholdLeast10 failed: not enough elements in list")
	}

	return es[len(es)-10:], nil
}

// newEntry creats an entry struct from Top X Amt key/value pair
func newEntry(k string, v float32) *donations.Entry {
	return &donations.Entry{ID: k, Total: v}
}

// TEMPORARILY DEPRECATED
// setThresholdLeast5 sets a threshold with the smallest 5 values in the Top x
// sorted greatest -> smallest
/* func setThresholdLeast5(es Entries) (Entries, error) {
	if len(es) < 5 {
		return nil, fmt.Errorf("=etThresholdLeast5 failed: not enough elements in list")
	}
	return es[len(es)-5:], nil
} */
