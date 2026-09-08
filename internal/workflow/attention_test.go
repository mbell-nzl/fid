package workflow

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable(t *testing.T) {
	options, activation, factID := finalizedStopConditionActivation(t)
	workspace, migrationsDir := options.WorkspaceRoot, options.MigrationsDir
	request := stopTestRequest()
	request.Subject.ID = activation.ID
	request.ResolutionCondition.Subject.ID = activation.ID
	first, err := RaiseStop(options, request)
	if err != nil {
		t.Fatalf("RaiseStop(first) error = %v", err)
	}
	options.IdempotencyKey = testIdempotencyKey()
	second, err := RaiseStop(options, request)
	if err != nil {
		t.Fatalf("RaiseStop(second) error = %v", err)
	}
	options.IdempotencyKey = testIdempotencyKey()
	firstResolution, err := ResolveStop(options, StopResolveRequest{StopID: first.Stop.ID, FactEventID: factID})
	if err != nil {
		t.Fatalf("ResolveStop(first) error = %v", err)
	}
	replayedStop, err := ResolveStop(options, StopResolveRequest{StopID: first.Stop.ID, FactEventID: factID})
	if err != nil || !replayedStop.Replayed || replayedStop.EventID != firstResolution.EventID || replayedStop.Stop.ID != firstResolution.Stop.ID {
		t.Fatalf("ResolveStop(replay) = (%#v, %v), want original resolution", replayedStop, err)
	}
	options.IdempotencyKey = testIdempotencyKey()
	secondResolution, err := ResolveStop(options, StopResolveRequest{StopID: second.Stop.ID, FactEventID: factID})
	if err != nil {
		t.Fatalf("ResolveStop(second) error = %v", err)
	}
	ids := []string{firstResolution.Stop.Resolution.EventID, secondResolution.Stop.Resolution.EventID}
	sort.Strings(ids)
	before := readWorkflowTestFile(t, workflowJournalPath(t, workspace))
	var unrelated struct {
		ID string `json:"event_id"`
	}
	if err := json.NewDecoder(strings.NewReader(string(before))).Decode(&unrelated); err != nil || unrelated.ID == "" {
		t.Fatalf("load unrelated project Event identity: %v", err)
	}
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer closeWorkflowDB(conn)
	var beforeIdempotency int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM idempotency_results`).Scan(&beforeIdempotency); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []struct {
		name string
		ids  []string
	}{
		{"missing disposition", []string{ids[0]}},
		{"extra reference", []string{ids[0], ids[1], factID}},
		{"duplicate reference", []string{ids[0], ids[0], ids[1]}},
		{"wrong project subject", []string{ids[0], unrelated.ID}},
		{"underlying fact only", []string{factID}},
	} {
		sort.Strings(invalid.ids)
		options.IdempotencyKey = testIdempotencyKey()
		if _, err := ResolveAttention(options, AttentionResolveRequest{AttentionID: first.Attention.ID, ExpectedVersion: 1, FactEventIDs: invalid.ids}); err == nil {
			t.Fatalf("ResolveAttention(%s) error = nil", invalid.name)
		}
		if after := readWorkflowTestFile(t, workflowJournalPath(t, workspace)); string(after) != string(before) {
			t.Fatalf("rejected %s changed journal bytes", invalid.name)
		}
		attention, err := loadAttentionForStop(conn, first.Stop.ID)
		if err != nil || attention.Status != "open" || attention.Version != 1 {
			t.Fatalf("rejected %s changed Attention: %#v, %v", invalid.name, attention, err)
		}
		var dispositions, idempotency int
		if err := conn.QueryRow(`SELECT (SELECT COUNT(*) FROM stop_resolutions), (SELECT COUNT(*) FROM idempotency_results)`).Scan(&dispositions, &idempotency); err != nil || dispositions != 2 || idempotency != beforeIdempotency {
			t.Fatalf("rejected %s changed disposition/idempotency state: %d/%d, %v", invalid.name, dispositions, idempotency, err)
		}
	}
	options.IdempotencyKey = testIdempotencyKey()
	resolved, err := ResolveAttention(options, AttentionResolveRequest{AttentionID: first.Attention.ID, ExpectedVersion: 1, FactEventIDs: ids})
	if err != nil {
		t.Fatalf("ResolveAttention(exact dispositions) error = %v", err)
	}
	if resolved.Attention.Status != "resolved" || resolved.Attention.Version != 2 {
		t.Fatalf("resolved Attention = %#v, want resolved version 2", resolved.Attention)
	}
	replayedAttention, err := ResolveAttention(options, AttentionResolveRequest{AttentionID: first.Attention.ID, ExpectedVersion: 1, FactEventIDs: ids})
	if err != nil || !replayedAttention.Replayed || replayedAttention.EventID != resolved.EventID || replayedAttention.Attention.ID != resolved.Attention.ID {
		t.Fatalf("ResolveAttention(replay) = (%#v, %v), want original resolution", replayedAttention, err)
	}
	for _, id := range []string{first.Stop.ID, second.Stop.ID} {
		stop, err := GetStop(workspace, migrationsDir, id)
		if err != nil || stop.Resolution == nil || stop.Code != "operator-stop" || string(stop.Details) != `{"reason":"operator intervention","reference":"DEC-188"}` {
			t.Fatalf("immutable resolved Stop %q = %#v, %v", id, stop, err)
		}
	}
}

func TestAttentionResolutionRevalidatesCurrentBudget(t *testing.T) {
	for _, exhaust := range []bool{false, true} {
		t.Run(map[bool]string{false: "known budget", true: "later exhausted budget"}[exhaust], func(t *testing.T) {
			options, _ := recoveryFixture(t)
			observeBudgetAuthorize(t, options.WorkspaceRoot, "input-tokens", 5)
			authorizeStopConditionRun(t, options)
			value := int64(3)
			sampleRequest := RunBudgetObserveRequest{RunID: "RUN-prior", ExpectedVersion: 1, BudgetName: "external", SourceSampleID: "attention-baseline", SampleKind: "baseline", Value: &value, ObservedAt: time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}
			observed, err := ObserveRunBudget(options, sampleRequest)
			if err != nil {
				t.Fatal(err)
			}
			stopOptions := StopOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, ActorID: options.ActorID, IdempotencyKey: testIdempotencyKey()}
			request := stopTestRequest()
			request.Subject = Subject{Type: "run", ID: "RUN-prior"}
			name := "external"
			request.ResolutionCondition = Condition{Type: "budget-observed", Subject: request.Subject, BudgetName: &name}
			raised, err := RaiseStop(stopOptions, request)
			if err != nil {
				t.Fatal(err)
			}
			stopOptions.IdempotencyKey = testIdempotencyKey()
			disposed, err := ResolveStop(stopOptions, StopResolveRequest{StopID: raised.Stop.ID, FactEventID: observed.EventID})
			if err != nil {
				t.Fatal(err)
			}
			if exhaust {
				options.IdempotencyKey = testIdempotencyKey()
				options.Now = time.Now().UTC()
				value = 2
				sampleRequest.ExpectedVersion, sampleRequest.SourceSampleID, sampleRequest.SampleKind = 2, "attention-later", "delta"
				sampleRequest.PreviousSampleID = &observed.Sample.ID
				sampleRequest.ObservedAt = options.Now.Format(time.RFC3339Nano)
				if _, err := ObserveRunBudget(options, sampleRequest); err != nil {
					t.Fatal(err)
				}
			}
			before := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
			stopOptions.IdempotencyKey = testIdempotencyKey()
			resolved, err := ResolveAttention(stopOptions, AttentionResolveRequest{AttentionID: raised.Attention.ID, ExpectedVersion: 1, FactEventIDs: []string{disposed.EventID}})
			if !exhaust {
				if err != nil || resolved.Attention.Status != "resolved" {
					t.Fatalf("known current budget resolution = %#v, %v", resolved, err)
				}
				return
			}
			var rejected *ResolutionInvalidError
			if !errors.As(err, &rejected) {
				t.Fatalf("later exhausted budget resolution = %#v, %v; want resolution-invalid", resolved, err)
			}
			if after := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot)); string(after) != string(before) {
				t.Fatal("stale underlying budget resolution changed journal")
			}
		})
	}
}

func TestAttentionResolutionDeduplicatesContainingEventAndRevalidatesEveryStop(t *testing.T) {
	for _, invalid := range []int{-1, 0, 1} {
		t.Run([]string{"qualifying shared Event", "invalid first fact", "invalid second fact"}[invalid+1], func(t *testing.T) {
			options, attention, factID, wrongFactID := sharedActivationDispositionFixture(t)
			conn, err := openWorkflowDB(options.WorkspaceRoot)
			if err != nil {
				t.Fatal(err)
			}
			defer closeWorkflowDB(conn)
			if invalid >= 0 {
				if _, err := conn.Exec(`UPDATE stop_resolutions SET fact_event_id = ? WHERE stop_id = ?`, wrongFactID, attention.StopIDs[invalid]); err != nil {
					t.Fatal(err)
				}
			}
			stops := make([]Stop, 0, 2)
			for _, id := range attention.StopIDs {
				stop, err := loadStop(conn, id)
				if err != nil {
					t.Fatal(err)
				}
				stops = append(stops, stop)
			}
			before := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot))
			request := AttentionResolveRequest{AttentionID: attention.ID, ExpectedVersion: 1, FactEventIDs: []string{factID}}
			resolved, err := ResolveAttention(options, request)
			if invalid < 0 {
				if err != nil || resolved.Attention.Status != "resolved" || resolved.Attention.Version != 2 {
					t.Fatalf("one containing Event for two Stops = %#v, %v; want resolved", resolved, err)
				}
				replayed, err := ResolveAttention(options, request)
				if err != nil || !replayed.Replayed || replayed.EventID != resolved.EventID {
					t.Fatalf("shared disposition replay = %#v, %v", replayed, err)
				}
			} else {
				var rejected *ResolutionInvalidError
				if !errors.As(err, &rejected) {
					t.Fatalf("invalid underlying fact on Stop %d: %v; want resolution-invalid", invalid, err)
				}
				if after := readWorkflowTestFile(t, workflowJournalPath(t, options.WorkspaceRoot)); string(after) != string(before) {
					t.Fatal("rejected shared disposition changed journal")
				}
			}
			for _, original := range stops {
				retained, err := loadStop(conn, original.ID)
				if err != nil || !reflect.DeepEqual(retained, original) {
					t.Fatalf("Attention resolution changed Stop/disposition: %#v, %v", retained, err)
				}
			}
		})
	}
}

func sharedActivationDispositionFixture(t *testing.T) (StopOptions, Attention, string, string) {
	t.Helper()
	options, prior, _ := finalizedStopConditionActivation(t)
	resources := make([]ActivationResourceInput, 0, len(prior.Resources))
	for _, resource := range prior.Resources {
		old := resource.NewIdentity
		resources = append(resources, ActivationResourceInput{Component: resource.Component, OldIdentity: &old, NewIdentity: old + "-upgrade"})
	}
	prepared, err := PrepareActivation(ActivationPrepareOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: testIdempotencyKey(), Request: ActivationPrepareRequest{KitIdentity: prior.KitIdentity + "-upgrade", BackendIdentity: prior.BackendIdentity + "-upgrade", Phase: 21, Direction: "activate", ExpectedPriorActivationID: &prior.ID, ExpectedPriorVersion: &prior.Version, Resources: resources}})
	if err != nil {
		t.Fatal(err)
	}
	transition := ActivationTransitionOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, BootstrapActor: "bootstrap:test", IdempotencyKey: testIdempotencyKey()}
	selected, err := RecoverActivation(transition, ActivationRecoveryRequest{ActivationID: prepared.Activation.ID, ExpectedVersion: prepared.Activation.Version, ExpectedStatus: "prepared", Direction: "restore-prior", Reason: "shared Stop disposition fixture"})
	if err != nil {
		t.Fatal(err)
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer closeWorkflowDB(conn)
	var firstID string
	if err := conn.QueryRow(`SELECT id FROM stops WHERE subject_id = ?`, prepared.Activation.ID).Scan(&firstID); err != nil {
		t.Fatal(err)
	}
	stop, err := loadStop(conn, firstID)
	if err != nil {
		t.Fatal(err)
	}
	attention, err := loadAttentionForStop(conn, firstID)
	if err != nil {
		t.Fatal(err)
	}
	// Seed a second equivalent immutable fact through the real Stop Event
	// producer; the actual activation completion must dispose both atomically.
	if err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "test equivalent activation Stop"}, func(lock *journal.WriteLock) error {
		stop.ID, err = newStopID(conn)
		if err != nil {
			return err
		}
		stop.RaisedAt = time.Now().UTC().Format(time.RFC3339Nano)
		event := stopRaisedEventForDomain(stop, attention, false, "")
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		_, err = journal.CommitEvent(lock, conn, stopRaisedEventForDomain(stop, attention, false, prepared.EventID))
		return err
	}); err != nil {
		t.Fatal(err)
	}
	activation := observeActivationResourcesForRecoveryTest(t, options.WorkspaceRoot, options.MigrationsDir, selected.Activation, "old")
	transition.IdempotencyKey = testIdempotencyKey()
	completed, err := CompleteActivationRecovery(transition, ActivationRecoveryCompleteRequest{ActivationID: activation.ID, ExpectedVersion: activation.Version, SelectedEventID: selected.EventID, ResourceVersions: activationResourceVersions(activation.Resources)})
	if err != nil {
		t.Fatal(err)
	}
	attention, err = loadAttentionForStop(conn, firstID)
	if err != nil || len(attention.StopIDs) != 2 {
		t.Fatalf("shared Attention = %#v, %v; want two linked Stops", attention, err)
	}
	options.IdempotencyKey = testIdempotencyKey()
	return options, attention, completed.EventID, selected.EventID
}
