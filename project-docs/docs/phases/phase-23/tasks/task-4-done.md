# Phase 23 - Task 4: Immutable Review Records and computed independence - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The Review independence admission boundary is implemented, but immutable Review Event creation and its strict journal projection path remain incomplete.
- Recovery attempts: 3
- Next action: Continue Task 4 from the green strict-journal tracer, implement the atomic `review.record` domain mutation and CLI result/read surface, then complete its full validation matrix and required gates.
- Human attention required: no

## Summary of what was implemented

- Added the mandated vertical red tracer for `review.record` and retained its original absent-command failure evidence.
- Added `fid review record` and `internal/review`'s write-locked independence admission. It derives conflicts solely from durable mutating Claim and Run history for the reviewed Run's Issue and Envelope revision, and returns the exact conflicting subject and Actor.
- Removed only the `review.recorded` availability fence and registered Task 4's closed `review.record` capability in both authoritative capability allowlists.
- Began the required strict journal convergence: its closed Review value shapes, JSON columns, typed result grammar, result re-derivation, projected Finding/Review fills, and Task 4-only repeated Finding/disposition cardinality are registered and tested. The actual atomic Review mutation has not yet emitted an Event through this path.

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

- VAL-014: UNVERIFIED - Review recording and immutable bound-field readback are not implemented; `internal/journal/current_event_validation_test.go:284` now proves the strict journal has the distinct Review/Coverage/Finding/JudgmentSubject registrations required when creation is added.
- VAL-015: UNVERIFIED - `cmd/fid/review_test.go:12` now proves the mutating-Claim-holder rejection at the real CLI boundary. The required Run-holder, independent-Actor, and rework-cycle cases remain to be implemented and tested.
- VAL-003: UNVERIFIED - The `review.recorded` fence is removed, while the retained Task 1 fence-parity assertion has not yet been converged to a still-fenced Event or run in this recovery.
- VAL-016: UNVERIFIED - repository re-observation is not implemented.
- VAL-017: UNVERIFIED - verdict/Finding validation is not implemented.
- VAL-018: UNVERIFIED - one-Review-per-Run enforcement is not implemented.
- VAL-019: UNVERIFIED - Review idempotency is not implemented.

## Tests run

- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordRejectsNonIndependentActor$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: `review.record` is absent, so no typed `review-not-independent` rejection can be emitted.
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordRejectsNonIndependentActor$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 8.269s` (terminated `exit_code: 0`)
- The focused green response contains `review-not-independent` and names the mutating Claim and `actor-builder`.
- `go test -race -timeout 30m ./internal/journal/ -run '^TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: `Review` was absent from `currentValueShapes` after the Event fence was removed.
- `go test -race -timeout 30m ./internal/journal/ -run '^TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 1.016s` (terminated `exit_code: 0`).
- `go test -race -timeout 30m ./internal/journal/ -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/journal 159.995s` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - `not run` (no observed termination); the atomic Review mutation and its typed validation matrix are incomplete.
- `golangci-lint run ./...` - BLOCKED - `not run` (no observed termination); the atomic Review mutation and its typed validation matrix are incomplete.
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `not run` (no observed termination); the atomic Review mutation and its typed validation matrix are incomplete.
- `go tool cover -func=coverage.out` - BLOCKED - `not run` (no observed termination); no fresh full-suite `coverage.out` exists for this incomplete implementation.
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `not run` (no observed termination); the atomic Review mutation and its real-process coverage are incomplete.
- Observed wall times: unavailable; no required full gate was run in this session. The recovery must record `cmd/fid`, `internal/workflow`, and `cmd/fid/integration` wall times against the 600 s / 600 s / 900 s budgets.

## Deviations from spec

- Focused implementation attempt initially failed to compile with `undefined: config.DatabaseName` (terminated `exit_code: 1`). The in-scope repair used the existing `config.DatabaseFileName`; the later focused tracer passed. No task command, assertion, threshold, or policy changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
