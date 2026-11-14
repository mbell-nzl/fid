package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
)

const configuredStatusPolicy = `
[lanes.routine]
checks = []
[lanes.standard]
checks = []
[lanes.assured]
checks = []
[review]
test_adequacy = []
`

func TestStatusGoldenPolicy(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, policy string
		kinds        []string
	}{
		{"no-origin", "", []string{"checks-not-configured", "test-paths-not-configured"}},
		{"no-file", "", []string{"checks-not-configured", "test-paths-not-configured"}},
		{"nothing-set", "", []string{"checks-not-configured", "test-paths-not-configured"}},
		{"some-lanes", "\n[lanes.standard]\nchecks = []\n[review]\ntest_adequacy = []\n", []string{"checks-not-configured"}},
		{"one-lane", "\n[lanes.routine]\nchecks = []\n[lanes.assured]\nchecks = []\n[review]\ntest_adequacy = []\n", []string{"checks-not-configured"}},
		{"all-set", strings.ReplaceAll(strings.ReplaceAll(configuredStatusPolicy, "checks = []", "checks = ['unit']"), "test_adequacy = []", "test_adequacy = ['**/*_test.go']") + "\n[checks.unit]\nrun = ['true']\ntimeout = '1m'\n", []string{}},
		{"empty-lists", configuredStatusPolicy, []string{}},
		{"test-paths", strings.Split(configuredStatusPolicy, "\n[review]")[0], []string{"test-paths-not-configured"}},
		{"unknown-default", "", []string{"policy-unavailable"}},
		{"deleted-default", "", []string{"policy-unavailable"}},
		{"invalid-file", "\n[lanes.standard]\npaths = []\n", []string{"policy-unavailable"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newStatusFixture(t, 100*time.Hour)
			fixture.setTarget()
			fixture.writeRecord(statusRecordLine{"pushed", 3 * time.Hour, ""})
			// The ordinary status fixtures have a configured policy. These tests
			// deliberately replace that trusted policy to exercise health warnings.
			identityGit(t, fixture.workspace, "remote", "remove", "origin")
			path := filepath.Join(fixture.workspace, ".fid.toml")
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			identityOnly := strings.Split(string(original), "\n[lanes.")[0]
			if err := os.WriteFile(path, []byte(identityOnly+test.policy), 0o600); err != nil {
				t.Fatal(err)
			}
			if test.name != "no-origin" {
				// Reuse the existing branch, with a local ref; fid itself never fetches.
				policyCommit(t, fixture.workspace)
				identityGit(t, fixture.workspace, "remote", "add", "origin", "/unused-origin")
				identityGit(t, fixture.workspace, "update-ref", "refs/remotes/origin/trunk", "HEAD")
				identityGit(t, fixture.workspace, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
			}
			if test.name == "no-file" {
				identityGit(t, fixture.workspace, "rm", "-q", ".fid.toml")
				identityGit(t, fixture.workspace, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-q", "-m", "remove policy")
				identityGit(t, fixture.workspace, "update-ref", "refs/remotes/origin/trunk", "HEAD")
			}
			if test.name == "unknown-default" {
				identityGit(t, fixture.workspace, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
			}
			if test.name == "deleted-default" {
				identityGit(t, fixture.workspace, "update-ref", "-d", "refs/remotes/origin/trunk")
			}
			if err := os.WriteFile(path, original, 0o600); err != nil {
				t.Fatal(err)
			} // worktree policy has no authority
			if test.name == "unknown-default" || test.name == "deleted-default" || test.name == "invalid-file" {
				fixture.seed(statusWork())
			}
			before := readAllPolicyJournal(t, fixture.workspace)
			for _, format := range []string{"json", "txt"} {
				args := []string{"status"}
				if format == "json" {
					args = []string{"--json", "status"}
				}
				result := fixture.run(args...)
				if result.err != nil || result.stderr != "" {
					t.Fatalf("status = %#v", result)
				}
				if format == "json" {
					var reply struct {
						Result struct {
							Warnings []map[string]any `json:"warnings"`
						} `json:"result"`
					}
					if err := json.Unmarshal([]byte(result.stdout), &reply); err != nil {
						t.Fatal(err)
					}
					kinds := []string{}
					for _, warning := range reply.Result.Warnings {
						kinds = append(kinds, warning["kind"].(string))
					}
					if reply.Result.Warnings == nil || !reflect.DeepEqual(kinds, test.kinds) {
						t.Fatalf("warnings = %#v, want %v", reply.Result.Warnings, test.kinds)
					}
					if test.name == "deleted-default" {
						want := map[string]any{"kind": "policy-unavailable", "problem": deletedPolicyDefaultProblem,
							"ref": "origin/trunk", "commit": nil, "summary": "Policy unavailable: " + deletedPolicyDefaultProblem,
							"next_action": "git remote set-head origin <branch>"}
						if !reflect.DeepEqual(reply.Result.Warnings[0], want) {
							t.Errorf("deleted default warning = %#v, want %#v", reply.Result.Warnings[0], want)
						}
					}
				}
				output := fixture.placeholders(result.stdout)
				if test.name != "no-origin" {
					output = strings.ReplaceAll(output, strings.TrimSpace(identityGit(t, fixture.workspace, "rev-parse", "HEAD")), "<COMMIT>")
				}
				assertGolden(t, filepath.Join("status", "policy-"+test.name+"."+format), output)
			}
			if after := readAllPolicyJournal(t, fixture.workspace); after != before {
				t.Fatal("status wrote an event")
			}
			if snapshots := policySnapshots(t, fixture.workspace); len(snapshots) != 0 {
				t.Fatalf("status wrote snapshots: %#v", snapshots)
			}
		})
	}
}

const deletedPolicyDefaultProblem = "origin's default branch origin/trunk does not exist"

func TestNextPolicyUnavailable(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"unknown-default", "invalid-file", "deleted-default"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, _ := policyWorkspace(t)
			policyOrigin(t, root)
			switch name {
			case "invalid-file":
				appendPolicy(t, root, "\n[lanes.standard]\npaths = []\n")
				policyCommit(t, root)
				identityGit(t, root, "update-ref", "refs/remotes/origin/trunk", "HEAD")
			case "unknown-default":
				identityGit(t, root, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
			case "deleted-default":
				identityGit(t, root, "update-ref", "-d", "refs/remotes/origin/trunk")
			}
			before := readAllPolicyJournal(t, root)
			for _, format := range []string{"json", "txt"} {
				args := []string{"next", "--claim", "--actor", "agent:builder"}
				if format == "json" {
					args = append([]string{"--json"}, args...)
				}
				result := executeCommandForTest(root, args...)
				if cli.ExitCode(result.err) != 4 || result.stdout != "" {
					t.Fatalf("next policy refusal = %#v", result)
				}
				if name == "deleted-default" && format == "json" {
					var reply struct {
						Error cli.Failure `json:"error"`
					}
					if err := json.Unmarshal([]byte(result.stderr), &reply); err != nil {
						t.Fatal(err)
					}
					want := map[string]any{"problem": deletedPolicyDefaultProblem, "ref": "origin/trunk", "commit": nil}
					if reply.Error.Code != cli.CodePolicyUnavailable || !reflect.DeepEqual(reply.Error.Details, want) {
						t.Errorf("deleted default failure = %#v, want details %#v", reply.Error, want)
					}
				}
				output := strings.ReplaceAll(result.stderr, strings.TrimSpace(identityGit(t, root, "rev-parse", "HEAD")), "<COMMIT>")
				assertGolden(t, filepath.Join("status", "next-policy-"+name+"."+format), output)
			}
			if readAllPolicyJournal(t, root) != before {
				t.Fatal("policy refusal wrote an event")
			}
		})
	}
}
