package doneeval

import (
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/mergerecord"
)

func evaluateRemote(o *Outcome, s Snapshot) {
	pr, h := s.GitHub.PR, s.Handoff
	if pr == nil {
		o.set(PRGate, Missing, PRUnavailable)
		return
	}
	expected, err := github.ParsePRURL(h.PR)
	actual, actualErr := github.ParsePRURL(pr.URL)
	o.set(PRGate, Satisfied, EvidenceSatisfied)
	switch {
	case err != nil || actualErr != nil || actual.Number != expected.Number || pr.Number != expected.Number:
		o.set(PRGate, RefusedState, PRMismatch)
	case !expected.Repository.Equal(s.Target.Repository) || !actual.Repository.Equal(s.Target.Repository) || !pr.BaseRepository.Equal(s.Target.Repository):
		o.set(PRGate, RefusedState, RepositoryMismatch)
	case s.Target.DefaultBranch == "" || pr.BaseRef != s.Target.DefaultBranch:
		o.set(PRGate, RefusedState, BranchMismatch)
	case !pr.Merged || pr.State != "MERGED":
		o.set(PRGate, WaitingState, NotMerged)
	case pr.MergeCommit == "":
		o.set(PRGate, RefusedState, MergeInvalid)
	}
	if !pr.Merged || pr.State != "MERGED" {
		for _, g := range []Gate{TreeGate, MergeGate, CIGate} {
			o.set(g, WaitingState, NotMerged)
		}
		return
	}
	o.set(TreeGate, Satisfied, EvidenceSatisfied)
	if s.GitHub.HeadTree == "" || s.GitHub.HeadTree != h.Tree {
		o.set(TreeGate, Stale, HeadTreeMismatch)
	}
	o.set(MergeGate, Satisfied, EvidenceSatisfied)
	m := s.Merge
	if m == nil {
		o.set(MergeGate, Missing, MergeMissing)
	} else {
		recorded, err := github.ParsePRURL(m.PR)
		if err != nil || recorded != expected || m.TaskID != s.Task.ID || m.HandoffID != h.ID || m.Number != pr.Number || m.Repository != s.Target.Repository.String() || m.BaseBranch != pr.BaseRef || m.MergeCommit != pr.MergeCommit || m.HeadSHA != pr.HeadSHA || m.HeadTree != s.GitHub.HeadTree {
			o.set(MergeGate, Stale, MergeMismatch)
		}
	}
	evaluateCI(o, s)
}

func evaluateCI(o *Outcome, s Snapshot) {
	ci, pr := s.GitHub.CI, s.GitHub.PR
	if ci == nil {
		o.set(CIGate, WaitingState, CIWaiting)
		return
	}
	if ci.Commit != pr.MergeCommit {
		o.set(CIGate, RefusedState, CIMismatch)
		return
	}
	// Reuse the adapter's pure rule, never a second attempt-selection algorithm.
	selected := github.EvaluateCI(ci.Commit, ci.Runs)
	switch selected.Verdict {
	case github.Waiting:
		o.set(CIGate, WaitingState, CIWaiting)
	case github.NotGreen:
		o.set(CIGate, RefusedState, CINotGreen)
	case github.Green:
		o.set(CIGate, Satisfied, EvidenceSatisfied)
		for _, run := range selected.Runs {
			if run.Ignored != "" {
				continue
			}
			recorded := false
			for _, r := range s.CIResults {
				recorded = recorded || matchesCI(s, r, run)
			}
			if !recorded {
				o.set(CIGate, Missing, CIUnrecorded)
			}
		}
	}
}

func matchesCI(s Snapshot, r mergerecord.CIResult, run github.WorkflowRun) bool {
	return s.Merge != nil && r.ID != "" && r.TaskID == s.Task.ID && r.MergeID == s.Merge.ID && r.Repository == s.Target.Repository.String() && r.RunID == run.ID && r.Attempt == run.Attempt && r.WorkflowID == run.WorkflowID && r.Name == run.Name && r.Path == run.Path && r.Event == run.Event && r.HeadSHA == run.HeadSHA && r.Status == run.Status && r.Conclusion == run.Conclusion && r.URL == run.URL
}
