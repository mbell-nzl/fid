package doneeval

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func TestReviewSelectionUsesDoneGates(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		edit func(*Snapshot)
		want string
	}{
		{"approved", func(*Snapshot) {}, ""},
		{"primary missing", func(s *Snapshot) { s.Rounds = s.Rounds[1:] }, "primary"},
		{"primary stale", func(s *Snapshot) { s.Rounds[0].Current = false }, "primary"},
		{"primary profile", func(s *Snapshot) { s.Rounds[0].Profiles = []string{"code"} }, "primary"},
		{"person closes majors as in V5", func(s *Snapshot) { s.Rounds[0].Verdict = "CHANGES" }, ""},
		{"separate missing", func(s *Snapshot) { s.Rounds = s.Rounds[:1] }, "separate"},
		{"separate stale", func(s *Snapshot) { s.Rounds[1].Current = false }, "separate"},
		{"separate session", func(s *Snapshot) { s.Rounds[1].Reviewer.SessionID = s.Rounds[0].Reviewer.SessionID }, "separate"},
		{"open major", func(s *Snapshot) { s.OpenMajors = []reviewrecord.Finding{{ID: "FND-open"}} }, ""},
		{"rework reviewed", func(s *Snapshot) { s.Task.State = "in_progress" }, ""},
		{"rework", func(s *Snapshot) { s.Task.State = "in_progress"; s.Rounds = nil }, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", true)
			tc.edit(&s)
			got := DueReview(s)
			if tc.want == "" {
				if got != nil {
					t.Fatalf("unexpected review: %v", got)
				}
			} else if got == nil || got.Stream != tc.want || len(got.RequiredProfiles) != 2 {
				t.Fatalf("review=%v want %s / both profiles", got, tc.want)
			}
			outcome := Evaluate(s)
			wantDone := outcome.Gates[PrimaryGate].State == Satisfied && outcome.Gates[SeparateGate].State == Satisfied && outcome.Gates[MajorsGate].State == Satisfied && s.Task.State == "in_review"
			if ReviewsDone(s) != wantDone {
				t.Fatalf("review completion disagrees with done: %v", outcome)
			}
		})
	}
}
