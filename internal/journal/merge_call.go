package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/mergerecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

// Like issue_creators, this is a derived projection, never caller-supplied
// Effects. INSERT OR IGNORE preserves the first call even within one second.
func projectMergeCall(tx *sql.Tx, event Event) error {
	if event.EventType != "task.approved" && event.EventType != "merge.observed" && event.EventType != "task.held" {
		return nil
	}
	var installed int
	if err := tx.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='merge_calls'`).Scan(&installed); err != nil || installed == 0 {
		return err
	}
	c := reviewrecord.PersonCall{Verdict: "MERGE", Actor: "human:mark", At: event.OccurredAt.UTC().Format(time.RFC3339)}
	e := event.Effects[0]
	c.SourceID, _ = e.PK["id"].(string)
	switch event.EventType {
	case "task.approved":
		var b approvalrecord.Binding
		raw, _ := e.Values["binding_json"].(string)
		if err := json.Unmarshal([]byte(raw), &b); err != nil {
			return err
		}
		c.HandoffID = b.HandoffID
		c.Actor, _ = e.Values["actor"].(string)
		c.Reason, _ = e.Values["reason"].(string)
	case "merge.observed":
		var m mergerecord.Merge
		raw, _ := e.Values["record_json"].(string)
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			return err
		}
		// Approval usually supplies the Assured call. A surviving approval may
		// belong to an earlier hand-off; this hand-off's observed merge still
		// supplies its own call, without replacing any earlier call here.
		c.HandoffID = m.HandoffID
		c.Reason = "Observed person's merge"
	case "task.held":
		c.Verdict = "HOLD"
		c.HandoffID, _ = e.Values["handoff_id"].(string)
		c.Actor, _ = e.Values["actor"].(string)
		c.Reason, _ = e.Values["reason"].(string)
	}
	_, err := tx.Exec(`INSERT OR IGNORE INTO merge_calls(handoff_id,verdict,actor,reason,source_id,at) VALUES(?,?,?,?,?,?)`, c.HandoffID, c.Verdict, c.Actor, c.Reason, c.SourceID, c.At)
	return err
}

func validateMergeHold(tx *sql.Tx, event Event) error {
	if event.CommandName != "task.hold" || event.Aggregate.Type != "decision" || event.Aggregate.Sequence != 1 || event.Repository != nil || event.ActorV2.ActorType != "human" {
		return fmt.Errorf("hold requires a person's decision")
	}
	n := 1
	if event.Idempotency != nil {
		n++
		if err := validateCurrentIdempotency(event); err != nil {
			return err
		}
	}
	if len(event.Effects) != n {
		return fmt.Errorf("hold requires exactly one decision")
	}
	e := event.Effects[0]
	if err := validateCurrentEffect(event, e); err != nil {
		return err
	}
	if !currentKeys(e.Values, "task_id handoff_id actor reason at") || e.PK["id"] != event.Aggregate.ID || e.Values["actor"] != "human:"+event.ActorV2.ActorID || e.Values["at"] != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("hold differs from Event")
	}
	reason, _ := e.Values["reason"].(string)
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("hold reason required")
	}
	task, _ := e.Values["task_id"].(string)
	var kind, state string
	if err := tx.QueryRow(`SELECT kind,status FROM issues WHERE id=?`, task).Scan(&kind, &state); err != nil {
		return err
	}
	h, err := reviewrecord.LatestHandoff(tx, task)
	if err != nil {
		return err
	}
	if kind != "task" || state != "in_review" || e.Values["handoff_id"] != h.ID || e.Values["at"].(string) < h.At {
		return fmt.Errorf("hold requires latest in-review hand-off")
	}
	if n == 2 {
		if err := validateCurrentEffect(event, event.Effects[1]); err != nil {
			return err
		}
		raw, _ := event.Effects[1].Values["result_json"].(string)
		var r struct {
			Fields   map[string]any `json:"fields"`
			Replayed bool           `json:"replayed"`
		}
		if err := decodeCheckJSON(raw, &r); err != nil {
			return err
		}
		want := map[string]any{"decision_id": event.Aggregate.ID, "task_id": task, "handoff_id": h.ID, "actor": e.Values["actor"], "reason": reason, "at": e.Values["at"], "verdict": "HOLD"}
		if r.Replayed || !currentValuesEqual(r.Fields, want) {
			return fmt.Errorf("hold retry differs from decision")
		}
	}
	return nil
}

// ReadMergeHoldRetry checks the entire cached response against its immutable
// source, as approval and review retry readers do for their own records.
func ReadMergeHoldRetry(conn reviewrecord.Queryer, raw, entity string) (map[string]any, error) {
	var stored struct {
		Fields   map[string]any `json:"fields"`
		Replayed bool           `json:"replayed"`
	}
	if err := decodeCheckJSON(raw, &stored); err != nil {
		return nil, err
	}
	var task, handoff, actor, reason, at string
	if err := conn.QueryRow(`SELECT task_id,handoff_id,actor,reason,at FROM merge_holds WHERE id=?`, entity).Scan(&task, &handoff, &actor, &reason, &at); err != nil {
		return nil, err
	}
	want := map[string]any{"decision_id": entity, "task_id": task, "handoff_id": handoff, "actor": actor, "reason": reason, "at": at, "verdict": "HOLD"}
	// ReplayIdempotent marks the served response replayed; the Event validator
	// separately requires the originally stored response to be unreplayed.
	if !currentValuesEqual(stored.Fields, want) {
		return nil, fmt.Errorf("hold retry differs from immutable decision")
	}
	return stored.Fields, nil
}
