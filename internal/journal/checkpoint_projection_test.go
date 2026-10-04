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
	t.Parallel()
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
	t.Parallel()
	event := recoveredProjectionEvent()

	live := openReplayDB(t, filepath.Join(t.TempDir(), "live-recovered.sqlite"))
	defer closeProjectDB(t, live)
	seedRecoveredProjection(t, live)
	hash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}
	if err := ProjectEvent(live, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}
	liveState := recoveredProjectionState(t, live)

	workspace := newReplayWorkspace(t)
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
	replayed := openReplayDB(t, filepath.Join(t.TempDir(), "replay-recovered.sqlite"))
	defer closeProjectDB(t, replayed)
	seedRecoveredProjection(t, replayed)
	if _, err := ReplayJournal(replayed, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if replayState := recoveredProjectionState(t, replayed); replayState != liveState {
		t.Fatalf("recovered projection replay state = %q, want live state %q", replayState, liveState)
	}
}

func TestBudgetObservedEventProjectsIdenticallyLiveAndReplay(t *testing.T) {
	t.Parallel()
	event := budgetObservedProjectionEvent()

	live := openReplayDB(t, filepath.Join(t.TempDir(), "live-budget.sqlite"))
	defer closeProjectDB(t, live)
	seedBudgetProjection(t, live)
	hash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}
	if err := ProjectEvent(live, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}
	liveState := budgetProjectionState(t, live)

	workspace := newReplayWorkspace(t)
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
	replayed := openReplayDB(t, filepath.Join(t.TempDir(), "replay-budget.sqlite"))
	defer closeProjectDB(t, replayed)
	seedBudgetProjection(t, replayed)
	if _, err := ReplayJournal(replayed, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if replayState := budgetProjectionState(t, replayed); replayState != liveState {
		t.Fatalf("budget projection replay state = %q, want live state %q", replayState, liveState)
	}
}

func TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay(t *testing.T) {
	t.Parallel()
	event := executionBlockedProjectionEvent()

	live := openReplayDB(t, filepath.Join(t.TempDir(), "live-blocked.sqlite"))
	defer closeProjectDB(t, live)
	seedExecutionBlockedProjection(t, live)
	hash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}
	if err := ProjectEvent(live, ProjectionInput{Event: event, EventHash: hash, File: "2026-09-06.jsonl", Line: 1, Offset: 1}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}
	liveState := executionBlockedProjectionState(t, live)

	workspace := newReplayWorkspace(t)
	writeReplayEvents(t, filepath.Join(workspace, ".fid", "journal", "2026-09-06.jsonl"), event)
	replayed := openReplayDB(t, filepath.Join(t.TempDir(), "replay-blocked.sqlite"))
	defer closeProjectDB(t, replayed)
	seedExecutionBlockedProjection(t, replayed)
	if _, err := ReplayJournal(replayed, filepath.Join(workspace, ".fid", "journal"), nil); err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if replayState := executionBlockedProjectionState(t, replayed); replayState != liveState {
		t.Fatalf("execution-blocked projection replay state = %q, want live state %q", replayState, liveState)
	}
}

func TestRunClosureEventsProjectIdenticallyLiveAndReplay(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	updatedAt := "2026-09-06T00:01:00Z"
	activationUpdate := Effect{Op: "update", Table: "kit_activations", PK: map[string]any{"id": "ACT-projection"}, Preconditions: map[string]any{"version": 1}, Prior: map[string]any{"version": 1, "updated_at": "2026-09-06T00:00:00Z"}, Values: map[string]any{"version": 2, "updated_at": updatedAt}}
	switch eventType {
	case "activation.resource-observed":
		event.Effects = []Effect{activationUpdate, {Op: "update", Table: "activation_resources", PK: map[string]any{"activation_id": "ACT-projection", "component": "kit"}, Preconditions: map[string]any{"version": 1}, Prior: map[string]any{"version": 1, "observed_identity": nil, "state": "pending", "evidence_ref": nil, "error_code": nil}, Values: map[string]any{"version": 2, "observed_identity": "kit:new", "state": "new", "evidence_ref": "test:kit", "error_code": nil}}}
	case "activation.recovery-selected":
		event.CommandName = "activation.recover"
		event.Effects = []Effect{activationUpdate,
			{Op: "insert", Table: "activation_recoveries", PK: map[string]any{"activation_id": "ACT-projection"}, Values: map[string]any{"direction": "restore-prior", "actor_attribution": "fid-bootstrap", "reason": "test", "selected_event_id": event.EventID, "completed_event_id": nil}},
			activationStopEffect(event.EventID), activationAttentionEffect(), activationStopLinkEffect(),
		}
	case "activation.recovery-completed":
		event.CommandName = "activation.recover-complete"
		activationUpdate.Values["status"] = "aborted"
		activationUpdate.Prior["status"] = "prepared"
		event.Effects = []Effect{activationUpdate,
			{Op: "update", Table: "activation_recoveries", PK: map[string]any{"activation_id": "ACT-projection"}, Prior: map[string]any{"completed_event_id": nil}, Values: map[string]any{"completed_event_id": event.EventID}},
			{Op: "insert", Table: "stop_resolutions", PK: map[string]any{"stop_id": "STP-activation"}, Values: map[string]any{"actor_id": "fid-bootstrap", "fact_event_id": event.EventID, "resolved_at": updatedAt, "created_event_id": event.EventID}},
		}
	}
	return currentProjectionIdempotency(event, "ACT-projection")
}

func seedActivationRecoveryProjection(t *testing.T, conn *sql.DB, eventType string) {
	t.Helper()
	if _, err := conn.Exec(`INSERT INTO kit_activations (id, kit_identity, backend_identity, phase, status, point_of_no_return, direction, created_at, updated_at, version) VALUES ('ACT-projection','kit','fid',21,'prepared',0,'activate','2026-09-06T00:00:00Z','2026-09-06T00:00:00Z',1)`); err != nil {
		t.Fatalf("seed activation: %v", err)
	}
	for _, component := range []string{"kit", "overlay", "provider", "project", "backend"} {
		if _, err := conn.Exec(`INSERT INTO activation_resources (activation_id, component, version, old_identity, new_identity, observed_identity, state, evidence_ref, error_code) VALUES (?, ?, 1, NULL, ?, NULL, 'pending', NULL, NULL)`, "ACT-projection", component, component+":new"); err != nil {
			t.Fatalf("seed activation resource %s: %v", component, err)
		}
	}
	if eventType == "activation.recovery-completed" {
		if _, err := conn.Exec(`UPDATE activation_resources SET state='old',evidence_ref='test:prior-agreement' WHERE activation_id='ACT-projection'`); err != nil {
			t.Fatal(err)
		}
		if _, err := conn.Exec(`INSERT INTO activation_recoveries (activation_id, direction, actor_attribution, reason, selected_event_id, completed_event_id) VALUES ('ACT-projection','restore-prior','fid-bootstrap','test','EVT-selected',NULL)`); err != nil {
			t.Fatalf("seed activation recovery: %v", err)
		}
		if _, err := conn.Exec(`INSERT INTO stops (id, subject_type, subject_id, code, blocked_actions_json, details_json, resolution_condition_json, resolver_capability, actor_id, raised_at, created_event_id) VALUES ('STP-activation','activation','ACT-projection','activation-incomplete','["claim.acquire"]','{"activation_id":"ACT-projection","components":["kit","overlay","provider","project","backend"],"point_of_no_return":false}','{"subject":{"id":"ACT-projection","type":"activation"},"type":"activation-recovered","budget_name":null,"capability_id":null}','stop.resolve','fid-bootstrap','2026-09-06T00:00:00Z','EVT-selected')`); err != nil {
			t.Fatalf("seed activation Stop: %v", err)
		}
		if _, err := conn.Exec(`INSERT INTO attention_items (id, semantic_key, subject_type, subject_id, code, status, version, raised_at, resolution_event_id) VALUES ('ATT-activation','activation-incomplete:ACT-projection','activation','ACT-projection','activation-incomplete','open',1,'2026-09-06T00:00:00Z',NULL)`); err != nil {
			t.Fatalf("seed activation Attention: %v", err)
		}
		if _, err := conn.Exec(`INSERT INTO attention_stop_links (attention_id, stop_id) VALUES ('ATT-activation','STP-activation')`); err != nil {
			t.Fatalf("seed activation Stop link: %v", err)
		}
	}
}

func activationRecoveryProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()
	var state string
	if err := conn.QueryRow(`SELECT status || ':' || version || ':' || updated_at || ':' || (SELECT group_concat(component || ':' || version || ':' || state || ':' || COALESCE(observed_identity,''), ',') FROM (SELECT component, version, state, observed_identity FROM activation_resources WHERE activation_id = 'ACT-projection' ORDER BY component)) || ':' || COALESCE((SELECT direction || ':' || selected_event_id || ':' || COALESCE(completed_event_id,'') FROM activation_recoveries WHERE activation_id = 'ACT-projection'),'') || ':' || (SELECT COUNT(*) FROM attention_stop_links l JOIN stops s ON s.id = l.stop_id WHERE s.subject_id = 'ACT-projection') || ':' || (SELECT COUNT(*) FROM stop_resolutions WHERE stop_id = 'STP-activation') FROM kit_activations WHERE id = 'ACT-projection'`).Scan(&state); err != nil {
		t.Fatalf("read activation projection: %v", err)
	}
	return state
}

func activationStopEffect(eventID string) Effect {
	return Effect{Op: "insert", Table: "stops", PK: map[string]any{"id": "STP-activation"}, Values: map[string]any{
		"subject_type": "activation", "subject_id": "ACT-projection", "code": "activation-incomplete", "blocked_actions_json": `["claim.acquire","envelope.complete","gate.run","issue.transition-terminal-success","run.start"]`,
		"details_json": `{"activation_id":"ACT-projection","components":["kit","overlay","provider","project","backend"],"point_of_no_return":false}`, "resolution_condition_json": `{"subject":{"id":"ACT-projection","type":"activation"},"type":"activation-recovered","budget_name":null,"capability_id":null}`,
		"resolver_capability": "stop.resolve", "actor_id": "fid-bootstrap", "raised_at": "2026-09-06T00:01:00Z", "created_event_id": eventID,
	}}
}

func activationAttentionEffect() Effect {
	return Effect{Op: "insert", Table: "attention_items", PK: map[string]any{"id": "ATT-activation"}, Values: map[string]any{
		"semantic_key": "sha256:" + strings.Repeat("1", 64), "subject_type": "activation", "subject_id": "ACT-projection", "code": "activation-incomplete", "status": "open", "version": 1, "raised_at": "2026-09-06T00:01:00Z", "resolution_event_id": nil,
	}}
}

func activationStopLinkEffect() Effect {
	return Effect{Op: "insert", Table: "attention_stop_links", PK: map[string]any{"attention_id": "ATT-activation", "stop_id": "STP-activation"}}
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
			{Op: "update", Table: "runs", PK: map[string]any{"id": "RUN-closure"}, Values: map[string]any{"status": status, "terminal_repository_json": projectionRepositoryJSON, "ended_at": "2026-09-06T00:00:00Z"}, Preconditions: map[string]any{"status": "active"}, Prior: map[string]any{"status": "active", "terminal_repository_json": nil, "ended_at": nil}},
			{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-closure"}, Values: map[string]any{"version": 2}, Preconditions: map[string]any{"version": 1}, Prior: map[string]any{"version": 1}},
		},
	}
	if eventType == "run.abandoned" {
		event.CommandName = "run.abandon"
		event.Effects = append(event.Effects, Effect{Op: "update", Table: "claims", PK: map[string]any{"id": "CLM-closure"}, Values: map[string]any{"version": 2, "status": "released", "released_at": "2026-09-06T00:00:00Z", "release_reason": "abandoned"}, Preconditions: map[string]any{"status": "active", "version": 1}, Prior: map[string]any{"version": 1, "status": "active", "released_at": nil, "release_reason": nil}})
	}
	return currentProjectionIdempotency(event, "RUN-closure")
}

func seedRunClosureProjection(t *testing.T, conn *sql.DB) {
	t.Helper()
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable foreign keys for closure projection seed: %v", err)
	}
	statements := []string{
		`INSERT INTO claims (id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at) VALUES ('CLM-closure', 'ISS-closure', 'ENV-closure', 1, 'agent:closure', 'WT-closure', '{"paths":[]}', 'mutation', 1, 'active', '2026-09-06T00:00:00Z', '2026-09-07T00:00:00Z')`,
		`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, status, started_at) VALUES ('RUN-closure', 'ISS-closure', 'ENV-closure', 1, 'CLM-closure', 1, 'agent:closure', '{"base_oid":null,"capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":null,"project_id":"PRJ-projection","tree_oid":null,"worktree_id":"WT-projection"}', 'active', '2026-09-06T00:00:00Z')`,
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

func recoveredProjectionEvent() Event {
	event := Event{EventID: "EVT-20260906-000003", EventSchema: EventSchemaV2, EventType: "run.recovered", OccurredAt: time.Date(2026, time.September, 6, 0, 1, 0, 0, time.UTC), ActorV2: EventActor{ActorID: "agent:recovery", ActorType: "agent", Producer: "activated-harness", ProvenanceRef: "test:recovery"}, Aggregate: Aggregate{ID: "RUN-prior", Type: "run", Sequence: 1}, CommandName: "run.recover"}
	event.Effects = []Effect{
		{Op: "update", Table: "runs", PK: map[string]any{"id": "RUN-prior"}, Preconditions: map[string]any{"status": "active"}, Prior: map[string]any{"status": "active", "terminal_repository_json": nil, "ended_at": nil}, Values: map[string]any{"status": "interrupted", "terminal_repository_json": projectionRepositoryJSON, "ended_at": "2026-09-06T00:01:00Z"}},
		{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-prior"}, Preconditions: map[string]any{"version": 1}, Prior: map[string]any{"version": 1}, Values: map[string]any{"version": 2}},
		{Op: "update", Table: "claims", PK: map[string]any{"id": "CLM-prior"}, Preconditions: map[string]any{"status": "active", "version": 1}, Prior: map[string]any{"version": 1, "status": "active", "released_at": nil, "release_reason": nil}, Values: map[string]any{"version": 2, "status": "superseded", "released_at": "2026-09-06T00:01:00Z", "release_reason": "recovery-superseded"}},
		{Op: "insert", Table: "claims", PK: map[string]any{"id": "CLM-successor"}, Values: claimProjectionValues("ISS-projection", "ENV-projection", "agent:recovery", "WT-projection")},
		{Op: "insert", Table: "runs", PK: map[string]any{"id": "RUN-successor"}, Values: runProjectionValues("ISS-projection", "ENV-projection", "CLM-successor", "agent:recovery", "active", "2026-09-06T00:01:00Z")},
		{Op: "insert", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-successor"}, Values: map[string]any{"version": 1, "executor_id": "executor-successor", "obligation_id": nil}},
		{Op: "insert", Table: "run_recoveries", PK: map[string]any{"id": "REC-projection"}, Values: map[string]any{"prior_run_id": "RUN-prior", "successor_run_id": "RUN-successor", "prior_claim_id": "CLM-prior", "successor_claim_id": "CLM-successor", "observation_id": "OBS-liveness", "reconciliation_json": projectionFixtureJSON(projectionReconciliation()), "adopted_paths_json": `[]`, "authority_json": `{"mode":"ordinary","actor_id":"agent:recovery","reason":"test","reference":null,"grant_ref":null}`, "reason": "test", "created_event_id": event.EventID}},
	}
	return currentProjectionIdempotency(event, "REC-projection")
}

func budgetObservedProjectionEvent() Event {
	event := Event{EventID: "EVT-20260906-000004", EventSchema: EventSchemaV2, EventType: "run.budget-observed", OccurredAt: time.Date(2026, time.September, 6, 0, 1, 0, 0, time.UTC), ActorV2: EventActor{ActorID: "agent:budget", ActorType: "agent", Producer: "activated-harness", ProvenanceRef: "test:budget"}, Aggregate: Aggregate{ID: "RUN-budget", Type: "run", Sequence: 1}, CommandName: "run.observe-budget"}
	event.Effects = []Effect{
		{Op: "insert", Table: "budget_samples", PK: map[string]any{"id": "BUD-projection"}, Values: map[string]any{"envelope_id": "ENV-projection", "envelope_revision": 1, "budget_name": "tokens", "run_id": "RUN-budget", "actor_id": "agent:budget", "source": "activated-harness", "source_sample_id": "sample-1", "sample_kind": "baseline", "previous_sample_id": nil, "value": 0, "measurement_ref_json": nil, "observed_at": "2026-09-06T00:01:00Z", "created_event_id": event.EventID}},
		{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-budget"}, Preconditions: map[string]any{"version": 1}, Prior: map[string]any{"version": 1}, Values: map[string]any{"version": 2}},
	}
	return currentProjectionIdempotency(event, "BUD-projection")
}

func executionBlockedProjectionEvent() Event {
	event := Event{EventID: "EVT-20260906-000005", EventSchema: EventSchemaV2, EventType: "run.execution-blocked", OccurredAt: time.Date(2026, time.September, 6, 0, 1, 0, 0, time.UTC), ActorV2: EventActor{ActorID: "agent:gate", ActorType: "agent", Producer: "fid", ProvenanceRef: "test:gate"}, Aggregate: Aggregate{ID: "RUN-blocked", Type: "run", Sequence: 1}, CommandName: "gate.run"}
	event.Effects = []Effect{
		{Op: "insert", Table: "stops", PK: map[string]any{"id": "STP-blocked"}, Values: map[string]any{"subject_type": "run", "subject_id": "RUN-blocked", "code": "budget-data-unavailable", "blocked_actions_json": `["gate.run","run.recover","run.start"]`, "details_json": `{"envelope_id":"ENV-projection","budget_name":"tokens","unit":"input-tokens","maximum":10,"consumed":null,"remaining":null}`, "resolution_condition_json": `{"subject":{"id":"RUN-blocked","type":"run"},"type":"budget-observed","budget_name":"tokens","capability_id":null}`, "resolver_capability": "stop.resolve", "actor_id": "agent:gate", "raised_at": "2026-09-06T00:01:00Z", "created_event_id": event.EventID}},
		{Op: "insert", Table: "attention_items", PK: map[string]any{"id": "ATT-blocked"}, Values: map[string]any{"semantic_key": "sha256:" + strings.Repeat("2", 64), "subject_type": "run", "subject_id": "RUN-blocked", "code": "budget-data-unavailable", "status": "open", "version": 1, "raised_at": "2026-09-06T00:01:00Z", "resolution_event_id": nil}},
		{Op: "insert", Table: "attention_stop_links", PK: map[string]any{"attention_id": "ATT-blocked", "stop_id": "STP-blocked"}},
		{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-blocked"}, Preconditions: map[string]any{"version": 1}, Prior: map[string]any{"version": 1}, Values: map[string]any{"version": 2}},
	}
	return currentProjectionIdempotency(event, "STP-blocked")
}

func claimProjectionValues(issueID, envelopeID, actorID, worktreeID string) map[string]any {
	return map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "actor_id": actorID, "worktree_id": worktreeID, "scope_json": `{"paths":[]}`, "mode": "mutation", "version": 1, "status": "active", "acquired_at": "2026-09-06T00:01:00Z", "expires_at": "2026-09-07T00:01:00Z", "released_at": nil, "release_reason": nil}
}

func runProjectionValues(issueID, envelopeID, claimID, actorID, status, startedAt string) map[string]any {
	return map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": claimID, "claim_version": 1, "actor_id": actorID, "start_repository_json": projectionRepositoryJSON, "terminal_repository_json": nil, "status": status, "started_at": startedAt, "ended_at": nil}
}

func seedRecoveredProjection(t *testing.T, conn *sql.DB) {
	t.Helper()
	seedProjectionStatements(t, conn, []string{
		`INSERT INTO claims (id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at) VALUES ('CLM-prior','ISS-projection','ENV-projection',1,'agent:recovery','WT-prior','{"paths":[]}','mutation',1,'active','2026-09-06T00:00:00Z','2026-09-07T00:00:00Z')`,
		`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, status, started_at) VALUES ('RUN-prior','ISS-projection','ENV-projection',1,'CLM-prior',1,'agent:recovery','{"base_oid":null,"capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":null,"project_id":"PRJ-projection","tree_oid":null,"worktree_id":"WT-projection"}','active','2026-09-06T00:00:00Z')`,
		`INSERT INTO run_runtime (run_id, version, executor_id, obligation_id) VALUES ('RUN-prior',1,'executor-prior',NULL)`,
		`INSERT INTO run_observations (id, run_id, kind, actor_id, observed_at, observation_json, created_event_id) VALUES ('OBS-liveness','RUN-prior','liveness','agent:recovery','2026-09-06T00:00:00Z','{"id":"OBS-liveness","kind":"liveness","run_id":"RUN-prior"}','EVT-liveness')`,
	})
}

func seedBudgetProjection(t *testing.T, conn *sql.DB) {
	t.Helper()
	seedProjectionStatements(t, conn, []string{
		`INSERT INTO claims (id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at) VALUES ('CLM-budget','ISS-projection','ENV-projection',1,'agent:budget','WT-budget','{"paths":[]}','mutation',1,'active','2026-09-06T00:00:00Z','2026-09-07T00:00:00Z')`,
		`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, status, started_at) VALUES ('RUN-budget','ISS-projection','ENV-projection',1,'CLM-budget',1,'agent:budget','{"base_oid":null,"capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":null,"project_id":"PRJ-projection","tree_oid":null,"worktree_id":"WT-projection"}','active','2026-09-06T00:00:00Z')`,
		`INSERT INTO run_runtime (run_id, version, executor_id, obligation_id) VALUES ('RUN-budget',1,'executor-budget',NULL)`,
	})
	seedProjectionBudget(t, conn, "RUN-budget")
}

func seedExecutionBlockedProjection(t *testing.T, conn *sql.DB) {
	t.Helper()
	seedProjectionStatements(t, conn, []string{
		`INSERT INTO claims (id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at) VALUES ('CLM-blocked','ISS-projection','ENV-projection',1,'agent:gate','WT-blocked','{"paths":[]}','mutation',1,'active','2026-09-06T00:00:00Z','2026-09-07T00:00:00Z')`,
		`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, status, started_at) VALUES ('RUN-blocked','ISS-projection','ENV-projection',1,'CLM-blocked',1,'agent:gate','{"base_oid":null,"capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":null,"project_id":"PRJ-projection","tree_oid":null,"worktree_id":"WT-projection"}','active','2026-09-06T00:00:00Z')`,
		`INSERT INTO run_runtime (run_id, version, executor_id, obligation_id) VALUES ('RUN-blocked',1,'executor-blocked',NULL)`,
	})
	seedProjectionBudget(t, conn, "RUN-blocked")
}

func seedProjectionStatements(t *testing.T, conn *sql.DB, statements []string) {
	t.Helper()
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable foreign keys for projection seed: %v", err)
	}
	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("seed projection: %v", err)
		}
	}
}

func recoveredProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()
	var state string
	if err := conn.QueryRow(`SELECT (SELECT status || ':' || COALESCE(ended_at,'') FROM runs WHERE id = 'RUN-prior') || '|' || (SELECT status || ':' || version || ':' || COALESCE(release_reason,'') FROM claims WHERE id = 'CLM-prior') || '|' || (SELECT status || ':' || claim_id || ':' || actor_id || ':' || executor_id FROM runs JOIN run_runtime ON run_runtime.run_id = runs.id WHERE runs.id = 'RUN-successor') || '|' || (SELECT prior_run_id || ':' || successor_run_id || ':' || prior_claim_id || ':' || successor_claim_id || ':' || observation_id || ':' || created_event_id FROM run_recoveries WHERE id = 'REC-projection') || '|' || (SELECT result_entity_id || ':' || result_event_id FROM idempotency_results)`).Scan(&state); err != nil {
		t.Fatalf("read recovered projection: %v", err)
	}
	return state
}

func budgetProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()
	var state string
	if err := conn.QueryRow(`SELECT (SELECT id || ':' || run_id || ':' || budget_name || ':' || source || ':' || sample_kind || ':' || value || ':' || created_event_id FROM budget_samples WHERE id = 'BUD-projection') || '|' || (SELECT version FROM run_runtime WHERE run_id = 'RUN-budget') || '|' || (SELECT result_entity_id || ':' || result_event_id FROM idempotency_results)`).Scan(&state); err != nil {
		t.Fatalf("read budget projection: %v", err)
	}
	return state
}

func executionBlockedProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()
	var state string
	if err := conn.QueryRow(`SELECT (SELECT id || ':' || subject_type || ':' || subject_id || ':' || code || ':' || created_event_id FROM stops WHERE id = 'STP-blocked') || '|' || (SELECT id || ':' || semantic_key || ':' || subject_id || ':' || status || ':' || version FROM attention_items WHERE id = 'ATT-blocked') || '|' || (SELECT attention_id || ':' || stop_id FROM attention_stop_links WHERE stop_id = 'STP-blocked') || '|' || (SELECT version FROM run_runtime WHERE run_id = 'RUN-blocked') || '|' || (SELECT result_entity_id || ':' || result_event_id FROM idempotency_results)`).Scan(&state); err != nil {
		t.Fatalf("read execution-blocked projection: %v", err)
	}
	return state
}

func checkpointedProjectionEvent() Event {
	event := Event{
		EventID:     "EVT-20260906-000001",
		EventSchema: EventSchemaV2,
		EventType:   "run.checkpointed",
		OccurredAt:  time.Date(2026, time.September, 6, 0, 0, 0, 0, time.UTC),
		ActorV2:     EventActor{ActorID: "agent:checkpoint", ActorType: "agent", Producer: "activated-harness", ProvenanceRef: "test:checkpoint"},
		Aggregate:   Aggregate{ID: "RUN-checkpoint", Type: "run", Sequence: 1},
		CommandName: "run.checkpoint",
		Effects: []Effect{
			{Op: "insert", Table: "run_observations", PK: map[string]any{"id": "OBS-checkpoint"}, Values: map[string]any{"run_id": "RUN-checkpoint", "kind": "checkpoint", "actor_id": "agent:checkpoint", "observed_at": "2026-09-06T00:00:00Z", "observation_json": projectionFixtureJSON(projectionCheckpoint()), "created_event_id": "EVT-20260906-000001"}},
			{Op: "update", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-checkpoint"}, Values: map[string]any{"version": 2}, Preconditions: map[string]any{"version": 1}, Prior: map[string]any{"version": 1}},
		},
	}
	return currentProjectionIdempotency(event, "OBS-checkpoint")
}

func currentProjectionIdempotency(event Event, entityID string) Event {
	event.Idempotency = &EventIdempotency{Key: "projection-fixture", RequestDigest: "sha256:" + strings.Repeat("1", 64), Scope: event.CommandName + ":PRJ-projection"}
	event.Effects = append(event.Effects, Effect{Op: "insert", Table: "idempotency_results", PK: map[string]any{"operation_scope": event.Idempotency.Scope, "idempotency_key": event.Idempotency.Key}, Values: map[string]any{"request_digest": event.Idempotency.RequestDigest, "request_field_digests_json": "{}", "result_entity_id": entityID, "result_event_id": event.EventID, "aggregate_sequence": event.Aggregate.Sequence, "result_json": projectionFixtureJSON(projectionResult(event)), "created_at": event.OccurredAt.Format(time.RFC3339Nano)}})
	return event
}

func seedCheckpointProjectionRun(t *testing.T, conn *sql.DB) {
	t.Helper()
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable foreign keys for projection seed: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, status, started_at) VALUES ('RUN-checkpoint', 'ISS-checkpoint', 'ENV-checkpoint', 1, 'CLM-checkpoint', 1, 'agent:checkpoint', '{"base_oid":null,"capture_source":"git","dirty":false,"dirty_digest":null,"head_oid":null,"project_id":"PRJ-projection","tree_oid":null,"worktree_id":"WT-projection"}', 'active', '2026-09-06T00:00:00Z')`); err != nil {
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
