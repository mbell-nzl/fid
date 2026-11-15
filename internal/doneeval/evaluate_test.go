package doneeval

import (
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/mergerecord"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func completeSnapshot(lane string, triggered bool) Snapshot {
	tree, head, base, merge := strings.Repeat("a", 40), strings.Repeat("b", 40), strings.Repeat("c", 40), strings.Repeat("d", 40)
	check := config.Check{Run: []string{"true"}, Timeout: time.Minute}
	p := policy.Resolve(config.ProjectFile{Checks: map[string]config.Check{"test": check}, Lanes: map[string]config.Lane{lane: {Checks: []string{"test"}, ChecksSet: true}}, TestAdequacySet: true, TestAdequacy: []string{"**/*_test.go"}})
	h := &handoffrecord.Record{View: handoffrecord.View{ID: "HOF-one", PR: "https://github.com/owner/repo/pull/1", Head: head, Base: base, RequiredChecks: []string{"test"}, ResultIDs: []string{"CHK-one"}}, TaskID: "FID-one", Tree: tree, Lane: lane, PolicySnapshot: 1, DefinitionDigests: map[string]string{"test": checkrecord.Digest(checkrecord.DefinitionFor(check))}, ChangedPaths: []string{"notes.md"}}
	if triggered {
		h.ChangedPaths = []string{"internal/x_test.go"}
	}
	pr := &github.PullRequest{Number: 1, URL: h.PR, State: "MERGED", Merged: true, MergeCommit: merge, HeadSHA: head, BaseRef: "main", BaseRepository: github.Repository{Owner: "owner", Name: "repo"}}
	profiles, _ := reviewrecord.RequiredProfiles(*h, p)
	s := Snapshot{Task: Task{ID: h.TaskID, Kind: "task", State: "in_review", Lane: lane}, Handoff: h, Policy: &p, ArtifactsVerified: true, Target: Target{Repository: pr.BaseRepository, DefaultBranch: "main"}, GitHub: GitHubFacts{PR: pr, HeadTree: tree}, Approval: approvalrecord.View{State: "current"}}
	s.Checks = []checkrecord.Record{{Start: checkrecord.Start{ID: "CHK-one", TaskID: h.TaskID, Tree: checkrecord.Tree{Tree: tree}, Name: "test", DefinitionDigest: h.DefinitionDigests["test"]}, Result: &checkrecord.Result{Counts: true, Outcome: "passed"}}}
	s.Rounds = []reviewrecord.RoundView{{Round: reviewrecord.Round{ID: "RND-one", TaskID: h.TaskID, Stream: "primary", Profiles: profiles, Reviewer: reviewrecord.Reviewer{SessionID: "first"}, Verdict: "APPROVE"}, Current: true, OtherModel: true}}
	s.GitHub.PR.Body = "## Rollback\nRevert."
	if lane == "assured" {
		s.Rounds = append(s.Rounds, reviewrecord.RoundView{Round: reviewrecord.Round{ID: "RND-two", TaskID: h.TaskID, Stream: "separate", Profiles: profiles, Reviewer: reviewrecord.Reviewer{SessionID: "second"}, Verdict: "APPROVE"}, Current: true})
	}
	s.Merge = &mergerecord.Merge{ID: "MRG-one", TaskID: h.TaskID, HandoffID: h.ID, PR: h.PR, Number: 1, Repository: "owner/repo", BaseBranch: "main", MergeCommit: merge, HeadSHA: head, HeadTree: tree}
	run := github.WorkflowRun{ID: 1, Attempt: 2, WorkflowID: 1, Name: "ci", Path: ".github/workflows/ci.yml", Event: "push", HeadSHA: merge, Status: "completed", Conclusion: "success", URL: "https://github.com/owner/repo/actions/runs/1"}
	s.GitHub.CI = &github.CI{Commit: merge, Runs: []github.WorkflowRun{run}}
	s.CIResults = []mergerecord.CIResult{{ID: "CIR-one", TaskID: h.TaskID, MergeID: s.Merge.ID, Repository: "owner/repo", RunID: 1, Attempt: 2, WorkflowID: 1, Name: run.Name, Path: run.Path, Event: run.Event, HeadSHA: merge, Status: run.Status, Conclusion: run.Conclusion, URL: run.URL}}
	return s
}

type matrixCase struct {
	name    string
	edit    func(*Snapshot)
	gate    Gate
	reason  Reason
	state   State
	outcome Status
}

func TestLaneMatrix(t *testing.T) {
	t.Parallel()
	for _, lane := range []string{"routine", "routine-tests", "standard", "assured"} {
		t.Run(lane, func(t *testing.T) {
			t.Parallel()
			actual := strings.TrimSuffix(lane, "-tests")
			cases := []matrixCase{
				{"complete", func(*Snapshot) {}, "", "", "", Complete},
				{"missing task", func(s *Snapshot) { s.Task.ID = "" }, TaskGate, TaskUnavailable, Missing, Refused},
				{"not in review", func(s *Snapshot) { s.Task.State = "in_progress" }, TaskGate, TaskNotInReview, RefusedState, Refused},
				{"missing handoff", func(s *Snapshot) { s.Handoff = nil }, HandoffGate, HandoffMissing, Missing, Refused},
				{"lane raised", func(s *Snapshot) {
					s.Task.Lane = "assured"
					if actual == "assured" {
						s.Task.Lane = "standard"
					}
				}, LaneGate, LaneChanged, Stale, Refused},
				{"missing policy", func(s *Snapshot) { s.Policy = nil }, PolicyGate, PolicyMissing, Missing, Refused},
				{"checks unset", func(s *Snapshot) { l := s.Policy.Lanes[actual]; l.ChecksSet = false; s.Policy.Lanes[actual] = l }, PolicyGate, ChecksNotConfigured, Missing, Refused},
				{"test paths unset", func(s *Snapshot) { s.Policy.TestAdequacySet = false }, PolicyGate, TestPathsNotConfigured, Missing, Refused},
				{"missing local result", func(s *Snapshot) { s.Checks = nil }, ChecksGate, ChecksUnsatisfied, Missing, Refused},
				{"later failure", func(s *Snapshot) {
					r := s.Checks[0]
					r.Result = &checkrecord.Result{Counts: true, Outcome: "failed"}
					s.Checks = append(s.Checks, r)
				}, ChecksGate, ChecksUnsatisfied, Missing, Refused},
				{"wrong local tree", func(s *Snapshot) { s.Checks[0].Start.Tree.Tree = "other" }, ChecksGate, ChecksUnsatisfied, Missing, Refused},
				{"wrong definition", func(s *Snapshot) { s.Checks[0].Start.DefinitionDigest = "other" }, ChecksGate, ChecksUnsatisfied, Missing, Refused},
				{"kept bytes missing", func(s *Snapshot) { s.ArtifactProblems = []checkread.ArtifactProblem{{ArtifactID: "ART-one"}} }, ArtifactsGate, KeptArtifactUnavailable, Missing, Refused},
				{"kept bytes not observed", func(s *Snapshot) { s.ArtifactsVerified = false }, ArtifactsGate, ArtifactsNotVerified, Missing, Refused},
				{"open major", func(s *Snapshot) { s.OpenMajors = []reviewrecord.Finding{{ID: "major"}} }, MajorsGate, OpenMajor, RefusedState, Refused},
				{"not merged", func(s *Snapshot) {
					s.GitHub.PR.Merged = false
					s.GitHub.PR.State = "OPEN"
					s.GitHub.PR.MergeCommit = ""
					s.Merge = nil
					s.GitHub.CI = nil
				}, PRGate, NotMerged, WaitingState, Waiting},
				{"closed not merged", func(s *Snapshot) {
					s.GitHub.PR.Merged = false
					s.GitHub.PR.State = "CLOSED"
					s.GitHub.PR.MergeCommit = ""
					s.Merge = nil
					s.GitHub.CI = nil
				}, PRGate, NotMerged, WaitingState, Waiting},
				{"wrong PR", func(s *Snapshot) { s.GitHub.PR.Number = 2; s.GitHub.PR.URL = "https://github.com/owner/repo/pull/2" }, PRGate, PRMismatch, RefusedState, Refused},
				{"wrong repository", func(s *Snapshot) { s.GitHub.PR.BaseRepository.Name = "other" }, PRGate, RepositoryMismatch, RefusedState, Refused},
				{"wrong branch", func(s *Snapshot) { s.GitHub.PR.BaseRef = "other" }, PRGate, BranchMismatch, RefusedState, Refused},
				{"head tree mismatch", func(s *Snapshot) { s.GitHub.HeadTree = "other" }, TreeGate, HeadTreeMismatch, Stale, Refused},
				{"no merge receipt", func(s *Snapshot) { s.Merge = nil }, MergeGate, MergeMissing, Missing, Refused},
				{"stale merge receipt", func(s *Snapshot) { s.Merge.HandoffID = "old" }, MergeGate, MergeMismatch, Stale, Refused},
				{"CI none", func(s *Snapshot) { s.GitHub.CI.Runs = nil }, CIGate, CIWaiting, WaitingState, Waiting},
				{"CI pending", func(s *Snapshot) {
					r := s.GitHub.CI.Runs[0]
					r.Attempt++
					r.Status = "in_progress"
					r.Conclusion = ""
					s.GitHub.CI.Runs = append(s.GitHub.CI.Runs, r)
				}, CIGate, CIWaiting, WaitingState, Waiting},
				{"CI failed", func(s *Snapshot) { s.GitHub.CI.Runs[0].Conclusion = "failure" }, CIGate, CINotGreen, RefusedState, Refused},
				{"CI wrong commit", func(s *Snapshot) { s.GitHub.CI.Commit = "other" }, CIGate, CIMismatch, RefusedState, Refused},
				{"CI receipt absent", func(s *Snapshot) { s.CIResults = nil }, CIGate, CIUnrecorded, Missing, Refused},
				{"CI green after failure", func(s *Snapshot) {
					r := s.GitHub.CI.Runs[0]
					r.Attempt = 1
					r.Conclusion = "failure"
					s.GitHub.CI.Runs = append(s.GitHub.CI.Runs, r)
				}, "", "", "", Complete},
				{"unfinished local start", func(s *Snapshot) { r := s.Checks[0]; r.Result = nil; s.Checks = append(s.Checks, r) }, "", "", "", Complete},
			}
			if actual != "routine" || lane == "routine-tests" {
				cases = append(cases, matrixCase{"review missing", func(s *Snapshot) { s.Rounds = nil }, PrimaryGate, ReviewMissing, Missing, Refused},
					matrixCase{"review stale", func(s *Snapshot) { s.Rounds[0].Current = false }, PrimaryGate, ReviewStale, Stale, Refused},
					matrixCase{"profile missing", func(s *Snapshot) { s.Rounds[0].Profiles = nil }, PrimaryGate, ProfileMissing, Missing, Refused})
			}
			if actual == "assured" {
				cases = append(cases, []matrixCase{
					{"separate missing", func(s *Snapshot) { s.Rounds = s.Rounds[:1] }, SeparateGate, ReviewMissing, Missing, Refused},
					{"separate stale", func(s *Snapshot) { s.Rounds[1].Current = false }, SeparateGate, ReviewStale, Stale, Refused},
					{"separate same session", func(s *Snapshot) { s.Rounds[1].Reviewer.SessionID = "first" }, SeparateGate, SeparateSessionRequired, RefusedState, Refused},
					{"approval missing", func(s *Snapshot) { s.Approval.State = "missing" }, ApprovalGate, ApprovalMissing, Missing, Refused},
					{"approval stale", func(s *Snapshot) { s.Approval.State = "stale" }, ApprovalGate, ApprovalStale, Stale, Refused},
					{"rollback missing", func(s *Snapshot) { s.GitHub.PR.Body = "" }, RollbackGate, RollbackMissing, Missing, Refused},
					{"rollback fenced", func(s *Snapshot) { s.GitHub.PR.Body = "```markdown\n## Rollback\n```" }, RollbackGate, RollbackMissing, Missing, Refused},
				}...)
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					s := completeSnapshot(actual, lane == "routine-tests")
					tc.edit(&s)
					got := Evaluate(s)
					if got.Status != tc.outcome {
						t.Fatalf("status=%s want %s; gates=%+v", got.Status, tc.outcome, got.Gates)
					}
					if tc.gate != "" {
						g := got.Gates[tc.gate]
						if g.Reason != tc.reason || g.State != tc.state {
							t.Fatalf("gate %s=%+v want %s/%s", tc.gate, g, tc.state, tc.reason)
						}
					}
				})
			}
		})
	}
}

func TestReviewSelectionAndModelFallback(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		edit   func(*Snapshot)
		status Status
		flag   bool
		reason Reason
	}{
		{"other family", func(*Snapshot) {}, Complete, false, EvidenceSatisfied},
		{"same family fallback", func(s *Snapshot) {
			for i := range s.Rounds {
				s.Rounds[i].OtherModel = false
			}
		}, Complete, true, EvidenceSatisfied},
		{"stale other family", func(s *Snapshot) {
			old := s.Rounds[0]
			old.Current = false
			for i := range s.Rounds {
				s.Rounds[i].OtherModel = false
			}
			s.Rounds = append(s.Rounds, old)
		}, Complete, true, EvidenceSatisfied},
		{"merge stream is shadow", func(s *Snapshot) { s.Rounds[0].Stream = "merge" }, Refused, true, ReviewMissing},
		{"latest current loses profile", func(s *Snapshot) {
			r := s.Rounds[0]
			r.ID = "RND-later"
			r.Profiles = []string{"test-adequacy"}
			s.Rounds = append(s.Rounds, r)
		}, Refused, false, ProfileMissing},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", true)
			tc.edit(&s)
			got := Evaluate(s)
			if got.Status != tc.status || got.NoOtherModelReview != tc.flag || got.Gates[PrimaryGate].Reason != tc.reason {
				t.Fatalf("outcome=%+v", got)
			}
		})
	}
}

func TestEveryGateIsReportedAndEmptyChecksAreConfigured(t *testing.T) {
	t.Parallel()
	s := completeSnapshot("routine", false)
	s.Rounds = nil
	s.Checks = nil
	s.Handoff.RequiredChecks = []string{}
	s.Handoff.ResultIDs = []string{}
	s.Handoff.DefinitionDigests = map[string]string{}
	l := s.Policy.Lanes["routine"]
	l.Checks = []string{}
	s.Policy.Lanes["routine"] = l
	got := Evaluate(s)
	if got.Status != Complete || len(got.Gates) != 15 {
		t.Fatalf("outcome=%+v", got)
	}
	for gate, g := range got.Gates {
		if g.State != Satisfied {
			t.Fatalf("gate %s=%+v", gate, g)
		}
	}
	// Refusal wins over a retryable wait, independently of map iteration order.
	s.GitHub.CI = nil
	s.OpenMajors = []reviewrecord.Finding{{ID: "open"}}
	for range 50 {
		if got := Evaluate(s); got.Status != Refused || got.Gates[CIGate].State != WaitingState {
			t.Fatalf("outcome=%+v", got)
		}
	}
}

func TestAdditionalBindingsAndCI(t *testing.T) {
	t.Parallel()
	cases := []matrixCase{
		{"missing PR", func(s *Snapshot) { s.GitHub.PR = nil }, PRGate, PRUnavailable, Missing, Refused},
		{"hand-off belongs elsewhere", func(s *Snapshot) { s.Handoff.TaskID = "other" }, HandoffGate, HandoffMismatch, RefusedState, Refused},
		{"frozen policy definition mismatch", func(s *Snapshot) { s.Handoff.DefinitionDigests["test"] = "other" }, PolicyGate, PolicyMismatch, RefusedState, Refused},
		{"frozen required set mismatch", func(s *Snapshot) { s.Handoff.RequiredChecks = nil }, PolicyGate, PolicyMismatch, RefusedState, Refused},
		{"wrong result task", func(s *Snapshot) { s.Checks[0].Start.TaskID = "other" }, ChecksGate, ChecksUnsatisfied, Missing, Refused},
		{"separate profile missing", func(s *Snapshot) { s.Rounds[1].Profiles = []string{"code"} }, SeparateGate, ProfileMissing, Missing, Refused},
		{"merge commit absent", func(s *Snapshot) { s.GitHub.PR.MergeCommit = "" }, PRGate, MergeInvalid, RefusedState, Refused},
		{"CI no fresh observation", func(s *Snapshot) { s.GitHub.CI = nil }, CIGate, CIWaiting, WaitingState, Waiting},
		{"CI only skipped", func(s *Snapshot) { s.GitHub.CI.Runs[0].Conclusion = "skipped" }, CIGate, CIWaiting, WaitingState, Waiting},
		{"CI different event", func(s *Snapshot) { s.GitHub.CI.Runs[0].Event = "pull_request" }, CIGate, CIWaiting, WaitingState, Waiting},
		{"CI different head", func(s *Snapshot) { s.GitHub.CI.Runs[0].HeadSHA = "different" }, CIGate, CIWaiting, WaitingState, Waiting},
		{"CI later failure", func(s *Snapshot) {
			r := s.GitHub.CI.Runs[0]
			r.Attempt++
			r.Conclusion = "failure"
			s.GitHub.CI.Runs = append(s.GitHub.CI.Runs, r)
		}, CIGate, CINotGreen, RefusedState, Refused},
		{"CI receipt wrong merge", func(s *Snapshot) { s.CIResults[0].MergeID = "other" }, CIGate, CIUnrecorded, Missing, Refused},
		{"CI receipt wrong attempt", func(s *Snapshot) { s.CIResults[0].Attempt = 1 }, CIGate, CIUnrecorded, Missing, Refused},
		{"CI receipt wrong repository", func(s *Snapshot) { s.CIResults[0].Repository = "other/repo" }, CIGate, CIUnrecorded, Missing, Refused},
		{"CI receipt wrong conclusion", func(s *Snapshot) { s.CIResults[0].Conclusion = "failure" }, CIGate, CIUnrecorded, Missing, Refused},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", true)
			tc.edit(&s)
			got := Evaluate(s)
			if got.Status != tc.outcome || got.Gates[tc.gate].Reason != tc.reason || got.Gates[tc.gate].State != tc.state {
				t.Fatalf("outcome=%+v", got)
			}
		})
	}
}
