package persist

import (
	"fmt"

	"github.com/elections/source/donations"
	"github.com/elections/source/protobuf"

	"github.com/golang/protobuf/proto"
)

// ecnodeCmte/decodeCmte encodes/decodes Committee structs as protocol buffers
func encodeCmte(cmte donations.Committee) ([]byte, error) { // move conversions to protobuf package?
	entry := &protobuf.Committee{
		ID:           cmte.ID,
		Name:         cmte.Name,
		TresName:     cmte.TresName,
		City:         cmte.City,
		State:        cmte.State,
		Zip:          cmte.Zip,
		Designation:  cmte.Designation,
		Type:         cmte.Type,
		Party:        cmte.Party,
		FilingFreq:   cmte.FilingFreq,
		OrgType:      cmte.OrgType,
		ConnectedOrg: cmte.ConnectedOrg,
		CandID:       cmte.CandID,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeCmte failed: ", err)
		return nil, fmt.Errorf("encodeCmte failed: %v", err)
	}
	return data, nil
}

func encodeCmteTxData(data donations.CmteTxData) ([]byte, error) {
	entry := &protobuf.CmteTxData{
		CmteID:                         data.CmteID,
		CandID:                         data.CandID,
		Party:                          data.Party,
		ContributionsInAmt:             data.ContributionsInAmt,
		ContributionsInTxs:             data.ContributionsInTxs,
		AvgContributionIn:              data.AvgContributionIn,
		OtherReceiptsInAmt:             data.OtherReceiptsInAmt,
		OtherReceiptsInTxs:             data.OtherReceiptsInTxs,
		AvgOtherIn:                     data.AvgOtherIn,
		TotalIncomingAmt:               data.TotalIncomingAmt,
		TotalIncomingTxs:               data.TotalIncomingTxs,
		AvgIncoming:                    data.AvgIncoming,
		TransfersAmt:                   data.TransfersAmt,
		TransfersTxs:                   data.TransfersTxs,
		AvgTransfer:                    data.AvgTransfer,
		ExpendituresAmt:                data.ExpendituresAmt,
		ExpendituresTxs:                data.ExpendituresTxs,
		AvgExpenditure:                 data.AvgExpenditure,
		TotalOutgoingAmt:               data.TotalOutgoingAmt,
		TotalOutgoingTxs:               data.TotalOutgoingTxs,
		AvgOutgoing:                    data.AvgOutgoing,
		NetBalance:                     data.NetBalance,
		TopIndvContributorsAmt:         data.TopIndvContributorsAmt,
		TopIndvContributorsTxs:         data.TopIndvContributorsTxs,
		TopIndvContributorThreshold:    encodeCmteThreshold(data.TopIndvContributorThreshold),
		TopCmteOrgContributorsAmt:      data.TopCmteOrgContributorsAmt,
		TopCmteOrgContributorsTxs:      data.TopCmteOrgContributorsTxs,
		TopCmteOrgContributorThreshold: encodeCmteThreshold(data.TopCmteOrgContributorThreshold),
		TransferRecsAmt:                data.TransferRecsAmt,
		TransferRecsTxs:                data.TransferRecsTxs,
		TopExpRecipientsAmt:            data.TopExpRecipientsAmt,
		TopExpRecipientsTxs:            data.TopExpRecipientsTxs,
		TopExpThreshold:                encodeCmteThreshold(data.TopExpThreshold),
	}
	bytes, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeCmte failed: ", err)
		return nil, fmt.Errorf("encodeCmte failed: %v", err)
	}
	return bytes, nil
}

func encodeCmteThreshold(entries []interface{}) []*protobuf.CmteEntry {
	var es []*protobuf.CmteEntry
	for _, intf := range entries {
		e := intf.(*donations.Entry)
		entry := &protobuf.CmteEntry{
			ID:    e.ID,
			Total: e.Total,
		}
		es = append(es, entry)
	}
	return es
}

func decodeCmte(data []byte) (donations.Committee, error) {
	cmte := &protobuf.Committee{}
	err := proto.Unmarshal(data, cmte)
	if err != nil {
		fmt.Println("decodeCmte failed: ", err)
		return donations.Committee{}, fmt.Errorf("decodeCmte failed: %v", err)
	}

	entry := donations.Committee{
		ID:           cmte.GetID(),
		Name:         cmte.GetName(),
		TresName:     cmte.GetTresName(),
		City:         cmte.GetCity(),
		State:        cmte.GetState(),
		Zip:          cmte.GetZip(),
		Designation:  cmte.GetDesignation(),
		Type:         cmte.GetType(),
		Party:        cmte.GetParty(),
		FilingFreq:   cmte.GetFilingFreq(),
		OrgType:      cmte.GetOrgType(),
		ConnectedOrg: cmte.GetConnectedOrg(),
		CandID:       cmte.GetCandID(),
	}

	return entry, nil
}

func decodeCmteTxData(input []byte) (donations.CmteTxData, error) {
	data := &protobuf.CmteTxData{}
	err := proto.Unmarshal(input, data)
	if err != nil {
		fmt.Println("decodeCmteTxData failed: ", err)
		return donations.CmteTxData{}, fmt.Errorf("decodeCmteTxData failed: %v", err)
	}

	entry := donations.CmteTxData{
		CmteID:                         data.GetCmteID(),
		CandID:                         data.GetCandID(),
		Party:                          data.GetParty(),
		ContributionsInAmt:             data.GetContributionsInAmt(),
		ContributionsInTxs:             data.GetContributionsInTxs(),
		AvgContributionIn:              data.GetAvgContributionIn(),
		OtherReceiptsInAmt:             data.GetOtherReceiptsInAmt(),
		OtherReceiptsInTxs:             data.GetOtherReceiptsInTxs(),
		AvgOtherIn:                     data.GetAvgOtherIn(),
		TotalIncomingAmt:               data.GetTotalIncomingAmt(),
		TotalIncomingTxs:               data.GetTotalIncomingTxs(),
		AvgIncoming:                    data.GetAvgIncoming(),
		TransfersAmt:                   data.GetTransfersAmt(),
		TransfersTxs:                   data.GetTransfersTxs(),
		AvgTransfer:                    data.GetAvgTransfer(),
		ExpendituresAmt:                data.GetExpendituresAmt(),
		ExpendituresTxs:                data.GetExpendituresTxs(),
		AvgExpenditure:                 data.GetAvgExpenditure(),
		TotalOutgoingAmt:               data.GetTotalOutgoingAmt(),
		TotalOutgoingTxs:               data.GetTotalOutgoingTxs(),
		AvgOutgoing:                    data.GetAvgOutgoing(),
		NetBalance:                     data.GetNetBalance(),
		TopIndvContributorsAmt:         data.GetTopIndvContributorsAmt(),
		TopIndvContributorsTxs:         data.GetTopIndvContributorsTxs(),
		TopIndvContributorThreshold:    decodeCmteThreshold(data.GetTopIndvContributorThreshold()),
		TopCmteOrgContributorsAmt:      data.GetTopCmteOrgContributorsAmt(),
		TopCmteOrgContributorsTxs:      data.GetTopCmteOrgContributorsTxs(),
		TopCmteOrgContributorThreshold: decodeCmteThreshold(data.GetTopCmteOrgContributorThreshold()),
		TransferRecsAmt:                data.GetTransferRecsAmt(),
		TransferRecsTxs:                data.GetTransferRecsTxs(),
		TopExpRecipientsAmt:            data.GetTopExpRecipientsAmt(),
		TopExpRecipientsTxs:            data.GetTopExpRecipientsTxs(),
		TopExpThreshold:                decodeCmteThreshold(data.GetTopExpThreshold()),
	}

	return entry, nil
}

func decodeCmteThreshold(es []*protobuf.CmteEntry) []interface{} {
	var entries []interface{}
	for _, e := range es {
		entry := donations.Entry{
			ID:    e.GetID(),
			Total: e.GetTotal(),
		}
		entries = append(entries, &entry)
	}
	return entries
}
