package databuilder

import (
	"fmt"

	"github.com/elections/donations"
)

type comparison struct {
	RecID     string
	RecAmts   map[string]float32
	RecTxs    map[string]float32
	Threshold []interface{}
	SenID     string
	SenAmts   map[string]float32
	SenTxs    map[string]float32
}

// update contributor maps for incoming transactions posted by filing committee
func updateTopDonors(receiver *donations.CmteTxData, sender interface{}, tx *donations.Contribution, transfer bool) (comparison, error) {
	comp := comparison{}
	switch t := sender.(type) {
	case *donations.Individual:
		comp = comparison{
			RecID:     receiver.CmteID,
			RecAmts:   receiver.TopIndvContributorsAmt,
			RecTxs:    receiver.TopIndvContributorsTxs,
			Threshold: receiver.TopIndvContributorThreshold,
			SenID:     t.ID,
			SenAmts:   sender.(*donations.Individual).RecipientsAmt,
			SenTxs:    sender.(*donations.Individual).RecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Candidate:
		comp = comparison{
			RecID:     receiver.CmteID,
			RecAmts:   receiver.TopIndvContributorsAmt,
			RecTxs:    receiver.TopIndvContributorsTxs,
			Threshold: receiver.TopIndvContributorThreshold,
			SenID:     t.ID,
			SenAmts:   sender.(*donations.Candidate).DirectRecipientsAmts,
			SenTxs:    sender.(*donations.Candidate).DirectRecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Organization:
		comp = comparison{
			RecID:     receiver.CmteID,
			RecAmts:   receiver.TopCmteOrgContributorsAmt,
			RecTxs:    receiver.TopCmteOrgContributorsTxs,
			Threshold: receiver.TopCmteOrgContributorThreshold,
			SenID:     t.ID,
			SenAmts:   sender.(*donations.Organization).RecipientsAmt,
			SenTxs:    sender.(*donations.Organization).RecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.CmteTxData:
		// Contribution || Other transfer
		comp = cmteCompGen(receiver, sender.(*donations.CmteTxData), tx.TxType, transfer)
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
func updateTopRecipients(sender *donations.CmteTxData, receiver interface{}, tx *donations.Contribution, transfer bool) (comparison, error) {
	comp := comparison{}
	switch t := receiver.(type) {
	case *donations.Individual:
		comp = comparison{
			RecID:     t.ID,
			RecAmts:   receiver.(*donations.Individual).SendersAmt,
			RecTxs:    receiver.(*donations.Individual).SendersTxs,
			Threshold: sender.TopExpThreshold,
			SenID:     sender.CmteID,
			SenAmts:   sender.TopExpRecipientsAmt,
			SenTxs:    sender.TopExpRecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Candidate:
		comp = comparison{
			RecID:     t.ID,
			RecAmts:   receiver.(*donations.Candidate).DirectSendersAmts,
			RecTxs:    receiver.(*donations.Candidate).DirectSendersTxs,
			Threshold: sender.TopExpThreshold,
			SenID:     sender.CmteID,
			SenAmts:   sender.TopExpRecipientsAmt,
			SenTxs:    sender.TopExpRecipientsTxs,
		}
		if transfer {
			comp.SenAmts = sender.TransferRecsAmt
			comp.SenTxs = sender.TransferRecsTxs
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.Organization:
		comp = comparison{
			RecID:     t.ID,
			RecAmts:   receiver.(*donations.Organization).SendersAmt,
			RecTxs:    receiver.(*donations.Organization).SendersTxs,
			Threshold: sender.TopExpThreshold,
			SenID:     sender.CmteID,
			SenAmts:   sender.TopExpRecipientsAmt,
			SenTxs:    sender.TopExpRecipientsTxs,
		}
		err := compare(&comp)
		if err != nil {
			fmt.Println("updateTopDonors failed: ", err)
			return comparison{}, fmt.Errorf("updateTopDonors failed: %v", err)
		}
	case *donations.CmteTxData:
		// Contribution || Other transfer
		comp = cmteCompGen(receiver.(*donations.CmteTxData), sender, tx.TxType, transfer)
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

// generates comparison object corresponding to incoming/outoing transaction
func cmteCompGen(filer, other *donations.CmteTxData, txType string, transfer bool) comparison {
	// determine tx type (incoming / outgoing / memo)
	var comp comparison
	if txType < "20" || (txType >= "30" && txType < "33") {
		// incoming
		comp = comparison{
			// values determined by tx type
			RecID:     filer.CmteID,
			RecAmts:   filer.TopCmteOrgContributorsAmt,
			RecTxs:    filer.TopCmteOrgContributorsTxs,
			Threshold: filer.TopCmteOrgContributorThreshold,
			SenID:     other.CmteID,
			SenAmts:   other.TopExpRecipientsAmt,
			SenTxs:    other.TopExpRecipientsTxs,
		}
		if transfer {
			comp.SenAmts = other.TransferRecsAmt
			comp.SenTxs = other.TransferRecsTxs
		}
	} else {
		// outgoing
		comp = comparison{
			// values determined by tx type
			RecID:     other.CmteID,
			RecAmts:   other.TopCmteOrgContributorsAmt,
			RecTxs:    other.TopCmteOrgContributorsTxs,
			Threshold: other.TopCmteOrgContributorThreshold,
			SenID:     filer.CmteID,
			SenAmts:   filer.TopExpRecipientsAmt,
			SenTxs:    filer.TopExpRecipientsTxs,
		}
		if transfer {
			comp.SenAmts = filer.TransferRecsAmt
			comp.SenTxs = filer.TransferRecsTxs
		}
	}
	return comp
}

// compare compares the maps set in the comparison object to the threshold
func compare(comp *comparison) error {
	var least Entries
	var err error

	// if Threshold list is exhausted
	if len(comp.Threshold) == 0 {
		// sort Amts map and take bottom 10 as threshold list
		es := sortTopX(comp.RecAmts)
		least, err = setThresholdLeast10(es)
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

	// if amount sent to receiver is > receiver's threshold
	if comp.SenAmts[comp.RecID] > threshold {
		// create new threshold entry for sender & amount contributed by sender
		new := newEntry(comp.SenID, comp.SenAmts[comp.RecID])
		// reSort threshold list w/ new entry and retreive deletion key for obj below threshold
		delID := reSortLeast(new, &least)
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

// DEPRECATED

/*
func updateTopCmteTotals(sen, rec *donations.Committee) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(rec.TopCmteDonorThreshold) == 0 {
		es := sortTopX(rec.TopCmteDonorsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	} else {
		for _, entry := range rec.TopCmteDonorThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// compare new sender's total to receiver's threshold value
	threshold := least[len(least)-1].Total // last/smallest obj in least
	if sen.AffiliatesAmt[rec.ID] > threshold {
		new := newEntry(sen.ID, sen.AffiliatesAmt[rec.ID])
		delID := reSortLeast(new, &least)
		delete(rec.TopCmteDonorsAmt, delID)
		delete(rec.TopCmteDonorsTxs, delID)
		rec.TopCmteDonorsAmt[sen.ID] = sen.AffiliatesAmt[rec.ID]
		rec.TopCmteDonorsTxs[sen.ID] = sen.AffiliatesTxs[rec.ID]
	}

	return nil
}

func updateTopIndvTotals(sen *donations.Individual, rec *donations.Committee) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(rec.TopIndvDonorThreshold) == 0 {
		es := sortTopX(rec.TopIndvDonorsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	} else {
		for _, entry := range rec.TopIndvDonorThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// compare new sender's total to receiver's threshold value
	threshold := least[len(least)-1].Total // last/smallest obj in least
	if sen.RecipientsAmt[rec.ID] > threshold {
		new := newEntry(sen.ID, sen.RecipientsAmt[rec.ID])
		delID := reSortLeast(new, &least)
		delete(rec.TopIndvDonorsAmt, delID)
		delete(rec.TopIndvDonorsTxs, delID)
		rec.TopIndvDonorsAmt[sen.ID] = sen.RecipientsAmt[rec.ID]
		rec.TopIndvDonorsTxs[sen.ID] = sen.RecipientsTxs[rec.ID]
	}

	return nil
}

func updateTopDisbRecTotals(sen *donations.Committee, rec *donations.DisbRecipient) error {
	var least Entries
	var err error

	if len(sen.TopRecThreshold) == 0 {
		es := sortTopX(sen.TopDisbRecipientsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopDisbRecTotals failed: ", err)
			return fmt.Errorf("updateTopDisbRecTotals failed: %v", err)
		}
	} else {
		for _, entry := range sen.TopRecThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// compare new sender's total to receiver's threshold value
	threshold := least[len(least)-1].Total // last/smallest obj in least
	if rec.SendersAmt[sen.ID] > threshold {
		new := newEntry(rec.ID, rec.SendersAmt[sen.ID])
		delID := reSortLeast(new, &least)
		delete(sen.TopDisbRecipientsAmt, delID)
		delete(sen.TopDisbRecipientsTxs, delID)
		sen.TopDisbRecipientsAmt[rec.ID] = rec.SendersAmt[sen.ID]
		sen.TopDisbRecipientsTxs[rec.ID] = rec.SendersTxs[sen.ID]
	}

	return nil
}

func updateCandTopCmteTotals(sen *donations.Committee, cand *donations.Candidate) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(cand.TopCDThreshold) == 0 {
		es := sortTopX(cand.TopCmteDonorsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	} else {
		for _, entry := range cand.TopCDThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// compare new sender's total to receiver's threshold value
	threshold := least[len(least)-1].Total // last/smallest obj in least
	if sen.AffiliatesAmt[cand.ID] > threshold {
		new := newEntry(sen.ID, sen.AffiliatesAmt[cand.ID])
		delID := reSortLeast(new, &least)
		delete(cand.TopCmteDonorsAmt, delID)
		delete(cand.TopCmteDonorsTxs, delID)
		cand.TopCmteDonorsAmt[sen.ID] = sen.AffiliatesAmt[cand.ID]
		cand.TopCmteDonorsTxs[sen.ID] = sen.AffiliatesTxs[cand.ID]
	}

	return nil
}

func updateCandTopIndvTotals(sen *donations.Individual, cand *donations.Candidate) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(cand.TopIDThreshold) == 0 {
		es := sortTopX(cand.TopIndvDonorsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	} else {
		for _, entry := range cand.TopIDThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// compare new sender's total to receiver's threshold value
	threshold := least[len(least)-1].Total // last/smallest obj in least
	if sen.RecipientsAmt[cand.ID] > threshold {
		new := newEntry(sen.ID, sen.RecipientsAmt[cand.ID])
		delID := reSortLeast(new, &least)
		delete(cand.TopIndvDonorsAmt, delID)
		delete(cand.TopIndvDonorsTxs, delID)
		cand.TopIndvDonorsAmt[sen.ID] = sen.RecipientsAmt[cand.ID]
		cand.TopIndvDonorsTxs[sen.ID] = sen.RecipientsTxs[cand.ID]
	}

	return nil
}

func updateCandTopDisbRecTotals(sen *donations.Candidate, rec *donations.DisbRecipient) error {
	var least Entries
	var err error

	if len(sen.TopDRThreshold) == 0 {
		es := sortTopX(sen.TopDisbRecsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopDisbRecTotals failed: ", err)
			return fmt.Errorf("updateTopDisbRecTotals failed: %v", err)
		}
	} else {
		for _, entry := range sen.TopDRThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// compare new sender's total to receiver's threshold value
	threshold := least[len(least)-1].Total // last/smallest obj in least
	if rec.SendersAmt[sen.ID] > threshold {
		new := newEntry(rec.ID, rec.SendersAmt[sen.ID])
		delID := reSortLeast(new, &least)
		delete(sen.TopDisbRecsAmt, delID)
		delete(sen.TopDisbRecsTxs, delID)
		sen.TopDisbRecsAmt[rec.ID] = rec.SendersAmt[sen.ID]
		sen.TopDisbRecsTxs[rec.ID] = rec.SendersTxs[sen.ID]
	}

	return nil
}
*/
