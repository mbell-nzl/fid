//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestEpicDoneOfflineJourney(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepoWithArtifacts(t, true, true)
	created := p21Run(t, binary, root, "--json", "issue", "create", "--kind", "epic", "--scheduling", "prioritised", "--title", "Close locally", "--actor", "human:mark", "--idempotency-key", "epic")
	epic := p21PathString(t, created, "result", "issue", "id")
	created = p21Run(t, binary, root, "--json", "issue", "create", "--kind", "task", "--epic", epic, "--title", "Finished child", "--actor", "human:mark", "--idempotency-key", "child")
	child := p21PathString(t, created, "result", "issue", "id")
	p21Run(t, binary, root, "--json", "issue", "transition", child, "--to", "cancelled", "--actor", "human:mark", "--reason", "Not needed")
	bin := t.TempDir()
	log := filepath.Join(bin, "network.log")
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	for name, script := range map[string]string{
		"gh":  "#!/bin/sh\necho gh >> \"$FID_TEST_NETWORK_LOG\"\nexit 97\n",
		"git": "#!/bin/sh\nfor arg do\n if [ \"$arg\" = fetch ]; then echo fetch >> \"$FID_TEST_NETWORK_LOG\"; exit 97; fi\ndone\nexec '" + git + "' \"$@\"\n",
	} {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(script), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	run := func(args ...string) string {
		t.Helper()
		cmd := exec.CommandContext(t.Context(), binary, args...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "FID_TEST_NETWORK_LOG="+log)
		raw, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%v=%v: %s", args, err, raw)
		}
		return string(raw)
	}
	if got := run("--json", "status"); !strings.Contains(got, `"kind":"epic-ready-to-close"`) {
		t.Fatalf("status=%s", got)
	}
	args := []string{"--json", "done", epic, "--actor", "human:mark", "--reason", "All tasks finished", "--idempotency-key", "close"}
	original := run(args...)
	if !strings.Contains(original, `"kind":"close-epic"`) {
		t.Fatalf("close=%s", original)
	}
	if retry := run(args...); retry != original {
		t.Fatalf("retry=%s want %s", retry, original)
	}
	run("rebuild")
	if retry := run(args...); retry != original {
		t.Fatalf("rebuilt retry=%s want %s", retry, original)
	}
	if got := run("--json", "status"); strings.Contains(got, `"kind":"epic-ready-to-close"`) {
		t.Fatalf("closed epic still needs person: %s", got)
	}
	if _, err := os.Stat(log); !os.IsNotExist(err) {
		t.Fatalf("offline journey invoked network: %v", err)
	}
}

func epicDoneCodePair(code string) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		root, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
		p21Run(t, binary, root, "init", "--prefix", "FID", "--json")
		created := p21Run(t, binary, root, "--json", "issue", "create", "--kind", "epic", "--scheduling", "prioritised", "--title", "Epic", "--actor", "human:mark", "--idempotency-key", "epic")
		epic := p21PathString(t, created, "result", "issue", "id")
		if code == "done-epic-not-open" {
			p21Run(t, binary, root, "--json", "done", epic, "--actor", "human:mark", "--reason", "Finished")
		} else {
			p21Run(t, binary, root, "--json", "issue", "create", "--kind", "task", "--epic", epic, "--title", "Open task", "--actor", "human:mark", "--idempotency-key", "child")
		}
		args := []string{"done", epic, "--actor", "human:mark", "--reason", "Finished"}
		out, stderr, err := closureRunCommand(root, binary, append([]string{"--json"}, args...)...)
		if err == nil || out != "" || processExitCode(t, err) != 4 {
			t.Fatalf("JSON %s=%v %s %s", code, err, out, stderr)
		}
		envelope := p23Envelope(t, stderr)
		if p21PathString(t, envelope, "error", "code") != code || envelope["operation"] != "task.done" {
			t.Fatalf("refusal=%v", envelope)
		}
		out, stderr, err = closureRunCommand(root, binary, args...)
		want := p24ExpectedFailureText(t, envelope, "task.done", false) + "next: fid show " + epic + "\n"
		if err == nil || out != "" || processExitCode(t, err) != 4 || stderr != want {
			t.Fatalf("text %s=%v %s %q want %q", code, err, out, stderr, want)
		}
	}
}
