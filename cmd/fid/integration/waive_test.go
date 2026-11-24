//go:build integration && linux

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWaiverPublicJourney(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepo(t, true)
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, root, append([]string{"--json"}, args...)...)
	}
	task := p21PathString(t, call("issue", "create", "--kind", "task", "--title", "Waiver journey", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "create"), "result", "issue", "id")
	hold := p21PathString(t, call("takeover", task, "--actor", "human:mark", "--reason", "Build journey"), "result", "hold_id")
	call("start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "start")
	if err := os.WriteFile(filepath.Join(root, "notes.md"), []byte("Waive this change\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Task change")
	call("check", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "check")
	call("handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/owner/repo/pull/1", "--idempotency-key", "handoff")
	args := []string{"--json", "waive", task, "--lane", "routine", "--actor", "human:mark", "--reason", "Reviewed risk", "--idempotency-key", "waive"}
	first, stderr, err := approvalCorpusCommand(t, root, binary, args...)
	if err != nil || stderr != "" {
		t.Fatal(first, stderr, err)
	}
	again, stderr, err := approvalCorpusCommand(t, root, binary, args...)
	if err != nil || stderr != "" || again != first {
		t.Fatal("retry changed reply", again, stderr, err)
	}
	for _, rebuild := range []bool{false, true} {
		if rebuild {
			call("rebuild")
		}
		var envelope struct {
			Result struct {
				Waivers []struct {
					Current bool `json:"current"`
				} `json:"waivers"`
			} `json:"result"`
		}
		raw, err := json.Marshal(call("show", task))
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(raw, &envelope); err != nil {
			t.Fatal(err)
		}
		if len(envelope.Result.Waivers) != 1 || !envelope.Result.Waivers[0].Current {
			t.Fatal("current waiver lost", envelope)
		}
	}
}

func waiverCodePair(code string) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		root, binary := handoffJourneyRepo(t, true)
		call := func(args ...string) map[string]any {
			t.Helper()
			return p21Run(t, binary, root, append([]string{"--json"}, args...)...)
		}
		if code == "waiver-hard-floor" {
			path := filepath.Join(root, ".fid.toml")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, append(raw, []byte("\n[lanes.assured]\nchecks=[]\n")...), 0o600); err != nil {
				t.Fatal(err)
			}
			closureRun(t, root, "git", "add", ".")
			closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Assured policy")
			closureRun(t, root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
		}
		task := p21PathString(t, call("issue", "create", "--kind", "task", "--title", "Waiver refusal", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "create"), "result", "issue", "id")
		switch code {
		case "waiver-handoff-required":
			conn := p22FaultDB(t, root)
			_, err := conn.Exec(`UPDATE issues SET status='in_review' WHERE id=?`, task)
			p22FaultOK(t, err)
			p22FaultOK(t, conn.Close())
		case "waiver-not-lower", "waiver-hard-floor":
			hold := p21PathString(t, call("takeover", task, "--actor", "human:mark", "--reason", "Build refusal"), "result", "hold_id")
			call("start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "start")
			path := "notes.md"
			if code == "waiver-hard-floor" {
				path = ".github/workflows/task.yml"
				if err := os.MkdirAll(filepath.Join(root, ".github/workflows"), 0o700); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.WriteFile(filepath.Join(root, path), []byte("task\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			closureRun(t, root, "git", "add", ".")
			closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Task change")
			call("check", "test", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "check")
			call("handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/owner/repo/pull/1", "--idempotency-key", "handoff")
		}
		args := []string{"waive", task, "--lane", "standard", "--reason", "Reviewed risk", "--actor", "human:mark"}
		out, stderr, err := approvalCorpusCommand(t, root, binary, append([]string{"--json"}, args...)...)
		if err == nil || out != "" {
			t.Fatal("JSON waiver refusal", out, stderr, err)
		}
		payload := p23Envelope(t, stderr)
		if p21PathString(t, payload, "error", "code") != code || payload["operation"] != "task.waive" {
			t.Fatal(payload)
		}
		status := processExitCode(t, err)
		out, stderr, err = approvalCorpusCommand(t, root, binary, args...)
		want := p24ExpectedFailureText(t, payload, "task.waive", false)
		if err == nil || processExitCode(t, err) != status || out != "" || stderr != want {
			t.Fatalf("human refusal: %v %q want %q", err, stderr, want)
		}
	}
}
