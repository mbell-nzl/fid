## Summary of what was implemented

Implemented Phase 9 Task 1: `fid concern raise`.

- Added the `internal/issues.RaiseConcern` write path with validate-before-append guards for known role, issue existence, completion existence, and completion-to-issue ownership.
- Added the `fid concern raise` CLI under a new `fid concern` namespace with `--on`, `--by`, `--in-completion`, `--summary`, and optional `--body-md`.
- Verified real concern writes round-trip through `fid issue get --section concerns` and flip the already-shipped `concerns_raised` / `open_concerns_count` fields on both issue get and issue list reads.

## Files created or modified

- `cmd/fid/concern.go`
- `cmd/fid/concern_test.go`
- `cmd/fid/root.go`
- `internal/issues/concern_raise.go`
- `internal/issues/concern_raise_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-9/tasks/task-1-done.md`

## Documentation updated

- Updated `project-docs/specs/06_DESIGN_PLAN_CLI.md` in the `fid concern raise` section to state the validate-before-append rejection rules.
- Added this completion report at `project-docs/docs/phases/phase-9/tasks/task-1-done.md`.

## Validation Contract Coverage

- `VAL-001` PASS.
  Evidence:
  `TestRaiseConcernCreatesOpenConcernAndFlipsDerivedFlags`
  `TestConcernRaiseCreatesConcernAndPrintsConfirmation`
  These prove `CON-` ID creation, persisted concern fields, `status="open"`, and read-back via `fid issue get --section concerns`.
- `VAL-002` PASS.
  Evidence:
  `TestRaiseConcernRejectsBeforeAppend`
  `TestConcernRaiseRequiresSummaryFlag`
  Rejected inputs leave journal line counts unchanged and keep `concerns` row count at zero; the valid path appends one concern row and becomes visible on read-back.
- `VAL-003` PASS.
  Evidence:
  `TestRaiseConcernCreatesOpenConcernAndFlipsDerivedFlags`
  `TestConcernRaiseCreatesConcernAndPrintsConfirmation`
  Before the raise, issue reads show `concerns_raised=false` and `open_concerns_count=0`; after the raise, both issue get and issue list show `true` and `1`.

## Tests run

- `go test -race ./internal/issues -run TestRaiseConcern -count=1` PASS
- `go test -race ./cmd/fid -run TestConcernRaise -count=1` PASS
- `gofmt -l cmd internal` PASS
- `golangci-lint run ./...` PASS
- `go test -race -coverprofile=coverage.out ./...` PASS
- `go tool cover -func=coverage.out` PASS
- `go build -o /tmp/fid ./cmd/fid` PASS
- `go test -tags integration -race ./...` N/A
  Reason:
  Task 1 added no `integration`-tagged tests.

## Exact quality-gate results

- `gofmt -l cmd internal` PASS
- `golangci-lint run ./...` PASS
- `go test -race -coverprofile=coverage.out ./...` PASS
  Coverage highlights:
  `cmd/fid` 78.9%
  `internal/issues` 81.8%
  `internal/journal` 80.3%
  `internal/policy` 85.0%
  `internal/workflow` 95.5%
- `go tool cover -func=coverage.out` PASS
  Overall statement coverage: 77.1%
- `go build -o /tmp/fid ./cmd/fid` PASS
- `go test -tags integration -race ./...` N/A
  Reason:
  Task 1 added no `integration`-tagged tests.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `PENDING_COMMIT_HASH`

## Commit subject

- `PENDING_COMMIT_SUBJECT`
