package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

func decodeCheckJSON(raw any, target any) error {
	text, ok := raw.(string)
	if !ok {
		return fmt.Errorf("check record requires JSON text")
	}
	if _, err := decodeCurrentJSON(text); err != nil {
		return err
	}
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func validateLeanCheck(tx *sql.Tx, event Event) error {
	rule, ok := retainedCurrentRules[event.EventType]
	if !ok || event.CommandName != "check.execute" || event.Aggregate.Type != "check" || event.Repository != nil || event.ActorV2.ActorType != "agent" && event.ActorV2.ActorType != "human" {
		return fmt.Errorf("invalid lean check family, command or actor")
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	if err := validateCurrentStages(event, rule); err != nil {
		return err
	}
	for _, effect := range event.Effects {
		if effect.Table == "holds" {
			continue
		}
		if err := validateCurrentEffect(event, effect); err != nil {
			return err
		}
		if effect.Table == "check_requests" && effect.Op == "update" {
			if effect.Preconditions["response_json"] != nil || effect.Values["response_json"] == nil {
				return fmt.Errorf("check request completes exactly once")
			}
		}
	}
	if event.EventType == "check.completed" {
		return validateEmptyCheck(tx, event)
	}
	if event.EventType == "check.stopped" {
		return validateCheckStopped(tx, event)
	}
	if event.EventType == "check.started" {
		return validateCheckStart(tx, event)
	}
	var raw, request string
	if err := tx.QueryRow(`SELECT start_json,request_id FROM check_starts WHERE id=?`, event.Aggregate.ID).Scan(&raw, &request); err != nil {
		return fmt.Errorf("read check start: %w", err)
	}
	var start checkrecord.Start
	if err := decodeCheckJSON(raw, &start); err != nil {
		return err
	}
	if start.Actor != event.ActorV2.ActorType+":"+event.ActorV2.ActorID {
		return fmt.Errorf("check actor differs from start")
	}
	var finished int
	if err := tx.QueryRow(`SELECT count(*) FROM check_results WHERE check_id=?`, start.ID).Scan(&finished); err != nil {
		return err
	}
	if finished != 0 {
		return fmt.Errorf("check already has a result")
	}
	for _, effect := range event.Effects {
		if id, ok := effect.PK["check_id"]; ok && id != start.ID {
			return fmt.Errorf("check Effect differs from aggregate")
		}
		if effect.Table == "check_requests" && effect.PK["id"] != request {
			return fmt.Errorf("check response differs from request")
		}
	}
	switch event.EventType {
	case "check.process-started":
		if len(event.Effects) != 1 || event.Idempotency != nil {
			return fmt.Errorf("process identity requires one observation")
		}
		var process checkrecord.Process
		if err := decodeCheckJSON(event.Effects[0].Values["process_json"], &process); err != nil {
			return err
		}
		return validateCheckProcess(process)
	case "check.hold-renewed":
		return validateCheckRenewal(tx, event, start)
	case "check.artifact-kept", "check.artifact-refused":
		return validateCheckFile(tx, event, start)
	case "check.result-recorded":
		return validateCheckResult(tx, event, start)
	}
	return fmt.Errorf("unsupported check Event")
}

func validateCheckProcess(process checkrecord.Process) error {
	if process.PID <= 0 || process.Machine == "" || process.StartTime == "" {
		return fmt.Errorf("incomplete check process identity")
	}
	return nil
}

func validateCheckTree(tree checkrecord.Tree) error {
	if !policyCommitPattern.MatchString(tree.Head) || !policyCommitPattern.MatchString(tree.Tree) || tree.Dirty != (tree.DirtyDigest != nil) || tree.DirtyDigest != nil && !currentDigest(*tree.DirtyDigest) {
		return fmt.Errorf("invalid check tree observation")
	}
	return nil
}

func checkHoldCurrent(tx *sql.Tx, start checkrecord.Start, at time.Time) (bool, error) {
	var task, actor, since, expires string
	var ended sql.NullString
	if err := tx.QueryRow(`SELECT task_id,actor,since,expires_at,ended_at FROM holds WHERE id=?`, start.HoldID).Scan(&task, &actor, &since, &expires, &ended); err != nil {
		return false, err
	}
	from, err := time.Parse(time.RFC3339, since)
	if err != nil {
		return false, err
	}
	until, err := time.Parse(time.RFC3339, expires)
	if err != nil {
		return false, err
	}
	return task == start.TaskID && actor == start.Actor && !ended.Valid && !at.Before(from) && at.Before(until), nil
}

func validateCheckStart(tx *sql.Tx, event Event) error {
	if len(event.Effects) != 2 && len(event.Effects) != 3 {
		return fmt.Errorf("check start requires one start and its retry state")
	}
	starts, requests := 0, 0
	for _, effect := range event.Effects {
		if effect.Table == "check_starts" {
			starts++
		}
		if effect.Table == "check_requests" {
			requests++
		}
	}
	if starts != 1 || requests > 1 {
		return fmt.Errorf("check start bundle has duplicate records")
	}
	if err := validateCurrentIdempotency(event); err != nil {
		return err
	}
	var effect Effect
	var requestEffect *Effect
	for i, e := range event.Effects {
		if e.Table == "check_starts" {
			effect = e
		}
		if e.Table == "check_requests" {
			requestEffect = &event.Effects[i]
		}
	}
	var start checkrecord.Start
	if err := decodeCheckJSON(effect.Values["start_json"], &start); err != nil {
		return err
	}
	if start.Definition.Artifacts == nil || start.ID != event.Aggregate.ID || effect.PK["id"] != start.ID || start.TaskID != effect.Values["task_id"] || start.Name != effect.Values["name"] || start.Actor != event.ActorV2.ActorType+":"+event.ActorV2.ActorID || start.StartedAt != event.OccurredAt.UTC().Format(time.RFC3339Nano) || start.DefinitionDigest != checkrecord.Digest(start.Definition) || !currentDigest(start.EnvironmentDigest) {
		return fmt.Errorf("check start bindings differ")
	}
	if err := validateCheckProcess(start.Owner); err != nil {
		return err
	}
	if err := validateCheckTree(start.Tree); err != nil {
		return err
	}
	current, err := checkHoldCurrent(tx, start, event.OccurredAt)
	if err != nil {
		return err
	}
	if !current {
		return fmt.Errorf("check start requires its current hold")
	}
	var status string
	var lane sql.NullString
	if err := tx.QueryRow(`SELECT status,lane FROM issues WHERE id=? AND kind='task'`, start.TaskID).Scan(&status, &lane); err != nil {
		return err
	}
	if status != "in_progress" || !lane.Valid {
		return fmt.Errorf("check requires a started task with a lane")
	}
	var policyJSON string
	if err := tx.QueryRow(`SELECT effective_json FROM policy_snapshots WHERE sequence=?`, start.PolicySnapshot).Scan(&policyJSON); err != nil {
		return err
	}
	var effective policy.Effective
	if err := json.Unmarshal([]byte(policyJSON), &effective); err != nil {
		return err
	}
	check, defined := effective.Checks[start.Name]
	if !defined || !effective.Lanes[lane.String].ChecksSet || !reflect.DeepEqual(start.Definition.Run, check.Run) || start.Definition.Timeout != check.Timeout || !slices.Equal(start.Definition.Artifacts, check.Artifacts) || !slices.Equal(start.Definition.Environment, checkrecord.EnvironmentNames()) {
		return fmt.Errorf("check definition differs from trusted policy")
	}
	var names []string
	if requestEffect != nil {
		if requestEffect.PK["id"] != start.ID || effect.Values["request_id"] != start.ID || requestEffect.Values["task_id"] != start.TaskID || requestEffect.Values["hold_id"] != start.HoldID || requestEffect.Values["actor"] != start.Actor || requestEffect.Values["response_json"] != nil {
			return fmt.Errorf("invalid check request binding")
		}
		if err := decodeCheckJSON(requestEffect.Values["names_json"], &names); err != nil {
			return err
		}
	} else {
		var task, hold, actor, raw string
		var response sql.NullString
		if err := tx.QueryRow(`SELECT task_id,hold_id,actor,names_json,response_json FROM check_requests WHERE id=?`, effect.Values["request_id"]).Scan(&task, &hold, &actor, &raw, &response); err != nil {
			return err
		}
		if response.Valid || task != start.TaskID || hold != start.HoldID || actor != start.Actor {
			return fmt.Errorf("check request changed or finished")
		}
		if err := decodeCheckJSON(raw, &names); err != nil {
			return err
		}
	}
	if len(names) == 0 || !slices.IsSorted(names) || len(slices.Compact(slices.Clone(names))) != len(names) {
		return fmt.Errorf("check request requires unique ordered names")
	}
	for _, name := range names {
		if _, ok := effective.Checks[name]; !ok {
			return fmt.Errorf("check request has an undefined name")
		}
	}
	var count, results int
	if requestEffect == nil {
		if err := tx.QueryRow(`SELECT count(*),count(r.check_id) FROM check_starts s LEFT JOIN check_results r ON r.check_id=s.id WHERE s.request_id=?`, effect.Values["request_id"]).Scan(&count, &results); err != nil {
			return err
		}
	}
	if count != results || count >= len(names) || names[count] != start.Name {
		return fmt.Errorf("checks must start in name order after the previous result")
	}
	return nil
}

func validateCheckRenewal(tx *sql.Tx, event Event, start checkrecord.Start) error {
	if len(event.Effects) != 1 || event.Idempotency != nil {
		return fmt.Errorf("invalid check renewal bundle")
	}
	effect := event.Effects[0]
	current, err := checkHoldCurrent(tx, start, event.OccurredAt)
	if err != nil {
		return err
	}
	var expires string
	if err := tx.QueryRow(`SELECT expires_at FROM holds WHERE id=?`, start.HoldID).Scan(&expires); err != nil {
		return err
	}
	if !current || effect.PK["id"] != start.HoldID || !currentKeys(effect.PK, "id") || !currentKeys(effect.Values, "last_heartbeat_at expires_at") || !currentKeys(effect.Preconditions, "ended_at expires_at") || effect.Preconditions["ended_at"] != nil || effect.Preconditions["expires_at"] != expires || effect.Values["last_heartbeat_at"] != event.OccurredAt.UTC().Format(time.RFC3339) || effect.Values["expires_at"] != event.OccurredAt.Add(time.Hour).UTC().Format(time.RFC3339) {
		return fmt.Errorf("check renewal differs from its current hold")
	}
	return nil
}

func validateCheckFile(tx *sql.Tx, event Event, start checkrecord.Start) error {
	if event.Idempotency != nil {
		return fmt.Errorf("file observations do not carry retry state")
	}
	file := event.Effects[len(event.Effects)-1]
	path, _ := file.PK["path"].(string)
	if !slices.Contains(start.Definition.Artifacts, path) {
		return fmt.Errorf("kept file was not listed")
	}
	if event.EventType == "check.artifact-refused" {
		reason, _ := file.Values["reason"].(string)
		if len(event.Effects) != 1 || file.Values["artifact_id"] != nil || !slices.Contains([]string{"outside-worktree", "missing", "non-regular", "oversized", "changed-during-copy"}, reason) {
			return fmt.Errorf("invalid file refusal")
		}
		return nil
	}
	if file.Values["reason"] != nil || file.Values["artifact_id"] == nil || len(event.Effects) > 2 {
		return fmt.Errorf("invalid kept file")
	}
	if len(event.Effects) == 2 {
		artifact := event.Effects[0]
		if artifact.Op == "update" {
			if artifact.PK["id"] != file.Values["artifact_id"] || artifact.Values["availability"] != "available" || artifact.Preconditions["availability"] == "available" {
				return fmt.Errorf("invalid restored Artifact availability")
			}
			var size int64
			if err := tx.QueryRow(`SELECT size_bytes FROM referenced_artifacts WHERE id=?`, artifact.PK["id"]).Scan(&size); err != nil {
				return err
			}
			if size < 0 || size > checkrecord.ArtifactLimit {
				return fmt.Errorf("check Artifact exceeds its limit")
			}
			return nil
		}
		if artifact.PK["id"] != file.Values["artifact_id"] || !currentDigest(artifact.Values["digest"].(string)) || artifact.Values["availability"] != "available" || artifact.Values["location_policy"] != "content-addressed" || artifact.Values["retention_class"] != "evidence" || artifact.Values["expires_at"] != nil || artifact.Values["created_at"] != event.OccurredAt.UTC().Format(time.RFC3339Nano) {
			return fmt.Errorf("invalid check Artifact binding")
		}
		var size int64
		data, _ := json.Marshal(artifact.Values["size_bytes"])
		if json.Unmarshal(data, &size) != nil || size < 0 || size > checkrecord.ArtifactLimit {
			return fmt.Errorf("check Artifact exceeds its limit")
		}
	} else {
		var size int64
		if err := tx.QueryRow(`SELECT size_bytes FROM referenced_artifacts WHERE id=? AND availability='available'`, file.Values["artifact_id"]).Scan(&size); err != nil {
			return err
		}
		if size > checkrecord.ArtifactLimit {
			return fmt.Errorf("check Artifact exceeds its limit")
		}
	}
	return nil
}

func validateCheckResult(tx *sql.Tx, event Event, start checkrecord.Start) error {
	if len(event.Effects) != 2 && len(event.Effects) != 3 {
		return fmt.Errorf("check result requires one result and its retry state")
	}
	results, requests := 0, 0
	for _, effect := range event.Effects {
		if effect.Table == "check_results" {
			results++
		}
		if effect.Table == "check_requests" {
			requests++
		}
	}
	if results != 1 || requests > 1 {
		return fmt.Errorf("check result bundle has duplicate records")
	}
	if err := validateCurrentIdempotency(event); err != nil {
		return err
	}
	var result checkrecord.Result
	if err := decodeCheckJSON(event.Effects[0].Values["result_json"], &result); err != nil {
		return err
	}
	if result.FinishedAt != event.OccurredAt.UTC().Format(time.RFC3339Nano) || !slices.Contains([]string{"passed", "failed", "timed-out", "cancelled", "spawn-failed", "signaled"}, result.Outcome) || result.OutputBytes < 0 || len(result.OutputExcerpt) > checkrecord.ExcerptLimit || !currentDigest(result.OutputDigest) || result.OutputTruncated != (result.OutputBytes > checkrecord.OutputLimit) || result.Reasons == nil || result.Artifacts == nil || result.RefusedArtifacts == nil {
		return fmt.Errorf("invalid check result observation")
	}
	if err := validateCheckTree(result.Tree); err != nil {
		return err
	}
	if result.Outcome == "passed" && (result.ExitCode == nil || *result.ExitCode != 0 || result.Signal != nil) || result.Outcome == "failed" && (result.ExitCode == nil || *result.ExitCode <= 0) || result.Outcome == "signaled" && (result.Signal == nil || result.ExitCode != nil) || result.Outcome == "spawn-failed" && (result.ExitCode != nil || result.Signal != nil) {
		return fmt.Errorf("check outcome differs from exit status")
	}
	var processCount int
	if err := tx.QueryRow(`SELECT count(*) FROM check_processes WHERE check_id=?`, start.ID).Scan(&processCount); err != nil {
		return err
	}
	if result.Outcome != "spawn-failed" && processCount != 1 {
		return fmt.Errorf("check result is missing its child identity")
	}
	if result.Outcome != "spawn-failed" && (result.Toolchain.Executable == "" || !filepath.IsAbs(result.Toolchain.Executable) || !currentDigest(result.Toolchain.Digest)) {
		return fmt.Errorf("check toolchain is incomplete")
	}
	if result.Toolchain.OS == "" || result.Toolchain.Arch == "" || result.Toolchain.FidVersion == "" {
		return fmt.Errorf("check toolchain platform is incomplete")
	}
	started, err := time.Parse(time.RFC3339Nano, start.StartedAt)
	if err != nil || event.OccurredAt.Before(started) {
		return fmt.Errorf("check finish precedes start")
	}
	current, err := checkHoldCurrent(tx, start, event.OccurredAt)
	if err != nil {
		return err
	}
	reasons := checkrecord.Qualification(start.Tree, result.Tree, result.Outcome, current, len(result.RefusedArtifacts) == 0)
	if !slices.Equal(reasons, result.Reasons) || result.Counts != (result.Outcome == "passed" && len(reasons) == 0) {
		return fmt.Errorf("check qualification differs from observations")
	}
	return validateCheckResultFiles(tx, start, result)
}

func validateCheckStopped(tx *sql.Tx, event Event) error {
	if len(event.Effects) != 2 || event.Effects[0].PK["id"] != event.Aggregate.ID {
		return fmt.Errorf("invalid stopped check request")
	}
	if err := validateCurrentIdempotency(event); err != nil {
		return err
	}
	var actor string
	if err := tx.QueryRow(`SELECT actor FROM check_requests WHERE id=?`, event.Aggregate.ID).Scan(&actor); err != nil {
		return err
	}
	if actor != event.ActorV2.ActorType+":"+event.ActorV2.ActorID {
		return fmt.Errorf("stopped check actor differs from request")
	}
	var starts, results int
	if err := tx.QueryRow(`SELECT count(*),count(r.check_id) FROM check_starts s LEFT JOIN check_results r ON r.check_id=s.id WHERE s.request_id=?`, event.Aggregate.ID).Scan(&starts, &results); err != nil {
		return err
	}
	if starts == 0 || starts != results {
		return fmt.Errorf("a request with an unfinished start cannot invent a terminal response")
	}
	return nil
}
