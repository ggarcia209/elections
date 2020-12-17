package persist

import (
	"fmt"
	"testing"

	"github.com/elections/source/donations"
)

// TestEncodeCand implements both persist.encodeCand & persist.decodeCand functions sequentially.
// Test passess if all values of the decoded protobuf byte slice match all values of
// the previously encoded donations.Candidate object.
// TestEncodeCand is sufficient to induce functionality and accuracy of the other
// persist.encodeObj & persist.decodeObj functions, given each field of the donations.Object
// types correctly matches the corresponding values for each protobuf set/get operation.
// Functionality and accuracy are assumed as each persist.encodeObj & persist.decodeObj
// function implements the same logic to generate each protobuf.Object & donations.Object
// type from the object to be encoded and corresponding encoded byte slice.
func TestEncodeCand(t *testing.T) {
	var tests = []struct {
		input []interface{}
		want  []interface{} // assert want == input after encode/decode
	}{
		{[]interface{}{
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
				ID:               "H0AZ01259",
				Name:             "GOSAR, PAUL DR.",
				Party:            "REP",
				ElectnYr:         "2018",
				OfficeState:      "AZ",
				Office:           "H",
				PCC:              "C00461806",
				City:             "PRESCOTT",
				State:            "AZ",
				TotalDirectInAmt: 90000,
			},
			&donations.Candidate{
				ID:                "H0AZ01333",
				Name:              "GRESSLEY, FORREST DAYL",
				Party:             "REP",
				ElectnYr:          "2010",
				OfficeState:       "AZ",
				Office:            "H",
				PCC:               "C00481267",
				City:              "GILBERT",
				State:             "AZ",
				DirectSendersAmts: map[string]float32{"test1": 1000.0, "test2": 2000.0},
			},
			&donations.Candidate{
				ID:                "H0AZ02166",
				Name:              "SCHMIDT II, JAMES A MR.",
				Party:             "REP",
				ElectnYr:          "2020",
				OfficeState:       "AZ",
				Office:            "H",
				PCC:               "",
				City:              "DRAGOON",
				State:             "AZ",
				TotalDirectInAmt:  90000,
				DirectSendersAmts: map[string]float32{"test1": 1000.0, "test2": 2000.0},
			},
		}, []interface{}{
			&donations.Candidate{
				ID:          "H0AZ01184",
				Name:        "FLAKE, JEFF MR.",
				Party:       "REP",
				OfficeState: "AZ",
				Office:      "H",
				PCC:         "C00347260",
				City:        "MESA",
				State:       "AZ",
			},
			&donations.Candidate{
				ID:               "H0AZ01259",
				Name:             "GOSAR, PAUL DR.",
				Party:            "REP",
				OfficeState:      "AZ",
				Office:           "H",
				PCC:              "C00461806",
				City:             "PRESCOTT",
				State:            "AZ",
				TotalDirectInAmt: 90000,
			},
			&donations.Candidate{
				ID:                "H0AZ01333",
				Name:              "GRESSLEY, FORREST DAYL",
				Party:             "REP",
				OfficeState:       "AZ",
				Office:            "H",
				PCC:               "C00481267",
				City:              "GILBERT",
				State:             "AZ",
				DirectSendersAmts: map[string]float32{"test1": 1000.0, "test2": 2000.0},
			},
			&donations.Candidate{
				ID:                "H0AZ02166",
				Name:              "SCHMIDT II, JAMES A MR.",
				Party:             "REP",
				OfficeState:       "AZ",
				Office:            "H",
				PCC:               "",
				City:              "DRAGOON",
				State:             "AZ",
				TotalDirectInAmt:  90000,
				DirectSendersAmts: map[string]float32{"test1": 1000.0, "test2": 2000.0},
			},
		}},
	}

	for _, test := range tests {
		for i, intf := range test.input {
			obj := *intf.(*donations.Candidate)
			data, err := encodeCand(obj)
			if err != nil {
				t.Errorf("encode/decode failed - err: %v", err)
				continue
			}

			cand, err := decodeCand(data)
			if err != nil {
				t.Errorf("decodeCand failed - err: %v", err)
				continue
			}
			wantCand := test.want[i].(*donations.Candidate)

			for k, v := range cand.DirectSendersAmts {
				if v != wantCand.DirectSendersAmts[k] {
					t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.DirectSendersAmts, wantCand.DirectSendersAmts)
				}
			}

			switch {
			case cand.ID != wantCand.ID:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.ID, wantCand.ID)
				fallthrough
			case cand.Name != wantCand.Name:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.Name, wantCand.Name)
				fallthrough
			case cand.Party != wantCand.Party:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.Party, wantCand.Party)
				fallthrough
			case cand.ElectnYr != wantCand.ElectnYr:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.ElectnYr, wantCand.ElectnYr)
				fallthrough
			case cand.OfficeState != wantCand.OfficeState:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.OfficeState, wantCand.OfficeState)
				fallthrough
			case cand.Office != wantCand.Office:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.Office, wantCand.Office)
				fallthrough
			case cand.PCC != wantCand.PCC:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.PCC, wantCand.PCC)
				fallthrough
			case cand.City != wantCand.City:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.City, wantCand.City)
				fallthrough
			case cand.State != wantCand.State:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.State, wantCand.State)
				fallthrough
			case cand.TotalDirectInAmt != wantCand.TotalDirectInAmt:
				t.Errorf("encode/decode failed - data - cand: %v; want: %v", cand.TotalDirectInAmt, wantCand.TotalDirectInAmt)
				fallthrough
			default:
				fmt.Println("pass")
			}
		}
	}
}
