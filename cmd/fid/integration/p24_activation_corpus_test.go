//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

var p24ActivationRecoveryFields = map[string]string{"activation.recover": "reason", "activation.recover-complete": "selected_event_id"}
var p24ActivationAbortFields = map[string]string{"activation.abort": "reason"}

func TestP24ActivationRecoveryKeyedLeafCorpus(t *testing.T) {
	t.Parallel()
	body, err := os.ReadFile(filepath.Join(p22CorpusDir, "activation-cases.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cases map[string][]p22ConformanceCase
	if err := json.Unmarshal(body, &cases); err != nil {
		t.Fatal(err)
	}
	for _, status := range []string{"prepared", "committed"} {
		a := p22NewFIDActivationAdapter(t, status).(*p22FIDRecoveryAdapter)
		p24StartKeyedCapture(a.workspace)
		for _, fixture := range cases[status] {
			response := a.Execute(t, fixture)
			if response.ExitCode != fixture.ExitCode {
				t.Fatalf("%s setup = exit %d stderr=%q; want %d", fixture.Name, response.ExitCode, response.Stderr, fixture.ExitCode)
			}
		}
		captured := p24EndKeyedCapture(a.workspace)
		for operation, name := range p24ActivationRecoveryFields {
			entry, found := captured[operation]
			if !found {
				t.Fatalf("%s fixture did not capture %s", status, operation)
			}
			p24ExerciseCapturedLeaf(t, a.workspace, a.binary, operation, name, entry)
		}
		p24AssertRebuildMapParity(t, a.workspace, a.binary)
	}
}

func TestP24ActivationAbortKeyedLeafCorpus(t *testing.T) {
	t.Parallel()
	a := p22NewFIDActivationAdapter(t, "prepared").(*p22FIDRecoveryAdapter)
	p24StartKeyedCapture(a.workspace)
	request := map[string]any{"activation_id": a.activation["id"], "expected_status": "prepared", "expected_version": a.activation["version"], "reason": "operator-abort"}
	response := p23RunInput(t, a, "activation", "abort", a.actor, "p24-activation-abort", request)
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("activation.abort setup = %d %q", response.ExitCode, response.Stderr)
	}
	entry, found := p24EndKeyedCapture(a.workspace)["activation.abort"]
	if !found {
		t.Fatal("activation.abort was not captured")
	}
	p24ExerciseCapturedLeaf(t, a.workspace, a.binary, "activation.abort", p24ActivationAbortFields["activation.abort"], entry)
	p24AssertRebuildMapParity(t, a.workspace, a.binary)
}
