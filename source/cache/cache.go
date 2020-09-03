package cache

import (
	"fmt"
	"strings"

	"github.com/elections/source/donations"
	"github.com/elections/source/idhash"
	"github.com/elections/source/persist"
)

// CreateCache creates a temporary in-memory cache of filer and other objects derived from a list of Contributions or Disbursements
func CreateCache(year string, txQueue interface{}) (map[string]map[string]interface{}, error) {
	var cache map[string]map[string]interface{}
	var err error

	_, cont := txQueue.([]*donations.Contribution)

	if cont {
		cache, err = createCacheFromContribution(year, txQueue.([]*donations.Contribution))
		if err != nil {
			fmt.Println("CreateCache failed: ", err)
			return nil, fmt.Errorf("CreateCache failed: %v", err)
		}
	} else {
		cache, err = createCacheFromDisbursement(year, txQueue.([]*donations.Disbursement))
		if err != nil {
			fmt.Println("CreateCache failed: ", err)
			return nil, fmt.Errorf("CreateCache failed: %v", err)
		}
	}

	return cache, nil
}

func createCacheFromContribution(year string, txQueue []*donations.Contribution) (map[string]map[string]interface{}, error) {
	if len(txQueue) == 0 {
		return nil, nil
	}

	cache := map[string]map[string]interface{}{
		"individuals":  make(map[string]interface{}),
		"committees":   make(map[string]interface{}),
		"cmte_tx_data": make(map[string]interface{}),
		"candidates":   make(map[string]interface{}),
	}
	filerIDs := []string{}
	otherIDs := map[string][]string{
		"individuals":  []string{},
		"cmte_tx_data": []string{},
		"candidates":   []string{},
	}
	seen := make(map[string]bool)

	for _, tx := range txQueue {
		// edge case - earmarked transaction type, treat as incoming transaction type to OtherID cmte
		if tx.TxType == "24I" {
			if tx.OtherID != "" {
				tx.CmteID = tx.OtherID
				tx.TxType = "15E"
			} else { // edge case
				fmt.Println("WARNING: NIL RECIPIENT - txID: ", tx.TxID)
				// transaction will be treated as incoming/memo transaction
				// from individual to intermediary (filing cmte)
				tx.TxType = "15E"
			}
		}

		// get filer IDs for obj lookup
		if !seen[tx.CmteID] {
			filerIDs = append(filerIDs, tx.CmteID)
			seen[tx.CmteID] = true
		}

		em := earmark(tx.TxType)

		// get otherIDs
		// initialize placeholder objects for Individuals
		if tx.OtherID == "" || em { // not registered filer or earmarked transaction - get Individual Object
			if tx.Occupation != "" { // find Individual by name/job
				id := idhash.NewHash(idhash.FormatIndvInput(tx.Name, tx.Employer, tx.Occupation, tx.Zip))
				if cache["individuals"][id] == nil { // not in cache
					other := createIndv(id, tx)
					cache["individuals"][id] = other
				}
				tx.OtherID = id
			} else { // find Individual by name/zip
				id := idhash.NewHash(idhash.FormatOrgInput(tx.Name, tx.Zip))
				if cache["individuals"][id] == nil { // not in cache
					other := createOrg(id, tx)
					cache["individuals"][id] = other
				}
				tx.OtherID = id
			}
		}
		if !seen[tx.OtherID] {
			bkt := getBucket(tx.OtherID)
			otherIDs[bkt] = append(otherIDs[bkt], tx.OtherID)
			seen[tx.OtherID] = true
		}
	}

	// batch get filer objs, add to cache
	bucket := "cmte_tx_data"
	filers, nilIDs, err := persist.BatchGetByID(year, bucket, filerIDs)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("createCacheFromContribution failed: %v", err)
	}
	for _, f := range filers {
		ID := f.(*donations.CmteTxData).CmteID
		cache[bucket][ID] = f
	}
	for _, nID := range nilIDs {
		if cache[bucket][nID] == nil {
			unk, txData := createCmte(nID)
			cache["committees"][nID] = unk
			cache[bucket][nID] = txData
		}
	}

	// batch get other objs, add to cache
	for bkt, objIDs := range otherIDs {
		others, nilIDs, err := persist.BatchGetByID(year, bkt, objIDs)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("createCacheFromContribution failed: %v", err)
		}
		for _, o := range others {
			ID := getObjID(o)
			cache[bkt][ID] = o // overwrites placeholder Individuals if any
		}

		if bkt == "individuals" {
			continue
		}

		// create placeholder objects for registered filers not listed in current election cycle's data
		for _, nID := range nilIDs {
			if cache[bkt][nID] == nil {
				if bkt == "cmte_tx_data" { // edge case - committee not registered in current election cycle
					var txData *donations.CmteTxData
					var unk *donations.Committee
					unk, txData = createCmte(nID)
					cache["committees"][nID] = unk
					cache[bkt][nID] = txData
				}
				if bkt == "candidates" { // edge case - candidate not registered in current election cycle
					var unk *donations.Candidate
					unk = createCand(nID)
					cache[bkt][nID] = unk
				}
			}
		}
	}

	return cache, nil
}

func createCacheFromDisbursement(year string, txQueue []*donations.Disbursement) (map[string]map[string]interface{}, error) {
	if len(txQueue) == 0 {
		return nil, nil
	}

	cache := map[string]map[string]interface{}{
		"individuals":  make(map[string]interface{}),
		"cmte_tx_data": make(map[string]interface{}),
		"candidates":   make(map[string]interface{}),
		"committees":   make(map[string]interface{}),
	}
	seen := make(map[string]bool)

	filerIDs := []string{}
	otherIDs := []string{}

	// add all filing committees and other objects if not already in cache
	for _, tx := range txQueue {
		// get filer IDs for obj lookup
		if !seen[tx.CmteID] {
			filerIDs = append(filerIDs, tx.CmteID)
			seen[tx.CmteID] = true
		}

		// initialize placeholder objects for Individuals
		id := idhash.NewHash(idhash.FormatOrgInput(tx.Name, tx.Zip))
		if cache["individuals"][id] == nil { // not in cache
			other := createOrg(id, tx)
			tx.RecID = other.ID
			cache["individuals"][tx.RecID] = other
		} else {
			tx.RecID = id
		}

		// get otherIDs
		if !seen[tx.RecID] {
			otherIDs = append(otherIDs, tx.RecID)
			seen[tx.RecID] = true
		}
	}

	// get filing cmte objects & add to cache
	filers, nilIDs, err := persist.BatchGetByID(year, "cmte_tx_data", filerIDs)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("createCacheFromDisbursement failed: %v", err)
	}
	for _, f := range filers {
		ID := f.(*donations.CmteTxData).CmteID
		cache["cmte_tx_data"][ID] = f
	}
	for _, nID := range nilIDs {
		if cache["cmte_tx_data"][nID] == nil {
			unk, txData := createCmte(nID)
			cache["committees"][nID] = unk
			cache["cmte_tx_data"][nID] = txData
		}
	}

	// get other objects and add to cache
	others, _, err := persist.BatchGetByID(year, "individuals", otherIDs)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("createCacheFromDisbursement failed: %v", err)
	}
	for _, o := range others {
		ID := o.(*donations.Individual).ID
		cache["cmte_tx_data"][ID] = o
	}

	return cache, nil
}

// SerializeCache converts a cache from map[string]map[string]interface{} to []interface{}
func SerializeCache(cache map[string]map[string]interface{}) []interface{} {
	objs := []interface{}{}

	for _, cat := range cache {
		for _, data := range cat {
			objs = append(objs, data)
		}
	}

	return objs
}

func getBucket(otherID string) string {
	if otherID == "" {
		return "individuals"
	}

	var bucket string

	// determine contributor/receiver type - derive from OtherID
	IDss := strings.Split(otherID, "")
	idCode := IDss[0]
	switch {
	case idCode == "C":
		if len(otherID) < 16 {
			bucket = "cmte_tx_data"
		} else {
			bucket = "individuals" // edge case - hashes beginning with C
		}
	case idCode == "H" || idCode == "S" || idCode == "P":
		if len(otherID) < 16 {
			bucket = "candidates"
		} else {
			bucket = "individuals" // edge case - hashes beginning with H, S, P
		}
	default:
		bucket = "individuals"
	}

	return bucket
}

func getObjID(obj interface{}) string {
	switch t := obj.(type) {
	case *donations.Individual:
		return t.ID
	case *donations.CmteTxData:
		return t.CmteID
	case *donations.Candidate:
		return t.ID
	default:
		return ""
	}
}

func earmark(code string) bool {
	check := map[string]bool{
		"15E": true,
		"15I": true,
		"15T": true,
		"24I": true,
	}
	return check[code]
}
