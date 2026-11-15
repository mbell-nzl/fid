package main

import (
	"strings"
	"testing"
)

func TestReviewRetryPreservesExactIntegerLocation(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	in := reviewInput()
	finding := reviewFinding("minor")
	in["findings"] = []any{finding}
	location := finding["location"].(map[string]any)
	// These adjacent, valid integers collide if request decoding uses float64.
	location["line"] = int64(9007199254740992)
	first := recordReview(t, f, in, "integer-line", "human:courier")
	lifecycleResult(t, first, "review.record")
	if !strings.Contains(first.stdout, `"line":9007199254740992`) {
		t.Fatal("first line was rounded")
	}
	before := reviewProjection(t, f)
	location["line"] = int64(9007199254740993)
	got := recordReview(t, f, in, "integer-line", "human:courier")
	if got.err == nil {
		t.Fatal("different integer line replayed the original review")
	}
	assertJSONFailureCode(t, got.stderr, "review.record", "idempotency-conflict")
	if !strings.Contains(got.stderr, `"differing_fields":["findings"]`) || reviewProjection(t, f) != before {
		t.Fatalf("integer conflict=%#v", got)
	}
}
