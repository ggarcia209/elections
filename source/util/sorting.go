package util

import (
	"sort"
)

// SortedKV repesents a string/float pair representing an object's ID & relevant $ total
type SortedKV struct {
	ID    string
	Total float32
}

// SortedTotalsMap is a sorted list of SortedKV objects.
type SortedTotalsMap []SortedKV

func (s SortedTotalsMap) Len() int           { return len(s) }
func (s SortedTotalsMap) Less(i, j int) bool { return s[i].Total > s[j].Total }
func (s SortedTotalsMap) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// SortObjectTotals sorts a string/float32 map.
func SortMapObjectTotals(m map[string]float32) SortedTotalsMap {
	var es SortedTotalsMap
	for k, v := range m {
		es = append(es, SortedKV{ID: k, Total: v})
	}
	sort.Sort(es)

	return es
}
