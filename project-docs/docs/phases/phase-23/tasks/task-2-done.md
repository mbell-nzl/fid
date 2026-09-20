# Phase 23 - Task 2: Typed multi-Obligation Envelope contract - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: MF-01 QA repair is complete; all Task 2 required gates pass.
- Recovery attempts: 1
- Next action: Phase Runner review.
- Human attention required: no

## Summary of what was implemented

- Envelope authoring now accepts 1..16 typed gate-evidence, review, and human-outcome Obligations and 0..16 protected Decisions.
- Evaluation returns ID-ordered results with distinct validity, freshness, and Artifact availability dimensions.
- The CLI and portable integration model serialize accepted null registrations and typed authoring rows.
- This review-2 extension is based on prior repair commit `d71618aae4d8f3f56fee98468d2400a896a7a760`; its earlier matrix evidence and failed attempts remain recorded below.

## Files created or modified

- `internal/workflow/envelope.go`, `internal/workflow/evaluation.go`, and matching workflow tests.
- `internal/journal/current_event_catalog.go` and `cmd/fid/{envelope_test.go,task18_completion_resume_test.go,integration/p22_portable_schema_test.go}`.
- `project-docs/docs/phases/phase-23/tasks/task-2-done.md`.

## Documentation updated

- Completion report only; implementation conforms to the accepted Phase 23 contracts.

## Validation Contract Coverage

- VAL-006: `internal/workflow/envelope_test.go:TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix` proves 0/1/16/17 Obligation boundaries; `TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations` covers every accepted type and each forbidden registration, discriminator, evidence-type, trust-source, freshness, Artifact-policy, and independence combination; `TestEnvelopeAuthoringProtectedDecisionMatrix` and `cmd/fid/envelope_test.go:TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries` prove 0/16/17 protected-Decision boundaries; `TestEnvelopeCreateCommandRejectsEveryForbiddenTypedFieldCombinationWithoutPersistenceChange` sends all 24 forbidden typed-field cases through `envelope create`.
- VAL-007: `internal/workflow/envelope_test.go:TestEnvelopeAuthoringProtectedDecisionMatrix` covers duplicate class, duplicate ordinal, and unguarded action; `cmd/fid/envelope_test.go:TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange` covers all five request-level negatives and every listed protected-Decision rejection, asserting unchanged journal and projection bytes through `snapshotWorkspacePersistence`.
- VAL-008: `internal/workflow/envelope_test.go:TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults` proves Obligation-ID order, distinct validity/freshness/Artifact dimensions, and all unsatisfied reasons; `TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState` proves the qualifying state.
- VAL-009: `internal/workflow/envelope_test.go:TestEnvelopeActivationEvaluationIncludesEveryObligation` proves activation evaluates every typed Obligation; `cmd/fid/envelope_test.go:TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions` proves the command surface preserves all three authored rows.
- VAL-010: `internal/workflow/envelope_test.go:TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState` proves gate Evidence is `available` in the qualifying valid/fresh state and review/Human metadata-sufficient rows are exactly `not-applicable`; the existing independent-state test preserves `missing` for absent gate Evidence.

## Tests run

- Focused red: `go test -race -timeout 30m ./internal/workflow/ -run '^TestEnvelopeCreateAcceptsTypedObligationSet$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- Focused green: `go test -race -timeout 30m ./cmd/fid/ -run '^TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 6.856s` (terminated `exit_code: 0`)
- Portable fixture: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/ -run '^TestP22(CurrentEventFamiliesLiveReplayValidation|CurrentEventClosedValuesAndConditionalEffects|PortableFakeAndFIDRejectUnknownEffects)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 335.024s` (terminated `exit_code: 0`)
- Required workflow pre-edit: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 192.569s coverage: 62.6% of statements` (terminated `exit_code: 0`)
- Required workflow post-edit, uncached observation: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 193.138s coverage: 63.8% of statements` (terminated `exit_code: 0`)
- Real entrypoint coverage: the tagged `cmd/fid/integration` gate starts the `fid` CLI process and passed in `650.206s` (terminated `exit_code: 0`).
- QA repair red: `go test -race -timeout 30m ./internal/workflow/ -run '^(TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix|TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations|TestEnvelopeAuthoringProtectedDecisionMatrix|TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`; test fixture initially produced 17 rows for the 16-row case).
- QA repair green: the same byte-identical focused domain command - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 3.382s` (terminated `exit_code: 0`).
- QA repair command-surface red: `go test -race -timeout 30m ./cmd/fid/ -run '^(TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`; the unguarded-action fixture still left the action guarded).
- QA repair command-surface green: the same byte-identical focused command - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 7.110s` (terminated `exit_code: 0`).
- QA repair named matrix: `go test -race -timeout 30m ./internal/workflow/ -run '^(TestEnvelopeCreateAcceptsTypedObligationSet|TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring|TestEnvelopeCreatePersistsEveryTypedObligationAndProtectedDecision|TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults|TestEnvelopeActivationEvaluationIncludesEveryObligation|TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix|TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations|TestEnvelopeAuthoringProtectedDecisionMatrix|TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 3.520s` (terminated `exit_code: 0`).
- QA repair named command surface: `go test -race -timeout 30m ./cmd/fid/ -run '^(TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions|TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 7.247s` (terminated `exit_code: 0`).
- QA review-2 remaining repair command surface: `go test -race -timeout 30m ./cmd/fid/ -run '^(TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions|TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange|TestEnvelopeCreateCommandRejectsEveryForbiddenTypedFieldCombinationWithoutPersistenceChange)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 7.687s` (terminated `exit_code: 0`). The new case table covers all 24 forbidden typed-field combinations and asserts typed rejection, empty stdout, and byte-identical persistence.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 63.8% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.9%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 650.206s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 63.8% of statements` (terminated `exit_code: 0`)
- Observed wall times: `cmd/fid` `440.185s`; `cmd/fid/integration` `650.206s`; `internal/workflow` `193.138s` for the current-session post-edit uncached required run (the final rerun was cached).
- QA repair attempt 1 exact rerun: `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 203.449s coverage: 64.1% of statements` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.9%` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 681.075s` (terminated `exit_code: 0`)
- QA repair attempt 1 exact rerun: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 198.163s coverage: 64.1% of statements` (terminated `exit_code: 0`)
- QA repair attempt 1 wall times: `cmd/fid` `451.009s`; `cmd/fid/integration` `681.075s`; `internal/workflow` `198.163s` for the final dedicated workflow run (`203.449s` in the full race/coverage gate).
- QA review-2 remaining repair exact rerun: `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.9%` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 666.085s` (terminated `exit_code: 0`)
- QA review-2 remaining repair exact rerun: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- QA review-2 wall times: `cmd/fid` `448.937s`; `cmd/fid/integration` `666.085s`; `internal/workflow` retained noncached package wall time `198.163s` from the prior repair gate, with the current exact dedicated rerun cached and completing in `4.131s` process wall time.

## Deviations from spec

- First full race attempt failed on the displaced single-map completion assertion; its exact token was `FAIL` (terminated `exit_code: 1`). Task 2 converged that assertion to require one list item for its legacy fixture, then the byte-identical gate passed.
- First integration attempt failed on accepted new portable columns; its exact token was `FAIL` (terminated `exit_code: 1`). The Task 2 integration fixture now recognizes the accepted schema columns, and the byte-identical gate passed.
- QA repair attempt 1 preserved two focused red attempts: the new cardinality fixture initially returned 17 instead of 16, and the command fixture initially left the supposedly unguarded action guarded. Both test-only fixture defects were corrected within the bounded review scope; their subsequent byte-identical focused commands passed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
