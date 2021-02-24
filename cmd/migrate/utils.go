package main

import (
	"github.com/borodyadka/db-migrate"
	"sort"
)

// returns entries thats exists only in left
// complexity is M*N, no need to optimize it
func subEntries(left []migrate.Entry, right []migrate.Entry) []migrate.Entry {
	var result []migrate.Entry
outer:
	for _, l := range left {
		for _, r := range right {
			if l.Timestamp == r.Timestamp && l.Name == r.Name {
				continue outer
			}
		}
		result = append(result, l)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return (result[i].Timestamp < result[j].Timestamp) && (result[i].Name < result[j].Name)
	})
	return result
}
