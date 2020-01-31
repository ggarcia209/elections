package persist

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/protobuf"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

// InitialCacheCCL stores CmteLink objects
func InitialCacheCCL(year string, objs []*donations.CmteLink, start bool) (bool, error) {
	if start {
		err := createBucket(year, "cmte_links")
		if err != nil {
			fmt.Println("InitialCacheCCL failed: ", err)
			return true, fmt.Errorf("InitialCacheCCL failed: %v", err)
		}
	}

	for _, obj := range objs {
		err := putCCL(year, obj)
		if err != nil {
			fmt.Println("InitialCacheCCL failed: putCCL failed: ", err)
			return true, fmt.Errorf("InitialCacheCCL failed: putCCL failed: %v", err)
		}
	}
	return false, nil
}

func putCCL(year string, ccl *donations.CmteLink) error {
	// convert obj to protobuf
	data, err := encodeCCL(*ccl)
	if err != nil {
		fmt.Println("encodeCCL failed: ", err)
		return fmt.Errorf("encodeCCL failed: %v", err)
	}
	// open/create bucket in db/offline_db.db
	// put protobuf item and use cand.ID as key
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: putCCL failed: 'offline_db.db' failed to open")
		return fmt.Errorf("putCCL failed: 'offline_db.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte("cmte_links"))
		if err := b.Put([]byte(ccl.CmteID), data); err != nil { // serialize k,v
			fmt.Printf("putCCL failed: offline_db.db': failed to store candidate-committee link: %s\n", ccl.LinkID)
			return fmt.Errorf("putCCL failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: putCCL failed: 'offline_db.db': 'cmte_links' bucket failed to open")
		return fmt.Errorf("putCCL failed: 'offline_db.db': 'cmte_links' bucket failed to open: %v", err)
	}

	return nil
}

// GetCCL returns a pointer to an Candidate obj stored on disk
func GetCCL(year, id string) (*donations.CmteLink, error) {
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: GetCCL failed: 'offline_db.db' failed to open")
		return nil, fmt.Errorf("GetCCL failed: 'offline_db.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(year)).Bucket([]byte("cmte_links")).Get([]byte(id))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetCCL failed: 'offline_db.db': 'candidates' bucket failed to open")
		return nil, fmt.Errorf("GetCCL failed: 'offline_db.db': 'candidates' bucket failed to open: %v", err)
	}

	ccl, err := decodeCCL(data)
	if err != nil {
		fmt.Println("GetCCL failed: decodeCCL failed: ", err)
		return nil, fmt.Errorf("GetCCL failed: decodeCCL failed: %v", err)
	}

	return &ccl, nil
}

// encodeCCL encodes CmteLink structs as protocol buffers
func encodeCCL(ccl donations.CmteLink) ([]byte, error) {
	entry := &protobuf.CmteLink{
		CandID:   ccl.CandID,
		CmteID:   ccl.CmteID,
		CmteType: ccl.CmteType,
		CmteDsgn: ccl.CmteDsgn,
		LinkID:   ccl.LinkID,
	}

	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeCCL failed: ", err)
		return nil, fmt.Errorf("encodeCCL failed: %v", err)
	}
	return data, nil
}

func decodeCCL(data []byte) (donations.CmteLink, error) {
	ccl := &protobuf.CmteLink{}
	err := proto.Unmarshal(data, ccl)
	if err != nil {
		fmt.Println("decodeCCL failed: ", err)
		return donations.CmteLink{}, fmt.Errorf("decodeCCL failed: %v", err)
	}

	entry := donations.CmteLink{CandID: ccl.GetCandID(), CmteID: ccl.GetCmteID(), CmteType: ccl.GetCmteType(),
		CmteDsgn: ccl.GetCmteDsgn(), LinkID: ccl.GetLinkID()}

	return entry, nil
}
