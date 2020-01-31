package parse

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/idhash"
	"github.com/elections/indexing"
	"github.com/elections/persist"
)

// findOrg finds an Organization from the given transaction type and creates the Organization object if it doesn't exist.
func findOrg(year string, tx interface{}, orgs map[string]map[string]interface{}) (*donations.Organization, error) {
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

func findOrgFromContribution(year string, cont *donations.Contribution, orgs map[string]map[string]interface{}) (*donations.Organization, error) {
	idEntry, err := indexing.LookupIDByZip(cont.Zip, cont.Name)
	if err != nil {
		fmt.Println("FindOrgFromIndvCont failed: ", err)
		return nil, fmt.Errorf("FindOrgFromIndvCont failed : %v", err)

	}

	// create new Indvidual obj if non-existent
	if orgs[cont.Name][cont.Zip] == nil && idEntry == nil { // nonexistent
		id := idhash.NewHash(idhash.FormatOrgInput(cont.Name, cont.Zip))
		org := &donations.Organization{
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
		o, err := persist.GetObject(year, "organizations", lookupID)
		org = o.(*donations.Organization)
		orgs[cont.Name][cont.Zip] = o.(*donations.Organization)
		if err != nil {
			fmt.Println("FindOrgFromIndvCont failed: ", err)
			return nil, fmt.Errorf("FindOrgFromIndvCont failed: %v", err)
		}
	}

	return org.(*donations.Organization), nil
}

func findOrgFromDisbursement(year string, disb *donations.Disbursement, orgs map[string]map[string]interface{}) (*donations.Organization, error) {
	idEntry, err := indexing.LookupIDByZip(disb.Zip, disb.Name)
	if err != nil {
		fmt.Println("FindOrgFromDisbursement failed: ", err)
		return nil, fmt.Errorf("FindOrgFromDisbursement failed : %v", err)

	}

	// create new Indvidual obj if non-existent
	if orgs[disb.Name][disb.Zip] == nil && idEntry == nil { // nonexistent
		id := idhash.NewHash(idhash.FormatOrgInput(disb.Name, disb.Zip))
		org := &donations.Organization{
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
		o, err := persist.GetObject(year, "organizations", lookupID)
		org = o.(*donations.Organization)
		orgs[disb.Name][disb.Zip] = o.(*donations.Organization)
		if err != nil {
			fmt.Println("FindOrgFromDisbursement failed: ", err)
			return nil, fmt.Errorf("FindOrgFromDisbursement failed: %v", err)
		}
	}

	return org.(*donations.Organization), nil
}

// DEPRECATED
/*
// findPerson returns a pointer to an Indvidual object. The object is created if it doesn't exist.
func findOrg(year string, disb *donations.Disbursement, orgs map[string]map[string]*donations.Organization) (*donations.Organization, error) {
	name := disb.Name
	zip := disb.Zip
	zipMap, err := persist.LookupRecIDByName(name)
	if err != nil {
		fmt.Println("FindOrg failed: LookupRecIDByName failed:", err)
		return nil, fmt.Errorf("FindOrg failed: LookupRecIDByName failed: %v", err)

	}

	// create new Indvidual obj if non-existent
	if orgs[name][zip] == nil && (zipMap == nil || zipMap[zip] == "") { // nonexistent
		toHash := name + " - " + zip
		id := idhash.NewHash(idhash.FormatOrgInput(disb.Name, disb.Zip))
		org := &donations.Organization{
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

		if orgs[name] == nil {
			orgs[name] = make(map[string]*donations.Organization)
		}

		return org, nil
	}

	if orgs[name] == nil {
		orgs[name] = make(map[string]*donations.Organization)
	}

	// find existing donor
	org := orgs[name][zip]
	if org == nil { // if not in memory, retrieve from disk
		lookupID := zipMap[zip]
		o, err := persist.GetObject(year, "organizations", lookupID)
		org = o.(*donations.Organization)
		orgs[name][zip] = o.(*donations.Organization)
		if err != nil {
			fmt.Println("FindRecipient failed: ", err)
			return nil, fmt.Errorf("FindRecipient failed: %v", err)
		}
	}

	return org, nil
*/
