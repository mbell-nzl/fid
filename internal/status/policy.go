package status

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/mbell-nzl/fid/internal/policy"
)

// Warning describes a problem fid could not confirm needs a person.
// details holds only the fields belonging to this warning's kind.
type Warning struct {
	Kind       string
	Summary    string
	NextAction string
	details    map[string]any
}

func (warning Warning) MarshalJSON() ([]byte, error) {
	fields := map[string]any{"kind": warning.Kind, "summary": warning.Summary, "next_action": warning.NextAction}
	for key, value := range warning.details {
		fields[key] = value
	}
	return json.Marshal(fields)
}

// NewWarning shares the closed warning shape with other read commands.
func NewWarning(kind, summary, next string, details map[string]any) Warning {
	return Warning{Kind: kind, Summary: summary, NextAction: next, details: details}
}

func PolicyUnavailableWarning(cause error) Warning {
	var source *policy.SourceError
	if !errors.As(cause, &source) {
		source = &policy.SourceError{Problem: cause.Error()}
	}
	return NewWarning("policy-unavailable", "Policy unavailable: "+source.Problem, source.NextAction(), source.Details())
}

func ChecksNotConfiguredWarning(lanes []string) Warning {
	names, tables := []string{}, []string{}
	for _, lane := range lanes {
		names = append(names, strings.ToUpper(lane[:1])+lane[1:])
		tables = append(tables, "[lanes."+lane+"]")
	}
	return NewWarning("checks-not-configured", "No checks set for "+strings.Join(names, ", ")+" in .fid.toml",
		"add checks to "+joinTables(tables)+", then merge .fid.toml and git fetch", map[string]any{"lanes": lanes})
}

func readPolicyWarnings(root string) []Warning {
	warnings := []Warning{}
	loaded, err := policy.Load(root)
	if err != nil {
		return append(warnings, PolicyUnavailableWarning(err))
	}
	lanes := []string{}
	for _, lane := range []string{"routine", "standard", "assured"} {
		if !loaded.Effective.Lanes[lane].ChecksSet {
			lanes = append(lanes, lane)
		}
	}
	if len(lanes) > 0 {
		warnings = append(warnings, ChecksNotConfiguredWarning(lanes))
	}
	if !loaded.Effective.TestAdequacySet {
		warnings = append(warnings, Warning{Kind: "test-paths-not-configured", Summary: "No test paths set in .fid.toml",
			NextAction: "add [review] test_adequacy, then merge .fid.toml and git fetch"})
	}
	return warnings
}

func joinTables(tables []string) string {
	if len(tables) == 1 {
		return tables[0]
	}
	return strings.Join(tables[:len(tables)-1], ", ") + " and " + tables[len(tables)-1]
}
