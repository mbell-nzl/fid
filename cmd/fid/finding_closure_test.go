package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func closureFixture(t *testing.T) (*statusFixture, []string, string) {
	t.Helper()
	f, task := canonicalReviewFixture(t)
	in := reviewInput()
	in["verdict"] = "CHANGES"
	in["findings"] = []any{reviewFinding("major"), reviewFinding("major"), reviewFinding("minor"), reviewFinding("nit")}
	fields, _ := lifecycleResult(t, recordReviewTask(t, f, in, "initial", "human:courier", task), "review.record")
	ids := []string{}
	for _, finding := range fields["round"].(map[string]any)["findings"].([]any) {
		ids = append(ids, finding.(map[string]any)["id"].(string))
	}
	return f, ids, task
}

func TestFindingClosureAcrossFingerprints(t *testing.T) {
	t.Parallel()
	f, ids, task := closureFixture(t)
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "Fix the reviewed fault\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Fix the findings")
	prepareSecondReviewHandoff(t, f, task)
	reviewHandoffTask(t, f, "fixed-handoff", task)
	in := reviewInput()
	in["dispositions"] = []any{disposition(ids[0], "resolved"), disposition(ids[1], "not-applicable")}
	result := recordReviewTask(t, f, in, "fixed-review", "human:courier", task)
	lifecycleResult(t, result, "review.record")
	views, err := reviewwork.List(f.workspace, task)
	if err != nil || len(views) != 2 || views[0].Current || !views[1].Current || views[0].Findings[0].Open || views[0].Binding.Fingerprint == views[1].Binding.Fingerprint {
		t.Fatalf("cross-fingerprint closure=%v/%v", views, err)
	}
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	open, err := reviewrecord.OpenMajors(conn, task)
	if err != nil || len(open) != 0 {
		t.Fatalf("open majors=%v/%v", open, err)
	}
}

func disposition(id, state string) map[string]any {
	return map[string]any{"finding_id": id, "disposition": state, "reason": "Verified the regression and accepted the fix."}
}

func closeFinding(f *statusFixture, id, actor, reason, key string) commandExecutionResult {
	return f.run("--json", "finding", "close", id, "--actor", actor, "--reason", reason, "--idempotency-key", key)
}

func TestFindingClosuresAndRetries(t *testing.T) {
	t.Parallel()
	f, ids, task := closureFixture(t)
	in := reviewInput()
	in["verdict"] = "CHANGES"
	in["dispositions"] = []any{disposition(ids[0], "resolved"), disposition(ids[2], "not-applicable"), disposition(ids[3], "open")}
	first := recordReviewTask(t, f, in, "dispose", "human:courier", task)
	fields, _ := lifecycleResult(t, first, "review.record")
	roundID := fields["round"].(map[string]any)["id"]
	before := reviewProjection(t, f)
	if got := recordReviewTask(t, f, in, "dispose", "human:courier", task); got.err != nil || got.stdout != first.stdout || reviewProjection(t, f) != before {
		t.Fatalf("round retry=%#v", got)
	}
	closed := closeFinding(f, ids[1], "human:mark", "Accepted risk", "person-close")
	decision, _ := lifecycleResult(t, closed, "finding.close")
	before = reviewProjection(t, f)
	if got := closeFinding(f, ids[1], "human:mark", "Accepted risk", "person-close"); got.err != nil || got.stdout != closed.stdout || reviewProjection(t, f) != before {
		t.Fatalf("decision retry=%#v", got)
	}
	conflict := closeFinding(f, ids[1], "human:mark", "Different reason", "person-close")
	assertJSONFailureCode(t, conflict.stderr, "finding.close", "idempotency-conflict")
	for _, rebuild := range []bool{false, true} {
		if rebuild {
			if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
				t.Fatal(err)
			}
		}
		views, err := reviewwork.List(f.workspace, task)
		if err != nil {
			t.Fatal(err)
		}
		raw := f.run("--json", "review", "list", task).stdout
		if views[0].Findings[0].Open || views[0].Findings[1].Open || views[0].Findings[2].Open || !views[0].Findings[3].Open || !strings.Contains(raw, fmt.Sprint(roundID)) || !strings.Contains(raw, fmt.Sprint(decision["decision_id"])) || !strings.Contains(raw, "Accepted risk") {
			t.Fatalf("closures missing: %s", raw)
		}
	}
	// A person's close changes the open-major gate without rewriting CHANGES.
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "approve", "human:courier", task), "review.record")
}

func TestFindingClosureRefusals(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, actor, reason string
		twice               bool
	}{
		{"agent", "agent:builder", "Fixed", false}, {"service", "service:ci", "Fixed", false},
		{"blank reason", "human:mark", " ", false}, {"double close", "human:mark", "Accepted", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, ids, _ := closureFixture(t)
			if tc.twice {
				lifecycleResult(t, closeFinding(f, ids[0], tc.actor, tc.reason, "first"), "finding.close")
			}
			before := reviewProjection(t, f)
			got := closeFinding(f, ids[0], tc.actor, tc.reason, "second")
			if got.err == nil || reviewProjection(t, f) != before {
				t.Fatalf("refusal=%#v", got)
			}
		})
	}
}

func TestFindingClosureRefusesOversizedEvent(t *testing.T) {
	t.Parallel()
	f, ids, _ := closureFixture(t)
	before := reviewProjection(t, f)
	got := closeFinding(f, ids[0], "human:mark", strings.Repeat("x", 600000), "oversized")
	if got.err == nil || !strings.Contains(got.stderr, "1 MiB") || reviewProjection(t, f) != before {
		t.Fatalf("oversize refusal error=%v, state changed=%t", got.err, reviewProjection(t, f) != before)
	}
}
