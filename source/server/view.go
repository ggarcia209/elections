package server

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/elections/source/persist"

	"github.com/elections/source/donations"
	"github.com/elections/source/dynamo"
	"github.com/elections/source/indexing"
	"github.com/elections/source/util"
)

// RankingsMap stores references to each rankings list by year
type RankingsMap map[string]map[string]RankingsData

// YrTotalsMap stores references to each yearly total by year
type YrTotalsMap map[string]map[string]YrTotalData

// SearchDataMap stores references to SearchData objects
type SearchDataMap map[string]indexing.SearchData

// IndexData wraps and encapsulates the indexing.IndexData object
type IndexData indexing.IndexData

// InitServerDiskCache creates the ../db directory on the local disk
func InitServerDiskCache() {
	persist.InitDiskCache()
	indexing.OUTPUT_PATH = ".."
	fmt.Println("local disk cache created")
}

// InitDynamo initialized a Dynamo session with default settings
func InitDynamo() (*dynamo.DbInfo, error) {
	// init sesh and db with default options
	db, err := initDynamoDbDefault()
	if err != nil {
		fmt.Println(err)
		return db, fmt.Errorf("Upload failed: %v", err)
	}
	return db, nil
}

// GetIndexData retreives the encapsulated IndexData object from disk
func GetIndexData() (*IndexData, error) {
	id, err := indexing.GetIndexData()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetIndexData failed: %v", err)
	}
	wrap := &IndexData{
		TermsSize:      id.TermsSize,
		LookupSize:     id.LookupSize,
		LastUpdated:    id.LastUpdated,
		YearsCompleted: id.YearsCompleted,
		Completed:      id.Completed,
		Shards:         id.Shards,
	}
	return wrap, nil
}

// SearchData takes a user query as a string and
// finds the results matching each word in query.
func SearchData(id *IndexData, txt string) ([]string, error) {
	// get query from user / return & print results
	indexing.OUTPUT_PATH = "/Volumes/T7/processed" // CHANGE TO LOCAL DIR
	terms := formatTerms(strings.Split(txt, " "))
	wrap := &indexing.IndexData{
		Shards: make(indexing.ShardMap),
	}
	for _, t := range terms {
		wrap.Shards[t] = id.Shards[t]
	}

	q := indexing.CreateQuery(txt, "user")
	common, err := indexing.GetResultsFromShards(wrap, q)
	if err != nil {
		if err.Error() == "MAX_LENGTH" {
			return []string{}, err
		}
		if err.Error() == "NO_RESULTS" {
			fmt.Println("search data: ", err)
			return []string{}, err
		}
		fmt.Println(err)
		return []string{}, fmt.Errorf("QueryData failed: %v", err)
	}
	return common, nil
}

// GetSearchResults returns the SearchData object for the given IDs
func GetSearchResults(db *dynamo.DbInfo, ids []string, cache SearchDataMap) ([]indexing.SearchData, error) {
	indexing.OUTPUT_PATH = "/Volumes/T7/processed" // CHANGE TO LOCAL DIR
	nilIDs, frmCache := indexing.LookupSearchDataFromCache(ids, cache)
	frmDisk, err := indexing.LookupSearchDataFromDynamo(db, nilIDs) // refactor to read from DynamoDB
	if err != nil {
		fmt.Println(err)
		return []indexing.SearchData{}, fmt.Errorf("GetSearchResults failed: %v", err)
	}
	sds := indexing.ConsolidateSearchData(ids, frmCache, frmDisk)

	return sds, nil
}

// LookupByID finds an entity by ID
func LookupByID(IDs []string) ([]indexing.SearchData, error) {
	indexing.OUTPUT_PATH = "/Volumes/T7/processed" // CHANGE TO LOCAL DIR
	sds, err := indexing.LookupSearchData(IDs)
	if err != nil {
		fmt.Println(err)
		return []indexing.SearchData{}, fmt.Errorf("LookupByID failed: %v", err)
	}
	return sds, nil
}

// GetObjectFromDisk gets object from disk and returns pointer to obj as interface{}
func GetObjectFromDisk(year, ID, bucket string) (interface{}, error) {
	obj, err := persist.GetObject(year, bucket, ID)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("GetObjFromDisk failed: %v", err)
	}

	// re-encode object to intermediate type sent thru gRPC
	var intf interface{}
	switch bucket {
	case "individuals":
		indv := obj.(*donations.Individual)
		new := Individual{
			ID:            indv.ID,
			Name:          indv.Name,
			City:          indv.City,
			State:         indv.State,
			Occupation:    indv.Occupation,
			Employer:      indv.Employer,
			TotalOutAmt:   indv.TotalOutAmt,
			TotalOutTxs:   indv.TotalOutTxs,
			AvgTxOut:      indv.AvgTxOut,
			TotalInAmt:    indv.TotalInAmt,
			TotalInTxs:    indv.TotalInTxs,
			AvgTxIn:       indv.AvgTxIn,
			NetBalance:    indv.NetBalance,
			RecipientsAmt: indv.RecipientsAmt,
			RecipientsTxs: indv.RecipientsTxs,
			SendersAmt:    indv.SendersAmt,
			SendersTxs:    indv.SendersTxs,
		}
		intf = new
	case "committees":
		cmte := obj.(*donations.Committee)
		designation, cmteType, party := getCmteCodes(cmte.Designation, cmte.Type, cmte.Party)
		new := Committee{
			ID:           cmte.ID,
			Name:         cmte.Name,
			TresName:     cmte.TresName,
			City:         cmte.City,
			State:        cmte.State,
			Zip:          cmte.Zip,
			Designation:  designation,
			Type:         cmteType,
			Party:        party,
			FilingFreq:   cmte.FilingFreq,
			OrgType:      cmte.OrgType,
			ConnectedOrg: cmte.ConnectedOrg,
			CandID:       cmte.CandID,
		}
		intf = new
	case "cmte_tx_data":
		cmte := obj.(*donations.CmteTxData)
		new := CmteTxData{
			CmteID:                    cmte.CmteID,
			CandID:                    cmte.CandID,
			ContributionsInAmt:        cmte.ContributionsInAmt,
			ContributionsInTxs:        cmte.ContributionsInTxs,
			AvgContributionIn:         cmte.AvgContributionIn,
			OtherReceiptsInAmt:        cmte.OtherReceiptsInAmt,
			OtherReceiptsInTxs:        cmte.OtherReceiptsInTxs,
			AvgOtherIn:                cmte.AvgOtherIn,
			TotalIncomingAmt:          cmte.TotalIncomingAmt,
			TotalIncomingTxs:          cmte.TotalIncomingTxs,
			AvgIncoming:               cmte.AvgIncoming,
			TransfersAmt:              cmte.TransfersAmt,
			TransfersTxs:              cmte.TransfersTxs,
			AvgTransfer:               cmte.AvgTransfer,
			ExpendituresAmt:           cmte.ExpendituresAmt,
			ExpendituresTxs:           cmte.ExpendituresTxs,
			AvgExpenditure:            cmte.AvgExpenditure,
			TotalOutgoingAmt:          cmte.TotalOutgoingAmt,
			TotalOutgoingTxs:          cmte.TotalOutgoingTxs,
			AvgOutgoing:               cmte.AvgOutgoing,
			NetBalance:                cmte.NetBalance,
			TopIndvContributorsAmt:    cmte.TopIndvContributorsAmt,
			TopIndvContributorsTxs:    cmte.TopIndvContributorsTxs,
			TopCmteOrgContributorsAmt: cmte.TopCmteOrgContributorsAmt,
			TopCmteOrgContributorsTxs: cmte.TopCmteOrgContributorsTxs,
			TransferRecsAmt:           cmte.TransferRecsAmt,
			TransferRecsTxs:           cmte.TransferRecsTxs,
			TopExpRecipientsAmt:       cmte.TopExpRecipientsAmt,
			TopExpRecipientsTxs:       cmte.TopExpRecipientsTxs,
		}
		intf = new
	case "cmte_fin":
		cmte := obj.(*donations.CmteFinancials)
		new := CmteFinancials{
			CmteID:          cmte.CmteID,
			TotalReceipts:   cmte.TotalReceipts,
			TxsFromAff:      cmte.TxsFromAff,
			IndvConts:       cmte.IndvConts,
			OtherConts:      cmte.OtherConts,
			CandCont:        cmte.CandCont,
			TotalLoans:      cmte.TotalLoans,
			TotalDisb:       cmte.TotalDisb,
			TxToAff:         cmte.TxToAff,
			IndvRefunds:     cmte.IndvRefunds,
			OtherRefunds:    cmte.OtherRefunds,
			LoanRepay:       cmte.LoanRepay,
			CashBOP:         cmte.CashBOP,
			CashCOP:         cmte.CashCOP,
			DebtsOwed:       cmte.DebtsOwed,
			NonFedTxsRecvd:  cmte.NonFedTxsRecvd,
			ContToOtherCmte: cmte.ContToOtherCmte,
			IndExp:          cmte.IndExp,
			PartyExp:        cmte.PartyExp,
			NonFedSharedExp: cmte.NonFedSharedExp,
		}
		intf = new
	case "candidates":
		cand := obj.(*donations.Candidate)
		party, office := getCandCodes(cand.Party, cand.Office)
		new := Candidate{
			ID:                   cand.ID,
			Name:                 cand.Name,
			Party:                party,
			OfficeState:          cand.OfficeState,
			Office:               office,
			PCC:                  cand.PCC,
			City:                 cand.City,
			State:                cand.State,
			Zip:                  cand.Zip,
			OtherAffiliates:      cand.OtherAffiliates,
			TransactionsList:     cand.TransactionsList,
			TotalDirectInAmt:     cand.TotalDirectInAmt,
			TotalDirectInTxs:     cand.TotalDirectInTxs,
			AvgDirectIn:          cand.AvgDirectIn,
			TotalDirectOutAmt:    cand.TotalDirectOutAmt,
			TotalDirectOutTxs:    cand.TotalDirectOutTxs,
			AvgDirectOut:         cand.AvgDirectOut,
			NetBalanceDirectTx:   cand.NetBalanceDirectTx,
			DirectRecipientsAmts: cand.DirectRecipientsAmts,
			DirectRecipientsTxs:  cand.DirectRecipientsTxs,
			DirectSendersAmts:    cand.DirectSendersAmts,
			DirectSendersTxs:     cand.DirectSendersTxs,
		}
		intf = new
	case "cmpn_fin":
		cf := obj.(*donations.CmpnFinancials)
		new := CmpnFinancials{
			CandID:         cf.CandID,
			Name:           cf.Name,
			PartyCd:        cf.PartyCd,
			Party:          cf.Party,
			TotalReceipts:  cf.TotalReceipts,
			TransFrAuth:    cf.TransFrAuth,
			TotalDisbsmts:  cf.TotalDisbsmts,
			TransToAuth:    cf.TransToAuth,
			COHBOP:         cf.COHBOP,
			COHCOP:         cf.COHCOP,
			CandConts:      cf.CandConts,
			CandLoans:      cf.CandLoans,
			OtherLoans:     cf.OtherLoans,
			CandLoanRepay:  cf.CandLoanRepay,
			OtherLoanRepay: cf.OtherLoanRepay,
			DebtsOwedBy:    cf.DebtsOwedBy,
			TotalIndvConts: cf.TotalIndvConts,
			SpecElection:   cf.SpecElection,
			PrimElection:   cf.PrimElection,
			RunElection:    cf.RunElection,
			GenElection:    cf.GenElection,
			GenElectionPct: cf.GenElectionPct,
			OtherCmteConts: cf.OtherCmteConts,
			PtyConts:       cf.PtyConts,
			IndvRefunds:    cf.IndvRefunds,
			CmteRefunds:    cf.CmteRefunds,
		}
		intf = new
	default:
		fmt.Println("INVALID_TYPE")
		return nil, fmt.Errorf("INVALID_TYPE")
	}

	return intf, nil
}

// GetObjectFromDynamo returns the yearly datasets
// for the queried object and the given years
func GetObjectFromDynamo(db *dynamo.DbInfo, query *dynamo.Query, bucket string, years []string) ([]interface{}, error) {
	datasets := []interface{}{}

	// retreive item's datasets for each year from db
	refObj := getRefObj(bucket)
	for _, yr := range years {
		tName := "cf-" + yr + "-" + bucket
		if db.Tables[tName] == nil {
			fmt.Println("TABLE_NOT_FOUND")
			return nil, fmt.Errorf("TABLE_NOT_FOUND")
		}
		obj, err := dynamo.GetItem(db.Svc, query, db.Tables[tName], refObj)
		if err != nil {
			fmt.Println(err)
			return datasets, fmt.Errorf("GetObjectFromDynamo failed: %v", err)
		}
		wrap := wrapObj(obj)
		datasets = append(datasets, wrap)
	}

	return datasets, nil
}

// CreateQueryFromSearchData returns a Dynamo Query object from SearchData info
func CreateQueryFromSearchData(sd indexing.SearchData) *dynamo.Query {
	pk := ""

	if sd.Bucket == "cmte_tx_data" {
		pk = sd.Employer
	} else {
		pk = sd.State
	}

	return &dynamo.Query{PrimaryValue: pk, SortValue: sd.ID}
}

// GetRankingsFromDynamo retrieves the TopOvearll datasets
// for the given year from Dynamo to store in memory
// TEST ONLY
func GetRankingsFromDisk() (RankingsMap, error) {
	persist.OUTPUT_PATH = "/Volumes/T7/processed"
	rankings := make(RankingsMap)
	years := []string{"2020"}
	/*years := []string{
		"2020", "2018", "2016", "2014", "2012", "2010", "2008", "2006", "2004", "2002",
		"2000", "1998", "1996", "1994", "1992", "1990", "1988", "1986", "1984", "1982",
		"1980",
	} */

	for _, yr := range years {
		fmt.Printf("getting rankings for %s...\n", yr)
		// get list of object IDs for the year,
		odl, err := persist.GetTopOverall(yr)
		if err != nil {
			fmt.Println(err)
			return rankings, fmt.Errorf("GetRankingsFromDisk failed: %v", err)
		}
		for _, od := range odl {
			full := od.(*donations.TopOverallData)
			// add rankiings list to map
			if full.Bucket == "individuals" {
				// clip individuals lists to 500 entries
				sorted := util.SortMapObjectTotals(full.Amts)
				clip := make(map[string]float32)
				for i, e := range sorted {
					if i == 500 {
						break
					}
					clip[e.ID] = e.Total
				}
				full.Amts = clip
			}
			// add rankiings list to map
			if rankings[yr] == nil {
				rankings[yr] = make(map[string]RankingsData)
			}
			wrap := RankingsData{
				ID:       full.ID,
				Year:     full.Year,
				Bucket:   full.Bucket,
				Category: full.Category,
				Party:    full.Party,
				Rankings: full.Amts,
			}
			rankings[yr][wrap.ID] = wrap

			// create preview list and add preview object to map
			pre := createRankingsPreview(wrap)
			rankings[yr][pre.ID] = pre
		}
	}
	return rankings, nil
}

// GetRankingsFromDynamo retrieves the TopOvearll datasets
// for the given year from Dynamo to store in memory
func GetRankingsFromDynamo(db *dynamo.DbInfo) (RankingsMap, error) {
	persist.OUTPUT_PATH = "/Volumes/T7/processed"
	rankings := make(RankingsMap)
	years := []string{"2020"}
	/*years := []string{
		"2020", "2018", "2016", "2014", "2012", "2010", "2008", "2006", "2004", "2002",
		"2000", "1998", "1996", "1994", "1992", "1990", "1988", "1986", "1984", "1982",
		"1980",
	} */

	for _, yr := range years {
		fmt.Printf("getting rankings for %s...\n", yr)
		// get list of object IDs for the year,
		queries := []*dynamo.Query{}
		names := createRankingsNames(yr)
		for _, n := range names {
			ss := strings.Split(n, "-")
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
				full := od.(*donations.TopOverallData)
				// add rankiings list to map
				if full.Bucket == "individuals" {
					// clip individuals lists to 500 entries
					sorted := util.SortMapObjectTotals(full.Amts)
					clip := make(map[string]float32)
					for i, e := range sorted {
						if i == 500 {
							break
						}
						clip[e.ID] = e.Total
					}
					full.Amts = clip
				}
				// add rankiings list to map
				if rankings[yr] == nil {
					rankings[yr] = make(map[string]RankingsData)
				}
				wrap := RankingsData{
					ID:       full.ID,
					Year:     full.Year,
					Bucket:   full.Bucket,
					Category: full.Category,
					Party:    full.Party,
					Rankings: full.Amts,
				}
				rankings[yr][wrap.ID] = wrap

				// create preview list and add preview object to map
				pre := createRankingsPreview(wrap)
				rankings[yr][pre.ID] = pre
			}
		}
	}
	return rankings, nil
}

// GetRankingsFromDynamo retrieves the TopOvearll datasets
// for the given year from Dynamo to store in memory
// TEST ONLY
func GetYrTotalsFromDisk() (YrTotalsMap, error) {
	persist.OUTPUT_PATH = "/Volumes/T7/processed"
	totals := make(YrTotalsMap)
	years := []string{"2020"}
	/* years := []string{
		"2020", "2018", "2016", "2014", "2012", "2010", "2008", "2006", "2004", "2002",
		"2000", "1998", "1996", "1994", "1992", "1990", "1988", "1986", "1984", "1982",
		"1980",
	} */
	cats := []string{"rec", "donor", "exp"}

	for _, yr := range years {
		fmt.Printf("getting yearly totals for %s...\n", yr)
		for _, cat := range cats {
			// get list of object IDs for the year,
			ytl, err := persist.GetYearlyTotals(yr, cat)
			if err != nil {
				fmt.Println(err)
				return totals, fmt.Errorf("GetRankingsFromDisk failed: %v", err)
			}
			for _, yt := range ytl {
				// add rankiings list to map
				if totals[yr] == nil {
					totals[yr] = make(map[string]YrTotalData)
				}
				total := yt.(*donations.YearlyTotal)
				wrap := YrTotalData{
					ID:       total.ID,
					Year:     total.Year,
					Category: total.Category,
					Party:    total.Party,
					Total:    total.Total,
				}
				totals[yr][wrap.ID] = wrap
			}
		}
	}
	return totals, nil
}

// CreateSearchCache creates a cache of SearchData objects for every unique entity listed in rankings
func CreateSearchCache(rankings RankingsMap) (SearchDataMap, error) {
	indexing.OUTPUT_PATH = "/Volumes/T7/processed" // change to local dir
	cache := make(SearchDataMap)
	if len(rankings) == 0 {
		return cache, fmt.Errorf("CreateSearchCache failed: empty Rankings cache")
	}
	ids := make(map[string]bool)
	idList := []string{}

	for yr, set := range rankings {
		if len(set) == 0 {
			return cache, fmt.Errorf("CreateSearchCache failed: empty set in year %s", yr)
		}
		for id, data := range set {
			if len(data.Rankings) == 0 {
				fmt.Printf("empty obj in year %s - ID: %s\n", yr, id)
				continue
			}
			for objID := range data.Rankings {
				if !ids[objID] {
					ids[objID] = true
					idList = append(idList, objID)
				}
			}
		}
	}

	sds, err := indexing.LookupSearchData(idList)
	if err != nil {
		return cache, fmt.Errorf("CreateSearchCache failed: %v", err)
	}
	for _, sd := range sds {
		cache[sd.ID] = sd
	}
	return cache, nil
}

// initDynamoDbDefault initializes a dynamo.DbInfo object with default DynamoDB session settings
func initDynamoDbDefault() (*dynamo.DbInfo, error) {
	// init DbInfo object and session
	db := dynamo.InitDbInfo()
	db.SetSvc(dynamo.InitSesh())
	db.SetFailConfig(dynamo.DefaultFailConfig)

	years := []string{"2020"}

	for _, yr := range years {
		// create Table objects
		initTableObjs(db, yr)
	}

	err := initDynamoTables(db)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("initDynamoDbDefault failed: %v", err)
	}

	// list tables currently in DB
	_, t, err := dynamo.ListTables(db.Svc)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("initDynamoDbDefault failed: %v", err)
	}

	fmt.Println("Total tables: ", t)
	fmt.Println()

	return db, nil
}

// initDynamoTables initializes Tables for each object category for the given year
// and adds the corresponding Table object to the db.Tables field.
// TableName format: "cf-%s-individuals", year
//                   "cf-%s-candidates", year
//                   "cf-s-committees", year
//                   "cf-%s-cmte_tx_data", year
//                   "cf-%s-cmte_financials", year
//                   "cf-%s-top_ovearll", year
func initDynamoTables(db *dynamo.DbInfo) error {
	for _, t := range db.Tables {
		err := dynamo.CreateTable(db.Svc, t)
		if err != nil {
			if err.Error() == "ResourceInUseException" {
				fmt.Printf("table '%s' already exists\n", t.TableName)
				continue
			}
			fmt.Println(err)
			return fmt.Errorf("initDynamoTables failed: %v", err)
		}
		ss := strings.Split(t.TableName, "-")
		year := ""
		bucket := ""
		if len(ss) == 3 {
			year = ss[1]   // year
			bucket = ss[2] // object bucket
		} else {
			year = "all-time"
			bucket = ss[1] // index/lookup
		}

		// reset BatchGetSequential startKeys for new tables
		switch bucket {
		case "index":
			fmt.Println("resetting key for table: ", t.TableName)
			err := persist.LogKey(year, "index-partitions", "")
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("initDynamoTables failed: %v", err)
			}
		case "lookup":
			fmt.Println("resetting key for table: ", t.TableName)
			err := persist.LogKey(year, bucket, "")
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("initDynamoTables failed: %v", err)
			}
		default:
			fmt.Println("resetting key for table: ", t.TableName)
			err := persist.LogKey(year, bucket, "")
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("initDynamoTables failed: %v", err)
			}
		}
	}

	return nil
}

// initTableObjs creates dynamo.Table objects for given year in memory only and
// adds them to the db.Tables field. See InitDynamoTables description for TableName format.
func initTableObjs(db *dynamo.DbInfo, year string) {
	indv := "cf-" + year + "-individuals"      // pk = First Letter of Name
	cand := "cf-" + year + "-candidates"       // pk = First Letter of Name
	cmte := "cf-" + year + "-committees"       // pk = First Letter of Name
	cmteData := "cf-" + year + "-cmte_tx_data" // pk = First Letter of Name
	// cmteFin := "cf-" + year + "-cmte_financials" // pk = First Letter of Name
	topOverall := "cf-" + year + "-top_overall" // pk = Year
	yrTotals := "cf-" + year + "-yearly_totals" // pk = Year
	index := "cf-index"                         // pk = Index Partition
	lookup := "cf-lookup"                       // pk = First Letter of Name

	// create object tables
	t := dynamo.CreateNewTableObj(indv, "State", "string", "ID", "string")
	db.AddTable(t)

	// create object tables
	t = dynamo.CreateNewTableObj(cand, "State", "string", "ID", "string")
	db.AddTable(t)

	// create object tables
	t = dynamo.CreateNewTableObj(cmte, "State", "string", "ID", "string")
	db.AddTable(t)

	// create object tables
	t = dynamo.CreateNewTableObj(cmteData, "Party", "string", "CmteID", "string")
	db.AddTable(t)

	// create TopOverall table
	t = dynamo.CreateNewTableObj(topOverall, "Year", "string", "ID", "string")
	db.AddTable(t)

	// create YearlyTotals table
	t = dynamo.CreateNewTableObj(yrTotals, "Year", "string", "ID", "string")
	db.AddTable(t)

	// create Index table
	t = dynamo.CreateNewTableObj(index, "Partition", "string", "Term", "string")
	db.AddTable(t)

	// create Index table
	t = dynamo.CreateNewTableObj(lookup, "Partition", "string", "ID", "string")
	db.AddTable(t)

	return
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

func createRankingsPreview(full RankingsData) RankingsData {
	// create preview object
	pre := RankingsData{
		ID:       full.ID + "-pre",
		Year:     full.Year,
		Bucket:   full.Bucket,
		Category: full.Category,
		Party:    full.Party,
		Rankings: make(map[string]float32),
	}

	// sort map and get top 10 results
	sorted := util.SortMapObjectTotals(full.Rankings)
	for i, e := range sorted {
		if i == 10 {
			break
		}
		pre.Rankings[e.ID] = e.Total
	}
	return pre
}

func getCmteCodes(d, t, p string) (string, string, string) {
	var dsgn = map[string]string{
		"A": "Authorized by Candidate", "B": "Lobbyist/Registrant PAC", "D": "Leadership PAC",
		"J": "Joint Fundraiser", "P": "Principal Campaign", "U": "Unauthorized",
	}
	var types = map[string]string{
		"C": "Corporation", "L": "Labor Organization", "M": "Membership Organization",
		"T": "Trade Association", "V": "Cooperative", "W": "Corporation without Capital Stock",
	}
	var parties = map[string]string{
		"ACE": "Ace Party",
		"AKI": "Alaskan Independence Party",
		"AIC": "American Independent Conservative",
		"AIP": "American Independent Party",
		"AMP": "American Party",
		"APF": "American People's Freedom Party",
		"AE":  "Americans Elect",
		"CIT": "Citizens' Party",
		"CMD": "Commandments Party",
		"CMP": "Commonwealth Party of the U.S.",
		"COM": "Communist Party",
		"CNC": "Concerned Citizens Party Of Connecticut",
		"CRV": "Conservative Party",
		"CON": "Constitution Party",
		"CST": "Constitutional",
		"COU": "Country",
		"DCG": "D.C. Statehood Green Party",
		"DNL": "Democratic-Nonpartisan League",
		"DEM": "Democratic Party",
		"D/C": "Democratic/Conservative",
		"DFL": "Democratic-Farmer-Labor",
		"DGR": "Desert Green Party",
		"FED": "Federalist",
		"FLP": "Freedom Labor Party",
		"FRE": "Freedom Party",
		"GWP": "George Wallace Party",
		"GRT": "Grassroots",
		"GRE": "Green Party",
		"GR":  "Green-Rainbow",
		"HRP": "Human Rights Party",
		"IDP": "Independence Party",
		"IND": "Independent",
		"IAP": "Independent American Party",
		"ICD": "Independent Conservative Democratic",
		"IGR": "Independent Green",
		"IP":  "Independent Party",
		"IDE": "Indepenent Party of Delaware",
		"IGD": "Industrial Government Party",
		"JCN": "Jewish/Christian National",
		"JUS": "Justice Party",
		"LRU": "La Raza Unida",
		"LBR": "Labor Party",
		"LFT": "Less Federal Taxes",
		"LBL": "Liberal Party",
		"LIB": "Libertarian Party",
		"LBU": "Liberty Union Party",
		"MTP": "Mountain Party",
		"NDP": "National Democratic Party",
		"NLP": "Natural Law Party",
		"NA":  "New Alliance",
		"NJC": "New Jersey Conservative Party",
		"NPP": "New Progressive Party",
		"NPA": "No Party Affiliation",
		"NOP": "No Party Preference",
		"NNE": "None",
		"N":   "Nonpartisan",
		"NON": "Non-Party",
		"OE":  "One Earth Party",
		"OTH": "Other",
		"PG":  "Pacific Green",
		"PSL": "Party for Socialism and Liberation",
		"PAF": "Peace And Freedom",
		"PFP": "Peace And Freedom Party",
		"PFD": "Peace Freedom Party",
		"POP": "People Over Politics",
		"PPY": "People's Party",
		"PCH": "Personal Choice Party",
		"PPD": "Popular Democratic Party",
		"PRO": "Progressive Party",
		"NAP": "Prohibition Party",
		"PRI": "Puerto Rican Independence Party",
		"RUP": "Raza Unida Party",
		"REF": "Reform Party",
		"REP": "Republican Party",
		"RES": "Resource Party",
		"RTL": "Right To Life",
		"SEP": "Socialist Equality Party",
		"SLP": "Socialist Labor Party",
		"SUS": "Socialist Party",
		"SOC": "Socialist Party U.S.A.",
		"SWP": "Socialist Workers Party",
		"TX":  "Taxpayers",
		"TWR": "Taxpayers Without Representation",
		"TEA": "Tea Party",
		"THD": "Theo-Democratic",
		"LAB": "U.S. Labor Party",
		"USP": "U.S. People's Party",
		"UST": "U.S. Taxpayers Party",
		"UN":  "Unaffiliated",
		"UC":  "United Citizen",
		"UNI": "United Party",
		"UNK": "Unknown",
		"VET": "Veterans Party",
		"WTP": "We the People",
		"W":   "Write-In",
	}
	return dsgn[d], types[t], parties[p]
}

func getCandCodes(p, o string) (string, string) {
	var parties = map[string]string{
		"ACE": "Ace Party",
		"AKI": "Alaskan Independence Party",
		"AIC": "American Independent Conservative",
		"AIP": "American Independent Party",
		"AMP": "American Party",
		"APF": "American People's Freedom Party",
		"AE":  "Americans Elect",
		"CIT": "Citizens' Party",
		"CMD": "Commandments Party",
		"CMP": "Commonwealth Party of the U.S.",
		"COM": "Communist Party",
		"CNC": "Concerned Citizens Party Of Connecticut",
		"CRV": "Conservative Party",
		"CON": "Constitution Party",
		"CST": "Constitutional",
		"COU": "Country",
		"DCG": "D.C. Statehood Green Party",
		"DNL": "Democratic-Nonpartisan League",
		"DEM": "Democratic Party",
		"D/C": "Democratic/Conservative",
		"DFL": "Democratic-Farmer-Labor",
		"DGR": "Desert Green Party",
		"FED": "Federalist",
		"FLP": "Freedom Labor Party",
		"FRE": "Freedom Party",
		"GWP": "George Wallace Party",
		"GRT": "Grassroots",
		"GRE": "Green Party",
		"GR":  "Green-Rainbow",
		"HRP": "Human Rights Party",
		"IDP": "Independence Party",
		"IND": "Independent",
		"IAP": "Independent American Party",
		"ICD": "Independent Conservative Democratic",
		"IGR": "Independent Green",
		"IP":  "Independent Party",
		"IDE": "Indepenent Party of Delaware",
		"IGD": "Industrial Government Party",
		"JCN": "Jewish/Christian National",
		"JUS": "Justice Party",
		"LRU": "La Raza Unida",
		"LBR": "Labor Party",
		"LFT": "Less Federal Taxes",
		"LBL": "Liberal Party",
		"LIB": "Libertarian Party",
		"LBU": "Liberty Union Party",
		"MTP": "Mountain Party",
		"NDP": "National Democratic Party",
		"NLP": "Natural Law Party",
		"NA":  "New Alliance",
		"NJC": "New Jersey Conservative Party",
		"NPP": "New Progressive Party",
		"NPA": "No Party Affiliation",
		"NOP": "No Party Preference",
		"NNE": "None",
		"N":   "Nonpartisan",
		"NON": "Non-Party",
		"OE":  "One Earth Party",
		"OTH": "Other",
		"PG":  "Pacific Green",
		"PSL": "Party for Socialism and Liberation",
		"PAF": "Peace And Freedom",
		"PFP": "Peace And Freedom Party",
		"PFD": "Peace Freedom Party",
		"POP": "People Over Politics",
		"PPY": "People's Party",
		"PCH": "Personal Choice Party",
		"PPD": "Popular Democratic Party",
		"PRO": "Progressive Party",
		"NAP": "Prohibition Party",
		"PRI": "Puerto Rican Independence Party",
		"RUP": "Raza Unida Party",
		"REF": "Reform Party",
		"REP": "Republican Party",
		"RES": "Resource Party",
		"RTL": "Right To Life",
		"SEP": "Socialist Equality Party",
		"SLP": "Socialist Labor Party",
		"SUS": "Socialist Party",
		"SOC": "Socialist Party U.S.A.",
		"SWP": "Socialist Workers Party",
		"TX":  "Taxpayers",
		"TWR": "Taxpayers Without Representation",
		"TEA": "Tea Party",
		"THD": "Theo-Democratic",
		"LAB": "U.S. Labor Party",
		"USP": "U.S. People's Party",
		"UST": "U.S. Taxpayers Party",
		"UN":  "Unaffiliated",
		"UC":  "United Citizen",
		"UNI": "United Party",
		"UNK": "Unknown",
		"VET": "Veterans Party",
		"WTP": "We the People",
		"W":   "Write-In",
	}
	var office = map[string]string{
		"H": "House",
		"S": "Senate",
		"P": "President",
	}
	return parties[p], office[o]
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

// filter generic terms & edge cases ("the", "for", "of", "",)
// returns true if term meets filter criteria
func filter(term string) bool {
	f := map[string]bool{
		"for":  true,
		"the":  true,
		"of":   true,
		"":     true,
		"and":  true,
		"mr":   true,
		"mr.":  true,
		"mrs":  true,
		"mrs.": true,
		"ms.":  true,
		"to":   true,
	}
	return f[term]
}

// wrap object from package donations
func wrapObj(obj interface{}) interface{} {
	var wrap interface{}
	switch t := obj.(type) {
	case *donations.Individual:
		w := Individual{
			ID:            obj.(*donations.Individual).ID,
			Name:          obj.(*donations.Individual).Name,
			City:          obj.(*donations.Individual).City,
			State:         obj.(*donations.Individual).State,
			Zip:           obj.(*donations.Individual).Zip,
			Occupation:    obj.(*donations.Individual).Occupation,
			TotalOutAmt:   obj.(*donations.Individual).TotalOutAmt,
			TotalOutTxs:   obj.(*donations.Individual).TotalOutTxs,
			AvgTxOut:      obj.(*donations.Individual).AvgTxOut,
			TotalInAmt:    obj.(*donations.Individual).TotalInAmt,
			TotalInTxs:    obj.(*donations.Individual).TotalInTxs,
			AvgTxIn:       obj.(*donations.Individual).AvgTxIn,
			NetBalance:    obj.(*donations.Individual).NetBalance,
			RecipientsAmt: obj.(*donations.Individual).RecipientsAmt,
			RecipientsTxs: obj.(*donations.Individual).RecipientsTxs,
			SendersAmt:    obj.(*donations.Individual).SendersAmt,
			SendersTxs:    obj.(*donations.Individual).SendersTxs,
		}
		wrap = w
	case *donations.Committee:
		w := Committee{
			ID:           obj.(*donations.Committee).ID,
			Name:         obj.(*donations.Committee).Name,
			TresName:     obj.(*donations.Committee).TresName,
			City:         obj.(*donations.Committee).City,
			State:        obj.(*donations.Committee).State,
			Designation:  obj.(*donations.Committee).Designation,
			Type:         obj.(*donations.Committee).Type,
			Party:        obj.(*donations.Committee).Party,
			FilingFreq:   obj.(*donations.Committee).FilingFreq,
			OrgType:      obj.(*donations.Committee).OrgType,
			ConnectedOrg: obj.(*donations.Committee).ConnectedOrg,
			CandID:       obj.(*donations.Committee).CandID,
		}
		wrap = w
	case *donations.CmteTxData:
		w := CmteTxData{
			CmteID:                    obj.(*donations.CmteTxData).CmteID,
			CandID:                    obj.(*donations.CmteTxData).CandID,
			Party:                     obj.(*donations.CmteTxData).Party,
			ContributionsInAmt:        obj.(*donations.CmteTxData).ContributionsInAmt,
			ContributionsInTxs:        obj.(*donations.CmteTxData).ContributionsInTxs,
			AvgContributionIn:         obj.(*donations.CmteTxData).AvgContributionIn,
			OtherReceiptsInAmt:        obj.(*donations.CmteTxData).OtherReceiptsInAmt,
			OtherReceiptsInTxs:        obj.(*donations.CmteTxData).OtherReceiptsInTxs,
			AvgOtherIn:                obj.(*donations.CmteTxData).AvgOtherIn,
			TotalIncomingAmt:          obj.(*donations.CmteTxData).TotalIncomingAmt,
			TotalIncomingTxs:          obj.(*donations.CmteTxData).TotalIncomingTxs,
			AvgIncoming:               obj.(*donations.CmteTxData).AvgIncoming,
			TransfersAmt:              obj.(*donations.CmteTxData).TransfersAmt,
			TransfersTxs:              obj.(*donations.CmteTxData).TransfersTxs,
			AvgTransfer:               obj.(*donations.CmteTxData).AvgTransfer,
			ExpendituresAmt:           obj.(*donations.CmteTxData).ExpendituresAmt,
			ExpendituresTxs:           obj.(*donations.CmteTxData).ExpendituresTxs,
			AvgExpenditure:            obj.(*donations.CmteTxData).AvgExpenditure,
			TotalOutgoingAmt:          obj.(*donations.CmteTxData).TotalOutgoingAmt,
			TotalOutgoingTxs:          obj.(*donations.CmteTxData).TotalOutgoingTxs,
			AvgOutgoing:               obj.(*donations.CmteTxData).AvgOutgoing,
			NetBalance:                obj.(*donations.CmteTxData).NetBalance,
			TopIndvContributorsAmt:    obj.(*donations.CmteTxData).TopIndvContributorsAmt,
			TopIndvContributorsTxs:    obj.(*donations.CmteTxData).TopIndvContributorsTxs,
			TopCmteOrgContributorsAmt: obj.(*donations.CmteTxData).TopCmteOrgContributorsAmt,
			TopCmteOrgContributorsTxs: obj.(*donations.CmteTxData).TopCmteOrgContributorsTxs,
			TransferRecsAmt:           obj.(*donations.CmteTxData).TransferRecsAmt,
			TransferRecsTxs:           obj.(*donations.CmteTxData).TransferRecsTxs,
			TopExpRecipientsAmt:       obj.(*donations.CmteTxData).TopExpRecipientsAmt,
			TopExpRecipientsTxs:       obj.(*donations.CmteTxData).TopExpRecipientsTxs,
		}
		wrap = w
	case *donations.Candidate:
		w := Candidate{
			ID:                   obj.(*donations.Candidate).ID,
			Name:                 obj.(*donations.Candidate).Name,
			Party:                obj.(*donations.Candidate).Party,
			ElectnYr:             obj.(*donations.Candidate).ElectnYr,
			OfficeState:          obj.(*donations.Candidate).OfficeState,
			Office:               obj.(*donations.Candidate).Office,
			PCC:                  obj.(*donations.Candidate).PCC,
			City:                 obj.(*donations.Candidate).City,
			State:                obj.(*donations.Candidate).State,
			TotalDirectInAmt:     obj.(*donations.Candidate).TotalDirectInAmt,
			TotalDirectInTxs:     obj.(*donations.Candidate).TotalDirectInTxs,
			AvgDirectIn:          obj.(*donations.Candidate).AvgDirectIn,
			TotalDirectOutAmt:    obj.(*donations.Candidate).TotalDirectOutAmt,
			TotalDirectOutTxs:    obj.(*donations.Candidate).TotalDirectOutTxs,
			AvgDirectOut:         obj.(*donations.Candidate).AvgDirectOut,
			NetBalanceDirectTx:   obj.(*donations.Candidate).NetBalanceDirectTx,
			DirectRecipientsAmts: obj.(*donations.Candidate).DirectRecipientsAmts,
			DirectRecipientsTxs:  obj.(*donations.Candidate).DirectRecipientsTxs,
			DirectSendersAmts:    obj.(*donations.Candidate).DirectSendersAmts,
			DirectSendersTxs:     obj.(*donations.Candidate).DirectSendersTxs,
		}
		wrap = w
	default:
		_ = t
		wrap = nil
	}
	return wrap
}
