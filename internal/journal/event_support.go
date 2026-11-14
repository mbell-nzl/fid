package journal

import (
	"fmt"
	"strings"
)

// SchemaUnsupportedError identifies retired history before any Effects land.
// Generic version-1 issue history remains supported; wire version alone is not
// evidence that an Event belongs to a removed family.
type SchemaUnsupportedError struct {
	EventID   string
	EventType string
	Table     string
}

func (err *SchemaUnsupportedError) Error() string {
	if err.Table != "" {
		return fmt.Sprintf("schema-unsupported: Event %q family %q uses removed table %q", err.EventID, err.EventType, err.Table)
	}
	return fmt.Sprintf("schema-unsupported: Event %q has removed family %q", err.EventID, err.EventType)
}

// ValidateEventSupport is shared by journal projection and archive validation.
// Physical tables remain until RESET 1.5d2; they do not authorize old Effects.
func ValidateEventSupport(event Event) error {
	if removedEventFamilies[strings.ToLower(event.EventType)] {
		return &SchemaUnsupportedError{EventID: event.EventID, EventType: event.EventType}
	}
	for _, effect := range event.Effects {
		// The SQL executor accepts only bare ASCII identifiers. Validate here
		// too, so replay cannot apply a prefix before reaching invalid syntax.
		// Within this grammar, SQLite folds ASCII case for table resolution.
		table, err := validatedIdentifier(effect.Table)
		if err != nil {
			return fmt.Errorf("project event %q: %w", event.EventID, err)
		}
		if removedEffectTables[strings.ToLower(table)] {
			return &SchemaUnsupportedError{EventID: event.EventID, EventType: event.EventType, Table: effect.Table}
		}
	}
	return nil
}

var removedEventFamilies = map[string]bool{
	"activation.aborted":            true,
	"activation.committed":          true,
	"activation.component-staged":   true,
	"activation.finalized":          true,
	"activation.prepared":           true,
	"activation.recovery-completed": true,
	"activation.recovery-selected":  true,
	"activation.resource-observed":  true,
	"attention.resolved":            true,
	"claim.acquired":                true,
	"claim.expired":                 true,
	"claim.released":                true,
	"claim.renewed":                 true,
	"concern.disposed":              true,
	"concern.recorded":              true,
	"decision.referenced":           true,
	"deviation.recorded":            true,
	"envelope.activated":            true,
	"envelope.completed":            true,
	"envelope.created":              true,
	"evidence.linked":               true,
	"finding.addressed":             true,
	"gate.execution-started":        true,
	"gate.receipt-recorded":         true,
	"gate.registered":               true,
	"human-outcome.recorded":        true,
	"requirement.linked":            true,
	"requirement.referenced":        true,
	"review.recorded":               true,
	"run.abandoned":                 true,
	"run.budget-observed":           true,
	"run.checkpointed":              true,
	"run.closed":                    true,
	"run.execution-blocked":         true,
	"run.observed":                  true,
	"run.recovered":                 true,
	"run.started":                   true,
	"stop.raised":                   true,
	"stop.resolved":                 true,
	"workflow.assigned":             true,
	"workflow.defined":              true,
}

var removedEffectTables = map[string]bool{
	"kit_activations":              true,
	"activation_actor_profiles":    true,
	"activation_components":        true,
	"activation_gate_definitions":  true,
	"activation_resources":         true,
	"activation_recoveries":        true,
	"work_envelopes":               true,
	"envelope_revisions":           true,
	"envelope_issue_selection":     true,
	"envelope_sources":             true,
	"envelope_authority_grants":    true,
	"envelope_obligations":         true,
	"envelope_protected_decisions": true,
	"envelope_budgets":             true,
	"claims":                       true,
	"runs":                         true,
	"run_runtime":                  true,
	"run_observations":             true,
	"run_recoveries":               true,
	"run_budget_bindings":          true,
	"budget_samples":               true,
	"gate_registrations":           true,
	"execution_attempts":           true,
	"execution_receipts":           true,
	"receipt_artifacts":            true,
	"evidence_records":             true,
	"evidence_artifact_links":      true,
	"stops":                        true,
	"stop_resolutions":             true,
	"attention_items":              true,
	"attention_stop_links":         true,
	"review_records":               true,
	"review_findings":              true,
	"finding_dispositions":         true,
	"concerns":                     true,
	"deviations":                   true,
	"human_outcomes":               true,
	"decision_refs":                true,
	"requirement_refs":             true,
	"requirement_issue_links":      true,
	"system_capabilities":          true,
	"workflows":                    true,
	"workflow_states":              true,
	"workflow_transitions":         true,
	"workflow_transition_requires": true,
	"workflow_assignments":         true,
}
