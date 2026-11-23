package main

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func TestWaiverPresentationAndCurrency(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"current", "unchanged", "stale"} {
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			f := approvalDisplayFixture(t)
			waiveSuccess(t, waiveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--lane", "standard"))
			if state != "current" {
				lifecycleResult(t, recordReview(t, f, changesInput("primary"), "changes", "human:courier"), "review.record")
				prepareSecondReviewHandoff(t, f, "FID-9f13aa")
				if state == "stale" {
					writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "edited change\n")
					runGitForTest(t, f.workspace, "add", ".")
					runGitForTest(t, f.workspace, "commit", "-m", "Edit waived change")
				}
				reviewHandoffTask(t, f, "next-handoff", "FID-9f13aa")
			}
			conn, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			snapshot, err := doneeval.GatherReviews(f.workspace, conn, "FID-9f13aa", nil)
			_ = conn.Close()
			if err != nil {
				t.Fatal(err)
			}
			gates := doneeval.Evaluate(snapshot).Gates
			if (gates[doneeval.ApprovalGate].Reason == doneeval.NotRequired) != (state != "stale") || (gates[doneeval.SeparateGate].Reason == doneeval.NotRequired) != (state != "stale") {
				t.Fatalf("%s gates=%v", state, gates)
			}
			before := reviewProjection(t, f)
			fields, _ := lifecycleResult(t, f.run("--json", "show", "FID-9f13aa"), "issue.show")
			waivers, ok := fields["waivers"].([]any)
			if !ok || len(waivers) != 1 || waivers[0].(map[string]any)["current"] != (state != "stale") {
				t.Fatalf("waiver display: %v", fields["waivers"])
			}
			if state != "stale" && fields["approval"].(map[string]any)["state"] != "not-required" {
				t.Fatal("approval still required", fields)
			}
			show, status := f.run("show", "FID-9f13aa"), f.run("status")
			if !strings.Contains(show.stdout, "waiver") || (state == "stale" && (!strings.Contains(show.stdout, "stale") || !strings.Contains(status.stdout, "stale waiver"))) {
				t.Fatalf("show=%s status=%s", show.stdout, status.stdout)
			}
			if reviewProjection(t, f) != before {
				t.Fatal("display changed ledger")
			}
			if state != "unchanged" {
				approvalDisplayGoldens(t, f, "waiver-"+state)
			}
		})
	}
}

func TestWaivedRoutineAcceptsOnlyRequiredTestAdequacyProfile(t *testing.T) {
	t.Parallel()
	f, task := waivedRoutineTestAdequacyFixture(t)
	in := reviewInput()
	in["profiles"] = []string{"test-adequacy"}
	lifecycleResult(t, recordReviewTask(t, f, in, "waived-profile", "human:courier", task), "review.record")
}

func waivedRoutineTestAdequacyFixture(t *testing.T) (*statusFixture, string) {
	t.Helper()
	f, task := canonicalReviewFixtureWithClock(t, "assured", false, time.Now, func(f *statusFixture) {
		trustedShowPolicy(t, f, func(raw string) string {
			updated := strings.Replace(raw, "test_adequacy = []", "test_adequacy = ['**/*.md']", 1)
			if updated == raw {
				t.Fatal("fixture did not configure Markdown paths for test adequacy")
			}
			return updated
		})
	})
	waiveSuccess(t, waiveForTest(f, task, task+"\n", true))
	return f, task
}

func TestLatestCurrentWaiverSetsEffectiveLane(t *testing.T) {
	t.Parallel()
	f := approvalDisplayFixture(t)
	waiveSuccess(t, waiveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--lane", "standard"))
	latest := waiveSuccess(t, waiveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--lane", "routine"))
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	s, err := doneeval.GatherReviews(f.workspace, conn, "FID-9f13aa", nil)
	if err != nil || s.EffectiveLane() != "routine" || s.Waiver.Decision.ID != latest["decision"].(map[string]any)["decision_id"] {
		t.Fatalf("effective waiver=%+v/%v", s.Waiver, err)
	}
	got := f.run("show", "FID-9f13aa")
	if got.err != nil || !strings.Contains(got.stdout, "waiver to routine") {
		t.Fatal("show chose an earlier waiver", got)
	}
}
