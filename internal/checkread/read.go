package checkread

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// Read returns finished records in finish journal order, then unfinished
// starts in start journal order. Projection insertion order survives replay.
func Read(conn fiddb.Queryer, task string) ([]checkrecord.Record, error) {
	rows, err := conn.Query(`SELECT s.start_json,p.process_json,r.result_json FROM check_starts s LEFT JOIN check_processes p ON p.check_id=s.id LEFT JOIN check_results r ON r.check_id=s.id WHERE s.task_id=? ORDER BY r.rowid IS NULL,r.rowid,s.rowid`, task)
	if err != nil {
		return nil, fmt.Errorf("read checks: %w", err)
	}
	defer func() { _ = rows.Close() }()
	results := []checkrecord.Record{}
	for rows.Next() {
		var start string
		var process, result sql.NullString
		var record checkrecord.Record
		if err := rows.Scan(&start, &process, &result); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(start), &record.Start); err != nil {
			return nil, err
		}
		if process.Valid {
			if err := json.Unmarshal([]byte(process.String), &record.Process); err != nil {
				return nil, err
			}
		}
		if result.Valid {
			if err := json.Unmarshal([]byte(result.String), &record.Result); err != nil {
				return nil, err
			}
		}
		results = append(results, record)
	}
	return results, rows.Err()
}
