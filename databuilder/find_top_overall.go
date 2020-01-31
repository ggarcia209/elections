package databuilder

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/persist"
)

func updateTopOverall(e *donations.Entry, od *donations.TopOverallData) error {
	// add to Amts map if len(Amts) < Size Limit
	if len(od.Amts) < od.SizeLimit {
		od.Amts[e.ID] = e.Total
		return nil
	}

	// if len(Amts) == SizeLimit
	// set/reset least threshold list
	var least Entries
	var err error
	if len(od.Threshold) == 0 {
		es := sortTopX(od.Amts)
		least, err = SetThresholdLeast10(es)
		if err != nil {
			fmt.Println("updateOverall failed: ", err)
			return fmt.Errorf("updateOverall failed: %v", err)
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
		delID := reSortLeast(new, &least)
		delete(od.Amts, delID)
		od.Amts[e.ID] = e.Total
	}

	return nil
}

func updateTopCmteRecs(year string, cmte *donations.Committee) error {
	entry := &donations.Entry{ID: cmte.ID, Total: cmte.TotalReceived}

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

func updateTopCmteDonors(year string, cmte *donations.Committee) error {
	entry := &donations.Entry{ID: cmte.ID, Total: cmte.TotalTransferred}

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

func updateTopCmteExp(year string, cmte *donations.Committee) error {
	entry := &donations.Entry{ID: cmte.ID, Total: cmte.TotalDisbursed}

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

func updateTopCandidates(year string, cand *donations.Candidate) error {
	entry := &donations.Entry{ID: cand.ID, Total: cand.TotalRaised}

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

func updateTopCandExp(year string, cand *donations.Candidate) error {
	entry := &donations.Entry{ID: cand.ID, Total: cand.TotalDisbursed}

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

func updateTopIndviduals(year string, indv *donations.Individual) error {
	entry := &donations.Entry{ID: indv.ID, Total: indv.TotalDonated}

	category := "indv"
	err := updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopIndv failed: ", err)
		return fmt.Errorf("updateTopIndv failed: %v", err)
	}

	return nil
}

// get TopOverallData obj per year/category, update, & save
func updateAndSave(year, category string, entry *donations.Entry) error {
	od, err := persist.GetObject(year, "top_overall", category)
	if err != nil {
		fmt.Println("updateAndSave failed: ", err)
		return fmt.Errorf("updateAndSave failed: %v", err)
	}
	err = updateTopOverall(entry, od.(*donations.TopOverallData))
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
}

// DEPRECATED

/*
func updateTopDisbRecs(year string, rec *donations.DisbRecipient) error {
	entry := &donations.Entry{ID: rec.ID, Total: rec.TotalReceived}

	category := "disb_rec"
	err := updateAndSave(year, category, entry)
	if err != nil {
		fmt.Println("updateTopDisbRecs failed: ", err)
		return fmt.Errorf("updateTopDisbRecs failed: %v", err)
	}

	return nil
}
*/
