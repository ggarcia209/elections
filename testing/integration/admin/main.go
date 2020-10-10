package main

import (
	"fmt"
	"os"

	"github.com/elections/source/persist"
	"github.com/elections/source/ui"
)

func main() {
	fmt.Println("testing admin.viewBucket()")
	err := viewBucket()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func viewBucket() error {
	year := ui.GetYear()
	opts := []string{"individuals", "committees", "candidates", "top_overall", "yearly_totals", "cancel"}
	menu := ui.CreateMenu("view-data-by-bucket", opts)
	start := ""   // start at first key in bucket
	curr := start // initialize starting key of next batch
	cont := false // continue to print next batch

	for {
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewBucket failed: %v", err)
		}
		switch {
		case menu.OptionsMap[ch] == "individuals":
			for {
				curr, cont, err = viewNext(year, menu.OptionsMap[ch], curr)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("viewBucket failed: %v", err)
				}
				if !cont {
					fmt.Println("Returning to menu...")
					break
				}
			}
		case menu.OptionsMap[ch] == "committees":
			for {
				curr, cont, err = viewNext(year, menu.OptionsMap[ch], curr)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("viewBucket failed: %v", err)
				}
				if !cont {
					fmt.Println("Returning to menu...")
					break
				}
			}
		case menu.OptionsMap[ch] == "candidates":
			for {
				curr, cont, err = viewNext(year, menu.OptionsMap[ch], curr)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("viewBucket failed: %v", err)
				}
				if !cont {
					fmt.Println("Returning to menu...")
					break
				}
			}
		case menu.OptionsMap[ch] == "top_overall":
			for {
				curr, cont, err = viewNext(year, menu.OptionsMap[ch], curr)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("viewBucket failed: %v", err)
				}
				if !cont {
					fmt.Println("Returning to menu...")
					break
				}
			}
		case menu.OptionsMap[ch] == "yearly_totals":
			for {
				curr, cont, err = viewNext(year, menu.OptionsMap[ch], curr)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("viewBucket failed: %v", err)
				}
				if !cont {
					fmt.Println("Returning to menu...")
					break
				}
			}
		case menu.OptionsMap[ch] == "cancel":
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

// print 1000 items from databse, ask user if continue
func viewNext(year, bucket, start string) (string, bool, error) {
	cont := false
	curr, err := persist.ViewDataByBucket(year, bucket, start)
	if err != nil {
		fmt.Println(err)
		return "", false, fmt.Errorf("viewNext failed: %v", err)
	}
	if curr == "" { // list exhausted
		return curr, false, nil
	}
	fmt.Println()
	fmt.Println(">>> Scan finished - print next 1000 objects?")
	yes := ui.Ask4confirm()
	if yes {
		cont = true
	}
	return curr, cont, nil
}
