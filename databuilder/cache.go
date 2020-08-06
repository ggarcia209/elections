package databuilder

import (
	"fmt"
	"strings"

	"github.com/elections/donations"
	"github.com/elections/idhash"
	"github.com/elections/persist"
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
		"cmte_tx_data": make(map[string]interface{}),
		"candidates":   make(map[string]interface{}),
		"top_overall":  make(map[string]interface{}),
	}

	// add all filing committees and other objects if not already in cache
	for _, tx := range txQueue {
		// get filing committee object
		filer, err := persist.GetObject(year, "cmte_tx_data", tx.CmteID)
		if err != nil {
			fmt.Println("createCacheFromContribution failed: ", err)
			return nil, fmt.Errorf("createCacheFromContribution failed: %v", err)
		}
		cache["cmte_tx_data"][tx.CmteID] = filer

		// get linked candidate if any
		if filer.(*donations.CmteTxData).CandID != "" {
			cand, err := persist.GetObject(year, "candidates", filer.(*donations.CmteTxData).CandID)
			if err != nil {
				fmt.Println("createCacheFromContribution failed: ", err)
				return nil, fmt.Errorf("createCacheFromContribution failed: %v", err)
			}
			cache["candidates"][filer.(*donations.CmteTxData).CandID] = cand

			// nil object returned - linked candidate object does not exist
			if cand.(*donations.Candidate).ID == "" {
				cand = createCand(filer.(*donations.CmteTxData).CandID)
				cache["candidates"][filer.(*donations.CmteTxData).CandID] = cand
			}

		}

		// get other object if not in cache
		var other interface{}

		if tx.OtherID == "" {
			// not registered filer - get Individual Object
			if tx.Occupation != "" { // find Individual by name/job
				id := idhash.NewHash(idhash.FormatIndvInput(tx.Name, tx.Employer, tx.Occupation, tx.Zip))
				if cache["indviduals"][id] == nil { // not in cache
					// check database
					other, err = persist.GetObject(year, "individuals", id)
					if err != nil {
						fmt.Println("createCacheFromContribution failed: ", err)
						return nil, fmt.Errorf("createCacheFromContribution failed: %v", err)
					}
					if other.(*donations.Individual).ID == "" { // object does not exist - create new obj
						other = createIndv(id, tx)
					}
					tx.OtherID = other.(*donations.Individual).ID
					cache["individuals"][tx.OtherID] = other
				}
			} else { // find Individual by name/zip
				id := idhash.NewHash(idhash.FormatOrgInput(tx.Name, tx.Zip))
				if cache["indviduals"][id] == nil { // not in cache
					// check database
					other, err = persist.GetObject(year, "individuals", id)
					if err != nil {
						fmt.Println("createCacheFromContribution failed: ", err)
						return nil, fmt.Errorf("createCacheFromContribution failed: %v", err)
					}
					if other.(*donations.Individual).ID == "" { // object does not exist - create new obj
						other = createOrg(id, tx)
					}
				}
			}
			tx.OtherID = other.(*donations.Individual).ID
			cache["individuals"][tx.OtherID] = other
		} else {
			// registered filer - create if not in DB
			bucket := getBucket(tx.OtherID)
			if cache[bucket][tx.OtherID] == nil {
				other, err := persist.GetObject(year, bucket, tx.OtherID)
				if err != nil {
					fmt.Println("createCacheFromContribution failed: ", err)
					return nil, fmt.Errorf("createCacheFromContribution failed: %v", err)
				}
				cache[bucket][tx.OtherID] = other

				// no record in cycle's bulk data file - nil object returned
				if bucket == "cmte_tx_data" && other.(*donations.CmteTxData).CmteID == "" { // edge case - Other Registered filers not registered in current election cycle
					other = createCmte(tx.OtherID)
					cache[bucket][tx.OtherID] = other
				}
				if bucket == "candidates" && other.(*donations.Candidate).ID == "" { // edge case - Other Registered filers not registered in current election cycle
					other = createCand(tx.OtherID)
					cache[bucket][tx.OtherID] = other
				}
			}

		}
	}

	// add top overall objects
	overall, err := persist.GetTopOverall(year)
	if err != nil {
		fmt.Println("createCacheFromContribution failed: ", err)
		return nil, fmt.Errorf("createCacheFromContribution failed: %v", err)
	}
	for _, od := range overall {
		cache["top_overall"][od.(*donations.TopOverallData).Category] = od
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
		"top_overall":  make(map[string]interface{}),
	}

	// add all filing committees and other objects if not already in cache
	for _, tx := range txQueue {
		// get filing committee object
		filer, err := persist.GetObject(year, "cmte_tx_data", tx.CmteID)
		if err != nil {
			fmt.Println("createCacheFromDisbursement failed: ", err)
			return nil, fmt.Errorf("createCacheFromDisbursement failed: %v", err)
		}
		cache["cmte_tx_data"][tx.CmteID] = filer

		// get linked candidate if any
		if filer.(*donations.CmteTxData).CandID != "" {
			cand, err := persist.GetObject(year, "candidates", filer.(*donations.CmteTxData).CandID)
			if err != nil {
				fmt.Println("createCacheFromDisbursement failed: ", err)
				return nil, fmt.Errorf("createCacheFromDisbursement failed: %v", err)
			}
			cache["candidates"][filer.(*donations.CmteTxData).CandID] = cand

			// nil object returned - linked candidate object does not exist
			if cand.(*donations.Candidate).ID == "" {
				cand = createCand(filer.(*donations.CmteTxData).CandID)
				cache["candidates"][filer.(*donations.CmteTxData).CandID] = cand
			}

		}

		// get other object if not in cache
		var other interface{}

		id := idhash.NewHash(idhash.FormatOrgInput(tx.Name, tx.Zip))
		if cache["indviduals"][id] == nil { // not in cache
			// check database
			other, err = persist.GetObject(year, "individuals", id)
			if err != nil {
				fmt.Println("createCacheFromContribution failed: ", err)
				return nil, fmt.Errorf("createCacheFromContribution failed: %v", err)
			}
			if other.(*donations.Individual).ID == "" { // object does not exist - create new obj
				other = createOrg(id, tx)
			}
			tx.RecID = other.(*donations.Individual).ID
			cache["individuals"][tx.RecID] = other
		}
	}

	// add top overall objects
	overall, err := persist.GetTopOverall(year)
	if err != nil {
		fmt.Println("createCacheFromDisbursement failed: ", err)
		return nil, fmt.Errorf("createCacheFromDisbursement failed: %v", err)
	}
	for _, od := range overall {
		cache["top_overall"][od.(*donations.TopOverallData).Category] = od
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
	var bucket string

	// determine contributor/receiver type - derive from OtherID
	IDss := strings.Split(otherID, "")
	idCode := IDss[0]
	switch {
	case idCode == "C":
		bucket = "cmte_tx_data"
	case idCode == "H" || idCode == "S" || idCode == "P":
		bucket = "candidates"
	default:
		bucket = "individuals"
	}

	return bucket
}
