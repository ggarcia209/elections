package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/elections/donations"
)

/* TEST TYPES */

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
	// TopIndvContributorsAmt:      map[string]float32{"indv1": 100, "indv2": 150, "indv3": 80, "indv4": 200, "indv5": 40, "indv6": 400, "indv7": 120, "indv8": 100, "indv9": 225},
	// TopIndvContributorsTxs:      map[string]float32{"indv1": 2, "indv2": 3, "indv3": 1, "indv4": 3, "indv5": 1, "indv6": 4, "indv7": 3, "indv8": 1, "indv9": 5},
	// TopIndvContributorThreshold: []interface{}{},
}

var cmte01 = donations.CmteTxData{
	CmteID: "cmte01",
}

var cand01 = donations.Candidate{
	ID: "PCand01",
}

var cand02 = donations.Candidate{
	ID: "SCand02",
}

var cand03 = donations.Candidate{
	ID: "HCand03",
}

var org01 = donations.Organization{
	ID: "org01",
}

var org02 = donations.Organization{
	ID: "org02",
}

var indv4 = donations.Individual{
	ID:            "indv4",
	// RecipientsAmt: map[string]float32{"cmte1": 50, "cmte00": 200, "cmte2": 100},
	// RecipientsTxs: map[string]float32{"cmte1": 1, "cmte00": 3, "cmte2": 2},
}

var indv8 = donations.Individual{
	ID:            "indv8",
	// RecipientsAmt: map[string]float32{"cmte1": 50, "cmte00": 100, "cmte2": 100},
	// RecipientsTxs: map[string]float32{"cmte1": 1, "cmte00": 1, "cmte2": 2},
}

var indv10 = donations.Individual{
	ID:            "indv10",
	// RecipientsAmt: map[string]float32{"cmte1": 40, "cmte2": 200},
	// RecipientsTxs: map[string]float32{"cmte1": 1, "cmte2": 2},
}

var indv11 = donations.Individual{
	ID:            "indv11",
	// RecipientsAmt: map[string]float32{"cmte1": 60, "cmte2": 50},
	// RecipientsTxs: map[string]float32{"cmte1": 1, "cmte2": 2},
}

var indv12 = donations.Individual{
	ID:            "indv12",
	// RecipientsAmt: map[string]float32{"cmte2": 60, "cmte3": 50},
	// RecipientsTxs: map[string]float32{"cmte2": 1, "cmte3": 2},
}

var indv13 = donations.Individual{
	ID:            "indv13",
	// RecipientsAmt: make(map[string]float32),
	// RecipientsTxs: make(map[string]float32),
}

var indv14 = donations.Individual{
	ID:            "indv14",
	// RecipientsAmt: make(map[string]float32),
	// RecipientsTxs: make(map[string]float32),
}

var indv15 = donations.Individual{
	ID:            "indv15",
	// RecipientsAmt: make(map[string]float32),
	// RecipientsTxs: make(map[string]float32),
}

var indv16 = donations.Individual{
	ID:            "indv16",
	// RecipientsAmt: make(map[string]float32),
	// RecipientsTxs: make(map[string]float32),
}

var tx1 = donations.Contribution{
	CmteID: "cmte00"
	OtherID:     "indv10",
	TxAmt: 175,
	TxType: "",
	MemoCode: "",
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

// derive transaction type from contribution data
func deriveTxTypes(cont *donations.Contribution) (string, bool, bool, bool) {
	// initialize return values
	incoming := false
	memo := false
	transfer := false               // indicates transfer (true) vs. expense (false)
	transferMap := map[string]bool{ // transfer tx codes
		"15Z": true,
		"16R": true,
		"18G": true,
		"18J": true,
		"18K": true,
		"19J": true,
		"22H": true,
		"24G": true,
		"24H": true,
		"24K": true,
		"24U": true,
		"24Z": true,
		"24I": true, // verify
		"24T": true, // verify
		"30K": true,
		"30G": true,
		"30F": true,
		"31K": true,
		"31G": true,
		"31F": true,
		"32K": true,
		"32G": true,
		"32F": true,
	}
	var bucket string

	// determine tx type (incoming / outgoing / transfer, memo)
	numCode := cont.TxType
	if numCode < "20" || (numCode >= "30" && numCode < "33") {
		incoming = true
	}

	// determine if transfer or expense
	if transferMap[cont.TxType] == true {
		transfer = true
	}

	if cont.MemoCode == "X" {
		memo = true
	}

	// determine contributor/receiver type - derive from OtherID
	IDss := strings.Split(cont.OtherID, "")
	idCode := IDss[0]
	switch {
	case idCode == "C":
		bucket = "cmte_tx_data"
	case idCode == "H" || idCode == "S" || idCode == "P":
		bucket = "candidates"
	default:
		if cont.Occupation == "" {
			bucket = "organizations"
		} else {
			bucket = "individuals"
		}
	}
	return bucket, incoming, transfer, memo
}

type comparison struct {
	RefID        string             // reference object
	RefAmts      map[string]float32 // marginal amount added to reference amount if compare amount > threshold
	RefTxs       map[string]float32
	RefThreshold []interface{}      // compare smallest amount in reference threshold list against compare amount
	CompID       string             // object being compared to reference object
	CompAmts     map[string]float32 // marginal amount included before comparison
	CompTxs      map[string]float32
}

// update contributor maps for incoming transactions posted by filing committee
func updateTopDonors(receiver *donations.CmteTxData, sender interface{}, cont *donations.Contribution, transfer bool) (comparison, error) {
	comp := comparison{}
	switch t := sender.(type) {
	case *donations.Individual:
		comp = comparison{
			RefID:        receiver.CmteID,
			RefAmts:      receiver.TopIndvContributorsAmt,
			RefTxs:       receiver.TopIndvContributorsTxs,
			RefThreshold: receiver.TopIndvContributorThreshold,
			CompID:       t.ID,
			CompAmts:     sender.(*donations.Individual).RecipientsAmt,
			CompTxs:      sender.(*donations.Individual).RecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Candidate:
		comp = comparison{
			RefID:        receiver.CmteID,
			RefAmts:      receiver.TopIndvContributorsAmt,
			RefTxs:       receiver.TopIndvContributorsTxs,
			RefThreshold: receiver.TopIndvContributorThreshold,
			CompID:       t.ID,
			CompAmts:     sender.(*donations.Candidate).DirectRecipientsAmts,
			CompTxs:      sender.(*donations.Candidate).DirectRecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Organization:
		comp = comparison{
			RefID:        receiver.CmteID,
			RefAmts:      receiver.TopCmteOrgContributorsAmt,
			RefTxs:       receiver.TopCmteOrgContributorsTxs,
			RefThreshold: receiver.TopCmteOrgContributorThreshold,
			CompID:       t.ID,
			CompAmts:     sender.(*donations.Organization).RecipientsAmt,
			CompTxs:      sender.(*donations.Organization).RecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.CmteTxData:
		// Contribution || Other transfer
		comp = cmteCompGen(receiver, sender.(*donations.CmteTxData), cont, transfer)
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	default:
		return comparison{}, fmt.Errorf("updateTopDonors failed: wrong interface type")
	}
	return comp, nil
}

// update contributor maps for outgoing transactions posted by filing committee
func updateTopRecipients(sender *donations.CmteTxData, receiver interface{}, cont *donations.Contribution, transfer bool) (comparison, error) {
	comp := comparison{}
	switch t := receiver.(type) {
	case *donations.Individual:
		comp = comparison{
			RefID:        sender.CmteID,
			RefAmts:      sender.TopExpRecipientsAmt,
			RefTxs:       sender.TopExpRecipientsTxs,
			RefThreshold: sender.TopExpThreshold,
			CompID:       t.ID,
			CompAmts:     receiver.(*donations.Individual).SendersAmt,
			CompTxs:      receiver.(*donations.Individual).SendersTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Candidate:
		comp = comparison{
			RefID:        sender.CmteID,
			RefAmts:      sender.TopExpRecipientsAmt,
			RefTxs:       sender.TopExpRecipientsTxs,
			RefThreshold: sender.TopExpThreshold,
			CompID:       t.ID,
			CompAmts:     receiver.(*donations.Candidate).DirectSendersAmts,
			CompTxs:      receiver.(*donations.Candidate).DirectSendersTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Organization:
		comp = comparison{
			RefID:        sender.CmteID,
			RefAmts:      sender.TopExpRecipientsAmt,
			RefTxs:       sender.TopExpRecipientsTxs,
			RefThreshold: sender.TopExpThreshold,
			CompID:       t.ID,
			CompAmts:     receiver.(*donations.Organization).SendersAmt,
			CompTxs:      receiver.(*donations.Organization).SendersTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.CmteTxData:
		// Contribution || Other transfer
		comp = cmteCompGen(receiver.(*donations.CmteTxData), sender, cont, transfer)
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	default:
		return comparison{}, fmt.Errorf("updateTopDonors failed: wrong interface type")
	}
	return comp, nil
}

func updateOpExpRecipients(sender *donations.CmteTxData, receiver *donations.Organization) (comparison, error) {
	comp := comparison{
		RefID:        sender.CmteID,
		RefAmts:      sender.TopExpRecipientsAmt,
		RefTxs:       sender.TopExpRecipientsTxs,
		RefThreshold: sender.TopExpThreshold,
		CompID:       receiver.ID,
		CompAmts:     receiver.SendersAmt,
		CompTxs:      receiver.SendersTxs,
	}
	err := compare(&comp)
	if err != nil {
		fmt.Println("updateOpExpRecipients failed: ", err)
		return comparison{}, fmt.Errorf("updateOpExpRecipients failed: %v", err)
	}
	return comp, nil
}

// generates comparison object corresponding to incoming/outoing transaction
func cmteCompGen(filer, other *donations.CmteTxData, cont *donations.Contribution, transfer bool) comparison {
	// determine tx type (incoming / outgoing / memo)
	var comp comparison
	if cont.TxType < "20" || (cont.TxType >= "30" && cont.TxType < "33") {
		// incoming
		comp = comparison{
			// values determined by tx type
			RefID:        filer.CmteID,
			RefAmts:      filer.TopCmteOrgContributorsAmt,
			RefTxs:       filer.TopCmteOrgContributorsTxs,
			RefThreshold: filer.TopCmteOrgContributorThreshold,
			CompID:       other.CmteID,
			CompAmts:     other.TopExpRecipientsAmt,
			CompTxs:      other.TopExpRecipientsTxs,
		}
		if transfer {
			// marginal value added before call to updateTop functions
			comp.CompAmts = other.TransferRecsAmt
			comp.CompTxs = other.TransferRecsTxs
		}
		// add marginal value from new transaction
		comp.CompAmts[comp.RefID] += cont.TxAmt
		comp.CompTxs[comp.RefID]++
	} else {
		// outgoing
		comp = comparison{
			// values determined by tx type
			RefID:        filer.CmteID,
			RefAmts:      filer.TopExpRecipientsAmt,
			RefTxs:       filer.TopExpRecipientsTxs,
			RefThreshold: filer.TopExpThreshold,
			CompID:       other.CmteID,
			CompAmts:     other.TopCmteOrgContributorsAmt,
			CompTxs:      other.TopCmteOrgContributorsTxs,
		}

		if transfer {
			comp.CompAmts = other.TransferRecsAmt
			comp.CompTxs = other.TransferRecsTxs
		}
		// add marginal value from new transaction
		comp.CompAmts[comp.RefID] += cont.TxAmt
		comp.CompTxs[comp.RefID]++
	}
	return comp
}

// compare compares the maps set in the comparison object to the threshold
func compare(comp *comparison) error {
	var least Entries
	var err error

	// if Threshold list is exhausted
	if len(comp.RefThreshold) == 0 {
		// sort Amts map and take bottom 10 as threshold list
		es := sortTopX(comp.RefAmts)
		least, err = setThresholdLeast10(es)
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
