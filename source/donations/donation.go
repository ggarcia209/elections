package donations

import (
	"time"
)

// Contribution represents a contribution
type Contribution struct {
	CmteID     string // filing committee
	AmndtInd   string // ammendment indicator
	ReportType string
	TxPGI      string // transaction primary-general indicator
	ImgNum     string // image number
	TxType     string
	EntityType string
	Name       string
	City       string
	State      string
	Zip        string
	Employer   string
	Occupation string
	TxDate     time.Time
	TxAmt      float32 // transaction amount
	OtherID    string  // Cmte/Cand/Org/Indv ID for recipient/sender
	TxID       string
	FileNum    int
	MemoCode   string
	MemoText   string
	SubID      int // FEC record number, unique row ID
}

// Disbursement represents a disbursement
type Disbursement struct {
	CmteID       string
	amndtInd     string
	RptYr        int
	RptTp        string
	ImgNum       string
	LineNum      string
	FormTp       string
	SchedTp      string
	Name         string
	City         string
	State        string
	Zip          string
	TxDate       time.Time
	TxAmt        float32
	TxPGI        string
	Purpose      string
	Category     string
	CategoryDesc string
	MemoCode     string
	MemoTxt      string
	EntityType   string
	SubID        int
	FileNum      int
	TxID         string
	BackRefTxID  string
	RecID        string
}

// DEPRECATED

/*
// Donation represents a contribution from either an Individual or a Committee
type Donation interface{}

// IndvContribution reperesents a contribution from an individual
type IndvContribution struct {
	CmteID     string
	AmndtInd   string // ammendment indicator
	ReportType string
	TxPGI      string // transaction primary-general indicator
	imgNum     string // image number
	TxType     string
	EntityType string
	Name       string
	City       string
	State      string
	Zip        string
	Employer   string
	Occupation string
	TxDate     time.Time
	TxAmt      float32 // transaction amount
	OtherID    string
	TxID       string
	FileNum    int
	MemoCode   string
	MemoText   string
	SubID      int // FEC record number, unique row ID
	DonorID    string
}

// CmteContribution represents a contribution from a committee
type CmteContribution struct {
	CmteID     string
	AmndtInd   string // ammendment indicator
	ReportType string
	TxPGI      string // transaction primary-general indicator
	imgNum     string // image number
	TxType     string
	EntityType string
	Name       string
	City       string
	State      string
	Zip        string
	Employer   string
	Occupation string
	TxDate     time.Time
	TxAmt      float32 // transaction amount
	OtherID    string
	CandID     string // candidate ID
	TxID       string
	FileNum    int
	MemoCode   string
	MemoText   string
	SubID      int // FEC record number, unique row ID
}
*/
