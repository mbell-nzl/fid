package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func createAcceptanceTask(t *testing.T, f *statusFixture, actor, kind, parent string) string {
	t.Helper()
	args := []string{"--json", "issue", "create", "--kind", kind, "--title", "Discovered work", "--actor", actor, "--idempotency-key", actor + kind + parent}
	if kind == "task" {
		args = append(args, "--lane", "standard")
		if strings.HasPrefix(actor, "agent:") {
			args = append(args, "--priority", "0")
		} else {
			args = append(args, "--priority", "1")
		}
	}
	if parent != "" {
		args = append(args, "--parent", parent)
	}
	fields, _ := lifecycleResult(t, f.run(args...), "issue.create")
	return fields["issue"].(map[string]any)["id"].(string)
}

func TestAcceptanceSelectionAndStatus(t *testing.T) {
	t.Parallel()
	f := newStatusFixture(t, time.Hour)
	f.clock = time.Now
	agent := createAcceptanceTask(t, f, "agent:orchestrator", "task", "")
	human := createAcceptanceTask(t, f, "human:mark", "task", "")
	// The rule also applies to work created before the projection is rebuilt.
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
	next := commandResultFields(t, f, "next", "--explain")
	if next["task"].(map[string]any)["id"] != human {
		t.Fatalf("served unaccepted work: %v", next)
	}
	found := false
	for _, raw := range next["queue"].([]any) {
		entry := raw.(map[string]any)
		if entry["id"] == agent {
			found = entry["reason"] == "not-accepted"
		}
	}
	if !found {
		t.Fatalf("no typed acceptance reason: %v", next)
	}
	status := commandResultFields(t, f, "status")
	if !strings.Contains(fmt.Sprint(status["needs_you"]), agent+", created by agent:orchestrator, waits for you to accept it") || !strings.Contains(fmt.Sprint(status["needs_you"]), "fid accept "+agent) {
		t.Fatal(status)
	}
	before := startEventCount(t, f.workspace)
	denied := f.run("--json", "accept", agent, "--actor", "agent:orchestrator")
	assertJSONFailureCode(t, denied.stderr, "task.accept", "authority-denied")
	if startEventCount(t, f.workspace) != before {
		t.Fatal("agent acceptance changed journal")
	}
	accepted, _ := lifecycleResult(t, f.run("--json", "accept", agent, "--actor", "human:mark"), "task.accept")
	if !strings.Contains(fmt.Sprint(accepted), agent) {
		t.Fatal(accepted)
	}
	next = commandResultFields(t, f, "next")
	if next["task"].(map[string]any)["id"] != agent {
		t.Fatal(next)
	}
	status = commandResultFields(t, f, "status")
	for _, raw := range status["needs_you"].([]any) {
		if raw.(map[string]any)["kind"] == "accept" {
			t.Fatal(status)
		}
	}
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
	if got := commandResultFields(t, f, "next")["task"].(map[string]any)["id"]; got != agent {
		t.Fatal(got)
	}
}

func TestAcceptSeveralIDsAndSubtask(t *testing.T) {
	t.Parallel()
	f := newStatusFixture(t, time.Hour)
	f.clock = time.Now
	task := createAcceptanceTask(t, f, "agent:orchestrator", "task", "")
	subtask := createAcceptanceTask(t, f, "agent:orchestrator", "subtask", task)
	status := commandResultFields(t, f, "status")
	if !strings.Contains(fmt.Sprint(status["needs_you"]), subtask) {
		t.Fatal(status)
	}
	fields, _ := lifecycleResult(t, f.run("--json", "accept", task, subtask, "--actor", "human:mark", "--reason", "Approved discovered work"), "task.accept")
	if len(fields["decisions"].([]any)) != 2 {
		t.Fatal(fields)
	}
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	var count int
	if err := conn.QueryRow(`SELECT count(*) FROM person_decisions WHERE kind='accept' AND actor='human:mark' AND reason='Approved discovered work'`).Scan(&count); err != nil || count != 2 {
		t.Fatalf("decisions=%d/%v", count, err)
	}
	before := startEventCount(t, f.workspace)
	lifecycleResult(t, f.run("--json", "accept", task, subtask, "--actor", "human:mark"), "task.accept")
	if startEventCount(t, f.workspace) != before {
		t.Fatal("accepting accepted work appended decisions")
	}
}

func TestWorkDecisionGoldens(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"txt", "json"} {
		t.Run(format, func(t *testing.T) {
			t.Parallel()
			f := newLifecycleFixture(t)
			nextWork(f)
			f.setTarget()
			startSQL(t, f, `INSERT INTO issue_creators VALUES('FID-9f13aa','agent:orchestrator')`)
			golden := func(name string, args ...string) {
				t.Helper()
				if format == "json" {
					args = append([]string{"--json"}, args...)
				}
				got := f.run(args...)
				if got.err != nil {
					t.Fatal(got)
				}
				lifecycleGolden(t, f, "../work-decisions/"+name+"."+format, got.stdout)
			}
			golden("accept-status", "status")
			golden("accept-next", "next", "--explain")
			golden("accept", "accept", "FID-9f13aa", "--actor", "human:mark", "--reason", "Approved discovered work")
			f = reviewFixture(t)
			golden("allowance", "rework", "allow", "FID-9f13aa", "--rounds", "2", "--actor", "human:mark", "--reason", "Cover both edge cases")
			lifecycleResult(t, recordReview(t, f, changesInput("primary"), "changes", "human:courier"), "review.record")
			golden("allowed-next", "next")
		})
	}
}

func TestWorkDecisionRequiresOpenSubject(t *testing.T) {
	t.Parallel()
	for _, operation := range []string{"accept", "rework"} {
		for _, state := range []string{"done", "cancelled"} {
			t.Run(operation+"/"+state, func(t *testing.T) {
				t.Parallel()
				f := newLifecycleFixture(t)
				nextWork(f)
				setIssueStatus(t, f.workspace, "FID-9f13aa", state)
				args := []string{"--json", "accept", "FID-9f13aa", "--actor", "human:mark"}
				op := "task.accept"
				if operation == "rework" {
					args = []string{"--json", "rework", "allow", "FID-9f13aa", "--actor", "human:mark", "--rounds", "1", "--reason", "Continue"}
					op = "task.rework-allow"
				}
				before := startEventCount(t, f.workspace)
				got := f.run(args...)
				assertJSONFailureCode(t, got.stderr, op, "invalid-argument")
				if startEventCount(t, f.workspace) != before {
					t.Fatal("closed subject wrote a decision")
				}
			})
		}
	}
}

func TestCancelledDiscoveryClearsAcceptanceItem(t *testing.T) {
	t.Parallel()
	f := newStatusFixture(t, time.Hour)
	f.clock = time.Now
	task := createAcceptanceTask(t, f, "agent:orchestrator", "task", "")
	lifecycleResult(t, f.run("--json", "issue", "transition", task, "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed", "--idempotency-key", "cancel"), "issue.transition")
	if got := selectionItem(t, f, "accept"); got != nil {
		t.Fatal("cancelled discovery still needs acceptance", got)
	}
	if got := commandResultFields(t, f, "next"); got["kind"] != "none" {
		t.Fatal(got)
	}
}
