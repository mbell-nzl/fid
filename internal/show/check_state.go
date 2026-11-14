package show

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/status"
)

type Check struct {
	Name      string               `json:"name"`
	Start     *checkrecord.Start   `json:"start"`
	Process   *checkrecord.Process `json:"process"`
	Result    *checkrecord.Result  `json:"result"`
	Required  *bool                `json:"required"`
	Satisfied *bool                `json:"satisfied"`
}

func (result *Result) readChecks(root string) {
	result.Checks, result.Warnings = []Check{}, []status.Warning{}
	if result.Issue.Kind != "task" {
		return
	}
	loaded, err := policy.Load(root)
	if err != nil {
		result.Warnings = append(result.Warnings, status.PolicyUnavailableWarning(err))
		result.selectChecks(nil, nil, false, checkread.Evaluation{})
		return
	}
	required, configured := result.requiredChecks(root, loaded.Effective)
	tree := ""
	if len(result.IssueView.Checks) > 0 {
		tree, err = policy.HeadTree(root)
		if err != nil {
			result.comparisonWarning(err)
		}
	}
	names := []string{}
	for name, needed := range required {
		if needed != nil && *needed {
			names = append(names, name)
		}
	}
	evaluation := checkread.Evaluate(result.Issue.ID, tree, names, loaded.Effective.Checks, result.IssueView.Checks)
	result.selectChecks(loaded.Effective.Checks, required, configured, evaluation)
}

func (result *Result) requiredChecks(root string, effective policy.Effective) (map[string]*bool, bool) {
	required := map[string]*bool{}
	if result.Issue.Lane == nil {
		return required, false
	}
	lane := *result.Issue.Lane
	configured := effective.Lanes[lane]
	if !configured.ChecksSet {
		result.Warnings = append(result.Warnings, status.ChecksNotConfiguredWarning([]string{lane}))
		return required, false
	}
	conditional := false
	for _, name := range configured.Checks {
		required[name] = boolPointer(true)
		if len(effective.Checks[name].WhenChanged) > 0 {
			required[name], conditional = nil, true
		}
	}
	if !conditional {
		return required, true
	}
	comparison, err := policy.Compare(root, result.Issue.ID)
	if err != nil {
		result.comparisonWarning(err)
		return required, true
	}
	names, err := policy.Required(effective, lane, result.Issue.ID, comparison.Paths)
	if err != nil {
		result.comparisonWarning(err)
		return required, true
	}
	result.ChecksBase = &comparison.Base
	for _, name := range configured.Checks {
		required[name] = boolPointer(slices.Contains(names, name))
	}
	return required, true
}

func (result *Result) selectChecks(definitions map[string]config.Check, required map[string]*bool, configured bool, evaluation checkread.Evaluation) {
	seen := map[string]bool{}
	for _, record := range checkread.Latest(result.IssueView.Checks) {
		name := record.Start.Name
		if chosen, ok := evaluation.Latest[name]; ok && record.Result != nil {
			record = chosen
		}
		entry := Check{Name: name, Start: &record.Start, Process: record.Process, Result: record.Result}
		if _, known := definitions[name]; known {
			entry.Satisfied = boolPointer(evaluation.Satisfied(name))
		}
		entry.Required = required[name]
		if _, listed := required[name]; configured && !listed {
			entry.Required = boolPointer(false)
		}
		result.Checks = append(result.Checks, entry)
		seen[name] = true
	}
	for name, needed := range required {
		if !seen[name] && (needed == nil || *needed) {
			result.Checks = append(result.Checks, Check{Name: name, Required: needed, Satisfied: boolPointer(false)})
		}
	}
	slices.SortStableFunc(result.Checks, func(a, b Check) int { return strings.Compare(a.Name, b.Name) })
}

func (result *Result) comparisonWarning(cause error) {
	for _, warning := range result.Warnings {
		if warning.Kind == "comparison-unavailable" {
			return
		}
	}
	details := map[string]any{"ref": nil, "problem": cause.Error()}
	var refused *checkrecord.Refusal
	if errors.As(cause, &refused) && refused.Code == "comparison-unavailable" {
		details = refused.Details
	}
	result.Warnings = append(result.Warnings, status.NewWarning("comparison-unavailable", "Comparison unavailable: "+fmt.Sprint(details["problem"]),
		"check origin's default branch and fetch enough history to find a common ancestor", details))
}

func boolPointer(value bool) *bool { return &value }
