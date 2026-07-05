# Phase 16 - Task 4: `fid sql` read-only escape hatch

## Objective

Ship the read-only SQL escape hatch (PRD-OUT-008, TRD-SEC-002): run an arbitrary
read query over the projected DB and render rows in all four formats, while
rejecting anything that is not SELECT/EXPLAIN/PRAGMA/WITH **before execution**,
with a read-only connection as the hard backstop (PI-3).

- New `internal/sqlquery` read package (PI-4): a statement classifier
  (leading-keyword allowlist, comment/whitespace-stripped, reject
  multi-statement) + `Run(options)` that opens the DB **read-only** and returns
  generic column/row results.
- New `cmd/fid/sql.go`: `fid sql <query>` cobra command, all four formats,
  read-only (no `--by`, PI-1), the documented rejection error on stderr.

## Requirements

- PRD refs: PRD-OUT-008
- TRD refs: TRD-SEC-002 (reject non-SELECT/EXPLAIN/PRAGMA/WITH; integration
  test: attempt INSERT/UPDATE/DELETE/DROP, assert rejection); `02_TRD.md` §§2-3
  (add `internal/sqlquery`)
- Design refs: `06_DESIGN_PLAN_CLI.md` §15 `fid sql`
- Quality refs: `08_QUALITY_PLAN.md` §5 (security review, read-only enforcement)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-16/qa/validation-contract.md`.

- VAL-008: `fid sql` accepts SELECT/EXPLAIN/read-PRAGMA/`WITH ... SELECT` and
  renders rows in all four formats.
- VAL-009: `fid sql` rejects non-read statements before execution (non-zero exit,
  documented error) and the read-only connection prevents any mutation even if
  classification is bypassed.

## Scope

### Blocker Addendum: Activity CLI Fixture Monotonicity (DEC-044)

Task 4 is resumed from the blocker report in `task-4-done.md`. Before closing
Task 4, fix the existing `cmd/fid/activity_test.go` failure that seeds a
`2026-07-04` journal event after `fid init` has created a later journal file.
That failure is out of the `fid sql` feature scope, but it blocks the required
global gates and violates `08_QUALITY_PLAN.md` TDD-006.

Allowed remediation scope:

- Edit only `cmd/fid/activity_test.go` and any directly local test helper needed
  to keep the Task 1 CLI activity fixture monotonic with the workspace journal.
- Preserve VAL-001..004 behavior: filter include/exclude semantics, chronological
  ordering, structured `details`, and all-format rendering.
- Preserve PRD-ERR-014 behavior. Do not weaken the journal out-of-order guard,
  bypass production append paths, or narrow the Task 4 gate.
- Record the remediation in the final Task 4 completion report as a
  DEC-044-backed deviation from the original Task 4 file list.

### Files To Create Or Modify

- `internal/sqlquery/sqlquery.go` (new package; PI-4) - `Classify(query)`:
  strip leading line/block comments + whitespace, take the leading keyword,
  allow only `SELECT`/`EXPLAIN`/`PRAGMA`/`WITH` (case-insensitive); reject
  multiple statements (a second non-empty statement after `;`); return a typed
  rejection error carrying the offending statement for the documented message.
  `Run(options)`: `journal.EnsureProjectionFresh` → open the DB **read-only**
  (e.g. `query_only`/`mode=ro`) → `Query` → return ordered columns + `[]map`
  rows (or `[][]any` + columns) for rendering.
- `internal/sqlquery/sqlquery_test.go` (new) - classifier table (accept
  SELECT/EXPLAIN/PRAGMA/WITH incl. lowercase + leading comment; reject INSERT/
  UPDATE/DELETE/DROP/CREATE/REPLACE/`WITH ... DELETE`/multi-statement) and a
  read-only run that returns rows; a mutation attempt against the read-only
  connection fails.
- `cmd/fid/sql.go` (new) - `newSQLCommand()`: `Args: cobra.ExactArgs(1)`,
  `--format`; on classify rejection print the documented stderr error
  (`error: fid sql accepts only SELECT, EXPLAIN, PRAGMA query, and WITH
  statements` + `rejected: <stmt>`) and return non-zero; on success
  `fidformat.Write` the rows.
- `cmd/fid/sql_test.go` (new) - CLI accept/reject for each statement kind and
  format dispatch.
- `cmd/fid/sql_integration_test.go` (new, `//go:build integration`) - the
  TRD-SEC-002 test: attempt INSERT/UPDATE/DELETE/DROP via the command; assert
  each rejected non-zero and the projected DB row counts are unchanged.
- `cmd/fid/root.go` - register `newSQLCommand()`.
- `project-docs/specs/02_TRD.md` - §§2-3: add the `internal/sqlquery` read-only
  escape-hatch package. Bump "Last updated".

### Files To Read For Context

- `internal/db/db.go` (DB open; how to obtain a read-only connection /
  `query_only` pragma for `modernc.org/sqlite`)
- `internal/stats/stats.go`, `cmd/fid/stats.go` (read freshness + cmd pattern)
- `internal/format/format.go` (rendering a generic row list)
- `06_DESIGN_PLAN_CLI.md` §15 (exact output + error text)
- `03_TRD_SCHEMA.md` §13 note (triggers as backstop for `fid sql`) and §12 read
  guidance

### Out Of Scope

- A query DSL (`fid query`, conditional v0.3).
- Saved views (`fid view`, v0.4).
- Any write path or permissions row (read-only, PI-1).

## Acceptance Criteria

- [ ] `fid sql "SELECT ..."` returns rows in json/toon/markdown/yaml; `EXPLAIN`,
  a read `PRAGMA`, and `WITH ... SELECT` are accepted. (VAL-008)
- [ ] `fid sql` rejects INSERT/UPDATE/DELETE/DROP and other non-read statements
  before execution with a non-zero exit and the documented error. (VAL-009)
- [ ] The underlying connection is read-only: a mutation that bypasses the
  classifier (`WITH ... DELETE`, multi-statement) cannot alter the DB. (VAL-009)
- [ ] Integration test (TRD-SEC-002) attempts INSERT/UPDATE/DELETE/DROP and
  asserts rejection + unchanged DB. (VAL-009)
- [ ] `02_TRD.md` records the `internal/sqlquery` boundary. (doc)
- [ ] DEC-044 remediation complete: `cmd/fid/activity_test.go` no longer seeds
  backdated journal events after `fid init`, while Task 1 VAL-001..004 behavior
  remains covered. (quality blocker)

## TDD Plan

- First failing test: `Classify("INSERT INTO issues ...")` returns a rejection —
  fails red (package absent).
- Additional tests: accept SELECT/EXPLAIN/PRAGMA/WITH (mixed case, leading
  comment); reject DML/DDL/`WITH ... DELETE`/multi-statement; `Run` returns rows
  for a SELECT; a mutation against the read-only connection errors; CLI
  accept/reject + the documented error text; integration mutation-rejection.
- Regression risks: defense in depth — the classifier alone is not sufficient
  (`WITH ... DELETE`), so the read-only connection is required; ensure
  `EnsureProjectionFresh` still runs before the read; do not leak SQL errors as
  raw strings where the documented message is expected; functions ≤60 lines.

## Documentation Updates

- `project-docs/specs/02_TRD.md` - add `internal/sqlquery` (§§2-3).
- Completion report: record the read-only enforcement mechanism (classifier +
  read-only connection) under the security note (PI-3).
- Completion report: record the DEC-044 activity-test fixture remediation and
  the exact file(s) touched for it.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: error text must match `06_DESIGN_PLAN_CLI.md` §15.
- Audit or release-readiness impact: TRD-SEC-002 is a security requirement — the
  integration rejection test is the release-readiness evidence; record it PASS.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command in
the completion report as `PASS`, or `BLOCKED` with reason. Record the
`internal/sqlquery` and `cmd/fid` coverage percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` — **required**: this task adds the
  TRD-SEC-002 integration mutation-rejection test (VAL-009); record PASS/BLOCKED.
- `go build -o /tmp/fid ./cmd/fid` — required: CLI build used by the sql tests.

## Notes For Developer

- PI-3: the leading-keyword allowlist gives the friendly error; the **read-only
  connection** is the hard guarantee. Both are required — do not rely on the
  classifier alone.
- Reject multi-statement input (`SELECT 1; DROP TABLE issues`) and `WITH`
  CTEs that lead into DML.
- For `modernc.org/sqlite`, prefer a documented read-only mechanism
  (`?_pragma=query_only(1)` / `?mode=ro` / a `PRAGMA query_only=ON` on the
  connection); record the exact mechanism in the completion report.
- This is a read command — no `--by`, no matrix row (PI-1).
