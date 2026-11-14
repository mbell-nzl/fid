//go:build integration

package integration

import (
	"reflect"
	"strings"
	"testing"
)

// A real retained file's contents never enter a generic verification failure.
func TestP24GenericFallbackSensitiveProcess(t *testing.T) {
	t.Parallel()
	binary, workspace := integrationFIDBinary(t), t.TempDir()
	const sentinel = "P24-GENERIC-SECRET-SENTINEL"
	leanArtifactFixture(t, binary, workspace, sentinel)
	input := p21Input(t, workspace, "sensitive-failure", map[string]any{"artifact_id": "ART-missing"})
	args := []string{"artifact", "verify", "--input", input, "--actor", "agent:builder", "--idempotency-key", "sensitive-failure"}
	before := p24FailureFilesystem(t, workspace)
	jsonOut, jsonErr, jsonExit := closureRunCommand(workspace, binary, append([]string{"--json"}, args...)...)
	if jsonExit == nil || processExitCode(t, jsonExit) != 2 || jsonOut != "" {
		t.Fatalf("verification failure: %v %q %q", jsonExit, jsonOut, jsonErr)
	}
	f := p23Envelope(t, jsonErr)["error"].(map[string]any)
	details := map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted gate registration input"}}
	message := `verify Artifact: Artifact "ART-missing" was not found`
	if f["code"] != "invalid-argument" || f["category"] != "invalid-input" || f["retryable"] != false || f["subject"] != nil || f["message"] != message || !reflect.DeepEqual(f["details"], details) {
		t.Fatalf("generic fallback=%#v", f)
	}
	humanOut, humanErr, humanExit := closureRunCommand(workspace, binary, args...)
	want := "error: " + message + "\nreason: invalid-argument; argument=input, allowed=accepted gate registration input\nnext: see fid artifact verify --help; template: fid artifact verify --example\n"
	if humanExit == nil || processExitCode(t, humanExit) != 2 || humanOut != "" || humanErr != want {
		t.Fatalf("human failure: %v %q %q want %q", humanExit, humanOut, humanErr, want)
	}
	if strings.Contains(jsonOut+jsonErr+humanOut+humanErr, sentinel) {
		t.Fatal("retained content leaked into failure")
	}
	if !reflect.DeepEqual(before, p24FailureFilesystem(t, workspace)) {
		t.Fatal("refusal changed filesystem")
	}
}
