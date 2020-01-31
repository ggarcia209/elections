// Package persist contains functions to persist data stored in memory to the local disk
package persist

// DEPRECATED

/*
import (
	"fmt"
	"projects/elections/donations"
	"projects/elections/protobuf"

	"github.com/golang/protobuf/proto"
)

// encodeDisbRecipient encodes LogData structs as protocol buffers
func encodeDisbRecipient(indv donations.DisbRecipient) ([]byte, error) { // move conversions to protobuf package?
	entry := &protobuf.DisbRecipient{
		ID:                 indv.ID,
		Name:               indv.Name,
		City:               indv.City,
		State:              indv.State,
		Zip:                indv.Zip,
		Disbursements:      indv.Disbursements,
		TotalDisbursements: indv.TotalDisbursements,
		TotalReceived:      indv.TotalReceived,
		AvgReceived:        indv.AvgReceived,
		SendersAmt:         indv.SendersAmt,
		SendersTxs:         indv.SendersTxs,
	}

	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeDisbRecipient failed: ", err)
		return nil, fmt.Errorf("encodeDisbRecipient failed: %v", err)
	}
	return data, nil
}

func decodeDisbRecipient(data []byte) (donations.DisbRecipient, error) {
	indv := &protobuf.DisbRecipient{}
	err := proto.Unmarshal(data, indv)
	if err != nil {
		fmt.Println("decodeDisbRecipient failed: ", err)
		return donations.DisbRecipient{}, fmt.Errorf("decodeDisbRecipient failed: %v", err)
	}

	entry := donations.DisbRecipient{
		ID:                 indv.GetID(),
		Name:               indv.GetName(),
		City:               indv.GetCity(),
		State:              indv.GetState(),
		Zip:                indv.GetZip(),
		Disbursements:      indv.GetDisbursements(),
		TotalDisbursements: indv.GetTotalDisbursements(),
		TotalReceived:      indv.GetTotalReceived(),
		AvgReceived:        indv.GetAvgReceived(),
		SendersAmt:         indv.GetSendersAmt(),
		SendersTxs:         indv.GetSendersTxs(),
	}

	return entry, nil
}

// CacheAndPersistDisbRecipient persists a list of of DisbRecipient objects to the on-disk cache
func CacheAndPersistDisbRecipient(year string, objs []*donations.DisbRecipient, start bool) error {
	if start {
		err := createBucket(year, "disb_recipients")
		if err != nil {
			fmt.Println("CacheAndPersistDisbRecipient failed: ", err)
			return fmt.Errorf("CacheAndPersistDisbRecipient failed: %v", err)
		}
	}

	// for each obj
	for _, obj := range objs {
		err := PutDisbRecipient(year, obj)
		if err != nil {
			fmt.Println("CacheAndPersistDisbRecipient failed: ", err)
			return fmt.Errorf("CacheAndPersistDisbRecipient failed: %v", err)
		}
	}
	return nil
}

// PutDisbRecipient puts a new DisbRecipient object in the database bucket for the given year
func PutDisbRecipient(year string, rec *donations.DisbRecipient) error {
	// get name & zip
	name, zip := rec.Name, rec.Zip

	// convert obj to protobuf
	data, err := encodeDisbRecipient(*rec)
	if err != nil {
		fmt.Println("PutDisbRecipient failed: ", err)
		return fmt.Errorf("PutDisbRecipient failed: %v", err)
	}
	// open/create bucket in db/offline_db.db
	// put protobuf item and use rec.ID as key
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: PutDisbRecipient failed: 'offline_db.db' failed to open")
		return fmt.Errorf("PutDisbRecipient failed: 'offline_db.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte("disb_recipients"))
		if err := b.Put([]byte(rec.ID), data); err != nil { // serialize k,v
			fmt.Printf("PutDisbRecipient failed: offline_db.db': failed to store recipient: %s\n", rec.ID)
			return fmt.Errorf("PutDisbRecipient failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: PutDisbRecipient failed: 'offline_db.db': 'disb_recipients' bucket failed to open")
		return fmt.Errorf("PutDisbRecipient failed: 'offline_db.db': 'disb_recipients' bucket failed to open: %v", err)
	}

	err = RecordIDByZip(name, zip, rec.ID)
	if err != nil {
		fmt.Println("putDisbRecipient failed: RecordIDByJob failed")
		return fmt.Errorf("putDisbRecipient failed: %v", err)
	}

	return nil
}

// GetDisbRecipient returns a pointer to an DisbRecipient obj stored on disk
func GetDisbRecipient(year, id string) (*donations.DisbRecipient, error) {
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: GetDisbRecipient failed: 'offline_db.db' failed to open")
		return nil, fmt.Errorf("GetDisbRecipient failed: 'offline_db.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(year)).Bucket([]byte("disb_recipients")).Get([]byte(id))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetDisbRecipient failed: 'offline_db.db': 'disb_recipients' bucket failed to open")
		return nil, fmt.Errorf("GetDisbRecipient failed: 'offline_db.db': 'disb_recipients' bucket failed to open: %v", err)
	}

	rec, err := decodeDisbRecipient(data)
	if err != nil {
		fmt.Println("GetDisbRecipient failed: convProtoToIndv failed: ", err)
		return nil, fmt.Errorf("GetDisbRecipient failed: convProtoToIndv failed: %v", err)
	}

	return &rec, nil
}

func updateDisbRecipient(new *donations.DisbRecipient) error {
	// get old value
	old, err := GetDisbRecipient(new.ID)
	if err != nil {
		fmt.Println("updateDisbRecipient failed: ", err)
		return fmt.Errorf("updateDisbRecipient failed: %v", err)
	}

	// Add old values to new struct
	for _, disbursement := range old.Disbursements {
		new.Disbursements = append(new.Disbursements, disbursement)
	}
	new.TotalDisbursements += old.TotalDisbursements
	new.TotalReceived += old.TotalReceived
	// replace/overwrite old value
	err = PutDisbRecipient(new)
	if err != nil {
		fmt.Println("updateDisbRecipient failed: putIndvDonor failed: ", err)
		return fmt.Errorf("updateDisbRecipient failed: putIndvDonor failed: %v", err)
	}
	return nil
} */
