package workflow

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/project"
)

func TestRunCheckpointRecordsCurrentReconciliationAndReplays(t *testing.T) {
	t.Parallel()
	options, request := checkpointFixture(t)
	result, err := RunCheckpoint(options, request)
	if err != nil {
		t.Fatalf("RunCheckpoint() error = %v", err)
	}
	if result.Replayed || result.Checkpoint.Kind != "checkpoint" || result.Checkpoint.RunID != request.RunID || result.Run.Version != 2 || result.Checkpoint.Reconciliation.Digest != request.ReconciliationDigest {
		t.Fatalf("RunCheckpoint() = %#v, want checkpointed Run at version 2", result)
	}
	replayed, err := RunCheckpoint(options, request)
	if err != nil {
		t.Fatalf("RunCheckpoint(replay) error = %v", err)
	}
	if !replayed.Replayed || replayed.Checkpoint.ID != result.Checkpoint.ID || replayed.Run.Version != result.Run.Version {
		t.Fatalf("RunCheckpoint(replay) = %#v, want original checkpoint", replayed)
	}
	changed := request
	changed.ReconciliationDigest = "sha256:changed"
	if _, err := RunCheckpoint(options, changed); err == nil {
		t.Fatal("RunCheckpoint(changed same key) error = nil, want idempotency conflict")
	} else {
		var conflict *journal.IdempotencyConflictError
		if !errors.As(err, &conflict) {
			t.Fatalf("RunCheckpoint(changed same key) error = %v, want idempotency conflict", err)
		}
	}
}

func TestRunCheckpointRejectsPreconditionsAndAuthority(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		prepare   func(*testing.T, RunOptions, *RunCheckpointRequest)
		wantError any
	}{
		{
			name: "stale version",
			prepare: func(_ *testing.T, _ RunOptions, request *RunCheckpointRequest) {
				request.ExpectedVersion++
			},
			wantError: &RunVersionConflictError{},
		},
		{
			name: "inactive Run",
			prepare: func(t *testing.T, options RunOptions, _ *RunCheckpointRequest) {
				checkpointSetRunStatus(t, options.WorkspaceRoot, "failed")
			},
			wantError: &RunInactiveError{},
		},
		{
			name: "authority denied",
			prepare: func(t *testing.T, options RunOptions, _ *RunCheckpointRequest) {
				checkpointSetActorProducer(t, options.WorkspaceRoot, "test")
			},
			wantError: &policy.EnvelopeAuthorityDeniedError{},
		},
		{
			name: "conflicting observation",
			prepare: func(t *testing.T, options RunOptions, request *RunCheckpointRequest) {
				checkpointWrite(t, options.WorkspaceRoot, "base.txt", "changed\n")
			},
			wantError: &RecoveryStateConflictError{},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			options, request := checkpointFixture(t)
			testCase.prepare(t, options, &request)
			_, err := RunCheckpoint(options, request)
			if err == nil {
				t.Fatal("RunCheckpoint() error = nil, want rejection")
			}
			if !errors.As(err, &testCase.wantError) {
				t.Fatalf("RunCheckpoint() error = %T %v, want %T", err, err, testCase.wantError)
			}
		})
	}
}

func checkpointFixture(t *testing.T) (RunOptions, RunCheckpointRequest) {
	t.Helper()
	options, recovery := recoveryFixture(t)
	options.IdempotencyKey = "checkpoint-lost-response"
	checkpointAuthorize(t, options.WorkspaceRoot)
	current, err := project.ObserveRepository(options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("ObserveRepository() error = %v", err)
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	run, err := loadRun(conn, recovery.PriorRunID)
	if err != nil {
		t.Fatalf("load Run: %v", err)
	}
	reconciliation, err := ReconcileRepository(run.StartRepository, current, nil, nil)
	if err != nil {
		t.Fatalf("ReconcileRepository() error = %v", err)
	}
	return options, RunCheckpointRequest{RunID: run.ID, ExpectedVersion: run.Version, ReconciliationDigest: reconciliation.Digest}
}

func checkpointAuthorize(t *testing.T, workspace string) {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	for index, statement := range []string{
		`UPDATE activation_actor_profiles SET producer = 'activated-harness', capabilities_json = '["run.checkpoint","run.recover"]' WHERE actor_id = 'agent:recovery'`,
		`UPDATE envelope_authority_grants SET actor_selector_json = '{"actor_id":"agent:recovery","required_capabilities":["run.checkpoint","run.recover"]}', actions_json = '["run.checkpoint","run.recover"]' WHERE envelope_id = 'ENV-recovery' AND ordinal = 0`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("authorize checkpoint statement %d: %v", index, err)
		}
	}
}

func checkpointSetRunStatus(t *testing.T, workspace, status string) {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	if _, err := conn.Exec(`UPDATE runs SET status = ? WHERE id = 'RUN-prior'`, status); err != nil {
		t.Fatalf("set Run status: %v", err)
	}
}

func checkpointSetActorProducer(t *testing.T, workspace, producer string) {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	if _, err := conn.Exec(`UPDATE activation_actor_profiles SET producer = ? WHERE actor_id = 'agent:recovery'`, producer); err != nil {
		t.Fatalf("set actor producer: %v", err)
	}
}

func checkpointWrite(t *testing.T, workspace, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(workspace, name), []byte(contents), 0o644); err != nil {
		t.Fatalf("write conflicting repository state: %v", err)
	}
}
