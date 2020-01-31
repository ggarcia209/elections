package databuilder

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/persist"
)

// MergeData merges multi-year data sets into one interface object
func MergeData(years []string, ID, bucket string) (interface{}, error) {
	var merged interface{}
	switch {
	case bucket == "individuals":
		set, err := createMergeSet(years, bucket, ID)
		if err != nil {
			fmt.Println("MergeData failed: ", err)
			return nil, fmt.Errorf("MergeData failed: %v", err)
		}

		merged = set[years[0]]
		for year, obj := range set {
			if year == years[0] {
				continue
			}
			indvTotalsMerge(merged.(*donations.Individual), obj.(*donations.Individual))
			indvMapMerge(merged.(*donations.Individual), obj.(*donations.Individual))
		}
	case bucket == "organizations":
		set, err := createMergeSet(years, bucket, ID)
		if err != nil {
			fmt.Println("MergeData failed: ", err)
			return nil, fmt.Errorf("MergeData failed: %v", err)
		}

		merged = set[years[0]]
		for year, obj := range set {
			if year == years[0] {
				continue
			}
			orgTotalsMerge(merged.(*donations.Organization), obj.(*donations.Organization))
			orgMapMerge(merged.(*donations.Organization), obj.(*donations.Organization))
		}
	case bucket == "cmte_tx_data":
		set, err := createMergeSet(years, bucket, ID)
		if err != nil {
			fmt.Println("MergeData failed: ", err)
			return nil, fmt.Errorf("MergeData failed: %v", err)
		}

		merged = set[years[0]]
		for year, obj := range set {
			if year == years[0] {
				continue
			}
			cmteTxTotalsMerge(merged.(*donations.CmteTxData), obj.(*donations.CmteTxData))
			cmteTxMapMerge(merged.(*donations.CmteTxData), obj.(*donations.CmteTxData))
		}
	case bucket == "candidates":
		set, err := createMergeSet(years, bucket, ID)
		if err != nil {
			fmt.Println("MergeData failed: ", err)
			return nil, fmt.Errorf("MergeData failed: %v", err)
		}

		merged = set[years[0]]
		for year, obj := range set {
			if year == years[0] {
				continue
			}
			candTotalsMerge(merged.(*donations.Candidate), obj.(*donations.Candidate))
			candMapMerge(merged.(*donations.Candidate), obj.(*donations.Candidate))
		}
	default:
		return nil, fmt.Errorf("MergeData failed: invalid bucket type")
	}

	return merged, nil
}

func createMergeSet(years []string, bucket, ID string) (map[string]interface{}, error) {
	set := make(map[string]interface{})
	for _, year := range years {
		obj, err := persist.GetObject(year, bucket, ID)
		if err != nil {
			fmt.Println("createMergeSet failed: ", err)
			return nil, fmt.Errorf("createMergeSet failed: %v", err)
		}
		set[year] = obj
	}
	return set, nil
}

func indvTotalsMerge(merge, indv *donations.Individual) {
	merge.Transactions = append(merge.Transactions, indv.Transactions...)
	merge.TotalOutAmt += indv.TotalOutAmt
	merge.TotalOutTxs += indv.TotalOutTxs
	merge.AvgTxOut = merge.TotalOutAmt / merge.TotalOutTxs
	merge.TotalInAmt += indv.TotalInAmt
	merge.TotalInTxs += indv.TotalInTxs
	merge.AvgTxIn = merge.TotalInAmt / merge.TotalInTxs
	merge.NetBalance = merge.TotalInAmt - merge.TotalOutAmt
}

func indvMapMerge(merge, indv *donations.Individual) {
	for k, v := range indv.RecipientsAmt {
		merge.RecipientsAmt[k] += v
		merge.RecipientsTxs[k] += indv.RecipientsTxs[k]
	}
	for k, v := range indv.SendersAmt {
		merge.SendersAmt[k] += v
		merge.SendersTxs[k] += indv.SendersTxs[k]
	}
}

func orgTotalsMerge(merge, org *donations.Organization) {
	merge.Transactions = append(merge.Transactions, org.Transactions...)
	merge.TotalOutAmt += org.TotalOutAmt
	merge.TotalOutTxs += org.TotalOutTxs
	merge.AvgTxOut = merge.TotalOutAmt / merge.TotalOutTxs
	merge.TotalInAmt += org.TotalInAmt
	merge.TotalInTxs += org.TotalInTxs
	merge.AvgTxIn = merge.TotalInAmt / merge.TotalInTxs
	merge.NetBalance = merge.TotalInAmt - merge.TotalOutAmt
}

func orgMapMerge(merge, org *donations.Organization) {
	for k, v := range org.RecipientsAmt {
		merge.RecipientsAmt[k] += v
		merge.RecipientsTxs[k] += org.RecipientsTxs[k]
	}
	for k, v := range org.SendersAmt {
		merge.SendersAmt[k] += v
		merge.SendersTxs[k] += org.SendersTxs[k]
	}
}

func cmteTxTotalsMerge(merge, cmte *donations.CmteTxData) {
	merge.ContributionsInAmt += cmte.ContributionsInAmt
	merge.ContributionsInTxs += cmte.ContributionsInTxs
	merge.AvgContributionIn = merge.ContributionsInAmt / merge.ContributionsInTxs
	merge.OtherReceiptsInAmt += cmte.OtherReceiptsInAmt
	merge.OtherReceiptsInTxs += cmte.OtherReceiptsInTxs
	merge.AvgOtherIn = merge.OtherReceiptsInAmt / merge.OtherReceiptsInTxs
	merge.TotalIncomingAmt = merge.ContributionsInAmt + merge.OtherReceiptsInAmt
	merge.TotalIncomingTxs = merge.ContributionsInTxs + merge.OtherReceiptsInTxs
	merge.AvgIncoming = merge.TotalIncomingAmt / merge.TotalIncomingTxs

	merge.TransfersAmt += cmte.TransfersAmt
	merge.TransfersTxs += cmte.TransfersTxs
	merge.AvgTransfer = merge.TransfersAmt / merge.TransfersTxs
	merge.ExpendituresAmt += cmte.ExpendituresAmt
	merge.ExpendituresTxs += cmte.ExpendituresTxs
	merge.AvgExpenditure = merge.ExpendituresAmt / merge.ExpendituresTxs
	merge.TotalOutgoingAmt = merge.TransfersAmt + merge.ExpendituresAmt
	merge.TotalOutgoingTxs = merge.TransfersTxs + merge.ExpendituresTxs
	merge.AvgOutgoing = merge.TotalOutgoingAmt / merge.TotalOutgoingTxs

	merge.NetBalance = merge.TotalIncomingAmt - merge.TotalOutgoingAmt
}

func cmteTxMapMerge(merge, cmte *donations.CmteTxData) {
	// Top Individual Contribtors
	for k, v := range cmte.TopIndvContributorsAmt {
		if len(merge.TopIndvContributorsAmt) < 1000 {
			// add directly to map
			if merge.TopIndvContributorsAmt[k] != 0 {
				// add to existing amounts
				merge.TopIndvContributorsAmt[k] += v
				merge.TopIndvContributorsTxs[k] += cmte.TopIndvContributorsTxs[k]
			} else {
				// create new entry in map
				merge.TopIndvContributorsAmt[k] = v
				merge.TopIndvContributorsTxs[k] = cmte.TopIndvContributorsTxs[k]
			}
		} else {
			// check values against threshold
			mergeTopTotals(merge.TopIndvContributorsAmt, merge.TopIndvContributorsTxs, cmte.TopIndvContributorsAmt, cmte.TopIndvContributorsTxs, &merge.TopIndvContributorThreshold)
		}
	}

	// Top Committee/Organization Contributors
	for k, v := range cmte.TopCmteOrgContributorsAmt {
		if len(merge.TopCmteOrgContributorsAmt) < 1000 {
			// add directly to map
			if merge.TopCmteOrgContributorsAmt[k] != 0 {
				// add to existing amounts
				merge.TopCmteOrgContributorsAmt[k] += v
				merge.TopCmteOrgContributorsTxs[k] += cmte.TopCmteOrgContributorsTxs[k]
			} else {
				// create new entry in map
				merge.TopCmteOrgContributorsAmt[k] = v
				merge.TopCmteOrgContributorsTxs[k] = cmte.TopCmteOrgContributorsTxs[k]
			}
		} else {
			// check values against threshold
			mergeTopTotals(merge.TopCmteOrgContributorsAmt, merge.TopCmteOrgContributorsTxs, cmte.TopCmteOrgContributorsAmt, cmte.TopCmteOrgContributorsTxs, &merge.TopCmteOrgContributorThreshold)
		}
	}

	// Transfers Recipients
	for k, v := range cmte.TransferRecsAmt {
		merge.TransferRecsAmt[k] += v
		merge.TransferRecsTxs[k] += cmte.TransferRecsTxs[k]
	}

	// Top Expenditure Recipients
	for k, v := range cmte.TopExpRecipientsAmt {
		if len(merge.TopExpRecipientsAmt) < 1000 {
			// add directly to map
			if merge.TopExpRecipientsAmt[k] != 0 {
				// add to existing amounts
				merge.TopExpRecipientsAmt[k] += v
				merge.TopExpRecipientsTxs[k] += cmte.TopExpRecipientsTxs[k]
			} else {
				// create new entry in map
				merge.TopExpRecipientsAmt[k] = v
				merge.TopExpRecipientsTxs[k] = cmte.TopExpRecipientsTxs[k]
			}
		} else {
			// check values against threshold
			mergeTopTotals(merge.TopExpRecipientsAmt, merge.TopExpRecipientsTxs, cmte.TopExpRecipientsAmt, cmte.TopExpRecipientsTxs, &merge.TopExpThreshold)
		}
	}
}

func mergeTopTotals(mergeAmts, mergeTxs, mAmts, mTxs map[string]float32, mergeTh *[]interface{}) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(*mergeTh) == 0 {
		es := sortTopX(mergeAmts)
		least, err = setThresholdLeast10(es)
		if err != nil {
			fmt.Println("mergeTopTotals failed: ", err)
			return fmt.Errorf("mergeTopTotals failed: %v", err)
		}
	} else {
		for _, entry := range merge {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// merge TopIndvDonor maps
	threshold := least[len(least)-1].Total // last/smallest obj in least
	for k, v := range mAmts {
		// update existing entrie's totals
		if mergeAmts[k] != 0 {
			mergeAmts[k] += v
			mergeTxs[k] += mTxs[k]
			continue
		}

		if mergeAmts[k] == 0 && v > threshold {
			new := newEntry(k, v)
			delID := reSortLeast(new, &least)
			delete(mergeAmts, delID)
			delete(mergeTxs, delID)
			mergeAmts[k] = v
			mergeTxs[k] = mTxs[k]
		}
	}

	// update object's threshold list
	th := []interface{}{}
	for _, entry := range least {
		th = append(th, entry)
	}
	mergeTh = append(*mergeTh[:0], th...)

	return nil
}

func candTotalsMerge(merge, cand *donations.Candidate) {
	merge.OtherAffiliates = append(merge.OtherAffiliates, cand.OtherAffiliates...)
	merge.TotalDirectInAmt += cand.TotalDirectInAmt
	merge.TotalDirectInTxs += cand.TotalDirectInTxs
	merge.AvgDirectIn = merge.TotalDirectInAmt / merge.TotalDirectInTxs
	merge.TotalDirectOutAmt += cand.TotalDirectOutAmt
	merge.TotalDirectOutTxs += cand.TotalDirectOutTxs
	merge.AvgDirectOut = merge.TotalDirectOutAmt / merge.TotalDirectOutTxs
	merge.NetBalanceDirectTx = merge.TotalDirectInAmt - merge.TotalDirectOutAmt
}

func candMapMerge(merge, cand *donations.Candidate) {
	for k, v := range cand.DirectRecipientsAmt {
		merge.DirectRecipientsAmt[k] += v
		merge.DirectRecipientsTxs[k] += cand.RecipientsTxs[k]
	}
	for k, v := range cand.DirectSendersAmt {
		merge.DirectSendersAmt[k] += v
		merge.DirectSendersTxs[k] += cand.DirectSendersTxs[k]
	}
}

// DEPRECATED

/*
// MergeIndvData merges multi-year data sets into one Individual object
func MergeIndvData(indvID string, years []string) (*donations.Individual, error) {
	set := make(map[string]*donations.Individual)

	for _, year := range years {
		indv, err := persist.GetObject(year, "individuals", indvID)
		if err != nil {
			fmt.Println("mergeIndvData failed: ", err)
			return nil, fmt.Errorf("mergeIndvData failed: %v", err)
		}
		set[year] = indv.(*donations.Individual)
	}

	merged := set[years[0]]
	for year, indv := range set {
		if year == years[0] {
			continue
		}
		merged.Donations = append(merged.Donations, indv.Donations...)
		merged.TotalDonations += indv.TotalDonations
		merged.TotalDonated += indv.TotalDonated
		merged.AvgDonation = merged.TotalDonated / merged.TotalDonations
		indvMapMerge(merged, indv)
	}

	return merged, nil
}

// MergeCmteData merges multi-year data sets into one Committee object
func MergeCmteData(cmteID string, years []string) (*donations.Committee, error) {
	set := make(map[string]*donations.Committee)

	for _, year := range years {
		cmte, err := persist.GetObject(year, "committees", cmteID)
		if err != nil {
			fmt.Println("mergeCmteData failed: ", err)
			return nil, fmt.Errorf("mergeCmteData failed: %v", err)
		}
		set[year] = cmte.(*donations.Committee)
	}

	merged := set[years[0]]
	for year, cmte := range set {
		if year == years[0] {
			continue
		}
		// Donations
		merged.TotalReceived += cmte.TotalReceived
		merged.TotalDonations += cmte.TotalDonations
		merged.AvgDonation = merged.TotalReceived / merged.TotalDonations

		// Individual Contributions
		for k, v := range cmte.TopIndvDonorsAmt {
			if len(merged.TopIndvDonorsAmt) < 1000 {
				merged.TopIndvDonorsAmt[k] += v
				merged.TopIndvDonorsTxs[k] += cmte.TopIndvDonorsTxs[k]
				delete(cmte.TopIndvDonorsAmt, k)
			} else {
				err := mergeTopIndvTotals(merged, cmte)
				if err != nil {
					fmt.Println("mergeCmteData failed: ", err)
					return nil, fmt.Errorf("mergeCmteData failed: %v", err)
				}
			}
		}

		// Committee Contributions
		for k, v := range cmte.TopCmteDonorsAmt {
			if len(merged.TopCmteDonorsAmt) < 1000 {
				merged.TopCmteDonorsAmt[k] += v
				merged.TopCmteDonorsTxs[k] += cmte.TopCmteDonorsTxs[k]
				delete(cmte.TopCmteDonorsAmt, k)
			} else {
				err := mergeTopCmteTotals(merged, cmte)
				if err != nil {
					fmt.Println("mergeCmteData failed: ", err)
					return nil, fmt.Errorf("mergeCmteData failed: %v", err)
				}
			}
		}

		// Transfers to other committees
		merged.TotalTransferred += cmte.TotalTransferred
		merged.TotalTransfers += cmte.TotalTransfers
		merged.AvgTransfer += merged.TotalTransferred / merged.TotalTransfers
		for k, v := range cmte.AffiliatesAmt {
			merged.AffiliatesAmt[k] += v
			merged.AffiliatesTxs[k] += cmte.AffiliatesTxs[k]
		}

		// Disbursements
		merged.TotalDisbursed += cmte.TotalDisbursed
		merged.TotalDisbursements += cmte.TotalDisbursements
		merged.AvgDisbursed = merged.TotalDisbursed / merged.TotalDisbursements
		for k, v := range cmte.TopDisbRecipientsAmt {
			if len(merged.TopDisbRecipientsAmt) < 1000 {
				merged.TopDisbRecipientsAmt[k] += v
				merged.TopDisbRecipientsTxs[k] += cmte.TopDisbRecipientsTxs[k]
				delete(cmte.TopDisbRecipientsAmt, k)
			} else {
				err := mergeTopDisbRecTotals(merged, cmte)
				if err != nil {
					fmt.Println("mergeCmteData failed: ", err)
					return nil, fmt.Errorf("mergeCmteData failed: %v", err)
				}
			}
		}
	}

	return merged, nil
}

// MergeCandData merges multi-year data sets into one Candidate object
func MergeCandData(candID string, years []string) (*donations.Candidate, error) {
	set := make(map[string]*donations.Candidate)

	for _, year := range years {
		cand, err := persist.GetObject(year, "candidates", candID)
		if err != nil {
			fmt.Println("mergeCandData failed: ", err)
			return nil, fmt.Errorf("mergeCandData failed: %v", err)
		}
		set[year] = cand.(*donations.Candidate)
	}

	merged := set[years[0]]
	for year, cand := range set {
		if year == years[0] {
			continue
		}
		// Donations
		merged.TotalRaised += cand.TotalRaised
		merged.TotalDonations += cand.TotalDonations
		merged.AvgDonation = merged.TotalRaised / merged.TotalDonations

		// Individual donors
		for k, v := range cand.TopIndvDonorsAmt {
			if len(merged.TopIndvDonorsAmt) < 1000 {
				merged.TopIndvDonorsAmt[k] += v
				merged.TopIndvDonorsTxs[k] += cand.TopIndvDonorsTxs[k]
				delete(cand.TopIndvDonorsAmt, k)
			} else {
				err := mergeCandTopIndvTotals(merged, cand)
				if err != nil {
					fmt.Println("mergeCandData failed: ", err)
					return nil, fmt.Errorf("mergeCandData failed: %v", err)
				}
			}
		}

		// Committee donors
		for k, v := range cand.TopCmteDonorsAmt {
			if len(merged.TopCmteDonorsAmt) < 1000 {
				merged.TopCmteDonorsAmt[k] += v
				merged.TopCmteDonorsTxs[k] += cand.TopCmteDonorsTxs[k]
				delete(cand.TopCmteDonorsAmt, k)
			} else {
				err := mergeCandTopCmteTotals(merged, cand)
				if err != nil {
					fmt.Println("mergeCandData failed: ", err)
					return nil, fmt.Errorf("mergeCandData failed: %v", err)
				}
			}
		}

		// Disbursement Recipients
		for k, v := range cand.TopDisbRecsAmt {
			if len(merged.TopDisbRecsAmt) < 1000 {
				merged.TopDisbRecsAmt[k] += v
				merged.TopDisbRecsTxs[k] += cand.TopDisbRecsTxs[k]
				delete(cand.TopDisbRecsAmt, k)
			} else {
				err := mergeCandTopDisbRecTotals(merged, cand)
				if err != nil {
					fmt.Println("mergeCandData failed: ", err)
					return nil, fmt.Errorf("mergeCandData failed: %v", err)
				}
			}
		}
	}

	return merged, nil
}

// MergeDisbRecData merges multi-year data sets into one DisbRecipient object
func MergeDisbRecData(drID string, years []string) (*donations.DisbRecipient, error) {
	set := make(map[string]*donations.DisbRecipient)

	for _, year := range years {
		rec, err := persist.GetObject(year, "disbursement_recipients", drID)
		if err != nil {
			fmt.Println("mergeIndvData failed: ", err)
			return nil, fmt.Errorf("mergeIndvData failed: %v", err)
		}
		set[year] = rec.(*donations.DisbRecipient)
	}

	merged := set[years[0]]
	for year, rec := range set {
		if year == years[0] {
			continue
		}
		merged.Disbursements = append(merged.Disbursements, rec.Disbursements...)
		merged.TotalDisbursements += rec.TotalDisbursements
		merged.TotalReceived += rec.TotalReceived
		merged.AvgReceived = merged.TotalReceived / merged.TotalDisbursements
		drMapMerge(merged, rec)
	}

	return merged, nil
}

func indvMapMerge(merge, indv *donations.Individual) {
	for k, v := range indv.RecipientsAmt {
		merge.RecipientsAmt[k] += v
		merge.RecipientsTxs[k] += indv.RecipientsTxs[k]
	}

}

func mergeTopIndvTotals(merge, cmte *donations.Committee) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(merge.TopIndvDonorThreshold) == 0 {
		es := sortTopX(merge.TopIndvDonorsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopIndvTotals failed: ", err)
			return fmt.Errorf("updateTopIndvTotals failed: %v", err)
		}
	} else {
		for _, entry := range merge.TopIndvDonorThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// merge TopIndvDonor maps
	threshold := least[len(least)-1].Total // last/smallest obj in least
	for k, v := range cmte.TopIndvDonorsAmt {
		if merge.TopIndvDonorsAmt[k] != 0 {
			merge.TopIndvDonorsAmt[k] += v
			merge.TopIndvDonorsTxs[k] += cmte.TopIndvDonorsTxs[k]
			continue
		}

		if v > threshold {
			new := newEntry(cmte.ID, v)
			delID := reSortLeast(new, &least)
			delete(merge.TopIndvDonorsAmt, delID)
			delete(merge.TopIndvDonorsTxs, delID)
			merge.TopIndvDonorsAmt[cmte.ID] = v
			merge.TopIndvDonorsTxs[cmte.ID] = cmte.TopIndvDonorsTxs[k]
		}
	}

	return nil
}

func mergeTopCmteTotals(merge, cmte *donations.Committee) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(merge.TopCmteDonorThreshold) == 0 {
		es := sortTopX(merge.TopCmteDonorsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	} else {
		for _, entry := range merge.TopCmteDonorThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// merge TopCmteDonor maps
	threshold := least[len(least)-1].Total // last/smallest obj in least
	for k, v := range cmte.TopCmteDonorsAmt {
		if merge.TopCmteDonorsAmt[k] != 0 {
			merge.TopCmteDonorsAmt[k] += v
			merge.TopCmteDonorsTxs[k] += cmte.TopCmteDonorsTxs[k]
			continue
		}

		if v > threshold {
			new := newEntry(cmte.ID, v)
			delID := reSortLeast(new, &least)
			delete(merge.TopCmteDonorsAmt, delID)
			delete(merge.TopCmteDonorsTxs, delID)
			merge.TopCmteDonorsAmt[cmte.ID] = v
			merge.TopCmteDonorsTxs[cmte.ID] = cmte.TopCmteDonorsTxs[k]
		}
	}

	return nil
}

func mergeTopDisbRecTotals(merge, cmte *donations.Committee) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(merge.TopRecThreshold) == 0 {
		es := sortTopX(merge.TopDisbRecipientsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	} else {
		for _, entry := range merge.TopRecThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// merge TopDisbRecipient maps
	threshold := least[len(least)-1].Total // last/smallest obj in least
	for k, v := range cmte.TopDisbRecipientsAmt {
		if merge.TopDisbRecipientsAmt[k] != 0 {
			merge.TopDisbRecipientsAmt[k] += v
			merge.TopDisbRecipientsTxs[k] += cmte.TopDisbRecipientsTxs[k]
			continue
		}

		if v > threshold {
			new := newEntry(cmte.ID, v)
			delID := reSortLeast(new, &least)
			delete(merge.TopDisbRecipientsAmt, delID)
			delete(merge.TopDisbRecipientsTxs, delID)
			merge.TopDisbRecipientsAmt[cmte.ID] = v
			merge.TopDisbRecipientsTxs[cmte.ID] = cmte.TopDisbRecipientsTxs[k]
		}
	}

	return nil
}

func mergeCandTopIndvTotals(merge, cand *donations.Candidate) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(merge.TopIDThreshold) == 0 {
		es := sortTopX(merge.TopIndvDonorsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopIndvTotals failed: ", err)
			return fmt.Errorf("updateTopIndvTotals failed: %v", err)
		}
	} else {
		for _, entry := range merge.TopIDThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// merge TopIndvDonor maps
	threshold := least[len(least)-1].Total // last/smallest obj in least
	for k, v := range cand.TopIndvDonorsAmt {
		if merge.TopIndvDonorsAmt[k] != 0 {
			merge.TopIndvDonorsAmt[k] += v
			merge.TopIndvDonorsTxs[k] += cand.TopIndvDonorsTxs[k]
			continue
		}

		if v > threshold {
			new := newEntry(cand.ID, v)
			delID := reSortLeast(new, &least)
			delete(merge.TopIndvDonorsAmt, delID)
			delete(merge.TopIndvDonorsTxs, delID)
			merge.TopIndvDonorsAmt[cand.ID] = v
			merge.TopIndvDonorsTxs[cand.ID] = cand.TopIndvDonorsTxs[k]
		}
	}

	return nil
}

func mergeCandTopCmteTotals(merge, cand *donations.Candidate) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(merge.TopCDThreshold) == 0 {
		es := sortTopX(merge.TopCmteDonorsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	} else {
		for _, entry := range merge.TopCDThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// merge TopCmteDonor maps
	threshold := least[len(least)-1].Total // last/smallest obj in least
	for k, v := range cand.TopCmteDonorsAmt {
		if merge.TopCmteDonorsAmt[k] != 0 {
			merge.TopCmteDonorsAmt[k] += v
			merge.TopCmteDonorsTxs[k] += cand.TopCmteDonorsTxs[k]
			continue
		}

		if v > threshold {
			new := newEntry(cand.ID, v)
			delID := reSortLeast(new, &least)
			delete(merge.TopCmteDonorsAmt, delID)
			delete(merge.TopCmteDonorsTxs, delID)
			merge.TopCmteDonorsAmt[cand.ID] = v
			merge.TopCmteDonorsTxs[cand.ID] = cand.TopCmteDonorsTxs[k]
		}
	}

	return nil
}

func mergeCandTopDisbRecTotals(merge, cand *donations.Candidate) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(merge.TopDRThreshold) == 0 {
		es := sortTopX(merge.TopDisbRecsAmt)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	} else {
		for _, entry := range merge.TopDRThreshold {
			least = append(least, entry.(*donations.Entry))
		}
	}

	// merge TopDisbRecipient maps
	threshold := least[len(least)-1].Total // last/smallest obj in least
	for k, v := range cand.TopDisbRecsAmt {
		if merge.TopDisbRecsAmt[k] != 0 {
			merge.TopDisbRecsAmt[k] += v
			merge.TopDisbRecsTxs[k] += cand.TopDisbRecsTxs[k]
			continue
		}

		if v > threshold {
			new := newEntry(cand.ID, v)
			delID := reSortLeast(new, &least)
			delete(merge.TopDisbRecsAmt, delID)
			delete(merge.TopDisbRecsTxs, delID)
			merge.TopDisbRecsAmt[cand.ID] = v
			merge.TopDisbRecsTxs[cand.ID] = cand.TopDisbRecsTxs[k]
		}
	}

	return nil
}

func drMapMerge(merge, rec *donations.DisbRecipient) {
	for k, v := range rec.SendersAmt {
		merge.SendersAmt[k] += v
		merge.SendersTxs[k] += rec.SendersTxs[k]
	}

}
*/
