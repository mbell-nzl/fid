package handoffrecord

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Queryer interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
}

func Latest(conn Queryer, task string) (Record, error) {
	var raw string
	if err := conn.QueryRow(`SELECT record_json FROM handoffs WHERE task_id=? ORDER BY rowid DESC LIMIT 1`, task).Scan(&raw); err != nil {
		return Record{}, fmt.Errorf("review needs a hand-off: %w", err)
	}
	var h Record
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		return h, err
	}
	return h, nil
}
