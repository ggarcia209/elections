package persist

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/protobuf"

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
		fmt.Println("decodeCand failed: ", err)
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

// DEPRECATED
/*

// InitialCacheCand stores Candidate objects with empty dynamic fields before
// they are updated by each contribution record
func InitialCacheCand(year string, objs []*donations.Candidate, start bool) error {
	if start {
		err := createBucket(year, "candidates")
		if err != nil {
			fmt.Println("InitialCache failed: ", err)
			return fmt.Errorf("InitialCache failed: %v", err)
		}
	}

	for _, obj := range objs {
		err := PutCandidate(year, obj)
		if err != nil {
			fmt.Println("InitialCache failed: putCandidate failed: ", err)
			return fmt.Errorf("InitialCache failed: putCandidate failed: %v", err)
		}
	}
	return nil
}

// PutCandidate saves a Candidate obj to the database
func PutCandidate(year string, cand *donations.Candidate) error {
	// convert obj to protobuf
	data, err := encodeCand(*cand)
	if err != nil {
		fmt.Println("encodeCand failed: ", err)
		return fmt.Errorf("encodeCand failed: %v", err)
	}
	// open/create bucket in db/offline_db.db
	// put protobuf item and use cand.ID as key
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: putCandidate failed: 'offline_db.db' failed to open")
		return fmt.Errorf("putCandidate failed: 'offline_db.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte("candidates"))
		if err := b.Put([]byte(cand.ID), data); err != nil { // serialize k,v
			fmt.Printf("putCandidate failed: offline_db.db': failed to store candidate: %s\n", cand.ID)
			return fmt.Errorf("putCandidate failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: putCandidate failed: 'offline_db.db': 'candidates' bucket failed to open")
		return fmt.Errorf("putCandidate failed: 'offline_db.db': 'candidates' bucket failed to open: %v", err)
	}

	return nil
}

// GetCandidate returns a pointer to an Candidate obj stored on disk
func GetCandidate(year, id string) (*donations.Candidate, error) {
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: GetCandidate failed: 'offline_db.db' failed to open")
		return nil, fmt.Errorf("GetCandidate failed: 'offline_db.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(year)).Bucket([]byte("candidates")).Get([]byte(id))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetCandidate failed: 'offline_db.db': 'candidates' bucket failed to open")
		return nil, fmt.Errorf("GetCandidate failed: 'offline_db.db': 'candidates' bucket failed to open: %v", err)
	}

	cand, err := decodeCand(data)
	if err != nil {
		fmt.Println("GetCandidate failed: decodeCand failed: ", err)
		return nil, fmt.Errorf("GetCandidate failed: decodeCand failed: %v", err)
	}

	return &cand, nil
}


// CacheAndPersistCandidates persists a list of of Candidate objects to the on-disk cache
func CacheAndPersistCandidates(objs []*donations.Candidate, seen map[string]bool) error {
	if len(seen) == 0 {
		err := createBucket("candidates")
		if err != nil {
			fmt.Println("CacheAndPersistCandidates failed: ", err)
			return fmt.Errorf("CacheAndPersistCandidates failed: %v", err)
		}
	}

	// for each obj
	for _, obj := range objs {
		err := PutCandidate(obj)
		if err != nil {
			fmt.Println("CacheAndPersistCandidates failed: putCandidate failed: ", err)
			return fmt.Errorf("CacheAndPersistCandidates failed: putCandidate failed: %v", err)
		}
	}
	return nil
} */

// Move update logic to databuilder package
/* CORRECT LOGIC TO UPDATE FOR EACh CMTE CONTRIBUTION */
/* func UpdateCandidate(new *donations.Candidate) error {
	// get old value
	old, err := GetCandidate(new.ID)
	if err != nil {
		fmt.Println("updateCandidate failed: GetCandidate failed", err)
		return fmt.Errorf("updateCandidate failed: GetCandidate failed: %v", err)
	}

	// Add old values to new struct
	new.TotalDonations += old.TotalDonations
	new.TotalRaised += old.TotalRaised
	new.AvgDonation = new.TotalRaised / new.TotalDonations
	// replace/overwrite old value
	err = putCandidate(new)
	if err != nil {
		fmt.Println("updateCandidate failed: putIndvDonor failed: ", err)
		return fmt.Errorf("updateCandidate failed: putIndvDonor failed: %v", err)
	}
	return nil
} */
