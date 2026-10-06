package main

import (
	"reflect"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/review"
)

func reviewArray(name, meaning string, element cli.InputField, minimum int) cli.InputField {
	maximum := 64
	field := cli.InputField{Name: name, Type: "array", Requirement: "required", Meaning: meaning, MaxLen: &maximum, Element: &element}
	if minimum > 0 {
		field.MinLen = &minimum
	}
	return field
}

func reviewDescriptors() map[string]cli.InputDescriptor {
	boolField := func(name, meaning string) cli.InputField {
		return cli.InputField{Name: name, Type: "boolean", Requirement: "required", Meaning: meaning}
	}
	expected := inputObject("expected_repository", "Repository preconditions checked against one fid observation.",
		inputString("head_oid", "required", "Full output of git rev-parse HEAD.", "oid", "<head-oid>"),
		boolField("dirty", "True when git status --porcelain lists any path outside .fid/."))
	item := inputObject("item", "Covered Source and assessment.", judgmentSource("source"), boolField("assessed", "Whether the Source was assessed."))
	exclusion := inputObject("exclusion", "Excluded Source and reason.", judgmentSource("source"), judgmentText("reason", "Reason for exclusion."))
	coverage := inputObject("coverage", "Declared Review coverage.", reviewArray("items", "Covered Sources, unique by resolved digest.", item, 1), reviewArray("exclusions", "Sources excluded from coverage.", exclusion, 0))
	subjectType := judgmentText("type", "Finding subject type.")
	subjectType.Members = []string{"issue", "envelope", "run", "review", "finding", "requirement"}
	subject := inputObject("subject", "Finding subject.", subjectType, judgmentID("id", "Subject ID."))
	paths := inputArray("paths", "Paths in Finding scope.", "text", false)
	paths.Format = "sorted-unique"
	severity := judgmentText("severity", "Finding severity.")
	severity.Members = []string{"blocker", "major", "minor", "nit"}
	finding := inputObject("finding", "Finding declared by this Review.", judgmentID("id", "Finding ID."), subject,
		inputObject("scope", "Finding scope.", paths), severity, judgmentText("summary", "Finding summary."), boolField("blocking", "Whether this Finding blocks the Review."))
	inapplicable := inputObject("inapplicable", "Finding declared inapplicable.", judgmentID("finding_id", "Finding ID."), judgmentNullable("reason", "Reason it does not apply."))
	inapplicableArray := reviewArray("inapplicable_findings", "Findings declared inapplicable.", inapplicable, 0)
	inapplicableArray.Requirement = "optional"
	verdict := inputString("verdict", "required", "Review verdict.", "", "")
	verdict.Members = []string{"pass", "changes-required", "blocked"}
	return map[string]cli.InputDescriptor{
		"review record": {Operation: "review.record", Purpose: "Record an independent Review against the observed repository.", Type: reflect.TypeFor[review.RecordRequest](), Fields: []cli.InputField{
			judgmentID("run_id", "Review Run ID."), inputInteger("expected_run_version", "Expected Run version.", 1), judgmentText("profile", "Review profile."), expected,
			coverage, verdict,
			reviewArray("findings", "Findings from this Review.", finding, 0), inputArray("resolves_finding_ids", "Findings resolved by this Review.", "ID", false),
			inapplicableArray, reviewArray("evidence_refs", "Supporting Source references.", judgmentSource("source"), 0), inputArray("receipt_ids", "Receipts from the Review Run.", "ID", false)}},
		"finding address": {Operation: "finding.address", Purpose: "Record a Builder response to a Finding.", Type: reflect.TypeFor[review.AddressRequest](), Fields: []cli.InputField{
			judgmentID("finding_id", "Finding to address."), judgmentID("claim_id", "Mutating Claim ID."),
			inputInteger("expected_claim_version", "Expected Claim version.", 1), judgmentText("reason", "Response to the Finding.")}},
	}
}

func init() {
	for leaf, descriptor := range reviewDescriptors() {
		inputDescriptorRegistry[leaf] = descriptor
	}
}
