// Package persist contains functions to persist data stored in memory to the local disk
package persist

import (
	"fmt"
	"os"

	"github.com/elections/donations"

	"github.com/boltdb/bolt"
)

// Init inititalizes the program by creating the db directory, disk_cache.db, offline_db.db, and corresponding buckets
func Init(year string) error {
	createDB()

	err := createLookupBuckets()
	if err != nil {
		fmt.Println("Init failed: ", err)
		return fmt.Errorf("Init failed: %v", err)
	}

	err = createObjBuckets(year)
	if err != nil {
		fmt.Println("Init failed: ", err)
		return fmt.Errorf("Init failed: %v", err)
	}

	fmt.Println("Init Done")
	return nil
}

// StoreObjects persists a list of objects to the on-disk database
func StoreObjects(year string, objs []interface{}) error {
	for _, obj := range objs {
		err := PutObject(year, obj)
		if err != nil {
			fmt.Println("StoreObjects failed: ", err)
			return fmt.Errorf("StoreObjects failed: %v", err)
		}
	}
	return nil
}

// PutObject puts an object by year:bucket:key
func PutObject(year string, object interface{}) error {
	// encode object
	bucket, key, data, err := encodeToProto(object)
	if err != nil {
		fmt.Println("PutObject failed: ", err)
		return fmt.Errorf("PutObject failed: %v", err)
	}

	// open/create bucket in db/offline_db.db
	// put protobuf item and use donor.ID as key
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("PutObject failed: ", err)
		return fmt.Errorf("PutObject failed: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte(bucket))
		if err := b.Put([]byte(key), data); err != nil { // serialize k,v
			fmt.Printf("PutObject failed: failed to store object: %s\n", key)
			return fmt.Errorf("PutObject failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("PutObject failed: ", err)
		return fmt.Errorf("PutObject failed: %v", err)
	}
	return nil
}

// GetObject gets an object by year:bucket:key and returns it as an interface
func GetObject(year, bucket, key string) (interface{}, error) {
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("GetObject failed: ", err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(year)).Bucket([]byte(bucket)).Get([]byte(key))
		return nil
	}); err != nil {
		fmt.Println("PutObject failed: ", err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}

	obj, err := decodeFromProto(bucket, data) // change to decode
	if err != nil {
		fmt.Println("GetObject failed: ", err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}

	return &obj, nil
}

// encodeToProto encodes an object interface to protobuf
func encodeToProto(obj interface{}) (string, string, []byte, error) {
	switch obj.(type) {
	case nil:
		return "", "", nil, fmt.Errorf("encodeToProto failed: nil interface")
	case *donations.Individual:
		bucket := "individuals"
		key := obj.(*donations.Individual).ID
		data, err := encodeIndv(*obj.(*donations.Individual))
		if err != nil {
			fmt.Println("encodeToProto failed: ", err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.Organization:
		bucket := "organizations"
		key := obj.(*donations.Organization).ID
		data, err := encodeOrg(*obj.(*donations.Organization))
		if err != nil {
			fmt.Println("encodeToProto failed: ", err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.Committee:
		bucket := "committees"
		key := obj.(*donations.Committee).ID
		data, err := encodeCmte(*obj.(*donations.Committee))
		if err != nil {
			fmt.Println("encodeToProto failed: ", err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.Candidate:
		bucket := "candidates"
		key := obj.(*donations.Candidate).ID
		data, err := encodeCand(*obj.(*donations.Candidate))
		if err != nil {
			fmt.Println("encodeToProto failed: ", err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.CmteTxData:
		bucket := "cmte_tx_data"
		key := obj.(*donations.CmteTxData).CmteID
		data, err := encodeCmteTxData(*obj.(*donations.CmteTxData))
		if err != nil {
			fmt.Println("encodeToProto failed: ", err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.TopOverallData:
		bucket := "top_overall"
		key := obj.(*donations.TopOverallData).Category
		data, err := encodeOverallData(*obj.(*donations.TopOverallData))
		if err != nil {
			fmt.Println("encodeToProto failed: ", err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	default:
		return "", "", nil, fmt.Errorf("encodeToProto failed: invalid interface type")
	}
}

// decodeFromProto encodes an object interface to protobuf
func decodeFromProto(bucket string, data []byte) (interface{}, error) {
	switch bucket {
	case "":
		return nil, fmt.Errorf("decodeFromProto failed: nil bucket")
	case "indviduals":
		data, err := decodeIndv(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return data, nil
	case "organizations":
		data, err := decodeOrg(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return data, nil
	case "committees":
		data, err := decodeCmte(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return data, nil
	case "candidates":
		data, err := decodeCand(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return data, nil
	case "cmte_tx_data":
		data, err := decodeCmteTxData(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return data, nil
	case "top_overall":
		data, err := decodeOverallData(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return data, nil
	default:
		return nil, fmt.Errorf("decodeFromProto failed: invalid bucket")
	}
}

// CreateDB creates the database directory. CreateDB must be called
// before any other function in 'parse' package is called.
func createDB() {
	if _, err := os.Stat("db"); os.IsNotExist(err) {
		os.Mkdir("db", 0744)
		fmt.Println("CreateDB successful: 'db' directory created")
	}
}

func createObjBuckets(year string) error {
	buckets := []string{"individuals", "organizations", "committees", "candidates", "cmte_tx_data", "top_overall"}
	for _, bucket := range buckets {
		err := createBucket(year, bucket)
		if err != nil {
			fmt.Println("createObjBuckets failed: ", err)
			return fmt.Errorf("createObjBuckets failed: %v", err)
		}
	}
	return nil
}

func createBucket(year, name string) error {
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("createBucket failed: ", err)
		return fmt.Errorf("createBucket failed: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		yb, err := tx.CreateBucketIfNotExists([]byte(year))
		if err != nil {
			fmt.Println("createBucket failed: ", err)
			return fmt.Errorf("createBucket failed: %v", err)
		}
		_, err = yb.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			fmt.Println("createBucket failed: ", err)
			return fmt.Errorf("createBucet failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("createBucket failed: ", err)
		return fmt.Errorf("createBucket failed: %v", err)
	}

	return nil
}

// DEPRECATED

/*

// CacheAndPersistIndvDonor persists a list of of Individual donor objects to the on-disk cache
func CacheAndPersistIndvDonor(year string, objs []*donations.Individual, start bool) error {
	if start {
		err := createBucket(year, "indv_donors")
		if err != nil {
			fmt.Println("CacheAndPersistDonor failed: ", err)
			return fmt.Errorf("CacheAndPersistDonor failed: %v", err)
		}
	}

	// for each obj
	for _, obj := range objs {
		err := PutIndvDonor(year, obj)
		if err != nil {
			fmt.Println("CacheAndPersistIndvDonor failed: ", err)
			return fmt.Errorf("CacheAndPersistIndvDonor failed: %v", err)
		}
	}
	return nil
}

// PutIndvDonor puts an Individual object belonging to the specified year to the database
func PutIndvDonor(year string, donor *donations.Individual) error {
	// get name & job
	name, job := donor.Name, fmt.Sprintf("%s - %s", donor.Employer, donor.Occupation)

	// convert obj to protobuf
	data, err := encodeIndv(*donor)
	if err != nil {
		fmt.Println("PutIndvDonor failed: ", err)
		return fmt.Errorf("petIndvDonor failed: %v", err)
	}
	// open/create bucket in db/offline_db.db
	// put protobuf item and use donor.ID as key
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: PutIndvDonor failed: 'offline_db.db' failed to open")
		return fmt.Errorf("PutIndvDonor failed: 'offline_db.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte("indv_donors"))
		if err := b.Put([]byte(donor.ID), data); err != nil { // serialize k,v
			fmt.Printf("PutIndvDonor failed: offline_db.db': failed to store donor: %s\n", donor.ID)
			return fmt.Errorf("PutIndvDonor failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: PutIndvDonor failed: 'offline_db.db': 'indv_donors' bucket failed to open")
		return fmt.Errorf("PutIndvDonor failed: 'offline_db.db': 'indv_donors' bucket failed to open: %v", err)
	}

	err = RecordIDByJob(name, job, donor.ID)
	if err != nil {
		fmt.Println("PutIndvDonor failed: RecordIDByJob failed")
		return fmt.Errorf("PutIndvDonor failed: %v", err)
	}

	return nil
}

// GetIndvDonor returns a pointer to an Indvidual donor obj stored on disk
func GetIndvDonor(year, id string) (*donations.Individual, error) {
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: GetIndvDonor failed: 'offline_db.db' failed to open")
		return nil, fmt.Errorf("GetIndvDonor failed: 'offline_db.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(year)).Bucket([]byte("indv_donors")).Get([]byte(id))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetIndvDonor failed: 'offline_db.db': 'indv_donors' bucket failed to open")
		return nil, fmt.Errorf("GetIndvDonor failed: 'offline_db.db': 'indv_donors' bucket failed to open: %v", err)
	}

	indv, err := decodeIndv(data)
	if err != nil {
		fmt.Println("GetIndvDonor failed: convProtoToIndv failed: ", err)
		return nil, fmt.Errorf("GetIndvDonor failed: convProtoToIndv failed: %v", err)
	}

	return &indv, nil
}


func deriveID(s string) int32 {
	conv := strings.Split(s, "")
	var ID []string
	for i, n := range conv {
		if i > 1 && n != "0" {
			ID = conv[i:]
			break
		}
	}
	donorIDint, _ := strconv.Atoi(strings.Join(ID, ""))
	donorID := int32(donorIDint)

	return donorID
}

 func updateIndvDonor(new *donations.Individual) error {
	// get old value
	old, err := GetIndvDonor(new.ID)
	if err != nil {
		fmt.Println("updateIndvDonor failed: ", err)
		return fmt.Errorf("updateIndvDonor failed: %v", err)
	}

	// Add old values to new struct
	for _, donation := range old.Donations {
		new.Donations = append(new.Donations, donation)
	}
	fmt.Printf("%s TotalDonations: %d + %d = %d\n", new.ID, old.TotalDonations, new.TotalDonations, old.TotalDonations+new.TotalDonations)
	fmt.Printf("%s TotalDonated: %d + %d = %d\n", new.ID, old.TotalDonated, new.TotalDonated, old.TotalDonated+new.TotalDonated)
	fmt.Println()
	new.TotalDonations += old.TotalDonations
	new.TotalDonated += old.TotalDonated
	// replace/overwrite old value
	err = PutIndvDonor(new)
	if err != nil {
		fmt.Println("updateIndvDonor failed: PutIndvDonor failed: ", err)
		return fmt.Errorf("updateIndvDonor failed: PutIndvDonor failed: %v", err)
	}
	return nil
} */

// CacheAndPersistIndvDonor persists a list of of Individual donor objects to the on-disk cache
/* func cacheAndPersistIndvDonorOld(objs []*donations.Individual, seen map[string]bool) error {
	if len(seen) == 0 {
		err := createBucket("indv_donors")
		if err != nil {
			fmt.Println("CacheAndPersistDonor failed: ", err)
			return fmt.Errorf("CacheAndPersistDonor failed: %v", err)
		}
	}

	// for each obj
	for _, obj := range objs {
		// check if seen
		// check in memory cache map
		if !seen[obj.ID] {
			// if not seen, open View Tx
			item, err := GetIndvDonor(obj.ID)
			if err != nil {
				fmt.Println("CacheAndPersistDonor failed: GetIndvDonor failed: ", err)
				return fmt.Errorf("CacheAndPersistDonor failed: GetIndvDonor failed: %v", err)
			}
			if item == nil { // confirm item has not already been saved
				// if db[key] == nil, put item, update cache
				err := PutIndvDonor(obj)
				if err != nil {
					fmt.Println("CacheAndPersist failed: PutIndvDonor failed: ", err)
					return fmt.Errorf("CacheAndPersist failed: PutIndvDonor failed: %v", err)
				}
			} else { // item has been saved but does not appear in memory
				// else, update item
				err := updateIndvDonor(obj)
				if err != nil {
					fmt.Println("CacheAndPersistDonor failed: updateIndvDonor failed: ", err)
					return fmt.Errorf("CacheAndPersistDonor failed: updateIndvDonor failed: %v", err)
				}
			}
			seen[obj.ID] = true
		} else { // seen[obj.ID] == true
			err := updateIndvDonor(obj)
			if err != nil {
				fmt.Println("CacheAndPersistDonor failed: updateIndvDonor failed: ", err)
				return fmt.Errorf("CacheAndPersistDonor failed: updateIndvDonor failed: %v", err)
			}

		}
	}
	return nil
} */
