# Phase 1 - Task 1: SQLite Connection, Pragmas, And Minimal Migration Harness

## Objective

Stand up `internal/db/`: a connection helper that enforces the required pragmas
on every open, plus a `golang-migrate` file-based migration harness and the
v0.1a *minimal* base migration (`projection_meta`, `roles`, `issues` only).
This is the persistence substrate every later Phase 1 task builds on.

## Requirements

- PRD refs: PRD-OPS-009
- TRD refs: TRD §1 (golang-migrate, modernc.org/sqlite), TRD §3
  (`internal/db/`, `migrations/` boundaries), TRD §4.4 pragmas, TRD-ARCH-007
- Design refs: N/A
- Quality refs: Quality Plan §5.1; AGENTS.md Commands

## Validation Contract Coverage

- VAL-014: Every connection opens with `PRAGMA foreign_keys = ON` and
  `PRAGMA recursive_triggers = ON`.

This task also builds the `projection_meta`, `roles`, and `issues` tables that
VAL-001/003/005/006/008/012 depend on, but those assertions are owned by later
tasks that exercise the end-to-end path.

## Scope

### Files To Create Or Modify

- `internal/db/db.go` (connection open + pragmas)
- `internal/db/migrate.go` (migration harness wiring)
- `internal/db/*_test.go`
- `migrations/0001_v01a_base.up.sql`
- `migrations/0001_v01a_base.down.sql`
- `go.mod` / `go.sum` (add `modernc.org/sqlite`, `golang-migrate/migrate`)
- `project-docs/docs/phases/phase-1/tasks/task-1-completion.md`

### Files To Read For Context

- `project-docs/specs/02_TRD.md` §§1-4
- `project-docs/specs/03_TRD_SCHEMA.md` §1.2 (migration placement), §2.1
  (issues DDL), §9.5 (`projection_meta`), §9.7 (`roles`), §15 (pragmas)
- `project-docs/docs/phases/phase-1/qa/validation-contract.md`

### Out Of Scope

- The remaining v0.1 tables, stub tables, triggers, and `issues_enriched`
  view (Phase 2). Create only `projection_meta`, `roles`, and `issues`.
- Journal writer, replayer, ID generation, `fid init`, and any CLI command.
- SQLite triggers (Phase 2). v0.1 `issues` static rules ship as the CHECK
  constraints already in the schema DDL; do not add triggers here.

## Acceptance Criteria

- [ ] Opening a DB connection applies `PRAGMA foreign_keys = ON` and
  `PRAGMA recursive_triggers = ON`; a test queries both and asserts ON. (VAL-014)
- [ ] Migration `0001` up creates exactly `projection_meta`, `roles`, and
  `issues` with the column/CHECK definitions from the schema doc. (no VAL; substrate)
- [ ] Migration `0001` down cleanly drops those tables; up/down round-trips in a
  test. (no VAL; substrate)
- [ ] Applying migrations to a fresh temp DB succeeds and is idempotent on
  re-open (already-applied migrations are a no-op, not an error). (no VAL; substrate)
- [ ] Added runtime dependencies are limited to `modernc.org/sqlite` and
  `golang-migrate/migrate`, documented in the completion report with license. (no VAL)

## TDD Plan

- First failing test: open a connection helper and assert both pragmas report
  ON (fails until the helper sets them).
- Additional tests: migrate-up creates the three tables; migrate-down drops
  them; re-applying migrations is a no-op; an `issues` row violating a CHECK
  (e.g. epic with non-epic workflow) is rejected.
- Regression risks: pragma set per-connection vs once; migration embedding path;
  driver name registration for `modernc.org/sqlite`.

## Documentation Updates

- `project-docs/docs/phases/phase-1/tasks/task-1-completion.md` - completion
  report with VAL-014 evidence, gate results, and dependency-budget note.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: None.
- Audit or release-readiness impact: Dependency-budget note required
  (`modernc.org/sqlite`, `golang-migrate/migrate`).

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Every required command
must be recorded in the completion report as `PASS`, or as `BLOCKED` with reason
in a blocker report. Missing, silently skipped, or narrower substituted
commands are failed gates.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` only if migration tests are placed
  behind the `integration` build tag.

## Notes For Developer

- Keep business rules out of `internal/db/`; it provides connection management
  and migration plumbing only (TRD §3).
- Use the exact DDL from `03_TRD_SCHEMA.md` for the three tables, including all
  CHECK constraints already present in the `issues` DDL. Do not add the v0.2+
  columns/tables.
- `golang-migrate` and `modernc.org/sqlite` are Accepted in TRD §1; still record
  the dependency-budget note (need, alternatives, license) in the report.
