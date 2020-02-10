package databuilder

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/persist"
)

func updateTopOverall(year string, filer *donations.CmteTxData, other interface{}, incoming, transfer bool) error {
	switch t := other.(type) {
	case *donations.Individual:
		if !incoming {
			// update Top Committes by expenses only - individuals not ranked by funds received
			err := updateTopCmteExp(year, filer)
			if err != nil {
				fmt.Println("updateTopOverall failed: ", err)
				return fmt.Errorf("updateTopOverall failed: %v", err)
			}
			return nil
		}

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
		return fmt.Errorf("updateTopOverall failed: wrong interface type")
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

	category := "org"
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

	category := "org"
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
}

// compare obj to top overall threshold
func compareTopOverall(e *donations.Entry, od *donations.TopOverallData) error {
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
		least, err = setThresholdLeast10(es)
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
