package testDB

import (
	"github.com/elections/donations"
)

/* TEST OBJECTS */
// 7/24/20 - updated to incorporate year key & TopOverall
// Top Overall Categories:
//   cmte_recs_all, cmte_recs_r, cmte_recs_d, cmte_recs_na, cmte_recs_misc,
//   cmte_donors_all, cmte_donors_r, cmte_donors_d, cmte_donors_na, cmte_donors_misc,
//   cmte_exp_all, cmte_exp_r, cmte_exp_d, cmte_exp_na, cmte_exp_misc,
//   cand_all, cand_r, cand_d, cand_na, cand_misc,
//   cand_exp_all, cand_exp_r, cand_exp_d, cand_exp_na, cand_exp_misc,
//   indv, org_conts, org_recs,
var DbSim = map[string]map[string]map[string]interface{}{
	"2018": map[string]map[string]interface{}{
		"individuals":   map[string]interface{}{"indv8": &Indv8, "indv10": &Indv10, "indv11": &Indv11, "indv12": &Indv12, "indv13": &Indv13, "indv14": &Indv14, "indv15": &Indv15, "indv16": &Indv16},
		"organizations": map[string]interface{}{"org01": &Org01, "org02": &Org02, "org03": &Org03},
		"cmte_tx_data":  map[string]interface{}{"Cmte00": &Filer, "Cmte01": &Cmte01, "Cmte02": &Cmte02, "Cmte03": &Cmte03},
		"candidates":    map[string]interface{}{"Pcand00": &Cand00, "Pcand01": &Cand01, "Scand02": &Cand02, "Hcand03": &Cand03},
		"top_overall": map[string]interface{}{
			"cmte_recs_all": &donations.TopOverallData{}, "cmte_recs_r": &donations.TopOverallData{}, "cmte_recs_d": &donations.TopOverallData{}, "cmte_recs_na": &donations.TopOverallData{}, "cmte_recs_misc": &donations.TopOverallData{},
			"cmte_donors_all": &donations.TopOverallData{}, "cmte_donors_r": &donations.TopOverallData{}, "cmte_donors_d": &donations.TopOverallData{}, "cmte_donors_na": &donations.TopOverallData{}, "cmte_donors_misc": &donations.TopOverallData{},
			"cmte_exp_all": &donations.TopOverallData{}, "cmte_exp_r": &donations.TopOverallData{}, "cmte_exp_d": &donations.TopOverallData{}, "cmte_exp_na": &donations.TopOverallData{}, "cmte_exp_misc": &donations.TopOverallData{},
			"cand_all": &donations.TopOverallData{}, "cand_r": &donations.TopOverallData{}, "cand_d": &donations.TopOverallData{}, "cand_na": &donations.TopOverallData{}, "cand_misc": &donations.TopOverallData{},
			"cand_exp_all": &donations.TopOverallData{}, "cand_exp_r": &donations.TopOverallData{}, "cand_exp_d": &donations.TopOverallData{}, "cand_exp_na": &donations.TopOverallData{}, "cand_exp_misc": &donations.TopOverallData{},
			"indv": &donations.TopOverallData{}, "org_conts": &donations.TopOverallData{}, "org_recs": &donations.TopOverallData{},
		},
	},
}

var Conts = []*donations.Contribution{&tx1, &tx2, &tx3, &tx4, &tx5, &tx6, &tx7, &tx8, &tx9, &tx10, &tx11, &tx12, &tx13, &tx14, &tx15,
	&tx16, &tx165, &tx17, &tx175, &tx18, &tx185, &tx19, &tx195, &tx20, &tx205, &tx21, &tx215, &tx22, &tx225, &tx23, &tx24, &tx25, &tx26, &tx27, &tx28, &tx29}

var Disbs = []*donations.Disbursement{&tx30, &tx31, &tx32, &tx33}

var Transfers = []*donations.Contribution{&tx20, &tx21, &tx24}

var ExpensesConts = []*donations.Contribution{&tx2, &tx14, &tx22, &tx29}
var ExpensesDisbs = []*donations.Disbursement{&tx30, &tx31, &tx32, &tx33}

// SIZE LIMIT == 5 for object maps

// filing committee
var Filer = donations.CmteTxData{
	CmteID:                    "cmte00",
	CandID:                    "Pcand00",
	Party:                     "IND",
	TopIndvContributorsAmt:    instantiate(),
	TopIndvContributorsTxs:    instantiate(),
	TopCmteOrgContributorsAmt: instantiate(),
	TopCmteOrgContributorsTxs: instantiate(),
	TransferRecsAmt:           instantiate(),
	TransferRecsTxs:           instantiate(),
	TopExpRecipientsAmt:       instantiate(),
	TopExpRecipientsTxs:       instantiate(),
}

// committee contributors/recipients
var Cmte01 = donations.CmteTxData{
	CmteID:                    "Cmte01",
	Party:                     "DEM",
	TopIndvContributorsAmt:    instantiate(),
	TopIndvContributorsTxs:    instantiate(),
	TopCmteOrgContributorsAmt: instantiate(),
	TopCmteOrgContributorsTxs: instantiate(),
	TransferRecsAmt:           instantiate(),
	TransferRecsTxs:           instantiate(),
	TopExpRecipientsAmt:       instantiate(),
	TopExpRecipientsTxs:       instantiate(),
}

var Cmte02 = donations.CmteTxData{
	CmteID:                    "Cmte02",
	Party:                     "IND",
	TopIndvContributorsAmt:    instantiate(),
	TopIndvContributorsTxs:    instantiate(),
	TopCmteOrgContributorsAmt: instantiate(),
	TopCmteOrgContributorsTxs: instantiate(),
	TransferRecsAmt:           instantiate(),
	TransferRecsTxs:           instantiate(),
	TopExpRecipientsAmt:       instantiate(),
	TopExpRecipientsTxs:       instantiate(),
}

var Cmte03 = donations.CmteTxData{
	CmteID:                    "Cmte03",
	Party:                     "IND",
	TopIndvContributorsAmt:    instantiate(),
	TopIndvContributorsTxs:    instantiate(),
	TopCmteOrgContributorsAmt: instantiate(),
	TopCmteOrgContributorsTxs: instantiate(),
	TransferRecsAmt:           instantiate(),
	TransferRecsTxs:           instantiate(),
	TopExpRecipientsAmt:       instantiate(),
	TopExpRecipientsTxs:       instantiate(),
}

// Filing Committee Linked Candidate
var Cand00 = donations.Candidate{
	ID:                   "PCand00",
	Party:                "IND",
	DirectRecipientsAmts: instantiate(),
	DirectRecipientsTxs:  instantiate(),
	DirectSendersAmts:    instantiate(),
	DirectSendersTxs:     instantiate(),
}

// candidate contributors/recipients
var Cand01 = donations.Candidate{
	ID:                   "PCand01",
	Party:                "IND",
	DirectRecipientsAmts: instantiate(),
	DirectRecipientsTxs:  instantiate(),
	DirectSendersAmts:    instantiate(),
	DirectSendersTxs:     instantiate(),
}

var Cand02 = donations.Candidate{
	ID:                   "SCand02",
	Party:                "REP",
	DirectRecipientsAmts: instantiate(),
	DirectRecipientsTxs:  instantiate(),
	DirectSendersAmts:    instantiate(),
	DirectSendersTxs:     instantiate(),
}

var Cand03 = donations.Candidate{
	ID:                   "HCand03",
	Party:                "DEM",
	DirectRecipientsAmts: instantiate(),
	DirectRecipientsTxs:  instantiate(),
	DirectSendersAmts:    instantiate(),
	DirectSendersTxs:     instantiate(),
}

// organization contributors/recipients
var Org01 = donations.Organization{
	ID:            "org01",
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
}

var Org02 = donations.Organization{
	ID:            "org02",
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
}

var Org03 = donations.Organization{
	ID:            "org03",
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
}

// individual contributors/recipients
// Top 5 individuals: indv8: 300, indv14: 250, indv15: 220, indv12: 200, indv10: 175
var indv4 = donations.Individual{
	ID:            "indv4",
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
}

var Indv8 = donations.Individual{
	ID:            "indv8",
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
}

// total out = 175, total in = 175, bal = 0
var Indv10 = donations.Individual{
	ID:            "indv10",
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
}

var Indv11 = donations.Individual{
	ID:            "indv11",
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
}

var Indv12 = donations.Individual{
	ID:            "indv12",
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
}

var Indv13 = donations.Individual{
	ID:            "indv13",
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
}

var Indv14 = donations.Individual{
	ID:            "indv14",
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
}

var Indv15 = donations.Individual{
	ID:            "indv15",
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
}

var Indv16 = donations.Individual{
	ID:            "indv16",
	RecipientsAmt: instantiate(),
	RecipientsTxs: instantiate(),
	SendersAmt:    instantiate(),
	SendersTxs:    instantiate(),
}

// Individual contributions
var tx1 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv10",
	TxAmt:      175,
	TxType:     "15",
	TxID:       "tx1",
	MemoCode:   "",
}

// refund tx
var tx2 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv10",
	TxAmt:      175,
	TxType:     "22Y",
	TxID:       "tx2",
	MemoCode:   "",
}

var tx3 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv11",
	TxAmt:      100,
	TxType:     "15",
	TxID:       "tx3",
	MemoCode:   "",
}

var tx4 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv11",
	TxAmt:      150,
	TxType:     "15",
	TxID:       "tx4",
	MemoCode:   "",
}

var tx5 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv12",
	TxAmt:      200,
	TxType:     "15",
	TxID:       "tx5",
	MemoCode:   "",
}

var tx6 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv13",
	TxAmt:      50,
	TxType:     "15",
	TxID:       "tx6",
	MemoCode:   "",
}

var tx7 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv14",
	TxAmt:      250,
	TxType:     "15",
	TxID:       "tx7",
	MemoCode:   "",
}

var tx8 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv15",
	TxAmt:      220,
	TxType:     "15",
	TxID:       "tx8",
	MemoCode:   "",
}

var tx9 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv16",
	TxAmt:      80,
	TxType:     "15",
	TxID:       "tx9",
	MemoCode:   "",
}

var tx10 = donations.Contribution{
	CmteID:     "Cmte00",
	Occupation: "worker",
	OtherID:    "indv8",
	TxAmt:      300,
	TxType:     "15J",
	TxID:       "tx10",
	MemoCode:   "X",
}

// organization transactions
var tx11 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "org01",
	TxAmt:    300,
	TxType:   "15",
	TxID:     "tx11",
	MemoCode: "",
}

var tx12 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "org02",
	TxAmt:    200,
	TxType:   "15",
	TxID:     "tx12",
	MemoCode: "",
}

var tx13 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "org03",
	TxAmt:    100,
	TxType:   "15",
	TxID:     "tx13",
	MemoCode: "",
}

// refund tx
var tx14 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "org03",
	TxAmt:    100,
	TxType:   "22Y",
	TxID:     "tx14",
	MemoCode: "",
}

// memo tx
var tx15 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "org01",
	TxAmt:    100,
	TxType:   "15J",
	TxID:     "tx15",
	MemoCode: "X",
}

// Committee transactions
// transfer to filing cmte
var tx16 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Cmte01",
	TxAmt:    245,
	TxType:   "18G",
	TxID:     "tx16",
	MemoCode: "",
}

// tx 16 corresponding
var tx165 = donations.Contribution{
	CmteID:   "Cmte01",
	OtherID:  "Cmte00",
	TxAmt:    245,
	TxType:   "24G",
	TxID:     "tx165",
	MemoCode: "",
}

var tx17 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Cmte02",
	TxAmt:    150,
	TxType:   "18G",
	TxID:     "tx17",
	MemoCode: "",
}

var tx175 = donations.Contribution{
	CmteID:   "Cmte02",
	OtherID:  "Cmte00",
	TxAmt:    150,
	TxType:   "24G",
	TxID:     "tx175",
	MemoCode: "",
}

var tx18 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Cmte03",
	TxAmt:    350,
	TxType:   "18G",
	TxID:     "tx18",
	MemoCode: "",
}

// var 18 corresponding tx
var tx185 = donations.Contribution{
	CmteID:   "Cmte03",
	OtherID:  "Cmte00",
	TxAmt:    350,
	TxType:   "24G",
	TxID:     "tx185",
	MemoCode: "",
}

// refund to filing cmte
var tx19 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Cmte01",
	TxAmt:    250,
	TxType:   "17R",
	TxID:     "tx19",
	MemoCode: "",
}

// tx19 corresponding tx
var tx195 = donations.Contribution{
	CmteID:   "Cmte01",
	OtherID:  "Cmte00",
	TxAmt:    250,
	TxType:   "22Z",
	TxID:     "tx195",
	MemoCode: "",
}

// transfers/refunds from filing cmte
// tramsfer
var tx20 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Cmte03",
	TxAmt:    225,
	TxType:   "24G", // transfer
	TxID:     "tx20",
	MemoCode: "",
}

// tx20 corresponding tx
var tx205 = donations.Contribution{
	CmteID:   "Cmte03",
	OtherID:  "Cmte00",
	TxAmt:    225,
	TxType:   "18G", // transfer
	TxID:     "tx205",
	MemoCode: "",
}

// transfer
var tx21 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Cmte02",
	TxAmt:    150,
	TxType:   "24G", // transfer
	TxID:     "tx21",
	MemoCode: "",
}

// tx215 corresponding tx
var tx215 = donations.Contribution{
	CmteID:   "Cmte02",
	OtherID:  "Cmte00",
	TxAmt:    150,
	TxType:   "18G", // transfer
	TxID:     "tx215",
	MemoCode: "",
}

// expense
var tx22 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Cmte01",
	TxAmt:    110,
	TxType:   "40Z", // convention account disbursement
	TxID:     "tx22",
	MemoCode: "",
}

// tx225 corresponding tx
var tx225 = donations.Contribution{
	CmteID:   "Cmte01",
	OtherID:  "Cmte00",
	TxAmt:    110,
	TxType:   "18G",
	TxID:     "tx225",
	MemoCode: "",
}

// memo tx
var tx23 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Cmte03",
	TxAmt:    117,
	TxType:   "18J",
	TxID:     "tx23",
	MemoCode: "X",
}

// corresponding txs
var tx24 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Cmte01",
	TxAmt:    124,
	TxType:   "24G", // transfer out
	TxID:     "tx24",
	MemoCode: "",
}

var tx25 = donations.Contribution{
	CmteID:   "Cmte01",
	OtherID:  "Cmte00",
	TxAmt:    124,
	TxType:   "18G", // transfer in
	TxID:     "tx25",
	MemoCode: "",
}

// candidate txs
var tx26 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Pcand01",
	TxAmt:    325,
	TxType:   "15C",
	TxID:     "tx26",
	MemoCode: "",
}

var tx27 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Scand02",
	TxAmt:    225,
	TxType:   "15C",
	TxID:     "tx27",
	MemoCode: "",
}

var tx28 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Hcand03",
	TxAmt:    225,
	TxType:   "15C",
	TxID:     "tx28",
	MemoCode: "",
}

// cand refund
var tx29 = donations.Contribution{
	CmteID:   "Cmte00",
	OtherID:  "Hcand03",
	TxAmt:    225,
	TxType:   "22Z",
	TxID:     "tx29",
	MemoCode: "",
}

// disbursements
var tx30 = donations.Disbursement{
	CmteID: "Cmte00",
	RecID:  "org01",
	TxAmt:  115,
	TxID:   "tx30",
}

var tx31 = donations.Disbursement{
	CmteID: "Cmte00",
	RecID:  "org02",
	TxAmt:  215,
	TxID:   "tx31",
}

var tx32 = donations.Disbursement{
	CmteID: "Cmte00",
	RecID:  "org03",
	TxAmt:  315,
	TxID:   "tx32",
}

// transaction between TopExpThreshold range

// transaction added to existing entry in threshold range
var tx33 = donations.Disbursement{
	CmteID: "Cmte00",
	RecID:  "org01",
	TxAmt:  222,
	TxID:   "tx33",
}

/* END TEST OBJECTS */

func instantiate() map[string]float32 {
	m := make(map[string]float32)
	return m
}
