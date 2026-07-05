# Phase 16 - Task 4 Completion Report

## Summary of what was implemented

- Implemented the read-only `fid sql` escape hatch in `internal/sqlquery` and
  `cmd/fid`, including read-statement classification, multi-statement
  rejection, read-only SQLite access, and all-format output for SELECT,
  EXPLAIN, read PRAGMA, and `WITH ... SELECT`.
- Preserved the Task 4 SQL implementation from the prior blocked run and
  completed the DEC-044 unblock by fixing the `cmd/fid/activity_test.go`
  fixture/helper path so seeded activity events derive monotonic dates and IDs
  from the workspace journal instead of hardcoded backdated timestamps.
- Re-ran the exact Task 4 gate to green, then exercised the built `fid sql`
  entrypoint against a fresh workspace.

## Files created or modified

- `cmd/fid/activity_test.go`
- `cmd/fid/root.go`
- `cmd/fid/sql.go`
- `cmd/fid/sql_integration_test.go`
- `cmd/fid/sql_test.go`
- `internal/sqlquery/sqlquery.go`
- `internal/sqlquery/sqlquery_test.go`
- `project-docs/specs/02_TRD.md`
- `project-docs/docs/phases/phase-16/tasks/task-4-done.md`

## Documentation updated

- `project-docs/specs/02_TRD.md`
  - Recorded `internal/sqlquery` in the architecture tree and module-boundary
    table.
- `project-docs/docs/phases/phase-16/tasks/task-4-done.md`
  - Replaced the DEC-044 blocker report with this normal completion report.

## Validation Contract Coverage

- `VAL-008`
  - Covered by the passing Task 4 suite in `go test -race -coverprofile=coverage.out ./...`,
    including the `cmd/fid/sql_test.go` and `internal/sqlquery/sqlquery_test.go`
    acceptance paths for SELECT, EXPLAIN, read PRAGMA, `WITH ... SELECT`, and
    all output formats.
- `VAL-009`
  - Covered by the passing Task 4 suite in `go test -race -coverprofile=coverage.out ./...`
    plus the required integration gate `go test -tags integration -race ./...`,
    including the TRD-SEC-002 mutation-rejection cases in
    `cmd/fid/sql_integration_test.go`.
- DEC-044 regression preservation
  - `go test ./cmd/fid -run 'TestActivityList'` failed before the fixture fix
    with `system clock is behind latest journal date ... refusing to append out
    of order`, then passed after the helper change. This preserved Task 1
    VAL-001..004 behavior and PRD-ERR-014 without changing production journal
    ordering rules.
- Coverage percentages from the required coverage gate
  - `cmd/fid`: `77.1%`
  - `internal/sqlquery`: `65.3%`
- Security note
  - Friendly rejection is enforced by the SQL classifier.
  - The hard backstop remains the read-only SQLite DSN using `mode=ro` with
    `_pragma=query_only(1)` alongside the existing required pragmas.

## Tests run

- Red:
  - `go test ./cmd/fid -run 'TestActivityList'`
- Green:
  - `go test ./cmd/fid -run 'TestActivityList'`
  - `gofmt -l cmd internal`
  - `golangci-lint run ./...`
  - `go test -race -coverprofile=coverage.out ./...`
  - `go tool cover -func=coverage.out`
  - `go test -tags integration -race ./...`
  - `go build -o /tmp/fid ./cmd/fid`
- Real entrypoint exercise:
  - `/tmp/fid init`
  - `/tmp/fid sql "SELECT name FROM roles ORDER BY name LIMIT 1" --format json`
  - Observed stdout: `[{"name":"developer"}]`

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS
- `go test -tags integration -race ./...` - PASS
- `go build -o /tmp/fid ./cmd/fid` - PASS

## Deviations from spec

- DEC-044 added a constrained out-of-scope remediation window to unblock the
  required global gates. This session touched `cmd/fid/activity_test.go` only
  to keep the activity CLI fixture monotonic with the workspace journal.
- The Task 4 SQL implementation files under `cmd/fid/sql.go`,
  `cmd/fid/sql_test.go`, `cmd/fid/sql_integration_test.go`, and
  `internal/sqlquery/` were left intact during the unblock, per the resumed
  scope instruction.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `9cae372`

## Commit subject

- `feat: add read-only fid sql command`
