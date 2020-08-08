package persist

import (
	"fmt"
	"sync"

	"github.com/elections/source/util"

	"github.com/boltdb/bolt"
)

var mu = &sync.Mutex{}

// LogOffset records the byte offset value in the database
func LogOffset(year, key string, offset int64) error {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("LogOffset failed: %v", err)
	}
	defer db.Close()

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("offsets"))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		y, err := b.CreateBucketIfNotExists([]byte(year))
		if err := y.Put([]byte(key), util.Itob(int(offset))); err != nil { // serialize k,v
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("LogOffset failed: %v", err)
	}
	return nil
}

// GetOffset retreives the offset value from the database in the event of failure
func GetOffset(year, key string) (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	var val int
	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("GetOffset failed: %v", err)
	}
	defer db.Close()

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("offsets"))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		y, err := b.CreateBucketIfNotExists([]byte(year))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		val = util.Btoi(y.Get([]byte(key)))
		return nil
	}); err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("GetOffset failed: %v", err)
	}
	return int64(val), nil
}

// LogKey logs the key of the last object uploaded to DynamoDB for the given year/bucket
func LogKey(year, bucket, key string) error {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("LogKeys failed: %v", err)
	}
	defer db.Close()

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("upload_keys"))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		y, err := b.CreateBucketIfNotExists([]byte(year))
		if err := y.Put([]byte(bucket), []byte(key)); err != nil { // serialize k,v
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("LogKeys failed: %v", err)
	}
	return nil
}

// GetKey retreives the key of the last object uploaded to DynamoDB for the given year/bucket
// returns nil if none
func GetKey(year, bucket string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	var key string
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("GetKeys failed: %v", err)
	}
	defer db.Close()

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("upload_keys"))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		y, err := b.CreateBucketIfNotExists([]byte(year))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		key = string(y.Get([]byte(bucket)))
		return nil
	}); err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("GetKeys failed: %v", err)
	}
	return key, nil
}
