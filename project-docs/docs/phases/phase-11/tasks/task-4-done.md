## Summary of what was implemented

- Added a started-state claim gate in `internal/issues/transition.go` so entering a started-category state now requires a live claim held by the caller.
- Added started-exit auto-release in the same transition event/transaction by reusing a shared claim-release effect builder from `internal/claims`.
- Extended domain, CLI, and affected integration tests to cover no-claim, wrong-holder, stale-holder, human-holder, matching-holder, and started-exit release behavior.

## Files created or modified

- `internal/issues/transition.go`
- `internal/issues/transition_test.go`
- `internal/claims/shared.go`
- `internal/claims/release.go`
- `internal/claims/claims_test.go`
- `cmd/fid/issue_transition_test.go`
- `cmd/fid/e2e_v01c_ac_test.go`
- `cmd/fid/e2e_v01c_complete_test.go`
- `project-docs/docs/phases/phase-11/tasks/task-4-done.md`

## Documentation updated

- Added this completion report: `project-docs/docs/phases/phase-11/tasks/task-4-done.md`.
- `project-docs/specs/05_DESIGN_PLAN.md`: verified no change required.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: verified no change required.
- `project-docs/docs/phases/phase-11/decisions.md`: no update required because claim-gate precedence stayed aligned with existing edge/rule ordering and no new DEC was needed.

## Validation Contract Coverage

- `VAL-010` PASS: CLI and domain tests cover a `human` claim blocking a developer transition into `in_progress`.
- `VAL-011` PASS: domain and CLI tests cover no-claim, wrong-holder, stale-holder, and matching-holder transitions into started states with zero-append rejection proof.
- `VAL-012` PASS: domain test proves the started-exit transition appends one combined status-change plus claim-delete event; CLI test proves read-back shows the new status and no current claim.

## Tests run

- `go test ./internal/issues -run TestTransitionIntoStartedRequiresLiveCallerClaim -count=1` (`red`, then `green`)
- `go test ./internal/issues -run TestTransitionOutOfStartedReleasesCallerClaimInSameEvent -count=1` (`red`)
- `go test ./internal/issues -run 'TestTransitionIntoStartedRequiresLiveCallerClaim|TestTransitionIntoStartedRejectsWrongHolderAndStaleClaims|TestTransitionOutOfStartedReleasesCallerClaimInSameEvent' -count=1`
- `go test ./cmd/fid -run TestIssueTransition -count=1`
- `go test ./internal/issues -run TestTransition -count=1`
- `go test ./internal/claims -count=1`
- `go test -race ./internal/claims -run TestAcquireIsAtomicAcrossConcurrentAttempts -count=1 -v`
- `go test -race ./internal/claims -run TestAcquireIsAtomicAcrossConcurrentAttempts -count=20`
- `go test -race ./internal/claims -count=1`
- `go test ./internal/issues -run TestTransition -count=1`
- `go test ./cmd/fid -run TestIssueTransition -count=1`

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS
- `go test -tags integration -race ./...` — PASS
- `go build -o /tmp/fid ./cmd/fid` — PASS
- `go test -bench . -benchmem ./...` — N/A, transition benchmark fixtures and journal scale helpers were not touched.

## Deviations from spec

- Updated `cmd/fid/e2e_v01c_ac_test.go` and `cmd/fid/e2e_v01c_complete_test.go` even though they were not named in the task scope list, because the required integration gate exercises the same public transition path and their fixtures had to acquire claims before entering `in_progress`.

## Concerns, technical debt, or notes for the Tech Lead

- `internal/claims.Acquire` and related claim writes still call `journal.EnsureProjectionFresh` from inside their write lock path. During test work this surfaced nested-lock failures when fixtures deliberately left the journal ahead of projection. Runtime paths normally keep projection current, so this task did not widen scope to refactor claim-write lock ordering, but the seam is worth reviewing.

## Dependency-budget note

- No dependency changes.

## Commit hash

- TBC after commit.

## Commit subject

- TBC after commit.
