//go:build integration

package integration

import (
	"strings"
	"testing"
)

func TestP23ConcernListHumanProjection(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "concern.record", "envelope.activate", "envelope.create", "gate.register", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions)
	response := p23ConcernMutation(t, a, "record", "human-list", map[string]any{
		"id": "CON-human-list", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"scope": map[string]any{"paths": []string{}}, "severity": "major",
		"summary": "private concern prose", "reasoning": "private reasoning", "references": []any{},
	})
	p23AssertReviewSuccess(t, "record generic Concern", response)
	jsonResult := a.p22Run(t, "concern", "list", "--subject", "issue:"+a.issueID)
	if jsonResult.ExitCode != 0 || jsonResult.Stderr != "" {
		t.Fatalf("generic Concern list = %#v", jsonResult)
	}
	items := p21Path(t, p23Envelope(t, jsonResult.Stdout), "result", "items").([]any)
	if len(items) != 1 || items[0].(map[string]any)["id"] != "CON-human-list" {
		t.Fatalf("generic Concern list = %#v", items)
	}
	human, stderr, err := closureRunCommand(a.workspace, a.binary, "--project", a.workspace, "concern", "list", "--subject", "issue:"+a.issueID)
	if err != nil || stderr != "" || !strings.Contains(human, "Purpose: Inspect a raised concern") || !strings.Contains(human, "CON-human-list") || strings.Contains(human, "private concern prose") {
		t.Fatalf("generic Concern Human list = %q / %q / %v", human, stderr, err)
	}
}
