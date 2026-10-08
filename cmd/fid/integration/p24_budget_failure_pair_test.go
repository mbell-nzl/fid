//go:build integration

package integration

import (
	"reflect"
	"strings"
	"testing"
)

// Phase 22 committed admission denials are singleton terminal fid.live.v1
// frames (no Attempt or process). Keep them separate from stderr envelopes
// and from the gate-unsuccessful row with streamed process output.
func p24BudgetTerminalPair(t *testing.T) {
	t.Helper()
	code := t.Name()[strings.LastIndex(t.Name(), "/")+1:]
	pair := p24GateCodePair(t, code)
	binary := integrationFIDBinary(t)
	args := append([]string{"--project", pair.root}, pair.args...)
	var payload map[string]any
	for attempt := 0; attempt < 2; attempt++ {
		out, stderr, err := closureRunCommand(pair.root, binary, append([]string{"--json"}, args...)...)
		if err == nil || processExitCode(t, err) != 4 || stderr != "" {
			t.Fatalf("JSON budget terminal: %v %q %q", err, out, stderr)
		}
		frame := p24JSONDecode(t, []byte(out))
		data := frame["data"].(map[string]any)
		payload = data["response"].(map[string]any)
		if frame["frame_schema"] != "fid.live.v1" || frame["sequence"] != float64(1) || data["outcome"] != "fid-error" || data["attempt_id"] != nil || data["receipt_id"] != nil || data["replayed"] != (attempt == 1) || payload["operation"] != "gate.run" || payload["error"].(map[string]any)["code"] != code {
			t.Fatalf("committed budget terminal: %#v", frame)
		}
	}
	before := p24FailureFilesystem(t, pair.root)
	out, stderr, err := closureRunCommand(pair.root, binary, args...)
	want := p24ExpectedFailureText(t, payload, "gate.run", true)
	if err == nil || processExitCode(t, err) != 4 || out != "" || stderr != want {
		t.Fatalf("Human budget terminal: %v stdout=%q stderr=%q want=%q", err, out, stderr, want)
	}
	if !reflect.DeepEqual(before, p24FailureFilesystem(t, pair.root)) {
		t.Fatal("budget denial replay changed filesystem")
	}
}
