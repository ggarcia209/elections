package databuilder

import (
	"fmt"

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

// check to see if previous total of entry is in threshold range when updating existing entry
func checkThreshold(newID string, newTotal float32, tr *[]interface{}) {

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
