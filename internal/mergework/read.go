package mergework

import (
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/mergerecord"
)

type Reader = mergerecord.Reader

func LatestMerge(conn Reader, task string) (mergerecord.Merge, error) {
	return mergerecord.LatestMerge(conn, task)
}
func CIResults(conn Reader, task, merge string) ([]mergerecord.CIResult, error) {
	return mergerecord.CIResults(conn, task, merge)
}
func readMerge(conn Reader, query string, args ...any) (mergerecord.Merge, error) {
	var raw string
	var record mergerecord.Merge
	if err := conn.QueryRow(query, args...).Scan(&raw); err != nil {
		return record, err
	}
	if err := json.Unmarshal([]byte(raw), &record); err != nil {
		return record, fmt.Errorf("decode merge observation: %w", err)
	}
	return record, record.Validate()
}
