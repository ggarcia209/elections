package testDB

import (
	"fmt"
	"sort"

	"github.com/elections/donations"
)

/* Unit Tests */
var TestMap = map[string]float32{"indv01": 100, "indv02": 200, "indv03": 50, "indv04": 0, "indv05": 150}

func testSortTopX(m map[string]float32) Entries {
	fmt.Println("sortTopX test")
	es := sortTopX(m)
	printEntries(es)
	fmt.Println()
	return es
}

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

func testReSortLeast(es Entries) string {
	fmt.Println("reSortLeast test")
	fmt.Println("outside range -1")
	copy := es
	new := newEntry("indv06", 175)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID := reSortLeast(new, &copy)
	fmt.Println("threshold: ")
	printEntries(copy)
	fmt.Println("	delID: ", delID)

	fmt.Println("within range -1 + 1")
	new = newEntry("indv07", 75)
	copy = es
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID = reSortLeast(new, &copy)
	fmt.Println("threshold: ")
	printEntries(copy)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	fmt.Println("within range -1 + 1")
	new = newEntry("indv08", 90)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID = reSortLeast(new, &copy)
	fmt.Println("threshold: ")
	printEntries(copy)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	fmt.Println("outside range -1")
	new = newEntry("indv09", 200)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID = reSortLeast(new, &copy)
	fmt.Println("threshold: ")
	printEntries(copy)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	fmt.Println("outside range -1")
	new = newEntry("indv10", 125)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID = reSortLeast(new, &copy)
	fmt.Println("threshold: ")
	printEntries(copy)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	fmt.Println("outside range -1")
	new = newEntry("indv11", 160)
	fmt.Printf("\tnew: %s: %v\n", new.ID, new.Total)
	delID = reSortLeast(new, &copy)
	fmt.Println("threshold: ")
	printEntries(copy)
	fmt.Println("	delID: ", delID)
	fmt.Println()

	return delID
}

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

func printEntries(es Entries) {
	for _, e := range es {
		fmt.Printf("%s: %v\t", e.ID, e.Total)
	}
	fmt.Println()
}

func TestCompare() error {
	// above least/empty threshold range
	fmt.Println("comp1: above least/empty threshold")
	comp1 := &comparison{
		RefID:        "cmte01",
		RefAmts:      map[string]float32{"indv01": 100, "indv02": 200, "indv03": 50, "indv04": 40, "indv05": 150},
		RefTxs:       map[string]float32{"indv01": 1, "indv02": 1, "indv03": 1, "indv04": 1, "indv05": 1},
		RefThreshold: []interface{}{},
		CompID:       "indv06",
		CompAmts:     map[string]float32{"cmte01": 300, "cmte02": 40},
		CompTxs:      map[string]float32{"cmte01": 2, "cmte02": 1},
	}
	err := compare(comp1)
	if err != nil {
		return err
	}
	printComp(comp1)

	// below least/empty threshold range
	fmt.Println("comp2: below least/empty threshold")
	comp2 := &comparison{
		RefID:        "cmte01",
		RefAmts:      map[string]float32{"indv01": 100, "indv02": 200, "indv03": 50, "indv04": 40, "indv05": 150},
		RefTxs:       map[string]float32{"indv01": 1, "indv02": 1, "indv03": 1, "indv04": 1, "indv05": 1},
		RefThreshold: []interface{}{},
		CompID:       "indv06",
		CompAmts:     map[string]float32{"cmte01": 10, "cmte02": 40},
		CompTxs:      map[string]float32{"cmte01": 1, "cmte02": 1},
	}
	err = compare(comp2)
	if err != nil {
		return err
	}
	printComp(comp2)

	// above least/non-empty threshold range
	fmt.Println("comp3: above least/non-empty threshold")
	comp3 := &comparison{
		RefID:        "cmte01",
		RefAmts:      map[string]float32{"indv01": 100, "indv02": 200, "indv03": 50, "indv04": 40, "indv05": 150},
		RefTxs:       map[string]float32{"indv01": 1, "indv02": 1, "indv03": 1, "indv04": 1, "indv05": 1},
		RefThreshold: []interface{}{&donations.Entry{"indv01", 100}, &donations.Entry{"indv03", 50}, &donations.Entry{"indv04", 0}},
		CompID:       "indv06",
		CompAmts:     map[string]float32{"cmte01": 300, "cmte02": 40},
		CompTxs:      map[string]float32{"cmte01": 2, "cmte02": 1},
	}
	err = compare(comp3)
	if err != nil {
		return err
	}
	printComp(comp3)

	// below least/non-empty threshold range
	fmt.Println("comp4: below least/non-empty threshold")
	comp4 := &comparison{
		RefID:        "cmte01",
		RefAmts:      map[string]float32{"indv01": 100, "indv02": 200, "indv03": 50, "indv04": 40, "indv05": 150},
		RefTxs:       map[string]float32{"indv01": 1, "indv02": 1, "indv03": 1, "indv04": 1, "indv05": 1},
		RefThreshold: []interface{}{&donations.Entry{"indv01", 100}, &donations.Entry{"indv03", 50}, &donations.Entry{"indv04", 0}},
		CompID:       "indv06",
		CompAmts:     map[string]float32{"cmte01": 10, "cmte02": 40},
		CompTxs:      map[string]float32{"cmte01": 1, "cmte02": 1},
	}
	err = compare(comp4)
	if err != nil {
		return err
	}
	printComp(comp4)

	return nil
}

func printComp(c *comparison) {
	fmt.Println("RefID: ", c.RefID)
	fmt.Println("RefAmts: ", c.RefAmts)
	fmt.Println("RefTxs: ", c.RefTxs)
	for _, e := range c.RefThreshold {
		fmt.Printf("\t%s: %v", e.(*donations.Entry).ID, e.(*donations.Entry).Total)
	}
	fmt.Println()
	fmt.Println("CompID: ", c.CompID)
	fmt.Println("CompAmts: ", c.CompAmts)
	fmt.Println("CompTxs: ", c.CompTxs)
	fmt.Println()
}

/* Working Code */

type comparison struct {
	RefID        string             // reference object
	RefAmts      map[string]float32 // marginal amount added to reference amount if compare amount > threshold
	RefTxs       map[string]float32
	RefThreshold []interface{}      // compare smallest amount in reference threshold list against compare amount
	CompID       string             // object being compared to reference object
	CompAmts     map[string]float32 // marginal amount included before comparison
	CompTxs      map[string]float32
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
	return del
}

// compare compares the maps set in the comparison object to the threshold
func compare(comp *comparison) error {
	var least Entries
	var err error

	// if Threshold list is exhausted
	if len(comp.RefThreshold) == 0 {
		// sort Amts map and take bottom 10 as threshold list
		es := sortTopX(comp.RefAmts)
		least, err = setThresholdLeast3(es) // test only
		if err != nil {
			fmt.Println("compare failed: ", err)
			return fmt.Errorf("compare failed: %v", err)
		}
	} else {
		for _, entry := range comp.RefThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// compare new sender's total to receiver's threshold value
	threshold := least[len(least)-1].Total // last/smallest obj in least
	fmt.Println("least: ", threshold)

	// if amount sent to receiver is > receiver's threshold
	if comp.CompAmts[comp.RefID] > threshold {
		// create new threshold entry for sender & amount contributed by sender
		new := newEntry(comp.CompID, comp.CompAmts[comp.RefID])
		// reSort threshold list w/ new entry and retreive deletion key for obj below threshold
		delID := reSortLeast(new, &least)
		// delete the records for obj below threshold
		delete(comp.RefAmts, delID)
		delete(comp.RefTxs, delID)
		// add new obj data to records
		comp.RefAmts[comp.CompID] = comp.CompAmts[comp.RefID]
		comp.RefTxs[comp.CompID] = comp.CompTxs[comp.RefID]
	} else {
		// sender/value does not qualify -- return and continue
		return nil
	}

	// update object's threshold list
	th := []interface{}{}
	for _, entry := range least {
		th = append(th, entry)
	}
	comp.RefThreshold = append(comp.RefThreshold[:0], th...)

	return nil
}

// check to see if previous total of entry is in threshold range when updating existing entry
func checkThreshold(newID string, m map[string]float32, th []interface{}) ([]interface{}, error) {
	inRange := false
	check := map[string]bool{newID: true}
	for _, e := range th {
		if check[e.(*donations.Entry).ID] == true {
			inRange = true
		}
	}
	if inRange {
		es := sortTopX(m)
		newRange, err := setThresholdLeast3(es)
		if err != nil {
			fmt.Println("checkThreshold failed: ", err)
			return []interface{}{}, fmt.Errorf("checkThreshold failed: %v", err)
		}
		// update object's threshold list
		newTh := []interface{}{}
		for _, entry := range newRange {
			newTh = append(newTh, entry)
		}
		return newTh, nil
	}
	return th, nil
}

// reSortLeast re-sorts the least 5 or 10 values when a new value breaks the threshold (least[len(least)-1].Total)
// and returns the ID of the key to be deleted and the new sorted list of least values
func reSortLeast(new *donations.Entry, es *Entries) string {
	copy := *es
	// if new.Total >= largest value in threshold list
	if new.Total >= copy[0].Total {
		// update original list of entries by overwriting it with new copy
		// es = &copy
		// pop smallest value and get it's ID to delete from records
		delID := es.popLeast().ID
		fmt.Println("resortLeast: delID: ", delID)
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
	fmt.Println("resortLeast: delID: ", delID)

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
func printThreshold(es Entries) {
	for _, e := range es {
		fmt.Printf("\tID: %v\tTotal: %v\n", e.ID, e.Total)
	}
}
