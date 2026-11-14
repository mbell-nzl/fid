package checkread

import (
	"slices"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/config"
)

// Evaluation is the shared answer for hand-off and done. Required is the
// selected policy's set; Missing names those without a counting pass. Latest
// holds the chosen finished records, and LatestResultIDs names them directly.
type Evaluation struct {
	Required, Missing []string
	LatestResultIDs   map[string]string
	Latest            map[string]checkrecord.Record
	passedBefore      map[string]bool
}

// Evaluate selects by task, starting tree and check-definition digest. Records
// must be in finish journal order, as Read returns them; caller times are never
// ordering evidence. An unfinished start cannot replace a finished result.
// Callers supply the current required set at hand-off or the recorded hand-off
// policy's set at done, together with the tree being evaluated.
func Evaluate(task, tree string, required []string, definitions map[string]config.Check, records []checkrecord.Record) Evaluation {
	result := Evaluation{Required: append([]string{}, required...), Missing: []string{}, LatestResultIDs: map[string]string{}, Latest: map[string]checkrecord.Record{}, passedBefore: map[string]bool{}}
	slices.Sort(result.Required)
	result.Required = slices.Compact(result.Required)
	for _, record := range records {
		start := record.Start
		definition, exists := definitions[start.Name]
		if record.Result == nil || start.TaskID != task || start.Tree.Tree != tree || !exists || start.DefinitionDigest != checkrecord.Digest(checkrecord.DefinitionFor(definition)) {
			continue
		}
		if previous, ok := result.Latest[start.Name]; ok && previous.Result.Counts {
			result.passedBefore[start.Name] = true
		}
		result.Latest[start.Name] = record
		result.LatestResultIDs[start.Name] = start.ID
	}
	for _, name := range result.Required {
		if latest, exists := result.Latest[name]; !exists || !latest.Result.Counts || latest.Result.Outcome != "passed" {
			result.Missing = append(result.Missing, name)
		}
	}
	return result
}

func (result Evaluation) Refusal(task string) *checkrecord.Refusal {
	if len(result.Missing) == 0 {
		return nil
	}
	ids := []string{}
	for _, name := range result.Missing {
		if id := result.LatestResultIDs[name]; id != "" {
			ids = append(ids, id)
		}
	}
	message := "required checks have no counting pass on this tree; check them before hand-off"
	if len(result.Missing) == 1 {
		name := result.Missing[0]
		if latest, ok := result.Latest[name]; ok && latest.Result.Outcome == "failed" && result.passedBefore[name] {
			message = name + " failed on this tree after passing; check it again before hand-off"
		}
	}
	return &checkrecord.Refusal{Code: "checks-unsatisfied", TaskID: task, Message: message, Details: map[string]any{"checks": slices.Clone(result.Missing), "latest_result_ids": ids}, Next: ""}
}

// Satisfied also describes explicitly named, optional checks on this tree.
func (result Evaluation) Satisfied(name string) bool {
	latest, ok := result.Latest[name]
	return ok && latest.Result.Counts && latest.Result.Outcome == "passed"
}

// RequiredName is used by renderers without another required-set evaluator.
func (result Evaluation) RequiredName(name string) bool {
	_, found := slices.BinarySearch(result.Required, name)
	return found
}
