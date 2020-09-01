package persist

import (
	"fmt"

	"github.com/elections/source/donations"
	"github.com/elections/source/protobuf"
	"github.com/golang/protobuf/proto"
)

// encode to protobuf
func encodeYrTotal(yt donations.YearlyTotal) ([]byte, error) { // move conversions to protobuf package?
	entry := &protobuf.YearlyTotal{
		ID:       yt.ID,
		Year:     yt.Year,
		Category: yt.Category,
		Party:    yt.Party,
		Total:    yt.Total,
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeOverallData failed: ", err)
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
