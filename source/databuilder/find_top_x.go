// Package databuilder conatins operations for updating datasets in memory.
// This package is primarily used by the admin service to create the
// primary datasets from the raw input, followed by the secondary
// datasets.
// This file contains operations for creating the Committee
// specific rankings lists.
package databuilder

import (
	"fmt"

	"github.com/elections/source/donations"
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

// Update contributor maps for incoming transactions posted by filing committee.
func updateTopDonors(receiver *donations.CmteTxData, sender interface{}, cont *donations.Contribution) (comparison, error) {
	comp := comparison{}
	switch t := sender.(type) {
	case *donations.Individual:
		comp = comparison{
			RefID:        receiver.CmteID,
			RefAmts:      receiver.TopIndvContributorsAmt,
			RefTxs:       receiver.TopIndvContributorsTxs,
			RefThreshold: receiver.TopIndvContributorThreshold,
			CompID:       sender.(*donations.Individual).ID,
			CompAmts:     sender.(*donations.Individual).RecipientsAmt,
			CompTxs:      sender.(*donations.Individual).RecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println(err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Candidate:
		comp = comparison{
			RefID:        receiver.CmteID,
			RefAmts:      receiver.TopIndvContributorsAmt,
			RefTxs:       receiver.TopIndvContributorsTxs,
			RefThreshold: receiver.TopIndvContributorThreshold,
			CompID:       sender.(*donations.Candidate).ID,
			CompAmts:     sender.(*donations.Candidate).DirectRecipientsAmts,
			CompTxs:      sender.(*donations.Candidate).DirectRecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println(err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.CmteTxData:
		// Contribution || Other transfer
		comp = cmteCompGen(receiver, sender.(*donations.CmteTxData), cont)
		err := compare(&comp)
		if err != nil {
			fmt.Println(err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	default:
		_ = t
		return comparison{}, fmt.Errorf("updateTopDonors failed: wrong interface type")
	}
	return comp, nil
}

// Update contributor maps for outgoing transactions posted by filing committee.
func updateTopRecipients(sender *donations.CmteTxData, receiver interface{}) (comparison, error) {
	comp := comparison{}
	switch t := receiver.(type) {
	case *donations.Individual:
		comp = comparison{
			RefID:        sender.CmteID,
			RefAmts:      sender.TopExpRecipientsAmt,
			RefTxs:       sender.TopExpRecipientsTxs,
			RefThreshold: sender.TopExpThreshold,
			CompID:       receiver.(*donations.Individual).ID,
			CompAmts:     receiver.(*donations.Individual).SendersAmt,
			CompTxs:      receiver.(*donations.Individual).SendersTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println(err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Candidate:
		comp = comparison{
			RefID:        sender.CmteID,
			RefAmts:      sender.TopExpRecipientsAmt,
			RefTxs:       sender.TopExpRecipientsTxs,
			RefThreshold: sender.TopExpThreshold,
			CompID:       receiver.(*donations.Candidate).ID,
			CompAmts:     receiver.(*donations.Candidate).DirectSendersAmts,
			CompTxs:      receiver.(*donations.Candidate).DirectSendersTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println(err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.CmteTxData:
		// Contribution || Other transfer
		return comparison{}, fmt.Errorf("type CmteTxData invalid - debit filer's TransferRecsAmt directly for outgoing transactions to other committee")
	default:
		_ = t
		return comparison{}, fmt.Errorf("updateTopDonors failed: wrong interface type")
	}
	return comp, nil
}

func updateOpExpRecipients(sender *donations.CmteTxData, receiver *donations.Individual) (comparison, error) {
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
		fmt.Println(err)
		return comparison{}, fmt.Errorf("updateOpExpRecipients failed: %v", err)
	}
	return comp, nil
}

// Generates comparison object corresponding to incoming/outoing transaction.
func cmteCompGen(filer, other *donations.CmteTxData, cont *donations.Contribution) comparison {
	comp := comparison{
		RefID:        filer.CmteID,
		RefAmts:      filer.TopCmteOrgContributorsAmt,
		RefTxs:       filer.TopCmteOrgContributorsTxs,
		RefThreshold: filer.TopCmteOrgContributorThreshold,
		CompID:       other.CmteID,
		CompAmts:     other.TransferRecsAmt,
		CompTxs:      other.TransferRecsTxs,
	}

	// add marginal value from new transaction
	comp.CompAmts[comp.RefID] += cont.TxAmt
	comp.CompTxs[comp.RefID]++

	return comp
}

// compare compares the maps set in the comparison object to the threshold.
func compare(comp *comparison) error {
	var least Entries
	var err error

	// if Threshold list is exhausted
	if len(comp.RefThreshold) == 0 {
		// sort Amts map and take bottom 10 as threshold list
		es := sortTopX(comp.RefAmts)
		least, err = setThresholdLeast10(es)
		if err != nil {
			fmt.Println(err)
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
		delID, newEntries := reSortLeast(new, least)
		least = newEntries
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

// Check to see if previous total of entry is in threshold range when updating existing entry.
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
		newRange, err := setThresholdLeast10(es)
		if err != nil {
			fmt.Println(err)
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
