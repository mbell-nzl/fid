package policy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"slices"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
)

// BuiltInRulesVersion changes whenever built-in policy semantics change.
const BuiltInRulesVersion = 2

type ActionSet string

const (
	WorkSet         ActionSet = "work"
	PersonDecisions ActionSet = "person-decisions"
	ActorPattern              = `^(agent|human|service):[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`
)

var actorPattern = regexp.MustCompile(ActorPattern)

func actorRules() map[string][]ActionSet {
	return map[string][]ActionSet{"agent": {WorkSet}, "human": {WorkSet, PersonDecisions}, "service": {}}
}

// Allows is the lean core's built-in trust by actor type. The kernel's
// Envelope authority remains independent until its removal in 1.5c/1.5d.
func Allows(actor string, set ActionSet) bool {
	if !actorPattern.MatchString(actor) {
		return false
	}
	kind, _, _ := strings.Cut(actor, ":")
	return slices.Contains(actorRules()[kind], set)
}

type Effective struct {
	RulesVersion    int                     `json:"rules_version"`
	HardAssured     []string                `json:"hard_assured"`
	ActorRules      map[string][]ActionSet  `json:"actor_rules"`
	ActionSets      map[ActionSet][]string  `json:"action_sets"`
	Checks          map[string]config.Check `json:"checks"`
	Lanes           map[string]config.Lane  `json:"lanes"`
	TestAdequacy    []string                `json:"test_adequacy"`
	TestAdequacySet bool                    `json:"test_adequacy_set"`
}

// Resolve combines the validated file with the rules a project cannot waive.
// Set-valued lists are sorted and deduplicated; argv keeps execution order.
func Resolve(file config.ProjectFile) Effective {
	result := Effective{RulesVersion: BuiltInRulesVersion, HardAssured: []string{".fid.toml", ".github/workflows/**"},
		ActorRules: actorRules(), ActionSets: map[ActionSet][]string{
			WorkSet:         {"take", "start", "check", "park", "finish", "review"},
			PersonDecisions: {"cancel", "approve", "waive", "takeover", "allow-rework", "close-finding", "accept-work", "approve-epic", "answer", "close-epic"},
		}, Checks: map[string]config.Check{}, Lanes: map[string]config.Lane{}, TestAdequacy: canonicalSet(file.TestAdequacy), TestAdequacySet: file.TestAdequacySet}
	for name, check := range file.Checks {
		check.Run = slices.Clone(check.Run)
		check.DefinedBy, check.WhenChanged = canonicalSet(check.DefinedBy), canonicalSet(check.WhenChanged)
		result.Checks[name] = check
	}
	for name, lane := range file.Lanes {
		lane.Paths, lane.Checks = canonicalSet(lane.Paths), canonicalSet(lane.Checks)
		result.Lanes[name] = lane
	}
	return result
}

func canonicalSet(values []string) []string {
	result := append([]string{}, values...)
	slices.Sort(result)
	return slices.Compact(result)
}

// Canonical uses JSON's stable struct field and sorted map key encoding.
func (effective Effective) Canonical() ([]byte, error) {
	return json.Marshal(effective)
}

func (effective Effective) Digest() string {
	// Effective contains only JSON-encodable scalars, slices, maps and structs.
	encoded, _ := effective.Canonical()
	digest := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(digest[:])
}
