package databuilder

import (
	"fmt"
	"strings"

	"github.com/elections/donations"
)

func deriveTxTypes(cont *donations.Contribution) (string, bool, bool) {
	// initialize return values
	incoming := false
	memo := false
	var bucket string

	// determine tx type (incoming / outgoing / memo)
	numCode := cont.TxType
	if numCode < "20" || (numCode >= "30" && numCode < "33") {
		incoming = true
	}
	if cont.MemoCode == "X" {
		memo = true
	}

	// determine contributor/receiver type - derive from OtherID
	IDss := strings.Split(cont.OtherID, "")
	idCode := IDss[0]
	switch {
	case idCode == "C":
		bucket = "committees"
	case idCode == "H" || idCode == "S" || idCode == "P":
		bucket = "candidates"
	default:
		if cont.Occupation == "" {
			bucket = "organizations"
		} else {
			bucket = "individuals"
		}
	}
	return bucket, incoming, memo
}

// ContributionUpdate updates the filing committee's transaction values
// and the transaction values of the contributing Individual/Organization if applicable.
// Committees identified in the contribution's OtherID field are not updated.
/* func ContributionUpdate(year string, conts []*donations.Contribution) error {
	for _, cont := range conts {
		// get tx type info
		bucket, incoming, memo := deriveTxTypes(cont)

		// get sender individual objects
		filer, err := persist.GetObject(year, "committees", cont.CmteID)
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}

		// get sender/receiver objects
		receiver, err := persist.GetObject(year, bucket, cont.OtherID)
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}

		// update sender & receiver's values
		err = updateIndvDonations(cont, sender.(*donations.Individual), receiver.(*donations.Committee))
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}

		// update Top Individual Donors Overall
		err = updateTopIndviduals(year, sender.(*donations.Individual))
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}

		// update Top Committes Overall by funds raised
		err = updateTopCmteRecs(year, receiver.(*donations.Committee))
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}

		// update linked Candidate, if any
		if receiver.(*donations.Committee).CandID != "" {
			cand, err := persist.GetObject(year, "candidates", receiver.(*donations.Committee).CandID)
			if err != nil {
				fmt.Println("ContributionUpdate failed: ", err)
				return fmt.Errorf("ContributionUpdate failed: %v", err)
			}

			err = updateCandIndvDonations(cont, sender.(*donations.Individual), cand.(*donations.Candidate))
			if err != nil {
				fmt.Println("ContributionUpdate failed: ", err)
				return fmt.Errorf("ContributionUpdate failed: %v", err)
			}

			err = updateTopCandidates(year, cand.(*donations.Candidate))
			if err != nil {
				fmt.Println("ContributionUpdate failed: ", err)
				return fmt.Errorf("ContributionUpdate failed: %v", err)
			}

			err = persist.PutObject(year, cand)
			if err != nil {
				fmt.Println("ContributionUpdate failed: ", err)
				return fmt.Errorf("ContributionUpdate failed: %v", err)
			}

		}

		// persist updated cmte objects
		err = persist.PutObject(year, receiver)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}

		// persist indv donor objects
		err = persist.PutObject(year, sender)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}
	}

	return nil
} */

// TX CRITERIA
// if sender == committee: do not credit sender's accounts -- account for corresponding tx
// 	all filers are committees only and *should* have corresponding transaction for sender
// if sender != committee: credit accounts -- no corresponding tx
func incomingTxUpdate(cont *donations.Contribution, filerData *donations.CmteTxData, sender interface{}, memo bool) error {
	// update maps only if memo == true
	// account for percentage of amounts received by memo transactions but do not add to totals
	if memo {
		// update maps
		err := mapUpdate(cont, filerData, sender)
		if err != nil {
			fmt.Println("incomingTxUpdate failed: ", err)
			return fmt.Errorf("incomingTxUpdate failed: %v", err)
		}
		return nil
	}

	// debit Contributions or OtherReceipts and TotalIncoming
	if cont.TxType > "15" && cont.TxType < "18" {
		filerData.OtherReceiptsInAmt += cont.TxAmt
		filerData.OtherReceiptsInTxs++
		filerData.AvgOtherIn = filerData.OtherReceiptsInAmt / filerData.OtherReceiptsInTxs
	} else {
		filerData.ContributionsInAmt += cont.TxAmt
		filerData.ContributionsInTxs++
		filerData.AvgContributionIn = filerData.ContributionsInAmt / filerData.ContributionsInTxs
	}
	filerData.TotalIncomingAmt = filerData.ContributionsInAmt + filerData.OtherReceiptsInAmt
	filerData.TotalIncomingTxs = filerData.ContributionsInTxs + filerData.OtherReceiptsInTxs
	filerData.AvgIncoming = filerData.TotalIncomingAmt / filerData.TotalIncomingTxs
	filerData.NetBalance = filerData.TotalIncomingAmt - filerData.TotalOutgoingAmt

	// credit sender account
	// ADD LOGIC TO ACCOUNT FOR CORRESPONDING TXs
	switch t := sender.(type) {
	case *donations.Individual:
		sender.(*donations.Individual).TotalOutAmt += cont.TxAmt
		sender.(*donations.Individual).TotalOutTxs++
		sender.(*donations.Individual).AvgTxOut = sender.(*donations.Individual).TotalOutAmt / sender.(*donations.Individual).TotalOutTxs
		sender.(*donations.Individual).NetBalance = sender.(*donations.Individual).TotalInAmt - sender.(*donations.Individual).TotalOutAmt
	case *donations.Organization:
		sender.(*donations.Organization).TotalOutAmt += cont.TxAmt
		sender.(*donations.Organization).TotalOutTxs++
		sender.(*donations.Organization).AvgTxOut = sender.(*donations.Organization).TotalOutAmt / sender.(*donations.Organization).TotalOutTxs
		sender.(*donations.Organization).NetBalance = sender.(*donations.Organization).TotalInAmt - sender.(*donations.Organization).TotalOutAmt
	case *donations.Candidate:
		sender.(*donations.Candidate).TotalDirectOutAmt += cont.TxAmt
		sender.(*donations.Candidate).TotalDirectOutTxs++
		sender.(*donations.Candidate).AvgDirectOut = sender.(*donations.Candidate).TotalDirectOutAmt / sender.(*donations.Candidate).TotalDirectOutTxs
		sender.(*donations.Candidate).NetBalanceDirectTx = sender.(*donations.Candidate).TotalDirectInAmt - sender.(*donations.Candidate).TotalDirectOutAmt
	case *donations.CmteTxData:
		// do nothing -- accounted for by sender's corresponding outgoing transaction
	default:
		return fmt.Errorf("incomingTxUpdate failed: wrong interface type")
	}

	// update maps
	err := mapUpdate(cont, filerData, sender)
	if err != nil {
		fmt.Println("incomingTxUpdate failed: ", err)
		return fmt.Errorf("incomingTxUpdate failed: %v", err)
	}
	return nil
}

func outgoingTxUpdate(cont *donations.Contribution, filerData *donations.CmteTxData, sender interface{}, memo bool) error {
	// update maps only if memo == true
	// account for percentage of amounts received by memo transactions but do not add to totals
	if memo {
		// update maps
		err := mapUpdate(cont, filerData, sender)
		if err != nil {
			fmt.Println("outgoingTxUpdate failed: ", err)
			return fmt.Errorf("outgoingTxUpdate failed: %v", err)
		}
		return nil
	}

	// credit Transfers or Expenditures and TotalOutgoing
	if cont.TxType > "15" && cont.TxType < "18" { // REFACTOR to transfers/contributions only
		// transfers tx types: 22H, 24G, 24H, 24K, 24U, 24Z, 24I*, 24T*,
		// *unverified through official filings
		filerData.TransfersAmt += cont.TxAmt
		filerData.TransfersTxs++
		filerData.AvgTransfer = filerData.TransfersAmt / filerData.TransfersTxs
	} else {
		filerData.ExpendituresAmt += cont.TxAmt
		filerData.ExpendituresTxs++
		filerData.AvgContributionIn = filerData.ExpendituresAmt / filerData.ExpendituresTxs
	}
	filerData.TotalOutgoingAmt = filerData.TransfersAmt + filerData.ExpendituresAmt
	filerData.TotalOutgoingTxs = filerData.TransfersTxs + filerData.ExpendituresTxs
	filerData.AvgOutgoing = filerData.TotalOutgoingAmt / filerData.TotalOutgoingTxs
	filerData.NetBalance = filerData.TotalIncomingAmt - filerData.TotalOutgoingAmt

	// debit sender account
	// ADD LOGIC TO ACCOUNT FOR CORRESPONDING TXs
	switch t := sender.(type) {
	case *donations.Individual:
		sender.(*donations.Individual).TotalInAmt += cont.TxAmt
		sender.(*donations.Individual).TotalInTxs++
		sender.(*donations.Individual).AvgTxIn = sender.(*donations.Individual).TotalInAmt / sender.(*donations.Individual).TotalInTxs
		sender.(*donations.Individual).NetBalance = sender.(*donations.Individual).TotalInAmt - sender.(*donations.Individual).TotalOutAmt
	case *donations.Organization:
		sender.(*donations.Organization).TotalInAmt += cont.TxAmt
		sender.(*donations.Organization).TotalInTxs++
		sender.(*donations.Organization).AvgTxIn = sender.(*donations.Organization).TotalInAmt / sender.(*donations.Organization).TotalInTxs
		sender.(*donations.Organization).NetBalance = sender.(*donations.Organization).TotalInAmt - sender.(*donations.Organization).TotalOutAmt
	case *donations.Candidate:
		sender.(*donations.Candidate).TotalDirectInAmt += cont.TxAmt
		sender.(*donations.Candidate).TotalDirectInTxs++
		sender.(*donations.Candidate).AvgDirectIn = sender.(*donations.Candidate).TotalDirectInAmt / sender.(*donations.Candidate).TotalDirectInTxs
		sender.(*donations.Candidate).NetBalanceDirectTx = sender.(*donations.Candidate).TotalDirectInAmt - sender.(*donations.Candidate).TotalDirectOutAmt
	case *donations.CmteTxData:
		// special case -- pass sender as filer (receiving committee)
		// and filer as sender (sending committee) to mapUpdate()
		err := mapUpdate(cont, sender.(*donations.CmteTxData), filerData)
		if err != nil {
			fmt.Println("outgoingTxUpdate failed: ", err)
			return fmt.Errorf("outgoingTxUpdate failed: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("outgoingTxUpdate failed: wrong interface type")
	}

	// update maps
	err := mapUpdate(cont, filerData, sender)
	if err != nil {
		fmt.Println("outgoingTxUpdate failed: ", err)
		return fmt.Errorf("outgoingTxUpdate failed: %v", err)
	}
	return nil
}

// REFACTOR TO ACCOUNT FOR INCOMING VS OUTGOING TXs
func mapUpdate(cont *donations.Contribution, filerData *donations.CmteTxData, sender interface{}) error {
	switch t := sender.(type) {
	case *donations.Individual:
		// re-initialize maps if nil
		if len(filerData.TopIndvContributorsAmt) == 0 {
			filerData.TopIndvContributorsAmt = make(map[string]float32)
			filerData.TopIndvContributorsTxs = make(map[string]float32)
		}
		if len(sender.(*donations.Individual).RecipientsAmt) == 0 {
			sender.(*donations.Individual).RecipientsAmt = make(map[string]float32)
			sender.(*donations.Individual).RecipientsTxs = make(map[string]float32)
		}

		// update sender's Recipients maps
		sender.(*donations.Individual).RecipientsAmt[filerData.CmteID] += cont.TxAmt
		sender.(*donations.Individual).RecipientsTxs[filerData.CmteID]++

		// update filing committee's Top Contributors maps
		if len(filerData.TopIndvContributorsAmt) < 1000 {
			filerData.TopIndvContributorsAmt[sender.(*donations.Individual).ID] += cont.TxAmt
			filerData.TopIndvContributorsTxs[sender.(*donations.Individual).ID]++
		} else {
			// update filer's top contributors maps
			err := updateTopDonors(filerData, sender, cont)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
		}
	case *donations.Organization:
		// re-initialize maps if nil
		if len(filerData.TopCmteOrgContributorsAmt) == 0 {
			filerData.TopCmteOrgContributorsAmt = make(map[string]float32)
			filerData.TopCmteOrgContributorsTxs = make(map[string]float32)
		}
		if len(sender.(*donations.Organization).RecipientsAmt) == 0 {
			sender.(*donations.Organization).RecipientsAmt = make(map[string]float32)
			sender.(*donations.Organization).RecipientsTxs = make(map[string]float32)
		}

		// update sender's Recipients maps
		sender.(*donations.Organization).RecipientsAmt[filerData.CmteID] += cont.TxAmt
		sender.(*donations.Organization).RecipientsTxs[filerData.CmteID]++

		// update filing committee's Top Contributors maps
		if len(filerData.TopCmteOrgContributorsAmt) < 1000 {
			filerData.TopCmteOrgContributorsAmt[sender.(*donations.Organization).ID] += cont.TxAmt
			filerData.TopCmteOrgContributorsTxs[sender.(*donations.Organization).ID]++
		} else {
			// update filer's top contributors maps
			err := updateTopDonors(filerData, sender, cont)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
		}
	case *donations.CmteTxData:
		// ADD LOGIC TO DISTINGUISH BETWEEN ACCOUNT TYPES
		// re-initialize maps if nil
		if len(filerData.TopCmteOrgContributorsAmt) == 0 {
			filerData.TopCmteOrgContributorsAmt = make(map[string]float32)
			filerData.TopCmteOrgContributorsTxs = make(map[string]float32)
		}
		if len(sender.(*donations.CmteTxData).TransferRecsAmt) == 0 {
			sender.(*donations.CmteTxData).TransferRecsAmt = make(map[string]float32)
			sender.(*donations.CmteTxData).TransferRecsTxs = make(map[string]float32)
		}

		// update sender's Affiliates maps
		sender.(*donations.CmteTxData).TransferRecsAmt[filerData.CmteID] += cont.TxAmt
		sender.(*donations.CmteTxData).TransferRecsTxs[filerData.CmteID]++

		// update filing committee's Top Contributors maps
		if len(filerData.TopCmteOrgContributorsAmt) < 1000 {
			filerData.TopCmteOrgContributorsAmt[sender.(*donations.CmteTxData).CmteID] += cont.TxAmt
			filerData.TopCmteOrgContributorsTxs[sender.(*donations.CmteTxData).CmteID]++
		} else {
			// update filer's top contributors maps
			err := updateTopDonors(filerData, sender, cont)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
		}
	case *donations.Candidate:
		// re-initialize maps if nil
		if len(filerData.TopIndvContributorsAmt) == 0 {
			filerData.TopIndvContributorsAmt = make(map[string]float32)
			filerData.TopIndvContributorsTxs = make(map[string]float32)
		}
		if len(sender.(*donations.Candidate).DirectRecipientsAmts) == 0 {
			sender.(*donations.Candidate).DirectRecipientsAmts = make(map[string]float32)
			sender.(*donations.Candidate).DirectRecipientsTxs = make(map[string]float32)
		}

		// update sender's Recipients maps
		sender.(*donations.Candidate).DirectRecipientsAmts[filerData.CmteID] += cont.TxAmt
		sender.(*donations.Candidate).DirectRecipientsTxs[filerData.CmteID]++

		// update filing committee's Top Contributors maps
		if len(filerData.TopIndvContributorsAmt) < 1000 {
			filerData.TopIndvContributorsAmt[sender.(*donations.Candidate).ID] += cont.TxAmt
			filerData.TopIndvContributorsTxs[sender.(*donations.Candidate).ID]++
		} else {
			// update filer's top contributors maps
			err := updateTopDonors(filerData, sender, cont)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
		}
	default:
		fmt.Errorf("mapUpdate failed: wrong interface type")
	}
	return nil
}

// DEPRECATED
/*

func updateCmteDonations(cont *donations.CmteContribution, sen, rec *donations.Committee) error {
	// update receiver's donations values
	rec.TotalReceived += cont.TxAmt
	rec.TotalDonations++
	rec.AvgDonation = rec.TotalReceived / rec.TotalDonations

	// update senders' transfer values
	sen.TotalTransferred += cont.TxAmt
	sen.TotalTransfers++
	sen.AvgTransfer = sen.TotalTransferred / sen.TotalTransfers

	// update senders's affilates values
	if len(sen.AffiliatesAmt) == 0 {
		sen.AffiliatesAmt = make(map[string]float32)
		sen.AffiliatesTxs = make(map[string]float32)
	}
	sen.AffiliatesAmt[rec.ID] += cont.TxAmt
	sen.AffiliatesTxs[rec.ID]++

	// update receiver's TopCmteDonors values
	if len(rec.TopCmteDonorsAmt) == 0 {
		rec.TopCmteDonorsAmt = make(map[string]float32)
		rec.TopCmteDonorsTxs = make(map[string]float32)
	}
	if len(rec.TopCmteDonorsAmt) < 1000 {
		rec.TopCmteDonorsAmt[sen.ID] = sen.AffiliatesAmt[rec.ID]
		rec.TopCmteDonorsTxs[sen.ID]++
	} else {
		err := updateTopCmteTotals(sen, rec)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	}

	return nil
}

// CmteDisbUpdate updates the dynamic values of both sender committee and reciever DisbRecipient objects for each Disbursement
func CmteDisbUpdate(year string, disbs []*donations.Disbursement) error {
	for _, disb := range disbs {
		// get sender committee object
		sen, err := persist.GetObject(year, "committees", disb.CmteID)
		if err != nil {
			fmt.Println("CmteDisbUpdate failed: ", err)
			return fmt.Errorf("CmteDisbUpdate failed: %v", err)
		}

		// get DisbRecipeint object
		rec, err := persist.GetObject(year, "disbursement_recipients", disb.RecID)
		if err != nil {
			fmt.Println("CmteDisbUpdate failed: ", err)
			return fmt.Errorf("CmteDisbUpdate failed: %v", err)
		}

		// update values
		err = updateDisbursement(disb, sen.(*donations.Committee), rec.(*donations.DisbRecipient))
		if err != nil {
			fmt.Println("CmteDisbUpdate failed: ", err)
			return fmt.Errorf("CmteDisbUpdate failed: %v", err)
		}

		// update Top Committees Overall by funds disbursed
		err = updateTopCmteExp(year, sen.(*donations.Committee))
		if err != nil {
			fmt.Println("CmteDisbUpdate failed: ", err)
			return fmt.Errorf("CmteDisbUpdate failed: %v", err)
		}

		// update Top Disbursement Recipients
		err = updateTopDisbRecs(year, rec.(*donations.DisbRecipient))
		if err != nil {
			fmt.Println("CmteDisbUpdate failed: ", err)
			return fmt.Errorf("CmteDisbUpdate failed: %v", err)
		}

		// update linked Candidate, if any
		if sen.(*donations.Committee).CandID != "" {
			cand, err := persist.GetObject(year, "candidates", sen.(*donations.Committee).CandID)
			if err != nil {
				fmt.Println("CmteDisbUpdate failed: ", err)
				return fmt.Errorf("CmteDisbUpdate failed: %v", err)
			}

			err = updateCandDisbursements(disb, rec.(*donations.DisbRecipient), cand.(*donations.Candidate))
			if err != nil {
				fmt.Println("CmteDisbUpdate failed: ", err)
				return fmt.Errorf("CmteDisbUpdate failed: %v", err)
			}

			// update Top Candidates by funds disbursed
			err = updateTopCandExp(year, cand.(*donations.Candidate))
			if err != nil {
				fmt.Println("CmteDisbUpdate failed: ", err)
				return fmt.Errorf("CmteDisbUpdate failed: %v", err)
			}

			err = persist.PutObject(year, cand)
			if err != nil {
				fmt.Println("CmteDisbUpdate failed: ", err)
				return fmt.Errorf("CmteDisbUpdate failed: %v", err)
			}
		}

		// persist updated cmte object
		err = persist.PutObject(year, sen)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}

		// persist updated DisbRecipient object
		err = persist.PutObject(year, rec)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}
	}

	return nil
}



// IndvContUpdate updates the dynamic values of both sender/receiver committees for each IndvContribution
func IndvContUpdate(year string, conts []*donations.IndvContribution) error {
	for _, cont := range conts {
		// get sender individual objects
		sender, err := persist.GetObject(year, "individuals", cont.DonorID)
		if err != nil {
			fmt.Println("IndvContUpdate failed: ", err)
			return fmt.Errorf("IndvContUpdate failed: %v", err)
		}

		// get receiver committee objects
		receiver, err := persist.GetObject(year, "committees", cont.CmteID)
		if err != nil {
			fmt.Println("IndvContUpdate failed: ", err)
			return fmt.Errorf("IndvContUpdate failed: %v", err)
		}

		// update sender & receiver's values
		err = updateIndvDonations(cont, sender.(*donations.Individual), receiver.(*donations.Committee))
		if err != nil {
			fmt.Println("IndvContUpdate failed: ", err)
			return fmt.Errorf("IndvContUpdate failed: %v", err)
		}

		// update Top Individual Donors Overall
		err = updateTopIndviduals(year, sender.(*donations.Individual))
		if err != nil {
			fmt.Println("IndvContUpdate failed: ", err)
			return fmt.Errorf("IndvContUpdate failed: %v", err)
		}

		// update Top Committes Overall by funds raised
		err = updateTopCmteRecs(year, receiver.(*donations.Committee))
		if err != nil {
			fmt.Println("IndvContUpdate failed: ", err)
			return fmt.Errorf("IndvContUpdate failed: %v", err)
		}

		// update linked Candidate, if any
		if receiver.(*donations.Committee).CandID != "" {
			cand, err := persist.GetObject(year, "candidates", receiver.(*donations.Committee).CandID)
			if err != nil {
				fmt.Println("IndvContUpdate failed: ", err)
				return fmt.Errorf("IndvContUpdate failed: %v", err)
			}

			err = updateCandIndvDonations(cont, sender.(*donations.Individual), cand.(*donations.Candidate))
			if err != nil {
				fmt.Println("IndvContUpdate failed: ", err)
				return fmt.Errorf("IndvContUpdate failed: %v", err)
			}

			err = updateTopCandidates(year, cand.(*donations.Candidate))
			if err != nil {
				fmt.Println("IndvContUpdate failed: ", err)
				return fmt.Errorf("IndvContUpdate failed: %v", err)
			}

			err = persist.PutObject(year, cand)
			if err != nil {
				fmt.Println("IndvContUpdate failed: ", err)
				return fmt.Errorf("IndvContUpdate failed: %v", err)
			}

		}

		// persist updated cmte objects
		err = persist.PutObject(year, receiver)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}

		// persist indv donor objects
		err = persist.PutObject(year, sender)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}
	}

	return nil
}

// CmteContUpdate updates the dynamic values of both sender/receiver committees for each CmteContribution
func CmteContUpdate(year string, conts []*donations.CmteContribution) error {
	for _, cont := range conts {

		// get sender & receiver committee objects
		// refactor to account for non-cmte senders (hash ID) and candidates
		sender, err := persist.GetObject(year, sendID, cont.OtherID)
		if err != nil {
			fmt.Println("CmteContUpdate failed: ", err)
			return fmt.Errorf("CmteContUpdate failed: %v", err)
		}

		receiver, err := persist.GetObject(year, "committees", cont.CmteID)
		if err != nil {
			fmt.Println("CmteContUpdate failed: ", err)
			return fmt.Errorf("CmteContUpdate failed: %v", err)
		}
		// testing only
		fmt.Println("cont.TxID: ", cont.TxID)
		fmt.Println("cont.CmteID (receiver): ", cont.CmteID)
		fmt.Println("cont.OtherID (sender): ", cont.OtherID)
		fmt.Println("cmte: ", receiver)

		// update sender & receiver's values
		// refactor to account for non-cmte senders (hash ID)
		err = updateCmteDonations(cont, sender.(*donations.Committee), receiver.(*donations.Committee))
		if err != nil {
			fmt.Println("CmteContUpdate failed: ", err)
			return fmt.Errorf("CmteContUpdate failed: %v", err)
		}

		// update Top Committees Overall by funds transferred
		err = updateTopCmteDonors(year, sender.(*donations.Committee))
		if err != nil {
			fmt.Println("CmteContUpdate failed: ", err)
			return fmt.Errorf("CmteContUpdate failed: %v", err)
		}

		// update Top Committees Overall by funds raised
		err = updateTopCmteRecs(year, receiver.(*donations.Committee))
		if err != nil {
			fmt.Println("CmteContUpdate failed: ", err)
			return fmt.Errorf("CmteContUpdate failed: %v", err)
		}

		// update linked Candidate, if any
		if receiver.(*donations.Committee).CandID != "" {
			cand, err := persist.GetObject(year, "candidates", receiver.(*donations.Committee).CandID)
			if err != nil {
				fmt.Println("CmteContUpdate failed: ", err)
				return fmt.Errorf("CmteContUpdate failed: %v", err)
			}

			err = updateCandCmteDonations(cont, sender.(*donations.Committee), cand.(*donations.Candidate))
			if err != nil {
				fmt.Println("CmteContUpdate failed: ", err)
				return fmt.Errorf("CmteContUpdate failed: %v", err)
			}

			// update Top Candidates Overall by funds raised
			err = updateTopCandidates(year, cand.(*donations.Candidate))
			if err != nil {
				fmt.Println("CmteContUpdate failed: ", err)
				return fmt.Errorf("CmteContUpdate failed: %v", err)
			}

			err = persist.PutObject(year, cand)
			if err != nil {
				fmt.Println("CmteContUpdate failed: ", err)
				return fmt.Errorf("CmteContUpdate failed: %v", err)
			}
		}

		// persist updated cmte objects
		err = persist.PutObject(year, receiver)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}

		// account for non-cmte senders
		err = persist.PutObject(year, sender)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}
	}

	return nil
}

func updateCmteDonations(cont *donations.CmteContribution, sen, rec *donations.Committee) error {
	// update receiver's donations values
	rec.TotalReceived += cont.TxAmt
	rec.TotalDonations++
	rec.AvgDonation = rec.TotalReceived / rec.TotalDonations

	// update senders' transfer values
	sen.TotalTransferred += cont.TxAmt
	sen.TotalTransfers++
	sen.AvgTransfer = sen.TotalTransferred / sen.TotalTransfers

	// update senders's affilates values
	if len(sen.AffiliatesAmt) == 0 {
		sen.AffiliatesAmt = make(map[string]float32)
		sen.AffiliatesTxs = make(map[string]float32)
	}
	sen.AffiliatesAmt[rec.ID] += cont.TxAmt
	sen.AffiliatesTxs[rec.ID]++

	// update receiver's TopCmteDonors values
	if len(rec.TopCmteDonorsAmt) == 0 {
		rec.TopCmteDonorsAmt = make(map[string]float32)
		rec.TopCmteDonorsTxs = make(map[string]float32)
	}
	if len(rec.TopCmteDonorsAmt) < 1000 {
		rec.TopCmteDonorsAmt[sen.ID] = sen.AffiliatesAmt[rec.ID]
		rec.TopCmteDonorsTxs[sen.ID]++
	} else {
		err := updateTopCmteTotals(sen, rec)
		if err != nil {
			fmt.Println("updateTopCmteTotals failed: ", err)
			return fmt.Errorf("updateTopCmteTotals failed: %v", err)
		}
	}

	return nil
}

func updateIndvDonations(cont *donations.IndvContribution, sen *donations.Individual, rec *donations.Committee) error {
	// update receiver's donations values
	rec.TotalReceived += cont.TxAmt
	rec.TotalDonations++
	rec.AvgDonation = rec.TotalReceived / rec.TotalDonations

	// add TxID to sender's list of donations
	sen.Donations = append(sen.Donations, cont.TxID)

	// update senders' transfer values
	sen.TotalDonated += cont.TxAmt
	sen.TotalDonations++
	sen.AvgDonation = sen.TotalDonated / sen.TotalDonations

	// update senders's affilates values
	if len(sen.RecipientsAmt) == 0 { // reinitialze if nil map
		sen.RecipientsAmt = make(map[string]float32)
		sen.RecipientsTxs = make(map[string]float32)
	}
	sen.RecipientsAmt[rec.ID] += cont.TxAmt
	sen.RecipientsTxs[rec.ID]++

	// update receiver's TopCmteDonors values
	if len(rec.TopIndvDonorThreshold) == 0 {
		rec.TopIndvDonorsAmt = make(map[string]float32)
		rec.TopIndvDonorsTxs = make(map[string]float32)
	}
	if len(rec.TopIndvDonorsAmt) < 1000 {
		rec.TopIndvDonorsAmt[sen.ID] = sen.RecipientsAmt[rec.ID]
		rec.TopIndvDonorsTxs[sen.ID]++
	} else {
		err := updateTopIndvTotals(sen, rec)
		if err != nil {
			fmt.Println("updateTopIndvTotals failed: ", err)
			return fmt.Errorf("updateTopIndvTotals failed: %v", err)
		}
	}

	return nil
}

func updateDisbursement(disb *donations.Disbursement, sen *donations.Committee, rec *donations.DisbRecipient) error {
	// update receiver's donations values
	rec.TotalReceived += disb.TxAmt
	rec.TotalDisbursements++
	rec.AvgReceived = rec.TotalReceived / float32(rec.TotalDisbursements)

	// add TxID to recipients list of disbursements
	rec.Disbursements = append(rec.Disbursements, disb.TxID)

	// update receiver's senders values
	if len(rec.SendersAmt) == 0 {
		rec.SendersAmt = make(map[string]float32)
		rec.SendersTxs = make(map[string]float32)
	}
	rec.SendersAmt[sen.ID] += disb.TxAmt
	rec.SendersTxs[sen.ID]++

	// update senders' disbursement values
	sen.TotalDisbursed += disb.TxAmt
	sen.TotalDisbursements++
	sen.AvgDisbursed = sen.TotalDisbursed / sen.TotalDisbursements

	// update senders's TopDisbRecipients values
	if len(sen.TopDisbRecipientsAmt) == 0 {
		sen.TopDisbRecipientsAmt = make(map[string]float32)
		sen.TopDisbRecipientsTxs = make(map[string]float32)
	}
	if len(sen.TopDisbRecipientsAmt) < 1000 {
		sen.TopDisbRecipientsAmt[rec.ID] = rec.SendersAmt[sen.ID]
		sen.TopDisbRecipientsTxs[rec.ID]++
	} else {
		err := updateTopDisbRecTotals(sen, rec)
		if err != nil {
			fmt.Println("updateTopDisbRecTotals failed: ", err)
			return fmt.Errorf("updateTopDisbRecTotals failed: %v", err)
		}
	}

	return nil
}

func updateCandCmteDonations(cont *donations.CmteContribution, sen *donations.Committee, cand *donations.Candidate) error {
	// update receiver's donations values
	cand.TotalRaised += cont.TxAmt
	cand.TotalDonations++
	cand.AvgDonation = cand.TotalRaised / cand.TotalDonations

	// update receiver's TopCmteDonors values
	if len(cand.TopCmteDonorsAmt) == 0 {
		cand.TopCmteDonorsAmt = make(map[string]float32)
		cand.TopCmteDonorsTxs = make(map[string]float32)
	}
	if len(cand.TopCmteDonorsAmt) < 1000 {
		cand.TopCmteDonorsAmt[sen.ID] = sen.AffiliatesAmt[cand.ID]
		cand.TopCmteDonorsTxs[sen.ID]++
	} else {
		err := updateCandTopCmteTotals(sen, cand)
		if err != nil {
			fmt.Println("updateCandTopCmteTotals failed: ", err)
			return fmt.Errorf("updateCandTopCmteTotals failed: %v", err)
		}
	}

	return nil
}

func updateCandIndvDonations(cont *donations.IndvContribution, sen *donations.Individual, cand *donations.Candidate) error {
	// update receiver's donations values
	cand.TotalRaised += cont.TxAmt
	cand.TotalDonations++
	cand.AvgDonation = cand.TotalRaised / cand.TotalDonations

	// update receiver's TopCmteDonors values
	if len(cand.TopIndvDonorsAmt) == 0 {
		cand.TopIndvDonorsAmt = make(map[string]float32)
		cand.TopIndvDonorsTxs = make(map[string]float32)
	}
	if len(cand.TopIndvDonorsAmt) < 1000 {
		cand.TopIndvDonorsAmt[sen.ID] = sen.RecipientsAmt[cand.ID]
		cand.TopIndvDonorsTxs[sen.ID]++
	} else {
		err := updateCandTopIndvTotals(sen, cand)
		if err != nil {
			fmt.Println("updateCandTopIndvTotals failed: ", err)
			return fmt.Errorf("updateCandTopIndvTotals failed: %v", err)
		}
	}

	return nil
}

func updateCandDisbursements(cont *donations.Disbursement, rec *donations.DisbRecipient, cand *donations.Candidate) error {
	// update receiver's donations values
	cand.TotalDisbursed += cont.TxAmt
	cand.TotalDisbursements++
	cand.AvgDisbursement = cand.TotalDisbursed / cand.TotalDisbursements

	// update receiver's TopCmteDonors values
	if len(cand.TopDisbRecsAmt) == 0 {
		cand.TopDisbRecsAmt = make(map[string]float32)
		cand.TopDisbRecsTxs = make(map[string]float32)
	}
	if len(cand.TopIndvDonorsAmt) < 1000 {
		cand.TopIndvDonorsAmt[rec.ID] += cont.TxAmt
		cand.TopIndvDonorsTxs[rec.ID]++
	} else {
		err := updateCandTopDisbRecTotals(cand, rec)
		if err != nil {
			fmt.Println("updateCandTopIndvTotals failed: ", err)
			return fmt.Errorf("updateCandTopIndvTotals failed: %v", err)
		}
	}

	return nil
}
*/
