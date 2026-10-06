package workflow

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestAbandonRunRequiresExactCurrentAuthorityAndObservation(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"ordinary", "force"} {
		cases := []string{"valid", "unknown", "missing-observation", "nonexistent-observation", "wrong-run", "wrong-executor", "wrong-kind", "stale-run", "stale-claim", "missing-grant", "closed-envelope", "wrong-actor", "expired-claim"}
		if mode == "ordinary" {
			cases = append(cases, "live")
		} else {
			cases = append(cases, "terminated", "agent-type", "missing-capability", "missing-reference", "missing-grant-ref", "wrong-envelope-ref", "wrong-revision-ref", "wrong-ordinal", "referenced-wrong-actor", "referenced-wrong-action", "referenced-wrong-scope", "referenced-missing-capability")
		}
		for _, kind := range cases {
			t.Run(mode+"/"+kind, func(t *testing.T) {
				options, recovery := recoveryFixture(t)
				options.IdempotencyKey = "abandon-matrix"
				if mode == "force" {
					recovery.Authority = exactRecoveryForce(t, options)
				}
				closeAuthorize(t, options.WorkspaceRoot)
				activateRecoveryClaim(t, options.WorkspaceRoot)
				want := "invalid"
				if kind == "expired-claim" {
					conn, err := openWorkflowDB(options.WorkspaceRoot)
					mustRecoveryTest(t, err)
					_, err = conn.Exec(`UPDATE claims SET status='expired'`)
					mustRecoveryTest(t, err)
					closeWorkflowDB(conn)
				} else {
					want = mutateRecoveryAuthorityCase(t, options, &recovery, kind)
				}
				request := RunAbandonRequest{RunID: recovery.PriorRunID, ExpectedVersion: recovery.ExpectedRunVersion, ExpectedClaimVersion: recovery.ExpectedClaimVersion, ObservationID: recovery.ObservationID, Authority: recovery.Authority, ExpectedRepository: expectedRunRepositoryForTest(t, options.WorkspaceRoot)}
				before := recoveryProjectionSnapshot(t, options.WorkspaceRoot)
				history := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
				result, err := AbandonRun(options, request)
				if want != "success" {
					assertRecoveryErrorKind(t, err, want)
					if !reflect.DeepEqual(before, recoveryProjectionSnapshot(t, options.WorkspaceRoot)) || string(history) != string(readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))) {
						t.Fatal("rejected abandon changed durable state")
					}
					return
				}
				mustRecoveryTest(t, err)
				assertAbandonProvenanceAndReplay(t, options, request, result, before)
			})
		}
	}
}

func assertAbandonProvenanceAndReplay(t *testing.T, options RunOptions, request RunAbandonRequest, result RunAbandonResult, before map[string][]string) {
	t.Helper()
	if result.Claim.ScopeJSON != `{"paths":["base.txt"]}` || result.Claim.ID != "CLM-prior" || result.Claim.ActorID != options.ActorID || result.Claim.WorktreeID != repositoryForRun(t, options, request.RunID).WorktreeID {
		t.Fatal("abandon changed immutable Claim identity or scope")
	}
	if result.Replayed || result.Run.Status != "abandoned" || result.Run.Version != request.ExpectedVersion+1 || result.Run.EndedAt == nil || result.Claim.Status != "released" || result.Claim.Version != request.ExpectedClaimVersion+1 || result.Claim.ReleaseReason == nil || *result.Claim.ReleaseReason != "abandoned" || result.Claim.ReleasedAt == nil || *result.Claim.ReleasedAt != *result.Run.EndedAt {
		t.Fatalf("wrong abandonment provenance: %#v", result)
	}
	after := recoveryProjectionSnapshot(t, options.WorkspaceRoot)
	for table, rows := range before {
		switch table {
		case "runs", "run_runtime", "claims", "idempotency_results", "applied_events", "projection_meta":
			continue
		}
		if !reflect.DeepEqual(rows, after[table]) {
			t.Fatalf("abandon changed %s proof/context", table)
		}
	}
	if len(after["runs"]) != len(before["runs"]) || len(after["claims"]) != len(before["claims"]) {
		t.Fatal("abandon created successor identity")
	}
	history := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
	found := false
	for _, line := range strings.Split(strings.TrimSpace(string(history)), "\n") {
		var event journal.Event
		mustRecoveryTest(t, json.Unmarshal([]byte(line), &event))
		if event.EventID == result.EventID {
			found = true
			if event.EventType != "run.abandoned" || event.ActorV2.ActorID != options.ActorID || event.Aggregate.ID != request.RunID {
				t.Fatal("abandon Event lost Actor/Run provenance")
			}
		}
	}
	if !found {
		t.Fatal("abandon Event is not durable")
	}
	replay, err := AbandonRun(options, request)
	mustRecoveryTest(t, err)
	if !replay.Replayed || replay.EventID != result.EventID || !reflect.DeepEqual(replay.Run, result.Run) || !reflect.DeepEqual(mustRunJSON(replay.Claim), mustRunJSON(result.Claim)) {
		t.Fatal("abandon replay changed identity or terminal facts")
	}
	if !reflect.DeepEqual(after, recoveryProjectionSnapshot(t, options.WorkspaceRoot)) || string(history) != string(readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))) {
		t.Fatal("abandon replay changed durable state")
	}
}

func TestCloseRunRejectsUnqualifiedObservationVersionAndAuthority(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"unknown", "live", "missing-observation", "nonexistent-observation", "wrong-run", "wrong-executor", "wrong-kind", "stale-run", "missing-grant", "closed-envelope"} {
		t.Run(kind, func(t *testing.T) {
			options, recovery := recoveryFixture(t)
			closeAuthorize(t, options.WorkspaceRoot)
			want := mutateRecoveryAuthorityCase(t, options, &recovery, kind)
			request := RunCloseRequest{RunID: recovery.PriorRunID, ExpectedVersion: recovery.ExpectedRunVersion, ObservationID: recovery.ObservationID, ExpectedRepository: expectedRunRepositoryForTest(t, options.WorkspaceRoot)}
			before := recoveryProjectionSnapshot(t, options.WorkspaceRoot)
			history := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
			_, err := CloseRun(options, request)
			assertRecoveryErrorKind(t, err, want)
			if !reflect.DeepEqual(before, recoveryProjectionSnapshot(t, options.WorkspaceRoot)) || string(history) != string(readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))) {
				t.Fatal("rejected close changed state")
			}
		})
	}
}
