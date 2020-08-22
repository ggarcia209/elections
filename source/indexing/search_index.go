package indexing

import (
	"fmt"
	"sort"
	"strings"
	"time"

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

// Entries represents a sorted map
type Entries []Entry

func (s Entries) Len() int           { return len(s) }
func (s Entries) Less(i, j int) bool { return s[i].ID < s[j].ID }
func (s Entries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

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
func GetResults(q Query) ([]SearchData, error) {
	if q.Text == "" {
		return nil, nil
	}

	// find corresponding SearchData objects
	terms := formatTerms(strings.Split(q.Text, " "))

	// get IDs for single term
	if len(terms) == 1 {
		results, err := getSearchEntry(terms[0])
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("GetResults failed: %v", err)
		}

		return results, nil
	}

	// get IDs for multiple terms
	resMap, err := getRefs(terms)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetResults failed: %v", err)
	}

	// Sort lists by smallest to largest
	sorted := sortMap(resMap)

	// Compare and find all common IDs in all terms
	// Start by finding the common IDs in the 2 smallest lists
	// then compare the next list to the previous comparison's intersection
	s1, s2 := sorted[0].Value, sorted[1].Value
	common := intersection(s1, s2)
	for i := 2; i < len(sorted); i++ {
		common = intersection(common, sorted[i].Value)
	}

	// Get data for the common IDs
	results, err := LookupSearchData(common)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetResults failed: %v", err)
	}

	return results, nil
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
	db, err := bolt.Open("../../db/search_index.db", 0644, nil)
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

// getRefs finds the references for each term in query
func getRefs(q []string) (map[string][]string, error) {
	var resultMap = make(map[string][]string)
	var result []string

	db, err := bolt.Open("../../db/search_index.db", 0644, nil)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("getRefs failed: %v", err)
	}
	defer db.Close()

	// Get index list for each term in query - use map
	for _, v := range q {
		if filter(v) {
			continue
		}
		prt := getPartition(v)
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(prt))
			v = strings.TrimSpace(v)
			data := b.Get([]byte(v))
			ids, err := decodeResultsList(data)
			if err != nil {
				return fmt.Errorf("tx failed: %v", err)
			}
			result = ids
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("getRefs failed: %s", err)
		}
		resultMap[v] = result
	}
	return resultMap, nil
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
	checkMap := map[string]bool{}
	common := []string{}
	for _, v := range s1 {
		checkMap[v] = true
	}
	for _, v := range s2 {
		if v > s1[len(s1)-1] {
			break // break if v.ID > largest ID value in smaller slice
		}
		if _, ok := checkMap[v]; ok { // common to both Entries
			common = append(common, v)
		}
	}
	return common
}

// retreive corresponding SearchData obj for ID
func LookupSearchData(ids []string) ([]SearchData, error) {
	db, err := bolt.Open("../../db/search_index.db", 0644, nil)
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
