package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/elections/donations"
)

/* TEST TYPES */

// comparison is used to compare corresponding map types (receiver/sender)
// passed to compare() function
type comparison struct {
	RecID     string
	RecAmts   map[string]float32
	RecTxs    map[string]float32
	Threshold []interface{}
	SenID     string
	SenAmts   map[string]float32
	SenTxs    map[string]float32
}

// test only -- simulated transaction
type tx struct {
	ID     string
	Amount float32
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

/* END TEST TYPES */

/* TEST OBJECTS */

var filer = donations.CmteTxData{
	CmteID:                      "cmte00",
	TopIndvContributorsAmt:      map[string]float32{"indv1": 100, "indv2": 150, "indv3": 80, "indv4": 200, "indv5": 40, "indv6": 400, "indv7": 120, "indv8": 100, "indv9": 225},
	TopIndvContributorsTxs:      map[string]float32{"indv1": 2, "indv2": 3, "indv3": 1, "indv4": 3, "indv5": 1, "indv6": 4, "indv7": 3, "indv8": 1, "indv9": 5},
	TopIndvContributorThreshold: []interface{}{},
}

var indv4 = donations.Individual{
	ID:            "indv4",
	RecipientsAmt: map[string]float32{"cmte1": 50, "cmte00": 200, "cmte2": 100},
	RecipientsTxs: map[string]float32{"cmte1": 1, "cmte00": 3, "cmte2": 2},
}

var indv8 = donations.Individual{
	ID:            "indv8",
	RecipientsAmt: map[string]float32{"cmte1": 50, "cmte00": 100, "cmte2": 100},
	RecipientsTxs: map[string]float32{"cmte1": 1, "cmte00": 1, "cmte2": 2},
}

var indv10 = donations.Individual{
	ID:            "indv10",
	RecipientsAmt: map[string]float32{"cmte1": 40, "cmte2": 200},
	RecipientsTxs: map[string]float32{"cmte1": 1, "cmte2": 2},
}

var indv11 = donations.Individual{
	ID:            "indv11",
	RecipientsAmt: map[string]float32{"cmte1": 60, "cmte2": 50},
	RecipientsTxs: map[string]float32{"cmte1": 1, "cmte2": 2},
}

var indv12 = donations.Individual{
	ID:            "indv12",
	RecipientsAmt: map[string]float32{"cmte2": 60, "cmte3": 50},
	RecipientsTxs: map[string]float32{"cmte2": 1, "cmte3": 2},
}

var indv13 = donations.Individual{
	ID:            "indv13",
	RecipientsAmt: make(map[string]float32),
	RecipientsTxs: make(map[string]float32),
}

var indv14 = donations.Individual{
	ID:            "indv14",
	RecipientsAmt: make(map[string]float32),
	RecipientsTxs: make(map[string]float32),
}

var indv15 = donations.Individual{
	ID:            "indv15",
	RecipientsAmt: make(map[string]float32),
	RecipientsTxs: make(map[string]float32),
}

var indv16 = donations.Individual{
	ID:            "indv16",
	RecipientsAmt: make(map[string]float32),
	RecipientsTxs: make(map[string]float32),
}

var tx1 = tx{
	ID:     "indv10",
	Amount: 175,
}

var tx2 = tx{
	ID:     "indv4",
	Amount: 50,
}

var tx3 = tx{
	ID:     "indv11",
	Amount: 90,
}

var tx4 = tx{
	ID:     "indv12",
	Amount: 125,
}

var tx5 = tx{
	ID:     "indv13",
	Amount: 135,
}

var tx6 = tx{
	ID:     "indv14",
	Amount: 110,
}

var tx7 = tx{
	ID:     "indv15",
	Amount: 200,
}

var tx8 = tx{
	ID:     "indv16",
	Amount: 30,
}

var tx9 = tx{
	ID:     "indv15",
	Amount: 100,
}

var tx10 = tx{
	ID:     "indv8",
	Amount: 200,
}

/* END TEST OBJECTS */

func main() {
	// values at start
	fmt.Println("***** PRE *****")
	fmt.Println("filer: ", filer.TopIndvContributorsAmt, filer.TopIndvContributorsTxs, filer.TopIndvContributorThreshold)
	fmt.Println("indv4: ", indv4.RecipientsAmt, indv4.RecipientsTxs)
	fmt.Println("indv8: ", indv8.RecipientsAmt, indv8.RecipientsTxs)
	fmt.Println("indv10: ", indv10.RecipientsAmt, indv10.RecipientsTxs)
	fmt.Println("indv11: ", indv11.RecipientsAmt, indv11.RecipientsTxs)
	fmt.Println("indv12: ", indv12.RecipientsAmt, indv12.RecipientsTxs)
	fmt.Println("indv13: ", indv13.RecipientsAmt, indv13.RecipientsTxs)
	fmt.Println("indv14: ", indv14.RecipientsAmt, indv14.RecipientsTxs)
	fmt.Println("indv15: ", indv15.RecipientsAmt, indv15.RecipientsTxs)
	fmt.Println("indv16: ", indv16.RecipientsAmt, indv16.RecipientsTxs)
	fmt.Println()

	// current max == 10
	// current len == 9
	err := updateCmte(tx1, &filer, &indv10)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// at max -- add to existing entry
	err = updateCmte(tx2, &filer, &indv4)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// at max -- replace least (indv5: 40) -- between bottom/upper thresholds (40, 80, 100)
	err = updateCmte(tx3, &filer, &indv11)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Threshold: ")
	for i, th := range filer.TopIndvContributorThreshold {
		fmt.Printf("%d) ID: %s, Total: %v\n", i, th.(*donations.Entry).ID, th.(*donations.Entry).Total)
	}

	// at max -- replace least (indv3: 80) -- above upper threshold
	err = updateCmte(tx4, &filer, &indv12)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Threshold: ")
	for i, th := range filer.TopIndvContributorThreshold {
		fmt.Printf("%d) ID: %s, Total: %v\n", i, th.(*donations.Entry).ID, th.(*donations.Entry).Total)
	}

	// at max -- replace least (indv11: 90) -- above upper threshold
	err = updateCmte(tx5, &filer, &indv13)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Threshold: ")
	for i, th := range filer.TopIndvContributorThreshold {
		fmt.Printf("%d) ID: %s, Total: %v\n", i, th.(*donations.Entry).ID, th.(*donations.Entry).Total)
	}

	// at max -- replace least (indv8: 100) -- above upper threshold
	err = updateCmte(tx6, &filer, &indv14)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Threshold: ")
	for i, th := range filer.TopIndvContributorThreshold {
		fmt.Printf("%d) ID: %s, Total: %v\n", i, th.(*donations.Entry).ID, th.(*donations.Entry).Total)
	}

	// at max -- replace least (indv8: 100) -- reset threshold
	err = updateCmte(tx7, &filer, &indv15)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Threshold: ")
	for i, th := range filer.TopIndvContributorThreshold {
		fmt.Printf("%d) ID: %s, Total: %v\n", i, th.(*donations.Entry).ID, th.(*donations.Entry).Total)
	}

	// at max -- below threshold / does not qualify
	err = updateCmte(tx8, &filer, &indv16)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Threshold: ")
	for i, th := range filer.TopIndvContributorThreshold {
		fmt.Printf("%d) ID: %s, Total: %v\n", i, th.(*donations.Entry).ID, th.(*donations.Entry).Total)
	}

	// at max -- add to existing value
	err = updateCmte(tx9, &filer, &indv15)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Threshold: ")
	for i, th := range filer.TopIndvContributorThreshold {
		fmt.Printf("%d) ID: %s, Total: %v\n", i, th.(*donations.Entry).ID, th.(*donations.Entry).Total)
	}

	// at max -- add to existing value -- add to donor previously below threshold
	err = updateCmte(tx10, &filer, &indv8)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Threshold: ")
	for i, th := range filer.TopIndvContributorThreshold {
		fmt.Printf("%d) ID: %s, Total: %v\n", i, th.(*donations.Entry).ID, th.(*donations.Entry).Total)
	}

	// values at end
	fmt.Println("***** POST *****")
	fmt.Println("filer: ", filer.TopIndvContributorsAmt, filer.TopIndvContributorsTxs, filer.TopIndvContributorThreshold)
	for i, th := range filer.TopIndvContributorThreshold {
		fmt.Printf("%d) ID: %s, Total: %v\n", i, th.(*donations.Entry).ID, th.(*donations.Entry).Total)
	}
	fmt.Println("indv4: ", indv4.RecipientsAmt, indv4.RecipientsTxs)
	fmt.Println("indv8: ", indv8.RecipientsAmt, indv8.RecipientsTxs)
	fmt.Println("indv10: ", indv10.RecipientsAmt, indv10.RecipientsTxs)
	fmt.Println("indv11: ", indv11.RecipientsAmt, indv11.RecipientsTxs)
	fmt.Println("indv12: ", indv12.RecipientsAmt, indv12.RecipientsTxs)
	fmt.Println("indv13: ", indv13.RecipientsAmt, indv13.RecipientsTxs)
	fmt.Println("indv14: ", indv14.RecipientsAmt, indv14.RecipientsTxs)
	fmt.Println("indv15: ", indv15.RecipientsAmt, indv15.RecipientsTxs)
	fmt.Println("indv16: ", indv16.RecipientsAmt, indv16.RecipientsTxs)
	fmt.Println()

}

func updateCmte(t tx, filer *donations.CmteTxData, sender *donations.Individual) error {
	// updating existing entry
	if filer.TopIndvContributorsAmt[t.ID] != 0 || len(filer.TopIndvContributorsAmt) < 10 {
		filer.TopIndvContributorsAmt[t.ID] += t.Amount
		filer.TopIndvContributorsTxs[t.ID]++
		sender.RecipientsAmt[filer.CmteID] += t.Amount
		sender.RecipientsTxs[filer.CmteID]++
		return nil
	}

	sender.RecipientsAmt[filer.CmteID] += t.Amount
	sender.RecipientsTxs[filer.CmteID]++

	// compare new entries if list is full
	comp := comparison{
		RecID:     filer.CmteID,
		RecAmts:   filer.TopIndvContributorsAmt,
		RecTxs:    filer.TopIndvContributorsTxs,
		Threshold: filer.TopIndvContributorThreshold,
		SenID:     sender.ID,
		SenAmts:   sender.RecipientsAmt,
		SenTxs:    sender.RecipientsTxs,
	}

	err := compare(&comp)
	if err != nil {
		fmt.Println("updateCmte failed: ", err)
		return fmt.Errorf("updateCmte failed: %v", err)
	}

	// update maps and threshold from comparison object
	filer.TopIndvContributorsAmt = comp.RecAmts
	filer.TopIndvContributorsTxs = comp.RecTxs
	filer.TopIndvContributorThreshold = comp.Threshold

	sender.RecipientsAmt = comp.SenAmts
	sender.RecipientsTxs = comp.SenTxs

	return nil
}

// compare compares the maps set in the comparison object to the threshold
func compare(comp *comparison) error {
	var least Entries
	var err error

	// if Threshold list is exhausted
	if len(comp.Threshold) == 0 {
		// sort Amts map and take bottom 10 as threshold list
		es := sortTopX(comp.RecAmts)
		least, err = setThresholdLeast3(es) // 3 - TEST ONLY
		if err != nil {
			fmt.Println("compare failed: ", err)
			return fmt.Errorf("compare failed: %v", err)
		}
	} else {
		for _, entry := range comp.Threshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// compare new sender's total to receiver's threshold value
	threshold := least[len(least)-1].Total // last/smallest obj in least
	fmt.Println("smallest: ", threshold)
	fmt.Println("compare value: ", comp.SenAmts[comp.RecID])

	// if amount sent to receiver is > receiver's threshold
	if comp.SenAmts[comp.RecID] > threshold {
		// create new threshold entry for sender & amount contributed by sender
		new := newEntry(comp.SenID, comp.SenAmts[comp.RecID])
		// reSort threshold list w/ new entry and retreive deletion key for obj below threshold
		delID := reSortLeast(new, &least)
		fmt.Println("delID: ", delID)
		// delete the records for obj below threshold
		delete(comp.RecAmts, delID)
		delete(comp.RecTxs, delID)
		// add new obj data to records
		comp.RecAmts[comp.SenID] = comp.SenAmts[comp.RecID]
		comp.RecTxs[comp.SenID] = comp.SenTxs[comp.RecID]
	} else {
		// sender/value does not qualify -- return and continue
		return nil
	}

	// update object's threshold list
	th := []interface{}{}
	for _, entry := range least {
		th = append(th, entry)
	}
	comp.Threshold = append(comp.Threshold[:0], th...)

	return nil
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
