package journal

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/epicrecord"
)

// Protect the decision and the state move even through legacy/case aliases.
func validateEpicCloseOwnership(tx *sql.Tx, event Event, e Effect) error {
	if strings.EqualFold(e.Table, "person_decisions") {
		for _, row := range []map[string]any{e.PK, e.Values, e.Prior, e.Preconditions} {
			for column, value := range row {
				if strings.EqualFold(column, "kind") && value == "close-epic" && (!event.IsV2() || event.EventType != "epic.closed" || e.Op != "insert") {
					return fmt.Errorf("epic close decisions require their immutable Event-v2 family")
				}
			}
		}
		if e.Op == "update" || e.Op == "delete" {
			values, found, err := readCurrentValues(tx, "person_decisions", e.PK, []string{"kind"})
			if err != nil {
				return err
			}
			if found && values["kind"] == "close-epic" {
				return fmt.Errorf("epic close decisions are immutable")
			}
		}
	}
	if !strings.EqualFold(e.Table, "issues") || e.Op != "update" {
		return nil
	}
	for column, value := range e.Values {
		if !strings.EqualFold(column, "status") || value != "done" {
			continue
		}
		values, found, err := readCurrentValues(tx, "issues", e.PK, []string{"kind"})
		if err != nil {
			return err
		}
		if found && values["kind"] == "epic" && (!event.IsV2() || event.EventType != "epic.closed") {
			return fmt.Errorf("epic completion requires its person's decision Event")
		}
	}
	return nil
}
func validateEpicClose(tx *sql.Tx, event Event) error {
	encoded, err := MarshalCanonicalJSON(event)
	if err != nil {
		return err
	}
	if len(encoded)+1 >= 1<<20 {
		return fmt.Errorf("epic close Event exceeds the journal's 1 MiB line budget")
	}
	if !event.IsV2() || event.CommandName != "task.done" || event.Aggregate.Type != "epic" || event.Repository != nil || event.ActorV2.ActorType != "human" {
		return fmt.Errorf("epic close requires a person's decision and its command")
	}
	if err := validateV2Event(event); err != nil {
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
		return fmt.Errorf("epic close requires exactly one decision, one state move and optional retry result")
	}
	if err := validateCurrentStages(event, retainedCurrentRules[event.EventType]); err != nil {
		return err
	}
	decision, move := event.Effects[0], event.Effects[1]
	if decision.Table != "person_decisions" || decision.Op != "insert" || move.Table != "issues" || move.Op != "update" {
		return fmt.Errorf("epic close requires its decision and state move")
	}
	if err := validateCurrentEffect(event, decision); err != nil {
		return err
	}
	if !currentKeys(decision.Values, "kind actor subject_type subject_id reason binding_id decided_at") {
		return fmt.Errorf("epic close requires a complete decision")
	}
	d := decision.Values
	at := event.OccurredAt.UTC().Format(time.RFC3339)
	if d["kind"] != "close-epic" || d["subject_type"] != "epic" || d["subject_id"] != event.Aggregate.ID || d["binding_id"] != event.Aggregate.ID || d["actor"] != "human:"+event.ActorV2.ActorID || d["decided_at"] != at {
		return fmt.Errorf("epic close decision differs from Event")
	}
	reason, ok := d["reason"].(string)
	if !ok || strings.TrimSpace(reason) == "" {
		return fmt.Errorf("epic close requires a reason")
	}
	eligible, err := epicrecord.ReadEligibility(tx, event.Aggregate.ID)
	if err != nil {
		return err
	}
	if !eligible.CanClose() {
		return fmt.Errorf("epic close requires an open epic with no open task")
	}
	if !currentKeys(move.PK, "id") || move.PK["id"] != event.Aggregate.ID || !currentKeys(move.Values, "status revision updated_at") || !currentKeys(move.Preconditions, "status revision") || len(move.Prior) != 0 || move.Values["status"] != "done" || move.Values["updated_at"] != at || !currentValuesEqual(move.Values["revision"], eligible.Revision+1) || move.Preconditions["status"] != eligible.State || !currentValuesEqual(move.Preconditions["revision"], eligible.Revision) {
		return fmt.Errorf("epic close must move exactly its open epic with a revision bump")
	}
	if event.Idempotency != nil {
		last := event.Effects[2]
		if err := validateCurrentEffect(event, last); err != nil {
			return err
		}
		var r epicrecord.Result
		if err := decodeCheckJSON(last.Values["result_json"], &r); err != nil {
			return err
		}
		want := epicrecord.Decision{ID: decision.PK["id"].(string), Kind: "close-epic", Actor: d["actor"].(string), SubjectType: "epic", SubjectID: event.Aggregate.ID, BindingID: event.Aggregate.ID, Reason: reason, At: at}
		if r.Replayed || r.Status != "done" || r.Revision != eligible.Revision+1 || r.Decision != want || last.Values["result_entity_id"] != event.Aggregate.ID {
			return fmt.Errorf("epic close retry differs from decision")
		}
	}
	return nil
}
