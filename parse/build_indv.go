// update - Removed DonorID global variable and added as argument to FindPerson function
// update - removed in memory map and implement lookup from disc
// update - removed update logic; FindPerson now returns a donor obj in both cases; updating done outside FindPerson
// update - moved to parse package
// 8/3/20 - Removed Organization cases and references
package parse

import (
	"fmt"
	"time"

	"github.com/elections/donations"
	"github.com/elections/idhash"
	"github.com/elections/indexing"
	"github.com/elections/persist"
)

// findPerson returns a pointer to an Indvidual object. The object is created if it doesn't exist.
func findPerson(year string, cont *donations.Contribution, donors map[string]map[string]interface{}) (*donations.Individual, error) {
	start := time.Now()

	job := cont.Employer + " - " + cont.Occupation
	idEntry, err := indexing.LookupIDByJob(cont.Name, cont.Employer, cont.Occupation)
	if err != nil {
		fmt.Println("findPerson failed:  failed:", err)
		return nil, fmt.Errorf("findPerson failed: failed: %v", err)
	}

	fin := time.Since(start)
	fmt.Println("      id lookup time: ", fin)

	// create new Indvidual obj if non-existent
	if donors[cont.Name][job] == nil && idEntry.ID == "" { // nonexistent
		start = time.Now()
		id := idhash.NewHash(idhash.FormatIndvInput(cont.Name, cont.Employer, cont.Occupation, cont.Zip))
		fin = time.Since(start)
		fmt.Println("      hash id gen time: ", fin)

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

		fin = time.Since(start)
		fmt.Println("      invidual obj gen time: ", fin)

		return donor, nil
	}

	// find existing donor
	if donors[cont.Name] == nil {
		donors[cont.Name] = make(map[string]interface{})
	}

	donor := donors[cont.Name][job]
	if donor == nil { // if not in memory, retrieve from disk
		start = time.Now()
		lookupID := idEntry.ID
		d, err := persist.GetObject(year, "individuals", lookupID)
		donor = d.(*donations.Individual)
		if err != nil {
			fmt.Println("FindPerson failed: ", err)
			return nil, fmt.Errorf("FindPerson failed: %v", err)
		}
		fin = time.Since(start)
		fmt.Println("      GetObject time: ", fin)
	}

	return donor.(*donations.Individual), nil
}

// findOrg finds an Organization from the given transaction type and creates the Organization object if it doesn't exist.
func findOrg(year string, tx interface{}, orgs map[string]map[string]interface{}) (*donations.Individual, error) {
	// 	fmt.Println("* find org *")
	switch t := tx.(type) {
	case *donations.Contribution:
		return findOrgFromContribution(year, t, orgs)
	case *donations.Disbursement:
		return findOrgFromDisbursement(year, t, orgs)
	default:
		fmt.Println("FindOrg failed: Invalid interface type")
		return nil, fmt.Errorf("FindOrg failed: Invalid interface type")
	}
}

func findOrgFromContribution(year string, cont *donations.Contribution, orgs map[string]map[string]interface{}) (*donations.Individual, error) {
	idEntry, err := indexing.LookupIDByZip(cont.Zip, cont.Name)
	if err != nil {
		fmt.Println("FindOrgFromIndvCont failed: ", err)
		return nil, fmt.Errorf("FindOrgFromIndvCont failed : %v", err)

	}

	// create new Indvidual obj if non-existent
	if orgs[cont.Name][cont.Zip] == nil && idEntry.ID == "" { // nonexistent
		id := idhash.NewHash(idhash.FormatOrgInput(cont.Name, cont.Zip))
		org := &donations.Individual{
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
			AvgTxIn:       0.0,
			RecipientsAmt: make(map[string]float32),
			RecipientsTxs: make(map[string]float32),
			SendersAmt:    make(map[string]float32),
			SendersTxs:    make(map[string]float32),
		}

		if orgs[cont.Name] == nil {
			orgs[cont.Name] = make(map[string]interface{})
		}
		orgs[cont.Name][cont.Zip] = org

		return org, nil
	}

	if orgs[cont.Name] == nil {
		orgs[cont.Name] = make(map[string]interface{})
	}

	// find existing donor
	org := orgs[cont.Name][cont.Zip]
	if org == nil { // if not in memory, retrieve from disk
		lookupID := idEntry.ID
		o, err := persist.GetObject(year, "individuals", lookupID)
		org = o.(*donations.Individual)
		orgs[cont.Name][cont.Zip] = o.(*donations.Individual)
		if err != nil {
			fmt.Println("findOrgFromContribution failed: ", err)
			return nil, fmt.Errorf("findOrgFromContribution failed: %v", err)
		}
	}

	/* fmt.Println("name: ", cont.Name)
	fmt.Println("zip: ", cont.Zip)
	fmt.Println("ID Entry: ", idEntry.ID)
	fmt.Println("OrgID: ", org.(*donations.Individual).ID)
	fmt.Println("TxID: ", cont.TxID) */

	return org.(*donations.Individual), nil
}

func findOrgFromDisbursement(year string, disb *donations.Disbursement, orgs map[string]map[string]interface{}) (*donations.Individual, error) {
	idEntry, err := indexing.LookupIDByZip(disb.Zip, disb.Name)
	if err != nil {
		fmt.Println("findOrgFromDisbursement failed: ", err)
		return nil, fmt.Errorf("findOrgFromDisbursement failed : %v", err)

	}

	// create new Indvidual obj if non-existent
	if orgs[disb.Name][disb.Zip] == nil && idEntry.ID == "" { // nonexistent
		id := idhash.NewHash(idhash.FormatOrgInput(disb.Name, disb.Zip))
		org := &donations.Individual{
			ID:            id,
			Name:          disb.Name,
			City:          disb.City,
			State:         disb.State,
			Zip:           disb.Zip,
			Occupation:    "",
			Employer:      "",
			Transactions:  []string{},
			TotalOutAmt:   0.0,
			TotalOutTxs:   0.0,
			AvgTxOut:      0.0,
			TotalInAmt:    0.0,
			TotalInTxs:    0.0,
			AvgTxIn:       0.0,
			RecipientsAmt: make(map[string]float32),
			RecipientsTxs: make(map[string]float32),
			SendersAmt:    make(map[string]float32),
			SendersTxs:    make(map[string]float32),
		}

		if orgs[disb.Name] == nil {
			orgs[disb.Name] = make(map[string]interface{})
		}
		orgs[disb.Name][disb.Zip] = org

		return org, nil
	}

	if orgs[disb.Name] == nil {
		orgs[disb.Name] = make(map[string]interface{})
	}

	// find existing donor
	org := orgs[disb.Name][disb.Zip]
	if org == nil { // if not in memory, retrieve from disk
		lookupID := idEntry.ID
		o, err := persist.GetObject(year, "individuals", lookupID)
		org = o.(*donations.Individual)
		orgs[disb.Name][disb.Zip] = o.(*donations.Individual)
		if err != nil {
			fmt.Println("findOrgFromDisbursement failed: ", err)
			return nil, fmt.Errorf("findOrgFromDisbursement failed: %v", err)
		}
	}

	return org.(*donations.Individual), nil
}
