package main

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestReviewGolden(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	in := reviewInput()
	in["verdict"] = "CHANGES"
	in["findings"] = []any{reviewFinding("major"), reviewFinding("minor"), reviewFinding("nit")}
	first := recordReview(t, f, in, "golden", "human:courier")
	lifecycleResult(t, first, "review.record")
	normalize := func(s string) string {
		s = f.placeholders(s)
		s = regexp.MustCompile(`RND-[0-9a-f]{6}`).ReplaceAllString(s, "RND-<new>")
		s = regexp.MustCompile(`HOF-[0-9a-f]{6}`).ReplaceAllString(s, "HOF-<new>")
		s = regexp.MustCompile(`\b(?:[0-9a-f]{64}|[0-9a-f]{40})\b`).ReplaceAllString(s, "<fingerprint-or-commit>")
		return s
	}
	assertGolden(t, filepath.Join("review", "record.json"), normalize(first.stdout))
	for _, jsonMode := range []bool{false, true} {
		args, ext := []string{"review", "list", "FID-9f13aa"}, "txt"
		if jsonMode {
			args = append([]string{"--json"}, args...)
			ext = "json"
		}
		got := f.run(args...)
		if got.err != nil {
			t.Fatal(got)
		}
		assertGolden(t, filepath.Join("review", "list."+ext), normalize(got.stdout))
	}
	// The house list is a 72-column box; external review prose is not a label.
	text := f.run("review", "list", "FID-9f13aa")
	if !strings.HasPrefix(text.stdout, "╭ fid review") || strings.Contains(text.stdout, "claim") || strings.Contains(text.stdout, "runs") {
		t.Fatalf("house text=%q", text.stdout)
	}
}
