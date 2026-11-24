package donework

import (
	"context"
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

// TaskOutcome is one ordinary Finish result, without turning a task's refusal
// into a failure of the reconciliation command.
type TaskOutcome struct {
	TaskID     string                   `json:"task_id"`
	Status     string                   `json:"status"`
	Completion *completionrecord.Record `json:"completion,omitempty"`
	Refusal    *Refusal                 `json:"refusal,omitempty"`
}

// Reconcile takes a stable list, then uses exactly the single-task path. Finish
// refreshes/revalidates each subject independently after its GitHub observation.
func Reconcile(ctx context.Context, o Options) ([]TaskOutcome, error) {
	if !policy.Allows(o.Actor, policy.WorkSet) {
		return nil, &persondecision.DeniedError{Actor: o.Actor}
	}
	if _, err := journal.EnsureProjectionFresh(o.WorkspaceRoot, o.MigrationsFS); err != nil {
		return nil, err
	}
	tasks, err := reviewTasks(o.WorkspaceRoot)
	if err != nil {
		return nil, err
	}
	outcomes := make([]TaskOutcome, 0, len(tasks))
	o.IdempotencyKey = ""
	for _, task := range tasks {
		o.TaskID = task
		result, err := Finish(ctx, o)
		outcome := TaskOutcome{TaskID: task, Status: "completed"}
		if err == nil {
			outcome.Completion = &result.Completion
		} else {
			var refused *Refusal
			if !errors.As(err, &refused) {
				return outcomes, fmt.Errorf("reconcile %s: %w", task, err)
			}
			outcome.Status, outcome.Refusal = "refused", refused
			switch refused.Code {
			case "done-not-merged", "done-ci-waiting":
				outcome.Status = "waiting"
			case "github-unavailable":
				outcome.Status = "github-unavailable"
			}
		}
		outcomes = append(outcomes, outcome)
	}
	return outcomes, nil
}

func reviewTasks(root string) ([]string, error) {
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	rows, err := conn.Query(`SELECT id FROM issues WHERE kind='task' AND status='in_review' AND EXISTS(SELECT 1 FROM handoffs WHERE task_id=issues.id) ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	tasks := []string{}
	for rows.Next() {
		var task string
		if err := rows.Scan(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}
