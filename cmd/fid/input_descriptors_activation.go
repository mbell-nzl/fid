package main

import (
	"reflect"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/evidence"
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

func activationPrepareIdentity(name, meaning string) cli.InputField {
	field := activationInputString(name, meaning, 1, 128)
	field.ScalarBounds = true
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
		activationPrepareIdentity("kit_identity", "Kit identity."), activationPrepareIdentity("backend_identity", "Backend identity."), phase,
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
	}
	inputDescriptorRegistry["activation component"] = activationComponentDescriptor()
}

func activationArray(name, meaning string, element cli.InputField, minimum int, sorted bool) cli.InputField {
	field := cli.InputField{Name: name, Type: "array", Requirement: "required", Meaning: meaning, Element: &element}
	if minimum > 0 {
		field.MinLen = &minimum
	}
	if sorted {
		field.Format = "sorted-unique"
	}
	return field
}

func activationBoundedInteger(name, meaning string, minimum, maximum int64) cli.InputField {
	field := inputInteger(name, meaning, minimum)
	field.MaxInt = &maximum
	return field
}

func activationActorFields() []cli.InputField {
	session := bindingNullable("session_ref", "Session reference, or null.")
	session.Format = ""
	return []cli.InputField{
		judgmentID("actor_id", "Actor ID; profiles are sorted by Actor ID."), activationInputEnum("actor_type", "Actor type.", "agent", "human", "service"),
		activationComponentString("producer", "Actor producer.", 1, 128), activationComponentString("provenance_ref", "Actor provenance reference.", 1, 512), session,
		activationArray("capabilities", "Sorted unique nonempty capabilities.", activationComponentString("capability", "Capability name.", 1, 128), 1, true),
	}
}

func activationGateFields() []cli.InputField {
	environmentName := activationInputString("name", "Environment name: uppercase letters, digits and underscores, starting with a letter or underscore.", 1, 128)
	environmentName.Format = "environment name"
	value := bindingNullable("value", "Literal value for literal sources; null for credential and inherit. Sensitive.")
	value.Format = ""
	value.Sensitive = true
	environment := inputObject("entry", "Environment entry.", environmentName, activationInputEnum("source", "Value source.", "credential", "inherit", "literal"), value)
	literals := activationArray("redact_literals", "Sorted unique literals to redact. Sensitive.", inputString("literal", "required", "Sensitive literal.", "", "<literal>"), 0, true)
	literals.Sensitive, literals.Element.Sensitive = true, true
	argv := activationArray("argv", "Command arguments.", activationComponentString("argument", "One command argument.", 1, 4096), 1, false)
	maximum := 128
	argv.MaxLen = &maximum
	shell := cli.InputField{Name: "shell", Type: "boolean", Requirement: "required", Meaning: "Shell execution is disabled.", Members: []string{"false"}}
	return []cli.InputField{
		judgmentID("definition_id", "Gate Definition ID; definitions are sorted by ID and version."), activationBoundedInteger("definition_version", "Gate Definition version.", 1, 2147483647),
		argv, inputString("cwd", "required", "Relative working directory; use . for the workspace root.", "cwd path", "."), shell,
		activationArray("environment", "Environment entries sorted by name.", environment, 0, false),
		activationBoundedInteger("timeout_ms", "Execution timeout in milliseconds.", 1, 600000), activationBoundedInteger("output_bytes", "Output byte limit.", 1, 10485760),
		activationBoundedInteger("excerpt_bytes", "Excerpt byte limit, also no greater than output_bytes.", 0, 16384),
		activationArray("success_exit_codes", "Sorted unique successful exit codes.", activationBoundedInteger("code", "Successful exit code.", -255, 255), 1, true),
		activationArray("credential_names", "Sorted unique credential names matching credential environment entries.", inputString("credential", "required", "Credential name.", "", "<credential>"), 0, true), literals,
		activationArray("artifact_paths", "Sorted unique relative artifact paths.", inputString("path", "required", "Artifact path.", "artifact path", "<artifact-path>"), 0, true),
		activationBoundedInteger("artifact_bytes", "Artifact byte limit.", 1, 1073741824), activationInputEnum("sandbox", "Sandbox policy.", "none"),
		activationInputEnum("network", "Network policy.", "allowed"), activationInputEnum("write", "Write policy.", "allowed"),
	}
}

func activationComponentString(name, meaning string, minimum, maximum int) cli.InputField {
	field := activationInputString(name, meaning, minimum, maximum)
	field.Format = ""
	return field
}

func activationComponentDescriptor() cli.InputDescriptor {
	content := activationArray("content", "Nonempty typed content selected by component_type.", inputObject("entry", "Selected component entry."), 1, false)
	content.Discriminator = "component_type"
	content.Variants = map[string][]cli.InputField{"authority-profile": activationActorFields(), "gate-definitions": activationGateFields()}
	content.VariantTypes = map[string]reflect.Type{"authority-profile": reflect.TypeFor[workflow.ActorProfile](), "gate-definitions": reflect.TypeFor[evidence.GateDefinition]()}
	return cli.InputDescriptor{Operation: "activation.component", Purpose: "Stage an untrusted activation component.", Type: reflect.TypeFor[workflow.ActivationComponentRequest](), Fields: []cli.InputField{
		judgmentID("activation_id", "Activation ID."), activationInputEnum("component_type", "Component content type.", "authority-profile", "gate-definitions"), judgmentID("component_id", "Component ID."),
		inputInteger("component_version", "Component version.", 1), inputString("expected_prior_digest", "nullable", "Prior component digest, or null.", "digest", "<digest>"), content,
		inputInteger("expected_version", "Expected Activation version.", 1),
	}}
}
