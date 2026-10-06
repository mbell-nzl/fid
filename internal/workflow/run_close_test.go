package workflow

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/project"
)

func TestActivationAdmitsRunCloseAndAbandonActions(t *testing.T) {
	t.Parallel()
	if !sortedUniqueCapabilities([]string{"run.abandon", "run.close"}) {
		t.Fatal("sortedUniqueCapabilities() rejected task-8 Run actions")
	}
}

func TestCloseRunDerivesTerminalStateFromTrustedObservation(t *testing.T) {
	t.Parallel()
	options, recovery := recoveryFixture(t)
	closeAuthorize(t, options.WorkspaceRoot)
	request := RunCloseRequest{
		RunID: recovery.PriorRunID, ExpectedVersion: recovery.ExpectedRunVersion,
		ObservationID: recovery.ObservationID, ExpectedRepository: expectedRunRepositoryForTest(t, options.WorkspaceRoot),
	}

	result, err := CloseRun(options, request)
	if err != nil {
		t.Fatalf("CloseRun() error = %v", err)
	}
	if result.Replayed || result.Run.Status != "succeeded" || result.Run.Version != 2 || result.Run.EndedAt == nil {
		t.Fatalf("CloseRun() = %#v, want a derived succeeded terminal Run", result)
	}
	replayed, err := CloseRun(options, request)
	if err != nil || !replayed.Replayed || replayed.Run.ID != result.Run.ID || replayed.EventID != result.EventID {
		t.Fatalf("CloseRun(replay) = (%#v, %v), want original terminal Run", replayed, err)
	}
}

func TestAbandonRunRequiresTerminalProofOrExactHumanForceAuthority(t *testing.T) {
	t.Parallel()
	t.Run("ordinary terminal proof releases the Claim without Evidence", func(t *testing.T) {
		options, recovery := recoveryFixture(t)
		closeAuthorize(t, options.WorkspaceRoot)
		activateRecoveryClaim(t, options.WorkspaceRoot)
		request := RunAbandonRequest{
			RunID: recovery.PriorRunID, ExpectedVersion: recovery.ExpectedRunVersion, ExpectedClaimVersion: recovery.ExpectedClaimVersion,
			ObservationID: recovery.ObservationID, Authority: RecoveryAuthority{Mode: "ordinary", ActorID: options.ActorID, Reason: "operator stopped work"},
			ExpectedRepository: expectedRunRepositoryForTest(t, options.WorkspaceRoot),
		}

		result, err := AbandonRun(options, request)
		if err != nil {
			t.Fatalf("AbandonRun() error = %v", err)
		}
		if result.Replayed || result.Run.Status != "abandoned" || result.Claim.Status != "released" || result.Claim.ReleaseReason == nil || *result.Claim.ReleaseReason != "abandoned" {
			t.Fatalf("AbandonRun() = %#v, want intentional abandonment and released Claim provenance", result)
		}
		if evidenceCount(t, options.WorkspaceRoot) != 0 {
			t.Fatal("AbandonRun() created Evidence")
		}
	})

	t.Run("unknown liveness rejects even force", func(t *testing.T) {
		options, recovery := recoveryFixture(t)
		closeAuthorize(t, options.WorkspaceRoot)
		setRecoveryObservationState(t, options.WorkspaceRoot, "unknown")
		reference := "HUM-1"
		request := RunAbandonRequest{
			RunID: recovery.PriorRunID, ExpectedVersion: recovery.ExpectedRunVersion, ExpectedClaimVersion: recovery.ExpectedClaimVersion,
			ObservationID: recovery.ObservationID, Authority: RecoveryAuthority{Mode: "force", ActorID: options.ActorID, Reason: "operator stopped work", Reference: &reference, GrantRef: &RecoveryGrantRef{EnvelopeID: "ENV-recovery", Revision: 1, Ordinal: 0}},
			ExpectedRepository: expectedRunRepositoryForTest(t, options.WorkspaceRoot),
		}
		_, err := AbandonRun(options, request)
		var unknown *ExecutorLivenessUnknownError
		if !errors.As(err, &unknown) || unknown.ObservationID != request.ObservationID {
			t.Fatalf("AbandonRun() error = %v, want executor-liveness-unknown", err)
		}
	})
}

func closeAuthorize(t *testing.T, workspace string) {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	if _, err := conn.Exec(`UPDATE activation_actor_profiles SET capabilities_json = '["run.abandon","run.close","run.force-recover","run.recover"]' WHERE actor_id = 'agent:recovery'`); err != nil {
		t.Fatalf("authorize profile: %v", err)
	}
	if _, err := conn.Exec(`UPDATE envelope_authority_grants SET actor_selector_json = '{"actor_id":"agent:recovery","required_capabilities":["run.abandon","run.close","run.force-recover","run.recover"]}', actions_json = '["run.abandon","run.close","run.force-recover","run.recover"]' WHERE envelope_id = 'ENV-recovery' AND ordinal = 0`); err != nil {
		t.Fatalf("authorize grant: %v", err)
	}
}

func activateRecoveryClaim(t *testing.T, workspace string) {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	if _, err := conn.Exec(`UPDATE claims SET status = 'active', expires_at = ? WHERE id = 'CLM-prior'`, time.Now().UTC().Add(time.Minute).Format(time.RFC3339Nano)); err != nil {
		t.Fatalf("activate claim: %v", err)
	}
}

func repositoryForRun(t *testing.T, options RunOptions, runID string) Repository {
	t.Helper()
	run, err := GetRun(options.WorkspaceRoot, options.MigrationsDir, runID)
	if err != nil {
		t.Fatalf("GetRun(%q) error = %v", runID, err)
	}
	return run.StartRepository
}

func evidenceCount(t *testing.T, workspace string) int {
	t.Helper()
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM evidence_records`).Scan(&count); err != nil {
		t.Fatalf("count Evidence: %v", err)
	}
	return count
}

func TestP24AbandonRecordsObservedBase(t *testing.T) {
	t.Parallel()
	options, recovery := recoveryFixture(t)
	closeAuthorize(t, options.WorkspaceRoot)
	activateRecoveryClaim(t, options.WorkspaceRoot)
	repository, err := resumeObservedRepository(options.WorkspaceRoot, options.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	repository.BaseOID = nil
	data, err := json.Marshal(map[string]any{"run_id": recovery.PriorRunID, "expected_version": recovery.ExpectedRunVersion, "expected_claim_version": recovery.ExpectedClaimVersion, "observation_id": recovery.ObservationID, "authority": recovery.Authority, "repository": repository, "expected_repository": map[string]any{"head_oid": repository.HeadOID, "dirty": repository.Dirty}})
	if err != nil {
		t.Fatal(err)
	}
	var request RunAbandonRequest
	if err := json.Unmarshal(data, &request); err != nil {
		t.Fatal(err)
	}
	result, err := AbandonRun(options, request)
	if err != nil {
		t.Fatal(err)
	}
	if result.Run.TerminalRepository.BaseOID == nil || *result.Run.TerminalRepository.BaseOID != repository.HeadOID {
		t.Fatalf("terminal base = %v, want HEAD %s", result.Run.TerminalRepository.BaseOID, repository.HeadOID)
	}
}

func expectedRunRepositoryForTest(t *testing.T, root string) project.ExpectedRepository {
	t.Helper()
	observed, err := project.ObserveRepository(root)
	if err != nil {
		t.Fatal(err)
	}
	return project.ExpectedRepository{HeadOID: observed.HeadOID, Dirty: observed.Dirty}
}
