package indexing

import (
	"fmt"
	"sort"
	"strconv"
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
func saveIndex(id *IndexData, index indexMap, lookup lookupPairs) (int, int, error) {
	var t int64 // new terms
	var u int64 // updated terms
	var n int64 // new IDs/Search Data pairs wrote
	wrote := make(map[string]bool)
	shards := make(map[string]float32) // number of shards for sharded term - keys will appears as: 'term', 'term.1', 'term.2', etc...
	maxSize := 90000                   // # of IDs (max size @ 4b/ID)
	fmt.Println("save index - writing objects to db/search_index.db")

	// open/create bucket in db/offline_db.db
	// put protobuf item and use donor.ID as key
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return 0, 0, fmt.Errorf("saveIndex failed: %v", err)
	}

	// persist inverted index
	// add/update each term for each partition
	for prt, terms := range index {
		fmt.Printf("writing partition '%s' - %d items...\n", prt, len(terms))
		newShards := make(map[string][]string)

		// 1st tx set
		if err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(prt))
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("tx failed failed: %v", err)
			}

			// begin sharding logic
			// merge new list of IDs with existing data for each term
			// create/re-create shards if necessary
			// existing shards are recreated from aggregate total of existingShards[i:]
			for term, ids := range terms {
				// begin aggregate existing data logic
				si := 0 // shard index
				sort.Strings(ids)
				comp := ids[0] // compare to greatest value in each shard to find corresponding index

				if id.Shards[term] != 0 { // shards exist for term
					fmt.Printf("existing shards found for '%s': %.0f\n", term, id.Shards[term])
					totalPrev := []string{} // aggregate total from shards
					key := term             // shardID
					prTotal := 0            // previous # of IDs in shards[i:]

					// get aggregate data from existing shards
					for i := 0; i < int(id.Shards[term]+1.0); i++ {
						if i > 0 {
							key = term + "." + strconv.Itoa(i)
							fmt.Println("existing shard found: ", key)
						}
						// get previous data & check each partition for appropriate index of new key
						// (store key in byte-sorted/alphabetical order)
						prevData := b.Get([]byte(key))
						prev, err := decodeResultsList(prevData)
						if err != nil {
							fmt.Println(err)
							return fmt.Errorf("tx failed failed: %v", err)
						}
						max := prev[len(prev)-1]
						fmt.Printf("\tcomp: %s\tmax: %s\n", comp, max)
						if comp > max && len(prev) >= maxSize {
							// if term out of range && partition is full - skip
							// item indexes of preceeding shards remain unchanged
							fmt.Println("skipping partition: ", si)
							si++ // increment for every preceeding shard
							continue
						}
						// add shard to aggregate total for re-ordering
						totalPrev = append(totalPrev, prev...)
						prTotal += len(prev)
						fmt.Printf("added %d IDs for '%s'\n", len(prev), term)
					}

					// create set of new/existing IDs and update index
					update := mergeIDs(ids, totalPrev)
					terms[term] = update
					shards[term] = float32(si) // new shards created at this index
					fmt.Println("previous total: ", prTotal)
					fmt.Println("new total: ", len(update))
					fmt.Printf("creating new shards for term '%s' at index %d\n", term, si)
					u++
				} else { // no existing shards
					// get previous data
					prevData := b.Get([]byte(term))
					if prevData != nil { // existing term entry
						prev, err := decodeResultsList(prevData)
						if err != nil {
							fmt.Println(err)
							return fmt.Errorf("tx failed failed: %v", err)
						}
						// create ID set and update index
						update := mergeIDs(ids, prev)
						terms[term] = update
						u++
					} else { // new term entry, no shards; sort IDs
						empty := []string{}
						sorted := mergeIDs(ids, empty)
						terms[term] = sorted
						t++
					}
				} // end aggregate existing data logic

				// create new shards if len(ids) > maxSize
				// recreate shards at index if new ID's added to existing shard
				l := len(terms[term])
				if l > maxSize {
					fmt.Printf("maxSize exceeded for term '%s' (%d)\n", term, l)
					orig := terms[term]
					shard := term
					if shards[term] > 0 {
						shard = shard + "." + strconv.Itoa(int(shards[term]))
						delete(terms, term) // delete unsharded term from map if not overwriting in current function call
					}
					j := 0
					for i, ID := range orig {
						// add IDs to each shard; incremement shardID for every maxSize items
						if i == maxSize*(j+1) {
							j++
							shards[term]++
							shard = term + "." + strconv.Itoa(int(shards[term])) // new shard
							fmt.Println("saveIndex: new list shard ", shard)
						}
						newShards[shard] = append(newShards[shard], ID)
					}
				} else if shards[term] > 0 { // shard index > 0; total items < maxSize
					shard := term + "." + strconv.Itoa(int(shards[term]))
					newShards[shard] = terms[term]
					delete(terms, term)
				}
			} // end shard creation logic

			// merge shards with index
			for shard, list := range newShards {
				index[prt][shard] = list
			}
			// end sharding logic

			// encode ID lists and save to disk
			for term, ids := range terms {
				r := resultList{Results: ids}
				data, err := encodeResultsList(r)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("tx failed failed: %v", err)
				}
				if err := b.Put([]byte(term), data); err != nil { // serialize k,v
					return fmt.Errorf("tx failed failed: %v", err)
				}
				// delete(terms, term) // delete from memory once persisted; drain pressure on memory
			}
			return nil
		}); err != nil {
			fmt.Println(err)
			return 0, 0, fmt.Errorf("saveIndex failed: %v", err)
		}
	}
	// end 1st tx set

	// write lookup objects to disk (90,000 max per iteration)
	// 2nd transaction set
	k := 0
	sds := []*SearchData{}
	for _, sd := range lookup {
		limit := 90000 // limit to 90000 items per tx (90% of recommended max per BoltDB docs)
		if k == limit {
			sn, err := batchWriteLookup(db, sds, n, wrote)
			if err != nil {
				fmt.Println(err)
				return 0, 0, fmt.Errorf("saveIndex failed: %v", err)
			}
			sds = []*SearchData{} // reset after batch write
			k = 0                 // reset
			n += sn
		}
		sds = append(sds, sd)
		k++
	}
	// write remainder of SearchData objects
	sn, err := batchWriteLookup(db, sds, n, wrote)
	if err != nil {
		fmt.Println(err)
		return 0, 0, fmt.Errorf("saveIndex failed: %v", err)
	}
	n += sn

	// merge shard counts with id.Shards
	for shard, count := range shards {
		if id.Shards == nil {
			id.Shards = make(map[string]float32)
		}
		id.Shards[shard] = count
	}

	fmt.Println("index saved!")
	fmt.Println("new terms wrote: ", t)
	fmt.Println("existing terms: ", u)
	fmt.Println("IDs recorded/updated: ", n)
	fmt.Println()
	return int(n), int(t), nil
}

// batch writes list of *SearchData objects to disk in single transaction
func batchWriteLookup(db *bolt.DB, sds []*SearchData, n int64, wrote map[string]bool) (int64, error) {
	// tx
	if err := db.Update(func(tx *bolt.Tx) error {
		for _, sd := range sds {
			// check if wrote to disk previously
			if wrote[sd.ID] == false { // not wrote to disk in parent function call
				luB, err := tx.CreateBucketIfNotExists([]byte("lookup"))
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("tx failed failed: %v", err)
				}
				prevData := luB.Get([]byte(sd.ID))
				if prevData != nil {
					prev, err := decodeSearchData(prevData)
					if err != nil {
						fmt.Println()
						return fmt.Errorf("tx failed: %v", err)
					}
					// reconcile "Unknown" records (entity not registered in current year)
					sd.Years = append(sd.Years, prev.Years...)
					if sd.Name == "Unknown" {
						sd.Name = prev.Name
						sd.City = prev.City
						sd.State = prev.State
						sd.Employer = prev.Employer
						sd.Bucket = prev.Bucket
					}
				} // else { n++ }
				data, err := encodeSearchData(sd)
				if err != nil {
					fmt.Println()
					return fmt.Errorf("tx failed: %v", err)
				}
				if err := luB.Put([]byte(sd.ID), data); err != nil { // serialize k,v
					return fmt.Errorf("tx failed failed: %v", err)
				}
				wrote[sd.ID] = true
			}
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return n, fmt.Errorf("batchWriteLookup failed: %v", err)
	}
	return n, nil
}

// get lookup pairs from disk for given term
func getSearchEntry(term string) ([]SearchData, error) {
	// retreive lookupPairs from disk
	prt := getPartition(term)
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return []SearchData{}, fmt.Errorf("getSearchEntry failed: %v", err)
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

	if len(ids) > 200 {
		return []SearchData{}, fmt.Errorf("MAX_LENGTH")
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
			if sd.Bucket != "" { // not nil item
				results = append(results, *sd)
			}
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
		ID:       sd.ID,
		Name:     sd.Name,
		City:     sd.City,
		State:    sd.State,
		Bucket:   sd.Bucket,
		Employer: sd.Employer,
		Years:    sd.Years,
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
		ID:       sr.GetID(),
		Name:     sr.GetName(),
		City:     sr.GetCity(),
		State:    sr.GetState(),
		Employer: sr.GetEmployer(),
		Bucket:   sr.GetBucket(),
		Years:    sr.GetYears(),
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
	sort.Strings(update)
	return update
}

// encode IndexData to protobuf
func encodeIndexData(id *IndexData) ([]byte, error) {
	entry := &protobuf.IndexData{
		TermsSize:      float32(id.TermsSize),
		LookupSize:     float32(id.LookupSize),
		Completed:      id.Completed,
		Shards:         id.Shards,
		YearsCompleted: id.YearsCompleted,
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
		TermsSize:      int(data.GetTermsSize()),
		LookupSize:     int(data.GetLookupSize()),
		Completed:      data.GetCompleted(),
		Shards:         data.GetShards(),
		YearsCompleted: data.GetYearsCompleted(),
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
