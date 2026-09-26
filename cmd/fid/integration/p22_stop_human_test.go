//go:build integration

package integration

import (
	"strings"
	"testing"
)

func TestP22StopAttentionHumanListsNameReasonAndResolution(t *testing.T) {
	t.Parallel()
	adapter := p22NewFIDRecoveryAdapter(t)
	fixture := p22ConformanceCase{Name: "human-stop", Operation: "stop.raise", Kind: "accepted"}
	raised := adapter.Execute(t, fixture)
	if raised.ExitCode != 0 {
		t.Fatalf("raise: %s", raised.Stderr)
	}
	for _, group := range []string{"stop", "attention"} {
		stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, "--project", adapter.workspace, group, "list")
		if err != nil {
			t.Fatalf("%s list: %v %s", group, err, stderr)
		}
		for _, want := range []string{adapter.stopID, "conformance operator intervention", "work-ended", "run:" + adapter.runID} {
			if !strings.Contains(stdout, want) {
				t.Fatalf("%s Human list %q missing %q", group, stdout, want)
			}
		}
		if group == "attention" && !strings.Contains(stdout, adapter.attentionID) {
			t.Fatalf("Attention ID absent: %q", stdout)
		}
	}
}
