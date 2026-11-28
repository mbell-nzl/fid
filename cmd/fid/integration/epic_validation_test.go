//go:build integration && linux

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestEpicValidationLoopCompletionJourney(t *testing.T) {
	t.Parallel()
	j := newDoneJourney(t, "standard", true)
	for _, v := range []struct{ claim, verifier, recipe string }{{"Default branch checks", "check:test", ""}, {"Inspect recipe", "human", "recipe.md"}} {
		args := []string{"--json", "epic", "validation", "add", j.epic, "--claim", v.claim, "--verifier", v.verifier, "--actor", "agent:builder", "--idempotency-key", v.verifier}
		if v.recipe != "" {
			args = append(args, "--recipe", v.recipe)
		}
		p21Run(t, j.binary, j.root, args...)
	}
	_, stderr, err := approvalCorpusCommand(t, j.root, j.binary, "--json", "epic", "approve", j.epic, "--actor", "human:mark", "--reason", "Approve validations", "--idempotency-key", "approve-plan")
	if err != nil {
		t.Fatalf("epic approval: %v %s", err, stderr)
	}
	j.review(t, "primary", "APPROVE", "review")
	head := strings.TrimSpace(closureRunOutput(t, j.root, "git", "rev-parse", "HEAD"))
	j.pr["mergeCommit"] = map[string]any{"oid": head}
	j.ci[0].(map[string]any)["head_sha"] = head
	j.requests[strings.Replace(doneCIRequest, strings.Repeat("a", 40), head, 1)] = map[string]any{"body": map[string]any{"total_count": len(j.ci), "runs": j.ci}}
	if out, stderr, err := j.run(t, "--json", "done", j.task, "--actor", "agent:builder"); err != nil {
		t.Fatalf("done=%v %s %s", err, out, stderr)
	}
	// The locally known default does not yet contain this task's merge.
	if out, stderr, err := j.run(t, "--json", "epic", "validation", "verify", j.epic, "--actor", "agent:builder", "--idempotency-key", "wrong-commit"); err == nil || !strings.Contains(stderr, "does not contain epic merge") {
		t.Fatalf("wrong commit=%v %s %s", err, out, stderr)
	}
	closureRun(t, j.root, "git", "update-ref", "refs/remotes/origin/trunk", head)
	// Run the actual loop with the built fid. Agents and gh are isolated stubs;
	// no task remains to invoke an agent or query GitHub.
	stub := filepath.Join(j.bin, "agent")
	if err := os.WriteFile(stub, []byte("#!/bin/sh\nif [ \"$1\" = --version ]; then echo stub-1; exit 0; fi\nexit 97\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	loop, err := filepath.Abs("../../../loops/loop.py")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(t.Context(), "python3", loop)
	cmd.Dir = j.root
	cmd.Env = append(os.Environ(), "PATH="+j.bin+string(os.PathListSeparator)+os.Getenv("PATH"), "LOOP_FID="+j.binary, "LOOP_EPIC="+j.epic, "LOOP_BUILDER_COMMAND="+stub, "LOOP_REVIEWER_COMMAND="+stub, "LOOP_JUDGE_COMMAND="+stub, "LOOP_BUILDER_USAGE=none", "LOOP_REVIEWER_USAGE=none", "LOOP_JUDGE_USAGE=none", "LOOP_REVIEW_WORKTREE="+filepath.Join(t.TempDir(), "review"), "XDG_CONFIG_HOME="+t.TempDir())
	raw, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(raw), "epic validation:") || !strings.Contains(string(raw), `"outcome": "passed"`) {
		t.Fatalf("loop=%v %s", err, raw)
	}
	if out, stderr, err := j.run(t, "--json", "done", j.epic, "--actor", "human:mark", "--reason", "Finished"); err == nil || !strings.Contains(stderr, "V2 human: missing") {
		t.Fatalf("human pass gate=%v %s %s", err, out, stderr)
	}
	p21Run(t, j.binary, j.root, "--json", "epic", "validation", "pass", j.epic, "V2", "--actor", "human:mark", "--reason", "Followed recipe", "--idempotency-key", "human-pass")
	p21Run(t, j.binary, j.root, "--json", "done", j.epic, "--actor", "human:mark", "--reason", "All validations passed")
}

func TestEpicVerificationSignalCleansUp(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepo(t, true)
	pidFile := filepath.Join(t.TempDir(), "check.pid")
	argv, _ := json.Marshal([]string{"sh", "-c", "echo $$ > \"$1\"; sleep 30", "check", pidFile})
	file, err := os.OpenFile(filepath.Join(root, ".fid.toml"), os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, writeErr := fmt.Fprintf(file, "\n[checks.slow]\nrun=%s\ntimeout='40s'\n", argv)
	closeErr := file.Close()
	if writeErr != nil || closeErr != nil {
		t.Fatalf("policy=%v %v", writeErr, closeErr)
	}
	closureRun(t, root, "git", "add", ".fid.toml")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Slow validation")
	closureRun(t, root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	e := p21Run(t, binary, root, "--json", "issue", "create", "--kind", "epic", "--title", "Interrupted validation", "--scheduling", "prioritised", "--actor", "human:mark", "--idempotency-key", "epic")
	epic := p21PathString(t, e, "result", "issue", "id")
	p21Run(t, binary, root, "--json", "epic", "validation", "add", epic, "--claim", "Slow check", "--verifier", "check:slow", "--actor", "agent:builder", "--idempotency-key", "claim")
	before := closureRunOutput(t, root, "git", "worktree", "list", "--porcelain")
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, "--json", "epic", "validation", "verify", epic, "--actor", "agent:builder", "--idempotency-key", "verify")
	cmd.Dir = root
	var out, errout bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errout
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		if data, err := os.ReadFile(pidFile); err == nil {
			if pid, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil && pid > 1 {
				_ = syscall.Kill(-pid, syscall.SIGKILL)
			}
		}
		data, _, _ := closureRunCommand(root, "git", "worktree", "list", "--porcelain")
		for _, line := range strings.Split(data, "\n") {
			if path, ok := strings.CutPrefix(line, "worktree "); ok && path != root && strings.Contains(path, "fid-epic-verify-") {
				_, _, _ = closureRunCommand(root, "git", "worktree", "remove", "--force", path)
				_ = os.RemoveAll(filepath.Dir(path))
			}
		}
	})
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(pidFile); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("check did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("verification did not record interruption: %v %s", err, errout.String())
	}
	if !strings.Contains(out.String(), `"outcome":"cancelled"`) {
		t.Fatalf("cancelled result=%s", out.String())
	}
	if after := closureRunOutput(t, root, "git", "worktree", "list", "--porcelain"); after != before {
		t.Fatalf("verification worktree leaked: %s", after)
	}
}

func closureRunOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, stderr, err := closureRunCommand(root, args[0], args[1:]...)
	if err != nil {
		t.Fatalf("%v: %v %s", args, err, stderr)
	}
	return out
}

func TestEpicValidationPassStalesOnCompletedMerge(t *testing.T) {
	t.Parallel()
	j := newDoneJourney(t, "standard", true)
	p21Run(t, j.binary, j.root, "--json", "epic", "validation", "add", j.epic, "--claim", "Review complete", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "claim")
	approveJourneyEpic(t, j.root, j.epic)
	p21Run(t, j.binary, j.root, "--json", "epic", "validation", "pass", j.epic, "V1", "--actor", "human:mark", "--reason", "Reviewed current main")
	if out, stderr, err := j.run(t, "--json", "epic", "validation", "verify", j.epic, "--actor", "agent:builder", "--idempotency-key", "open"); err == nil || !strings.Contains(stderr, "no open tasks") {
		t.Fatalf("open epic=%v %s %s", err, out, stderr)
	}
	j.review(t, "primary", "APPROVE", "review")
	head := strings.TrimSpace(closureRunOutput(t, j.root, "git", "rev-parse", "HEAD"))
	j.pr["mergeCommit"] = map[string]any{"oid": head}
	j.ci[0].(map[string]any)["head_sha"] = head
	j.requests[strings.Replace(doneCIRequest, strings.Repeat("a", 40), head, 1)] = map[string]any{"body": map[string]any{"total_count": len(j.ci), "runs": j.ci}}
	if out, stderr, err := j.run(t, "--json", "done", j.task, "--actor", "agent:builder"); err != nil {
		t.Fatalf("done=%v %s %s", err, out, stderr)
	}
	closureRun(t, j.root, "git", "update-ref", "refs/remotes/origin/trunk", head)
	if out, stderr, err := j.run(t, "--json", "done", j.epic, "--actor", "human:mark", "--reason", "Finished"); err == nil || !strings.Contains(stderr, "V1 review:code: stale") {
		t.Fatalf("old pass=%v %s %s", err, out, stderr)
	}
	p21Run(t, j.binary, j.root, "--json", "epic", "validation", "pass", j.epic, "V1", "--actor", "human:mark", "--reason", "Reviewed merged change")
	p21Run(t, j.binary, j.root, "--json", "done", j.epic, "--actor", "human:mark", "--reason", "Finished")
}

func TestEpicPassCoversMergeRecordedAfterPass(t *testing.T) {
	t.Parallel()
	j := newDoneJourney(t, "standard", true)
	p21Run(t, j.binary, j.root, "--json", "epic", "validation", "add", j.epic, "--claim", "Review complete", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", "claim")
	approveJourneyEpic(t, j.root, j.epic)
	j.review(t, "primary", "APPROVE", "review")
	merge := strings.TrimSpace(closureRunOutput(t, j.root, "git", "rev-parse", "HEAD"))
	closureRun(t, j.root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "--allow-empty", "-qm", "Later unrelated main commit")
	closureRun(t, j.root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	p21Run(t, j.binary, j.root, "--json", "epic", "validation", "pass", j.epic, "V1", "--actor", "human:mark", "--reason", "Reviewed main including merge")
	j.pr["mergeCommit"] = map[string]any{"oid": merge}
	j.ci[0].(map[string]any)["head_sha"] = merge
	j.requests[strings.Replace(doneCIRequest, strings.Repeat("a", 40), merge, 1)] = map[string]any{"body": map[string]any{"total_count": len(j.ci), "runs": j.ci}}
	if out, stderr, err := j.run(t, "--json", "done", j.task, "--actor", "agent:builder"); err != nil {
		t.Fatalf("task done=%v %s %s", err, out, stderr)
	}
	shown := p21Run(t, j.binary, j.root, "--json", "show", j.epic)
	validation := p21Path(t, shown, "result", "validations").([]any)[0].(map[string]any)
	if validation["state"] != "current" {
		t.Fatalf("contained merge staled a pass: %v", validation)
	}
	p21Run(t, j.binary, j.root, "--json", "done", j.epic, "--actor", "human:mark", "--reason", "All evidence current", "--idempotency-key", "close")
	objects := strings.TrimSpace(closureRunOutput(t, j.root, "git", "rev-parse", "--git-path", "objects"))
	if !filepath.IsAbs(objects) {
		objects = filepath.Join(j.root, objects)
	}
	if err := os.RemoveAll(objects); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(objects, 0o700); err != nil {
		t.Fatal(err)
	}
	p21Run(t, j.binary, j.root, "--json", "rebuild")
	shown = p21Run(t, j.binary, j.root, "--json", "show", j.epic)
	validation = p21Path(t, shown, "result", "validations").([]any)[0].(map[string]any)
	if validation["state"] != "current" {
		t.Fatalf("offline close evidence lost currency: %v", validation)
	}
}
