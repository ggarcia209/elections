package cache

import (
	"fmt"

	"github.com/elections/source/donations"
)

// createIndv creates a new Individual object and returns the pointer
func createIndv(id string, cont *donations.Contribution) *donations.Individual {
	donor := donations.Individual{
		ID:            id,
		Name:          cont.Name,
		City:          cont.City,
		State:         cont.State,
		Zip:           cont.Zip,
		Occupation:    cont.Occupation,
		Employer:      cont.Employer,
		Transactions:  []string{},
		TotalOutAmt:   0.0,
		TotalOutTxs:   0.0,
		AvgTxOut:      0.0,
		TotalInAmt:    0.0,
		TotalInTxs:    0.0,
		NetBalance:    0.0,
		RecipientsAmt: make(map[string]float32),
		RecipientsTxs: make(map[string]float32),
		SendersAmt:    make(map[string]float32),
		SendersTxs:    make(map[string]float32),
	}

	return &donor
}

// createOrg creates a new Individual object and returns the pointer
func createOrg(id string, tx interface{}) *donations.Individual {
	// 	fmt.Println("* find org *")
	switch t := tx.(type) {
	case *donations.Contribution:
		return findOrgFromContribution(id, t)
	case *donations.Disbursement:
		return findOrgFromDisbursement(id, t)
	default:
		fmt.Println("FindOrg failed: Invalid interface type")
		return nil
	}
}

func findOrgFromContribution(id string, cont *donations.Contribution) *donations.Individual {
	org := donations.Individual{
		ID:            id,
		Name:          cont.Name,
		City:          cont.City,
		State:         cont.State,
		Zip:           cont.Zip,
		Occupation:    cont.Occupation,
		Employer:      cont.Employer,
		Transactions:  []string{},
		TotalOutAmt:   0.0,
		TotalOutTxs:   0.0,
		AvgTxOut:      0.0,
		TotalInAmt:    0.0,
		TotalInTxs:    0.0,
		AvgTxIn:       0.0,
		RecipientsAmt: make(map[string]float32),
		RecipientsTxs: make(map[string]float32),
		SendersAmt:    make(map[string]float32),
		SendersTxs:    make(map[string]float32),
	}

	return &org
}

func findOrgFromDisbursement(id string, disb *donations.Disbursement) *donations.Individual {
	// create new Indvidual obj if non-existent
	org := donations.Individual{
		ID:            id,
		Name:          disb.Name,
		City:          disb.City,
		State:         disb.State,
		Zip:           disb.Zip,
		Occupation:    "",
		Employer:      "",
		Transactions:  []string{},
		TotalOutAmt:   0.0,
		TotalOutTxs:   0.0,
		AvgTxOut:      0.0,
		TotalInAmt:    0.0,
		TotalInTxs:    0.0,
		AvgTxIn:       0.0,
		RecipientsAmt: make(map[string]float32),
		RecipientsTxs: make(map[string]float32),
		SendersAmt:    make(map[string]float32),
		SendersTxs:    make(map[string]float32),
	}

	return &org
}

// createCmte creates a new CmteTxData object and returns the pointer
func createCmte(ID string) (*donations.Committee, *donations.CmteTxData) {
	cmte := donations.Committee{
		ID:    ID,
		Name:  "Unknown",
		City:  "???",
		State: "???",
		Party: "UNK",
	}
	txData := donations.CmteTxData{
		CmteID:                    ID,
		Party:                     "UNK",
		TopIndvContributorsAmt:    make(map[string]float32),
		TopIndvContributorsTxs:    make(map[string]float32),
		TopCmteOrgContributorsAmt: make(map[string]float32),
		TopCmteOrgContributorsTxs: make(map[string]float32),
		TransferRecsAmt:           make(map[string]float32),
		TransferRecsTxs:           make(map[string]float32),
		TopExpRecipientsAmt:       make(map[string]float32),
		TopExpRecipientsTxs:       make(map[string]float32),
	}
	return &cmte, &txData
}

// createCand creates a new Candidate object and returns the pointer
func createCand(ID string) *donations.Candidate {
	cand := donations.Candidate{
		ID:                   ID,
		Name:                 "Unknown",
		City:                 "???",
		State:                "???",
		Party:                "UNK",
		DirectRecipientsAmts: make(map[string]float32),
		DirectRecipientsTxs:  make(map[string]float32),
		DirectSendersAmts:    make(map[string]float32),
		DirectSendersTxs:     make(map[string]float32),
	}
	return &cand
}
