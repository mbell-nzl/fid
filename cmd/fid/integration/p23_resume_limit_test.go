//go:build integration

package integration

import (
	"fmt"
	"strings"
	"testing"
)

func TestP23ResumeRejectsJudgmentSectionOverflow(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "concern.record", "envelope.activate", "envelope.create", "gate.register", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions)
	for index := 0; index < 101; index++ {
		response := p23ConcernMutation(t, a, "record", fmt.Sprintf("limit-%03d", index), map[string]any{
			"id":      fmt.Sprintf("CON-limit-%03d", index),
			"subject": map[string]any{"type": "issue", "id": a.issueID},
			"scope":   map[string]any{"paths": []string{}}, "severity": "minor",
			"summary": "bounded judgment", "reasoning": "retained judgment", "references": []any{},
		})
		p23AssertReviewSuccess(t, "record bounded Concern", response)
	}
	result := a.p22Run(t, "resume", a.issueID, "--repository", a.workspace)
	assertP22ProcessFailure(t, "projection.resume", "result-limit-exceeded", 6, result)
	if !strings.Contains(result.Stderr, `"collection":"concerns"`) || !strings.Contains(result.Stderr, `"limit":100`) {
		t.Fatalf("resume overflow details = %q, want complete Concern bound", result.Stderr)
	}
}
