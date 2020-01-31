package persist

import (
	"fmt"
	"sync"

	"github.com/elections/util"

	"github.com/boltdb/bolt"
)

var mu = &sync.Mutex{}

// logOffset records the byte offset value in the database
func LogOffset(key string, offset int64) error {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	if err != nil {
		fmt.Println("FATAL: logOffset failed: 'disk_cache.db' failed to open")
		return fmt.Errorf("logOffset failed: 'disk_cache.db' failed to open: %v", err)
	}
	defer db.Close()

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("offsets"))
		if err != nil {
			fmt.Println("FATAL: logOffset failed: 'disk_cache.db': 'offsets' bucket failed to open")
			return fmt.Errorf("'main': FATAL: 'disk_cache.db': 'offsets' bucket failed to open: %v", err)
		}
		if err := b.Put([]byte(key), util.Itob(int(offset))); err != nil { // serialize k,v
			fmt.Printf("logOffset failed: disk_cache.db': '%s': '%v' failed to store\n", key, offset)
			return fmt.Errorf("logOffset failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: logOffset failed: 'disk_cache.db': 'offsets' bucket failed to open")
		return fmt.Errorf("logOffset failed: 'disk_cache.db': 'offsets' bucket failed to open: %v", err)
	}
	return nil
}

// GetOffset retreives the offset value from the database in the event of failure
func GetOffset(key string) (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	db, err := bolt.Open("db/disk_cache.db", 0644, nil)
	var val int
	if err != nil {
		fmt.Println("FATAL: GetOffset failed: 'disk_cache.db' failed to open")
		return 0, fmt.Errorf("GetOffset failed: 'disk_cache.db' failed to open: %v", err)
	}
	defer db.Close()

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("offsets"))
		if err != nil {
			fmt.Println("FATAL: GetOffset failed: 'disk_cache.db': 'offsets' bucket failed to open")
			return fmt.Errorf("'main': FATAL: 'disk_cache.db': 'offsets' bucket failed to open: %v", err)
		}
		val = util.Btoi(b.Get([]byte(key)))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetOffset failed: 'disk_cache.db': 'offsets' bucket failed to open")
		return 0, fmt.Errorf("GetOffset failed: 'disk_cache.db': 'offsets' bucket failed to open: %v", err)
	}
	return int64(val), nil
}
