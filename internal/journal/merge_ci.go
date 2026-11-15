package journal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/mergerecord"
)

func validateMergeCI(tx *sql.Tx, event Event) error {
	rule := retainedCurrentRules[event.EventType]
	command := "merge.observe"
	if event.EventType == "ci.result-recorded" {
		command = "ci.record"
	}
	if event.CommandName != command || event.Aggregate.Type != "merge" || event.Repository != nil || event.Idempotency != nil || event.ActorV2.ActorType != "agent" && event.ActorV2.ActorType != "human" {
		return fmt.Errorf("invalid merge/CI command, actor or aggregate")
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	if err := validateCurrentStages(event, rule); err != nil {
		return err
	}
	if event.EventType == "merge.observed" && len(event.Effects) != 1 {
		return fmt.Errorf("merge observation requires one insert")
	}
	for _, effect := range event.Effects {
		if err := validateCurrentEffect(event, effect); err != nil {
			return err
		}
		if event.EventType == "merge.observed" {
			if err := validateMergeInsert(tx, event, effect); err != nil {
				return err
			}
		} else if err := validateCIInsert(tx, event, effect); err != nil {
			return err
		}
	}
	return nil
}

func validateMergeInsert(tx *sql.Tx, event Event, effect Effect) error {
	var record mergerecord.Merge
	if err := decodeCheckJSON(effect.Values["record_json"], &record); err != nil {
		return err
	}
	if err := record.Validate(); err != nil {
		return err
	}
	want := map[string]any{"id": record.ID, "task_id": record.TaskID, "handoff_id": record.HandoffID, "merge_commit": record.MergeCommit, "record_json": effect.Values["record_json"]}
	row, err := mergeEffectValues(effect.PK, effect.Values)
	if err != nil {
		return err
	}
	if !currentValuesEqual(row, want) || event.Aggregate.ID != record.ID || record.ObservedAt != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("merge observation differs from Event")
	}
	if err := validateObservationTask(tx, record.TaskID); err != nil {
		return err
	}
	var raw, task string
	if err := tx.QueryRow(`SELECT task_id,record_json FROM handoffs WHERE id=?`, record.HandoffID).Scan(&task, &raw); err != nil {
		return fmt.Errorf("merge hand-off: %w", err)
	}
	var handoff handoffrecord.Record
	if err := decodeCheckJSON(raw, &handoff); err != nil {
		return err
	}
	identity, err := github.ParsePRURL(handoff.PR)
	if err != nil {
		return err
	}
	if task != record.TaskID || handoff.TaskID != record.TaskID || handoff.ID != record.HandoffID || identity.Repository.String() != record.Repository || identity.Number != record.Number || handoff.Tree != record.HeadTree {
		return fmt.Errorf("merge observation does not bind to its task's hand-off PR and tree")
	}
	return nil
}

func validateCIInsert(tx *sql.Tx, event Event, effect Effect) error {
	var record mergerecord.CIResult
	if err := decodeCheckJSON(effect.Values["record_json"], &record); err != nil {
		return err
	}
	if err := record.Validate(); err != nil {
		return err
	}
	want := map[string]any{"id": record.ID, "task_id": record.TaskID, "merge_id": record.MergeID, "repository": record.Repository, "run_id": record.RunID, "attempt": record.Attempt, "record_json": effect.Values["record_json"]}
	row, err := mergeEffectValues(effect.PK, effect.Values)
	if err != nil {
		return err
	}
	if !currentValuesEqual(row, want) || event.Aggregate.ID != record.MergeID || record.ObservedAt != event.OccurredAt.UTC().Format(time.RFC3339) {
		return fmt.Errorf("CI result differs from Event")
	}
	if err := validateObservationTask(tx, record.TaskID); err != nil {
		return err
	}
	var raw string
	if err := tx.QueryRow(`SELECT record_json FROM merge_observations WHERE id=?`, record.MergeID).Scan(&raw); err != nil {
		return fmt.Errorf("CI merge observation: %w", err)
	}
	var merge mergerecord.Merge
	if err := decodeCheckJSON(raw, &merge); err != nil {
		return err
	}
	if merge.TaskID != record.TaskID || merge.Repository != record.Repository || merge.MergeCommit != record.HeadSHA {
		return fmt.Errorf("CI result does not bind to its merge observation")
	}
	return nil
}

func validateObservationTask(tx *sql.Tx, id string) error {
	var kind string
	if err := tx.QueryRow(`SELECT kind FROM issues WHERE id=?`, id).Scan(&kind); err != nil {
		return fmt.Errorf("observation task: %w", err)
	}
	if kind != "task" {
		return fmt.Errorf("merge/CI observation requires a task")
	}
	return nil
}
