package next

import (
	"slices"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func TestReviewOrderFinishesOpenWorkFirst(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"resume", "rework", "review", "reserved-review"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			work := issues.Work{Epics: map[string]issues.WorkEpic{}, Tasks: []issues.WorkTask{
				task("new", "", "unstarted", false), task("unblock", "", "unstarted", true), task("review", "", "started", false), task("reserved", "", "started", false),
			}}
			work.Tasks[2].Status, work.Tasks[3].Status = "in_review", "in_review"
			reviews := map[string]reviewCandidate{"review": {work: &doneeval.ReviewWork{}}, "reserved": {work: &doneeval.ReviewWork{}}}
			var open []holds.Hold
			want := []string{"review", "reserved", "unblock", "new"}
			if mode == "resume" || mode == "rework" {
				work.Tasks = append(work.Tasks, task("build", "", "started", false))
				h := hold("HLD-1", "build", here, 10*time.Minute)
				if mode == "rework" {
					stamp, reason := holds.Timestamp(orderNow), "review"
					h.EndedAt, h.ReworkAt, h.EndReason, h.Kept = &stamp, &stamp, &reason, true
				}
				open = []holds.Hold{h}
				want = append([]string{"build"}, want...)
			}
			if mode == "resume" {
				// All five kinds coexist; the foreign rework stays reserved
				// while this worktree resumes its own started task first.
				work.Tasks = append(work.Tasks, task("foreign-rework", "", "started", false))
				foreign := hold("HLD-foreign", "foreign-rework", "sha256:other", time.Minute)
				stamp, reason := holds.Timestamp(orderNow), "review"
				foreign.EndedAt, foreign.ReworkAt, foreign.EndReason, foreign.Kept = &stamp, &stamp, &reason, true
				open = append(open, foreign)
				want = append(want, "foreign-rework")
			}

			if mode == "reserved-review" {
				reviews["reserved"] = reviewCandidate{work: &doneeval.ReviewWork{}, claim: &reviewrecord.Claim{ID: "RCL-1", TaskID: "reserved", WorktreeID: here, ExpiresAt: holds.Timestamp(orderNow.Add(time.Hour))}}
				want = []string{"reserved", "review", "unblock", "new"}
			}
			got := plan(work, open, here, orderNow, reviews)
			ids := []string{}
			for _, entry := range got.queue {
				ids = append(ids, entry.ID)
			}
			kind := mode
			if mode == "reserved-review" {
				kind = "review"
			}
			if !slices.Equal(ids, want) || got.served != 0 || got.kind != kind {
				t.Fatalf("order=%v kind=%s served=%d; want %v / %s", ids, got.kind, got.served, want, kind)
			}
		})
	}
}
