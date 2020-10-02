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
		"View Index Metadata",
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
		case menu.OptionsMap[ch] == "View Index Metadata":
			err := viewIndexData()
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
		res, err := indexing.GetResultsFromShards(q)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("searchData failed: %v", err)
		}
		sds, err := indexing.LookupSearchData(res)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("searchData failed: %v", err)
		}
		resMap := make(map[string]indexing.SearchData)

		// crate submenus for selecting dataset from search results
		options := []string{"exit"}
		for _, sd := range sds {
			optn := sd.ID + " " + sd.Name + " " + sd.City + ", " + sd.State
			options = append(options, optn)
			resMap[sd.ID] = sd
		}
		idsSubmenu := ui.CreateMenu("admin-search-results", options)
		for {
			printResults(sds)
			fmt.Println("Choose an ID from the list to view more info or choose 'exit' to return")
			chID, err := ui.Ask4MenuChoice(idsSubmenu)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}
			objID := idsSubmenu.OptionsMap[chID]
			ss := strings.Split(objID, " ")
			objID = strings.TrimSpace(ss[0])
			if objID == "exit" {
				fmt.Println("Returning to menu...")
				return nil
			}

			// select year for given object
			yrs := resMap[objID].Years
			yrs = append(yrs, "Return")
			yrsSubMenu := ui.CreateMenu("admin-search-result-years", yrs)
			fmt.Println("Choose a year to view the objects data for that year: ")
			chYr, err := ui.Ask4MenuChoice(yrsSubMenu)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}
			year := yrsSubMenu.OptionsMap[chYr]
			if year == "Return" {
				fmt.Println("Returning to menu...")
				return nil
			}
			bucket := resMap[objID].Bucket

			// get year/obj dataset from disk
			obj, err := persist.GetObject(year, bucket, objID)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}

			err = printEntity(year, obj)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
			}

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
	return nil
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
			err = printEntity(year, obj)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("searchData failed: %v", err)
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
		"Individual Donors":     []string{"cancel", "ALL"},
		"Individual Recipients": []string{"cancel", "ALL"},
		"Committee Donors":      []string{"cancel", "ALL", "DEM", "REP", "IND", "OTH", "UNK"},
		"Committee Recipients":  []string{"cancel", "ALL", "DEM", "REP", "IND", "OTH", "UNK"},
		"Committee Spenders":    []string{"cancel", "ALL", "DEM", "REP", "IND", "OTH", "UNK"},
		"Candidate Recipients":  []string{"cancel", "ALL", "DEM", "REP", "IND", "OTH", "UNK"},
		"Candidate Donors":      []string{"cancel", "ALL", "DEM", "REP", "IND", "OTH", "UNK"},
		"Candidate Spenders":    []string{"cancel", "ALL", "DEM", "REP", "IND", "OTH", "UNK"},
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
		"cancel":                     "cancel",
		"All":                        "ALL",
		"Democrat":                   "DEM",
		"Republican":                 "REP",
		"Independent/Non-Affiliated": "IND",
		"Other":                      "OTH",
		"Unknown":                    "UNK",
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

func viewIndexData() error {
	id, err := indexing.GetIndexData()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("viewIndexData failed: %v", err)
	}
	fmt.Println("Index Data")
	fmt.Println("Terms Size: ", id.TermsSize)
	fmt.Println("Lookup Size: ", id.LookupSize)
	fmt.Println("Last Updated: ", id.LastUpdated)
	fmt.Println("Categories Completed for latest build: ", id.Completed)
	fmt.Println("Years completed: ", id.YearsCompleted)
	fmt.Println("Shards: ")
	for k, v := range id.Shards {
		fmt.Printf("Term: %s\tShards Created: %v\n", k, v)
	}
	fmt.Println()
	return nil
}

// prints search results
func printResults(res []indexing.SearchData) {
	fmt.Println("SEARCH RESULTS: ")
	for i, r := range res {
		fmt.Printf(
			"%d)  ID: %s\n\tName: %s\n\tCit: %s\n\tState: %s\n\tBucket: %s\n\tYears: %s\n",
			i+1, r.ID, r.Name, r.City, r.State, r.Bucket, r.Years,
		)
	}
	fmt.Println()
}

func printEntity(year string, ent interface{}) error {
	switch t := ent.(type) {
	case *donations.Individual:
		err := printIndividual(ent.(*donations.Individual))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("printEntity failed: %v", err)
		}
	case *donations.Committee:
		err := printCommittee(year, ent.(*donations.Committee))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("printEntity failed: %v", err)
		}
	case *donations.Candidate:
		err := printCandidate(ent.(*donations.Candidate))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("printEntity failed: %v", err)
		}
	default:
		_ = t
		return fmt.Errorf("wrong interface type")
	}
	return nil
}

func printIndividual(indv *donations.Individual) error {
	fmt.Println("Viewing data for Individual: ", indv.ID)
	fmt.Println("Name: ", indv.Name)
	fmt.Println("City: ", indv.City)
	fmt.Println("State: ", indv.State)
	fmt.Println("Zip: ", indv.Zip)
	fmt.Println("Occupation: ", indv.Occupation)
	fmt.Println("Employer: ", indv.Employer)
	fmt.Println()
	fmt.Println("Total Outgoing $: ", indv.TotalOutAmt)
	fmt.Println("Total Outgoing Txs: ", indv.TotalOutTxs)
	fmt.Println("Avg. Outgoing Tx: ", indv.AvgTxOut)
	fmt.Println("Total Incoming $: ", indv.TotalInAmt)
	fmt.Println("Total Incoming Txs: ", indv.TotalInTxs)
	fmt.Println("Avg. Incoming Tx: ", indv.AvgTxIn)
	fmt.Println()
	fmt.Println("Recipients: ")
	recSrt := util.SortMapObjectTotals(indv.RecipientsAmt)
	err := printSortedEntities(recSrt, indv.RecipientsAmt, indv.RecipientsTxs)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printIndvidual failed: %v", err)
	}
	fmt.Println("Senders: ")
	sendSrt := util.SortMapObjectTotals(indv.SendersAmt)
	err = printSortedEntities(sendSrt, indv.SendersAmt, indv.SendersTxs)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printIndvidual failed: %v", err)
	}
	fmt.Println()
	return nil
}

func printCandidate(cand *donations.Candidate) error {
	fmt.Println("Viewing data for Candidate: ", cand.ID)
	fmt.Println("Name: ", cand.Name)
	fmt.Println("Party: ", cand.Party)
	fmt.Println("Election Year: ", cand.ElectnYr)
	fmt.Println("Office: ", cand.Office)
	fmt.Println("Office State: ", cand.OfficeState)
	fmt.Println("Primary Committee ID: ", cand.PCC)
	fmt.Println("City: ", cand.City)
	fmt.Println("State: ", cand.State)
	fmt.Println("Zip: ", cand.Zip)
	fmt.Println()
	fmt.Println("Total Direct Outgoing $: ", cand.TotalDirectOutAmt)
	fmt.Println("Total Direct Outgoing Txs: ", cand.TotalDirectOutTxs)
	fmt.Println("Avg. Direct Outgoing Tx: ", cand.AvgDirectOut)
	fmt.Println("Total Direct Incoming $: ", cand.TotalDirectInAmt)
	fmt.Println("Total Direct Incoming Txs: ", cand.TotalDirectInTxs)
	fmt.Println("Avg. Direct Incoming Tx: ", cand.AvgDirectIn)
	fmt.Println()
	fmt.Println("Recipients: ")
	recSrt := util.SortMapObjectTotals(cand.DirectRecipientsAmts)
	err := printSortedEntities(recSrt, cand.DirectRecipientsAmts, cand.DirectRecipientsTxs)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printCandidate failed: %v", err)
	}
	fmt.Println("Senders: ")
	sendSrt := util.SortMapObjectTotals(cand.DirectSendersAmts)
	err = printSortedEntities(sendSrt, cand.DirectSendersAmts, cand.DirectSendersTxs)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printCandidate failed: %v", err)
	}
	fmt.Println()
	return nil
}

func printCommittee(year string, obj *donations.Committee) error {
	intf, err := persist.GetObject(year, "cmte_tx_data", obj.ID)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printCommittee failed: %v", err)
	}
	txd := intf.(*donations.CmteTxData)
	fmt.Println("Viewing data for Committee: ", obj.ID)
	fmt.Println("Name: ", obj.Name)
	fmt.Println("Party: ", obj.Party)
	fmt.Println("Designation: ", obj.Designation)
	fmt.Println("Type: ", obj.Type)
	fmt.Println("Treasurer Name: ", obj.TresName)
	fmt.Println("City: ", obj.City)
	fmt.Println("State: ", obj.State)
	fmt.Println("Zip: ", obj.Zip)
	fmt.Println("Filing Frequency: ", obj.FilingFreq)
	fmt.Println("Organization Type: ", obj.OrgType)
	fmt.Println("Connected Organization: ", obj.ConnectedOrg)
	fmt.Println("Candidate ID: ", obj.CandID)
	fmt.Println()
	fmt.Println("Contributions $: ", txd.ContributionsInAmt)
	fmt.Println("Contributions Txs: ", txd.ContributionsInTxs)
	fmt.Println("Avg. Contribution: ", txd.AvgContributionIn)
	fmt.Println("Other Receipts $: ", txd.OtherReceiptsInAmt)
	fmt.Println("Other Receipts Txs: ", txd.OtherReceiptsInTxs)
	fmt.Println("Avg. Other: ", txd.AvgOtherIn)
	fmt.Println("Total Incoming $: ", txd.TotalIncomingAmt)
	fmt.Println("Total Incoming Txs: ", txd.TotalIncomingTxs)
	fmt.Println("Avg. Incoming: ", txd.AvgIncoming)
	fmt.Println()
	fmt.Println("Transfers $: ", txd.TransfersAmt)
	fmt.Println("Transfers Txs: ", txd.TransfersTxs)
	fmt.Println("Avg. Transfer: ", txd.AvgTransfer)
	fmt.Println("Expenditures $: ", txd.ExpendituresAmt)
	fmt.Println("Expenditures Txs: ", txd.ExpendituresTxs)
	fmt.Println("Avg. Expenditure: ", txd.AvgExpenditure)
	fmt.Println("Total Outgoing $: ", txd.TotalOutgoingAmt)
	fmt.Println("Total Outgoing Txs: ", txd.TotalOutgoingTxs)
	fmt.Println("Avg. Outgoing: ", txd.AvgOutgoing)
	fmt.Println()
	fmt.Println("Top Indvidual Contributors: ")
	indvSrt := util.SortMapObjectTotals(txd.TopIndvContributorsAmt)
	err = printSortedEntities(indvSrt, txd.TopIndvContributorsAmt, txd.TopIndvContributorsTxs)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printCommittee failed: %v", err)
	}
	fmt.Println("Top Committee Contributors: ")
	cmteSrt := util.SortMapObjectTotals(txd.TopCmteOrgContributorsAmt)
	err = printSortedEntities(cmteSrt, txd.TopCmteOrgContributorsAmt, txd.TopCmteOrgContributorsTxs)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printCommittee failed: %v", err)
	}
	fmt.Println("Transfer Recipients: ")
	trsSrt := util.SortMapObjectTotals(txd.TransferRecsAmt)
	err = printSortedEntities(trsSrt, txd.TransferRecsAmt, txd.TransferRecsTxs)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printCommittee failed: %v", err)
	}
	fmt.Println("Top Expenditure Recipients: ")
	expSrt := util.SortMapObjectTotals(txd.TopExpRecipientsAmt)
	err = printSortedEntities(expSrt, txd.TopExpRecipientsAmt, txd.TopExpRecipientsTxs)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printCommittee failed: %v", err)
	}
	fmt.Println()
	return nil
}

// lookup corresponding SearchData object for each ID in rankings and print data
func printSortedEntities(sorted util.SortedTotalsMap, orig, txs map[string]float32) error {
	ids := []string{}
	for _, e := range sorted {
		ids = append(ids, e.ID)
	}
	sds, err := indexing.LookupSearchData(ids)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printSortedEntities failed: %v", err)
	}
	for i, sd := range sds {
		ID := sd.ID
		amt := orig[ID]
		tx := txs[ID]
		avg := amt / tx
		fmt.Printf("Rank %d)  %s - %s (%s, %s):\n\tTotal $: %.2f\t# Txs: %f\tAvg Tx $: %.2f\n", i+1, ID, sd.Name, sd.City, sd.State, amt, tx, avg)
	}
	fmt.Println()

	return nil
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
		ids = append(ids, e.ID)
	}
	sds, err := indexing.LookupSearchData(ids)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("printSortedRankings failed: %v", err)
	}
	for i, sd := range sds {
		fmt.Printf("Rank %d)  %s - %s (%s, %s): %.2f\n", i, sd.ID, sd.Name, sd.City, sd.State, r.Amts[sd.ID])
		if sd.Bucket == "individuals" && i == 499 {
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
