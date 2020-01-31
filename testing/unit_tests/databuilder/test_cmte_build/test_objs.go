package main

import (
	"fmt"
	"os"

	"github.com/elections/databuilder"
	"github.com/elections/donations"
	"github.com/elections/persist"
)

func main() {
	persist.Init()

	// save objs
	_, err := persist.InitialCacheCand(cands, true)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	_, err = persist.InitialCacheCmte(cmtes, true)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	err = persist.CacheAndPersistIndvDonor(donors)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	err = persist.CacheAndPersistDisbRecipient(recs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// update cmtes
	err = databuilder.IndvContUpdate(ics)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	err = databuilder.CmteContUpdate(ccs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	err = databuilder.CmteDisbUpdate(dbs)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	d00, err := persist.GetIndvDonor("D00")
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	pac00, err := persist.GetCommittee("PAC00")
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}
	pcc00, err := persist.GetCommittee("PCC00")
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	fmt.Println(d00)
	fmt.Println()
	fmt.Println(pac00)
	fmt.Println(pac00.TotalTransfers)
	fmt.Println()
	fmt.Println(pcc00)
	fmt.Println()
	pct := databuilder.FindDonationDirectPct(donor0.RecipientsAmt, pcc00)
	pct2, _ := databuilder.FindDonationTotalPct(donor0.RecipientsAmt, pcc00)
	fmt.Println("direct %: ", pct)
	fmt.Println("total %: ", pct2)

}

func printDonor(id string) error {
	donor, err := persist.GetIndvDonor(id)
	if err != nil {
		fmt.Println("print donor failed: ", err)
		return fmt.Errorf("print donor failed: %v,", err)
	}

	fmt.Printf("ID: %s\nName: %s\nDonations: %v\nTotalDonations: %v\tTotalDonated: %v\tAvgDonation: %v\nRecipentsAmt: %v\nRecipientTxs: %v\n",
		donor.ID, donor.Name, donor.Donations, donor.TotalDonations, donor.TotalDonated, donor.AvgDonation, donor.RecipientsAmt, donor.RecipientsTxs)

	return nil
}

/* initialzie objects as returned from Parse operations */

// individual donors
var donor0 = &donations.Individual{
	ID:             "D00",
	Name:           "John Smith",
	City:           "New York",
	State:          "New York",
	Zip:            "00001",
	Occupation:     "Drug Dealer",
	Employer:       "The Trap",
	Donations:      []string{},
	TotalDonations: 0.0,
	TotalDonated:   0.0,
	AvgDonation:    0.0,
	RecipientsAmt:  make(map[string]float32),
	RecipientsTxs:  make(map[string]float32),
}

var donor1 = &donations.Individual{
	ID:             "D01",
	Name:           "Jose Conseco",
	City:           "Oakland",
	State:          "California",
	Zip:            "94608",
	Occupation:     "Scammer",
	Employer:       "Self-Employed",
	Donations:      []string{},
	TotalDonations: 0.0,
	TotalDonated:   0.0,
	AvgDonation:    0.0,
	RecipientsAmt:  make(map[string]float32),
	RecipientsTxs:  make(map[string]float32),
}

// committee donors
var pac00 = &donations.Committee{
	ID:                    "PAC00",
	Name:                  "Foos Gone Wild",
	City:                  "Los Angeles",
	State:                 "California",
	TotalReceived:         0.0,
	TotalDonations:        0.0,
	AvgDonation:           0.0,
	TopIndvDonorsAmt:      make(map[string]float32),
	TopIndvDonorsTxs:      make(map[string]float32),
	TopIndvDonorThreshold: []interface{}{},
	TopCmteDonorsAmt:      make(map[string]float32),
	TopCmteDonorsTxs:      make(map[string]float32),
	TopCmteDonorThreshold: []interface{}{},
	TotalTransferred:      0.0,
	TotalTransfers:        0.0,
	AvgTransfer:           0.0,
	AffiliatesAmt:         make(map[string]float32),
	AffiliatesTxs:         make(map[string]float32),
	TotalDisbursed:        0.0,
	TotalDisbursements:    0.0,
	AvgDisbursed:          0.0,
	TopDisbRecipientsAmt:  make(map[string]float32),
	TopDisbRecipientsTxs:  make(map[string]float32),
	TopRecThreshold:       []interface{}{},
}

var pac01 = &donations.Committee{
	ID:                    "PAC01",
	Name:                  "Make America Hyphy Again",
	City:                  "Oakland",
	State:                 "California",
	TotalReceived:         0.0,
	TotalDonations:        0.0,
	AvgDonation:           0.0,
	TopIndvDonorsAmt:      make(map[string]float32),
	TopIndvDonorsTxs:      make(map[string]float32),
	TopIndvDonorThreshold: []interface{}{},
	TopCmteDonorsAmt:      make(map[string]float32),
	TopCmteDonorsTxs:      make(map[string]float32),
	TopCmteDonorThreshold: []interface{}{},
	TotalTransferred:      0.0,
	TotalTransfers:        0.0,
	AvgTransfer:           0.0,
	AffiliatesAmt:         make(map[string]float32),
	AffiliatesTxs:         make(map[string]float32),
	TotalDisbursed:        0.0,
	TotalDisbursements:    0.0,
	AvgDisbursed:          0.0,
	TopDisbRecipientsAmt:  make(map[string]float32),
	TopDisbRecipientsTxs:  make(map[string]float32),
	TopRecThreshold:       []interface{}{},
}

var pac02 = &donations.Committee{
	ID:                    "PAC02",
	Name:                  "Epstein Didn't Kill Himself Super PAC",
	City:                  "New York",
	State:                 "New York",
	TotalReceived:         0.0,
	TotalDonations:        0.0,
	AvgDonation:           0.0,
	TopIndvDonorsAmt:      make(map[string]float32),
	TopIndvDonorsTxs:      make(map[string]float32),
	TopIndvDonorThreshold: []interface{}{},
	TopCmteDonorsAmt:      make(map[string]float32),
	TopCmteDonorsTxs:      make(map[string]float32),
	TopCmteDonorThreshold: []interface{}{},
	TotalTransferred:      0.0,
	TotalTransfers:        0.0,
	AvgTransfer:           0.0,
	AffiliatesAmt:         make(map[string]float32),
	AffiliatesTxs:         make(map[string]float32),
	TotalDisbursed:        0.0,
	TotalDisbursements:    0.0,
	AvgDisbursed:          0.0,
	TopDisbRecipientsAmt:  make(map[string]float32),
	TopDisbRecipientsTxs:  make(map[string]float32),
	TopRecThreshold:       []interface{}{},
}

// candidiate committees
var pcc00 = &donations.Committee{
	ID:                    "PCC00",
	Name:                  "Dr. Foo for America",
	City:                  "Seattle",
	State:                 "Washington",
	CandID:                "CAND00",
	TotalReceived:         0.0,
	TotalDonations:        0.0,
	AvgDonation:           0.0,
	TopIndvDonorsAmt:      make(map[string]float32),
	TopIndvDonorsTxs:      make(map[string]float32),
	TopIndvDonorThreshold: []interface{}{},
	TopCmteDonorsAmt:      make(map[string]float32),
	TopCmteDonorsTxs:      make(map[string]float32),
	TopCmteDonorThreshold: []interface{}{},
	TotalTransferred:      0.0,
	TotalTransfers:        0.0,
	AvgTransfer:           0.0,
	AffiliatesAmt:         make(map[string]float32),
	AffiliatesTxs:         make(map[string]float32),
	TotalDisbursed:        0.0,
	TotalDisbursements:    0.0,
	AvgDisbursed:          0.0,
	TopDisbRecipientsAmt:  make(map[string]float32),
	TopDisbRecipientsTxs:  make(map[string]float32),
	TopRecThreshold:       []interface{}{},
}

var pcc01 = &donations.Committee{
	ID:                    "PCC01",
	Name:                  "Christian Genius Billionaire Kanye for President",
	City:                  "Chicago",
	State:                 "Illinois",
	CandID:                "CAND01",
	TotalReceived:         0.0,
	TotalDonations:        0.0,
	AvgDonation:           0.0,
	TopIndvDonorsAmt:      make(map[string]float32),
	TopIndvDonorsTxs:      make(map[string]float32),
	TopIndvDonorThreshold: []interface{}{},
	TopCmteDonorsAmt:      make(map[string]float32),
	TopCmteDonorsTxs:      make(map[string]float32),
	TopCmteDonorThreshold: []interface{}{},
	TotalTransferred:      0.0,
	TotalTransfers:        0.0,
	AvgTransfer:           0.0,
	AffiliatesAmt:         make(map[string]float32),
	AffiliatesTxs:         make(map[string]float32),
	TotalDisbursed:        0.0,
	TotalDisbursements:    0.0,
	AvgDisbursed:          0.0,
	TopDisbRecipientsAmt:  make(map[string]float32),
	TopDisbRecipientsTxs:  make(map[string]float32),
	TopRecThreshold:       []interface{}{},
}

// candidiates
var cand00 = &donations.Candidate{
	ID:               "CAND00",
	Name:             "Dr. Foo",
	PCC:              "PCC00",
	TotalDonations:   0.0,
	TotalRaised:      0.0,
	AvgDonation:      0.0,
	TopIndvDonorsAmt: make(map[string]float32),
	TopIndvDonorsTxs: make(map[string]float32),
	TopIDThreshold:   []interface{}{},
	TopCmteDonorsAmt: make(map[string]float32),
	TopCmteDonorsTxs: make(map[string]float32),
	TopCDThreshold:   []interface{}{},
}

var cand01 = &donations.Candidate{
	ID:               "CAND01",
	Name:             "Christian Genius Billionaire Kanye West",
	PCC:              "PCC01",
	TotalDonations:   0.0,
	TotalRaised:      0.0,
	AvgDonation:      0.0,
	TopIndvDonorsAmt: make(map[string]float32),
	TopIndvDonorsTxs: make(map[string]float32),
	TopIDThreshold:   []interface{}{},
	TopCmteDonorsAmt: make(map[string]float32),
	TopCmteDonorsTxs: make(map[string]float32),
	TopCDThreshold:   []interface{}{},
}

// disbursement recipients
var drec00 = &donations.DisbRecipient{
	ID:                 "DR00",
	Name:               "O Block Security",
	City:               "Chicago",
	State:              "Illinois",
	Disbursements:      []string{},
	TotalDisbursements: 0.0,
	TotalReceived:      0.0,
	AvgReceived:        0.0,
	SendersAmt:         make(map[string]float32),
	SendersTxs:         make(map[string]float32),
}

var drec01 = &donations.DisbRecipient{
	ID:                 "DR01",
	Name:               "Popeyes",
	City:               "Baton Rouge",
	State:              "Louisianna",
	Disbursements:      []string{},
	TotalDisbursements: 0.0,
	TotalReceived:      0.0,
	AvgReceived:        0.0,
	SendersAmt:         make(map[string]float32),
	SendersTxs:         make(map[string]float32),
}

var drec02 = &donations.DisbRecipient{
	ID:                 "DR02",
	Name:               "Magic City",
	City:               "Atlanta",
	State:              "Georgia",
	Disbursements:      []string{},
	TotalDisbursements: 0.0,
	TotalReceived:      0.0,
	AvgReceived:        0.0,
	SendersAmt:         make(map[string]float32),
	SendersTxs:         make(map[string]float32),
}

// individual contributions
var ic00 = &donations.IndvContribution{
	CmteID:  "PAC01",
	Name:    "John Smith",
	TxAmt:   1000.00,
	TxID:    "IC00",
	DonorID: "D00",
}

var ic01 = &donations.IndvContribution{
	CmteID:  "PAC01",
	Name:    "John Smith",
	TxAmt:   500.00,
	TxID:    "IC01",
	DonorID: "D00",
}

var ic02 = &donations.IndvContribution{
	CmteID:  "PCC01",
	Name:    "John Smith",
	TxAmt:   500.00,
	TxID:    "IC02",
	DonorID: "D00",
}

var ic03 = &donations.IndvContribution{
	CmteID:  "PCC00",
	Name:    "Jose Conseco",
	TxAmt:   250.00,
	TxID:    "IC03",
	DonorID: "D01",
}

var ic04 = &donations.IndvContribution{
	CmteID:  "PAC00",
	Name:    "Jose Conseco",
	TxAmt:   500.00,
	TxID:    "IC04",
	DonorID: "D01",
}

var ic05 = &donations.IndvContribution{
	CmteID:  "PAC00",
	Name:    "Jose Conseco",
	TxAmt:   750.00,
	TxID:    "IC05",
	DonorID: "D01",
}

var ic06 = &donations.IndvContribution{
	CmteID:  "PAC02",
	Name:    "Jose Conseco",
	TxAmt:   250.00,
	TxID:    "IC06",
	DonorID: "D01",
}

var ic07 = &donations.IndvContribution{
	CmteID:  "PAC02",
	Name:    "John Smith",
	TxAmt:   750.00,
	TxID:    "IC07",
	DonorID: "D00",
}

var cc00 = &donations.CmteContribution{
	CmteID:  "PCC00",
	Name:    "Foos Gone Wild",
	TxAmt:   500.00,
	OtherID: "PAC00",
	CandID:  "CAND00",
	TxID:    "CC00",
}

var cc01 = &donations.CmteContribution{
	CmteID:  "PAC02",
	Name:    "Foos Gone Wild",
	TxAmt:   250.00,
	OtherID: "PAC00",
	TxID:    "CC01",
}

var cc02 = &donations.CmteContribution{
	CmteID:  "PAC01",
	Name:    "Foos Gone Wild",
	TxAmt:   250.00,
	OtherID: "PAC00",
	TxID:    "CC02",
}

var cc03 = &donations.CmteContribution{
	CmteID:  "PAC02",
	Name:    "Make America Hyphy Again",
	TxAmt:   250.00,
	OtherID: "PAC01",
	TxID:    "CC03",
}

var cc04 = &donations.CmteContribution{
	CmteID:  "PAC00",
	Name:    "Make America Hyphy Again",
	TxAmt:   250.00,
	OtherID: "PAC01",
	TxID:    "CC04",
}

var cc05 = &donations.CmteContribution{
	CmteID:  "PAC02",
	Name:    "Make America Hyphy Again",
	TxAmt:   250.00,
	OtherID: "PAC01",
	TxID:    "CC05",
}

var cc06 = &donations.CmteContribution{
	CmteID:  "PCC00",
	Name:    "Make America Hyphy Again",
	TxAmt:   250.00,
	OtherID: "PAC01",
	TxID:    "CC06",
}

var cc07 = &donations.CmteContribution{
	CmteID:  "PCC01",
	Name:    "Make America Hyphy Again",
	TxAmt:   250.00,
	OtherID: "PAC01",
	TxID:    "CC07",
}

var cc08 = &donations.CmteContribution{
	CmteID:  "PCC01",
	Name:    "Epstein Didn't Kill Himself Super PAC",
	TxAmt:   250.00,
	OtherID: "PAC02",
	TxID:    "CC08",
}

var cc09 = &donations.CmteContribution{
	CmteID:  "PCC00",
	Name:    "Epstein Didn't Kill Himself Super PAC",
	TxAmt:   250.00,
	OtherID: "PAC02",
	TxID:    "CC09",
}

var db00 = &donations.Disbursement{
	CmteID: "PCC00",
	Name:   "O Block Security",
	TxAmt:  500.00,
	RecID:  "DR00",
}

var db01 = &donations.Disbursement{
	CmteID: "PCC00",
	Name:   "Popeyes",
	TxAmt:  200.00,
	RecID:  "DR01",
}

var db02 = &donations.Disbursement{
	CmteID: "PCC01",
	Name:   "O Block Security",
	TxAmt:  500.00,
	RecID:  "DR00",
}

var db03 = &donations.Disbursement{
	CmteID: "PCC01",
	Name:   "Popeyes",
	TxAmt:  100.00,
	RecID:  "DR01",
}

var db04 = &donations.Disbursement{
	CmteID: "PCC01",
	Name:   "Magic City",
	TxAmt:  500.00,
	RecID:  "DR02",
}

var donors = []*donations.Individual{donor0, donor1}

var cmtes = []*donations.Committee{pac00, pac01, pac02, pcc00, pcc01}

var cands = []*donations.Candidate{cand00, cand01}

var recs = []*donations.DisbRecipient{drec00, drec01, drec02}

var ics = []*donations.IndvContribution{ic00, ic01, ic02, ic03, ic04, ic05, ic06, ic07}

var ccs = []*donations.CmteContribution{cc00, cc01, cc02, cc03, cc04, cc05, cc06, cc07, cc08, cc09}

var dbs = []*donations.Disbursement{db00, db01, db02, db03, db04}
