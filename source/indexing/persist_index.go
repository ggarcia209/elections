package indexing

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/boltdb/bolt"
	"github.com/elections/source/protobuf"
	"github.com/golang/protobuf/proto"
)

// save new Index entries
func saveIndex(index indexMap) (int, error) {
	i := 0 // new terms
	u := 0 // updated terms
	fmt.Println("save index - writing objects to db/search_index.db")

	// open/create bucket in db/offline_db.db
	// put protobuf item and use donor.ID as key
	db, err := bolt.Open("../../db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return i, fmt.Errorf("saveIndex failed: %v", err)
	}

	// add/update each term for each partition
	for prt, terms := range index {
		// tx
		if err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(prt))
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed failed: %v", err)
			}

			for term, lookup := range terms {
				// get previous data
				prevData := b.Get([]byte(term))

				if prevData != nil { // previous data exists
					prev, err := decodeSearchEntry(prevData)
					if err != nil {
						fmt.Println(err)
						return fmt.Errorf("tx failed failed: %v", err)
					}

					// update & overwrite with updated copy
					mergeLookup(lookup, prev)
					data, err := encodeSearchEntry(prev)
					if err := b.Put([]byte(term), data); err != nil { // serialize k,v
						return fmt.Errorf("tx failed failed: %v", err)
					}
					u++
				} else { // new term
					data, err := encodeSearchEntry(lookup)
					if err != nil {
						fmt.Println(err)
						return fmt.Errorf("tx failed failed: %v", err)
					}
					if err := b.Put([]byte(term), data); err != nil { // serialize k,v
						return fmt.Errorf("tx failed failed: %v", err)
					}
					i++
				}
			}
			return nil
		}); err != nil {
			fmt.Println(err)
			return i, fmt.Errorf("saveIndex failed: %v", err)
		}
	}
	fmt.Println("index saved!")
	fmt.Println("new items wrote: ", i)
	fmt.Println("existing items updated: ", u)

	fmt.Println()
	return i, nil
}

// get lookup pairs from disk for given term
func getSearchEntry(term string) (lookupPairs, error) {
	// retreive lookupPairs from disk
	prt := getPartition(term)
	db, err := bolt.Open("../../db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return lookupPairs{}, fmt.Errorf("getIndexData failed: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(prt)).Get([]byte(term))
		return nil
	}); err != nil {
		fmt.Println(err)
		return lookupPairs{}, fmt.Errorf("getIndexData failed: %v", err)
	}

	lookup, err := decodeSearchEntry(data)
	if err != nil {
		fmt.Println(err)
		return lookupPairs{}, fmt.Errorf("getIndexData failed: %v", err)
	}
	return lookup, nil
}

// persist IndexData object to disk
func saveIndexData(index *IndexData) error {
	db, err := bolt.Open("../../db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("saveIndexData failed: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("index_data"))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed failed: %v", err)
		}
		data, err := encodeIndexData(index)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed failed: %v", err)
		}
		if err := b.Put([]byte("data"), data); err != nil {
			return fmt.Errorf("saveIndexData failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("saveIndexData failed: %v", err)
	}
	return nil
}

// retreive IndexData from disk
func getIndexData() (*IndexData, error) {
	db, err := bolt.Open("../../db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return &IndexData{}, fmt.Errorf("getIndexData failed: %v", err)
	}

	var data []byte

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		d, err := tx.CreateBucketIfNotExists([]byte("index_data"))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed: %v", err)
		}
		data = d.Get([]byte("data"))
		return nil
	}); err != nil {
		fmt.Println(err)
		return &IndexData{}, fmt.Errorf("getIndexData failed: %v", err)
	}

	id, err := decodeIndexData(data)
	if err != nil {
		fmt.Println(err)
		return &IndexData{}, fmt.Errorf("getIndexData failed: %v", err)
	}

	return id, nil
}

// save PartitionMap
func savePartitionMap(pm map[string]bool) error {
	db, err := bolt.Open("../../db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("saveIndexData failed: %v", err)
	}

	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("index_data"))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed failed: %v", err)
		}
		data, err := encodePartitionMap(pm)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("tx failed failed: %v", err)
		}
		if err := b.Put([]byte("partition_map"), data); err != nil {
			return fmt.Errorf("saveIndexData failed: %v", err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("saveIndexData failed: %v", err)
	}
	return nil
}

// get PartitionMap
func getPartitionMap() (map[string]bool, error) {
	db, err := bolt.Open("../../db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("getPartition failed: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte("index_data")).Get([]byte("partition_map"))
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("getPartitionMap failed: %v", err)
	}

	pm, err := decodePartitionMap(data)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("getPartitionMap failed: %v", err)
	}
	return pm, nil
}

// encode lookupPairs to protobuf
func encodeSearchEntry(m lookupPairs) ([]byte, error) {
	lookup := make(map[string]*protobuf.SearchResult)
	for k, sd := range m {
		lookup[k] = encodeSearchData(sd)
	}
	entry := &protobuf.SearchEntry{
		Lookup: lookup,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeSearchEntry failed: %v", err)
	}
	return data, nil
}

// encode SearchData to protobuf
func encodeSearchData(sd *SearchData) *protobuf.SearchResult {
	entry := &protobuf.SearchResult{
		ID:     sd.ID,
		Name:   sd.Name,
		City:   sd.City,
		State:  sd.State,
		Bucket: sd.Bucket,
		Years:  sd.Years,
	}
	return entry
}

// add new lookupPairs to existing data
func mergeLookup(new, prev lookupPairs) {
	for k, v := range new {
		prev[k] = v
	}
}

// decode protobuf to lookupPairs
func decodeSearchEntry(input []byte) (lookupPairs, error) {
	data := &protobuf.SearchEntry{}
	err := proto.Unmarshal(input, data)
	if err != nil {
		fmt.Println(err)
		return make(lookupPairs), fmt.Errorf("decodeSearchEntry failed: %v", err)
	}

	pairs := make(lookupPairs)
	lookup := data.GetLookup()
	for k, v := range lookup {
		pairs[k] = decodeSearchData(v)
	}

	return pairs, nil
}

// decode protobuf to SearchData
func decodeSearchData(sr *protobuf.SearchResult) *SearchData {
	sd := &SearchData{
		ID:     sr.GetID(),
		Name:   sr.GetName(),
		City:   sr.GetCity(),
		State:  sr.GetState(),
		Bucket: sr.GetBucket(),
		Years:  sr.GetYears(),
	}
	return sd
}

// encode IndexData to protobuf
func encodeIndexData(id *IndexData) ([]byte, error) {
	entry := &protobuf.IndexData{
		Size:      float32(id.Size),
		Completed: id.Completed,
	}
	ts, err := ptypes.TimestampProto(id.LastUpdated)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeIndexData failed: %v", err)
	}
	entry.LastUpdated = ts
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeIndexData failed: %v", err)
	}
	return data, nil
}

// decode protobuf to IndexData
func decodeIndexData(raw []byte) (*IndexData, error) {
	data := &protobuf.IndexData{}
	err := proto.Unmarshal(raw, data)
	if err != nil {
		fmt.Println(err)
		return &IndexData{}, fmt.Errorf("decodeIndexData failed: %v", err)
	}
	id := &IndexData{
		Size:      int(data.GetSize()),
		Completed: data.GetCompleted(),
	}
	if data.GetLastUpdated() == nil {
		id.LastUpdated = time.Now()
	} else {
		ts, err := ptypes.Timestamp(data.GetLastUpdated())
		if err != nil {
			fmt.Println(err)
			return &IndexData{}, fmt.Errorf("decodeIndexData failed: %v", err)
		}
		id.LastUpdated = ts
	}

	return id, nil
}

// encode PartitionMap to protobuf
func encodePartitionMap(m map[string]bool) ([]byte, error) {
	entry := &protobuf.PartitionMap{
		Partitions: m,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodePartitionMap failed: %v", err)
	}
	return data, nil
}

// decode PartitionMap from protobuf
func decodePartitionMap(raw []byte) (map[string]bool, error) {
	data := &protobuf.PartitionMap{}
	err := proto.Unmarshal(raw, data)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("decodePartitionMap failed: %v", err)
	}
	m := data.GetPartitions()
	return m, nil

}
