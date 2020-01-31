package main

import (
	"fmt"
	"github.com/elections/donations"
	"github.com/elections/persist"
)

type Merge struct {
	AggAmt     float32
	AggTxs     float32
	AggAvg     float32
	AggAmtsMap map[string]float32
	AggTxsMap  map[string]float32
	AggTotal   float32
}

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
			cmteTxTotalsMerge(merged.(*donations.CmteTxData), obj.(*donations.CmteTxData))
			cmteTxMapMerge(merged.(*donations.CmteTxData), obj.(*donations.CmteTxData))
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
			mergeTopTotals(&merge.TopIndvContributorsAmt, &merge.TopIndvContributorsTxs, &cmte.TopIndvContributorsAmt, &cmte.TopIndvContributorsTxs, &merge.TopIndvContributorThreshold)
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
			mergeTopTotals(&merge.TopCmteOrgContributorsAmt, &merge.TopCmteOrgContributorsTxs, &cmte.TopCmteOrgContributorsAmt, &cmte.TopCmteOrgContributorsTxs, &merge.TopCmteOrgContributorThreshold)
		}
	}

	// Transfers Recipients
	for k, v := range cmte.TransferRecsAmt {
		merge.TransferRecsAmt[k] += v
		merge.TransferRecsTxs[k] += cmte.TransferRecsTxs[k]
	}

	// Top Expenditure Recipients
	for k, v := range cmte.TopExpRecipientsAmt{
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
			mergeTopTotals(&merge.TopExpRecipientsAmt, &merge.TopExpRecipientsTxs, &cmte.TopExpRecipientsAmt, &cmte.TopExpRecipientsTxs, &merge.TopExpThreshold)
		}
	}
}


func mergeTopTotals(mergeAmts, mergeTxs, mAmts, mTxs *map[string]float32, mergeTh *[]interface{}) error {
	// set/reset least threshold list
	var least Entries
	var err error
	if len(mergeTh) == 0 {
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
			new := newEntry(cmte.ID, v)
			delID := reSortLeast(new, &least)
			delete(mergeAmts, delID)
			delete(mergeTxs, delID)
			mergeAmts[k] = v
			mergeTxs[k] = mTxs[k]
		}
	}

	// update object's threshold list
	th := []interface{}
	for _, entry := range least {
		th = append(th, entry)
	}
	mergeTh = append(mergeTh[:0], th...)

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

