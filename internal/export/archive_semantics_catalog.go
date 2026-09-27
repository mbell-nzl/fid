package export

import "strings"

type semanticEventRule struct {
	aggregate string
	stages    []semanticEffectStage
}

type semanticEffectStage struct {
	minimum int
	maximum int
	allowed map[string]struct{}
}

func semanticRule(aggregate string, stages ...semanticEffectStage) semanticEventRule {
	return semanticEventRule{aggregate: aggregate, stages: stages}
}

func semanticStage(minimum int, operations ...string) semanticEffectStage {
	allowed := make(map[string]struct{}, len(operations))
	for _, operation := range operations {
		allowed[operation] = struct{}{}
	}
	return semanticEffectStage{minimum: minimum, maximum: -1, allowed: allowed}
}

func semanticEffectKey(table, operation string) string { return table + "\x00" + operation }

func (stage semanticEffectStage) accepts(table, operation string) bool {
	_, found := stage.allowed[semanticEffectKey(table, operation)]
	return found
}

// semanticEventRules is the closed Event-v2 family catalog. Stages preserve
// the accepted projection order while allowing the documented optional and
// repeated rows within a stage. An idempotency result, where present, is last.
var semanticEventRules = map[string]semanticEventRule{
	"review.recorded":                semanticRule("review", semanticStage(1, semanticEffectKey("review_records", "insert")), semanticStage(0, semanticEffectKey("review_findings", "insert"), semanticEffectKey("finding_dispositions", "insert"), semanticEffectKey("budget_samples", "insert"), semanticEffectKey("stops", "insert"), semanticEffectKey("attention_items", "insert"), semanticEffectKey("attention_stop_links", "insert"), semanticEffectKey("stop_resolutions", "insert"), semanticEffectKey("attention_items", "update"), semanticEffectKey("activity_log", "activity"), semanticEffectKey("idempotency_results", "insert"))),
	"finding.addressed":              semanticRule("finding", semanticStage(1, semanticEffectKey("finding_dispositions", "insert")), semanticStage(0, semanticEffectKey("activity_log", "activity"), semanticEffectKey("idempotency_results", "insert"))),
	"concern.recorded":               semanticRule("concern", semanticStage(1, semanticEffectKey("judgment_concerns", "insert")), semanticStage(0, semanticEffectKey("activity_log", "activity"), semanticEffectKey("idempotency_results", "insert"))),
	"concern.disposed":               semanticRule("concern", semanticStage(1, semanticEffectKey("judgment_concerns", "update")), semanticStage(0, semanticEffectKey("activity_log", "activity"), semanticEffectKey("idempotency_results", "insert"))),
	"deviation.recorded":             semanticRule("deviation", semanticStage(1, semanticEffectKey("deviations", "insert")), semanticStage(0, semanticEffectKey("activity_log", "activity"), semanticEffectKey("idempotency_results", "insert"))),
	"decision.referenced":            semanticRule("decision", semanticStage(1, semanticEffectKey("decision_refs", "insert")), semanticStage(0, semanticEffectKey("stop_resolutions", "insert"), semanticEffectKey("attention_items", "update"), semanticEffectKey("activity_log", "activity"), semanticEffectKey("idempotency_results", "insert"))),
	"human-outcome.recorded":         semanticRule("human-outcome", semanticStage(1, semanticEffectKey("human_outcomes", "insert")), semanticStage(0, semanticEffectKey("stop_resolutions", "insert"), semanticEffectKey("attention_items", "update"), semanticEffectKey("activity_log", "activity"), semanticEffectKey("idempotency_results", "insert"))),
	"requirement.referenced":         semanticRule("requirement", semanticStage(1, semanticEffectKey("requirement_refs", "insert")), semanticStage(0, semanticEffectKey("requirement_refs", "update"), semanticEffectKey("activity_log", "activity"), semanticEffectKey("idempotency_results", "insert"))),
	"requirement.linked":             semanticRule("requirement", semanticStage(1, semanticEffectKey("requirement_issue_links", "insert")), semanticStage(0, semanticEffectKey("activity_log", "activity"), semanticEffectKey("idempotency_results", "insert"))),
	"project.initialized":            semanticRule("project"),
	"activation.prepared":            semanticRule("activation", semanticStage(1, semanticEffectKey("kit_activations", "insert")), semanticStage(1, semanticEffectKey("activation_resources", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"activation.component-staged":    semanticRule("activation", semanticStage(1, semanticEffectKey("kit_activations", "update")), semanticStage(1, semanticEffectKey("activation_components", "insert")), semanticStage(1, semanticEffectKey("activation_actor_profiles", "insert"), semanticEffectKey("activation_gate_definitions", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"activation.committed":           semanticRule("activation", semanticStage(1, semanticEffectKey("kit_activations", "update")), semanticStage(0, semanticEffectKey("activation_components", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"activation.finalized":           semanticRule("activation", semanticStage(1, semanticEffectKey("kit_activations", "update")), semanticStage(0, semanticEffectKey("activation_components", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"activation.aborted":             semanticRule("activation", semanticStage(1, semanticEffectKey("kit_activations", "update")), semanticStage(0, semanticEffectKey("activation_components", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"activation.resource-observed":   semanticRule("activation", semanticStage(1, semanticEffectKey("kit_activations", "update")), semanticStage(1, semanticEffectKey("activation_resources", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"activation.recovery-selected":   semanticRule("activation", semanticStage(1, semanticEffectKey("kit_activations", "update")), semanticStage(1, semanticEffectKey("activation_recoveries", "insert")), semanticStage(0, semanticEffectKey("stops", "insert")), semanticStage(0, semanticEffectKey("attention_items", "insert")), semanticStage(0, semanticEffectKey("attention_stop_links", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"activation.recovery-completed":  semanticRule("activation", semanticStage(1, semanticEffectKey("kit_activations", "update")), semanticStage(1, semanticEffectKey("activation_recoveries", "update")), semanticStage(0, semanticEffectKey("activation_components", "update")), semanticStage(0, semanticEffectKey("stop_resolutions", "insert")), semanticStage(0, semanticEffectKey("attention_items", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"gate.registered":                semanticRule("gate-registration", semanticStage(1, semanticEffectKey("gate_registrations", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"envelope.created":               semanticRule("envelope", semanticStage(1, semanticEffectKey("work_envelopes", "insert")), semanticStage(1, semanticEffectKey("envelope_revisions", "insert")), semanticStage(1, semanticEffectKey("envelope_issue_selection", "insert")), semanticStage(0, semanticEffectKey("envelope_sources", "insert")), semanticStage(0, semanticEffectKey("envelope_authority_grants", "insert")), semanticStage(1, semanticEffectKey("envelope_obligations", "insert")), semanticStage(0, semanticEffectKey("envelope_protected_decisions", "insert"), semanticEffectKey("envelope_budgets", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"envelope.activated":             semanticRule("envelope", semanticStage(1, semanticEffectKey("work_envelopes", "update")), semanticStage(1, semanticEffectKey("envelope_revisions", "update")), semanticStage(0, semanticEffectKey("stops", "insert"), semanticEffectKey("attention_items", "insert"), semanticEffectKey("attention_stop_links", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"envelope.completed":             semanticRule("envelope", semanticStage(1, semanticEffectKey("work_envelopes", "update")), semanticStage(1, semanticEffectKey("envelope_revisions", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"claim.acquired":                 semanticRule("claim", semanticStage(1, semanticEffectKey("claims", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"claim.renewed":                  semanticRule("claim", semanticStage(1, semanticEffectKey("claims", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"claim.released":                 semanticRule("claim", semanticStage(1, semanticEffectKey("claims", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"claim.expired":                  semanticRule("claim", semanticStage(1, semanticEffectKey("claims", "update")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"run.started":                    semanticRule("run", semanticStage(1, semanticEffectKey("runs", "insert")), semanticStage(1, semanticEffectKey("run_runtime", "insert")), semanticStage(0, semanticEffectKey("run_budget_bindings", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"run.observed":                   semanticRule("run", semanticStage(1, semanticEffectKey("run_observations", "insert")), semanticStage(1, semanticEffectKey("run_runtime", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"run.checkpointed":               semanticRule("run", semanticStage(1, semanticEffectKey("run_observations", "insert")), semanticStage(1, semanticEffectKey("run_runtime", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"run.recovered":                  semanticRule("run", semanticStage(0, semanticEffectKey("runs", "update")), semanticStage(0, semanticEffectKey("run_runtime", "update")), semanticStage(0, semanticEffectKey("claims", "update")), semanticStage(1, semanticEffectKey("claims", "insert")), semanticStage(1, semanticEffectKey("runs", "insert")), semanticStage(1, semanticEffectKey("run_runtime", "insert")), semanticStage(0, semanticEffectKey("run_budget_bindings", "insert")), semanticStage(1, semanticEffectKey("run_recoveries", "insert")), semanticStage(0, semanticEffectKey("budget_samples", "insert")), semanticStage(0, semanticEffectKey("stops", "insert")), semanticStage(0, semanticEffectKey("attention_items", "insert")), semanticStage(0, semanticEffectKey("attention_stop_links", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"run.budget-observed":            semanticRule("run", semanticStage(1, semanticEffectKey("budget_samples", "insert")), semanticStage(0, semanticEffectKey("stops", "insert")), semanticStage(0, semanticEffectKey("attention_items", "insert")), semanticStage(0, semanticEffectKey("attention_stop_links", "insert")), semanticStage(1, semanticEffectKey("run_runtime", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"run.execution-blocked":          semanticRule("run", semanticStage(1, semanticEffectKey("stops", "insert")), semanticStage(0, semanticEffectKey("attention_items", "insert")), semanticStage(0, semanticEffectKey("attention_stop_links", "insert")), semanticStage(1, semanticEffectKey("run_runtime", "update")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"run.closed":                     semanticRule("run", semanticStage(1, semanticEffectKey("runs", "update")), semanticStage(1, semanticEffectKey("run_runtime", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"run.abandoned":                  semanticRule("run", semanticStage(1, semanticEffectKey("runs", "update")), semanticStage(1, semanticEffectKey("run_runtime", "update")), semanticStage(1, semanticEffectKey("claims", "update")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"gate.execution-started":         semanticRule("attempt", semanticStage(1, semanticEffectKey("execution_attempts", "insert")), semanticStage(1, semanticEffectKey("idempotency_results", "insert")), semanticStage(0, semanticEffectKey("budget_samples", "insert")), semanticStage(0, semanticEffectKey("stops", "insert")), semanticStage(0, semanticEffectKey("attention_items", "insert")), semanticStage(0, semanticEffectKey("attention_stop_links", "insert")), semanticStage(0, semanticEffectKey("run_runtime", "update"))),
	"gate.receipt-recorded":          semanticRule("receipt", semanticStage(0, semanticEffectKey("referenced_artifacts", "insert")), semanticStage(1, semanticEffectKey("execution_receipts", "insert")), semanticStage(0, semanticEffectKey("receipt_artifacts", "insert")), semanticStage(0, semanticEffectKey("budget_samples", "insert")), semanticStage(0, semanticEffectKey("stops", "insert")), semanticStage(0, semanticEffectKey("attention_items", "insert")), semanticStage(0, semanticEffectKey("attention_stop_links", "insert")), semanticStage(0, semanticEffectKey("run_runtime", "update")), semanticStage(1, semanticEffectKey("idempotency_results", "update"))),
	"evidence.linked":                semanticRule("evidence", semanticStage(1, semanticEffectKey("evidence_records", "insert")), semanticStage(0, semanticEffectKey("evidence_artifact_links", "insert")), semanticStage(0, semanticEffectKey("stops", "insert"), semanticEffectKey("attention_items", "insert"), semanticEffectKey("attention_stop_links", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"artifact.availability-observed": semanticRule("artifact", semanticStage(0, semanticEffectKey("referenced_artifacts", "update")), semanticStage(0, semanticEffectKey("stops", "insert"), semanticEffectKey("attention_items", "insert"), semanticEffectKey("attention_stop_links", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"stop.raised":                    semanticRule("stop", semanticStage(1, semanticEffectKey("stops", "insert")), semanticStage(0, semanticEffectKey("attention_items", "insert")), semanticStage(1, semanticEffectKey("attention_stop_links", "insert")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"stop.resolved":                  semanticRule("stop", semanticStage(1, semanticEffectKey("stop_resolutions", "insert")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"attention.resolved":             semanticRule("attention", semanticStage(1, semanticEffectKey("attention_items", "update")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"issue.created":                  semanticRule("issue", semanticStage(1, semanticEffectKey("issues", "insert")), semanticStage(0, semanticEffectKey("issue_labels", "insert")), semanticStage(0, semanticEffectKey("epic_metadata", "insert")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.updated":                  semanticRule("issue", semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(0, semanticEffectKey("issue_labels", "insert"), semanticEffectKey("issue_labels", "delete")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.transitioned":             semanticRule("issue", semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.ac-added":                 semanticRule("issue", semanticStage(1, semanticEffectKey("acceptance_criteria", "insert")), semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.label-added":              semanticRule("issue", semanticStage(1, semanticEffectKey("issue_labels", "insert")), semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"issue.label-removed":            semanticRule("issue", semanticStage(1, semanticEffectKey("issue_labels", "delete")), semanticStage(1, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"label.added":                    semanticRule("label", semanticStage(1, semanticEffectKey("labels", "insert")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
	"link.added":                     semanticRule("link", semanticStage(1, semanticEffectKey("links", "insert")), semanticStage(0, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"link.removed":                   semanticRule("link", semanticStage(1, semanticEffectKey("links", "delete")), semanticStage(0, semanticEffectKey("issues", "update")), semanticStage(1, semanticEffectKey("activity_log", "activity")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"workflow.defined":               semanticRule("workflow", semanticStage(1, semanticEffectKey("workflows", "insert")), semanticStage(1, semanticEffectKey("workflow_states", "insert")), semanticStage(0, semanticEffectKey("workflow_transitions", "insert")), semanticStage(0, semanticEffectKey("idempotency_results", "insert"))),
	"workflow.assigned":              semanticRule("workflow-assignment", semanticStage(0, semanticEffectKey("workflow_assignments", "delete")), semanticStage(1, semanticEffectKey("workflow_assignments", "insert")), semanticStage(1, semanticEffectKey("idempotency_results", "insert"))),
}

var semanticPost14Events = map[string]bool{
	"activation.resource-observed": true, "activation.recovery-selected": true, "activation.recovery-completed": true,
	"claim.expired": true, "run.observed": true, "run.checkpointed": true, "run.recovered": true, "run.budget-observed": true, "run.execution-blocked": true, "run.closed": true, "run.abandoned": true,
	"stop.raised": true, "stop.resolved": true, "attention.resolved": true,
	"review.recorded": true, "finding.addressed": true, "concern.recorded": true,
	"concern.disposed": true, "deviation.recorded": true, "decision.referenced": true,
	"human-outcome.recorded": true, "requirement.referenced": true, "requirement.linked": true,
}

func init() {
	for eventType, rule := range semanticEventRules {
		if rule.aggregate == "" {
			panic("semantic archive validation has no aggregate for " + eventType)
		}
		for _, stage := range rule.stages {
			for operation := range stage.allowed {
				if len(strings.Split(operation, "\x00")) != 2 {
					panic("semantic archive validation has malformed Effect operation")
				}
			}
		}
	}
}
