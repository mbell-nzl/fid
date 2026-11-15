package donework

import "github.com/mbell-nzl/fid/internal/doneeval"

var gateOrder = []doneeval.Gate{doneeval.TaskGate, doneeval.HandoffGate, doneeval.LaneGate, doneeval.PolicyGate, doneeval.ChecksGate, doneeval.ArtifactsGate, doneeval.PrimaryGate, doneeval.MajorsGate, doneeval.SeparateGate, doneeval.ApprovalGate, doneeval.RollbackGate, doneeval.PRGate, doneeval.TreeGate, doneeval.MergeGate, doneeval.CIGate}

func outcomeRefusal(o doneeval.Outcome) error {
	for _, waiting := range []bool{false, true} {
		for _, gate := range gateOrder {
			g := o.Gates[gate]
			if g.State == doneeval.Satisfied || (g.State == doneeval.WaitingState) != waiting || g.Reason == doneeval.PrerequisiteMissing {
				continue
			}
			return refusal(gate, string(g.Reason), g.IDs)
		}
	}
	return refusal(doneeval.HandoffGate, string(doneeval.PrerequisiteMissing), nil)
}

type guidance struct{ message, next string }

var repairs = map[string]guidance{
	"task-not-in-review":        {"the task must be in review before completion", "start, check and hand off the task before fid done"},
	"handoff-missing":           {"the task has no hand-off", "hand off the task before fid done"},
	"handoff-mismatch":          {"the hand-off belongs to another task", "hand off again"},
	"lane-changed":              {"the task's lane changed after hand-off", "hand off again at the current lane"},
	"policy-missing":            {"the hand-off's policy is missing", "restore the recorded policy or hand off again"},
	"policy-mismatch":           {"the hand-off's policy and required checks disagree", "hand off again"},
	"checks-not-configured":     {"the hand-off's lane has no checks configured", "configure checks on the default branch and hand off again"},
	"test-paths-not-configured": {"the hand-off has no test paths configured", "configure test paths on the default branch and hand off again"},
	"checks-unsatisfied":        {"required local checks are missing or unsuccessful", "check the change and hand off again"},
	"kept-artifact-unavailable": {"kept Artifacts are missing or changed", "restore their bytes or check and hand off again"},
	"artifacts-not-verified":    {"kept Artifacts have not been verified", "restore their bytes and retry fid done"},
	"review-missing":            {"a required review is missing", "record the required review profiles"},
	"review-stale":              {"the review covers an earlier change", "record a review of the current change"},
	"profile-missing":           {"the review is missing a required profile", "record every required profile in one review session"},
	"separate-session-required": {"Assured work needs a separate review session", "record a separate review in another session"},
	"open-major":                {"the task has an open major finding", "close or resolve the major finding"},
	"approval-missing":          {"the change needs a person's approval", "approve the change with fid approve"},
	"approval-stale":            {"the approval covers an earlier change", "approve the current change with fid approve"},
	"rollback-missing":          {"the PR needs a Rollback section", "add a level-two Rollback heading to the PR body"},
	"pr-unavailable":            {"the PR has not been confirmed", "check the linked PR and retry fid done"},
	"pr-mismatch":               {"the PR differs from the hand-off", "hand off again with the correct PR"},
	"repository-mismatch":       {"the PR must belong to origin's GitHub repository", "check origin and the hand-off's PR"},
	"branch-mismatch":           {"the PR must target origin's default branch", "check origin/HEAD and the PR's base branch"},
	"not-merged":                {"the PR has not been merged", "merge the reviewed PR, then retry fid done"},
	"merge-invalid":             {"the PR has no valid merge commit", "check the PR merge and retry fid done"},
	"head-tree-mismatch":        {"the merged PR head differs from the hand-off's tree", "check the merged change and hand off again"},
	"merge-missing":             {"the merge observation is missing", "retry fid done to record the merge"},
	"merge-mismatch":            {"the merge observation differs from the latest hand-off", "hand off again"},
	"ci-waiting":                {"merged, waiting for CI", "wait for CI to finish, then retry fid done"},
	"ci-not-green":              {"CI on the merge commit is not green", "repair the failing CI check, then retry fid done"},
	"ci-mismatch":               {"CI belongs to another commit", "retry fid done for CI on the merge commit"},
	"ci-unrecorded":             {"the completed CI check has not been recorded", "retry fid done"},
	"fingerprint-unavailable":   {"the hand-off's change fingerprint is unavailable", "restore the recorded base and head objects locally, then retry fid done"},
	"prerequisite-missing":      {"completion evidence is missing", "check the hand-off's required evidence"},
	"observation-stale":         {"the task changed while GitHub was being observed", "inspect the current task, then retry fid done"},
}

func refusal(gate doneeval.Gate, reason string, ids []string) *Refusal {
	if ids == nil {
		ids = []string{}
	}
	g := repairs[reason]
	code := publicCode(reason)
	if gate == doneeval.SeparateGate && (code == "done-review-required" || reason == "separate-session-required") {
		code = "done-separate-review-required"
	}
	return &Refusal{Code: code, Message: g.message, Next: g.next, Details: map[string]any{"gate": string(gate), "reason": reason, "record_ids": ids}}
}

// Public codes group repairs; details retain the evaluator's exact gate/reason.
func publicCode(reason string) string {
	switch reason {
	case "task-not-in-review", "checks-unsatisfied", "kept-artifact-unavailable", "open-major", "approval-missing", "approval-stale", "rollback-missing", "not-merged", "ci-waiting", "ci-not-green", "observation-stale":
		return "done-" + reason
	case "handoff-missing", "handoff-mismatch", "lane-changed", "pr-mismatch", "repository-mismatch", "branch-mismatch", "merge-invalid", "head-tree-mismatch":
		return "done-handoff-required"
	case "review-missing", "review-stale", "profile-missing":
		return "done-review-required"
	case "separate-session-required":
		return "done-separate-review-required"
	default:
		return "done-evidence-required"
	}
}
