//go:build linux

package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

// The child is fid itself, with real terminal file descriptors. It receives a
// response only after its prompt has reached the terminal. stdout stays a pipe.
func approvePublicForTest(t *testing.T, root, task, answer string, inputTTY, errorTTY, jsonMode bool, extra ...string) commandExecutionResult {
	t.Helper()
	return personPublicForTest(t, root, task, answer, inputTTY, errorTTY, jsonMode, "approve", extra...)
}

func personPublicForTest(t *testing.T, root, task, answer string, inputTTY, errorTTY, jsonMode bool, action string, extra ...string) commandExecutionResult {
	t.Helper()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	args := []string{action, task, "--actor", "human:mark", "--reason", "Reviewed the full change"}
	if jsonMode {
		args = append([]string{"--json"}, args...)
	}
	args = append(args, extra...)
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), runAsFIDVariable+"=1")
	if testing.CoverMode() != "" {
		cmd.Env = append(cmd.Env, "GOCOVERDIR="+t.TempDir())
	}
	leader, follower := openTerminal(t, display.Screen{Width: 120, Height: 60})
	var out, errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	cmd.Stdin = strings.NewReader(answer)
	if inputTTY {
		cmd.Stdin = follower
	}
	if errorTTY {
		cmd.Stderr = follower
	}
	reached := make(chan []byte, 1)
	go func() {
		var read bytes.Buffer
		sent := false
		chunk := make([]byte, 4096)
		for {
			n, err := leader.Read(chunk)
			if n > 0 {
				read.Write(chunk[:n])
				if !sent && strings.Contains(read.String(), "Type "+task+" to "+action+": ") {
					sent = true
					_, _ = io.WriteString(leader, answer)
				}
			}
			if err != nil {
				break
			}
		}
		reached <- read.Bytes()
	}()
	err = cmd.Start()
	_ = follower.Close()
	if err == nil {
		err = cmd.Wait()
	}
	terminal := bytes.ReplaceAll(<-reached, []byte("\r\n"), []byte("\n"))
	return commandExecutionResult{stdout: out.String(), stderr: errout.String() + string(terminal), err: err}
}

func assertApprovalState(t *testing.T, root, task, want string) approvalrecord.View {
	t.Helper()
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	h, err := reviewrecord.LatestHandoff(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	fp, err := fingerprint.Compute(root, h.Base, h.Head)
	if err != nil {
		t.Fatal(err)
	}
	got, err := approvalrecord.Current(conn, task, fp)
	if err != nil || got.State != want {
		t.Fatalf("approval=%+v/%v want %s", got, err, want)
	}
	return got
}

func TestV6ApprovalSurvivesCleanRebaseAndStalesOnEdit(t *testing.T) {
	t.Parallel()
	// The terminal child reads real time. Every fixture command must too; a
	// captured time from before UTC midnight would backdate a later Event.
	f, task := canonicalReviewFixtureWithClock(t, "assured", false, time.Now)
	assertApprovalState(t, f.workspace, task, "missing")
	first := approvePublicForTest(t, f.workspace, task, task+"\n", true, true, true, "--idempotency-key", "v6-first")
	approveSuccess(t, first)
	old := assertApprovalState(t, f.workspace, task, "current")
	lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "v6-rebase", "human:courier", task), "review.record")
	runGitForTest(t, f.workspace, "checkout", "-q", "trunk")
	writeBackupTestFile(t, filepath.Join(f.workspace, "unrelated.txt"), "upstream change\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Move base")
	runGitForTest(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	runGitForTest(t, f.workspace, "checkout", "-q", "task")
	runGitForTest(t, f.workspace, "rebase", "trunk")
	prepareSecondReviewHandoff(t, f, task)
	reviewHandoffTask(t, f, "v6-rebased-handoff", task)
	rebased := assertApprovalState(t, f.workspace, task, "current")
	if rebased.Decision.ID != old.Decision.ID {
		t.Fatal("rebase replaced approval")
	}
	in := changesInput("primary")
	in["findings"] = []any{}
	lifecycleResult(t, recordReviewTask(t, f, in, "v6-edit-request", "human:courier", task), "review.record")
	prepareSecondReviewHandoff(t, f, task)
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "edited review change\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Edit task")
	reviewHandoffTask(t, f, "v6-edited-handoff", task)
	assertApprovalState(t, f.workspace, task, "stale")
	again := approvePublicForTest(t, f.workspace, task, task+"\n", true, true, true, "--idempotency-key", "v6-second")
	approveSuccess(t, again)
	for _, part := range []string{"notes.md", "stale", "Open majors:", "Retry conflicts lack coverage", old.Decision.ID, "no other-model review"} {
		if !strings.Contains(again.stderr, part) {
			t.Errorf("reapproval missing %q: %s", part, again.stderr)
		}
	}
	latest := assertApprovalState(t, f.workspace, task, "current")
	if latest.Decision.ID == old.Decision.ID {
		t.Fatal("reapproval reused old decision")
	}
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
	assertApprovalState(t, f.workspace, task, "current")
}

func TestApproveRealTerminalBoundaries(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name            string
		input, errorTTY bool
		actor, answer   string
		json            bool
		want            string
	}{
		{"piped input", false, true, "human:mark", "", true, "standard input and standard error must be terminals"},
		{"piped prompt", true, false, "human:mark", "", true, "standard input and standard error must be terminals"},
		{"agent terminal", true, true, "agent:builder", "", true, "only a person"},
		{"service terminal", true, true, "service:ci", "", true, "only a person"},
		{"cancel terminal", true, true, "human:mark", "no\n", true, "Approval cancelled"},
		{"human text", true, true, "human:mark", "TASK", false, "Type "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithClock(t, "standard", false, time.Now)
			answer := tc.answer
			if answer == "TASK" {
				answer = task + "\n"
			}
			before := reviewProjection(t, f)
			got := approvePublicForTest(t, f.workspace, task, answer, tc.input, tc.errorTTY, tc.json, "--actor", tc.actor)
			if !strings.Contains(got.stderr, tc.want) {
				t.Fatal(got)
			}
			if tc.name == "human text" {
				if got.err != nil || !strings.Contains(got.stdout, "✓ Approved "+task) {
					t.Fatal(got)
				}
			} else if reviewProjection(t, f) != before || got.stdout != "" {
				t.Fatal("refusal changed state", got)
			}
		})
	}
}

func TestApproveRevertedChangeUsesOriginalDecision(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithClock(t, "standard", false, time.Now)
	approveSuccess(t, approvePublicForTest(t, f.workspace, task, task+"\n", true, true, true))
	first := assertApprovalState(t, f.workspace, task, "current")
	requestChanges := changesInput("primary")
	lifecycleResult(t, recordReviewTask(t, f, requestChanges, "edit-request", "human:courier", task), "review.record")
	prepareSecondReviewHandoff(t, f, task)
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "different change for approval\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Edit approved change")
	editHead := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
	reviewHandoffTask(t, f, "edited-handoff", task)
	assertApprovalState(t, f.workspace, task, "stale")
	approveSuccess(t, approvePublicForTest(t, f.workspace, task, task+"\n", true, true, true))
	second := assertApprovalState(t, f.workspace, task, "current")
	if second.Decision.ID == first.Decision.ID || second.Decision.Binding.Fingerprint == first.Decision.Binding.Fingerprint {
		t.Fatal("edited change must have a different approval and fingerprint")
	}
	requestChanges["findings"] = []any{}
	lifecycleResult(t, recordReviewTask(t, f, requestChanges, "revert-request", "human:courier", task), "review.record")
	prepareSecondReviewHandoff(t, f, task)
	runGitForTest(t, f.workspace, "revert", "--no-edit", editHead)
	reviewHandoffTask(t, f, "reverted-handoff", task)
	for _, rebuild := range []bool{false, true} {
		if rebuild {
			if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
				t.Fatal(err)
			}
		}
		restored := assertApprovalState(t, f.workspace, task, "current")
		if *restored.Decision != *first.Decision {
			t.Fatalf("rebuild=%t: reverted change approval=%+v, want original %+v", rebuild, restored.Decision, first.Decision)
		}
	}
}
