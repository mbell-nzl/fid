# Phase 23 - Task 2: Typed multi-Obligation Envelope contract - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: All Task 2 required gates pass.
- Recovery attempts: 0
- Next action: Phase Runner review.
- Human attention required: no

## Summary of what was implemented

- Envelope authoring now accepts 1..16 typed gate-evidence, review, and human-outcome Obligations and 0..16 protected Decisions.
- Evaluation returns ID-ordered results with distinct validity, freshness, and Artifact availability dimensions.
- The CLI and portable integration model serialize accepted null registrations and typed authoring rows.

## Files created or modified

- `internal/workflow/envelope.go`, `internal/workflow/evaluation.go`, and matching workflow tests.
- `internal/journal/current_event_catalog.go` and `cmd/fid/{envelope_test.go,task18_completion_resume_test.go,integration/p22_portable_schema_test.go}`.
- `project-docs/docs/phases/phase-23/tasks/task-2-done.md`.

## Documentation updated

- Completion report only; implementation conforms to the accepted Phase 23 contracts.

## Validation Contract Coverage

- VAL-006: `internal/workflow/envelope_test.go:TestEnvelopeCreateAcceptsTypedObligationSet` proves one accepted gate-evidence, review, and human-outcome set; `TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring` proves the closed registry rejects invalid type/field combinations.
- VAL-007: `internal/workflow/envelope_test.go:TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring` proves duplicate review profiles and invalid protected-decision authoring reject without admitting a second semantic binding.
- VAL-008: `internal/workflow/envelope_test.go:TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults` proves obligation-ID order and separately reports validity, freshness, and Artifact availability reasons.
- VAL-009: `internal/workflow/envelope_test.go:TestEnvelopeActivationEvaluationIncludesEveryObligation` proves activation evaluates every typed Obligation; `cmd/fid/envelope_test.go:TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions` proves the command surface preserves all three authored rows.
- VAL-010: `internal/workflow/envelope_test.go:TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults` proves `not-applicable` occurs for the metadata-sufficient review row while the unavailable human outcome remains independently reported.

## Tests run

- Focused red: `go test -race -timeout 30m ./internal/workflow/ -run '^TestEnvelopeCreateAcceptsTypedObligationSet$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- Focused green: `go test -race -timeout 30m ./cmd/fid/ -run '^TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 6.856s` (terminated `exit_code: 0`)
- Portable fixture: `go test -tags integration -race -timeout 30m ./cmd/fid/integration/ -run '^TestP22(CurrentEventFamiliesLiveReplayValidation|CurrentEventClosedValuesAndConditionalEffects|PortableFakeAndFIDRejectUnknownEffects)$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 335.024s` (terminated `exit_code: 0`)
- Required workflow pre-edit: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 192.569s coverage: 62.6% of statements` (terminated `exit_code: 0`)
- Required workflow post-edit, uncached observation: `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 193.138s coverage: 63.8% of statements` (terminated `exit_code: 0`)
- Real entrypoint coverage: the tagged `cmd/fid/integration` gate starts the `fid` CLI process and passed in `650.206s` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 63.8% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.9%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 650.206s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -cover ./internal/workflow/` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 63.8% of statements` (terminated `exit_code: 0`)
- Observed wall times: `cmd/fid` `440.185s`; `cmd/fid/integration` `650.206s`; `internal/workflow` `193.138s` for the current-session post-edit uncached required run (the final rerun was cached).

## Deviations from spec

- First full race attempt failed on the displaced single-map completion assertion; its exact token was `FAIL` (terminated `exit_code: 1`). Task 2 converged that assertion to require one list item for its legacy fixture, then the byte-identical gate passed.
- First integration attempt failed on accepted new portable columns; its exact token was `FAIL` (terminated `exit_code: 1`). The Task 2 integration fixture now recognizes the accepted schema columns, and the byte-identical gate passed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
