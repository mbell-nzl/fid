# Phase 16 Task 3 Completion Report

## Summary of what was implemented

- Finished the remaining CF-011 writer backfills by adding `gate_recorded` activity emission to `dev.gate.recorded`, keeping the existing single-event `gate_runs` contract intact and preserving zero-row rejection paths.
- Carried forward the existing partial Task 3 worktree for `concern raise`, `concern address`, `issue unclaim`, and `link add`, then verified those slices together with the new gate activity behavior.
- Added the Phase 16 VAL-007 lifecycle/rebuild integration proof so `fid activity list` now gets exercised across the full Task 3 writer set, including `dev gate`, and delete-and-rebuild parity is asserted against ordered `activity_log` rows.

## Files created or modified

- `internal/dev/gate.go`
- `internal/dev/gate_test.go`
- `internal/issues/concern_raise.go`
- `internal/issues/concern_address.go`
- `internal/claims/release.go`
- `internal/claims/release_test.go`
- `internal/issues/activity_test_helpers_test.go`
- `internal/issues/concern_raise_activity_test.go`
- `internal/issues/concern_address_activity_test.go`
- `internal/links/add.go`
- `internal/links/add_activity_test.go`
- `internal/links/remove_test.go`
- `cmd/fid/link_test.go`
- `cmd/fid/e2e_v01e_activity_test.go`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/docs/phases/phase-16/tasks/task-3-done.md`

## Documentation updated

- `project-docs/specs/03_TRD_SCHEMA.md`
  Confirmed and retained the `activity_log.event` vocabulary update for `concern_raised`, `concern_addressed`, `claim_released`, `link_added`, and `gate_recorded`.
- `project-docs/docs/phases/phase-16/tasks/task-3-done.md`
  Replaced the blocker report with this normal completion report.

## Validation Contract Coverage

- `VAL-006` PASS.
  Focused domain and CLI tests now prove `concern raise`, `concern address`, `issue unclaim`, `link add`, and `dev gate` each append the expected issue-scoped `activity_log` row inside their existing single journal event, while rejection paths still append zero activity rows.
- `VAL-007` PASS.
  `cmd/fid/e2e_v01e_activity_test.go` drives a real CLI lifecycle covering create, `dev start`, AC add/verify, `dev gate`, `dev finish`, `link add`, `issue unclaim`, and concern raise/address, then deletes `fid.sqlite`, runs `fid rebuild`, and compares ordered `activity_log` snapshots for parity.

## Tests run

- `go test ./internal/dev -run 'TestGateRecordsReportedResultOnActiveAttempt|TestGateRejectsWithoutLiveClaimBeforeAppend'`
  Red: failed with `missing gate_recorded activity effect`.
  Green: passed after `dev.gate.recorded` appended one `gate_recorded` activity row and used monotonic issue-scoped timing.
- `go test -tags integration ./cmd/fid -run TestV01EActivityLifecycleCoverageAndRebuildParity`
  Red: failed because the lifecycle timeline did not yet satisfy the new `dev gate` activity expectations.
  Green: passed after the `gate_recorded` activity row was emitted and the new Task 3 lifecycle/rebuild assertions were completed.
- `go test ./internal/claims ./internal/issues ./internal/links ./internal/dev ./cmd/fid -run 'TestReleaseProjectsClaimReleasedActivityInSameEvent|TestReleaseRejectsNonHolderWithoutActivityAppend|TestRaiseConcernProjectsConcernRaisedActivityInSameEvent|TestAddressConcernProjectsConcernAddressedActivityInSameEvent|TestAddressConcernRejectsWithoutActivityAppend|TestAddProjectsLinkAddedActivityForIssueAndRequirementLinks|TestAddRejectsDuplicateWithoutActivityAppend|TestGateRecordsReportedResultOnActiveAttempt|TestGateRejectsWithoutLiveClaimBeforeAppend|TestIssueUnclaimCommandReleasesHolderClaim|TestConcernRaiseCreatesConcernAndPrintsConfirmation|TestConcernAddressUpdatesConcernAndPrintsConfirmation|TestLinkAddCreatesDeclaredAndInverseRowsInOneJournalEvent|TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend'`
  PASS.
- Runtime surface exercised through the real CLI by `go test -tags integration ./cmd/fid -run TestV01EActivityLifecycleCoverageAndRebuildParity`.
  Observed result: `fid activity list --format json` included the Task 3 writer events (`gate_recorded`, `claim_released`, `link_added`, `concern_raised`, `concern_addressed`), and `fid rebuild --by workflow_runner` reproduced an identical ordered `activity_log`.

## Exact quality-gate results

- `gofmt -l cmd internal` PASS.
- `golangci-lint run ./...` PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` PASS.
  Relevant package coverage from the command output: `cmd/fid` `76.8%`, `internal/claims` `81.1%`, `internal/dev` `75.7%`, `internal/issues` `81.6%`, `internal/links` `79.3%`.
- `go tool cover -func=coverage.out` PASS.
  Total statement coverage: `75.9%`.
- `go test -tags integration -race ./...` PASS.
- `go build -o /tmp/fid ./cmd/fid` PASS.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `fd98a02`

## Commit subject

- `feat(activity): finish task 3 lifecycle coverage`
