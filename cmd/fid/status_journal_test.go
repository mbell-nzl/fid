package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// projectionFiles are the projection's files in the ledger.
var projectionFiles = []string{config.DatabaseFileName, config.DatabaseFileName + "-wal", config.DatabaseFileName + "-shm"}

// newJournalStatusFixture is a fresh project whose work fid writes through
// its journal, with a home and user config of its own.
func newJournalStatusFixture(t *testing.T) *statusFixture {
	t.Helper()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init", "--json")
	base := t.TempDir()
	return &statusFixture{t: t, workspace: workspace, configHome: filepath.Join(base, "config"), home: filepath.Join(base, "home")}
}

// createTask creates a task through fid and returns its ID.
func (fixture *statusFixture) createTask(title string) string {
	fixture.t.Helper()
	_, stderr := runFIDCommand(fixture.t, fixture.workspace, "issue", "create", "--kind", "task", "--title", title, "--actor", "human:mark")
	return parseCreatedIssueID(fixture.t, stderr)
}

func (fixture *statusFixture) start(id string) {
	fixture.t.Helper()
	runFIDCommand(fixture.t, fixture.workspace, "issue", "transition", id, "--to", "in_progress", "--actor", "human:mark")
}

// snapshotProjection copies the projection's files aside, and returns a
// function that puts that copy back in place of whatever is there then.
func (fixture *statusFixture) snapshotProjection() (restore func()) {
	fixture.t.Helper()
	ledger, saved := config.LedgerDir(fixture.workspace), fixture.t.TempDir()
	copyFiles(fixture.t, ledger, saved)
	return func() {
		fixture.t.Helper()
		fixture.removeProjection()
		copyFiles(fixture.t, saved, ledger)
	}
}

// copyFiles copies the projection's files that exist in from into to.
func copyFiles(t *testing.T, from, to string) {
	t.Helper()
	for _, name := range projectionFiles {
		payload, err := os.ReadFile(filepath.Join(from, name))
		if os.IsNotExist(err) {
			continue
		}
		if err == nil {
			err = os.WriteFile(filepath.Join(to, name), payload, 0o644)
		}
		if err != nil {
			t.Fatalf("copy %s: %v", name, err)
		}
	}
}

func (fixture *statusFixture) removeProjection() {
	fixture.t.Helper()
	for _, name := range projectionFiles {
		if err := os.Remove(filepath.Join(config.LedgerDir(fixture.workspace), name)); err != nil && !os.IsNotExist(err) {
			fixture.t.Fatal(err)
		}
	}
}

// statusTaskIDs is the IDs of the tasks fid --json status lists in progress
// and to do, with no listed epic.
func statusTaskIDs(t *testing.T, stdout string) (inProgress, toDo []string) {
	t.Helper()
	type section struct {
		Tasks []struct {
			ID string `json:"id"`
		} `json:"tasks"`
	}
	var envelope struct {
		Result struct {
			InProgress section `json:"in_progress"`
			ToDo       section `json:"to_do"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil {
		t.Fatalf("status --json = %v:\n%s", err, stdout)
	}
	inProgress, toDo = []string{}, []string{}
	for _, task := range envelope.Result.InProgress.Tasks {
		inProgress = append(inProgress, task.ID)
	}
	for _, task := range envelope.Result.ToDo.Tasks {
		toDo = append(toDo, task.ID)
	}
	return inProgress, toDo
}

// Status reads the work from the journal, through the projection: it
// rebuilds a missing projection and catches up a stale one.
func TestStatusShowsTheJournalsWork(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		// work writes the work, leaving the projection as the test needs,
		// and returns the tasks in progress and to do.
		work func(fixture *statusFixture) (inProgress, toDo []string)
	}{
		{"missing projection", func(fixture *statusFixture) ([]string, []string) {
			started, waiting := fixture.createTask("Started task"), fixture.createTask("Waiting task")
			fixture.start(started)
			fixture.removeProjection()
			return []string{started}, []string{waiting}
		}},
		{"stale projection", func(fixture *statusFixture) ([]string, []string) {
			started := fixture.createTask("Started task")
			restore := fixture.snapshotProjection()
			fixture.start(started)
			waiting := fixture.createTask("Waiting task")
			restore()
			return []string{started}, []string{waiting}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newJournalStatusFixture(t)
			wantInProgress, wantToDo := test.work(fixture)
			result := fixture.run("--json", "status")
			if result.err != nil || result.stderr != "" {
				t.Fatalf("status --json = %v, stderr %q", result.err, result.stderr)
			}
			inProgress, toDo := statusTaskIDs(t, result.stdout)
			if !slices.Equal(inProgress, wantInProgress) || !slices.Equal(toDo, wantToDo) {
				t.Fatalf("status lists %v in progress and %v to do, want %v and %v:\n%s", inProgress, toDo, wantInProgress, wantToDo, result.stdout)
			}
			text := fixture.run("status")
			for _, want := range append(wantInProgress, wantToDo...) {
				if text.err != nil || !strings.Contains(text.stdout, want) {
					t.Fatalf("status = %v, text doesn't list %s:\n%s", text.err, want, text.stdout)
				}
			}
		})
	}
}

// A journal whose tail isn't a whole Event stops status, as it stops every
// read, before it shows any work.
func TestStatusStopsOnACorruptJournalTail(t *testing.T) {
	t.Parallel()
	fixture := newJournalStatusFixture(t)
	fixture.createTask("Task before the damage")
	appendRawJournalLine(t, mustJournalPath(t, fixture.workspace), `{"event_id":"EVT-20261108-999999"`)
	result := fixture.run("--json", "status")
	if cli.ExitCode(result.err) != 6 || result.stdout != "" {
		t.Fatalf("status --json = %#v, want journal-replay-stopped", result)
	}
	assertJSONFailureCode(t, result.stderr, "project.status", "journal-replay-stopped")
	assertStoppedJournalFile(t, result.stderr, fixture.workspace)
	text := fixture.run("status")
	if cli.ExitCode(text.err) != 6 || text.stdout != "" || !strings.HasPrefix(text.stderr, "error: journal replay stopped\nreason: journal-replay-stopped; ") {
		t.Fatalf("status = %#v, want journal-replay-stopped", text)
	}
}
