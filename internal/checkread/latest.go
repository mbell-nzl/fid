package checkread

import (
	"slices"

	"github.com/mbell-nzl/fid/internal/checkrecord"
)

// Latest keeps each name's last result and every unfinished start.
func Latest(records []checkrecord.Record) []checkrecord.Record {
	latest := map[string]string{}
	for _, record := range records {
		if record.Result != nil {
			latest[record.Start.Name] = record.Start.ID
		}
	}
	result := []checkrecord.Record{}
	for _, record := range records {
		if record.Result == nil || latest[record.Start.Name] == record.Start.ID {
			result = append(result, record)
		}
	}
	slices.SortStableFunc(result, func(a, b checkrecord.Record) int {
		if a.Start.Name < b.Start.Name {
			return -1
		}
		if a.Start.Name > b.Start.Name {
			return 1
		}
		return 0
	})
	return result
}
