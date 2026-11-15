//go:build integration

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func doneCodePair(code string) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		lane := "standard"
		if strings.Contains(code, "approval") || code == "done-separate-review-required" || code == "done-rollback-missing" {
			lane = "assured"
		}
		j := newDoneJourney(t, lane)
		if code != "done-review-required" {
			j.review(t, "primary", "APPROVE", "primary")
		}
		if lane == "assured" && code != "done-separate-review-required" {
			j.review(t, "separate", "APPROVE", "separate")
		}
		if lane == "assured" && code != "done-approval-missing" && code != "done-separate-review-required" {
			j.approve(t, "approval")
		}
		switch code {
		case "done-task-not-in-review":
			p21Run(t, j.binary, j.root, "--json", "issue", "transition", j.task, "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed")
		case "done-handoff-required":
			p21Run(t, j.binary, j.root, "--json", "issue", "update", j.task, "--lane", "assured", "--actor", "human:mark", "--idempotency-key", "raised")
		case "done-checks-unsatisfied":
			conn := p22FaultDB(t, j.root)
			_, err := conn.Exec(`DELETE FROM check_results`)
			p22FaultOK(t, err)
			p22FaultOK(t, conn.Close())
		case "done-evidence-required":
			conn := p22FaultDB(t, j.root)
			_, err := conn.Exec(`UPDATE policy_snapshots SET effective_json=json_set(effective_json,'$.checks',json('{}'))`)
			p22FaultOK(t, err)
			p22FaultOK(t, conn.Close())
		case "done-kept-artifact-unavailable":
			files, err := os.ReadDir(filepath.Join(config.LedgerDir(j.root), "artifacts"))
			if err != nil || len(files) != 1 {
				t.Fatalf("files=%v/%v", files, err)
			}
			if err := os.Remove(filepath.Join(config.LedgerDir(j.root), "artifacts", files[0].Name())); err != nil {
				t.Fatal(err)
			}
		case "done-not-merged":
			j.pr["state"] = "OPEN"
			j.pr["mergedAt"] = nil
			j.pr["mergeCommit"] = nil
		case "done-ci-waiting":
			j.ci = []any{}
		case "done-ci-not-green":
			j.ci[0].(map[string]any)["conclusion"] = "failure"
		case "github-unavailable":
			j.requests[donePRRequest] = map[string]any{"error": "network unavailable"}
		case "done-rollback-missing":
			j.pr["body"] = "No heading"
		case "done-open-major":
			j.major(t)
			j.rework(t, false)
		case "done-approval-stale":
			j.major(t)
			j.rework(t, true)
			conn := p22FaultDB(t, j.root)
			var finding string
			p22FaultOK(t, conn.QueryRow(`SELECT id FROM round_findings LIMIT 1`).Scan(&finding))
			p22FaultOK(t, conn.Close())
			p21Run(t, j.binary, j.root, "--json", "finding", "close", finding, "--actor", "human:mark", "--reason", "Verified the fix", "--idempotency-key", "close")
			j.review(t, "primary", "APPROVE", "fixed-primary")
			j.review(t, "separate", "APPROVE", "fixed-separate")
		}
		args := []string{"done", j.task, "--actor", "agent:finisher", "--idempotency-key", "finish"}
		var envelope map[string]any
		for _, jsonMode := range []bool{true, false} {
			if code == "done-observation-stale" {
				j.action = []string{j.binary, "--project", j.root, "issue", "update", j.task, "--title", fmt.Sprint("Changed during observation ", jsonMode), "--actor", "human:mark", "--idempotency-key", fmt.Sprint("mutate-", jsonMode)}
			}
			invoke := args
			if jsonMode {
				invoke = append([]string{"--json"}, args...)
			}
			out, stderr, err := j.run(t, invoke...)
			if err == nil || out != "" || processExitCode(t, err) != 4 {
				t.Fatalf("%s=%v %s %s", code, err, out, stderr)
			}
			if jsonMode {
				envelope = p23Envelope(t, stderr)
				if p21PathString(t, envelope, "error", "code") != code || envelope["operation"] != "task.done" {
					t.Fatalf("failure=%v", envelope)
				}
			} else {
				expected := p24ExpectedFailureText(t, envelope, "task.done", false)
				reason, _ := p21Path(t, envelope, "error", "details", "reason").(string)
				next := doneRepairForTest[reason]
				if code == "github-unavailable" {
					next = "check gh authentication and connectivity, then retry fid done"
				}
				expected += "next: " + next + "\n"
				if stderr != expected {
					t.Fatalf("text=%q want=%q", stderr, expected)
				}
			}
		}
		conn := p22FaultDB(t, j.root)
		var n int
		p22FaultOK(t, conn.QueryRow(`SELECT count(*) FROM completions`).Scan(&n))
		p22FaultOK(t, conn.Close())
		if n != 0 {
			t.Fatal("refusal completed task")
		}
	}
}

var doneRepairForTest = map[string]string{
	"task-not-in-review":        "start, check and hand off the task before fid done",
	"lane-changed":              "hand off again at the current lane",
	"checks-unsatisfied":        "check the change and hand off again",
	"policy-mismatch":           "hand off again",
	"kept-artifact-unavailable": "restore their bytes or check and hand off again",
	"not-merged":                "merge the reviewed PR, then retry fid done",
	"ci-waiting":                "wait for CI to finish, then retry fid done",
	"ci-not-green":              "repair the failing CI check, then retry fid done",
	"review-missing":            "record the required review profiles",
	"approval-missing":          "approve the change with fid approve",
	"approval-stale":            "approve the current change with fid approve",
	"rollback-missing":          "add a level-two Rollback heading to the PR body",
	"open-major":                "close or resolve the major finding",
	"observation-stale":         "inspect the current task, then retry fid done",
}

func (j *doneJourney) major(t *testing.T) {
	t.Helper()
	in := map[string]any{"reviewer": map[string]any{"actor": "agent:reviewer", "model": "claude/opus-5.5", "session_id": "major-session"}, "stream": "primary", "profiles": []string{"code"}, "verdict": "CHANGES", "purpose": "Review the current change", "findings": []any{map[string]any{"profile": "code", "severity": "major", "summary": "Missing behavior", "location": nil, "evidence": "The task needs a fix", "fix": "Fix the behavior"}}, "verification": "Inspected the tests."}
	path := p21Input(t, t.TempDir(), "major", in)
	p21Run(t, j.binary, j.root, "--json", "review", "record", j.task, "--input", path, "--actor", "human:courier", "--idempotency-key", "major")
}
func (j *doneJourney) rework(t *testing.T, edit bool) {
	t.Helper()
	took := p21Run(t, j.binary, j.root, "--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol")
	hold := p21PathString(t, took, "result", "hold", "id")
	if edit {
		if err := os.WriteFile(filepath.Join(j.root, "change.txt"), []byte("Edited task change\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		closureRun(t, j.root, "git", "add", ".")
		closureRun(t, j.root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Fix the change")
	}
	p21Run(t, j.binary, j.root, "--json", "check", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "recheck")
	h := p21Run(t, j.binary, j.root, "--json", "handoff", "--hold", hold, "--actor", "agent:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "rehandoff")
	head := p21PathString(t, h, "result", "head")
	tree, _, err := closureRunCommand(j.root, "git", "rev-parse", head+"^{tree}")
	if err != nil {
		t.Fatal(err)
	}
	j.pr["headRefOid"] = head
	j.requests["api --hostname github.com --method GET repos/mbell-nzl/fid/git/commits/"+head+" --jq {sha, tree: .tree.sha}"] = map[string]any{"body": map[string]any{"sha": head, "tree": strings.TrimSpace(tree)}}
}

func TestDoneWaitingRetryReobservesThroughBuiltBinary(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"open", "no-ci", "pending", "failed"} {
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			j := newDoneJourney(t, "standard")
			j.review(t, "primary", "APPROVE", "primary")
			code := "done-ci-waiting"
			switch state {
			case "open":
				j.pr["state"] = "OPEN"
				j.pr["mergedAt"] = nil
				j.pr["mergeCommit"] = nil
				code = "done-not-merged"
			case "no-ci":
				j.ci = []any{}
			case "pending":
				j.ci[0].(map[string]any)["status"] = "in_progress"
				j.ci[0].(map[string]any)["conclusion"] = ""
			case "failed":
				j.ci[0].(map[string]any)["conclusion"] = "failure"
				code = "done-ci-not-green"
			}
			for range 2 {
				_, stderr, err := j.run(t, "--json", "done", j.task, "--actor", "agent:finisher", "--idempotency-key", "retry")
				if err == nil || p21PathString(t, p23Envelope(t, stderr), "error", "code") != code {
					t.Fatalf("waiting=%v %s", err, stderr)
				}
			}
			log, err := os.ReadFile(j.log)
			if err != nil {
				t.Fatal(err)
			}
			want := 6
			if state == "open" {
				want = 2
			}
			if len(strings.Split(strings.TrimSpace(string(log)), "\n")) != want {
				t.Fatalf("calls=%s", log)
			}
		})
	}
}
