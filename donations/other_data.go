package donations

// TopOverallData stores an entry for the Top 100 Overall data types.
// Each instance / enty represents a different category.
// (Ex: Top Candidates by Funds Raised, Top Committees by Funds Transferred
type TopOverallData struct {
	Category  string
	Amts      map[string]float32
	Threshold []*Entry
	SizeLimit int
}

// TopOverallData Categories:

// Individual Donors - "indv"
// Disbursement Recipients "disb_rec"

// Committee Donors - All - "cmte_donors_all"
// Committee Donors - Republican - "cmte_donors_r"
// Committee Donors - Democrat - "cmte_donors_d"
// Committee Donors - Independent/Non-Affiliated - "cmte_donors_na"
// Committee Donors - All other parties - "cmte_donors_misc"

// Committee Recipients - All - "cmte_recs_all"
// Committee Recipients - Republican - "cmte_recs_r"
// Committee Recipients - Democrat - "cmte_recs_d"
// Committee Recipients - Independent/Non-Affiliated - "cmte_recs_na"
// Committee Recipients - All other parties - "cmte_recs_misc"

// Committee Spenders - All - "cmte_exp_all"
// Committee Spenders - Republican - "cmte_exp_r"
// Committee Spenders - Democrat - "cmte_exp_d"
// Committee Spenders - Independent/Non-Affiliated - "cmte_exp_na"
// Committee Spenders - All other parties - "cmte_exp_misc"

// Candidate Recipients - All - "cand_all"
// Candidate Recipients - Republican - "cand_r"
// Candidate Recipients - Democrat - "cand_d"
// Candidate Recipients - Independent/Non-Affiliated - "cand_na"
// Candidate Recipients - All other parties - "cand_misc"

// Candidate Spenders - All - "cand_exp_all"
// Candidate Spenders - Republican - "cand_exp_r"
// Candidate Spenders - Democrat - "cand_exp_d"
// Candidate Spenders - Independent/Non-Affiliated - "cand_exp_na"
// Candidate Spenders - All other parties - "cand_exp_misc"

// InitTopOverallDataObjs creates a TopOverallData object for each
// category and returns the objects in a list
func InitTopOverallDataObjs(limit int) []interface{} {
	od := []interface{}{}
	// indv & disb_rec
	indv := &TopOverallData{"indv", make(map[string]float32), nil, limit}
	od = append(od, indv)
	indvRec := &TopOverallData{"indv_rec", make(map[string]float32), nil, limit}
	od = append(od, indvRec)

	// cmte_donors
	cmteAll := &TopOverallData{"cmte_donors_all", make(map[string]float32), nil, limit}
	od = append(od, cmteAll)
	cmteR := &TopOverallData{"cmte_donors_r", make(map[string]float32), nil, limit}
	od = append(od, cmteR)
	cmteD := &TopOverallData{"cmte_donors_d", make(map[string]float32), nil, limit}
	od = append(od, cmteD)
	cmteNa := &TopOverallData{"cmte_donors_na", make(map[string]float32), nil, limit}
	od = append(od, cmteNa)
	cmteOth := &TopOverallData{"cmte_donors_misc", make(map[string]float32), nil, limit}
	od = append(od, cmteOth)

	// cmte_recs
	cmteRecAll := &TopOverallData{"cmte_recs_all", make(map[string]float32), nil, limit}
	od = append(od, cmteRecAll)
	cmteRecR := &TopOverallData{"cmte_recs_r", make(map[string]float32), nil, limit}
	od = append(od, cmteRecR)
	cmteRecD := &TopOverallData{"cmte_recs_d", make(map[string]float32), nil, limit}
	od = append(od, cmteRecD)
	cmteRecNa := &TopOverallData{"cmte_recs_na", make(map[string]float32), nil, limit}
	od = append(od, cmteRecNa)
	cmteRecOth := &TopOverallData{"cmte_recs_misc", make(map[string]float32), nil, limit}
	od = append(od, cmteRecOth)

	// cmte_exp
	cmteExpAll := &TopOverallData{"cmte_exp_all", make(map[string]float32), nil, limit}
	od = append(od, cmteExpAll)
	cmteExpR := &TopOverallData{"cmte_exp_r", make(map[string]float32), nil, limit}
	od = append(od, cmteExpR)
	cmteExpD := &TopOverallData{"cmte_exp_d", make(map[string]float32), nil, limit}
	od = append(od, cmteExpD)
	cmteExpNa := &TopOverallData{"cmte_exp_na", make(map[string]float32), nil, limit}
	od = append(od, cmteExpNa)
	cmteExpOth := &TopOverallData{"cmte_exp_misc", make(map[string]float32), nil, limit}
	od = append(od, cmteExpOth)

	// cand
	candAll := &TopOverallData{"cand_all", make(map[string]float32), nil, limit}
	od = append(od, candAll)
	candR := &TopOverallData{"cand_r", make(map[string]float32), nil, limit}
	od = append(od, candR)
	candD := &TopOverallData{"cand_d", make(map[string]float32), nil, limit}
	od = append(od, candD)
	candNa := &TopOverallData{"cand_na", make(map[string]float32), nil, limit}
	od = append(od, candNa)
	candOth := &TopOverallData{"cand_misc", make(map[string]float32), nil, limit}
	od = append(od, candOth)

	// cand_exp
	candExpAll := &TopOverallData{"cand_exp_all", make(map[string]float32), nil, limit}
	od = append(od, candExpAll)
	candExpR := &TopOverallData{"cand_exp_r", make(map[string]float32), nil, limit}
	od = append(od, candExpR)
	candExpD := &TopOverallData{"cand_exp_d", make(map[string]float32), nil, limit}
	od = append(od, candExpD)
	candExpNa := &TopOverallData{"cand_exp_na", make(map[string]float32), nil, limit}
	od = append(od, candExpNa)
	candExpOth := &TopOverallData{"cand_exp_misc", make(map[string]float32), nil, limit}
	od = append(od, candExpOth)

	return od
}

// Entry represents a key/value pair from a Top X map and is used to sort and update the map.
type Entry struct {
	ID    string
	Total float32
}
