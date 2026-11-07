package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
// naming the caller's repository. Backup's git commands must ignore them.
func TestBackupIgnoresInheritedGitRepositoryVariables(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()
	fixture.setTarget(fixture.target)
	objects := filepath.Join(fixture.base, "foreign-objects")
	before := fixture.git("-C", workspace, "count-objects")
	fixture.env = append(fixture.env,
		"GIT_DIR="+filepath.Join(workspace, ".git"), "GIT_WORK_TREE="+workspace,
		"GIT_INDEX_FILE="+filepath.Join(fixture.base, "foreign-index"), "GIT_OBJECT_DIRECTORY="+objects,
		"GIT_CONFIG_PARAMETERS='core.hookspath'='"+filepath.Join(fixture.base, "nowhere")+"'")

	record, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	fixture.env = backupTestEnvironment(fixture.home, fixture.configHome, fixture.cacheHome)
	if record.Outcome != "pushed" || *record.Commit != fixture.targetMain() {
		t.Fatalf("backup = %+v, want a push to the target", record)
	}
	if after := fixture.git("-C", workspace, "count-objects"); after != before {
		t.Fatalf("backup wrote into the project's repository: %q, then %q", before, after)
	}
	for _, path := range []string{objects, filepath.Join(fixture.base, "foreign-index")} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("backup used the inherited %s: %v", path, err)
		}
	}
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
