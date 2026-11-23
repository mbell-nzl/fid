package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/defaultconflict"
	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

// DefaultConflictEffects records the observation and returns the kept hand-off
// reservation atomically. The following hold.taken uses the usual rework path.
func DefaultConflictEffects(conn reviewrecord.Queryer, task, hold, id string, evidence defaultconflict.Record, at time.Time) ([]Effect, error) {
	var revision int
	if err := conn.QueryRow(`SELECT revision FROM issues WHERE id=?`, task).Scan(&revision); err != nil {
		return nil, err
	}
	raw, err := json.Marshal(evidence)
	if err != nil {
		return nil, err
	}
	stamp := at.UTC().Format(time.RFC3339)
	return []Effect{
		{Op: "insert", Table: "default_conflicts", PK: map[string]any{"id": id}, Values: map[string]any{"task_id": task, "hold_id": hold, "record_json": string(raw)}},
		{Op: "update", Table: "issues", PK: map[string]any{"id": task}, Values: map[string]any{"status": "in_progress", "revision": revision + 1, "updated_at": stamp}, Preconditions: map[string]any{"status": "in_review", "revision": revision}},
		{Op: "update", Table: "holds", PK: map[string]any{"id": hold}, Values: map[string]any{"rework_at": stamp}, Preconditions: map[string]any{"kept": true, "end_reason": "review", "rework_at": nil}},
	}, nil
}

// Admission verifies the recorded evidence and its ownership, never rerunning
// git: historical objects and refs may be gone on replay or import.
func validateDefaultConflict(tx *sql.Tx, event Event) error {
	if event.CommandName != "next.claim" || event.Aggregate.Sequence != 1 || event.Repository != nil || event.Idempotency != nil || len(event.Effects) != 3 {
		return fmt.Errorf("default conflict requires next.claim and three atomic Effects")
	}
	insert := event.Effects[0]
	if err := validateCurrentEffect(event, insert); err != nil {
		return err
	}
	var evidence defaultconflict.Record
	if err := decodeCheckJSON(insert.Values["record_json"], &evidence); err != nil {
		return err
	}
	if insert.PK["id"] != event.Aggregate.ID || evidence.Cause != "default-conflict" || !strings.HasPrefix(evidence.DefaultRef, "origin/") || evidence.DefaultRef == "origin/HEAD" || !oidForHandoff(evidence.DefaultCommit) || !oidForHandoff(evidence.Head) || len(evidence.ConflictingPaths) == 0 {
		return fmt.Errorf("invalid default conflict evidence")
	}
	if err := validateConflictPaths(evidence.Evidence); err != nil {
		return err
	}
	task, _ := insert.Values["task_id"].(string)
	hold, _ := insert.Values["hold_id"].(string)
	h, err := reviewrecord.LatestHandoff(tx, task)
	if err != nil {
		return err
	}
	if evidence.HandoffID != h.ID || evidence.Head != h.Head || hold != h.HoldID {
		return fmt.Errorf("default conflict must name the latest hand-off and its hold/head")
	}
	fp := evidence.Binding
	if fp.Algorithm != fingerprint.Algorithm || !hexFingerprint(fp.Fingerprint) || fp.Base != h.Base || fp.Head != h.Head {
		return fmt.Errorf("invalid conflict fingerprint binding")
	}
	s, err := doneeval.GatherReviews("", tx, task, &fp)
	if err != nil {
		return err
	}
	awaiting, err := doneeval.AwaitingMerge(tx, s)
	if err != nil {
		return err
	}
	if !awaiting {
		return fmt.Errorf("default conflict requires completed reviews")
	}
	var n int
	if err := tx.QueryRow(`SELECT count(*) FROM holds WHERE id=? AND task_id=? AND kept=1 AND end_reason='review' AND ended_at IS NOT NULL AND rework_at IS NULL`, hold, task).Scan(&n); err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("default conflict requires the kept hand-off reservation")
	}
	expected, err := DefaultConflictEffects(tx, task, hold, event.Aggregate.ID, evidence, event.OccurredAt)
	if err != nil {
		return err
	}
	if !currentValuesEqual(event.Effects, expected) {
		return fmt.Errorf("default conflict Effects differ from the atomic return")
	}
	return nil
}

func hexFingerprint(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, ch := range value {
		if ch < '0' || ch > '9' && ch < 'a' || ch > 'f' {
			return false
		}
	}
	return true
}

func validateConflictPaths(e defaultconflict.Evidence) error {
	paths := e.ConflictingPaths
	if e.PathBytes != nil {
		if len(e.PathBytes) != len(paths) {
			return fmt.Errorf("conflicting path bytes differ from paths")
		}
		for i, path := range e.PathBytes {
			if string([]rune(string(path))) != string([]rune(paths[i])) {
				return fmt.Errorf("conflicting path bytes differ from paths")
			}
		}
		paths = policy.PathsFromBytes(e.PathBytes)
	}
	if !slices.IsSorted(paths) {
		return fmt.Errorf("conflicting paths must be sorted")
	}
	for i, path := range paths {
		if path == "" || i > 0 && path == paths[i-1] {
			return fmt.Errorf("invalid conflicting paths")
		}
	}
	return nil
}
