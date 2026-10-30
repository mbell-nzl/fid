//go:build integration

package integration

import (
	"strings"
	"testing"
)

func TestP23ConcernListHumanProjection(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "concern.dispose", "concern.record", "envelope.activate", "envelope.create", "gate.register", "run.observe", "run.start"}
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
	if err != nil || stderr != "" || !strings.Contains(human, "Purpose: Inspect a raised concern") || !strings.Contains(human, "CON-human-list") || !strings.Contains(human, "Summary: private concern prose") || !strings.Contains(human, "Reasoning: private reasoning") {
		t.Fatalf("generic Concern Human list = %q / %q / %v", human, stderr, err)
	}
	second := p23ConcernMutation(t, a, "record", "human-list-minor", map[string]any{
		"id": "CON-human-list-minor", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"scope": map[string]any{"paths": []string{}}, "severity": "minor",
		"summary": "other private prose", "reasoning": "other private reasoning", "references": []any{},
	})
	p23AssertReviewSuccess(t, "record second generic Concern", second)
	for _, test := range []struct {
		name string
		args []string
		want string
	}{
		{"severity", []string{"--severity", "major"}, "CON-human-list"},
		{"status", []string{"--status", "open", "--severity", "minor"}, "CON-human-list-minor"},
		{"subject", []string{"--subject", "requirement:REQ-other"}, ""},
		{"disposed", []string{"--status", "disposed"}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			result := a.p22Run(t, append([]string{"concern", "list"}, test.args...)...)
			p23AssertReviewSuccess(t, "filtered Concern list", result)
			items := p21Path(t, p23Envelope(t, result.Stdout), "result", "items").([]any)
			if test.want == "" && len(items) == 0 {
				return
			}
			if len(items) != 1 || items[0].(map[string]any)["id"] != test.want {
				t.Fatalf("filtered Concern list = %#v, want %q", items, test.want)
			}
		})
	}
	before := p21JournalEventCount(t, a.workspace)
	limited := a.p22Run(t, "concern", "list", "--limit", "1")
	if limited.ExitCode == 0 || !strings.Contains(limited.Stderr, `"code":"result-limit-exceeded"`) || limited.Stdout != "" {
		t.Fatalf("bounded Concern list = %#v", limited)
	}
	if after := p21JournalEventCount(t, a.workspace); after != before {
		t.Fatalf("bounded read appended Event(s): %d -> %d", before, after)
	}
	mixed := a.p22Run(t, "concern", "list", "--issue", a.issueID, "--severity", "major")
	if mixed.ExitCode == 0 || mixed.Stdout != "" || !strings.Contains(mixed.Stderr, `"code":"invalid-argument"`) {
		t.Fatalf("mixed V1 and generic Concern filters = %#v", mixed)
	}
	if after := p21JournalEventCount(t, a.workspace); after != before {
		t.Fatalf("mixed filter read appended Event(s): %d -> %d", before, after)
	}
	disposed := p23ConcernMutation(t, a, "dispose", "human-list-dispose", map[string]any{
		"concern_id": "CON-human-list", "action": "rejected", "issue_id": nil,
		"supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "reviewed",
	})
	p23AssertReviewSuccess(t, "dispose listed Concern", disposed)
	for _, test := range []struct{ status, want string }{{"open", "CON-human-list-minor"}, {"disposed", "CON-human-list"}} {
		result := a.p22Run(t, "concern", "list", "--status", test.status)
		p23AssertReviewSuccess(t, "status-filtered Concern list", result)
		items := p21Path(t, p23Envelope(t, result.Stdout), "result", "items").([]any)
		if len(items) != 1 || items[0].(map[string]any)["id"] != test.want {
			t.Fatalf("%s Concern list = %#v, want %q", test.status, items, test.want)
		}
	}
}
