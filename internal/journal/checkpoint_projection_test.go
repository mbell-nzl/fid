package journal

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCheckpointedEventProjectsIdenticallyLiveAndReplay(t *testing.T) {
	event := checkpointedProjectionEvent()

	live := openReplayDB(t, filepath.Join(t.TempDir(), "live.sqlite"))
	defer closeProjectDB(t, live)
	seedCheckpointProjectionRun(t, live)
	hash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}
	if err := ProjectEvent(live, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}
	liveState := checkpointProjectionState(t, live)

	workspace := newReplayWorkspace(t)
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
	replayed := openReplayDB(t, filepath.Join(t.TempDir(), "replay.sqlite"))
	defer closeProjectDB(t, replayed)
	seedCheckpointProjectionRun(t, replayed)
	if _, err := ReplayJournal(replayed, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if replayState := checkpointProjectionState(t, replayed); replayState != liveState {
		t.Fatalf("checkpoint projection replay state = %q, want live state %q", replayState, liveState)
	}
}

func TestRecoveredEventProjectsIdenticallyLiveAndReplay(t *testing.T) {
	event := checkpointedProjectionEvent()
	event.EventID = "EVT-20260906-000003"
	event.EventType = "run.recovered"
	event.CommandName = "run.recover"
	event.Effects[0].Values["created_event_id"] = event.EventID

	live := openReplayDB(t, filepath.Join(t.TempDir(), "live-recovered.sqlite"))
	defer closeProjectDB(t, live)
	seedCheckpointProjectionRun(t, live)
	hash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}
	if err := ProjectEvent(live, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}
	liveState := checkpointProjectionState(t, live)

	workspace := newReplayWorkspace(t)
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
	replayed := openReplayDB(t, filepath.Join(t.TempDir(), "replay-recovered.sqlite"))
	defer closeProjectDB(t, replayed)
	seedCheckpointProjectionRun(t, replayed)
	if _, err := ReplayJournal(replayed, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if replayState := checkpointProjectionState(t, replayed); replayState != liveState {
		t.Fatalf("recovered projection replay state = %q, want live state %q", replayState, liveState)
	}
}

func TestBudgetObservedEventProjectsIdenticallyLiveAndReplay(t *testing.T) {
	event := checkpointedProjectionEvent()
	event.EventID = "EVT-20260906-000004"
	event.EventType = "run.budget-observed"
	event.CommandName = "run.observe-budget"
	event.Effects[0].Values["created_event_id"] = event.EventID

	live := openReplayDB(t, filepath.Join(t.TempDir(), "live-budget.sqlite"))
	defer closeProjectDB(t, live)
	seedCheckpointProjectionRun(t, live)
	hash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}
	if err := ProjectEvent(live, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}
	liveState := checkpointProjectionState(t, live)

	workspace := newReplayWorkspace(t)
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
	replayed := openReplayDB(t, filepath.Join(t.TempDir(), "replay-budget.sqlite"))
	defer closeProjectDB(t, replayed)
	seedCheckpointProjectionRun(t, replayed)
	if _, err := ReplayJournal(replayed, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if replayState := checkpointProjectionState(t, replayed); replayState != liveState {
		t.Fatalf("budget projection replay state = %q, want live state %q", replayState, liveState)
	}
}

func TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay(t *testing.T) {
	event := checkpointedProjectionEvent()
	event.EventID = "EVT-20260906-000005"
	event.EventType = "run.execution-blocked"
	event.CommandName = "gate.run"
	event.Effects[0].Values["created_event_id"] = event.EventID

	live := openReplayDB(t, filepath.Join(t.TempDir(), "live-blocked.sqlite"))
	defer closeProjectDB(t, live)
	seedCheckpointProjectionRun(t, live)
	hash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}
	if err := ProjectEvent(live, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}
	liveState := checkpointProjectionState(t, live)

	workspace := newReplayWorkspace(t)
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
	replayed := openReplayDB(t, filepath.Join(t.TempDir(), "replay-blocked.sqlite"))
	defer closeProjectDB(t, replayed)
	seedCheckpointProjectionRun(t, replayed)
	if _, err := ReplayJournal(replayed, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if replayState := checkpointProjectionState(t, replayed); replayState != liveState {
		t.Fatalf("execution-blocked projection replay state = %q, want live state %q", replayState, liveState)
	}
}

func TestRunClosureEventsProjectIdenticallyLiveAndReplay(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		status    string
	}{
		{name: "closed", eventType: "run.closed", status: "succeeded"},
		{name: "abandoned", eventType: "run.abandoned", status: "abandoned"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			event := runClosureProjectionEvent(testCase.eventType, testCase.status)
			live := openReplayDB(t, filepath.Join(t.TempDir(), "live-closure.sqlite"))
			defer closeProjectDB(t, live)
			seedRunClosureProjection(t, live)
			hash, err := HashEvent(event)
			if err != nil {
				t.Fatalf("HashEvent() error = %v", err)
			}
			if err := ProjectEvent(live, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1}); err != nil {
				t.Fatalf("ProjectEvent() error = %v", err)
			}
			liveState := runClosureProjectionState(t, live)

			workspace := newReplayWorkspace(t)
			writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
			replayed := openReplayDB(t, filepath.Join(t.TempDir(), "replay-closure.sqlite"))
			defer closeProjectDB(t, replayed)
			seedRunClosureProjection(t, replayed)
			if _, err := ReplayJournal(replayed, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
				t.Fatalf("ReplayJournal() error = %v", err)
			}
			if replayState := runClosureProjectionState(t, replayed); replayState != liveState {
				t.Fatalf("closure replay state = %q, want live state %q", replayState, liveState)
			}
		})
	}
}

func TestRemainingRecoveryEventsRemainFenced(t *testing.T) {
	for _, eventType := range []string{
		"activation.resource-observed",
		"activation.recovery-selected",
		"activation.recovery-completed",
	} {
		t.Run(eventType, func(t *testing.T) {
			event := checkpointedProjectionEvent()
			event.EventType = eventType
			err := DryRunProjectEvent(openReplayDB(t, filepath.Join(t.TempDir(), "fence.sqlite")), event)
			if err == nil || !strings.Contains(err.Error(), fmt.Sprintf(`Event-v2 %q does not permit effect application before its recovery domain is available`, eventType)) {
				t.Fatalf("DryRunProjectEvent() error = %v, want %q fence", err, eventType)
			}
		})
	}
}

func runClosureProjectionEvent(eventType, status string) Event {
	event := Event{
		EventID:     "EVT-20260906-000006",
		EventSchema: EventSchemaV2,
		EventType:   eventType,
		OccurredAt:  time.Date(2026, time.September, 6, 0, 0, 0, 0, time.UTC),
		ActorV2:     EventActor{ActorID: "agent:closure", ActorType: "agent", Producer: "fid", ProvenanceRef: "test:closure"},
		Aggregate:   Aggregate{ID: "RUN-closure", Type: "run", Sequence: 1},
		CommandName: strings.TrimSuffix(eventType, "d"),
		Effects: []Effect{
			{Op: "update", Table: "runs", PK: map[string]any{"id": "RUN-closure"}, Values: map[string]any{"status": status, "terminal_repository_json": `{}`, "ended_at": "2026-09-06T00:00:00Z"}, Preconditions: map[string]any{"status": "active"}},
			{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-closure"}, Values: map[string]any{"version": 2}, Preconditions: map[string]any{"version": 1}},
		},
	}
	if eventType == "run.abandoned" {
		event.CommandName = "run.abandon"
		event.Effects = append(event.Effects, Effect{Op: "update", Table: "claims", PK: map[string]any{"id": "CLM-closure"}, Values: map[string]any{"version": 2, "status": "released", "released_at": "2026-09-06T00:00:00Z", "release_reason": "abandoned"}, Preconditions: map[string]any{"status": "active", "version": 1}})
	}
	return event
}

func seedRunClosureProjection(t *testing.T, conn *sql.DB) {
	t.Helper()
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable foreign keys for closure projection seed: %v", err)
	}
	statements := []string{
		`INSERT INTO claims (id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at) VALUES ('CLM-closure', 'ISS-closure', 'ENV-closure', 1, 'agent:closure', 'WT-closure', '{}', 'mutation', 1, 'active', '2026-09-06T00:00:00Z', '2026-09-07T00:00:00Z')`,
		`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, status, started_at) VALUES ('RUN-closure', 'ISS-closure', 'ENV-closure', 1, 'CLM-closure', 1, 'agent:closure', '{}', 'active', '2026-09-06T00:00:00Z')`,
		`INSERT INTO run_runtime (run_id, version, executor_id, obligation_id) VALUES ('RUN-closure', 1, 'executor-closure', NULL)`,
	}
	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("seed closure projection: %v", err)
		}
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatalf("enable foreign keys after closure projection seed: %v", err)
	}
}

func runClosureProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()
	var status, terminal, ended, claimStatus, released, reason string
	var version, claimVersion int
	if err := conn.QueryRow(`SELECT status, terminal_repository_json, ended_at FROM runs WHERE id = 'RUN-closure'`).Scan(&status, &terminal, &ended); err != nil {
		t.Fatalf("load closure Run: %v", err)
	}
	if err := conn.QueryRow(`SELECT version FROM run_runtime WHERE run_id = 'RUN-closure'`).Scan(&version); err != nil {
		t.Fatalf("load closure runtime: %v", err)
	}
	if err := conn.QueryRow(`SELECT status, version, COALESCE(released_at, ''), COALESCE(release_reason, '') FROM claims WHERE id = 'CLM-closure'`).Scan(&claimStatus, &claimVersion, &released, &reason); err != nil {
		t.Fatalf("load closure Claim: %v", err)
	}
	return fmt.Sprintf("%s|%s|%s|%d|%s|%d|%s|%s", status, terminal, ended, version, claimStatus, claimVersion, released, reason)
}

func checkpointedProjectionEvent() Event {
	return Event{
		EventID:     "EVT-20260906-000001",
		EventSchema: EventSchemaV2,
		EventType:   "run.checkpointed",
		OccurredAt:  time.Date(2026, time.September, 6, 0, 0, 0, 0, time.UTC),
		ActorV2:     EventActor{ActorID: "agent:checkpoint", ActorType: "agent", Producer: "activated-harness", ProvenanceRef: "test:checkpoint"},
		Aggregate:   Aggregate{ID: "RUN-checkpoint", Type: "run", Sequence: 1},
		CommandName: "run.checkpoint",
		Effects: []Effect{
			{Op: "insert", Table: "run_observations", PK: map[string]any{"id": "OBS-checkpoint"}, Values: map[string]any{"run_id": "RUN-checkpoint", "kind": "checkpoint", "actor_id": "agent:checkpoint", "observed_at": "2026-09-06T00:00:00Z", "observation_json": `{"id":"OBS-checkpoint","kind":"checkpoint","run_id":"RUN-checkpoint"}`, "created_event_id": "EVT-20260906-000001"}},
			{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-checkpoint"}, Values: map[string]any{"version": 2}, Preconditions: map[string]any{"version": 1}},
		},
	}
}

func seedCheckpointProjectionRun(t *testing.T, conn *sql.DB) {
	t.Helper()
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable foreign keys for projection seed: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, status, started_at) VALUES ('RUN-checkpoint', 'ISS-checkpoint', 'ENV-checkpoint', 1, 'CLM-checkpoint', 1, 'agent:checkpoint', '{}', 'active', '2026-09-06T00:00:00Z')`); err != nil {
		t.Fatalf("seed run: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO run_runtime (run_id, version, executor_id, obligation_id) VALUES ('RUN-checkpoint', 1, 'executor-checkpoint', NULL)`); err != nil {
		t.Fatalf("seed runtime: %v", err)
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatalf("enable foreign keys after projection seed: %v", err)
	}
}

func checkpointProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()
	var id, runID, kind, actorID, observedAt, observation, eventID string
	if err := conn.QueryRow(`SELECT id, run_id, kind, actor_id, observed_at, observation_json, created_event_id FROM run_observations WHERE id = 'OBS-checkpoint'`).Scan(&id, &runID, &kind, &actorID, &observedAt, &observation, &eventID); err != nil {
		t.Fatalf("load checkpoint observation: %v", err)
	}
	var version int
	if err := conn.QueryRow(`SELECT version FROM run_runtime WHERE run_id = 'RUN-checkpoint'`).Scan(&version); err != nil {
		t.Fatalf("load checkpoint runtime: %v", err)
	}
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%d", id, runID, kind, actorID, observedAt, observation, eventID, version)
}
