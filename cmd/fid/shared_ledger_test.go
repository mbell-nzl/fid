package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// TestTwoWorktreesShareOneLedger is RESET task 1.1d3's done-when. The ledger
// lives in the git common directory (ADR-0013), so a repository's main
// worktree and a linked one use the same ledger: init from either creates it
// once, and a write from one is visible from the other.
func TestTwoWorktreesShareOneLedger(t *testing.T) {
	t.Parallel()
	mainWorktree := gittest.NewWorkspace(t)
	linked := gittest.AddWorktree(t, mainWorktree)
	ledger := filepath.Join(gitCommonDir(t, mainWorktree), "fid")

	initialized := executeCommandForTest(linked, "--json", "init", "--prefix", "FID")
	if initialized.err != nil {
		t.Fatalf("init from the linked worktree = %#v", initialized)
	}
	journalPath, _ := decodeCanonicalResult(t, initialized.stdout)["journal_path"].(string)
	if !strings.HasPrefix(journalPath, filepath.ToSlash(ledger)+"/journal/") {
		t.Fatalf("init journal_path = %q, want a journal file in the shared ledger %q", journalPath, ledger)
	}
	for _, root := range []string{mainWorktree, linked} {
		if got := config.LedgerDir(gittest.ResolvedPath(t, root)); got != ledger {
			t.Fatalf("LedgerDir(%q) = %q, want the shared ledger %q", root, got, ledger)
		}
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatal(err)
		}
		var names []string
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		want := ".git"
		if root == linked {
			want = ".fid.toml,.git"
		}
		if got := strings.Join(names, ","); got != want {
			t.Fatalf("worktree %q holds %q, want %q", root, got, want)
		}
	}
	again := executeCommandForTest(mainWorktree, "--json", "init", "--prefix", "FID")
	if cli.ExitCode(again.err) == 0 || !strings.Contains(again.stderr, `"code":"`+string(cli.CodeProjectAlreadyInit)+`"`) {
		t.Fatalf("init from the main worktree after the linked one = %#v, want %s", again, cli.CodeProjectAlreadyInit)
	}

	for _, write := range []struct{ from, to, title string }{
		{from: mainWorktree, to: linked, title: "Written in the main worktree"},
		{from: linked, to: mainWorktree, title: "Written in the linked worktree"},
	} {
		created := executeCommandForTest(write.from, "--json", "issue", "create", "--kind", "task", "--title", write.title, "--actor", "agent:developer")
		if created.err != nil {
			t.Fatalf("issue create %q = %#v", write.title, created)
		}
		issueID, _ := decodeCanonicalResult(t, created.stdout)["issue"].(map[string]any)["id"].(string)
		if issueID == "" {
			t.Fatalf("issue create %q returned no issue id: %s", write.title, created.stdout)
		}
		listed := executeCommandForTest(write.to, "--json", "issue", "list")
		if listed.err != nil || !strings.Contains(listed.stdout, `"id":"`+issueID+`"`) || !strings.Contains(listed.stdout, write.title) {
			t.Fatalf("issue list in the other worktree = %#v, want %s %q", listed, issueID, write.title)
		}
	}
}

// gitCommonDir returns the absolute git common directory of the repository
// that holds directory, as git reports it, with symlinks resolved.
func gitCommonDir(t *testing.T, directory string) string {
	t.Helper()
	output, err := exec.Command("git", "-C", directory, "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	if err != nil {
		t.Fatalf("git rev-parse --git-common-dir in %q: %v", directory, err)
	}
	return strings.TrimSpace(string(output))
}
