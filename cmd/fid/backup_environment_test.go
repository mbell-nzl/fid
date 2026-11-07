package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// The user's git config takes part in a backup only for transport and
// credentials. Hooks, signing, filters, templates and a missing identity in
// it must not stop or change a backup.
func TestBackupIgnoresHostileGlobalGitConfig(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, id := fixture.project()
	fixture.setTarget(fixture.target)
	// The global hooks path also reaches the local target's receive-pack,
	// which is the server's business, so the hooks fail only in fid's clone.
	marker := filepath.Join(fixture.base, "hook-ran")
	hook := "#!/bin/sh\ncase \"$(git rev-parse --absolute-git-dir)\" in\n'" + fixture.cacheHome + "'/*) touch '" + marker + "'; exit 1;;\nesac\n"
	for _, directory := range []string{filepath.Join(fixture.home, "hooks"), filepath.Join(fixture.home, "template", "hooks")} {
		for _, name := range []string{"pre-commit", "commit-msg", "post-commit", "pre-push", "reference-transaction", "post-checkout"} {
			writeBackupTestFile(t, filepath.Join(directory, name), hook)
			if err := os.Chmod(filepath.Join(directory, name), 0o755); err != nil {
				t.Fatal(err)
			}
		}
	}
	// A template can carry any file into a new repository, not just hooks.
	writeBackupTestFile(t, filepath.Join(fixture.home, "template", "fid-template-marker"), "from the user's template\n")
	writeBackupTestFile(t, filepath.Join(fixture.home, "attributes"), "* filter=broken\n")
	writeBackupTestFile(t, filepath.Join(fixture.home, ".gitconfig"), strings.Join([]string{
		"[user]", "\tuseConfigOnly = true",
		"[commit]", "\tgpgSign = true",
		"[gpg]", "\tprogram = false",
		"[core]", "\thooksPath = " + filepath.Join(fixture.home, "hooks"), "\tattributesFile = " + filepath.Join(fixture.home, "attributes"),
		"[filter \"broken\"]", "\tclean = false", "\tsmudge = false", "\trequired = true",
		"[init]", "\ttemplateDir = " + filepath.Join(fixture.home, "template"), "\tdefaultBranch = trunk",
		"",
	}, "\n"))

	record, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	if record.Outcome != "pushed" || *record.Commit != fixture.targetMain() {
		t.Fatalf("backup = %+v, want a push to main despite the global config", record)
	}
	if _, err := os.Lstat(marker); !os.IsNotExist(err) {
		t.Fatalf("a hook from the user's git config ran: %v", err)
	}
	if copied, err := filepath.Glob(filepath.Join(fixture.cacheHome, "fid", "backup", "*", "repo.git", "fid-template-marker")); err != nil || len(copied) != 0 {
		t.Fatalf("fid's clone was made from the user's init template: %v, %v", copied, err)
	}
	if got := fixture.targetGit("log", "-1", "--format=%an <%ae> %G?", "main"); got != "fid backup <backup@fid.invalid> N" {
		t.Fatalf("backup commit = %q, want fid's identity, unsigned", got)
	}
	archive, err := fixture.gitBytes("--git-dir="+fixture.target, "cat-file", "blob", "main:"+id+"/export.zip")
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(archive)
	if "sha256:"+hex.EncodeToString(digest[:]) != *record.ArchiveDigest {
		t.Fatal("stored archive differs from the export: a filter changed it")
	}
}

// fid may run inside a git hook or a git alias, where git exports variables
// naming the caller's repository, or where the user's shell sets a commit
// identity. Backup's git commands must ignore each of them.
func TestBackupIgnoresInheritedGitRepositoryVariables(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		env  func(fixture *backupFixture, workspace string) []string
	}{
		{"GIT_DIR", func(_ *backupFixture, workspace string) []string {
			return []string{"GIT_DIR=" + filepath.Join(workspace, ".git")}
		}},
		{"GIT_WORK_TREE", func(_ *backupFixture, workspace string) []string { return []string{"GIT_WORK_TREE=" + workspace} }},
		{"GIT_COMMON_DIR", func(_ *backupFixture, workspace string) []string {
			return []string{"GIT_COMMON_DIR=" + filepath.Join(workspace, ".git")}
		}},
		{"GIT_OBJECT_DIRECTORY", func(fixture *backupFixture, _ string) []string {
			return []string{"GIT_OBJECT_DIRECTORY=" + filepath.Join(fixture.base, "foreign-objects")}
		}},
		{"GIT_NAMESPACE", func(*backupFixture, string) []string { return []string{"GIT_NAMESPACE=intruder"} }},
		{"GIT_QUARANTINE_PATH", func(fixture *backupFixture, _ string) []string {
			return []string{"GIT_QUARANTINE_PATH=" + filepath.Join(fixture.base, "quarantine")}
		}},
		{"GIT_CONFIG_PARAMETERS", func(fixture *backupFixture, _ string) []string {
			return []string{"GIT_CONFIG_PARAMETERS='url." + filepath.Join(fixture.base, "missing.git") + ".insteadof'='" + fixture.target + "'"}
		}},
		{"GIT_CONFIG_COUNT", func(fixture *backupFixture, _ string) []string {
			return []string{"GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=url." + filepath.Join(fixture.base, "missing.git") + ".insteadOf", "GIT_CONFIG_VALUE_0=" + fixture.target}
		}},
		{"commit identity", func(*backupFixture, string) []string {
			return []string{"GIT_AUTHOR_NAME=intruder", "GIT_AUTHOR_EMAIL=intruder@example.test", "GIT_COMMITTER_NAME=intruder", "GIT_COMMITTER_EMAIL=intruder@example.test"}
		}},
		// A date left in the environment would date the backup's commit
		// wrongly, or stop the backup when it isn't a date at all.
		{"GIT_AUTHOR_DATE", func(*backupFixture, string) []string { return []string{"GIT_AUTHOR_DATE=@978307200 +0000"} }},
		{"GIT_COMMITTER_DATE", func(*backupFixture, string) []string { return []string{"GIT_COMMITTER_DATE=@978307200 +0000"} }},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newBackupFixture(t)
			workspace, _ := fixture.project()
			fixture.setTarget(fixture.target)
			before := backupRepositoryState(fixture, workspace)
			started := time.Now().Unix()

			result := runBackupWith(append(append([]string{}, fixture.env...), test.env(fixture, workspace)...), workspace, time.Minute, "--json")
			record, _ := decodeBackupSuccess(t, result)
			main := fixture.targetMain()
			if record.Outcome != "pushed" || main == "" || *record.Commit != main {
				t.Fatalf("backup = %+v, want a push to the target's main %q", record, main)
			}
			if refs := fixture.targetGit("for-each-ref", "--format=%(refname)"); refs != "refs/heads/main" {
				t.Fatalf("target refs = %q, want only refs/heads/main", refs)
			}
			if tracked := fixture.git("--git-dir="+fixture.clone(), "rev-parse", "refs/remotes/origin/main"); tracked != main {
				t.Fatalf("fid's clone tracks %s, want %s", tracked, main)
			}
			if got := fixture.targetGit("log", "-1", "--format=%an <%ae>|%cn <%ce>", "main"); got != "fid backup <backup@fid.invalid>|fid backup <backup@fid.invalid>" {
				t.Fatalf("backup commit identity = %q", got)
			}
			for _, format := range []string{"%at", "%ct"} {
				if dated, err := strconv.ParseInt(fixture.targetGit("log", "-1", "--format="+format, "main"), 10, 64); err != nil || dated < started {
					t.Fatalf("backup commit %s = %d, want this run's time, from %d: %v", format, dated, started, err)
				}
			}
			if after := backupRepositoryState(fixture, workspace); after != before {
				t.Fatalf("backup changed the project's repository:\n%s\nthen\n%s", before, after)
			}
			for _, path := range []string{filepath.Join(fixture.base, "foreign-objects"), filepath.Join(fixture.base, "quarantine")} {
				if _, err := os.Lstat(path); !os.IsNotExist(err) {
					t.Fatalf("backup used the inherited %s: %v", path, err)
				}
			}
		})
	}
}

// backupRepositoryState describes the project's own repository: its refs,
// object counts and config.
func backupRepositoryState(fixture *backupFixture, workspace string) string {
	fixture.t.Helper()
	config, err := os.ReadFile(filepath.Join(workspace, ".git", "config"))
	if err != nil {
		fixture.t.Fatal(err)
	}
	return fixture.git("-C", workspace, "for-each-ref") + "\n" + fixture.git("-C", workspace, "count-objects", "-v") + "\n" + string(config)
}

// Without the test seam, fid backup reads the real process environment:
// user config from XDG_CONFIG_HOME or ~/.config, and its clone under
// XDG_CACHE_HOME or ~/.cache.
func TestBackupFindsConfigAndCacheThroughXDGDirectories(t *testing.T) { //nolint:paralleltest // t.Setenv replaces HOME and the XDG directories.
	base := t.TempDir()
	home := filepath.Join(base, "home")
	target := filepath.Join(base, "backups.git")
	fixture := &backupFixture{t: t, base: base, home: home, env: backupTestEnvironment(home, "", "")}
	fixture.git("init", "--bare", "-q", "--template=", target)
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("XDG_CACHE_HOME", "")
	if err := config.SetBackupTarget(filepath.Join(home, ".config", "fid", "config.toml"), target); err != nil {
		t.Fatal(err)
	}
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init", "--json")
	sum := sha256.Sum256([]byte(target))
	key := hex.EncodeToString(sum[:])[:32]

	if result := executeCommandForTest(workspace, "--json", "backup"); result.err != nil {
		t.Fatalf("backup = %#v", result)
	}
	clone := filepath.Join(home, ".cache", "fid", "backup", key, "repo.git")
	if got := fixture.git("--git-dir="+clone, "rev-parse", "--is-bare-repository"); got != "true" {
		t.Fatalf("fid's clone at %s is not a bare repository: %q", clone, got)
	}

	cacheHome := filepath.Join(base, "cache")
	t.Setenv("XDG_CACHE_HOME", cacheHome)
	if result := executeCommandForTest(workspace, "--json", "backup"); result.err != nil || !strings.Contains(result.stdout, `"outcome":"unchanged"`) {
		t.Fatalf("backup = %#v, want unchanged through a new clone", result)
	}
	if _, err := os.Stat(filepath.Join(cacheHome, "fid", "backup", key, "repo.git", "HEAD")); err != nil {
		t.Fatalf("XDG_CACHE_HOME clone: %v", err)
	}
}

// Pushes use the user's own credential setup, so the variables that carry it
// reach git and the transports it starts. But git must never prompt for a
// password, even where the user's environment allows prompts: the loop and
// cron would hang. A remote helper on PATH records the environment git gives
// it, then fails like an unreachable target.
func TestBackupPassesTheUsersCredentialSetupButNeverPrompts(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()
	bin, seen := filepath.Join(fixture.base, "bin"), filepath.Join(fixture.base, "helper-environment")
	writeBackupTestFile(t, filepath.Join(bin, "git-remote-fidprobe"), "#!/bin/sh\nenv > '"+seen+"'\nexit 1\n")
	if err := os.Chmod(filepath.Join(bin, "git-remote-fidprobe"), 0o755); err != nil {
		t.Fatal(err)
	}
	fixture.setTarget("fidprobe::ledgers")
	credentials := []string{"GIT_ASKPASS=/fid-test/askpass", "SSH_AUTH_SOCK=/fid-test/agent.sock", "GIT_SSH_COMMAND=ssh -o fid-test=yes"}
	environ := append(append(append([]string{}, fixture.env...), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "GIT_TERMINAL_PROMPT=1"), credentials...)

	failure := decodeBackupFailure(t, runBackupWith(environ, workspace, time.Minute, "--json"), 4)
	if failure.Code != "repository-unavailable" || failure.Details["reason"] != "not-a-repository" {
		t.Fatalf("failure = %+v, want the helper's failure as an unreachable target", failure)
	}
	payload, err := os.ReadFile(seen)
	if err != nil {
		t.Fatalf("git never ran the remote helper: %v", err)
	}
	lines := strings.Split(string(payload), "\n")
	for _, want := range append(credentials, "GIT_TERMINAL_PROMPT=0") {
		if !slices.Contains(lines, want) {
			t.Errorf("git's environment lacks %q:\n%s", want, payload)
		}
	}
	if slices.Contains(lines, "GIT_TERMINAL_PROMPT=1") {
		t.Error("git's environment still allows terminal prompts")
	}
}
