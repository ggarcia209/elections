// update - Removed DonorID global variable and added as argument to FindPerson function
// update - removed in memory map and implement lookup from disc
// update - removed update logic; FindPerson now returns a donor obj in both cases; updating done outside FindPerson
// update - moved to parse package
package parse

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/idhash"
	"github.com/elections/indexing"
	"github.com/elections/persist"
)

// findPerson returns a pointer to an Indvidual object. The object is created if it doesn't exist.
func findPerson(year string, cont *donations.Contribution, donors map[string]map[string]interface{}) (*donations.Individual, error) {
	job := cont.Employer + " - " + cont.Occupation
	idEntry, err := indexing.LookupIDByJob(cont.Name, cont.Employer, cont.Occupation)
	if err != nil {
		fmt.Println("FindPerson failed: LookupIDByName failed:", err)
		return nil, fmt.Errorf("FindPerson failed: LookupIDByName failed: %v", err)

	}

	// create new Indvidual obj if non-existent
	if donors[cont.Name][job] == nil && idEntry == nil { // nonexistent
		id := idhash.NewHash(idhash.FormatIndvInput(cont.Name, cont.Employer, cont.Occupation, cont.Zip))
		donor := &donations.Individual{
			ID:            id,
			Name:          cont.Name,
			City:          cont.City,
			State:         cont.State,
			Zip:           cont.Zip,
			Occupation:    cont.Occupation,
			Employer:      cont.Employer,
			Transactions:  []string{},
			TotalOutAmt:   0.0,
			TotalOutTxs:   0.0,
			AvgTxOut:      0.0,
			TotalInAmt:    0.0,
			TotalInTxs:    0.0,
			NetBalance:    0.0,
			RecipientsAmt: make(map[string]float32),
			RecipientsTxs: make(map[string]float32),
			SendersAmt:    make(map[string]float32),
			SendersTxs:    make(map[string]float32),
		}

		if donors[cont.Name] == nil {
			donors[cont.Name] = make(map[string]interface{})
		}

		return donor, nil
	}

	// find existing donor
	if donors[cont.Name] == nil {
		donors[cont.Name] = make(map[string]interface{})
	}

	donor := donors[cont.Name][job]
	if donor == nil { // if not in memory, retrieve from disk
		lookupID := idEntry.ID
		d, err := persist.GetObject(year, "individuals", lookupID)
		donor = d.(*donations.Individual)
		if err != nil {
			fmt.Println("FindPerson failed: ", err)
			return nil, fmt.Errorf("FindPerson failed: %v", err)
		}
	}

	return donor.(*donations.Individual), nil
}

// DEPRECATED
/*
// findPerson returns a pointer to an Indvidual object. The object is created if it doesn't exist.
func findPerson(year string, cont *donations.IndvContribution, donors map[string]map[string]*donations.Individual) (*donations.Individual, error) {
	name := cont.Name
	job := cont.Employer + " - " + cont.Occupation
	jobMap, err := persist.LookupIDByName(name)
	if err != nil {
		fmt.Println("FindPerson failed: LookupIDByName failed:", err)
		return nil, fmt.Errorf("FindPerson failed: LookupIDByName failed: %v", err)

	}

	// create new Indvidual obj if non-existent
	if donors[name][job] == nil && (jobMap == nil || jobMap[job] == "") { // nonexistent
		id := idhash.NewHash(idhash.FormatIndvInput(cont.Name, cont.Employer, cont.Occupation, cont.Zip))
		donor := &donations.Individual{
			ID:            id,
			Name:          cont.Name,
			City:          cont.City,
			State:         cont.State,
			Zip:           cont.Zip,
			Occupation:    cont.Occupation,
			Employer:      cont.Employer,
			Transactions:  []string{},
			TotalOutAmt:   0.0,
			TotalOutTxs:   0.0,
			AvgTxOut:      0.0,
			TotalInAmt:    0.0,
			TotalInTxs:    0.0,
			NetBalance:    0.0,
			RecipientsAmt: make(map[string]float32),
			RecipientsTxs: make(map[string]float32),
			SendersAmt:    make(map[string]float32),
			SendersTxs:    make(map[string]float32),
		}

		if donors[name] == nil {
			donors[name] = make(map[string]*donations.Individual)
		}

		return donor, nil
	}

	// find existing donor
	if donors[name] == nil {
		donors[name] = make(map[string]*donations.Individual)
	}

	donor := donors[name][job]
	if donor == nil { // if not in memory, retrieve from disk
		lookupID := jobMap[job]
		d, err := persist.GetObject(year, "individuals", lookupID)
		donor = d.(*donations.Individual)
		donors[name][job] = d.(*donations.Individual)
		if err != nil {
			fmt.Println("FindPerson failed: ", err)
			return nil, fmt.Errorf("FindPerson failed: %v", err)
		}
	}

	return donor, nil
}
*/
