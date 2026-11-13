// Package show assembles fid show: one issue with its epic and task, what
// blocks it and what it blocks, its criteria and its children. Text and
// JSON are both rendered from one Result.
package show

import (
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/project"
)

type Options struct {
	WorkspaceRoot string
	MigrationsDir string
	ID            string
	Now           time.Time
}

// Result is everything fid show shows. Later tasks add the lane, checks,
// reviews, approval, what needs a person, and what's next.
type Result struct {
	issues.IssueView
	Working *holds.Working
	Kept    *holds.Kept
	Now     time.Time
}

// Build reads the issue whose ID is exactly options.ID.
func Build(options Options) (Result, error) {
	view, err := issues.ReadIssueView(issues.GetOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, ID: options.ID})
	if err != nil {
		return Result{}, fmt.Errorf("read show: %w", err)
	}
	result := Result{IssueView: view, Now: options.Now}
	if hold := view.OpenHold; hold != nil && (hold.Kept || !hold.Lapsed(options.Now) || view.Issue.StatusCategory == "started") {
		result.Working = holds.Describe(*hold, project.WorktreeID(options.WorkspaceRoot), options.Now)
		if hold.Kept || view.Issue.StatusCategory == "started" {
			result.Kept = holds.DescribeKept(*hold, project.WorktreeID(options.WorkspaceRoot), options.Now)
		}
	}
	return result, nil
}

// Fields is the result as the issue.show JSON result. Its blocks links
// include those whose other end is completed or cancelled.
func (result Result) Fields() map[string]any {
	return map[string]any{
		"questions":  result.Questions,
		"issue":      result.Issue,
		"state":      result.State,
		"epic":       result.Epic,
		"task":       result.Task,
		"blocked_by": result.BlockedBy,
		"blocks":     result.Blocks,
		"progress":   result.Progress,
		"children":   result.Children,
		"working":    result.Working,
		"kept":       result.Kept,
		"scheduling": result.Scheduling,
	}
}
