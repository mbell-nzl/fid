package workflow

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
)

func TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse(t *testing.T) {
	t.Parallel()
	options, request := recoveryFixture(t)
	beforeEnvelope, err := GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, "ENV-recovery")
	if err != nil {
		t.Fatalf("GetEnvelope(before) error = %v", err)
	}
	beforeIssue := recoveryIssueStatus(t, options.WorkspaceRoot)
	result, err := RecoverRun(options, request)
	if err != nil {
		t.Fatalf("RecoverRun() error = %v", err)
	}
	if result.Replayed || result.Run.ID == "" || result.Run.ID == request.PriorRunID || result.Recovery.SuccessorRunID != result.Run.ID || result.Run.PredecessorRunID == nil || *result.Run.PredecessorRunID != request.PriorRunID {
		t.Fatalf("RecoverRun() = %#v, want linked successor Run", result)
	}
	prior, err := GetRun(options.WorkspaceRoot, options.MigrationsDir, request.PriorRunID)
	if err != nil {
		t.Fatalf("GetRun(prior) error = %v", err)
	}
	if prior.Status != "interrupted" || prior.SuccessorRunID == nil || *prior.SuccessorRunID != result.Run.ID {
		t.Fatalf("prior Run = %#v, want interrupted linked predecessor", prior)
	}
	interrupted, err := ListRuns(RunListOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, Status: "interrupted", Limit: 10})
	if err != nil || len(interrupted) != 1 || interrupted[0].ID != prior.ID {
		t.Fatalf("ListRuns(interrupted) = %#v, %v; want predecessor %q", interrupted, err, prior.ID)
	}
	afterEnvelope, err := GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, "ENV-recovery")
	if err != nil {
		t.Fatalf("GetEnvelope(after) error = %v", err)
	}
	if !reflect.DeepEqual(afterEnvelope, beforeEnvelope) || recoveryIssueStatus(t, options.WorkspaceRoot) != beforeIssue {
		t.Fatalf("recovery changed immutable Envelope or Issue: before=%#v after=%#v", beforeEnvelope, afterEnvelope)
	}
	replayed, err := RecoverRun(options, request)
	if err != nil {
		t.Fatalf("RecoverRun(replay) error = %v", err)
	}
	if !replayed.Replayed || replayed.Run.ID != result.Run.ID || replayed.Recovery.ID != result.Recovery.ID {
		t.Fatalf("RecoverRun(replay) = %#v, want original successor", replayed)
	}
	changed := request
	changed.TTLSeconds++
	if _, err := RecoverRun(options, changed); err == nil {
		t.Fatal("RecoverRun(changed same key) error = nil, want idempotency conflict")
	} else {
		var conflict *journal.IdempotencyConflictError
		if !errors.As(err, &conflict) {
			t.Fatalf("RecoverRun(changed same key) error = %v, want idempotency conflict", err)
		}
	}
}

func recoveryIssueStatus(t *testing.T, workspace string) string {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	var status string
	if err := conn.QueryRow(`SELECT status FROM issues WHERE id = 'ISS-recovery'`).Scan(&status); err != nil {
		t.Fatalf("load recovery Issue: %v", err)
	}
	return status
}

func TestRecoverRunRejectsUnknownLivenessForOrdinaryAndForce(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"ordinary", "force"} {
		t.Run(mode, func(t *testing.T) {
			options, request := recoveryFixture(t)
			setRecoveryObservationState(t, options.WorkspaceRoot, "unknown")
			if mode == "force" {
				configureForceRecoveryAuthority(t, options.WorkspaceRoot)
				reference := "HUM-1"
				request.Authority = RecoveryAuthority{Mode: "force", ActorID: "agent:recovery", Reason: "operator takeover", Reference: &reference, GrantRef: &RecoveryGrantRef{EnvelopeID: "ENV-recovery", Revision: 1, Ordinal: 0}}
			}
			_, err := RecoverRun(options, request)
			var unknown *ExecutorLivenessUnknownError
			if !errors.As(err, &unknown) || unknown.ObservationID != request.ObservationID {
				t.Fatalf("RecoverRun() error = %v, want executor-liveness-unknown", err)
			}
		})
	}
}

func TestRecoverRunRejectsApplicableStopWithoutEffects(t *testing.T) {
	t.Parallel()
	options, request := recoveryFixture(t)
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	if err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "test recovery Stop"}, func(lock *journal.WriteLock) error {
		event := recoveryStopEvent("", "STOP-recovery", "ATT-recovery")
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		_, err = journal.CommitEvent(lock, conn, recoveryStopEvent(prepared.EventID, "STOP-recovery", "ATT-recovery"))
		return err
	}); err != nil {
		t.Fatalf("commit recovery Stop: %v", err)
	}
	before := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
	_, err = RecoverRun(options, request)
	var blocked *StopBlockedError
	if !errors.As(err, &blocked) || len(blocked.Stops) != 1 || blocked.Stops[0].ID != "STOP-recovery" {
		t.Fatalf("RecoverRun(applicable Stop) error = %v, want STOP-recovery", err)
	}
	if after := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot)); string(after) != string(before) {
		t.Fatal("Stop-blocked recovery appended an Event")
	}
}

// TestRecoverRunAppliesStopsBeforeOrdinaryOrForceMutation covers both recovery
// authorities with admitted Stop facts. Stop producer and resolution validation
// are covered separately; this matrix isolates the recovery guard and verifies
// that rejected requests do not append a recovery Event or successor Run.
func TestRecoverRunAppliesStopsBeforeOrdinaryOrForceMutation(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"ordinary", "force"} {
		for _, state := range []string{"applicable", "unrelated-subject", "unrelated-action", "resolved"} {
			t.Run(mode+"/"+state, func(t *testing.T) {
				options, request := recoveryFixture(t)
				if mode == "force" {
					setRecoveryObservationState(t, options.WorkspaceRoot, "live")
					configureForceRecoveryAuthority(t, options.WorkspaceRoot)
					reference := "HUM-stop-matrix"
					request.Authority = RecoveryAuthority{Mode: "force", ActorID: options.ActorID, Reason: "operator takeover", Reference: &reference, GrantRef: &RecoveryGrantRef{EnvelopeID: "ENV-recovery", Revision: 1, Ordinal: 0}}
					conn, err := openWorkflowDB(options.WorkspaceRoot)
					if err != nil {
						t.Fatalf("open force recovery projection: %v", err)
					}
					if _, err := conn.Exec(`UPDATE claims SET status = 'active', expires_at = ? WHERE id = 'CLM-prior'`, time.Now().UTC().Add(time.Minute).Format(time.RFC3339Nano)); err != nil {
						closeWorkflowDB(conn)
						t.Fatalf("activate force recovery Claim: %v", err)
					}
					closeWorkflowDB(conn)
				}
				wantIDs := seedRecoveryStopMatrix(t, options.WorkspaceRoot, state)
				beforeJournal := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
				beforeSuccessors := recoverySuccessorCount(t, options.WorkspaceRoot)

				result, err := RecoverRun(options, request)
				if len(wantIDs) == 0 {
					if err != nil || result.Run.PredecessorRunID == nil || *result.Run.PredecessorRunID != request.PriorRunID {
						t.Fatalf("nonqualifying Stop recovery = %#v, %v; want linked successor", result, err)
					}
					if recoverySuccessorCount(t, options.WorkspaceRoot) != beforeSuccessors+1 {
						t.Fatal("nonqualifying Stop did not create exactly one successor")
					}
					return
				}
				var blocked *StopBlockedError
				if !errors.As(err, &blocked) || blocked.Action != runRecoverAction || len(blocked.Stops) != len(wantIDs) {
					t.Fatalf("applicable Stop recovery error = %v, want ordered Stops %v", err, wantIDs)
				}
				for index, stop := range blocked.Stops {
					if stop.ID != wantIDs[index] || stop.Code != "operator-stop" {
						t.Fatalf("Stop %d = %#v, want %q operator-stop", index, stop, wantIDs[index])
					}
				}
				if after := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot)); string(after) != string(beforeJournal) {
					t.Fatal("Stop-blocked recovery appended an Event")
				}
				if recoverySuccessorCount(t, options.WorkspaceRoot) != beforeSuccessors {
					t.Fatal("Stop-blocked recovery created a successor Run")
				}
			})
		}
	}
}

func seedRecoveryStopMatrix(t *testing.T, workspace, state string) []string {
	t.Helper()
	subjectID, actions := "RUN-prior", `["run.recover"]`
	if state == "unrelated-subject" {
		subjectID = "RUN-unrelated"
	}
	if state == "unrelated-action" {
		actions = `["gate.run"]`
	}
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open recovery Stop projection: %v", err)
	}
	defer closeWorkflowDB(conn)
	if _, err := conn.Exec(`INSERT INTO stops (id, subject_type, subject_id, code, blocked_actions_json, details_json, resolution_condition_json, resolver_capability, actor_id, raised_at, created_event_id) VALUES ('STOP-recovery-matrix', 'run', ?, 'operator-stop', ?, '{"reason":"recovery guard matrix","reference":"policy:matrix"}', '{"type":"work-ended","subject":{"type":"run","id":"RUN-prior"}}', 'stop.resolve', 'agent:recovery', '2026-01-01T00:00:00Z', 'EVT-recovery-matrix-raised')`, subjectID, actions); err != nil {
		t.Fatalf("seed recovery Stop: %v", err)
	}
	if state == "resolved" {
		if _, err := conn.Exec(`INSERT INTO stop_resolutions (stop_id, actor_id, fact_event_id, resolved_at, created_event_id) VALUES ('STOP-recovery-matrix', 'agent:recovery', 'EVT-recovery-matrix-fact', '2026-01-01T00:00:01Z', 'EVT-recovery-matrix-resolved')`); err != nil {
			t.Fatalf("seed resolved recovery Stop: %v", err)
		}
	}
	if state == "applicable" {
		return []string{"STOP-recovery-matrix"}
	}
	return nil
}

func recoverySuccessorCount(t *testing.T, workspace string) int {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open recovery successor projection: %v", err)
	}
	defer closeWorkflowDB(conn)
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM run_recoveries WHERE prior_run_id = 'RUN-prior'`).Scan(&count); err != nil {
		t.Fatalf("count recovery successors: %v", err)
	}
	return count
}

func recoveryStopEvent(eventID, stopID, attentionID string) journal.Event {
	event := stopRaisedEvent(eventID, stopID, attentionID, true)
	event.Effects[0].Values["subject_id"] = "RUN-prior"
	event.Effects[0].Values["resolution_condition_json"] = `{"type":"work-ended","subject":{"type":"run","id":"RUN-prior"},"budget_name":null,"capability_id":null}`
	event.Effects[1].Values["subject_id"] = "RUN-prior"
	return refreshStopFixtureResult(event)
}

func TestRecoverRunForceRequiresCurrentHumanAuthorityAndFencesPriorRun(t *testing.T) {
	t.Parallel()
	options, request := recoveryFixture(t)
	setRecoveryObservationState(t, options.WorkspaceRoot, "live")
	reference := "HUM-1"
	request.Authority = RecoveryAuthority{Mode: "force", ActorID: options.ActorID, Reason: "operator takeover", Reference: &reference, GrantRef: &RecoveryGrantRef{EnvelopeID: "ENV-recovery", Revision: 1, Ordinal: 0}}
	if _, err := RecoverRun(options, request); err == nil {
		t.Fatal("RecoverRun(non-Human force) error = nil, want rejection")
	}
	configureForceRecoveryAuthority(t, options.WorkspaceRoot)
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	if _, err := conn.Exec(`UPDATE claims SET status = 'active', expires_at = ? WHERE id = 'CLM-prior'`, time.Now().UTC().Add(time.Minute).Format(time.RFC3339Nano)); err != nil {
		t.Fatalf("activate prior claim: %v", err)
	}
	closeWorkflowDB(conn)
	result, err := RecoverRun(options, request)
	if err != nil {
		t.Fatalf("RecoverRun(Human force) error = %v", err)
	}
	prior, err := GetRun(options.WorkspaceRoot, options.MigrationsDir, request.PriorRunID)
	if err != nil {
		t.Fatalf("GetRun(prior) error = %v", err)
	}
	if prior.Status != "interrupted" || result.Run.PredecessorRunID == nil || *result.Run.PredecessorRunID != prior.ID {
		t.Fatalf("force recovery Run chain = prior %#v successor %#v, want fenced predecessor", prior, result.Run)
	}
	if _, err := CloseRun(options, RunCloseRequest{RunID: prior.ID, ExpectedVersion: prior.Version, ObservationID: "OBS-fenced", ExpectedRepository: expectedRunRepositoryForTest(t, options.WorkspaceRoot)}); err == nil {
		t.Fatal("CloseRun(fenced predecessor) error = nil, want inactive predecessor")
	}
}

func recoveryFixture(t *testing.T) (RunOptions, RunRecoveryRequest) {
	t.Helper()
	workspace := t.TempDir()
	for _, arguments := range [][]string{{"init", "-q"}, {"config", "user.email", "fid@example.test"}, {"config", "user.name", "fid test"}} {
		runReconciliationGit(t, workspace, arguments...)
	}
	if err := os.WriteFile(filepath.Join(workspace, "base.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write repository base: %v", err)
	}
	runReconciliationGit(t, workspace, "add", "base.txt")
	runReconciliationGit(t, workspace, "commit", "-qm", "base")
	migrationsDir := workflowMigrationsDir(t)
	if err := cloneWorkflowPristine(workspace); err != nil {
		t.Fatalf("clone initialized recovery workspace: %v", err)
	}
	observation, err := project.ObserveRepository(workspace)
	if err != nil {
		t.Fatalf("ObserveRepository() error = %v", err)
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(workspace), "config.toml"))
	if err != nil {
		t.Fatalf("ReadProjectConfig() error = %v", err)
	}
	repository := repositoryFromObservation(observation)
	repository.CaptureSource, repository.ProjectID = "git", projectConfig.ProjectID
	now := time.Now().UTC()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	seedRecoveryFixture(t, conn, repository, now)
	reconciliation, err := ReconcileRepository(repository, observation, nil, nil)
	if err != nil {
		t.Fatalf("ReconcileRepository() error = %v", err)
	}
	return RunOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "agent:recovery", IdempotencyKey: "recover-lost-response", Now: now}, RunRecoveryRequest{
		PriorRunID: "RUN-prior", ExpectedRunVersion: 1, ExpectedClaimVersion: 1, ObservationID: "OBS-terminated", ReconciliationDigest: reconciliation.Digest, AdoptedPaths: []string{},
		Authority: RecoveryAuthority{Mode: "ordinary", ActorID: "agent:recovery", Reason: "trusted termination"}, Executor: Executor{ID: "executor-successor", Source: "activated-harness"}, TTLSeconds: 60,
	}
}

func seedRecoveryFixture(t *testing.T, conn *sql.DB, repository Repository, now time.Time) {
	t.Helper()
	timestamp := now.Format(time.RFC3339Nano)
	repositoryJSON := string(mustRunJSON(repository))
	observation := RunObservation{ID: "OBS-terminated", RunID: "RUN-prior", Kind: "liveness", ActorID: "agent:observer", ObservedAt: timestamp, Liveness: Liveness{ExecutorID: "executor-prior", State: "terminated", SourceRef: "receipt:terminated", ObservedAt: timestamp}}
	observationJSON := string(mustRunJSON(observation))
	if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable fixture foreign keys: %v", err)
	}
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO kit_activations (id, kit_identity, backend_identity, phase, status, point_of_no_return, direction, created_at, updated_at, version) VALUES (?, 'kit', 'backend', 22, 'finalized', 0, 'activate', ?, ?, 1)`, []any{"ACT-recovery", timestamp, timestamp}},
		{`INSERT INTO activation_actor_profiles (activation_id, actor_id, actor_type, producer, provenance_ref, capabilities_json) VALUES ('ACT-recovery', 'agent:recovery', 'agent', 'test', 'test:recovery', '["run.recover"]')`, nil},
		{`INSERT INTO issues (id, kind, title, workflow_id, status) VALUES ('ISS-recovery', 'task', 'recovery fixture', 'task-default', 'todo')`, nil},
		{`INSERT INTO work_envelopes (id, current_revision, status, created_at) VALUES ('ENV-recovery', 1, 'active', ?)`, []any{timestamp}},
		{`INSERT INTO envelope_revisions (envelope_id, revision, prior_revision, outcome, non_goals_json, risk_result_json, repository_scope_json, status, created_actor_id, approved_actor_id, created_at, activated_at, completed_at) VALUES ('ENV-recovery', 1, NULL, 'recover', '[]', '{}', ?, 'active', 'agent:recovery', 'agent:recovery', ?, ?, NULL)`, []any{repositoryJSON, timestamp, timestamp}},
		{`INSERT INTO envelope_issue_selection (envelope_id, revision, issue_id, primary_issue) VALUES ('ENV-recovery', 1, 'ISS-recovery', 1)`, nil},
		{`INSERT INTO envelope_authority_grants (envelope_id, revision, ordinal, actor_selector_json, actions_json, scope_json) VALUES ('ENV-recovery', 1, 0, '{"actor_id":"agent:recovery","required_capabilities":["run.recover"]}', '["run.recover"]', '{"issue_id":"ISS-recovery"}')`, nil},
		{`INSERT INTO claims (id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at, released_at, release_reason) VALUES ('CLM-prior', 'ISS-recovery', 'ENV-recovery', 1, 'agent:recovery', ?, '{"paths":["base.txt"]}', 'mutation', 1, 'expired', ?, ?, NULL, NULL)`, []any{repository.WorktreeID, timestamp, now.Add(-time.Second).Format(time.RFC3339Nano)}},
		{`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, terminal_repository_json, status, started_at, ended_at) VALUES ('RUN-prior', 'ISS-recovery', 'ENV-recovery', 1, 'CLM-prior', 1, 'agent:recovery', ?, NULL, 'active', ?, NULL)`, []any{repositoryJSON, timestamp}},
		{`INSERT INTO run_runtime (run_id, version, executor_id, obligation_id) VALUES ('RUN-prior', 1, 'executor-prior', NULL)`, nil},
		{`INSERT INTO run_observations (id, run_id, kind, actor_id, observed_at, observation_json, created_event_id) VALUES ('OBS-terminated', 'RUN-prior', 'liveness', 'agent:observer', ?, ?, 'EVT-observed')`, []any{timestamp, observationJSON}},
	}
	for index, statement := range statements {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed recovery fixture statement %d: %v", index, err)
		}
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatalf("enable fixture foreign keys: %v", err)
	}
}

func setRecoveryObservationState(t *testing.T, workspace, state string) {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	observation := RunObservation{ID: "OBS-terminated", RunID: "RUN-prior", Kind: "liveness", ActorID: "agent:observer", ObservedAt: now, Liveness: Liveness{ExecutorID: "executor-prior", State: state, SourceRef: "receipt:liveness", ObservedAt: now}}
	if _, err := conn.Exec(`UPDATE run_observations SET observation_json = ? WHERE id = 'OBS-terminated'`, string(mustRunJSON(observation))); err != nil {
		t.Fatalf("update liveness observation: %v", err)
	}
}

func configureForceRecoveryAuthority(t *testing.T, workspace string) {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	if _, err := conn.Exec(`UPDATE activation_actor_profiles SET actor_type = 'human', capabilities_json = '["run.force-recover","run.recover"]' WHERE actor_id = 'agent:recovery'`); err != nil {
		t.Fatalf("update Human profile: %v", err)
	}
	if _, err := conn.Exec(`UPDATE envelope_authority_grants SET actor_selector_json = '{"actor_id":"agent:recovery","required_capabilities":["run.force-recover","run.recover"]}', actions_json = '["run.force-recover","run.recover"]' WHERE envelope_id = 'ENV-recovery' AND ordinal = 0`); err != nil {
		t.Fatalf("update force grant: %v", err)
	}
}
