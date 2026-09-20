# Phase 23 - Task 4: Immutable Review Records and computed independence - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: The Review independence admission boundary is implemented, but immutable Review Event creation and its strict journal projection path remain incomplete.
- Recovery attempts: 1
- Next action: Continue Task 4 from the green independence tracer, implement atomic `review.recorded` projection/result parity and the remaining accepted validation matrix, then complete the required gates.
- Human attention required: no

## Summary of what was implemented

- Added the mandated vertical red tracer for `review.record` and retained its original absent-command failure evidence.
- Added `fid review record` and `internal/review`'s write-locked independence admission. It derives conflicts solely from durable mutating Claim and Run history for the reviewed Run's Issue and Envelope revision, and returns the exact conflicting subject and Actor.

## Files created or modified

- `cmd/fid/review_test.go`
- `cmd/fid/review.go`
- `internal/review/review.go`
- `internal/cli/codes.go`
- `cmd/fid/root.go`
- `project-docs/docs/phases/phase-23/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-23/tasks/task-4-done.md`

## Validation Contract Coverage

- VAL-014: UNVERIFIED - Review recording and immutable bound-field readback are not implemented.
- VAL-015: UNVERIFIED - `cmd/fid/review_test.go:12` now proves the mutating-Claim-holder rejection at the real CLI boundary. The required Run-holder, independent-Actor, and rework-cycle cases remain to be implemented and tested.
- VAL-003: UNVERIFIED - `review.recorded` remains fenced; no convergence has started.
- VAL-016: UNVERIFIED - repository re-observation is not implemented.
- VAL-017: UNVERIFIED - verdict/Finding validation is not implemented.
- VAL-018: UNVERIFIED - one-Review-per-Run enforcement is not implemented.
- VAL-019: UNVERIFIED - Review idempotency is not implemented.

## Tests run

- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordRejectsNonIndependentActor$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red state: `review.record` is absent, so no typed `review-not-independent` rejection can be emitted.
- `go test -race -timeout 30m ./cmd/fid/ -run '^TestReviewRecordRejectsNonIndependentActor$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/cmd/fid 8.269s` (terminated `exit_code: 0`); the response contains `review-not-independent` and names the mutating Claim and `actor-builder`.

## Exact quality-gate results

- `gofmt -l cmd internal` - BLOCKED - `not run` (no observed termination); the implementation is intentionally incomplete and the required red tracer remains failing.
- `golangci-lint run ./...` - BLOCKED - `not run` (no observed termination); the implementation is intentionally incomplete and the required red tracer remains failing.
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `not run` (no observed termination); the implementation is intentionally incomplete and the required red tracer remains failing.
- `go tool cover -func=coverage.out` - BLOCKED - `not run` (no observed termination); no fresh `coverage.out` exists for this incomplete implementation.
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `not run` (no observed termination); the mandatory tagged gate is deferred only to the same-task recovery after implementation.
- Observed wall times: unavailable; no required full gate was run in this session. The recovery must record `cmd/fid`, `internal/workflow`, and `cmd/fid/integration` wall times against the 600 s / 600 s / 900 s budgets.

## Deviations from spec

- Focused implementation attempt initially failed to compile with `undefined: config.DatabaseName` (terminated `exit_code: 1`). The in-scope repair used the existing `config.DatabaseFileName`; the later focused tracer passed. No task command, assertion, threshold, or policy changed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `45a2a168e0abcf1d4d43dd37138f5e3b6024c485`

## Commit subject

- `feat(review): enforce recorded review independence`
