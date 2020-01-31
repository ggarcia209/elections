package persist

import (
	"fmt"

	"github.com/elections/protobuf"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

// IDLookup is used to marshall/unmarshall protobufs
// used for both DonorID and RecID lookups; "job" string interchangeable with "zip"
type IDLookup struct {
	LookupByJob map[string]string
}

// createLookupBucket initializes the 'id_lookup' bucket in disk_cache.db
func createLookupBuckets() error {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("createLookupBucket failed: ", err)
		return fmt.Errorf("createLookupBucket failed: %v", err)
	}

	// indv ID lookup
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("id_lookup"))
		if err != nil {
			fmt.Println("createLookupBucket failed: ", err)
			return fmt.Errorf("createLookupBucket failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("createLookupBucket failed: ", err)
		return fmt.Errorf("createLookupBucket failed: %v", err)
	}

	// disb_rec ID lookup
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("rec_id_lookup"))
		if err != nil {
			fmt.Println("createLookupBucket failed: ", err)
			return fmt.Errorf("createLookupBucket failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("createLookupBucket failed: ", err)
		return fmt.Errorf("createLookupBucket failed: %v", err)
	}

	return nil
}

// RecordIDByJob saves ID Lookup entries to disk
// get old values and update map before recording
func RecordIDByJob(name, job string, donorID string) error { // map is assumed to be current values (new + old)
	m, err := LookupIDByName(name)
	if m == nil {
		m = make(map[string]string)
	}
	m[job] = donorID

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: RecordIDByJob failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("RecordIDByJob failed: 'disk_cache.db' failed to open: %v", err)
	}

	if err != nil {
		fmt.Printf("RecordIDByJob failed: LookupIDByName failed: %v\n", err)
		return fmt.Errorf("RecordIDByJob failed: LookupIDByName failed: %v", err)
	}

	mapBytes, err := encodeLookupMap(m)
	if err != nil {
		fmt.Println("RecordIDByJob failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("RecordIDByJob failed: 'disk_cache.db' failed to open: %v", err)
	}

	// tx
	if err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("id_lookup"))
		if err != nil {
			fmt.Println("FATAL: RecordIDByJob failed: 'disk_cache.db': 'id_lookup' bucket failed to open")
			return fmt.Errorf("'main': FATAL: 'disk_cache.db': 'id_lookup' bucket failed to open: %v", err)
		}
		if err := b.Put([]byte(name), mapBytes); err != nil { // serialize k,v
			fmt.Printf("RecordIDByJob failed: disk_cache.db': 'id_lookup': '%s' failed to store\n", name)
			return fmt.Errorf("RecordIDByJob failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: RecordIDByJob failed: 'disk_cache.db': 'id_lookup' bucket failed to open")
		return fmt.Errorf("RecordIDByJob failed: 'disk_cache.db': 'id_lookup' bucket failed to open: %v", err)
	}
	return nil
}

// LookupIDByName finds the ID for the given name and job
func LookupIDByName(name string) (map[string]string, error) {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: LookupIDByName failed: 'disk_cache.db' failed to open")
		return nil, fmt.Errorf("LookupIDByName failed: 'disk_cache.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("id_lookup")).Get([]byte(name))
		return nil
	}); err != nil {
		fmt.Println("FATAL: LookupIDByName failed: 'disk_cache.db': 'id_lookup' bucket failed to open")
		return nil, fmt.Errorf("LookupIDByName failed: 'disk_cache.db': 'id_lookup' bucket failed to open: %v", err)
	}

	if data == nil {
		return nil, nil
	}

	m, err := decodeLookupMap(data)
	if err != nil {
		fmt.Println("LookupIDByName failed: decodeLookupMap failed ", err)
		return nil, fmt.Errorf("LookupIDByName failed: decodeLookupMap failed %v", err)
	}

	return m, nil
}

func encodeLookupMap(m map[string]string) ([]byte, error) {
	entry := &protobuf.Lookup{
		DonorID: m,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeLookupMap failed: ", err)
		return nil, fmt.Errorf("encodeLookupMap failed: %v", err)
	}
	return data, nil
}

func decodeLookupMap(data []byte) (map[string]string, error) {
	m := &protobuf.Lookup{}
	err := proto.Unmarshal(data, m)
	if err != nil {
		fmt.Println("decodeLookupMap failed: ", err)
		return nil, fmt.Errorf("decodeLookupMap failed: %v", err)
	}

	entry := IDLookup{m.GetDonorID()}
	return entry.LookupByJob, nil
}

/* DEPRECATED */
/*

// CreateIDBucket creates the "ID's" bucket in the "disk_cache.db" file
func CceateIDBucket() error {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: CreateIDBucket failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("CreateIDBucket failed: 'disk_cache.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("IDs"))
		if err != nil {
			fmt.Println("FATAL: CreateIDBucket failed: 'disk_cache.db': 'IDs' bucket failed to open")
			return fmt.Errorf("'main': FATAL: 'disk_cache.db': 'IDs' bucket failed to open: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: CreateIDBucket failed: 'disk_cache.db': 'IDs' bucket failed to open")
		return fmt.Errorf("CreateIDBucket failed: 'disk_cache.db': 'IDs' bucket failed to open: %v", err)
	}

	return nil
}

// IDLookup is used to marshall/unmarshall protobufs
// used for both DonorID and RecID lookups; "job" string interchangeable with "zip"
type IDLookup struct {
	LookupByJob map[string]int32
}

// RecordIDByJob saves ID Lookup entries to disk
// get old values and update map before recording
// create bucket at top level if ID < 1
func RecordIDByJob(name, job string, donorID int32) error { // map is assumed to be current values (new + old)
	m, err := LookupIDByName(name)
	if m == nil {
		m = make(map[string]int32)
	}
	m[job] = donorID

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: RecordIDByJob failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("RecordIDByJob failed: 'disk_cache.db' failed to open: %v", err)
	}

	if err != nil {
		fmt.Printf("RecordIDByJob failed: LookupIDByName failed: %v\n", err)
		return fmt.Errorf("RecordIDByJob failed: LookupIDByName failed: %v", err)
	}

	mapBytes, err := encodeLookupMap(m)
	if err != nil {
		fmt.Println("RecordIDByJob failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("RecordIDByJob failed: 'disk_cache.db' failed to open: %v", err)
	}

	// tx
	if err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("id_lookup"))
		if err != nil {
			fmt.Println("FATAL: RecordIDByJob failed: 'disk_cache.db': 'id_lookup' bucket failed to open")
			return fmt.Errorf("'main': FATAL: 'disk_cache.db': 'id_lookup' bucket failed to open: %v", err)
		}
		if err := b.Put([]byte(name), mapBytes); err != nil { // serialize k,v
			fmt.Printf("RecordIDByJob failed: disk_cache.db': 'id_lookup': '%s' failed to store\n", name)
			return fmt.Errorf("RecordIDByJob failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: RecordIDByJob failed: 'disk_cache.db': 'id_lookup' bucket failed to open")
		return fmt.Errorf("RecordIDByJob failed: 'disk_cache.db': 'id_lookup' bucket failed to open: %v", err)
	}
	return nil
}

// LookupIDByName finds the ID for the given name and job
func LookupIDByName(name string) (map[string]int32, error) {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: LookupIDByName failed: 'disk_cache.db' failed to open")
		return nil, fmt.Errorf("LookupIDByName failed: 'disk_cache.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("id_lookup")).Get([]byte(name))
		return nil
	}); err != nil {
		fmt.Println("FATAL: LookupIDByName failed: 'disk_cache.db': 'id_lookup' bucket failed to open")
		return nil, fmt.Errorf("LookupIDByName failed: 'disk_cache.db': 'id_lookup' bucket failed to open: %v", err)
	}

	if data == nil {
		return nil, nil
	}

	m, err := decodeLookupMap(data)
	if err != nil {
		fmt.Println("LookupIDByName failed: decodeLookupMap failed ", err)
		return nil, fmt.Errorf("LookupIDByName failed: decodeLookupMap failed %v", err)
	}

	return m, nil
}


// StoreDonorID stores DonorID for future calls to databuild.FindPerson()
func StoreDonorID(id int) error {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: storeDonorID failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("storeDonorID failed: 'disk_cache.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("IDs"))
		if err := b.Put([]byte("donorID"), util.Itob(id)); err != nil { // serialize k,v
			fmt.Printf("storeDonorID failed: disk_cache.db': 'donorID': '%d' failed to store\n", id)
			return fmt.Errorf("storeDonorID failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: storeDonorID failed: 'disk_cache.db': 'IDs' bucket failed to open")
		return fmt.Errorf("storeDonorID failed: 'disk_cache.db': 'IDs' bucket failed to open: %v", err)
	}
	return nil
}

// GetDonorID returns the last-saved DonorID value
func GetDonorID() (int, error) {
	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: GetDonorID failed: 'disk_cache.db' failed to open")
		return 0, fmt.Errorf("GetDonorID failed: 'disk_cache.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("IDs")).Get([]byte("donorID"))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetDonorID failed: 'disk_cache.db': 'IDs' bucket failed to open")
		return 0, fmt.Errorf("GetDonorID failed: 'disk_cache.db': 'IDs' bucket failed to open: %v", err)
	}

	if data == nil {
		return 1, nil
	}

	id := util.Btoi(data)
	return id, nil
}

func encodeLookupMap(m map[string]int32) ([]byte, error) {
	entry := &protobuf.Lookup{
		DonorID: m,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeLookupMap failed: ", err)
		return nil, fmt.Errorf("encodeLookupMap failed: %v", err)
	}
	return data, nil
}

func decodeLookupMap(data []byte) (map[string]int32, error) {
	m := &protobuf.Lookup{}
	err := proto.Unmarshal(data, m)
	if err != nil {
		fmt.Println("decodeLookupMap failed: ", err)
		return nil, fmt.Errorf("decodeLookupMap failed: %v", err)
	}

	entry := IDLookup{m.GetDonorID()}
	return entry.LookupByJob, nil
}

*/
