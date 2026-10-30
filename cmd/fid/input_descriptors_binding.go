package main

import (
	"reflect"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func bindingArray(name, meaning string, element cli.InputField, minimum int) cli.InputField {
	maximum := 10000
	field := cli.InputField{Name: name, Type: "array", Requirement: "required", Meaning: meaning, Element: &element, MaxLen: &maximum}
	if minimum > 0 {
		field.MinLen = &minimum
	}
	return field
}

func bindingString(name, meaning string) cli.InputField {
	return inputString(name, "required", meaning, "text", "<"+name+">")
}

func bindingNullable(name, meaning string) cli.InputField {
	return inputString(name, "nullable", meaning, "text", "")
}

func bindingObligation() cli.InputField {
	registration := inputObject("registration", "Gate Registration used by gate-evidence Obligations.",
		judgmentID("id", "Registration ID."), inputInteger("version", "Registration version.", 1))
	registration.Requirement = "nullable"
	subject := inputObject("subject", "Issue guarded by this Obligation.", bindingString("type", "Subject type (issue)."), judgmentID("id", "Issue ID."))
	reviewProfile := bindingNullable("review_profile", "Review profile, if applicable.")
	reviewProfile.Requirement = "optional"
	outcomeClass := bindingNullable("outcome_class", "Human Outcome class, if applicable.")
	outcomeClass.Requirement = "optional"
	return inputObject("obligation", "One required Obligation.",
		judgmentID("id", "Obligation ID."), bindingString("type", "Obligation type."), bindingString("applicability", "Applicability (required)."),
		inputArray("guarded_actions", "Actions guarded by this Obligation.", "text", true), subject, registration,
		inputArray("accepted_evidence_types", "Accepted Evidence types.", "text", true), inputArray("accepted_trust_sources", "Accepted trust sources.", "text", true),
		bindingString("freshness_target", "Freshness target."), bindingString("artifact_policy", "Artifact policy."), bindingString("independence", "Independence rule."),
		bindingNullable("required_lane", "Required lane, if any."), reviewProfile,
		outcomeClass, judgmentSource("source"))
}

func bindingBudget() cli.InputField {
	warning := inputInteger("warning_threshold", "Optional warning threshold.", 0)
	warning.Requirement = "nullable"
	consumed := inputInteger("consumed", "Reserved null counter.", 0)
	consumed.Requirement = "nullable"
	remaining := inputInteger("remaining", "Reserved null counter.", 0)
	remaining.Requirement = "nullable"
	state := bindingString("state", "Output-only budget state.")
	state.Requirement = "optional"
	return inputObject("budget", "One immutable Budget.", judgmentID("name", "Budget name."), bindingString("scope", "Budget scope."), bindingString("unit", "Budget unit."),
		inputInteger("maximum", "Maximum amount.", 1), warning, bindingString("exhaustion_action", "Action on exhaustion."), consumed, remaining, state)
}

func bindingDescriptors() map[string]cli.InputDescriptor {
	source := judgmentSource("source")
	actions := inputArray("actions", "Sorted authorized actions.", "text", true)
	selector := inputObject("actor_selector", "Selected Actor and capabilities.", judgmentID("actor_id", "Actor ID."), inputArray("required_capabilities", "Required capabilities.", "ID", false))
	grant := inputObject("grant", "Authority grant.", actions, selector, inputObject("scope", "Issue scope.", judgmentID("issue_id", "Issue ID.")))
	decision := inputObject("decision", "Protected Decision rule.", inputInteger("ordinal", "Zero-based ordinal.", 0), bindingString("decision_class", "Decision class."),
		inputArray("required_capabilities", "Capabilities required.", "ID", false), inputArray("affected_actions", "Actions affected.", "text", false))
	decisions := bindingArray("protected_decisions", "Protected Decision rules.", decision, 0)
	decisions.Requirement = "optional"
	paths := inputArray("paths", "Sorted paths in Claim scope.", "text", true)
	paths.Format = "sorted-unique"
	mode := bindingString("mode", "Claim mode: mutation or review; omitted defaults to mutation.")
	mode.Requirement = "optional"
	mode.Members = []string{"mutation", "review"}
	return map[string]cli.InputDescriptor{
		"gate register": {Operation: "gate.register", Purpose: "Register a Gate against the observed worktree.", Type: reflect.TypeFor[evidence.RegisterRequest](), Fields: []cli.InputField{
			judgmentID("definition_id", "Activated Gate Definition ID."), inputInteger("definition_version", "Definition version.", 1), inputObject("scope", "Issue in scope.", judgmentID("issue_id", "Issue ID."))}},
		"envelope create": {Operation: "envelope.create", Purpose: "Create an Envelope using the observed repository and checked Sources.", Type: reflect.TypeFor[workflow.EnvelopeCreateRequest](), Fields: []cli.InputField{
			bindingString("outcome", "Intended outcome."), bindingArray("non_goals", "Explicit non-goals.", bindingString("non_goal", "Non-goal."), 0),
			inputObject("risk_result", "Risk assessment Source and digest.", source, inputString("result_digest", "required", "Risk result digest.", "digest", "<sha256-digest>")),
			bindingArray("sources", "Sources supporting the Envelope.", source, 1), judgmentID("primary_issue_id", "Primary Issue ID."),
			inputArray("related_issue_ids", "Related Issues (empty in this revision).", "ID", false), bindingArray("authority_grants", "Authority grants.", grant, 1),
			bindingArray("obligations", "Required Obligations.", bindingObligation(), 1), decisions,
			bindingArray("budgets", "Immutable Budget declarations.", bindingBudget(), 0)}},
		"envelope activate": {Operation: "envelope.activate", Purpose: "Activate a draft Envelope.", Type: reflect.TypeFor[workflow.EnvelopeActivateRequest](), Fields: []cli.InputField{
			judgmentID("envelope_id", "Envelope ID."), inputInteger("expected_revision", "Expected revision.", 1)}},
		"envelope complete": {Operation: "envelope.complete", Purpose: "Complete an eligible Envelope against the observed repository.", Type: reflect.TypeFor[workflow.EnvelopeCompleteRequest](), Fields: []cli.InputField{
			judgmentID("envelope_id", "Envelope ID."), inputInteger("expected_revision", "Expected revision.", 1)}},
		"claim acquire": {Operation: "claim.acquire", Purpose: "Acquire one Claim. fid observes the selected worktree.", Type: reflect.TypeFor[claims.AcquireRequest](), Fields: []cli.InputField{
			judgmentID("issue_id", "Issue ID."), judgmentID("envelope_id", "Envelope ID."), inputInteger("envelope_revision", "Envelope revision.", 1),
			inputObject("scope", "Claim path scope.", paths), mode, inputInteger("ttl_seconds", "Claim lifetime in seconds (1..86400).", 1)}},
		"claim renew": {Operation: "claim.renew", Purpose: "Renew an active Claim.", Type: reflect.TypeFor[claims.RenewRequest](), Fields: []cli.InputField{
			judgmentID("claim_id", "Claim ID."), inputInteger("expected_version", "Expected Claim version.", 1), inputInteger("ttl_seconds", "New lifetime in seconds.", 1)}},
		"claim release": {Operation: "claim.release", Purpose: "Release a Claim.", Type: reflect.TypeFor[claims.ReleaseRequest](), Fields: []cli.InputField{
			judgmentID("claim_id", "Claim ID."), inputInteger("expected_version", "Expected Claim version.", 1), bindingString("reason", "Release reason.")}},
		"claim expire": {Operation: "claim.expire", Purpose: "Expire a Claim.", Type: reflect.TypeFor[claims.ExpireRequest](), Fields: []cli.InputField{
			judgmentID("claim_id", "Claim ID."), inputInteger("expected_version", "Expected Claim version.", 1)}},
	}
}

func init() {
	for leaf, descriptor := range bindingDescriptors() {
		inputDescriptorRegistry[leaf] = descriptor
	}
}
