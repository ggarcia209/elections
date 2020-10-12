// The admin service enables a local user to create the datasets and view the data on the command line like a web client.
// However, the admin service is distinct in that it is ran locally, and provides a much greater in-depth view of the data.
// The primary features include:
// - Building primary and secondary datasets from the FEC's bulk data files
// - Building search index and lookup data
// - Searching & viewing all data, including complete datasets by year/category,
//   complete search index, and application metadata.
// - Uploading datasets to DynamoDB tables
// - Deleting data from disk & DynamoDB.
package main

import (
	"fmt"
	"os"

	"github.com/elections/source/admin"
	"github.com/elections/source/persist"
	"github.com/elections/source/ui"
)

func main() {
	persist.InitDiskCache()
	delete := false
	opts := []string{
		"Process Raw Data",
		"Build/Update Search Index",
		"View/Seach Datasets",
		"Upload Data to DynamoDB",
		"Delete Data from Disk",
		"Exit Admin Console",
	}

	menu := ui.CreateMenu("admin-main", opts)
	// menu
	fmt.Println("***** Welcome! *****")

	for {
		if delete {
			// check ../db if delete operation was compeleted
			if _, err := os.Stat("../db"); os.IsNotExist(err) {
				os.Mkdir("../db", 0744)
				fmt.Printf("CreateDB successful: '../db' directory created")
			}
			delete = false
		}

		// get input from user
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		switch {
		case menu.OptionsMap[ch] == "Process Raw Data": // process new records
			err := admin.ProcessData()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case menu.OptionsMap[ch] == "Build/Update Search Index": // build index
			err := admin.BuildIndexFromYear()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case menu.OptionsMap[ch] == "View/Seach Datasets": // view data
			err := admin.ViewMenu()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case menu.OptionsMap[ch] == "Upload Data to DynamoDB": // upload
			err := admin.Upload()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case menu.OptionsMap[ch] == "Delete Data from Disk": // delete
			err := admin.Delete()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			delete = true
		case menu.OptionsMap[ch] == "Exit Admin Console": // exit
			fmt.Println("Terminating Admin console...")
			os.Exit(1)
		}
	}

}
