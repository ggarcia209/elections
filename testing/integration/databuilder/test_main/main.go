package main

import (
	"fmt"
	"os"

	"github.com/elections/donations"
	"github.com/elections/testing/integration/databuilder/testDB"
)

func main() {
	/* err := testDB.TransactionUpdate(testDB.Conts)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	err = testDB.TransactionUpdate(testDB.Disbs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	} */
	err := testDB.TransactionUpdate(testDB.Transfers)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	printCmte(&testDB.Filer)

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
	fmt.Printf("TopIndividualsAmts/Txs: %v\t%v\n", c.TopIndvContributorsAmt, c.TopIndvContributorsTxs)
	fmt.Println("TopIndvTreshold:")
	printThreshold(c.TopIndvContributorThreshold)
	fmt.Printf("TopCmteOrgsAmts/Txs: %v\t%v\n", c.TopCmteOrgContributorsAmt, c.TopCmteOrgContributorsTxs)
	fmt.Println("TopCmteOrgThreshold:")
	printThreshold(c.TopCmteOrgContributorThreshold)
	fmt.Printf("TransferRecsAmts/Txs: %v\t%v\n", c.TransferRecsAmt, c.TransferRecsTxs)
	fmt.Printf("TopExpRecipientsAmts/Txs: %v\t%v\n", c.TopExpRecipientsAmt, c.TopExpRecipientsTxs)
	fmt.Println("TopExpThreshold:")
	printThreshold(c.TopExpThreshold)
}

func printThreshold(es []interface{}) {
	for _, e := range es {
		fmt.Printf("\tID: %v\tTotal: %v\n", e.(*donations.Entry).ID, e.(*donations.Entry).Total)
	}
}
