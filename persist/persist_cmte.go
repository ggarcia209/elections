package persist

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/protobuf"

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

// DEPRECATED

/* CORRECT LOGIC TO UPDATE FOR EACh CMTE CONTRIBUTION - REMOVE IF UNNECESSARY */
/*

// InitialCacheCmte stores Committee objects with empty dynamic fields before
// they are updated by each contribution record
func InitialCacheCmte(year string, objs []*donations.Committee, start bool) error {
	if start {
		err := createBucket(year, "committees")
		if err != nil {
			fmt.Println("InitialCacheCmte failed: ", err)
			return fmt.Errorf("InitialCacheCmte failed: %v", err)
		}
	}

	for _, obj := range objs {
		err := PutCommittee(year, obj)
		if err != nil {
			fmt.Println("InitialCacheCmte failed: putCandidate failed: ", err)
			return fmt.Errorf("InitialCacheCmte failed: putCandidate failed: %v", err)
		}
	}
	return nil
}

// PutCommittee puts a Committee object belonging to the specified year to the database
func PutCommittee(year string, cmte *donations.Committee) error {
	// convert obj to protobuf
	data, err := encodeCmte(*cmte)
	if err != nil {
		fmt.Println("encodeCand failed: ", err)
		return fmt.Errorf("encodeCand failed: %v", err)
	}
	// open/create bucket in db/offline_db.db
	// put protobuf item and use cand.ID as key
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("Database failed to open: ", err)
		return fmt.Errorf("Database failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte("committees"))
		if err := b.Put([]byte(cmte.ID), data); err != nil { // serialize k,v
			fmt.Printf("Put failed to store object: %s\n", cmte.ID)
			return fmt.Errorf("Put failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("Update transaction failed: ", err)
		return fmt.Errorf("Update transaction failed: %v", err)
	}

	return nil
}

// GetCommittee returns a pointer to an Committee obj stored on disk
func GetCommittee(year, id string) (*donations.Committee, error) {
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: GetCommittee failed: 'offline_db.db' failed to open")
		return nil, fmt.Errorf("GetCommittee failed: 'offline_db.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(year)).Bucket([]byte("committees")).Get([]byte(id))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetCommittee failed: 'offline_db.db': 'committees' bucket failed to open")
		return nil, fmt.Errorf("GetCommittee failed: 'offline_db.db': 'committees' bucket failed to open: %v", err)
	}

	cmte, err := decodeCmte(data)
	if err != nil {
		fmt.Println("GetCommittee failed: decodeCand failed: ", err)
		return nil, fmt.Errorf("GetCommittee failed: decodeCand failed: %v", err)
	}

	return &cmte, nil
}


// CacheAndPersistCommittee persists a list of of Candidate objects to the on-disk cache
func CacheAndPersistCommittee(objs []*donations.Committee) error {
	err := createBucket("committees")
	if err != nil {
		fmt.Println("CacheAndPersistCommittee failed: ", err)
		return fmt.Errorf("CacheAndPersistCommittee failed: %v", err)
	}

	// for each obj
	for _, obj := range objs {
		err := PutCommittee(obj)
		if err != nil {
			fmt.Println("CacheAndPersistCommittee failed: putCommittee failed: ", err)
			return fmt.Errorf("CacheAndPersistCommittee failed: putCommittee failed: %v", err)
		}
	}
	return nil
} */

/* func updateCmteTxIn(cont *donations.IndvContribution) error {
	// get old value
	cmte, err := GetCommittee(cont.CmteID)
	if err != nil {
		fmt.Println("updateCmteIndvCont failed: GetCommittee failed", err)
		return fmt.Errorf("updateCmteIndvCont failed: GetCommittee failed: %v", err)
	}

	// update old values
	cmte.Donors = append(cmte.Donors)
	cmte.DonationsRcvd = append(cmte.DonationsRcvd, cont.TxID)
	cmte.TotalDonationsRcvd++
	cmte.TotalReceived += cont.TxAmt
	cmte.AvgDonationRcvd = cmte.TotalReceived / cmte.TotalDonationsRcvd

	// replace/overwrite old obj
	err = putCommittee(cmte)
	if err != nil {
		fmt.Println("updateCmteIndvCont failed: putCommittee failed: ", err)
		return fmt.Errorf("updateCmteIndvCont failed: putCommittee failed: %v", err)
	}
	return nil
}

func updateCmteTxOut(cont *donations.CmteContribution) {
	// get old value
	cmte, err := GetCommittee(cont.CmteID)
	if err != nil {
		fmt.Println("updateCmteIndvCont failed: GetCommittee failed", err)
		return fmt.Errorf("updateCmteIndvCont failed: GetCommittee failed: %v", err)
	}

	// update old values

} */
