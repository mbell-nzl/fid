package main

import (
	"path/filepath"
	"regexp"
	"testing"
)

func TestFindingClosureGolden(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	in := reviewInput()
	in["verdict"] = "CHANGES"
	in["findings"] = []any{reviewFinding("major"), reviewFinding("major"), reviewFinding("minor")}
	fields, _ := lifecycleResult(t, recordReview(t, f, in, "first", "human:courier"), "review.record")
	findings := fields["round"].(map[string]any)["findings"].([]any)
	ids := []string{}
	for _, finding := range findings {
		ids = append(ids, finding.(map[string]any)["id"].(string))
	}
	in = reviewInput()
	in["verdict"] = "CHANGES"
	in["dispositions"] = []any{disposition(ids[0], "resolved"), disposition(ids[2], "not-applicable")}
	lifecycleResult(t, recordReview(t, f, in, "later", "human:courier"), "review.record")
	normalize := func(s string) string {
		s = f.placeholders(s)
		s = regexp.MustCompile(`RND-[0-9a-f]{6}`).ReplaceAllString(s, "RND-<new>")
		s = regexp.MustCompile(`PDN-[0-9a-f]{6}`).ReplaceAllString(s, "PDN-<new>")
		s = regexp.MustCompile(`HOF-[0-9a-f]{6}`).ReplaceAllString(s, "HOF-<new>")
		return regexp.MustCompile(`\b(?:[0-9a-f]{64}|[0-9a-f]{40})\b`).ReplaceAllString(s, "<fingerprint-or-commit>")
	}
	for _, jsonMode := range []bool{true, false} {
		args := []string{"finding", "close", ids[1], "--actor", "human:mark", "--reason", "Accepted risk", "--idempotency-key", "close"}
		ext := "txt"
		if jsonMode {
			args = append([]string{"--json"}, args...)
			ext = "json"
		}
		got := f.run(args...)
		if got.err != nil {
			t.Fatal(got)
		}
		assertGolden(t, filepath.Join("finding", "close."+ext), normalize(got.stdout))
		args = []string{"review", "list", "FID-9f13aa"}
		if jsonMode {
			args = append([]string{"--json"}, args...)
		}
		got = f.run(args...)
		if got.err != nil {
			t.Fatal(got)
		}
		assertGolden(t, filepath.Join("finding", "list."+ext), normalize(got.stdout))
	}
}
