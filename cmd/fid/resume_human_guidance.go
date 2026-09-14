package main

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/workflow"
)

func resumeNextStepLines(resume workflow.Resume) []string {
	lines := []string{"", "Next steps"}
	if len(resume.NextActions) == 0 {
		return append(lines, "  Next actions: none.")
	}
	for _, action := range resume.NextActions {
		if action.Operation == "gate.run" && len(action.Reasons) == 0 {
			lines = append(lines, "  To collect proof for unfinished work, run the configured Gate (gate.run).",
				"  A successful execution can supply Evidence; it does not itself complete the Issue.")
		}
	}
	lines = append(lines, "  Candidate operations below are not an execution sequence.",
		"  Ready means the recorded prerequisites permit a request; supplied inputs and authority are checked when it runs.")
	for _, action := range resume.NextActions {
		state := "ready"
		if len(action.Reasons) > 0 {
			state = "blocked"
		}
		lines = append(lines, fmt.Sprintf("  Action %s %s:%s: %s", action.Operation, action.Subject.Type, action.Subject.ID, state),
			"    Purpose: "+resumeActionPurpose(action.Operation))
		if len(action.RequiredCapabilities) > 0 {
			lines = append(lines, "    Requires: "+strings.Join(action.RequiredCapabilities, ", "))
		}
		lines = append(lines, evaluationReasonHumanLines(action.Reasons)...)
	}
	return lines
}

func resumeActionPurpose(operation string) string {
	switch operation {
	case "run.inspect":
		return "Compare this repository with the Run's recorded start or checkpoint; inspect changed paths without recording a checkpoint."
	case "run.observe":
		return "Have the activated harness record a fresh observation of whether this Run's executor is live, terminated or unknown."
	case "gate.run":
		return "Execute the configured check to produce Evidence for the current work."
	case "claim.renew":
		return "Extend an eligible work lease after checking authority."
	case "claim.release":
		return "Give up an eligible work lease after checking authority; this does not prove the work complete."
	default:
		return "Inspect the operation's requirements before submitting a request."
	}
}
