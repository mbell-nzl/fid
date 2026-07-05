# Phase 16 Task 2 Completion Report

## Summary of what was implemented

- Backfilled `issue create`, `issue ac add`, `issue ac verify`, and `issue complete` so each appends exactly one issue-scoped `activity_log` row inside its existing single journal event.
- Moved completion activity emission into the shared `issues.BuildCompletionEffects` path so `fid dev finish` and `fid dev block` inherit `task_completed` while remaining one event each.
- Extended domain, CLI, activity, and link-adjacent tests to reflect the new timeline behavior without weakening rejection-path assertions.

## Files created or modified

- `internal/issues/create.go`
- `internal/issues/ac.go`
- `internal/issues/ac_verify.go`
- `internal/issues/complete.go`
- `internal/issues/write_helpers.go`
- `internal/issues/issue_test.go`
- `internal/issues/ac_test.go`
- `internal/issues/complete_test.go`
- `internal/issues/transition_test.go`
- `internal/dev/finish.go`
- `internal/dev/block.go`
- `internal/dev/finish_test.go`
- `internal/dev/block_test.go`
- `cmd/fid/issue_test.go`
- `cmd/fid/issue_ac_test.go`
- `cmd/fid/issue_complete_test.go`
- `cmd/fid/activity_test.go`
- `internal/activity/activity_test.go`
- `internal/links/remove_test.go`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/docs/phases/phase-16/tasks/task-2-done.md`

## Documentation updated

- `project-docs/specs/03_TRD_SCHEMA.md`
  Updated the `activity_log.event` vocabulary comment to include `issue_created`, `ac_added`, `ac_verified`, and `task_completed`, and bumped the document timestamp to `2026-07-05`.

## Validation Contract Coverage

- `VAL-005` PASS.
  Domain tests now prove `issue create`, `issue ac add`, `issue ac verify`, and `issue complete` each emit one named `activity_log` row in the same event, while existing rejection-path tests continue to prove zero journal append on rejected writes.
- PI-6 confirmation recorded.
  `internal/dev` tests prove `fid dev finish` and `fid dev block` still emit one event each and now carry both `status_changed` and `task_completed` activity rows via the shared completion builder.

## Tests run

- `go test ./internal/issues -run TestCreateAppendsIssueCreatedActivityInSameEvent`
  Red: failed with `missing issue_created activity effect`.
  Green: passed after `issue create` emitted the `issue_created` activity row.
- `go test ./internal/issues -run 'TestAddAcceptanceCriterionPersistsAndReadsBackInOrdinalOrder|TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent'`
  Red: failed with `missing ac_added activity effect` and `missing ac_verified activity effect`.
  Green: passed after both AC writers emitted their activity rows.
- `go test ./internal/issues ./internal/dev -run 'TestCompleteCreatesAttemptOneAndPersistsCompletionFields|TestCompleteStoresHumanInterventionStateAndProjectsActivity|TestFinishClosesRunCompletesTransitionsAndReleasesClaimInOneEvent|TestBlockRecordsFailedGateCompletionTransitionsAndReleasesClaim'`
  Red: failed on missing `task_completed` activity evidence.
  Green: passed after `BuildCompletionEffects` emitted `task_completed` and the dev composites inherited it.
- `go test ./cmd/fid -run 'TestIssueCreateAppendsOneEventAndProjectsOneIssue|TestIssueACAddCommandAddsAcceptanceCriterionAndPrintsConfirmation|TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation|TestIssueCompleteCommandCreatesAttemptOneAndReturnsJSONCompletionFields'`
  PASS.
- `go test ./internal/issues ./internal/dev ./cmd/fid`
  First run failed on impacted expectation drift in `internal/activity` and `internal/links` tests caused by the new `issue_created` rows.
  Second run PASS after tightening those impacted tests to ignore fixture noise or scope assertions to the target event.
- Runtime smoke via real entrypoint:
  Built `/tmp/fid`, initialized a temp workspace, ran `issue create`, `issue ac add`, `issue ac verify`, `issue complete`, then `/tmp/fid activity list --on issue/TMP.R82HZ0MP5X-951488 --format json`.
  Observed four activity entries for the issue: `issue_created`, `ac_added`, `ac_verified`, and `task_completed`.

## Exact quality-gate results

- `gofmt -l cmd internal` PASS.
- `golangci-lint run ./...` PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` PASS.
  Relevant package coverage from the command output: `cmd/fid` `76.8%`, `internal/issues` `81.6%`, `internal/dev` `75.4%`.
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

- `931d83a`

## Commit subject

- `feat(activity): backfill issue lifecycle writers`
