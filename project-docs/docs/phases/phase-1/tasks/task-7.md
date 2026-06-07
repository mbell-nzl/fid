# Phase 1 - Task 7: Enforce SQLite Pragmas On Every Physical Connection

## Objective

Close ESC-001 / TRG-003. `internal/db.Open` currently applies
`PRAGMA foreign_keys = ON` and `PRAGMA recursive_triggers = ON` by `Exec` on the
pooled `*sql.DB`, which only configures whichever single physical connection the
pool hands out at open time. SQLite pragma state is per physical connection and
`database/sql` opens further physical connections on demand, so VAL-014's "both
pragmas ON on every connection path" is not actually guaranteed. Move pragma
enforcement into the connection string so the driver applies it to every
physical connection, and prove it with a multi-connection regression test.

## Requirements

- PRD refs: PRD-OPS-009
- TRD refs: TRD §4.4 (required pragmas), TRD-ARCH-007, TRD §3
  (`internal/db/` boundary)
- Design refs: N/A
- Quality refs: Quality Plan §5.1; AGENTS.md Commands; QA plan QA-01

## Validation Contract Coverage

- VAL-014: Every SQLite connection opens with `PRAGMA foreign_keys = ON` and
  `PRAGMA recursive_triggers = ON`. Pass requires both pragmas ON on **every**
  connection path; fail is either pragma OFF/unset on **any** connection path.

## Scope

### Files To Create Or Modify

- `internal/db/db.go` (move pragma enforcement into the DSN/connector path)
- `internal/db/db_test.go` (multi-connection regression test)
- `project-docs/docs/phases/phase-1/tasks/task-7-completion.md`

### Files To Read For Context

- `internal/db/db.go`, `internal/db/db_test.go`
- `project-docs/docs/phases/phase-1/decisions.md` (DEC-006)
- `project-docs/docs/current/tl_escalations_current.md` (ESC-001)
- `project-docs/docs/phases/phase-1/qa/review-1.md` (Task 1 Blocking finding)
- `project-docs/docs/phases/phase-1/qa/validation-contract.md` (VAL-014)
- `project-docs/specs/02_TRD.md` §4.4; `project-docs/specs/03_TRD_SCHEMA.md` §15

## Out Of Scope

- Any schema, migration, journal, replay, init, or issue behavior change.
- Changing the `db.Open` public signature or its callers.
- Tuning pool size for performance; only set pool limits where a test must
  force multiple physical connections.

## Acceptance Criteria

- [ ] `db.Open` enforces both pragmas via the `modernc.org/sqlite` connection
  string (`_pragma=foreign_keys(1)`, `_pragma=recursive_triggers(1)`) so they
  are applied by the driver on every physical connection it opens. The file
  path is correctly encoded into the DSN. (VAL-014)
- [ ] A regression test forces multiple distinct physical connections (e.g.
  `SetMaxOpenConns(N)` with N concurrently held `*sql.Conn` via `db.Conn(ctx)`,
  N ≥ 2) and asserts `PRAGMA foreign_keys` and `PRAGMA recursive_triggers` each
  report ON on every one of those connections. The test must fail against the
  old single-`Exec` implementation. (VAL-014)
- [ ] Existing `internal/db` tests (pragmas, migrations, CHECK constraint) still
  pass unchanged in behavior. (no VAL; regression guard)
- [ ] `db.Open` public signature and all existing callers are unchanged. (no VAL)
- [ ] No new runtime dependency is added; if one is unavoidable, record a
  dependency-budget note (need, alternatives, license) in the completion
  report. (no VAL)

## TDD Plan

- First failing test: open the DB, `SetMaxOpenConns(N)` with N ≥ 2, acquire N
  `*sql.Conn` concurrently and hold them so the pool creates N real
  connections, then assert both pragmas ON on each. This fails against the
  current per-pool single-`Exec` code.
- Make it green by moving pragma enforcement into the DSN `_pragma` parameters.
- Additional tests: keep `TestOpenEnablesRequiredPragmas`; confirm migration
  up/down/round-trip and the CHECK-constraint test still pass.
- Regression risks: DSN path encoding for absolute/relative paths and paths with
  special characters; ensuring the in-memory/temp-file path used by tests still
  resolves; not accidentally serializing the pool to a single connection (the
  test must prove >1 physical connection).

## Documentation Updates

- `project-docs/docs/phases/phase-1/tasks/task-7-completion.md` - completion
  report with VAL-014 multi-connection evidence (test name + output showing
  both pragmas ON across N connections), full gate results, and confirmation
  that `db.Open` signature and callers are unchanged. Cite DEC-006 and ESC-001.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: None.
- Audit or release-readiness impact: Resolves ESC-001 / TRG-003. If a dependency
  is added, a dependency-budget note is required.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. A focused test command
may be used during the red/green loop, but it does not substitute for the full
task gate at close. Every required command must be recorded in the completion
report as `PASS`, or as `BLOCKED` with reason in a blocker report. Missing,
silently skipped, or narrower substituted commands are failed gates.

Conditional commands:

```text
go test -tags integration -race ./...
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` only if the multi-connection
  regression test is placed behind the `integration` build tag.

## Notes For Developer

- Keep business rules out of `internal/db/`; this is connection management only
  (TRD §3).
- `modernc.org/sqlite` v1.51.0 (already in `go.mod`) parses `_pragma=` query
  parameters and runs them on each new physical connection. Prefer this over a
  single `Exec`, a `SetMaxOpenConns(1)` pin, or a global connection hook
  (see DEC-006 for the rejected alternatives).
- The proof that matters for VAL-014 is the multi-connection test: it must
  observe >1 distinct physical connection and assert both pragmas ON on each.
  Do not satisfy the assertion by constraining the pool to one connection.
