package databuilder

// Refactored 8/3/20 - removed Organization cases and references

import (
	"fmt"
	"strings"

	"github.com/elections/source/donations"
)

// TransactionUpdate updates each sender/receiver data for each transaction in a list of transactions.
func TransactionUpdate(year string, txs interface{}, cache map[string]map[string]interface{}) error {
	_, cont := txs.([]*donations.Contribution)

	if cont { // tx type is standard contribution/disbursement type
		err := contributionUpdate(year, txs.([]*donations.Contribution), cache)
		if err != nil {
			fmt.Println("TransactionUpdate failed: ", err)
			return fmt.Errorf("TransactionUpdate failed: %v", err)
		}
	} else { // tx type is operating expense disbursement
		err := opExpensesUpdate(year, txs.([]*donations.Disbursement), cache)
		if err != nil {
			fmt.Println("TransactionUpdate failed: ", err)
			return fmt.Errorf("TransactionUpdate failed: %v", err)
		}
	}
	return nil
}

// update data from Contribution transactiosn derived from contribution files
func contributionUpdate(year string, conts []*donations.Contribution, cache map[string]map[string]interface{}) error {
	for _, cont := range conts {
		// get tx type info
		bucket, incoming, transfer, memo := deriveTxTypes(cont)

		// get filer object
		filer := cache["cmte_tx_data"][cont.CmteID]

		// get sender/receiver objects
		other := cache[bucket][cont.OtherID]

		// start := time.Now()

		// update incoming/outgoing tx data
		if incoming {
			err := incomingTxUpdate(cont, filer.(*donations.CmteTxData), other, transfer, memo)
			if err != nil {
				fmt.Println("contributionUpdate failed: ", err)
				return fmt.Errorf("contributionUpdate failed: %v", err)
			}
		} else {
			err := outgoingTxUpdate(cont, filer.(*donations.CmteTxData), other, transfer, memo)
			if err != nil {
				fmt.Println("contributionUpdate failed: ", err)
				return fmt.Errorf("contributionUpdate failed: %v", err)
			}
		}

		// update TopOverall rankings if not memo transaction
		if !memo {
			// update top individuals, organizations and committees
			err := updateTopOverall(filer.(*donations.CmteTxData), other, incoming, transfer, cache)
			if err != nil {
				fmt.Println("contributionUpdate failed: ", err)
				return fmt.Errorf("contributionUpdate failed: %v", err)
			}
		}
	}

	return nil
}

// update data from Disbursement transactions derived from operating expenses files
func opExpensesUpdate(year string, disbs []*donations.Disbursement, cache map[string]map[string]interface{}) error {
	for _, disb := range disbs {
		// get filing committee
		// get filer object
		filer := cache["cmte_tx_data"][disb.CmteID]

		// get sender/receiver objects
		receiver := cache["individuals"][disb.RecID]

		// update object account totals
		err := disbursementTxUpdate(disb, filer.(*donations.CmteTxData), receiver.(*donations.Individual))
		if err != nil {
			fmt.Println("opExpensesUpdate failed: ", err)
			return fmt.Errorf("opExpensesUpdate failed: %v", err)
		}

		// update TopOverall rankings
		// update top individuals, organizations and committees
		err = updateTopOverall(filer.(*donations.CmteTxData), receiver, false, false, cache)
		if err != nil {
			fmt.Println("opExpensesUpdate failed: ", err)
			return fmt.Errorf("opExpensesUpdate failed: %v", err)
		}

	}
	return nil
}

// derive transaction type from contribution data
func deriveTxTypes(cont *donations.Contribution) (string, bool, bool, bool) {
	// initialize return values
	incoming := false
	memo := false
	transfer := false               // indicates transfer (true) vs. expense (false)
	transferMap := map[string]bool{ // transfer tx codes
		"15Z": true,
		"16R": true,
		"17R": true,
		"18G": true,
		"18J": true,
		"18K": true,
		"19J": true,
		"22H": true,
		"22Z": true,
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
		"40Z": true,
		"41Z": true,
		"42Z": true,
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
		bucket = "individuals"
	}
	return bucket, incoming, transfer, memo
}

/*
	TX CRITERIA
	If sender/receiver (OtherID) == committee: do not debit/credit other's accounts -- accounted for  by corresponding tx.
	All filers are committees only and *should* have corresponding transaction filed by sender/receiver.
	if sender/receiver != committee: credit/debit accounts -- no corresponding tx.
	Receiver's of transactions from operating expense files are always treated as organizations.
		- individual recipients are treated as sole-proprietor or single member llc business entities
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
	if cont.TxType < "16" || cont.TxType > "18" {
		filerData.ContributionsInAmt += cont.TxAmt
		filerData.ContributionsInTxs++
		filerData.AvgContributionIn = filerData.ContributionsInAmt / filerData.ContributionsInTxs
	} else {
		filerData.OtherReceiptsInAmt += cont.TxAmt
		filerData.OtherReceiptsInTxs++
		filerData.AvgOtherIn = filerData.OtherReceiptsInAmt / filerData.OtherReceiptsInTxs
	}
	filerData.TotalIncomingAmt = filerData.ContributionsInAmt + filerData.OtherReceiptsInAmt
	filerData.TotalIncomingTxs = filerData.ContributionsInTxs + filerData.OtherReceiptsInTxs
	filerData.AvgIncoming = filerData.TotalIncomingAmt / filerData.TotalIncomingTxs
	filerData.NetBalance = filerData.TotalIncomingAmt - filerData.TotalOutgoingAmt

	// debit sender account
	switch t := sender.(type) {
	case *donations.Individual:
		sender.(*donations.Individual).Transactions = append(sender.(*donations.Individual).Transactions, cont.TxID)
		sender.(*donations.Individual).TotalOutAmt += cont.TxAmt
		sender.(*donations.Individual).TotalOutTxs++
		sender.(*donations.Individual).AvgTxOut = sender.(*donations.Individual).TotalOutAmt / sender.(*donations.Individual).TotalOutTxs
		sender.(*donations.Individual).NetBalance = sender.(*donations.Individual).TotalInAmt - sender.(*donations.Individual).TotalOutAmt
	case *donations.Candidate:
		sender.(*donations.Candidate).TotalDirectOutAmt += cont.TxAmt
		sender.(*donations.Candidate).TotalDirectOutTxs++
		sender.(*donations.Candidate).AvgDirectOut = sender.(*donations.Candidate).TotalDirectOutAmt / sender.(*donations.Candidate).TotalDirectOutTxs
		sender.(*donations.Candidate).NetBalanceDirectTx = sender.(*donations.Candidate).TotalDirectInAmt - sender.(*donations.Candidate).TotalDirectOutAmt
	case *donations.CmteTxData:
		// do nothing -- accounted for by sender's corresponding outgoing transaction
	case nil:
		return fmt.Errorf("incomingTxUpdate failed: nil interface")
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
		filerData.AvgExpenditure = filerData.ExpendituresAmt / filerData.ExpendituresTxs
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
	case *donations.Candidate:
		receiver.(*donations.Candidate).TotalDirectInAmt += cont.TxAmt
		receiver.(*donations.Candidate).TotalDirectInTxs++
		receiver.(*donations.Candidate).AvgDirectIn = receiver.(*donations.Candidate).TotalDirectInAmt / receiver.(*donations.Candidate).TotalDirectInTxs
		receiver.(*donations.Candidate).NetBalanceDirectTx = receiver.(*donations.Candidate).TotalDirectInAmt - receiver.(*donations.Candidate).TotalDirectOutAmt
	case *donations.CmteTxData:
		// special case -- pass sender as filer (receiving committee)
		// and filer as sender (sending committee) to mapUpdate()
		err := mapUpdateOutgoing(cont, filerData, receiver, transfer)
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

func disbursementTxUpdate(disb *donations.Disbursement, filer *donations.CmteTxData, receiver *donations.Individual) error {
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
		- individual recipients are treated as sole-proprietor or single member llc business entities
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

			if len(filerData.TopIndvContributorsAmt) >= 100 {
				th, err := checkThreshold(sender.(*donations.Individual).ID, filerData.TopIndvContributorsAmt, filerData.TopIndvContributorThreshold)
				if err != nil {
					fmt.Println("mapUpdateIncoming failed: ", err)
					return fmt.Errorf("mapUpdateIncoming failed: %v", err)
				}
				filerData.TopIndvContributorThreshold = th
			}
		} else {
			// update filer's top contributors maps by comparison
			comp, err := updateTopDonors(filerData, sender, cont)
			if err != nil {
				fmt.Println("mapUpdate failed: ", err)
				return fmt.Errorf("mapUpdate failed: %v", err)
			}
			filerData.TopIndvContributorsAmt = comp.RefAmts
			filerData.TopIndvContributorsTxs = comp.RefTxs
			filerData.TopIndvContributorThreshold = comp.RefThreshold
		}
	case *donations.CmteTxData:
		// re-initialize maps if nil
		if len(filerData.TopCmteOrgContributorsAmt) == 0 {
			filerData.TopCmteOrgContributorsAmt = make(map[string]float32)
			filerData.TopCmteOrgContributorsTxs = make(map[string]float32)
		}
		if len(sender.(*donations.CmteTxData).TransferRecsAmt) == 0 {
			sender.(*donations.CmteTxData).TransferRecsAmt = make(map[string]float32)
			sender.(*donations.CmteTxData).TransferRecsTxs = make(map[string]float32)
		}

		// update filing committee's Top Contributors maps only --  sending committee's maps updated in corresponding tx
		if len(filerData.TopCmteOrgContributorsAmt) < 100 || filerData.TopCmteOrgContributorsAmt[sender.(*donations.CmteTxData).CmteID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopCmteOrgContributorsAmt[sender.(*donations.CmteTxData).CmteID] += cont.TxAmt
			filerData.TopCmteOrgContributorsTxs[sender.(*donations.CmteTxData).CmteID]++

			if len(filerData.TopCmteOrgContributorsAmt) >= 100 {
				th, err := checkThreshold(sender.(*donations.CmteTxData).CmteID, filerData.TopCmteOrgContributorsAmt, filerData.TopCmteOrgContributorThreshold)
				if err != nil {
					fmt.Println("mapUpdateIncoming failed: ", err)
					return fmt.Errorf("mapUpdateIncoming failed: %v", err)
				}
				filerData.TopCmteOrgContributorThreshold = th
			}
		} else {
			// update by comparison
			comp, err := updateTopDonors(filerData, sender, cont)
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
			if len(filerData.TopIndvContributorsAmt) >= 100 {
				th, err := checkThreshold(sender.(*donations.Candidate).ID, filerData.TopIndvContributorsAmt, filerData.TopIndvContributorThreshold)
				if err != nil {
					fmt.Println("mapUpdateIncoming failed: ", err)
					return fmt.Errorf("mapUpdateIncoming failed: %v", err)
				}
				filerData.TopIndvContributorThreshold = th
			}
		} else {
			// update by comparison
			comp, err := updateTopDonors(filerData, sender, cont)
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

			if len(filerData.TopExpRecipientsAmt) >= 100 {
				th, err := checkThreshold(receiver.(*donations.Individual).ID, filerData.TopExpRecipientsAmt, filerData.TopExpThreshold)
				if err != nil {
					fmt.Println("mapUpdateOutgoing failed: ", err)
					return fmt.Errorf("mapUpdateOutgoing failed: %v", err)
				}
				filerData.TopExpThreshold = th
			}
		} else {
			// update by comparison
			comp, err := updateTopRecipients(filerData, receiver)
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
		if len(filerData.TransferRecsAmt) == 0 {
			filerData.TransferRecsAmt = make(map[string]float32)
			filerData.TransferRecsTxs = make(map[string]float32)
		}
		if len(receiver.(*donations.CmteTxData).TopCmteOrgContributorsAmt) == 0 {
			receiver.(*donations.CmteTxData).TopCmteOrgContributorsAmt = make(map[string]float32)
			receiver.(*donations.CmteTxData).TopCmteOrgContributorsTxs = make(map[string]float32)
		}

		// update filer's TransfersRecs maps only -- receiving committees maps updated in corresponding tx
		// all outgoing transactions between committees are considered transfers
		filerData.TransferRecsAmt[receiver.(*donations.CmteTxData).CmteID] += cont.TxAmt
		filerData.TransferRecsTxs[receiver.(*donations.CmteTxData).CmteID]++
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
				if len(filerData.TopExpRecipientsAmt) >= 100 {
					th, err := checkThreshold(receiver.(*donations.Candidate).ID, filerData.TopExpRecipientsAmt, filerData.TopExpThreshold)
					if err != nil {
						fmt.Println("mapUpdateOutgoing failed: ", err)
						return fmt.Errorf("mapUpdateOutgoing failed: %v", err)
					}
					filerData.TopExpThreshold = th
				}
			} else {
				// update by comparison
				comp, err := updateTopRecipients(filerData, receiver)
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
func mapUpdateOpExp(disb *donations.Disbursement, filer *donations.CmteTxData, receiver *donations.Individual) error {
	// re-initialize maps if nil
	if len(filer.TopExpRecipientsAmt) == 0 {
		filer.TopExpRecipientsAmt = make(map[string]float32)
		filer.TopExpRecipientsTxs = make(map[string]float32)
	}
	if len(receiver.SendersAmt) == 0 {
		receiver.SendersAmt = make(map[string]float32)
		receiver.SendersTxs = make(map[string]float32)
	}

	// update receiver's sender's maps
	receiver.SendersAmt[filer.CmteID] += disb.TxAmt
	receiver.SendersTxs[filer.CmteID]++

	// update filer's top expenditure recipient maps
	if len(filer.TopExpRecipientsAmt) < 100 || filer.TopExpRecipientsAmt[receiver.ID] != 0 {
		// update by direct add
		filer.TopExpRecipientsAmt[receiver.ID] += disb.TxAmt
		filer.TopExpRecipientsTxs[receiver.ID]++

		if len(filer.TopExpRecipientsAmt) >= 100 {
			th, err := checkThreshold(receiver.ID, filer.TopExpRecipientsAmt, filer.TopExpThreshold)
			if err != nil {
				fmt.Println("mapUpdateOpExp failed: ", err)
				return fmt.Errorf("mapUpdateOpExp failed: %v", err)
			}
			filer.TopExpThreshold = th
		}
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

	return nil
}
