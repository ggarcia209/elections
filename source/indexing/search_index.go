package indexing

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/boltdb/bolt"
	"github.com/elections/source/dynamo"
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

// SearchEntry is used to store & retreive a Search Index entry from DynamoDB
type SearchEntry struct {
	Partition string
	Term      string
	IDs       []string
}

// LookupEntry is used to store & retreive SearchData objects from DynamoDB
type LookupEntry struct {
	Partition string // first 2 chars of FEC ID / first 3 letters of MD5 hash ID
	ID        string
	Name      string
	City      string
	State     string
	Bucket    string
	Employer  string
	Years     []string
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

type chanResult struct {
	SD  SearchData
	Err error
}

// CreateQuery returns a Query object with the given data
func CreateQuery(text, UID string) Query {
	q := Query{
		Text:      text,
		UserID:    UID,
		TimeStamp: time.Now(),
	}
	return q
}

// GetResultsFromShards returns search results from the sharded index stored on disk
func GetResultsFromShards(id *IndexData, q Query) ([]string, error) {
	st := time.Now()
	if q.Text == "" {
		return []string{}, nil
	}

	terms := formatTerms(strings.Split(q.Text, " ")) // normalize search terms input
	common := []string{}                             // aggegate total of intersections for every shard
	maxResultsSize := 500                            // max number of SearchResults returned before throwing MAX_LENGTH error

	// open db
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("getShard failed: %v", err)
	}
	defer db.Close()

	// get IDs for single term
	if len(terms) == 1 {
		t := terms[0]
		if id.Shards[t] != nil && id.Shards[t].Shards > 0 {
			fmt.Println("GetResults failed: MAX_LENGTH exceeded")
			return []string{}, fmt.Errorf("MAX_LENGTH")
		}
		t = strings.TrimSpace(t)

		// check if lookup by ID
		lookup := checkForID(t)
		if lookup {
			return []string{strings.ToUpper(t)}, nil
		}

		ids, err := getShard(db, t)
		if err != nil {
			fmt.Println(err)
			return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
		}
		if len(ids) == 0 { // no IDs found for term - no results
			fmt.Println("GetResults failed: NO_RESULTS")
			return []string{}, fmt.Errorf("NO_RESULTS")
		}
		// add results to common aggregate total; throw error if maxResultsSize reached
		for i, ID := range ids {
			if i >= maxResultsSize {
				fmt.Println("GetResults failed: MAX_LENGTH exceeded")
				return []string{}, fmt.Errorf("MAX_LENGTH")
			}
			common = append(common, ID)
		}
		fmt.Println("finish - ", time.Since(st))
		return common, nil
	}

	// sort terms by # of shards, least to greatest
	termsSrt := shardEntries{}
	for _, t := range terms {
		sh := float32(0)
		if id.Shards[t] != nil {
			sh = id.Shards[t].Shards
		}
		se := shardEntry{ID: t, Shards: sh}
		termsSrt = append(termsSrt, se)
	}
	sort.Sort(termsSrt)

	// find intersections of two smallest records - terms[0] & terms[1] (t0, t1)
	i0, i1 := termsSrt[0].ID, termsSrt[1].ID // terms
	k0, k1 := i0, i1                         // shard keys
	ct0, ct1 := float32(0.0), float32(0.0)   // number of shards created for each term
	if id.Shards[i0] != nil {
		ct0 = id.Shards[i0].Shards
	}
	if id.Shards[i1] != nil {
		ct1 = id.Shards[i1].Shards
	}
	// compare each shard (s0) in t0 to each shard in t1 (s1) within min/max range
	for x := 0; x < int((ct0 + 1)); x++ {
		s0 := []string{} // shard IDs - used to store IDs after first disk read call in inner loop
		if x > 0 {
			k0 = i0 + "." + strconv.Itoa(x)
		}

		min0, max0 := "", ""
		if ct0 == 0 {
			st1 := time.Now()
			s0, err = getShard(db, k0)
			if err != nil {
				fmt.Println(err)
				return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
			}
			if len(s0) == 0 { // no IDs found for term - no results
				fmt.Println("GetResults failed: NO_RESULTS")
				return []string{}, fmt.Errorf("NO_RESULTS")
			}
			fmt.Println("s0 read time: ", time.Since(st1))
			min0, max0 = s0[0], s0[len(s0)-1]
		} else {
			r0 := id.Shards[i0].Ranges[k0]
			min0, max0 = r0.Range[0], r0.Range[1] // min, max value of each shard in t1 compared to these values
		}

		// get each shard in t1 (s1), compare to s0 and find intersection
		for y := 0; y < int((ct1 + 1)); y++ {
			s1 := []string{}
			if y > 0 {
				k1 = i1 + "." + strconv.Itoa(y)
			}

			min1, max1 := "", ""
			if ct1 == 0 {
				s1, err := getShard(db, k1)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
				}
				min1, max1 = s1[0], s1[len(s1)-1]
			} else {
				r1 := id.Shards[i1].Ranges[k1]
				min1, max1 = r1.Range[0], r1.Range[1]
			}

			if min0 > max1 {
				fmt.Println("skipping shard ", y)
				continue // skip
			}
			if max0 < min1 {
				fmt.Println("breaking at shard ", y)
				break // stop
			}

			if len(s1) == 0 {
				s1, err = getShard(db, k1)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("getRefs failed: %v", err)
				}
			}

			if len(s0) == 0 {
				s0, err = getShard(db, k0)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("getRefs failed: %v", err)
				}
			}

			// find intersection of s0, s1, add to common aggregate total
			intsec := intersection(s0, s1)
			common = append(common, intsec...)
		}
	}

	// find intersections of each remaining term
	for t := 2; t < len(termsSrt); t++ {
		if len(common) == 0 {
			return common, nil
		}
		i2 := termsSrt[t].ID // next term in list
		k2 := i2
		buffer := common    // use to find values in common within range of current s1 shard
		common = []string{} // reset for each term

		// get each shard in t1 (s1), compare to sc and find intersection
		sc := []string{} // shard created from common IDs within range of s2
		ct2 := float32(0.0)
		if id.Shards[i2] != nil {
			ct2 = id.Shards[i2].Shards
		}
		for z := 0; z < int((ct2 + 1)); z++ {
			s2 := []string{}
			if z > 0 {
				k2 = i2 + "." + strconv.Itoa(z)
			}

			min2, max2 := "", ""
			if ct2 == 0 {
				s2, err = getShard(db, k2)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
				}
				min2, max2 = s2[0], s2[len(s2)-1]
			} else {
				r2 := id.Shards[i2].Ranges[k2]
				min2, max2 = r2.Range[0], r2.Range[1]
			}

			// create comparision shard from buffer IDs within range
			for _, v := range buffer {
				if v >= max2 { // no more common values exist
					break // stop
				}
				if v >= min2 {
					sc = append(sc, v)
				}
			}

			if len(s2) == 0 {
				s2, err = getShard(db, k2)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
				}
			}

			// find intersection of s0, s2, add to common aggregate total
			intsec := intersection(sc, s2)
			common = append(common, intsec...)
		}
		sort.Slice(common, func(i, j int) bool { return common[i] < common[j] })
	}

	fmt.Println("len(common): ", len(common))
	if len(common) > maxResultsSize {
		fmt.Println("GetResults failed: MAX_LENGTH exceeded")
		return []string{}, fmt.Errorf("MAX_LENGTH")
	}

	fmt.Println("finish - ", time.Since(st))
	return common, nil
}

// GetResultsFromDynamo returns search results from the sharded index stored in DynamoDB
func GetResultsFromDynamo(db *dynamo.DbInfo, id *IndexData, q Query) ([]string, error) {
	st := time.Now()
	if q.Text == "" {
		return []string{}, nil
	}

	terms := formatTerms(strings.Split(q.Text, " ")) // normalize search terms input
	common := []string{}                             // aggegate total of intersections for every shard
	maxResultsSize := 500                            // max number of SearchResults returned before throwing MAX_LENGTH error
	var err error

	// get IDs for single term
	if len(terms) == 1 {
		t := terms[0]
		if id.Shards[t] != nil && id.Shards[t].Shards > 0 {
			fmt.Println("GetResults failed: MAX_LENGTH exceeded")
			return []string{}, fmt.Errorf("MAX_LENGTH")
		}
		t = strings.TrimSpace(t)

		// check if lookup by ID
		lookup := checkForID(t)
		if lookup {
			return []string{strings.ToUpper(t)}, nil
		}

		st1 := time.Now()
		ids, err := getShardFromDynamo(db, t)
		if err != nil {
			fmt.Println(err)
			return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
		}
		fmt.Println("get shard time: ", time.Since(st1))
		// add results to common aggregate total; throw error if maxResultsSize reached
		for i, ID := range ids {
			if i >= maxResultsSize {
				fmt.Println("GetResults failed: MAX_LENGTH exceeded")
				return []string{}, fmt.Errorf("MAX_LENGTH")
			}
			common = append(common, ID)
		}
		fmt.Println("finish - ", time.Since(st))
		return common, nil
	}

	// sort terms by # of shards, least to greatest
	termsSrt := shardEntries{}
	for _, t := range terms {
		sh := float32(0)
		if id.Shards[t] != nil {
			sh = id.Shards[t].Shards
		}
		se := shardEntry{ID: t, Shards: sh}
		termsSrt = append(termsSrt, se)
	}
	sort.Sort(termsSrt)

	// find intersections of two smallest records - terms[0] & terms[1] (t0, t1)
	i0, i1 := termsSrt[0].ID, termsSrt[1].ID // terms
	k0, k1 := i0, i1                         // shard keys
	ct0, ct1 := float32(0.0), float32(0.0)   // number of shards created for each term
	if id.Shards[i0] != nil {
		ct0 = id.Shards[i0].Shards
	}
	if id.Shards[i1] != nil {
		ct1 = id.Shards[i1].Shards
	}
	// compare each shard (s0) in t0 to each shard in t1 (s1) within min/max range
	for x := 0; x < int((ct0 + 1)); x++ {
		s0 := []string{} // shard IDs - used to store IDs after first disk read call in inner loop
		if x > 0 {
			k0 = i0 + "." + strconv.Itoa(x)
		}

		min0, max0 := "", ""
		if ct0 == 0 {
			st1 := time.Now()
			s0, err = getShardFromDynamo(db, k0)
			if err != nil {
				fmt.Println(err)
				return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
			}
			fmt.Println("s0 read time: ", time.Since(st1))
			min0, max0 = s0[0], s0[len(s0)-1]
		} else {
			r0 := id.Shards[i0].Ranges[k0]
			min0, max0 = r0.Range[0], r0.Range[1] // min, max value of each shard in t1 compared to these values
		}

		// get each shard in t1 (s1), compare to s0 and find intersection
		for y := 0; y < int((ct1 + 1)); y++ {
			s1 := []string{}
			if y > 0 {
				k1 = i1 + "." + strconv.Itoa(y)
			}

			min1, max1 := "", ""
			if ct1 == 0 {
				s1, err := getShardFromDynamo(db, k1)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
				}
				min1, max1 = s1[0], s1[len(s1)-1]
			} else {
				r1 := id.Shards[i1].Ranges[k1]
				min1, max1 = r1.Range[0], r1.Range[1]
			}

			if min0 > max1 {
				fmt.Println("skipping shard ", y)
				continue // skip
			}
			if max0 < min1 {
				fmt.Println("breaking at shard ", y)
				break // stop
			}

			if len(s1) == 0 {
				s1, err = getShardFromDynamo(db, k1)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("getRefs failed: %v", err)
				}
			}

			if len(s0) == 0 {
				s0, err = getShardFromDynamo(db, k0)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("getRefs failed: %v", err)
				}
			}

			// find intersection of s0, s1, add to common aggregate total
			intsec := intersection(s0, s1)
			common = append(common, intsec...)
		}
	}

	fmt.Println("common - s0, s1: ", len(common))

	// find intersections of each remaining term
	for t := 2; t < len(termsSrt); t++ {
		if len(common) == 0 {
			return common, nil
		}
		i2 := termsSrt[t].ID // next term in list
		k2 := i2
		buffer := common    // use to find values in common within range of current s1 shard
		common = []string{} // reset for each term

		// get each shard in t1 (s1), compare to sc and find intersection
		sc := []string{} // shard created from common IDs within range of s2
		ct2 := float32(0.0)
		if id.Shards[i2] != nil {
			ct2 = id.Shards[i2].Shards
		}
		for z := 0; z < int((ct2 + 1)); z++ {
			s2 := []string{}
			if z > 0 {
				k2 = i2 + "." + strconv.Itoa(z)
			}

			min2, max2 := "", ""
			if ct2 == 0 {
				s2, err = getShardFromDynamo(db, k2)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
				}
				min2, max2 = s2[0], s2[len(s2)-1]
			} else {
				r2 := id.Shards[i2].Ranges[k2]
				min2, max2 = r2.Range[0], r2.Range[1]
			}

			// create comparision shard from buffer IDs within range
			for _, v := range buffer {
				if v >= max2 { // no more common values exist
					break // stop
				}
				if v >= min2 {
					sc = append(sc, v)
				}
			}

			if len(s2) == 0 {
				s2, err = getShardFromDynamo(db, k2)
				if err != nil {
					fmt.Println(err)
					return []string{}, fmt.Errorf("GetResultsFromShards failed: %v", err)
				}
			}

			// find intersection of s0, s2, add to common aggregate total
			intsec := intersection(sc, s2)
			common = append(common, intsec...)
		}
		sort.Slice(common, func(i, j int) bool { return common[i] < common[j] })
	}

	fmt.Println("len(common): ", len(common))
	if len(common) > maxResultsSize {
		fmt.Println("GetResults failed: MAX_LENGTH exceeded")
		return []string{}, fmt.Errorf("MAX_LENGTH")
	}

	fmt.Println("finish - ", time.Since(st))
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

// LookupSearchDataFromDynamo retreives corresponding SearchData obj for ID from DynamoDB
func LookupSearchDataFromDynamo(db *dynamo.DbInfo, ids []string) ([]SearchData, error) {
	tn := "cf-lookup"
	results := []SearchData{}
	dataChan := make(chan chanResult)
	var wg sync.WaitGroup

	for _, ID := range ids {
		wg.Add(1)
		go getSearchDataAsync(db, ID, tn, dataChan, &wg)
	}
	for i := 0; i < len(ids); i++ {
		msg := <-dataChan
		obj, err := msg.SD, msg.Err
		if err != nil {
			fmt.Println(err)
		}
		sd := SearchData{
			ID:       obj.ID,
			Name:     obj.Name,
			City:     obj.City,
			State:    obj.State,
			Employer: obj.Employer,
			Bucket:   obj.Bucket,
			Years:    obj.Years,
		}
		results = append(results, sd)
	}
	wg.Wait()

	return results, nil
}

func getSearchDataAsync(db *dynamo.DbInfo, ID, tn string, dc chan chanResult, wg *sync.WaitGroup) {
	defer wg.Done()
	prt := getLookupDynamoPrt(ID)
	q := dynamo.CreateNewQueryObj(prt, ID)
	refObj := LookupEntry{}
	sd := SearchData{
		ID:       "",
		Name:     "Unknown",
		City:     "???",
		State:    "???",
		Employer: "",
		Bucket:   "indviduals",
		Years:    []string{"all-time"},
	}

	retries := 0
	maxRetries := 5
	for { // retry loop - breaks if successful / returns if maxRetries exceeded
		obj, err := dynamo.GetItem(db.Svc, q, db.Tables[tn], refObj)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() == "RequestError" {
					// wait and retry
					fmt.Println("Request failed - retrying...")
					time.Sleep(50 * time.Millisecond)
					retries++
					if retries > maxRetries {
						msg := "MAX_RETRIES_EXCEEDED"
						fmt.Println(msg)
						err = fmt.Errorf(msg)
						cr := chanResult{sd, err}
						dc <- cr
						return
					}
					continue
				}
				fmt.Println("getSearchDataAsync failed: ", awsErr.Code(), awsErr.Message())
			} else {
				fmt.Println("getSearchDataAsync failed: ", err.Error())
				cr := chanResult{sd, err}
				dc <- cr
				return
			}
		}
		if err != nil {
			fmt.Println("getSearchDataAsync failed: ", err)
			cr := chanResult{sd, err}
			dc <- cr
			return
		}
		intf := obj.(map[string]interface{})

		id, name, city, state, employer, bucket := "", "", "", "", "", ""
		years := []string{}
		if intf["ID"] == nil {
			id = "000000000"
		} else {
			id = intf["ID"].(string)
		}
		if intf["Name"] == nil {
			name = "Unknown"
		} else {
			name = intf["Name"].(string)
		}
		if intf["City"] == nil {
			city = "???"
		} else {
			city = intf["City"].(string)
		}
		if intf["State"] == nil {
			state = "???"
		} else {
			state = intf["State"].(string)
		}
		if intf["Employer"] == nil {
			employer = ""
		} else {
			employer = intf["Employer"].(string)
		}
		if intf["Bucket"] == nil {
			bucket = "committees"
		} else {
			bucket = intf["Bucket"].(string)
		}
		if intf["Years"] == nil {
			years = append(years, "all-time")
		} else {
			for _, yr := range intf["Years"].([]interface{}) {
				years = append(years, yr.(string))
			}
		}
		sd = SearchData{
			ID:       id,
			Name:     name,
			City:     city,
			State:    state,
			Employer: employer,
			Bucket:   bucket,
			Years:    years,
		}
		cr := chanResult{sd, nil}
		dc <- cr
		break
	}

	return
}

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

	// sort results by most recent year
	sort.Slice(agg, func(i, j int) bool {
		if agg[i].Years[0] != agg[j].Years[0] {
			return agg[i].Years[0] > agg[j].Years[0]
		}
		return len(agg[i].Years) > len(agg[j].Years)
	})

	return agg
}

// ViewIndex displays the index
func ViewIndex() error {
	ct := 0

	// get partition map
	partitions, err := GetPartitionMap()
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

	totalWCU := 0.0
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
				fmt.Printf("Partition: %s\tTerm: %s\tReferences: %d\tSize: %d\n", prt.Prt, k, len(d), len(v))
				ct++
				wcuPerShard := float64((len(v) / 1000)) // 1KB = 1 WCU
				totalWCU += wcuPerShard
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ViewIndex failed: %v", err)
	}

	rate := 1.25 / 1000000.0
	price := totalWCU * rate
	fmt.Println()
	fmt.Println("***** SCAN FINISHED *****")
	fmt.Println()
	fmt.Printf("total price: %.2f\n", price)

	fmt.Println("\nTotal entries: ", ct)
	return nil
}

// ViewIndex displays the index
func ViewLookup() error {
	ct := 0

	// open db
	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ViewIndex failed: %v", err)
	}
	defer db.Close()

	totalWCU := 0.0
	i := 0
	// iterate through each key in each partiton in alphabetical order
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lookup"))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			fmt.Printf("%d.\t%s\n", i, k)
			ct++
			totalWCU++
			i++
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ViewIndex failed: %v", err)
	}

	rate := 1.25 / 1000000.0
	price := totalWCU * rate
	fmt.Println()
	fmt.Println("***** SCAN FINISHED *****")
	fmt.Println()
	fmt.Printf("total price: %.2f\n", price)

	fmt.Println("\nTotal entries: ", ct)
	return nil
}

// BatchGetSequential retrieves a sequential list of n objects from the database starting at the given key.
// Returns []SearchData and []SearchEntry objects. SearchEntry objects are partitioned for upload to DynamoDB.
func BatchGetSequential(bucket, startKey string, n int) ([]interface{}, string, error) {
	objs := []interface{}{}
	currKey := startKey

	db, err := bolt.Open(OUTPUT_PATH+"/db/search_index.db", 0644, nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return nil, currKey, fmt.Errorf("BatchGetSequential failed: %v", err)
	}

	if err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))

		c := b.Cursor()

		if startKey == "" {
			skBytes, _ := c.First()
			startKey = string(skBytes)
		}

		for k, v := c.Seek([]byte(startKey)); k != nil; k, v = c.Next() {
			if bucket == "lookup" {
				sd, err := decodeSearchData(v)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("tx failed: %v", err)
				}
				prt := getLookupDynamoPrt(sd.ID)
				obj := LookupEntry{
					Partition: prt,
					ID:        sd.ID,
					Name:      sd.Name,
					City:      sd.City,
					State:     sd.State,
					Employer:  sd.Employer,
					Bucket:    sd.Bucket,
					Years:     sd.Years,
				}
				objs = append(objs, obj)
			} else {
				res, err := decodeResultsList(v)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("tx failed: %v", err)
				}
				// get dynamo partiton name (first letter + partiton number)
				prt := getDynamoPrt(string(k))
				obj := SearchEntry{Partition: prt, Term: string(k), IDs: res}
				objs = append(objs, obj)
			}
			currKey = string(k)
			if len(objs) == n {
				break
			}
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil, currKey, fmt.Errorf("BatchGetSequential failed: %v", err)
	}
	if len(objs) < n {
		currKey = ""
	}
	return objs, currKey, nil
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

/*
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
		for i := 0; i < int(id.Shards[v].Shards+1.0); i++ {
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
} */

// check if single term query is ID for lookup
func checkForID(id string) bool {
	fmt.Println("ID: ", id)
	if len(id) == 32 {
		fmt.Println("ID found - len 32")
		return true
	}
	if len(id) != 9 {
		fmt.Println("not ID - len != 9 || 32")
		return false
	}

	fecIDs := map[string]bool{
		"C": true, "c": true, "H": true, "h": true,
		"S": true, "s": true, "P": true, "p": true,
	}
	nums := map[string]bool{
		"0": true, "1": true, "2": true, "3": true, "4": true,
		"5": true, "6": true, "7": true, "8": true, "9": true,
	}
	ss := strings.Split(id, "")
	for i, s := range ss {
		if i == 0 {
			if !fecIDs[s] {
				fmt.Println("not FEC ID code")
				return false
			}
		}
		if i > 0 && !nums[s] {
			fmt.Println("not id - alphanumeric sequence found")
			return false
		}
	}
	fmt.Println("ID found")
	return true
}

// get ids in single shard
func getShard(db *bolt.DB, id string) ([]string, error) {
	result := []string{}
	prt := getPartition(id)

	err := db.View(func(tx *bolt.Tx) error {
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

// get ids in single shard from DynamoDB
func getShardFromDynamo(db *dynamo.DbInfo, id string) ([]string, error) {
	dPrt := getDynamoPrt(id)
	dq := dynamo.CreateNewQueryObj(dPrt, id)
	fmt.Println("dq: ", dq)
	fmt.Println("tables: ", db.Tables)

	shard, err := dynamo.GetItem(db.Svc, dq, db.Tables["cf-index"], &SearchEntry{})
	if err != nil {
		fmt.Println(err)
		return []string{}, fmt.Errorf("getShardFromDynamo failed: ", err)
	}

	ids := shard.(*SearchEntry).IDs

	return ids, nil
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
	for i := 0; i < int(id.Shards[t].Shards+1.0); i++ {
		t = strings.TrimSpace(t)
		key = t
		if i > 0 {
			key = t + "." + strconv.Itoa(int(id.Shards[t].Shards))
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
		ids := obj.(SearchEntry).IDs // check
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

func getDynamoPrt(term string) string {
	ss := strings.Split(term, ".")
	pn := "" // partition number
	if len(ss) > 1 {
		pn = "." + ss[1]
	}
	sss := strings.Split(ss[0], "")
	prt := sss[0] + pn
	return prt
}

func getLookupDynamoPrt(id string) string {
	prt := ""
	if len(id) == 9 { // get last 2 chars of FEC ID
		ss := strings.Split(id, "")
		st := ss[len(id)-2:]
		prt := strings.Join(st, "")
		return prt
	}
	if len(id) == 32 { // get first 2 chars of hash ID
		ss := strings.Split(id, "")
		st := ss[:2]
		prt := strings.Join(st, "")
		return prt
	}
	return prt
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
