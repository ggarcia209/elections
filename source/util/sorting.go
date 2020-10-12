// Package util contains operations for basic utility functions.
// This file contains operations for sorting different map types,
// by both key & value.
import (
	"sort"
)

// SortedKV repesents a string/float pair representing an object's ID & relevant $ total.
type SortedKV struct {
	ID    string
	Total float32
}

// SortedTotalsMap is a sorted list of SortedKV objects.
type SortedTotalsMap []SortedKV

func (s SortedTotalsMap) Len() int           { return len(s) }
func (s SortedTotalsMap) Less(i, j int) bool { return s[i].Total > s[j].Total }
func (s SortedTotalsMap) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// SortedCheckKV repesents a string/bool pair used to sort a set by key created from a string/bool map.
type SortedCheckKV struct {
	Key   string
	Check bool
}

// SortedCheckMap is a sorted list of SortedCheckKV objects.
type SortedCheckMap []SortedCheckKV

func (s SortedCheckMap) Len() int           { return len(s) }
func (s SortedCheckMap) Less(i, j int) bool { return s[i].Key < s[j].Key }
func (s SortedCheckMap) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// SortMapObjectTotals sorts a string/float32 map.
func SortMapObjectTotals(m map[string]float32) SortedTotalsMap {
	var es SortedTotalsMap
	for k, v := range m {
		es = append(es, SortedKV{ID: k, Total: v})
	}
	sort.Sort(es)
	return es
}

// SortCheckMap sorts a set created from string/bool map by ID.
func SortCheckMap(m map[string]bool) SortedCheckMap {
	var es SortedCheckMap
	for k, v := range m {
		es = append(es, SortedCheckKV{Key: k, Check: v})
	}
	sort.Sort(es)
	return es
}
