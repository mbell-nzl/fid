package status

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/decisionview"
	"github.com/mbell-nzl/fid/internal/fingerprint"
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

func (result *Result) addApproval(task string) {
	approval := result.approvals[task]
	if approval.Decision == nil || approval.Decision.State == "current" {
		return
	}
	summary := task + " approval " + approval.Decision.State
	if flag := approval.NoOtherModelReview; flag != nil && *flag {
		summary += ", no other-model review"
	}
	action := fmt.Sprintf("fid approve %s --actor human:<name> --reason <text>", task)
	result.NeedsYou = append(result.NeedsYou, NeedsYouItem{Kind: "approval", Summary: summary, NextAction: &action, Subject: &Subject{Type: "task", ID: task}})
}
