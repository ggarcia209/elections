package admin

import (
	"fmt"
	"sort"
	"strings"

	"github.com/elections/source/donations"
	"github.com/elections/source/indexing"
	"github.com/elections/source/persist"
	"github.com/elections/source/ui"
	"github.com/elections/source/util"
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
		"View Data by Year/Bucket",
		"View Top Rankings",
		"View Yearly Totals",
		"View Search Index",
		"Query DyanamoDB",
		"Return to Main Menu",
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
	persist.OUTPUT_PATH = output
	for {
		ch, err := ui.Ask4MenuChoice(menu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("ViewMenu failed: %v", err)
		}

		switch {
		case menu.OptionsMap[ch] == "Search Data":
			err := searchData()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ViewMenu failed: %v", err)
			}
		case menu.OptionsMap[ch] == "View Data by Year/Bucket":
			err := viewBucket()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ViewMenu failed: %v", err)
			}
		case menu.OptionsMap[ch] == "View Top Rankings":
			err := viewRankings()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ViewMenu failed: %v", err)
			}
		case menu.OptionsMap[ch] == "View Yearly Totals":
			err := viewYrTotals()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ViewMenu failed: %v", err)
			}
		case menu.OptionsMap[ch] == "View Search Index":
			err := indexing.ViewIndex()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ViewMenu failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Query DyanamoDB":
			err := QueryDynamoDB()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("ViewMenu failed: %v", err)
			}
		case menu.OptionsMap[ch] == "Return to Main Menu":
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
	cats := []string{
		"Individual Donors", "Individual Recipients", "Committee Donors", "Committee Recipients",
		"Committee Spenders", "Candidate Recipients", "Candidate Donors", "Candidate Spenders",
	}
	catMap := map[string][]string{
		"Individual Donors":     []string{"individuals", "donor"},
		"Individual Recipients": []string{"individuals", "rec"},
		"Committee Donors":      []string{"cmte_tx_data", "donor"},
		"Committee Recipients":  []string{"cmte_tx_data", "rec"},
		"Committee Spenders":    []string{"cmte_tx_data", "exp"},
		"Candidate Recipients":  []string{"candidates", "rec"},
		"Candidate Donors":      []string{"candidates", "donor"},
		"Candidate Spenders":    []string{"candidates", "exp"},
	}
	catMenu := ui.CreateMenu("admin-rankings-cats", cats)
	ptyMap := map[string][]string{
		"Individual Donors":     []string{"ALL", "cancel"},
		"Individual Recipients": []string{"ALL", "cancel"},
		"Committee Donors":      []string{"ALL", "DEM", "REP", "IND", "OTH", "UNK", "cancel"},
		"Committee Recipients":  []string{"ALL", "DEM", "REP", "IND", "OTH", "UNK", "cancel"},
		"Committee Spenders":    []string{"ALL", "DEM", "REP", "IND", "OTH", "UNK", "cancel"},
		"Candidate Recipients":  []string{"ALL", "DEM", "REP", "IND", "OTH", "UNK", "cancel"},
		"Candidate Donors":      []string{"ALL", "DEM", "REP", "IND", "OTH", "UNK", "cancel"},
		"Candidate Spenders":    []string{"ALL", "DEM", "REP", "IND", "OTH", "UNK", "cancel"},
	}

	for {
		// get Year
		year := ui.GetYear()

		// get category
		ch, err := ui.Ask4MenuChoice(catMenu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewRankings failed: %v", err)
		}
		cat := catMenu.OptionsMap[ch]
		b, c := catMap[cat][0], catMap[cat][1]

		// get party sub category
		ptys := ptyMap[cat]
		ptysMenu := ui.CreateMenu("admin-rankings-pty", ptys)
		ch, err = ui.Ask4MenuChoice(ptysMenu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewRankings failed: %v", err)
		}

		// get party
		pty := ptysMenu.OptionsMap[ch]
		fmt.Println("selection: ", pty)
		if pty == "cancel" {
			fmt.Println("Returning to menu...")
			return nil
		}

		// get object
		id := year + "-" + b + "-" + c + "-" + pty
		obj, err := persist.GetObject(year, "top_overall", id)
		rankings := obj.(*donations.TopOverallData)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewRankings failed: %v", err)
		}

		// print sorted list of top overall entities
		sorted := util.SortMapObjectTotals(rankings.Amts)
		err = printSortedRankings(rankings, sorted)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewRankings failed: %v", err)
		}

		// option to lookup object form list by ID
		fmt.Println("Lookup by ID?")
		yes := ui.Ask4confirm()
		if yes {
			err := lookupByID()
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("viewRankings failed: %v", err)
			}
		}

		// option to view new top overall category
		fmt.Println("View new category?")
		yes = ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}

}

// routine for viewing rankings by year/category/party
func viewYrTotals() error {
	// menu options
	cats := []string{"Funds Received", "Funds Donated or Transferred", "Funds Expensed"}
	catMap := map[string]string{
		"Funds Received":               "rec",
		"Funds Donated or Transferred": "donor",
		"Funds Expensed":               "exp",
	}
	catMenu := ui.CreateMenu("admin-rankings-cats", cats)
	ptys := []string{"All", "Democrat", "Republican", "Independent/Non-Affiliated", "Other", "Unknown", "cancel"}
	ptyMap := map[string]string{
		"All":                        "ALL",
		"Democrat":                   "DEM",
		"Republican":                 "REP",
		"Independent/Non-Affiliated": "IND",
		"Other":                      "OTH",
		"Unknown":                    "UNK",
		"cancel":                     "cancel",
	}

	for {
		// get Year
		year := ui.GetYear()

		// get category
		ch, err := ui.Ask4MenuChoice(catMenu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewYrTotals failed: %v", err)
		}
		c := catMenu.OptionsMap[ch]
		cat := catMap[c]

		// get party sub category
		ptysMenu := ui.CreateMenu("admin-rankings-pty", ptys)
		ch, err = ui.Ask4MenuChoice(ptysMenu)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewYrTotals failed: %v", err)
		}

		// get party
		p := ptysMenu.OptionsMap[ch]
		pty := ptyMap[p]
		fmt.Println("selection: ", pty)
		if pty == "cancel" {
			fmt.Println("Returning to menu...")
			return nil
		}

		// get object
		id := year + "-" + cat + "-" + pty
		obj, err := persist.GetObject(year, "yearly_totals", id)
		yt := obj.(*donations.YearlyTotal)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("viewYrTotals failed: %v", err)
		}

		// print sorted list of top overall entities
		fmt.Println("Yearly Total:")
		fmt.Printf("%s\t%s\t%s\n\tTotal: %d\n", yt.Year, yt.Category, yt.Party, int(yt.Total))

		// option to view new top overall category
		fmt.Println("View new category?")
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}

}

func viewBucket() error {
	year := ui.GetYear()
	opts := []string{"individuals", "committees", "cmte_tx_data", "cmte_fin", "candidates", "top_overall", "yearly_totals", "cancel"}
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
		case menu.OptionsMap[ch] == "cmte_tx_data":
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
		case menu.OptionsMap[ch] == "cmte_fin":
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
func printSortedRankings(r *donations.TopOverallData, sorted util.SortedTotalsMap) error {
	fmt.Println("Top Rankings")
	fmt.Printf("%s\t%s\t%s\n", r.Year, r.Category, r.Party)
	fmt.Println()

	ids := []string{}
	for _, e := range sorted {
		// sfmt.Printf("ID: %s\tTotal: %.2f\n", e.ID, e.Total)
		ids = append(ids, e.ID)
	}
	sds, err := indexing.LookupSearchData(ids)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printSortedRankings failed: %v", err)
	}
	for i, sd := range sds {
		fmt.Printf("Rank %d)  %s - %s (%s, %s): %.2f\n", i, sd.ID, sd.Name, sd.City, sd.State, r.Amts[sd.ID])
		if sd.Bucket == "individuals" && i == 99 {
			break
		}
	}
	fmt.Println()

	return nil
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
