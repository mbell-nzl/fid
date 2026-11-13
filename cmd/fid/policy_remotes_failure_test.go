package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
)

const policyRemotesProblem = "couldn't list git remotes: git remote: exit status 1: remotes unavailable for test"

func TestPolicyRemotesUnavailable(t *testing.T) { //nolint:paralleltest // Overrides process PATH to fail only git remote.
	fixture := newStatusFixture(t, 100*time.Hour)
	fixture.setTarget()
	fixture.writeRecord(statusRecordLine{"pushed", 3 * time.Hour, ""})
	fixture.seed(statusWork())
	failPolicyRemotes(t)
	before := readAllPolicyJournal(t, fixture.workspace)
	for _, operation := range []string{"status", "next"} {
		for _, format := range []string{"json", "txt"} {
			t.Run(operation+"/"+format, func(t *testing.T) {
				args := []string{operation}
				if operation == "next" {
					args = append(args, "--claim", "--actor", "agent:builder")
				}
				if format == "json" {
					args = append([]string{"--json"}, args...)
				}
				result := fixture.run(args...)
				output, golden := result.stdout, "policy-remotes-unavailable."
				if operation == "next" {
					if cli.ExitCode(result.err) != 4 || result.stdout != "" {
						t.Fatalf("next refusal = %#v, want exit 4 and empty stdout", result)
					}
					output, golden = result.stderr, "next-"+golden
				} else if result.err != nil || result.stderr != "" {
					t.Fatalf("status = %#v", result)
				}
				if format == "json" {
					var reply struct {
						Error  cli.Failure `json:"error"`
						Result struct {
							Warnings []map[string]any `json:"warnings"`
						} `json:"result"`
					}
					if err := json.Unmarshal([]byte(output), &reply); err != nil {
						t.Fatal(err)
					}
					want := map[string]any{"problem": policyRemotesProblem, "ref": nil, "commit": nil}
					if operation == "next" {
						if reply.Error.Code != cli.CodePolicyUnavailable || !reflect.DeepEqual(reply.Error.Details, want) {
							t.Errorf("failure = %#v, want details %#v", reply.Error, want)
						}
					} else {
						want["kind"], want["summary"], want["next_action"] = "policy-unavailable", "Policy unavailable: "+policyRemotesProblem, "check the repository: git remote -v"
						if !reflect.DeepEqual(reply.Result.Warnings, []map[string]any{want}) {
							t.Errorf("warnings = %#v, want %#v", reply.Result.Warnings, want)
						}
					}
				}
				assertGolden(t, filepath.Join("status", golden+format), fixture.placeholders(output))
			})
		}
	}
	if readAllPolicyJournal(t, fixture.workspace) != before {
		t.Fatal("policy read failure wrote an event")
	}
}

func failPolicyRemotes(t *testing.T) {
	t.Helper()
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	script := "#!/bin/sh\n" +
		"if [ \"$#\" -eq 3 ] && [ \"$1\" = '-C' ] && [ \"$3\" = 'remote' ]; then\n" +
		"  echo 'remotes unavailable for test' >&2\n  exit 1\nfi\n" +
		"exec \"$FID_TEST_REAL_GIT\" \"$@\"\n"
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FID_TEST_REAL_GIT", git)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
}
