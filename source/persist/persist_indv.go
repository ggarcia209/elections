// Package persist contains operations for reading and writing disk data.
// Most operations in this package are intended to be performed on the
// admin local machine and are not intended to be used in the service logic.
// This file contains operations for encoding/decoding protobufs for the
// donations.Individual object.
package persist

import (
	"fmt"

	"github.com/elections/source/donations"
	"github.com/elections/source/protobuf"

	"github.com/golang/protobuf/proto"
)

func encodeIndv(indv donations.Individual) ([]byte, error) {
	entry := &protobuf.Individual{
		ID:            indv.ID,
		Name:          indv.Name,
		City:          indv.City,
		State:         indv.State,
		Zip:           indv.Zip,
		Occupation:    indv.Occupation,
		Employer:      indv.Employer,
		Transactions:  indv.Transactions,
		TotalOutAmt:   indv.TotalOutAmt,
		TotalOutTxs:   indv.TotalOutTxs,
		AvgTxOut:      indv.AvgTxOut,
		TotalInAmt:    indv.TotalInAmt,
		TotalInTxs:    indv.TotalInTxs,
		AvgTxIn:       indv.AvgTxIn,
		NetBalance:    indv.NetBalance,
		RecipientsAmt: indv.RecipientsAmt,
		RecipientsTxs: indv.RecipientsTxs,
		SendersAmt:    indv.SendersAmt,
		SendersTxs:    indv.SendersTxs,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeIndv failed: %v", err)
	}
	return data, nil
}

func decodeIndv(data []byte) (donations.Individual, error) {
	indv := &protobuf.Individual{}
	err := proto.Unmarshal(data, indv)
	if err != nil {
		fmt.Println(err)
		return donations.Individual{}, fmt.Errorf("decodeIndv failed: %v", err)
	}

	entry := donations.Individual{
		ID:            indv.GetID(),
		Name:          indv.GetName(),
		City:          indv.GetCity(),
		State:         indv.GetState(),
		Zip:           indv.GetZip(),
		Occupation:    indv.GetOccupation(),
		Employer:      indv.GetEmployer(),
		Transactions:  indv.GetTransactions(),
		TotalInAmt:    indv.GetTotalInAmt(),
		TotalInTxs:    indv.GetTotalInTxs(),
		AvgTxIn:       indv.GetAvgTxIn(),
		TotalOutAmt:   indv.GetTotalOutAmt(),
		TotalOutTxs:   indv.GetTotalOutTxs(),
		AvgTxOut:      indv.GetAvgTxOut(),
		NetBalance:    indv.GetNetBalance(),
		RecipientsAmt: indv.GetRecipientsAmt(),
		RecipientsTxs: indv.GetRecipientsTxs(),
		SendersAmt:    indv.GetSendersAmt(),
		SendersTxs:    indv.GetSendersTxs(),
	}

	return entry, nil
}
