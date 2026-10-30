package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// flagDomain is a leaf/spelling declaration. The same declaration validates
// product input and enumerates the malformed-value corpus (DEC-280).
type flagDomain struct {
	Allowed  []string
	Citation string
	Accept   func(string) bool
}

func enumDomain(citation string, members ...string) flagDomain {
	sort.Strings(members)
	return flagDomain{members, citation, func(value string) bool {
		i := sort.SearchStrings(members, value)
		return i < len(members) && members[i] == value
	}}
}

func integerDomain(citation string, low, high int64) flagDomain {
	return flagDomain{[]string{fmt.Sprintf("%d..%d", low, high)}, citation, func(value string) bool {
		parsed, err := strconv.ParseInt(value, 0, 64)
		return err == nil && parsed >= low && parsed <= high
	}}
}

func patternDomain(citation, pattern string) flagDomain {
	compiled := regexp.MustCompile(pattern)
	return flagDomain{[]string{pattern}, citation, compiled.MatchString}
}

func subjectDomain(citation string, members ...string) flagDomain {
	domain := enumDomain(citation, members...)
	return flagDomain{[]string{"<type>:<id>"}, citation, func(value string) bool {
		parts := strings.Split(value, ":")
		return len(parts) == 2 && domain.Accept(parts[0]) && opaqueFlagDomain(citation).Accept(parts[1])
	}}
}

// Only the newest accepted row is cited. Generic pflag domains are registered
// mechanically below; these declarations supply the stronger row constraints.
var declaredStringDomains = map[string]map[string]flagDomain{
	"init": {"prefix": patternDomain("P21 §4.1 project.init", "^[A-Z][A-Z0-9]{0,15}$")},
	"issue create": {
		"kind":       enumDomain("P21 §4.1 issue.create", "epic", "subtask", "task"),
		"scheduling": enumDomain("P21 §4.1 issue.create", "prioritised", "sequential"),
	},
	"issue list": {
		"kind":            enumDomain("P21 §4.1 issue.list; P20 §3.2.1", "epic", "subtask", "task"),
		"status-category": enumDomain("P21 §4.1 issue.list; P20 §3.2.1", "cancelled", "completed", "started", "unstarted"),
		"order":           enumDomain("P21 §4.1 issue.list", "created", "priority", "sequence", "updated"),
		"priority":        {[]string{"0..5"}, "P21 §4.1 issue.list; P20 §3.2.1", validPrioritySet},
	},
	"workflow define": {
		"name":        patternDomain("P21 §4.1 workflow.define", "^[a-z][a-z0-9-]{0,63}$"),
		"applies-to":  enumDomain("P21 §4.1 workflow.define", "any", "epic", "subtask", "task"),
		"states":      {[]string{"string"}, "P21 §4.1 workflow.define", validStateList},
		"transitions": {[]string{"string"}, "P21 §4.1 workflow.define", validTransitionList},
		"categories":  {[]string{"cancelled", "completed", "started", "unstarted"}, "P21 §4.1 workflow.define", validCategoryList},
	},
	"workflow assign": {"name": patternDomain("P21 §4.1 workflow.define, workflow.assign", "^[a-z][a-z0-9-]{0,63}$"), "kind": enumDomain("P21 §4.1 workflow.assign", "epic", "subtask", "task")},
	"workflow list":   {"kind": enumDomain("P21 §4.1 workflow.list", "epic", "subtask", "task")},
	"link add": {
		"from-type": enumDomain("P21 §4.1 link.add", "issue", "requirement"),
		"to-type":   enumDomain("P21 §4.1 link.add", "issue", "requirement"),
	},
	"link remove": {
		"from-type": enumDomain("P21 §4.1 link.remove", "issue", "requirement"),
		"to-type":   enumDomain("P21 §4.1 link.remove", "issue", "requirement"),
	},
	"link list":    {"type": enumDomain("P21 §4.1 link.list", "issue", "requirement")},
	"stats cost":   {"group-by": enumDomain("P21 §4.1 stats.cost", "actor", "epic", "model")},
	"stats tokens": {"group-by": enumDomain("P21 §4.1 stats.tokens", "actor", "epic", "model")},
	"claim list":   {"status": enumDomain("P22 §2 Claim (amends P21 §4.2 claim.list)", "active", "expired", "released", "superseded")},
	"run list":     {"status": enumDomain("P22 §2 Run status / §4 (amends P21 §4.2)", "abandoned", "active", "failed", "interrupted", "succeeded")},
	"receipt list": {"outcome": enumDomain("P22 §5 cancellation / schema §1 (amends P21 §4.2 receipt.list)", "cancelled", "failed", "signaled", "spawn-failed", "succeeded", "timed-out")},
	"stop list": {
		"status":  enumDomain("P22 §4 stop.list", "open", "resolved"),
		"subject": subjectDomain("P22 §4 stop.list", "run", "envelope", "activation"),
	},
	"attention list": {
		"status":  enumDomain("P22 §4 attention.list", "open", "resolved"),
		"subject": subjectDomain("P22 §4 attention.list", "run", "envelope", "activation"),
	},
	"review list": {"verdict": enumDomain("P23 §4 review.list / §2 Review", "blocked", "changes-required", "pass")},
	"finding list": {
		"state":    enumDomain("P23 §4 finding.list / §2 Finding", "addressed", "inapplicable", "open", "resolved"),
		"severity": enumDomain("P23 §4 finding.list / §2 Finding", "blocker", "major", "minor", "nit"),
	},
	"concern list": {
		"status":   enumDomain("P23 §4 concern.list / §2 Concern", "disposed", "open"),
		"severity": enumDomain("P23 §4 concern.list / §2 Concern", "blocker", "major", "minor", "nit"),
		"subject":  subjectDomain("P23 §4 concern.list / §2 JudgmentSubject", "issue", "envelope", "run", "review", "finding", "requirement"),
	},
	"deviation list": {
		"disposition": enumDomain("P23 §4 deviation.list / §2 Deviation", "accepted", "pending-decision", "rejected"),
		"subject":     subjectDomain("P23 §4 deviation.list / §2 JudgmentSubject", "issue", "envelope", "run", "review", "finding", "requirement"),
	},
	"decision list":      {"outcome": enumDomain("P23 §4 decision.list / §2 DecisionRef", "approved", "deferred", "rejected")},
	"human-outcome list": {"result": enumDomain("P23 §4 human-outcome.list / §2 HumanOutcome", "acknowledged", "approved", "changes-required", "rejected")},
	"requirement list":   {"status": enumDomain("P23 §4 requirement.list / §2 RequirementRef", "active", "superseded", "withdrawn")},
	"activity list": {"event": enumDomain("P21 §3 ActivityEvent amended by P23 §2 ActivityEvent",
		"ac-added", "claim-released", "claim-renewed", "claim-superseded", "comment-added", "gate-recorded", "issue-created", "issue-label-added", "issue-label-removed", "issue-updated", "link-added", "link-removed", "status-changed",
		"concern-recorded", "concern-disposed", "review-recorded", "finding-addressed", "finding-resolved", "deviation-recorded", "decision-referenced", "human-outcome-recorded", "requirement-referenced", "requirement-linked")},
}

var declaredIntegerDomains = map[string]map[string]flagDomain{
	"issue create": {
		"priority": integerDomain("P21 §4.1; P20 §3.2.1", 0, 5),
		"sequence": integerDomain("P21 §4.1; P20 §3.2.1", 1, 2147483647),
	},
	"issue update": {"priority": integerDomain("P21 §4.1; P20 §3.2.1", 0, 5)},
	"issue ac add": {"ordinal": integerDomain("P21 §4.1; P20 §3.2.1", 1, 2147483647)},
	"gate resolve": {"version": integerDomain("P21 §4.2 gate.resolve", 1, 9223372036854775807)},
}

func validPrioritySet(value string) bool {
	parts := strings.Split(value, ",")
	if len(parts) > 6 {
		return false
	}
	seen := map[int64]bool{}
	for _, part := range parts {
		number, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil || number < 0 || number > 5 || seen[number] {
			return false
		}
		seen[number] = true
	}
	return true
}

func validCategoryList(value string) bool {
	parts := strings.Split(value, ",")
	if len(parts) > 64 {
		return false
	}
	domain := enumDomain("P21 §4.1 workflow.define", "cancelled", "completed", "started", "unstarted")
	for _, part := range parts {
		if !domain.Accept(strings.TrimSpace(part)) {
			return false
		}
	}
	return true
}

var declaredBooleanFlags = map[string]string{
	"issue list":        "has-concerns",
	"maintenance check": "warn-as-error",
	"rebuild":           "verbose",
	"workflow assign":   "project-wide",
	"workflow list":     "active",
}

func flagDomainFor(cmd *cobra.Command, flag *pflag.Flag) flagDomain {
	key := rootCommandKey(cmd)
	if domain, ok := declaredStringDomains[key][flag.Name]; ok {
		return domain
	}
	if domain, ok := declaredTextDomains[key][flag.Name]; ok {
		return domain
	}
	if domain, ok := declaredIntegerDomains[key][flag.Name]; ok {
		return domain
	}
	if flag.Name == "limit" {
		return integerDomain("P21 §4.2; P22 §4; P23 §1,4 lists", 1, 10000)
	}
	if flag.Name == "expected-revision" {
		domain := integerDomain("P21 §4.1; P20 §3.2.1", 1, 9223372036854775807)
		domain.Accept = func(raw string) bool {
			value, err := strconv.ParseInt(raw, 10, 64)
			return err == nil && value >= 1 && strings.Trim(raw, "0123456789") == ""
		}
		return domain
	}
	if declaredBooleanFlags[key] == flag.Name {
		return flagDomain{[]string{"bool"}, "P21 §4.1 boolean leaf flag", func(string) bool { return true }}
	}
	if flag.Name == "example" {
		if _, adopted := inputDescriptorRegistry[key]; adopted {
			return flagDomain{[]string{"bool"}, "P24 §8 --example", func(string) bool { return true }}
		}
	}
	if flag.Name == "idempotency-key" {
		return flagDomain{[]string{"string"}, "P21 §2 idempotency key", validIdempotencyKey}
	}
	if flag.Name == "actor" && !strings.HasPrefix(key, "activation ") && cmd.Flags().Lookup("input") != nil {
		return opaqueFlagDomain("P21 §2 activated Actor / §4.2 mutation row")
	}
	if flag.Name == "actor" && !strings.HasPrefix(key, "activation ") && cmd.Flags().Lookup("input") == nil {
		return patternDomain("P21 §2 standalone Actor", `^(agent|human|service):[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)
	}
	return flagDomain{[]string{flag.Value.Type()}, "pflag type; P24 §8 for --example", func(string) bool { return true }}
}

type declaredFlagValue struct {
	pflag.Value
	domain flagDomain
}

func (value declaredFlagValue) Set(raw string) error {
	if err := value.Value.Set(raw); err != nil {
		return err
	}
	if !value.domain.Accept(raw) {
		return &flagDomainError{Raw: raw}
	}
	return nil
}

type flagDomainError struct{ Raw string }

func (err *flagDomainError) Error() string { return "value outside declared flag domain" }

func installFlagDomains(cmd *cobra.Command) {
	cmd.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "help" {
			return
		}
		domain := flagDomainFor(cmd, flag)
		if domain.Citation == "pflag type; P24 §8 for --example" || flag.Value.Type() == "bool" {
			return
		}
		flag.Value = declaredFlagValue{flag.Value, domain}
	})
}

func validIdempotencyKey(value string) bool {
	if len(value) < 1 || len(value) > 128 || strings.TrimSpace(value) != value {
		return false
	}
	for _, char := range value {
		if char < ' ' || char > '~' {
			return false
		}
	}
	return true
}

func validOpaqueState(value string) bool {
	if !utf8.ValidString(value) {
		return false
	}
	count := 0
	for _, char := range value {
		if !unicode.IsPrint(char) {
			return false
		}
		count++
	}
	return count >= 1 && count <= 128
}

func validStateList(value string) bool {
	parts := strings.Split(value, ",")
	if len(parts) > 64 {
		return false
	}
	seen := map[string]bool{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if !validOpaqueState(part) || seen[part] {
			return false
		}
		seen[part] = true
	}
	return true
}

func validTransitionList(value string) bool {
	parts := strings.Split(value, ",")
	if len(parts) > 4096 {
		return false
	}
	seen := map[string]bool{}
	for _, part := range parts {
		endpoints := strings.Split(part, "->")
		if len(endpoints) != 2 {
			return false
		}
		from, to := strings.TrimSpace(endpoints[0]), strings.TrimSpace(endpoints[1])
		if !validOpaqueState(from) || !validOpaqueState(to) || seen[from+"->"+to] {
			return false
		}
		seen[from+"->"+to] = true
	}
	return true
}
