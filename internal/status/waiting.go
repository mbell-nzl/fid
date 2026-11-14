package status

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/next"
	"github.com/mbell-nzl/fid/internal/project"
)

func (result *Result) addWaiting(planned next.Result) {
	for _, q := range planned.Questions {
		action := fmt.Sprintf(`fid answer %s "<answer>" --actor human:<name>`, q.TaskID)
		result.NeedsYou = append(result.NeedsYou, NeedsYouItem{Kind: "waiting-for-you", Subject: &Subject{Type: "task", ID: q.TaskID}, Summary: q.TaskID + " is waiting for you: " + q.Text, NextAction: &action})
	}
}

func (result *Result) checkWorktrees(root string, planned next.Result) {
	for _, entry := range planned.Queue {
		if entry.Kept == nil {
			continue
		}
		worktree := entry.Kept.Worktree
		present, err := project.CheckKeptWorktree(root, worktree)
		if err != nil {
			result.Warnings = append(result.Warnings, Warning{Kind: "worktree-unverified", Summary: fmt.Sprintf("Couldn't check the worktree kept for %s: %s", entry.ID, worktree), NextAction: "look at " + worktree + ", then fid status", details: map[string]any{"task_id": entry.ID, "worktree": worktree}})
		} else if !present {
			action := fmt.Sprintf(`fid takeover %s --actor human:<name> --reason "<why>"`, entry.ID)
			result.NeedsYou = append(result.NeedsYou, NeedsYouItem{Kind: "worktree-missing", Subject: &Subject{Type: "task", ID: entry.ID}, Summary: fmt.Sprintf("%s is kept for a worktree that no longer exists: %s", entry.ID, worktree), NextAction: &action})
		}
	}
}
