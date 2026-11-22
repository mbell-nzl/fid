//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func (c *loopContract) failure(name, binary, root string, args ...string) {
	c.t.Helper()
	out, stderr, err := closureRunCommand(root, binary, append([]string{"--json"}, args...)...)
	if err == nil || out != "" || stderr == "" {
		c.t.Fatalf("%s: %v %q %q", name, err, out, stderr)
	}
	c.observe(name, p23Envelope(c.t, stderr), processExitCode(c.t, err))
}

func TestLoopContractFailures(t *testing.T) {
	t.Parallel()
	c := &loopContract{t: t, decl: readLoopDeclaration(t)}
	root, binary := handoffJourneyRepo(t, true)
	for _, command := range [][]string{{"next"}, {"next", "--claim"}, {"start"}, {"heartbeat"}, {"park"}, {"answer"}, {"check"}, {"handoff"}, {"review", "record"}, {"done"}} {
		c.failure(strings.Join(command, "-")+"-invalid", binary, root, append(command, "--bogus")...)
	}
	c.failure("next-unknown-epic", binary, root, "next", "--epic", "FID-missing")
	c.failure("next-denied", binary, root, "next", "--claim", "--actor", "service:ci")
	c.failure("heartbeat-stale", binary, root, "heartbeat", "--hold", "HLD-missing", "--actor", "agent:builder")
	// Journal corruption and an unusable projection stay isolated in their
	// own temporary repositories. They pin transport exits 6 and 7.
	for _, mode := range []string{"integrity", "internal"} {
		broken, _ := handoffJourneyRepo(t, true)
		ledger := config.LedgerDir(broken)
		if mode == "integrity" {
			entries, err := os.ReadDir(filepath.Join(ledger, "journal"))
			if err != nil {
				t.Fatal(err)
			}
			file, err := os.OpenFile(filepath.Join(ledger, "journal", entries[len(entries)-1].Name()), os.O_APPEND|os.O_WRONLY, 0o600)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := file.WriteString("not-json\n"); err != nil {
				t.Fatal(err)
			}
			if err := file.Close(); err != nil {
				t.Fatal(err)
			}
		} else {
			path := filepath.Join(ledger, config.DatabaseFileName)
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(path, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		c.failure("next-"+mode, binary, broken, "next")
	}
	loopCheckFailure(t, c)
}

func loopCheckFailure(t *testing.T, c *loopContract) {
	t.Helper()
	root, binary := handoffJourneyRepo(t, true)
	path := filepath.Join(root, ".fid.toml")
	raw := readP21File(t, path)
	if err := os.WriteFile(path, []byte(strings.Replace(string(raw), "run=['true']", "run=['false']", 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Failing check")
	closureRun(t, root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	w := loopWorker{t, binary, root, "agent:builder"}
	w.create("task", "Failing check task", "")
	hold := p21PathString(t, w.next("", true), "result", "hold", "id")
	w.call("start", "--hold", hold, "--actor", w.actor, "--idempotency-key", "start")
	c.failure("check-unsuccessful", binary, root, "check", "--hold", hold, "--actor", w.actor, "--idempotency-key", "check")
}
