package doneeval

import "testing"

func TestMergeObservationPRIdentity(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		edit   func(*Snapshot)
		status Status
		state  State
		reason Reason
	}{
		{"same URL", func(*Snapshot) {}, Complete, Satisfied, EvidenceSatisfied},
		{"mixed-case hand-off", func(s *Snapshot) { s.Handoff.PR = "https://github.com/OWNER/Repo/pull/1" }, Complete, Satisfied, EvidenceSatisfied},
		{"mixed-case merge URL", func(s *Snapshot) { s.Merge.PR = "https://github.com/Owner/REPO/pull/1" }, Complete, Satisfied, EvidenceSatisfied},
		{"mixed-case observed URL", func(s *Snapshot) { s.GitHub.PR.URL = "https://github.com/OWNER/REPO/pull/1" }, Complete, Satisfied, EvidenceSatisfied},
		{"invalid merge URL", func(s *Snapshot) { s.Merge.PR = "invalid" }, Refused, Stale, MergeMismatch},
		{"different PR number in URL", func(s *Snapshot) { s.Merge.PR = "https://github.com/owner/repo/pull/2" }, Refused, Stale, MergeMismatch},
		{"different owner in URL", func(s *Snapshot) { s.Merge.PR = "https://github.com/other/repo/pull/1" }, Refused, Stale, MergeMismatch},
		{"different repository in URL", func(s *Snapshot) { s.Merge.PR = "https://github.com/owner/other/pull/1" }, Refused, Stale, MergeMismatch},
		{"different recorded number", func(s *Snapshot) { s.Merge.Number++ }, Refused, Stale, MergeMismatch},
		{"different task", func(s *Snapshot) { s.Merge.TaskID = "other" }, Refused, Stale, MergeMismatch},
		{"different hand-off", func(s *Snapshot) { s.Merge.HandoffID = "other" }, Refused, Stale, MergeMismatch},
		{"different recorded repository", func(s *Snapshot) { s.Merge.Repository = "other/repo" }, Refused, Stale, MergeMismatch},
		{"different base branch", func(s *Snapshot) { s.Merge.BaseBranch = "other" }, Refused, Stale, MergeMismatch},
		{"different merge commit", func(s *Snapshot) { s.Merge.MergeCommit = "other" }, Refused, Stale, MergeMismatch},
		{"different head commit", func(s *Snapshot) { s.Merge.HeadSHA = "other" }, Refused, Stale, MergeMismatch},
		{"different head tree", func(s *Snapshot) { s.Merge.HeadTree = "other" }, Refused, Stale, MergeMismatch},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", true)
			tc.edit(&s)
			got := Evaluate(s)
			gate := got.Gates[MergeGate]
			if got.Status != tc.status || gate.State != tc.state || gate.Reason != tc.reason {
				t.Fatalf("status=%s merge=%+v; want %s %s/%s", got.Status, gate, tc.status, tc.state, tc.reason)
			}
		})
	}
}
