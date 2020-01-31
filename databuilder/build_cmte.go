package databuilder

import (
	"fmt"
	"strings"

	"github.com/elections/donations"
	"github.com/elections/persist"
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
func ContributionUpdate(year string, conts []*donations.Contribution) error {
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
}

func outgoingTxUpdate(cont *donations.Contribution, filerData *donations.CmteTxData, sender interface{}) {
	// corresponding txs - 15Z, 18G, 18K, 30K, 30G, 31K, 31G, 32K, 32G
	// credit Contributions or OtherReceipts and TotalIncoming
	if cont.TxType > 15 && cont.TxType < 18 { // correct to transfers/contributions only
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
	filerData.AvgOutgoing = filerData.TotalOutgoingAmt / filerData.TotalOutgingTxs
	filerData.NetBalance = filerData.TotalIncomingAmt - filerData.TotalOutgoingAmt

	// debit sender account
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
		sender.(*donations.Candidate).NetBalance = sender.(*donations.Candidate).TotalDirectInAmt - sender.(*donations.Candidate).TotalDirectOutAmt
	}

	// update maps
	mapUpdate(cont, filerData, sender)
}

func incomingTxUpdate(cont *donations.Contribution, filerData *donations.CmteTxData, sender interface{}) {
	// corresponding txs - 15Z, 18G, 18K, 30K, 30G, 31K, 31G, 32K, 32G
	// debit Contributions or OtherReceipts and TotalIncoming
	if cont.TxType > 15 && cont.TxType < 18 {
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
		sender.(*donations.Candidate).NetBalance = sender.(*donations.Candidate).TotalDirectInAmt - sender.(*donations.Candidate).TotalDirectOutAmt
	}

	// update maps
	mapUpdate(cont, filerData, sender)
}

func mapUpdate(cont *donations.Contribution, filerData *donations.CmteTxData, sender interface{}) {
	switch t := sender.(type) {
	case *donations.Individual:
		// re-initialize maps if nil
		if len(filerData.TopIndvContributorsAmt) == 0 {
			filerData.TopIndvContributorsAmt = make(map[string]float32)
			filerData.TopIndvContributorsTxs = make(map[string]float32)
		}
		if len(sender.RecipientsAmt) == 0 {
			sender.RecipientsAmt = make(map[string]float32)
			sender.RecipientsTxs = make(map[string]float32)
		}

		// update sender's Recipients maps
		sender.RecipientsAmt[filerData.CmteID] += cont.TxAmt
		sender.RecipientsTxs[filerData.CmteID]++

		// update filing committee's Top Contributors maps
		if len(filerData.TopIndvContributorsAmt) < 1000 {
			filerData.TopIndvContributorsAmt[sender.ID] += cont.TxAmt
			filerData.TopIndvContributorsTxs[sender.ID]++
		} else {
			// updateTop function
		}
	case *donations.Organization:
		// re-initialize maps if nil
		if len(filerData.TopCmteOrgContributorsAmt) == 0 {
			filerData.TopCmteOrgContributorsAmt = make(map[string]float32)
			filerData.TopCmteOrgContributorsTxs = make(map[string]float32)
		}
		if len(sender.RecipientsAmt) == 0 {
			sender.RecipientsAmt = make(map[string]float32)
			sender.RecipientsTxs = make(map[string]float32)
		}

		// update sender's Recipients maps
		sender.RecipientsAmt[filerData.CmteID] += cont.TxAmt
		sender.RecipientsTxs[filerData.CmteID]++

		// update filing committee's Top Contributors maps
		if len(filerData.TopCmteOrgContributorsAmt) < 1000 {
			filerData.TopCmteOrgContributorsAmt[sender.ID] += cont.TxAmt
			filerData.TopCmteOrgContributorsTxs[sender.ID]++
		} else {
			// updateTop function
		}
	case *donations.CmteTxData:
		// re-initialize maps if nil
		if len(filerData.TopCmteOrgContributorsAmt) == 0 {
			filerData.TopCmteOrgContributorsAmt = make(map[string]float32)
			filerData.TopCmteOrgContributorsTxs = make(map[string]float32)
		}
		if len(sender.AffiliatesAmt) == 0 {
			sender.AffiliatesAmt = make(map[string]float32)
			sender.AffiliatesTxs = make(map[string]float32)
		}

		// update sender's Affiliates maps
		sender.AffiliatesAmt[filerData.CmteID] += cont.TxAmt
		sender.AfiliatesTxs[filerData.CmteID]++

		// update filing committee's Top Contributors maps
		if len(filerData.TopCmteOrgContributorsAmt) < 1000 {
			filerData.TopCmteOrgContributorsAmt[sender.ID] += cont.TxAmt
			filerData.TopCmteOrgContributorsTxs[sender.ID]++
		} else {
			// updateTop function
		}
	case *donations.Candidate:
		// re-initialize maps if nil
		if len(filerData.TopIndvContributorsAmt) == 0 {
			filerData.TopIndvContributorsAmt = make(map[string]float32)
			filerData.TopIndvContributorsTxs = make(map[string]float32)
		}
		if len(sender.DirectRecipientsAmt) == 0 {
			sender.DirectRecipientsAmt = make(map[string]float32)
			sender.DirectRecipientsTxs = make(map[string]float32)
		}

		// update sender's Recipients maps
		sender.DirectRecipientsAmt[filerData.CmteID] += cont.TxAmt
		sender.DirectRecipientsTxs[filerData.CmteID]++

		// update filing committee's Top Contributors maps
		if len(filerData.TopIndvContributorsAmt) < 1000 {
			filerData.TopIndvContributorsAmt[sender.ID] += cont.TxAmt
			filerData.TopIndvContributorsTxs[sender.ID]++
		} else {
			// updateTop function
		}
	}
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

// DEPRECATED
/*

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
