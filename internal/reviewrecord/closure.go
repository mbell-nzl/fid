package reviewrecord

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Closure names the immutable authority that closed a finding. Open dispositions
// remain in their review round; only terminal dispositions have a Closure.
type Closure struct {
	RoundID     string `json:"round_id,omitempty"`
	DecisionID  string `json:"decision_id,omitempty"`
	Disposition string `json:"disposition"`
	Reason      string `json:"reason"`
}

func FindingTask(conn Queryer, id string) (string, error) {
	var task string
	err := conn.QueryRow(`SELECT r.task_id FROM round_findings f JOIN review_rounds r ON r.id=f.round_id WHERE f.id=?`, id).Scan(&task)
	if err != nil {
		return "", fmt.Errorf("finding must already exist: %w", err)
	}
	return task, nil
}

func ReadClosure(conn Queryer, id string) (*Closure, error) {
	var installed int
	if err := conn.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='finding_closures'`).Scan(&installed); err != nil {
		return nil, err
	}
	if installed == 0 {
		return nil, nil
	}
	var c Closure
	var round, decision sql.NullString
	err := conn.QueryRow(`SELECT round_id,decision_id,disposition,reason FROM finding_closures WHERE finding_id=?`, id).Scan(&round, &decision, &c.Disposition, &c.Reason)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	c.RoundID, c.DecisionID = round.String, decision.String
	return &c, nil
}

func AdmitDispositions(conn Queryer, r Round) error {
	for _, d := range r.Dispositions {
		// Existing projection rows establish journal order. A caller's time or
		// fingerprint cannot make an unrecorded or same-round finding earlier.
		for _, f := range r.Findings {
			if d.FindingID == f.ID {
				return fmt.Errorf("a round cannot dispose of its own finding")
			}
		}
		task, err := FindingTask(conn, d.FindingID)
		if err != nil {
			return err
		}
		if task != r.TaskID {
			return fmt.Errorf("disposition requires an earlier finding of the same task")
		}
		c, err := ReadClosure(conn, d.FindingID)
		if err != nil {
			return err
		}
		if c != nil {
			return fmt.Errorf("finding is already closed")
		}
	}
	return nil
}

// OpenMajors is the single task-wide gate reader for admission and done. It
// deliberately ignores fingerprints: an edit never closes an earlier finding.
func OpenMajors(conn Queryer, task string) ([]Finding, error) {
	var installed int
	if err := conn.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='finding_closures'`).Scan(&installed); err != nil {
		return nil, err
	}
	join, filter := "", ""
	if installed != 0 {
		join = " LEFT JOIN finding_closures c ON c.finding_id=f.id"
		filter = " AND c.finding_id IS NULL"
	}
	rows, err := conn.Query(`SELECT f.record_json FROM round_findings f JOIN review_rounds r ON r.id=f.round_id `+join+` WHERE r.task_id=?`+filter+` ORDER BY r.rowid,f.rowid`, task)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	findings := []Finding{}
	for rows.Next() {
		var raw string
		var f Finding
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(raw), &f); err != nil {
			return nil, err
		}
		if f.Severity == "major" {
			findings = append(findings, f)
		}
	}
	return findings, rows.Err()
}
