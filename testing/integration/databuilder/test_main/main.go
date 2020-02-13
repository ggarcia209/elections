package main

import (
	"fmt"
	"os"

	"github.com/elections/donations"
	"github.com/elections/testing/integration/databuilder/testDB"
)

func main() {
	err := testDB.TransactionUpdate(testDB.ExpensesConts)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	err = testDB.TransactionUpdate(testDB.ExpensesDisbs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	fmt.Println("FILER")
	printCmte(&testDB.Filer)
	/* fmt.Println("indv10")
	printIndv(&testDB.Indv10)
	fmt.Println("org01")
	printOrg(&testDB.Org01)
	fmt.Println("org02")
	printOrg(&testDB.Org02)
	fmt.Println("org03")
	printOrg(&testDB.Org03)
	fmt.Println("indv11")
	printIndv(&testDB.Indv11)
	fmt.Println("indv8")
	printIndv(&testDB.Indv8)
	fmt.Println("org01")
	printOrg(&testDB.Org01)
	fmt.Println("org03")
	printOrg(&testDB.Org03) */
	/* 	fmt.Println("cmte01")
	   	printCmte(&testDB.Cmte01)
	   	fmt.Println("cmte03")
	   	printCmte(&testDB.Cmte03)
	   	fmt.Println("cand01")
	   	printCand(&testDB.Cand01)
	   	fmt.Println("cand03")
	   	printCand(&testDB.Cand03)*/

}

func printCmte(c *donations.CmteTxData) {
	fmt.Println("*** Committee Data ***")
	fmt.Printf("ID: %s\n", c.CmteID)
	fmt.Println("* Transactions *")
	fmt.Printf("ContributionsInAmt: %v\tContributionsInTxs: %v\tAvgContributionIn: %v\n", c.ContributionsInAmt, c.ContributionsInTxs, c.AvgContributionIn)
	fmt.Printf("OtherInAmt: %v\tOtherInTxs: %v\tAvgOtherIn: %v\n", c.OtherReceiptsInAmt, c.OtherReceiptsInTxs, c.AvgOtherIn)
	fmt.Printf("TotalIncomingAmt: %v\tTotalIncomingTxs: %v\tAvgIncoming: %v\n", c.TotalIncomingAmt, c.TotalIncomingTxs, c.AvgIncoming)
	fmt.Printf("TransfersAmt: %v\tTransfersTxs: %v\tAvgTransfer: %v\n", c.TransfersAmt, c.TransfersTxs, c.AvgTransfer)
	fmt.Printf("ExpendituresAmt: %v\tExpendituresTxs: %v\tAvgExpenditure: %v\n", c.ExpendituresAmt, c.ExpendituresTxs, c.AvgExpenditure)
	fmt.Printf("TotalOutgoingAmt: %v\tTotalOutgoingTxs: %v\tAvgOutgoing: %v\n", c.TotalOutgoingAmt, c.TotalOutgoingTxs, c.AvgOutgoing)
	fmt.Println("NetBalance: ", c.NetBalance)

	fmt.Println("* Maps *")
	fmt.Printf("TopIndividualsAmts/Txs: \n\t%v\n\t%v\n", c.TopIndvContributorsAmt, c.TopIndvContributorsTxs)
	fmt.Println("TopIndvTreshold:")
	printThreshold(c.TopIndvContributorThreshold)
	fmt.Printf("TopCmteOrgsAmts/Txs: \n\t%v\n\t%v\n", c.TopCmteOrgContributorsAmt, c.TopCmteOrgContributorsTxs)
	fmt.Println("TopCmteOrgThreshold:")
	printThreshold(c.TopCmteOrgContributorThreshold)
	fmt.Printf("TransferRecsAmts/Txs: \n\t%v\n\t%v\n", c.TransferRecsAmt, c.TransferRecsTxs)
	fmt.Printf("TopExpRecipientsAmts/Txs: \n\t%v\n\t%v\n", c.TopExpRecipientsAmt, c.TopExpRecipientsTxs)
	fmt.Println("TopExpThreshold:")
	printThreshold(c.TopExpThreshold)
}

func printIndv(i *donations.Individual) {
	fmt.Println("*** Individual Data ***")
	fmt.Println("ID: ", i.ID)
	fmt.Printf("TotalInAmt/Txs: %v\t%v\tAvgTxIn: %v\n", i.TotalInAmt, i.TotalInTxs, i.AvgTxIn)
	fmt.Printf("TotalOutAm/Txs: %v\t%v\tAvgtxOut: %v\n", i.TotalOutAmt, i.TotalOutTxs, i.AvgTxOut)
	fmt.Printf("NetBalance: %v\n", i.NetBalance)
	fmt.Printf("RecipientsAmt/Txs: \n\t%v\n\t%v\n", i.RecipientsAmt, i.RecipientsTxs)
	fmt.Printf("SendersAmt/Txs: \n\t%v\n\t%v\n", i.SendersAmt, i.SendersTxs)
	fmt.Println()
}

func printOrg(i *donations.Organization) {
	fmt.Println("*** Organization Data ***")
	fmt.Println("ID: ", i.ID)
	fmt.Printf("TotalInAmt/Txs: %v\t%v\tAvgTxIn: %v\n", i.TotalInAmt, i.TotalInTxs, i.AvgTxIn)
	fmt.Printf("TotalOutAm/Txs: %v\t%v\tAvgtxOut: %v\n", i.TotalOutAmt, i.TotalOutTxs, i.AvgTxOut)
	fmt.Printf("NetBalance: %v\n", i.NetBalance)
	fmt.Printf("RecipientsAmt/Txs: \n\t%v\n\t%v\n", i.RecipientsAmt, i.RecipientsTxs)
	fmt.Printf("SendersAmt/Txs: \n\t%v\n\t%v\n", i.SendersAmt, i.SendersTxs)
	fmt.Println()
}

func printCand(i *donations.Candidate) {
	fmt.Println("*** Candidate Data ***")
	fmt.Println("ID: ", i.ID)
	fmt.Printf("TotalInAmt/Txs: %v\t%v\tAvgTxIn: %v\n", i.TotalDirectInAmt, i.TotalDirectInTxs, i.AvgDirectIn)
	fmt.Printf("TotalOutAm/Txs: %v\t%v\tAvgtxOut: %v\n", i.TotalDirectOutAmt, i.TotalDirectOutTxs, i.AvgDirectOut)
	fmt.Printf("NetBalance: %v\n", i.NetBalanceDirectTx)
	fmt.Printf("RecipientsAmt/Txs: \n\t%v\n\t%v\n", i.DirectRecipientsAmts, i.DirectRecipientsTxs)
	fmt.Printf("SendersAmt/Txs: \n\t%v\n\t%v\n", i.DirectSendersAmts, i.DirectSendersTxs)
	fmt.Println()
}

func printThreshold(es []interface{}) {
	for _, e := range es {
		fmt.Printf("\tID: %v\tTotal: %v\n", e.(*donations.Entry).ID, e.(*donations.Entry).Total)
	}
}
