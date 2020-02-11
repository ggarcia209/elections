package databuilder

import (
	"fmt"
	"strings"

	"github.com/elections/donations"
	"github.com/elections/persist"
)

// ContributionUpdate updates the filing committee's transaction values
// and the transaction values of the contributing Individual/Organization if applicable.
// Committees identified in the contribution's OtherID field are not updated.

func TransactionUpdate(year string, txs []interface{}) error {
	_, cont := txs.([]*donations.Contribution)
	if cont { // tx type is standard contribution/disbursement type
		err := contributionUpdate(year, txs.([]*donations.Contribution))
		if err != nil {
			fmt.Println("TransactionUpdate failed: ", err)
			return fmt.Errorf("TransactionUpdate failed: %v", err)
		}
	} else { // tx type is operating expense disbursement
		err := opExpensesUpdate(year, txs.([]*donations.Disbursement))
		if err != nil {
			fmt.Println("TransactionUpdate failed: ", err)
			return fmt.Errorf("TransactionUpdate failed: %v", err)
		}
	}
	return nil
}

// update data from Contribution transactiosn derived from contribution files
func contributionUpdate(year string, conts []*donations.Contribution) error {
	for _, cont := range conts {
		// get tx type info
		bucket, incoming, transfer, memo := deriveTxTypes(cont)

		// get filer object
		filer, err := persist.GetObject(year, "cmte_tx_data", cont.CmteID)
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}

		// get sender/receiver objects
		other, err := persist.GetObject(year, bucket, cont.OtherID)
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}

		// update incoming/outgoing tx data
		if incoming {
			err := incomingTxUpdate(cont, filer.(*donations.CmteTxData), other, transfer, memo)
			if err != nil {
				fmt.Println("ContributionUpdate failed: ", err)
				return fmt.Errorf("ContributionUpdate failed: %v", err)
			}
		} else {
			err := outgoingTxUpdate(cont, filer.(*donations.CmteTxData), other, transfer, memo)
			if err != nil {
				fmt.Println("ContributionUpdate failed: ", err)
				return fmt.Errorf("ContributionUpdate failed: %v", err)
			}
		}

		// update TopOverall rankings if not memo transaction
		if !memo {
			// update top individuals, organizations and committees
			err := updateTopOverall(year, filer.(*donations.CmteTxData), other, incoming, transfer)
			if err != nil {
				fmt.Println("ContributionUpdate failed: ", err)
				return fmt.Errorf("ContributionUpdate failed: %v", err)
			}
			// update top candidates by funds received if candidate linked to filing committee
			if filer.(*donations.CmteTxData).CandID != "" {
				// get linked candidate
				cand, err := persist.GetObject(year, "candidates", filer.(*donations.CmteTxData).CandID)
				if err != nil {
					fmt.Println("ContributionUpdate failed: ", err)
					return fmt.Errorf("ContributionUpdate failed: %v", err)
				}
				// update top candidates by total funds incoming/outgoing
				err = updateTopCandidates(year, cand.(*donations.Candidate), filer.(*donations.CmteTxData), incoming)
				if err != nil {
					fmt.Println("ContributionUpdate failed: ", err)
					return fmt.Errorf("ContributionUpdate failed: %v", err)
				}
			}
		}

		// persist updated cmte objects
		err = persist.PutObject(year, filer)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}

		// persist indv donor objects
		err = persist.PutObject(year, other)
		if err != nil {
			fmt.Println("UpdateCmte failed: ", err)
			return fmt.Errorf("UpdateCmte failed: %v", err)
		}
	}

	return nil
}

// update data from Disbursement transactions derived from operating expenses files
func opExpensesUpdate(year string, disbs []*donations.Disbursement) error {
	for _, disb := range disbs {
		// get  filing committee
		filer, err := persist.GetObject(year, "cmte_tx_data", disb.CmteID)
		if err != nil {
			fmt.Println("OpExpensesUpdate failed: ", err)
			return fmt.Errorf("OpExpensesUpdate failed: %v", err)
		}
		// get receiving organization
		receiver, err := persist.GetObject(year, "organizations", disb.RecID)
		if err != nil {
			fmt.Println("OpExpensesUpdate failed: ", err)
			return fmt.Errorf("OpExpensesUpdate failed: %v", err)
		}

		// update object account totals
		err = disbursementTxUpdate(disb, filer.(*donations.CmteTxData), receiver.(*donations.Organization))
		if err != nil {
			fmt.Println("OpExpensesUpdate failed: ", err)
			return fmt.Errorf("OpExpensesUpdate failed: %v", err)
		}
	}
	return nil
}

/*
	TX TYPES
	Incoming - funds received by filing committee
	Outgoing - funds disbursed by filing committee
	Transfer - contribution or loan to committee or candidate
	Memo - funds not credited/debited to receiver/senders' account totals
	but are credited/debited to receivers/senders maps only.
		- account for share of funds given to joint-fundraising commmittees and received by filing committee
		- amounts are not included in filing committee's and sender's itemization totals
		- itemization totals are accounted for by separate transactions between JFC and sender
*/

// derive transaction type from contribution data
func deriveTxTypes(cont *donations.Contribution) (string, bool, bool, bool) {
	// initialize return values
	incoming := false
	memo := false
	transfer := false               // indicates transfer (true) vs. expense (false)
	transferMap := map[string]bool{ // transfer tx codes
		"15Z": true,
		"16R": true,
		"18G": true,
		"18J": true,
		"18K": true,
		"19J": true,
		"22H": true,
		"24G": true,
		"24H": true,
		"24K": true,
		"24U": true,
		"24Z": true,
		"24I": true, // verify
		"24T": true, // verify
		"30K": true,
		"30G": true,
		"30F": true,
		"31K": true,
		"31G": true,
		"31F": true,
		"32K": true,
		"32G": true,
		"32F": true,
	}
	var bucket string

	// determine tx type (incoming / outgoing / transfer, memo)
	numCode := cont.TxType
	if numCode < "20" || (numCode >= "30" && numCode < "33") {
		incoming = true
	}

	// determine if transfer or expense
	if transferMap[cont.TxType] == true {
		transfer = true
	}

	if cont.MemoCode == "X" {
		memo = true
	}

	// determine contributor/receiver type - derive from OtherID
	IDss := strings.Split(cont.OtherID, "")
	idCode := IDss[0]
	switch {
	case idCode == "C":
		bucket = "cmte_tx_data"
	case idCode == "H" || idCode == "S" || idCode == "P":
		bucket = "candidates"
	default:
		if cont.Occupation == "" {
			bucket = "organizations"
		} else {
			bucket = "individuals"
		}
	}
	return bucket, incoming, transfer, memo
}

/*
	TX CRITERIA
	If sender/receiver (OtherID) == committee: do not debit/credit other's accounts -- accounted for  by corresponding tx.
	All filers are committees only and *should* have corresponding transaction filed by sender/receiver.
	if sender/receiver != committee: credit/debit accounts -- no corresponding tx.
	Receiver's of transactions from operating expense files are always treated as organizations.
		- individual recipients are treated as sole-propritor or single member llc business entities
		  rendering services or other material contributions for which they are paid by committees.
*/

// update filing committee and sender object data for incoming transactions
func incomingTxUpdate(cont *donations.Contribution, filerData *donations.CmteTxData, sender interface{}, transfer, memo bool) error {
	// update maps only if memo == true
	// account for percentage of amounts received by memo transactions but do not add to totals
	if memo {
		// update maps
		err := mapUpdateIncoming(cont, filerData, sender, transfer)
		if err != nil {
			fmt.Println("incomingTxUpdate failed: ", err)
			return fmt.Errorf("incomingTxUpdate failed: %v", err)
		}
		return nil
	}

	// credit Contributions or OtherReceipts and TotalIncoming
	if !transfer {
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

	// debit sender account
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
		_ = t // discard unused variable
		return fmt.Errorf("incomingTxUpdate failed: wrong interface type")
	}

	// update maps
	err := mapUpdateIncoming(cont, filerData, sender, transfer)
	if err != nil {
		fmt.Println("incomingTxUpdate failed: ", err)
		return fmt.Errorf("incomingTxUpdate failed: %v", err)
	}
	return nil
}

// update filing committe and receiver object data for outgoing transactions
func outgoingTxUpdate(cont *donations.Contribution, filerData *donations.CmteTxData, receiver interface{}, transfer, memo bool) error {
	// update maps only if memo == true
	// account for percentage of amounts received by memo transactions but do not add to totals
	if memo {
		// update maps
		err := mapUpdateOutgoing(cont, filerData, receiver, transfer)
		if err != nil {
			fmt.Println("outgoingTxUpdate failed: ", err)
			return fmt.Errorf("outgoingTxUpdate failed: %v", err)
		}
		return nil
	}

	// credit Transfers or Expenditures and TotalOutgoing
	if transfer { // tx is transfer
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

	// debit receiver accounts
	switch t := receiver.(type) {
	case *donations.Individual:
		receiver.(*donations.Individual).TotalInAmt += cont.TxAmt
		receiver.(*donations.Individual).TotalInTxs++
		receiver.(*donations.Individual).AvgTxIn = receiver.(*donations.Individual).TotalInAmt / receiver.(*donations.Individual).TotalInTxs
		receiver.(*donations.Individual).NetBalance = receiver.(*donations.Individual).TotalInAmt - receiver.(*donations.Individual).TotalOutAmt
	case *donations.Organization:
		receiver.(*donations.Organization).TotalInAmt += cont.TxAmt
		receiver.(*donations.Organization).TotalInTxs++
		receiver.(*donations.Organization).AvgTxIn = receiver.(*donations.Organization).TotalInAmt / receiver.(*donations.Organization).TotalInTxs
		receiver.(*donations.Organization).NetBalance = receiver.(*donations.Organization).TotalInAmt - receiver.(*donations.Organization).TotalOutAmt
	case *donations.Candidate:
		receiver.(*donations.Candidate).TotalDirectInAmt += cont.TxAmt
		receiver.(*donations.Candidate).TotalDirectInTxs++
		receiver.(*donations.Candidate).AvgDirectIn = receiver.(*donations.Candidate).TotalDirectInAmt / receiver.(*donations.Candidate).TotalDirectInTxs
		receiver.(*donations.Candidate).NetBalanceDirectTx = receiver.(*donations.Candidate).TotalDirectInAmt - receiver.(*donations.Candidate).TotalDirectOutAmt
	case *donations.CmteTxData:
		// special case -- pass sender as filer (receiving committee)
		// and filer as sender (sending committee) to mapUpdate()
		err := mapUpdateOutgoing(cont, receiver.(*donations.CmteTxData), filerData, transfer)
		if err != nil {
			fmt.Println("outgoingTxUpdate failed: ", err)
			return fmt.Errorf("outgoingTxUpdate failed: %v", err)
		}
		return nil
	default:
		_ = t // discard unused variable
		return fmt.Errorf("outgoingTxUpdate failed: wrong interface type")
	}

	// update maps
	err := mapUpdateOutgoing(cont, filerData, receiver, transfer)
	if err != nil {
		fmt.Println("outgoingTxUpdate failed: ", err)
		return fmt.Errorf("outgoingTxUpdate failed: %v", err)
	}
	return nil
}

func disbursementTxUpdate(disb *donations.Disbursement, filer *donations.CmteTxData, receiver *donations.Organization) error {
	// debit filer's expense account
	filer.ExpendituresAmt += disb.TxAmt
	filer.ExpendituresTxs++
	filer.AvgExpenditure = filer.ExpendituresAmt / filer.ExpendituresTxs
	filer.TotalOutgoingAmt = filer.TransfersAmt + filer.ExpendituresAmt
	filer.TotalOutgoingTxs = filer.TransfersTxs + filer.ExpendituresTxs
	filer.AvgOutgoing = filer.TotalOutgoingAmt / filer.TotalOutgoingTxs
	filer.NetBalance = filer.TotalIncomingAmt - filer.TotalOutgoingAmt

	// credit receiver's accounts
	receiver.TotalInAmt += disb.TxAmt
	receiver.TotalInTxs++
	receiver.AvgTxIn = receiver.TotalInAmt / receiver.TotalInTxs
	receiver.NetBalance = receiver.TotalInAmt - receiver.TotalOutAmt

	// update filer's expense recipient maps and receiving org's sender's maps
	err := mapUpdateOpExp(disb, filer, receiver)
	if err != nil {
		fmt.Println("disbursementTxUpdate failed: ", err)
		return fmt.Errorf("disbursementTxUpdate failed: %v", err)
	}
	return nil
}

/*
	MAP UPDATE CRITERIA
	Maps are updated for every transaction.
	'transfer' variable determines which committee maps are updated.
	Receiver's of transactions from operating expense files are always treated as organizations.
		- individual recipients are treated as sole-propritor or single member llc business entities
		  rendering services or other material contributions for which they are paid by committees.
*/

// update maps for incoming transactions posted by filing committee
func mapUpdateIncoming(cont *donations.Contribution, filerData *donations.CmteTxData, sender interface{}, transfer bool) error {
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
		if len(filerData.TopIndvContributorsAmt) < 100 || filerData.TopIndvContributorsAmt[sender.(*donations.Individual).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopIndvContributorsAmt[sender.(*donations.Individual).ID] += cont.TxAmt
			filerData.TopIndvContributorsTxs[sender.(*donations.Individual).ID]++
		} else {
			// update filer's top contributors maps by comparison
			comp, err := updateTopDonors(filerData, sender, cont, transfer)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
			filerData.TopIndvContributorsAmt = comp.RefAmts
			filerData.TopIndvContributorsTxs = comp.RefTxs
			filerData.TopIndvContributorThreshold = comp.RefThreshold
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
		if len(filerData.TopCmteOrgContributorsAmt) < 100 || filerData.TopCmteOrgContributorsAmt[sender.(*donations.Organization).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopCmteOrgContributorsAmt[sender.(*donations.Organization).ID] += cont.TxAmt
			filerData.TopCmteOrgContributorsTxs[sender.(*donations.Organization).ID]++
		} else {
			// update by comparison
			comp, err := updateTopDonors(filerData, sender, cont, transfer)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
			filerData.TopCmteOrgContributorsAmt = comp.RefAmts
			filerData.TopCmteOrgContributorsTxs = comp.RefTxs
			filerData.TopCmteOrgContributorThreshold = comp.RefThreshold
		}
	case *donations.CmteTxData:
		// re-initialize maps if nil
		if len(filerData.TopCmteOrgContributorsAmt) == 0 {
			filerData.TopCmteOrgContributorsAmt = make(map[string]float32)
			filerData.TopCmteOrgContributorsTxs = make(map[string]float32)
		}
		if len(sender.(*donations.CmteTxData).TransferRecsAmt) == 0 && transfer {
			sender.(*donations.CmteTxData).TransferRecsAmt = make(map[string]float32)
			sender.(*donations.CmteTxData).TransferRecsTxs = make(map[string]float32)
		}
		if len(sender.(*donations.CmteTxData).TopExpRecipientsAmt) == 0 && !transfer {
			sender.(*donations.CmteTxData).TopExpRecipientsAmt = make(map[string]float32)
			sender.(*donations.CmteTxData).TopExpRecipientsTxs = make(map[string]float32)
		}

		// update filing committee's Top Contributors maps only --  sending committee's maps updated in corresponding tx
		if len(filerData.TopCmteOrgContributorsAmt) < 100 || filerData.TopCmteOrgContributorsAmt[sender.(*donations.CmteTxData).CmteID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopCmteOrgContributorsAmt[sender.(*donations.CmteTxData).CmteID] += cont.TxAmt
			filerData.TopCmteOrgContributorsTxs[sender.(*donations.CmteTxData).CmteID]++
		} else {
			// update by comparison
			comp, err := updateTopDonors(filerData, sender, cont, transfer)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
			filerData.TopCmteOrgContributorsAmt = comp.RefAmts
			filerData.TopCmteOrgContributorsTxs = comp.RefTxs
			filerData.TopCmteOrgContributorThreshold = comp.RefThreshold
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
		if len(filerData.TopIndvContributorsAmt) < 100 || filerData.TopIndvContributorsAmt[sender.(*donations.Candidate).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopIndvContributorsAmt[sender.(*donations.Candidate).ID] += cont.TxAmt
			filerData.TopIndvContributorsTxs[sender.(*donations.Candidate).ID]++
		} else {
			// update by comparison
			comp, err := updateTopDonors(filerData, sender, cont, transfer)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
			filerData.TopIndvContributorsAmt = comp.RefAmts
			filerData.TopIndvContributorsTxs = comp.RefTxs
			filerData.TopIndvContributorThreshold = comp.RefThreshold
		}
	default:
		_ = t // discard unused variable
		return fmt.Errorf("mapUpdate failed: wrong interface type")
	}
	return nil
}

// update maps for outgoing transaction posted by filing committee - transactions to individuals
// or organizations are always counted as expenses (loan repayments, refunds, independent expendiures, etc..)
func mapUpdateOutgoing(cont *donations.Contribution, filerData *donations.CmteTxData, receiver interface{}, transfer bool) error {
	switch t := receiver.(type) {
	case *donations.Individual:
		// re-initialize maps if nil
		if len(filerData.TopExpRecipientsAmt) == 0 {
			filerData.TopExpRecipientsAmt = make(map[string]float32)
			filerData.TopExpRecipientsTxs = make(map[string]float32)
		}
		if len(receiver.(*donations.Individual).SendersAmt) == 0 {
			receiver.(*donations.Individual).SendersAmt = make(map[string]float32)
			receiver.(*donations.Individual).SendersTxs = make(map[string]float32)
		}

		// update receiver's Senders maps
		receiver.(*donations.Individual).SendersAmt[filerData.CmteID] += cont.TxAmt
		receiver.(*donations.Individual).SendersTxs[filerData.CmteID]++

		// update filing committee's Top Expense Recipients maps
		if len(filerData.TopExpRecipientsAmt) < 100 || filerData.TopExpRecipientsAmt[receiver.(*donations.Individual).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopExpRecipientsAmt[receiver.(*donations.Individual).ID] += cont.TxAmt
			filerData.TopExpRecipientsTxs[receiver.(*donations.Individual).ID]++
		} else {
			// update by comparison
			comp, err := updateTopRecipients(filerData, receiver, cont, transfer)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
			filerData.TopExpRecipientsAmt = comp.RefAmts
			filerData.TopExpRecipientsTxs = comp.RefTxs
			filerData.TopExpThreshold = comp.RefThreshold
		}
	case *donations.Organization:
		// re-initialize maps if nil
		if len(filerData.TopExpRecipientsAmt) == 0 {
			filerData.TopExpRecipientsAmt = make(map[string]float32)
			filerData.TopExpRecipientsTxs = make(map[string]float32)
		}
		if len(receiver.(*donations.Organization).SendersAmt) == 0 {
			receiver.(*donations.Organization).SendersAmt = make(map[string]float32)
			receiver.(*donations.Organization).SendersTxs = make(map[string]float32)
		}

		// update sender's Recipients maps
		receiver.(*donations.Organization).SendersAmt[filerData.CmteID] += cont.TxAmt
		receiver.(*donations.Organization).SendersTxs[filerData.CmteID]++

		// update filing committee's Top Contributors maps
		if len(filerData.TopExpRecipientsAmt) < 100 || filerData.TopExpRecipientsAmt[receiver.(*donations.Organization).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopExpRecipientsAmt[receiver.(*donations.Organization).ID] += cont.TxAmt
			filerData.TopExpRecipientsTxs[receiver.(*donations.Organization).ID]++
		} else {
			// update by comparison
			comp, err := updateTopRecipients(filerData, receiver, cont, transfer)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
			filerData.TopExpRecipientsAmt = comp.RefAmts
			filerData.TopExpRecipientsTxs = comp.RefTxs
			filerData.TopExpThreshold = comp.RefThreshold
		}
	case *donations.CmteTxData:
		// re-initialize maps if nil
		if len(filerData.TransferRecsAmt) == 0 && transfer {
			filerData.TransferRecsAmt = make(map[string]float32)
			filerData.TransferRecsTxs = make(map[string]float32)
		}
		if len(filerData.TopExpRecipientsAmt) == 0 && !transfer {
			filerData.TopExpRecipientsAmt = make(map[string]float32)
			filerData.TopExpRecipientsTxs = make(map[string]float32)
		}
		if len(receiver.(*donations.CmteTxData).TopCmteOrgContributorsAmt) == 0 {
			receiver.(*donations.CmteTxData).TopCmteOrgContributorsAmt = make(map[string]float32)
			receiver.(*donations.CmteTxData).TopCmteOrgContributorsTxs = make(map[string]float32)
		}

		// update filer's TransfersRecs or TopExpRecipients maps only -- receiving committees maps updated in corresponding tx
		if transfer {
			filerData.TransferRecsAmt[receiver.(*donations.CmteTxData).CmteID] += cont.TxAmt
			filerData.TransferRecsTxs[receiver.(*donations.CmteTxData).CmteID]++
		} else {
			if len(filerData.TopExpRecipientsAmt) < 100 || filerData.TopExpRecipientsAmt[receiver.(*donations.CmteTxData).CmteID] != 0 {
				// add new entry directly or update existing entry
				filerData.TopExpRecipientsAmt[receiver.(*donations.CmteTxData).CmteID] += cont.TxAmt
				filerData.TopExpRecipientsTxs[receiver.(*donations.CmteTxData).CmteID]++
			} else {
				// update by comparison
				comp, err := updateTopRecipients(filerData, receiver, cont, transfer)
				if err != nil {
					fmt.Println("mapUpdate failed: ", err)
					return fmt.Errorf("mapUpdate failed: %v", err)
				}
				filerData.TopExpRecipientsAmt = comp.RefAmts
				filerData.TopExpRecipientsTxs = comp.RefTxs
				filerData.TopExpThreshold = comp.RefThreshold
			}
		}
	case *donations.Candidate:
		// re-initialize maps if nil
		if len(filerData.TransferRecsAmt) == 0 && transfer {
			filerData.TransferRecsAmt = make(map[string]float32)
			filerData.TransferRecsTxs = make(map[string]float32)
		}
		if len(filerData.TopExpRecipientsAmt) == 0 && !transfer {
			filerData.TopExpRecipientsAmt = make(map[string]float32)
			filerData.TopExpRecipientsTxs = make(map[string]float32)
		}
		if len(receiver.(*donations.Candidate).DirectSendersAmts) == 0 {
			receiver.(*donations.Candidate).DirectSendersAmts = make(map[string]float32)
			receiver.(*donations.Candidate).DirectSendersTxs = make(map[string]float32)
		}

		// update receiver's senders maps
		receiver.(*donations.Candidate).DirectSendersAmts[filerData.CmteID] += cont.TxAmt
		receiver.(*donations.Candidate).DirectSendersTxs[filerData.CmteID]++

		// update filing committee's transfers or top expenditures maps
		if transfer {
			filerData.TransferRecsAmt[receiver.(*donations.Candidate).ID] += cont.TxAmt
			filerData.TransferRecsTxs[receiver.(*donations.Candidate).ID]++
		} else {
			if len(filerData.TopExpRecipientsAmt) < 100 || filerData.TopExpRecipientsAmt[receiver.(*donations.Candidate).ID] != 0 {
				// add new entry directly or update existing entry
				filerData.TopExpRecipientsAmt[receiver.(*donations.Candidate).ID] += cont.TxAmt
				filerData.TopExpRecipientsTxs[receiver.(*donations.Candidate).ID]++
			} else {
				// update by comparison
				comp, err := updateTopRecipients(filerData, receiver, cont, transfer)
				if err != nil {
					fmt.Println("mapUpdate failed: ", err)
					return fmt.Errorf("mapUpdate failed: %v", err)
				}
				filerData.TopExpRecipientsAmt = comp.RefAmts
				filerData.TopExpRecipientsTxs = comp.RefTxs
				filerData.TopExpThreshold = comp.RefThreshold
			}
		}
	default:
		_ = t // discard unused variable
		return fmt.Errorf("mapUpdate failed: wrong interface type")
	}
	return nil
}

// update maps for expenditures listed in operating expenses files
func mapUpdateOpExp(disb *donations.Disbursement, filer *donations.CmteTxData, receiver *donations.Organization) error {
	// re-initialize maps if nil
	if len(filer.TopExpRecipientsAmt) == 0 {
		filer.TopExpRecipientsAmt = make(map[string]float32)
		filer.TopExpRecipientsTxs = make(map[string]float32)
	}
	if len(receiver.SendersAmt) == 0 {
		receiver.SendersAmt = make(map[string]float32)
		receiver.SendersTxs = make(map[string]float32)
	}

	// update filer's top expenditure recipient maps
	if len(filer.TopExpRecipientsAmt) < 100 || filer.TopExpRecipientsAmt[receiver.ID] != 0 {
		// update by direct add
		filer.TopExpRecipientsAmt[receiver.ID] += disb.TxAmt
		filer.TopExpRecipientsTxs[receiver.ID]++
	} else {
		// update by comparison
		comp, err := updateOpExpRecipients(filer, receiver)
		if err != nil {
			fmt.Println("mapUpdateOpExp failed: ", err)
			return fmt.Errorf("mapUpdateOpExp failed: %v", err)
		}
		filer.TopExpRecipientsAmt = comp.RefAmts
		filer.TopExpRecipientsTxs = comp.RefTxs
		filer.TopExpThreshold = comp.RefThreshold
	}

	// update receiver's sender's maps
	receiver.SendersAmt[filer.CmteID] += disb.TxAmt
	receiver.SendersTxs[filer.CmteID]++

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
