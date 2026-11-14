package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// showDetail is what an issue holds beyond its row: labels, a spec path
// and acceptance criteria.
type showDetail struct {
	id       string
	labels   []string
	spec     string
	criteria []string
}

// detail adds labels, spec paths and acceptance criteria to seeded issues.
// A criterion's ID is AC-<issue ID>-<ordinal>.
func (fixture *statusFixture) detail(details ...showDetail) {
	fixture.t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(fixture.workspace), config.DatabaseFileName))
	if err != nil {
		fixture.t.Fatal(err)
	}
	defer closeCommandDB(fixture.t, conn)
	exec := func(query string, args ...any) {
		if _, err := conn.Exec(query, args...); err != nil {
			fixture.t.Fatalf("seed %v: %v", args, err)
		}
	}
	for _, detail := range details {
		for _, label := range detail.labels {
			exec(`INSERT OR IGNORE INTO labels (name) VALUES (?)`, label)
			exec(`INSERT INTO issue_labels (issue_id, label) VALUES (?, ?)`, detail.id, label)
		}
		if detail.spec != "" {
			exec(`UPDATE issues SET spec_md_path = ? WHERE id = ?`, detail.spec, detail.id)
		}
		for index, text := range detail.criteria {
			exec(`INSERT INTO acceptance_criteria (id, issue_id, ordinal, text) VALUES (?, ?, ?, ?)`,
				fmt.Sprintf("AC-%s-%d", detail.id, index+1), detail.id, index+1, text)
		}
	}
}

func showSubtask(id, title, status string, task *string) readyCommandRow {
	return readyCommandRow{id: id, kind: "subtask", title: title, status: status, parentID: task}
}

// showControls is text fid didn't write that a terminal would act on.
const showControls = "\x1b[2J\nINJECTED"

// showWork is the approved mock-ups' project: an epic with tasks done,
// started, unstarted, blocked and cancelled; a task with subtasks, blocked
// by a done task and blocking two others, one link recorded from each end;
// an epic whose own state is todo but has a started task; a task whose
// title and criteria wrap; and a task holding text fid didn't write.
func showWork(fixture *statusFixture) {
	fixture.t.Helper()
	n, s := intPointer, stringPointer
	ops, loop, export := s("FID-e3a1c0"), s("FID-8a1b2c"), s("FID-71a0c2")
	epic := statusEpic("FID-e3a1c0", "Operating experience", "in_progress", n(1))
	epic.priority = n(1)
	fixture.seed([]readyCommandRow{
		epic,
		statusTask("FID-1a2b3c", "Escape controls in status", "done", ops, n(1), nil),
		statusTask("FID-2b3c4d", "Port the display package", "done", ops, n(2), nil),
		statusTask("FID-71a0c2", "Export validation", "in_progress", ops, n(3), n(1)),
		statusTask("FID-3b9e01", "Tidy resume output", "in_review", ops, n(4), nil),
		statusTask("FID-9f13aa", "Harden login rate limiting", "todo", ops, n(5), nil),
		statusTask("FID-c2d4e5", "Retry a rejected backup push", "todo", ops, n(6), n(2)),
		statusTask("FID-0d4e8b", "Document backup age", "todo", ops, n(7), nil),
		statusTask("FID-4c5d6e", "Drop the V1 resume view", "cancelled", ops, n(8), nil),
		showSubtask("FID-aa01b2", "Check the manifest schema", "done", export),
		showSubtask("FID-aa02b3", "Validate member digests", "in_progress", export),
		showSubtask("FID-aa03b4", "Refuse a truncated archive", "todo", export),
		statusEpic("FID-8a1b2c", "Loop contract", "todo", n(2)),
		statusTask("FID-1f2e3d", "Contract fields and exit codes", "in_progress", loop, nil, n(1)),
		statusTask("FID-2f2e3d", "Golden tests for the contract and every exit code it names", "todo", loop, nil, n(2)),
		statusTask("FID-5d2f77", "Wrap the heading of fid show when an issue's title is far too long for one line", "todo", nil, nil, n(3)),
		statusTask("FID-c0c0c0", "Title"+showControls, "todo", nil, nil, nil),
	}, [][3]string{
		{"FID-71a0c2", "blocked_by", "FID-1a2b3c"},
		{"FID-71a0c2", "blocks", "FID-c2d4e5"}, {"FID-c2d4e5", "blocked_by", "FID-71a0c2"},
		{"FID-0d4e8b", "blocked_by", "FID-71a0c2"},
		{"FID-1f2e3d", "blocks", "FID-2f2e3d"},
	})
	fixture.detail(
		showDetail{"FID-e3a1c0", []string{"ux"}, "docs/specs/07_DESIGN_PLAN_FLOWS.md", nil},
		showDetail{"FID-71a0c2", []string{"backup", "export"}, "docs/specs/02_TRD.md", []string{
			"fid export validate refuses an archive whose member digests don't match its manifest, naming the first bad member",
			"A truncated archive is refused, not read to its end"}},
		showDetail{"FID-c2d4e5", []string{"backup", "reliability"}, "", []string{
			"A push the backup repo rejects is retried once, after fetching and rebasing fid's clone, and the retry is recorded as its own line",
			"A second rejection is a warning in fid status, naming the failure kind"}},
		showDetail{"FID-aa02b3", nil, "", []string{"Each member's SHA-256 is checked against the manifest before import reads it"}},
		showDetail{"FID-5d2f77", []string{"display", "documentation", "golden-tests", "human-output", "status", "wrapping"}, "", []string{
			"A title too long for one line wraps at spaces, its continuation lines indented to where it began, and is never cut",
			"A criterion wraps the same way",
			"A path such as docs/specs/a-very-long-directory-name/another-long-directory/06_DESIGN_PLAN_CLI.md stands whole past the edge"}},
		showDetail{"FID-c0c0c0", []string{"label" + showControls}, "docs/" + showControls + ".md", []string{"Criterion" + showControls}},
	)
}

func TestShowGolden(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, id string
	}{
		{"task", "FID-71a0c2"},
		{"blocked", "FID-c2d4e5"},
		{"epic", "FID-e3a1c0"},
		{"epic-todo-started", "FID-8a1b2c"},
		{"subtask", "FID-aa02b3"},
		{"wrap", "FID-5d2f77"},
		{"controls", "FID-c0c0c0"},
	} {
		for _, format := range []string{"txt", "json"} {
			t.Run(test.name+"."+format, func(t *testing.T) {
				t.Parallel()
				fixture := newStatusFixture(t, 100*time.Hour)
				showWork(fixture)
				args := []string{"show", test.id}
				if format == "json" {
					args = []string{"--json", "show", test.id}
				}
				fixture.golden(filepath.Join("show", test.name+"."+format), args...)
			})
		}
	}
}

// An unknown ID is issue-not-found. The match is exact, so another case, a
// prefix and an acceptance criterion's ID are unknown too; the text names
// where to find IDs.
func TestShowRefusesAnUnknownID(t *testing.T) {
	t.Parallel()
	fixture := newStatusFixture(t, 100*time.Hour)
	showWork(fixture)
	for _, format := range []string{"txt", "json"} {
		args := []string{"show", "FID-nope"}
		if format == "json" {
			args = []string{"--json", "show", "FID-nope"}
		}
		result := fixture.run(args...)
		if cli.ExitCode(result.err) != 3 || result.stdout != "" {
			t.Fatalf("show %v = %#v, want issue-not-found", args, result)
		}
		assertGolden(t, filepath.Join("show", "not-found."+format), result.stderr)
	}
	for _, id := range []string{"fid-71a0c2", "FID-71a0", "AC-FID-71a0c2-1"} {
		result := fixture.run("--json", "show", id)
		if cli.ExitCode(result.err) != 3 || result.stdout != "" {
			t.Fatalf("show %s = %#v, want issue-not-found", id, result)
		}
		assertJSONFailureCode(t, result.stderr, "issue.show", "issue-not-found")
	}
}

func TestShowRefusesWithoutALedger(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	result := executeCommandForTest(workspace, "--json", "show", "FID-71a0c2")
	if cli.ExitCode(result.err) != 3 || result.stdout != "" {
		t.Fatalf("show without a ledger = %#v, want project-not-found", result)
	}
	assertJSONFailureCode(t, result.stderr, "issue.show", "project-not-found")
}

// The JSON's issue is exactly issue.get's.
func TestShowJSONHoldsIssueGetsIssue(t *testing.T) {
	t.Parallel()
	fixture := newStatusFixture(t, 100*time.Hour)
	showWork(fixture)
	issue := func(args ...string) string {
		t.Helper()
		result := fixture.run(args...)
		var envelope struct {
			Result struct {
				Issue json.RawMessage `json:"issue"`
			} `json:"result"`
		}
		if result.err != nil || json.Unmarshal([]byte(result.stdout), &envelope) != nil || len(envelope.Result.Issue) == 0 {
			t.Fatalf("%v = %v:\n%s%s", args, result.err, result.stdout, result.stderr)
		}
		return string(envelope.Result.Issue)
	}
	for _, id := range []string{"FID-71a0c2", "FID-e3a1c0", "FID-aa02b3"} {
		if shown, got := issue("--json", "show", id), issue("--json", "issue", "get", id); shown != got {
			t.Fatalf("show's issue =\n%s\nwant issue.get's\n%s", shown, got)
		}
	}
}
