package show

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/decisionview"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/mbell-nzl/fid/internal/status"
)

func (result *Result) readApproval(root string) error {
	if result.doneTask() || result.Issue.Kind != "task" || result.Handoff == nil {
		return nil
	}
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
	h, err := reviewrecord.LatestHandoff(tx, result.Issue.ID)
	if err != nil {
		return err
	}
	result.Approval, err = decisionview.ReadApproval(root, tx, h)
	var failure *fingerprint.Failure
	if errors.As(err, &failure) {
		result.Warnings = append(result.Warnings, status.ApprovalUnavailableWarning(h.TaskID, err))
		return nil
	}
	return err
}

func (text textWriter) approval(result Result) []string {
	lines := []string{}
	if decision := result.Approval.Decision; decision != nil {
		lines = text.value("Approval", decision.Text("approved an earlier change"))
		for i := range lines {
			prefix := valueIndent
			if i == 0 {
				prefix = text.label("Approval")
			}
			lines[i] = prefix + text.renderer.Colorize(decision.Role(), lines[i][len(prefix):])
		}
	}
	if flag := result.Approval.NoOtherModelReview; flag != nil && *flag {
		lines = append(lines, valueIndent+text.renderer.Colorize(display.RoleWarning, "no other-model review"))
	}
	return lines
}
