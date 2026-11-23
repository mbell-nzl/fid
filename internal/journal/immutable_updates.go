package journal

import "strings"

// The immutable facts of these lean tables are shared across all retained
// profiles. In particular, profiles before person_decisions still need this
// check: they predate the lifecycle validator's installed-table guard.
var immutableUpdateColumns = map[string]string{
	"check_requests":       "response_json",
	"referenced_artifacts": "availability",
	"idempotency_results":  "result_entity_id result_event_id aggregate_sequence result_json",
	"questions":            "resolved_at resolution decision_id",
	"holds":                "ended_at end_reason last_heartbeat_at expires_at kept park_reason rework_at",
}

// ValidateEventUpdates checks mutable fields before projection on every path.
// Import also calls it before schema validation to retain its effect-invalid
// classification; both calls use this same rule set.
func ValidateEventUpdates(event Event) error {
	for _, effect := range event.Effects {
		if effect.Op != "update" {
			continue
		}
		if effect.Table == "holds" && event.EventType == "issue.transitioned" {
			if reason, present := effect.Values["end_reason"]; present && reason != "task-cancelled" {
				return rowFailure(event, "holds", "invalid task hold end")
			}
		}
		columns, limited := immutableUpdateColumns[effect.Table]
		if !limited {
			continue
		}
		for column := range effect.Values {
			if !strings.Contains(" "+columns+" ", " "+column+" ") {
				return rowFailure(event, effect.Table, column+" is immutable in this Event family")
			}
		}
	}
	return nil
}
