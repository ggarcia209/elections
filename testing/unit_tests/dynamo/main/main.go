package main

/* TEST NOTES */
// CreateTable call for existing table returns error
// Calling CreateTable(t) and PutItem(t) in same routine returns resource not found error
//

import (
	"fmt"
	"os"

	"github.com/elections/source/donations"
	"github.com/elections/source/dynamo"
)

var item0 = &donations.Individual{
	ID:         "indv00",
	Name:       "Guzman, El Chapo",
	City:       "Culiacan",
	State:      "Sinaloa",
	TotalInAmt: 20000000000,
	SendersAmt: map[string]float32{"the_crips_llc": 2000000, "cia": 100000000},
}

var item1 = &donations.Individual{
	ID:         "indv01",
	Name:       "Trump, Donald",
	City:       "New York",
	State:      "New York",
	TotalInAmt: 5000000000,
	SendersAmt: map[string]float32{"russia": 12500000, "team_trump": 5000000, "lock_her_up_pac": 333333},
}

var item2 = &donations.Individual{
	ID:         "indv02",
	Name:       "Biden, Joe",
	City:       "Scranton",
	State:      "Pennsylvania",
	TotalInAmt: 350000000,
	SendersAmt: map[string]float32{"china": 2000000, "burisma": 250000},
}

var itemList = []interface{}{item0, item1, item2}

func main() {
	db := dynamo.InitDbInfo()
	db.SetSvc(dynamo.InitSesh())
	db.AddTable(dynamo.CreateNewTableObj("test_table_2", "State", "string", "ID", "string"))
	db.AddTable(dynamo.CreateNewTableObj("test_table_3", "SizeLimit", "int", "Category", "string")) // topOverall obj schema
	db.SetFailConfig(dynamo.DefaultFailConfig)

	/* q0 := dynamo.CreateNewQueryObj(item0.State, item0.ID)
	q1 := dynamo.CreateNewQueryObj(item1.State, item1.ID)
	q2 := dynamo.CreateNewQueryObj(item2.State, item2.ID)
	queries := []*dynamo.Query{q0, q1, q2} */

	item := &donations.TopOverallData{
		Category:  "indv",
		SizeLimit: 5,
	}

	/* err := testCreateItem(db, item)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	} */

	q := dynamo.CreateNewQueryObj(item.SizeLimit, item.Category)
	err := testGetItem(db, q)
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	fmt.Println("main done")
}

// test InitSesh() & List Tables
func testInit() error {
	svc := dynamo.InitSesh()
	err := dynamo.ListTables(svc)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("TestInit failed: %v", err)
	}

	fmt.Println("testInit done")
	fmt.Println()

	return nil
}

func testCreateTable(db *dynamo.DbInfo) error {
	/* err := dynamo.CreateTable(db.Svc, db.Tables["test_table_2"])
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("testCreateTable failed: %v", err)
	} */

	err := dynamo.CreateTable(db.Svc, db.Tables["test_table_3"])
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("testCreateTable failed: %v", err)
	}

	err = dynamo.ListTables(db.Svc)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("TestInit failed: %v", err)
	}

	fmt.Println("testCreateTable done")

	return nil
}

func testCreateItem(db *dynamo.DbInfo, item interface{}) error {
	err := dynamo.CreateItem(db.Svc, item, db.Tables["test_table_3"])
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("TestInit failed: %v", err)
	}

	fmt.Println("testCreateItem done")
	fmt.Println()

	return nil
}

func testGetItem(db *dynamo.DbInfo, q *dynamo.Query) error {
	// q := dynamo.CreateNewQueryObj("New York", "indv01")
	obj := &donations.TopOverallData{}

	item, err := dynamo.GetItem(db.Svc, q, db.Tables["test_table_3"], obj)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("testGetItem failed: %v", err)
	}

	fmt.Println("item: ", item)
	// fmt.Println("name: ", item.(*donations.Individual).Name)
	fmt.Println("testGetItem done")
	fmt.Println()

	return nil
}

func testUpdateItem(db *dynamo.DbInfo, q *dynamo.Query) error {
	// svc := dynamo.InitSesh()
	// t := dynamo.CreateNewTableObj("test_table_2", "State", "string", "ID", "string")
	// q := dynamo.CreateNewQueryObj("New York", "indv00")
	// q.UpdateCurrent("NetBalance", 5000)
	t := db.Tables["test_table_2"]

	err := dynamo.UpdateItem(db.Svc, q, t)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("testUpdateItem failed: %v", err)
	}

	fmt.Println("testUpdateItem done")
	fmt.Println()

	return nil
}

func testDeleteItem(db *dynamo.DbInfo, q *dynamo.Query) error {
	t := db.Tables["test_table_2"]
	// q := dynamo.CreateNewQueryObj("New York", "indv01")

	err := dynamo.DeleteItem(db.Svc, q, t)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("testUpdateItem failed: %v", err)
	}

	fmt.Println("testDeleteItem done")
	fmt.Println()

	return nil
}

func testBatchWriteCreate(db *dynamo.DbInfo, items []interface{}) error {
	t := db.Tables["test_table_2"]
	err := dynamo.BatchWriteCreate(db.Svc, t, db.FailConfig, items)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("testBatchWriteCreate failed: %v", err)
	}

	fmt.Println("testBatchWriteCreate done")
	fmt.Println()

	return nil
}

func testBatchGet(db *dynamo.DbInfo, items []*dynamo.Query) error {
	t := db.Tables["test_table_2"]
	refs := []interface{}{}

	for i := 0; i < len(items); i++ {
		r := &donations.Individual{}
		refs = append(refs, r)
	}

	results, err := dynamo.BatchGet(db.Svc, t, db.FailConfig, items, refs)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("testBatchWriteDelete failed: %v", err)
	}

	fmt.Println("items: ", results)
	for _, r := range results {
		fmt.Println("Result: ", r)
	}

	fmt.Println("testBatchGet done")
	fmt.Println()

	return nil
}

func testBatchWriteDelete(db *dynamo.DbInfo, items []*dynamo.Query) error {
	t := db.Tables["test_table_2"]
	err := dynamo.BatchWriteDelete(db.Svc, t, db.FailConfig, items)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("testBatchWriteDelete failed: %v", err)
	}

	fmt.Println("testBatchWriteDelete done")
	fmt.Println()

	return nil
}
