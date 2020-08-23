package admin

import (
	"fmt"

	"github.com/elections/source/indexing"
	"github.com/elections/source/persist"
	"github.com/elections/source/ui"
)

// Delete provides menu options for delete operations
func Delete() error {
	opts := []string{
		"Delete Dataset by Year",
		"Delete Dataset by Category",
		"Delete Database",
		"Delete SearchIndex",
		"Delete Metadata",
		"Delete All Data on Disk",
		"Delete DynamoDB Table",
		"Return",
	}
	menu := ui.CreateMenu("admin-delete", opts)
	path, err := getPath(false)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("deleteYr failed: %v", err)
	}
	persist.OUTPUT_PATH = path
	indexing.OUTPUT_PATH = path

	for {
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Delete failed: %v", err)
		}
		switch {
		case menu.OptionsMap[ch] == "Delete Dataset by Year":
			err := deleteYr()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Delete failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Delete Dataset by Category":
			err := deleteCat()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Delete failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Delete Database":
			err := delDB()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Delete failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Delete SearchIndex":
			err := delIndex()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Delete failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Delete Metadata":
			err := delMeta()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Delete failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Delete All Data on Disk":
			err := delAll()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Delete failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Delete DynamoDB Table":
			err := DeleteDynamoTable()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Delete failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Return":
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

// get user input and delete by year
func deleteYr() error {
	fmt.Println("Delete by year: ")
	year := ui.GetYear()
	if year == "cancel" {
		fmt.Println("Returning to menu...")
		return nil
	}

	fmt.Printf("Are you sure you want to delete all data for the year %s?\n", year)
	yes := ui.Ask4confirm()
	if !yes {
		fmt.Println("Returning to menu...")
		return nil
	}

	err := persist.DeleteYear(year)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("deleteYr failed %v", err)
	}
	fmt.Println("deleted: ", year)
	return nil
}

// get user input and delete by year/category
func deleteCat() error {
	fmt.Println("Delete by year: ")
	year := ui.GetYear()
	if year == "cancel" {
		fmt.Println("Returning to menu...")
		return nil
	}

	opts := []string{"individuals", "committees", "candidates", "top_overall", "cancel"}
	menu := ui.CreateMenu("admin-delete-bycat", opts)
	fmt.Printf("Choose a cateogry to delete (year: %s):\n", year)
	ch, err := ui.Ask4MenuChoice(menu)
	cat := menu.OptionsMap[ch]
	if cat == "cancel" {
		fmt.Println("Returning to menu...")
		return nil
	}

	fmt.Printf("Are you sure you want to delete all data for the year/category %s - %s?\n", year, cat)
	yes := ui.Ask4confirm()
	if !yes {
		fmt.Println("Returning to menu...")
		return nil
	}

	err = persist.DeleteCategory(year, cat)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("deleteCat failed: %v", err)
	}
	fmt.Printf("deleted %s - %s\n", year, cat)
	return nil
}

func delDB() error {
	fmt.Printf("Are you sure you want to delete the database at %s?\n", persist.OUTPUT_PATH)
	yes := ui.Ask4confirm()
	if !yes {
		fmt.Println("Returning to menu...")
		return nil
	}
	err := persist.DeleteDatabase()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("delDB failed %v", err)
	}
	return nil
}

func delIndex() error {
	fmt.Printf("Are you sure you want to delete the search index at %s?\n", persist.OUTPUT_PATH)
	yes := ui.Ask4confirm()
	if !yes {
		fmt.Println("Returning to menu...")
		return nil
	}
	err := persist.DeleteSearchIndex()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("delIndex failed %v", err)
	}
	return nil
}

func delMeta() error {
	fmt.Printf("Are you sure you want to delete all metadata at ../db?\n")
	yes := ui.Ask4confirm()
	if !yes {
		fmt.Println("Returning to menu...")
		return nil
	}
	err := persist.DeleteMetaData()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("delMeta failed %v", err)
	}
	return nil
}

func delAll() error {
	fmt.Println("Are you sure you want to delete ALL data on disk?")
	yes := ui.Ask4confirm()
	if !yes {
		fmt.Println("Returning to menu...")
		return nil
	}
	err := persist.DeleteAll()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("delAll failed: %v", err)
	}
	return nil
}
