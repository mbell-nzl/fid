package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

// ValidateHandoffResult protects the CLI's stored retry response, even when it
// is returned directly from an already-current projection.
func ValidateHandoffResult(raw string) error {
	var stored struct {
		Fields     map[string]any `json:"fields"`
		OccurredAt string         `json:"occurred_at"`
	}
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		return fmt.Errorf("decode hand-off response: %w", err)
	}
	fields := stored.Fields
	if len(fields) != 8 || fields["status"] != "in_review" {
		return fmt.Errorf("invalid hand-off response fields")
	}
	for _, name := range []string{"task_id", "handoff_id", "pr", "head", "base"} {
		if text, ok := fields[name].(string); !ok || text == "" {
			return fmt.Errorf("invalid hand-off response %s", name)
		}
	}
	var view handoffrecord.View
	encoded, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(encoded, &view); err != nil {
		return err
	}
	if !handoffrecord.ValidPR(view.PR) || !oidForHandoff(view.Head) || !oidForHandoff(view.Base) || view.RequiredChecks == nil || view.ResultIDs == nil || len(view.RequiredChecks) != len(view.ResultIDs) || !slices.IsSorted(view.RequiredChecks) {
		return fmt.Errorf("invalid hand-off response evidence")
	}
	for i, name := range view.RequiredChecks {
		if name == "" || view.ResultIDs[i] == "" || i > 0 && name == view.RequiredChecks[i-1] {
			return fmt.Errorf("invalid hand-off response checks")
		}
	}
	if _, err := time.Parse(time.RFC3339, stored.OccurredAt); err != nil {
		return fmt.Errorf("hand-off response time: %w", err)
	}
	return nil
}

func oidForHandoff(oid string) bool {
	if len(oid) != 40 && len(oid) != 64 {
		return false
	}
	for _, ch := range oid {
		if (ch < '0' || ch > '9') && (ch < 'a' || ch > 'f') {
			return false
		}
	}
	return true
}

func validateHandoff(tx *sql.Tx, event Event) error {
	if event.CommandName != "task.handoff" || event.Aggregate.Type != "task" || event.Repository != nil || event.ActorV2.ActorType != "agent" && event.ActorV2.ActorType != "human" || len(event.Effects) != 4 {
		return fmt.Errorf("hand-off requires its command, work actor and four atomic Effects")
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	if err := validateCurrentStages(event, retainedCurrentRules[event.EventType]); err != nil {
		return err
	}
	if err := validateCurrentIdempotency(event); err != nil {
		return err
	}
	insert := event.Effects[0]
	if err := validateCurrentEffect(event, insert); err != nil {
		return err
	}
	var record handoffrecord.Record
	if err := decodeCheckJSON(insert.Values["record_json"], &record); err != nil {
		return err
	}
	if err := record.Validate(); err != nil {
		return err
	}
	if record.ID != insert.PK["id"] || record.TaskID != insert.Values["task_id"] || record.HoldID != insert.Values["hold_id"] || !currentValuesEqual(record.PolicySnapshot, insert.Values["policy_snapshot"]) || record.TaskID != event.Aggregate.ID || record.At != event.OccurredAt.UTC().Format(time.RFC3339) || record.Actor != event.ActorV2.ActorType+":"+event.ActorV2.ActorID {
		return fmt.Errorf("hand-off receipt differs from Event")
	}
	if err := validateHandoffEvidence(tx, record); err != nil {
		return err
	}
	if err := validateHandoffChanges(tx, event, record); err != nil {
		return err
	}
	raw, ok := event.Effects[3].Values["result_json"].(string)
	if !ok {
		return fmt.Errorf("hand-off requires retry response")
	}
	if err := ValidateHandoffResult(raw); err != nil {
		return err
	}
	var stored struct {
		Fields     map[string]any `json:"fields"`
		OccurredAt string         `json:"occurred_at"`
	}
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		return err
	}
	if !currentValuesEqual(stored.Fields, record.Fields()) || stored.OccurredAt != record.At {
		return fmt.Errorf("hand-off response differs from receipt")
	}
	return nil
}

func validateHandoffEvidence(tx *sql.Tx, record handoffrecord.Record) error {
	var raw string
	var latest int
	if err := tx.QueryRow(`SELECT max(sequence) FROM policy_snapshots`).Scan(&latest); err != nil {
		return err
	}
	if record.PolicySnapshot != latest {
		return fmt.Errorf("hand-off must use policy in force")
	}
	if err := tx.QueryRow(`SELECT effective_json FROM policy_snapshots WHERE sequence=?`, record.PolicySnapshot).Scan(&raw); err != nil {
		return err
	}
	var effective policy.Effective
	if err := json.Unmarshal([]byte(raw), &effective); err != nil {
		return err
	}
	if !effective.TestAdequacySet {
		return fmt.Errorf("hand-off requires configured test paths")
	}
	required, err := policy.Required(effective, record.Lane, record.TaskID, record.ChangedPaths)
	if err != nil {
		return err
	}
	if !slices.Equal(required, record.RequiredChecks) {
		return fmt.Errorf("hand-off required checks differ from policy")
	}
	records, err := checkread.Read(tx, record.TaskID)
	if err != nil {
		return err
	}
	evaluation := checkread.Evaluate(record.TaskID, record.Tree, required, effective.Checks, records)
	if len(evaluation.Missing) != 0 {
		return fmt.Errorf("hand-off checks are unsatisfied")
	}
	for i, name := range required {
		if record.DefinitionDigests[name] != checkrecord.Digest(checkrecord.DefinitionFor(effective.Checks[name])) || record.ResultIDs[i] != evaluation.LatestResultIDs[name] {
			return fmt.Errorf("hand-off does not name the counting result and definition")
		}
	}
	return nil
}

func validateHandoffChanges(tx *sql.Tx, event Event, record handoffrecord.Record) error {
	var state, kind string
	var lane sql.NullString
	var revision int
	if err := tx.QueryRow(`SELECT kind,status,lane,revision FROM issues WHERE id=?`, record.TaskID).Scan(&kind, &state, &lane, &revision); err != nil {
		return err
	}
	if kind != "task" || state != "in_progress" || !lane.Valid || lane.String != record.Lane {
		return fmt.Errorf("hand-off requires its started task and lane")
	}
	want := Effect{Op: "update", Table: "issues", PK: map[string]any{"id": record.TaskID}, Values: map[string]any{"status": "in_review", "revision": revision + 1, "updated_at": record.At}, Preconditions: map[string]any{"status": "in_progress", "revision": revision}}
	if !currentValuesEqual(event.Effects[1], want) {
		return fmt.Errorf("hand-off task change differs from receipt")
	}
	var task, actor, since, expiry string
	var ended sql.NullString
	if err := tx.QueryRow(`SELECT task_id,actor,since,expires_at,ended_at FROM holds WHERE id=?`, record.HoldID).Scan(&task, &actor, &since, &expiry, &ended); err != nil {
		return err
	}
	if task != record.TaskID || actor != record.Actor || ended.Valid || record.At < since || record.At >= expiry {
		return fmt.Errorf("hand-off requires the exact current hold")
	}
	want = Effect{Op: "update", Table: "holds", PK: map[string]any{"id": record.HoldID}, Values: map[string]any{"ended_at": record.At, "end_reason": "review", "kept": true}, Preconditions: map[string]any{"ended_at": nil}}
	if !currentValuesEqual(event.Effects[2], want) {
		return fmt.Errorf("hand-off must end and keep its builder hold")
	}
	return nil
}
