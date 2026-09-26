package workflow

import (
	"encoding/json"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestStopResolutionRequiresTerminatedExecutorFact(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"live", "unknown", "terminated", "wrong Run", "missing fact", "changed executor", "changed liveness"} {
		t.Run(name, func(t *testing.T) {
			state := name
			if name != "live" && name != "unknown" {
				state = "terminated"
			}
			options, recovery := recoveryFixture(t)
			authorizeStopConditionRun(t, options)
			setRecoveryObservationState(t, options.WorkspaceRoot, state)
			options.Now = time.Now().UTC()
			observation, err := ObserveRun(options, RunObserveRequest{RunID: "RUN-prior", ExpectedVersion: 1, Liveness: Liveness{ExecutorID: "executor-prior", State: state, SourceRef: "process:condition-matrix", ObservedAt: options.Now.Format(time.RFC3339Nano)}})
			if err != nil {
				t.Fatalf("record %s observation: %v", state, err)
			}
			stopOptions := StopOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, ActorID: options.ActorID, IdempotencyKey: testIdempotencyKey()}
			request := stopTestRequest()
			request.Subject = Subject{Type: "run", ID: "RUN-prior"}
			if name == "wrong Run" {
				options.IdempotencyKey = testIdempotencyKey()
				recovery.ExpectedRunVersion, recovery.ObservationID = observation.Run.Version, observation.Observation.ID
				successor, err := RecoverRun(options, recovery)
				if err != nil {
					t.Fatal(err)
				}
				request.Subject.ID = successor.Run.ID
			}
			request.ResolutionCondition = Condition{Type: "executor-terminated", Subject: request.Subject}
			raised, err := RaiseStop(stopOptions, request)
			if err != nil {
				t.Fatalf("raise executor Stop: %v", err)
			}
			if name == "changed executor" || name == "changed liveness" {
				conn, err := openWorkflowDB(options.WorkspaceRoot)
				if err != nil {
					t.Fatal(err)
				}
				if name == "changed executor" {
					_, err = conn.Exec(`UPDATE run_runtime SET executor_id = 'another-executor' WHERE run_id = 'RUN-prior'`)
				} else {
					_, err = conn.Exec(`UPDATE run_observations SET observation_json = json_set(observation_json, '$.liveness.state', 'unknown') WHERE id = ?`, observation.Observation.ID)
				}
				closeWorkflowDB(conn)
				if err != nil {
					t.Fatal(err)
				}
			}
			factID := observation.EventID
			if name == "missing fact" {
				factID = "EVT-19000101-999999"
			}
			before := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
			stopOptions.IdempotencyKey = testIdempotencyKey()
			resolved, err := ResolveStop(stopOptions, StopResolveRequest{StopID: raised.Stop.ID, FactEventID: factID})
			if name == "terminated" {
				if err != nil || resolved.Stop.Resolution == nil || resolved.Stop.Resolution.FactEventID != observation.EventID {
					t.Fatalf("terminated fact resolution = %#v, %v", resolved, err)
				}
				return
			}
			var invalid *ResolutionInvalidError
			if !errors.As(err, &invalid) {
				t.Fatalf("%s observation resolution = %#v, %v; want resolution-invalid without Effects", state, resolved, err)
			}
			if after := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot)); string(after) != string(before) {
				t.Fatal("nonqualifying fact changed journal")
			}
			retained, err := GetStop(options.WorkspaceRoot, options.MigrationsDir, raised.Stop.ID)
			if err != nil || retained.Resolution != nil {
				t.Fatalf("rejected Stop = %#v, %v; want unresolved", retained, err)
			}
		})
	}
}

func TestStopResolutionChecksCurrentActivationAgreementAndCapability(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		condition  string
		capability string
		mutation   string
		wantValid  bool
	}{
		{"agreed activation", "activation-recovered", "", "", true},
		{"pending component", "activation-recovered", "", `UPDATE activation_resources SET state = 'pending' WHERE component = 'provider'`, false},
		{"mismatched component", "activation-recovered", "", `UPDATE activation_resources SET observed_identity = 'wrong' WHERE component = 'provider'`, false},
		{"missing component", "activation-recovered", "", `DELETE FROM activation_resources WHERE component = 'provider'`, false},
		{"nonterminal activation", "activation-recovered", "", `UPDATE kit_activations SET status = 'committed'`, false},
		{"available capability", "capability-available", "workflow.state.v1", "", true},
		{"unavailable capability", "capability-available", "human.outcomes.v1", "", false},
		{"unknown capability", "capability-available", "invented.v1", "", false},
		{"superseded selection", "capability-available", "workflow.state.v1", `UPDATE kit_activations SET status = 'superseded'`, false},
		{"different backend", "capability-available", "workflow.state.v1", `UPDATE kit_activations SET backend_identity = 'another-backend'`, false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			options, activation, factID := finalizedStopConditionActivation(t)
			conn, err := openWorkflowDB(options.WorkspaceRoot)
			if err != nil {
				t.Fatal(err)
			}
			defer closeWorkflowDB(conn)
			if test.mutation != "" {
				if _, err := conn.Exec(test.mutation); err != nil {
					t.Fatal(err)
				}
			}
			condition := Condition{Type: test.condition, Subject: Subject{Type: "activation", ID: activation.ID}}
			if test.capability != "" {
				condition.CapabilityID = &test.capability
			}
			before := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
			err = qualifyingFact(options.WorkspaceRoot, conn, Stop{ID: "STP-condition", ResolutionCondition: condition}, factID)
			if (err == nil) != test.wantValid {
				t.Fatalf("qualifying %s fact: %v; want valid=%v", test.name, err, test.wantValid)
			}
			if after := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot)); string(after) != string(before) {
				t.Fatal("condition check changed the journal")
			}
		})
	}
}

// This fixture produces a real finalized Event; condition tests vary only the
// current projection to isolate the required present-state revalidation.
func finalizedStopConditionActivation(t *testing.T) (StopOptions, Activation, string) {
	t.Helper()
	workspace, migrations := initializedWorkflowWorkspace(t)
	transition := ActivationTransitionOptions{WorkspaceRoot: workspace, MigrationsDir: migrations, BootstrapActor: "bootstrap:test", IdempotencyKey: testIdempotencyKey()}
	prepared, err := PrepareActivation(ActivationPrepareOptions{WorkspaceRoot: workspace, MigrationsDir: migrations, BootstrapActor: transition.BootstrapActor, IdempotencyKey: testIdempotencyKey(), Request: ActivationPrepareRequest{KitIdentity: "kit:test", BackendIdentity: "backend:test", Phase: 21, Direction: "activate", Resources: activationRecoveryTestResources()}})
	if err != nil {
		t.Fatal(err)
	}
	activation := prepared.Activation
	contents := map[string]json.RawMessage{
		"authority-profile": json.RawMessage(`[{"actor_id":"agent:stop-condition","actor_type":"agent","producer":"activated-harness","provenance_ref":"test:stop-condition","session_ref":null,"capabilities":["attention.resolve","stop.raise","stop.resolve"]}]`),
		"gate-definitions":  json.RawMessage(`[{"definition_id":"gate-condition","definition_version":1,"argv":["/bin/true"],"artifact_bytes":1024,"artifact_paths":[],"credential_names":[],"cwd":".","environment":[],"excerpt_bytes":0,"network":"allowed","output_bytes":1024,"redact_literals":[],"sandbox":"none","shell":false,"success_exit_codes":[0],"timeout_ms":1000,"write":"allowed"}]`),
	}
	digests := []string{}
	for _, kind := range []string{"authority-profile", "gate-definitions"} {
		staged, err := StageActivationComponent(ActivationComponentOptions{WorkspaceRoot: workspace, MigrationsDir: migrations, BootstrapActor: transition.BootstrapActor, IdempotencyKey: testIdempotencyKey(), Request: ActivationComponentRequest{ActivationID: activation.ID, ComponentType: kind, ComponentID: "condition-" + kind, ComponentVersion: 1, Content: contents[kind], ExpectedVersion: activation.Version}})
		if err != nil {
			t.Fatal(err)
		}
		activation = staged.Activation
		digests = append(digests, staged.Component.Digest)
	}
	for _, resource := range activation.Resources {
		transition.IdempotencyKey = testIdempotencyKey()
		identity := resource.NewIdentity
		observed, err := ObserveActivationResource(transition, ActivationResourceObservationRequest{ActivationID: activation.ID, Component: resource.Component, ExpectedActivationVersion: activation.Version, ExpectedVersion: resource.Version, ObservedIdentity: &identity, State: "new", EvidenceRef: "test:condition"})
		if err != nil {
			t.Fatal(err)
		}
		activation = observed.Activation
	}
	sort.Strings(digests)
	transition.IdempotencyKey = testIdempotencyKey()
	committed, err := CommitActivation(transition, ActivationCommitRequest{ActivationID: activation.ID, ExpectedStatus: "staged", ExpectedVersion: activation.Version, ComponentDigests: digests})
	if err != nil {
		t.Fatal(err)
	}
	transition.IdempotencyKey = testIdempotencyKey()
	finalized, err := FinalizeActivation(transition, ActivationFinalizeRequest{ActivationID: activation.ID, ExpectedStatus: "committed", ExpectedVersion: committed.Activation.Version})
	if err != nil {
		t.Fatal(err)
	}
	return StopOptions{WorkspaceRoot: workspace, MigrationsDir: migrations, ActorID: "agent:stop-condition", IdempotencyKey: testIdempotencyKey()}, finalized.Activation, finalized.EventID
}

func authorizeStopConditionRun(t *testing.T, options RunOptions) {
	t.Helper()
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer closeWorkflowDB(conn)
	capabilities := `["attention.resolve","envelope.complete","run.observe","run.observe-budget","run.recover","stop.raise","stop.resolve"]`
	if _, err := conn.Exec(`UPDATE activation_actor_profiles SET producer = 'activated-harness', capabilities_json = ? WHERE actor_id = ?`, capabilities, options.ActorID); err != nil {
		t.Fatal(err)
	}
	selector := `{"actor_id":"agent:recovery","required_capabilities":["attention.resolve","envelope.complete","run.observe","run.observe-budget","run.recover","stop.raise","stop.resolve"]}`
	if _, err := conn.Exec(`UPDATE envelope_authority_grants SET actor_selector_json = ?, actions_json = ? WHERE envelope_id = 'ENV-recovery'`, selector, capabilities); err != nil {
		t.Fatal(err)
	}
}

func TestEnvelopeStopResolutionUsesPrimaryIssueGrantAndCurrentBudget(t *testing.T) {
	t.Parallel()
	options, _ := recoveryFixture(t)
	observeBudgetAuthorize(t, options.WorkspaceRoot, "input-tokens", 5)
	authorizeStopConditionRun(t, options)
	value := int64(3)
	observed, err := ObserveRunBudget(options, RunBudgetObserveRequest{RunID: "RUN-prior", ExpectedVersion: 1, BudgetName: "external", SourceSampleID: "envelope-condition", SampleKind: "baseline", Value: &value, ObservedAt: time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)})
	if err != nil {
		t.Fatal(err)
	}
	stopOptions := StopOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, ActorID: options.ActorID, IdempotencyKey: testIdempotencyKey()}
	request := stopTestRequest()
	request.Subject = Subject{Type: "envelope", ID: "ENV-recovery"}
	name := "external"
	request.ResolutionCondition = Condition{Type: "budget-observed", Subject: request.Subject, BudgetName: &name}
	raised, err := RaiseStop(stopOptions, request)
	if err != nil {
		t.Fatalf("raise authorized Envelope Stop: %v", err)
	}
	stopOptions.IdempotencyKey = testIdempotencyKey()
	resolved, err := ResolveStop(stopOptions, StopResolveRequest{StopID: raised.Stop.ID, FactEventID: observed.EventID})
	if err != nil || resolved.Stop.Resolution == nil {
		t.Fatalf("resolve authorized Envelope Stop: %#v, %v", resolved, err)
	}
}

func TestStopResolutionRecoveryFactNamesThePredecessor(t *testing.T) {
	t.Parallel()
	options, request := recoveryFixture(t)
	recovered, err := RecoverRun(options, request)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer closeWorkflowDB(conn)
	for _, id := range []string{request.PriorRunID, recovered.Run.ID, "RUN-unrelated"} {
		condition := Condition{Type: "run-recovered", Subject: Subject{Type: "run", ID: id}}
		err := qualifyingFact(options.WorkspaceRoot, conn, Stop{ID: "STP-recovery-condition", ResolutionCondition: condition}, recovered.EventID)
		if (err == nil) != (id == request.PriorRunID) {
			t.Errorf("recovery fact for %s: %v; only predecessor %s qualifies", id, err, request.PriorRunID)
		}
	}
}

func TestStopResolutionRequiresNamedKnownNonExhaustedBudget(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"known Run", "known Envelope", "unknown", "exhausted", "wrong budget", "wrong Envelope", "later exhausted"} {
		t.Run(name, func(t *testing.T) {
			options, recovery := recoveryFixture(t)
			observeBudgetAuthorize(t, options.WorkspaceRoot, "input-tokens", 5)
			value := int64(3)
			request := RunBudgetObserveRequest{RunID: recovery.PriorRunID, ExpectedVersion: 1, BudgetName: "external", SourceSampleID: "condition-baseline", SampleKind: "baseline", Value: &value, ObservedAt: time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}
			if name == "unknown" {
				request.Value = nil
			}
			if name == "exhausted" {
				value = 5
			}
			observed, err := ObserveRunBudget(options, request)
			if err != nil {
				t.Fatal(err)
			}
			if name == "later exhausted" {
				options.IdempotencyKey = testIdempotencyKey()
				options.Now = time.Now().UTC()
				value = 2
				request.ExpectedVersion, request.SourceSampleID, request.SampleKind = 2, "condition-later", "delta"
				request.PreviousSampleID = &observed.Sample.ID
				request.ObservedAt = options.Now.Format(time.RFC3339Nano)
				if _, err := ObserveRunBudget(options, request); err != nil {
					t.Fatal(err)
				}
			}
			budgetName := "external"
			condition := Condition{Type: "budget-observed", Subject: Subject{Type: "run", ID: recovery.PriorRunID}, BudgetName: &budgetName}
			if name == "known Envelope" || name == "wrong Envelope" {
				condition.Subject = Subject{Type: "envelope", ID: "ENV-recovery"}
			}
			if name == "wrong Envelope" {
				condition.Subject.ID = "ENV-unrelated"
			}
			if name == "wrong budget" {
				budgetName = "another-budget"
			}
			conn, err := openWorkflowDB(options.WorkspaceRoot)
			if err != nil {
				t.Fatal(err)
			}
			defer closeWorkflowDB(conn)
			err = qualifyingFact(options.WorkspaceRoot, conn, Stop{ID: "STP-budget-condition", ResolutionCondition: condition}, observed.EventID)
			want := name == "known Run" || name == "known Envelope"
			if (err == nil) != want {
				t.Fatalf("budget condition %s: %v; want valid=%v", name, err, want)
			}
		})
	}
}

func TestStopResolutionWorkEndedRequiresCurrentTerminalSubject(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"terminal Run", "terminal Envelope", "active Run", "Envelope has active Run", "wrong Run"} {
		t.Run(name, func(t *testing.T) {
			options, recovery := recoveryFixture(t)
			closeAuthorize(t, options.WorkspaceRoot)
			closed, err := CloseRun(options, RunCloseRequest{RunID: recovery.PriorRunID, ExpectedVersion: 1, ObservationID: recovery.ObservationID, Repository: repositoryForRun(t, options, recovery.PriorRunID)})
			if err != nil {
				t.Fatal(err)
			}
			conn, err := openWorkflowDB(options.WorkspaceRoot)
			if err != nil {
				t.Fatal(err)
			}
			defer closeWorkflowDB(conn)
			condition := Condition{Type: "work-ended", Subject: Subject{Type: "run", ID: recovery.PriorRunID}}
			if name == "terminal Envelope" || name == "Envelope has active Run" {
				condition.Subject = Subject{Type: "envelope", ID: "ENV-recovery"}
			}
			if name == "active Run" || name == "Envelope has active Run" {
				if _, err := conn.Exec(`UPDATE runs SET status = 'active', ended_at = NULL, terminal_repository_json = NULL WHERE id = 'RUN-prior'`); err != nil {
					t.Fatal(err)
				}
			}
			if name == "wrong Run" {
				condition.Subject.ID = "RUN-unrelated"
			}
			err = qualifyingFact(options.WorkspaceRoot, conn, Stop{ID: "STP-ended-condition", ResolutionCondition: condition}, closed.EventID)
			want := name == "terminal Run" || name == "terminal Envelope"
			if (err == nil) != want {
				t.Fatalf("work-ended condition %s: %v; want valid=%v", name, err, want)
			}
		})
	}
}

func TestStopResolutionWorkEndedAcceptsCompletedEnvelopeFact(t *testing.T) {
	t.Parallel()
	options, recovery := recoveryFixture(t)
	closeAuthorize(t, options.WorkspaceRoot)
	if _, err := CloseRun(options, RunCloseRequest{RunID: recovery.PriorRunID, ExpectedVersion: 1, ObservationID: recovery.ObservationID, Repository: repositoryForRun(t, options, recovery.PriorRunID)}); err != nil {
		t.Fatal(err)
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer closeWorkflowDB(conn)
	// Admit the exact closed completed-Envelope Event as a predicate fixture.
	// recoveryFixture intentionally has no Obligation; it cannot stand in for
	// the separate complete-envelope workflow/evidence integration journey.
	var completed journal.Event
	if err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "test completed Envelope fact"}, func(lock *journal.WriteLock) error {
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "envelope.completed", Aggregate: journal.Aggregate{ID: "ENV-recovery", Type: "envelope", Sequence: 1}, CommandName: "envelope.complete", ActorV2: journal.EventActor{ActorID: options.ActorID, ActorType: "agent", Producer: "activated-harness", ProvenanceRef: "test:completed-envelope-condition"}, Effects: []journal.Effect{
			{Op: "update", Table: "work_envelopes", PK: map[string]any{"id": "ENV-recovery"}, Values: map[string]any{"status": "completed"}, Preconditions: map[string]any{"status": "active"}},
			{Op: "update", Table: "envelope_revisions", PK: map[string]any{"envelope_id": "ENV-recovery", "revision": 1}, Values: map[string]any{"status": "completed", "completed_at": time.Now().UTC().Format(time.RFC3339Nano)}, Preconditions: map[string]any{"status": "active"}},
		}}
		completed, err = journal.CommitEvent(lock, conn, event)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	condition := Condition{Type: "work-ended", Subject: Subject{Type: "envelope", ID: "ENV-recovery"}}
	stop := Stop{ID: "STP-envelope-complete", ResolutionCondition: condition}
	if err := qualifyingFact(options.WorkspaceRoot, conn, stop, completed.EventID); err != nil {
		t.Fatalf("completed Envelope fact did not qualify: %v", err)
	}
	if _, err := conn.Exec(`UPDATE runs SET status = 'active', ended_at = NULL, terminal_repository_json = NULL WHERE id = 'RUN-prior'`); err != nil {
		t.Fatal(err)
	}
	if err := qualifyingFact(options.WorkspaceRoot, conn, stop, completed.EventID); err == nil {
		t.Fatal("completed Envelope fact ignored a current active Run")
	}
}
