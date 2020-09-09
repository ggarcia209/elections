package server

import (
	"fmt"
	"strings"

	"github.com/elections/source/persist"

	"github.com/elections/source/donations"
	"github.com/elections/source/dynamo"

	"github.com/elections/source/indexing"
)

// RankingsMap stores references to each rankings list by year
type RankingsMap map[string]map[string]donations.TopOverallData

// YrTotalsMap stores references to each yearly total by year
type YrTotalsMap map[string]map[string]donations.YearlyTotal

// InitServerDiskCache creates the ../db directory on the local disk
func InitServerDiskCache() {
	persist.InitDiskCache()
	indexing.OUTPUT_PATH = ".."
	fmt.Println("local disk cache created")
}

// InitDynamo initialized a Dynamo session with default settings
func InitDynamo() (*dynamo.DbInfo, error) {
	db, err := initDynamoDbDefault()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetRankingsFromDisk failed: %v", err)
	}
	return db, nil
}

// SearchData takes a user query as a string and
// finds the results matching each word in query.
func SearchData(txt string) ([]indexing.SearchData, error) {
	// get query from user / return & print results
	q := indexing.CreateQuery(txt, "user")
	res, err := indexing.GetResults(q)
	if err != nil {
		fmt.Println(err)
		return []indexing.SearchData{}, fmt.Errorf("QueryData failed: %v", err)
	}
	return res, nil
}

// GetObjectFromDynamo returns the yearly datasets
// for the queried object and the given years
func GetObjectFromDynamo(db *dynamo.DbInfo, query *dynamo.Query, bucket string, years []string) ([]interface{}, error) {
	datasets := []interface{}{}

	err := initDynamoTables(db, years)
	if err != nil {
		fmt.Println(err)
		return datasets, fmt.Errorf("GetObjectFromDynamo failed: %v", err)
	}

	// retreive item's datasets for each year from db
	refObj := getRefObj(bucket)
	for _, yr := range years {
		tName := "cf-" + yr + "-" + bucket
		obj, err := dynamo.GetItem(db.Svc, query, db.Tables[tName], refObj)
		if err != nil {
			fmt.Println(err)
			return datasets, fmt.Errorf("GetObjectFromDynamo failed: %v", err)
		}
		datasets = append(datasets, obj)
	}

	return datasets, nil
}

// ADMIN - CHANGE UPLOAD PARTITION KEYS TO FIRST LETTER OF LAST NAME
// CreateQueryFromSearchData returns a Dynamo Query object from SearchData info
func CreateQueryFromSearchData(sd indexing.SearchData) *dynamo.Query {
	ns := strings.Split(sd.Name, "")
	pk := ns[0]
	return &dynamo.Query{PrimaryValue: pk, SortValue: sd.ID}
}

// GetRankingsFromDynamo retrieves the TopOvearll datasets
// for the given year from Dynamo to store in memory
func GetRankingsFromDynamo(db *dynamo.DbInfo) (RankingsMap, error) {
	rankings := make(RankingsMap)
	queries := []*dynamo.Query{}
	years := []string{
		"2020", "2018", "2016", "2014", "2012", "2010", "2008", "2006", "2004", "2002",
		"2000", "1998", "1996", "1994", "1992", "1990", "1988", "1986", "1984", "1982",
		"1980", "all-time",
	}

	// create table references for each year
	err := initDynamoTables(db, years)
	if err != nil {
		fmt.Println(err)
		return rankings, fmt.Errorf("GetRankingsFromDisk failed: %v", err)
	}

	for _, yr := range years {
		// get list of object IDs for the year,
		// create query from name
		names := createRankingsNames(yr)
		for _, n := range names {
			rpl := strings.ReplaceAll(n, "-", " ")
			ss := strings.Split(rpl, " ")
			prt := ss[1]
			query := &dynamo.Query{PrimaryValue: prt, SortValue: n}
			queries = append(queries, query)
		}
		for _, q := range queries {
			// get rankings list for the year
			odl, err := GetObjectFromDynamo(db, q, "top_overall", []string{yr})
			if err != nil {
				fmt.Println(err)
				return rankings, fmt.Errorf("GetRankingsFromDisk failed: %v", err)
			}
			for _, od := range odl {
				// add rankiings list to map
				if rankings[yr] == nil {
					rankings[yr] = make(map[string]donations.TopOverallData)
				}
				rankings[yr][od.(*donations.TopOverallData).ID] = *od.(*donations.TopOverallData)
			}
		}
	}
	return rankings, nil
}

// initDynamoDbDefault initializes a dynamo.DbInfo object with default DynamoDB session settings
func initDynamoDbDefault() (*dynamo.DbInfo, error) {
	// init DbInfo object and session
	db := dynamo.InitDbInfo()
	db.SetSvc(dynamo.InitSesh())
	db.SetFailConfig(dynamo.DefaultFailConfig)
	return db, nil
}

func initDynamoTables(db *dynamo.DbInfo, years []string) error {
	// get existing tables
	tables := make(map[string]bool)
	names, _, err := dynamo.ListTables(db.Svc)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("InitDynamoTables failed: %v", err)
	}

	// check tables for every given year exist
	for _, n := range names {
		tables[n] = true
	}
	for _, yr := range years {
		ts := initTableObjs(db, yr) // create in memory references to dynamo tables
		for _, t := range ts {
			if !tables[t.TableName] { // tables for given year do not exist
				err := createTables(db, ts) // create table set for year
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("initDynamoTables failed: %v", err)
				}
				continue
			}
		}
	}
	return nil
}

func createTables(db *dynamo.DbInfo, ts []*dynamo.Table) error {
	for _, t := range ts {
		err := dynamo.CreateTable(db.Svc, t)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("createTables failed: %v", err)
		}
	}
	return nil
}

// initTableObjs creates dynamo.Table objects for given year in memory only and
// adds them to the db.Tables field. See InitDynamoTables description for TableName format.
func initTableObjs(db *dynamo.DbInfo, year string) []*dynamo.Table {
	indv := "cf-" + year + "-individuals"        // pk = State
	cand := "cf-" + year + "-candidates"         // pk = State
	cmte := "cf-" + year + "-committees"         // pk = State
	cmteData := "cf-" + year + "-cmte_tx_data"   // pk = Name
	cmteFin := "cf-" + year + "-cmte_financials" // pk = Name
	topOverall := "cf-" + year + "-top_overall"  // pk = SizeLimit
	yrTotals := "cf-" + year + "-yearly_totals"  // pk = SizeLimit
	tables := []*dynamo.Table{}

	// create object tables
	t := dynamo.CreateNewTableObj(indv, "State", "string", "ID", "string")
	db.AddTable(t)
	tables = append(tables, t)

	// create object tables
	t = dynamo.CreateNewTableObj(cand, "State", "string", "ID", "string")
	db.AddTable(t)
	tables = append(tables, t)

	// create object tables
	t = dynamo.CreateNewTableObj(cmte, "State", "string", "ID", "string")
	db.AddTable(t)
	tables = append(tables, t)

	// create object tables
	t = dynamo.CreateNewTableObj(cmteData, "Name", "string", "ID", "string")
	db.AddTable(t)
	tables = append(tables, t)

	// create object tables
	t = dynamo.CreateNewTableObj(cmteFin, "Name", "string", "ID", "string")
	db.AddTable(t)
	tables = append(tables, t)

	// create TopOverall table
	t = dynamo.CreateNewTableObj(topOverall, "Category", "int", "ID", "string")
	db.AddTable(t)
	tables = append(tables, t)

	// create TopOverall table
	t = dynamo.CreateNewTableObj(yrTotals, "Category", "string", "ID", "string")
	db.AddTable(t)
	tables = append(tables, t)

	return tables
}

func getRefObj(bucket string) interface{} {
	var refObj interface{}
	switch {
	case bucket == "individuals":
		refObj = &donations.Individual{}
	case bucket == "candidates":
		refObj = &donations.Candidate{}
	case bucket == "committees":
		refObj = &donations.Committee{}
	case bucket == "cmte_tx_data":
		refObj = &donations.CmteTxData{}
	case bucket == "cmte_financials":
		refObj = &donations.CmteFinancials{}
	case bucket == "top_overall":
		refObj = &donations.TopOverallData{}
	default:
		refObj = nil
	}
	return refObj
}

func createRankingsNames(year string) []string {
	names := []string{}
	buckets := []string{"candidates", "cmte_tx_data", "individuals"}
	cats := []string{"rec", "donor", "exp"}
	ptys := []string{"ALL", ""}

	for _, b := range buckets {
		for _, c := range cats {
			for _, p := range ptys {
				n := year + "-" + b + "-" + c + "-" + p
				names = append(names, n)
			}
		}
	}

	return names
}
