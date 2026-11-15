package journal

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type FindingCloseFields struct {
	TaskID      string `json:"task_id"`
	FindingID   string `json:"finding_id"`
	DecisionID  string `json:"decision_id"`
	Disposition string `json:"disposition"`
	Reason      string `json:"reason"`
}

type findingCloseResult struct {
	Fields     FindingCloseFields `json:"fields"`
	OccurredAt string             `json:"occurred_at"`
	Replayed   bool               `json:"replayed"`
	Worktree   string             `json:"worktree"`
}

func decodeFindingCloseResult(raw string) (findingCloseResult, error) {
	var result findingCloseResult
	if err := decodeCheckJSON(raw, &result); err != nil {
		return result, err
	}
	f := result.Fields
	if f.TaskID == "" || f.FindingID == "" || f.DecisionID == "" || f.Disposition != "closed" || strings.TrimSpace(f.Reason) == "" || result.Worktree != "" {
		return result, fmt.Errorf("invalid finding close response")
	}
	if _, err := time.Parse(time.RFC3339, result.OccurredAt); err != nil {
		return result, err
	}
	return result, nil
}

// ValidateFindingCloseRetry binds the response to the requested finding and
// stored entity before comparing it with immutable closure evidence.
func ValidateFindingCloseRetry(conn reviewrecord.Queryer, raw, findingID, entityID string) error {
	r, err := decodeFindingCloseResult(raw)
	if err != nil {
		return err
	}
	f := r.Fields
	if f.FindingID != findingID || entityID != findingID {
		return fmt.Errorf("finding retry differs from requested finding or stored entity")
	}
	var task, reason, at string
	err = conn.QueryRow(`SELECT subject_id,reason,decided_at FROM person_decisions WHERE id=? AND kind='close-finding' AND binding_id=?`, f.DecisionID, f.FindingID).Scan(&task, &reason, &at)
	if err != nil {
		return err
	}
	c, err := reviewrecord.ReadClosure(conn, f.FindingID)
	if err != nil {
		return err
	}
	if c == nil || c.DecisionID != f.DecisionID || c.Reason != f.Reason || task != f.TaskID || reason != f.Reason || at != r.OccurredAt {
		return fmt.Errorf("finding retry differs from recorded closure")
	}
	return nil
}

func validateFindingClosed(tx *sql.Tx, event Event) error {
	encoded, err := MarshalCanonicalJSON(event)
	if err != nil {
		return err
	}
	if len(encoded)+1 >= 1<<20 {
		return fmt.Errorf("finding close Event exceeds the journal's 1 MiB line budget")
	}

	if event.CommandName != "finding.close" || event.Aggregate.Type != "finding" || event.Aggregate.Sequence != 1 || event.Repository != nil || event.ActorV2.ActorType != "human" || !policy.Allows("human:"+event.ActorV2.ActorID, policy.PersonDecisions) {
		return fmt.Errorf("finding closure requires a person's decision")
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	if err := validateCurrentStages(event, retainedCurrentRules[event.EventType]); err != nil {
		return err
	}
	n := 2
	if event.Idempotency != nil {
		n++
		if err := validateCurrentIdempotency(event); err != nil {
			return err
		}
	}
	if len(event.Effects) != n {
		return fmt.Errorf("finding close requires one decision and one closure")
	}
	for _, e := range event.Effects {
		if err := validateCurrentEffect(event, e); err != nil {
			return err
		}
	}
	d, c := event.Effects[0], event.Effects[1]
	task, err := reviewrecord.FindingTask(tx, event.Aggregate.ID)
	if err != nil {
		return err
	}
	prior, err := reviewrecord.ReadClosure(tx, event.Aggregate.ID)
	if err != nil {
		return err
	}
	if prior != nil {
		return fmt.Errorf("finding is already closed")
	}
	reason, ok := d.Values["reason"].(string)
	if !ok || strings.TrimSpace(reason) == "" || d.Values["kind"] != "close-finding" || d.Values["actor"] != "human:"+event.ActorV2.ActorID || d.Values["subject_type"] != "task" || d.Values["subject_id"] != task || d.Values["binding_id"] != event.Aggregate.ID || d.Values["decided_at"] != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("person's closure decision differs from finding or Event")
	}
	if c.PK["finding_id"] != event.Aggregate.ID || c.Values["round_id"] != nil || c.Values["decision_id"] != d.PK["id"] || c.Values["disposition"] != "closed" || c.Values["reason"] != reason {
		return fmt.Errorf("finding closure differs from person's decision")
	}
	if event.Idempotency != nil {
		raw, ok := event.Effects[2].Values["result_json"].(string)
		if !ok {
			return fmt.Errorf("missing finding close response")
		}
		r, err := decodeFindingCloseResult(raw)
		if err != nil {
			return err
		}
		want := FindingCloseFields{task, event.Aggregate.ID, d.PK["id"].(string), "closed", reason}
		if r.Fields != want || r.OccurredAt != d.Values["decided_at"] {
			return fmt.Errorf("finding response differs from decision")
		}
	}
	return nil
}

// Preserve the new decision's immutability even through legacy Issue Effects
// and SQLite case aliases. Other decision kinds retain their existing rules.
func validateFindingDecisionEffect(tx *sql.Tx, event Event, e Effect) error {
	row := e.Values
	if e.Op == "insert" {
		var err error
		row, err = mergeEffectValues(e.PK, e.Values)
		if err != nil {
			return err
		}
	}
	for column, value := range row {
		if strings.EqualFold(column, "kind") && value == "close-finding" && (!event.IsV2() || event.EventType != "finding.closed") {
			return fmt.Errorf("finding decisions require their Event-v2 family")
		}
	}
	if e.Op == "update" || e.Op == "delete" {
		values, found, err := readCurrentValues(tx, "person_decisions", e.PK, []string{"kind"})
		if err != nil {
			return err
		}
		if found && values["kind"] == "close-finding" {
			return fmt.Errorf("finding closure decisions are immutable")
		}
	}
	return nil
}
