# Phase 9 Task 4 Completion Report

## Summary of what was implemented

- Exposed `addressed_in_completion_id` on both public concern read surfaces:
  `fid concern list` and `fid issue get --section concerns`.
- Changed concern read ordering on both surfaces to deterministic append order
  with `ORDER BY rowid DESC`, matching DEC-024 and clearing the ESC-006
  benchmark blocker without changing write behavior.
- Added regression coverage for open-vs-addressed read-back, same-second
  ordering, CLI parity, and the v0.1c concern e2e read-back path.

## Files created or modified

- `cmd/fid/concern_test.go`
- `cmd/fid/e2e_v01c_concern_test.go`
- `internal/issues/concern_list.go`
- `internal/issues/concern_list_test.go`
- `internal/issues/detail.go`
- `internal/issues/issue_test.go`
- `project-docs/docs/phases/phase-9/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-9/tasks/task-4-done.md`
- No additional spec-doc edits were required; `project-docs/specs/06_DESIGN_PLAN_CLI.md`
  already matched DEC-023 and DEC-024 as written.

## Validation Contract Coverage

- `VAL-004`
  - PASS. Evidence: `cmd/fid/concern_test.go`
    `TestConcernAddressUpdatesConcernAndPrintsConfirmation` now asserts
    `addressed_in_completion_id` read-back on both `fid issue get --section concerns`
    and `fid concern list`; `cmd/fid/e2e_v01c_concern_test.go`
    `TestV01CClosureSmoke` asserts the same field on the addressed concern read-back;
    `internal/issues/issue_test.go` and `internal/issues/concern_list_test.go`
    assert `null` for open concerns and the addressing `TC-` id once addressed.
- `VAL-008`
  - PASS and preserved. Evidence: `internal/issues/concern_list_test.go`
    keeps issue/status/composed-filter coverage and now asserts deterministic
    newest-first ordering; `cmd/fid/concern_test.go`
    `TestConcernListFiltersByIssueAndStatus` now checks returned ids in exact
    order rather than sorted order. The required benchmark gate
    `go test -bench . -benchmem ./...` exited `0`.
- `VAL-009`
  - PASS and preserved. Evidence: `cmd/fid/e2e_v01c_concern_test.go`
    `TestV01CConcernRebuildParity` still compares full `fid concern list`,
    `fid issue get`, and `fid issue list` JSON before and after delete-DB +
    `fid rebuild`; parity held with the new field exposed and `rowid DESC`
    ordering intact.
- Coverage thresholds
  - `cmd/fid`: `77.8%`
  - `internal/issues`: `81.8%`
  - `internal/journal`: `80.3%`
  - `internal/policy`: `85.0%`
  - `internal/workflow`: `95.5%`

## Tests run

- `RED` `go test ./cmd/fid -run TestConcernAddressUpdatesConcernAndPrintsConfirmation -count=1`
  - Expected failure before the fix: `addressed_in_completion_id = (*string)(nil)`.
- `GREEN` `go test ./cmd/fid -run TestConcernAddressUpdatesConcernAndPrintsConfirmation -count=1`
  - PASS.
- `RED` `go test ./internal/issues -run TestListConcernsOrdersSameSecondRowsByAppendOrder -count=1`
  - Expected failure before the fix: same-second rows returned in id order,
    not append order.
- `GREEN` `go test ./internal/issues -run TestListConcernsOrdersSameSecondRowsByAppendOrder -count=1`
  - PASS.
- `GREEN` `go test ./internal/issues -run 'TestGetReturnsEnrichedConcernFields|TestGetTreatsAddressedConcernsAsNotOpen|TestGetConcernsOrderSameSecondRowsByAppendOrder|TestListConcernsFiltersByIssueAndStatus|TestListConcernsOrdersSameSecondRowsByAppendOrder' -count=1`
  - PASS.
- `GREEN` `go test ./cmd/fid -run 'TestConcernAddressUpdatesConcernAndPrintsConfirmation|TestConcernListFiltersByIssueAndStatus' -count=1`
  - PASS.
- `GREEN` `go test -tags integration ./cmd/fid -run TestV01CClosureSmoke -count=1`
  - PASS.
- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
- `PASS` `go test -tags integration -race ./...`
- `PASS` `go test -bench . -benchmem ./...`
  - Exit `0`; `cmd/fid` package passed, so the benchmark gate no longer fails on
    `TestConcernListFiltersByIssueAndStatus`.
- `PASS` `go build -o /tmp/fid ./cmd/fid`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
  - No files listed.
- `PASS` `golangci-lint run ./...`
  - Output: `0 issues.`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
  - `cmd/fid` coverage: `77.8%`
  - `internal/issues` coverage: `81.8%`
  - `internal/journal` coverage: `80.3%`
  - `internal/policy` coverage: `85.0%`
  - `internal/workflow` coverage: `95.5%`
- `PASS` `go tool cover -func=coverage.out`
  - Thresholds satisfied. Package rows include
    `internal/issues/detail.go` `83.3%`,
    `internal/issues/concern_list.go` `75.9%`,
    `cmd/fid/concern.go` read-command functions exercised, and overall
    statement coverage `76.5%`.
- `PASS` `go test -tags integration -race ./...`
  - All integration-tagged packages passed, including `cmd/fid` and `internal/issues`.
- `PASS` `go test -bench . -benchmem ./...`
  - Exit `0`.
  - `cmd/fid`: `PASS`.
  - `internal/journal` benchmarks completed:
    `BenchmarkRebuildWorkspace1kEvents-6` and `BenchmarkRebuildWorkspace10kEvents-6`.
- `PASS` `go build -o /tmp/fid ./cmd/fid`
  - Build completed successfully and produced `/tmp/fid`.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `c927de0`

## Commit subject

- `fix: expose concern address read-back deterministically`
