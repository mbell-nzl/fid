package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
)

// nextWork is the --explain mock-up's project: an in-progress epic whose
// started task blocks one of its tasks, a sequential epic whose first task
// another worktree is working, and a task with no epic.
func nextWork(fixture *statusFixture) {
	fixture.t.Helper()
	n, s := intPointer, stringPointer
	ops, loop := s("FID-e3a1c0"), s("FID-8a1b2c")
	fixture.seed([]readyCommandRow{
		statusEpic("FID-e3a1c0", "Operating experience", "in_progress", n(1)),
		statusTask("FID-71a0c2", "Export validation", "in_progress", ops, nil, n(1)),
		statusTask("FID-c2d4e5", "Retry a rejected backup push", "todo", ops, nil, n(1)),
		statusTask("FID-9f13aa", "Harden login rate limiting", "todo", ops, nil, n(1)),
		statusTask("FID-0d4e8b", "Document backup age", "todo", ops, nil, n(2)),
		statusEpic("FID-8a1b2c", "Loop contract", "todo", n(2)),
		statusTask("FID-1f2e3d", "Contract fields and exit codes", "in_progress", loop, n(1), nil),
		statusTask("FID-2f2e3d", "Golden tests for the contract", "todo", loop, n(2), nil),
		statusTask("FID-3f2e3d", "Scenario test for rework", "todo", loop, n(3), nil),
		statusTask("FID-77aa01", "Tidy help text", "todo", nil, nil, n(2)),
		statusTask("FID-77aa03", "Old finished chore", "done", nil, nil, nil),
	}, [][3]string{{"FID-c2d4e5", "blocked_by", "FID-71a0c2"}})
	fixture.hold("HLD-77b2c0", "FID-1f2e3d", "/home/mark/projects/fid-wt2", 20*time.Minute)
}

// hold seeds an open hold on task, taken ago before statusNow, in the
// worktree at root; "" is the fixture's own worktree.
func (fixture *statusFixture) hold(id, task, root string, ago time.Duration) {
	fixture.t.Helper()
	if root == "" {
		root = fixture.root()
	}
	held := holds.New(id, task, project.WorktreeID(root), root, "agent:claude-builder", statusNow.Add(-ago), nil)
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(fixture.workspace), config.DatabaseFileName))
	if err != nil {
		fixture.t.Fatal(err)
	}
	defer closeCommandDB(fixture.t, conn)
	insert := held.Insert()
	if _, err := conn.Exec(`INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, task, insert.Values["worktree_id"], root, insert.Values["actor"], insert.Values["since"], insert.Values["expires_at"]); err != nil {
		fixture.t.Fatalf("seed hold %s: %v", id, err)
	}
}

// root is the fixture's worktree as fid finds it.
func (fixture *statusFixture) root() string {
	fixture.t.Helper()
	root, err := project.Resolve(fixture.workspace)
	if err != nil {
		fixture.t.Fatal(err)
	}
	return root
}

var newHoldID = regexp.MustCompile(`"HLD-[0-9a-f]{6}"`)

// seededHolds are the holds the goldens seed; fid never generates an ID
// that is taken, so any other is one the run took.
var seededHolds = map[string]bool{"HLD-77b2c0": true, "HLD-0a0a0a": true, "HLD-0b0b0b": true, "HLD-0c0c0c": true,
	"HLD-0d0d0d": true, "HLD-0e0e0e": true, "HLD-0f0f0f": true, "HLD-0a0b0c": true}

// nextGolden runs fid and compares its output to testdata/next/<name>, with
// any hold ID the run took named HLD-<new>.
func (fixture *statusFixture) nextGolden(name string, args ...string) {
	fixture.t.Helper()
	result := fixture.run(args...)
	if result.err != nil || result.stderr != "" {
		fixture.t.Fatalf("next %v = %v, stderr %q", args, result.err, result.stderr)
	}
	output := newHoldID.ReplaceAllStringFunc(fixture.placeholders(result.stdout), func(id string) string {
		if seededHolds[strings.Trim(id, `"`)] {
			return id
		}
		return `"HLD-<new>"`
	})
	assertGolden(fixture.t, filepath.Join("next", name), output)
}

func TestNextGolden(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name  string
		setup func(*statusFixture)
		args  []string
	}{
		{"next", nil, []string{"next"}},
		{"explain", nil, []string{"next", "--explain"}},
		{"taken", nil, []string{"next", "--claim", "--actor", "agent:claude-builder"}},
		{"taken-explain", nil, []string{"next", "--claim", "--actor", "agent:claude-builder", "--explain"}},
		{"resume", func(fixture *statusFixture) { fixture.hold("HLD-0a0a0a", "FID-9f13aa", "", 20*time.Minute) },
			[]string{"next", "--claim", "--actor", "agent:claude-builder"}},
		{"resume-read", func(fixture *statusFixture) { fixture.hold("HLD-0a0a0a", "FID-9f13aa", "", 20*time.Minute) }, []string{"next"}},
		{"interrupted", func(fixture *statusFixture) { fixture.hold("HLD-0b0b0b", "FID-71a0c2", "", 75*time.Minute) },
			[]string{"next", "--claim", "--actor", "agent:claude-builder"}},
		{"interrupted-read", func(fixture *statusFixture) { fixture.hold("HLD-0b0b0b", "FID-71a0c2", "", 75*time.Minute) }, []string{"next"}},
		{"lapsed-here", func(fixture *statusFixture) { fixture.hold("HLD-0c0c0c", "FID-0d4e8b", "", 26*time.Hour) },
			[]string{"next", "--claim", "--actor", "agent:claude-builder", "--explain"}},
		{"lapsed-elsewhere", func(fixture *statusFixture) {
			fixture.hold("HLD-0d0d0d", "FID-9f13aa", "/home/mark/projects/fid-wt3", time.Minute)
			fixture.hold("HLD-0e0e0e", "FID-0d4e8b", "/home/mark/projects/fid-wt3-tests", time.Minute)
			fixture.hold("HLD-0f0f0f", "FID-77aa01", "/home/mark/projects/fid-wt4", 2*time.Hour)
			fixture.hold("HLD-0a0b0c", "FID-71a0c2", "/home/mark/projects/fid-wt5", 2*time.Hour)
		}, []string{"next", "--claim", "--actor", "agent:claude-builder", "--explain"}},
		{"nothing", func(fixture *statusFixture) {
			fixture.hold("HLD-0d0d0d", "FID-9f13aa", "/home/mark/projects/fid-wt3", time.Minute)
			fixture.hold("HLD-0e0e0e", "FID-0d4e8b", "/home/mark/projects/fid-wt3-tests", time.Minute)
			fixture.hold("HLD-0f0f0f", "FID-77aa01", "/home/mark/projects/fid-wt4", time.Minute)
			fixture.hold("HLD-0a0b0c", "FID-71a0c2", "/home/mark/projects/fid-wt5", 2*time.Hour)
		}, []string{"next"}},
		{"nothing-explain", func(fixture *statusFixture) {
			fixture.hold("HLD-0d0d0d", "FID-9f13aa", "/home/mark/projects/fid-wt3", time.Minute)
			fixture.hold("HLD-0e0e0e", "FID-0d4e8b", "/home/mark/projects/fid-wt3-tests", time.Minute)
			fixture.hold("HLD-0a0b0c", "FID-71a0c2", "/home/mark/projects/fid-wt5", 2*time.Hour)
			fixture.hold("HLD-0f0f0f", "FID-77aa01", "/home/mark/projects/fid-wt4", time.Minute)
		}, []string{"next", "--claim", "--actor", "agent:claude-builder", "--explain"}},
	} {
		for _, format := range []string{"txt", "json"} {
			t.Run(test.name+"."+format, func(t *testing.T) {
				t.Parallel()
				fixture := newStatusFixture(t, 100*time.Hour)
				nextWork(fixture)
				if test.setup != nil {
					test.setup(fixture)
				}
				args := test.args
				if format == "json" {
					args = append([]string{"--json"}, args...)
				}
				fixture.nextGolden(test.name+"."+format, args...)
			})
		}
	}
	for _, format := range []string{"txt", "json"} {
		t.Run("empty."+format, func(t *testing.T) {
			t.Parallel()
			args := []string{"next", "--claim", "--actor", "human:mark"}
			if format == "json" {
				args = append([]string{"--json"}, args...)
			}
			newStatusFixture(t, 100*time.Hour).nextGolden("empty."+format, args...)
		})
	}
}

// decodeNext reads a fid next success envelope.
func decodeNext(t *testing.T, result commandExecutionResult, operation string) (map[string]any, any) {
	t.Helper()
	var envelope struct {
		OK        bool           `json:"ok"`
		Operation string         `json:"operation"`
		Result    map[string]any `json:"result"`
		Event     any            `json:"event"`
	}
	if result.err != nil || result.stderr != "" || json.Unmarshal([]byte(result.stdout), &envelope) != nil || !envelope.OK || envelope.Operation != operation {
		t.Fatalf("next = %#v, want %s", result, operation)
	}
	return envelope.Result, envelope.Event
}

// A task that blocks an epic is unblocking work only while the epic has
// open tasks, which wait on it through their epic (Mark, 12 November
// 2026): blocking an empty epic doesn't put a P5 task ahead of a P0 one.
func TestNextServesAnEpicsBlockerFirstOnlyWhileTheEpicHasOpenTasks(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runFIDCommand(t, workspace, "init", "--prefix", "FID")
	create := func(args ...string) string {
		_, created := runFIDCommand(t, workspace, append([]string{"issue", "create"}, append(args, "--actor", "human:mark")...)...)
		return parseCreatedIssueID(t, created)
	}
	epic := create("--kind", "epic", "--title", "Empty for now", "--scheduling", "prioritised")
	blocker := create("--kind", "task", "--title", "Blocks the epic", "--priority", "5")
	urgent := create("--kind", "task", "--title", "Urgent", "--priority", "0")
	runFIDCommand(t, workspace, "link", "add", blocker, "blocks", epic, "--from-type", "issue", "--to-type", "issue", "--actor", "human:mark")
	served := func() string {
		t.Helper()
		result, _ := decodeNext(t, executeCommandForTest(workspace, "--json", "next"), "next.get")
		task, _ := result["task"].(map[string]any)
		return fmt.Sprint(task["id"])
	}
	if got := served(); got != urgent {
		t.Fatalf("next with an empty blocked epic = %s, want the P0 task %s", got, urgent)
	}
	create("--kind", "task", "--title", "Waits through its epic", "--epic", epic)
	if got := served(); got != blocker {
		t.Fatalf("next once the blocked epic has a task = %s, want its blocker %s", got, blocker)
	}
}

// next's JSON names the served task's epic, with its progress as fid
// status gives it, and each queue entry's epic blockers (Mark, 12 November
// 2026). The goldens show a null epic for a task with none and for none.
func TestNextJSONNamesTheEpic(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runFIDCommand(t, workspace, "init", "--prefix", "FID")
	create := func(args ...string) string {
		_, created := runFIDCommand(t, workspace, append([]string{"issue", "create"}, append(args, "--actor", "human:mark")...)...)
		return parseCreatedIssueID(t, created)
	}
	served, waiting := create("--kind", "epic", "--title", "Served epic", "--scheduling", "prioritised"), create("--kind", "epic", "--title", "Waiting epic", "--scheduling", "prioritised")
	blocker := create("--kind", "task", "--title", "Blocks the waiting epic", "--epic", served)
	waits := create("--kind", "task", "--title", "Waits through its epic", "--epic", waiting)
	runFIDCommand(t, workspace, "link", "add", blocker, "blocks", waiting, "--from-type", "issue", "--to-type", "issue", "--actor", "human:mark")
	result, _ := decodeNext(t, executeCommandForTest(workspace, "--json", "next"), "next.get")
	epic, _ := json.Marshal(result["epic"])
	if want := `{"id":"` + served + `","progress":{"done":0,"total":1},"title":"Served epic"}`; string(epic) != want {
		t.Fatalf("next's epic = %s, want %s", epic, want)
	}
	got := map[string]string{}
	for _, entry := range result["queue"].([]any) {
		entry := entry.(map[string]any)
		blockers, _ := json.Marshal(entry["epic_blocked_by"])
		got[entry["id"].(string)] = fmt.Sprint(entry["reason"]) + " " + string(blockers)
	}
	if want := map[string]string{blocker: "<nil> []", waits: `epic-blocked ["` + blocker + `"]`}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("next's queue = %v, want %v", got, want)
	}
}

// Re-running --claim in the same worktree returns the hold it already has,
// with no Event; a read never writes one.
func TestNextClaimTwiceKeepsTheHold(t *testing.T) {
	t.Parallel()
	fixture := newStatusFixture(t, 100*time.Hour)
	nextWork(fixture)
	read, event := decodeNext(t, fixture.run("--json", "next"), "next.get")
	if read["kind"] != "new" || read["hold"] != nil || event != nil {
		t.Fatalf("next = %v with event %v, want new with no hold or Event", read, event)
	}
	first, event := decodeNext(t, fixture.run("--json", "next", "--claim", "--actor", "agent:one"), "next.claim")
	hold, _ := first["hold"].(map[string]any)
	if first["kind"] != "new" || hold == nil || hold["task_id"] != "FID-9f13aa" || event == nil {
		t.Fatalf("next --claim = %v with event %v, want a new hold on FID-9f13aa", first, event)
	}
	again, event := decodeNext(t, fixture.run("--json", "next", "--claim", "--actor", "agent:one"), "next.claim")
	if again["kind"] != "resume" || event != nil || again["hold"].(map[string]any)["id"] != hold["id"] {
		t.Fatalf("next --claim again = %v with event %v, want resume of %v with no Event", again, event, hold["id"])
	}
	read, event = decodeNext(t, fixture.run("--json", "next"), "next.get")
	if read["kind"] != "resume" || read["hold"] != nil || event != nil {
		t.Fatalf("next after taking = %v with event %v, want resume with no hold or Event", read, event)
	}
}

// Every refusal named in the brief: a flag is never accepted and ignored.
// Human text names no kernel word; JSON keeps each code's details.
func TestNextRefusals(t *testing.T) {
	t.Parallel()
	fixture := newStatusFixture(t, 100*time.Hour)
	nextWork(fixture)
	for _, test := range []struct {
		name      string
		args      []string
		exit      int
		operation string
		code      string
	}{
		{"--actor without --claim", []string{"next", "--actor", "agent:one"}, 2, "next.get", "invalid-argument"},
		{"--claim without --actor", []string{"next", "--claim"}, 2, "next.claim", "invalid-argument"},
		{"a malformed actor", []string{"next", "--claim", "--actor", "builder"}, 2, "next.claim", "invalid-argument"},
		{"an empty actor", []string{"next", "--claim", "--actor="}, 2, "next.claim", "invalid-argument"},
		{"a service actor", []string{"next", "--claim", "--actor", "service:ci"}, 5, "next.claim", "authority-denied"},
		{"an argument", []string{"next", "FID-9f13aa"}, 2, "next.get", "invalid-argument"},
	} {
		t.Run(test.name, func(t *testing.T) {
			result := fixture.run(append([]string{"--json"}, test.args...)...)
			if cli.ExitCode(result.err) != test.exit || result.stdout != "" {
				t.Fatalf("%v = %#v, want exit %d", test.args, result, test.exit)
			}
			assertJSONFailureCode(t, result.stderr, test.operation, test.code)
			if test.code == "authority-denied" && !strings.Contains(result.stderr, `"details":{"actor_id":"service:ci","operation":"next.claim","required_actions":["next.claim"],"required_capabilities":["next.claim"]}`) {
				t.Fatalf("%v = %s, want authority-denied's details", test.args, result.stderr)
			}
			human := fixture.run(test.args...)
			if cli.ExitCode(human.err) != test.exit || !strings.HasPrefix(human.stderr, "error: ") {
				t.Fatalf("%v = %#v, want a refusal", test.args, human)
			}
			if found := kernelWords.FindAllString(human.stderr, -1); len(found) > 0 {
				t.Fatalf("%v refusal names kernel words %q:\n%s", test.args, found, human.stderr)
			}
		})
	}
	if open := fixture.run("--json", "next"); !strings.Contains(open.stdout, `"kind":"new"`) {
		t.Fatalf("after the refusals, next = %#v, want nothing taken", open)
	}
	outside := gittest.NewWorkspace(t)
	result := executeCommandForTest(outside, "--json", "next", "--claim", "--actor", "agent:one")
	if cli.ExitCode(result.err) != 3 {
		t.Fatalf("next --claim without a ledger = %#v, want project-not-found", result)
	}
	assertJSONFailureCode(t, result.stderr, "next.claim", "project-not-found")
}

// lastJournalEvent is the journal's last Event, decoded.
func lastJournalEvent(t *testing.T, workspace string) map[string]any {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(config.LedgerDir(workspace), config.JournalDirName, "*.jsonl"))
	if err != nil || len(files) == 0 {
		t.Fatalf("journal files = %v, %v", files, err)
	}
	lines := strings.Split(strings.TrimSpace(string(mustReadCommandFile(t, files[len(files)-1]))), "\n")
	var event map[string]any
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &event); err != nil {
		t.Fatal(err)
	}
	return event
}

// A task completed or cancelled ends its hold in the same Event as the
// transition. The journal replays: a rebuild gives back the same holds, and
// an export of the ledger validates.
func TestNextHoldEndsInTheEventThatClosesItsTask(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runFIDCommand(t, workspace, "init", "--prefix", "FID")
	runFIDCommand(t, workspace, "workflow", "define", "--name", "quick", "--applies-to", "task", "--states", "todo,done,cancelled",
		"--categories", "unstarted,completed,cancelled", "--transitions", "todo->done,todo->cancelled", "--actor", "human:mark")
	runFIDCommand(t, workspace, "workflow", "assign", "--name", "quick", "--kind", "task", "--actor", "human:mark")
	for _, title := range []string{"Cancelled while held", "Done while held", "Left held"} {
		runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", title, "--actor", "human:mark")
	}
	take := func() string {
		t.Helper()
		result, event := decodeNext(t, executeCommandForTest(workspace, "--json", "next", "--claim", "--actor", "agent:builder"), "next.claim")
		if result["kind"] != "new" || event == nil {
			t.Fatalf("next --claim = %v, want a new hold", result)
		}
		return result["task"].(map[string]any)["id"].(string)
	}
	for _, to := range []string{"cancelled", "done"} {
		task := take()
		runFIDCommand(t, workspace, "issue", "transition", task, "--to", to, "--actor", "human:mark")
		event := lastJournalEvent(t, workspace)
		ended := false
		for _, effect := range event["effects"].([]any) {
			effect := effect.(map[string]any)
			if effect["table"] == "holds" {
				values := effect["values"].(map[string]any)
				ended = effect["op"] == "update" && values["end_reason"] == "task-"+map[string]string{"cancelled": "cancelled", "done": "completed"}[to]
			}
		}
		if event["event_type"] != "issue.transitioned" || !ended {
			t.Fatalf("moving %s to %s = %v, want the transition to end its hold", task, to, event)
		}
	}
	take()
	holdsQuery := "SELECT task_id, actor, end_reason, replaces_hold_id FROM holds ORDER BY since, task_id"
	before := executeCommandForTest(workspace, "sql", holdsQuery)
	if before.err != nil || strings.Count(before.stdout, `"end_reason":"task-`) != 2 || strings.Count(before.stdout, `"end_reason":null`) != 1 {
		t.Fatalf("holds = %#v, want one ended by each transition and one open", before)
	}
	if rebuilt := executeCommandForTest(workspace, "rebuild"); rebuilt.err != nil {
		t.Fatalf("rebuild = %#v", rebuilt)
	}
	if after := executeCommandForTest(workspace, "sql", holdsQuery); after.stdout != before.stdout {
		t.Fatalf("holds after a rebuild = %s, want %s", after.stdout, before.stdout)
	}
	archive := filepath.Join(t.TempDir(), "export.zip")
	if created := executeCommandForTest(workspace, "export", "create", archive); created.err != nil {
		t.Fatalf("export create = %#v", created)
	}
	if validated := executeCommandForTest(workspace, "--json", "export", "validate", archive); validated.err != nil {
		t.Fatalf("export validate = %#v", validated)
	}
}

// --claim reads the clock once it holds the write lock: a hold that lapses
// while fid waits for the lock is read as lapsed, so its started task is
// resumed under a new hold that replaces it.
func TestNextClaimReadsTheClockUnderTheLock(t *testing.T) {
	t.Parallel()
	fixture := newStatusFixture(t, 100*time.Hour)
	nextWork(fixture)
	// This worktree's hold on the started FID-71a0c2 lapses a minute after
	// statusNow.
	fixture.hold("HLD-0b0b0b", "FID-71a0c2", "", holds.Lifetime-time.Minute)
	var read, released atomic.Bool
	fixture.clock = func() time.Time {
		read.Store(true)
		if released.Load() {
			return statusNow.Add(2 * time.Minute)
		}
		return statusNow
	}
	done := make(chan commandExecutionResult)
	if err := journal.WithLock(fixture.workspace, journal.LockOwner{Command: "test"}, func(*journal.WriteLock) error {
		go func() { done <- fixture.run("--json", "next", "--claim", "--actor", "agent:one") }()
		// fid is waiting for the lock once it has read the clock, or after
		// a while if it reads the clock only once it holds the lock.
		for deadline := time.Now().Add(500 * time.Millisecond); !read.Load() && time.Now().Before(deadline); {
			time.Sleep(10 * time.Millisecond)
		}
		released.Store(true)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	result, event := decodeNext(t, <-done, "next.claim")
	hold, _ := result["hold"].(map[string]any)
	if result["kind"] != "resume" || event == nil || hold == nil || hold["replaces"] != "HLD-0b0b0b" {
		t.Fatalf("next --claim after the hold lapsed = %v with event %v, want a resume under a new hold replacing HLD-0b0b0b", result, event)
	}
}

// While another command holds the write lock past fid's wait, --claim
// refuses as retryable; a read doesn't need the lock.
func TestNextClaimRefusesWhileTheLockIsHeld(t *testing.T) {
	t.Parallel()
	fixture := newStatusFixture(t, 100*time.Hour)
	nextWork(fixture)
	var claimed, read commandExecutionResult
	if err := journal.WithLock(fixture.workspace, journal.LockOwner{Command: "test"}, func(*journal.WriteLock) error {
		claimed = fixture.run("--json", "next", "--claim", "--actor", "agent:one")
		read = fixture.run("--json", "next")
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if cli.ExitCode(claimed.err) != 4 || !strings.Contains(claimed.stderr, `"retryable":true`) {
		t.Fatalf("next --claim under a held lock = %#v, want retryable lock-unavailable", claimed)
	}
	assertJSONFailureCode(t, claimed.stderr, "next.claim", "lock-unavailable")
	decodeNext(t, read, "next.get")
}
