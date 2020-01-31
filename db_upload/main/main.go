package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"projects/elections/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/boltdb/bolt"
)

type Queue []string

func main() {
	// get flag for year
	year, err := getYear()
	if err != nil {
		fmt.Println("main failed: ", err)
		os.Exit(1)
	}

	// initialize session
	config := aws.NewConfig().WithRegion("us-west-1").WithMaxRetries(3)
	db := dynamo.InitSesh(config)

	// upload Candiate data
	// while unprocessed objects
	// iterate thru each item in the given year's dataset
	// create a queue of 1000 item id's max
	// while id's in queue
	// pop 25 item ids from queue - tx queue
	// for each item, get corresponding obj from database and add to tx queue
	// batch write 25 items to corresponding dynamoDB table
	// break when tx queue < 25
	// break when id queue < 1000

	// upload Committee data
	// upload Indvidual data
	// upload DisbRecipient data

}

func getYear() (string, error) {
	yearStr := "" // default return value

	for {
		// get flag and check validity by verifying if bucket for specified year exists
		yearFlag := flag.Int("year", 0, "'year' flag defines which election cycle's dataset to process")
		flag.Parse()
		year := *yearFlag
		if year == 0 {
			fmt.Println("'year' flag must be set to valid year")
			continue
		}

		// convert int value to string
		yearStr = strconv.Itoa(year)

		// open db and start view tx
		db, err := bolt.Open("db/offline_db.db", 0644, nil)
		defer db.Close()
		if err != nil {
			fmt.Println("FATAL: getYear failed: ")
			return "", fmt.Errorf("getYear failed: %v", err)
		}

		// check validity by searching for bucket corresponding to given year
		exists := true
		if err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(yearStr))
			if bucket == nil {
				exists = false
			}
			return nil
		}); err != nil {
			fmt.Println("FATAL: getYear failed: ", err)
			return "", fmt.Errorf("getYear failed: %v", err)
		}

		if exists == false {
			fmt.Printf("Invalid year: %d --- No dataset found!\n", year)
			continue
		} else {
			break
		}
	}
	return yearStr, nil
}
