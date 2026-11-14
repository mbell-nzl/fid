//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Check Effects as well as family names: a generic V1 seed must not hide a
// dependency on a retiring table. Rejected old bytes belong in refusal tests.
func assertLeanIntegrationJournal(t *testing.T, root string) {
	t.Helper()
	forbidden := strings.Fields("kit_activations activation_resources activation_components activation_actor_profiles activation_gate_definitions activation_recoveries work_envelopes envelope_revisions envelope_issue_selection envelope_sources envelope_authority_grants envelope_obligations envelope_budgets envelope_protected_decisions claims runs run_runtime run_budget_bindings run_observations run_recoveries gate_registrations execution_attempts execution_receipts receipt_artifacts evidence_records evidence_artifact_links stops attention_items attention_stop_links stop_resolutions review_records review_findings finding_dispositions concerns deviations decision_refs human_outcomes requirement_refs requirement_issue_links workflow_assignments")
	for _, raw := range bytes.Split(bytes.TrimSpace(p21JournalBytes(t, root)), []byte{'\n'}) {
		var event journal.Event
		if err := json.Unmarshal(raw, &event); err != nil {
			t.Fatal(err)
		}
		family, _, _ := strings.Cut(event.EventType, ".")
		switch family {
		case "activation", "envelope", "claim", "run", "gate", "evidence", "stop", "attention", "review", "finding", "concern", "deviation", "decision", "human-outcome", "requirement", "workflow":
			t.Fatalf("lean fixture appended kernel family %s", event.EventType)
		}
		for _, effect := range event.Effects {
			for _, table := range forbidden {
				if effect.Table == table {
					t.Fatalf("lean fixture appended %s Effect on %s", event.EventType, table)
				}
			}
		}
	}
}

func p24RetainedCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	binary := integrationFIDBinary(t)
	root := gittest.NewWorkspace(t)
	p21Run(t, binary, root, "init", "--prefix", "FID", "--json")
	switch code {
	case "repository-unavailable":
		if err := os.Rename(filepath.Join(root, ".git"), filepath.Join(root, "hidden-git")); err != nil {
			t.Fatal(err)
		}
		return p24CodePair{root, "project.init", []string{"init", "--prefix", "FID"}}
	case "authority-denied":
		return p24CodePair{root, "next.claim", []string{"next", "--claim", "--actor", "service:ci"}}
	case "duplicate-entity", "idempotency-conflict":
		p21Run(t, binary, root, "--json", "label", "add", "paired", "--colour", "red", "--actor", "agent:builder", "--idempotency-key", "pair-first")
		key := "pair-second"
		if code == "idempotency-conflict" {
			key = "pair-first"
		}
		assertLeanIntegrationJournal(t, root)
		return p24CodePair{root, "label.add", []string{"label", "add", "paired", "--colour", "blue", "--actor", "agent:builder", "--idempotency-key", key}}
	case "result-limit-exceeded":
		conn := p22FaultDB(t, root)
		conn.SetMaxOpenConns(1)
		_, err := conn.Exec("PRAGMA foreign_keys=OFF")
		p22FaultOK(t, err)
		_, err = conn.Exec("INSERT INTO labels(name) VALUES ('limit')")
		p22FaultOK(t, err)
		// Deliberately damaged derived state; no accepted kernel record or Event.
		_, err = conn.Exec(`WITH RECURSIVE n(x) AS (VALUES(1) UNION ALL SELECT x+1 FROM n WHERE x<1001) INSERT INTO issue_labels(issue_id,label) SELECT 'missing-'||x,'limit' FROM n`)
		p22FaultOK(t, err)
		p22FaultOK(t, conn.Close())
		return p24CodePair{root, "maintenance.check", []string{"maintenance", "check"}}
	}
	t.Fatalf("no retained code producer for %s", code)
	return p24CodePair{}
}

func TestPortableAssertCorruptEventsRejected(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	source := gittest.NewWorkspace(t)
	p21Run(t, binary, source, "init", "--prefix", "FID", "--json")
	made := p21Run(t, binary, source, "--json", "issue", "create", "--kind", "task", "--title", "Portable start", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "portable-task")
	task := p21PathString(t, made, "result", "issue", "id")
	taken := p21Run(t, binary, source, "--json", "takeover", task, "--actor", "human:mark", "--reason", "Portable fixture")
	hold := p21PathString(t, taken, "result", "hold_id")
	baseline := gittest.NewWorkspace(t)
	copyParityFixture(t, config.LedgerDir(source), config.LedgerDir(baseline))
	started := p21Run(t, binary, source, "--json", "start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "portable-start")
	found, err := journal.FindEvent(filepath.Join(config.LedgerDir(source), "journal"), p21PathString(t, started, "event", "id"))
	p22FaultOK(t, err)
	assertLeanIntegrationJournal(t, source)
	encoded, err := json.Marshal(found.Event)
	p22FaultOK(t, err)
	for _, kind := range []string{"healthy", "unknown-effect-table", "missing-required-effect", "malformed-body", "unknown-family"} {
		t.Run(kind, func(t *testing.T) {
			for _, mode := range []string{"live", "replay"} {
				t.Run(mode, func(t *testing.T) {
					root := gittest.NewWorkspace(t)
					copyParityFixture(t, config.LedgerDir(baseline), config.LedgerDir(root))
					conn := p22FaultDB(t, root)
					defer func() { p22FaultOK(t, conn.Close()) }()
					before := p23PortableAllRows(t, conn)
					journalBefore := p21JournalBytes(t, root)
					var event journal.Event
					p22FaultOK(t, json.Unmarshal(encoded, &event))
					switch kind {
					case "unknown-effect-table":
						event.Effects[0].Table = "unknown_retained_table"
					case "missing-required-effect":
						event.Effects = event.Effects[1:]
					case "malformed-body":
						event.Effects[0].Values["status"] = []any{"invalid scalar"}
					case "unknown-family":
						event.EventType = "unknown.retained"
					}
					hash, err := journal.HashEvent(event)
					p22FaultOK(t, err)
					if mode == "live" {
						err = journal.ProjectEvent(conn, journal.ProjectionInput{Event: event, EventHash: hash, File: "probe.jsonl", Line: 1, Offset: 1})
					} else {
						dir := t.TempDir()
						copyParityFixture(t, filepath.Join(config.LedgerDir(baseline), "journal"), dir)
						name := event.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
						prefix, readErr := os.ReadFile(filepath.Join(dir, name))
						if readErr != nil && !os.IsNotExist(readErr) {
							t.Fatal(readErr)
						}
						raw, marshalErr := json.Marshal(event)
						p22FaultOK(t, marshalErr)
						p22FaultOK(t, os.WriteFile(filepath.Join(dir, name), append(append(prefix, raw...), '\n'), 0o600))
						_, err = journal.ReplayJournal(conn, dir, &journal.ReplayOptions{StartFile: name, StartOffset: int64(len(prefix))})
					}
					if kind == "healthy" {
						if err != nil {
							t.Fatalf("healthy %s control: %v", mode, err)
						}
						return
					}
					if err == nil {
						t.Fatalf("%s %s accepted", kind, mode)
					}
					p23PortableAssertRows(t, kind, before, p23PortableAllRows(t, conn))
					if !bytes.Equal(journalBefore, p21JournalBytes(t, root)) {
						t.Fatal("refusal changed journal bytes")
					}

				})
			}
		})
	}
}
