// Package approvalrecord reads immutable approvals and validates their recorded
// evidence without git, a terminal or the network.
package approvalrecord

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type Binding struct {
	Algorithm   string `json:"algorithm"`
	Fingerprint string `json:"fingerprint"`
	HandoffID   string `json:"handoff_id"`
	Base        string `json:"base"`
	Head        string `json:"head"`
}

type Record struct {
	ID      string  `json:"decision_id"`
	TaskID  string  `json:"task_id"`
	Actor   string  `json:"actor"`
	Reason  string  `json:"reason"`
	At      string  `json:"at"`
	Binding Binding `json:"binding"`
}

type Result struct {
	Decision Record `json:"decision"`
	At       string `json:"occurred_at"`
	Replayed bool   `json:"replayed"`
	EventID  string `json:"-"`
}

var hex40 = regexp.MustCompile(`^[0-9a-f]{40}$`)
var hex64 = regexp.MustCompile(`^[0-9a-f]{64}$`)

func (r Record) Validate() error {
	b := r.Binding
	if r.ID == "" || r.TaskID == "" || !policy.Allows(r.Actor, policy.PersonDecisions) || strings.TrimSpace(r.Reason) == "" {
		return fmt.Errorf("approval requires a person, task and reason")
	}
	if b.Algorithm != fingerprint.Algorithm || !hex64.MatchString(b.Fingerprint) || b.HandoffID == "" || !hex40.MatchString(b.Base) || !hex40.MatchString(b.Head) {
		return fmt.Errorf("invalid approval fingerprint binding")
	}
	if _, err := time.Parse(time.RFC3339, r.At); err != nil {
		return fmt.Errorf("approval time: %w", err)
	}
	return nil
}

func Admit(conn reviewrecord.Queryer, r Record) error {
	if err := r.Validate(); err != nil {
		return err
	}
	var kind, state string
	if err := conn.QueryRow(`SELECT kind,status FROM issues WHERE id=?`, r.TaskID).Scan(&kind, &state); err != nil {
		return err
	}
	if kind != "task" || state != "in_review" {
		return fmt.Errorf("approval task must be in_review")
	}
	h, err := reviewrecord.LatestHandoff(conn, r.TaskID)
	if err != nil {
		return err
	}
	if h.ID != r.Binding.HandoffID || h.Base != r.Binding.Base || h.Head != r.Binding.Head || r.At < h.At {
		return fmt.Errorf("approval must bind to the latest hand-off's recorded base and head")
	}
	return nil
}

func Read(conn reviewrecord.Queryer, task string) ([]Record, error) {
	rows, err := conn.Query(`SELECT id,subject_id,actor,reason,decided_at,binding_json FROM person_decisions WHERE kind='approve' AND subject_id=? ORDER BY rowid`, task)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	result := []Record{}
	for rows.Next() {
		var r Record
		var binding string
		if err := rows.Scan(&r.ID, &r.TaskID, &r.Actor, &r.Reason, &r.At, &binding); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(binding), &r.Binding); err != nil {
			return nil, err
		}
		if err := r.Validate(); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

type View struct {
	State    string  `json:"state"` // current, stale or missing
	Decision *Record `json:"decision"`
}

// Current is the shared reader for done. The caller supplies the latest
// hand-off's fingerprint; endpoint and hand-off IDs explain, but do not bind it.
func Current(conn reviewrecord.Queryer, task string, fp fingerprint.Result) (View, error) {
	records, err := Read(conn, task)
	if err != nil {
		return View{}, err
	}
	if len(records) == 0 {
		return View{State: "missing"}, nil
	}
	// Any approval of this exact change remains evidence, even if later decisions
	// approved a different change which was subsequently reverted.
	for i := len(records) - 1; i >= 0; i-- {
		r := records[i]
		if r.Binding.Algorithm == fp.Algorithm && r.Binding.Fingerprint == fp.Fingerprint {
			return View{State: "current", Decision: &r}, nil
		}
	}
	return View{State: "stale", Decision: &records[len(records)-1]}, nil
}
