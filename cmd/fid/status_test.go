package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// statusNow is when every status golden is taken, shown in NZDT, 13 hours
// ahead of UTC, so local times and UTC differ.
var statusNow = time.Date(2026, 11, 8, 9, 0, 0, 0, time.UTC).In(time.FixedZone("NZDT", 13*60*60))

// statusFixture is a project whose genesis Event is backdated to a fixed
// time, with a home and user config of its own. Its commands see stdout as
// a terminal of screen's size when screen is set, environ besides HOME
// and XDG_CONFIG_HOME, and the time clock gives, or statusNow.
type statusFixture struct {
	t          *testing.T
	workspace  string
	configHome string
	home       string
	screen     *display.Screen
	environ    []string
	clock      func() time.Time
}

func newStatusFixture(t *testing.T, started time.Duration) *statusFixture {
	t.Helper()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init", "--json")
	backdateGenesis(t, workspace, statusNow.Add(-started))
	appendPolicy(t, workspace, configuredStatusPolicy)
	policyOrigin(t, workspace)
	base := t.TempDir()
	return &statusFixture{t: t, workspace: workspace, configHome: filepath.Join(base, "config"), home: filepath.Join(base, "home")}
}

// backdateGenesis moves the journal's only Event, fid init's genesis Event,
// to at, and rebuilds the projection from it.
func backdateGenesis(t *testing.T, workspace string, at time.Time) {
	t.Helper()
	ledger := config.LedgerDir(workspace)
	journalDir := filepath.Join(ledger, config.JournalDirName)
	entries, err := os.ReadDir(journalDir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("journal = %v, %v; want one file", entries, err)
	}
	original := filepath.Join(journalDir, entries[0].Name())
	payload, err := os.ReadFile(original)
	if err != nil {
		t.Fatal(err)
	}
	line := regexp.MustCompile(`"occurred_at":"[^"]*"`).ReplaceAllString(string(payload), `"occurred_at":"`+at.UTC().Format(time.RFC3339Nano)+`"`)
	line = regexp.MustCompile(`"event_id":"EVT-\d{8}-`).ReplaceAllString(line, `"event_id":"EVT-`+at.UTC().Format("20060102")+`-`)
	if err := os.Remove(original); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(journalDir, at.UTC().Format("2006-01-02")+".jsonl"), []byte(line), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{config.DatabaseFileName, config.DatabaseFileName + "-wal", config.DatabaseFileName + "-shm"} {
		if err := os.Remove(filepath.Join(ledger, name)); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
	if result := executeCommandForTest(workspace, "rebuild"); result.err != nil {
		t.Fatalf("rebuild = %v: %s", result.err, result.stderr)
	}
}

func (fixture *statusFixture) setTarget() {
	fixture.t.Helper()
	if err := config.SetBackupTarget(filepath.Join(fixture.configHome, "fid", "config.toml"), "/srv/backups.git"); err != nil {
		fixture.t.Fatal(err)
	}
}

func (fixture *statusFixture) writeUserConfig(contents string) {
	fixture.t.Helper()
	path := filepath.Join(fixture.configHome, "fid", "config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		fixture.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		fixture.t.Fatal(err)
	}
}

// statusRecordLine is a backup record line finished ago before statusNow.
type statusRecordLine struct {
	outcome string
	ago     time.Duration
	kind    string
}

func (fixture *statusFixture) writeRecord(lines ...statusRecordLine) {
	fixture.t.Helper()
	fixture.writeRecordAt(statusNow, lines...)
}

// writeRecordAt writes lines finished ago before now, not statusNow.
func (fixture *statusFixture) writeRecordAt(now time.Time, lines ...statusRecordLine) {
	fixture.t.Helper()
	var record strings.Builder
	for _, line := range lines {
		finished := now.Add(-line.ago).UTC()
		entry := map[string]any{"schema": "fid.backup-record.v1", "started_at": finished.Add(-4 * time.Second).Format(time.RFC3339Nano),
			"finished_at": finished.Format(time.RFC3339Nano), "outcome": line.outcome, "project_id": "project-0123456789abcdef0123456789abcdef",
			"commit": nil, "archive_digest": nil, "event_count": nil, "terminal_event_hash": nil, "failure": nil}
		if line.outcome == "failed" {
			entry["failure"] = map[string]any{"kind": line.kind, "message": "the backup failed"}
		} else {
			entry["commit"], entry["archive_digest"] = strings.Repeat("c", 40), "sha256:"+strings.Repeat("d", 64)
			entry["event_count"], entry["terminal_event_hash"] = 12, "sha256:"+strings.Repeat("e", 64)
		}
		encoded, err := json.Marshal(entry)
		if err != nil {
			fixture.t.Fatal(err)
		}
		record.Write(append(encoded, '\n'))
	}
	path := filepath.Join(config.LedgerDir(fixture.workspace), "backups.jsonl")
	if err := os.WriteFile(path, []byte(record.String()), 0o644); err != nil {
		fixture.t.Fatal(err)
	}
}

// seed adds issues to the projection, epics' scheduling chosen from their
// tasks, and links between them.
func (fixture *statusFixture) seed(rows []readyCommandRow, links [][3]string) {
	fixture.t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(fixture.workspace), config.DatabaseFileName))
	if err != nil {
		fixture.t.Fatal(err)
	}
	defer closeCommandDB(fixture.t, conn)
	for index := range rows {
		row := &rows[index]
		row.workflowID, row.ownerRole = row.kind+"-default", "developer"
		created := statusNow.Add(-time.Duration(200-index) * time.Hour).UTC().Format(time.RFC3339)
		row.createdAt, row.updatedAt = created, created
		if row.kind == "epic" {
			row.scheduling = stringPointer("prioritised")
			for _, task := range rows {
				if task.parentID != nil && *task.parentID == row.id && task.sequence != nil {
					row.scheduling = stringPointer("sequential")
				}
			}
		}
	}
	seedReadyCommandRows(fixture.t, conn, rows)
	for _, link := range links {
		if _, err := conn.Exec(`INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at) VALUES ('issue', ?, ?, 'issue', ?, ?)`,
			link[0], link[1], link[2], statusNow.UTC().Format(time.RFC3339)); err != nil {
			fixture.t.Fatalf("seed link %v: %v", link, err)
		}
	}
}

// run runs fid with args in the fixture's environment.
func (fixture *statusFixture) run(args ...string) commandExecutionResult {
	root := newRootCommandAt(fixture.workspace)
	now := fixture.clock
	if now == nil {
		now = func() time.Time { return statusNow }
	}
	environment := displayEnvironment{environ: append([]string{"HOME=" + fixture.home, "XDG_CONFIG_HOME=" + fixture.configHome}, fixture.environ...),
		now: now, screen: func(io.Writer) *display.Screen { return fixture.screen }}
	root.SetContext(context.WithValue(root.Context(), displayEnvironmentKey{}, environment))
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs(args)
	err := root.Execute()
	return commandExecutionResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
}

// golden runs fid and compares its output, with the fixture's temp paths
// named <ROOT> and <CONFIG>, to testdata/<name>.
func (fixture *statusFixture) golden(name string, args ...string) {
	fixture.t.Helper()
	result := fixture.run(args...)
	if result.err != nil || result.stderr != "" {
		fixture.t.Fatalf("status %v = %v, stderr %q", args, result.err, result.stderr)
	}
	assertGolden(fixture.t, name, fixture.placeholders(result.stdout))
}

// placeholders names the fixture's temp paths in output <ROOT> and
// <CONFIG>.
func (fixture *statusFixture) placeholders(output string) string {
	fixture.t.Helper()
	for _, path := range []string{fixture.workspace, fixture.configHome} {
		placeholder := "<ROOT>"
		if path == fixture.configHome {
			placeholder = "<CONFIG>"
		}
		output = strings.ReplaceAll(output, gittest.ResolvedPath(fixture.t, filepath.Dir(path))+"/"+filepath.Base(path), placeholder)
		output = strings.ReplaceAll(output, path, placeholder)
	}
	return output
}

func statusTask(id, title, status string, parent *string, sequence, priority *int) readyCommandRow {
	return readyCommandRow{id: id, kind: "task", title: title, status: status, parentID: parent, sequence: sequence, priority: priority}
}

func statusEpic(id, title, status string, sequence *int) readyCommandRow {
	return readyCommandRow{id: id, kind: "epic", title: title, status: status, sequence: sequence}
}

// statusWork is a project with in-progress and to-do epics, started,
// unstarted, blocked and standalone tasks, and twelve to-do task rows.
func statusWork() ([]readyCommandRow, [][3]string) {
	n := intPointer
	ops, loop, closed := stringPointer("FID-e3a1c0"), stringPointer("FID-8a1b2c"), stringPointer("FID-d0d0d0")
	rows := []readyCommandRow{
		statusEpic("FID-e3a1c0", "Operating experience", "in_progress", n(1)),
		statusTask("FID-71a0c2", "Export validation", "in_progress", ops, n(1), nil),
		statusTask("FID-3b9e01", "Tidy resume output", "in_review", ops, n(2), nil),
		statusTask("FID-9f13aa", "Harden login rate limiting", "todo", ops, n(3), nil),
		statusTask("FID-c2d4e5", "Retry a rejected backup push", "todo", ops, n(4), nil),
		statusTask("FID-0d4e8b", "Document backup age", "todo", ops, n(5), nil),
		statusTask("FID-a0b1c2", "Measure restore time", "done", ops, n(6), nil),
		statusTask("FID-a0b1c3", "Drop the old backup format", "done", ops, n(7), nil),
		statusTask("FID-a0b1c4", "Prototype cloud sessions", "cancelled", ops, n(8), nil),
		statusEpic("FID-8a1b2c", "Loop contract", "todo", n(2)),
		statusTask("FID-1f2e3d", "Contract fields and exit codes", "todo", loop, nil, n(1)),
		statusTask("FID-2f2e3d", "Golden tests for the contract and every exit code it names", "todo", loop, nil, n(2)),
		statusTask("FID-3f2e3d", "Scenario test for rework", "todo", loop, nil, n(2)),
		statusTask("FID-4f2e3d", "Scenario test for interruption", "todo", loop, nil, n(3)),
		statusTask("FID-5f2e3d", "Scenario test for two worktrees", "todo", loop, nil, n(3)),
		statusTask("FID-6f2e3d", "Version the contract", "todo", loop, nil, nil),
		statusEpic("FID-c0ffee", "Kit install", "todo", n(3)),
		statusEpic("FID-d0d0d0", "Finished groundwork", "done", n(4)),
		statusTask("FID-d1d1d1", "Leftover started task", "in_progress", closed, nil, nil),
		statusTask("FID-d2d2d2", "Leftover unstarted task", "todo", closed, nil, nil),
		statusTask("FID-5d2f77", "Fix journal replay line", "in_progress", nil, nil, n(1)),
		statusTask("FID-77aa01", "Tidy help text", "todo", nil, nil, n(2)),
		statusTask("FID-77aa02", "Refresh the changelog", "todo", nil, nil, nil),
		statusTask("FID-77aa03", "Old finished chore", "done", nil, nil, nil),
	}
	links := [][3]string{
		{"FID-71a0c2", "blocks", "FID-c2d4e5"}, {"FID-c2d4e5", "blocked_by", "FID-71a0c2"},
		{"FID-1f2e3d", "blocks", "FID-2f2e3d"}, {"FID-2f2e3d", "blocked_by", "FID-1f2e3d"},
		{"FID-9f13aa", "blocks", "FID-2f2e3d"}, {"FID-2f2e3d", "blocked_by", "FID-9f13aa"},
		{"FID-a0b1c2", "blocks", "FID-0d4e8b"}, {"FID-0d4e8b", "blocked_by", "FID-a0b1c2"},
	}
	return rows, links
}

func TestStatusGoldenWork(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		args []string
		work bool
	}{
		{"empty.txt", []string{"status"}, false},
		{"empty.json", []string{"--json", "status"}, false},
		{"work.txt", []string{"status"}, true},
		{"work-all.txt", []string{"status", "--all"}, true},
		{"work.json", []string{"--json", "status"}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newStatusFixture(t, 100*time.Hour)
			fixture.setTarget()
			fixture.writeRecord(statusRecordLine{"pushed", 3 * time.Hour, ""})
			if test.work {
				fixture.seed(statusWork())
			}
			fixture.golden(filepath.Join("status", test.name), test.args...)
		})
	}
}

// JSON is complete with or without --all.
func TestStatusJSONIgnoresAll(t *testing.T) {
	t.Parallel()
	fixture := newStatusFixture(t, 100*time.Hour)
	fixture.setTarget()
	fixture.seed(statusWork())
	capped, all := fixture.run("--json", "status"), fixture.run("--json", "status", "--all")
	if capped.err != nil || capped.stdout != all.stdout {
		t.Fatalf("status --json = %v\n%s\nwith --all\n%s", capped.err, capped.stdout, all.stdout)
	}
}

func TestStatusGoldenBackup(t *testing.T) {
	t.Parallel()
	hour, day := time.Hour, 24*time.Hour
	pushed := func(ago time.Duration) statusRecordLine { return statusRecordLine{"pushed", ago, ""} }
	failed := func(ago time.Duration, kind string) statusRecordLine { return statusRecordLine{"failed", ago, kind} }
	for _, test := range []struct {
		name    string
		started time.Duration
		target  bool
		config  string
		record  []statusRecordLine
	}{
		{"pushed", 100 * hour, true, "", []statusRecordLine{pushed(3 * hour)}},
		{"unchanged", 100 * hour, true, "", []statusRecordLine{pushed(30 * hour), {"unchanged", 2 * hour, ""}}},
		{"not-yet", 2 * hour, true, "", nil},
		{"over-24-hours", 100 * hour, true, "", []statusRecordLine{pushed(26 * hour)}},
		{"over-72-hours", 10 * day, true, "", []statusRecordLine{pushed(4*day + 3*hour)}},
		{"failed", 100 * hour, true, "", []statusRecordLine{pushed(2 * hour), failed(hour, "push-rejected")}},
		{"failures", 100 * hour, true, "", []statusRecordLine{pushed(8 * hour), failed(3*hour, "lock-unavailable"),
			failed(2*hour, "lock-unavailable"), failed(hour, "lock-unavailable")}},
		{"no-target", 2 * hour, false, "", nil},
		{"no-target-72-hours", 80 * hour, false, "", nil},
		{"unreadable-config", 100 * hour, false, "[backup]\ntarget = unquoted\n", []statusRecordLine{pushed(hour)}},
		{"restored", 40 * day, true, "", nil},
	} {
		for _, format := range []string{"txt", "json"} {
			t.Run(test.name+"."+format, func(t *testing.T) {
				t.Parallel()
				fixture := newStatusFixture(t, test.started)
				if test.target {
					fixture.setTarget()
				}
				if test.config != "" {
					fixture.writeUserConfig(test.config)
				}
				if test.record != nil {
					fixture.writeRecord(test.record...)
				}
				fixture.seed([]readyCommandRow{statusTask("FID-77aa01", "Tidy help text", "todo", nil, nil, nil)}, nil)
				args := []string{"status"}
				if format == "json" {
					args = []string{"--json", "status"}
				}
				fixture.golden(filepath.Join("status", "backup-"+test.name+"."+format), args...)
			})
		}
	}
}

func TestStatusRefusesWithoutALedger(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	result := executeCommandForTest(workspace, "--json", "status")
	if cli.ExitCode(result.err) != 3 || result.stdout != "" {
		t.Fatalf("status without a ledger = %#v, want project-not-found", result)
	}
	assertJSONFailureCode(t, result.stderr, "project.status", "project-not-found")
}
