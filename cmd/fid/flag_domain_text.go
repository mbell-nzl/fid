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
	"park":             {"hold": opaqueFlagDomain("RESET 1.5a2"), "reason": textDomain(4096), "waiting": textDomain(4096)},
	"heartbeat":        {"hold": opaqueFlagDomain("RESET 1.5a2")},
	"takeover":         {"reason": textDomain(4096)},
	"issue create":     {"title": textDomain(512), "epic": opaqueFlagDomain("P21 §4.1 issue.create / §2 ID"), "parent": opaqueFlagDomain("P21 §4.1 issue.create / §2 ID"), "label": labelListDomain(false)},
	"issue update":     {"title": textDomain(512), "labels": labelListDomain(true)},
	"issue ac add":     {"text": textDomain(4096)},
	"label add":        {"colour": textDomain(64), "description": textDomain(4096)},
	"issue list":       {"status": opaqueFlagDomain("P21 §4.1 issue.list opaque state"), "epic": opaqueFlagDomain("P21 §4.1 issue.list / §2 ID"), "label": textDomain(255)},
	"issue transition": {"to": opaqueFlagDomain("P21 §4.1 issue.transition opaque state"), "reason": textDomain(4096)},

	"next": {"actor": nextActorDomain},
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
