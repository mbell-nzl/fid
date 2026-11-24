package main

import (
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestReviewAdmissionUsesFrozenLaneAfterRaise(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, frozen, current, stream, refusal string
		profiles                               []string
	}{
		{"standard separate", "standard", "assured", "separate", "separate review requires an Assured hand-off", []string{"code"}},
		{"standard missing code", "standard", "assured", "primary", "missing required profile code", []string{}},
		{"routine raised to standard", "routine", "standard", "primary", "", []string{}},
		{"routine raised to assured", "routine", "assured", "primary", "", []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithLane(t, tc.frozen)
			lifecycleResult(t, f.run("--json", "issue", "update", task, "--lane", tc.current, "--actor", "human:mark", "--idempotency-key", "raise-lane"), "issue.update")
			// Exercise admission both before and after rebuilding canonical Events.
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
						t.Fatal(err)
					}
				}
				assertReviewFrozenLane(t, f, task, tc.frozen, tc.current)
				in := reviewInput()
				in["stream"], in["profiles"] = tc.stream, tc.profiles
				key := "before-rebuild"
				if rebuild {
					key = "after-rebuild"
				}
				before := snapshotWorkspacePersistence(t, f.workspace)
				got := recordReviewTask(t, f, in, key, "human:courier", task)
				if tc.refusal != "" {
					if got.err == nil || !strings.Contains(got.stderr, tc.refusal) {
						t.Fatalf("frozen=%s current=%s: err=%v stderr=%s", tc.frozen, tc.current, got.err, got.stderr)
					}
					if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
						t.Fatal("refusal changed journal or projection")
					}
					continue
				}
				lifecycleResult(t, got, "review.record")
				if after := startEventCount(t, f.workspace); after != before.eventCount()+1 {
					t.Fatalf("accepted review appended %d Events, want one", after-before.eventCount())
				}
			}
		})
	}
}

func assertReviewFrozenLane(t *testing.T, f *statusFixture, task, frozen, current string) {
	t.Helper()
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	h, err := reviewrecord.LatestHandoff(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	var lane string
	if err := conn.QueryRow(`SELECT lane FROM issues WHERE id=?`, task).Scan(&lane); err != nil {
		t.Fatal(err)
	}
	if h.Lane != frozen || lane != current || lane == h.Lane {
		t.Fatalf("hand-off lane=%s task lane=%s, want %s/%s", h.Lane, lane, frozen, current)
	}
}
