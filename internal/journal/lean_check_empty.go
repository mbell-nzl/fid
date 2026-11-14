package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

func validateEmptyCheck(tx *sql.Tx, event Event) error {
	if len(event.Effects) != 2 || event.Effects[0].PK["id"] != event.Aggregate.ID {
		return fmt.Errorf("empty check requires one completed request and its retry state")
	}
	if err := validateCurrentIdempotency(event); err != nil {
		return err
	}
	values := event.Effects[0].Values
	var names []string
	if err := decodeCheckJSON(values["names_json"], &names); err != nil {
		return err
	}
	var response checkrecord.Response
	if err := decodeCheckJSON(values["response_json"], &response); err != nil {
		return err
	}
	if names == nil || len(names) != 0 || response.TaskID != values["task_id"] || values["actor"] != event.ActorV2.ActorType+":"+event.ActorV2.ActorID || !response.Complete || response.Lane == "" {
		return fmt.Errorf("invalid empty check request")
	}
	hold, ok := values["hold_id"].(string)
	if !ok {
		return fmt.Errorf("empty check requires its hold")
	}
	current, err := checkHoldCurrent(tx, checkrecord.Start{TaskID: response.TaskID, HoldID: hold, Actor: values["actor"].(string)}, event.OccurredAt)
	if err != nil {
		return err
	}
	var state, lane string
	if err := tx.QueryRow(`SELECT status,lane FROM issues WHERE id=? AND kind='task'`, response.TaskID).Scan(&state, &lane); err != nil {
		return err
	}
	if !current || state != "in_progress" || lane != response.Lane {
		return fmt.Errorf("empty check requires a started task, lane and current hold")
	}
	var raw string
	if err := tx.QueryRow(`SELECT effective_json FROM policy_snapshots ORDER BY sequence DESC LIMIT 1`).Scan(&raw); err != nil {
		return err
	}
	var effective policy.Effective
	if err := json.Unmarshal([]byte(raw), &effective); err != nil {
		return err
	}
	if !effective.Lanes[lane].ChecksSet {
		return fmt.Errorf("empty check requires configured lane checks")
	}
	return nil
}
