package show

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/decisionview"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/mergerecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

// Completion expands only the immutable evidence named by fid done. Missing
// rows stay explicit, so an older receipt can still explain what is available.
type Completion struct {
	*completionrecord.Record
	Handoff   *handoffrecord.Record  `json:"handoff"`
	Merge     *mergerecord.Merge     `json:"merge_observation"`
	CIResults []mergerecord.CIResult `json:"ci_results"`
	Rounds    []reviewrecord.Round   `json:"rounds"`
	Approval  *approvalrecord.Record `json:"approval"`
	Waiver    *decisionview.Waiver   `json:"waiver,omitempty"`
	Missing   []string               `json:"missing"`

	requiredChecks []string
	resultIDs      map[string]string
}

func (result Result) doneTask() bool {
	return result.Issue.Kind == "task" && result.Issue.Status == "done"
}

func (result *Result) readCompletion(root string) error {
	if !result.doneTask() {
		return nil
	}
	c := &Completion{CIResults: []mergerecord.CIResult{}, Rounds: []reviewrecord.Round{}, Missing: []string{}}
	result.Completion = c
	result.Handoff = nil
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	record, err := completionrecord.Read(tx, result.Issue.ID)
	if errors.Is(err, sql.ErrNoRows) {
		c.missing("completion record", "")
		return nil
	}
	if err != nil {
		return err
	}
	c.Record = &record
	if c.Fingerprint.Fingerprint == "" {
		c.missing("fingerprint", "")
	}
	c.Handoff, err = completionRow[handoffrecord.Record](tx, c, "hand-off", `SELECT record_json FROM handoffs WHERE id=?`, c.HandoffID)
	if err != nil {
		return err
	}
	if err := result.readCompletionWaiver(tx); err != nil {
		return err
	}
	c.Merge, err = completionRow[mergerecord.Merge](tx, c, "merge", `SELECT record_json FROM merge_observations WHERE id=?`, c.MergeID)
	if err != nil {
		return err
	}
	for _, id := range c.CIResultIDs {
		ci, err := completionRow[mergerecord.CIResult](tx, c, "CI result", `SELECT record_json FROM ci_results WHERE id=?`, id)
		if err != nil {
			return err
		}
		if ci != nil {
			c.CIResults = append(c.CIResults, *ci)
		}
	}
	if err := c.readReviews(tx); err != nil {
		return err
	}
	if h := c.Handoff; h != nil {
		result.Handoff, result.Floor = &h.View, h.Floor
	}
	if c.assured() {
		result.Approval.NoOtherModelReview = &c.NoOtherModelReview
		if a := c.Approval; a != nil {
			result.Approval.Decision = &decisionview.View{State: "current", DecisionID: a.ID, Actor: a.Actor, At: a.At, Fingerprint: a.Binding.Fingerprint}
		}
	}
	return nil
}

func completionRow[T any](conn reviewrecord.Queryer, c *Completion, name, query, id string) (*T, error) {
	var raw string
	if err := conn.QueryRow(query, id).Scan(&raw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.missing(name, id)
			return nil, nil
		}
		return nil, fmt.Errorf("read completed %s: %w", name, err)
	}
	var value T
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return nil, fmt.Errorf("decode completed %s: %w", name, err)
	}
	return &value, nil
}

func (c *Completion) readReviews(conn reviewrecord.Queryer) error {
	rounds, err := reviewrecord.Read(conn, c.TaskID)
	if err != nil {
		return err
	}
	for _, id := range c.RoundIDs {
		found := false
		for _, r := range rounds {
			if r.ID == id {
				c.Rounds = append(c.Rounds, r)
				found = true
				break
			}
		}
		if !found {
			c.missing("review round", id)
		}
	}
	approvals, err := approvalrecord.Read(conn, c.TaskID)
	if err != nil {
		return err
	}
	for _, a := range approvals {
		if a.ID == c.ApprovalID {
			c.Approval = &a
			break
		}
	}
	if c.Approval == nil && (c.ApprovalID != "" || c.assured()) {
		c.missing("approval", c.ApprovalID)
	}
	return nil
}

func (c *Completion) assured() bool {
	if c.Waiver != nil {
		return false
	}
	return c.Handoff != nil && c.Handoff.Lane == "assured" || c.Record != nil && (c.ApprovalID != "" || c.NoOtherModelReview)
}

func (result *Result) completedChecks() {
	c := result.Completion
	if c == nil || c.Handoff == nil {
		return
	}
	for _, name := range c.requiredChecks {
		entry := Check{Name: name, Required: boolPointer(true), Satisfied: boolPointer(false)}
		id := c.resultIDs[name]
		for _, r := range result.IssueView.Checks {
			if r.Start.ID == id {
				entry.Start, entry.Process, entry.Result = &r.Start, r.Process, r.Result
				entry.Satisfied = boolPointer(r.Result != nil && r.Result.Counts && r.Result.Outcome == "passed")
				break
			}
		}
		if entry.Result == nil {
			c.missing("check result", id)
		}
		result.Checks = append(result.Checks, entry)
	}
}

func (c *Completion) missing(kind, id string) {
	message := kind + " missing"
	if strings.TrimSpace(id) != "" {
		message += " · " + id
	}
	c.Missing = append(c.Missing, message)
}
