package policy

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/pathglob"
)

// ValidateRecorded checks the recorded schema and value domains. It does not
// resolve defaults or compare built-in rule values with the running binary:
// historical snapshots must retain the rules that were actually in force.
func (effective Effective) ValidateRecorded() error {
	if effective.RulesVersion < 1 || effective.HardAssured == nil || effective.ActorRules == nil || effective.ActionSets == nil || effective.Checks == nil || len(effective.Lanes) != 3 || effective.TestAdequacy == nil {
		return fmt.Errorf("policy snapshot effective policy is incomplete")
	}
	if err := effective.validateRecordedTrust(); err != nil {
		return err
	}
	if err := recordedGlobs("hard_assured", effective.HardAssured); err != nil {
		return err
	}
	if !effective.TestAdequacySet && len(effective.TestAdequacy) != 0 {
		return fmt.Errorf("policy snapshot has test_adequacy paths that are not set up")
	}
	if err := recordedGlobs("test_adequacy", effective.TestAdequacy); err != nil {
		return err
	}
	for name, check := range effective.Checks {
		if len(check.Run) == 0 || check.Run[0] == "" || check.Timeout <= 0 {
			return fmt.Errorf("policy snapshot checks.%s requires non-empty argv with an executable and a positive timeout", name)
		}
		if err := recordedGlobs("checks."+name+".defined_by", check.DefinedBy); err != nil {
			return err
		}
		if err := recordedGlobs("checks."+name+".when_changed", check.WhenChanged); err != nil {
			return err
		}
	}
	return effective.validateRecordedLanes()
}

func (effective Effective) validateRecordedTrust() error {
	if len(effective.ActorRules) != 3 {
		return fmt.Errorf("policy snapshot requires agent, human and service actor rules")
	}
	for _, kind := range []string{"agent", "human", "service"} {
		sets, ok := effective.ActorRules[kind]
		if !ok || sets == nil {
			return fmt.Errorf("policy snapshot is missing %s actor rules", kind)
		}
		for _, set := range sets {
			if _, ok := effective.ActionSets[set]; !ok {
				return fmt.Errorf("policy snapshot actor %s refers to undefined action set %q", kind, set)
			}
		}
	}
	for _, set := range []ActionSet{WorkSet, PersonDecisions} {
		if _, ok := effective.ActionSets[set]; !ok {
			return fmt.Errorf("policy snapshot is missing action set %q", set)
		}
	}
	for set, actions := range effective.ActionSets {
		if set == "" || actions == nil {
			return fmt.Errorf("policy snapshot has an incomplete action set %q", set)
		}
		for _, action := range actions {
			if action == "" {
				return fmt.Errorf("policy snapshot action set %q has an empty action", set)
			}
		}
	}
	return nil
}

func (effective Effective) validateRecordedLanes() error {
	for _, name := range []string{"routine", "standard", "assured"} {
		lane, ok := effective.Lanes[name]
		if !ok {
			return fmt.Errorf("policy snapshot is missing lane %q", name)
		}
		if lane.Checks == nil || lane.ReworkCap < 0 || lane.TokenBudget < 0 || !lane.ChecksSet && len(lane.Checks) != 0 {
			return fmt.Errorf("policy snapshot lanes.%s has invalid checks, rework cap or token budget", name)
		}
		if name == "standard" && len(lane.Paths) != 0 || name == "assured" && lane.AutoMerge {
			return fmt.Errorf("policy snapshot lanes.%s has forbidden paths or auto_merge", name)
		}
		if err := recordedGlobs("lanes."+name+".paths", lane.Paths); err != nil {
			return err
		}
		for _, check := range lane.Checks {
			if _, ok := effective.Checks[check]; !ok {
				return fmt.Errorf("policy snapshot lanes.%s refers to undefined check %q", name, check)
			}
		}
	}
	return nil
}

func recordedGlobs(field string, globs []string) error {
	if globs == nil {
		return fmt.Errorf("policy snapshot %s requires a path list", field)
	}
	for _, pattern := range globs {
		if _, err := pathglob.Compile(pattern); err != nil {
			return fmt.Errorf("policy snapshot %s: %w", field, err)
		}
	}
	return nil
}
