package main

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func writeRunInspectHuman(cmd *cobra.Command, inspection workflow.RunInspection, context workflow.RecoveryReadContext) error {
	run, comparison := inspection.Run, inspection.Reconciliation
	lines := []string{"Run " + run.ID, "  Status: " + run.Status, "  Issue: " + run.IssueID,
		"  Actor ID: " + run.ActorID, "  Executor ID: " + run.Executor.ID}
	lines = append(lines, recoveryLocalContextLines(context, run.Status == "active")...)
	lines = append(lines, "", "Reconciliation (repository comparison)",
		"  Recorded repository: "+repositoryHumanSummary(run.StartRepository.HeadOID, run.StartRepository.Dirty),
		"  Current repository: "+repositoryHumanSummary(comparison.Repository.HeadOID, comparison.Repository.Dirty))
	relation, err := repositoryHumanRelation(cmd, run.StartRepository, comparison.Repository)
	if err != nil {
		return err
	}
	lines = append(lines, "  Worktree comparison: "+relation)
	base := "Run start (no accepted checkpoint)"
	if comparison.CheckpointID != nil {
		base = "checkpoint " + *comparison.CheckpointID
	}
	lines = append(lines, "  Comparison base: "+base, "  Ancestry: "+comparison.Ancestry,
		"  Repository operation: "+comparison.Operation, "  Changed paths (not proof of correctness):")
	if len(comparison.Paths) == 0 {
		lines = append(lines, "    None.")
	}
	for _, change := range comparison.Paths {
		lines = append(lines, "    "+change.Path+" — "+change.Classification)
	}
	if len(comparison.UnmergedPaths) > 0 {
		lines = append(lines, "  Unmerged paths: "+strings.Join(comparison.UnmergedPaths, ", "))
	}
	lines = append(lines, "  Reconciliation digest: "+comparison.Digest,
		"  This digest identifies this comparison; it is not test Evidence or executor liveness.",
		"  This inspection does not record a checkpoint or a liveness observation.")
	lines = append(lines, livenessHumanLines(context)...)
	_, err = fmt.Fprintln(cmd.OutOrStdout(), strings.Join(lines, "\n"))
	return err
}

func recoveryLocalContextLines(context workflow.RecoveryReadContext, activeRun bool) []string {
	if context.Import == nil {
		return nil
	}
	history := context.Import
	lines := []string{"", "Imported project: " + history.ProjectID,
		fmt.Sprintf("  History: %d imported Events preserved with their original identities and order.", history.EventCount),
		"  Source archive digest: " + history.ArchiveDigest, "  Destination: " + history.Destination,
		"  Import preserves work history; it does not grant mutation authority.",
		"  Completion Evidence is a separate requirement for finishing the Issue."}
	if !activeRun {
		return lines
	}
	lines = append(lines, "  Before continuing work here: inspect this destination (run.inspect); recovery still checks authority and repository state.")
	if context.Liveness == nil || context.LivenessHistorical || context.Liveness.State == "unknown" {
		lines = append(lines, "  Have the activated harness establish fresh destination liveness (run.observe) before using an observation for recovery.")
	} else {
		lines = append(lines, "  A destination liveness observation is recorded below; check its time and applicability before using it.")
	}
	return lines
}

func livenessHumanLines(context workflow.RecoveryReadContext) []string {
	lines := []string{"", "Liveness (executor observation)"}
	if context.Liveness == nil {
		return append(lines, "  No observation recorded for this Run; whether its executor is live is unknown.")
	}
	label := "Recorded liveness observation"
	if context.LivenessHistorical {
		label = "Historical liveness — imported history, not fresh destination evidence"
	} else if context.Import != nil {
		label = "Destination liveness observation"
	}
	observation := context.Liveness
	return append(lines, "  "+label+": "+observation.State,
		"  Observed at: "+observation.ObservedAt, "  Source reference: "+observation.SourceRef,
		"  An observation describes that time; this read does not establish that the executor is still in that state.")
}

func repositoryHumanSummary(head string, dirty bool) string {
	if len(head) > 12 {
		head = head[:12]
	}
	if head == "" {
		head = "unknown"
	}
	state := "clean"
	if dirty {
		state = "dirty"
	}
	return "HEAD " + head + "; " + state
}

func repositoryHumanRelation(cmd *cobra.Command, recorded, current workflow.Repository) (string, error) {
	root := ""
	if recorded.ProjectID != "" {
		var err error
		root, err = issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return "", err
		}
	}
	if workflow.SameWorktree(root, recorded, current) {
		return "same worktree", nil
	}
	return "different worktree", nil
}
