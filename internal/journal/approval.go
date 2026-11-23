package journal

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func approvalFromEffect(e Effect) (approvalrecord.Record, error) {
	var r approvalrecord.Record
	id, _ := e.PK["id"].(string)
	task, _ := e.Values["subject_id"].(string)
	actor, _ := e.Values["actor"].(string)
	reason, _ := e.Values["reason"].(string)
	at, _ := e.Values["decided_at"].(string)
	raw, _ := e.Values["binding_json"].(string)
	if err := decodeCheckJSON(raw, &r.Binding); err != nil {
		return r, err
	}
	r.ID, r.TaskID, r.Actor, r.Reason, r.At = id, task, actor, reason, at
	if e.Values["kind"] != "approve" || e.Values["subject_type"] != "task" || e.Values["binding_id"] != r.Binding.Fingerprint {
		return r, fmt.Errorf("approval decision differs from fingerprint binding")
	}
	return r, r.Validate()
}

func decodeApprovalResult(raw string) (approvalrecord.Result, error) {
	var r approvalrecord.Result
	if err := decodeCheckJSON(raw, &r); err != nil {
		return r, err
	}
	if err := r.Decision.Validate(); err != nil {
		return r, err
	}
	if r.At != r.Decision.At {
		return r, fmt.Errorf("approval response differs from decision")
	}
	return r, nil
}

func ReadApprovalRetry(conn reviewrecord.Queryer, raw, entity string) (approvalrecord.Result, error) {
	r, err := decodeApprovalResult(raw)
	if err != nil {
		return r, err
	}
	if r.Decision.ID != entity {
		return r, fmt.Errorf("approval retry differs from stored entity")
	}
	records, err := approvalrecord.Read(conn, r.Decision.TaskID)
	if err != nil {
		return r, err
	}
	for _, record := range records {
		if record == r.Decision {
			return r, nil
		}
	}
	return r, fmt.Errorf("approval retry differs from immutable decision")
}

func validateApproval(tx *sql.Tx, event Event) error {
	if event.CommandName != "task.approve" || event.Aggregate.Type != "decision" || event.Aggregate.Sequence != 1 || event.Repository != nil || event.ActorV2.ActorType != "human" {
		return fmt.Errorf("approval requires a person's decision")
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	n := 1
	if event.Idempotency != nil {
		n++
		if err := validateCurrentIdempotency(event); err != nil {
			return err
		}
	}
	if len(event.Effects) != n {
		return fmt.Errorf("approval requires one decision")
	}
	for _, e := range event.Effects {
		if err := validateCurrentEffect(event, e); err != nil {
			return err
		}
	}
	r, err := approvalFromEffect(event.Effects[0])
	if err != nil {
		return err
	}
	if event.Aggregate.ID != r.ID || r.Actor != "human:"+event.ActorV2.ActorID || r.At != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("approval decision differs from Event")
	}
	if err := approvalrecord.Admit(tx, r); err != nil {
		return err
	}
	if event.Idempotency != nil {
		raw, _ := event.Effects[1].Values["result_json"].(string)
		result, err := decodeApprovalResult(raw)
		if err != nil {
			return err
		}
		if result.Decision != r || result.Replayed {
			return fmt.Errorf("approval response differs from decision")
		}
	}
	return nil
}

func validateApprovalDecisionEffect(tx *sql.Tx, event Event, e Effect) error {
	for _, row := range []map[string]any{e.PK, e.Values, e.Preconditions, e.Prior} {
		for column, value := range row {
			if strings.EqualFold(column, "binding_json") && value != nil && (!event.IsV2() || (event.EventType != "task.approved" && event.EventType != "rework.allowed" && event.EventType != "task.waived") || e.Op != "insert") {
				return fmt.Errorf("approval bindings require their Event-v2 family")
			}
			if strings.EqualFold(column, "kind") && value == "waive" && (!event.IsV2() || event.EventType != "task.waived") {
				return fmt.Errorf("waiver decisions require their Event-v2 family")
			}
			if strings.EqualFold(column, "kind") && value == "approve" && (!event.IsV2() || event.EventType != "task.approved") {
				return fmt.Errorf("approval decisions require their Event-v2 family")
			}
		}
	}
	if e.Op == "update" || e.Op == "delete" {
		values, found, err := readCurrentValues(tx, "person_decisions", e.PK, []string{"kind"})
		if err != nil {
			return err
		}
		if found && (values["kind"] == "approve" || values["kind"] == "waive") {
			return fmt.Errorf("approval decisions are immutable")
		}
	}
	return nil
}
