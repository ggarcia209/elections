package testDB

import (
	"fmt"
	"strings"

	"github.com/elections/donations"
)

// SUCCESS
func TestTxUpdateInternalLogic() {
	txType := "17"
	memo := true
	if memo {
		fmt.Println("memo == true")
	}
	if txType < "16" || txType > "18" {
		fmt.Println("account == contributionsIn")
	} else {
		fmt.Println("account == otherReceipts")
	}
	// switch cases verifed in test_findTopX.go
}

// SUCCESS
func TestDeriveTxTypes() {
	cont1 := &donations.Contribution{
		TxType:   "15",
		MemoCode: "",
		OtherID:  "C0001",
	}
	cont2 := &donations.Contribution{
		TxType:   "32K",
		MemoCode: "X",
		OtherID:  "H0001",
	}
	cont3 := &donations.Contribution{
		TxType:     "20",
		MemoCode:   "",
		OtherID:    "I0001",
		Occupation: "worker",
	}
	cont4 := &donations.Contribution{
		TxType:   "24G",
		MemoCode: "",
		OtherID:  "O0001",
	}
	b, i, t, m := deriveTxTypes(cont1)
	fmt.Println(b, i, t, m)
	b, i, t, m = deriveTxTypes(cont2)
	fmt.Println(b, i, t, m)
	b, i, t, m = deriveTxTypes(cont3)
	fmt.Println(b, i, t, m)
	b, i, t, m = deriveTxTypes(cont4)
	fmt.Println(b, i, t, m)
}

// SUCCESS
func TestTransactionUpdate() error {
	year := "2018"

	/* Top Overall Objects */
	// Create test record in DbSim
	category := "cmte_recs_all"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_recs_na"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_recs_r"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_recs_d"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	// Create test record in DbSim
	category = "cmte_donors_all"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_donors_na"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_donors_d"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	// Create test record in DbSim
	category = "cmte_exp_all"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_exp_na"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_exp_d"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	// Create test record in DbSim
	category = "cand_all"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cand_na"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	// Create test record in DbSim
	category = "cand_exp_all"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cand_exp_na"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	// Create test record in DbSim
	category = "indv"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	// Create test record in DbSim
	category = "org_conts"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "org_recs"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	/* END TOP OVERALL OBJECTS */

	// test 1 - contributions
	fmt.Println("test 1: contributions")
	err := TransactionUpdate(year, Conts)
	if err != nil {
		fmt.Println("TestTransactionUpdate failed: ", err)
		return fmt.Errorf("TestTransactionUpdate failed: %v", err)
	}
	// printCmteTxData(Filer)

	// test 2 - disbursements
	fmt.Println("test 2 - disbursements")
	err = TransactionUpdate(year, Disbs)
	if err != nil {
		fmt.Println("TestTransactionUpdate failed: ", err)
		return fmt.Errorf("TestTransactionUpdate failed: %v", err)
	}

	// check filing and other committee objects - SUCCESS
	printCmteTxData(Filer)
	/* printCmteTxData(Cmte01)

	// check Top Overall Objects - SUCCESS
	for _, v := range DbSim[year]["top_overall"] {
		data := v.(*donations.TopOverallData)
		if data.Category == "" {
			continue
		}
		printOverallData(*data)
	}

	// check Individual, Organization, & Candidate Objects - SUCCESS
	printIndvData(Indv10)
	printIndvData(Indv8)

	printOrgData(Org01)
	printOrgData(Org02)
	printOrgData(Org03)

	printCandData(Cand00)
	printCandData(Cand01)
	printCandData(Cand02)
	printCandData(Cand03) */

	// verify data accuracy - SUCCESS

	return nil
}

func printCmteTxData(c donations.CmteTxData) {
	fmt.Println("***** COMMITTEE TRANSACTION DATA *****")
	fmt.Println("CmteID: ", c.CmteID)
	fmt.Println("CandID: ", c.CandID)
	fmt.Println("Party: ", c.Party)
	fmt.Println()
	fmt.Printf("Contributions: %v / %v - Avg: %v\n", c.ContributionsInAmt, c.ContributionsInTxs, c.AvgContributionIn)
	fmt.Printf("Other Receipts: %v / %v - Avg: %v\n", c.OtherReceiptsInAmt, c.OtherReceiptsInTxs, c.AvgOtherIn)
	fmt.Printf("Total Incoming: %v / %v - Avg: %v\n", c.TotalIncomingAmt, c.TotalIncomingTxs, c.AvgIncoming)
	fmt.Println()
	fmt.Printf("Transfers: %v / %v - Avg: %v\n", c.TransfersAmt, c.TransfersTxs, c.AvgTransfer)
	fmt.Printf("Expenditures: %v / %v - Avg: %v\n", c.ExpendituresAmt, c.ExpendituresTxs, c.AvgExpenditure)
	fmt.Printf("Total Outgoing: %v / %v - Avg: %v\n", c.TotalOutgoingAmt, c.TotalOutgoingTxs, c.AvgOutgoing)
	fmt.Println("Net Balance: ", c.NetBalance)
	fmt.Println()
	fmt.Println("Top Individual Contributors Amounts: ", c.TopIndvContributorsAmt)
	fmt.Println("Top Individual Contributors Transactions: ", c.TopIndvContributorsTxs)
	fmt.Println("Top Individuals Threshold: ")
	printThreshold(c.TopIndvContributorThreshold)
	fmt.Println()
	fmt.Println("Top Committee/Organization Contributors Amounts: ", c.TopCmteOrgContributorsAmt)
	fmt.Println("Top Committee/Organization Contributors Tranactions: ", c.TopCmteOrgContributorsTxs)
	fmt.Println("Top Committee/Organization Contributors Threshold: ")
	printThreshold(c.TopCmteOrgContributorThreshold)
	fmt.Println()
	fmt.Println("Transfers Recipients Amounts: ", c.TransferRecsAmt)
	fmt.Println("Transfers Recipients Transactions: ", c.TransferRecsTxs)
	fmt.Println()
	fmt.Println("Top Expense Recipients Amounts: ", c.TopExpRecipientsAmt)
	fmt.Println("Top Expense Recipients Transactions: ", c.TopExpRecipientsTxs)
	fmt.Println("Top Expense Recipients Threshold: ", c.TopExpThreshold)
	fmt.Println()
	fmt.Println("***** END COMMITTEE TRANSACTION DATA *****")
	fmt.Println()
	fmt.Println()

}

func printIndvData(i donations.Individual) {
	fmt.Println("***** INDIVIDUAL DATA *****")
	fmt.Println("ID: ", i.ID)
	fmt.Println("Transactions: ", i.Transactions)
	fmt.Printf("TotalOut: %v / %v - Avg: %v\n", i.TotalOutAmt, i.TotalOutTxs, i.AvgTxOut)
	fmt.Printf("TotalIn: %v / %v - Avg: %v\n", i.TotalInAmt, i.TotalInTxs, i.AvgTxIn)
	fmt.Println("Net Balance: ", i.NetBalance)
	fmt.Println()
	fmt.Println("Recipients Amounts: ", i.RecipientsAmt)
	fmt.Println("Recipients Transactions: ", i.RecipientsTxs)
	fmt.Println("Senders Amounts: ", i.SendersAmt)
	fmt.Println("Senders Transactions: ", i.SendersTxs)
	fmt.Println("***** END INDIVIDUAL DATA *****")
	fmt.Println()
	fmt.Println()
}

func printOrgData(o donations.Organization) {
	fmt.Println("***** ORGANIZATION DATA *****")
	fmt.Println("ID: ", o.ID)
	fmt.Println("Transactions: ", o.Transactions)
	fmt.Printf("TotalOut: %v / %v - Avg: %v\n", o.TotalOutAmt, o.TotalOutTxs, o.AvgTxOut)
	fmt.Printf("TotalIn: %v / %v - Avg: %v\n", o.TotalInAmt, o.TotalInTxs, o.AvgTxIn)
	fmt.Println("Net Balance: ", o.NetBalance)
	fmt.Println()
	fmt.Println("Recipients Amounts: ", o.RecipientsAmt)
	fmt.Println("Recipients Transactions: ", o.RecipientsTxs)
	fmt.Println("Senders Amounts: ", o.SendersAmt)
	fmt.Println("Senders Transactions: ", o.SendersTxs)
	fmt.Println("***** END ORGANIZATION DATA *****")
	fmt.Println()
	fmt.Println()
}

func printCandData(c donations.Candidate) {
	fmt.Println("***** CANDIDATE DATA *****")
	fmt.Println("ID: ", c.ID)
	fmt.Println("Other affiliates: ", c.OtherAffiliates)
	fmt.Println()
	fmt.Printf("Total Direct In: %v / %v - Avg: %v\n", c.TotalDirectInAmt, c.TotalDirectInTxs, c.AvgDirectIn)
	fmt.Printf("Total Direct Out: %v / %v - Avg: %v\n", c.TotalDirectOutAmt, c.TotalDirectOutTxs, c.AvgDirectOut)
	fmt.Println("Net Balance - Direct Transactions: ", c.NetBalanceDirectTx)
	fmt.Println()
	fmt.Println("Direct Recipients Amounts: ", c.DirectRecipientsAmts)
	fmt.Println("Direct Recipients Transactions: ", c.DirectRecipientsTxs)
	fmt.Println("Direct Senders Amounts: ", c.DirectSendersAmts)
	fmt.Println("Direct Senders Transactions: ", c.DirectSendersTxs)
	fmt.Println("***** END CANDIDATE DATA *****")
	fmt.Println()
	fmt.Println()
}

// TransactionUpdate updates each sender/receiver data for each transaction in a list of transactions.
// 7/22/20 - Correct logic to incorporate "year" database key - DONE
func TransactionUpdate(year string, txs interface{}) error {
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
	for i, cont := range conts {
		fmt.Printf("%d: filing cmte: %s\n", i, cont.CmteID)
		fmt.Printf("%d: otherID: %s\n", i, cont.OtherID)
		fmt.Println()
		// get tx type info
		bucket, incoming, transfer, memo := deriveTxTypes(cont)

		// get filer object
		filer := DbSim[year]["cmte_tx_data"][cont.CmteID]

		// get sender/receiver objects
		other := DbSim[year][bucket][cont.OtherID]

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
		// 7/22/20 - refactor & test
		if !memo {
			// update top individuals, organizations and committees
			err := updateTopOverall(year, filer.(*donations.CmteTxData), other, incoming, transfer)
			if err != nil {
				fmt.Println("ContributionUpdate failed: ", err)
				return fmt.Errorf("ContributionUpdate failed: %v", err)
			}
			/* Moving this block to within func updateTopOverall()
			// update top candidates by funds received if candidate linked to filing committee
			if filer.(*donations.CmteTxData).CandID != "" {
				// get linked candidate
				cand := DbSim[year]["candidates"][filer.(*donations.CmteTxData).CandID]

				// update top candidates by total funds incoming/outgoing
				err = updateTopCandidates(year, cand.(*donations.Candidate), filer.(*donations.CmteTxData), incoming)
				if err != nil {
					fmt.Println("ContributionUpdate failed: ", err)
					return fmt.Errorf("ContributionUpdate failed: %v", err)
				}
			} */
		}

		// persist objects
		DbSim[year]["cmte_tx_data"][cont.CmteID] = filer
		DbSim[year][bucket][cont.OtherID] = other
	}

	return nil
}

// update data from Disbursement transactions derived from operating expenses files
func opExpensesUpdate(year string, disbs []*donations.Disbursement) error {
	for _, disb := range disbs {
		// get  filing committee
		filer := DbSim[year]["cmte_tx_data"][disb.CmteID]
		// get receiving organization
		receiver := DbSim[year]["organizations"][disb.RecID]

		// update object account totals
		err := disbursementTxUpdate(disb, filer.(*donations.CmteTxData), receiver.(*donations.Organization))
		if err != nil {
			fmt.Println("OpExpensesUpdate failed: ", err)
			return fmt.Errorf("OpExpensesUpdate failed: %v", err)
		}

		// update TopOverall rankings
		// update top individuals, organizations and committees
		err = updateTopOverall(year, filer.(*donations.CmteTxData), receiver, false, false)
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}
		/* Moving this block to within func updateTopOverall()
		// update top candidates by funds received if candidate linked to filing committee
		if filer.(*donations.CmteTxData).CandID != "" {
			// get linked candidate
			cand := DbSim[year]["candidates"][filer.(*donations.CmteTxData).CandID]

			// update top candidates by total funds incoming/outgoing
			err = updateTopCandidates(year, cand.(*donations.Candidate), filer.(*donations.CmteTxData), false)
			if err != nil {
				fmt.Println("ContributionUpdate failed: ", err)
				return fmt.Errorf("ContributionUpdate failed: %v", err)
			}
		} */

		// persist objects
		DbSim[year]["cmte_tx_data"][disb.CmteID] = filer
		DbSim[year]["organizations"][disb.RecID] = receiver
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
		if len(filerData.TopIndvContributorsAmt) < 5 || filerData.TopIndvContributorsAmt[sender.(*donations.Individual).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopIndvContributorsAmt[sender.(*donations.Individual).ID] += cont.TxAmt
			filerData.TopIndvContributorsTxs[sender.(*donations.Individual).ID]++
			th, err := checkThreshold(sender.(*donations.Individual).ID, filerData.TopIndvContributorsAmt, filerData.TopIndvContributorThreshold)
			if err != nil {
				fmt.Println("mapUpdateIncoming failed: ", err)
				return fmt.Errorf("mapUpdateIncoming failed: %v", err)
			}
			filerData.TopIndvContributorThreshold = th
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
		if len(filerData.TopCmteOrgContributorsAmt) < 5 || filerData.TopCmteOrgContributorsAmt[sender.(*donations.Organization).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopCmteOrgContributorsAmt[sender.(*donations.Organization).ID] += cont.TxAmt
			filerData.TopCmteOrgContributorsTxs[sender.(*donations.Organization).ID]++
			th, err := checkThreshold(sender.(*donations.Organization).ID, filerData.TopCmteOrgContributorsAmt, filerData.TopCmteOrgContributorThreshold)
			if err != nil {
				fmt.Println("mapUpdateIncoming failed: ", err)
				return fmt.Errorf("mapUpdateIncoming failed: %v", err)
			}
			filerData.TopCmteOrgContributorThreshold = th
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
		if len(filerData.TopCmteOrgContributorsAmt) < 5 || filerData.TopCmteOrgContributorsAmt[sender.(*donations.CmteTxData).CmteID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopCmteOrgContributorsAmt[sender.(*donations.CmteTxData).CmteID] += cont.TxAmt
			filerData.TopCmteOrgContributorsTxs[sender.(*donations.CmteTxData).CmteID]++
			th, err := checkThreshold(sender.(*donations.CmteTxData).CmteID, filerData.TopCmteOrgContributorsAmt, filerData.TopCmteOrgContributorThreshold)
			if err != nil {
				fmt.Println("mapUpdateIncoming failed: ", err)
				return fmt.Errorf("mapUpdateIncoming failed: %v", err)
			}
			filerData.TopCmteOrgContributorThreshold = th
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
		if len(filerData.TopIndvContributorsAmt) < 5 || filerData.TopIndvContributorsAmt[sender.(*donations.Candidate).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopIndvContributorsAmt[sender.(*donations.Candidate).ID] += cont.TxAmt
			filerData.TopIndvContributorsTxs[sender.(*donations.Candidate).ID]++
			th, err := checkThreshold(sender.(*donations.Candidate).ID, filerData.TopIndvContributorsAmt, filerData.TopIndvContributorThreshold)
			if err != nil {
				fmt.Println("mapUpdateIncoming failed: ", err)
				return fmt.Errorf("mapUpdateIncoming failed: %v", err)
			}
			filerData.TopIndvContributorThreshold = th
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
		if len(filerData.TopExpRecipientsAmt) < 5 || filerData.TopExpRecipientsAmt[receiver.(*donations.Individual).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopExpRecipientsAmt[receiver.(*donations.Individual).ID] += cont.TxAmt
			filerData.TopExpRecipientsTxs[receiver.(*donations.Individual).ID]++
			th, err := checkThreshold(receiver.(*donations.Individual).ID, filerData.TopExpRecipientsAmt, filerData.TopExpThreshold)
			if err != nil {
				fmt.Println("mapUpdateOutgoing failed: ", err)
				return fmt.Errorf("mapUpdateOutgoing failed: %v", err)
			}
			filerData.TopExpThreshold = th
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
		if len(filerData.TopExpRecipientsAmt) < 5 || filerData.TopExpRecipientsAmt[receiver.(*donations.Organization).ID] != 0 {
			// add new entry directly or update existing entry
			filerData.TopExpRecipientsAmt[receiver.(*donations.Organization).ID] += cont.TxAmt
			filerData.TopExpRecipientsTxs[receiver.(*donations.Organization).ID]++
			th, err := checkThreshold(receiver.(*donations.Organization).ID, filerData.TopExpRecipientsAmt, filerData.TopExpThreshold)
			if err != nil {
				fmt.Println("mapUpdateOutgoing failed: ", err)
				return fmt.Errorf("mapUpdateOutgoing failed: %v", err)
			}
			filerData.TopExpThreshold = th
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
			if len(filerData.TopExpRecipientsAmt) < 5 || filerData.TopExpRecipientsAmt[receiver.(*donations.Candidate).ID] != 0 {
				// add new entry directly or update existing entry
				filerData.TopExpRecipientsAmt[receiver.(*donations.Candidate).ID] += cont.TxAmt
				filerData.TopExpRecipientsTxs[receiver.(*donations.Candidate).ID]++
				th, err := checkThreshold(receiver.(*donations.Candidate).ID, filerData.TopExpRecipientsAmt, filerData.TopExpThreshold)
				if err != nil {
					fmt.Println("mapUpdateOutgoing failed: ", err)
					return fmt.Errorf("mapUpdateOutgoing failed: %v", err)
				}
				filerData.TopExpThreshold = th
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

	// update receiver's sender's maps
	receiver.SendersAmt[filer.CmteID] += disb.TxAmt
	receiver.SendersTxs[filer.CmteID]++

	// update filer's top expenditure recipient maps
	if len(filer.TopExpRecipientsAmt) < 5 || filer.TopExpRecipientsAmt[receiver.ID] != 0 {
		// update by direct add
		filer.TopExpRecipientsAmt[receiver.ID] += disb.TxAmt
		filer.TopExpRecipientsTxs[receiver.ID]++
		th, err := checkThreshold(receiver.ID, filer.TopExpRecipientsAmt, filer.TopExpThreshold)
		if err != nil {
			fmt.Println("mapUpdateOpExp failed: ", err)
			return fmt.Errorf("mapUpdateOpExp failed: %v", err)
		}
		filer.TopExpThreshold = th
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
