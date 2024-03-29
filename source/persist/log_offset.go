// Package persist contains operations for reading and writing disk data.
// Most operations in this package are intended to be performed on the
// admin local machine and are not intended to be used in the service logic.
// This file contains operations for periodically persisting metadata for
// long running operations in the event of failure.
package persist

import (
	"fmt"
	"sync"

	"github.com/elections/source/util"

	"github.com/boltdb/bolt"
)

var mu = &sync.Mutex{}

// LogOffset records the byte offset of .txt input files while scanning
// and creating/updating datasets.
func LogOffset(year, key string, offset int64) error {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("../db/disk_cache.db", 0644, nil)
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
		if err := y.Put([]byte(key), util.Itob(offset)); err != nil { // serialize k,v
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

// GetOffset retreives the byte offset value from the database in the event of failure.
func GetOffset(year, key string) (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("../db/disk_cache.db", 0644, nil)
	var val int64
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
	return val, nil
}

// LogKey logs the key of the last object uploaded to DynamoDB for the given year/bucket.
// LogKey can also be implemented with other BatchWrite/BatchGet operations.
func LogKey(year, bucket, key string) error {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("../db/disk_cache.db", 0644, nil)
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
// returns nil if none.
func GetKey(year, bucket string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("../db/disk_cache.db", 0644, nil)
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

// LogPath saves the input/output file path set by admin to disk cache.
func LogPath(path string, input bool) error {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("../db/disk_cache.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("LogPath failed: %v", err)
	}
	defer db.Close()

	var key string
	if input {
		key = "input_path"
	} else {
		key = "output_path"
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("filepaths"))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		if err := b.Put([]byte(key), []byte(path)); err != nil { // serialize k,v
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("LogPath failed: %v", err)
	}
	return nil
}

// GetPath retreives the input/output filepaths.
// Returns nil if none.
func GetPath(input bool) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("../db/disk_cache.db", 0644, nil)
	var path string
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("GetPath failed: %v", err)
	}
	defer db.Close()

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("filepaths"))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}

		var key string
		if input {
			key = "input_path"
		} else {
			key = "output_path"
		}

		path = string(b.Get([]byte(key)))
		return nil
	}); err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("GetKeys failed: %v", err)
	}
	return path, nil
}
