// Package donations contains the base objects that are used throughout the application.
// Objects within this package are primarily used for creating, updating, and persisting
// the datasets derived from the input data.
// Datasets for objects in this file are derived from the procesed objects
// in the other files in this package, once all raw transaction data has been processed.
// This file also contains & exposes functions for instantiating
// these bjects outside of the package.
package donations

import "fmt"

// TopOverallData stores a map of the top x number of objects
// for a specific category (Top 500 by default).
// Each instance corresponds to a specific Year/Bucket/Cateogry/Party.
// Ex: ("all_time/cmte_tx_data/rec/ALL")
type TopOverallData struct {
	ID        string // hash(yr+bucket+cat+pty)
	Year      string // "2018"
	Bucket    string // "cmte_tx_data"
	Category  string // "rec"
	Party     string // "ALL"
	Amts      map[string]float32
	Threshold []*Entry
	SizeLimit int
}

// YearlyTotal contains the total sum of funds donated/transferred/spent
// for a given year and party (including Year: all_time & Party: ALL).
type YearlyTotal struct {
	ID       string  // hash(year+cat+pty)
	Year     string  // "2018"
	Category string  // "exp"
	Party    string  // "REP"
	Total    float32 // total sum
}

// Entry represents a key/value pair from a Top X map and is used to sort and update the map.
type Entry struct {
	ID    string
	Total float32
}

// InitSecondaryDataObjs initializes a set of TopOverall and YearlyTotal objects for the given year.
func InitSecondaryDataObjs(year string) ([]interface{}, []interface{}) {
	ods := initTopOverallDataObjs(year)
	yts := initYearlyTotalObjs(year)
	return ods, yts
}

// initTopOverallDataObjs creates a TopOverallData object for each
// category and returns the objects in a list.
// Formats object ID's as year-bucket-category-party.
func initTopOverallDataObjs(year string) []interface{} {
	limit := 500
	od := []interface{}{}
	// indv & disb_rec
	indv := initTopOverallObj(year, "individuals", "donor", "ALL", 100000)
	od = append(od, indv)
	indvRec := initTopOverallObj(year, "individuals", "rec", "ALL", 100000)
	od = append(od, indvRec)

	// cmte_recs
	cmteRecAll := initTopOverallObj(year, "cmte_tx_data", "rec", "ALL", limit)
	od = append(od, cmteRecAll)
	cmteRecR := initTopOverallObj(year, "cmte_tx_data", "rec", "REP", limit)
	od = append(od, cmteRecR)
	cmteRecD := initTopOverallObj(year, "cmte_tx_data", "rec", "DEM", limit)
	od = append(od, cmteRecD)
	cmteRecNa := initTopOverallObj(year, "cmte_tx_data", "rec", "IND", limit)
	od = append(od, cmteRecNa)
	cmteRecOth := initTopOverallObj(year, "cmte_tx_data", "rec", "OTH", limit)
	od = append(od, cmteRecOth)
	cmteRecUnk := initTopOverallObj(year, "cmte_tx_data", "rec", "UNK", limit)
	od = append(od, cmteRecUnk)

	// cmte_donors
	cmteAll := initTopOverallObj(year, "cmte_tx_data", "donor", "ALL", limit)
	od = append(od, cmteAll)
	cmteR := initTopOverallObj(year, "cmte_tx_data", "donor", "REP", limit)
	od = append(od, cmteR)
	cmteD := initTopOverallObj(year, "cmte_tx_data", "donor", "DEM", limit)
	od = append(od, cmteD)
	cmteNa := initTopOverallObj(year, "cmte_tx_data", "donor", "IND", limit)
	od = append(od, cmteNa)
	cmteOth := initTopOverallObj(year, "cmte_tx_data", "donor", "OTH", limit)
	od = append(od, cmteOth)
	cmteUnk := initTopOverallObj(year, "cmte_tx_data", "donor", "UNK", limit)
	od = append(od, cmteUnk)

	// cmte_exp
	cmteExpAll := initTopOverallObj(year, "cmte_tx_data", "exp", "ALL", limit)
	od = append(od, cmteExpAll)
	cmteExpR := initTopOverallObj(year, "cmte_tx_data", "exp", "REP", limit)
	od = append(od, cmteExpR)
	cmteExpD := initTopOverallObj(year, "cmte_tx_data", "exp", "DEM", limit)
	od = append(od, cmteExpD)
	cmteExpNa := initTopOverallObj(year, "cmte_tx_data", "exp", "IND", limit)
	od = append(od, cmteExpNa)
	cmteExpOth := initTopOverallObj(year, "cmte_tx_data", "exp", "OTH", limit)
	od = append(od, cmteExpOth)
	cmteExpUnk := initTopOverallObj(year, "cmte_tx_data", "exp", "UNK", limit)
	od = append(od, cmteExpUnk)

	// cand
	candAll := initTopOverallObj(year, "candidates", "rec", "ALL", limit)
	od = append(od, candAll)
	candR := initTopOverallObj(year, "candidates", "rec", "REP", limit)
	od = append(od, candR)
	candD := initTopOverallObj(year, "candidates", "rec", "DEM", limit)
	od = append(od, candD)
	candNa := initTopOverallObj(year, "candidates", "rec", "IND", limit)
	od = append(od, candNa)
	candOth := initTopOverallObj(year, "candidates", "rec", "OTH", limit)
	od = append(od, candOth)
	candUnk := initTopOverallObj(year, "candidates", "rec", "UNK", limit)
	od = append(od, candUnk)

	// cand_donor
	candDonorAll := initTopOverallObj(year, "candidates", "donor", "ALL", limit)
	od = append(od, candDonorAll)
	candDonorR := initTopOverallObj(year, "candidates", "donor", "REP", limit)
	od = append(od, candDonorR)
	candDonorD := initTopOverallObj(year, "candidates", "donor", "DEM", limit)
	od = append(od, candDonorD)
	candDonorNa := initTopOverallObj(year, "candidates", "donor", "IND", limit)
	od = append(od, candDonorNa)
	candDonorOth := initTopOverallObj(year, "candidates", "donor", "OTH", limit)
	od = append(od, candDonorOth)
	candDonorUnk := initTopOverallObj(year, "candidates", "donor", "UNK", limit)
	od = append(od, candDonorUnk)

	// cand_exp
	candExpAll := initTopOverallObj(year, "candidates", "exp", "ALL", limit)
	od = append(od, candExpAll)
	candExpR := initTopOverallObj(year, "candidates", "exp", "REP", limit)
	od = append(od, candExpR)
	candExpD := initTopOverallObj(year, "candidates", "exp", "DEM", limit)
	od = append(od, candExpD)
	candExpNa := initTopOverallObj(year, "candidates", "exp", "IND", limit)
	od = append(od, candExpNa)
	candExpOth := initTopOverallObj(year, "candidates", "exp", "OTH", limit)
	od = append(od, candExpOth)
	candExpUnk := initTopOverallObj(year, "candidates", "exp", "UNK", limit)
	od = append(od, candExpUnk)

	return od
}

// initYearlyTotalObjs initializes the YearlyTotal objects for the given year.
func initYearlyTotalObjs(year string) []interface{} {
	yts := []interface{}{}

	recAll := initYrTotalObj(year, "rec", "ALL")
	yts = append(yts, recAll)
	recR := initYrTotalObj(year, "rec", "REP")
	yts = append(yts, recR)
	recD := initYrTotalObj(year, "rec", "DEM")
	yts = append(yts, recD)
	recNa := initYrTotalObj(year, "rec", "IND")
	yts = append(yts, recNa)
	recOth := initYrTotalObj(year, "rec", "OTH")
	yts = append(yts, recOth)
	recUnk := initYrTotalObj(year, "rec", "UNK")
	yts = append(yts, recUnk)

	donorAll := initYrTotalObj(year, "donor", "ALL")
	yts = append(yts, donorAll)
	donorR := initYrTotalObj(year, "donor", "REP")
	yts = append(yts, donorR)
	donorD := initYrTotalObj(year, "donor", "DEM")
	yts = append(yts, donorD)
	donorNa := initYrTotalObj(year, "donor", "IND")
	yts = append(yts, donorNa)
	donorOth := initYrTotalObj(year, "donor", "OTH")
	yts = append(yts, donorOth)
	donorUnk := initYrTotalObj(year, "donor", "UNK")
	yts = append(yts, donorUnk)

	expAll := initYrTotalObj(year, "exp", "ALL")
	yts = append(yts, expAll)
	expR := initYrTotalObj(year, "exp", "REP")
	yts = append(yts, expR)
	expD := initYrTotalObj(year, "exp", "DEM")
	yts = append(yts, expD)
	expNa := initYrTotalObj(year, "exp", "IND")
	yts = append(yts, expNa)
	expOth := initYrTotalObj(year, "exp", "OTH")
	yts = append(yts, expOth)
	expUnk := initYrTotalObj(year, "exp", "UNK")
	yts = append(yts, expUnk)

	return yts
}

func initTopOverallObj(year, bucket, cat, pty string, limit int) *TopOverallData {
	id := year + "-" + bucket + "-" + cat + "-" + pty
	fmt.Println("created Top Overall: ", id)
	od := &TopOverallData{
		ID:        id,
		Year:      year,
		Bucket:    bucket,
		Category:  cat,
		Party:     pty,
		Amts:      make(map[string]float32),
		Threshold: nil,
		SizeLimit: limit,
	}
	return od
}

func initYrTotalObj(year, cat, pty string) *YearlyTotal {
	id := year + "-" + cat + "-" + pty
	fmt.Println("created Yearly Total: ", id)
	yt := &YearlyTotal{
		ID:       id,
		Year:     year,
		Category: cat,
		Party:    pty,
		Total:    0,
	}
	return yt
}
