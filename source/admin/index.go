// Package admin contains operations for running the local admin console service.
// Only the functions in this package are exposed to the admin service; lower
// level source packages remain encapsulated.
// This file contains operations for building and updating the search index
// and object lookup data from the processed datasets.
// NOTE: logic is not UX optimized and may contain unresolved errors.
package admin

import (
	"fmt"

	"github.com/elections/source/persist"

	"github.com/elections/source/indexing"
	"github.com/elections/source/ui"
)

// BuildIndexFromYear builds a new index from user-input year's
// dataset and adds it to the existing index on disk.
func BuildIndexFromYear() error {
	// get output path
	output, err := getPath(false)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("BuildIndexFromYear failed: %v", err)
	}
	indexing.OUTPUT_PATH = output
	persist.OUTPUT_PATH = output
	// create submenu
	opts := []string{"Build New Index", "Update Index", "Write Out Index", "Return"}
	menu := ui.CreateMenu("admin-index-options", opts)

	for {
		fmt.Println("Choose an option: ")
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("BuildIndexFromYear failed: %v", err)
		}
		choice := menu.OptionsMap[ch]
		if choice == "Return" {
			fmt.Println("Returning to menu...")
			return nil
		}
		fmt.Println("Choose year: ")
		year := ui.GetYear()

		switch {
		case choice == "Build New Index":
			err := build(year)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("BuildIndexFromYear failed: %v", err)
			}
		case choice == "Update Index":
			err := update(year)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("BuildIndexFromYear failed: %v", err)
			}
		case choice == "Write Out Index":
			subOpts := []string{"Write Out Search Index", "Write Out Index Data"}
			sub := ui.CreateMenu("index-sub-menu", subOpts)
			ch, err := ui.Ask4MenuChoice(sub)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("BuildIndexFromYear failed: %v", err)
			}
			err = indexing.WriteOutIndex(ch)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("BuildIndexFromYear failed: %v", err)
			}
		}

		fmt.Println("Continue?")
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}

}

func build(year string) error {
	err := indexing.BuildIndex(year)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("build failed: %v", err)
	}
	return nil
}

func update(year string) error {
	opts := []string{"individuals", "committees", "candidates", "Return"}
	menu := ui.CreateMenu("admin-index-update", opts)
	var category string
	for {
		fmt.Println("Choose category to read data from: ")
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("update failed: %v", err)
		}
		category = menu.OptionsMap[ch]
		if category == "Return" {
			fmt.Println("Returning to menu...")
			return nil
		}
		err = indexing.UpdateIndex(year, category)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("update failed: %v", err)
		}
		fmt.Println("Update another category?")
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}
