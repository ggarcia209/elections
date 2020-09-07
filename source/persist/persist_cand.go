package persist

import (
	"fmt"

	"github.com/elections/source/donations"
	"github.com/elections/source/protobuf"

	"github.com/golang/protobuf/proto"
)

// convIndvToProto encodes LogData structs as protocol buffers
func encodeCand(cand donations.Candidate) ([]byte, error) { // move conversions to protobuf package?
	entry := &protobuf.Candidate{
		ID:                   cand.ID,
		Name:                 cand.Name,
		Party:                cand.Party,
		OfficeState:          cand.OfficeState,
		Office:               cand.Office,
		PCC:                  cand.PCC,
		City:                 cand.City,
		State:                cand.State,
		Zip:                  cand.Zip,
		OtherAffiliates:      cand.OtherAffiliates,
		TransactionsList:     cand.TransactionsList,
		TotalDirectInAmt:     cand.TotalDirectInAmt,
		TotalDirectInTxs:     cand.TotalDirectInTxs,
		AvgDirectIn:          cand.AvgDirectIn,
		TotalDirectOutAmt:    cand.TotalDirectOutAmt,
		TotalDirectOutTxs:    cand.TotalDirectOutTxs,
		AvgDirectOut:         cand.AvgDirectOut,
		NetBalanceDirectTx:   cand.NetBalanceDirectTx,
		DirectRecipientsAmts: cand.DirectRecipientsAmts,
		DirectRecipientsTxs:  cand.DirectRecipientsTxs,
		DirectSendersAmts:    cand.DirectSendersAmts,
		DirectSendersTxs:     cand.DirectSendersTxs,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeCand failed: ", err)
		return nil, fmt.Errorf("encodeCand failed: %v", err)
	}
	return data, nil
}

func encodeCandThreshold(entries []interface{}) []*protobuf.CandEntry {
	var es []*protobuf.CandEntry
	for _, intf := range entries {
		e := intf.(*donations.Entry)
		entry := &protobuf.CandEntry{
			ID:    e.ID,
			Total: e.Total,
		}
		es = append(es, entry)
	}
	return es
}

func decodeCand(data []byte) (donations.Candidate, error) {
	cand := &protobuf.Candidate{}
	err := proto.Unmarshal(data, cand)
	if err != nil {
		fmt.Println(err)
		return donations.Candidate{}, fmt.Errorf("decodeCand failed: %v", err)
	}

	entry := donations.Candidate{
		ID:                   cand.GetID(),
		Name:                 cand.GetName(),
		Party:                cand.GetParty(),
		OfficeState:          cand.GetOfficeState(),
		Office:               cand.GetOffice(),
		PCC:                  cand.GetPCC(),
		City:                 cand.GetCity(),
		State:                cand.GetState(),
		Zip:                  cand.GetZip(),
		OtherAffiliates:      cand.GetOtherAffiliates(),
		TransactionsList:     cand.GetTransactionsList(),
		TotalDirectInAmt:     cand.GetTotalDirectInAmt(),
		TotalDirectInTxs:     cand.GetTotalDirectInTxs(),
		AvgDirectIn:          cand.GetAvgDirectIn(),
		TotalDirectOutAmt:    cand.GetTotalDirectOutAmt(),
		TotalDirectOutTxs:    cand.GetTotalDirectOutTxs(),
		AvgDirectOut:         cand.GetAvgDirectOut(),
		NetBalanceDirectTx:   cand.GetNetBalanceDirectTx(),
		DirectRecipientsAmts: cand.GetDirectRecipientsAmts(),
		DirectRecipientsTxs:  cand.GetDirectRecipientsTxs(),
		DirectSendersAmts:    cand.GetDirectSendersAmts(),
		DirectSendersTxs:     cand.GetDirectSendersTxs(),
	}

	return entry, nil
}

func decodeCandThreshold(es []*protobuf.CandEntry) []interface{} {
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

func encodeCmpnFinancials(cf donations.CmpnFinancials) ([]byte, error) {
	entry := &protobuf.CmpnFinancials{
		CandID:         cf.CandID,
		Name:           cf.Name,
		PartyCd:        cf.PartyCd,
		Party:          cf.Party,
		TotalReceipts:  cf.TotalReceipts,
		TransFrAuth:    cf.TransFrAuth,
		TotalDisbsmts:  cf.TotalDisbsmts,
		TransToAuth:    cf.TransToAuth,
		COHBOP:         cf.COHBOP,
		COHCOP:         cf.COHCOP,
		CandConts:      cf.CandConts,
		CandLoans:      cf.CandLoans,
		OtherLoans:     cf.OtherLoans,
		CandLoanRepay:  cf.CandLoanRepay,
		OtherLoanRepay: cf.OtherLoanRepay,
		DebtsOwedBy:    cf.DebtsOwedBy,
		TotalIndvConts: cf.TotalIndvConts,
		SpecElection:   cf.SpecElection,
		PrimElection:   cf.PrimElection,
		RunElection:    cf.RunElection,
		GenElection:    cf.GenElection,
		GenElectionPct: cf.GenElectionPct,
		OtherCmteConts: cf.OtherCmteConts,
		PtyConts:       cf.PtyConts,
		// CvgEndDate: cf.CvgEndDate,
		IndvRefunds: cf.IndvRefunds,
		CmteRefunds: cf.CmteRefunds,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeCmpnFinancials failed: %v", err)
	}
	return data, nil
}

func decodeCmpnFinancials(data []byte) (donations.CmpnFinancials, error) {
	cf := &protobuf.CmpnFinancials{}
	err := proto.Unmarshal(data, cf)
	if err != nil {
		fmt.Println(err)
		return donations.CmpnFinancials{}, fmt.Errorf("decodeCmpnFinancials failed: %v", err)
	}

	cmpn := donations.CmpnFinancials{
		CandID:         cf.GetCandID(),
		Name:           cf.GetName(),
		PartyCd:        cf.GetPartyCd(),
		Party:          cf.GetParty(),
		TransFrAuth:    cf.GetTransFrAuth(),
		TotalDisbsmts:  cf.GetTotalDisbsmts(),
		TransToAuth:    cf.GetTransToAuth(),
		COHBOP:         cf.GetCOHBOP(),
		COHCOP:         cf.GetCOHCOP(),
		CandConts:      cf.GetCandConts(),
		CandLoans:      cf.GetCandLoans(),
		OtherLoans:     cf.GetOtherLoans(),
		CandLoanRepay:  cf.GetCandLoanRepay(),
		OtherLoanRepay: cf.GetOtherLoanRepay(),
		DebtsOwedBy:    cf.GetDebtsOwedBy(),
		TotalIndvConts: cf.GetTotalIndvConts(),
		SpecElection:   cf.GetSpecElection(),
		PrimElection:   cf.GetPrimElection(),
		RunElection:    cf.GetRunElection(),
		GenElection:    cf.GetGenElection(),
		GenElectionPct: cf.GetGenElectionPct(),
		OtherCmteConts: cf.GetOtherCmteConts(),
		PtyConts:       cf.GetPtyConts(),
		// CvgEndDate: cf.GetCvgEndDate(),
		IndvRefunds: cf.GetIndvRefunds(),
		CmteRefunds: cf.GetCmteRefunds(),
	}

	return cmpn, nil
}
