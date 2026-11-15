package main

import (
	"testing"
)

func TestFindingClosedHistoryCannotBeReplaced(t *testing.T) {
	t.Parallel()
	for _, byRound := range []bool{false, true} {
		t.Run(map[bool]string{false: "person first", true: "round first"}[byRound], func(t *testing.T) {
			t.Parallel()
			f, ids, task := closureFixture(t)
			in := reviewInput()
			in["verdict"] = "CHANGES"
			in["dispositions"] = []any{disposition(ids[0], "resolved")}
			if byRound {
				lifecycleResult(t, recordReviewTask(t, f, in, "round", "human:courier", task), "review.record")
			} else {
				lifecycleResult(t, closeFinding(f, ids[0], "human:mark", "Accepted", "person"), "finding.close")
			}
			before := reviewProjection(t, f)
			if got := closeFinding(f, ids[0], "human:mark", "Another decision", "second"); got.err == nil || reviewProjection(t, f) != before {
				t.Fatal("person replaced closed history")
			}
			if got := recordReviewTask(t, f, in, "second-round", "human:courier", task); got.err == nil || reviewProjection(t, f) != before {
				t.Fatal("round replaced closed history")
			}
		})
	}
}

func TestFindingCloseOptionalKeyAndUnknownFinding(t *testing.T) {
	t.Parallel()
	f, ids, _ := closureFixture(t)
	lifecycleResult(t, f.run("--json", "finding", "close", ids[0], "--actor", "human:mark", "--reason", "Accepted"), "finding.close")
	before := reviewProjection(t, f)
	if got := closeFinding(f, "RND-ffffff-1", "human:mark", "Missing", "missing"); got.err == nil || reviewProjection(t, f) != before {
		t.Fatal("closed nonexistent finding")
	}
}

func TestFindingCloseRetryRefusesCorruptStoredResponse(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"reason", "decision_id", "occurred_at"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			f, ids, _ := closureFixture(t)
			lifecycleResult(t, closeFinding(f, ids[0], "human:mark", "Accepted", "close"), "finding.close")
			mutateStoredRetry(t, f, "close", func(result map[string]any) {
				if field == "occurred_at" {
					delete(result, field)
				} else {
					result["fields"].(map[string]any)[field] = "Forged"
				}
			})
			before := reviewProjection(t, f)
			if got := closeFinding(f, ids[0], "human:mark", "Accepted", "close"); got.err == nil || reviewProjection(t, f) != before {
				t.Fatal("rendered corrupt closure response")
			}
		})
	}
}

func TestFindingLaterChangesCannotBeMaskedByEarlierPass(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "pass", "human:courier", task), "review.record")
	in := reviewInput()
	in["verdict"] = "CHANGES"
	in["findings"] = []any{reviewFinding("major")}
	lifecycleResult(t, recordReviewTask(t, f, in, "changes", "human:courier", task), "review.record")
	before := reviewProjection(t, f)
	if got := recordReviewTask(t, f, reviewInput(), "conceal", "human:courier", task); got.err == nil || reviewProjection(t, f) != before {
		t.Fatal("earlier pass concealed later major")
	}
}
