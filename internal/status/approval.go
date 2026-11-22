package status

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/decisionview"
	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func ApprovalUnavailableWarning(task string, cause error) Warning {
	return NewWarning("approval-unavailable", task+" approval unavailable: "+cause.Error(),
		"restore the hand-off's recorded base and head objects locally, then fid show "+task,
		map[string]any{"task_id": task, "problem": cause.Error()})
}

func (result *Result) readApproval(root string, conn reviewrecord.Queryer, task string) error {
	h, err := reviewrecord.LatestHandoff(conn, task)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	approval, err := decisionview.ReadApproval(root, conn, h)
	var failure *fingerprint.Failure
	if errors.As(err, &failure) {
		result.Warnings = append(result.Warnings, ApprovalUnavailableWarning(task, err))
	} else if err != nil {
		return err
	}
	result.approvals[task] = approval
	return nil
}

func (result *Result) addReviewAction(root string, conn reviewrecord.Queryer, task string) error {
	s, err := doneeval.GatherReviews(root, conn, task, nil)
	var failure *fingerprint.Failure
	if errors.As(err, &failure) {
		if s.Handoff != nil && s.Handoff.Lane != "assured" {
			result.Warnings = append(result.Warnings, NewWarning("review-unavailable", task+" review unavailable: "+err.Error(), "restore the hand-off's recorded base and head objects locally, then fid show "+task, map[string]any{"task_id": task, "problem": err.Error()}))
		} else if s.Handoff != nil {
			// readApproval already warns when existing evidence needs the
			// fingerprint. The first required round needs that warning too.
			for _, warning := range result.Warnings {
				if warning.Kind == "approval-unavailable" && warning.details["task_id"] == task {
					return nil
				}
			}
			result.Warnings = append(result.Warnings, ApprovalUnavailableWarning(task, err))
		}
		return nil
	}
	if err != nil {
		return err
	}
	if !doneeval.ReviewsDone(s) || result.merges[task] != nil {
		return nil
	}
	kind, summary, action := "merge", task+": merge PR #"+handoffrecord.PRNumber(s.Handoff.PR), s.Handoff.PR
	if s.Handoff.Lane == "assured" && s.Approval.State != "current" {
		kind, summary = "approval", task+": approve, then merge PR #"+handoffrecord.PRNumber(s.Handoff.PR)
		if reviewrecord.NoOtherModelReview(task, s.Rounds) {
			summary += ", no other-model review"
		}
		action = fmt.Sprintf("fid approve %s --actor human:<name> --reason <text>", task)
	}
	result.NeedsYou = append(result.NeedsYou, NeedsYouItem{Kind: kind, Summary: summary, NextAction: &action, Subject: &Subject{Type: "task", ID: task}})
	return nil
}
