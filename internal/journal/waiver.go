package journal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

func waiverFromEffect(e Effect) (waiverrecord.Record, error) {
	var r waiverrecord.Record
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
	if e.Values["kind"] != "waive" || e.Values["subject_type"] != "task" || e.Values["binding_id"] != r.Binding.Fingerprint {
		return r, fmt.Errorf("waiver decision differs from fingerprint binding")
	}
	return r, r.Validate()
}

func decodeWaiverResult(raw string) (waiverrecord.Result, error) {
	var r waiverrecord.Result
	if err := decodeCheckJSON(raw, &r); err != nil {
		return r, err
	}
	if err := r.Decision.Validate(); err != nil {
		return r, err
	}
	if r.At != r.Decision.At {
		return r, fmt.Errorf("waiver response differs from decision")
	}
	return r, nil
}

func ReadWaiverRetry(conn reviewrecord.Queryer, raw, entity string) (waiverrecord.Result, error) {
	r, err := decodeWaiverResult(raw)
	if err != nil {
		return r, err
	}
	if r.Decision.ID != entity {
		return r, fmt.Errorf("waiver retry differs from stored entity")
	}
	records, err := waiverrecord.Read(conn, r.Decision.TaskID)
	if err != nil {
		return r, err
	}
	for _, record := range records {
		if record == r.Decision {
			return r, nil
		}
	}
	return r, fmt.Errorf("waiver retry differs from immutable decision")
}

func validateWaiver(tx *sql.Tx, event Event) error {
	if event.CommandName != "task.waive" || event.Aggregate.Type != "decision" || event.Aggregate.Sequence != 1 || event.Repository != nil || event.ActorV2.ActorType != "human" {
		return fmt.Errorf("waiver requires a person's decision")
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
		return fmt.Errorf("waiver requires one decision")
	}
	for _, e := range event.Effects {
		if err := validateCurrentEffect(event, e); err != nil {
			return err
		}
	}
	r, err := waiverFromEffect(event.Effects[0])
	if err != nil {
		return err
	}
	if event.Aggregate.ID != r.ID || r.Actor != "human:"+event.ActorV2.ActorID || r.At != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("waiver decision differs from Event")
	}
	if err := waiverrecord.Admit(tx, r); err != nil {
		return err
	}
	if event.Idempotency != nil {
		raw, _ := event.Effects[1].Values["result_json"].(string)
		result, err := decodeWaiverResult(raw)
		if err != nil {
			return err
		}
		if result.Decision != r || result.Replayed {
			return fmt.Errorf("waiver response differs from decision")
		}
	}
	return nil
}
