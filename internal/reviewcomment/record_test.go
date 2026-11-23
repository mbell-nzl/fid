package reviewcomment

import (
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func TestRoundCommentAttributionAndDispositions(t *testing.T) {
	t.Parallel()
	r := reviewrecord.Round{ID: "RND-one"}
	r.Reviewer = reviewrecord.Reviewer{Actor: "agent:reviewer", Model: "claude/opus-5.5", Effort: "xhigh", SessionID: "session-2"}
	r.Stream, r.Profiles, r.Verdict = "separate", []string{"code", "test-adequacy"}, "APPROVE"
	r.Purpose, r.Verification = "Check the fix", "Focused tests passed"
	r.Dispositions = []reviewrecord.Disposition{{FindingID: "RND-earlier-1", Disposition: "resolved", Reason: "Regression now covers the refusal"}}
	body := RenderRoundV1(r, 2)
	for _, want := range []string{"## Review: APPROVE", "agent:reviewer (claude/opus-5.5), effort xhigh", "Profiles: code, test-adequacy · Round 2", "Session: session-2 · Stream: separate · RND-one", "Claim: Check the fix", "### Findings\n\nNone.", "### Earlier findings", "RND-earlier-1 · resolved: Regression now covers the refusal", "### Verification\n\nFocused tests passed"} {
		t.Run(want, func(t *testing.T) {
			t.Parallel()
			if !strings.Contains(body, want) {
				t.Fatalf("comment lacks %q:\n%s", want, body)
			}
		})
	}
}
