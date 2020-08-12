// Package persist contains functions to persist data stored in memory to the local disk
package persist

import (
	"fmt"
	"os"

	"github.com/elections/source/donations"

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
	// open/create bucket in db/offline_db.db
	// put protobuf item and use donor.ID as key
	db, err := bolt.Open("../../db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("StoreObjects failed: ", err)
		return fmt.Errorf("StoreObjects failed: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		for _, obj := range objs {
			// encode object
			bucket, key, data, err := encodeToProto(obj)
			if err != nil {
				fmt.Println("StoreObjects failed: ", err)
				return fmt.Errorf("StoreObjects failed: %v", err)
			}

			b := tx.Bucket([]byte(year)).Bucket([]byte(bucket))
			if err := b.Put([]byte(key), data); err != nil { // serialize k,v
				return fmt.Errorf("StoreObjects failed: %v", err)
			}
		}
		return nil
	}); err != nil {
		fmt.Println("StoreObjects failed: ", err)
		return fmt.Errorf("StoreObjecs failed: %v", err)
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
	db, err := bolt.Open("../../db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("PutObject failed: ", err)
		return fmt.Errorf("PutObject failed: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte(bucket))
		if err := b.Put([]byte(key), data); err != nil { // serialize k,v
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
	db, err := bolt.Open("../../db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(year)).Bucket([]byte(bucket)).Get([]byte(key))
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}

	obj, err := decodeFromProto(bucket, data) // change to decode
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}

	return obj, nil
}

// BatchGetSequential retrieves a sequential list of n objects from the database starting at the given key.
func BatchGetSequential(year, bucket, startKey string, n int) ([]interface{}, string, error) {
	objs := []interface{}{}
	currKey := startKey

	db, err := bolt.Open("../../db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return nil, currKey, fmt.Errorf("BatchGetSequential failed: %v", err)
	}

	if err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(year)).Bucket([]byte(bucket))

		c := b.Cursor()

		if startKey == "" {
			skBytes, _ := c.First()
			startKey = string(skBytes)
		}

		for k, v := c.Seek([]byte(startKey)); k != nil; k, v = c.Next() {
			obj, err := decodeFromProto(bucket, v)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed: %v", err)
			}
			objs = append(objs, obj)
			ckBytes, _ := c.Next()
			currKey = string(ckBytes)
			if len(objs) == n {
				break
			}
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, currKey, fmt.Errorf("BatchGetSequential failed: %v", err)
	}

	return objs, currKey, nil
}

// BatchGetByID returns a list of objects for the given IDs contained within the given bucket.
// Returns a list of non-nil objects of same type and list of IDs that returned nil objects.
func BatchGetByID(year, bucket string, IDs []string) ([]interface{}, []string, error) {
	objs := []interface{}{}
	nilIDs := []string{}

	db, err := bolt.Open("../../db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return nil, nil, fmt.Errorf("BatchGetByID failed: %v", err)
	}

	if err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(year)).Bucket([]byte(bucket))

		for _, id := range IDs {
			data := b.Get([]byte(id))
			obj, err := decodeFromProto(bucket, data)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed: %v", err)
			}
			if obj == nil {
				nilIDs = append(nilIDs, id)
			} else {
				objs = append(objs, obj)
			}
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, nil, fmt.Errorf("BatchGetByID failed: %v", err)
	}

	return objs, nilIDs, nil
}

// GetTopOverall retreives the TopOverall objects from disk to store in memory
func GetTopOverall(year string) ([]interface{}, error) {
	objs := []interface{}{}

	db, err := bolt.Open("../../db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("GetObject failed: ", err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}

	var data [][]byte

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(year)).Bucket([]byte("top_overall"))

		c := b.Cursor()

		i := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			data = append(data, v)
			i++
		}

		return nil
	})

	for _, bs := range data {
		obj, err := decodeFromProto("top_overall", bs)
		if err != nil {
			fmt.Println("GetTopOverall failed: ", err)
			return nil, fmt.Errorf("GetTopOverall failed: %v", err)
		}
		objs = append(objs, obj)
	}

	return objs, nil
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
	case "individuals":
		data, err := decodeIndv(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return &data, nil
	case "committees":
		data, err := decodeCmte(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return &data, nil
	case "candidates":
		data, err := decodeCand(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return &data, nil
	case "cmte_tx_data":
		data, err := decodeCmteTxData(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return &data, nil
	case "top_overall":
		data, err := decodeOverallData(data)
		if err != nil {
			fmt.Println("decodeFromProto failed: ", err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return &data, nil
	default:
		return nil, fmt.Errorf("decodeFromProto failed: invalid bucket")
	}
}

// CreateDB creates the database directory. CreateDB must be called
// before any other function in 'parse' package is called.
func createDB() {
	if _, err := os.Stat("../../db"); os.IsNotExist(err) {
		os.Mkdir("../../db", 0744)
		fmt.Println("CreateDB successful: 'db' directory created")
	}
}

func createObjBuckets(year string) error {
	buckets := []string{"individuals", "committees", "candidates", "cmte_tx_data", "top_overall"}
	for _, bucket := range buckets {
		err := createBucket(year, bucket)
		if err != nil {
			fmt.Println("createObjBuckets failed: ", err)
			return fmt.Errorf("createObjBuckets failed: %v", err)
		}
	}
	return nil
}

// createLookupBucket initializes the 'id_lookup' bucket in disk_cache.db
func createLookupBuckets() error {
	db, err := bolt.Open("../../db/disk_cache.db", 0644, nil)
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

func createBucket(year, name string) error {
	db, err := bolt.Open("../../db/offline_db.db", 0644, nil)
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
