// Package persist contains operations for reading and writing disk data.
// Most operations in this package are intended to be performed on the
// admin local machine and are not intended to be used in the service logic.
// This file contains operations for encoding/decoding protobufs for the
// donations.YearlyTotal object.
package persist

import (
	"fmt"

	"github.com/elections/source/donations"
	"github.com/elections/source/protobuf"
	"github.com/golang/protobuf/proto"
)

func encodeYrTotal(yt donations.YearlyTotal) ([]byte, error) {
	entry := &protobuf.YearlyTotal{
		ID:       yt.ID,
		Year:     yt.Year,
		Category: yt.Category,
		Party:    yt.Party,
		Total:    yt.Total,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("encodeOverallData failed: %v", err)
	}
	return data, nil
}

func decodeYrTotal(data []byte) (donations.YearlyTotal, error) {
	pb := &protobuf.YearlyTotal{}
	err := proto.Unmarshal(data, pb)
	if err != nil {
		fmt.Println(err)
		return donations.YearlyTotal{}, fmt.Errorf("decodeYearlyTotal failed: %v", err)
	}

	yt := donations.YearlyTotal{
		ID:       pb.GetID(),
		Year:     pb.GetYear(),
		Category: pb.GetCategory(),
		Party:    pb.GetParty(),
		Total:    pb.GetTotal(),
	}
	return yt, nil
}
