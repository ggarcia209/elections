package persist

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/protobuf"

	"github.com/golang/protobuf/proto"
)

// convIndvToProto encodes LogData structs as protocol buffers
func encodeIndv(indv donations.Individual) ([]byte, error) { // move conversions to protobuf package?
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
		fmt.Println("encodeIndv failed: ", err)
		return nil, fmt.Errorf("encodeIndv failed: %v", err)
	}
	return data, nil
}

func decodeIndv(data []byte) (donations.Individual, error) {
	indv := &protobuf.Individual{}
	err := proto.Unmarshal(data, indv)
	if err != nil {
		fmt.Println("decodeIndv failed: ", err)
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

// convIndvToProto encodes LogData structs as protocol buffers
func encodeOrg(indv donations.Organization) ([]byte, error) { // move conversions to protobuf package?
	entry := &protobuf.Organization{
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
		fmt.Println("encodeIndv failed: ", err)
		return nil, fmt.Errorf("encodeIndv failed: %v", err)
	}
	return data, nil
}

func decodeOrg(data []byte) (donations.Organization, error) {
	indv := &protobuf.Organization{}
	err := proto.Unmarshal(data, indv)
	if err != nil {
		fmt.Println("decodeIndv failed: ", err)
		return donations.Organization{}, fmt.Errorf("decodeIndv failed: %v", err)
	}

	entry := donations.Organization{
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
