package persist

import (
	"os"
	"testing"

	"github.com/elections/source/donations"
)

// var tests is used to test encode/decode, put/get, and batch write/read operations.
// Tests in this file test encode/decode and write/read operations in each function.
// Test objects are first encoded and saved to disk as protocol buffer objects then
// read from disk and decoded. Data integrity is verified by matching the each
// objects decoded ID to test.want[i].ID. All other fields can be assumed to be valid
// given successful tests in persist_cand_test.go.
var tests = []struct {
	IDs   []string
	input []interface{}
	want  []interface{}
}{
	{[]string{"test007", "indv00", "indv01"},
		[]interface{}{
			&donations.Individual{
				ID:            "test007",
				Name:          "James Bond",
				City:          "New York",
				State:         "NY",
				TotalOutAmt:   100000,
				RecipientsAmt: map[string]float32{"cmte00": 50000, "cmte01": 50000},
			},
			&donations.Individual{
				ID:          "indv00",
				Name:        "Eric Cartman",
				City:        "South Park",
				State:       "CO",
				TotalOutAmt: 100000,
			},
			&donations.Individual{
				ID:          "indv01",
				Name:        "Goofy",
				City:        "Disneyland",
				State:       "CA",
				TotalOutAmt: 100000,
			},
			&donations.Committee{
				ID:    "cmte00",
				Name:  "CIA",
				City:  "Langley",
				State: "VA",
			},
			&donations.CmteTxData{
				CmteID:             "cmte00",
				ContributionsInAmt: 50000,
				ContributionsInTxs: 1,
			},
			&donations.Candidate{
				ID:                "cand00",
				Name:              "Kanye West",
				City:              "Chicago",
				State:             "IL",
				TotalDirectInAmt:  20000000,
				DirectSendersAmts: map[string]float32{"indv01": 5000000, "indv02": 15000000},
			},
		},
		[]interface{}{
			&donations.Individual{
				ID:            "test007",
				Name:          "James Bond",
				City:          "New York",
				State:         "NY",
				TotalOutAmt:   100000,
				RecipientsAmt: map[string]float32{"cmte00": 50000, "cmte01": 50000},
			},
			&donations.Individual{
				ID:          "indv00",
				Name:        "Eric Cartman",
				City:        "South Park",
				State:       "CO",
				TotalOutAmt: 100000,
			},
			&donations.Individual{
				ID:          "indv01",
				Name:        "Goofy",
				City:        "Disneyland",
				State:       "CA",
				TotalOutAmt: 100000,
			},
			&donations.Committee{
				ID:    "cmte00",
				Name:  "CIA",
				City:  "Langley",
				State: "VA",
			},
			&donations.CmteTxData{
				CmteID:             "cmte00",
				ContributionsInAmt: 50000,
				ContributionsInTxs: 1,
			},
			&donations.Candidate{
				ID:                "cand00",
				Name:              "Kanye West",
				City:              "Chicago",
				State:             "IL",
				TotalDirectInAmt:  20000000,
				DirectSendersAmts: map[string]float32{"indv01": 5000000, "indv02": 15000000},
			},
		},
	},
}

// TestInit tests the create folder, databise file, and database bucket operatinos.
// Test passes if no error is returned. Idempotency is guaranteed by implementing
// methods from the built-in os package and BoltDB. Succesive calls to Init() will
// not overwrite the existing data.
func TestInit(t *testing.T) {
	err := Init("2020")
	if err != nil {
		t.Errorf("Init failed - err: %v", err)
	}
}

// TestEncodeToProto tests both encodeToProto and decodetoProto functions.
// Test passes if the ID of each encoded/decoded object matches test.want[i].ID
// All other fields can be assumed to be valid given succesful testing of the
// underlying encodeObject functions.
func TestEncodeToProto(t *testing.T) {
	for _, test := range tests {
		for i, obj := range test.want {
			bucket, ID, data, err := encodeToProto(obj)
			if err != nil {
				t.Errorf("encodeToProto failed - err: %v", err)
			}
			switch obj.(type) {
			case *donations.Individual:
				if bucket != "individuals" {
					t.Errorf("encodeToProto failed - bucket: %s; want: individuals", bucket)
				}
				if ID != test.want[i].(*donations.Individual).ID {
					t.Errorf("encodeToProto failed - ID: %s; want: %v", ID, test.want[i].(*donations.Individual).ID)
				}
			case *donations.Committee:
				if bucket != "committees" {
					t.Errorf("encodeToProto failed - bucket: %s; want: committees", bucket)
				}
				if ID != test.want[i].(*donations.Committee).ID {
					t.Errorf("encodeToProto failed - ID: %s; want: %v", ID, test.want[i].(*donations.Committee).ID)
				}
			case *donations.CmteTxData:
				if bucket != "cmte_tx_data" {
					t.Errorf("encodeToProto failed - bucket: %s; want: cmte_tx_data", bucket)
				}
				if ID != test.want[i].(*donations.CmteTxData).CmteID {
					t.Errorf("encodeToProto failed - ID: %s; want: %v", ID, test.want[i].(*donations.CmteTxData).CmteID)
				}
			case *donations.Candidate:
				if bucket != "candidates" {
					t.Errorf("encodeToProto failed - bucket: %s; want: candidates", bucket)
				}
				if ID != test.want[i].(*donations.Candidate).ID {
					t.Errorf("encodeToProto failed - ID: %s; want: %v", ID, test.want[i].(*donations.Candidate).ID)
				}
			case nil:
				t.Errorf("encodeToProto failed - nil interface")
			}

			res, err := decodeFromProto(bucket, data)
			if err != nil {
				t.Errorf("decodeFromProto failed - err: %v", err)
			}
			switch res.(type) {
			case *donations.Individual:
				if res.(*donations.Individual).ID != test.want[i].(*donations.Individual).ID {
					t.Errorf("encode/decode failed - data: %s; want: %s", res.(*donations.Individual).ID, test.want[i].(*donations.Individual).ID)
				}
			case *donations.Committee:
				if res.(*donations.Committee).ID != test.want[i].(*donations.Committee).ID {
					t.Errorf("encode/decode failed - data: %s; want: %s", res.(*donations.Committee).ID, test.want[i].(*donations.Committee).ID)
				}
			case *donations.CmteTxData:
				if res.(*donations.CmteTxData).CmteID != test.want[i].(*donations.CmteTxData).CmteID {
					t.Errorf("encode/decode failed - data: %s; want: %s", res.(*donations.CmteTxData).CmteID, test.want[i].(*donations.CmteTxData).CmteID)
				}
			case *donations.Candidate:
				if res.(*donations.Candidate).ID != test.want[i].(*donations.Candidate).ID {
					t.Errorf("encode/decode failed - data: %s; want: %s", res.(*donations.Candidate).ID, test.want[i].(*donations.Candidate).ID)
				}
			case nil:
				t.Errorf("decodeFromProto failed - nil interface returned")
			}
		}
	}
}

// TestStoreObjects tests the StoreObjects, BatchGetByID, BatchGetSequential,
// ViewDataByBucket functions.BatchGetByID, BatchGetSequential functions pass
// if objects returned match the corresponding objects in test.want.
// ViewDataByBucket passes if no errors returned and startKey value == "".
// startKey must always return "" when len(objs) < 1000.
func TestStoreObjects(t *testing.T) {
	OUTPUT_PATH = "."
	year := "2020"
	bucket := "individuals"
	Init(year)

	for _, test := range tests {
		// StoreObjects
		err := StoreObjects(year, test.input)
		if err != nil {
			t.Errorf("StoreObjects failed - err: %v", err)
		}

		// BatchGetByID
		objs, nilIDs, err := BatchGetByID(year, bucket, test.IDs)
		if err != nil {
			t.Errorf("BatchGetByID failed: %v", err)
		}
		if len(nilIDs) != 0 {
			t.Errorf("BatchGetByID: nilIDs found - verify test.input, test.IDs")
		}
		m := make(map[string]interface{})
		for i, obj := range objs {
			indv := obj.(*donations.Individual)
			if indv.ID != test.want[i].(*donations.Individual).ID {
				t.Errorf("BatchGetByID failed - data: ID: %s; want: %s", indv.ID, test.want[i].(*donations.Individual).ID)
			}
			m[indv.ID] = indv
		}

		// BatchGetSequential
		objs, startKey, err := BatchGetSequential(year, bucket, "", 100)
		if err != nil {
			t.Errorf("BatchGetSequential failed: %v", err)
		}
		if startKey != "" {
			t.Errorf("BatchGetSequential failed - startKey failed to reset: %s; want: ''", startKey)
		}
		for _, obj := range objs {
			indv := obj.(*donations.Individual)
			if m[indv.ID] == nil {
				t.Errorf("BatchGetSequential failed: ID '%s' not found", indv.ID)
			}
		}

		// ViewDataByBucket
		startKey, err = ViewDataByBucket(year, bucket, "")
		if err != nil {
			t.Errorf("ViewDataByBucket failed - err: %v", err)
		}
		if err != nil {
			t.Errorf("ViewDataByBucket failed - startKey failed to reset: %s; want: ''", startKey)
		}
	}

	err := os.RemoveAll("./db")
	if err != nil {
		t.Errorf("failed to remove ./db directory")
	}
}

// TestPutObject tests the PutObject and GetObject functions. Test passes if the ID of
// each object returned by GetObject() matches the ID of test.want[i].
func TestPutObject(t *testing.T) {
	bucket := ""
	ID := ""
	year := "2020"
	Init(year)

	for _, test := range tests {
		for i, obj := range test.input {
			err := PutObject(year, obj)
			if err != nil {
				t.Errorf("PutObject failed - err: %v", err)
			}
			switch obj.(type) {
			case *donations.Individual:
				bucket = "individuals"
				ID = obj.(*donations.Individual).ID
			case *donations.Committee:
				bucket = "committees"
				ID = obj.(*donations.Committee).ID
			case *donations.CmteTxData:
				bucket = "cmte_tx_data"
				ID = obj.(*donations.CmteTxData).CmteID
			case *donations.Candidate:
				bucket = "candidates"
				ID = obj.(*donations.Candidate).ID
			}
			obj, err = GetObject(year, bucket, ID)
			if err != nil {
				t.Errorf("GetObjects failed - err: %v", err)
			}
			switch obj.(type) {
			case *donations.Individual:
				if obj.(*donations.Individual).ID != test.want[i].(*donations.Individual).ID {
					t.Errorf("GetObject failed - data: %s; want: %s", obj.(*donations.Individual).ID, test.want[i].(*donations.Individual).ID)
				}
			case *donations.Committee:
				if obj.(*donations.Committee).ID != test.want[i].(*donations.Committee).ID {
					t.Errorf("GetObject failed - data: %s; want: %s", obj.(*donations.Committee).ID, test.want[i].(*donations.Committee).ID)
				}
			case *donations.CmteTxData:
				if obj.(*donations.CmteTxData).CmteID != test.want[i].(*donations.CmteTxData).CmteID {
					t.Errorf("GetObject failed - data: %s; want: %s", obj.(*donations.CmteTxData).CmteID, test.want[i].(*donations.CmteTxData).CmteID)
				}
			case *donations.Candidate:
				if obj.(*donations.Candidate).ID != test.want[i].(*donations.Candidate).ID {
					t.Errorf("GetObject failed - data: %s; want: %s", obj.(*donations.Candidate).ID, test.want[i].(*donations.Candidate).ID)
				}
			}
		}
	}
	err := os.RemoveAll("./db")
	if err != nil {
		t.Errorf("failed to remove ./db directory")
	}
}
