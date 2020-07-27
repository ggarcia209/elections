package testDB

import (
	"fmt"

	"github.com/elections/donations"
)

// SUCCESS
func TestCompareTopOverall() error {
	od1 := &donations.TopOverallData{
		Category:  "top_indv",
		Amts:      map[string]float32{"indv00": 100, "indv01": 200, "indv02": 50, "indv04": 150, "indv05": 250},
		Threshold: []*donations.Entry{},
		SizeLimit: 5,
	}
	od2 := &donations.TopOverallData{
		Category:  "top_indv",
		Amts:      map[string]float32{"indv00": 100, "indv01": 200, "indv02": 50},
		Threshold: []*donations.Entry{},
		SizeLimit: 5,
	}
	e := &donations.Entry{ID: "indv03", Total: 75}

	err := compareTopOverall(e, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("0: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	err = compareTopOverall(e, od2)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("0: od2")
	fmt.Println(od2)
	fmt.Println()

	e2 := &donations.Entry{ID: "indv06", Total: 300}
	err = compareTopOverall(e2, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("1: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e3 := &donations.Entry{ID: "indv07", Total: 120}
	err = compareTopOverall(e3, od1) // call removes lowest value but fails to update threshold
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("2: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e4 := &donations.Entry{ID: "indv08", Total: 135}
	err = compareTopOverall(e4, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("3: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e5 := &donations.Entry{ID: "indv09", Total: 175}
	err = compareTopOverall(e5, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("4: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e6 := &donations.Entry{ID: "indv10", Total: 210}
	err = compareTopOverall(e6, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("5: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	e7 := &donations.Entry{ID: "indv11", Total: 500}
	err = compareTopOverall(e7, od1)
	if err != nil {
		return fmt.Errorf("compareTopOverall failed: %v", err)
	}

	fmt.Println("6: od1")
	fmt.Println(od1)
	printODThreshold(od1.Threshold)
	fmt.Println()

	return nil
}

// SUCCESS
func TestUpdateAndSave() error {
	year := "2018"
	category := "cmte_recs_all"

	// Create test record in DbSim
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	fmt.Println("DbSim before update: ")
	fmt.Println(DbSim[year]["top_overall"][category])
	fmt.Println()

	// 0: first entry
	e0 := &donations.Entry{ID: "cmte00", Total: 0}
	err := updateAndSave(year, category, e0)
	if err != nil {
		fmt.Println("TestUpdateAndSave failed: ", err)
		return err
	}

	// 1: below size limit - direct add
	e1 := &donations.Entry{ID: "cmte01", Total: 10}
	err = updateAndSave(year, category, e1)
	if err != nil {
		fmt.Println("TestUpdateAndSave failed: ", err)
		return err
	}

	// 2: below size limit - direct add
	e2 := &donations.Entry{ID: "cmte02", Total: 20}
	err = updateAndSave(year, category, e2)
	if err != nil {
		fmt.Println("TestUpdateAndSave failed: ", err)
		return err
	}

	// 3: at size limit - create threshold / delete smallest value
	e3 := &donations.Entry{ID: "cmte03", Total: 30}
	err = updateAndSave(year, category, e3)
	if err != nil {
		fmt.Println("TestUpdateAndSave failed: ", err)
		return err
	}

	// 4: at limit - new value within threshold range
	e4 := &donations.Entry{ID: "cmte04", Total: 15}
	err = updateAndSave(year, category, e4)
	if err != nil {
		fmt.Println("TestUpdateAndSave failed: ", err)
		return err
	}

	// 5: at limit - new value greater than threshold range
	e5 := &donations.Entry{ID: "cmte05", Total: 50}
	err = updateAndSave(year, category, e5)
	if err != nil {
		fmt.Println("TestUpdateAndSave failed: ", err)
		return err
	}

	// 6: at limit - greater than threshold range
	e6 := &donations.Entry{ID: "cmte06", Total: 40}
	err = updateAndSave(year, category, e6)
	if err != nil {
		fmt.Println("TestUpdateAndSave failed: ", err)
		return err
	}

	// 7: at limit - instantiate new threshold range to check value - below range
	e7 := &donations.Entry{ID: "cmte07", Total: 10}
	err = updateAndSave(year, category, e7)
	if err != nil {
		fmt.Println("TestUpdateAndSave failed: ", err)
		return err
	}

	// 8: at limit - new value within new threshold range
	e8 := &donations.Entry{ID: "cmte08", Total: 35}
	err = updateAndSave(year, category, e8)
	if err != nil {
		fmt.Println("TestUpdateAndSave failed: ", err)
		return err
	}

	fmt.Println("DbSim after update:")
	fmt.Println(DbSim[year]["top_overall"][category])
	fmt.Println()

	return nil
}

// SUCCESS
func TestPartySwitchCases() error {
	test := func(party string) {
		// party specific committees
		switch {
		case party == "REP":
			// republican commitees
			category := "r"
			fmt.Println(category)
		case party == "DEM":
			// democrat committees
			category := "d"
			fmt.Println(category)
		case party == "IND" || party == "N" || party == "NPA" || party == "NOP" || party == "NNE" || party == "UN":
			// independent/non-affiliated committees
			category := "na"
			fmt.Println(category)
		default:
			// all other parties
			category := "misc"
			fmt.Println(category)
		}
	}

	test("REP")
	test("DEM")
	test("IND")
	test("N")
	test("NPA")
	test("NOP")
	test("NNE")
	test("UN")
	test("COM")
	test("")
	test("dem")

	return nil
}

// SUCCESS
func TestTopOverallInternalLogic() error {
	indv := &donations.Individual{}
	org := &donations.Organization{}
	cmte := &donations.CmteTxData{}
	cand := &donations.Candidate{}
	def := struct{}{}

	test := func(other interface{}, incoming, transfer bool, candID string) {
		switch t := other.(type) {
		case *donations.Individual:
			if incoming {
				fmt.Println("indvidual - incoming")
			} else {
				fmt.Println("indvidual - !incoming")
			}

		case *donations.Organization:
			if incoming {
				fmt.Println("org - incoming")
			} else {
				fmt.Println("org - !incoming")
			}
		case *donations.CmteTxData:
			if incoming {
				if transfer {
					fmt.Println("cmte - incoming - transfer")
				} else {
					fmt.Println("cmte - incoming - !transfer")
				}
			} else {
				if transfer {
					fmt.Println("cmte - !incoming - transfer")
				} else {
					fmt.Println("cmte - !incoming - !transfer")
				}
			}
		case *donations.Candidate:
			if incoming {
				fmt.Println("cand - incoming")
			} else {
				if transfer {
					fmt.Println("cand - !incoming - transfer")
				} else {
					fmt.Println("cand - !incoming - !transfer")
				}
			}
		default:
			_ = t
			fmt.Println("invalid interface")
		}

		if candID != "" {
			fmt.Println("linked cand")
		}

	}

	// individual cases
	test(indv, false, false, "cand00") // "indvidual - !incoming"
	test(indv, true, false, "")        // "indvidual - incoming"

	/// org cases
	test(org, true, false, "can00") // "org - incoming"
	test(org, false, false, "")     // "org - !incoming"

	// cmte cases
	test(cmte, true, true, "cand00") // "cmte - incoming - transfer"
	test(cmte, true, false, "")      // "cmte - incoming - !transfer"
	test(cmte, false, true, "")      // "cmte - !incoming - transfer"
	test(cmte, false, false, "")     // "cmte - !incoming - !transfer"

	// cand cases
	test(cand, true, false, "")       // "cand - incoming"
	test(cand, false, true, "cand00") // "cand - !incomoing - transfer"
	test(cand, false, false, "")      // "cand - !incomoing - !transfer"

	// default case
	test(def, false, false, "") // invalid

	return nil
}

// SUCESS
func TestUpdateTopOverall() error {
	year := "2018"

	/* Top Overall Objects */
	// Create test record in DbSim
	category := "cmte_recs_all"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_recs_na"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_recs_r"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	// Create test record in DbSim
	category = "cmte_donors_all"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cmte_donors_na"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	// Create test record in DbSim
	category = "cand_all"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}
	// Create test record in DbSim
	category = "cand_na"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	// Create test record in DbSim
	category = "indv"
	DbSim[year]["top_overall"][category] = &donations.TopOverallData{
		Category:  category,
		Amts:      make(map[string]float32),
		Threshold: []*donations.Entry{},
		SizeLimit: 3,
	}

	/* Filer / Other objects */
	// filer w/o linked cand
	filer00 := &donations.CmteTxData{
		CmteID:           "cmte00",
		CandID:           "",
		Party:            "IND",
		TotalIncomingAmt: 1000,
	}

	// filer w/ linked cand
	filer01 := &donations.CmteTxData{
		CmteID:           "cmte01",
		CandID:           "Pcand01", // ../test_vars.go/Cand01
		Party:            "IND",
		TotalIncomingAmt: 1200,
	}

	// donor
	indv00 := &donations.Individual{
		ID:          "indv00",
		TotalOutAmt: 250,
	}

	// donor
	indv01 := &donations.Individual{
		ID:          "indv01",
		TotalOutAmt: 200,
	}

	// donor
	indv02 := &donations.Individual{
		ID:          "indv02",
		TotalOutAmt: 100,
	}

	// donor
	indv03 := &donations.Individual{
		ID:          "indv03",
		TotalOutAmt: 150,
	}

	// recipient
	cmte02 := &donations.CmteTxData{
		CmteID:       "cmte02",
		CandID:       "",
		Party:        "REP",
		TransfersAmt: 2000,
	}

	/* Tests */

	// Test A: Cmte w/o linked cand / individual donors
	fmt.Println("***** TEST A: Cmte w/o cand / indv donors *****")

	filer00.TotalIncomingAmt += indv00.TotalOutAmt
	err := updateTopOverall(year, filer00, indv00, true, false)
	if err != nil {
		fmt.Println("updateTopOverall failed: ", err)
		return fmt.Errorf("updateTopOverall failed: %v", err)
	}
	fmt.Println("A0: ")
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_na"])
	fmt.Println(DbSim[year]["top_overall"]["indv"])
	fmt.Println()

	filer00.TotalIncomingAmt += indv01.TotalOutAmt
	err = updateTopOverall(year, filer00, indv01, true, false)
	if err != nil {
		fmt.Println("updateTopOverall failed: ", err)
		return fmt.Errorf("updateTopOverall failed: %v", err)
	}
	fmt.Println("A1: ")
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_na"])
	fmt.Println(DbSim[year]["top_overall"]["indv"])
	fmt.Println()

	filer00.TotalIncomingAmt += indv02.TotalOutAmt
	err = updateTopOverall(year, filer00, indv02, true, false)
	if err != nil {
		fmt.Println("updateTopOverall failed: ", err)
		return fmt.Errorf("updateTopOverall failed: %v", err)
	}
	fmt.Println("A2: ")
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_na"])
	fmt.Println(DbSim[year]["top_overall"]["indv"])
	fmt.Println()

	filer00.TotalIncomingAmt += indv03.TotalOutAmt
	err = updateTopOverall(year, filer00, indv03, true, false)
	if err != nil {
		fmt.Println("updateTopOverall failed: ", err)
		return fmt.Errorf("updateTopOverall failed: %v", err)
	}
	fmt.Println("A3: ")
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_na"])
	fmt.Println(DbSim[year]["top_overall"]["indv"])
	fmt.Println()

	fmt.Println("***** END TEST A *****")

	// Test B: Cmte w/ linked cand / individual donors
	fmt.Println("***** TEST B: Cmte w/o cand / indv donors *****")

	filer01.TotalIncomingAmt += indv00.TotalOutAmt
	indv00.TotalOutAmt += indv00.TotalOutAmt
	err = updateTopOverall(year, filer01, indv00, true, false)
	if err != nil {
		fmt.Println("updateTopOverall failed: ", err)
		return fmt.Errorf("updateTopOverall failed: %v", err)
	}
	fmt.Println("B0: ")
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_na"])
	fmt.Println(DbSim[year]["top_overall"]["cand_all"])
	fmt.Println(DbSim[year]["top_overall"]["cand_na"])
	fmt.Println(DbSim[year]["top_overall"]["indv"])
	fmt.Println()

	filer01.TotalIncomingAmt += indv01.TotalOutAmt
	indv01.TotalOutAmt += indv01.TotalOutAmt
	err = updateTopOverall(year, filer01, indv01, true, false)
	if err != nil {
		fmt.Println("updateTopOverall failed: ", err)
		return fmt.Errorf("updateTopOverall failed: %v", err)
	}
	fmt.Println("B1: ")
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_na"])
	fmt.Println(DbSim[year]["top_overall"]["cand_all"])
	fmt.Println(DbSim[year]["top_overall"]["cand_na"])
	fmt.Println(DbSim[year]["top_overall"]["indv"])
	fmt.Println()

	filer01.TotalIncomingAmt += indv02.TotalOutAmt
	indv02.TotalOutAmt += indv02.TotalOutAmt
	err = updateTopOverall(year, filer01, indv02, true, false)
	if err != nil {
		fmt.Println("updateTopOverall failed: ", err)
		return fmt.Errorf("updateTopOverall failed: %v", err)
	}
	fmt.Println("B2: ")
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_na"])
	fmt.Println(DbSim[year]["top_overall"]["cand_all"])
	fmt.Println(DbSim[year]["top_overall"]["cand_na"])
	fmt.Println(DbSim[year]["top_overall"]["indv"])
	fmt.Println()

	filer01.TotalIncomingAmt += indv03.TotalOutAmt
	indv03.TotalOutAmt += indv03.TotalOutAmt
	err = updateTopOverall(year, filer01, indv03, true, false)
	if err != nil {
		fmt.Println("updateTopOverall failed: ", err)
		return fmt.Errorf("updateTopOverall failed: %v", err)
	}
	fmt.Println("B3: ")
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_na"])
	fmt.Println(DbSim[year]["top_overall"]["cand_all"])
	fmt.Println(DbSim[year]["top_overall"]["cand_na"])
	fmt.Println(DbSim[year]["top_overall"]["indv"])
	fmt.Println()

	fmt.Println("***** END TEST B *****")

	// Test B: Cmte w/ linked cand / individual donors
	fmt.Println("***** TEST C: Outgoing transfer to recipient committee *****")

	filer00.TotalIncomingAmt += cmte02.TransfersAmt
	err = updateTopOverall(year, filer00, cmte02, false, true)
	if err != nil {
		fmt.Println("updateTopOverall failed: ", err)
		return fmt.Errorf("updateTopOverall failed: %v", err)
	}
	fmt.Println("C0: ")
	fmt.Println(DbSim[year]["top_overall"]["cmte_donors_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_donors_na"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_all"])
	fmt.Println(DbSim[year]["top_overall"]["cmte_recs_r"])
	fmt.Println()

	fmt.Println("***** END TEST C *****")

	return nil
}

func printODThreshold(th []*donations.Entry) {
	for i, e := range th {
		fmt.Printf("%d: ID: %s\tTotal: %v\n", i, e.ID, e.Total)
	}
}

// compare obj to top overall threshold
// 7/26/20 - refactored to include checkODThreshold logic
func compareTopOverall(e *donations.Entry, od *donations.TopOverallData) error {
	// add to Amts map if len(Amts) < Size Limit
	if len(od.Amts) < od.SizeLimit {
		od.Amts[e.ID] = e.Total
		return nil
	}

	// check threshold when updating existing entry
	if od.Amts[e.ID] != 0 {
		od.Amts[e.ID] = e.Total
		th, err := checkODThreshold(e.ID, od.Amts, od.Threshold)
		if err != nil {
			fmt.Println("compareTopOverall failed: ", err)
			return fmt.Errorf("compareTopOverall failed: %v", err)
		}
		od.Threshold = th
		return nil
	}

	// if len(Amts) == SizeLimit
	// set/reset least threshold list
	var least Entries
	var err error
	if len(od.Threshold) == 0 {
		es := sortTopX(od.Amts)
		least, err = setThresholdLeast3(es)
		if err != nil {
			fmt.Println("compareTopOverall failed: ", err)
			return fmt.Errorf("compareTopOverall failed: %v", err)
		}
	} else {
		for _, entry := range od.Threshold {
			least = append(least, entry)
		}
	}

	// compare sen cmte's total received value to threshold
	threshold := least[len(least)-1].Total // last/smallest obj in least
	if e.Total > threshold {
		new := newEntry(e.ID, e.Total)
		delID, newEntries := reSortLeast(new, least)
		least = newEntries
		delete(od.Amts, delID)
		od.Amts[e.ID] = e.Total
	} else {
		newTh := []*donations.Entry{}
		for _, e := range least {
			newTh = append(newTh, e)
		}
		od.Threshold = append(od.Threshold[:0], newTh...)
		return nil
	}

	// update threshold
	newTh := []*donations.Entry{}
	for _, e := range least {
		newTh = append(newTh, e)
	}
	od.Threshold = append(od.Threshold[:0], newTh...)

	return nil
}

func updateTopOverall(year string, filer *donations.CmteTxData, other interface{}, incoming, transfer bool) error {
	switch t := other.(type) {
	case *donations.Individual:
		if incoming {
			// update Top Individuals by funds contributed and Top Committees by funds received
			err := updateTopIndividuals(year, other.(*donations.Individual))
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
			err = updateTopCmteRecs(year, filer)
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
		} else {
			// update Top Committes by expenses only - individuals not ranked by funds received
			err := updateTopCmteExp(year, filer)
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
		}
	case *donations.Organization:
		if incoming {
			// update top orgs by contributions and top committees by funds received
			err := updateTopOrgsByContributions(year, other.(*donations.Organization))
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
			err = updateTopCmteRecs(year, filer)
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
		} else {
			// update top orgs by funds received and top committees by expenses
			err := updateTopOrgsByReceipts(year, other.(*donations.Organization))
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
			err = updateTopCmteExp(year, filer)
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
		}
	case *donations.Candidate:
		// update Top Candidates by total funds received in separate function
		if incoming {
			// update top committees by funds received
			err := updateTopCmteRecs(year, filer)
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
		} else {
			if transfer {
				// update top committees by funds contributed
				err := updateTopCmteDonors(year, filer)
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			} else {
				// update top committees by expenses
				err := updateTopCmteExp(year, filer)
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			}
		}
	case *donations.CmteTxData:
		// 7/26/20 - may need to refactor/remove transfer cases - all transactions between committees treated as transactions by default
		if incoming {
			// update top committees by funds received (filer)
			err := updateTopCmteRecs(year, filer)
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
			if transfer {
				// update top committees by funds contributed (sender)
				err := updateTopCmteDonors(year, other.(*donations.CmteTxData))
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			} else {
				// update top committees by expenses (sender)
				err := updateTopCmteExp(year, other.(*donations.CmteTxData))
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			}
		} else {
			if transfer {
				// update top committees by funds contributed (filer)
				err := updateTopCmteDonors(year, filer)
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			} else {
				// update top committees by expenses (filer)
				err := updateTopCmteExp(year, filer)
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			}
			// update top committees by funds received (receiver)
			err := updateTopCmteRecs(year, other.(*donations.CmteTxData))
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
		}
	default:
		_ = t // discard unused variable
		return fmt.Errorf("updateTopOverall failed: wrong interface type")
	}

	// update top candidates by funds received/transferred if candidate linked to filing committee
	if filer.CandID != "" {
		// get linked candidate
		cand := DbSim[year]["candidates"][filer.CandID]

		// update top candidates by total funds incoming/outgoing
		err := updateTopCandidates(year, cand.(*donations.Candidate), filer, incoming)
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}
	}

	return nil
}

func updateTopCandidates(year string, cand *donations.Candidate, pcc *donations.CmteTxData, incoming bool) error {
	if incoming {
		err := updateTopCandsByIncoming(year, cand, pcc)
		if err != nil {
			fmt.Println("updateTopCandidates failed: ", err)
			return fmt.Errorf("updateTopCandidates failed: %v", err)
		}
	} else {
		err := updateTopCandsByOutgoing(year, cand, pcc)
		if err != nil {
			fmt.Println("updateTopCandidates failed: ", err)
			return fmt.Errorf("updateTopCandidates failed: %v", err)
		}
	}
	return nil
}

// top committees by contributions received
func updateTopCmteRecs(year string, cmte *donations.CmteTxData) error {
	entry := &donations.Entry{ID: cmte.CmteID, Total: cmte.TotalIncomingAmt}

	// all committees
	category := "cmte_recs_all"
	err := updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopCmteRecs failed: ", err)
		return fmt.Errorf("updateTopCmteRecs failed: %v", err)
	}

	// party specific committees
	switch {
	case cmte.Party == "REP":
		// republican commitees
		category = "cmte_recs_r"
	case cmte.Party == "DEM":
		// democrat committees
		category = "cmte_recs_d"
	case cmte.Party == "IND" || cmte.Party == "N" || cmte.Party == "NPA" || cmte.Party == "NOP" || cmte.Party == "NNE" || cmte.Party == "UN":
		// independent/non-affiliated committees
		category = "cmte_recs_na"
	default:
		// all other parties
		category = "cmte_recs_misc"
	}

	err = updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopCmteRecs failed: ", err)
		return fmt.Errorf("updateTopCmteRecs failed: %v", err)
	}
	return nil
}

func updateTopCmteDonors(year string, cmte *donations.CmteTxData) error {
	entry := &donations.Entry{ID: cmte.CmteID, Total: cmte.TransfersAmt}

	// all committees
	category := "cmte_donors_all"
	err := updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopCmteDonors failed: ", err)
		return fmt.Errorf("updateTopCmteDonors failed: %v", err)
	}

	// party specific committees
	switch {
	case cmte.Party == "REP":
		// republican commitees
		category = "cmte_donors_r"
	case cmte.Party == "DEM":
		// democrat committees
		category = "cmte_donors_d"
	case cmte.Party == "IND" || cmte.Party == "N" || cmte.Party == "NPA" || cmte.Party == "NOP" || cmte.Party == "NNE" || cmte.Party == "UN":
		// independent/non-affiliated committees
		category = "cmte_donors_na"
	default:
		// all other parties
		category = "cmte_donors_misc"
	}

	err = updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopCmteDonors failed: ", err)
		return fmt.Errorf("updateTopCmteDonors failed: %v", err)
	}

	return nil
}

func updateTopCmteExp(year string, cmte *donations.CmteTxData) error {
	entry := &donations.Entry{ID: cmte.CmteID, Total: cmte.ExpendituresAmt}

	// all committees
	category := "cmte_exp_all"
	err := updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopCmteExp failed: ", err)
		return fmt.Errorf("updateTopCmteExp failed: %v", err)
	}

	// party specific committees
	switch {
	case cmte.Party == "REP":
		// republican commitees
		category = "cmte_exp_r"
	case cmte.Party == "DEM":
		// democrat committee
		category = "cmte_exp_d"
	case cmte.Party == "IND" || cmte.Party == "N" || cmte.Party == "NPA" || cmte.Party == "NOP" || cmte.Party == "NNE" || cmte.Party == "UN":
		// independent/non-affiliated committees
		category = "cmte_exp_na"
	default:
		// all other parties
		category = "cmte_exp_misc"
	}

	err = updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopCmteExp failed: ", err)
		return fmt.Errorf("updateTopCmteExp failed: %v", err)
	}

	return nil
}

// Top Candidates by total funds received
func updateTopCandsByIncoming(year string, cand *donations.Candidate, pcc *donations.CmteTxData) error {
	entry := &donations.Entry{ID: cand.ID, Total: cand.TotalDirectInAmt + pcc.TotalIncomingAmt}

	// all
	category := "cand_all"
	err := updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopCandidates failed: ", err)
		return fmt.Errorf("updateTopCandidates failed: %v", err)
	}

	// party specific
	switch {
	case cand.Party == "REP":
		// republican
		category = "cand_r"
	case cand.Party == "DEM":
		// democrat
		category = "cand_d"
	case cand.Party == "IND" || cand.Party == "N" || cand.Party == "NPA" || cand.Party == "NOP" || cand.Party == "NNE" || cand.Party == "UN":
		// independent/non-affiliated
		category = "cand_na"
	default:
		// all other parties
		category = "cand_misc"
	}

	err = updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopCandidates failed: ", err)
		return fmt.Errorf("updateTopCandidates failed: %v", err)
	}
	return nil
}

// Top Candidates by total funds disbursed
func updateTopCandsByOutgoing(year string, cand *donations.Candidate, pcc *donations.CmteTxData) error {
	entry := &donations.Entry{ID: cand.ID, Total: cand.TotalDirectOutAmt + pcc.TotalOutgoingAmt}

	// all
	category := "cand_exp_all"
	err := updateAndSave(year, category, entry)

	// party specific
	switch {
	case cand.Party == "REP":
		// republican
		category = "cand_exp_r"
	case cand.Party == "DEM":
		// democrat
		category = "cand_exp_d"
	case cand.Party == "IND" || cand.Party == "N" || cand.Party == "NPA" || cand.Party == "NOP" || cand.Party == "NNE" || cand.Party == "UN":
		// independent/non-affiliated
		category = "cand_exp_na"
	default:
		// all other parties
		category = "cand_exp_misc"
	}

	err = updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopCandExp failed: ", err)
		return fmt.Errorf("updateTopCandExp failed: %v", err)
	}
	return nil
}

// Top Individual donors by total funds contributed
func updateTopIndividuals(year string, indv *donations.Individual) error {
	entry := &donations.Entry{ID: indv.ID, Total: indv.TotalOutAmt}

	category := "indv"
	err := updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopIndv failed: ", err)
		return fmt.Errorf("updateTopIndv failed: %v", err)
	}

	return nil
}

// Top Organizations by funds contributed
func updateTopOrgsByContributions(year string, org *donations.Organization) error {
	entry := &donations.Entry{ID: org.ID, Total: org.TotalOutAmt}

	category := "org_conts"
	err := updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopOrganizations failed: ", err)
		return fmt.Errorf("updateTopOrganizations failed: %v", err)
	}

	return nil
}

// Top Organizations by funds received
func updateTopOrgsByReceipts(year string, org *donations.Organization) error {
	entry := &donations.Entry{ID: org.ID, Total: org.TotalInAmt}

	category := "org_recs"
	err := updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopOrganizations failed: ", err)
		return fmt.Errorf("updateTopOrganizations failed: %v", err)
	}

	return nil
}

// get TopOverallData obj per year/category, update, & save
func updateAndSave(year, category string, entry *donations.Entry) error {
	// "top_overall" objects are initialzed before first call to function
	// and are never returned as nil objects
	od := DbSim[year]["top_overall"][category]
	err := compareTopOverall(entry, od.(*donations.TopOverallData))
	if err != nil {
		fmt.Println("updateAndSave failed: ", err)
		return fmt.Errorf("updateAndSave failed: %v", err)
	}

	// TEST ONLY
	// fmt.Println("Overall Data:")
	// fmt.Println(od)
	// printODThreshold(od.(*donations.TopOverallData).Threshold)
	// fmt.Println()

	DbSim[year]["top_overall"][category] = od
	return nil
}

// get TopOverallData obj per year/category, update, & save
/* func updateAndSave(year, category string, entry *donations.Entry) error {
	// "top_overall" objects are initialzed before first call to function
	// and are never returned as nil objects
	od, err := persist.GetObject(year, "top_overall", category)
	if err != nil {
		fmt.Println("updateAndSave failed: ", err)
		return fmt.Errorf("updateAndSave failed: %v", err)
	}
	err = compareTopOverall(entry, od.(*donations.TopOverallData))
	if err != nil {
		fmt.Println("updateAndSave failed: ", err)
		return fmt.Errorf("updateAndSave failed: %v", err)
	}
	err = persist.PutObject(year, od)
	if err != nil {
		fmt.Println("updateAndSave failed: ", err)
		return fmt.Errorf("updateAndSave failed: %v", err)
	}
	return nil
} */

// check to see if previous total of entry is in threshold range when updating existing entry
func checkODThreshold(newID string, m map[string]float32, th []*donations.Entry) ([]*donations.Entry, error) {
	inRange := false
	check := map[string]bool{newID: true}
	for _, e := range th {
		if check[e.ID] == true {
			inRange = true
		}
	}
	if inRange {
		es := sortTopX(m)
		newRange, err := setThresholdLeast3(es)
		if err != nil {
			fmt.Println("checkODThreshold failed: ", err)
			return []*donations.Entry{}, fmt.Errorf("checkODThreshold failed: %v", err)
		}
		// update object's threshold list
		newTh := []*donations.Entry{}
		for _, entry := range newRange {
			newTh = append(newTh, entry)
		}
		return newTh, nil
	}
	return th, nil
}
