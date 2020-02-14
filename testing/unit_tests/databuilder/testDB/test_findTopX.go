package testDB

import (
	"fmt"

	"github.com/elections/donations"
)

// SUCCESS
func TestCompare() error {
	// above least/empty threshold range
	fmt.Println("map: ", map[string]float32{"indv01": 100, "indv02": 200, "indv03": 50, "indv04": 40, "indv05": 150})
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
		RefThreshold: []interface{}{&donations.Entry{"indv01", 100}, &donations.Entry{"indv03", 50}, &donations.Entry{"indv04", 40}},
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
		RefThreshold: []interface{}{&donations.Entry{"indv01", 100}, &donations.Entry{"indv03", 50}, &donations.Entry{"indv04", 40}},
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
	fmt.Printf("RefTreshold: ")
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
		// after updating object's threshold list
		th := []interface{}{}
		for _, entry := range least {
			th = append(th, entry)
		}
		comp.RefThreshold = append(comp.RefThreshold[:0], th...)
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
