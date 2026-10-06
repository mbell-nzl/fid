package main

import (
	"reflect"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func runExpectedRepository() cli.InputField {
	return inputObject("expected_repository", "Preconditions compared with fid's repository observation.",
		inputString("head_oid", "required", "Full output of git rev-parse HEAD.", "oid", "<head-oid>"),
		cli.InputField{Name: "dirty", Type: "boolean", Requirement: "required", Meaning: "True when git status --porcelain lists any path outside .fid/."})
}

func runExecutor() cli.InputField {
	source := bindingString("source", "Executor provenance.")
	source.Members = []string{"activated-harness"}
	return inputObject("executor", "Executor identity bound to the Run.", judgmentID("id", "Executor ID."), source)
}

func runAuthority() cli.InputField {
	mode := bindingString("mode", "Ordinary or Human force authority.")
	mode.Members = []string{"ordinary", "force"}
	grant := inputObject("grant_ref", "Exact Envelope grant for force authority.", judgmentID("envelope_id", "Envelope ID."), inputInteger("revision", "Envelope revision.", 1), inputInteger("ordinal", "Zero-based grant ordinal.", 0))
	grant.Requirement = "nullable"
	return inputObject("authority", "Authority for abandonment or recovery.", mode, judgmentID("actor_id", "Acting Actor (must equal --actor)."), bindingString("reason", "Reason for acting."), judgmentNullable("reference", "Human reference for force authority."), grant)
}

func runDescriptors() map[string]cli.InputDescriptor {
	id := judgmentID("run_id", "Run ID.")
	version := inputInteger("expected_version", "Expected current Run version.", 1)
	observation := judgmentID("observation_id", "Trusted executor Observation ID.")
	claimVersion := inputInteger("expected_claim_version", "Expected current Claim version.", 1)
	digest := inputString("reconciliation_digest", "required", "Digest from Run inspection.", "digest", "<sha256-digest>")
	state := bindingString("state", "Observed executor liveness.")
	state.Members = []string{"live", "terminated", "unknown"}
	liveness := inputObject("liveness", "Trusted activated-harness observation.", judgmentID("executor_id", "Bound executor ID."), state, bindingString("source_ref", "Harness observation reference."), inputString("observed_at", "required", "Observation timestamp.", "timestamp", "<timestamp>"))
	kind := bindingString("sample_kind", "Baseline starts a chain; delta extends it.")
	kind.Members = []string{"delta", "baseline"}
	value := inputInteger("value", "Measured amount, or null when unavailable.", 0)
	value.Requirement = "nullable"
	measurement := judgmentSource("measurement_ref")
	measurement.Requirement = "nullable"
	obligation := judgmentID("obligation_id", "Optional bound Obligation ID, or null.")
	obligation.Requirement = "nullable"
	ttl := inputInteger("ttl_seconds", "Successor Claim lifetime in seconds (1..86400).", 1)
	maximum := int64(86400)
	ttl.MaxInt = &maximum
	paths := inputArray("adopted_paths", "Sorted unique paths adopted by recovery.", "text", false)
	paths.Format = "sorted-unique"
	return map[string]cli.InputDescriptor{
		"run start":          {Operation: "run.start", Purpose: "Start a Run in the Envelope's worktree, checking the current repository expectation.", Type: reflect.TypeFor[workflow.RunStartRequest](), Fields: []cli.InputField{judgmentID("issue_id", "Issue ID."), judgmentID("envelope_id", "Active Envelope ID."), inputInteger("envelope_revision", "Envelope revision.", 1), judgmentID("claim_id", "Active Claim ID."), inputInteger("claim_version", "Current Claim version.", 1), runExecutor(), obligation, runExpectedRepository()}},
		"run close":          {Operation: "run.close", Purpose: "Close a Run and evaluate Obligations against the observed repository.", Type: reflect.TypeFor[workflow.RunCloseRequest](), Fields: []cli.InputField{id, version, observation, runExpectedRepository()}},
		"run abandon":        {Operation: "run.abandon", Purpose: "Abandon a Run and release its Claim with trusted executor proof.", Type: reflect.TypeFor[workflow.RunAbandonRequest](), Fields: []cli.InputField{id, version, claimVersion, observation, runAuthority(), runExpectedRepository()}},
		"run observe-budget": {Operation: "run.observe-budget", Purpose: "Record a Budget measurement with an optional checked SourceInput.", Type: reflect.TypeFor[workflow.RunBudgetObserveRequest](), Fields: []cli.InputField{id, version, judgmentID("budget_name", "Budget name."), judgmentID("source_sample_id", "Harness sample ID."), kind, inputString("previous_sample_id", "nullable", "Previous sample ID for a delta.", "ID", "<sample-id>"), value, measurement, inputString("observed_at", "required", "Measurement timestamp.", "timestamp", "<timestamp>")}},
		"run observe":        {Operation: "run.observe", Purpose: "Record trusted executor liveness.", Type: reflect.TypeFor[workflow.RunObserveRequest](), Fields: []cli.InputField{id, version, liveness}},
		"run checkpoint":     {Operation: "run.checkpoint", Purpose: "Record a trusted checkpoint for an inspected repository.", Type: reflect.TypeFor[workflow.RunCheckpointRequest](), Fields: []cli.InputField{id, version, digest}},
		"run recover":        {Operation: "run.recover", Purpose: "Recover a terminated or Human-forced Run into one successor.", Type: reflect.TypeFor[workflow.RunRecoveryRequest](), Fields: []cli.InputField{judgmentID("prior_run_id", "Prior Run ID."), inputInteger("expected_run_version", "Expected prior Run version.", 1), claimVersion, observation, digest, paths, runAuthority(), runExecutor(), ttl}},
	}
}

func init() {
	for leaf, descriptor := range runDescriptors() {
		inputDescriptorRegistry[leaf] = descriptor
		delete(pendingInputDescriptors, leaf)
	}
}
