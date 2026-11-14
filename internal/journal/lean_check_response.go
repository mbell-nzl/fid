package journal

import (
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
)

func validateCheckResultFiles(tx *sql.Tx, start checkrecord.Start, result checkrecord.Result) error {
	rows, err := tx.Query(`SELECT f.path,f.artifact_id,f.reason,a.digest,a.size_bytes FROM check_files f LEFT JOIN referenced_artifacts a ON a.id=f.artifact_id WHERE f.check_id=? ORDER BY f.path`, start.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	kept, refused := []checkrecord.Artifact{}, []checkrecord.RefusedArtifact{}
	for rows.Next() {
		var path string
		var id, reason, digest sql.NullString
		var size sql.NullInt64
		if err := rows.Scan(&path, &id, &reason, &digest, &size); err != nil {
			return err
		}
		if id.Valid {
			kept = append(kept, checkrecord.Artifact{ID: id.String, Path: path, Digest: digest.String, SizeBytes: size.Int64})
		} else {
			refused = append(refused, checkrecord.RefusedArtifact{Path: path, Reason: reason.String})
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(kept)+len(refused) != len(start.Definition.Artifacts) || !reflect.DeepEqual(kept, result.Artifacts) || !reflect.DeepEqual(refused, result.RefusedArtifacts) {
		return fmt.Errorf("check result differs from its listed file observations")
	}
	return nil
}

func validateCheckResponse(tx *sql.Tx, event Event) error {
	var effect *Effect
	for i := range event.Effects {
		if event.Effects[i].Table == "idempotency_results" {
			effect = &event.Effects[i]
		}
	}
	if effect == nil {
		return nil
	}
	var response checkrecord.Response
	if err := decodeCheckJSON(effect.Values["result_json"], &response); err != nil {
		return err
	}
	var requestID string
	if err := tx.QueryRow(`SELECT request_id FROM check_starts WHERE id=?`, event.Aggregate.ID).Scan(&requestID); err != nil {
		return err
	}
	var task, hold, namesJSON string
	var final sql.NullString
	if err := tx.QueryRow(`SELECT task_id,hold_id,names_json,response_json FROM check_requests WHERE id=?`, requestID).Scan(&task, &hold, &namesJSON, &final); err != nil {
		return err
	}
	var names []string
	if err := decodeCheckJSON(namesJSON, &names); err != nil {
		return err
	}
	if response.TaskID != task || !slices.Equal(response.Names, names) || response.Replayed || response.CheckIDs == nil || response.Checks == nil || response.NotStarted == nil {
		return fmt.Errorf("check retry response differs from its request")
	}
	rows, err := tx.Query(`SELECT s.id,s.name,r.result_json FROM check_starts s LEFT JOIN check_results r ON r.check_id=s.id WHERE s.request_id=? ORDER BY s.rowid`, requestID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	ids := []string{}
	summaries := []checkrecord.Summary{}
	for rows.Next() {
		var id, name string
		var raw sql.NullString
		if err := rows.Scan(&id, &name, &raw); err != nil {
			return err
		}
		ids = append(ids, id)
		if raw.Valid {
			var result checkrecord.Result
			if err := decodeCheckJSON(raw.String, &result); err != nil {
				return err
			}
			artifacts := []string{}
			for _, file := range result.Artifacts {
				artifacts = append(artifacts, file.ID)
			}
			summaries = append(summaries, checkrecord.Summary{ID: id, Name: name, Outcome: result.Outcome, ExitCode: result.ExitCode, Counts: result.Counts, Reasons: result.Reasons, ArtifactIDs: artifacts})
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if !slices.Equal(ids, response.CheckIDs) || !reflect.DeepEqual(summaries, response.Checks) {
		return fmt.Errorf("check retry response differs from recorded results")
	}
	if response.Complete {
		if !final.Valid || !currentValuesEqualJSON(final.String, effect.Values["result_json"]) || response.EventID != event.EventID || response.OccurredAt != event.OccurredAt.UTC().Format(time.RFC3339Nano) || len(summaries) != len(ids) || !slices.Equal(response.NotStarted, names[len(ids):]) {
			return fmt.Errorf("invalid completed check response")
		}
		if len(response.NotStarted) == 0 {
			expected := checkrecord.FailureFor(response)
			if response.Failure == nil || response.Failure.Code != "hold-not-current" {
				if checkrecord.Digest(expected) != checkrecord.Digest(response.Failure) {
					return fmt.Errorf("check failure differs from its results")
				}
			}
		} else if event.EventType == "check.stopped" {
			if response.Failure == nil || !slices.Contains([]string{"hold-not-current", "check-still-running", "internal-error"}, response.Failure.Code) {
				return fmt.Errorf("invalid stopped request refusal")
			}
		} else if response.Failure == nil || response.Failure.Code != "hold-not-current" {
			if len(summaries) == 0 || summaries[len(summaries)-1].Outcome != "cancelled" || checkrecord.Digest(checkrecord.FailureFor(response)) != checkrecord.Digest(response.Failure) {
				return fmt.Errorf("checks stopped without an ownership refusal or cancellation")
			}
		}
	} else if final.Valid || response.EventID != "" || response.OccurredAt != "" || response.Failure != nil || len(response.NotStarted) != 0 {
		return fmt.Errorf("unfinished check request contains a terminal response")
	}
	return validateCheckOwnershipFailure(tx, response, hold, event.OccurredAt)
}

func currentValuesEqualJSON(a string, b any) bool {
	left, err := decodeCurrentJSON(a)
	if err != nil {
		return false
	}
	right, ok := b.(string)
	if !ok {
		return false
	}
	value, err := decodeCurrentJSON(right)
	return err == nil && currentValuesEqual(left, value)
}

// Keep the stopped request's ownership error closed and bound to its hold.
func validateCheckOwnershipFailure(tx *sql.Tx, response checkrecord.Response, hold string, at time.Time) error {
	if response.Failure == nil || response.Failure.Code != "hold-not-current" {
		return nil
	}
	failure := response.Failure
	var ended sql.NullString
	var expires string
	if err := tx.QueryRow(`SELECT ended_at,expires_at FROM holds WHERE id=?`, hold).Scan(&ended, &expires); err != nil {
		return err
	}
	reason := "expired"
	if ended.Valid {
		reason = "ended"
	}
	until, err := time.Parse(time.RFC3339, expires)
	if err != nil {
		return err
	}
	if !ended.Valid && at.Before(until) || failure.TaskID != response.TaskID || failure.Details["hold_id"] != hold || failure.Details["reason"] != reason || len(failure.Details) != 2 || !strings.Contains(failure.Message, "hold") {
		return fmt.Errorf("check ownership failure differs from its hold")
	}
	return nil
}
