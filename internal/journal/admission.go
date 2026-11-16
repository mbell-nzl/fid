package journal

import (
	"fmt"
	"strings"
	"time"
)

// Admission applies one rule set to every Event, whichever path brings it
// into the ledger: a live append (checked before the write), replay
// (catch-up and rebuild) and archive import. This file holds the rules that
// need no projection state; the family validators in validateCurrentEvent
// follow them, inside the projection transaction.

// journalLineBudget is the journal's line cap: the aggregate and sequence
// scanners, and archive import, read lines shorter than this, newline
// included. A longer line would stop every later write.
const journalLineBudget = 1 << 20

func withinLineBudget(lineLength int) bool { return lineLength+1 < journalLineBudget }

// ValidateEventShape is the path-independent part of admission, shared by
// projection, replay's pre-scan, export and archive import.
func ValidateEventShape(event Event) error {
	if event.EventSchema != EventSchemaV1 && event.EventSchema != EventSchemaV2 {
		return fmt.Errorf("project event %q: unsupported Event schema %d", event.EventID, event.EventSchema)
	}
	encoded, err := MarshalCanonicalJSON(event)
	if err != nil {
		return err
	}
	return validateEventShape(event, len(encoded))
}

func validateEventShape(event Event, encodedLength int) error {
	if err := ValidateEventSupport(event); err != nil {
		return err
	}
	if err := validateDatedEventID(event.EventID); err != nil {
		return err
	}
	if err := validateLowercaseIdentifiers(event); err != nil {
		return err
	}
	if !withinLineBudget(encodedLength) {
		return fmt.Errorf("%s Event exceeds the journal's 1 MiB line budget", eventSizeLabel(event.EventType))
	}
	if !event.IsV2() {
		return validateLegacyEffects(event)
	}
	return validateFamilyShape(event)
}

// maxEventSequence bounds a day's Event sequence far beyond any real day,
// so the next sequence always exists.
const maxEventSequence = 999999999

// The writers allocate an Event ID by reading every ID of its day
// (nextSequence): each ID of three dash-separated parts whose middle part is
// that date, whatever its prefix. An ID the scanner cannot parse, or that
// leaves no next sequence, would stop every write on that date, so each ID
// dated that way must parse as the scanner parses it. IDs without a date
// there are never read.
func validateDatedEventID(id string) error {
	parts := strings.Split(id, "-")
	if len(parts) != 3 {
		return nil
	}
	if _, err := time.Parse("20060102", parts[1]); err != nil {
		return nil
	}
	if sequence, err := parseEventSequence(id); err != nil || sequence > maxEventSequence {
		return fmt.Errorf("Event ID %q needs the form EVT-<date>-<sequence>, with a sequence from 1 to %d", id, maxEventSequence)
	}
	return nil
}

// The families that checked their own size before it was shared keep
// their wording.
func eventSizeLabel(family string) string {
	labels := map[string]string{"epic.closed": "epic close", "task.finished": "completion", "task.approved": "approval", "finding.closed": "finding close", "review.round-recorded": "review"}
	if label, found := labels[family]; found {
		return label
	}
	return fmt.Sprintf("%q", family)
}

// SQLite resolves table and column names without regard to ASCII case, so
// "HOLDS" or "LANE" would reach the same storage as "holds" or "lane" while
// slipping past every exact-name rule. fid's schema is all lowercase.
func validateLowercaseIdentifiers(event Event) error {
	for _, effect := range event.Effects {
		if strings.ToLower(effect.Table) != effect.Table {
			return fmt.Errorf("project event %q: table %q must be lowercase", event.EventID, effect.Table)
		}
		for _, row := range []map[string]any{effect.PK, effect.Values, effect.Prior, effect.Preconditions} {
			for column := range row {
				if strings.ToLower(column) != column {
					return fmt.Errorf("project event %q: column %q must be lowercase", event.EventID, column)
				}
			}
		}
	}
	return nil
}

// Event-v1 is Issue history from before the lean core. No writer produces
// it now and import refuses it; replay still reads it, but only on the
// tables of the schema that history wrote (migrations 0001 and 0002), never
// the lean core's or the projection's own.
var legacyIssueTables = map[string]bool{"issues": true, "roles": true, "epic_metadata": true, "requirements": true, "decisions": true, "escalations": true, "triage_events": true, "task_completions": true, "acceptance_criteria": true, "ac_verifications": true, "agent_runs": true, "gate_runs": true, "labels": true, "issue_labels": true, "links": true, "activity_log": true}

func validateLegacyEffects(event Event) error {
	for _, effect := range event.Effects {
		if !legacyIssueTables[effect.Table] {
			return fmt.Errorf("Event-v1 %q may change only Issue tables, not %q", event.EventID, effect.Table)
		}
	}
	return nil
}

// ValidateEventGrammar checks an Event-v2's family, aggregate type and
// ordered Effect grammar: the tables and operations its family owns.
// Archive import reports it apart from the rules that bind the Event's
// subject, actor and retry result.
func ValidateEventGrammar(event Event) error {
	rule, found := retainedCurrentRules[event.EventType]
	if !found {
		return fmt.Errorf("current Event: unsupported Event family %q", event.EventType)
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	if event.Aggregate.Type != rule.aggregate {
		return fmt.Errorf("current Event: retained aggregate mismatch")
	}
	for _, effect := range event.Effects {
		switch effect.Op {
		case "insert", "update", "delete", "activity":
		default:
			return fmt.Errorf("project event %q: unsupported effect op %q", event.EventID, effect.Op)
		}
	}
	return validateCurrentStages(event, rule)
}

func validateFamilyShape(event Event) error {
	if err := ValidateEventGrammar(event); err != nil {
		return err
	}
	rule := retainedCurrentRules[event.EventType]
	if !strings.Contains(" "+rule.actors+" ", " "+event.ActorV2.ActorType+" ") {
		return fmt.Errorf("current Event %s: %s actor not admitted", event.EventType, event.ActorV2.ActorType)
	}
	// An Issue Event changes only its own Issue row.
	for _, effect := range event.Effects {
		if id, keyed := effect.PK["id"]; keyed && rule.aggregate == "issue" && effect.Table == "issues" && id != event.Aggregate.ID {
			return fmt.Errorf("issues identity differs from Event aggregate")
		}
	}
	// A retry result belongs to the Event's own request, and only then.
	retried := event.Idempotency != nil
	for _, effect := range event.Effects {
		retried = retried || effect.Table == "idempotency_results"
	}
	if retried {
		_, err := validateIdempotencyBinding(event)
		return err
	}
	return nil
}
