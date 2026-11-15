// Package completionrecord holds the immutable evidence accepted by fid done.
package completionrecord

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/mergerecord"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type Record struct {
	ID                  string               `json:"id"`
	TaskID              string               `json:"task_id"`
	HandoffID           string               `json:"handoff_id"`
	MergeID             string               `json:"merge_id"`
	ObservedCIResultIDs []string             `json:"observed_ci_result_ids"`
	CIResultIDs         []string             `json:"ci_result_ids"`
	RoundIDs            []string             `json:"round_ids"`
	ApprovalID          string               `json:"approval_id"`
	Fingerprint         reviewrecord.Binding `json:"fingerprint"`
	RollbackPresent     bool                 `json:"rollback_present"`
	NoOtherModelReview  bool                 `json:"no_other_model_review"`
	Actor               string               `json:"actor"`
	At                  string               `json:"at"`
}

type Result struct {
	Completion Record `json:"completion"`
	Status     string `json:"status"`
	Replayed   bool   `json:"replayed"`
	EventID    string `json:"-"`
}

func From(s doneeval.Snapshot, o doneeval.Outcome) Record {
	r := Record{TaskID: s.Task.ID, HandoffID: s.Handoff.ID, MergeID: s.Merge.ID, CIResultIDs: []string{}, ObservedCIResultIDs: []string{}, RoundIDs: []string{}, Fingerprint: reviewrecord.Bind(s.Fingerprint), NoOtherModelReview: o.NoOtherModelReview, RollbackPresent: o.Gates[doneeval.RollbackGate].Reason == doneeval.EvidenceSatisfied}
	for _, gate := range []doneeval.Gate{doneeval.PrimaryGate, doneeval.SeparateGate} {
		r.RoundIDs = append(r.RoundIDs, o.Gates[gate].IDs...)
	}
	if o.Gates[doneeval.ApprovalGate].Reason == doneeval.EvidenceSatisfied && s.Approval.Decision != nil {
		r.ApprovalID = s.Approval.Decision.ID
	}
	for _, run := range s.GitHub.CI.Runs {
		if run.Status != "completed" {
			continue
		}
		for _, c := range s.CIResults {
			if c.RunID == run.ID && c.Attempt == run.Attempt {
				r.ObservedCIResultIDs = append(r.ObservedCIResultIDs, c.ID)
			}
		}
	}
	selected := github.EvaluateCI(s.GitHub.CI.Commit, s.GitHub.CI.Runs)
	for _, run := range selected.Runs {
		if run.Ignored != "" {
			continue
		}
		for _, c := range s.CIResults {
			if c.RunID == run.ID && c.Attempt == run.Attempt {
				r.CIResultIDs = append(r.CIResultIDs, c.ID)
			}
		}
	}
	return r
}

func (r Record) Validate() error {
	if r.ID == "" || r.TaskID == "" || r.HandoffID == "" || r.MergeID == "" || len(r.CIResultIDs) == 0 || len(r.ObservedCIResultIDs) == 0 || r.RoundIDs == nil || !policy.Allows(r.Actor, policy.WorkSet) {
		return fmt.Errorf("incomplete task completion")
	}
	if _, err := time.Parse(time.RFC3339, r.At); err != nil {
		return fmt.Errorf("completion time: %w", err)
	}
	for _, ids := range [][]string{r.CIResultIDs, r.ObservedCIResultIDs, r.RoundIDs} {
		seen := map[string]bool{}
		for _, id := range ids {
			if id == "" || seen[id] {
				return fmt.Errorf("invalid completion evidence IDs")
			}
			seen[id] = true
		}
	}
	if r.Fingerprint.Algorithm != "" || r.Fingerprint.Fingerprint != "" || r.Fingerprint.Base != "" || r.Fingerprint.Head != "" {
		if r.Fingerprint.Algorithm != fingerprint.Algorithm || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(r.Fingerprint.Fingerprint) {
			return fmt.Errorf("invalid completion fingerprint")
		}
	}
	return nil
}

func Read(conn reviewrecord.Queryer, task string) (Record, error) {
	var r Record
	var raw string
	if err := conn.QueryRow(`SELECT record_json FROM completions WHERE task_id=?`, task).Scan(&raw); err != nil {
		return r, err
	}
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		return r, err
	}
	return r, r.Validate()
}

// Admit replays the single evaluator using recorded present-evidence facts.
// It does not reacquire mutable external evidence or rewrite any receipt.
func Admit(conn reviewrecord.Queryer, r Record) error {
	if err := r.Validate(); err != nil {
		return err
	}
	s, err := doneeval.Recorded(conn, r.TaskID, doneeval.Target{}, doneeval.GitHubFacts{}, fingerprint.Result{Algorithm: r.Fingerprint.Algorithm, Fingerprint: r.Fingerprint.Fingerprint, Base: r.Fingerprint.Base, Head: r.Fingerprint.Head})
	if err != nil {
		return err
	}
	if s.Handoff == nil || s.Merge == nil || s.Handoff.ID != r.HandoffID || s.Merge.ID != r.MergeID || s.Merge.HandoffID != r.HandoffID || s.Merge.ObservedAt > r.At {
		return fmt.Errorf("completion must name its latest hand-off and merge")
	}
	if r.Fingerprint.Algorithm != "" && (r.Fingerprint.Base != s.Handoff.Base || r.Fingerprint.Head != s.Handoff.Head) {
		return fmt.Errorf("completion fingerprint differs from latest hand-off endpoints")
	}
	m := s.Merge
	repo, err := github.ParseOrigin("https://github.com/" + m.Repository)
	if err != nil {
		return err
	}
	at, err := time.Parse(time.RFC3339, m.MergedAt)
	if err != nil {
		return err
	}
	pr := github.PullRequest{Number: m.Number, URL: m.PR, State: "MERGED", Merged: true, MergedAt: &at, MergeCommit: m.MergeCommit, HeadSHA: m.HeadSHA, BaseRef: m.BaseBranch, BaseRepository: repo}
	if r.RollbackPresent {
		pr.Body = "## Rollback\n"
	}
	s.Target = doneeval.Target{Repository: repo, DefaultBranch: m.BaseBranch}
	ci := github.CI{Commit: m.MergeCommit, Runs: []github.WorkflowRun{}}
	selected := make([]mergerecord.CIResult, 0)
	for _, id := range r.ObservedCIResultIDs {
		found := false
		for _, c := range s.CIResults {
			if c.ID == id {
				if c.ObservedAt > r.At {
					return fmt.Errorf("completion precedes CI evidence")
				}
				found = true
				ci.Runs = append(ci.Runs, github.WorkflowRun{ID: c.RunID, Attempt: c.Attempt, WorkflowID: c.WorkflowID, Name: c.Name, Path: c.Path, Event: c.Event, HeadSHA: c.HeadSHA, Status: c.Status, Conclusion: c.Conclusion, URL: c.URL})
				selected = append(selected, c)
				break
			}
		}
		if !found {
			return fmt.Errorf("completion CI result belongs to another task or merge")
		}
	}
	s.CIResults = selected
	s.GitHub = doneeval.GitHubFacts{PR: &pr, HeadTree: m.HeadTree, CI: &ci}
	o := doneeval.Evaluate(s)
	if o.Status != doneeval.Complete {
		return fmt.Errorf("completion evidence does not satisfy done: %+v", o.Gates)
	}
	expected := From(s, o)
	if !reflect.DeepEqual(expected.ObservedCIResultIDs, r.ObservedCIResultIDs) || !reflect.DeepEqual(expected.CIResultIDs, r.CIResultIDs) || !reflect.DeepEqual(expected.RoundIDs, r.RoundIDs) || expected.ApprovalID != r.ApprovalID || expected.NoOtherModelReview != r.NoOtherModelReview {
		return fmt.Errorf("completion differs from counting evidence")
	}
	return nil
}
