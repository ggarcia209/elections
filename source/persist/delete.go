// Package persist contains operations for reading and writing disk data.
// Most operations in this package are intended to be performed on the
// admin local machine and are not intended to be used in the service logic.
// This file contains operations for deleting various datasets from disk.
// Operations in this file should not be executed without out explicit
// confirmation from end user.
package persist

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

// DeleteDatabase deletes the entire database file.
func DeleteDatabase() error {
	path := OUTPUT_PATH + "/db/offline_db.db"
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteDatabse failed: %v", err)
	}
	return nil
}

// DeleteSearchIndex deletes the search index.
func DeleteSearchIndex() error {
	path := OUTPUT_PATH + "/db/search_index.db"
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteSearchIndex failed: %v", err)
	}
	return nil
}

// DeleteMetaData deletes the application & database metadata.
func DeleteMetaData() error {
	path := "../db/disk_cache.db"
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteSearchIndex failed: %v", err)
	}
	return nil
}

// DeleteYear deletes the data set for the given year.
func DeleteYear(year string) error {
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteYear failed: %v", err)
	}
	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket([]byte(year)); err != nil {
			return fmt.Errorf("tx failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteYear failed: %v", err)
	}

	// delete corresponding offsets
	db, err = bolt.Open("../db/disk_cache.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteYear failed: %v", err)
	}
	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("offsets"))
		if err := b.DeleteBucket([]byte(year)); err != nil {
			return fmt.Errorf("tx failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteYear failed: %v", err)
	}
	return nil
}

// DeleteCategory deletes the selected category for the given year.
func DeleteCategory(year, category string) error {
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteCategory failed: %v", err)
	}
	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte(category))
		if err := b.DeleteBucket([]byte(category)); err != nil { // serialize k,v
			return fmt.Errorf("tx failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteCategory failed: %v", err)
	}
	return nil
}

// DeleteAll deletes the output and disk cache directories
func DeleteAll() error {
	err := os.RemoveAll(OUTPUT_PATH + "/db")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteAll failed: %v", err)
	}
	err = os.RemoveAll("../db")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteAll failed: %v", err)
	}
	return nil
}
