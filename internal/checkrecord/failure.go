package checkrecord

import (
	"fmt"
	"strings"
)

func FailureFor(response Response) *Refusal {
	ids, outcomes, names, reasons := []string{}, []string{}, []string{}, []string{}
	for _, check := range response.Checks {
		if check.Outcome != "passed" {
			ids = append(ids, check.ID)
			outcomes = append(outcomes, check.Outcome)
			names = append(names, check.Name)
		}
	}
	if len(ids) > 0 {
		return &Refusal{Code: "check-unsuccessful", TaskID: response.TaskID, Message: fmt.Sprintf("%d of %d checks failed: %s", len(ids), len(response.Checks), strings.Join(names, ", ")), Details: map[string]any{"check_ids": ids, "outcomes": outcomes}}
	}
	for _, check := range response.Checks {
		if !check.Counts {
			ids = append(ids, check.ID)
			reasons = append(reasons, check.Reasons...)
		}
	}
	if len(ids) > 0 {
		return &Refusal{Code: "check-not-qualified", TaskID: response.TaskID, Message: "checks passed, but some don't count for hand-off", Details: map[string]any{"check_ids": ids, "reasons": reasons}}
	}
	return nil
}
