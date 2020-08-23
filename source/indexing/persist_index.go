package indexing

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/boltdb/bolt"
	"github.com/elections/source/protobuf"
	"github.com/golang/protobuf/proto"
)

// OUTPUT_PATH is used to set the directory to write the search index to
var OUTPUT_PATH string

type resultList struct {
	Results []string
}

// save new Index entries
func saveIndex(index indexMap, lookup lookupPairs) (int, error) {
	i := 0 // new terms
	u := 0 // updated terms
	n := 0 // new IDs/Search Data pairs wrote
	wrote := make(map[string]bool)
	fmt.Println("save index - writing objects to db/search_index.db")

	// open/create bucket in db/offline_db.db
	// put protobuf item and use donor.ID as key
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return i, fmt.Errorf("saveIndex failed: %v", err)
	}

	// persist inverted index
	// add/update each term for each partition
	for prt, terms := range index {
		// tx
		if err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(prt))
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed failed: %v", err)
			}

			for term, ids := range terms {
				// get previous data
				prevData := b.Get([]byte(term))

				if prevData != nil { // previous data exists
					prev, err := decodeResultsList(prevData)
					if err != nil {
						fmt.Println(err)
						return fmt.Errorf("tx failed failed: %v", err)
					}

					// update & overwrite with updated copy
					update := mergeIDs(ids, prev)
					r := resultList{Results: update}
					data, err := encodeResultsList(r)
					if err := b.Put([]byte(term), data); err != nil { // serialize k,v
						return fmt.Errorf("tx failed failed: %v", err)
					}
					u++
				} else { // new term - no existing data
					r := resultList{Results: ids}
					data, err := encodeResultsList(r)
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

	// check if SearchData persisted for corresponding ID; write new if not exist
	// 2nd transaction
	if err := db.Update(func(tx *bolt.Tx) error {
		for ID, sd := range lookup {
			// check if wrote to disk previously
			if wrote[ID] == false { // not wrote to disk in this function call
				luB, err := tx.CreateBucketIfNotExists([]byte("lookup"))
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("tx failed failed: %v", err)
				}
				data, err := encodeSearchData(sd)
				if err != nil {
					fmt.Println()
					return fmt.Errorf("tx failed: %v", err)
				}
				if err := luB.Put([]byte(ID), data); err != nil { // serialize k,v
					return fmt.Errorf("tx failed failed: %v", err)
				}
				n++
				wrote[ID] = true
			}
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return i, fmt.Errorf("saveIndex failed: %v", err)
	}
	fmt.Println("index saved!")
	fmt.Println("new terms wrote: ", i)
	fmt.Println("existing terms: ", u)
	fmt.Println("IDs recorded/updated: ", n)
	fmt.Println()
	return i, nil
}

// get lookup pairs from disk for given term
func getSearchEntry(term string) ([]SearchData, error) {
	// retreive lookupPairs from disk
	prt := getPartition(term)
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return []SearchData{}, fmt.Errorf("getIndexData failed: %v", err)
	}

	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(prt)).Get([]byte(term))
		return nil
	}); err != nil {
		fmt.Println(err)
		return []SearchData{}, fmt.Errorf("getSearchEntry failed: %v", err)
	}

	ids, err := decodeResultsList(data)
	if err != nil {
		fmt.Println(err)
		return []SearchData{}, fmt.Errorf("getSearchEntry failed: %v", err)
	}

	// get SeachData objects from IDS
	results, err := getSearchData(db, ids)
	if err != nil {

	}

	return results, nil
}

func getSearchData(db *bolt.DB, ids []string) ([]SearchData, error) {
	results := []SearchData{}
	var data []byte

	// tx
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lookup"))
		for _, ID := range ids {
			data = b.Get([]byte(ID))
			sd, err := decodeSearchData(data)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed: %v", err)
			}
			results = append(results, *sd)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return []SearchData{}, fmt.Errorf("getSearchData failed: %v", err)
	}

	return results, nil
}

// persist IndexData object to disk
func saveIndexData(index *IndexData) error {
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
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
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
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
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
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
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
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

// encode SearchData to protobuf
func encodeSearchData(sd *SearchData) ([]byte, error) {
	entry := &protobuf.SearchResult{
		ID:     sd.ID,
		Name:   sd.Name,
		City:   sd.City,
		State:  sd.State,
		Bucket: sd.Bucket,
		Years:  sd.Years,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeSearchEntry failed: %v", err)
	}
	return data, nil
}

// decode protobuf to SearchData
func decodeSearchData(input []byte) (*SearchData, error) {
	sr := &protobuf.SearchResult{}
	err := proto.Unmarshal(input, sr)
	if err != nil {
		fmt.Println(err)
		return &SearchData{}, fmt.Errorf("decodeSearchData failed: %v", err)
	}
	sd := &SearchData{
		ID:     sr.GetID(),
		Name:   sr.GetName(),
		City:   sr.GetCity(),
		State:  sr.GetState(),
		Bucket: sr.GetBucket(),
		Years:  sr.GetYears(),
	}
	return sd, nil
}

func encodeResultsList(r resultList) ([]byte, error) {
	entry := &protobuf.ResultList{
		IDs: r.Results,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeSearchEntry failed: %v", err)
	}
	return data, nil
}

func decodeResultsList(input []byte) ([]string, error) {
	r := &protobuf.ResultList{}
	err := proto.Unmarshal(input, r)
	if err != nil {
		fmt.Println(err)
		return []string{}, fmt.Errorf("decodeSearchEntry failed: %v", err)
	}
	results := r.GetIDs()
	return results, nil

}

// add new lookupPairs to existing data
func mergeLookup(new, prev lookupPairs) {
	for k, v := range new {
		prev[k] = v
	}
}

func mergeIDs(new, prev []string) []string {
	update := []string{}
	set := make(map[string]bool)
	for _, ID := range new {
		set[ID] = true
	}
	for _, ID := range prev {
		set[ID] = true
	}
	for ID := range set {
		update = append(update, ID)
	}
	return update
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
