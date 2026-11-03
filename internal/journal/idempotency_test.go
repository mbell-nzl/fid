package journal

import (
	"bytes"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestCommitIdempotentReplaysStoredResultAfterRebuild(t *testing.T) {
	t.Parallel()
	workspace, conn := newIdempotencyWorkspace(t)
	request := testIdempotencyRequest()

	first, err := commitTestIdempotentActivation(workspace, conn, request)
	if err != nil {
		t.Fatalf("first CommitIdempotent() error = %v", err)
	}
	var liveFields string
	if err := conn.QueryRow(`SELECT request_field_digests_json FROM idempotency_results WHERE operation_scope=? AND idempotency_key=?`, request.OperationScope, request.IdempotencyKey).Scan(&liveFields); err != nil {
		t.Fatal(err)
	}
	if liveFields != "{}" {
		t.Fatalf("live field digests = %q, want {}", liveFields)
	}
	if first.Replayed {
		t.Fatal("first idempotency outcome replayed = true, want false")
	}
	if got, want := first.ResultJSON, testActivationResultJSON("key-1", false); got != want {
		t.Fatalf("first result JSON = %s, want %s", got, want)
	}

	closeIdempotencyDB(t, conn)
	if err := os.Remove(filepath.Join(config.LedgerDir(workspace), "fid.sqlite")); err != nil {
		t.Fatalf("Remove(projection) error = %v", err)
	}
	if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}
	conn = openIdempotencyDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIdempotencyDB(t, conn)
	var rebuiltFields string
	if err := conn.QueryRow(`SELECT request_field_digests_json FROM idempotency_results WHERE operation_scope=? AND idempotency_key=?`, request.OperationScope, request.IdempotencyKey).Scan(&rebuiltFields); err != nil {
		t.Fatal(err)
	}
	if rebuiltFields != liveFields {
		t.Fatalf("rebuilt field digests = %q, want %q", rebuiltFields, liveFields)
	}

	replayed, err := commitTestIdempotentActivation(workspace, conn, request)
	if err != nil {
		t.Fatalf("replay CommitIdempotent() error = %v", err)
	}
	if !replayed.Replayed {
		t.Fatal("replay outcome replayed = false, want true")
	}
	if replayed.EventID != first.EventID || replayed.AggregateSequence != first.AggregateSequence {
		t.Fatalf("replay identity = (%s, %d), want (%s, %d)", replayed.EventID, replayed.AggregateSequence, first.EventID, first.AggregateSequence)
	}
	if got, want := replayed.ResultJSON, testActivationResultJSON("key-1", true); got != want {
		t.Fatalf("replay result JSON = %s, want %s", got, want)
	}
	if got := idempotencyJournalEventCount(t, workspace); got != 1 {
		t.Fatalf("journal Event count after replay = %d, want 1", got)
	}
}

func TestCommitIdempotentRejectsSameKeyWithDifferentDigest(t *testing.T) {
	t.Parallel()
	workspace, conn := newIdempotencyWorkspace(t)
	defer closeIdempotencyDB(t, conn)

	if _, err := commitTestIdempotentActivation(workspace, conn, testIdempotencyRequest()); err != nil {
		t.Fatalf("first CommitIdempotent() error = %v", err)
	}

	conflicting := testIdempotencyRequest()
	conflicting.RequestDigest = "sha256:2222222222222222222222222222222222222222222222222222222222222222"
	_, err := commitTestIdempotentActivation(workspace, conn, conflicting)
	var conflict *IdempotencyConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("conflicting CommitIdempotent() error = %T %v, want *IdempotencyConflictError", err, err)
	}
	if conflict.OperationScope != conflicting.OperationScope || conflict.IdempotencyKey != conflicting.IdempotencyKey ||
		conflict.StoredRequestDigest != testIdempotencyRequest().RequestDigest || conflict.SuppliedRequestDigest != conflicting.RequestDigest {
		t.Fatalf("conflict = %+v, want stored and supplied digest details", conflict)
	}
	if got := idempotencyJournalEventCount(t, workspace); got != 1 {
		t.Fatalf("journal Event count after conflict = %d, want 1", got)
	}
}

func TestCommitIdempotentSerializesConcurrentSameKeyCalls(t *testing.T) {
	t.Parallel()
	workspace, conn := newIdempotencyWorkspace(t)
	closeIdempotencyDB(t, conn)

	start := make(chan struct{})
	results := make(chan IdempotencyOutcome, 2)
	errs := make(chan error, 2)
	var wait sync.WaitGroup
	for range 2 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			connection := openIdempotencyDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
			defer closeIdempotencyDB(t, connection)
			result, err := commitTestIdempotentActivation(workspace, connection, testIdempotencyRequest())
			if err == nil {
				results <- result
			}
			errs <- err
		}()
	}

	close(start)
	wait.Wait()
	close(results)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent CommitIdempotent() error = %v", err)
		}
	}

	var outcomes []IdempotencyOutcome
	for result := range results {
		outcomes = append(outcomes, result)
	}
	if len(outcomes) != 2 {
		t.Fatalf("concurrent outcomes = %d, want 2", len(outcomes))
	}
	if outcomes[0].EventID != outcomes[1].EventID || outcomes[0].AggregateSequence != outcomes[1].AggregateSequence {
		t.Fatalf("concurrent result identities = %+v, want one committed Event", outcomes)
	}
	if outcomes[0].Replayed == outcomes[1].Replayed {
		t.Fatalf("concurrent replay states = (%t, %t), want one first result and one replay", outcomes[0].Replayed, outcomes[1].Replayed)
	}
	if got := idempotencyJournalEventCount(t, workspace); got != 1 {
		t.Fatalf("journal Event count after concurrent calls = %d, want 1", got)
	}
}

func commitTestIdempotentActivation(workspace string, conn *sql.DB, request IdempotencyRequest) (IdempotencyOutcome, error) {
	var outcome IdempotencyOutcome
	err := WithLock(workspace, LockOwner{Command: "test idempotent activation"}, func(lock *WriteLock) error {
		var commitErr error
		outcome, commitErr = CommitIdempotent(lock, conn, request, testIdempotentActivationEvent(), func(Event) (IdempotencyMutation, error) {
			return IdempotencyMutation{
				EntityID:   "ACT-1",
				ResultJSON: testActivationResultJSON("key-1", false),
			}, nil
		})
		return commitErr
	})
	return outcome, err
}

func testIdempotencyRequest() IdempotencyRequest {
	return IdempotencyRequest{
		OperationScope: "activation.prepare:project-1",
		IdempotencyKey: "key-1",
		RequestDigest:  "sha256:1111111111111111111111111111111111111111111111111111111111111111",
	}
}

func testIdempotentActivationEvent() Event {
	event := Event{
		EventSchema: EventSchemaV2,
		EventType:   "activation.prepared",
		ActorV2: EventActor{
			ActorID:       "actor-builder",
			ActorType:     "agent",
			Producer:      "kit",
			ProvenanceRef: "activation:ACT-1/profile:builder",
		},
		Aggregate:   Aggregate{ID: "ACT-1", Sequence: 1, Type: "activation"},
		CommandName: "activation.prepare",
		Effects: []Effect{{
			Op:    "insert",
			Table: "kit_activations",
			PK:    map[string]any{"id": "ACT-1"},
			Values: map[string]any{
				"version":             1,
				"prior_activation_id": nil,
				"kit_identity":        "kit-1",
				"backend_identity":    "fid-1",
				"phase":               1,
				"status":              "prepared",
				"point_of_no_return":  false,
				"direction":           "activate",
				"created_at":          "2026-08-22T00:00:00Z",
				"updated_at":          "2026-08-22T00:00:00Z",
			},
		}},
	}
	for _, component := range []string{"backend", "kit", "overlay", "project", "provider"} {
		event.Effects = append(event.Effects, Effect{Op: "insert", Table: "activation_resources", PK: map[string]any{"activation_id": "ACT-1", "component": component}, Values: map[string]any{"version": 1, "old_identity": nil, "new_identity": component + ":new", "observed_identity": nil, "state": "pending", "evidence_ref": nil, "error_code": nil}})
	}
	return event
}

func newIdempotencyWorkspace(t *testing.T) (string, *sql.DB) {
	t.Helper()
	workspace := t.TempDir()
	if err := os.MkdirAll(filepath.Join(config.LedgerDir(workspace), "journal"), 0o755); err != nil {
		t.Fatalf("MkdirAll(journal) error = %v", err)
	}
	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}
	return workspace, openIdempotencyDB(t, dbPath)
}

func openIdempotencyDB(t *testing.T, path string) *sql.DB {
	t.Helper()
	conn, err := fiddb.Open(path)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", path, err)
	}
	return conn
}

func closeIdempotencyDB(t *testing.T, conn *sql.DB) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func idempotencyJournalEventCount(t *testing.T, workspace string) int {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(config.LedgerDir(workspace), "journal"))
	if err != nil {
		t.Fatalf("ReadDir(journal) error = %v", err)
	}
	count := 0
	for _, entry := range entries {
		payload, err := os.ReadFile(filepath.Join(config.LedgerDir(workspace), "journal", entry.Name()))
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", entry.Name(), err)
		}
		count += bytes.Count(payload, []byte{'\n'})
	}
	return count
}
