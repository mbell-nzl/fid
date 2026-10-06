package main

import (
	"reflect"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func activationInputString(name, meaning string, minimum, maximum int) cli.InputField {
	field := bindingString(name, meaning)
	field.MinLen, field.MaxLen = &minimum, &maximum
	return field
}

func activationInputEnum(name, meaning string, members ...string) cli.InputField {
	field := inputString(name, "required", meaning, "", "")
	field.Members = members
	return field
}

func activationPrepareDescriptor() cli.InputDescriptor {
	phase := inputInteger("phase", "Activation protocol phase (21).", 21)
	maximum := int64(21)
	phase.MaxInt = &maximum
	priorVersion := inputInteger("expected_prior_version", "Prior Activation version, or null for the first activation.", 1)
	priorVersion.Requirement = "nullable"
	resources := bindingArray("resources", "Five distinct activation resources, sorted by component.",
		inputObject("resource", "One resource binding.",
			activationInputEnum("component", "Resource component.", "backend", "kit", "overlay", "project", "provider"),
			bindingNullable("old_identity", "Prior resource identity, or null."), bindingString("new_identity", "New resource identity.")), 5)
	maximumLength := 5
	resources.MaxLen = &maximumLength
	return cli.InputDescriptor{Operation: "activation.prepare", Purpose: "Prepare an untrusted activation.", Type: reflect.TypeFor[workflow.ActivationPrepareRequest](), Fields: []cli.InputField{
		activationInputString("kit_identity", "Kit identity.", 1, 128), activationInputString("backend_identity", "Backend identity.", 1, 128), phase,
		activationInputEnum("direction", "Activation direction.", "activate"),
		inputString("expected_prior_activation_id", "nullable", "Prior Activation ID, or null for the first activation.", "ID", "<activation-id>"), priorVersion, resources,
	}}
}

func activationSimpleDescriptors() map[string]cli.InputDescriptor {
	id := judgmentID("activation_id", "Activation ID.")
	version := inputInteger("expected_version", "Expected Activation version.", 1)
	digests := inputArray("component_digests", "Exactly two sorted unique staged component digests.", "digest", true)
	minimum, maximum := 2, 2
	digests.MinLen, digests.MaxLen, digests.Format = &minimum, &maximum, "sorted-unique"
	return map[string]cli.InputDescriptor{
		"activation prepare": activationPrepareDescriptor(),
		"activation commit": {Operation: "activation.commit", Purpose: "Commit a staged activation.", Type: reflect.TypeFor[workflow.ActivationCommitRequest](), Fields: []cli.InputField{
			id, activationInputEnum("expected_status", "Expected Activation status.", "staged"), digests, version,
		}},
		"activation abort": {Operation: "activation.abort", Purpose: "Abort a prepared or staged activation.", Type: reflect.TypeFor[workflow.ActivationAbortRequest](), Fields: []cli.InputField{
			id, activationInputEnum("expected_status", "Expected Activation status.", "prepared", "staged"), activationInputString("reason", "Reason for aborting.", 1, 256), version,
		}},
		"activation finalize": {Operation: "activation.finalize", Purpose: "Finalize a committed activation.", Type: reflect.TypeFor[workflow.ActivationFinalizeRequest](), Fields: []cli.InputField{
			id, activationInputEnum("expected_status", "Expected Activation status.", "committed"), version,
		}},
		"activation resource": {Operation: "activation.resource", Purpose: "Record an activation resource observation.", Type: reflect.TypeFor[workflow.ActivationResourceObservationRequest](), Fields: []cli.InputField{
			id, bindingString("component", "Resource component."), inputInteger("expected_activation_version", "Expected Activation version.", 1),
			inputInteger("expected_version", "Expected resource version.", 1), bindingNullable("observed_identity", "Observed resource identity, or null."),
			activationInputEnum("state", "Observed resource state.", "old", "new", "error"), bindingString("evidence_ref", "Observation evidence reference."), bindingNullable("error_code", "Observation error code, or null."),
		}},
		"activation recover": {Operation: "activation.recover", Purpose: "Select an activation recovery direction.", Type: reflect.TypeFor[workflow.ActivationRecoveryRequest](), Fields: []cli.InputField{
			id, version, bindingString("expected_status", "Expected Activation status."), activationInputEnum("direction", "Recovery direction.", "restore-prior", "complete-forward"), bindingString("reason", "Reason for recovery."),
		}},
		"activation recover-complete": {Operation: "activation.recover-complete", Purpose: "Complete the selected activation recovery.", Type: reflect.TypeFor[workflow.ActivationRecoveryCompleteRequest](), Fields: []cli.InputField{
			id, version, judgmentID("selected_event_id", "Recovery selection Event ID."), bindingArray("resource_versions", "Resource versions sorted by component.",
				inputObject("resource", "One resource version.", bindingString("component", "Resource component."), inputInteger("version", "Resource version.", 1)), 0),
		}},
	}
}

func init() {
	for leaf, descriptor := range activationSimpleDescriptors() {
		inputDescriptorRegistry[leaf] = descriptor
		delete(pendingInputDescriptors, leaf)
	}
}
