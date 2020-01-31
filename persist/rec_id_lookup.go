package persist

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// RecordIDByZip saves ID Lookup entries to disk
// get old values and update map before recording
// create bucket at top level if ID < 1
func RecordIDByZip(name, zip, recID string) error { // map is assumed to be current values (new + old)
	m, err := LookupRecIDByName(name)
	if m == nil {
		m = make(map[string]string)
	}
	m[zip] = recID

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("RecordIDByZip failed: ", err)
		return fmt.Errorf("RecordIDByJob failed: %v", err)
	}

	if err != nil {
		fmt.Println("RecordIDByZip failed: ", err)
		return fmt.Errorf("RecordIDByJob failed: %v", err)
	}

	mapBytes, err := encodeLookupMap(m)
	if err != nil {
		fmt.Println("RecordIDByZip failed: ", err)
		return fmt.Errorf("RecordIDByJob failed: %v", err)
	}

	// tx
	if err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("rec_id_lookup"))
		if err != nil {
			fmt.Println("RecordIDByZip failed: ", err)
			return fmt.Errorf("RecordIDByJob failed: %v", err)
		}
		if err := b.Put([]byte(name), mapBytes); err != nil { // serialize k,v
			fmt.Println("RecordIDByZip failed: ", err)
			return fmt.Errorf("RecordIDByJob failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("RecordIDByZip failed: ", err)
		return fmt.Errorf("RecordIDByJob failed: %v", err)
	}
	return nil
}

// LookupRecIDByName finds the ID for the given name and zip
func LookupRecIDByName(name string) (map[string]string, error) {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("LookupRecIDByName failed: ", err)
		return nil, fmt.Errorf("LookupRecIDByName failed: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("rec_id_lookup")).Get([]byte(name))
		return nil
	}); err != nil {
		fmt.Println("LookupRecIDByName failed: ", err)
		return nil, fmt.Errorf("LookupRecIDByName failed: %v", err)
	}

	if data == nil {
		return nil, nil
	}

	m, err := decodeLookupMap(data)
	if err != nil {
		fmt.Println("LookupRecIDByName failed: ", err)
		return nil, fmt.Errorf("LookupRecIDByName failed: %v", err)
	}

	return m, nil
}

/* DEPRECATED */

/*

// CreateRecIDLookupBucket initializes the 'id_lookup' bucket in disk_cache.db
func CreateRecIDLookupBucket() error {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: CreateRecIDLookupBucket failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("CreateRecIDLookupBucket failed: 'disk_cache.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("rec_id_lookup"))
		if err != nil {
			fmt.Println("FATAL: CreateRecIDLookupBucket failed: 'disk_cache.db': 'id_lookup' bucket failed to open")
			return fmt.Errorf("'main': FATAL: 'disk_cache.db': 'id_lookup' bucket failed to open: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: CreateRecIDLookupBucket failed: 'disk_cache.db': 'id_lookup' bucket failed to open")
		return fmt.Errorf("CreateRecIDLookupBucket failed: 'disk_cache.db': 'id_lookup' bucket failed to open: %v", err)
	}

	return nil
}

// RecordIDByZip saves ID Lookup entries to disk
// get old values and update map before recording
// create bucket at top level if ID < 1
func RecordIDByZip(name, zip string, recID int32) error { // map is assumed to be current values (new + old)
	m, err := LookupRecIDByName(name)
	if m == nil {
		m = make(map[string]int32)
	}
	m[zip] = recID

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: RecordIDByJob failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("RecordIDByJob failed: 'disk_cache.db' failed to open: %v", err)
	}

	if err != nil {
		fmt.Printf("RecordIDByJob failed: %v\n", err)
		return fmt.Errorf("RecordIDByJob failed: %v", err)
	}

	mapBytes, err := encodeLookupMap(m)
	if err != nil {
		fmt.Println("RecordIDByJob failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("RecordIDByJob failed: 'disk_cache.db' failed to open: %v", err)
	}

	// tx
	if err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("rec_id_lookup"))
		if err != nil {
			fmt.Println("FATAL: RecordIDByJob failed: 'disk_cache.db': 'rec_id_lookup' bucket failed to open")
			return fmt.Errorf("'main': FATAL: 'disk_cache.db': 'rec_id_lookup' bucket failed to open: %v", err)
		}
		if err := b.Put([]byte(name), mapBytes); err != nil { // serialize k,v
			fmt.Printf("RecordIDByJob failed: disk_cache.db': 'rec_id_lookup': '%s' failed to store\n", name)
			return fmt.Errorf("RecordIDByJob failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: RecordIDByJob failed: 'disk_cache.db': 'rec_id_lookup' bucket failed to open")
		return fmt.Errorf("RecordIDByJob failed: 'disk_cache.db': 'rec_id_lookup' bucket failed to open: %v", err)
	}
	return nil
}



// LookupRecIDByName finds the ID for the given name and zip
func LookupRecIDByName(name string) (map[string]int32, error) {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: LookupRecIDByName failed: 'disk_cache.db' failed to open")
		return nil, fmt.Errorf("LookupRecIDByName failed: 'disk_cache.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("rec_id_lookup")).Get([]byte(name))
		return nil
	}); err != nil {
		fmt.Println("FATAL: LookupRecIDByName failed: 'disk_cache.db': 'rec_id_lookup' bucket failed to open")
		return nil, fmt.Errorf("LookupRecIDByName failed: 'disk_cache.db': 'rec_id_lookup' bucket failed to open: %v", err)
	}

	if data == nil {
		return nil, nil
	}

	m, err := decodeLookupMap(data)
	if err != nil {
		fmt.Println("LookupRecIDByName failed: decodeLookupMap failed ", err)
		return nil, fmt.Errorf("LookupRecIDByName failed: decodeLookupMap failed %v", err)
	}

	return m, nil
}


func RecordIDByZipOld(name string, m map[string]int32) error { // map is assumed to be current values (new + old)
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: RecordIDByZip failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("RecordIDByZip failed: 'disk_cache.db' failed to open: %v", err)
	}

	mapBytes, err := encodeLookupMap(m)
	if err != nil {
		fmt.Println("RecordIDByZip failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("RecordIDByZip failed: 'disk_cache.db' failed to open: %v", err)
	}

	// tx
	if err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("rec_id_lookup"))
		if err != nil {
			fmt.Println("FATAL: RecordIDByZip failed: 'disk_cache.db': 'rec_id_lookup' bucket failed to open")
			return fmt.Errorf("'main': FATAL: 'disk_cache.db': 'rec_id_lookup' bucket failed to open: %v", err)
		}
		if err := b.Put([]byte(name), mapBytes); err != nil { // serialize k,v
			fmt.Printf("RecordIDByZip failed: disk_cache.db': 'rec_id_lookup': '%s' failed to store\n", name)
			return fmt.Errorf("RecordIDByZip failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: RecordIDByZip failed: 'disk_cache.db': 'rec_id_lookup' bucket failed to open")
		return fmt.Errorf("RecordIDByZip failed: 'disk_cache.db': 'rec_id_lookup' bucket failed to open: %v", err)
	}
	return nil
}

// StoreRecID stores RecID for future calls to databuild.FindRecipient()
func StoreRecID(id int) error {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: StoreRecID failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("StoreRecID failed: 'disk_cache.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("IDs"))
		if err := b.Put([]byte("recID"), util.Itob(id)); err != nil { // serialize k,v
			fmt.Printf("StoreRecID failed: disk_cache.db': 'recID': '%d' failed to store\n", id)
			return fmt.Errorf("StoreRecID failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: StoreRecID failed: 'disk_cache.db': 'IDs' bucket failed to open")
		return fmt.Errorf("StoreRecID failed: 'disk_cache.db': 'IDs' bucket failed to open: %v", err)
	}
	return nil
}

// GetRecID returns the last-saved RecID value
func GetRecID() (int, error) {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: GetRecID failed: 'disk_cache.db' failed to open")
		return 0, fmt.Errorf("GetRecID failed: 'disk_cache.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("IDs")).Get([]byte("recID"))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetRecID failed: 'disk_cache.db': 'IDs' bucket failed to open")
		return 0, fmt.Errorf("GetRecID failed: 'disk_cache.db': 'IDs' bucket failed to open: %v", err)
	}

	if data == nil {
		return 1, nil
	}

	id := util.Btoi(data)
	return id, nil
}

*/
