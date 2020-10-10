package admin

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/elections/source/donations"
	"github.com/elections/source/indexing"
	"github.com/elections/source/ui"
	"github.com/elections/source/util"

	"github.com/elections/source/dynamo"
	"github.com/elections/source/persist"
)

// Upload uploads the user-input year/category to DynamoDB
func Upload() error {
	path, err := persist.GetPath(false)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Upload failed: %v", err)
	}
	if path != "" {
		path, err = getPath(false)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Upload failed: %v", err)
		}
	}
	persist.OUTPUT_PATH = path
	indexing.OUTPUT_PATH = path

	opts := []string{"individuals", "committees", "cmte_tx_data", "candidates", "top_overall", "yearly_totals", "all", "index", "lookup", "Return"}
	menu := ui.CreateMenu("admin-upload-category", opts)

	fmt.Println("Choose year: ")
	year := ui.GetYear()

	// init sesh and db with default options
	db, err := initDynamoDbDefault(year)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Upload failed: %v", err)
	}

	for {
		fmt.Println("Choose a category: ")
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Upload failed: %v", err)
		}
		cat := menu.OptionsMap[ch]

		switch cat {
		case "Return":
			fmt.Println("Returning to menu...")
			return nil
		case "all":
			// upload all categories for given year; return when complete
			for _, cat := range opts[:7] {
				err := uploadFromDisk(db, year, cat, 1000)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("Upload failed: %v", err)
				}
			}
			fmt.Printf("Year %s uploaded. Returning to menu...\n", year)
			return nil
		case "index":
			total := 0
			new := false
			// get partiton map; sort
			pm, err := indexing.GetPartitionMap()
			if err != nil {
				fmt.Println(err)
				fmt.Println("total items scanned: ", total)
				return fmt.Errorf("Upload failed: %v", err)
			}
			prtSrt := util.SortCheckMap(pm)

			// get partition start key
			startPrt, err := persist.GetKey("all-time", "index-partitions")
			if err != nil {
				fmt.Println(err)
				fmt.Println("total items scanned: ", total)
				return fmt.Errorf("Upload failed: %v", err)
			}
			fmt.Println("partition map")
			fmt.Println(prtSrt)
			fmt.Println("starting partition: ", startPrt)
			if startPrt == "" {
				new = true
				startPrt = prtSrt[0].Key
			}

			// upload each partition
			for _, prt := range prtSrt {
				if startPrt != prtSrt[0].Key && prt.Key <= startPrt {
					continue // skip if partiton already uploaded
				}
				fmt.Println("uploading partition: ", prt)
				ct, err := uploadIndex(db, prt.Key, 1000, new)
				if err != nil {
					fmt.Println(err)
					fmt.Println("total items scanned: ", total)
					return fmt.Errorf("Upload failed: %v", err)
				}
				// log completed partition
				err = persist.LogKey("all-time", "index-partitions", prt.Key)
				if err != nil {
					fmt.Println(err)
					fmt.Println("total items scanned: ", total)
					return fmt.Errorf("Upload failed: %v", err)
				}
				total += ct
			}
			// reset once all partitions uploaded complete
			err = persist.LogKey("all-time", "index-partitions", "")
			if err != nil {
				fmt.Println(err)
				fmt.Println("total items scanned: ", total)
				return fmt.Errorf("Upload failed: %v", err)
			}
			fmt.Println("total items scanned: ", total)
		case "lookup":
			ct, err := uploadIndex(db, "lookup", 1000, false)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Upload failed: %v", err)
			}
			fmt.Println("total items scanned: ", ct)
		default:
			// upload single category
			err = uploadFromDisk(db, year, cat, 1000)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Upload failed: %v", err)
			}
		}

		// coninue/return
		fmt.Printf("year %s - %s uploaded. Continue?\n", year, cat)
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

// QueryDynamoDB retreives an object from DynamoDB per the specified input
func QueryDynamoDB() error {
	// get dataset year
	year := ui.GetYear()
	// init sesh and db with default options
	db, err := initDynamoDbDefault(year)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Query failed: %v", err)
	}
	var refObj interface{}

	tables := []string{}
	for t := range db.Tables {
		tables = append(tables, t)
	}

	menu := ui.CreateMenu("dynamo-query-tables", tables)
	for {
		// get table
		fmt.Println("Select a DynamoDB table to query: ")
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Query failed: %v", err)
		}
		tName := menu.OptionsMap[ch]

		// determine reference object type
		tName = strings.ReplaceAll(tName, "-", " ")
		ss := strings.Split(tName, " ")
		ty := strings.TrimSpace(ss[2])
		yr := strings.TrimSpace(ss[1])
		fmt.Printf("year - category: %s - %s\n", yr, ty)
		switch {
		case ty == "individuals":
			refObj = &donations.Individual{}
		case ty == "candidates":
			refObj = &donations.Candidate{}
		case ty == "committees":
			refObj = &donations.Committee{}
		case ty == "cmte_tx_data":
			refObj = &donations.CmteTxData{}
		case ty == "cmte_financials":
			refObj = &donations.CmteFinancials{}
		case ty == "top_overall":
			refObj = &donations.TopOverallData{}
		}

		// get partition/sort keys
		q := ui.GetDynamoQuery()
		query := &dynamo.Query{}
		for k, v := range q {
			query = dynamo.CreateNewQueryObj(k, v)
		}

		// retreive item from DB
		obj, err := dynamo.GetItem(db.Svc, query, db.Tables[tName], refObj)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Query failed: %v", err)
		}
		fmt.Printf("%#v\n", obj)
		// err = dynamo.DeleteTable(db, db.Tables[tName])
		fmt.Println()
		fmt.Println("New query?")
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

// DeleteDynamoTable provides options for Deleting DynamoDB tables
func DeleteDynamoTable() error {
	opts := []string{"Delete DynamoDB Table by Year", "Delete Table - Manual Input", "Return"}
	menu := ui.CreateMenu("delete-dynamo-table-main", opts)
	for {
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("DeleteDynamoTable failed: %v", err)
		}
		switch {
		case menu.OptionsMap[ch] == "Delete DynamoDB Table by Year":
			err := deleteTableByYr()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("DeleteDynamoTable failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Delete Table - Manual Input":
			err := deleteTableManual()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("DeleteDynamoTable failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Return":
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

// DeleteTable deletes a specified table from DynamoDB
func deleteTableByYr() error {
	year := ui.GetYear()
	db, err := initDynamoDbDefault(year)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("deleteTableByYr failed: %v", err)
	}
	fmt.Printf("Are you sure you want to delete the all tables for Year %s?\n", year)
	fmt.Println("Index and Lookup tables will NOT be deleted - use manual deete.")
	fmt.Println("StartKeys for each table in the given year will be reset.")
	yes := ui.Ask4confirm()
	if !yes {
		fmt.Println("Returning to menu...")
		return nil
	}
	for _, t := range db.Tables {
		if t.TableName == "cf-index" || t.TableName == "cf-lookup" {
			continue
		}
		err := dynamo.DeleteTable(db.Svc, t)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("deleteTableByYr failed: %v", err)
		}
		ss := strings.Split(t.TableName, "-")
		bucket := ss[2]
		persist.LogKey(year, bucket, "")
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("deleteTableByYr failed: %v", err)
		}
	}
	fmt.Println("Returning to menu...")
	return nil
}

func deleteTableManual() error {
	fmt.Println("*** Manual Table Delete ***")
	fmt.Println("Choose any year to intitialize the DynamoDB session.")
	fmt.Println("Upload start key for corresponding table will be reset.")
	year := ui.GetYear()
	db, err := initDynamoDbDefault(year)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("deleteTableManual failed: %v", err)
	}

	names, _, err := dynamo.ListTables(db.Svc)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("deleteTableManual failed: %v", err)
	}
	nameMap := make(map[string]bool)
	for _, n := range names {
		nameMap[n] = true
	}
	for {
		fmt.Println("Select database to delete.")
		q := ui.GetQuery()
		if !nameMap[q] {
			fmt.Println("Invalid Table name - Please try again.")
			continue
		}
		fmt.Printf("Are you sure you want to delete the Table %s?\n", q)
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
		t := dynamo.CreateNewTableObj(q, "", "", "", "")
		err := dynamo.DeleteTable(db.Svc, t)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("deleteTableManual failed: %v", err)
		}
		ss := strings.Split(t.TableName, "-")
		bucket := ""
		if len(ss) == 3 {
			bucket = ss[2]
		} else {
			year = "all-time"
			bucket = ss[1]
		}
		persist.LogKey(year, bucket, "")
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("deleteTableManual failed: %v", err)
		}
		fmt.Println("Delete another Table?")
		yes = ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

// initDynamoDbDefault initializes a dynamo.DbInfo object with default DynamoDB session settings
func initDynamoDbDefault(year string) (*dynamo.DbInfo, error) {
	// init DbInfo object and session
	db := dynamo.InitDbInfo()
	db.SetSvc(dynamo.InitSesh())
	db.SetFailConfig(dynamo.DefaultFailConfig)

	// create Table objects
	initTableObjs(db, year)

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

// uploadFromDisk intitializes the batch upload process to DynamoDB for the specified year: bucket
// Each call uploads n items to the bucket's correspondi
func uploadFromDisk(db *dynamo.DbInfo, year, bucket string, n int) error {
	i := 0
	final := false
	maxRetries := 10
	fmt.Printf("starting upload for %s - %s\n", year, bucket)

	/*
		// upload Top 100,000 individuals only
		if bucket == "individuals" {
			err := uploadTopIndv(db, year, bucket, n)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("UploadFromDisk failed: %v", err)
			}
			return nil

		} */

	startKey, err := persist.GetKey(year, bucket)
	fmt.Println("startKey: ", startKey)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("UploadFromDisk failed: %v", err)
	}

	for {
		// get next batch of objects
		objs, currKey, err := persist.BatchGetSequential(year, bucket, startKey, n)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("UploadFromDisk failed: %v", err)
		}

		if len(objs) == 0 {
			break
		}

		// check if Top Overall Data, clip Amts maps to 500 entries
		if bucket == "top_overall" {
			for _, obj := range objs {
				od := obj.(*donations.TopOverallData)
				if len(od.Amts) > 500 {
					full := od.Amts
					// clip individuals lists to 500 entries
					sorted := util.SortMapObjectTotals(full)
					clip := make(map[string]float32)
					for i, e := range sorted {
						if i == 500 {
							break
						}
						clip[e.ID] = e.Total
					}
					od.Amts = clip
				}
			}
		}

		// batch write n returned objects, 25 (max) per iteration
		retries := 0
		for {
			if len(objs) == 0 {
				break
			}
			if len(objs) < 25 { // final batch write
				tn := getTableName(year, bucket)
				err := dynamo.BatchWriteCreate(db.Svc, db.Tables[tn], db.FailConfig, objs)
				if err != nil {
					if awsErr, ok := err.(awserr.Error); ok {
						if awsErr.Code() == "RequestError" {
							// wait and retry
							fmt.Println("Request failed - retrying...")
							time.Sleep(250 * time.Millisecond)
							retries++
							if retries > maxRetries {
								msg := "MAX_RETRIES_EXCEEDED"
								fmt.Println(msg)
								return fmt.Errorf(msg)
							}
							continue
						}
						fmt.Println("uploadFromDisk failed:", awsErr.Code(), awsErr.Message())
					} else {
						fmt.Println(err.Error())
						return fmt.Errorf("uploadFromDisk failed: %v", err)
					}
				}
				i += len(objs)
				final = true
				retries = 0
				break
			}

			// batch write 25 objects from stack
			tn := getTableName(year, bucket)
			data := objs[len(objs)-25:]
			err = dynamo.BatchWriteCreate(db.Svc, db.Tables[tn], db.FailConfig, data)
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok {
					if awsErr.Code() == "RequestError" {
						// wait and retry
						retries++
						if retries > maxRetries {
							msg := "MAX_RETRIES_EXCEEDED"
							fmt.Println(msg)
							return fmt.Errorf(msg)
						}
						fmt.Println("Request failed - retrying...")
						time.Sleep(250 * time.Millisecond)
						continue
					}
					fmt.Println("uploadFromDisk failed:", awsErr.Code(), awsErr.Message())
				} else {
					fmt.Println(err.Error())
					return fmt.Errorf("uploadFromDisk failed: %v", err)
				}
			}

			// remove uploaded objects from stack
			objs = objs[:len(objs)-25]
			i += 25
			retries = 0
			fmt.Println("items scanned: ", i)
		}

		// update startKey & log currKey value
		startKey = currKey
		err = persist.LogKey(year, bucket, currKey)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("UploadFromDisk failed: %v", err)
		}
		fmt.Println("items scanned: ", i)
		if final == true {
			break
		}
	}

	// reset key log for next call
	err = persist.LogKey(year, bucket, "")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("UploadFromDisk failed: %v", err)
	}

	fmt.Println("***** UPLOAD FINSIHED *****")
	tn := getTableName(year, bucket)
	fmt.Printf("wrote %d items to table %s\n", i, tn)
	fmt.Println()

	return nil
}

func uploadIndex(db *dynamo.DbInfo, bucket string, n int, new bool) (int, error) {
	fmt.Println("uploading index...")
	i := 0
	fmt.Println("starting upload for index")
	year := "all-time"
	final := false
	maxRetries := 10
	var err error

	startKey := ""
	if new != true {
		startKey, err = persist.GetKey(year, bucket)
		if err != nil {
			fmt.Println(err)
			return i, fmt.Errorf("UploadFromDisk failed: %v", err)
		}
	}
	fmt.Println("start key: ", startKey)

	for {
		// get next batch of objects
		objs, currKey, err := indexing.BatchGetSequential(bucket, startKey, n)
		if err != nil {
			fmt.Println(err)
			return i, fmt.Errorf("UploadFromDisk failed: %v", err)
		}

		if len(objs) == 0 {
			break
		}

		// batch write n returned objects, 25 (max) per iteration
		fmt.Println("writing objects to database...")
		retries := 0
		for {
			if len(objs) == 0 {
				break
			}
			if len(objs) < 25 { // final batch write
				tn := getTableName(year, bucket)
				if tn == "" {
					tn = "cf-index"
				}
				err := dynamo.BatchWriteCreate(db.Svc, db.Tables[tn], db.FailConfig, objs)
				if err != nil {
					if awsErr, ok := err.(awserr.Error); ok {
						if awsErr.Code() == "RequestError" {
							// wait and retry
							fmt.Println("Request failed - retrying...")
							time.Sleep(250 * time.Millisecond)
							retries++
							if retries > maxRetries {
								msg := "MAX_RETRIES_EXCEEDED"
								fmt.Println(msg)
								return i, fmt.Errorf(msg)
							}
							continue
						}
						fmt.Println("uploadFromDisk failed:", awsErr.Code(), awsErr.Message())
					} else {
						fmt.Println(err.Error())
						return i, fmt.Errorf("uploadFromDisk failed: %v", err)
					}
				}
				i += len(objs)
				fmt.Println("items scanned: ", i)
				retries = 0
				final = true
				break
			}

			// batch write 25 objects from stack
			tn := getTableName(year, bucket)
			if tn == "" {
				tn = "cf-index"
			}
			data := objs[len(objs)-25:]
			err := dynamo.BatchWriteCreate(db.Svc, db.Tables[tn], db.FailConfig, data)
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok {
					if awsErr.Code() == "RequestError" {
						// wait and retry
						fmt.Println("Request failed - retrying...")
						time.Sleep(250 * time.Millisecond)
						retries++
						if retries > maxRetries {
							msg := "MAX_RETRIES_EXCEEDED"
							fmt.Println(msg)
							return i, fmt.Errorf(msg)
						}
						continue
					}
					fmt.Println("uploadFromDisk failed:", awsErr.Code(), awsErr.Message())
				} else {
					fmt.Println(err.Error())
					return i, fmt.Errorf("uploadFromDisk failed: %v", err)
				}
			}

			// remove uploaded objects from stack
			objs = objs[:len(objs)-25]
			i += 25
			retries = 0
			fmt.Println("items scanned: ", i)
		}

		// update startKey & log currKey value
		startKey = currKey

		err = persist.LogKey(year, bucket, currKey)
		if err != nil {
			fmt.Println(err)
			return i, fmt.Errorf("UploadFromDisk failed: %v", err)
		}

		fmt.Println("items scanned: ", i)

		// last batch of objects wrote to table
		if final == true {
			break
		}

	}

	// reset key log for next call
	err = persist.LogKey(year, bucket, "")
	if err != nil {
		fmt.Println(err)
		return i, fmt.Errorf("UploadFromDisk failed: %v", err)
	}

	fmt.Println("***** UPLOAD FINSIHED *****")
	tn := getTableName(year, bucket)
	fmt.Printf("wrote %d items to table %s\n", i, tn)
	fmt.Println()

	return i, nil

}

// initTableObjs creates dynamo.Table objects for given year in memory only and
// adds them to the db.Tables field. See InitDynamoTables description for TableName format.
func initTableObjs(db *dynamo.DbInfo, year string) {
	indv := "cf-" + year + "-individuals"      // pk = State
	cand := "cf-" + year + "-candidates"       // pk = State
	cmte := "cf-" + year + "-committees"       // pk = State
	cmteData := "cf-" + year + "-cmte_tx_data" // pk = Party
	// cmteFin := "cf-" + year + "-cmte_financials" // pk = First Letter of Name
	topOverall := "cf-" + year + "-top_overall" // pk = Year
	yrTotals := "cf-" + year + "-yearly_totals" // pk = Year
	index := "cf-index"                         // pk = Index Partition + shard number
	lookup := "cf-lookup"                       // pk = truncated ID (first 2 chars hash ID / last 2 chars FEC ID)

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

func uploadTopIndv(db *dynamo.DbInfo, year, bucket string, n int) error {
	ids := []string{}
	i := 0

	// get Top Individuals by incoming & outgoing funds
	topIndv, err := persist.GetObject(year, "top_overall", "indv")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("UploadFromDisk failed: %v", err)
	}
	topIndvRec, err := persist.GetObject(year, "top_overall", "indv_rec")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("UploadFromDisk failed: %v", err)
	}

	// create list of IDs for BatchGetByID
	for k := range topIndv.(*donations.TopOverallData).Amts {
		ids = append(ids, k)
	}
	for k := range topIndvRec.(*donations.TopOverallData).Amts {
		ids = append(ids, k)
	}

	for {
		// pop n IDs from stack and return corresponding objects
		if len(ids) < n { // queue exhausted - last write
			n = len(ids) // set starting index to 0
		}
		objs, _, err := persist.BatchGetByID(year, bucket, ids[len(ids)-n:])
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("UploadFromDisk failed: %v", err)
		}
		if len(objs) == 0 {
			break
		}

		// batch write n returned objects, 25 (max) per iteration
		for {
			if len(objs) == 0 {
				break
			}
			if len(objs) < 25 { // final batch write
				tn := getTableName(year, bucket)
				err := dynamo.BatchWriteCreate(db.Svc, db.Tables[tn], db.FailConfig, objs)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("UploadFromDisk failed: %v", err)
				}
				i += len(objs)
				break
			}
			fmt.Println(objs)

			// batch write 25 objects from stack
			data := objs[len(objs)-25:]
			err := dynamo.BatchWriteCreate(db.Svc, db.Tables[bucket], db.FailConfig, data)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("UploadFromDisk failed: %v", err)
			}

			// remove uploaded objects from stack
			objs = objs[:len(objs)-25]

			i += 25
		}
		// remove processed IDs from stack
		ids = ids[:len(ids)-n]
		fmt.Println("objects wrote to table: ", i)
	}

	tn := getTableName(year, bucket)
	fmt.Println("***** UPLOAD FINSIHED *****")
	fmt.Printf("wrote %d items to table %s\n", i, tn)
	fmt.Println()

	return nil
}

func getTableName(year, bucket string) string {
	tables := map[string]string{
		"individuals":   "cf-" + year + "-individuals",
		"candidates":    "cf-" + year + "-candidates",
		"committees":    "cf-" + year + "-committees",
		"cmte_tx_data":  "cf-" + year + "-cmte_tx_data",
		"cmte_fin":      "cf-" + year + "-cmte_financials",
		"top_overall":   "cf-" + year + "-top_overall",
		"yearly_totals": "cf-" + year + "-yearly_totals",
		"index":         "cf-index",
		"lookup":        "cf-lookup",
	}
	return tables[bucket]
}
