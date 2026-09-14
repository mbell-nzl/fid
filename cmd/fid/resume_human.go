package main

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func writeRunRecoveryHuman(cmd *cobra.Command, result workflow.RunRecoveryResult) error {
	lines := recoveryHumanLines(result.Recovery)
	lines = append(lines, budgetHumanLines(result.Run.Budgets)...)
	_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(lines, "\n"))
	return err
}

func writeResumeHuman(cmd *cobra.Command, resume workflow.Resume, context workflow.RecoveryReadContext) error {
	lines := resumeContextHumanLines(resume)
	lines = append(lines, recoveryLocalContextLines(context, resume.Run != nil && resume.Run.Status == "active")...)
	if context.Import != nil {
		lines = append(lines, livenessHumanLines(context)...)
	}
	lines = append(lines, resumeNextStepLines(resume)...)
	if len(resume.Recoveries) > 0 {
		lines = append(lines, "", "Retained work and recovery history")
		for _, recovery := range resume.Recoveries {
			lines = append(lines, recoveryHumanLines(recovery)...)
		}
	}
	if len(resume.Stops)+len(resume.Attention) > 0 {
		lines = append(lines, "", "Blockers and attention")
		for _, stop := range resume.Stops {
			lines = append(lines, stopHumanDescription(stop))
		}
		for _, attention := range resume.Attention {
			lines = append(lines, fmt.Sprintf("Attention %s: %s; %s; linked Stops: %s", attention.ID, attention.Status, attention.Code, strings.Join(attention.StopIDs, ", ")))
		}
	}
	lines = append(lines, resumeCompletionHumanLines(resume)...)
	if len(resume.Budgets) > 0 {
		lines = append(lines, "", "Budgets")
		lines = append(lines, budgetHumanLines(resume.Budgets)...)
	}
	_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(lines, "\n"))
	return err
}

func resumeContextHumanLines(resume workflow.Resume) []string {
	lines := []string{
		fmt.Sprintf("Issue %s", resume.Issue.ID),
		"  Status: " + resume.Issue.Status,
		"  Completion: " + resume.Completion.Status,
		fmt.Sprintf("  Envelope: %s (revision %d, %s)", resume.Envelope.ID, resume.Envelope.CurrentRevision, resume.Envelope.Status),
	}
	if run := resume.Run; run != nil {
		lines = append(lines, "", "Current Run: "+run.ID, "  Status: "+run.Status,
			"  Actor ID: "+run.ActorID, "  Executor ID: "+run.Executor.ID)
		if run.PredecessorRunID != nil {
			lines = append(lines, "  Prior Run: "+*run.PredecessorRunID)
		}
		if run.SuccessorRunID != nil {
			lines = append(lines, "  Successor Run: "+*run.SuccessorRunID)
		}
	}
	if claim := resume.Claim; claim != nil {
		lines = append(lines, "", "Claim "+claim.ID,
			"  Lease status now: "+claim.EffectiveStatus,
			"  Recorded status: "+claim.Status, "  Expires: "+claim.ExpiresAt,
			"  An active lease alone does not authorize an operation; its authority and repository checks still apply.")
	} else {
		lines = append(lines, "", "Current Claim: none.")
	}
	lines = append(lines, "", "Repository",
		"  Project: "+resume.Repository.ProjectID,
		"  Current worktree: "+resume.Repository.WorktreeID,
		"  HEAD: "+resume.Repository.HeadOID,
		fmt.Sprintf("  Uncommitted changes: %t", resume.Repository.Dirty))
	if resume.Run != nil {
		lines = append(lines, "  Run's recorded worktree: "+resume.Run.StartRepository.WorktreeID)
	}
	return lines
}

func resumeCompletionHumanLines(resume workflow.Resume) []string {
	lines := []string{"", "Completion requirements", "  These describe what is needed to finish the Issue, separately from recovery checks above."}
	for _, obligation := range resume.Obligations {
		lines = append(lines, fmt.Sprintf("  Obligation %s: %s", obligation.ObligationID, obligation.Status),
			fmt.Sprintf("    Evidence validity: %s; freshness: %s; Artifacts: %s", obligation.Validity, obligation.Freshness, obligation.ArtifactAvailability))
		lines = append(lines, evaluationReasonHumanLines(obligation.Reasons)...)
	}
	lines = append(lines, "  Overall completion: "+resume.Completion.Status)
	lines = append(lines, evaluationReasonHumanLines(resume.Completion.Reasons)...)
	if !resume.Findings.Available {
		lines = append(lines, "  Findings: unavailable; this view does not establish completed Review support.")
	}
	return lines
}

func evaluationReasonHumanLines(reasons []workflow.EvaluationReason) []string {
	lines := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		lines = append(lines, fmt.Sprintf("  %s for %s:%s: %s", reason.Code, reason.SubjectType, reason.SubjectID, reason.Detail))
	}
	return lines
}

func recoveryHumanLines(recovery workflow.Recovery) []string {
	lines := []string{
		fmt.Sprintf("Recovery %s: Run %s -> %s", recovery.ID, recovery.PriorRunID, recovery.SuccessorRunID),
		fmt.Sprintf("Claims: %s -> %s", recovery.PriorClaimID, recovery.SuccessorClaimID),
		fmt.Sprintf("Authority: %s; Actor ID: %s; reason: %s", recovery.Authority.Mode, recovery.Authority.ActorID, recovery.Authority.Reason),
	}
	if recovery.Authority.Reference != nil {
		lines = append(lines, "Authority reference: "+*recovery.Authority.Reference)
	}
	if grant := recovery.Authority.GrantRef; grant != nil {
		lines = append(lines, fmt.Sprintf("Authority grant: %s revision %d ordinal %d", grant.EnvelopeID, grant.Revision, grant.Ordinal))
	}
	if recovery.Authority.Mode == "force" {
		lines = append(lines, "Force recovery records a known-live takeover; residual filesystem risk remains.")
	}
	paths := "none"
	if len(recovery.AdoptedPaths) != 0 {
		paths = strings.Join(recovery.AdoptedPaths, ", ")
	}
	return append(lines, "Adopted paths (unproven work): "+paths,
		"  Adoption retains existing uncommitted changes for the successor Run.",
		"  Adoption does not establish that these changes were tested, reviewed, approved or proved correct.",
		"  Missing Evidence alone does not tell us whether a Gate was interrupted.")
}

func budgetHumanLines(budgets []workflow.Budget) []string {
	lines := make([]string, 0, len(budgets))
	for _, budget := range budgets {
		remaining, consumed := "unknown", "unknown"
		if budget.Remaining != nil {
			remaining = fmt.Sprint(*budget.Remaining)
		}
		if budget.Consumed != nil {
			consumed = fmt.Sprint(*budget.Consumed)
		}
		line := fmt.Sprintf("Budget %s: remaining %s %s; consumed %s; maximum %d", budget.Name, remaining, budget.Unit, consumed, budget.Maximum)
		if budget.State != "" {
			line += "; state " + budget.State
		}
		lines = append(lines, line)
	}
	return lines
}
