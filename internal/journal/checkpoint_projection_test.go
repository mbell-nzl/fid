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

func TestActivationRecoveryEventsProjectIdenticallyLiveAndReplay(t *testing.T) {
	for _, eventType := range []string{
		"activation.resource-observed",
		"activation.recovery-selected",
		"activation.recovery-completed",
	} {
		t.Run(eventType, func(t *testing.T) {
			event := activationRecoveryProjectionEvent(eventType)
			live := openReplayDB(t, filepath.Join(t.TempDir(), "live-activation.sqlite"))
			defer closeProjectDB(t, live)
			seedActivationRecoveryProjection(t, live, eventType)
			hash, err := HashEvent(event)
			if err != nil {
				t.Fatalf("HashEvent() error = %v", err)
			}
			if err := ProjectEvent(live, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1}); err != nil {
				t.Fatalf("ProjectEvent() error = %v", err)
			}
			liveState := activationRecoveryProjectionState(t, live)
			workspace := newReplayWorkspace(t)
			writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
			replayed := openReplayDB(t, filepath.Join(t.TempDir(), "replay-activation.sqlite"))
			defer closeProjectDB(t, replayed)
			seedActivationRecoveryProjection(t, replayed, eventType)
			if _, err := ReplayJournal(replayed, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
				t.Fatalf("ReplayJournal() error = %v", err)
			}
			if replayState := activationRecoveryProjectionState(t, replayed); replayState != liveState {
				t.Fatalf("activation replay state = %q, want live state %q", replayState, liveState)
			}
		})
	}
}

func activationRecoveryProjectionEvent(eventType string) Event {
	event := Event{EventID: "EVT-20260906-000099", EventSchema: EventSchemaV2, EventType: eventType, OccurredAt: time.Date(2026, time.September, 6, 0, 0, 0, 0, time.UTC), ActorV2: EventActor{ActorID: "fid-bootstrap", ActorType: "service", Producer: "fid", ProvenanceRef: "test:activation"}, Aggregate: Aggregate{ID: "ACT-projection", Type: "activation", Sequence: 1}, CommandName: "activation.resource"}
	switch eventType {
	case "activation.resource-observed":
		event.Effects = []Effect{{Op: "update", Table: "kit_activations", PK: map[string]any{"id": "ACT-projection"}, Preconditions: map[string]any{"version": 1}, Values: map[string]any{"version": 2}}, {Op: "update", Table: "activation_resources", PK: map[string]any{"activation_id": "ACT-projection", "component": "kit"}, Preconditions: map[string]any{"version": 1}, Values: map[string]any{"version": 2, "observed_identity": "kit:new", "state": "new", "evidence_ref": "test:kit", "error_code": nil}}}
	case "activation.recovery-selected":
		event.CommandName = "activation.recover"
		event.Effects = []Effect{{Op: "update", Table: "kit_activations", PK: map[string]any{"id": "ACT-projection"}, Preconditions: map[string]any{"version": 1}, Values: map[string]any{"version": 2}}, {Op: "insert", Table: "activation_recoveries", PK: map[string]any{"activation_id": "ACT-projection"}, Values: map[string]any{"direction": "restore-prior", "actor_attribution": "fid-bootstrap", "reason": "test", "selected_event_id": "EVT-20260906-000099", "completed_event_id": nil}}}
	case "activation.recovery-completed":
		event.CommandName = "activation.recover-complete"
		event.Effects = []Effect{{Op: "update", Table: "kit_activations", PK: map[string]any{"id": "ACT-projection"}, Preconditions: map[string]any{"version": 1}, Values: map[string]any{"version": 2, "status": "aborted"}}, {Op: "update", Table: "activation_recoveries", PK: map[string]any{"activation_id": "ACT-projection"}, Values: map[string]any{"completed_event_id": "EVT-20260906-000099"}}}
	}
	return event
}

func seedActivationRecoveryProjection(t *testing.T, conn *sql.DB, eventType string) {
	t.Helper()
	if _, err := conn.Exec(`INSERT INTO kit_activations (id, kit_identity, backend_identity, phase, status, point_of_no_return, direction, created_at, updated_at, version) VALUES ('ACT-projection','kit','fid',21,'prepared',0,'activate','2026-09-06T00:00:00Z','2026-09-06T00:00:00Z',1)`); err != nil {
		t.Fatalf("seed activation: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO activation_resources (activation_id, component, version, old_identity, new_identity, state) VALUES ('ACT-projection','kit',1,NULL,'kit:new','pending')`); err != nil {
		t.Fatalf("seed activation resource: %v", err)
	}
	if eventType == "activation.recovery-completed" {
		if _, err := conn.Exec(`INSERT INTO activation_recoveries (activation_id, direction, actor_attribution, reason, selected_event_id, completed_event_id) VALUES ('ACT-projection','restore-prior','fid-bootstrap','test','EVT-selected',NULL)`); err != nil {
			t.Fatalf("seed activation recovery: %v", err)
		}
	}
}

func activationRecoveryProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()
	var state string
	if err := conn.QueryRow(`SELECT status || ':' || version || ':' || COALESCE((SELECT state FROM activation_resources WHERE activation_id = 'ACT-projection' AND component = 'kit'),'') || ':' || COALESCE((SELECT direction || ':' || COALESCE(completed_event_id,'') FROM activation_recoveries WHERE activation_id = 'ACT-projection'),'') FROM kit_activations WHERE id = 'ACT-projection'`).Scan(&state); err != nil {
		t.Fatalf("read activation projection: %v", err)
	}
	return state
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
