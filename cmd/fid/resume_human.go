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

func writeResumeHuman(cmd *cobra.Command, resume workflow.Resume) error {
	lines := resumeContextHumanLines(resume)
	for _, recovery := range resume.Recoveries {
		lines = append(lines, recoveryHumanLines(recovery)...)
	}
	lines = append(lines, budgetHumanLines(resume.Budgets)...)
	for _, stop := range resume.Stops {
		lines = append(lines, stopHumanDescription(stop))
	}
	for _, attention := range resume.Attention {
		lines = append(lines, fmt.Sprintf("Attention %s: %s; %s; linked Stops: %s", attention.ID, attention.Status, attention.Code, strings.Join(attention.StopIDs, ", ")))
	}
	for _, obligation := range resume.Obligations {
		lines = append(lines, fmt.Sprintf("Obligation %s: %s; validity %s; freshness %s; Artifacts %s", obligation.ObligationID, obligation.Status, obligation.Validity, obligation.Freshness, obligation.ArtifactAvailability))
		lines = append(lines, evaluationReasonHumanLines(obligation.Reasons)...)
	}
	lines = append(lines, evaluationReasonHumanLines(resume.Completion.Reasons)...)
	if !resume.Findings.Available {
		lines = append(lines, "Findings: unavailable; this view does not establish completed Review support.")
	}
	for _, action := range resume.NextActions {
		state := "ready"
		if len(action.Reasons) != 0 {
			state = "blocked"
		}
		lines = append(lines, fmt.Sprintf("Action %s %s:%s: %s", action.Operation, action.Subject.Type, action.Subject.ID, state))
		if len(action.RequiredCapabilities) != 0 {
			lines = append(lines, "  Requires: "+strings.Join(action.RequiredCapabilities, ", "))
		}
		lines = append(lines, evaluationReasonHumanLines(action.Reasons)...)
	}
	if len(resume.NextActions) == 0 {
		lines = append(lines, "Next actions: none.")
	}
	_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(lines, "\n"))
	return err
}

func resumeContextHumanLines(resume workflow.Resume) []string {
	lines := []string{
		fmt.Sprintf("Issue %s: %s; completion %s", resume.Issue.ID, resume.Issue.Status, resume.Completion.Status),
		fmt.Sprintf("Envelope %s revision %d: %s", resume.Envelope.ID, resume.Envelope.CurrentRevision, resume.Envelope.Status),
		fmt.Sprintf("Repository %s: HEAD %s; tree %s; dirty %t", resume.Repository.WorktreeID, resume.Repository.HeadOID, resume.Repository.TreeOID, resume.Repository.Dirty),
	}
	if resume.Repository.DirtyDigest != nil {
		lines = append(lines, "Dirty target: "+*resume.Repository.DirtyDigest)
	}
	if run := resume.Run; run != nil {
		lines = append(lines, fmt.Sprintf("Run %s: %s; executor %s", run.ID, run.Status, run.Executor.ID))
		if run.PredecessorRunID != nil {
			lines = append(lines, "Prior Run: "+*run.PredecessorRunID)
		}
		if run.SuccessorRunID != nil {
			lines = append(lines, "Successor Run: "+*run.SuccessorRunID)
		}
	}
	if claim := resume.Claim; claim != nil {
		lines = append(lines, fmt.Sprintf("Claim %s: effective %s; recorded %s; expires %s", claim.ID, claim.EffectiveStatus, claim.Status, claim.ExpiresAt))
	} else {
		lines = append(lines, "Current Claim: none.")
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
		fmt.Sprintf("Authority: %s by %s; %s", recovery.Authority.Mode, recovery.Authority.ActorID, recovery.Authority.Reason),
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
	return append(lines, "Adopted paths (unproven work): "+paths)
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
