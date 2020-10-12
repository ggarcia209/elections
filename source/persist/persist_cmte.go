// Package persist contains operations for reading and writing disk data.
// Most operations in this package are intended to be performed on the
// admin local machine and are not intended to be used in the service logic.
// This file contains operations for encoding/decoding protobufs for the
// donations.Committee, CmteTxData, and CmteFinancials objects.
package persist

import (
	"fmt"

	"github.com/elections/source/donations"
	"github.com/elections/source/protobuf"

	"github.com/golang/protobuf/proto"
)

// ecnodeCmte/decodeCmte encodes/decodes Committee structs as protocol buffers
func encodeCmte(cmte donations.Committee) ([]byte, error) {
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
		fmt.Println(err)
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
		TransfersList:                  data.TransfersList,
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
		fmt.Println(err)
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
		fmt.Println(err)
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
		fmt.Println(err)
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
		TransfersList:                  data.GetTransfersList(),
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

// ecnodeCmte/decodeCmte encodes/decodes Committee structs as protocol buffers
func encodeCmteFinancials(cmte donations.CmteFinancials) ([]byte, error) { // move conversions to protobuf package?
	entry := &protobuf.CmteFinancials{
		CmteID:          cmte.CmteID,
		TotalReceipts:   cmte.TotalReceipts,
		TxsFromAff:      cmte.TxsFromAff,
		IndvConts:       cmte.IndvConts,
		OtherConts:      cmte.OtherConts,
		CandCont:        cmte.CandCont,
		TotalLoans:      cmte.TotalLoans,
		TotalDisb:       cmte.TotalDisb,
		TxToAff:         cmte.TxToAff,
		IndvRefunds:     cmte.IndvRefunds,
		OtherRefunds:    cmte.OtherRefunds,
		LoanRepay:       cmte.LoanRepay,
		CashBOP:         cmte.CashBOP,
		CashCOP:         cmte.CashCOP,
		DebtsOwed:       cmte.DebtsOwed,
		NonFedTxsRecvd:  cmte.NonFedTxsRecvd,
		ContToOtherCmte: cmte.ContToOtherCmte,
		IndExp:          cmte.IndExp,
		PartyExp:        cmte.PartyExp,
		NonFedSharedExp: cmte.NonFedSharedExp,
	}
	// Temp. deprecated - error parsing time / time not used in current version.
	/* ts, err := ptypes.TimestampProto(cmte.CovgEndDate)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeCmteFinancials failed: %v", err)
	}
	entry.CovgEndDate = ts */
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeCmteFinancials failed: %v", err)
	}
	return data, nil
}

func decodeCmteFinancials(data []byte) (donations.CmteFinancials, error) {
	cmte := &protobuf.CmteFinancials{}
	err := proto.Unmarshal(data, cmte)
	if err != nil {
		fmt.Println(err)
		return donations.CmteFinancials{}, fmt.Errorf("decodeCmteFinancials failed: %v", err)
	}

	entry := donations.CmteFinancials{
		CmteID:          cmte.GetCmteID(),
		TotalReceipts:   cmte.GetTotalReceipts(),
		TxsFromAff:      cmte.GetTxsFromAff(),
		IndvConts:       cmte.GetIndvConts(),
		OtherConts:      cmte.GetOtherConts(),
		CandCont:        cmte.GetCandCont(),
		TotalLoans:      cmte.GetTotalLoans(),
		TotalDisb:       cmte.GetTotalDisb(),
		TxToAff:         cmte.GetTxToAff(),
		IndvRefunds:     cmte.GetIndvRefunds(),
		OtherRefunds:    cmte.GetOtherRefunds(),
		LoanRepay:       cmte.GetLoanRepay(),
		CashBOP:         cmte.GetCashBOP(),
		CashCOP:         cmte.GetCashCOP(),
		DebtsOwed:       cmte.GetDebtsOwed(),
		NonFedTxsRecvd:  cmte.GetNonFedTxsRecvd(),
		ContToOtherCmte: cmte.GetContToOtherCmte(),
		IndExp:          cmte.GetIndExp(),
		PartyExp:        cmte.GetPartyExp(),
		NonFedSharedExp: cmte.GetNonFedSharedExp(),
	}
	/* ts, err := ptypes.Timestamp(cmte.GetCovgEndDate())
	if err != nil {
		fmt.Println(err)
		return donations.CmteFinancials{}, fmt.Errorf("decodeCmteFinancials failed: %v", err)
	}
	entry.CovgEndDate = ts */

	return entry, nil
}
