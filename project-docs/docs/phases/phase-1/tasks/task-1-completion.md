## Summary of what was implemented

Implemented Phase 1 Task 1 only:

- Added `internal/db/Open` to open SQLite connections through `modernc.org/sqlite`
  and enforce `PRAGMA foreign_keys = ON` plus
  `PRAGMA recursive_triggers = ON` on every open.
- Added a file-based migration harness in `internal/db/migrate.go` using
  `golang-migrate/migrate` with the SQLite database driver and local
  `migrations/` files.
- Added migration `0001_v01a_base` with the exact in-scope schema objects only:
  `projection_meta`, `roles`, and `issues`, plus the matching down migration.
- Added focused tests for pragma enforcement, migration creation/drop
  behavior, up/down/up round-trip behavior, idempotent re-apply, and one
  `issues` CHECK constraint rejection path.

## Files created or modified

- `internal/db/db.go`
- `internal/db/migrate.go`
- `internal/db/db_test.go`
- `migrations/0001_v01a_base.up.sql`
- `migrations/0001_v01a_base.down.sql`
- `go.mod`
- `go.sum`
- `project-docs/docs/phases/phase-1/tasks/task-1-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-1/tasks/task-1-completion.md`

## Validation Contract Coverage

- `VAL-014` covered by `TestOpenEnablesRequiredPragmas`, which opens a fresh
  SQLite connection through `internal/db.Open` and asserts both
  `foreign_keys` and `recursive_triggers` pragmas report `1`.
- The migration harness and `0001` DDL provide substrate for later assertions
  (`VAL-001`, `VAL-003`, `VAL-005`, `VAL-006`, `VAL-008`, `VAL-012`) but do
  not claim those validations in this task.

## Tests run

- `go test ./internal/db -run TestOpenEnablesRequiredPragmas`
  Red: failed first on missing `Open` / helper symbols.
  Green: passed after adding the connection helper and runtime dependencies.
- `go test ./internal/db -run 'TestMigrate|TestIssuesCheckConstraintRejectsInvalidWorkflowForKind'`
- `go test ./internal/db`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal`
- `PASS` `golangci-lint run ./...`
- `PASS` `go test -race -coverprofile=coverage.out ./...`
- `PASS` `go tool cover -func=coverage.out`
- `Not required` `go test -tags integration -race ./...` because migration tests
  were not placed behind the `integration` build tag.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- `modernc.org/sqlite`
  Need: provides the accepted pure-Go SQLite driver required by TRD §1 and by
  the per-connection pragma helper in `internal/db`.
  Alternatives considered: `mattn/go-sqlite3` (rejected because it adds CGo
  build complexity the TRD explicitly avoids).
  License: BSD-3-Clause.
- `github.com/golang-migrate/migrate/v4`
  Need: provides the accepted file-based SQL migration harness required by TRD
  §1 and this task's `0001` up/down migration execution.
  Alternatives considered: `pressly/goose` and hand-rolled migration logic
  (rejected because the TRD already accepts `golang-migrate/migrate`).
  License: MIT.

## Commit hash

- Implementation commit: `afbd61b`

## Commit subject

- Implementation subject: `feat(db): add sqlite migration foundation`
