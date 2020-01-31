package persist

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/protobuf"

	"github.com/golang/protobuf/proto"
)

// encode to protobuf
func encodeOverallData(od donations.TopOverallData) ([]byte, error) { // move conversions to protobuf package?
	entry := &protobuf.TopOverallData{
		Category:  od.Category,
		Amts:      od.Amts,
		Threshold: encodeThreshold(od.Threshold),
		SizeLimit: int32(od.SizeLimit),
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeOverallData failed: ", err)
		return nil, fmt.Errorf("encodeOverallData failed: %v", err)
	}
	return data, nil
}

func encodeThreshold(entries []*donations.Entry) []*protobuf.Entry {
	var es []*protobuf.Entry
	for _, e := range entries {
		entry := &protobuf.Entry{
			ID:    e.ID,
			Total: e.Total,
		}
		es = append(es, entry)
	}
	return es
}

// decode from protobuf
func decodeOverallData(data []byte) (donations.TopOverallData, error) {
	od := &protobuf.TopOverallData{}
	err := proto.Unmarshal(data, od)
	if err != nil {
		fmt.Println("convProtoToIndv failed: ", err)
		return donations.TopOverallData{}, fmt.Errorf("convProtoToIndv failed: %v", err)
	}

	entry := donations.TopOverallData{
		Category:  od.GetCategory(),
		Amts:      od.GetAmts(),
		Threshold: decodeThreshold(od.GetThreshold()),
		SizeLimit: int(od.GetSizeLimit()),
	}

	if len(entry.Amts) == 0 {
		entry.Amts = make(map[string]float32)
	}

	return entry, nil
}

func decodeThreshold(es []*protobuf.Entry) []*donations.Entry {
	var entries []*donations.Entry
	for _, e := range es {
		entry := donations.Entry{
			ID:    e.GetID(),
			Total: e.GetTotal(),
		}
		entries = append(entries, &entry)
	}
	return entries
}

// DEPRECATED
/*
// CacheTopOverall stores multipe TopOverallData objects
func CacheTopOverall(year string, objs []*donations.TopOverallData) error {
	err := createBucket(year, "top_overall")
	if err != nil {
		fmt.Println("CacheTopOverall failed: ", err)
		return fmt.Errorf("CacheTopOverall failed: %v", err)
	}

	for _, obj := range objs {
		err := PutTopOverallData(year, obj)
		if err != nil {
			fmt.Println("CacheTopOverall failed: putCandidate failed: ", err)
			return fmt.Errorf("CacheTopOverall failed: putCandidate failed: %v", err)
		}
	}
	return nil
}

// PutTopOverallData saves a TopOverallData obj to the database
func PutTopOverallData(year string, od *donations.TopOverallData) error {
	// convert obj to protobuf
	data, err := encodeOverallData(*od)
	if err != nil {
		fmt.Println("PutTopOverallData failed: ", err)
		return fmt.Errorf("PutTopOverallData failed: %v", err)
	}
	// open/create bucket in db/offline_db.db
	// put protobuf item and use cand.ID as key
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: PutTopOverallData failed: 'offline_db.db' failed to open")
		return fmt.Errorf("PutTopOverallData failed: 'offline_db.db' failed to open: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(year)).Bucket([]byte("top_overall"))
		if err := b.Put([]byte(od.Category), data); err != nil { // serialize k,v
			fmt.Printf("PutTopOverallData failed: offline_db.db': failed to store candidate: %s\n", od.Category)
			return fmt.Errorf("PutTopOverallData failed: could not update:\n%v", err)
		}
		return nil
	}); err != nil {
		fmt.Println("FATAL: PutTopOverallData failed: 'offline_db.db': 'top_overall' bucket failed to open")
		return fmt.Errorf("PutTopOverallData failed: 'offline_db.db': 'top_overall' bucket failed to open: %v", err)
	}

	return nil
}

// GetTopOverallData returns a pointer to a TopOverallData obj stored on disk
func GetTopOverallData(year, cat string) (*donations.TopOverallData, error) {
	db, err := bolt.Open("db/offline_db.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println("FATAL: GetTopOverallData failed: 'offline_db.db' failed to open")
		return nil, fmt.Errorf("GetTopOverallData failed: 'offline_db.db' failed to open: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(year)).Bucket([]byte("top_overall")).Get([]byte(cat))
		return nil
	}); err != nil {
		fmt.Println("FATAL: GetTopOverallData failed: 'offline_db.db': 'top_overall' bucket failed to open")
		return nil, fmt.Errorf("GetTopOverallData failed: 'offline_db.db': 'top_overall' bucket failed to open: %v", err)
	}

	od, err := decodeOverallData(data)
	if err != nil {
		fmt.Println("GetTopOverallData failed: decodeOverallData failed: ", err)
		return nil, fmt.Errorf("GetTopOverallData failed: decodeOverallData failed: %v", err)
	}

	return &od, nil
}

*/
