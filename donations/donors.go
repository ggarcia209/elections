package donations

import "time"

// Individual donor represents an individual donor.
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

// Committee represents a committee
// Commitee objects both receive and send donations
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

// CmteFinancials represents the financial data of a political action committee
type CmteFinancials struct {
	CmteID      string
	name        string
	Type        string
	designation string
	filingFreq  string
	// all following int values represent dollar values
	TotalReceipts   int       // total receipts
	TxsFromAff      int       //  transfers from affilliates ($)
	IndvConts       int       // individual contributions ($)
	OtherConts      int       // Other political committee contributions ($)
	CandCont        int       // contributions from candidate
	CandLoans       int       // candidate loans
	TotalLoans      int       // total loans received
	TotalDisb       int       // total disbursements
	TxToAff         int       // transfers to affiliates
	IndvRefunds     int       // Refunds to individuals
	OtherRefunds    int       // other political committee refunds
	LoanRepay       int       // candidate loan repayments
	CashBOP         int       // cash at beginning of period
	CashCOP         int       // cash at end of period
	DebtsOwed       int       // debts owed by
	NonFedTxsRecvd  int       // non federal transfers received
	ContToOtherCmte int       // contributions to other committess
	IndExp          int       // independent expenditures
	PartyExp        int       // party coordinated expenditures
	NonFedSharedExp int       // non-federal shared expenditures
	CovgEndDate     time.Time // coverage end date
}

// DEPRECATED
/*

// DisbRecipient representst a recipient of a committee's disbursement
type DisbRecipient struct {
	ID                 string
	Name               string
	City               string
	State              string
	Zip                string
	Disbursements      []string // IN <- Disb. TxID's
	TotalDisbursements float32
	TotalReceived      float32
	AvgReceived        float32
	SendersAmt         map[string]float32 // IN <- disbursements from committees
	SendersTxs         map[string]float32 // IN <- disbursements from committees
}

*/
