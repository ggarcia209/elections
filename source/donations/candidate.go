// Package donations contains the base objects that are used throughout the application.
// Objects within this package are primarily used for creating, updating, and persisting
// the datasets derived from the input data.
// Struct fields for objects in this file are populated from the data parsed
// from the input .txt files. One object is created for each record.
package donations

import "time"

// Candidate represents a candidate for federal office (House, Senate, President).
type Candidate struct {
	ID                   string
	Name                 string
	Party                string
	ElectnYr             string
	OfficeState          string
	Office               string
	officeDist           string
	ici                  string
	candStatus           string
	PCC                  string // principal campaign committee
	st1                  string // mailing address
	st2                  string // mailing address 2
	City                 string
	State                string
	Zip                  string
	OtherAffiliates      []string // ID's of other affiliated committees
	TransactionsList     []string // all direct incoming/outgoing transactions
	TotalDirectInAmt     float32
	TotalDirectInTxs     float32
	AvgDirectIn          float32
	TotalDirectOutAmt    float32
	TotalDirectOutTxs    float32
	AvgDirectOut         float32
	NetBalanceDirectTx   float32
	DirectRecipientsAmts map[string]float32 // Direct recipients receive funds directly from the candidate
	DirectRecipientsTxs  map[string]float32
	DirectSendersAmts    map[string]float32 // DirectSenders send funds directly to the candidate
	DirectSendersTxs     map[string]float32
}

// CmpnFinancials contains financial data reported by a candidate's campaign.
type CmpnFinancials struct {
	CandID         string
	Name           string
	ici            string
	PartyCd        string
	Party          string
	TotalReceipts  float32
	TransFrAuth    float32
	TotalDisbsmts  float32
	TransToAuth    float32
	COHBOP         float32
	COHCOP         float32
	CandConts      float32
	CandLoans      float32
	OtherLoans     float32
	CandLoanRepay  float32
	OtherLoanRepay float32
	DebtsOwedBy    float32
	TotalIndvConts float32
	OfficeState    string
	OfficeDistrict string
	SpecElection   string
	PrimElection   string
	RunElection    string
	GenElection    string
	GenElectionPct float32
	OtherCmteConts float32
	PtyConts       float32
	CvgEndDate     time.Time
	IndvRefunds    float32
	CmteRefunds    float32
}
