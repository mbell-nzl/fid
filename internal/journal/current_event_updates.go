package journal

import (
	"fmt"
	"sort"
	"strings"
)

func validateCurrentUpdate(event Event, effect Effect) error {
	if effect.Op != "update" {
		return nil
	}
	values, guards := "", ""
	switch effect.Table {
	case "run_runtime":
		values, guards = "version", "version"
	case "runs":
		values, guards = "status terminal_repository_json ended_at", "status"
	case "claims":
		values, guards = "version status released_at release_reason", "version status"
		if event.EventType == "claim.expired" {
			values = "version status"
			guards = "version status actor_id"
		}
	case "kit_activations":
		values, guards = "version updated_at status", "version"
		if event.EventType == "activation.resource-observed" || event.EventType == "activation.recovery-selected" {
			values = "version updated_at"
		}
		if event.EventType == "activation.committed" {
			values += " point_of_no_return"
		}
	case "activation_components":
		values = "status"
	case "activation_resources":
		values, guards = "version observed_identity state evidence_ref error_code", "version"
	case "activation_recoveries":
		values = "completed_event_id"
	case "attention_items":
		values, guards = "version status resolution_event_id", "version status"
	case "idempotency_results":
		values, guards = "result_entity_id result_event_id aggregate_sequence result_json", "request_digest"
	default:
		return fmt.Errorf("current Event update table not admitted")
	}
	if !currentKeys(effect.Values, values) || !currentKeys(effect.Preconditions, guards) {
		return fmt.Errorf("current Event %s: wrong changed fields or guards for %s", event.EventType, effect.Table)
	}
	return validateCurrentTerminalValues(event, effect)
}

func currentKeys(values map[string]any, names string) bool {
	want := strings.Fields(names)
	sort.Strings(want)
	return strings.Join(sortedMapKeys(values), " ") == strings.Join(want, " ")
}

func validateCurrentTerminalValues(event Event, effect Effect) error {
	status := effect.Values["status"]
	valid := true
	switch effect.Table {
	case "runs":
		valid = effect.Preconditions["status"] == "active"
		switch event.EventType {
		case "run.closed":
			valid = valid && (status == "succeeded" || status == "failed")
		case "run.abandoned":
			valid = valid && status == "abandoned"
		case "run.recovered":
			valid = valid && status == "interrupted"
		}
	case "claims":
		valid = effect.Preconditions["status"] == "active"
		switch event.EventType {
		case "claim.expired":
			valid = valid && status == "expired"
		case "run.abandoned":
			valid = valid && status == "released" && effect.Values["release_reason"] == "abandoned"
		case "run.recovered":
			valid = valid && status == "superseded" && effect.Values["release_reason"] == "recovery-superseded"
		}
	case "attention_items":
		valid = status == "resolved" && effect.Preconditions["status"] == "open" && effect.Values["resolution_event_id"] == event.EventID
	case "activation_recoveries":
		valid = effect.Values["completed_event_id"] == event.EventID
	case "kit_activations", "activation_components":
		switch event.EventType {
		case "activation.component-staged":
			valid = status == "staged"
		case "activation.committed":
			valid = status == "committed"
		case "activation.finalized":
			valid = status == "finalized" || (effect.Table == "kit_activations" && effect.PK["id"] != event.Aggregate.ID && status == "superseded")
		case "activation.aborted":
			valid = status == "aborted"
		case "activation.recovery-completed":
			valid = status == "finalized" || status == "aborted" || (effect.Table == "kit_activations" && effect.PK["id"] != event.Aggregate.ID && status == "superseded")
		}
	}
	if !valid {
		return fmt.Errorf("current Event %s: terminal values do not match family", event.EventType)
	}
	return nil
}
