// Package donations contains the base objects that are used throughout the application.
// Objects within this package are primarily used for creating, updating, and persisting
// the datasets derived from the input data.
// Individual object data is derived from the Contribution and Disbursement objects.
// Struct fields for other objects in this file are populated from the data parsed
// from the input .txt files. One object is created for each record.
package donations

import "time"

// Individual donor represents an individual donor,
// business, or other private or non-federal govt. entity.
type Individual struct {
	ID            string
	Name          string
	City          string
	State         string
	Zip           string
	Occupation    string
	Employer      string
	Transactions  []string           // List of all incoming/outgoing transactions
	TotalOutAmt   float32            // Total $ Vale of Outgoing Transactions
	TotalOutTxs   float32            // Total # of Contributions/Loans To/etc
	AvgTxOut      float32            // Average value of outgoing transactions
	TotalInAmt    float32            // Total Amount of Incoming Transactions
	TotalInTxs    float32            // Total # of Refunds/Repayments/etc
	AvgTxIn       float32            // Average value of incoming transactions
	NetBalance    float32            // TotalInAmt - TotalOutAmt (negative balance indicates funds out > funds in)
	RecipientsAmt map[string]float32 // # of Txs to each committee
	RecipientsTxs map[string]float32 // $ Value contributed to each committee
	SendersTxs    map[string]float32 // # of Txs from each committee
	SendersAmt    map[string]float32 // $ Value returned from each committee
}

// Committee represents a federal politcal committee
// registered with the FEC.
// Commitee objects both receive and send transactions.
type Committee struct {
	ID           string
	Name         string
	TresName     string // treasurer name
	st1          string // address
	st2          string
	City         string
	State        string
	Zip          string
	Designation  string
	Type         string
	Party        string
	FilingFreq   string
	OrgType      string // interest group category
	ConnectedOrg string
	CandID       string // null if Type != "H", "S", or "P"
}

// CmteTxData contains incoming/outgoing cashflow data, top contributors/recipiens of cashflows,
// and the corresponding total $ values/# of transactions for each contributor/recipient.
// Candidate data is derived by aggregating all affiliated committees into one CmteTxData object.
type CmteTxData struct {
	CmteID                         string             // ID of committee directly linked to data
	CandID                         string             // ID of candidate indirectly linked through Candidate PCC ID (nil if non-affiliated committee)
	Party                          string             // Committee's political party
	ContributionsInAmt             float32            // $ value of incoming contributions
	ContributionsInTxs             float32            // # contributions from individuals, organizations, committees, and candidates
	AvgContributionIn              float32            // Average $ value of incoming contributions
	OtherReceiptsInAmt             float32            // $ value of loans from/refunds from/other incoming transactions
	OtherReceiptsInTxs             float32            // # of loans from/refunds from/other incoming transactions
	AvgOtherIn                     float32            // Average $ value of other incoming receipts
	TotalIncomingAmt               float32            // Total $ value of incoming transactions
	TotalIncomingTxs               float32            // Total # of incoming transactions
	AvgIncoming                    float32            // Average $ value of incoming transactions
	TransfersAmt                   float32            // $ value of contributions/transfers/loans to other committees
	TransfersTxs                   float32            // # of contributions/transfers/loans to other committees
	AvgTransfer                    float32            // Average value of transfers to other committees
	TransfersList                  []string           // list of transfer tx ID's
	ExpendituresAmt                float32            // $ value of expenditure transactions (operating expenses/loan repayments/refunds/etc)
	ExpendituresTxs                float32            // # of expenditure transactions (operating expenses/loan repayments/refunds/etc)
	AvgExpenditure                 float32            // Average value of expenditures
	TotalOutgoingAmt               float32            // Total outgoing $ Value (TransfersAmt + ExpendituresAmt)
	TotalOutgoingTxs               float32            // Total # of outgoing transactions (TransfersTxs + ExpendituresTxs)
	AvgOutgoing                    float32            // Average outgoing transaction
	NetBalance                     float32            // NetBalance = TotalIncomingAmt - TotalOutgoingAmt
	TopIndvContributorsAmt         map[string]float32 // Top Individuals by $ value contributed
	TopIndvContributorsTxs         map[string]float32 // # of transactions for each top contributor by $ value
	TopIndvContributorThreshold    []interface{}      // Minimum values to be in Top x Contributors
	TopCmteOrgContributorsAmt      map[string]float32 // Top Committee and Organization contributors by $ value contributed
	TopCmteOrgContributorsTxs      map[string]float32 // Number of transactions for each top contributor by $ value
	TopCmteOrgContributorThreshold []interface{}      // Minimum values to be in Top x Contributors
	TransferRecsAmt                map[string]float32 // total $ value of transactions to each recipient committee
	TransferRecsTxs                map[string]float32 // # of transactions for each recipient committee
	TopExpRecipientsAmt            map[string]float32 // Top expenditure recipients by $ value
	TopExpRecipientsTxs            map[string]float32 // # of transactions for each top recipient by $ value
	TopExpThreshold                []interface{}      // Minimum values to be in Top x Recipients
}

// CmteFinancials represents the financial data of a political action committee.
type CmteFinancials struct {
	CmteID      string
	name        string
	Type        string
	designation string
	filingFreq  string
	// all following int values represent dollar values
	TotalReceipts   float32   // total receipts
	TxsFromAff      float32   //  transfers from affilliates ($)
	IndvConts       float32   // individual contributions ($)
	OtherConts      float32   // Other political committee contributions ($)
	CandCont        float32   // contributions from candidate
	CandLoans       float32   // candidate loans
	TotalLoans      float32   // total loans received
	TotalDisb       float32   // total disbursements
	TxToAff         float32   // transfers to affiliates
	IndvRefunds     float32   // Refunds to individuals
	OtherRefunds    float32   // other political committee refunds
	LoanRepay       float32   // candidate loan repayments
	CashBOP         float32   // cash at beginning of period
	CashCOP         float32   // cash at end of period
	DebtsOwed       float32   // debts owed by
	NonFedTxsRecvd  float32   // non federal transfers received
	ContToOtherCmte float32   // contributions to other committess
	IndExp          float32   // independent expenditures
	PartyExp        float32   // party coordinated expenditures
	NonFedSharedExp float32   // non-federal shared expenditures
	CovgEndDate     time.Time // coverage end date
}
