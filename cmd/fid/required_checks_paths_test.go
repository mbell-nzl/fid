package main

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/policy"
)

func TestRequiredChecksConservativelyMatchUnsupportedGitPaths(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, path string
		all        bool
	}{
		{"newline", "new\nfile.go", true},
		{"backslash", `new\file.go`, true},
		{"invalid UTF-8", "new\xfffile.go", true},
		{"Unicode control", "new\u0085file.go", true},
		{"space and Unicode", "notes space 日本語.md", false},
		{"literal glob operators", "notes[ab]?.md", false},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "['go', 'always', 'empty']")
			path := filepath.Join(f.workspace, ".fid.toml")
			raw := strings.Replace(string(mustReadCommandFile(t, path)), "checks = ['go', 'always', 'empty']", "checks = ['go', 'docs', 'always', 'empty']", 1)
			raw += "\n[checks.docs]\nrun=['true']\ntimeout='1m'\nwhen_changed=['docs/**']\n"
			writeBackupTestFile(t, path, raw)
			runGitForTest(t, f.workspace, "add", ".fid.toml")
			runGitForTest(t, f.workspace, "commit", "-m", "Independent conditional requirements")
			identityGit(t, f.workspace, "push", "-q", "origin", "HEAD:trunk")
			writeBackupTestFile(t, filepath.Join(f.workspace, row.path), "committed task file\n")
			runGitForTest(t, f.workspace, "add", ".")
			runGitForTest(t, f.workspace, "commit", "-m", "Task path")
			comparison, err := policy.Compare(f.workspace, "FID-9f13aa")
			if err != nil || !slices.Equal(comparison.Paths, []string{row.path}) {
				t.Fatalf("committed path bytes=%q/%v, want %q", comparison.Paths, err, row.path)
			}
			want := []string{"always", "empty"}
			if row.all {
				want = []string{"always", "docs", "empty", "go"}
			}
			shown := readShownChecks(t, f)
			names := []string{}
			for _, check := range shown.Checks {
				checkBool(t, check.Name+" required", check.Required, true)
				checkBool(t, check.Name+" satisfied before checking", check.Satisfied, false)
				names = append(names, check.Name)
			}
			if !slices.Equal(names, want) || shown.Base == nil || len(shown.Warnings) != 0 {
				t.Errorf("show=%+v, want required=%v with a known base and no warnings", shown, want)
			}
			if got := requiredNames(t, checkRequired(t, f, "path")); !slices.Equal(got, want) {
				t.Errorf("checks=%v, want %v", got, want)
			}
		})
	}
}
