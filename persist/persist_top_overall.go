package persist

import (
	"fmt"

	"github.com/elections/donations"
	"github.com/elections/protobuf"

	"github.com/golang/protobuf/proto"
)

// encode to protobuf
func encodeOverallData(od donations.TopOverallData) ([]byte, error) { // move conversions to protobuf package?
	entry := &protobuf.TopOverallData{
		Category:  od.Category,
		Amts:      od.Amts,
		Threshold: encodeThreshold(od.Threshold),
		SizeLimit: int32(od.SizeLimit),
	}
	data, err := proto.Marshal(entry)
	if err != nil {
		fmt.Println("encodeOverallData failed: ", err)
		return nil, fmt.Errorf("encodeOverallData failed: %v", err)
	}
	return data, nil
}

func encodeThreshold(entries []*donations.Entry) []*protobuf.Entry {
	var es []*protobuf.Entry
	for _, e := range entries {
		entry := &protobuf.Entry{
			ID:    e.ID,
			Total: e.Total,
		}
		es = append(es, entry)
	}
	return es
}

// decode from protobuf
func decodeOverallData(data []byte) (donations.TopOverallData, error) {
	od := &protobuf.TopOverallData{}
	err := proto.Unmarshal(data, od)
	if err != nil {
		fmt.Println("convProtoToIndv failed: ", err)
		return donations.TopOverallData{}, fmt.Errorf("convProtoToIndv failed: %v", err)
	}

	entry := donations.TopOverallData{
		Category:  od.GetCategory(),
		Amts:      od.GetAmts(),
		Threshold: decodeThreshold(od.GetThreshold()),
		SizeLimit: int(od.GetSizeLimit()),
	}

	if len(entry.Amts) == 0 {
		entry.Amts = make(map[string]float32)
	}

	return entry, nil
}

func decodeThreshold(es []*protobuf.Entry) []*donations.Entry {
	var entries []*donations.Entry
	for _, e := range es {
		entry := donations.Entry{
			ID:    e.GetID(),
			Total: e.GetTotal(),
		}
		entries = append(entries, &entry)
	}
	return entries
}
