package donations

// Candidate represents a candidate
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

// CmteLink represents the link between a candidate and their primary committee
type CmteLink struct {
	CandID   string
	ElectnYr int
	fecYr    int
	CmteID   string
	CmteType string
	CmteDsgn string // Committee Designation
	LinkID   string
}
