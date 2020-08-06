package databuilder

// Refactored 8/3/20 - removed Organization functions, cases, and references

import (
	"fmt"

	"github.com/elections/donations"
)

// compare obj to top overall threshold
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
		least, err = setThresholdLeast10(es)
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

func updateTopOverall(filer *donations.CmteTxData, other interface{}, incoming, transfer bool, cache map[string]map[string]interface{}) error {
	switch t := other.(type) {
	case *donations.Individual:
		if incoming {
			// update Top Individuals by funds contributed and Top Committees by funds received
			err := updateTopIndividuals(other.(*donations.Individual), cache["top_overall"])
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
			err = updateTopCmteRecs(filer, cache["top_overall"])
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
		} else {
			// update Top Committes by expenses only - individuals not ranked by funds received
			err := updateTopCmteExp(filer, cache["top_overall"])
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
		}
	case *donations.Candidate:
		// update Top Candidates by total funds received in separate function
		if incoming {
			// update top committees by funds received
			err := updateTopCmteRecs(filer, cache["top_overall"])
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
		} else {
			if transfer {
				// update top committees by funds contributed
				err := updateTopCmteDonors(filer, cache["top_overall"])
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			} else {
				// update top committees by expenses
				err := updateTopCmteExp(filer, cache["top_overall"])
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
			err := updateTopCmteRecs(filer, cache["top_overall"])
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
			if transfer {
				// update top committees by funds contributed (sender)
				err := updateTopCmteDonors(other.(*donations.CmteTxData), cache["top_overall"])
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			} else {
				// update top committees by expenses (sender)
				err := updateTopCmteExp(other.(*donations.CmteTxData), cache["top_overall"])
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			}
		} else {
			if transfer {
				// update top committees by funds contributed (filer)
				err := updateTopCmteDonors(filer, cache["top_overall"])
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			} else {
				// update top committees by expenses (filer)
				err := updateTopCmteExp(filer, cache["top_overall"])
				if err != nil {
					fmt.Println("updateTopOverall failed: ", err)
					return fmt.Errorf("updateTopOverall failed: %v", err)
				}
			}
			// update top committees by funds received (receiver)
			err := updateTopCmteRecs(other.(*donations.CmteTxData), cache["top_overall"])
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
		cand := cache["candidates"][filer.CandID] // DbSim[year]["candidates"][filer.CandID]

		// update top candidates by total funds incoming/outgoing
		err := updateTopCandidates(cand.(*donations.Candidate), filer, incoming, cache["top_overall"])
		if err != nil {
			fmt.Println("ContributionUpdate failed: ", err)
			return fmt.Errorf("ContributionUpdate failed: %v", err)
		}
	}

	return nil
}

func updateTopCandidates(cand *donations.Candidate, pcc *donations.CmteTxData, incoming bool, cache map[string]interface{}) error {
	if incoming {
		err := updateTopCandsByIncoming(cand, pcc, cache)
		if err != nil {
			fmt.Println("updateTopCandidates failed: ", err)
			return fmt.Errorf("updateTopCandidates failed: %v", err)
		}
	} else {
		err := updateTopCandsByOutgoing(cand, pcc, cache)
		if err != nil {
			fmt.Println("updateTopCandidates failed: ", err)
			return fmt.Errorf("updateTopCandidates failed: %v", err)
		}
	}
	return nil
}

// top committees by contributions received
func updateTopCmteRecs(cmte *donations.CmteTxData, cache map[string]interface{}) error {
	entry := &donations.Entry{ID: cmte.CmteID, Total: cmte.TotalIncomingAmt}

	// all committees
	category := "cmte_recs_all"
	err := updateAndSave(category, entry, cache)
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

	err = updateAndSave(category, entry, cache)
	if err != nil {
		fmt.Println("updateTopCmteRecs failed: ", err)
		return fmt.Errorf("updateTopCmteRecs failed: %v", err)
	}
	return nil
}

func updateTopCmteDonors(cmte *donations.CmteTxData, cache map[string]interface{}) error {
	entry := &donations.Entry{ID: cmte.CmteID, Total: cmte.TransfersAmt}

	// all committees
	category := "cmte_donors_all"
	err := updateAndSave(category, entry, cache)
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

	err = updateAndSave(category, entry, cache)
	if err != nil {
		fmt.Println("updateTopCmteDonors failed: ", err)
		return fmt.Errorf("updateTopCmteDonors failed: %v", err)
	}

	return nil
}

func updateTopCmteExp(cmte *donations.CmteTxData, cache map[string]interface{}) error {
	entry := &donations.Entry{ID: cmte.CmteID, Total: cmte.ExpendituresAmt}

	// all committees
	category := "cmte_exp_all"
	err := updateAndSave(category, entry, cache)
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

	err = updateAndSave(category, entry, cache)
	if err != nil {
		fmt.Println("updateTopCmteExp failed: ", err)
		return fmt.Errorf("updateTopCmteExp failed: %v", err)
	}

	return nil
}

// Top Candidates by total funds received
func updateTopCandsByIncoming(cand *donations.Candidate, pcc *donations.CmteTxData, cache map[string]interface{}) error {
	entry := &donations.Entry{ID: cand.ID, Total: cand.TotalDirectInAmt + pcc.TotalIncomingAmt}

	// all
	category := "cand_all"
	err := updateAndSave(category, entry, cache)
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

	err = updateAndSave(category, entry, cache)
	if err != nil {
		fmt.Println("updateTopCandidates failed: ", err)
		return fmt.Errorf("updateTopCandidates failed: %v", err)
	}
	return nil
}

// Top Candidates by total funds disbursed
func updateTopCandsByOutgoing(cand *donations.Candidate, pcc *donations.CmteTxData, cache map[string]interface{}) error {
	entry := &donations.Entry{ID: cand.ID, Total: cand.TotalDirectOutAmt + pcc.TotalOutgoingAmt}

	// all
	category := "cand_exp_all"
	err := updateAndSave(category, entry, cache)

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

	err = updateAndSave(category, entry, cache)
	if err != nil {
		fmt.Println("updateTopCandExp failed: ", err)
		return fmt.Errorf("updateTopCandExp failed: %v", err)
	}
	return nil
}

// Top Individual donors by total funds contributed
func updateTopIndividuals(indv *donations.Individual, cache map[string]interface{}) error {
	entry := &donations.Entry{ID: indv.ID, Total: indv.TotalOutAmt}

	category := "indv"
	err := updateAndSave(category, entry, cache)
	if err != nil {
		fmt.Println("updateTopIndv failed: ", err)
		return fmt.Errorf("updateTopIndv failed: %v", err)
	}

	return nil
}

// get TopOverallData obj per year/category, update, & save
func updateAndSave(category string, entry *donations.Entry, cache map[string]interface{}) error {
	// "top_overall" objects are initialzed before first call to function
	// and are never returned as nil objects
	od := cache[category]

	err := compareTopOverall(entry, od.(*donations.TopOverallData))
	if err != nil {
		fmt.Println("updateAndSave failed: ", err)
		return fmt.Errorf("updateAndSave failed: %v", err)
	}

	return nil
}

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
		newRange, err := setThresholdLeast10(es)
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
