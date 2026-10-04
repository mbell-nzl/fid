package main

import (
	"reflect"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/review"
)

func judgmentID(name, meaning string) cli.InputField {
	return inputString(name, "required", meaning, "ID", "<"+name+">")
}

func judgmentText(name, meaning string) cli.InputField {
	return inputString(name, "required", meaning, "text", "<"+name+">")
}

func judgmentNullable(name, meaning string) cli.InputField {
	return inputString(name, "nullable", meaning, "text", "")
}

func judgmentSource(name string) cli.InputField {
	kind := inputString("source_type", "required", "Project file or external reference.", "", "")
	kind.Members = []string{"project-file", "external-reference"}
	label := judgmentText("display_label", "Human-readable label for this source.")
	minimum, maximum := 1, 128
	label.MinLen, label.MaxLen = &minimum, &maximum
	digest := inputString("digest", "optional", "SHA-256 digest; required for external references and optional for project files.", "digest", "<sha256-digest>")
	return inputObject(name, "SourceInput checked against a project file or recorded as an external reference.",
		kind, label, judgmentNullable("path", "Relative project-file path under the project root."),
		judgmentNullable("section", "Section within the referenced material, if any."),
		judgmentNullable("external_id", "External reference ID, if any."), digest)
}

func judgmentReferences() cli.InputField {
	maximum := 32
	source := judgmentSource("reference")
	return cli.InputField{Name: "references", Type: "array", Requirement: "required", Meaning: "Sources supporting this judgment.",
		MaxLen: &maximum, Element: &source}
}

func judgmentSubject() cli.InputField {
	kind := judgmentText("type", "Subject kind.")
	kind.Members = []string{"issue", "envelope", "run", "review", "finding", "requirement"}
	return inputObject("subject", "Subject of this judgment.", kind, judgmentID("id", "Subject ID."))
}

func judgmentScope() cli.InputField {
	paths := inputArray("paths", "Sorted paths in scope.", "text", false)
	paths.Format = "sorted-unique"
	return inputObject("scope", "Paths in scope.", paths)
}

func judgmentDecisionScope() cli.InputField {
	return inputObject("scope", "Issue in scope.", judgmentID("issue_id", "Issue ID."))
}

func judgmentDescriptors() map[string]cli.InputDescriptor {
	severity := judgmentText("severity", "Concern severity.")
	severity.Members = []string{"blocker", "major", "minor", "nit"}
	result := judgmentText("result", "Human judgment result.")
	result.Members = []string{"approved", "rejected", "changes-required", "acknowledged"}
	outcome := judgmentText("outcome", "Decision outcome.")
	outcome.Members = []string{"approved", "rejected", "deferred"}
	status := judgmentText("status", "Requirement Reference status.")
	status.Members = []string{"active", "withdrawn"}
	return map[string]cli.InputDescriptor{
		"decision record-reference": {Operation: "decision.record-reference", Purpose: "Record a Decision Reference to a checked project file.", Type: reflect.TypeFor[review.DecisionReferenceRequest](), Fields: []cli.InputField{
			judgmentID("id", "Decision Reference ID."), judgmentText("class", "Decision class."), judgmentSource("source"),
			inputObject("authority", "Actor and capabilities authorizing this decision.", judgmentID("actor_id", "Actor ID."), inputArray("required_capabilities", "Sorted required capabilities.", "ID", true)),
			outcome, judgmentDecisionScope()}},
		"requirement record-reference": {Operation: "requirement.record-reference", Purpose: "Record a Requirement Reference to a checked project file.", Type: reflect.TypeFor[review.RequirementReferenceRequest](), Fields: []cli.InputField{
			judgmentID("id", "Requirement ID."), judgmentSource("source"), status}},
		"requirement link": {Operation: "requirement.link", Purpose: "Link a Requirement to an Issue.", Type: reflect.TypeFor[review.RequirementLinkRequest](), Fields: []cli.InputField{
			judgmentID("requirement_id", "Requirement ID."), judgmentID("issue_id", "Issue ID."), judgmentText("relation", "Link relation.")}},
		"human-outcome record": {Operation: "human-outcome.record", Purpose: "Record an authorized Human Outcome with checked references.", Type: reflect.TypeFor[review.HumanOutcomeRequest](), Fields: []cli.InputField{
			judgmentID("id", "Human Outcome ID."), judgmentText("class", "Outcome class."), judgmentSubject(), result,
			judgmentDecisionScope(), judgmentText("reasoning", "Recorded Human reasoning."), judgmentReferences()}},
		"concern record": {Operation: "concern.record", Purpose: "Record a Concern with checked references.", Type: reflect.TypeFor[review.ConcernRecordRequest](), Fields: []cli.InputField{
			judgmentID("id", "Concern ID."), judgmentSubject(), judgmentScope(), severity, judgmentText("summary", "Concern summary."),
			judgmentText("reasoning", "Recorded reasoning."), judgmentReferences()}},
		"concern dispose": {Operation: "concern.dispose", Purpose: "Dispose a recorded Concern.", Type: reflect.TypeFor[review.ConcernDisposeRequest](), Fields: []cli.InputField{
			judgmentID("concern_id", "Concern ID."), judgmentText("action", "Disposition action."), judgmentNullable("issue_id", "Related Issue ID, if any."),
			judgmentNullable("supersedes_concern_id", "Superseded Concern ID, if any."), judgmentNullable("supersedes_finding_id", "Superseded Finding ID, if any."),
			judgmentText("reason", "Reason for disposition.")}},
		"deviation record": {Operation: "deviation.record", Purpose: "Record a Deviation with checked references.", Type: reflect.TypeFor[review.DeviationRecordRequest](), Fields: []cli.InputField{
			judgmentID("id", "Deviation ID."), judgmentSubject(), judgmentText("expected", "Expected behavior."), judgmentText("observed", "Observed behavior."),
			inputObject("authority", "Authority accepting the disposition.", judgmentID("actor_id", "Actor ID."), judgmentNullable("reference", "Authority reference, if any.")),
			judgmentText("disposition", "Disposition."), judgmentText("reasoning", "Recorded reasoning."), judgmentReferences(),
			judgmentNullable("supersedes_deviation_id", "Superseded Deviation ID, if any.")}},
	}
}

func init() {
	for leaf, descriptor := range judgmentDescriptors() {
		inputDescriptorRegistry[leaf] = descriptor
		delete(pendingInputDescriptors, leaf)
	}
}
