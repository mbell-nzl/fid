package main

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

// withRemoteHelper installs a git remote helper called git-remote-<name> on
// PATH and returns the environment that finds it.
func withRemoteHelper(t *testing.T, fixture *backupFixture, name, script string) []string {
	t.Helper()
	bin := filepath.Join(fixture.base, "bin")
	writeBackupTestFile(t, filepath.Join(bin, "git-remote-"+name), script)
	if err := os.Chmod(filepath.Join(bin, "git-remote-"+name), 0o755); err != nil {
		t.Fatal(err)
	}
	return append(append([]string{}, fixture.env...), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// A ledger can appear after restore's first check: another restore, or a
// fid init, finishing while this one fetches. Restore then refuses as if the
// ledger had been there all along, and leaves it as it is. A remote helper
// stands in for the target's transport and puts the ledger in place as the
// restore connects.
func TestRestoreRefusesALedgerThatAppearsWhileItFetches(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		// appear readies a ledger beside the workspace's, and returns the
		// shell command that moves it into place and the ledger's contents.
		appear func(t *testing.T, fixture *backupFixture, workspace string) (string, map[string]workspaceSnapshotEntry)
	}{
		{"the same backup restored meanwhile", func(t *testing.T, fixture *backupFixture, workspace string) (string, map[string]workspaceSnapshotEntry) {
			decodeRestoreSuccess(t, fixture.restore(workspace, "--json"))
			aside := filepath.Join(fixture.base, "aside")
			if err := os.Rename(config.LedgerDir(workspace), aside); err != nil {
				t.Fatal(err)
			}
			return "mv '" + aside + "' '" + config.LedgerDir(workspace) + "'", snapshotWorkspace(t, aside)
		}},
		{"another ledger made meanwhile", func(t *testing.T, fixture *backupFixture, workspace string) (string, map[string]workspaceSnapshotEntry) {
			return "mkdir '" + config.LedgerDir(workspace) + "'", map[string]workspaceSnapshotEntry{}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newBackupFixture(t)
			original, id := fixture.project()
			fixture.setTarget(fixture.target)
			decodeBackupSuccess(t, fixture.backup(original, "--json"))
			workspace := fixture.restoreWorkspace(id)
			move, contents := test.appear(t, fixture, workspace)
			marker := filepath.Join(fixture.base, "appeared")
			environ := withRemoteHelper(t, fixture, "fidappear", strings.Join([]string{
				"#!/bin/sh",
				"while read -r command; do",
				"  case \"$command\" in",
				"  capabilities) printf 'connect\\n\\n' ;;",
				"  'connect git-upload-pack')",
				"    if [ ! -e '" + marker + "' ]; then touch '" + marker + "'; " + move + " || exit 1; fi",
				"    printf '\\n'; exec git upload-pack \"$2\" ;;",
				"  *) exit 1 ;;",
				"  esac",
				"done",
				"",
			}, "\n"))
			fixture.setTarget("fidappear::" + fixture.target)

			failure := decodeRestoreFailure(t, restoreWith(t, environ, workspace, "--json"), 4)
			if failure.Code != "project-already-initialized" {
				t.Fatalf("failure = %+v, want project-already-initialized", failure)
			}
			if _, err := os.Lstat(marker); err != nil {
				t.Fatalf("the ledger never appeared: %v", err)
			}
			if appeared := snapshotWorkspace(t, config.LedgerDir(workspace)); !reflect.DeepEqual(appeared, contents) {
				t.Fatal("restore changed the ledger that appeared")
			}
		})
	}
}

// Restore fetches into a temporary bare repository of its own: private,
// made without the user's init template, deleted afterwards, and never
// fid's backup clone. Git runs there with the user's credential setup but
// never prompts. A remote helper records what git gives it, then fails like
// an unreachable target.
func TestRestoreFetchesIntoAPrivateTemporaryRepository(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	seen := filepath.Join(fixture.base, "helper-saw")
	environ := withRemoteHelper(t, fixture, "fidprobe", "#!/bin/sh\n"+
		"{ env; echo \"PARENT=$(ls -ld \"$(dirname \"$GIT_DIR\")\" | cut -c1-10)\"; echo \"CONTENTS=$(ls -A \"$GIT_DIR\" | tr '\\n' ' ')\"; } > '"+seen+"'\nexit 1\n")
	writeBackupTestFile(t, filepath.Join(fixture.home, "template", "fid-template-marker"), "from the user's template\n")
	writeBackupTestFile(t, filepath.Join(fixture.home, "template", "hooks", "reference-transaction"), "#!/bin/sh\nexit 1\n")
	writeBackupTestFile(t, filepath.Join(fixture.home, ".gitconfig"), "[init]\n\ttemplateDir = "+filepath.Join(fixture.home, "template")+"\n")
	fixture.setTarget("fidprobe::ledgers")
	temporary := filepath.Join(t.TempDir(), "does-not-exist")
	environ = append(environ, "TMPDIR="+temporary, "GIT_TERMINAL_PROMPT=1", "SSH_AUTH_SOCK=/fid-test/agent.sock")
	workspace := fixture.restoreWorkspace("project-" + strings.Repeat("e", 32))

	failure := decodeRestoreFailure(t, runRestoreWith(environ, workspace, "--json"), 4)
	if failure.Code != "repository-unavailable" || failure.Details["reason"] != "not-a-repository" {
		t.Fatalf("failure = %+v, want the helper's failure as an unreachable target", failure)
	}
	payload, err := os.ReadFile(seen)
	if err != nil {
		t.Fatalf("git never ran the remote helper: %v", err)
	}
	lines := strings.Split(string(payload), "\n")
	gitDir := ""
	for _, line := range lines {
		if value, ok := strings.CutPrefix(line, "GIT_DIR="); ok {
			gitDir = value
		}
	}
	common := filepath.Dir(config.LedgerDir(workspace))
	if parent := filepath.Dir(gitDir); filepath.Dir(parent) != common || !strings.HasPrefix(filepath.Base(parent), "fid-restore-") || filepath.Base(gitDir) != "repo.git" {
		t.Fatalf("restore fetched into %q, want a fid-restore-* directory in %s", gitDir, common)
	}
	for _, want := range []string{"PARENT=drwx------", "GIT_TERMINAL_PROMPT=0", "SSH_AUTH_SOCK=/fid-test/agent.sock"} {
		if !slices.Contains(lines, want) {
			t.Errorf("the helper did not see %q:\n%s", want, payload)
		}
	}
	for _, line := range lines {
		if contents, ok := strings.CutPrefix(line, "CONTENTS="); ok && (strings.Contains(contents, "fid-template-marker") || strings.Contains(contents, "hooks")) {
			t.Errorf("the temporary repository was made from the user's template: %s", contents)
		}
	}
	if slices.Contains(lines, "GIT_TERMINAL_PROMPT=1") {
		t.Error("git's environment still allows terminal prompts")
	}
	if left, err := filepath.Glob(filepath.Join(common, "fid-restore-*")); err != nil || len(left) != 0 {
		t.Fatalf("restore left scratch directories %v: %v", left, err)
	}
	if _, err := os.Lstat(temporary); !os.IsNotExist(err) {
		t.Fatalf("restore used TMPDIR: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(fixture.cacheHome, "fid")); !os.IsNotExist(err) {
		t.Fatalf("restore used fid's cache: %v", err)
	}
}

// Restore fetches only the tip of the target's main, so its time doesn't
// grow with the backup history. A remote helper stands in for the target's
// transport. Git closes the connection once the fetched pack and the shallow
// boundary are in the temporary repository, and only then does the helper
// record which commits the repository holds and where its history stops.
func TestRestoreFetchesOnlyTheTipOfMain(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	original, id := fixture.project()
	fixture.setTarget(fixture.target)
	for backups := 1; backups <= 3; backups++ {
		if backups > 1 {
			runFIDCommand(t, original, "issue", "create", "--kind", "task", "--title", "Backed up again", "--actor", "agent:developer")
		}
		if pushed, _ := decodeBackupSuccess(t, fixture.backup(original, "--json")); pushed.Outcome != "pushed" {
			t.Fatalf("backup %d = %+v, want a new commit on main", backups, pushed)
		}
	}
	tip := fixture.targetMain()
	if count := fixture.targetGit("rev-list", "--count", "main"); count != "3" {
		t.Fatalf("target main has %s commits, want the 3 backups", count)
	}
	seen := filepath.Join(fixture.base, "fetched")
	environ := withRemoteHelper(t, fixture, "fidtip", strings.Join([]string{
		"#!/bin/sh",
		"while read -r command; do",
		"  case \"$command\" in",
		"  capabilities) printf 'connect\\n\\n' ;;",
		"  'connect git-upload-pack')",
		"    printf '\\n'; git upload-pack \"$2\" || exit 1",
		"    cat >/dev/null",
		"    { echo \"SHALLOW=$(git --git-dir=\"$GIT_DIR\" rev-parse --is-shallow-repository)\"",
		"      echo \"BOUNDARY=$(cat \"$GIT_DIR/shallow\" 2>/dev/null | tr '\\n' ' ')\"",
		"      echo \"COMMITS=$(git --git-dir=\"$GIT_DIR\" cat-file --batch-all-objects --batch-check='%(objecttype) %(objectname)' | sed -n 's/^commit //p' | tr '\\n' ' ')\"",
		"    } > '" + seen + "'",
		"    exit 0 ;;",
		"  *) exit 1 ;;",
		"  esac",
		"done",
		"",
	}, "\n"))
	fixture.setTarget("fidtip::" + fixture.target)

	restored, _ := decodeRestoreSuccess(t, restoreWith(t, environ, fixture.restoreWorkspace(id), "--json"))
	if restored.Commit != tip {
		t.Fatalf("restore = %+v, want the target's main %s", restored, tip)
	}
	payload, err := os.ReadFile(seen)
	if err != nil {
		t.Fatalf("git never fetched through the remote helper: %v", err)
	}
	fetched := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(payload)), "\n") {
		if name, value, ok := strings.Cut(line, "="); ok {
			fetched[name] = strings.TrimSpace(value)
		}
	}
	if fetched["SHALLOW"] != "true" || fetched["BOUNDARY"] != tip || fetched["COMMITS"] != tip {
		t.Fatalf("restore's temporary repository after the fetch:\n%s\nwant a shallow repository holding only main's tip %s, its history stopping there", payload, tip)
	}
}

// Hooks and a missing identity in the user's git config must not stop or
// change a restore.
func TestRestoreIgnoresHostileGlobalGitConfig(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	original, id := fixture.project()
	fixture.setTarget(fixture.target)
	decodeBackupSuccess(t, fixture.backup(original, "--json"))
	marker := filepath.Join(fixture.base, "hook-ran")
	for _, name := range []string{"reference-transaction", "post-checkout", "post-merge", "pre-auto-gc"} {
		writeBackupTestFile(t, filepath.Join(fixture.home, "hooks", name), "#!/bin/sh\ntouch '"+marker+"'\nexit 1\n")
		if err := os.Chmod(filepath.Join(fixture.home, "hooks", name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeBackupTestFile(t, filepath.Join(fixture.home, ".gitconfig"), strings.Join([]string{
		"[user]", "\tuseConfigOnly = true",
		"[core]", "\thooksPath = " + filepath.Join(fixture.home, "hooks"),
		"[init]", "\tdefaultBranch = trunk",
		"",
	}, "\n"))

	restored, _ := decodeRestoreSuccess(t, fixture.restore(fixture.restoreWorkspace(id), "--json"))
	if restored.ProjectID != id || restored.Commit != fixture.targetMain() {
		t.Fatalf("restore = %+v, want %s at the target's main", restored, id)
	}
	if _, err := os.Lstat(marker); !os.IsNotExist(err) {
		t.Fatalf("a hook from the user's git config ran: %v", err)
	}
}

// fid may run inside a git hook or alias, where git exports variables naming
// the caller's repository. Restore's git commands ignore them.
func TestRestoreIgnoresInheritedGitRepositoryVariables(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		env  func(fixture *backupFixture, workspace string) []string
	}{
		{"GIT_DIR", func(_ *backupFixture, workspace string) []string {
			return []string{"GIT_DIR=" + filepath.Join(workspace, ".git")}
		}},
		{"GIT_OBJECT_DIRECTORY", func(fixture *backupFixture, _ string) []string {
			return []string{"GIT_OBJECT_DIRECTORY=" + filepath.Join(fixture.base, "foreign-objects")}
		}},
		{"GIT_CONFIG_PARAMETERS", func(fixture *backupFixture, _ string) []string {
			return []string{"GIT_CONFIG_PARAMETERS='url." + filepath.Join(fixture.base, "missing.git") + ".insteadof'='" + fixture.target + "'"}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newBackupFixture(t)
			original, id := fixture.project()
			fixture.setTarget(fixture.target)
			decodeBackupSuccess(t, fixture.backup(original, "--json"))
			workspace := fixture.restoreWorkspace(id)
			before := backupRepositoryState(fixture, workspace)

			restored, _ := decodeRestoreSuccess(t, restoreWith(t, append(append([]string{}, fixture.env...), test.env(fixture, workspace)...), workspace, "--json"))
			if restored.Commit != fixture.targetMain() {
				t.Fatalf("restore = %+v, want the target's main", restored)
			}
			if after := backupRepositoryState(fixture, workspace); !reflect.DeepEqual(after, before) {
				t.Fatalf("restore changed the project's repository:\n%s\nthen\n%s", before, after)
			}
			if _, err := os.Lstat(filepath.Join(fixture.base, "foreign-objects")); !os.IsNotExist(err) {
				t.Fatalf("restore used the inherited object directory: %v", err)
			}
		})
	}
}
