package journal

// Current families are the Phase 22 schema section 5 replacements/extensions
// and their incorporated activation lifecycle. Historical archive grammar is separate.
type currentEventRule struct {
	aggregate string
	stages    []currentEffectStage
}
type currentEffectStage struct {
	minimum int
	allowed map[string]bool
}

func currentRule(aggregate string, stages ...currentEffectStage) currentEventRule {
	return currentEventRule{aggregate, stages}
}
func currentStage(minimum int, operations ...string) currentEffectStage {
	allowed := map[string]bool{}
	for _, op := range operations {
		allowed[op] = true
	}
	return currentEffectStage{minimum, allowed}
}
func currentEffectKey(table, op string) string { return table + "/" + op }

var currentEventRules = map[string]currentEventRule{
	"review.recorded":               currentRule("review", currentStage(1, currentEffectKey("review_records", "insert")), currentStage(0, currentEffectKey("review_findings", "insert"), currentEffectKey("finding_dispositions", "insert"), currentEffectKey("budget_samples", "insert"), currentEffectKey("stops", "insert"), currentEffectKey("attention_items", "insert"), currentEffectKey("attention_stop_links", "insert"), currentEffectKey("stop_resolutions", "insert"), currentEffectKey("attention_items", "update"), currentEffectKey("activity_log", "activity"), currentEffectKey("idempotency_results", "insert"))),
	"finding.addressed":             currentRule("finding", currentStage(1, currentEffectKey("finding_dispositions", "insert")), currentStage(0, currentEffectKey("activity_log", "activity"), currentEffectKey("idempotency_results", "insert"))),
	"concern.recorded":              currentRule("concern", currentStage(1, currentEffectKey("concerns", "insert")), currentStage(0, currentEffectKey("activity_log", "activity"), currentEffectKey("idempotency_results", "insert"))),
	"concern.disposed":              currentRule("concern", currentStage(1, currentEffectKey("concerns", "update")), currentStage(0, currentEffectKey("activity_log", "activity"), currentEffectKey("idempotency_results", "insert"))),
	"deviation.recorded":            currentRule("deviation", currentStage(1, currentEffectKey("deviations", "insert")), currentStage(0, currentEffectKey("activity_log", "activity"), currentEffectKey("idempotency_results", "insert"))),
	"decision.referenced":           currentRule("decision", currentStage(1, currentEffectKey("decision_refs", "insert")), currentStage(0, currentEffectKey("stop_resolutions", "insert"), currentEffectKey("attention_items", "update"), currentEffectKey("activity_log", "activity"), currentEffectKey("idempotency_results", "insert"))),
	"human-outcome.recorded":        currentRule("human-outcome", currentStage(1, currentEffectKey("human_outcomes", "insert")), currentStage(0, currentEffectKey("stop_resolutions", "insert"), currentEffectKey("attention_items", "update"), currentEffectKey("activity_log", "activity"), currentEffectKey("idempotency_results", "insert"))),
	"requirement.referenced":        currentRule("requirement", currentStage(0, currentEffectKey("requirement_refs", "update")), currentStage(1, currentEffectKey("requirement_refs", "insert")), currentStage(0, currentEffectKey("activity_log", "activity"), currentEffectKey("idempotency_results", "insert"))),
	"requirement.linked":            currentRule("requirement", currentStage(1, currentEffectKey("requirement_issue_links", "insert")), currentStage(0, currentEffectKey("activity_log", "activity"), currentEffectKey("idempotency_results", "insert"))),
	"activation.prepared":           currentRule("activation", currentStage(1, currentEffectKey("kit_activations", "insert")), currentStage(1, currentEffectKey("activation_resources", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"activation.component-staged":   currentRule("activation", currentStage(1, currentEffectKey("kit_activations", "update")), currentStage(1, currentEffectKey("activation_components", "insert")), currentStage(1, currentEffectKey("activation_actor_profiles", "insert"), currentEffectKey("activation_gate_definitions", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"activation.committed":          currentRule("activation", currentStage(1, currentEffectKey("kit_activations", "update")), currentStage(0, currentEffectKey("activation_components", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"activation.finalized":          currentRule("activation", currentStage(1, currentEffectKey("kit_activations", "update")), currentStage(0, currentEffectKey("activation_components", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"activation.aborted":            currentRule("activation", currentStage(1, currentEffectKey("kit_activations", "update")), currentStage(0, currentEffectKey("activation_components", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"activation.resource-observed":  currentRule("activation", currentStage(1, currentEffectKey("kit_activations", "update")), currentStage(1, currentEffectKey("activation_resources", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"activation.recovery-selected":  currentRule("activation", currentStage(1, currentEffectKey("kit_activations", "update")), currentStage(1, currentEffectKey("activation_recoveries", "insert")), currentStage(0, currentEffectKey("stops", "insert")), currentStage(0, currentEffectKey("attention_items", "insert")), currentStage(0, currentEffectKey("attention_stop_links", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"activation.recovery-completed": currentRule("activation", currentStage(1, currentEffectKey("kit_activations", "update")), currentStage(1, currentEffectKey("activation_recoveries", "update")), currentStage(0, currentEffectKey("activation_components", "update")), currentStage(0, currentEffectKey("stop_resolutions", "insert")), currentStage(0, currentEffectKey("attention_items", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"claim.expired":                 currentRule("claim", currentStage(1, currentEffectKey("claims", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"run.started":                   currentRule("run", currentStage(1, currentEffectKey("runs", "insert")), currentStage(1, currentEffectKey("run_runtime", "insert")), currentStage(0, currentEffectKey("run_budget_bindings", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"run.observed":                  currentRule("run", currentStage(1, currentEffectKey("run_observations", "insert")), currentStage(1, currentEffectKey("run_runtime", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"run.checkpointed":              currentRule("run", currentStage(1, currentEffectKey("run_observations", "insert")), currentStage(1, currentEffectKey("run_runtime", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"run.recovered":                 currentRule("run", currentStage(0, currentEffectKey("runs", "update")), currentStage(0, currentEffectKey("run_runtime", "update")), currentStage(0, currentEffectKey("claims", "update")), currentStage(1, currentEffectKey("claims", "insert")), currentStage(1, currentEffectKey("runs", "insert")), currentStage(1, currentEffectKey("run_runtime", "insert")), currentStage(0, currentEffectKey("run_budget_bindings", "insert")), currentStage(1, currentEffectKey("run_recoveries", "insert")), currentStage(0, currentEffectKey("budget_samples", "insert")), currentStage(0, currentEffectKey("stops", "insert")), currentStage(0, currentEffectKey("attention_items", "insert")), currentStage(0, currentEffectKey("attention_stop_links", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"run.budget-observed":           currentRule("run", currentStage(1, currentEffectKey("budget_samples", "insert")), currentStage(0, currentEffectKey("stops", "insert")), currentStage(0, currentEffectKey("attention_items", "insert")), currentStage(0, currentEffectKey("attention_stop_links", "insert")), currentStage(1, currentEffectKey("run_runtime", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"run.execution-blocked":         currentRule("run", currentStage(1, currentEffectKey("stops", "insert")), currentStage(0, currentEffectKey("attention_items", "insert")), currentStage(0, currentEffectKey("attention_stop_links", "insert")), currentStage(1, currentEffectKey("run_runtime", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"run.closed":                    currentRule("run", currentStage(1, currentEffectKey("runs", "update")), currentStage(1, currentEffectKey("run_runtime", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"run.abandoned":                 currentRule("run", currentStage(1, currentEffectKey("runs", "update")), currentStage(1, currentEffectKey("run_runtime", "update")), currentStage(1, currentEffectKey("claims", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"gate.execution-started":        currentRule("attempt", currentStage(1, currentEffectKey("execution_attempts", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert")), currentStage(0, currentEffectKey("budget_samples", "insert")), currentStage(0, currentEffectKey("stops", "insert")), currentStage(0, currentEffectKey("attention_items", "insert")), currentStage(0, currentEffectKey("attention_stop_links", "insert")), currentStage(0, currentEffectKey("run_runtime", "update"))),
	"gate.receipt-recorded":         currentRule("receipt", currentStage(0, currentEffectKey("referenced_artifacts", "insert")), currentStage(0, currentEffectKey("budget_samples", "insert")), currentStage(0, currentEffectKey("stops", "insert")), currentStage(0, currentEffectKey("attention_items", "insert")), currentStage(0, currentEffectKey("attention_stop_links", "insert")), currentStage(0, currentEffectKey("run_runtime", "update")), currentStage(1, currentEffectKey("execution_receipts", "insert")), currentStage(0, currentEffectKey("receipt_artifacts", "insert")), currentStage(1, currentEffectKey("idempotency_results", "update"))),
	"stop.raised":                   currentRule("stop", currentStage(1, currentEffectKey("stops", "insert")), currentStage(0, currentEffectKey("attention_items", "insert")), currentStage(1, currentEffectKey("attention_stop_links", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"stop.resolved":                 currentRule("stop", currentStage(1, currentEffectKey("stop_resolutions", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"attention.resolved":            currentRule("attention", currentStage(1, currentEffectKey("attention_items", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
}
var currentEventCommands = map[string]string{
	"review.recorded":               "review.record",
	"finding.addressed":             "finding.address",
	"concern.recorded":              "concern.record",
	"concern.disposed":              "concern.dispose",
	"deviation.recorded":            "deviation.record",
	"decision.referenced":           "decision.record-reference",
	"human-outcome.recorded":        "human-outcome.record",
	"requirement.referenced":        "requirement.record-reference",
	"requirement.linked":            "requirement.link",
	"run.started":                   "run.start",
	"claim.expired":                 "claim.expire",
	"run.observed":                  "run.observe",
	"run.checkpointed":              "run.checkpoint",
	"run.recovered":                 "run.recover",
	"stop.raised":                   "stop.raise",
	"stop.resolved":                 "stop.resolve",
	"attention.resolved":            "attention.resolve",
	"run.budget-observed":           "run.observe-budget",
	"gate.execution-started":        "gate.run",
	"gate.receipt-recorded":         "gate.run",
	"run.execution-blocked":         "gate.run",
	"run.closed":                    "run.close",
	"run.abandoned":                 "run.abandon",
	"activation.prepared":           "activation.prepare",
	"activation.component-staged":   "activation.component",
	"activation.committed":          "activation.commit",
	"activation.finalized":          "activation.finalize",
	"activation.aborted":            "activation.abort",
	"activation.resource-observed":  "activation.resource",
	"activation.recovery-selected":  "activation.recover",
	"activation.recovery-completed": "activation.recover-complete",
}

// Retained current Event-v2 families keep the accepted Phase 21 Effect
// grammar. They cannot carry a Phase 22 Effect under an unrelated family.
// The lean core's families (ADR-0014) share this grammar, and their Effects
// are checked against their stages in order too (leanEventFamilies).
var retainedCurrentRules = map[string]currentEventRule{
	"task.started":                   currentRule("task", currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"policy.snapshotted":             currentRule("policy", currentStage(1, currentEffectKey("policy_snapshots", "insert"))),
	"hold.heartbeat":                 currentRule("hold", currentStage(1, currentEffectKey("holds", "update"))),
	"hold.waiting":                   currentRule("hold", currentStage(1, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("questions", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"task.answered":                  currentRule("task", currentStage(1, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("questions", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"hold.parked":                    currentRule("hold", currentStage(1, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"task.taken-over":                currentRule("task", currentStage(0, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("holds", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"hold.taken":                     currentRule("hold", currentStage(0, currentEffectKey("holds", "update")), currentStage(1, currentEffectKey("holds", "insert"))),
	"issue.label.added":              currentRule("issue", currentStage(1, currentEffectKey("issue_labels", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.label.removed":            currentRule("issue", currentStage(1, currentEffectKey("issue_labels", "delete")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"project.initialized":            currentRule("project"),
	"gate.registered":                currentRule("gate-registration", currentStage(1, currentEffectKey("gate_registrations", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"envelope.created":               currentRule("envelope", currentStage(1, currentEffectKey("work_envelopes", "insert")), currentStage(1, currentEffectKey("envelope_revisions", "insert")), currentStage(1, currentEffectKey("envelope_issue_selection", "insert")), currentStage(0, currentEffectKey("envelope_sources", "insert")), currentStage(0, currentEffectKey("envelope_authority_grants", "insert")), currentStage(1, currentEffectKey("envelope_obligations", "insert")), currentStage(0, currentEffectKey("envelope_protected_decisions", "insert"), currentEffectKey("envelope_budgets", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"envelope.activated":             currentRule("envelope", currentStage(1, currentEffectKey("work_envelopes", "update")), currentStage(1, currentEffectKey("envelope_revisions", "update")), currentStage(0, currentEffectKey("stops", "insert"), currentEffectKey("attention_items", "insert"), currentEffectKey("attention_stop_links", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"envelope.completed":             currentRule("envelope", currentStage(1, currentEffectKey("work_envelopes", "update")), currentStage(1, currentEffectKey("envelope_revisions", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"claim.acquired":                 currentRule("claim", currentStage(1, currentEffectKey("claims", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"claim.renewed":                  currentRule("claim", currentStage(1, currentEffectKey("claims", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"claim.released":                 currentRule("claim", currentStage(1, currentEffectKey("claims", "update")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"evidence.linked":                currentRule("evidence", currentStage(1, currentEffectKey("evidence_records", "insert")), currentStage(0, currentEffectKey("evidence_artifact_links", "insert")), currentStage(0, currentEffectKey("stops", "insert"), currentEffectKey("attention_items", "insert"), currentEffectKey("attention_stop_links", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"artifact.availability-observed": currentRule("artifact", currentStage(0, currentEffectKey("referenced_artifacts", "update")), currentStage(0, currentEffectKey("stops", "insert"), currentEffectKey("attention_items", "insert"), currentEffectKey("attention_stop_links", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.created":                  currentRule("issue", currentStage(1, currentEffectKey("issues", "insert")), currentStage(0, currentEffectKey("issue_labels", "insert")), currentStage(0, currentEffectKey("epic_metadata", "insert")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.updated":                  currentRule("issue", currentStage(1, currentEffectKey("issues", "update")), currentStage(0, currentEffectKey("issue_labels", "insert"), currentEffectKey("issue_labels", "delete")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.transitioned":             currentRule("issue", currentStage(1, currentEffectKey("issues", "update")), currentStage(0, currentEffectKey("holds", "update")), currentStage(0, currentEffectKey("questions", "update")), currentStage(0, currentEffectKey("person_decisions", "insert")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.ac-added":                 currentRule("issue", currentStage(1, currentEffectKey("acceptance_criteria", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.label-added":              currentRule("issue", currentStage(1, currentEffectKey("issue_labels", "insert")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"issue.label-removed":            currentRule("issue", currentStage(1, currentEffectKey("issue_labels", "delete")), currentStage(1, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"label.added":                    currentRule("label", currentStage(1, currentEffectKey("labels", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
	"link.added":                     currentRule("link", currentStage(1, currentEffectKey("links", "insert")), currentStage(0, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"link.removed":                   currentRule("link", currentStage(1, currentEffectKey("links", "delete")), currentStage(0, currentEffectKey("issues", "update")), currentStage(1, currentEffectKey("activity_log", "activity")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"workflow.defined":               currentRule("workflow", currentStage(1, currentEffectKey("workflows", "insert")), currentStage(1, currentEffectKey("workflow_states", "insert")), currentStage(0, currentEffectKey("workflow_transitions", "insert")), currentStage(0, currentEffectKey("idempotency_results", "insert"))),
	"workflow.assigned":              currentRule("workflow-assignment", currentStage(0, currentEffectKey("workflow_assignments", "delete")), currentStage(1, currentEffectKey("workflow_assignments", "insert")), currentStage(1, currentEffectKey("idempotency_results", "insert"))),
}

// leanEventFamilies are the lean core's families. Every Effect they carry
// is checked against their stages, in order, not only those on a domain
// table another family owns.
var leanEventFamilies = map[string]bool{"task.started": true, "hold.waiting": true, "task.answered": true, "hold.heartbeat": true, "hold.parked": true, "task.taken-over": true, "hold.taken": true, "policy.snapshotted": true}
