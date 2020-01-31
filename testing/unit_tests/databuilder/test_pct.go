package main

import (
	"fmt"
	"os"

	"github.com/elections/donations"
)

var pcc00 = donations.Committee{
	ID:            "pcc00",
	TotalReceived: 2000,
}

var pac00 = donations.Committee{
	ID:            "pac00",
	TotalReceived: 1000,
	AffiliatesAmt: map[string]float32{
		"pcc00": 500,
		"pac01": 100,
	},
}

var pac01 = donations.Committee{
	ID:            "pac01",
	TotalReceived: 2000,
	AffiliatesAmt: map[string]float32{
		"pcc00": 400,
		"pac00": 250,
	},
}

var d0 = donations.Committee{
	ID: "d0",
	AffiliatesAmt: map[string]float32{
		"pcc00": 500,
		"pac00": 200,
		"pac01": 200,
	},
}

var memCache = make(map[string]*donations.Committee)

func main() {

	memCache[d0.ID] = &d0
	memCache[pcc00.ID] = &pcc00
	memCache[pac00.ID] = &pac00
	memCache[pac01.ID] = &pac01

	pct, err := FindTotalPct(&d0, &pcc00)
	if err != nil {
		fmt.Println("failed: ", err)
		os.Exit(1)
	}

	fmt.Println(pct)
}

func FindTotalPct(source interface{}, target *donations.Committee) (float32, error) {
	seen := make(map[string]bool)

	switch s := source.(type) {
	case *donations.Individual:
		seen[s.ID] = true
		seen[target.ID] = true
		return findDonationPct(s.RecipientsAmt, target, seen)
	case *donations.Committee:
		seen[s.ID] = true
		seen[target.ID] = true
		return findDonationPct(s.AffiliatesAmt, target, seen)
	default:
		fmt.Println("FindTotalPct failed: wrong interface type")
		return 0.0, fmt.Errorf("FindTotalPct failed: wrong interface type")
	}
}

func findDonationPct(recs map[string]float32, target *donations.Committee, seen map[string]bool) (float32, error) {
	// find direct contribution %
	direct := recs[target.ID] / target.TotalReceived

	// find indirect %
	indir := float32(0.0)
	for affID := range recs {
		// check if seen; copy and update seen map
		if seen[affID] {
			continue
		}
		newSeen := make(map[string]bool)
		for k := range seen {
			newSeen[k] = true
		}
		newSeen[affID] = true

		// get affilate cmte obj
		aff := memCache[affID]
		/* aff, err := persist.GetCommittee(affID)
		if err != nil {
			fmt.Println("findCmteCmtePct failed: ", err)
			return 0.0, fmt.Errorf("findCmteCmtePct failed: %v", err)
		} */

		// calculate indirect %
		i, err := findDonationPct(aff.AffiliatesAmt, target, newSeen)
		if err != nil {
			fmt.Println("findCmteCmtePct failed: ", err)
			return 0.0, fmt.Errorf("findCmteCmtePct failed: %v", err)
		}
		indir += (i * (float32(recs[affID]) / float32(aff.TotalReceived)))
	}

	// return total
	return direct + indir, nil
}
