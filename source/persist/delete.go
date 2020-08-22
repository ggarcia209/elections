package persist

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

// DeleteDatabase deletes the database file
func DeleteDatabase() error {
	path := OUTPUT_PATH + "/db/offline_db.db"
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteDatabse failed: %v", err)
	}
	return nil
}

// DeleteSearchIndex deletes the search index
func DeleteSearchIndex() error {
	path := OUTPUT_PATH + "/db/search_index.db"
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteSearchIndex failed: %v", err)
	}
	return nil
}

// DeleteMetaData deletes the application & database metadata
func DeleteMetaData() error {
	path := "../db/disk_cache.db"
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeleteSearchIndex failed: %v", err)
	}
	return nil
}

// DeleteYear deletes the data set for the given year
func DeleteYear(year string) error {
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("DeleteYear failed: ", err)
		return fmt.Errorf("DeleteYear failed: %v", err)
	}
	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year))
		if err := b.DeleteBucket([]byte(year)); err != nil { // serialize k,v
			return fmt.Errorf("tx failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("DeleteYear failed: ", err)
		return fmt.Errorf("DeleteYear failed: %v", err)
	}
	return nil
}

// DeleteCategory deletes the selected category for the given year
func DeleteCategory(year, category string) error {
	db, err := bolt.Open(OUTPUT_PATH+"/db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("DeleteCategory failed: ", err)
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
		fmt.Println("DeleteCategory failed: ", err)
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
