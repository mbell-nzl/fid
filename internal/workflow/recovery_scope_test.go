package workflow

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/project"
)

func TestRecoveryUsesImmutableClaimScopeForOrdinaryAndForceAdoption(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"ordinary-in-scope", "ordinary-outside", "ordinary-forged-scope", "force-outside", "force-protected", "force-stale-digest", "force-omitted-path"} {
		t.Run(kind, func(t *testing.T) {
			options, request := recoveryFixture(t)
			force := kind[:5] == "force"
			if force {
				request.Authority = exactRecoveryForce(t, options)
			}
			name := "outside.txt"
			if kind == "ordinary-in-scope" {
				name = "base.txt"
			}
			if kind == "force-protected" {
				name = ".fid/recovery-protected"
			}
			mustRecoveryTest(t, os.WriteFile(filepath.Join(options.WorkspaceRoot, name), []byte("precise adopted bytes"), 0o644))
			current, err := project.ObserveRepository(options.WorkspaceRoot)
			mustRecoveryTest(t, err)
			prior, err := GetRun(options.WorkspaceRoot, options.MigrationsDir, request.PriorRunID)
			mustRecoveryTest(t, err)
			approved := []string{"base.txt"}
			if kind == "ordinary-forged-scope" {
				approved = []string{name}
			}
			reconciliation, err := ReconcileRepository(prior.StartRepository, current, nil, approved)
			mustRecoveryTest(t, err)
			request.AdoptedPaths, request.ReconciliationDigest = []string{name}, reconciliation.Digest
			if kind == "force-stale-digest" {
				request.ReconciliationDigest = "sha256:stale"
			}
			if kind == "force-omitted-path" {
				request.AdoptedPaths = []string{}
			}
			before := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
			beforeEnvelope, err := GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, prior.EnvelopeID)
			mustRecoveryTest(t, err)
			result, err := RecoverRun(options, request)
			wantSuccess := kind == "ordinary-in-scope" || kind == "force-outside"
			if !wantSuccess {
				var conflict *RecoveryStateConflictError
				if !errors.As(err, &conflict) {
					t.Fatalf("recovery error = %v, want scope conflict", err)
				}
				if string(readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))) != string(before) || recoverySuccessorCount(t, options.WorkspaceRoot) != 0 {
					t.Fatal("rejected adoption changed journal or successor")
				}
				return
			}
			mustRecoveryTest(t, err)
			wantClass := "in-scope"
			if force {
				wantClass = "outside-scope"
			}
			if len(result.Recovery.Reconciliation.Paths) != 1 || result.Recovery.Reconciliation.Paths[0].Classification != wantClass || result.Claim.ScopeJSON != `{"paths":["base.txt"]}` {
				t.Fatalf("adoption relabeled or widened Claim scope: %#v", result)
			}
			afterEnvelope, err := GetEnvelope(options.WorkspaceRoot, options.MigrationsDir, prior.EnvelopeID)
			mustRecoveryTest(t, err)
			if !reflect.DeepEqual(beforeEnvelope, afterEnvelope) {
				t.Fatal("adoption changed immutable Envelope")
			}
		})
	}
}

func TestRecoveryUsesLatestDurableCheckpoint(t *testing.T) {
	t.Parallel()
	options, request := recoveryFixture(t)
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	mustRecoveryTest(t, err)
	_, err = conn.Exec(`UPDATE activation_actor_profiles SET producer = 'activated-harness', capabilities_json = '["run.checkpoint","run.recover"]'; UPDATE envelope_authority_grants SET actor_selector_json = '{"actor_id":"agent:recovery","required_capabilities":["run.checkpoint","run.recover"]}', actions_json = '["run.checkpoint","run.recover"]'`)
	mustRecoveryTest(t, err)
	closeWorkflowDB(conn)
	inspection, err := InspectRun(options.WorkspaceRoot, options.MigrationsDir, request.PriorRunID, options.WorkspaceRoot)
	mustRecoveryTest(t, err)
	checkpointOptions := options
	checkpointOptions.IdempotencyKey = "accepted-checkpoint"
	checkpoint, err := RunCheckpoint(checkpointOptions, RunCheckpointRequest{RunID: request.PriorRunID, ExpectedVersion: request.ExpectedRunVersion, ReconciliationDigest: inspection.Reconciliation.Digest})
	mustRecoveryTest(t, err)
	inspection, err = InspectRun(options.WorkspaceRoot, options.MigrationsDir, request.PriorRunID, options.WorkspaceRoot)
	mustRecoveryTest(t, err)
	request.ExpectedRunVersion, request.ReconciliationDigest = checkpoint.Run.Version, inspection.Reconciliation.Digest
	result, err := RecoverRun(options, request)
	mustRecoveryTest(t, err)
	if result.Recovery.Reconciliation.CheckpointID == nil || *result.Recovery.Reconciliation.CheckpointID != checkpoint.Checkpoint.ID {
		t.Fatal("recovery discarded accepted checkpoint binding")
	}
}

func exactRecoveryForce(t *testing.T, options RunOptions) RecoveryAuthority {
	t.Helper()
	configureForceRecoveryAuthority(t, options.WorkspaceRoot)
	activateRecoveryClaim(t, options.WorkspaceRoot)
	setRecoveryObservationState(t, options.WorkspaceRoot, "live")
	reference := "HUM-exact-takeover"
	return RecoveryAuthority{Mode: "force", ActorID: options.ActorID, Reason: "precise Human takeover", Reference: &reference, GrantRef: &RecoveryGrantRef{EnvelopeID: "ENV-recovery", Revision: 1, Ordinal: 0}}
}

func mustRecoveryTest(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
