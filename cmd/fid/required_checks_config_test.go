package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestRequiredChecksIncludeGitlinksDespiteIgnoreConfig(t *testing.T) {
	t.Parallel()
	for _, setting := range []string{"diff.ignoreSubmodules", "submodule.dependency.ignore"} {
		t.Run(setting, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "['go']")
			dependency := filepath.Join(f.workspace, "dependency.go")
			if err := gittest.InitRepository(dependency); err != nil {
				t.Fatal(err)
			}
			runGitForTest(t, dependency, "commit", "--allow-empty", "-m", "Dependency base")
			writeBackupTestFile(t, filepath.Join(f.workspace, ".gitmodules"), "[submodule \"dependency\"]\npath = dependency.go\nurl = ./dependency\n")
			runGitForTest(t, f.workspace, "add", ".gitmodules", "dependency.go")
			runGitForTest(t, f.workspace, "commit", "-m", "Dependency base gitlink")
			identityGit(t, f.workspace, "push", "-q", "origin", "HEAD:trunk")
			runGitForTest(t, dependency, "commit", "--allow-empty", "-m", "Dependency change")
			runGitForTest(t, f.workspace, "add", "dependency.go")
			runGitForTest(t, f.workspace, "commit", "-m", "Task gitlink change")
			for _, mode := range []string{"before", "after"} {
				if mode == "after" {
					runGitForTest(t, f.workspace, "config", setting, "all")
				}
				if got := requiredNames(t, checkRequired(t, f, mode)); !slices.Equal(got, []string{"go"}) {
					t.Fatalf("%s %s: checks=%v, want [go]", setting, mode, got)
				}
			}
		})
	}
}

func TestRequiredChecksComparisonIgnoresLocalDiffConfig(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, key, value string
		subdirectory     bool
	}{
		{"relative", "diff.relative", "true", true},
		{"order file", "diff.orderFile", "missing-order-file", false},
		{"renames", "diff.renames", "copies", false},
		{"external", "diff.external", "false", false},
		{"external driver", "diff.hostile.command", "false", false},
		{"textconv", "diff.hostile.textconv", "false", false},
		{"quoted paths", "core.quotePath", "true", false},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "['go']")
			identityGit(t, f.workspace, "mv", "old.go", "renamed space.go")
			runGitForTest(t, f.workspace, "commit", "-m", "Task rename")
			writeBackupTestFile(t, filepath.Join(f.workspace, ".git", "info", "attributes"), "*.go diff=hostile\n")
			runGitForTest(t, f.workspace, "config", row.key, row.value)
			root := f.workspace
			if row.subdirectory {
				root = filepath.Join(root, "nested")
				if err := os.Mkdir(root, 0o700); err != nil {
					t.Fatal(err)
				}
			}
			comparison, err := policy.Compare(root, "FID-9f13aa")
			if err != nil {
				t.Fatal(err)
			}
			if want := []string{"old.go", "renamed space.go"}; !slices.Equal(comparison.Paths, want) {
				t.Fatalf("paths=%q, want %q", comparison.Paths, want)
			}
			if head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD")); comparison.Head != head {
				t.Fatalf("head=%s, want %s", comparison.Head, head)
			}
		})
	}
}
