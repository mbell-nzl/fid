# Phase 2 - Task 2: `issues_enriched` View, Indexes, And Enriched `fid issue get`

## Objective

Add the `issues_enriched` view and its supporting indexes (migration `0003`),
then wire `fid issue get` to read the derived `concerns_raised`,
`open_concerns_count`, and `triage_status` columns from the view instead of the
Phase-1 hard-coded zeros. This is the phase's headline end-to-end increment: a
real command surfaces real derived data flowing journal → projection → view →
JSON.

## Requirements

- PRD refs: PRD-ISSUE-006, PRD-ISSUE-018, PRD-DEV-007, PRD-OUT-001, PRD-OUT-004
- TRD refs: TRD §4.1 (`issues_enriched` derived view), TRD SCHEMA §2.2 (view DDL
  + required indexes); TRD-ARCH-008, TRD-PERF-004
- Design refs: `06_DESIGN_PLAN_CLI.md` `fid issue get`; `05_DESIGN_PLAN.md` §6.1
  (JSON output shape)
- Quality refs: Quality Plan §4.3 (COV-001 `internal/issues`, COV-002 `cmd/fid`)

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-2/qa/validation-contract.md`.

- VAL-005: `issues_enriched` open-concern count + `concerns_raised`.
- VAL-006: `issues_enriched.triage_status` from the latest triage verdict.
- VAL-007: `fid issue get [--format json]` returns enriched columns from the
  view, identical after rebuild (package-level `internal/issues` tests + CLI).
- VAL-008: enriched joins use the concern/triage indexes (no full scan).

## Scope

### Files To Create Or Modify

- `migrations/0003_v01_enriched_view.up.sql` (`issues_enriched` view +
  `idx_concerns_issue_status` + `idx_triage_issue_created`)
- `migrations/0003_v01_enriched_view.down.sql` (drop view + indexes)
- `internal/issues/issue.go` (`Get` reads from `issues_enriched`; populate
  `ConcernsRaised`, `OpenConcernsCount`, `TriageStatus`)
- `internal/issues/issue_test.go` (package-level enriched-read tests, both states)
- `internal/db/*_test.go` (view correctness + `EXPLAIN QUERY PLAN` index check)
- `cmd/fid/issue_test.go` (CLI/integration enriched JSON, both states + rebuild)
- `internal/db/projection_meta.go` (bump `ProjectionSchemaVersion` 2 → 3 if Task
  1 already bumped to 2; keep one monotonic bump per schema-changing migration)
- `project-docs/docs/phases/phase-2/tasks/task-2-completion.md`

### Files To Read For Context

- `project-docs/specs/03_TRD_SCHEMA.md` §2.2
- `project-docs/specs/02_TRD.md` §4.1, TRD-ARCH-008, TRD-PERF-004
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` (`fid issue get`),
  `project-docs/specs/05_DESIGN_PLAN.md` §6.1
- `internal/issues/issue.go` (current `Get` read path + `Issue` struct)
- `project-docs/docs/phases/phase-2/tasks/task-1.md` (tables `concerns`,
  `triage_events` shipped there)

### Out Of Scope

- `fid issue list` and the `--has-concerns` filter (v0.1b).
- TOON/markdown/YAML renderers and `--section` (Phase 3 / v0.1b).
- Triggers (Tasks 3-4).
- Any new write command to create concerns/triage (those ship in v0.1c/v0.2);
  this task drives them through journal events in tests only.

## Acceptance Criteria

- [ ] `issues_enriched` is created exactly per TRD SCHEMA §2.2 (open-concern
  count CTE, latest-triage CTE, `concerns_raised`/`open_concerns_count`/
  `triage_status` columns). (VAL-005, VAL-006)
- [ ] `idx_concerns_issue_status` on `concerns(issue_id, status)` and
  `idx_triage_issue_created` on `triage_events(issue_id, created_at)` exist.
  (VAL-008)
- [ ] View returns `open_concerns_count=N`, `concerns_raised=true` for an issue
  with N open concerns, and `0`/`false` for an issue whose concerns are all
  addressed or absent. (VAL-005)
- [ ] View returns `triage_status` of `'resumed'` (latest `RESUME`), `'blocked'`
  (latest `BLOCK`), and `NULL` (no triage). (VAL-006)
- [ ] `fid issue get <id>` and `--format json` return those enriched values
  sourced from the view; values are identical after delete-DB + `fid rebuild`;
  an issue with no concerns/triage returns `0`/`false`/`null`. (VAL-007)
- [ ] `EXPLAIN QUERY PLAN` for the enriched concern/triage joins shows index
  search on both indexes, not a full table scan. (VAL-008)

## TDD Plan

- First failing test (`internal/issues`): get an issue that has open concerns and
  assert `open_concerns_count > 0` — red while `Get` still hard-codes `0`.
- Additional tests: zero-concern and addressed-only states; triage mapping for
  all three states; CLI JSON for both states; delete-DB + rebuild parity;
  `EXPLAIN QUERY PLAN` index assertion.
- Regression risks: view column ordering vs the `Issue` scan; `Get` must read
  every base issue column from the view (it is `SELECT i.*, ...`); SQLite
  boolean coercion for `concerns_raised`; keeping COV-001 `internal/issues`
  ≥ 80% and COV-002 `cmd/fid` ≥ 70%.

## Documentation Updates

- `project-docs/docs/phases/phase-2/tasks/task-2-completion.md` - completion
  report citing VAL-005..008 evidence and gate outcomes.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: JSON `fid issue get` shape is a compatibility
  surface; enriched fields already exist in the `Issue` struct, so this populates
  them rather than adding fields. Confirm the JSON shape still matches
  `05_DESIGN_PLAN.md` §6.1.
- Audit or release-readiness impact: None.

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
in a blocker report.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go build -o /tmp/fid ./cmd/fid` is required: this task changes a user-facing
  command (`fid issue get`); confirm the binary builds and smoke the enriched
  JSON.
- `go test -tags integration -race ./...` if the CLI/rebuild parity tests are
  behind the `integration` build tag.

## Notes For Developer

- Read enriched values through `issues_enriched`, not by adding a second query;
  the view already LEFT JOINs the counts so a single `SELECT ... FROM
  issues_enriched WHERE id = ?` replaces the current `issues` query.
- Per Phase-1 DEC-007, VAL-007 behavior lives in `internal/issues` (a COV-001
  package): add package-level `internal/issues` tests, not only `cmd/fid`
  integration tests.
- Drive concerns/triage rows in tests via journal events (`ProjectEvent`/replay),
  consistent with VAL-003 — do not raw-INSERT them.
- Keep `cmd/fid/issue.go` a thin adapter (VAL-017): the read logic stays in
  `internal/issues`.
- Coordinate the `ProjectionSchemaVersion` bump with Task 1: each
  schema-changing migration should leave the constant equal to the highest
  applied schema generation so stale DBs rebuild. If Task 1 set it to `2`, this
  task sets it to `3`.
