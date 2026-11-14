package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestRetiredCommandsAreUnknown(t *testing.T) {
	t.Parallel()
	for _, name := range strings.Fields("activation attention capabilities claim concern decision deviation envelope evidence finding gate human-outcome ready receipt requirement resume review run stats stop workflow") {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, jsonMode := range []bool{false, true} {
				args := []string{name}
				if jsonMode {
					args = append([]string{"--json"}, args...)
				}
				result := executeCommandForTest(t.TempDir(), args...)
				if result.err == nil || result.stdout != "" {
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
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"`+row.code+`"`) {
				t.Fatalf("%s: %#v", row.actor, result)
			}
			if strings.HasPrefix(row.actor, "agent:") || strings.HasPrefix(row.actor, "human:") {
				if !strings.Contains(result.stderr, "was not found") {
					t.Fatalf("valid actor refused before lookup: %s", result.stderr)
				}
			}
			if got := startEventCount(t, workspace); got != before {
				t.Fatalf("refusal wrote %d Events", got-before)
			}
		})
	}
}
