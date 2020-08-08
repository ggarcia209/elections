package upload

import (
	"fmt"
	"projects/elections/dynamo"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// TableInit creates DynamoDB tables for
// 'individuals', 'committees', 'candidates', 'disbursement_recipients',
// using the year as the Primary Key and the object ID as the Sort Key
func TableInit(svc *dynamodb.DynamoDB) error {
	// create Table objects in memory
	it := initIndvTable()
	ct := initCmteTable()
	cat := initCandTable()
	dt := initDisbRecTable()
	tables := []*dynamo.Table{it, ct, cat, dt}

	// create DynamoDB table for each correspond Table object
	for _, t := range tables {
		err := dynamo.CreateTable(svc, t)
		if err != nil {
			fmt.Println("UploadInit failed: ", err)
			return fmt.Errorf("UploadInit failed: %v", err)
		}
	}

	return nil
}

func initIndvTable() *dynamo.Table {
	t := dynamo.CreateNewTableObj("individuals", "year", "int", "ID", "string")
	return t
}

func initCmteTable() *dynamo.Table {
	t := dynamo.CreateNewTableObj("committees", "year", "int", "ID", "string")
	return t
}

func initCandTable() *dynamo.Table {
	t := dynamo.CreateNewTableObj("candidates", "year", "int", "ID", "string")
	return t
}

func initDisbRecTable() *dynamo.Table {
	t := dynamo.CreateNewTableObj("disbursement_recipients", "year", "int", "ID", "string")
	return t
}
