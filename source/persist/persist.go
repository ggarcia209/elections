// Package persist contains operations for reading and writing disk data.
// Most operations in this package are intended to be performed on the
// admin local machine and are not intended to be used in the service logic.
// This file contains the exported functions that are called by higher level packages.
package persist

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/elections/source/donations"

	"github.com/boltdb/bolt"
)

// OUTPUT_PATH sets the output path for all database & search index storage & retreiveal operations.
var OUTPUT_PATH = "." // default value

// InitDiskCache creates the disk cache directory next to the application.
// This function is called by each of the admin, server, and index services.
func InitDiskCache() {
	// metadata - store in /admin_app
	if _, err := os.Stat("../db"); os.IsNotExist(err) {
		os.Mkdir("../db", 0744)
		fmt.Printf("InitDiskCache successful: '../db' directory created\n")
	}
}

// Init inititalizes the admin program by creating the db directory, disk_cache.db, offline_db.db, and corresponding buckets.
func Init(year string) error {
	createDB()

	err := createObjBuckets(year)
	if err != nil {
		fmt.Println("Init failed: ", err)
		return fmt.Errorf("Init failed: %v", err)
	}

	fmt.Println("Init Done")
	return nil
}

// StoreObjects persists a list of objects to the on-disk database as a batch write transaction.
func StoreObjects(year string, objs []interface{}) error {
	// open/create bucket in db/offline_db.db
	// put protobuf item and use donor.ID as key
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("StoreObjects failed: %v", err)
	}
	defer db.Close()

	// tx

	if err := db.Update(func(tx *bolt.Tx) error {
		for _, obj := range objs {
			// encode object
			bucket, key, data, err := encodeToProto(obj)
			if err != nil {
				fmt.Println(err)
				fmt.Println("obj: ", obj)
				return fmt.Errorf("tx failed: %v", err)
			}

			b := tx.Bucket([]byte(year)).Bucket([]byte(bucket))
			if err := b.Put([]byte(key), data); err != nil { // serialize k,v
				fmt.Println("obj: ", obj)
				return fmt.Errorf("tx failed: %v", err)
			}
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("StoreObjecs failed: %v", err)
	}

	return nil
}

// PutObject puts an object by year:bucket:key.
func PutObject(year string, object interface{}) error {
	// encode object
	bucket, key, data, err := encodeToProto(object)
	if err != nil {
		fmt.Println("PutObject failed: ", err)
		return fmt.Errorf("PutObject failed: %v", err)
	}

	// open/create bucket in db/offline_db.db
	// put protobuf item and use donor.ID as key
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	if err != nil {
		fmt.Println("PutObject failed: ", err)
		return fmt.Errorf("PutObject failed: %v", err)
	}
	defer db.Close()

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

// GetObject gets an object by year:bucket:key and returns it as an interface.
func GetObject(year, bucket, key string) (interface{}, error) {
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}
	defer db.Close()

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(year)).Bucket([]byte(bucket)).Get([]byte(key))
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}

	obj, err := decodeFromProto(bucket, data)
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

	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
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
			currKey = string(k)
			if len(objs) == n {
				break
			}

		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, currKey, fmt.Errorf("BatchGetSequential failed: %v", err)
	}
	if len(objs) < n {
		currKey = ""
	}
	return objs, currKey, nil
}

// BatchGetByID returns a list of objects for the given IDs contained within the given bucket.
// Returns a list of non-nil objects of same type and list of IDs that returned nil objects.
func BatchGetByID(year, bucket string, IDs []string) ([]interface{}, []string, error) {
	objs := []interface{}{}
	nilIDs := []string{}

	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
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
			if data == nil {
				nilIDs = append(nilIDs, id)
				continue
			}
			obj, err := decodeFromProto(bucket, data)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed: %v", err)
			}

			objs = append(objs, obj)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, nil, fmt.Errorf("BatchGetByID failed: %v", err)
	}

	return objs, nilIDs, nil
}

// GetTopOverall retreives the TopOverall objects from disk to store in memory.
func GetTopOverall(year string) ([]interface{}, error) {
	objs := []interface{}{}

	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}
	defer db.Close()

	if err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(year)).Bucket([]byte("top_overall"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v == nil {
				fmt.Printf("nil object: %s\n", string(k))
				continue
			}
			obj, err := decodeFromProto("top_overall", v)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed: %v", err)
			}
			objs = append(objs, obj)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetTopOverall failed: %v", err)
	}

	return objs, nil
}

// SaveTopOverall saves a list of TopOverall objects by year/bucket/category.
func SaveTopOverall(year, bucket string, ods []interface{}) error {
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("SaveTopOverall failed: %v", err)
	}
	defer db.Close()

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte("top_overall"))

		for _, od := range ods {
			_, key, data, err := encodeToProto(od)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed %v", err)
			}
			if err = b.Put([]byte(key), data); err != nil { // serialize k,v
				return fmt.Errorf("tx failed %v", err)
			}
		}

		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("SaveTopOverall failed: %v", err)
	}
	return nil
}

// GetYearlyTotals retreives the Yearly objects from disk for the given year/cateogry.
func GetYearlyTotals(year, cat string) ([]interface{}, error) {
	objs := []interface{}{}

	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}
	defer db.Close()

	if err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(year)).Bucket([]byte("yearly_totals"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			obj, err := decodeFromProto("yearly_totals", v)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed: %v", err)
			}
			objs = append(objs, obj)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetYearlyTotals failed: %v", err)
	}

	return objs, nil
}

// SaveYearlyTotals saves a list of YearlyTotal objects by year/category
func SaveYearlyTotals(year, cat string, yts []interface{}) error {
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("SaveTopOverall failed: %v", err)
	}
	defer db.Close()

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte("yearly_totals"))

		for _, yt := range yts {
			_, key, data, err := encodeToProto(yt)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed %v", err)
			}
			if err = b.Put([]byte(key), data); err != nil { // serialize k,v
				return fmt.Errorf("tx failed %v", err)
			}
		}

		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("SaveYearlyTotals failed: %v", err)
	}
	return nil
}

// ViewDataByBucket prints 1000 objects stored at at the given year/bucket
// per function call and returns the last key printed.
// Enter "" to start at first key.
func ViewDataByBucket(year, bucket, start string) (string, error) {
	fmt.Printf("Displaying data for %s - %s: \n", year, bucket)
	i := 0
	curr := start
	keyN := 0

	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("ViewDataByBucket failed: %v", err)
	}

	if err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(year)).Bucket([]byte(bucket))
		keyN = b.Stats().KeyN

		c := b.Cursor()

		if start == "" {
			skBytes, _ := c.First()
			start = string(skBytes)
		}

		for k, _ := c.Seek([]byte(start)); k != nil; k, _ = c.Next() {
			fmt.Printf("%d) %s - %s:\t%s\n", i, year, bucket, string(k))
			i++
			if i == 1000 {
				break
			}
			curr = string(k)
		}
		if i < 1000 {
			curr = ""
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("ViewDataByBucket failed: %v", err)
	}
	fmt.Printf("items in bucket %s: %d\n", bucket, keyN)
	fmt.Println("items scanned: ", i)
	return curr, nil
}

// encodeToProto encodes an object interface to protobuf.
// Returns bucket, key (ID), serialized data, and err.
func encodeToProto(obj interface{}) (string, string, []byte, error) {
	if obj == nil {
		fmt.Println("empty object passed to encodeProto")
		return "", "", nil, nil
	}
	switch obj.(type) {
	case nil:
		return "", "", nil, fmt.Errorf("encodeToProto failed: nil interface")
	case *donations.Individual:
		bucket := "individuals"
		key := obj.(*donations.Individual).ID
		data, err := encodeIndv(*obj.(*donations.Individual))
		if err != nil {
			fmt.Println(err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.Committee:
		bucket := "committees"
		key := obj.(*donations.Committee).ID
		data, err := encodeCmte(*obj.(*donations.Committee))
		if err != nil {
			fmt.Println(err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.Candidate:
		bucket := "candidates"
		key := obj.(*donations.Candidate).ID
		data, err := encodeCand(*obj.(*donations.Candidate))
		if err != nil {
			fmt.Println(err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.CmteTxData:
		bucket := "cmte_tx_data"
		key := obj.(*donations.CmteTxData).CmteID
		data, err := encodeCmteTxData(*obj.(*donations.CmteTxData))
		if err != nil {
			fmt.Println(err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.CmpnFinancials:
		bucket := "cmpn_fin"
		key := obj.(*donations.CmpnFinancials).CandID
		data, err := encodeCmpnFinancials(*obj.(*donations.CmpnFinancials))
		if err != nil {
			fmt.Println(err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.CmteFinancials:
		bucket := "cmte_fin"
		key := obj.(*donations.CmteFinancials).CmteID
		data, err := encodeCmteFinancials(*obj.(*donations.CmteFinancials))
		if err != nil {
			fmt.Println(err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.TopOverallData:
		bucket := "top_overall"
		key := obj.(*donations.TopOverallData).ID
		data, err := encodeOverallData(*obj.(*donations.TopOverallData))
		if err != nil {
			fmt.Println(err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	case *donations.YearlyTotal:
		bucket := "yearly_totals"
		key := obj.(*donations.YearlyTotal).ID
		data, err := encodeYrTotal(*obj.(*donations.YearlyTotal))
		if err != nil {
			fmt.Println(err)
			return "", "", nil, fmt.Errorf("encodeToProto failed: %v", err)
		}
		return bucket, key, data, nil
	default:
		return "", "", nil, fmt.Errorf("encodeToProto failed: invalid interface type")
	}
}

// decodeFromProto decodes a protobuf object.
// Returns pointer to deserialized object as interface{}
func decodeFromProto(bucket string, data []byte) (interface{}, error) {
	switch bucket {
	case "":
		return nil, fmt.Errorf("decodeFromProto failed: nil bucket")
	case "individuals":
		data, err := decodeIndv(data)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		if data.State == "" { // set DynamoDB partition key value if none
			data.State = "???"
		}
		return &data, nil
	case "committees":
		data, err := decodeCmte(data)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		if data.State == "" {
			data.State = "???"
		}
		return &data, nil
	case "candidates":
		data, err := decodeCand(data)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		if data.State == "" {
			data.State = "???"
		}
		return &data, nil
	case "cmte_tx_data":
		data, err := decodeCmteTxData(data)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		if data.Party == "" {
			data.Party = "???"
		}
		return &data, nil
	case "cmpn_fin":
		data, err := decodeCmpnFinancials(data)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return &data, nil
	case "cmte_fin":
		data, err := decodeCmteFinancials(data)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		return &data, nil
	case "top_overall":
		data, err := decodeOverallData(data)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		if data.Year == "" {
			data.Year = "0000"
		}
		return &data, nil
	case "yearly_totals":
		data, err := decodeYrTotal(data)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("decodeFromProto failed: %v", err)
		}
		if data.Year == "" {
			data.Year = "0000"
		}
		return &data, nil
	default:
		return nil, fmt.Errorf("decodeFromProto failed: invalid bucket")
	}
}

// CreateDB creates the database directory. CreateDB must be called
// before any other function in 'parse' package is called.
func createDB() {
	// create output directory
	if _, err := os.Stat(filepath.Join(OUTPUT_PATH, "db")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(OUTPUT_PATH, "db"), 0744)
		fmt.Printf("CreateDB successful: '%s/db' directory created\n", OUTPUT_PATH)
		fmt.Println()
	}
}

// initializes BoltDB buckets datasets are stored in on disk
func createObjBuckets(year string) error {
	buckets := []string{"individuals", "committees", "candidates", "cmte_tx_data", "cmte_fin", "cmpn_fin", "top_overall", "yearly_totals"}
	for _, bucket := range buckets {
		err := createBucket(year, bucket)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("createObjBuckets failed: %v", err)
		}
	}

	return nil
}

// create an individual boltDB bucket
func createBucket(year, name string) error {
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("createBucket failed: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		yb, err := tx.CreateBucketIfNotExists([]byte(year))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("createBucket failed: %v", err)
		}
		_, err = yb.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("createBucet failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("createBucket failed: %v", err)
	}

	return nil
}
