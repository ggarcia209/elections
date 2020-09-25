package indexing

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elections/source/dynamo"

	"github.com/boltdb/bolt"
)

// Query uses text from user input to look up corresponding results
type Query struct {
	Text      string
	UserID    string
	TimeStamp time.Time
}

// Data is used to find the DocID's common to all terms in query
type Data struct {
	Key   string
	Value []string // sorted by ID value
	Len   int
}

// Entry represents a k/v pair in a sorted map
type Entry struct {
	ID   string
	Data SearchData
}

// DynamoSearchEntry is used to retrieve a list of IDs from the Index table
type DynamoSearchEntry struct {
	Partition string
	IDs       []string
}

// Entries represents a sorted map
type Entries []Entry

func (s Entries) Len() int           { return len(s) }
func (s Entries) Less(i, j int) bool { return s[i].ID < s[j].ID }
func (s Entries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// shardEntry is used to sort terms by number of shards for the given term
type shardEntry struct {
	ID     string
	Shards float32
}

// Entries represents a sorted map
type shardEntries []shardEntry

func (s shardEntries) Len() int           { return len(s) }
func (s shardEntries) Less(i, j int) bool { return s[i].Shards < s[j].Shards }
func (s shardEntries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// prtEntry represents a k/v pair for a sorted partition map
type prtEntry struct {
	Prt string
	B   bool
}

// prtEntries represents a partition map to be sorted
type prtEntries []prtEntry

func (s prtEntries) Len() int           { return len(s) }
func (s prtEntries) Less(i, j int) bool { return s[i].Prt < s[j].Prt }
func (s prtEntries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// CreateQuery returns a Query object with the given data
func CreateQuery(text, UID string) Query {
	q := Query{
		Text:      text,
		UserID:    UID,
		TimeStamp: time.Now(),
	}
	return q
}

// GetResults returns search results to user from user query.
func GetResults(q Query) ([]string, error) {
	if q.Text == "" {
		return []string{}, nil
	}

	// get metadata, use to check # of shards for each term
	indexData, err := getIndexData()
	if err != nil {
		fmt.Println(err)
		return []string{}, fmt.Errorf("GetResults failed: %v", err)
	}

	// normalize search terms input
	terms := formatTerms(strings.Split(q.Text, " "))

	// get IDs for multiple terms
	resMap, x, err := getRefs(terms, indexData)
	if err != nil {
		fmt.Println(err)
		return []string{}, fmt.Errorf("GetResults failed: %v", err)
	}

	// results not found for 1+ more terms
	if x > 0 {
		return []string{}, nil
	}

	// get IDs for single term
	if len(terms) == 1 {
		common := []string{}
		for _, ss := range resMap {
			for i, ID := range ss {
				common = append(common, ID)
				if i > 200 {
					fmt.Println("GetResults failed: MAX_LENGTH exceeded")
					return []string{}, fmt.Errorf("MAX_LENGTH")
				}
			}
		}
		return common, nil
	}

	// Sort lists by smallest to largest
	sorted := sortMap(resMap)

	// Compare and find all common IDs in all terms
	// Start by finding the common IDs in the 2 smallest lists
	// then compare the next list to the previous comparison's intersection
	s0, s1 := sorted[0].Value, sorted[1].Value
	common := intersection(s0, s1)
	for i := 2; i < len(sorted); i++ {
		if len(common) == 0 {
			return common, nil
		}
		common = intersection(common, sorted[i].Value)
	}
	fmt.Println("len(common): ", len(common))
	if len(common) > 200 {
		fmt.Println("GetResults failed: MAX_LENGTH exceeded")
		return []string{}, fmt.Errorf("MAX_LENGTH")
	}

	return common, nil
}

// GetResults returns search results to user from user query.
func GetResultsFromShards(q Query) ([]string, error) {
	stTotal := time.Now()
	if q.Text == "" {
		return []string{}, nil
	}

	// get metadata, use to lookup # of shards for each term
	id, err := getIndexData()
	if err != nil {
		fmt.Println(err)
		return []string{}, fmt.Errorf("GetResults failed: %v", err)
	}

	terms := formatTerms(strings.Split(q.Text, " ")) // normalize search terms input
	common := []string{}                             // aggegate total of intersections for every shard
	maxResultsSize := 500                            // max number of SearchResults returned before throwing MAX_LENGTH error

	// get IDs for single term
	if len(terms) == 1 {
		t := terms[0]
		if id.Shards[t] > 0 {
			fmt.Println("GetResults failed: MAX_LENGTH exceeded")
			return []string{}, fmt.Errorf("MAX_LENGTH")
		}
		t = strings.TrimSpace(t)
		ids, err := getShard(t)
		if err != nil {
			fmt.Println(err)
			return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
		}
		// add results to common aggregate total; throw error if maxResultsSize reached
		for i, ID := range ids {
			if i >= maxResultsSize {
				fmt.Println("GetResults failed: MAX_LENGTH exceeded")
				return []string{}, fmt.Errorf("MAX_LENGTH")
			}
			common = append(common, ID)
		}
		fmt.Println("finish - ", time.Since(stTotal))
		return common, nil
	}

	// sort terms by # of shards, least to greatest
	st := time.Now()
	termsSrt := shardEntries{}
	for _, t := range terms {
		se := shardEntry{ID: t, Shards: id.Shards[t]}
		termsSrt = append(termsSrt, se)
	}
	sort.Sort(termsSrt)
	fmt.Println("terms sorted: ", time.Since(st))

	// find intersections of two smallest records - terms[0] & terms[1] (t0, t1)
	st = time.Now()
	i0, i1 := termsSrt[0].ID, termsSrt[1].ID
	k0, k1 := i0, i1
	// compare each shard (s0) in t0 to each shard in t1 (s1) within min/max range
	for x := 0; x < int((id.Shards[i0] + 1)); x++ {
		if x > 0 {
			k0 = i0 + "." + strconv.Itoa(x)
		}
		s0, err := getShard(k0)
		if err != nil {
			fmt.Println(err)
			return []string{}, fmt.Errorf("getRefs failed: %v", err)
		}
		min0, max0 := s0[0], s0[len(s0)-1] // min, max value of each shard in t1 compared to these values

		// get each shard in t1 (s1), compare to s0 and find intersection
		for y := 0; y < int((id.Shards[i1] + 1)); y++ {
			if y > 0 {
				k1 = i1 + "." + strconv.Itoa(y)
			}
			s1, err := getShard(k1)
			if err != nil {
				fmt.Println(err)
				return []string{}, fmt.Errorf("getRefs failed: %v", err)
			}

			// skip shards < min0 value, stop at shard > max0 value
			min1, max1 := s1[0], s1[len(s1)-1]
			if min0 > max1 {
				continue // skip
			}
			if max0 < min1 {
				break // stop
			}

			// find intersection of s0, s1, add to common aggregate total
			intsec := intersection(s0, s1)
			common = append(common, intsec...)
		}
	}
	fmt.Println("s0, s1 comparision: ", time.Since(st))

	// find intersections of each remaining term
	for t := 2; t < len(termsSrt); t++ {
		if len(common) == 0 {
			return common, nil
		}
		i2 := termsSrt[t].ID // next term in list
		k2 := i2
		buffer := common    // use to find values in common within range of current s1 shard
		common = []string{} // reset for each term
		errChan := make(chan error)
		var wg sync.WaitGroup
		st = time.Now()

		// get each shard in t1 (s1), compare to s0 and find intersection
		for z := 0; z < int((id.Shards[i2] + 1)); z++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				if z > 0 {
					k2 = i2 + "." + strconv.Itoa(z)
				}
				st = time.Now()
				s2, err := getShard(k2)
				if err != nil {
					fmt.Println(err)
					errChan <- err
					return
				}
				fmt.Printf("%s - getShard - %v\n", k2, time.Since(st))

				s0 := []string{} // shard created from common IDs within range of s2
				min1, max1 := s2[0], s2[len(s2)-1]

				st = time.Now()
				// create comparision shard from buffer IDs within range
				for _, v := range buffer {
					if v >= max1 { // no more common values exist
						break // stop
					}
					if v >= min1 {
						s0 = append(s0, v)
					}
				}
				fmt.Printf("%s - create buffer - %v\n", k2, time.Since(st))

				// find intersection of s0, s2, add to common aggregate total
				st = time.Now()
				intsec := intersection(s0, s2)
				mu.Lock()
				common = append(common, intsec...)
				errChan <- nil
				mu.Unlock()
				fmt.Printf("%s - intersection - %v\n", k2, time.Since(st))
				return

			}(&wg)
			err := <-errChan
			if err != nil {
				return []string{}, fmt.Errorf("getResultsFromShards failed: %v", err)
			}
		}
		// wait for goroutines to finish and sort results
		wg.Wait()
		sort.Strings(common)
		fmt.Println("s2 comparison: ", time.Since(st))
	}

	fmt.Println("len(common): ", len(common))
	if len(common) > maxResultsSize {
		fmt.Println("GetResults failed: MAX_LENGTH exceeded")
		return []string{}, fmt.Errorf("MAX_LENGTH")
	}

	fmt.Println("finish - ", time.Since(stTotal))
	return common, nil
}

// LookupSearchDataFromCache retreives SearchData objects from in memory cache
func LookupSearchDataFromCache(ids []string, cache map[string]SearchData) ([]string, []SearchData) {
	sds := []SearchData{}
	nilIDs := []string{}

	for _, ID := range ids {
		sd := cache[ID]
		if sd.ID != "" {
			sds = append(sds, sd)
		} else {
			nilIDs = append(nilIDs, ID)
		}
	}

	return nilIDs, sds
}

// LookupSearchData Retreives corresponding SearchData obj for ID from disk
func LookupSearchData(ids []string) ([]SearchData, error) {
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return []SearchData{}, fmt.Errorf("LookupSearchData failed: %v", err)
	}

	results, err := getSearchData(db, ids)
	if err != nil {
		fmt.Println(err)
		return []SearchData{}, fmt.Errorf("LookupSearchData failed: %v", err)
	}

	return results, nil
}

// GetResultsFromDynamo returns list of common IDs from DynamoDB
/* func GetResultsFromDynamo(db *dynamo.DbInfo, q Query) ([]string, error) {
	if q.Text == "" {
		return []string{}, nil
	}
	resMap := make(map[string][]string)
	shardsSrt := shardEntries{}
	indexData, err := getIndexData()
	if err != nil {
		fmt.Println(err)
		return []string{}, fmt.Errorf("GetResults failed: %v", err)
	}

	// find corresponding SearchData objects
	terms := formatTerms(strings.Split(q.Text, " "))

	// get IDs for single term
	if len(terms) == 1 {
		if indexData.Shards[terms[0]] > 0 {
			fmt.Println("GetResults failed: MAX_LENGTH exceeded")
			return []string{}, fmt.Errorf("MAX_LENGTH")
		}
		obj

		return common, nil
	}

	for _, term := range terms {
		// get IDs for multiple terms
		result, x, err := getRefsFromDynamo(terms, indexData)
		if err != nil {
			fmt.Println(err)
			return []string{}, fmt.Errorf("GetResults failed: %v", err)
		}

		// results not found for 1+ more terms
		if x > 0 {
			return []string{}, nil
		}
		resultMap[term] = result
	}

	// Sort lists by smallest to largest
	// get terms with least # of shards first
	for _, term := range terms {
		se := shardEntry{
			ID: term,
			Shards: indexData.Shards[term],
		}
		shardsSrt = append(shardsSrt, se)
	}
	sort.Sort(shardsSrt)



	// Compare and find all common IDs in all terms
	// Start by finding the common IDs in the 2 smallest lists
	// then compare the next list to the previous comparison's intersection
	s0, s1 := sorted[0].Value, sorted[1].Value
	common := intersection(s0, s1)
	for i := 2; i < len(sorted); i++ {
		if len(common) == 0 {
			return common, nil
		}
		common = intersection(common, sorted[i].Value)
	}
	fmt.Println("len(common): ", len(common))
	if len(common) > 200 {
		fmt.Println("GetResults failed: MAX_LENGTH exceeded")
		return []string{}, fmt.Errorf("MAX_LENGTH")
	}

	return common, nil
} */

// ConsolidateSearchData consolidates SearchData lists from cache and disk in their original order.
func ConsolidateSearchData(origIDs []string, frmCache, frmDisk []SearchData) []SearchData {
	sds := make(map[string]SearchData)
	agg := []SearchData{}
	for _, sd := range frmCache {
		sds[sd.ID] = sd
	}
	for _, sd := range frmDisk {
		sds[sd.ID] = sd
	}
	for _, ID := range origIDs {
		agg = append(agg, sds[ID])
	}
	return agg
}

// ViewIndex displays the index
func ViewIndex() error {
	ct := 0

	// get partition map
	partitions, err := getPartitionMap()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ViewIndex failed: %v", err)
	}
	if len(partitions) == 0 {
		fmt.Println("*** Index does not exist! ***")
		return nil
	}

	// sort partition map
	es := prtEntries{}
	for k, v := range partitions {
		e := prtEntry{
			Prt: k,
			B:   v,
		}
		es = append(es, e)
	}
	sort.Sort(&es)

	// open db
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ViewIndex failed: %v", err)
	}
	defer db.Close()

	// iterate through each key in each partiton in alphabetical order
	err = db.View(func(tx *bolt.Tx) error {
		for _, prt := range es {
			b := tx.Bucket([]byte(prt.Prt))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				d, err := decodeResultsList(v)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("tx failed: %v", err)
				}
				fmt.Printf("Partition: %s\tTerm: %s\tReferences: %d\n", prt.Prt, k, len(d))
				ct++
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ViewIndex failed: %v", err)
	}

	fmt.Println("\nTotal entries: ", ct)
	return nil
}

// GetIndexData retrieves the IndexData search index metadata object from disk
func GetIndexData() (*IndexData, error) {
	id, err := getIndexData()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("ViewIndexData failed: %v", err)
	}
	return id, nil
}

// getRefs finds the references for each term in query
func getRefs(q []string, id *IndexData) (map[string][]string, int, error) {
	var resultMap = make(map[string][]string)
	var noIDs bool
	x := 0 // value > 0 indicates 1+ terms in query returned no matching IDs; no results for query

	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return nil, 0, fmt.Errorf("getRefs failed: %v", err)
	}
	defer db.Close()

	// Get index list for each term in query - use map
	for _, v := range q {
		if filter(v) {
			continue
		}
		prt := getPartition(v)
		key := ""
		result := []string{}
		// get ids for each term and any shards
		for i := 0; i < int(id.Shards[v]+1.0); i++ {
			err := db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(prt))
				v = strings.TrimSpace(v)
				key = v
				if i > 0 {
					key = v + "." + strconv.Itoa(i)
				}
				data := b.Get([]byte(key))
				ids, err := decodeResultsList(data)
				if err != nil {
					return fmt.Errorf("tx failed: %v", err)
				}
				if len(ids) == 0 {
					fmt.Println("no IDs found for term: ", v)
					noIDs = true
					return nil
				}
				fmt.Printf("ids found for '%s': %d\n", v, len(ids))
				result = append(result, ids...)
				return nil
			})
			if err != nil {
				return nil, 0, fmt.Errorf("getRefs failed: %s", err)
			}
		}
		if noIDs {
			x++
			continue
		}
		resultMap[v] = result
	}
	return resultMap, x, nil
}

// get ids in single shard
func getShard(id string) ([]string, error) {
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("getShard failed: %v", err)
	}
	defer db.Close()

	result := []string{}
	prt := getPartition(id)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(prt))
		id = strings.TrimSpace(id)
		data := b.Get([]byte(id))
		ids, err := decodeResultsList(data)
		if err != nil {
			return fmt.Errorf("tx failed: %v", err)
		}
		result = ids
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("getShard failed: %s", err)
	}

	return result, nil
}

// getRefs finds the references for each term in query
func getRefsFromDynamo(db *dynamo.DbInfo, t string, id *IndexData) ([]string, int, error) {
	var result []string
	x := 0 // value > 0 indicates 1+ terms in query returned no matching IDs; no results for query

	// Get index list for each term
	if filter(t) { // term is filtered out; no results returned for term
		return result, x, nil
	}

	prt := getPartition(t)
	key := ""
	// get aggregate IDs from all shards
	for i := 0; i < int(id.Shards[t]+1.0); i++ {
		t = strings.TrimSpace(t)
		key = t
		if i > 0 {
			key = t + "." + strconv.Itoa(int(id.Shards[t]))
		}
		query := dynamo.CreateNewQueryObj(prt, key)

		// retreive item from DynamoDB
		tName := "cf-index"
		refObj := IndexData{}
		obj, err := dynamo.GetItem(db.Svc, query, db.Tables[tName], refObj)
		if err != nil {
			fmt.Println(err)
			return result, x, fmt.Errorf("Query failed: %v", err)
		}
		ids := obj.(DynamoSearchEntry).IDs // check
		if len(ids) == 0 {
			fmt.Println("no IDs found for term: ", t)
			x++
			return result, x, nil
		}
		fmt.Printf("ids found for '%s': %d\n", t, len(ids))
		result = append(result, ids...)
		if err != nil {
			return nil, 0, fmt.Errorf("getRefs failed: %s", err)
		}
	}
	return result, x, nil
}

// sortMap converts k:v pairs to struct, adds and sorts by len(v)
func sortMap(m map[string][]string) []Data {
	// []Data represnts inverted index and corresponding SearchData objects
	var ss []Data
	for k, v := range m {
		ss = append(ss, Data{k, sortIDs(v), len(v)}) // term, refs, len
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Len < ss[j].Len
	})

	return ss
}

// sort lookupPairs by ID (key) values
func sortIDs(lu []string) []string {
	sort.Slice(lu, func(i, j int) bool {
		return lu[i] < lu[j]
	})
	return lu
}

// intersection returns the intersection of two integer slices
func intersection(s1, s2 []string) []string {
	if len(s1) > len(s2) {
		// swap
		s1, s2 = s2, s1
	}
	if len(s1) == 0 {
		return s2
	}
	max := s1[len(s1)-1]
	checkMap := make(map[string]bool)
	common := []string{}

	for _, v := range s1 {
		checkMap[v] = true
	}
	for _, v := range s2 {
		if v > max {
			break // break if v.ID > largest ID value in smaller slice
		}
		if checkMap[v] { // common to both Entries
			common = append(common, v)
		}
	}
	return common
}

/* func main() {
	// command-line flags/if statements for choosing function
	update := flag.Bool("u", false, "update index")
	viewIndex := flag.Bool("vi", false, "view inverted index")
	viewData := flag.Bool("vd", false, "view data index")
	search := flag.Bool("s", false, "search index")

	flag.Parse()
	if *update != false {
		updateIndex()
	}
	if *viewIndex != false {
		viewInvertedIndex()
	}
	if *viewData != false {
		viewDataIndex()
	}
	if *search != false {
		err := searchIndex()
		if err != nil {
			fmt.Println(err)
		}
	}
} */
