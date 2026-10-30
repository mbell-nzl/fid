package main

import (
	"reflect"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func inputString(name, requirement, meaning, format, placeholder string) cli.InputField {
	return cli.InputField{Name: name, Type: "string", Requirement: requirement, Meaning: meaning, Format: format, Placeholder: placeholder}
}

func inputInteger(name, meaning string, min int64) cli.InputField {
	return cli.InputField{Name: name, Type: "integer", Requirement: "required", Meaning: meaning, MinInt: &min}
}

func inputObject(name, meaning string, children ...cli.InputField) cli.InputField {
	return cli.InputField{Name: name, Type: "object", Requirement: "required", Meaning: meaning, Children: children}
}

func inputArray(name, meaning, elementFormat string, nonempty bool) cli.InputField {
	maximum := 10000
	field := cli.InputField{Name: name, Type: "array", Requirement: "required", Meaning: meaning,
		MaxLen:  &maximum,
		Element: &cli.InputField{Name: name, Type: "string", Requirement: "required", Meaning: "Array member.", Format: elementFormat, Placeholder: "<" + name + ">"}}
	if nonempty {
		minimum := 1
		field.MinLen = &minimum
	}
	return field
}

func inputSubject() cli.InputField {
	kind := inputString("type", "required", "Subject kind.", "", "")
	kind.Members = []string{"run", "envelope", "activation"}
	return inputObject("subject", "The blocked subject.", kind, inputString("id", "required", "Subject ID.", "ID", "<subject-id>"))
}

func inputCondition() cli.InputField {
	kind := inputString("type", "required", "Qualifying fact type.", "", "")
	kind.Members = []string{"executor-terminated", "run-recovered", "activation-recovered", "capability-available", "budget-observed", "work-ended", "decision-satisfied", "human-outcome-recorded", "review-recorded", "finding-resolved"}
	condition := inputObject("resolution_condition", "Durable condition required to resolve the Stop.", kind, inputSubject())
	condition.Children = append(condition.Children,
		inputString("budget_name", "nullable", "Budget name when budget-observed.", "ID", "<budget-id>"),
		inputString("capability_id", "nullable", "Capability when capability-available.", "ID", "<capability-id>"),
		inputString("decision_class", "optional", "Decision class when decision-satisfied.", "ID", "<decision-class>"),
		inputString("outcome_class", "optional", "Human Outcome class when human-outcome-recorded.", "ID", "<outcome-class>"),
		inputString("finding_id", "optional", "Finding when finding-resolved.", "ID", "<finding-id>"))
	return condition
}

func inputDetails() cli.InputField {
	id := func(name string) cli.InputField {
		return inputString(name, "required", "Identity of "+name+".", "ID", "<"+strings.ReplaceAll(name, "_", "-")+">")
	}
	text := func(name string) cli.InputField {
		return inputString(name, "required", "The "+name+" text.", "text", "<"+strings.ReplaceAll(name, "_", "-")+">")
	}
	nullID := inputString("observation_id", "nullable", "Observation, if available.", "ID", "<observation-id>")
	nullInt := func(name string) cli.InputField {
		field := inputInteger(name, "Observed "+name+" when available.", 0)
		field.Requirement = "nullable"
		return field
	}
	boolean := func(name, meaning string) cli.InputField {
		return cli.InputField{Name: name, Type: "boolean", Requirement: "required", Meaning: meaning}
	}
	unit := text("unit")
	maximum := inputInteger("maximum", "Budget maximum.", 0)
	budget := []cli.InputField{id("envelope_id"), id("budget_name"), unit, maximum, nullInt("consumed"), nullInt("remaining")}
	budgetReview := []cli.InputField{id("envelope_id"), id("budget_name"), unit, maximum,
		inputInteger("consumed", "Consumed cycles.", 0), inputInteger("remaining", "Remaining cycles.", 0), id("last_review_id")}
	available := boolean("available", "Whether the capability is available (false).")
	available.Members = []string{"false"}
	owner := inputInteger("owner_phase", "Phase responsible for the capability.", 0)
	return cli.InputField{Name: "details", Type: "object", Requirement: "required", Meaning: "Code-specific Stop facts.", Variants: map[string][]cli.InputField{
		"executor-liveness-unknown":   {id("run_id"), id("executor_id"), nullID},
		"recovery-state-conflict":     {id("run_id"), inputString("observation_digest", "required", "Observation digest.", "digest", "<digest>"), inputArray("paths", "Conflicting paths.", "text", false), inputArray("reasons", "Conflict reasons.", "text", false)},
		"budget-exhausted":            budget,
		"budget-data-unavailable":     budget,
		"review-rework-exhausted":     budgetReview,
		"protected-decision-required": {id("envelope_id"), id("decision_class"), inputArray("required_capabilities", "Capabilities needed.", "ID", false), inputArray("affected_actions", "Actions blocked.", "text", false)},
		"human-outcome-required":      {id("envelope_id"), id("obligation_id"), id("outcome_class")},
		"activation-incomplete":       {id("activation_id"), boolean("point_of_no_return", "Whether activation passed the irreversible point."), inputArray("components", "Incomplete components.", "text", false)},
		"capability-unavailable":      {id("capability_id"), available, owner},
		"operator-stop":               {text("reason"), text("reference")},
	}}
}

func stopRaiseDescriptor() cli.InputDescriptor {
	code := inputString("code", "required", "Stop reason code.", "", "")
	code.Members = []string{"executor-liveness-unknown", "recovery-state-conflict", "budget-exhausted", "budget-data-unavailable", "review-rework-exhausted", "protected-decision-required", "human-outcome-required", "activation-incomplete", "capability-unavailable", "operator-stop"}
	actions := inputArray("blocked_actions", "Sorted unique actions blocked by this Stop.", "text", true)
	actions.Format = "sorted-unique"
	actions.Element.Members = []string{"claim.acquire", "run.start", "run.recover", "gate.run", "issue.transition-terminal-success", "envelope.complete"}
	return cli.InputDescriptor{Operation: "stop.raise", Purpose: "Raise a durable Stop for a subject.", Type: reflect.TypeFor[workflow.StopRaiseRequest](),
		Fields: []cli.InputField{inputSubject(), code, actions, inputDetails(), inputCondition(), inputString("resolver_capability", "required", "Capability required to resolve the Stop.", "ID", "<capability-id>")}}
}

func stopResolveDescriptor() cli.InputDescriptor {
	return cli.InputDescriptor{Operation: "stop.resolve", Purpose: "Resolve a Stop with a qualifying fact Event.", Type: reflect.TypeFor[workflow.StopResolveRequest](),
		Fields: []cli.InputField{inputString("stop_id", "required", "Stop to resolve.", "ID", "<stop-id>"), inputString("fact_event_id", "required", "Qualifying fact Event.", "ID", "<fact-event-id>")}}
}

func attentionResolveDescriptor() cli.InputDescriptor {
	ids := inputArray("fact_event_ids", "Sorted unique disposition Event IDs for linked Stops.", "ID", false)
	ids.Format = "sorted-unique"
	return cli.InputDescriptor{Operation: "attention.resolve", Purpose: "Resolve an Attention item after its linked Stops have dispositions.", Type: reflect.TypeFor[workflow.AttentionResolveRequest](),
		Fields: []cli.InputField{inputString("attention_id", "required", "Attention item to resolve.", "ID", "<attention-id>"), inputInteger("expected_version", "Expected Attention version.", 1), ids}}
}

var inputDescriptorRegistry = map[string]cli.InputDescriptor{
	"stop raise":        stopRaiseDescriptor(),
	"stop resolve":      stopResolveDescriptor(),
	"attention resolve": attentionResolveDescriptor(),
}
