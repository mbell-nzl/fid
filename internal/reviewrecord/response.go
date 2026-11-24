package reviewrecord

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/policy"
)

// Response is a builder's or person's immutable answer or discovered work.
// A deferral closes a finding; an answer never does.
type Response struct {
	ID           string        `json:"id"`
	FindingID    string        `json:"finding_id"`
	HandoffID    string        `json:"handoff_id"`
	TaskID       string        `json:"task_id"`
	Actor        string        `json:"actor"`
	Text         string        `json:"text"`
	At           string        `json:"at"`
	DeferredTask *DeferredTask `json:"deferred_task,omitempty"`
}

type DeferredTask struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	EpicID *string `json:"epic_id"`
	Lane   string  `json:"lane"`
}

func ReadFinding(conn Queryer, id string) (Finding, error) {
	var raw string
	var f Finding
	if err := conn.QueryRow(`SELECT record_json FROM round_findings WHERE id=?`, id).Scan(&raw); err != nil {
		return f, fmt.Errorf("finding must already exist: %w", err)
	}
	err := json.Unmarshal([]byte(raw), &f)
	return f, err
}

func Responses(conn Queryer, finding string) ([]Response, error) {
	result := []Response{}
	var installed int
	if err := conn.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='finding_responses'`).Scan(&installed); err != nil {
		return nil, err
	}
	if installed == 0 {
		return result, nil
	}
	rows, err := conn.Query(`SELECT record_json FROM finding_responses WHERE finding_id=? ORDER BY rowid`, finding)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var raw string
		var r Response
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(raw), &r); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

func (r Response) Admit(conn Queryer) error {
	if !strings.HasPrefix(r.ID, "FRS-") || strings.TrimSpace(r.Text) == "" || !actor.MatchString(r.Actor) || !policy.Allows(r.Actor, policy.WorkSet) {
		return fmt.Errorf("finding response requires identity, work actor and text")
	}
	if at, err := time.Parse(time.RFC3339, r.At); err != nil || at.IsZero() || !strings.HasSuffix(r.At, "Z") {
		return fmt.Errorf("invalid finding response time")
	}
	task, err := FindingTask(conn, r.FindingID)
	if err != nil {
		return err
	}
	if task != r.TaskID {
		return fmt.Errorf("response belongs to another task")
	}
	h, err := LatestHandoff(conn, task)
	if err != nil {
		return err
	}
	if r.HandoffID != h.ID {
		return fmt.Errorf("response must bind to the latest hand-off")
	}
	built, err := BuiltBy(conn, task, r.Actor)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(r.Actor, "human:") && !built {
		return fmt.Errorf("only a builder or person may answer or defer a finding")
	}
	closed, err := ReadClosure(conn, r.FindingID)
	if err != nil {
		return err
	}
	if closed != nil {
		return fmt.Errorf("finding is already closed")
	}
	if d := r.DeferredTask; d != nil {
		f, err := ReadFinding(conn, r.FindingID)
		if err != nil {
			return err
		}
		if f.Severity != "minor" && f.Severity != "nit" {
			return fmt.Errorf("only minor and nit findings may be deferred")
		}
		if d.ID == "" || d.ID == task || strings.TrimSpace(d.Title) == "" || (d.Lane != "standard" && d.Lane != "assured") {
			return fmt.Errorf("deferral requires a new task, title and at least Standard lane")
		}
		if r.Text != DeferralCriterion(f) {
			return fmt.Errorf("deferral criterion differs from finding")
		}
		floor, err := DeferralLane(conn, f, r.HandoffID)
		if err != nil {
			return err
		}
		if floor == "assured" && d.Lane != "assured" {
			return fmt.Errorf("deferral lane is below the finding's path floor")
		}
		if d.EpicID != nil {
			var kind string
			if err := conn.QueryRow(`SELECT kind FROM issues WHERE id=?`, *d.EpicID).Scan(&kind); err != nil {
				return err
			}
			if kind != "epic" {
				return fmt.Errorf("deferral parent must be an epic")
			}
		}
	}
	return nil
}

// DeferralLane applies the frozen policy to the finding's known location.
// The normal check/hand-off floor still covers every path eventually changed.
func DeferralLane(conn Queryer, f Finding, handoff string) (string, error) {
	if f.Location == nil {
		return "standard", nil
	}
	var raw string
	if err := conn.QueryRow(`SELECT p.effective_json FROM policy_snapshots p JOIN handoffs h ON h.policy_snapshot=p.sequence WHERE h.id=?`, handoff).Scan(&raw); err != nil {
		return "", err
	}
	var effective policy.Effective
	if err := json.Unmarshal([]byte(raw), &effective); err != nil {
		return "", err
	}
	floor, err := policy.ComputeFloor(effective, []string{f.Location.File})
	return floor.Raise("standard"), err
}

func DeferralCriterion(f Finding) string {
	return fmt.Sprintf("Address %s: %s\nEvidence: %s\nFix: %s", f.ID, f.Summary, f.Evidence, f.Fix)
}
