//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

var p24PlainCopyMessages = map[string]string{
	"envelope.create": "prior registration binding does not match envelope",
	"claim.acquire":   "claim acquire binding is invalid",
	"run.start":       "run claim binding is invalid",
	"gate.run":        "run gate: repository binding is invalid",
}

func TestP24PlainCopyFailurePairs(t *testing.T) {
	t.Parallel()
	for operation, message := range p24PlainCopyMessages {
		t.Run(operation, func(t *testing.T) {
			t.Parallel()
			binary, original := integrationFIDBinary(t), t.TempDir()
			p24StartKeyedCapture(original)
			a := p24RunAdapterAt(t, binary, original)
			captured := p24EndKeyedCapture(original)
			var input map[string]any
			if operation == "gate.run" {
				input = p24ObservedGateInput(t, a)
			} else {
				if err := json.Unmarshal(captured[operation].input, &input); err != nil {
					t.Fatal(err)
				}
			}
			// The copied Run retains an otherwise valid, active Claim and its exact
			// identity/version. Only the observed root differs; no snapshot mismatch.
			copied := t.TempDir()
			copyParityFixture(t, original, copied)
			// Copying changes file stat data. Refresh the fixture index before the
			// refusal snapshot, so normal Git observation is read-only thereafter.
			p22Git(t, copied, "update-index", "--refresh")
			path := p21Input(t, copied, "plain-copy", input)
			parts := strings.Split(operation, ".")
			args := []string{"--project", copied, parts[0], parts[1], "--actor", a.actor, "--input", path, "--idempotency-key", "plain-copy"}
			before := p24FailureFilesystem(t, copied)
			out, stderr, err := closureRunCommand(copied, binary, append([]string{"--json"}, args...)...)
			if err == nil || processExitCode(t, err) != 2 || out != "" {
				t.Fatalf("copy refusal: %v %q %q", err, out, stderr)
			}
			f := p23Envelope(t, stderr)["error"].(map[string]any)
			allowed := "accepted " + operation + " input"
			if operation == "gate.run" {
				allowed = "accepted gate registration input"
			}
			details := map[string]any{"argument": "input", "value": nil, "allowed": []any{allowed}}
			if f["code"] != "invalid-argument" || f["message"] != message || f["category"] != "invalid-input" || f["retryable"] != false || f["subject"] != nil || !reflect.DeepEqual(f["details"], details) {
				t.Fatalf("copy mapping: %#v", f)
			}
			if !reflect.DeepEqual(before, p24FailureFilesystem(t, copied)) {
				t.Fatal("JSON refusal changed filesystem")
			}
			out, stderr, err = closureRunCommand(copied, binary, args...)
			want := "error: " + message + "\nreason: invalid-argument; argument=input, allowed=" + allowed + "\nnext: see fid " + strings.Join(parts, " ") + " --help; template: fid " + strings.Join(parts, " ") + " --example\n"
			if err == nil || processExitCode(t, err) != 2 || out != "" || stderr != want {
				t.Fatalf("Human copy refusal: %v %q %q want %q", err, out, stderr, want)
			}
			if !reflect.DeepEqual(before, p24FailureFilesystem(t, copied)) {
				t.Fatal("Human refusal changed filesystem")
			}
			// Baseline replay succeeds in the original worktree. The existing
			// original/copy/import acceptance suites also exercise fresh import keys.
			if operation != "gate.run" {
				entry := captured[operation]
				out, stderr, err = closureRunCommand(original, binary, entry.args...)
				if err != nil || stderr != "" || out == "" {
					t.Fatalf("original control: %v %q %q", err, out, stderr)
				}
			}
		})
	}
}

func p24RunAdapterAt(t *testing.T, binary, root string) *p22FIDRecoveryAdapter {
	t.Helper()
	return p22NewFIDRecoveryAdapterAt(t, binary, root, "agent", 3600, []string{}, p22FinalActions)
}
