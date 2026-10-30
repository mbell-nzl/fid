package main

import (
	"reflect"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/evidence"
)

func gateRunDescriptor() cli.InputDescriptor {
	return cli.InputDescriptor{Operation: "gate.run", Purpose: "Execute a registered Gate against one pre-spawn repository observation.", Type: reflect.TypeFor[evidence.GateRunRequest](), Fields: []cli.InputField{
		judgmentID("registration_id", "Registration ID."), inputInteger("registration_version", "Registration version.", 1), judgmentID("obligation_id", "Bound Obligation ID."), judgmentID("run_id", "Active Run ID."), runExpectedRepository(),
	}}
}

func init() {
	inputDescriptorRegistry["gate run"] = gateRunDescriptor()
	kind := bindingString("evidence_type", "Type of the linked subject.")
	kind.Members = []string{"receipt", "review", "human-outcome"}
	nullID := func(name, meaning string) cli.InputField {
		field := judgmentID(name, meaning)
		field.Requirement = "nullable"
		return field
	}
	artifacts := inputArray("required_artifact_ids", "Sorted unique required Artifact IDs.", "ID", false)
	artifacts.Format = "sorted-unique"
	inputDescriptorRegistry["evidence link"] = cli.InputDescriptor{Operation: "evidence.link", Purpose: "Link trusted Evidence, deriving its repository from the subject.", Type: reflect.TypeFor[evidence.EvidenceLinkRequest](), Fields: []cli.InputField{
		judgmentID("obligation_id", "Obligation ID."), kind, nullID("receipt_id", "Receipt ID for receipt Evidence, otherwise null."), nullID("review_id", "Review ID for review Evidence, otherwise null."), nullID("human_outcome_id", "Human Outcome ID for human-outcome Evidence, otherwise null."), artifacts,
	}}
	obligation := nullID("obligation_id", "Evaluate this Obligation; use null to evaluate the active revision.")
	inputDescriptorRegistry["evidence evaluate"] = cli.InputDescriptor{Operation: "evidence.evaluate", Purpose: "Evaluate Evidence freshness against the current repository observation.", Type: reflect.TypeFor[evidence.EvidenceEvaluateRequest](), Fields: []cli.InputField{obligation}}
	inputDescriptorRegistry["artifact verify"] = cli.InputDescriptor{Operation: "artifact.verify", Purpose: "Observe retained Artifact availability.", Type: reflect.TypeFor[evidence.ArtifactVerifyRequest](), Fields: []cli.InputField{judgmentID("artifact_id", "Artifact ID.")}}
}
