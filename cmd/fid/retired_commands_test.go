package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
)

func TestRetiredCommandsAreUnknown(t *testing.T) {
	t.Parallel()
	for _, name := range strings.Fields("activation attention capabilities claim concern decision deviation envelope evidence gate human-outcome ready receipt requirement resume run stats stop workflow") {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, jsonMode := range []bool{false, true} {
				args := []string{name}
				if jsonMode {
					args = append([]string{"--json"}, args...)
				}
				result := executeCommandForTest(t.TempDir(), args...)
				if cli.ExitCode(result.err) != 2 || result.stdout != "" {
					t.Fatalf("retired %s = %#v", name, result)
				}
				message := fmt.Sprintf("unknown command %q for %q", name, "fid")
				if !jsonMode {
					if !strings.Contains(result.stderr, "error: "+message+"\n") {
						t.Fatalf("%s: %s", name, result.stderr)
					}
					continue
				}
				var failure struct {
					Operation string `json:"operation"`
					Error     struct {
						Code, Message string
						Details       struct {
							Argument, Value string
							Allowed         []string
						}
					}
				}
				if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil {
					t.Fatal(err)
				}
				if failure.Operation != "fid" || failure.Error.Code != "operation-unavailable" || failure.Error.Message != message || failure.Error.Details.Argument != "command" || failure.Error.Details.Value != name {
					t.Fatalf("%s: %s", name, result.stderr)
				}
				for _, allowed := range failure.Error.Details.Allowed {
					if allowed == name {
						t.Fatalf("retired alternative %s", name)
					}
				}
			}
		})
	}
}

func TestArtifactVerifyTypedActors(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ actor, code string }{{"untyped", "invalid-argument"}, {"service:robot", "authority-denied"}, {"agent:builder", "invalid-argument"}, {"human:mark", "invalid-argument"}} {
		t.Run(row.actor, func(t *testing.T) {
			t.Parallel()
			workspace := clonePristineWorkspace(t)
			input := t.TempDir() + "/verify.json"
			writeActivationInput(t, input, map[string]any{"artifact_id": "ART-missing"})
			before := startEventCount(t, workspace)
			result := executeCommandForTest(workspace, "--json", "artifact", "verify", "--input", input, "--actor", row.actor, "--idempotency-key", "verify")
			wantExit := 2
			if row.code == "authority-denied" {
				wantExit = 5
			}
			if cli.ExitCode(result.err) != wantExit || result.stdout != "" || !strings.Contains(result.stderr, `"code":"`+row.code+`"`) || !strings.Contains(result.stderr, `"operation":"artifact.verify"`) {
				t.Fatalf("%s: %#v", row.actor, result)
			}
			if strings.HasPrefix(row.actor, "agent:") || strings.HasPrefix(row.actor, "human:") {
				if !strings.Contains(result.stderr, "was not found") || !strings.Contains(result.stderr, `"argument":"input"`) || !strings.Contains(result.stderr, `"allowed":["accepted gate registration input"]`) {
					t.Fatalf("valid actor refused before lookup: %s", result.stderr)
				}
			}
			if got := startEventCount(t, workspace); got != before {
				t.Fatalf("refusal wrote %d Events", got-before)
			}
		})
	}
}

// The review root returns with a new lean schema in 1.6d. Its retired kernel
// inputs remain refused, alongside the unchanged retired Event/table tests.
func TestRetiredReviewInputRemainsRefused(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"envelope_id", "claim_id", "run_id"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			f := reviewFixture(t)
			in := reviewInput()
			in[field] = "retired-record"
			before := reviewProjection(t, f)
			got := recordReview(t, f, in, "old-input", "human:courier")
			if got.err == nil || !strings.Contains(got.stderr, "input."+field) || reviewProjection(t, f) != before {
				t.Fatalf("retired input=%#v", got)
			}
		})
	}
}

// RESET 2.1c restores defer; the other retired operations remain unavailable.
func TestRetiredFindingOperationsRemainRefused(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"address", "resolve", "record", "defer", "list", "show"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, jsonMode := range []bool{false, true} {
				args := []string{"finding", name}
				if jsonMode {
					args = append([]string{"--json"}, args...)
				}
				got := executeCommandForTest(t.TempDir(), args...)
				if cli.ExitCode(got.err) != 2 || got.stdout != "" {
					t.Fatalf("retired finding %s=%#v", name, got)
				}
				if name == "defer" {
					if jsonMode {
						assertJSONFailureCode(t, got.stderr, "finding.defer", "invalid-argument")
					} else if !strings.Contains(got.stderr, "accepts 1 arg(s), received 0") || !strings.Contains(got.stderr, "allowed=<finding>") {
						t.Fatal(got.stderr)
					}
					continue
				}
				if jsonMode {
					assertJSONFailureCode(t, got.stderr, "finding", "operation-unavailable")
				} else if !strings.Contains(got.stderr, fmt.Sprintf("unknown command %q for %q", name, "fid finding")) {
					t.Fatal(got.stderr)
				}
			}
		})
	}
}
