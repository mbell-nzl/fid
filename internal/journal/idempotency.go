package journal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type IdempotencyRequest struct {
	OperationScope      string
	IdempotencyKey      string
	RequestDigest       string
	RequestFieldDigests map[string]string
}

type IdempotencyMutation struct {
	EntityID   string
	ResultJSON string
}

type IdempotencyOutcome struct {
	EntityID          string
	EventID           string
	AggregateSequence int
	ResultJSON        string
	Replayed          bool
}

type IdempotencyConflictError struct {
	OperationScope        string
	IdempotencyKey        string
	StoredRequestDigest   string
	SuppliedRequestDigest string
	StoredFieldDigests    map[string]string
	SuppliedFieldDigests  map[string]string
}

// CommitEvent appends and projects one non-idempotent Event-v2 mutation.
func CommitEvent(lock *WriteLock, conn *sql.DB, event Event) (Event, error) {
	if lock == nil {
		return Event{}, fmt.Errorf("commit Event: lock is nil")
	}
	if conn == nil {
		return Event{}, fmt.Errorf("commit Event: connection is nil")
	}
	if !event.IsV2() {
		return Event{}, fmt.Errorf("commit Event: Event-v2 is required")
	}
	prepared, err := PrepareEvent(lock, event)
	if err != nil {
		return Event{}, err
	}
	if err := captureCurrentEventPrior(conn, &prepared); err != nil {
		return Event{}, err
	}
	if err := DryRunProjectEvent(conn, prepared); err != nil {
		return Event{}, err
	}
	appended, err := AppendEvent(lock, prepared)
	if err != nil {
		return Event{}, err
	}
	eventHash, err := HashEvent(appended)
	if err != nil {
		return Event{}, fmt.Errorf("commit Event: hash Event %q: %w", appended.EventID, err)
	}
	journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(lock.workspaceRoot, ".fid", "journal", journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return Event{}, fmt.Errorf("commit Event: stat journal %q: %w", journalPath, err)
	}
	if err := ProjectEvent(conn, ProjectionInput{Event: appended, EventHash: eventHash, File: journalFile, Offset: info.Size()}); err != nil {
		return Event{}, err
	}
	return appended, nil
}

func (err *IdempotencyConflictError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("idempotency conflict for %s key %q", err.OperationScope, err.IdempotencyKey)
}

func CommitIdempotent(lock *WriteLock, conn *sql.DB, request IdempotencyRequest, event Event, build func(Event) (IdempotencyMutation, error)) (IdempotencyOutcome, error) {
	if lock == nil {
		return IdempotencyOutcome{}, fmt.Errorf("commit idempotent: lock is nil")
	}
	if conn == nil {
		return IdempotencyOutcome{}, fmt.Errorf("commit idempotent: connection is nil")
	}
	if build == nil {
		return IdempotencyOutcome{}, fmt.Errorf("commit idempotent: result builder is nil")
	}
	if err := validateIdempotencyRequest(request); err != nil {
		return IdempotencyOutcome{}, err
	}

	if outcome, found, err := ReplayIdempotent(conn, request); err != nil {
		return IdempotencyOutcome{}, err
	} else if found {
		return outcome, nil
	}

	if !event.IsV2() {
		return IdempotencyOutcome{}, fmt.Errorf("commit idempotent: Event-v2 is required")
	}
	event.Idempotency = &EventIdempotency{request.IdempotencyKey, request.RequestDigest, request.OperationScope}
	prepared, err := PrepareEvent(lock, event)
	if len(event.Effects) == 0 {
		prepared, err = prepareIdempotentEvent(lock, event)
	}
	if err != nil {
		return IdempotencyOutcome{}, err
	}
	mutation, err := build(prepared)
	if err != nil {
		return IdempotencyOutcome{}, err
	}
	if strings.TrimSpace(mutation.EntityID) == "" {
		return IdempotencyOutcome{}, fmt.Errorf("commit idempotent: result entity ID is required")
	}
	resultJSON, err := canonicalNonReplayedResultJSON(mutation.ResultJSON)
	if err != nil {
		return IdempotencyOutcome{}, err
	}

	prepared.Effects = append(prepared.Effects, idempotencyResultEffect(request, prepared, mutation.EntityID, resultJSON))
	if err := captureCurrentEventPrior(conn, &prepared); err != nil {
		return IdempotencyOutcome{}, err
	}
	if err := validateV2Event(prepared); err != nil {
		return IdempotencyOutcome{}, err
	}
	if err := DryRunProjectEvent(conn, prepared); err != nil {
		return IdempotencyOutcome{}, err
	}
	appended, err := AppendEvent(lock, prepared)
	if err != nil {
		return IdempotencyOutcome{}, err
	}
	eventHash, err := HashEvent(appended)
	if err != nil {
		return IdempotencyOutcome{}, fmt.Errorf("commit idempotent: hash Event %q: %w", appended.EventID, err)
	}
	journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(lock.workspaceRoot, ".fid", "journal", journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return IdempotencyOutcome{}, fmt.Errorf("commit idempotent: stat journal %q: %w", journalPath, err)
	}
	if err := ProjectEvent(conn, ProjectionInput{Event: appended, EventHash: eventHash, File: journalFile, Offset: info.Size()}); err != nil {
		return IdempotencyOutcome{}, err
	}
	return IdempotencyOutcome{mutation.EntityID, appended.EventID, appended.Aggregate.Sequence, resultJSON, false}, nil
}

// ReplayIdempotent returns a previously committed result for the exact request.
func ReplayIdempotent(conn *sql.DB, request IdempotencyRequest) (IdempotencyOutcome, bool, error) {
	if conn == nil {
		return IdempotencyOutcome{}, false, fmt.Errorf("replay idempotent: connection is nil")
	}
	if err := validateIdempotencyRequest(request); err != nil {
		return IdempotencyOutcome{}, false, err
	}

	stored, found, err := findIdempotencyResult(conn, request)
	if err != nil {
		return IdempotencyOutcome{}, false, err
	}
	if !found {
		return IdempotencyOutcome{}, false, nil
	}
	if stored.requestDigest != request.RequestDigest {
		return IdempotencyOutcome{}, false, &IdempotencyConflictError{OperationScope: request.OperationScope, IdempotencyKey: request.IdempotencyKey, StoredRequestDigest: stored.requestDigest, SuppliedRequestDigest: request.RequestDigest, StoredFieldDigests: stored.fieldDigests, SuppliedFieldDigests: request.RequestFieldDigests}
	}
	replayedJSON, err := replayedResultJSON(stored.resultJSON)
	if err != nil {
		return IdempotencyOutcome{}, false, err
	}
	return IdempotencyOutcome{stored.entityID, stored.eventID, stored.aggregateSequence, replayedJSON, true}, true, nil
}

type storedIdempotencyResult struct {
	requestDigest     string
	fieldDigests      map[string]string
	entityID          string
	eventID           string
	aggregateSequence int
	resultJSON        string
}

func findIdempotencyResult(conn *sql.DB, request IdempotencyRequest) (storedIdempotencyResult, bool, error) {
	var result storedIdempotencyResult
	var fieldDigestsJSON string
	err := conn.QueryRow(`
SELECT request_digest, request_field_digests_json, result_entity_id, result_event_id, aggregate_sequence, result_json
FROM idempotency_results
WHERE operation_scope = ? AND idempotency_key = ?`, request.OperationScope, request.IdempotencyKey).Scan(
		&result.requestDigest,
		&fieldDigestsJSON,
		&result.entityID,
		&result.eventID,
		&result.aggregateSequence,
		&result.resultJSON,
	)
	if err == sql.ErrNoRows {
		return storedIdempotencyResult{}, false, nil
	}
	if err != nil {
		return storedIdempotencyResult{}, false, fmt.Errorf("commit idempotent: read stored result: %w", err)
	}
	if err := json.Unmarshal([]byte(fieldDigestsJSON), &result.fieldDigests); err != nil {
		return storedIdempotencyResult{}, false, fmt.Errorf("commit idempotent: decode stored field digests: %w", err)
	}
	return result, true, nil
}

func idempotencyResultEffect(request IdempotencyRequest, event Event, entityID, resultJSON string) Effect {
	fields := request.RequestFieldDigests
	if fields == nil {
		fields = map[string]string{}
	}
	fieldJSON, _ := MarshalCanonicalJSON(fields)
	return Effect{
		Op:    "insert",
		Table: "idempotency_results",
		PK: map[string]any{
			"operation_scope": request.OperationScope,
			"idempotency_key": request.IdempotencyKey,
		},
		Values: map[string]any{
			"request_digest":             request.RequestDigest,
			"request_field_digests_json": string(fieldJSON),
			"result_entity_id":           entityID,
			"result_event_id":            event.EventID,
			"aggregate_sequence":         event.Aggregate.Sequence,
			"result_json":                resultJSON,
			"created_at":                 event.OccurredAt.UTC().Format(time.RFC3339Nano),
		},
	}
}

func validateIdempotencyRequest(request IdempotencyRequest) error {
	if strings.TrimSpace(request.OperationScope) == "" {
		return fmt.Errorf("commit idempotent: operation scope is required")
	}
	if len(request.IdempotencyKey) < 1 || len(request.IdempotencyKey) > 128 {
		return fmt.Errorf("commit idempotent: idempotency key must contain 1..128 printable ASCII characters")
	}
	for _, character := range request.IdempotencyKey {
		if character < 0x20 || character > 0x7e {
			return fmt.Errorf("commit idempotent: idempotency key must contain 1..128 printable ASCII characters")
		}
	}
	if strings.TrimSpace(request.RequestDigest) == "" {
		return fmt.Errorf("commit idempotent: request digest is required")
	}
	return nil
}

func canonicalNonReplayedResultJSON(raw string) (string, error) {
	return resultJSONWithReplay(raw, false)
}

func replayedResultJSON(raw string) (string, error) {
	return resultJSONWithReplay(raw, true)
}

func resultJSONWithReplay(raw string, replayed bool) (string, error) {
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	var value map[string]any
	if err := decoder.Decode(&value); err != nil {
		return "", fmt.Errorf("idempotency result JSON: decode: %w", err)
	}
	if value == nil {
		return "", fmt.Errorf("idempotency result JSON: object is required")
	}
	if _, exists := value["replayed"]; !exists {
		return "", fmt.Errorf("idempotency result JSON: replayed is required")
	}
	value["replayed"] = replayed
	canonical, err := MarshalCanonicalJSON(value)
	if err != nil {
		return "", fmt.Errorf("idempotency result JSON: canonicalize: %w", err)
	}
	return string(canonical), nil
}
