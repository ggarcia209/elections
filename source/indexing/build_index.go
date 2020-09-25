package indexing

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/elections/source/donations"
	"github.com/elections/source/persist"
	"github.com/elections/source/ui"
)

// SearchData contains basic data for an object returned
// to the client as a result of search index query.
// Used to return basic data to user from local machine
// instead of making API call to DynamoDB table.
type SearchData struct {
	ID       string
	Name     string
	City     string
	State    string
	Employer string
	Bucket   string
	Years    []string
}

// IndexData type stores data related to the Index
type IndexData struct {
	TermsSize      int // number of terms in index
	LookupSize     int // number of lookup objects
	LastUpdated    time.Time
	Completed      map[string]bool // track categories completed in event of failure
	YearsCompleted []string
	Shards         map[string]float32
}

// inverted index
// Schema - partition: term: []objID
type indexMap map[string]map[string][]string

// k/v pairs containing obj ID & corresponding *SearchData object
// Schema - term: objID: *SearchData
type lookupPairs map[string]*SearchData

// DataMap contains the *SearchData objects inititialed for each database record
// Schema - ObjID: Obj
type DataMap map[string]*SearchData

var mu sync.Mutex

// BuildIndex creates a new search index from the objects in the db/offline_db.db
// TEST IDEMPOTENCY ACROSS YEARS
func BuildIndex(year string) error {
	var wg sync.WaitGroup

	indexData, err := getIndexData()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("BuildIndex failed: %v", err)
	}

	// get user confirmation if new index
	if len(indexData.YearsCompleted) == 0 {
		fmt.Println("*** Build new index? ***")
		y := ui.Ask4confirm()
		if !y {
			fmt.Println("operation stopped")
			fmt.Println("new items wrote: 0")
			return nil
		}
		indexData = &IndexData{ // initialize if nil
			TermsSize:   0,
			LookupSize:  0,
			LastUpdated: time.Now(),
			Completed:   map[string]bool{"individuals": false, "committees": false, "candidates": false},
		}
	}
	fmt.Println("index data: ", indexData)
	fmt.Println()

	// start goroutine to build index from each bucket not yet completed
	for k, v := range indexData.Completed {
		if v != true { // not completed
			if k == "individuals" {
				wg.Add(1)
				go indvRtn(year, indexData, &wg)
			}
			if k == "committees" {
				wg.Add(1)
				go cmteRtn(year, indexData, &wg)
			}
			if k == "candidates" {
				wg.Add(1)
				go candRtn(year, indexData, &wg)
			}
		}
	}

	wg.Wait()

	// reset map and save
	indexData.Completed = map[string]bool{"individuals": false, "committees": false, "candidates": false}
	indexData.YearsCompleted = append(indexData.YearsCompleted, year)
	err = saveIndexData(indexData)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("BuildIndex failed: %v", err)
	}

	fmt.Println("***** INDEX BUILD COMPLETE *****")
	fmt.Println("terms: ", indexData.TermsSize)
	fmt.Println("lookup items: ", indexData.LookupSize)
	fmt.Println()

	return nil
}

// UpdateIndex updates the Index with terms dervied from the given bucket
// TEST IDEMPOTENCY ACROSS YEARS
func UpdateIndex(year, bucket string) error {
	var wg sync.WaitGroup
	id, err := getIndexData()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("UpdateIndex failed: %v", err)
	}
	// re-initialize map if nil
	if len(id.Completed) == 0 {
		id.Completed = make(map[string]bool)
	}
	if len(id.Shards) == 0 {
		id.Shards = make(map[string]float32)
	}
	fmt.Printf("bucket '%s' updating...\n", bucket)

	if bucket == "individuals" {
		wg.Add(1)
		err := indvRtn(year, id, &wg)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("UpdateIndex failed: %v", err)
		}
	}
	wg.Add(1)
	if bucket == "committees" {
		err := cmteRtn(year, id, &wg)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("UpdateIndex failed: %v", err)
		}
	}
	wg.Add(1)
	if bucket == "candidates" {
		err := candRtn(year, id, &wg)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("UpdateIndex failed: %v", err)
		}
	}

	fmt.Println("***** UPDATE COMPLETE *****")
	fmt.Printf("bucket: '%s'\n", bucket)
	fmt.Println()

	return nil
}

// process Top Individual objects
func indvRtn(year string, id *IndexData, wg *sync.WaitGroup) error {
	defer wg.Done()
	// process individuals
	bucket := "individuals"
	index := make(indexMap)
	lookup := make(lookupPairs)

	err := getObjData(year, bucket, index, lookup)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("indvRtn failed: %v", err)
	}

	// update & save
	mu.Lock()
	newWrites, newTerms, err := saveIndex(id, index, lookup)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("indvRtn failed: %v", err)
	}
	updateIndexData(id, year, bucket, newWrites, newTerms)
	err = saveIndexData(id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("indvRtn failed: %v", err)
	}
	mu.Unlock()
	fmt.Println("individual data saved")
	fmt.Println("FINISHED -- ", year)
	fmt.Println()
	return nil
}

// process Committee objects
func cmteRtn(year string, id *IndexData, wg *sync.WaitGroup) error {
	defer wg.Done()
	// process committees
	bucket := "committees"
	index := make(indexMap) // reset in-memory index
	lookup := make(lookupPairs)

	err := getObjData(year, bucket, index, lookup)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("cmteRtn failed: %v", err)
	}

	// update & save
	mu.Lock()
	newWrites, newTerms, err := saveIndex(id, index, lookup)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("cmteRtn failed: %v", err)
	}
	updateIndexData(id, year, bucket, newWrites, newTerms)
	err = saveIndexData(id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("cmteRtn failed: %v", err)
	}
	mu.Unlock()
	fmt.Println("committee data saved")
	return nil
}

// process Candidates
func candRtn(year string, id *IndexData, wg *sync.WaitGroup) error {
	defer wg.Done()
	// process candidates
	bucket := "candidates"
	index := make(indexMap) // reset in-memory index
	lookup := make(lookupPairs)

	err := getObjData(year, bucket, index, lookup)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("candRtn failed: %v", err)
	}

	// update & save
	mu.Lock()
	newWrites, newTerms, err := saveIndex(id, index, lookup)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("candRtn failed: %v", err)
	}
	updateIndexData(id, year, bucket, newWrites, newTerms)
	err = saveIndexData(id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("candRtn failed: %v", err)
	}
	mu.Unlock()
	fmt.Println("candidate data saved")
	return nil
}

// TEMPORARILY DEPRECATED
// process top individuals by funds received/sent and add to Index
func getTopIndvData(year, bucket string, index indexMap, lookup lookupPairs) error {
	fmt.Println("processing top individuals...")
	ids := []string{}
	n := 20000
	x := n

	// get Top Individuals by incoming & outgoing funds
	odID := year + "-individuals-donor-ALL"
	topIndv, err := persist.GetObject(year, "top_overall", odID)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("getTopIndvData failed: %v", err)
	}
	odID = year + "-individuals-rec-ALL"
	topIndvRec, err := persist.GetObject(year, "top_overall", odID)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("getTopIndvData failed: %v", err)
	}

	fmt.Println("Got TopIndv Objects")

	// create list of IDs for BatchGetByID
	for k := range topIndv.(*donations.TopOverallData).Amts {
		ids = append(ids, k)
	}
	for k := range topIndvRec.(*donations.TopOverallData).Amts {
		ids = append(ids, k)
	}

	// get partition map
	pm, err := getPartitionMap()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("getTopIndvData failed: %v", err)
	}
	if len(pm) == 0 {
		pm = make(map[string]bool)
	}

	for {
		if len(ids) < n { // queue exhausted - last write
			n = len(ids) // set starting index to 0
		}
		// pop n IDs from stack and return corresponding objects
		objs, _, err := persist.BatchGetByID(year, bucket, ids[len(ids)-n:])
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("getTopIndvData failed: %v", err)
		}
		if len(objs) == 0 {
			break
		}

		// create SearchData objects (map[id]*SearchData)
		lu := createSearchData(year, objs)

		// proces SearchData objects and add terms to Index
		for k, sd := range lu {
			// derive search terms from object data
			terms := getTerms(sd)
			termsFmt := formatTerms(terms)
			for _, term := range termsFmt {
				if filter(term) {
					continue
				}
				prt := getPartition(term)
				if index[prt] == nil {
					index[prt] = make(map[string][]string)
				}

				index[prt][term] = append(index[prt][term], k)
				lookup[k] = sd
				pm[prt] = true
			}
		}

		if len(objs) < x || len(ids) == n { // last batch write complete
			break
		}

		// remove processed IDs from stack
		ids = ids[:len(ids)-n]
	}

	mu.Lock()
	err = savePartitionMap(pm)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("getTopIndvData failed: %v", err)
	}
	mu.Unlock()

	fmt.Println("partition map saved")

	fmt.Println("***** BUILD INDEX - 'individuals' FINSIHED *****")
	fmt.Println()

	return nil
}

// add Candidate/Committee object info to Index
func getObjData(year, bucket string, index indexMap, lookup lookupPairs) error {
	n := 10000
	startKey := ""

	// get partition map
	pm, err := getPartitionMap()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("getObjData failed: %v", err)
	}
	if len(pm) == 0 {
		pm = make(map[string]bool)
	}

	for {
		// get next batch of objects
		objs, currKey, err := persist.BatchGetSequential(year, bucket, startKey, n)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("getObjData failed: %v", err)
		}
		if len(objs) == 0 {
			break
		}

		// create SearchData objects
		lu := createSearchData(year, objs)

		// proces SearchData objects and add terms to Index
		for k, sd := range lu {
			// derive search terms from object data
			terms := getTerms(sd)
			termsFmt := formatTerms(terms)
			for _, term := range termsFmt {
				if filter(term) {
					continue
				}

				prt := getPartition(term)
				if index[prt] == nil {
					index[prt] = make(map[string][]string)
				}
				// update maps by reference
				index[prt][term] = append(index[prt][term], k)
				lookup[k] = sd
				pm[prt] = true
			}
		}

		if len(objs) < n { // last batch write complete
			break
		}

		startKey = currKey
	}

	mu.Lock()
	err = savePartitionMap(pm)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("getObjData failed: %v", err)
	}
	mu.Unlock()

	fmt.Printf("***** BUILD INDEX - '%s' FINSIHED *****\n", bucket)
	fmt.Println()

	return nil
}

// create list of SearchData objects from returned objects
func createSearchData(year string, objs []interface{}) lookupPairs {
	lookup := make(lookupPairs)

	for _, obj := range objs {
		switch t := obj.(type) {
		case *donations.Individual:
			sd := &SearchData{
				ID:       obj.(*donations.Individual).ID,
				Name:     obj.(*donations.Individual).Name,
				City:     obj.(*donations.Individual).City,
				State:    obj.(*donations.Individual).State,
				Employer: obj.(*donations.Individual).Employer,
				Bucket:   "individuals",
				Years:    []string{year},
			}
			lookup[sd.ID] = sd
		case *donations.Committee:
			sd := &SearchData{
				ID:       obj.(*donations.Committee).ID,
				Name:     obj.(*donations.Committee).Name,
				City:     obj.(*donations.Committee).City,
				State:    obj.(*donations.Committee).State,
				Employer: obj.(*donations.Committee).Party,
				Bucket:   "committees",
				Years:    []string{year},
			}
			lookup[sd.ID] = sd
		case *donations.Candidate:
			sd := &SearchData{
				ID:       obj.(*donations.Candidate).ID,
				Name:     obj.(*donations.Candidate).Name,
				City:     obj.(*donations.Candidate).City,
				State:    obj.(*donations.Candidate).State,
				Employer: obj.(*donations.Candidate).Party,
				Bucket:   "candidates",
				Years:    []string{year},
			}
			lookup[sd.ID] = sd
		default:
			_ = t
			fmt.Println("createSearchData err: invalid interface type")
			return lookup
		}
	}

	return lookup
}

// derive search terms from SearchData
func getTerms(sd *SearchData) []string {
	name := sd.Name
	city := sd.City
	states := map[string]string{
		"AL": "Alabama", "AK": "Alaska", "AZ": "Arizona", "AR": "Arkansas",
		"CA": "California", "CO": "Colorado", "CT": "Connecticut", "DE": "Delaware",
		"FL": "Florida", "GA": "Georgia", "HI": "Hawaii", "ID": "Idaho",
		"IL": "Illinois", "IN": "Indiana", "IA": "Iowa", "KS": "Kansas", "KY": "Kentucky",
		"LA": "Louisiana", "ME": "Maine", "MD": "Maryland", "MA": "Massachusetts",
		"MI": "Michigan", "MN": "Minnesota", "MS": "Mississippi", "MO": "Missouri",
		"MT": "Montana", "NE": "Nebraska", "NV": "Nevada", "NH": "New Hampshire",
		"NJ": "New Jersey", "NM": "New Mexico", "NY": "New York", "NC": "North Carolina",
		"ND": "North Dakota", "OH": "Ohio", "OK": "Oklahoma", "OR": "Oregon", "PA": "Pennsylvania",
		"RI": "Rhode Island", "SC": "South Carolina", "SD": "South Dakota", "TN": "Tennessee",
		"TX": "Texas", "UT": "Utah", "VT": "Vermont", "VA": "Virginia", "WA": "Washington",
		"WV": "West Virginia", "WI": "Wisconsin", "WY": "Wyoming",
	}
	state := states[sd.State]
	employer := sd.Employer

	return []string{name, city, state, employer}
}

// formatTerms derives and formats search terms from a SearchData object
// (ex; "Bush, George H.W. -> []string{"bush", "george", "hw")
func formatTerms(terms []string) []string {
	fmtStrs := []string{}
	for _, term := range terms {
		if filter(term) {
			continue
		}
		// remove & replace non-alpha-numeric characters and lowercase text
		reg, err := regexp.Compile("[^a-zA-Z0-9]+") // removes all non alpha-numeric characters
		if err != nil {
			log.Fatal(err)
		}
		rmApost := strings.Replace(term, "'", "", -1)    // don't split contractions (ex: 'can't' !-> "can", "t")
		rmComma := strings.Replace(rmApost, ",", "", -1) // don't split numerical values > 999 (ex: 20,000 !-> 20 000)
		lwr := strings.ToLower(rmComma)
		regged := reg.ReplaceAllString(lwr, " ")
		spl := strings.Split(regged, " ")
		for _, s := range spl {
			trim := strings.TrimSpace(s)
			if trim != "" {
				fmtStrs = append(fmtStrs, trim)
			}
		}
	}

	return fmtStrs
}

// derive partition (first letter) from term
func getPartition(term string) string {
	s := strings.Split(term, "")
	p := s[0]
	return p
}

// update IndexData object
func updateIndexData(id *IndexData, year, bucket string, newWrites, newTerms int) {
	id.TermsSize += newTerms
	id.LookupSize += (newWrites * 0) // temporarily disabled
	id.Completed[bucket] = true
	id.LastUpdated = time.Now()
	return
}

// filter generic terms & edge cases ("the", "for", "of", "",)
// returns true if term meets filter criteria
func filter(term string) bool {
	f := map[string]bool{
		"for": true,
		"the": true,
		"of":  true,
		"":    true,
	}
	return f[term]
}
