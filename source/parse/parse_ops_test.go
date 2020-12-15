package parse

import (
	"fmt"
	"os"
	"testing"

	"github.com/elections/source/donations"
)

const (
	indv     = "txt_samples/test_indv.txt"
	cmteCont = "txt_samples/test_cmte_cont.txt"
	cn1      = "txt_samples/test_cn_1.txt"
	ccl      = "txt_samples/test_ccl.txt"
	cm       = "txt_samples/test_cm.txt"
	disb     = "txt_samples/test_disb.txt"
	pac      = "txt_samples/test_pac.txt"
)

func TestScanRow(t *testing.T) {
	var tests = []struct {
		row  string
		want map[int]string
	}{
		{"test0|test1|test2|test3", map[int]string{0: "test0", 1: "test1", 2: "test2", 3: "test3"}},
		{"test0|test1|test2", map[int]string{0: "test0", 1: "test1", 2: "test2"}},
		{"test0||test2|test3", map[int]string{0: "test0", 1: "", 2: "test2", 3: "test3"}},
		{"|||", map[int]string{0: "", 1: "", 2: "", 3: ""}},
	}

	for _, test := range tests {
		m := make(map[int]string)
		result := scanRow(test.row, m)
		for k, v := range result {
			if v != test.want[k] {
				t.Errorf("scanRow failed - row: %s; result: %v; want: %v", test.row, result, test.want)
				break
			}
		}
	}
}

// TestScanCandidates is sufficient to induce functionality and accuracy of the other
// parse.ScanObject functions, given each field of the donations.Object types correctly
// matches the corresponding keys contained within the 'fieldmap' variable located within
// each parse.ScanObject function. Functionality and accuracy are assumed as each
// parse.ScanObject function implements the same logic to generate each donations.Object
// type from the source .txt files.

func TestScanCandidates(t *testing.T) {
	var tests = []struct {
		filepath string
		want     []interface{}
	}{
		{cn1, []interface{}{
			&donations.Candidate{
				ID:          "H0AZ01184",
				Name:        "FLAKE, JEFF MR.",
				Party:       "REP",
				ElectnYr:    "2012",
				OfficeState: "AZ",
				Office:      "H",
				PCC:         "C00347260",
				City:        "MESA",
				State:       "AZ",
			},
			&donations.Candidate{
				ID:          "H0AZ01259",
				Name:        "GOSAR, PAUL DR.",
				Party:       "REP",
				ElectnYr:    "2018",
				OfficeState: "AZ",
				Office:      "H",
				PCC:         "C00461806",
				City:        "PRESCOTT",
				State:       "AZ",
			},
			&donations.Candidate{
				ID:          "H0AZ01333",
				Name:        "GRESSLEY, FORREST DAYL",
				Party:       "REP",
				ElectnYr:    "2010",
				OfficeState: "AZ",
				Office:      "H",
				PCC:         "C00481267",
				City:        "GILBERT",
				State:       "AZ",
			},
			&donations.Candidate{
				ID:          "H0AZ02166",
				Name:        "SCHMIDT II, JAMES A MR.",
				Party:       "REP",
				ElectnYr:    "2020",
				OfficeState: "AZ",
				Office:      "H",
				PCC:         "",
				City:        "DRAGOON",
				State:       "AZ",
			},
		}}, // fail cases below
		{cn1, []interface{}{
			&donations.Candidate{
				ID:          "H0AZ01234",
				Name:        "FLAKO, JEFF MR.",
				Party:       "REP",
				ElectnYr:    "2012",
				OfficeState: "AZ",
				Office:      "H",
				PCC:         "C00347260",
				City:        "TABLE",
				State:       "AZ",
			},
			&donations.Candidate{
				ID:          "H0AZ01259",
				Name:        "GOOSAR, PAUL DR.",
				Party:       "REP",
				ElectnYr:    "2018",
				OfficeState: "AZ",
				Office:      "H",
				PCC:         "C00461806",
				City:        "PRESCOTT",
				State:       "AZ",
			},
			&donations.Candidate{
				ID:          "H0AZ01333",
				Name:        "GRESSLEY, FORREST DAYL",
				Party:       "REP",
				ElectnYr:    "2010",
				OfficeState: "AZ",
				Office:      "H",
				PCC:         "C00481267",
				City:        "MODESTO",
				State:       "CA",
			},
			&donations.Candidate{
				ID:          "H0AZ02166",
				Name:        "SCHMIDT II, JAMES A MR.",
				Party:       "REP",
				ElectnYr:    "2020",
				OfficeState: "AZ",
				Office:      "H",
				PCC:         "",
				City:        "DRAGOON",
				State:       "AZ",
			},
		}},
	}

	for _, test := range tests {
		file, err := os.Open(test.filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		start := int64(0)
		objs, offset, err := ScanCandidates(file, start)
		if err != nil {
			t.Errorf("ScanCandidates failed - err: %v", err)
		}
		fmt.Println("byte offset:", offset)

		for i, obj := range objs {
			cand, wantCand := obj.(*donations.Candidate), test.want[i].(*donations.Candidate)
			switch {
			case cand.ID != wantCand.ID:
				t.Errorf("ScanCandidates failed - data - cand: %v; want: %v", cand.ID, wantCand.ID)
				fallthrough
			case cand.Name != wantCand.Name:
				t.Errorf("ScanCandidates failed - data - cand: %v; want: %v", cand.Name, wantCand.Name)
				fallthrough
			case cand.Party != wantCand.Party:
				t.Errorf("ScanCandidates failed - data - cand: %v; want: %v", cand.Party, wantCand.Party)
				fallthrough
			case cand.ElectnYr != wantCand.ElectnYr:
				t.Errorf("ScanCandidates failed - data - cand: %v; want: %v", cand.ElectnYr, wantCand.ElectnYr)
				fallthrough
			case cand.OfficeState != wantCand.OfficeState:
				t.Errorf("ScanCandidates failed - data - cand: %v; want: %v", cand.OfficeState, wantCand.OfficeState)
				fallthrough
			case cand.Office != wantCand.Office:
				t.Errorf("ScanCandidates failed - data - cand: %v; want: %v", cand.Office, wantCand.Office)
				fallthrough
			case cand.PCC != wantCand.PCC:
				t.Errorf("ScanCandidates failed - data - cand: %v; want: %v", cand.PCC, wantCand.PCC)
				fallthrough
			case cand.City != wantCand.City:
				t.Errorf("ScanCandidates failed - data - cand: %v; want: %v", cand.City, wantCand.City)
				fallthrough
			case cand.State != wantCand.State:
				t.Errorf("ScanCandidates failed - data - cand: %v; want: %v", cand.State, wantCand.State)
			}
		}
	}
}
