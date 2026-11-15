//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

type doneJourney struct {
	root, binary, task, hold, bin, responses, log string
	pr                                            map[string]any
	ci                                            []any
	requests                                      map[string]any
	action                                        []string
}

func newDoneJourney(t *testing.T, lane string) *doneJourney {
	t.Helper()
	root, binary := handoffJourneyRepoWithArtifacts(t, true, true)
	file := filepath.Join(root, ".fid.toml")
	raw, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, append(raw, []byte("\n[lanes.assured]\nchecks=['test']\n")...), 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Assured test policy")
	closureRun(t, root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	created := p21Run(t, binary, root, "--json", "issue", "create", "--kind", "task", "--title", "Finish the task", "--lane", lane, "--actor", "agent:builder", "--idempotency-key", "create")
	task := p21PathString(t, created, "result", "issue", "id")
	took := p21Run(t, binary, root, "--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol")
	hold := p21PathString(t, took, "result", "hold", "id")
	p21Run(t, binary, root, "--json", "start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start")
	if err := os.WriteFile(filepath.Join(root, "change.txt"), []byte("Task change\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Task change")
	p21Run(t, binary, root, "--json", "check", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "check")
	h := p21Run(t, binary, root, "--json", "handoff", "--hold", hold, "--actor", "agent:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "handoff")
	return newDoneObservation(t, root, binary, task, hold, h)
}

func newDoneObservation(t *testing.T, root, binary, task, hold string, h map[string]any) *doneJourney {
	t.Helper()
	head := p21PathString(t, h, "result", "head")
	base, _, err := closureRunCommand(root, "git", "rev-parse", "refs/remotes/origin/trunk")
	if err != nil {
		t.Fatal(err)
	}
	tree, _, err := closureRunCommand(root, "git", "rev-parse", head+"^{tree}")
	if err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "remote", "set-url", "origin", "https://github.com/mbell-nzl/fid.git")
	bin := t.TempDir()
	script, err := os.ReadFile("../testdata/gh_stub.py")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bin, "gh"), script, 0o700); err != nil {
		t.Fatal(err)
	}
	realgit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	script = []byte("#!/bin/sh\nfor arg do\n if [ \"$arg\" = fetch ]; then echo 'fetch is forbidden' >&2; exit 97; fi\ndone\nexec '" + realgit + "' \"$@\"\n")
	if err := os.WriteFile(filepath.Join(bin, "git"), script, 0o700); err != nil {
		t.Fatal(err)
	}
	merge := strings.Repeat("a", 40)
	j := &doneJourney{root: root, binary: binary, task: task, hold: hold, bin: bin, responses: filepath.Join(bin, "responses.json"), log: filepath.Join(bin, "requests.log")}
	j.pr = map[string]any{"number": 63, "url": "https://github.com/mbell-nzl/fid/pull/63", "state": "MERGED", "mergedAt": time.Now().UTC().Format(time.RFC3339), "mergeCommit": map[string]any{"oid": merge}, "headRefOid": head, "headRefName": "task", "baseRefName": "trunk", "baseRefOid": strings.TrimSpace(base), "isCrossRepository": false, "body": "## Rollback\nRevert the change.\n"}
	j.ci = []any{map[string]any{"id": 101, "run_attempt": 1, "workflow_id": 7, "name": "CI", "path": ".github/workflows/ci.yml", "event": "push", "head_sha": merge, "status": "completed", "conclusion": "success", "html_url": "https://github.com/mbell-nzl/fid/actions/runs/101"}}
	j.requests = map[string]any{
		"api --hostname github.com --method GET repos/mbell-nzl/fid/git/commits/" + head + " --jq {sha, tree: .tree.sha}": map[string]any{"body": map[string]any{"sha": head, "tree": strings.TrimSpace(tree)}},
	}
	return j
}

const donePRRequest = "pr view 63 --repo mbell-nzl/fid --json number,url,state,mergedAt,mergeCommit,headRefOid,headRefName,baseRefName,baseRefOid,isCrossRepository,body"
const doneCIRequest = "api --hostname github.com --method GET repos/mbell-nzl/fid/actions/runs?head_sha=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa&per_page=100&page=1 --jq {total_count, runs: [.workflow_runs[] | {id, run_attempt, workflow_id, name, path, event, head_sha, status, conclusion, html_url}]}"

func (j *doneJourney) run(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	if j.requests[donePRRequest] == nil {
		j.requests[donePRRequest] = map[string]any{"body": j.pr}
	}
	j.requests[doneCIRequest] = map[string]any{"body": map[string]any{"total_count": len(j.ci), "runs": j.ci}}
	raw, err := json.Marshal(map[string]any{"requests": j.requests, "action": j.action})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(j.responses, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(t.Context(), j.binary, args...)
	cmd.Dir = j.root
	cmd.Env = append(os.Environ(), "PATH="+j.bin+string(os.PathListSeparator)+os.Getenv("PATH"), "FID_TEST_GH_RESPONSES="+j.responses, "FID_TEST_GH_LOG="+j.log)
	var out, errout bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errout
	err = cmd.Run()
	p24CaptureKeyedCommand(j.root, args, out.String(), errout.String(), err)
	return out.String(), errout.String(), err
}
func (j *doneJourney) review(t *testing.T, stream, verdict, key string) {
	t.Helper()
	in := map[string]any{"reviewer": map[string]any{"actor": "agent:reviewer", "model": "claude/opus-5.5", "session_id": key}, "stream": stream, "profiles": []string{"code"}, "verdict": verdict, "purpose": "Review the current change", "findings": []any{}, "verification": "Inspected the change and tests."}
	path := p21Input(t, t.TempDir(), "review", in)
	p21Run(t, j.binary, j.root, "--json", "review", "record", j.task, "--input", path, "--actor", "human:courier", "--idempotency-key", key)
}
func (j *doneJourney) approve(t *testing.T, key string) {
	t.Helper()
	_, stderr, err := approvalCorpusCommand(t, j.root, j.binary, "--json", "approve", j.task, "--actor", "human:mark", "--reason", "Reviewed the full change", "--idempotency-key", key)
	if err != nil {
		t.Fatalf("approve=%v %s", err, stderr)
	}
}

func TestDoneBuiltBinaryStandardAndAssuredJourney(t *testing.T) {
	t.Parallel()
	for _, lane := range []string{"standard", "assured"} {
		t.Run(lane, func(t *testing.T) {
			t.Parallel()
			j := newDoneJourney(t, lane)
			j.review(t, "primary", "APPROVE", "primary")
			if lane == "assured" {
				j.review(t, "separate", "APPROVE", "separate")
				j.approve(t, "approval")
			}
			p24StartKeyedCapture(j.root)
			args := []string{"--json", "done", j.task, "--actor", "agent:finisher", "--idempotency-key", "finish"}
			first, stderr, err := j.run(t, args...)
			if err != nil || stderr != "" {
				t.Fatalf("done=%v %s", err, stderr)
			}
			r := p23Envelope(t, first)
			if p21PathString(t, r, "result", "status") != "done" {
				t.Fatalf("done=%v", r)
			}
			log, err := os.ReadFile(j.log)
			if err != nil || len(strings.Split(strings.TrimSpace(string(log)), "\n")) != 3 {
				t.Fatalf("calls=%s/%v", log, err)
			}
			entry := p24EndKeyedCapture(j.root)["task.done"]
			assertDoneKeyedRetries(t, j, entry)
			p21Run(t, j.binary, j.root, "rebuild", "--json")
			shown := p21Run(t, j.binary, j.root, "--json", "show", j.task)
			if p21PathString(t, shown, "result", "issue", "status") != "done" {
				t.Fatalf("show=%v", shown)
			}
			next := p21Run(t, j.binary, j.root, "--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol")
			if p21Path(t, next, "result", "waiting") != nil {
				t.Fatalf("reservation=%v", next)
			}
		})
	}
}

func assertDoneKeyedRetries(t *testing.T, j *doneJourney, entry p24CapturedKeyedCommand) {
	t.Helper()
	args := entry.args
	p24AssertStoredFieldNames(t, j.root, "task.done", entry)
	before := p21JournalEventCount(t, j.root)
	log, err := os.ReadFile(j.log)
	if err != nil {
		t.Fatal(err)
	}
	retry, stderr, err := j.run(t, args...)
	if err != nil || stderr != "" || retry != entry.stdout {
		t.Fatalf("retry=%v %s %s", err, retry, stderr)
	}
	for _, field := range []string{"task", "--actor"} {
		changed := append([]string(nil), args...)
		for i, arg := range changed {
			if (field == "task" && arg == "done") || arg == field {
				changed[i+1] += "-changed"
				break
			}
		}
		out, stderr, err := j.run(t, changed...)
		if err == nil || out != "" || processExitCode(t, err) != 4 {
			t.Fatalf("conflict=%v %s %s", err, out, stderr)
		}
		envelope := p23Envelope(t, stderr)
		if p21PathString(t, envelope, "error", "code") != "idempotency-conflict" || !reflect.DeepEqual(p21Path(t, envelope, "error", "details", "differing_fields"), []any{field}) {
			t.Fatalf("conflict=%v", envelope)
		}
	}
	humanArgs := []string{}
	for _, arg := range args {
		if arg != "--json" {
			humanArgs = append(humanArgs, arg)
		}
	}
	human, stderr, err := j.run(t, humanArgs...)
	if err != nil || stderr != "" || human != "✓ Completed "+j.task+" · merged with CI green\n" {
		t.Fatalf("human success=%v %q %q", err, human, stderr)
	}
	if p21JournalEventCount(t, j.root) != before {
		t.Fatal("retries or conflicts appended Events")
	}
	log2, _ := os.ReadFile(j.log)
	if !bytes.Equal(log, log2) {
		t.Fatal("retry observed GitHub")
	}
}
