// Package waiverrecord reads immutable waivers and validates their recorded
// evidence without git, a terminal or the network.
package waiverrecord

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

type Binding struct {
	Lane        string `json:"lane"`
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
	if b.Lane != "routine" && b.Lane != "standard" {
		return &Refusal{Code: "waiver-not-lower", Message: "waiver needs routine or standard, strictly below the hand-off lane"}
	}
	if r.ID == "" || r.TaskID == "" || !policy.Allows(r.Actor, policy.PersonDecisions) || strings.TrimSpace(r.Reason) == "" {
		return fmt.Errorf("waiver requires a person, task and reason")
	}
	if b.Algorithm != fingerprint.Algorithm || !hex64.MatchString(b.Fingerprint) || b.HandoffID == "" || !hex40.MatchString(b.Base) || !hex40.MatchString(b.Head) {
		return fmt.Errorf("invalid waiver fingerprint binding")
	}
	if _, err := time.Parse(time.RFC3339, r.At); err != nil {
		return fmt.Errorf("waiver time: %w", err)
	}
	return nil
}

func Admit(conn handoffrecord.Queryer, r Record) error {
	if err := r.Validate(); err != nil {
		return err
	}
	_, _, err := Context(conn, r.TaskID, r.Binding.Lane)
	if err != nil {
		return err
	}
	h, err := handoffrecord.Latest(conn, r.TaskID)
	if err != nil {
		return err
	}
	if h.ID != r.Binding.HandoffID || h.Base != r.Binding.Base || h.Head != r.Binding.Head || r.At < h.At {
		return fmt.Errorf("waiver must bind to the latest hand-off's recorded base and head")
	}
	return nil
}

type Refusal struct {
	Code, Message string
	Paths         []string
}

func (e *Refusal) Error() string { return e.Message }

// Context validates the state and frozen policy without git, so every admission
// path applies the same hard floor to the latest hand-off.
func Context(conn handoffrecord.Queryer, task, lane string) (handoffrecord.Record, policy.Floor, error) {
	var h handoffrecord.Record
	var floor policy.Floor
	var kind, state string
	if err := conn.QueryRow(`SELECT kind,status FROM issues WHERE id=?`, task).Scan(&kind, &state); err != nil {
		return h, floor, err
	}
	if kind != "task" || state != "in_review" {
		return h, floor, &Refusal{Code: "waiver-not-in-review", Message: "waiver task must be in_review"}
	}
	h, err := handoffrecord.Latest(conn, task)
	if errors.Is(err, sql.ErrNoRows) {
		return h, floor, &Refusal{Code: "waiver-handoff-required", Message: "waiver needs a hand-off"}
	}
	if err != nil {
		return h, floor, err
	}
	if (lane != "routine" && lane != "standard") || lane == h.Lane || h.Lane == "routine" {
		return h, floor, &Refusal{Code: "waiver-not-lower", Message: "waiver lane must be strictly below the hand-off lane"}
	}
	var raw string
	if err := conn.QueryRow(`SELECT effective_json FROM policy_snapshots WHERE sequence=?`, h.PolicySnapshot).Scan(&raw); err != nil {
		return h, floor, err
	}
	var p policy.Effective
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return h, floor, err
	}
	floor, err = policy.ComputeFloor(p, h.Paths())
	if err != nil {
		return h, floor, err
	}
	if floor.HardFloor != nil && *floor.HardFloor == "assured" {
		return h, floor, &Refusal{Code: "waiver-hard-floor", Message: "waiver cannot go below hard floor assured: " + strings.Join(floor.HardPaths, ", "), Paths: floor.HardPaths}
	}
	return h, floor, nil
}

func Read(conn handoffrecord.Queryer, task string) ([]Record, error) {
	var count int
	if err := conn.QueryRow(`SELECT count(*) FROM person_decisions WHERE kind='waive' AND subject_id=?`, task).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		return []Record{}, nil
	}
	rows, err := conn.Query(`SELECT id,subject_id,actor,reason,decided_at,binding_json FROM person_decisions WHERE kind='waive' AND subject_id=? ORDER BY rowid`, task)
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
func Current(conn handoffrecord.Queryer, task string, fp fingerprint.Result) (View, error) {
	records, err := Read(conn, task)
	if err != nil {
		return View{}, err
	}
	if len(records) == 0 {
		return View{State: "missing"}, nil
	}
	// The latest matching decision applies; endpoint and hand-off IDs explain
	// history but never determine currency.
	for i := len(records) - 1; i >= 0; i-- {
		r := records[i]
		if r.Matches(fp) {
			return View{State: "current", Decision: &r}, nil
		}
	}
	return View{State: "stale", Decision: &records[len(records)-1]}, nil
}

// Matches is the sole currency predicate; endpoints explain the recorded act.
func (r Record) Matches(fp fingerprint.Result) bool {
	return r.Binding.Algorithm == fp.Algorithm && r.Binding.Fingerprint == fp.Fingerprint
}
