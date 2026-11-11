package main

import (
	"encoding/csv"
	"strings"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/validation"
)

// Retained text domains are inherited by P21 §4.1 from P20 §3.2.1.
// Opaque IDs and states use the newer P21 §2 / §4.1 scalar domain.
var declaredTextDomains = map[string]map[string]flagDomain{
	"issue create":       {"title": textDomain(512), "epic": opaqueFlagDomain("P21 §4.1 issue.create / §2 ID"), "parent": opaqueFlagDomain("P21 §4.1 issue.create / §2 ID"), "label": labelListDomain(false)},
	"issue update":       {"title": textDomain(512), "labels": labelListDomain(true)},
	"issue ac add":       {"text": textDomain(4096)},
	"label add":          {"colour": textDomain(64), "description": textDomain(4096)},
	"issue list":         {"status": opaqueFlagDomain("P21 §4.1 issue.list opaque state"), "epic": opaqueFlagDomain("P21 §4.1 issue.list / §2 ID"), "label": textDomain(255)},
	"issue transition":   {"to": opaqueFlagDomain("P21 §4.1 issue.transition opaque state")},
	"ready":              {"epic": opaqueFlagDomain("P21 §4.1 ready.get / §2 ID")},
	"next":               {"actor": nextActorDomain},
	"stats cost":         {"epic": opaqueFlagDomain("P21 §4.1 stats.cost / §2 ID")},
	"stats tokens":       {"epic": opaqueFlagDomain("P21 §4.1 stats.tokens / §2 ID")},
	"claim list":         {"issue": opaqueFlagDomain("P21 §4.2 claim.list / §2 ID"), "envelope": opaqueFlagDomain("P21 §4.2 claim.list / §2 ID")},
	"run list":           {"issue": opaqueFlagDomain("P22 §4 run.list / P21 §2 ID"), "envelope": opaqueFlagDomain("P22 §4 run.list / P21 §2 ID")},
	"receipt list":       {"run": opaqueFlagDomain("P21 §4.2 receipt.list / §2 ID"), "obligation": opaqueFlagDomain("P21 §4.2 receipt.list / §2 ID")},
	"review list":        {"issue": opaqueFlagDomain("P23 §4 review.list / §2 Review"), "envelope": opaqueFlagDomain("P23 §4 review.list / §2 Review")},
	"finding list":       {"issue": opaqueFlagDomain("P23 §4 finding.list / §2 Finding"), "review": opaqueFlagDomain("P23 §4 finding.list / §2 Finding")},
	"decision list":      {"issue": opaqueFlagDomain("P23 §4 decision.list / §2 DecisionRef"), "class": opaqueFlagDomain("P23 §4 decision.list / §2 DecisionRef")},
	"human-outcome list": {"issue": opaqueFlagDomain("P23 §4 human-outcome.list / §2 HumanOutcome"), "class": opaqueFlagDomain("P23 §4 human-outcome.list / §2 HumanOutcome")},
	"requirement list":   {"issue": opaqueFlagDomain("P23 §4 requirement.list / §2 RequirementRef")},
}

func textDomain(maximumBytes int) flagDomain {
	return flagDomain{[]string{"string"}, "P21 §4.1 retains P20 §3.2.1 text/label bounds", func(raw string) bool {
		_, err := validation.NormalizeText("flag", raw, maximumBytes)
		return err == nil
	}}
}

func opaqueFlagDomain(citation string) flagDomain {
	return flagDomain{[]string{"string"}, citation, func(raw string) bool { return utf8.ValidString(raw) && validOpaqueState(strings.TrimSpace(raw)) }}
}

func labelListDomain(allowEmpty bool) flagDomain {
	return flagDomain{[]string{"stringSlice"}, "P21 §4.1 retains P20 §3.2.1 repeated labels / empty --labels=", func(raw string) bool {
		if raw == "" {
			return allowEmpty
		}
		reader := csv.NewReader(strings.NewReader(raw))
		reader.FieldsPerRecord = -1
		records, err := reader.ReadAll()
		if err != nil || len(records) != 1 {
			return false
		}
		for _, label := range records[0] {
			if !textDomain(255).Accept(label) {
				return false
			}
		}
		return true
	}}
}
