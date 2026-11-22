package workrecord

import (
	"encoding/json"
	"fmt"
	"math"
)

type Allowance struct {
	DecisionID string `json:"decision_id"`
	Rounds     int64  `json:"rounds"`
	Reason     string `json:"reason"`
	Actor      string `json:"actor"`
	At         string `json:"at"`
}

// AllowanceRounds accepts only the persisted allowance payload, including on
// replay and import. It has no fingerprint: its subject is the whole task.
func AllowanceRounds(raw string) (int64, error) {
	var body map[string]json.RawMessage
	var rounds int64
	if json.Unmarshal([]byte(raw), &body) != nil || len(body) != 1 || json.Unmarshal(body["rounds"], &rounds) != nil || rounds < 1 {
		return 0, fmt.Errorf("allowance requires rounds as an integer of at least 1")
	}
	return rounds, nil
}

// Allowances adds every allowance and retains the last decision in journal
// order as guidance. Saturation is sufficient: a count cannot exceed MaxInt64.
func Allowances(conn Queryer, task string, cap int64) (int64, *Allowance, error) {
	rows, err := conn.Query(`SELECT id,binding_json,reason,actor,decided_at FROM person_decisions WHERE kind='rework-allowance' AND subject_id=? ORDER BY rowid`, task)
	if err != nil {
		return 0, nil, err
	}
	defer func() { _ = rows.Close() }()
	var latest *Allowance
	for rows.Next() {
		var a Allowance
		var raw string
		if err := rows.Scan(&a.DecisionID, &raw, &a.Reason, &a.Actor, &a.At); err != nil {
			return 0, nil, err
		}
		a.Rounds, err = AllowanceRounds(raw)
		if err != nil {
			return 0, nil, err
		}
		if a.Rounds > math.MaxInt64-cap {
			cap = math.MaxInt64
		} else {
			cap += a.Rounds
		}
		latest = &a
	}
	return cap, latest, rows.Err()
}
