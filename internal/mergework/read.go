// Package mergework records already-observed GitHub facts and reads their
// durable history. The caller performs all network requests before calling it.
package mergework

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/mergerecord"
)

type Reader interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
}

// LatestMerge uses journal insertion order, never remote timestamps.
func LatestMerge(conn Reader, task string) (mergerecord.Merge, error) {
	return readMerge(conn, `SELECT record_json FROM merge_observations WHERE task_id=? ORDER BY rowid DESC LIMIT 1`, task)
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

// CIResults returns immutable attempts in journal insertion order for the
// specified task and merge observation; it never evaluates local checks.
func CIResults(conn Reader, task, merge string) ([]mergerecord.CIResult, error) {
	rows, err := conn.Query(`SELECT record_json FROM ci_results WHERE task_id=? AND merge_id=? ORDER BY rowid`, task, merge)
	if err != nil {
		return nil, fmt.Errorf("read CI results: %w", err)
	}
	defer func() { _ = rows.Close() }()
	results := []mergerecord.CIResult{}
	for rows.Next() {
		var raw string
		var record mergerecord.CIResult
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(raw), &record); err != nil {
			return nil, err
		}
		if err := record.Validate(); err != nil {
			return nil, err
		}
		results = append(results, record)
	}
	return results, rows.Err()
}
