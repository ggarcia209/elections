// Package indexing contains operations for indexing object IDs.
// Indexed object IDs are not segregated by election year.
// Individual donors are indexed by Job (Employer + Occupation): Name: ID.
// All other objects are indexed by Zip: Name: ID.
package indexing

import (
	"fmt"

	"github.com/elections/protobuf"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

// IndexEntry contains an objects ID and a list of
// election years containing data for the object.
type IndexEntry struct {
	ID    string
	Years []string
}

// UpdateYears method adds a new year to the IndexEntry's Years field.
func (s *IndexEntry) UpdateYears(year string) {
	s.Years = append(s.Years, year)
}

// NewEntry creates a new IndexEntry object with ID ONLY and returns the pointer.
func NewEntry(id string) *IndexEntry {
	return &IndexEntry{ID: id, Years: []string{}}
}

// StoreIDByZip stores an IndexEntry object by Committee/Candidate name.
func StoreIDByZip(zip, name string, entry *IndexEntry) error {
	// encode IndexEntry
	data, err := encodeIndexEntry(*entry)
	if err != nil {
		fmt.Println("StoreIDByZip failed: ", err)
		return fmt.Errorf("StoreIDByZip failed: %v", err)
	}
	// open id_index.db and put object
	db, err := bolt.Open("db/id_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("StoreIDByZip failed: ", err)
		return fmt.Errorf("StoreIDByZip failed: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		// main bucket
		b, err := tx.CreateBucketIfNotExists([]byte("zip_lookup"))
		if err != nil {
			fmt.Println("StoreIDByZip failed: ", err)
			return fmt.Errorf("StoreIDByZip failed: %v", err)
		}
		// nested bucket corresponding to zip code
		nb, err := b.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			fmt.Println("StoreIDByZip failed: ", err)
			return fmt.Errorf("StoreIDByZip failed: %v", err)
		}
		if err := nb.Put([]byte(zip), []byte(data)); err != nil { // serialize k,v
			fmt.Printf("StoreIDByZip failed: failed to store object: %s\n", name)
			return fmt.Errorf("StoreIDByZip failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("StoreIDByZip failed: ", err)
		return fmt.Errorf("StoreIDByZip failed: %v", err)
	}
	return nil
}

// LookupIDByZip finds a Committee/Candidates ID by name.
func LookupIDByZip(zip, name string) (*IndexEntry, error) {
	db, err := bolt.Open("db/id_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("LookupIDByZip failed: ", err)
		return nil, fmt.Errorf("LookupIDByZip failed: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("name_lookup")).Bucket([]byte(name)).Get([]byte(zip))
		return nil
	}); err != nil {
		fmt.Println("LookupIDByZip failed: ", err)
		return nil, fmt.Errorf("LookupIDByZip failed: %v", err)
	}

	entry, err := decodeIndexEntry(data)
	if err != nil {
		fmt.Println("LookupIDByZip failed: ", err)
		return nil, fmt.Errorf("LookupIDByZip failed: %v", err)
	}

	return &entry, nil
}

// StoreIDByJob stores an IndexEntry object by Committee/Candidate name.
// Individual donors with same name are distinguished by job.
func StoreIDByJob(name, employer, occupation string, entry *IndexEntry) error {
	job := employer + " - " + occupation
	// encode IndexEntry
	data, err := encodeIndexEntry(*entry)
	if err != nil {
		fmt.Println("StoreIDByJob failed: ", err)
		return fmt.Errorf("StoreIDByJob failed: %v", err)
	}
	// open id_index.db and put object
	db, err := bolt.Open("db/id_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("StoreIDByJob failed: ", err)
		return fmt.Errorf("StoreIDByJob failed: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		// main bucket
		b, err := tx.CreateBucketIfNotExists([]byte("job_lookup"))
		if err != nil {
			fmt.Println("StoreIDByJob failed: ", err)
			return fmt.Errorf("StoreIDByJob failed: %v", err)
		}
		// nested bucket corresponding to name
		nb, err := b.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			fmt.Println("StoreIDByJob failed: ", err)
			return fmt.Errorf("StoreIDByJob failed: %v", err)
		}
		if err := nb.Put([]byte(job), []byte(data)); err != nil { // serialize k,v
			fmt.Printf("StoreIDByJob failed: failed to store object: %s\n", name)
			return fmt.Errorf("StoreIDByJob failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("StoreIDByJob failed: ", err)
		return fmt.Errorf("StoreIDByJob failed: %v", err)
	}
	return nil
}

// LookupIDByJob finds a Committee/Candidates ID by name.
func LookupIDByJob(name, employer, occupation string) (*IndexEntry, error) {
	job := employer + " - " + occupation
	db, err := bolt.Open("db/id_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("LookupIDByJob failed: ", err)
		return nil, fmt.Errorf("LookupIDByJob failed: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("name_lookup")).Bucket([]byte(name)).Get([]byte(job))
		return nil
	}); err != nil {
		fmt.Println("LookupIDByJob failed: ", err)
		return nil, fmt.Errorf("LookupIDByJob failed: %v", err)
	}

	entry, err := decodeIndexEntry(data)
	if err != nil {
		fmt.Println("LookupIDByJob failed: ", err)
		return nil, fmt.Errorf("LookupIDByJob failed: %v", err)
	}

	return &entry, nil
}

// convIndvToProto encodes LogData structs as protocol buffers
func encodeIndexEntry(e IndexEntry) ([]byte, error) { // move conversions to protobuf package?
	entry := &protobuf.IndexEntry{
		ID:    e.ID,
		Years: e.Years,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeIndexEntry failed: ", err)
		return nil, fmt.Errorf("encodeIndexEntry failed: %v", err)
	}
	return data, nil
}

func decodeIndexEntry(data []byte) (IndexEntry, error) {
	e := &protobuf.IndexEntry{}
	err := proto.Unmarshal(data, e)
	if err != nil {
		fmt.Println("decodeIndexEntry failed: ", err)
		return IndexEntry{}, fmt.Errorf("decodeIndexEntry failed: %v", err)
	}

	entry := IndexEntry{
		ID:    e.GetID(),
		Years: e.GetYears(),
	}

	return entry, nil
}
