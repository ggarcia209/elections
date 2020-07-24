package testDB

import (
	"fmt"
	"sort"

	"github.com/elections/donations"
)

/* Unit Tests */
var TestMap = map[string]float32{"indv01": 100, "indv02": 200, "indv03": 50, "indv04": 0, "indv05": 150}

// SUCCESS
func testSortTopX(m map[string]float32) Entries {
	fmt.Println("sortTopX test")
	es := sortTopX(m)
	printEntries(es)
	fmt.Println()
	return es
}

// SUCCESS
func testSetThresholdLeastX(es Entries) (Entries, error) {
	fmt.Println("setThresholdLeastX test")
	th, err := setThresholdLeast3(es)
	if err != nil {
		fmt.Println("TestSetThresholdLeastX failed: ", err)
		return nil, err
	}
	printEntries(th)
	fmt.Println()
	return th, nil
}

// SUCCESS
func testReSortLeast(es Entries) string {
	fmt.Println("reSortLeast test")
	fmt.Println("outside range -1")
	fmt.Println("")
	copy := es
	new := newEntry("indv06", 175)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID, newEntries := reSortLeast(new, copy)
	fmt.Println("threshold: ")
	printEntries(newEntries)
	fmt.Println("	delID: ", delID)

	fmt.Println("within range -1 + 1")
	new = newEntry("indv07", 75)
	copy = es
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID, newEntries = reSortLeast(new, copy)
	fmt.Println("threshold: ")
	printEntries(newEntries)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	fmt.Println("within range -1 + 1")
	new = newEntry("indv08", 90)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID, newEntries = reSortLeast(new, newEntries)
	fmt.Println("threshold: ")
	printEntries(newEntries)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	fmt.Println("outside range -1")
	new = newEntry("indv09", 200)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID, newEntries = reSortLeast(new, newEntries)
	fmt.Println("threshold: ")
	printEntries(newEntries)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	fmt.Println("outside range -1")
	new = newEntry("indv10", 125)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID, newEntries = reSortLeast(new, newEntries)
	fmt.Println("threshold: ")
	printEntries(newEntries)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	fmt.Println("outside range -1")
	new = newEntry("indv11", 160)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID, newEntries = reSortLeast(new, newEntries)
	fmt.Println("threshold: ")
	printEntries(newEntries)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	// panics at this point - index of empty slice newEntries

	/* fmt.Println("outside range -1")
	new = newEntry("indv12", 120)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID, newEntries = reSortLeast(new, newEntries)
	fmt.Println("threshold: ")
	printEntries(newEntries)
	fmt.Println("	delID: ", delID)
	fmt.Println() */

	return delID
}

// SUCCESS
func TestCompareUnits() error {
	fmt.Println("testMap: ", TestMap)
	es := testSortTopX(TestMap)
	th, err := testSetThresholdLeastX(es)
	if err != nil {
		fmt.Println("main failed: ")
		return err
	}
	_ = testReSortLeast(th)
	return nil
}

/* Working Code */

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
// REFACTOR 6/23/20 - return var least/do not modify in place
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

// TEST ONLY
func setThresholdLeast3(es Entries) (Entries, error) {
	if len(es) < 3 {
		return nil, fmt.Errorf("setThresholdLeast5 failed: not enough elements in list")
	}
	return es[len(es)-3:], nil
}

// setThresholdLeast5 sets a threshold with the smallest 5 values in the Top x
// sorted greatest -> smallest
func setThresholdLeast5(es Entries) (Entries, error) {
	if len(es) < 5 {
		return nil, fmt.Errorf("=etThresholdLeast5 failed: not enough elements in list")
	}
	return es[len(es)-5:], nil
}

// setThresholdLeast10 sets a threshold with the smallest 10 values in the Top x
func setThresholdLeast10(es Entries) (Entries, error) {
	if len(es) < 10 {
		return nil, fmt.Errorf("setThresholdLeast10 failed: not enough elements in list")
	}

	return es[len(es)-10:], nil
}

// newEntry creats an entry struct from Top X Amt key/value pair
func newEntry(k string, v float32) *donations.Entry {
	return &donations.Entry{ID: k, Total: v}
}

// TEST ONLY
func printThreshold(th []interface{}) {
	for _, e := range th {
		fmt.Printf("\tID: %v\tTotal: %v\n", e.(*donations.Entry).ID, e.(*donations.Entry).Total)
	}
}

func printEntries(es Entries) {
	for _, e := range es {
		fmt.Printf("%s: %v\t", e.ID, e.Total)
	}
	fmt.Println()
}
