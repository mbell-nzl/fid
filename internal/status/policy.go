package status

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/mbell-nzl/fid/internal/policy"
)

// Warning describes policy health without turning it into a NEEDS YOU item.
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

func readPolicyWarnings(root string) []Warning {
	warnings := []Warning{}
	loaded, err := policy.Load(root)
	if err != nil {
		var source *policy.SourceError
		if !errors.As(err, &source) {
			source = &policy.SourceError{Problem: err.Error()}
		}
		return append(warnings, Warning{Kind: "policy-unavailable", Summary: "Policy unavailable: " + source.Problem,
			NextAction: source.NextAction(), details: source.Details()})
	}
	lanes, names, tables := []string{}, []string{}, []string{}
	for _, lane := range []struct{ key, name string }{{"routine", "Routine"}, {"standard", "Standard"}, {"assured", "Assured"}} {
		if !loaded.Effective.Lanes[lane.key].ChecksSet {
			lanes = append(lanes, lane.key)
			names = append(names, lane.name)
			tables = append(tables, "[lanes."+lane.key+"]")
		}
	}
	if len(lanes) > 0 {
		warnings = append(warnings, Warning{Kind: "checks-not-configured",
			Summary:    "No checks set for " + strings.Join(names, ", ") + " in .fid.toml",
			NextAction: "add checks to " + joinTables(tables) + ", then merge .fid.toml and git fetch", details: map[string]any{"lanes": lanes}})
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
