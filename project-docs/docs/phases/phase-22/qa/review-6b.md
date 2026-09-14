# Phase 22 - Mid-Phase Review 6b (tasks 16-18)

Reviewer: qa-reviewer
Date: 2026-09-08
Tasks covered: phase-22/task-16, phase-22/task-17, phase-22/task-18
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-22/qa/review-6.md

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-053, VAL-066, VAL-067, VAL-068, VAL-069, VAL-070, VAL-071, VAL-072, VAL-073, VAL-074, VAL-075

## Scope and independent evidence

- Verified the recorded task commits with `git rev-parse --verify <hash>^{commit}`: task 16 `644f30e332a9e5099f0cb31c77fc6b9699fb5db8`, task 17 `ccc15f664c74494cf16219ca350fbec96c51808d`, task 18 initial product `24e3cbdb22dd9e43b55636da313dcef1ef633b6f`, and R6-01 test repair `e04bfbdd0194393b0bfee200477a14bc6c563fb3`. Each resolves, has its recorded subject and `Refs: phase-22/task-N` footer. `git diff --check` is clean for all four commits (terminated `exit_code: 0`).
- Independently ran `project-docs/tools/report-lint.sh` for each report/spec pair. All three returned `PASS report-lint completed with no failures` (terminated `exit_code: 0`). The reports reconcile all required task-spec and QA-plan command occurrences: task 16 12/12, task 17 17/17, and task 18 7/7. Task 18's seventh, retained-Human-fixture invocation is required by its accepted task spec in addition to the QA plan's six-command base.
- Reproduced the R6-01 affected closed-value/corpus command: `go test -tags integration ./cmd/fid/integration/... -run '^TestP22(ClosedValueSchemaCases|ClosedValuesObserveRealResponses|ConformanceCorpusHasEveryAcceptedOperationAndFailure)$' -count=1 -v` - PASS - `1.852s` (terminated `exit_code: 0`). It enumerated 50 accepted member cases, rejected 90 invalid or missing cases, and passed real-response checks for all 12 domains.
- Reproduced the affected final corpus: `go test -tags integration ./cmd/fid/integration/... -run '^TestP22(BudgetCrossingFrameOrder|CompleteOneIssueRecoveryJourney|PortableRecoveryConformance|HumanFixtures|AmendedFailureWireCorpus|ConformanceCorpusHasEveryAcceptedOperationAndFailure|FinalActivated|ClosedValueSchemaCases|ClosedValuesObserveRealResponses)$' -count=1` - PASS - `29.828s` (terminated `exit_code: 0`).
- Reproduced representative task 16/17 evidence: `go test ./internal/maintenance -run '^(TestMaintenanceReportsBrokenRecoveryLinkWithoutMutation|TestMaintenanceRecoveryIntegrityMatrixDoesNotMutate)$' -count=1` - PASS - `0.197s` (terminated `exit_code: 0`); `go test -tags integration ./cmd/fid/integration -run '^(TestP22QualifiedCapabilityMatrixMatchesFakeDiscoveryAndExport|TestRecoverySchemaFreshBinaryInitAndRebuild)$' -count=1` - PASS - `1.175s` (terminated `exit_code: 0`).
- Reproduced the unchanged required tagged integration gate: `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `101.464s` (terminated `exit_code: 0`). The process was polled until its own result reported termination.
- Inspected the retained Human-fixture manifest at `.cache/phase-22-human/checkpoints-4043242815/manifest.json`: it names clean source commit `e04bfbdd0194393b0bfee200477a14bc6c563fb3`, four checkpoints and a usable binary. The recorded `sha256:da111fbdc19a200ff8fe7d263e6fdd7f25915fa3120c83b554649b876dea1d9e` matches the retained binary. VAL-076 through VAL-078 remain Human-pending; VAL-080 remains for final independent QA after this review.

## Task 16 - Retire escalation state and duplicate recovery truth

### Must-fix

- None.

### Improvement

- None.

Observation: the accepted schema-removal evidence remains executable. `internal/db/escalation_removal_test.go:9` verifies latest-schema absence with retained tables and seed removal, while `internal/journal/escalation_removal_test.go:9` exercises both removed-Effect rejection and retained-Effect controls. The report's 12 required command occurrences, one independent reference sample, and per-row provenance support VAL-066 through VAL-068.

## Task 17 - Maintenance, capability qualification, and rebuild budget

### Must-fix

- None.

### Improvement

- None.

Observation: `internal/maintenance/check_test.go:188` covers healthy state, a later-superseded complete-forward activation, every named corruption class, and before/after journal and projection state. `cmd/fid/integration/p22_capability_matrix_test.go:12` structurally compares the fake's capability set with live discovery, exported manifest, and resume output; the focused reproduction passed. The task report records three separate sub-30-second reference samples and two distinct full benchmark commands, satisfying the unchanged VAL-071/QA-78 evidence obligation without substituting this review's focused checks.

## Task 18 - Final recovery conformance and Human validation fixture

### Must-fix

- None.

### Improvement

- None.

Observation: R6-01 is resolved by `cmd/fid/integration/p22_closed_values_test.go:18`. Its hard-coded, field-specific grammar is independently compiled: at run time it reads neither the documentation parser, inventory, nor product allowlists. `TestP22ClosedValueSchemaCases` verifies every one of the 50 documented members across the 15 fields, maintains field membership rather than accepting a merged-domain union, and rejects an unknown string, empty string, null, numeric, boolean, and missing state for every field. `TestP22ClosedValuesObserveRealResponses` obtains all 12 domain objects through public fid responses and applies that same oracle; it rejects an unaccepted value substituted into each captured field. This is the required output-conformance evidence. It does not assert that every valid derived output value is accepted as a lifecycle command input; existing operation/admission tests continue to own lifecycle eligibility.

The R6-01 diff is limited to `p22_closed_values_test.go`, `p22_conformance_test.go`, and `wire-inventory.json`. The corpus now requires a non-empty evidence mapping for all 12 closed domains and verifies that every mapping points to a non-empty Go test. It also corrects the inventory from an invented `claim.get` leaf to the Phase 21 contract's actual `claim.list` leaf: the final 31-operation inventory is 23 amended plus eight retained operations, with no accepted leaf removed. The resume response is correctly used for the DEC-201-only `effective_status` field.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-066, VAL-067, VAL-068 | task-16 | ON TRACK | Removal/replay tests, completion-report command reconciliation, and fresh report lint. |
| VAL-069, VAL-070, VAL-071, VAL-072 | task-17 | ON TRACK | Fresh maintenance/capability/schema reproductions; report retains all three reference samples and two full benchmark runs. |
| VAL-053, VAL-073, VAL-074, VAL-075 | task-18 | ON TRACK | Fresh affected corpus and unchanged tagged integration gate; R6-01's positive and negative closed-value evidence is resolved. |
| VAL-076, VAL-077, VAL-078 | task-18 / Human | DEFERRED | Fixture and manifest are ready; Human outcomes remain intentionally blank. |
| VAL-080 | final independent QA | DEFERRED | This re-review clears its prerequisite; final QA owns the as-built cross-product disposition. |

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| R6-01: the 12 closed-value rows lacked qualifying and rejecting behavioral evidence. | task-18 / VAL-075 | RESOLVED | `e04bfbdd0194393b0bfee200477a14bc6c563fb3` adds independent field grammar, all-member and invalid/missing cases, real public-response validation, per-field rejection mutation, and inventory-reference guards. Fresh closed/corpus evidence passed in 1.852s; the expanded final corpus passed in 29.828s; unchanged tagged integration passed in 101.464s, all with `exit_code: 0`. |

## Cross-cutting observations

- The binary or conditional VAL-075 evidence now contains both states at the field level: every accepted contract member qualifies under the independent grammar, while six invalid/missing alternatives reject for every one of 15 fields. The real-response test then binds that grammar to actual public output and exercises an unaccepted-value rejection for every captured field. This resolves the prior positive-only/evidence-reference weakness without recasting closed output enums as lifecycle inputs.
- No product, schema, public-interface, dependency, gate, threshold, or validation-contract definition changed in the R6-01 continuation. The test-only repair is bounded inside task 18, and the final report `3e960180a22059901037faa4f103ce367a8c6083` preserves its failed and passing history.

## Blocker Candidates

- None.

## Existing Escalation References

- None.

## Final Lines

Review written: project-docs/docs/phases/phase-22/qa/review-6b.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
