//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestApprovedPRDefaultConflictJourney(t *testing.T) {
	t.Parallel()
	first := newDoneJourney(t, "standard")
	first.review(t, "primary", "APPROVE", "approve-first")
	root := gittest.AddWorktree(t, first.root)
	closureRun(t, root, "git", "checkout", "-b", "task-two", "origin/trunk")
	call := func(args ...string) map[string]any {
		t.Helper()
		offline := *first
		offline.root = root
		stdout, stderr, err := offline.run(t, append([]string{"--json"}, args...)...)
		if err != nil {
			t.Fatalf("fid %v: %v: %s", args, err, stderr)
		}
		return p23Envelope(t, stdout)
	}
	git := func(args ...string) {
		t.Helper()
		closureRun(t, root, "git", append([]string{"-c", "user.name=fid test", "-c", "user.email=fid@example.test"}, args...)...)
	}
	created := call("issue", "create", "--kind", "task", "--title", "Second approved task", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "second-create")
	task := p21PathString(t, created, "result", "issue", "id")
	taken := call("next", "--claim", "--actor", "agent:second-builder", "--model", "codex/gpt-6.1-sol")
	hold := p21PathString(t, taken, "result", "hold", "id")
	call("start", "--hold", hold, "--actor", "agent:second-builder", "--idempotency-key", "second-start")
	if err := os.WriteFile(filepath.Join(root, "change.txt"), []byte("Second implementation\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	git("add", ".")
	git("commit", "-qm", "Second implementation")
	call("check", "--hold", hold, "--actor", "agent:second-builder", "--idempotency-key", "second-check")
	handed := call("handoff", "--hold", hold, "--actor", "agent:second-builder", "--pr", "https://github.com/mbell-nzl/fid/pull/64", "--idempotency-key", "second-handoff")
	second := newDoneObservation(t, root, first.binary, task, hold, handed)
	second.review(t, "primary", "APPROVE", "approve-second")
	clean := call("next")
	if p21PathString(t, clean, "result", "kind") != "none" {
		t.Fatalf("approved clean task returned: %v", clean)
	}
	// The first approved PR lands on the local default branch.
	main := gittest.AddWorktree(t, root)
	closureRun(t, main, "git", "checkout", "--detach", "origin/trunk")
	closureRun(t, main, "git", "merge", "--ff-only", first.pr["headRefOid"].(string))
	closureRun(t, main, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	stdout, stderr, err := first.run(t, "--json", "done", first.task, "--actor", "agent:finisher", "--idempotency-key", "first-done")
	if err != nil || p21PathString(t, p23Envelope(t, stdout), "result", "status") != "done" {
		t.Fatalf("first done=%v %s", err, stderr)
	}
	nextOut, nextErr, nextFailure := second.run(t, "--json", "next", "--claim", "--actor", "agent:second-builder", "--model", "codex/gpt-6.1-sol")
	if nextFailure != nil {
		t.Fatalf("offline conflict return=%v %s", nextFailure, nextErr)
	}
	updated := p23Envelope(t, nextOut)
	if p21PathString(t, updated, "result", "kind") != "rework" || p21PathString(t, updated, "result", "rework", "cause") != "default-conflict" || p21PathString(t, updated, "result", "hold", "replaces") != hold {
		t.Fatalf("conflict rework=%v", updated)
	}
	contract := loopContract{t: t, decl: readLoopDeclaration(t)}
	contract.observe("default-conflict", updated, 0)
	hold = p21PathString(t, updated, "result", "hold", "id")
	// A real merge conflicts. Resolve it, check and hand off through public CLI.
	_, _, mergeErr := closureRunCommand(root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "merge", "--no-edit", "origin/trunk")
	if mergeErr == nil {
		t.Fatal("fixture did not conflict")
	}
	if err := os.WriteFile(filepath.Join(root, "change.txt"), []byte("First and second implementation\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	git("add", "change.txt")
	git("commit", "-qm", "Resolve approved PR conflict")
	call("check", "--hold", hold, "--actor", "agent:second-builder", "--idempotency-key", "update-check")
	handed = call("handoff", "--hold", hold, "--actor", "agent:second-builder", "--pr", "https://github.com/mbell-nzl/fid/pull/64", "--idempotency-key", "update-handoff")
	rounds := call("review", "list", task)["result"].(map[string]any)["rounds"].([]any)
	if len(rounds) != 1 || rounds[0].(map[string]any)["current"] != false {
		t.Fatalf("changed conflict retained review: %v", rounds)
	}
	second = newDoneObservation(t, root, first.binary, task, hold, handed)
	second.review(t, "primary", "APPROVE", "approve-updated")
	second.pr["number"], second.pr["url"], second.pr["headRefName"] = 64, "https://github.com/mbell-nzl/fid/pull/64", "task-two"
	second.requests[strings.Replace(donePRRequest, "pr view 63", "pr view 64", 1)] = map[string]any{"body": second.pr}
	stdout, stderr, err = second.run(t, "--json", "done", task, "--actor", "agent:finisher", "--idempotency-key", "second-done")
	if err != nil || p21PathString(t, p23Envelope(t, stdout), "result", "status") != "done" {
		t.Fatalf("second done=%v %s", err, stderr)
	}
	// Include the new immutable receipt in the V8 table comparison, and prove
	// actual archive import with no comparison objects at its destination.
	snapshot := func(workspace string) map[string]string {
		t.Helper()
		tables := closureProjectedTableSnapshots(t, first.binary, workspace)
		meta, _ := closureRunFID(t, first.binary, workspace, "sql", "SELECT last_event_id,last_event_hash,event_schema,projection_schema_version,replay_status FROM projection_meta")
		tables["projection_meta"] = meta
		return tables
	}
	before := snapshot(root)
	archive := filepath.Join(t.TempDir(), "updated.zip")
	call("export", "create", archive)
	call("export", "validate", archive)
	destination := gittest.NewWorkspace(t)
	p21Run(t, first.binary, destination, "--json", "export", "import", archive, "--repository", destination, "--idempotency-key", "import-update")
	if got := snapshot(destination); !reflect.DeepEqual(got, before) {
		for name, want := range before {
			if got[name] != want {
				t.Errorf("import changed %s: got %s; want %s", name, got[name], want)
			}
		}
	}
	call("rebuild")
	if got := snapshot(root); !reflect.DeepEqual(got, before) {
		t.Fatal("replay changed update receipt")
	}
	// Read its durable reservation/evidence directly via the public read command.
	receipts, _ := closureRunFID(t, first.binary, root, "sql", "SELECT * FROM default_conflicts")
	if !strings.Contains(receipts, "default-conflict") {
		t.Fatalf("missing receipt: %v", receipts)
	}
}
