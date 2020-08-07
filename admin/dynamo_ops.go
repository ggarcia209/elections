package admin

import (
	"fmt"

	"github.com/elections/dynamo"
	"github.com/elections/persist"
)

// InitDynamoDbDefault initializes a dynamo.DbInfo object with default DynamoDB session settings
func InitDynamoDbDefault() *dynamo.DbInfo {
	db := dynamo.InitDbInfo()
	db.SetSvc(dynamo.InitSesh())
	db.SetFailConfig(dynamo.DefaultFailConfig)
	return db
}

// InitDynamoTables initializes Tables for each object category for the given year
// and adds the corresponding Table object to the db.Tables field.
// TableName format: "[ecf] %s: individuals", year
//                   "[ecf] %s: candidates", year
//                   "[ecf] %s: committees", year
//                   "[ecf] %s: cmte_tx_data", year
//                   "[ecf] %s: cmte_financials", year
//                   "[ecf] %s: top_ovearll", year
func InitDynamoTables(db *dynamo.DbInfo, year string) error {
	indv := fmt.Sprintf("[ecf] %s: individuals", year)        // pk = State
	cand := fmt.Sprintf("[ecf] %s: candidates", year)         // pk = State
	cmte := fmt.Sprintf("[ecf] %s: committees", year)         // pk = State
	cmteData := fmt.Sprintf("[ecf] %s: cmte_tx_data", year)   // pk = Name
	cmteFin := fmt.Sprintf("[ecf] %s: cmte_financials", year) // pk = Name
	topOverall := fmt.Sprintf("[ecf] %s: top_overall", year)  // pk = SizeLimit

	// create object tables
	t := dynamo.CreateNewTableObj(indv, "State", "string", "ID", "string")
	err := dynamo.CreateTable(db.Svc, t)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("InitDynamoTables failed: %v", err)
	}
	db.AddTable(t)

	// create object tables
	t = dynamo.CreateNewTableObj(cand, "State", "string", "ID", "string")
	err = dynamo.CreateTable(db.Svc, t)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("InitDynamoTables failed: %v", err)
	}
	db.AddTable(t)

	// create object tables
	t = dynamo.CreateNewTableObj(cmte, "State", "string", "ID", "string")
	err = dynamo.CreateTable(db.Svc, t)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("InitDynamoTables failed: %v", err)
	}
	db.AddTable(t)

	// create object tables
	t = dynamo.CreateNewTableObj(cmteData, "Name", "string", "ID", "string")
	err = dynamo.CreateTable(db.Svc, t)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("InitDynamoTables failed: %v", err)
	}
	db.AddTable(t)

	// create object tables
	t = dynamo.CreateNewTableObj(cmteFin, "Name", "string", "ID", "string")
	err = dynamo.CreateTable(db.Svc, t)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("InitDynamoTables failed: %v", err)
	}
	db.AddTable(t)

	// create TopOverall table
	t = dynamo.CreateNewTableObj(topOverall, "SizeLimit", "int", "Category", "string")
	err = dynamo.CreateTable(db.Svc, t)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("InitDynamoTables failed: %v", err)
	}
	db.AddTable(t)

	return nil
}

// InitTableObjsOnly creates dynamo.Table objects for given year in memory only and
// adds them to the db.Tables field. See InitDynamoTables description for TableName format.
func InitTableObjsOnly(db *dynamo.DbInfo, year string) {
	indv := fmt.Sprintf("[ecf] %s: individuals", year)        // pk = State
	cand := fmt.Sprintf("[ecf] %s: candidates", year)         // pk = State
	cmte := fmt.Sprintf("[ecf] %s: committees", year)         // pk = State
	cmteData := fmt.Sprintf("[ecf] %s: cmte_tx_data", year)   // pk = Name
	cmteFin := fmt.Sprintf("[ecf] %s: cmte_financials", year) // pk = Name
	topOverall := fmt.Sprintf("[ecf] %s: top_overall", year)  // pk = SizeLimit

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

// UploadFromDisk intitializes the batch upload process to DynamoDB for the specified year: bucket
// Each call uploads n items to the bucket's correspondi
func UploadFromDisk(db *dynamo.DbInfo, year, bucket string, n int) error {
	i := 0
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

	fmt.Println("***** UPLOAD FINSIHED *****")
	fmt.Printf("wrote %d items to table %s\n", i, db.Tables[bucket].TableName)
	fmt.Println()

	return nil
}
