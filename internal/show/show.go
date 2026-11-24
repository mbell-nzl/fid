// Package show assembles fid show: one issue with its epic and task, what
// blocks it and what it blocks, its criteria and its children. Text and
// JSON are both rendered from one Result.
package show

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/mbell-nzl/fid/internal/decisionview"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/status"
)

type Options struct {
	WorkspaceRoot string
	MigrationsFS  fs.FS
	ID            string
	Now           time.Time
}

// Result is everything fid show shows. Later tasks add reviews, what needs a
// person, and what's next.
type Result struct {
	issues.IssueView
	Working    *holds.Working
	Kept       *holds.Kept
	Now        time.Time
	Checks     []Check
	Floor      *policy.Floor
	ChecksBase *string
	Warnings   []status.Warning
	Approval   decisionview.Approval
	Completion *Completion
}

// Build reads the issue whose ID is exactly options.ID.
func Build(options Options) (Result, error) {
	view, err := issues.ReadIssueView(issues.GetOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsFS: options.MigrationsFS, ID: options.ID})
	if err != nil {
		return Result{}, fmt.Errorf("read show: %w", err)
	}
	result := Result{IssueView: view, Now: options.Now}
	if err := result.readCompletion(options.WorkspaceRoot); err != nil {
		return Result{}, fmt.Errorf("read show completion: %w", err)
	}
	result.readChecks(options.WorkspaceRoot)
	if err := result.readApproval(options.WorkspaceRoot); err != nil {
		return Result{}, fmt.Errorf("read show approval: %w", err)
	}
	if hold := view.OpenHold; hold != nil && (hold.Kept || hold.Waiting() || !hold.Lapsed(options.Now) || view.Issue.StatusCategory == "started") {
		result.Working = holds.Describe(*hold, project.WorktreeID(options.WorkspaceRoot), options.Now)
		if hold.Kept || hold.Waiting() || view.Issue.StatusCategory == "started" {
			result.Kept = holds.DescribeKept(*hold, project.WorktreeID(options.WorkspaceRoot), options.Now)
		}
	}
	return result, nil
}

// Fields is the result as the issue.show JSON result. Its blocks links
// include those whose other end is completed or cancelled.
func (result Result) Fields() map[string]any {
	fields := map[string]any{
		"completion":            result.Completion,
		"approval":              result.Approval.Decision,
		"no_other_model_review": result.Approval.NoOtherModelReview,
		"floor":                 result.Floor,
		"handoff":               result.Handoff,
		"questions":             result.Questions,
		"checks":                result.Checks,
		"checks_base":           result.ChecksBase,
		"warnings":              result.Warnings,
		"issue":                 result.Issue,
		"state":                 result.State,
		"epic":                  result.Epic,
		"task":                  result.Task,
		"blocked_by":            result.BlockedBy,
		"blocks":                result.Blocks,
		"progress":              result.Progress,
		"children":              result.Children,
		"working":               result.Working,
		"kept":                  result.Kept,
		"scheduling":            result.Scheduling,
	}
	if result.Approval.Waivers != nil {
		fields["waivers"] = result.Approval.Waivers
	}
	if result.Merge != nil {
		fields["merge"] = result.Merge
	}
	return fields
}
