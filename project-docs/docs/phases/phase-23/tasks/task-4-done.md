# Phase 23 - Task 4: Immutable Review Records and computed independence - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The atomic pass-Review Event and `review.get` readback now work through the real CLI, but coverage/source validation, non-empty Finding projection, prior-Finding resolution, and the remaining typed VAL-014 through VAL-019 matrix are incomplete.
- Recovery attempts: 5
- Next action: Continue from the green coverage-digest and two-Finding checkpoint; add the typed no-Event repository/Receipt, complete verdict matrix, one-Run conflict, replay/conflict, both independence directions and rework-cycle tests, then implement prior-Finding resolution.
- Human attention required: no

## Summary of what was implemented

- Added the mandated vertical red tracer for `review.record` and retained its original absent-command failure evidence.
- Added `fid review record` and `internal/review`'s write-locked independence admission. It derives conflicts solely from durable mutating Claim and Run history for the reviewed Run's Issue and Envelope revision, and returns the exact conflicting subject and Actor.
- Removed only the `review.recorded` availability fence and registered Task 4's closed `review.record` capability in both authoritative capability allowlists.
- Completed the first atomic `review.record` slice: it locks, refreshes the projection, validates the active review Run/Claim and authority grant, re-observes repository identity, appends a `review.recorded` Event with idempotency, projects `review_records`, and returns the canonical typed result.
- Added `fid review get <review-id>` and a real CLI record/get tracer. Converged the retained executable-leaf and Attention surface inventories for the accepted new Review leaves.
- Added closed-coverage validation that rejects an item and exclusion sharing one Source digest before Event creation. Expanded the real CLI tracer to record, strictly project, and read back two non-blocking Findings in one Event.

## Files created or modified

- `cmd/fid/review_test.go`
- `cmd/fid/review.go`
- `internal/review/review.go`
- `internal/cli/codes.go`
- `cmd/fid/root.go`
- `internal/journal/{current_event_cardinality.go,current_event_json.go,current_event_result_projection.go,current_event_result_values.go,current_event_validation_test.go,current_event_value_catalog.go,current_event_value_relations.go,project_effects.go}`
- `internal/workflow/activation.go`
- `project-docs/docs/phases/phase-23/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/tasks/task-4-done.md`

## Validation Contract Coverage

- VAL-014: PARTIAL - `TestReviewRecordCreatesImmutableReview` records and reads a pass Review with two Findings through `fid`; `TestValidateCoverageRejectsSourceInItemAndExclusion` proves the same-digest negative. Receipt ownership and full bound-field evidence remain required.
- VAL-015: UNVERIFIED - `cmd/fid/review_test.go:12` now proves the mutating-Claim-holder rejection at the real CLI boundary. The required Run-holder, independent-Actor, and rework-cycle cases remain to be implemented and tested.
- VAL-003: UNVERIFIED - The `review.recorded` fence is removed, while the retained Task 1 fence-parity assertion has not yet been converged to a still-fenced Event or run in this recovery.
- VAL-016: PARTIAL - implementation re-observes and compares Repository, but no typed no-Event mismatch test exists.
- VAL-017: PARTIAL - implementation checks verdict/Finding agreement; the required valid and invalid CLI matrix is absent.
- VAL-018: PARTIAL - implementation rejects an existing `review_records.run_id`; no typed test exists.
- VAL-019: PARTIAL - implementation uses `journal.CommitIdempotent`; replay/conflict tests are absent.

## Tests run

- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordRejectsNonIndependentActor$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: `review.record` is absent, so no typed `review-not-independent` rejection can be emitted.
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordRejectsNonIndependentActor$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 8.269s` (terminated `exit_code: 0`)
- The focused green response contains `review-not-independent` and names the mutating Claim and `actor-builder`.
- `go test -race -timeout 30m ./internal/journal/ -run '^TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: `Review` was absent from `currentValueShapes` after the Event fence was removed.
- `go test -race -timeout 30m ./internal/journal/ -run '^TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 1.016s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./internal/journal/ -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 159.995s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./internal/review/ -run '^TestValidateVerdictFindingAgreement$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: Task 4 Review/Finding types and verdict validation were absent.
- `go test -race -timeout 30m ./internal/review/ -run '^TestValidateVerdictFindingAgreement$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/review 1.017s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 9.350s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./internal/review/ -run '^TestValidateCoverageRejectsSourceInItemAndExclusion$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: the closed coverage validator was absent.
- `go test -race -timeout 30m ./internal/review/ -run '^TestValidateCoverageRejectsSourceInItemAndExclusion$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/review 1.017s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: non-empty Finding result values carried the Event timestamp placeholder.
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordCreatesImmutableReview$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 9.070s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.0%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 665.125s` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`); retained executable-leaf and Attention inventories needed Task 4 convergence.
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 64.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.0%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration 694.141s` (terminated `exit_code: 0`)
- Observed wall times: `cmd/fid` 542.160s failed attempt, 473.760s prior rerun, and 457.990s current rerun; `internal/workflow` 230.960s first attempt, cached on later reruns; `cmd/fid/integration` 694.141s prior and 665.125s current. All applicable 600s/900s budgets passed.

## Deviations from spec

- Focused implementation attempt initially failed to compile with `undefined: config.DatabaseName` (terminated `exit_code: 1`). The in-scope repair used the existing `config.DatabaseFileName`; the later focused tracer passed. No task command, assertion, threshold, or policy changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `9416d743f2d20dfaf6cd0e70d605e2a870dd15ac`

## Commit subject

- `feat(review): validate coverage and finding parity`
