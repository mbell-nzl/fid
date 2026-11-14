// Package handoffrecord defines the immutable receipt of a builder handing a
// task to review. It contains the evidence needed by the later done evaluator.
package handoffrecord

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
)

var prPattern = regexp.MustCompile(`^https://github\.com/[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?/[A-Za-z0-9_.-]+/pull/[1-9][0-9]*$`)
var oidPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)

func ValidPR(pr string) bool { return prPattern.MatchString(pr) }

func PRNumber(pr string) string { return pr[strings.LastIndex(pr, "/")+1:] }

type View struct {
	ID             string   `json:"id"`
	PR             string   `json:"pr"`
	Head           string   `json:"head"`
	Base           string   `json:"base"`
	RequiredChecks []string `json:"required_checks"`
	ResultIDs      []string `json:"result_ids"`
	At             string   `json:"at"`
}

type Record struct {
	View
	TaskID            string            `json:"task_id"`
	HoldID            string            `json:"hold_id"`
	Actor             string            `json:"actor"`
	Tree              string            `json:"tree"`
	Lane              string            `json:"lane"`
	PolicySnapshot    int               `json:"policy_snapshot"`
	DefinitionDigests map[string]string `json:"definition_digests"`
	ChangedPaths      []string          `json:"changed_paths"`
}

func (record Record) Validate() error {
	if record.ID == "" || record.TaskID == "" || record.HoldID == "" || record.Actor == "" || !ValidPR(record.PR) || !oidPattern.MatchString(record.Head) || !oidPattern.MatchString(record.Tree) || !oidPattern.MatchString(record.Base) || record.PolicySnapshot < 1 {
		return fmt.Errorf("incomplete hand-off identity, comparison or policy")
	}
	if record.Lane != "routine" && record.Lane != "standard" && record.Lane != "assured" {
		return fmt.Errorf("invalid hand-off lane")
	}
	if _, err := time.Parse(time.RFC3339, record.At); err != nil {
		return fmt.Errorf("hand-off time: %w", err)
	}
	if record.RequiredChecks == nil || record.ResultIDs == nil || record.DefinitionDigests == nil || record.ChangedPaths == nil || len(record.RequiredChecks) != len(record.ResultIDs) || len(record.DefinitionDigests) != len(record.RequiredChecks) || !slices.IsSorted(record.RequiredChecks) {
		return fmt.Errorf("incomplete hand-off evidence")
	}
	for i, name := range record.RequiredChecks {
		if name == "" || i > 0 && name == record.RequiredChecks[i-1] || record.ResultIDs[i] == "" || record.DefinitionDigests[name] == "" {
			return fmt.Errorf("invalid hand-off check binding")
		}
	}
	return nil
}

func (record Record) Fields() map[string]any {
	return map[string]any{"task_id": record.TaskID, "status": "in_review", "handoff_id": record.ID, "pr": record.PR, "head": record.Head, "base": record.Base, "required_checks": record.RequiredChecks, "result_ids": record.ResultIDs}
}
