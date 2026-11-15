// Package mergerecord defines immutable merge observations and remote CI check
// results. Neither record implies local execution or task completion.
package mergerecord

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/github"
)

type Merge struct {
	ID          string `json:"id"`
	TaskID      string `json:"task_id"`
	HandoffID   string `json:"handoff_id"`
	PR          string `json:"pr"`
	Number      int64  `json:"number"`
	Repository  string `json:"repository"`
	BaseBranch  string `json:"base_branch"`
	MergeCommit string `json:"merge_commit"`
	HeadSHA     string `json:"head_sha"`
	HeadTree    string `json:"head_tree"`
	MergedAt    string `json:"merged_at"`
	ObservedAt  string `json:"observed_at"`
}

type CIResult struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	MergeID    string `json:"merge_id"`
	Repository string `json:"repository"`
	RunID      int64  `json:"run_id"`
	Attempt    int64  `json:"attempt"`
	WorkflowID int64  `json:"workflow_id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Event      string `json:"event"`
	HeadSHA    string `json:"head_sha"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	URL        string `json:"url"`
	ObservedAt string `json:"observed_at"`
}

var shaPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)

func (m Merge) Validate() error {
	identity, err := github.ParsePRURL(m.PR)
	if err != nil || identity.Number != m.Number || identity.Repository.String() != m.Repository {
		return fmt.Errorf("merge observation PR identity mismatch")
	}
	if m.ID == "" || m.TaskID == "" || m.HandoffID == "" || strings.TrimSpace(m.BaseBranch) == "" || !shaPattern.MatchString(m.MergeCommit) || !shaPattern.MatchString(m.HeadSHA) || !shaPattern.MatchString(m.HeadTree) {
		return fmt.Errorf("incomplete merge observation identity or commit")
	}
	return validateTimes(m.MergedAt, m.ObservedAt)
}

func (c CIResult) Validate() error {
	if c.ID == "" || c.TaskID == "" || c.MergeID == "" || c.RunID < 1 || c.Attempt < 1 || c.WorkflowID < 1 || strings.TrimSpace(c.Name) == "" || strings.TrimSpace(c.Path) == "" || strings.TrimSpace(c.Event) == "" || !shaPattern.MatchString(c.HeadSHA) {
		return fmt.Errorf("incomplete CI result identity or commit")
	}
	identity, err := github.ParsePRURL("https://github.com/" + c.Repository + "/pull/1")
	if err != nil || identity.Repository.String() != c.Repository || !strings.EqualFold(c.URL, "https://github.com/"+c.Repository+"/actions/runs/"+strconv.FormatInt(c.RunID, 10)) {
		return fmt.Errorf("CI result repository or URL mismatch")
	}
	if c.Status != "completed" {
		return fmt.Errorf("CI result requires completed status")
	}
	switch c.Conclusion {
	case "success", "failure", "neutral", "cancelled", "skipped", "timed_out", "action_required", "stale", "startup_failure":
	default:
		return fmt.Errorf("invalid CI conclusion")
	}
	return validateTimes(c.ObservedAt)
}

func validateTimes(values ...string) error {
	for _, value := range values {
		stamp, err := time.Parse(time.RFC3339, value)
		if err != nil || stamp.IsZero() || !strings.HasSuffix(value, "Z") {
			return fmt.Errorf("invalid observation time %q", value)
		}
	}
	return nil
}
