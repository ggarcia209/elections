package testDB

import (
	"fmt"
	"sort"

	"github.com/elections/donations"
)

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
		least, err = setThresholdLeast3(es) // test only
		if err != nil {
			fmt.Println("compare failed: ", err)
			return fmt.Errorf("compare failed: %v", err)
		}
		fmt.Println("new threshold: ")
		printThreshold(least)
	} else {
		for _, entry := range comp.RefThreshold {
			least = append(least, entry.(*donations.Entry))
		}
		fmt.Println("current threshold: ")
		printThreshold(least)
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
		fmt.Println("updated threshold: ")
		printThreshold(least)
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

// TEST ONLY
func printThreshold(es Entries) {
	for _, e := range es {
		fmt.Printf("\tID: %v\tTotal: %v\n", e.ID, e.Total)
	}
}

// reSortLeast re-sorts the least 5 or 10 values when a new value breaks the threshold (least[len(least)-1].Total)
// and returns the ID of the key to be deleted and the new sorted list of least values
func reSortLeast(new *donations.Entry, es *Entries) string {
	copy := *es
	fmt.Println("resortLeast: new.Total: ", new.Total)
	fmt.Println("resortLeast: copy[0].Total: ", copy[0].Total)
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
