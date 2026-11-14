package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
)

func TestRequiredCheckGolden(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"no-names", "required", "conditional-empty", "missing-checks", "no-origin"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			names := "[]"
			switch name {
			case "required":
				names = "['empty', 'always']"
			case "conditional-empty":
				names = "['go']"
			}
			f := requiredCheckFixture(t, names)
			switch name {
			case "missing-checks":
				path := filepath.Join(f.workspace, ".fid.toml")
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				writeBackupTestFile(t, path, strings.Replace(string(raw), "[lanes.standard]\nchecks = []", "[lanes.standard]", 1))
				policyCommit(t, f.workspace)
				identityGit(t, f.workspace, "push", "-q", "origin", "HEAD:trunk")
			case "no-origin":
				identityGit(t, f.workspace, "remote", "remove", "origin")
			}
			for _, jsonMode := range []bool{false, true} {
				args := []string{"check", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "required-golden"}
				ext := ".txt"
				if jsonMode {
					args = append([]string{"--json"}, args...)
					ext = ".json"
				}
				got := f.run(args...)
				if name == "missing-checks" || name == "no-origin" {
					if cli.ExitCode(got.err) != 4 {
						t.Fatalf("refusal=%#v", got)
					}
				} else if got.err != nil {
					t.Fatalf("success=%#v", got)
				}
				output := got.stdout + got.stderr
				checkGolden(t, f, name+ext, output)
			}
		})
	}
}
