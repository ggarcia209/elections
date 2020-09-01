package admin

import (
	"fmt"

	"github.com/elections/source/databuilder"
	"github.com/elections/source/donations"
	"github.com/elections/source/persist"
	"github.com/elections/source/ui"
)

type odMapping map[string]map[string]map[string]*donations.TopOverallData
type ytMapping map[string]map[string]*donations.YearlyTotal

// CreateSecondaryDatasets processes objects created from raw data
// and creates the TopOverall and YearlyTotals datasets
func createSecondaryDatasets() error {
	fmt.Println("***** PROCESS SECONDARY DATA *****")
	path, err := getPath(false)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("createSecondaryDatasets failed: %v", err)
	}
	year := ui.GetYear()
	buckets := []string{"individuals", "cmte_tx_data", "candidates"}
	persist.OUTPUT_PATH = path

	// check if objects already exist for idempotency
	startCheck, err := persist.GetTopOverall(year)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ProcessNewRecords failed: %v", err)
	}
	if len(startCheck) != 0 {
		fmt.Println("Secondary already data exists. Overwrite with new data?")
		yes := ui.Ask4confirm()
		if !yes {
			fmt.Println("Returning to menu...")
			return nil
		}
	}

	// initialize TopOverallData objects & mappings
	topOverall, yearlyTotals := donations.InitSecondaryDataObjs(year)

	odMap := make(odMapping)
	for _, intf := range topOverall {
		od := intf.(*donations.TopOverallData)
		if odMap[od.Bucket] == nil {
			odMap[od.Bucket] = make(map[string]map[string]*donations.TopOverallData)
		}
		if odMap[od.Bucket][od.Category] == nil {
			odMap[od.Bucket][od.Category] = make(map[string]*donations.TopOverallData)
		}
		odMap[od.Bucket][od.Category][od.ID] = od
	}

	ytMap := make(ytMapping)
	for _, intf := range yearlyTotals {
		yt := intf.(*donations.YearlyTotal)
		if ytMap[yt.Category] == nil {
			ytMap[yt.Category] = make(map[string]*donations.YearlyTotal)
		}
		ytMap[yt.Category][yt.ID] = yt
	}

	// derive data
	fmt.Println("Creating Top Overall Rankings and Yearly Totals...")
	for _, b := range buckets {
		fmt.Println("Processing bucket: ", b)
		err := deriveDatabyBucket(year, b, odMap, ytMap)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("createSecondaryDatasets failed: %v", err)
		}
	}
	fmt.Println("Top Overall Rankings  and Yearly Totals complete!")

	return nil
}

// deriveDatabyCat finds and records the top rankings for each party for a given year/bucket/category
func deriveDatabyBucket(year, bucket string, odm odMapping, ytm ytMapping) error {
	cats := []string{"rec", "donor", "exp"}

	// scan every object each category
	// update TopOverall/YearlyTotal objects
	err := scanObjects(year, bucket, odm[bucket], ytm)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("DeriveTopOverall failed: %v", err)
	}

	for _, cat := range cats {
		ods := []interface{}{}
		yts := []interface{}{}
		for _, od := range odm[bucket][cat] {
			ods = append(ods, od)
		}
		for _, yt := range ytm[cat] {
			yts = append(yts, yt)
		}

		// save TopOverall objects
		// overwrite any previously existing data
		err = persist.SaveTopOverall(year, bucket, ods)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("DeriveTopOverall failed: %v", err)
		}
		err = persist.SaveYearlyTotals(year, cat, yts)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("DeriveTopOverall failed: %v", err)
		}
	}

	return nil
}

// scan each object and update TopRankings/Yearly Totals for each object
func scanObjects(year, bucket string, ods map[string]map[string]*donations.TopOverallData, yts map[string]map[string]*donations.YearlyTotal) error {
	n := 1000
	j := 0
	start := ""
	curr := start
	cmteTotals := map[string]map[string]float32{
		"rec":   make(map[string]float32),
		"donor": make(map[string]float32),
		"exp":   make(map[string]float32),
	}
	cats := []string{"rec", "donor", "exp"}

	for {
		objs, key, err := persist.BatchGetSequential(year, bucket, curr, n)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("scanObjects failed: %v", err)
		}
		curr = key

		// add funds raised by candidate's PCC to direct amts
		if bucket == "candidates" {
			// get corresponding CmteTxData for each candidate
			ids := []string{}
			for _, obj := range objs {
				ids = append(ids, obj.(*donations.Candidate).PCC)
			}
			cmtes, _, err := persist.BatchGetByID(year, "cmte_tx_data", ids)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("scanObjects failed: %v", err)
			}
			// get total for each cmte
			for _, c := range cmtes {
				for _, cat := range cats {
					id, _, total, err := deriveTotal(c, cat)
					if err != nil {
						fmt.Println(err)
						return fmt.Errorf("scanObjects failed: %v", err)
					}
					cmteTotals[cat][id] = total
				}
			}
		}

		// get totals and add/compare to rankings list
		// update yearly totals
		for _, obj := range objs {
			for _, cat := range cats {
				if bucket == "individuals" && cat == "exp" {
					continue // non-existent category
				}
				// update ALL
				all := ods[cat][year+"-"+bucket+"-"+cat+"-ALL"]
				id, pty, total, err := deriveTotal(obj, cat)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("scanObjects failed: %v", err)
				}
				if bucket == "candidates" {
					total += cmteTotals[cat][obj.(*donations.Candidate).PCC]
				}
				err = databuilder.CompareTopOverall(id, total, all)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("scanObjects failed: %v", err)
				}

				if bucket == "individuals" {
					continue // no party specific objects
				}

				// update Party-specific
				ptyOd := ods[cat][year+"-"+bucket+"-"+cat+"-"+pty]
				err = databuilder.CompareTopOverall(id, total, ptyOd)
				if err != nil {
					fmt.Println(err)
					return fmt.Errorf("scanObjects failed: %v", err)
				}

				// update yearly totals while processing cmtes
				if bucket == "cmte_tx_data" {
					ytAll := yts[cat][year+"-"+cat+"-ALL"]
					databuilder.UpdateYearlyTotal(total, ytAll)
					ytPty := yts[cat][year+"-"+cat+"-"+pty]
					databuilder.UpdateYearlyTotal(total, ytPty)
				}
				j++
			}

		}

		if len(objs) < n {
			j += len(objs)
			fmt.Println("objects scanned: ", j)
			break
		}

	}
	return nil
}

func deriveTotal(obj interface{}, cat string) (string, string, float32, error) {
	var ID string
	var pty string
	switch t := obj.(type) {
	case *donations.Individual:
		ID = obj.(*donations.Individual).ID
		pty = "ALL"
		switch {
		case cat == "rec":
			return ID, pty, obj.(*donations.Individual).TotalInAmt, nil
		case cat == "donor":
			return ID, pty, obj.(*donations.Individual).TotalOutAmt, nil
		default:
			return "", "", 0, fmt.Errorf("deriveTotal failed: invalid category")
		}
	case *donations.CmteTxData:
		ID = obj.(*donations.CmteTxData).CmteID
		pty = getParty(obj.(*donations.CmteTxData).Party)
		switch {
		case cat == "rec":
			return ID, pty, obj.(*donations.CmteTxData).TotalIncomingAmt, nil
		case cat == "donor":
			return ID, pty, obj.(*donations.CmteTxData).TransfersAmt, nil
		case cat == "exp":
			return ID, pty, obj.(*donations.CmteTxData).ExpendituresAmt, nil
		default:
			return "", "", 0, fmt.Errorf("deriveTotal failed: invalid category")
		}
	case *donations.Candidate:
		ID = obj.(*donations.Candidate).ID
		pty = getParty(obj.(*donations.Candidate).Party)
		switch {
		case cat == "rec":
			return ID, pty, obj.(*donations.Candidate).TotalDirectInAmt, nil
		case cat == "donor":
			return ID, pty, obj.(*donations.Candidate).TotalDirectOutAmt, nil
		case cat == "exp":
			return ID, pty, 0, nil // special case - derive from CmteTxData.ExpendituresAmt only
		default:
			return "", "", 0, fmt.Errorf("deriveTotal failed: invalid category")
		}
	default:
		_ = t
		return "", "", 0, fmt.Errorf("deriveTotal failed: invalid interface type")
	}
}

func getParty(pty string) string {
	// party specific
	switch {
	case pty == "REP":
		return pty
	case pty == "DEM":
		return pty
	case pty == "IND" || pty == "N" || pty == "NPA" || pty == "NOP" || pty == "NNE" || pty == "UN":
		return "IND"
	case pty == "":
		return "UNK"
	default:
		return "OTH"
	}
}
