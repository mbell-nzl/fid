package workflow

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestEquivalentStopsShareAttentionButKeepDistinctFacts(t *testing.T) {
	workspace, _ := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)

	seedStandaloneStopSubject(t, conn)
	first := stopRaisedEvent("", "STOP-1", "ATT-1", true)
	second := stopRaisedEvent("", "STOP-2", "ATT-1", false)
	if err := journal.WithLock(workspace, journal.LockOwner{Command: "test stop raise"}, func(lock *journal.WriteLock) error {
		prepared, err := journal.PrepareEvent(lock, first)
		if err != nil {
			return err
		}
		first = stopRaisedEvent(prepared.EventID, "STOP-1", "ATT-1", true)
		if _, err := journal.CommitEvent(lock, conn, first); err != nil {
			return err
		}
		prepared, err = journal.PrepareEvent(lock, second)
		if err != nil {
			return err
		}
		second = stopRaisedEvent(prepared.EventID, "STOP-2", "ATT-1", false)
		_, err = journal.CommitEvent(lock, conn, second)
		return err
	}); err != nil {
		t.Fatalf("commit equivalent stop events: %v", err)
	}

	var stops, attention, links int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM stops`).Scan(&stops); err != nil {
		t.Fatalf("count stops: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM attention_items WHERE status = 'open'`).Scan(&attention); err != nil {
		t.Fatalf("count open Attention: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM attention_stop_links WHERE attention_id = 'ATT-1'`).Scan(&links); err != nil {
		t.Fatalf("count Attention links: %v", err)
	}
	if stops != 2 || attention != 1 || links != 2 {
		t.Fatalf("stops/Attention/links = %d/%d/%d, want 2/1/2", stops, attention, links)
	}
}

func TestRaiseStopUsesOneOpenAttentionAndRejectsRoutineReadsWithoutEffects(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	installStopTestActor(t, conn)
	request := StopRaiseRequest{
		Subject: Subject{Type: "activation", ID: "ACT-stop-test"}, Code: "operator-stop",
		BlockedActions: []string{"run.recover"}, Details: []byte(`{"reason":"operator intervention","reference":"DEC-188"}`),
		ResolutionCondition: Condition{Type: "activation-recovered", Subject: Subject{Type: "activation", ID: "ACT-stop-test"}}, ResolverCapability: "stop.resolve",
	}
	options := StopOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "agent:stop-test", IdempotencyKey: testIdempotencyKey()}
	firstOptions := options
	first, err := RaiseStop(options, request)
	if err != nil {
		t.Fatalf("RaiseStop(first) error = %v", err)
	}
	options.IdempotencyKey = testIdempotencyKey()
	second, err := RaiseStop(options, request)
	if err != nil {
		t.Fatalf("RaiseStop(second) error = %v", err)
	}
	if first.Stop.ID == second.Stop.ID || first.Attention.ID != second.Attention.ID || len(second.Attention.StopIDs) != 2 {
		t.Fatalf("RaiseStop identities = %#v / %#v, want distinct Stops linked to one Attention", first, second)
	}
	if !strings.HasPrefix(first.Attention.SemanticKey, "sha256:") || len(first.Attention.SemanticKey) != 71 {
		t.Fatal("Attention semantic_key does not use the accepted digest grammar")
	}
	before := readWorkflowTestFile(t, workflowJournalPath(t, workspace))
	replayed, err := RaiseStop(firstOptions, request)
	if err != nil {
		t.Fatal(err)
	}
	expected := first
	expected.Replayed = true
	if !reflect.DeepEqual(replayed, expected) || len(replayed.Attention.StopIDs) != 1 {
		t.Fatalf("replay refreshed historical Attention with later links: got %#v, want %#v", replayed, expected)
	}
	if _, err := GetStop(workspace, migrationsDir, "missing"); err == nil {
		t.Fatal("GetStop(missing) error = nil, want rejected read")
	}
	after := readWorkflowTestFile(t, workflowJournalPath(t, workspace))
	if string(before) != string(after) {
		t.Fatal("rejected Stop read produced a Stop or Attention Event")
	}
}

func TestRoutineStatusLeavesAttentionAndJournalUnchangedUntilExplicitStop(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	installStopTestActor(t, conn)

	beforeJournal := readWorkflowTestFile(t, workflowJournalPath(t, workspace))
	assertOpenAttentionCount(t, conn, 0)
	if _, err := ActivationStatus(workspace, migrationsDir, "ACT-stop-test"); err != nil {
		t.Fatalf("ActivationStatus() error = %v", err)
	}
	assertOpenAttentionCount(t, conn, 0)
	if afterJournal := readWorkflowTestFile(t, workflowJournalPath(t, workspace)); string(afterJournal) != string(beforeJournal) {
		t.Fatal("routine status changed journal bytes")
	}

	result, err := RaiseStop(StopOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "agent:stop-test", IdempotencyKey: testIdempotencyKey()}, stopTestRequest())
	if err != nil {
		t.Fatalf("RaiseStop() error = %v", err)
	}
	if result.Stop.ID == "" || result.Attention.ID == "" {
		t.Fatalf("RaiseStop() = %#v, want Stop and Attention IDs", result)
	}
	assertOpenAttentionCount(t, conn, 1)
	if afterJournal := readWorkflowTestFile(t, workflowJournalPath(t, workspace)); string(afterJournal) == string(beforeJournal) {
		t.Fatal("explicit Stop did not append a journal Event")
	}
}

// TestRaiseStopAcceptsTheClosedFourTypeJudgmentConditionSet proves task 11's
// two remaining registrations, review-recorded and finding-resolved, are
// accepted end to end through RaiseStop's own strict journal write and
// stored verbatim: review-recorded carries no discriminator at all (DEC-236),
// while finding-resolved requires exactly one, and rejects when it is absent.
func TestRaiseStopAcceptsTheClosedFourTypeJudgmentConditionSet(t *testing.T) {
	options, _ := recoveryFixture(t)
	authorizeStopConditionRun(t, options)
	stopOptions := StopOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, ActorID: options.ActorID, IdempotencyKey: testIdempotencyKey()}

	findingID := "FND-condition"
	cases := []struct {
		name      string
		condition Condition
		wantValid bool
	}{
		{"review-recorded bare four-field form is accepted", Condition{Type: "review-recorded", Subject: Subject{Type: "run", ID: "RUN-prior"}}, true},
		{"review-recorded rejects a finding_id", Condition{Type: "review-recorded", Subject: Subject{Type: "run", ID: "RUN-prior"}, FindingID: &findingID}, false},
		{"finding-resolved requires a finding_id", Condition{Type: "finding-resolved", Subject: Subject{Type: "run", ID: "RUN-prior"}}, false},
		{"finding-resolved with a finding_id is accepted", Condition{Type: "finding-resolved", Subject: Subject{Type: "run", ID: "RUN-prior"}, FindingID: &findingID}, true},
		{"finding-resolved accepts an Envelope subject", Condition{Type: "finding-resolved", Subject: Subject{Type: "envelope", ID: "ENV-recovery"}, FindingID: &findingID}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stopOptions.IdempotencyKey = testIdempotencyKey()
			request := stopTestRequest()
			request.Subject = tc.condition.Subject
			request.ResolutionCondition = tc.condition
			raised, err := RaiseStop(stopOptions, request)
			if (err == nil) != tc.wantValid {
				t.Fatalf("RaiseStop with %s condition: %v; want valid=%v", tc.condition.Type, err, tc.wantValid)
			}
			if !tc.wantValid {
				return
			}
			retained, err := GetStop(stopOptions.WorkspaceRoot, stopOptions.MigrationsDir, raised.Stop.ID)
			if err != nil || !reflect.DeepEqual(retained.ResolutionCondition, tc.condition) {
				t.Fatalf("stored Condition round-trip = %#v, %v; want %#v", retained.ResolutionCondition, err, tc.condition)
			}
		})
	}
}

func assertOpenAttentionCount(t *testing.T, conn *sql.DB, want int) {
	t.Helper()
	var got int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM attention_items WHERE status = 'open'`).Scan(&got); err != nil {
		t.Fatalf("count open Attention: %v", err)
	}
	if got != want {
		t.Fatalf("open Attention count = %d, want %d", got, want)
	}
}

func stopTestRequest() StopRaiseRequest {
	return StopRaiseRequest{
		Subject: Subject{Type: "activation", ID: "ACT-stop-test"}, Code: "operator-stop",
		BlockedActions: []string{"run.recover"}, Details: []byte(`{"reason":"operator intervention","reference":"DEC-188"}`),
		ResolutionCondition: Condition{Type: "activation-recovered", Subject: Subject{Type: "activation", ID: "ACT-stop-test"}}, ResolverCapability: "stop.resolve",
	}
}

func installStopTestActor(t *testing.T, conn *sql.DB) {
	t.Helper()
	if _, err := conn.Exec(`INSERT INTO kit_activations (id, kit_identity, backend_identity, phase, status, point_of_no_return, direction, created_at, updated_at, version) VALUES ('ACT-stop-test', 'kit', 'backend', 22, 'finalized', 0, 'activate', '2026-09-05T00:00:00Z', '2026-09-05T00:00:00Z', 1)`); err != nil {
		t.Fatalf("seed activation: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO activation_actor_profiles (activation_id, actor_id, actor_type, producer, provenance_ref, capabilities_json) VALUES ('ACT-stop-test', 'agent:stop-test', 'agent', 'test', 'test:stop', '["attention.resolve","stop.raise","stop.resolve"]')`); err != nil {
		t.Fatalf("seed actor profile: %v", err)
	}
}

func stopRaisedEvent(eventID, stopID, attentionID string, createAttention bool) journal.Event {
	effects := []journal.Effect{{
		Op: "insert", Table: "stops", PK: map[string]any{"id": stopID},
		Values: map[string]any{
			"subject_type": "run", "subject_id": "RUN-1", "code": "operator-stop",
			"blocked_actions_json": `["run.recover"]`, "details_json": `{"reason":"operator intervention","reference":"DEC-188"}`,
			"resolution_condition_json": `{"type":"work-ended","subject":{"type":"run","id":"RUN-1"},"budget_name":null,"capability_id":null}`,
			"resolver_capability":       "stop.resolve", "actor_id": "agent:tl", "raised_at": "2026-09-05T00:00:00Z", "created_event_id": eventID,
		},
	}}
	if createAttention {
		effects = append(effects, journal.Effect{Op: "insert", Table: "attention_items", PK: map[string]any{"id": attentionID}, Values: map[string]any{
			"semantic_key": "sha256:" + strings.Repeat("1", 64), "subject_type": "run", "subject_id": "RUN-1", "code": "operator-stop", "status": "open", "version": 1, "raised_at": "2026-09-05T00:00:00Z", "resolution_event_id": nil,
		}})
	}
	effects = append(effects, journal.Effect{Op: "insert", Table: "attention_stop_links", PK: map[string]any{"attention_id": attentionID, "stop_id": stopID}})
	event := journal.Event{
		EventID: eventID, EventSchema: journal.EventSchemaV2, EventType: "stop.raised",
		ActorV2:   journal.EventActor{ActorID: "agent:tl", ActorType: "agent", Producer: "test", ProvenanceRef: "test"},
		Aggregate: journal.Aggregate{ID: stopID, Type: "stop", Sequence: 1}, CommandName: "stop.raise", Effects: effects,
	}
	event = stopFixtureIdempotency(event)
	condition := map[string]any{"type": "work-ended", "subject": Subject{Type: "run", ID: "RUN-1"}, "budget_name": nil, "capability_id": nil}
	stop := map[string]any{"id": stopID, "subject": Subject{Type: "run", ID: "RUN-1"}, "code": "operator-stop", "blocked_actions": []string{"run.recover"}, "details": map[string]any{"reason": "operator intervention", "reference": "DEC-188"}, "resolution_condition": condition, "resolver_capability": "stop.resolve", "actor_id": "agent:tl", "raised_at": "2026-09-05T00:00:00Z", "resolution": nil}
	stopIDs := []string{stopID}
	if !createAttention {
		stopIDs = []string{"STOP-1", stopID}
	}
	attention := map[string]any{"id": attentionID, "semantic_key": "sha256:" + strings.Repeat("1", 64), "subject": Subject{Type: "run", ID: "RUN-1"}, "code": "operator-stop", "status": "open", "version": 1, "stop_ids": stopIDs, "raised_at": "2026-09-05T00:00:00Z", "resolution_event_id": nil}
	raw, err := json.Marshal(map[string]any{"stop": stop, "attention": attention, "idempotency_key": event.Idempotency.Key, "replayed": false})
	if err != nil {
		panic(err)
	}
	event.Effects[len(event.Effects)-1].Values["result_json"] = string(raw)
	return event
}

func stopFixtureIdempotency(event journal.Event) journal.Event {
	event.Idempotency = &journal.EventIdempotency{Key: event.Aggregate.ID, Scope: event.CommandName + ":test-project", RequestDigest: strings.Repeat("1", 64)}
	event.Effects = append(event.Effects, journal.Effect{Op: "insert", Table: "idempotency_results", PK: map[string]any{"operation_scope": event.Idempotency.Scope, "idempotency_key": event.Idempotency.Key}, Values: map[string]any{"request_digest": event.Idempotency.RequestDigest, "result_entity_id": event.Aggregate.ID, "result_event_id": event.EventID, "aggregate_sequence": event.Aggregate.Sequence, "result_json": `{"replayed":false}`, "created_at": time.Now().UTC().Format(time.RFC3339Nano)}})
	return event
}

func seedStandaloneStopSubject(t *testing.T, conn *sql.DB) {
	t.Helper()
	// Only Stop/Attention projections are under test; the referenced Run is a
	// pre-existing fixture, with its unrelated Envelope/Claim graph out of scope.
	if _, err := conn.Exec(`PRAGMA foreign_keys=OFF`); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO runs (id,issue_id,envelope_id,envelope_revision,claim_id,claim_version,actor_id,start_repository_json,status,started_at) VALUES ('RUN-1','ISS-stop','ENV-stop',1,'CLM-stop',1,'agent:tl','{}','active','2026-09-05T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`PRAGMA foreign_keys=ON`); err != nil {
		t.Fatal(err)
	}
}

// Adapt the complete hand-authored Stop fixture while retaining its response
// snapshot; callers deliberately change the subject/actions before commit.
func refreshStopFixtureResult(event journal.Event) journal.Event {
	last := &event.Effects[len(event.Effects)-1]
	var result map[string]any
	if err := json.Unmarshal([]byte(last.Values["result_json"].(string)), &result); err != nil {
		panic(err)
	}
	stop := result["stop"].(map[string]any)
	subject := map[string]any{"type": event.Effects[0].Values["subject_type"], "id": event.Effects[0].Values["subject_id"]}
	stop["subject"] = subject
	result["attention"].(map[string]any)["subject"] = subject
	for column, field := range map[string]string{"resolution_condition_json": "resolution_condition", "blocked_actions_json": "blocked_actions"} {
		var value any
		if err := json.Unmarshal([]byte(event.Effects[0].Values[column].(string)), &value); err != nil {
			panic(err)
		}
		stop[field] = value
	}
	raw, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	last.Values["result_json"] = string(raw)
	return event
}
