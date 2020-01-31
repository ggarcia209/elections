package main

import (
	"fmt"

	"github.com/elections/donations"
)

type wrong struct {
	ID        string
	SomeField string
}

func main() {
	// encode test
	indv := donations.Individual{ID: "001"}
	cmte := donations.Committee{ID: "002"}
	cand := donations.Candidate{ID: "003"}
	disbRec := donations.DisbRecipient{ID: "004"}
	od := donations.TopOverallData{Category: "test"}
	w := wrong{ID: "12", SomeField: "poopity scoop"}

	objs := []interface{}{indv, cmte, cand, disbRec, od, nil, w}
	for i, obj := range objs {
		res, err := encodeToProto(obj)
		if err != nil {
			fmt.Println("failed: ", err)
			continue
		}
		fmt.Printf("%d: %s\n", i, res)
	}

	// decode test
	buckets := []string{"individuals", "committees", "candidates", "disbursement_recipients", "top_overall", "test", ""}
	for i, obj := range buckets {
		res, err := decodeFromProto(obj)
		if err != nil {
			fmt.Println("failed: ", err)
			continue
		}
		fmt.Printf("%d: %v\n", i, res)
	}

}

// EncodeToProto encodes an object interface to protobuf
func encodeToProto(obj interface{}) (string, error) {
	switch obj.(type) {
	case nil:
		return "nil", fmt.Errorf("EncodeToProto failed: nil interface")
	case donations.Individual:
		return "encodeToProto: Indvidiual", nil
	case donations.Committee:
		return "encodeToProto: Committee", nil
	case donations.Candidate:
		return "encodeToProto: Candidate", nil
	case donations.DisbRecipient:
		return "encodeToProto: DisbRecipient", nil
	case donations.TopOverallData:
		return "encodeToProto: TopOverallData", nil
	default:
		return "nil", fmt.Errorf("EncodeToProto failed: invalid interface type")
	}
}

// decodeFromProto encodes an object interface to protobuf
func decodeFromProto(bucket string) (interface{}, error) {
	switch bucket {
	case "":
		return nil, fmt.Errorf("decodeFromProto failed: nil bucket")
	case "individuals":
		data := donations.Individual{ID: "005"}
		return data, nil
	case "committees":
		data := donations.Committee{ID: "006"}
		return data, nil
	case "candidates":
		data := donations.Candidate{ID: "007"}
		return data, nil
	case "disbursement_recipients":
		data := donations.DisbRecipient{ID: "008"}
		return data, nil
	case "top_overall":
		data := donations.TopOverallData{Category: "test2"}
		return data, nil
	default:
		return nil, fmt.Errorf("decodeFromProto failed: invalid bucket")
	}
}
