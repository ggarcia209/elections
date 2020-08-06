package databuilder

import (
	"fmt"

	"github.com/elections/donations"
)

/* may move code to separate "aggregate" package */
// basic function tests - successful
// panics if year does not contain corresponding data (nil interface)
// refactor to implement DynamoDB API calls to retreive objects
// re-do unit tests w/ edge cases
// removed org references

// MergeData merges multi-year data sets into one interface object
func MergeData(years []string, ID, bucket string) (interface{}, error) {
	var merged interface{}
	switch {
	case bucket == "individuals":
		// get corresponding objects for each year
		set, err := createMergeSet(years, bucket, ID)
		if err != nil {
			fmt.Println("MergeData failed: ", err)
			return nil, fmt.Errorf("MergeData failed: %v", err)
		}

		// merge object values into one object
		mergedIndv := *set[years[0]].(*donations.Individual)
		for year, obj := range set {
			if year == years[0] {
				continue
			}
			if obj == nil {
				continue
			}
			compIndv := *obj.(*donations.Individual)
			indvTotalsMerge(&mergedIndv, &compIndv)
			indvMapMerge(&mergedIndv, &compIndv)
		}

		// Filter Top 100 entries
		if len(mergedIndv.RecipientsAmt) > 100 {
			mergedIndv.RecipientsAmt, mergedIndv.RecipientsTxs = sort100(mergedIndv.RecipientsAmt, mergedIndv.RecipientsTxs)
		}
		if len(mergedIndv.SendersAmt) > 100 {
			mergedIndv.SendersAmt, mergedIndv.SendersTxs = sort100(mergedIndv.SendersAmt, mergedIndv.SendersTxs)
		}

		merged = mergedIndv
	case bucket == "cmte_tx_data":
		set, err := createMergeSet(years, bucket, ID)
		if err != nil {
			fmt.Println("MergeData failed: ", err)
			return nil, fmt.Errorf("MergeData failed: %v", err)
		}

		mergedCmte := *set[years[0]].(*donations.CmteTxData)
		for year, obj := range set {
			if year == years[0] {
				continue
			}
			if obj == nil {
				continue
			}
			compCmte := *obj.(*donations.CmteTxData)
			cmteTxTotalsMerge(&mergedCmte, &compCmte)
			cmteTxMapMerge(&mergedCmte, &compCmte)
		}

		// Filter Top 100 entries
		if len(mergedCmte.TopIndvContributorsAmt) > 3 {
			mergedCmte.TopIndvContributorsAmt, mergedCmte.TopIndvContributorsTxs = sort3(mergedCmte.TopIndvContributorsAmt, mergedCmte.TopIndvContributorsTxs)
		}
		if len(mergedCmte.TopCmteOrgContributorsAmt) > 3 {
			mergedCmte.TopCmteOrgContributorsAmt, mergedCmte.TopCmteOrgContributorsTxs = sort3(mergedCmte.TopCmteOrgContributorsAmt, mergedCmte.TopCmteOrgContributorsTxs)
		}
		if len(mergedCmte.TransferRecsAmt) > 3 {
			mergedCmte.TransferRecsAmt, mergedCmte.TransferRecsTxs = sort3(mergedCmte.TransferRecsAmt, mergedCmte.TransferRecsTxs)
		}
		if len(mergedCmte.TopExpRecipientsAmt) > 3 {
			mergedCmte.TopExpRecipientsAmt, mergedCmte.TopExpRecipientsTxs = sort3(mergedCmte.TopExpRecipientsAmt, mergedCmte.TopExpRecipientsTxs)
		}

		merged = mergedCmte
	case bucket == "candidates":
		set, err := createMergeSet(years, bucket, ID)
		if err != nil {
			fmt.Println("MergeData failed: ", err)
			return nil, fmt.Errorf("MergeData failed: %v", err)
		}

		mergedCand := *set[years[0]].(*donations.Candidate)
		for year, obj := range set {
			if year == years[0] {
				continue
			}
			if obj == nil {
				continue
			}
			compCand := *obj.(*donations.Candidate)
			candTotalsMerge(&mergedCand, &compCand)
			candMapMerge(&mergedCand, &compCand)
		}

		if len(mergedCand.DirectRecipientsAmts) > 100 {
			mergedCand.DirectRecipientsAmts, mergedCand.DirectRecipientsTxs = sort100(mergedCand.DirectRecipientsAmts, mergedCand.DirectRecipientsTxs)
		}
		if len(mergedCand.DirectSendersAmts) > 100 {
			mergedCand.DirectSendersAmts, mergedCand.DirectSendersTxs = sort100(mergedCand.DirectSendersAmts, mergedCand.DirectSendersTxs)
		}

		merged = mergedCand
	default:
		return nil, fmt.Errorf("MergeData failed: invalid bucket type")
	}

	return merged, nil
}

func createMergeSet(years []string, bucket, ID string) (map[string]interface{}, error) {
	set := make(map[string]interface{})
	for _, year := range years {
		// OBJECT MUST BE RETREIVED FROM DynamoDB API CALL
		// obj := DbSim[year][bucket][ID] // test only
		obj := &donations.Individual{} // test only
		// verify obj != nil before adding to set
		set[year] = obj
	}
	return set, nil
}

func createMergeObj(obj interface{}) interface{} {
	merge := obj
	return merge
}

func mapMerge(merge, source map[string]float32) map[string]float32 {
	mergeMap := make(map[string]float32)
	for k, v := range merge {
		mergeMap[k] += v
	}
	for k, v := range source {
		// add directly to map
		mergeMap[k] += v
	}

	return mergeMap
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
	merge.RecipientsAmt = mapMerge(merge.RecipientsAmt, indv.RecipientsAmt)
	merge.RecipientsTxs = mapMerge(merge.RecipientsTxs, indv.RecipientsTxs)
	merge.SendersAmt = mapMerge(merge.SendersAmt, indv.SendersAmt)
	merge.SendersTxs = mapMerge(merge.SendersTxs, indv.SendersTxs)
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
	merge.TopIndvContributorsAmt = mapMerge(merge.TopIndvContributorsAmt, cmte.TopIndvContributorsAmt)
	merge.TopIndvContributorsTxs = mapMerge(merge.TopIndvContributorsTxs, cmte.TopIndvContributorsTxs)

	// Top Committee/Organization Contributors
	merge.TopCmteOrgContributorsAmt = mapMerge(merge.TopCmteOrgContributorsAmt, cmte.TopCmteOrgContributorsAmt)
	merge.TopCmteOrgContributorsTxs = mapMerge(merge.TopCmteOrgContributorsTxs, cmte.TopCmteOrgContributorsTxs)

	// Transfers Recipients
	merge.TransferRecsAmt = mapMerge(merge.TransferRecsAmt, cmte.TransferRecsAmt)
	merge.TransferRecsTxs = mapMerge(merge.TransferRecsTxs, cmte.TransferRecsTxs)

	// Top Expenditure Recipients
	merge.TopExpRecipientsAmt = mapMerge(merge.TopExpRecipientsAmt, cmte.TopExpRecipientsAmt)
	merge.TopExpRecipientsTxs = mapMerge(merge.TopExpRecipientsTxs, cmte.TopExpRecipientsTxs)
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
	merge.DirectRecipientsAmts = mapMerge(merge.DirectRecipientsAmts, cand.DirectRecipientsAmts)
	merge.DirectRecipientsTxs = mapMerge(merge.DirectRecipientsTxs, cand.DirectRecipientsTxs)
	merge.DirectSendersAmts = mapMerge(merge.DirectSendersAmts, cand.DirectSendersAmts)
	merge.DirectSendersTxs = mapMerge(merge.DirectSendersTxs, cand.DirectSendersTxs)
}

// Sort maps and derive top 100 entries by value
func sort100(amts, txs map[string]float32) (map[string]float32, map[string]float32) {
	topAmts := make(map[string]float32)
	topTxs := make(map[string]float32)
	es := sortTopX(amts)

	for _, e := range es[:100] {
		topAmts[e.ID] = e.Total
		topTxs[e.ID] = txs[e.ID]
	}

	return topAmts, topTxs
}

// TEST ONLY
// Sort maps and derive top 5 entries by value
func sort3(amts, txs map[string]float32) (map[string]float32, map[string]float32) {
	topAmts := make(map[string]float32)
	topTxs := make(map[string]float32)
	es := sortTopX(amts)

	for _, e := range es[:3] {
		topAmts[e.ID] = e.Total
		topTxs[e.ID] = txs[e.ID]
	}

	return topAmts, topTxs
}
