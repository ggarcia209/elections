// Package donations contains the base objects that are used throughout the application.
// Objects within this package are primarily used for creating, updating, and persisting
// the datasets derived from the input data.
// Struct fields for objects in this file are populated from the data parsed
// from the input .txt files. One object is created for each record.
package donations

import (
	"time"
)

// Contribution represents a contribution, expense, or other transaction
// from a contribution/transactions bulk input file.
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

// Disbursement represents a disbursement transaction
// from an operationg expenses bulk input file.
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
