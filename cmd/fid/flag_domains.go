package main

import (
	"fmt"
	"regexp"
	"slices"
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

// Only the newest accepted row is cited. Generic pflag domains are registered
// mechanically below; these declarations supply the stronger row constraints.
var declaredStringDomains = map[string]map[string]flagDomain{
	"init":         {"prefix": patternDomain("P21 §4.1 project.init", "^[A-Z][A-Z0-9]{0,15}$")},
	"issue update": {"lane": enumDomain("RESET 1.5a4 planning lane", "routine", "standard", "assured")},
	"issue create": {
		"lane":       enumDomain("RESET 1.5a4 planning lane", "routine", "standard", "assured"),
		"kind":       enumDomain("P21 §4.1 issue.create", "epic", "subtask", "task"),
		"scheduling": enumDomain("P21 §4.1 issue.create", "prioritised", "sequential"),
	},
	"issue list": {
		"kind":            enumDomain("P21 §4.1 issue.list; P20 §3.2.1", "epic", "subtask", "task"),
		"status-category": enumDomain("P21 §4.1 issue.list; P20 §3.2.1", "cancelled", "completed", "started", "unstarted"),
		"order":           enumDomain("P21 §4.1 issue.list", "created", "priority", "sequence", "updated"),
		"priority":        {[]string{"0..5"}, "P21 §4.1 issue.list; P20 §3.2.1", validPrioritySet},
	},

	"link add": {
		"from-type": enumDomain("P21 §4.1 link.add", "issue", "requirement"),
		"to-type":   enumDomain("P21 §4.1 link.add", "issue", "requirement"),
	},
	"link remove": {
		"from-type": enumDomain("P21 §4.1 link.remove", "issue", "requirement"),
		"to-type":   enumDomain("P21 §4.1 link.remove", "issue", "requirement"),
	},
	"link list": {"type": enumDomain("P21 §4.1 link.list", "issue", "requirement")},

	"activity list": {"event": enumDomain("P21 §3 ActivityEvent amended by P23 §2 ActivityEvent",
		"ac-added", "comment-added", "issue-created", "issue-label-added", "issue-label-removed", "issue-updated", "link-added", "link-removed", "status-changed")},
}

var declaredIntegerDomains = map[string]map[string]flagDomain{
	"issue create": {
		"priority": integerDomain("P21 §4.1; P20 §3.2.1", 0, 5),
		"sequence": integerDomain("P21 §4.1; P20 §3.2.1", 1, 2147483647),
	},
	"issue update": {"priority": integerDomain("P21 §4.1; P20 §3.2.1", 0, 5)},
	"issue ac add": {"ordinal": integerDomain("P21 §4.1; P20 §3.2.1", 1, 2147483647)},
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

// declaredBooleanFlags names each leaf's boolean flags, separated by spaces.
var declaredBooleanFlags = map[string]string{
	"maintenance check": "warn-as-error",
	"next":              "claim explain",
	"park":              "interrupted",
	"rebuild":           "verbose",
	"status":            "all",
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
	if flag.Name == "expected-revision" {
		domain := integerDomain("P21 §4.1; P20 §3.2.1", 1, 9223372036854775807)
		domain.Accept = func(raw string) bool {
			value, err := strconv.ParseInt(raw, 10, 64)
			return err == nil && value >= 1 && strings.Trim(raw, "0123456789") == ""
		}
		return domain
	}
	if slices.Contains(strings.Fields(declaredBooleanFlags[key]), flag.Name) {
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
	if flag.Name == "actor" {
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
