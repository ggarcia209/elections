package admin

import (
	"fmt"

	"github.com/elections/source/donations"
	"github.com/elections/source/util"

	"github.com/elections/source/dynamo"
	"github.com/elections/source/persist"
)

// Upload uploads the user-input year/category to DynamoDB
func Upload() error {
	opts := []string{"individuals", "committees", "candidates", "top_overall", "all", "Return"}
	menu := util.CreateMenu("admin - upload - category", opts)

	fmt.Println("Choose year: ")
	year := util.GetYear()

	// init sesh and db with default options
	db, err := initDynamoDbDefault(year)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Upload failed: %v", err)
	}

	for {
		fmt.Println("Choose a category: ")
		ch, err := util.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Upload failed: %v", err)
		}
		cat := menu.OptionsMap[ch]

		if cat == "Return" {
			fmt.Println("Returning to menu...")
			return nil
		}
		if cat == "all" {
			// upload all categories for given year; return when complete
			for _, cat := range opts[:4] {
				err := uploadFromDisk(db, year, cat, 1000)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("Upload failed: %v", err)
				}
			}
			fmt.Printf("Year %s uploaded. Returning to menu...\n", year)
			return nil
		}

		// upload single category
		err = uploadFromDisk(db, year, cat, 1000)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Upload failed: %v", err)
		}
		// coninue/return
		fmt.Printf("year %s - %s uploaded. Continue?\n", year, cat)
		yes := util.Ask4confirm()
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

	// list tables currently in DB
	t, err := dynamo.ListTables(db.Svc)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("InitDynamoDbDefault failed: %v", err)
	}
	if t == 0 { // create tables if none
		err := initDynamoTables(db)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("InitDynamoDbDefault failed: %v", err)
		}
	}

	return db, nil
}

// initDynamoTables initializes Tables for each object category for the given year
// and adds the corresponding Table object to the db.Tables field.
// TableName format: "[cf] %s: individuals", year
//                   "[cf] %s: candidates", year
//                   "[cf] %s: committees", year
//                   "[cf] %s: cmte_tx_data", year
//                   "[cf] %s: cmte_financials", year
//                   "[cf] %s: top_ovearll", year
func initDynamoTables(db *dynamo.DbInfo) error {
	for _, t := range db.Tables {
		err := dynamo.CreateTable(db.Svc, t)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("initDynamoTables failed: %v", err)
		}
	}

	return nil
}

// uploadFromDisk intitializes the batch upload process to DynamoDB for the specified year: bucket
// Each call uploads n items to the bucket's correspondi
func uploadFromDisk(db *dynamo.DbInfo, year, bucket string, n int) error {
	i := 0
	fmt.Printf("starting upload for %s - %s\n", year, bucket)
	// upload Top 100,000 individuals only
	if bucket == "individuals" {
		err := uploadTopIndv(db, year, bucket, n)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("UploadFromDisk failed: %v", err)
		}
		return nil

	}

	startKey, err := persist.GetKey(year, bucket)
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

		// batch write n returned objects, 25 (max) per iteration
		for {
			if len(objs) == 0 {
				break
			}
			if len(objs) < 25 { // final batch write
				err := dynamo.BatchWriteCreate(db.Svc, db.Tables[bucket], db.FailConfig, objs)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("UploadFromDisk failed: %v", err)
				}
				break
			}

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

		// update startKey & log currKey value
		startKey = currKey
		err = persist.LogKey(year, bucket, currKey)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("UploadFromDisk failed: %v", err)
		}

		fmt.Println("objects wrote to table: ", i)

		// last batch of objects wrote to table
		if len(objs) < n {
			fmt.Println("last batch wrote to table")
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
	fmt.Printf("wrote %d items to table %s\n", i, db.Tables[bucket].TableName)
	fmt.Println()

	return nil
}

// initTableObjs creates dynamo.Table objects for given year in memory only and
// adds them to the db.Tables field. See InitDynamoTables description for TableName format.
func initTableObjs(db *dynamo.DbInfo, year string) {
	indv := fmt.Sprintf("[cf] %s: individuals", year)        // pk = State
	cand := fmt.Sprintf("[cf] %s: candidates", year)         // pk = State
	cmte := fmt.Sprintf("[cf] %s: committees", year)         // pk = State
	cmteData := fmt.Sprintf("[cf] %s: cmte_tx_data", year)   // pk = Name
	cmteFin := fmt.Sprintf("[cf] %s: cmte_financials", year) // pk = Name
	topOverall := fmt.Sprintf("[cf] %s: top_overall", year)  // pk = SizeLimit

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
	t = dynamo.CreateNewTableObj(cmteData, "Name", "string", "ID", "string")
	db.AddTable(t)

	// create object tables
	t = dynamo.CreateNewTableObj(cmteFin, "Name", "string", "ID", "string")
	db.AddTable(t)

	// create TopOverall table
	t = dynamo.CreateNewTableObj(topOverall, "SizeLimit", "int", "Category", "string")
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
				err := dynamo.BatchWriteCreate(db.Svc, db.Tables[bucket], db.FailConfig, objs)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("UploadFromDisk failed: %v", err)
				}
				break
			}

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

	fmt.Println("***** UPLOAD FINSIHED *****")
	fmt.Printf("wrote %d items to table %s\n", i, db.Tables[bucket].TableName)
	fmt.Println()

	return nil
}
