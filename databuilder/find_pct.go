package databuilder

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/persist"
)

// FindTotalPct finds the total percentage of the target Committee's received contributions are owned by the source object
// The source object can be either an Individual or Committee donor. Both direct & indirect contribution percentages are totaled.
func FindTotalPct(year string, source interface{}, target *donations.Committee) (float32, error) {
	seen := make(map[string]bool)

	switch s := source.(type) {
	case *donations.Individual:
		seen[s.ID] = true
		seen[target.ID] = true
		return findDonationTotalPct(year, s.RecipientsAmt, target, seen)
	case *donations.Committee:
		seen[s.ID] = true
		seen[target.ID] = true
		return findDonationTotalPct(year, s.AffiliatesAmt, target, seen)
	default:
		fmt.Println("FindTotalPct failed: wrong interface type")
		return 0.0, fmt.Errorf("FindTotalPct failed: wrong interface type")
	}
}

// FindDirectPct finds the percentage of the target Committee's funds that are directly owned by the source object.
// The source object can be either an Individual or Committee donor. Indirect contributions are not included.
func FindDirectPct(source interface{}, target *donations.Committee) (float32, error) {
	switch s := source.(type) {
	case *donations.Individual:
		return findDonationDirectPct(s.RecipientsAmt, target), nil
	case *donations.Committee:
		return findDonationDirectPct(s.AffiliatesAmt, target), nil
	default:
		fmt.Println("FindTotalPct failed: wrong interface type")
		return 0.0, fmt.Errorf("FindTotalPct failed: wrong interface type")
	}
}

// FindDonationDirectPct finds the direct ownership percentage of a given committee
func findDonationDirectPct(recs map[string]float32, target *donations.Committee) float32 {
	return recs[target.ID] / target.TotalReceived
}

// FindDonationTotalPct finds the total percentage of a specified committee owned by a donor/committee
func findDonationTotalPct(year string, recs map[string]float32, target *donations.Committee, seen map[string]bool) (float32, error) {
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
		aff, err := persist.GetObject(year, "committees", affID)
		if err != nil {
			fmt.Println("findCmteCmtePct failed: ", err)
			return 0.0, fmt.Errorf("findCmteCmtePct failed: %v", err)
		}

		// calculate indirect %
		i, err := findDonationTotalPct(year, aff.(*donations.Committee).AffiliatesAmt, target, newSeen)
		if err != nil {
			fmt.Println("findCmteCmtePct failed: ", err)
			return 0.0, fmt.Errorf("findCmteCmtePct failed: %v", err)
		}
		indir += (i * (recs[affID] / aff.(*donations.Committee).TotalReceived))
	}

	// return total
	return direct + indir, nil
}
