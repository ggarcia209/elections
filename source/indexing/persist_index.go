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

// list wrapped in struct for protobuf encoding
type resultList struct {
	Results []string // new/updated ID references for given term encoded as Big Endian uint64s
}

type shardMap map[string]*shardRanges // schema: term: shards: range (min, max)

// data for shards for each term in shardMap
type shardRanges struct {
	Term   string                // original search term
	Shards float32               // number of shards for given term
	Ranges map[string]rangeTuple // min, max value for each shard (ID = string(shardIndex))
}

// list wrapped in struct for protobuf encoding
type rangeTuple struct {
	Range []string // (min, max)
}

// save new Index entries
func saveIndex(id *IndexData, index indexMap, lookup lookupPairs) (int, int, error) {
	var t int64 // new terms
	var u int64 // updated terms
	var n int64 // new IDs/Search Data pairs wrote
	// uints := make(map[string]map[string][]string) // store big endian encoded IDs for each shard
	wrote := make(map[string]bool)
	shards := make(shardMap) // shardMap buffer object
	maxSize := 1250          // # of IDs (max size @ 4b/ID)
	ns := "!"                // indicates max value not set for shard (shard incomplete)
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

				if id.Shards[term] != nil && id.Shards[term].Shards != 0 { // shards exist for term
					fmt.Printf("existing shards found for '%s': %.0f\n", term, id.Shards[term].Shards)
					totalPrev := []string{}  // aggregate total from shards
					key := term              // shardID
					prTotal := 0             // previous # of IDs in shards[i:]
					if shards[term] == nil { // re-initialize buffer object
						shards[term] = &shardRanges{
							Term:   term,
							Shards: 0,
							Ranges: make(map[string]rangeTuple),
						}
					}

					// get aggregate data from existing shards
					for i := 0; i < int(id.Shards[term].Shards+1.0); i++ {
						if i > 0 {
							key = term + "." + strconv.Itoa(i)
						}
						max := ns
						if id.Shards[term].Ranges[key].Range[1] != ns { // shard is full and max value set
							max = id.Shards[term].Ranges[key].Range[1]
						}

						// compare max value of full shards, skip if below threshold
						if max != ns && comp > max {
							// if term out of range && partition is full - skip
							// item indexes of preceeding shards remain unchanged
							fmt.Println("skipping partition: ", si)
							si++ // increment for every preceeding shard
							continue
						}

						prevData := b.Get([]byte(key))
						prev, err := decodeResultsList(prevData)
						if err != nil {
							fmt.Println(err)
							return fmt.Errorf("tx failed failed: %v", err)
						}

						// DEBUGGING
						if len(prev) == 0 {
							fmt.Println("!!! EMPTY SHARD")
						}
						if i < int(id.Shards[term].Shards) && len(prev) < maxSize {
							fmt.Println(" !!! INCOMPLETE SHARD")
						}
						/* repeats := make(map[string]bool)
						for _, ID := range prev {
							repeats[ID] = true
						}
						if len(repeats) < len(prev) {
							fmt.Println("!!! INTERSHARD REPEATS FOUND")
							fmt.Println("original len: ", len(prev))
							fmt.Println("intershard set", len(repeats))
						} */

						// add shard to aggregate total for re-ordering
						totalPrev = append(totalPrev, prev...)
						prTotal += len(prev)
					}

					// debugging
					/* repeats := make(map[string]bool)
					r := 0
					for _, ID := range totalPrev {
						if repeats[ID] == true {
							r++
						}
						repeats[ID] = true
					}
					if r != 0 {
						fmt.Println("!!! REPEATS: ", r)
					} */

					// create set of new/existing IDs and update index
					update := mergeIDs(ids, totalPrev)
					terms[term] = update
					shards[term].Shards = float32(si) // new shards created at this index

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
					if shards[term] == nil { // previous IDs = 0 shards; + newIDs = 1+ shards
						shards[term] = &shardRanges{
							Term:   term,
							Shards: 0,
							Ranges: make(map[string]rangeTuple),
						}
					}
					if shards[term].Shards > 0 { // new shard created at index > 0
						shard = shard + "." + strconv.Itoa(int(shards[term].Shards))
						delete(terms, term) // delete original term from map to prevent overwrite of index 0
					}
					j := 0
					min, max := "", ns // min, max values of every shard ("!" indicates max value not set/shard incomplete)
					for i, ID := range orig {
						if i == 0 { // set min value of first shard (shards[term])
							min = ID
							shards[term].Ranges[shard] = rangeTuple{[]string{min, max}}
							fmt.Println("new list shard ", shard)
						}
						// add IDs to each shard; incremement shardID for every maxSize items
						if i == maxSize*(j+1) {
							j++
							// update data for current shard
							max = orig[i-1] // set new max for current shard
							shards[term].Ranges[shard].Range[1] = max
							// create new shard
							shards[term].Shards++
							shard = term + "." + strconv.Itoa(int(shards[term].Shards)) // new shard
							min = orig[i]                                               // set new min value
							shards[term].Ranges[shard] = rangeTuple{[]string{min, ns}}
							fmt.Println("new list shard ", shard)
						}
						newShards[shard] = append(newShards[shard], ID)
					}
					fmt.Println()
				} else if shards[term] != nil && shards[term].Shards > 0 { // shard index > 0; total items < maxSize (partial shard added to existing shard)
					shard := term + "." + strconv.Itoa(int(shards[term].Shards)) // find current shard index
					min := terms[term][0]
					shards[term].Ranges[shard] = rangeTuple{[]string{min, ns}}
					newShards[shard] = terms[term]
					delete(terms, term) // delete shard from original to prevent overwrite of index 0
				} else {
					// new entries
					// no previous shards, new IDs < maxLength
					// no change
				}
			} // end shard creation logic

			// merge shards with index
			for shard, list := range newShards {
				terms[shard] = list
				delete(newShards, shard)
			}
			// end sharding logic

			// encode ID lists and save to disk
			for shard, ids := range terms {
				r := resultList{Results: ids}
				data, err := encodeResultsList(r)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("tx failed failed: %v", err)
				}
				if err := b.Put([]byte(shard), data); err != nil { // serialize k,v
					return fmt.Errorf("tx failed failed: %v", err)
				}
				delete(terms, shard)
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
	fmt.Printf("Writing %d lookup objects...\n", len(lookup))
	k := 0
	sds := []*SearchData{}
	for _, sd := range lookup {
		limit := 90000 // limit to 90000 items per tx (90% of recommended max per BoltDB docs)
		if k == limit {
			_, err := batchWriteLookup(db, sds, n, wrote)
			if err != nil {
				fmt.Println(err)
				return 0, 0, fmt.Errorf("saveIndex failed: %v", err)
			}
			sds = []*SearchData{} // reset after batch write
			k = 0                 // reset
		}
		sds = append(sds, sd)
		k++
	}
	// write remainder of SearchData objects
	_, err = batchWriteLookup(db, sds, n, wrote)
	if err != nil {
		fmt.Println(err)
		return 0, 0, fmt.Errorf("saveIndex failed: %v", err)
	}

	// merge shard counts and ranges with id.Shards
	for term, sr := range shards {
		if sr.Shards == 0 {
			continue
		}
		if id.Shards[term] == nil {
			id.Shards[term] = sr
			continue
		}
		for shard, tuple := range sr.Ranges {
			id.Shards[term].Ranges[shard] = tuple
		}
		id.Shards[term].Shards = sr.Shards
	}

	fmt.Println("index saved!")
	fmt.Println("new terms wrote: ", t)
	fmt.Println("existing terms: ", u)
	fmt.Println("IDs recorded/updated: ", len(lookup))
	fmt.Println()
	return int(n), int(t), nil
}

func getPreviousShards(id *IndexData, b *bolt.Bucket, term string, ids []string, terms map[string][]string, shards shardMap) (int, int, error) {
	// begin aggregate existing data logic
	si := 0 // shard index
	t := 0  // new terms recorded in func call
	u := 0  // existing terms updated in func call
	sort.Strings(ids)
	ns := "!"      // indicates max value for shard not set - shard is incomplete
	comp := ids[0] // compare to greatest value in each shard to find corresponding index

	if id.Shards[term] != nil && id.Shards[term].Shards != 0 { // shards exist for term
		fmt.Printf("existing shards found for '%s': %.0f\n", term, id.Shards[term].Shards)
		totalPrev := []string{}  // aggregate total from shards
		key := term              // shardID
		prTotal := 0             // previous # of IDs in shards[i:]
		if shards[term] == nil { // re-initialize buffer object
			shards[term] = &shardRanges{
				Term:   term,
				Shards: 0,
				Ranges: make(map[string]rangeTuple),
			}
		}

		// get aggregate data from existing shards
		for i := 0; i < int(id.Shards[term].Shards+1.0); i++ {
			if i > 0 {
				key = term + "." + strconv.Itoa(i)
			}
			max := ns
			if id.Shards[term].Ranges[key].Range[1] != ns { // shard is full and max value set
				max = id.Shards[term].Ranges[key].Range[1]
			}

			// compare max value of full shards, skip if below threshold
			if max != ns && comp > max {
				// if term out of range && partition is full - skip
				// item indexes of preceeding shards remain unchanged
				fmt.Println("skipping partition: ", si)
				si++ // increment for every preceeding shard
				continue
			}

			prevData := b.Get([]byte(key))
			prev, err := decodeResultsList(prevData)
			if err != nil {
				fmt.Println(err)
				return 0, 0, fmt.Errorf("getPreviousShards failed failed: %v", err)
			}

			// add shard to aggregate total for re-ordering
			totalPrev = append(totalPrev, prev...)
			prTotal += len(prev)
		}

		// create set of new/existing IDs and update index
		update := mergeIDs(ids, totalPrev)
		terms[term] = update
		shards[term].Shards = float32(si) // new shards created at this index

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
				return 0, 0, fmt.Errorf("getPreviousShards failed failed: %v", err)
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
	return t, u, nil
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

// get SearchData objs from ID references for given term
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

// get SearchData objects form list of IDs
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
func GetPartitionMap() (map[string]bool, error) {
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

// encode new IDs with big endian encoding and merge with existing IDs
func mergeIDs(new, prev []string) []string {
	// fmt.Printf("mergeIDs start: len new: %d, len prev: %d\n", len(new), len(prev))
	update := []string{}
	set := make(map[string]bool)
	for _, ID := range prev {
		set[ID] = true
	}
	for _, ID := range new {
		set[ID] = true
	}
	for ID := range set {
		update = append(update, ID)
	}

	sort.Slice(update, func(i, j int) bool { return update[i] < update[j] })

	return update
}

// encode IndexData to protobuf
func encodeIndexData(id *IndexData) ([]byte, error) {
	entry := &protobuf.IndexData{
		TermsSize:      float32(id.TermsSize),
		LookupSize:     float32(id.LookupSize),
		Completed:      id.Completed,
		YearsCompleted: id.YearsCompleted,
	}

	// encode shardMap nested objects
	shards := make(map[string]*protobuf.ShardRanges)
	for term, sr := range id.Shards {
		srPb := &protobuf.ShardRanges{
			Term:   sr.Term,
			Shards: sr.Shards,
			Ranges: make(map[string]*protobuf.Range),
		}
		for s, r := range sr.Ranges {
			rangePb := &protobuf.Range{
				Range: []string{},
			}
			for _, r := range r.Range {
				rangePb.Range = append(rangePb.Range, r)
			}
			srPb.Ranges[s] = rangePb
		}
		shards[term] = srPb
	}
	entry.Shards = shards

	// encode timestamp
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
		YearsCompleted: data.GetYearsCompleted(),
	}

	// decode shards
	shards := make(shardMap)
	shardsPb := data.GetShards()
	for shard, srPb := range shardsPb {
		sr := &shardRanges{
			Term:   srPb.GetTerm(),
			Shards: srPb.GetShards(),
		}
		rs := srPb.GetRanges()
		ranges := make(map[string]rangeTuple)
		for key, r := range rs {
			rt := rangeTuple{r.GetRange()}
			ranges[key] = rt
		}
		sr.Ranges = ranges
		shards[shard] = sr
	}
	id.Shards = shards

	// decode timestampe
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
