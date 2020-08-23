package admin

import (
	"fmt"
	"sort"
	"strings"

	"github.com/elections/source/donations"
	"github.com/elections/source/indexing"
	"github.com/elections/source/persist"
	"github.com/elections/source/ui"
)

// entry represents a k/v pair in a sorted map
type entry struct {
	ID    string
	Total float32
}

// entries represents a sorted map
type entries []entry

func (s entries) Len() int           { return len(s) }
func (s entries) Less(i, j int) bool { return s[i].Total > s[j].Total }
func (s entries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// ViewMenu provides a submenu for searching the data,
// viewing rankings by year/category/party,
// and viewing the search index
func ViewMenu() error {
	opts := []string{
		"Search Data",
		"View Top Rankings",
		"View Search Index",
		"Return to Main Menu",
		// Query DynamoDB
	}
	menu := ui.CreateMenu("admin-view-data", opts)

	fmt.Println("***** View Data *****")
	// get output path
	output, err := getPath(false)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ViewMenu failed: %v", err)
	}

	indexing.OUTPUT_PATH = output
	for {
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("ViewMenu failed: %v", err)
		}

		switch {
		case ch == 0: // Seach
			err := searchData()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ViewMenu failed: %v", err)
			}
		case ch == 1: // Rankings
			err := viewRankings()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ViewMenu failed: %v", err)
			}
		case ch == 2: // View Index
			err := indexing.ViewIndex()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ViewMenu failed: %v", err)
			}
		case ch == 3: // Return
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

func searchData() error {
	opts := []string{"Search Datasets", "Lookup by Object IDs", "Return"}
	menu := ui.CreateMenu("admin-search-sub", opts)
	for {
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("searchData failed: %v", err)
		}
		switch {
		case menu.OptionsMap[ch] == "Search Datasets":
			err := queryData()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Lookup by Object IDs":
			err := lookupByID()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Return":
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

// routine for searching data by query
// returns results and provides sub menu for
// selecting dataset by object ID and year
func queryData() error {
	for {
		// get query from user / return & print results
		fmt.Println("*** Search ***")
		txt := ui.GetQuery()
		q := indexing.CreateQuery(txt, "local_admin")
		res, err := indexing.GetResults(q)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("searchData failed: %v", err)
		}
		resMap := make(map[string]indexing.SearchData)
		printResults(res)

		// crate submenus for selecting dataset from search results
		ids := []string{}
		for _, r := range res {
			ids = append(ids, r.ID)
			resMap[r.ID] = r
		}
		ids = append(ids, "exit")
		idsSubmenu := ui.CreateMenu("admin-search-results", ids)
		for {
			fmt.Println("Choose an ID from the list to view more info or choose 'exit' to return")
			chID, err := ui.Ask4MenuChoice(idsSubmenu)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}
			objID := idsSubmenu.OptionsMap[chID]
			if objID == "exit" {
				fmt.Println("Returning to menu...")
				return nil
			}

			// select year for given object
			yrs := resMap[objID].Years
			yrsSubMenu := ui.CreateMenu("admin-search-result-years", yrs)
			fmt.Println("Choose a year to view the objects data for that year: ")
			chYr, err := ui.Ask4MenuChoice(yrsSubMenu)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}
			year := yrsSubMenu.OptionsMap[chYr]
			bucket := resMap[objID].Bucket

			// get year/obj dataset from disk
			obj, err := persist.GetObject(year, bucket, objID)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}

			fmt.Println("RESULT: ")
			fmt.Printf("%#v\n", obj)
			fmt.Println()

			fmt.Println("Return to results?")
			yes := ui.Ask4confirm()
			if yes {
				continue
			}
			break
		}

		fmt.Println("New Search?")
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

func lookupByID() error {
	for {
		// get query from user / return & print results
		fmt.Println("*** Lookup by ID ***")
		fmt.Println("Enter IDs seperated by spaces")
		lu := []string{}
		txt := ui.GetQuery()
		ss := strings.Split(txt, " ")
		for _, s := range ss {
			s = strings.TrimSpace(s)
			lu = append(lu, s)
		}

		res, err := indexing.LookupSearchData(lu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("searchData failed: %v", err)
		}
		resMap := make(map[string]indexing.SearchData)
		printResults(res)

		// crate submenus for selecting dataset from search results
		ids := []string{}
		for _, r := range res {
			ids = append(ids, r.ID)
			resMap[r.ID] = r
		}
		ids = append(ids, "exit")
		idsSubmenu := ui.CreateMenu("admin-search-results", ids)
		for {
			fmt.Println("Choose an ID from the list to view more info or choose 'exit' to return")
			chID, err := ui.Ask4MenuChoice(idsSubmenu)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}
			objID := idsSubmenu.OptionsMap[chID]
			if objID == "exit" {
				fmt.Println("Returning to menu...")
				break
			}

			// select year for given object
			yrs := resMap[objID].Years
			yrsSubMenu := ui.CreateMenu("admin-search-result-years", yrs)
			fmt.Println("Choose a year to view the objects data for that year: ")
			chYr, err := ui.Ask4MenuChoice(yrsSubMenu)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}
			year := yrsSubMenu.OptionsMap[chYr]
			bucket := resMap[objID].Bucket

			// get year/obj dataset from disk
			obj, err := persist.GetObject(year, bucket, objID)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}
			if bucket == "committees" {
				// get tx data for committee
				txData, err := persist.GetObject(year, "cmte_tx_data", objID)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("searchData failed: %v", err)
				}
				fmt.Println("RESULT: ")
				fmt.Printf("%#v\n", obj)
				fmt.Printf("%#v\n", txData)
				fmt.Println()
			} else {
				fmt.Println("RESULT: ")
				fmt.Printf("%#v\n", obj)
				fmt.Println()
			}
			fmt.Println("Return to results?")
			yes := ui.Ask4confirm()
			if yes {
				continue
			}
			break
		}
		fmt.Println("New search?")
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}
}

// routine for viewing rankings by year/category/party
func viewRankings() error {
	yrs := []string{
		"2020", "2018", "2016", "2014", "2012",
		"2010", "2008", "2006", "2004", "2002",
		"2000", "1998", "1996", "1994", "1992",
		"1990", "1988", "1986", "1984", "1982",
		"1980",
	}
	yrsMenu := ui.CreateMenu("admin-rankings-yrs", yrs)

	cats := []string{
		"indv", "indv_rec", "cmte_donors", "cmte_recs", "cmte_exp", "cand", "cand_exp",
	}
	catMenu := ui.CreateMenu("admin-rankings-cats", cats)
	ptyMap := map[string][]string{
		"indv":        []string{"indv", "cancel"},
		"indv_rec":    []string{"indv_rec", "cancel"},
		"cmte_donors": []string{"cmte_donors_all", "cmte_donors_d", "cmte_donors_r", "cmte_donors_na", "cmte_donors_misc", "cancel"},
		"cmte_recs":   []string{"cmte_rec_all", "cmte_rec_d", "cmte_rec_r", "cmte_rec_na", "cmte_rec_misc", "cancel"},
		"cmte_exp":    []string{"cmte_exp_all", "cmte_exp_d", "cmte_exp_r", "cmte_exp_na", "cmte_exp_misc", "cancel"},
		"cand":        []string{"cand_all", "cand_d", "cand_r", "cand_na", "cand_misc", "cancel"},
		"cand_exp":    []string{"cand_exp_all", "cand_exp_d", "cand_exp_r", "cand_exp_na", "cand_exp_misc", "cancel"},
	}

	for {
		// get Year
		ch, err := ui.Ask4MenuChoice(yrsMenu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewRankings failed: %v", err)
		}
		year := yrsMenu.OptionsMap[ch]

		// get category
		ch, err = ui.Ask4MenuChoice(catMenu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewRankings failed: %v", err)
		}
		cat := catMenu.OptionsMap[ch]

		// get party sub category
		ptys := ptyMap[cat]
		ptysMenu := ui.CreateMenu("admin-rankings-pty", ptys)
		ch, err = ui.Ask4MenuChoice(ptysMenu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewRankings failed: %v", err)
		}

		// display sorted rankings from selection
		selection := ptysMenu.OptionsMap[ch]
		fmt.Println("selection: ", selection)
		if selection == "cancel" {
			fmt.Println("Returning to menu...")
			return nil
		}
		rankings, err := persist.GetObject(year, "top_overall", selection)
		sorted := sortRankings(rankings.(*donations.TopOverallData).Amts)
		err = printSortedRankings(rankings.(*donations.TopOverallData), sorted)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewRankings failed: %v", err)
		}

		fmt.Println("Lookup by ID?")
		yes := ui.Ask4confirm()
		if yes {
			err := lookupByID()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("viewRankings failed: %v", err)
			}
		}

		fmt.Println("View new category?")
		yes = ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}

}

// prints search results
func printResults(res []indexing.SearchData) {
	fmt.Println("SEARCH RESULTS: ")
	for i, r := range res {
		fmt.Printf(
			"%d)  ID: %s\n\tName: %s\n\tCit: %s\n\tState: %s\n\tBucket: %s\n\tYears: %s\n",
			i, r.ID, r.Name, r.City, r.State, r.Bucket, r.Years,
		)
	}
	fmt.Println()
}

// sort rankings map by vale
func sortRankings(m map[string]float32) entries {
	var es entries
	for k, v := range m {
		e := entry{k, v}
		es = append(es, e)
	}
	sort.Sort(es)
	return es
}

// lookup corresponding SearchData object for each ID in rankings and print data
func printSortedRankings(r *donations.TopOverallData, sorted entries) error {
	fmt.Println("Top Rankings")
	fmt.Println("Category: ", r.Category)
	fmt.Println("Size Limit: ", r.SizeLimit)
	fmt.Printf("Top %d:\n", r.SizeLimit)
	ids := []string{}
	for _, e := range sorted {
		ids = append(ids, e.ID)
	}
	sds, err := indexing.LookupSearchData(ids)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printSortedRankings failed: %v", err)
	}
	for i, sd := range sds {
		fmt.Printf("Rank %d)  %s - %s (%s, %s): %v\n", i, sd.ID, sd.Name, sd.City, sd.State, r.Amts[sd.ID])
	}
	fmt.Println()

	return nil
}
