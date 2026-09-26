//go:build integration

package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

// The input is emitted by real CLI journeys, including conditional recovery,
// shared Attention, and first/upgrade activation. No relabeled positives.
func TestP22CurrentEventFamiliesLiveReplayValidation(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	a := p22CompleteOneIssueRecoveryJourney(t)
	p22AcceptedCall(t, a, "validation-close-observation", "run.observe", nil)
	p22AcceptedCall(t, a, "validation-close", "run.close", nil)
	p22ValidateActualJournal(t, a.workspace, seen)
	b, _ := p22PortableSource(t, "abandoned")
	p22ValidateActualJournal(t, b.workspace, seen)
	c, _ := p22PortableSource(t, "deduplicated-attention")
	p22ValidateActualJournal(t, c.workspace, seen)
	binary := integrationFIDBinary(t)
	for _, direction := range []string{"restore-prior", "complete-forward"} {
		workspace, _, activation := p22ActivationUpgradeFaultFixture(t, binary, direction)
		selected := p21Run(t, binary, workspace, "--json", "activation", "recover", "--input", p21Input(t, workspace, "validation-select", map[string]any{"activation_id": activation.ID, "expected_version": activation.Version, "expected_status": activation.Status, "direction": direction, "reason": "Event validation"}), "--idempotency-key", "validation-select")
		p22ObserveRecoveryResources(t, binary, workspace, activation.ID, direction)
		observed := p21ActivationStatus(t, workspace, activation.ID)
		versions := []any{}
		for _, r := range observed.Resources {
			versions = append(versions, map[string]any{"component": r.Component, "version": r.Version})
		}
		p21Run(t, binary, workspace, "--json", "activation", "recover-complete", "--input", p21Input(t, workspace, "validation-complete", map[string]any{"activation_id": activation.ID, "expected_version": observed.Version, "selected_event_id": p21PathString(t, selected, "event", "id"), "resource_versions": versions}), "--idempotency-key", "validation-complete")
		p22ValidateActualJournal(t, workspace, seen)
	}
	workspace, id := p22ActivationBoundaryFixture(t, binary, "prepared")
	p21Activation(t, binary, workspace, "abort", "validation-abort", map[string]any{"activation_id": id, "expected_status": "prepared", "expected_version": 1, "reason": "Event validation"})
	p22ValidateActualJournal(t, workspace, seen)
	for _, family := range strings.Fields(p22CurrentFamilies) {
		if !seen[family] {
			t.Errorf("missing genuine family %s", family)
		}
	}
}

const p22CurrentFamilies = "activation.prepared activation.component-staged activation.committed activation.finalized activation.aborted activation.resource-observed activation.recovery-selected activation.recovery-completed claim.expired run.started run.observed run.checkpointed run.recovered run.budget-observed run.execution-blocked run.closed run.abandoned gate.execution-started gate.receipt-recorded stop.raised stop.resolved attention.resolved"

func p22ValidateActualJournal(t *testing.T, source string, seen map[string]bool) {
	t.Helper()
	workspace := t.TempDir()
	p22FaultOK(t, os.MkdirAll(filepath.Join(workspace, ".fid"), 0o700))
	p22FaultOK(t, fiddb.MigrateUp(filepath.Join(workspace, ".fid", "fid.sqlite"), filepath.Join(integrationRepositoryRoot(), "migrations")))
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	p22FaultOK(t, err)
	defer conn.Close()
	// Observe committed state through an independent connection throughout this
	// workspace's corpus, without reopening and reparsing its schema per case.
	observer := p22FaultDB(t, workspace)
	defer observer.Close()
	dir := t.TempDir()
	path := filepath.Join(dir, "2026-09-09.jsonl")
	var prefix []byte
	entries, err := os.ReadDir(filepath.Join(source, ".fid", "journal"))
	p22FaultOK(t, err)
	for _, entry := range entries {
		raw, err := os.ReadFile(filepath.Join(source, ".fid", "journal", entry.Name()))
		p22FaultOK(t, err)
		for _, line := range bytes.Split(bytes.TrimSpace(raw), []byte{'\n'}) {
			var event journal.Event
			p22FaultOK(t, json.Unmarshal(line, &event))
			if strings.Contains(" "+p22CurrentFamilies+" ", " "+event.EventType+" ") {
				seen[event.EventType] = true
				t.Run(event.EventType+"/"+event.EventID, func(t *testing.T) { p22RejectEventMutations(t, conn, observer, workspace, path, prefix, event) })
			}
			hash, err := journal.HashEvent(event)
			p22FaultOK(t, err)
			prefix = append(append(prefix, line...), '\n')
			p22FaultOK(t, journal.ProjectEvent(conn, journal.ProjectionInput{Event: event, EventHash: hash, File: filepath.Base(path), Offset: int64(len(prefix))}))
		}
	}
	// Fresh replay of all original bytes must equal the incremental live apply.
	replay := t.TempDir()
	p22FaultOK(t, os.MkdirAll(filepath.Join(replay, ".fid"), 0o700))
	dbPath := filepath.Join(replay, ".fid", "fid.sqlite")
	p22FaultOK(t, fiddb.MigrateUp(dbPath, filepath.Join(integrationRepositoryRoot(), "migrations")))
	replayDB, err := fiddb.Open(dbPath)
	p22FaultOK(t, err)
	defer replayDB.Close()
	p22FaultOK(t, os.WriteFile(path, prefix, 0o600))
	_, err = journal.ReplayJournal(replayDB, dir, nil)
	p22FaultOK(t, err)
	if !reflect.DeepEqual(p22FaultProjection(t, workspace), p22FaultProjection(t, replay)) {
		t.Fatal("actual-family live/replay full projection differs")
	}
}

func p22RejectEventMutations(t *testing.T, conn, observer *sql.DB, workspace, path string, prefix []byte, event journal.Event) {
	t.Helper()
	p22FaultOK(t, journal.DryRunProjectEvent(conn, event))
	p22RejectNestedCurrentValues(t, conn, event)
	p22RejectMissingConditionalFacts(t, conn, event)
	before := p22FaultProjectionFromDB(t, observer)
	mutations := map[string]func(*journal.Event){
		"unrelated-family":      func(e *journal.Event) { e.EventType = "issue.updated"; e.Aggregate.Type = "issue" },
		"missing-domain-effect": func(e *journal.Event) { e.Effects = e.Effects[1:] },
		"reordered-effects": func(e *journal.Event) {
			e.Effects[0], e.Effects[len(e.Effects)-1] = e.Effects[len(e.Effects)-1], e.Effects[0]
		},
		"unknown-column": func(e *journal.Event) { e.Effects[0].Values["unknown_field"] = 1 },
		"wrong-scalar": func(e *journal.Event) {
			for k, v := range e.Effects[0].PK {
				if _, ok := v.(string); ok {
					e.Effects[0].PK[k] = 42
					break
				}
			}
		},
		"wrong-entity-id": func(e *journal.Event) {
			table := map[string]string{"gate.execution-started": "execution_attempts", "gate.receipt-recorded": "execution_receipts", "stop.raised": "stops", "stop.resolved": "stop_resolutions", "attention.resolved": "attention_items"}[e.EventType]
			if table == "" {
				e.Aggregate.ID = "OTHER-aggregate"
				return
			}
			for i := range e.Effects {
				if e.Effects[i].Table == table {
					key := "id"
					if table == "stop_resolutions" {
						key = "stop_id"
					}
					e.Effects[i].PK[key] = "OTHER-entity"
					return
				}
			}
		},
		"wrong-result-event": func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table == "idempotency_results" {
					e.Effects[i].Values["result_event_id"] = "EVT-wrong"
				}
			}
		},
	}

	for i, effect := range event.Effects {
		if effect.Op == "insert" {
			if _, exists := effect.Values["version"]; exists {
				mutations[fmt.Sprintf("initial-version-%d", i)] = func(e *journal.Event) { e.Effects[i].Values["version"] = 2 }
			}
		}
		if effect.Op == "update" {
			mutations[fmt.Sprintf("stale-prior-%d", i)] = func(e *journal.Event) { e.Effects[i].Prior = nil }
		}
	}

	for i, effect := range event.Effects {
		if effect.Table == "stop_resolutions" && event.EventType == "stop.resolved" {
			mutations[fmt.Sprintf("missing-fact-event-%d", i)] = func(e *journal.Event) { e.Effects[i].Values["fact_event_id"] = "EVT-missing" }
		}
		if effect.Table == "stops" {
			mutations[fmt.Sprintf("missing-subject-%d", i)] = func(e *journal.Event) {
				old := e.Effects[i].Values["subject_id"]
				for j := range e.Effects {
					if e.Effects[j].Values["subject_id"] == old {
						e.Effects[j].Values["subject_id"] = "SUBJECT-missing"
					}
				}
			}
		}
	}

	// Missing an optional prior supersession may still be a syntactically valid
	// transition; target a required primary Effect for those activation Events.
	mutations["missing-domain-effect"] = func(e *journal.Event) {
		index := 0
		if e.Effects[0].Table == "kit_activations" && e.Effects[0].PK["id"] != e.Aggregate.ID {
			index = 1
		}
		e.Effects = append(e.Effects[:index], e.Effects[index+1:]...)
	}
	for name, mutate := range mutations {
		for _, mode := range []string{"live", "replay"} {
			t.Run(name+"/"+mode, func(t *testing.T) {
				raw, err := json.Marshal(event)
				p22FaultOK(t, err)
				var changed journal.Event
				p22FaultOK(t, json.Unmarshal(raw, &changed))
				mutate(&changed)
				if mode == "live" {
					hash, hashErr := journal.HashEvent(changed)
					p22FaultOK(t, hashErr)
					err = journal.ProjectEvent(conn, journal.ProjectionInput{Event: changed, EventHash: hash})
				} else {
					raw, err = json.Marshal(changed)
					p22FaultOK(t, err)
					body := append(append(append([]byte{}, prefix...), raw...), '\n')
					p22FaultOK(t, os.WriteFile(path, body, 0o600))
					_, err = journal.ReplayJournal(conn, filepath.Dir(path), &journal.ReplayOptions{StartFile: filepath.Base(path), StartOffset: int64(len(prefix))})
				}
				if err == nil {
					t.Fatalf("%s accepted %s", mode, name)
				}
				if !reflect.DeepEqual(before, p22FaultProjectionFromDB(t, observer)) {
					t.Fatal("invalid Event changed any table")
				}
			})
		}
	}
	// Replay's sequence tracker reads the genuine prefix while skipping its
	// already applied bytes; rejection therefore tests the target Event's gap.
	changed := event
	changed.Aggregate.Sequence += 2
	raw, err := json.Marshal(changed)
	p22FaultOK(t, err)
	body := append(append(append([]byte{}, prefix...), raw...), '\n')
	p22FaultOK(t, os.WriteFile(path, body, 0o600))
	_, err = journal.ReplayJournal(conn, filepath.Dir(path), &journal.ReplayOptions{StartFile: filepath.Base(path), StartOffset: int64(len(prefix))})
	if err == nil || !strings.Contains(err.Error(), fmt.Sprintf("sequence %d", changed.Aggregate.Sequence)) {
		t.Fatalf("sequence gap: %v", err)
	}
	if !reflect.DeepEqual(before, p22FaultProjection(t, workspace)) {
		t.Fatal("sequence gap changed projection")
	}
}
