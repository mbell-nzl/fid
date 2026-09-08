package workflow

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/policy"
)

func TestRecoverRunRequiresExactObservationVersionsAndAuthority(t *testing.T) {
	for _, mode := range []string{"ordinary", "force"} {
		cases := []string{"valid", "unknown", "missing-observation", "nonexistent-observation", "wrong-run", "wrong-executor", "wrong-kind", "stale-run", "stale-claim", "missing-grant", "closed-envelope", "wrong-actor"}
		if mode == "ordinary" {
			cases = append(cases, "live", "active-claim", "released-claim")
		} else {
			cases = append(cases, "terminated", "agent-type", "missing-capability", "missing-reference", "missing-grant-ref", "wrong-envelope-ref", "wrong-revision-ref", "wrong-ordinal", "referenced-wrong-actor", "referenced-wrong-action", "referenced-wrong-scope", "referenced-missing-capability")
		}
		for _, kind := range cases {
			t.Run(mode+"/"+kind, func(t *testing.T) {
				options, request := recoveryFixture(t)
				if mode == "force" {
					request.Authority = exactRecoveryForce(t, options)
				}
				want := mutateRecoveryAuthorityCase(t, options, &request, kind)
				before := recoveryProjectionSnapshot(t, options.WorkspaceRoot)
				beforeJournal := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
				result, err := RecoverRun(options, request)
				if want == "success" {
					mustRecoveryTest(t, err)
					if result.Run.PredecessorRunID == nil || *result.Run.PredecessorRunID != request.PriorRunID {
						t.Fatal("missing predecessor binding")
					}
					if mode == "force" && !reflect.DeepEqual(result.Recovery.Authority, request.Authority) {
						t.Fatal("force authority provenance changed")
					}
					assertRecoveryPredecessorFenced(t, options, request)
					return
				}
				assertRecoveryErrorKind(t, err, want)
				if !reflect.DeepEqual(before, recoveryProjectionSnapshot(t, options.WorkspaceRoot)) || string(beforeJournal) != string(readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))) {
					t.Fatal("rejected recovery changed durable facts")
				}
			})
		}
	}
}

func mutateRecoveryAuthorityCase(t *testing.T, options RunOptions, request *RunRecoveryRequest, kind string) string {
	t.Helper()
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	mustRecoveryTest(t, err)
	defer closeWorkflowDB(conn)
	exec := func(query string, args ...any) { _, err := conn.Exec(query, args...); mustRecoveryTest(t, err) }
	want := "invalid"
	switch kind {
	case "valid":
		return "success"
	case "released-claim":
		exec(`UPDATE claims SET status = 'released', released_at = acquired_at, release_reason = 'operator-released'`)
		return "success"
	case "unknown", "live", "terminated":
		exec(`UPDATE run_observations SET observation_json = json_set(observation_json, '$.liveness.state', ?)`, kind)
		if kind == "unknown" {
			want = "unknown"
		}
		if kind == "live" {
			want = "live"
		}
	case "active-claim":
		exec(`UPDATE claims SET status = 'active'`)
	case "missing-observation":
		request.ObservationID = ""
	case "nonexistent-observation":
		request.ObservationID, want = "OBS-absent", "missing"
	case "wrong-run":
		exec(`UPDATE run_observations SET observation_json = json_set(observation_json, '$.run_id', 'RUN-unrelated')`)
	case "wrong-executor":
		exec(`UPDATE run_observations SET observation_json = json_set(observation_json, '$.liveness.executor_id', 'executor-unrelated')`)
	case "wrong-kind":
		exec(`UPDATE run_observations SET observation_json = json_set(observation_json, '$.kind', 'checkpoint')`)
	case "stale-run":
		request.ExpectedRunVersion, want = 99, "version"
	case "stale-claim":
		request.ExpectedClaimVersion = 99
	case "missing-grant":
		exec(`DELETE FROM envelope_authority_grants`)
		want = "authority"
	case "closed-envelope":
		exec(`UPDATE work_envelopes SET status = 'completed'`)
		want = "authority"
	case "wrong-actor":
		request.Authority.ActorID = "actor-unrelated"
	case "agent-type":
		exec(`UPDATE activation_actor_profiles SET actor_type = 'agent'`)
	case "missing-capability":
		exec(`UPDATE activation_actor_profiles SET capabilities_json = '["run.recover"]'`)
		want = "authority"
	case "missing-reference":
		request.Authority.Reference = nil
	case "missing-grant-ref":
		request.Authority.GrantRef = nil
	case "wrong-envelope-ref":
		request.Authority.GrantRef.EnvelopeID = "ENV-unrelated"
	case "wrong-revision-ref":
		request.Authority.GrantRef.Revision = 99
	case "wrong-ordinal":
		request.Authority.GrantRef.Ordinal = 99
	case "referenced-wrong-actor", "referenced-wrong-action", "referenced-wrong-scope", "referenced-missing-capability":
		// Keep ordinal 0 as a valid grant so broad authorization alone cannot
		// prove the separately referenced ordinal 1 is the exact force grant.
		exec(`INSERT INTO envelope_authority_grants SELECT envelope_id, revision, 1, actor_selector_json, actions_json, scope_json FROM envelope_authority_grants WHERE ordinal = 0`)
		request.Authority.GrantRef.Ordinal = 1
		switch kind {
		case "referenced-wrong-actor":
			exec(`UPDATE envelope_authority_grants SET actor_selector_json = json_set(actor_selector_json, '$.actor_id', 'other') WHERE ordinal = 1`)
		case "referenced-wrong-action":
			exec(`UPDATE envelope_authority_grants SET actions_json = '["run.recover"]' WHERE ordinal = 1`)
		case "referenced-wrong-scope":
			exec(`UPDATE envelope_authority_grants SET scope_json = '{"issue_id":"ISS-other"}' WHERE ordinal = 1`)
		case "referenced-missing-capability":
			exec(`UPDATE envelope_authority_grants SET actor_selector_json = '{"actor_id":"agent:recovery","required_capabilities":["run.recover"]}' WHERE ordinal = 1`)
		}
	}
	return want
}

func assertRecoveryErrorKind(t *testing.T, err error, want string) {
	t.Helper()
	var invalid *RunInvalidError
	var unknown *ExecutorLivenessUnknownError
	var live *ExecutorStillLiveError
	var missing *ObservationNotFoundError
	var version *RunVersionConflictError
	var authority *policy.EnvelopeAuthorityDeniedError
	valid := map[string]bool{"invalid": errors.As(err, &invalid), "unknown": errors.As(err, &unknown), "live": errors.As(err, &live), "missing": errors.As(err, &missing), "version": errors.As(err, &version), "authority": errors.As(err, &authority)}
	if !valid[want] {
		t.Fatalf("error = %v, want %s", err, want)
	}
}

func assertRecoveryPredecessorFenced(t *testing.T, options RunOptions, request RunRecoveryRequest) {
	t.Helper()
	prior, err := GetRun(options.WorkspaceRoot, options.MigrationsDir, request.PriorRunID)
	mustRecoveryTest(t, err)
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	mustRecoveryTest(t, err)
	_, err = conn.Exec(`UPDATE activation_actor_profiles SET producer = 'activated-harness', capabilities_json = '["run.observe"]'; UPDATE envelope_authority_grants SET actor_selector_json = '{"actor_id":"agent:recovery","required_capabilities":["run.observe"]}', actions_json = '["run.observe"]'`)
	mustRecoveryTest(t, err)
	closeWorkflowDB(conn)
	before := recoveryProjectionSnapshot(t, options.WorkspaceRoot)
	journalBefore := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
	options.IdempotencyKey = "fenced-observe"
	_, err = ObserveRun(options, RunObserveRequest{RunID: prior.ID, ExpectedVersion: prior.Version, Liveness: Liveness{ExecutorID: prior.Executor.ID, State: "terminated", SourceRef: "process:fenced", ObservedAt: options.Now.Format(time.RFC3339Nano)}})
	var inactive *RunInactiveError
	if !errors.As(err, &inactive) || inactive.RunID != prior.ID || inactive.Status != "interrupted" {
		t.Fatalf("otherwise valid predecessor observation = %v, want inactive interrupted predecessor", err)
	}
	if !reflect.DeepEqual(before, recoveryProjectionSnapshot(t, options.WorkspaceRoot)) || string(journalBefore) != string(readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))) {
		t.Fatal("fenced mutation changed durable facts")
	}
}

func recoveryProjectionSnapshot(t *testing.T, root string) map[string][]string {
	t.Helper()
	conn, err := openWorkflowDB(root)
	mustRecoveryTest(t, err)
	defer closeWorkflowDB(conn)
	names, err := conn.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	mustRecoveryTest(t, err)
	var tables []string
	for names.Next() {
		var name string
		mustRecoveryTest(t, names.Scan(&name))
		tables = append(tables, name)
	}
	mustRecoveryTest(t, names.Err())
	mustRecoveryTest(t, names.Close())
	result := map[string][]string{}
	for _, table := range tables {
		rows, err := conn.Query(`SELECT * FROM "` + table + `"`)
		mustRecoveryTest(t, err)
		columns, err := rows.Columns()
		mustRecoveryTest(t, err)
		result[table] = []string{}
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			mustRecoveryTest(t, rows.Scan(pointers...))
			encoded, err := json.Marshal(values)
			mustRecoveryTest(t, err)
			result[table] = append(result[table], string(encoded))
		}
		mustRecoveryTest(t, rows.Err())
		mustRecoveryTest(t, rows.Close())
		sort.Strings(result[table])
	}
	return result
}
