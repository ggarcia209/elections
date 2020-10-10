package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/elections/donations"
)

var filer = donations.CmteTxData{
	CmteID:                      "cmte00",
	TopIndvContributorsAmt:      map[string]float32{"indv1": 100, "indv2": 150, "indv3": 80, "indv4": 200, "indv5": 40, "indv6": 400, "indv7": 120, "indv8": 100, "indv9": 225},
	TopIndvContributorsTxs:      map[string]float32{"indv1": 2, "indv2": 3, "indv3": 1, "indv4": 3, "indv5": 1, "indv6": 4, "indv7": 3, "indv8": 1, "indv9": 5},
	TopIndvContributorThreshold: []interface{}{},
}

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
	fmt.Println("popLeast - len(s): ", len(a))
	return del
}

func main() {
	sortedEntries := sortTopX(filer.TopIndvContributorsAmt)
	least, err := setThresholdLeast3(sortedEntries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("least - init")
	for i, e := range least {
		printEntry(i, e)
	}
	fmt.Println()

	pop1 := least.popLeast()
	printEntry(0, pop1)
	fmt.Println("len least: ", len(least))

	pop2 := least.popLeast()
	printEntry(1, pop2)
	fmt.Println("len least: ", len(least))

	pop3 := least.popLeast()
	printEntry(2, pop3)
	fmt.Println("len least: ", len(least))
}

func printEntry(i int, e *donations.Entry) {
	fmt.Printf("%d) ID: %s / Total: %v\n", i, e.ID, e.Total)
}

// reSortLeast re-sorts the least 5 or 10 values when a new value breaks the threshold (least[len(least)-1].Total)
// and returns the ID of the key to be deleted and the new sorted list of least values
func reSortLeast(new *donations.Entry, es *Entries) string {
	copy := *es
	// if new.Total >= largest value in threshold list
	if new.Total >= copy[0].Total {
		// pop smallest value and get it's ID to delete from records
		delID := copy.popLeast().ID
		// update original list of entries by overwriting it with new copy
		es = &copy
		return delID
	}
	// value falls between threshold range:
	// add new value to copy of threshold list (# of items remains the same)
	// len + 1 (append) - 1 (popLeast)
	copy = append(copy, new)
	// update original list by overwriting it with copy
	es = &copy
	// reSort with new value included
	sort.Sort(es)
	// remove smallest item by value from list and return ID
	delID := es.popLeast().ID

	return delID
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

// TEST ONLY
func setThresholdLeast3(es Entries) (Entries, error) {
	if len(es) < 5 {
		return nil, fmt.Errorf("=etThresholdLeast5 failed: not enough elements in list")
	}
	return es[len(es)-3:], nil
}

// newEntry creats an entry struct from Top X Amt key/value pair
func newEntry(k string, v float32) *donations.Entry {
	return &donations.Entry{ID: k, Total: v}
}
