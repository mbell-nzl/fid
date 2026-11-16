package journal

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func ReadCompletionRetry(conn reviewrecord.Queryer, raw, entity string) (completionrecord.Result, error) {
	var r completionrecord.Result
	if err := decodeCheckJSON(raw, &r); err != nil {
		return r, err
	}
	if r.Status != "done" || r.Completion.TaskID != entity {
		return r, fmt.Errorf("invalid completion response")
	}
	stored, err := completionrecord.Read(conn, entity)
	if err != nil {
		return r, err
	}
	if !reflect.DeepEqual(stored, r.Completion) {
		return r, fmt.Errorf("completion response differs from immutable record")
	}
	return r, nil
}

func validateCompletion(tx *sql.Tx, event Event) error {
	if event.CommandName != "task.done" || event.Aggregate.Type != "task" || event.Repository != nil || event.ActorV2.ActorType != "agent" && event.ActorV2.ActorType != "human" {
		return fmt.Errorf("completion requires a work actor and its command")
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	var open int
	if err := tx.QueryRow(`SELECT count(*) FROM questions WHERE task_id=? AND resolved_at IS NULL`, event.Aggregate.ID).Scan(&open); err != nil {
		return err
	}
	// Stage grammar permits repeated Effects. Completion needs an exact
	// shape so an extra release or retry cannot replace a question resolution.
	tables := []string{"completions", "issues", "holds"}
	for i := 0; i < open; i++ {
		tables = append(tables, "questions")
	}
	if event.Idempotency != nil {
		tables = append(tables, "idempotency_results")
		if err := validateCurrentIdempotency(event); err != nil {
			return err
		}
	}
	if len(event.Effects) != len(tables) {
		return fmt.Errorf("completion requires its record, state move, reservation release and open question")
	}
	for i, e := range event.Effects {
		if e.Table != tables[i] {
			return fmt.Errorf("completion requires exactly its record, state move, reservation release, open question resolutions and optional retry result")
		}
		if e.Table != "issues" {
			if err := validateCurrentEffect(event, e); err != nil {
				return err
			}
		}
	}
	first := event.Effects[0]
	var r completionrecord.Record
	if err := decodeCheckJSON(first.Values["record_json"], &r); err != nil {
		return err
	}
	if r.TaskID != event.Aggregate.ID || r.Actor != event.ActorV2.ActorType+":"+event.ActorV2.ActorID || r.At != event.OccurredAt.UTC().Format(time.RFC3339) || first.PK["id"] != r.ID || first.Values["task_id"] != r.TaskID || first.Values["handoff_id"] != r.HandoffID || first.Values["merge_id"] != r.MergeID {
		return fmt.Errorf("completion record differs from Event")
	}
	if err := completionrecord.Admit(tx, r); err != nil {
		return err
	}
	var revision int
	if err := tx.QueryRow(`SELECT revision FROM issues WHERE id=?`, r.TaskID).Scan(&revision); err != nil {
		return err
	}
	move := event.Effects[1]
	if !currentKeys(move.PK, "id") || move.PK["id"] != r.TaskID || !currentKeys(move.Values, "status revision updated_at") || !currentKeys(move.Preconditions, "status revision") || len(move.Prior) != 0 || move.Values["status"] != "done" || move.Values["updated_at"] != r.At || !currentValuesEqual(move.Values["revision"], revision+1) || move.Preconditions["status"] != "in_review" || !currentValuesEqual(move.Preconditions["revision"], revision) {
		return fmt.Errorf("completion must move its in_review task to done")
	}
	var hold string
	if err := tx.QueryRow(`SELECT hold_id FROM handoffs WHERE id=?`, r.HandoffID).Scan(&hold); err != nil {
		return err
	}
	var kept bool
	var ended, reason sql.NullString
	var task string
	if err := tx.QueryRow(`SELECT task_id,kept,ended_at,end_reason FROM holds WHERE id=?`, hold).Scan(&task, &kept, &ended, &reason); err != nil {
		return err
	}
	release := event.Effects[2]
	if task != r.TaskID || !kept || !ended.Valid || reason.String != "review" || !currentKeys(release.PK, "id") || release.PK["id"] != hold || !currentKeys(release.Values, "kept") || !currentValuesEqual(release.Values["kept"], false) || !currentKeys(release.Preconditions, "kept") || !currentValuesEqual(release.Preconditions["kept"], true) || len(release.Prior) != 0 {
		return fmt.Errorf("completion must release its review reservation")
	}
	if err := validateQuestions(tx, event); err != nil {
		return err
	}
	if event.Idempotency != nil {
		var result completionrecord.Result
		if err := decodeCheckJSON(event.Effects[len(event.Effects)-1].Values["result_json"], &result); err != nil {
			return err
		}
		if result.Status != "done" || result.Replayed || !reflect.DeepEqual(result.Completion, r) {
			return fmt.Errorf("completion retry response differs from record")
		}
	}
	return nil
}
