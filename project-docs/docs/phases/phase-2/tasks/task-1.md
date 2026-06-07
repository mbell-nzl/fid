# Phase 2 - Task 1: Complete The v0.1 Base Schema (Migration 0002)

## Objective

Add the remaining v0.1 base tables and the v0.1 stub tables in a single additive
migration `0002`, completing the v0.1 schema started by Phase 1's `0001`
(`projection_meta`, `roles`, `issues`). Prove the new tables are reachable
through the proven journal/projection machinery — not via direct SQL — and bump
the projection schema version so existing derived DBs rebuild cleanly.

Triggers and the `issues_enriched` view are **out of scope** for this task
(Tasks 2-4). This task ships DDL only plus its verification.

## Requirements

- PRD refs: PRD-OPS-003, PRD-OPS-004 (rebuild/catch-up over the wider schema)
- TRD refs: TRD §4.1 (entity inventory), TRD §4.3 (projection_schema_version
  freshness), TRD §5 (schema versioned via numbered migrations); TRD SCHEMA
  §1.2 (v0.1 base list), §2.3, §5, §7, §9, §10, §11
- Design refs: N/A (no user-facing surface change in this task)
- Quality refs: Quality Plan §4.3 (COV-001/COV-002), §5.1

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-2/qa/validation-contract.md`.

- VAL-001: Migration `0002` up creates every remaining v0.1 base table + stubs
  with specified columns and static CHECKs.
- VAL-002: Migration `0002` down is clean and up→down→up round-trips.
- VAL-003: New tables are populated through journal effects and read back equal
  (not direct SQL).
- VAL-004: `ProjectionSchemaVersion` bump forces rebuild of a prior-version
  checkpoint; current matching checkpoint stays fresh.

## Scope

### Files To Create Or Modify

- `migrations/0002_v01_tables.up.sql` (DDL for the remaining v0.1 base tables +
  stub tables)
- `migrations/0002_v01_tables.down.sql` (drop in reverse dependency order)
- `internal/db/projection_meta.go` (bump `ProjectionSchemaVersion` 1 → 2)
- `internal/db/*_test.go` (migration up/down round-trip; schema introspection;
  freshness version-bump regression)
- `internal/journal/*_test.go` (journal-exercise of representative new tables)
- `project-docs/docs/phases/phase-2/tasks/task-1-completion.md`

### Files To Read For Context

- `project-docs/specs/03_TRD_SCHEMA.md` §1.2, §2.3, §5, §7, §9, §10, §11
- `project-docs/specs/02_TRD.md` §4.1, §4.3
- `migrations/0001_v01a_base.up.sql` / `.down.sql` (existing pattern)
- `internal/db/projection_meta.go`, `internal/db/migrate.go`
- `internal/journal/project.go`, `internal/journal/replay.go`
- `project-docs/docs/phases/phase-2/qa/validation-contract.md`

### Out Of Scope

- `issues_enriched` view and its indexes (Task 2).
- Any trigger (Tasks 3-4).
- Sequence-vs-scheduling backstop trigger, `activity` effect op, §14.6
  projection-only guard (deferred — see contract Scope Deferrals).
- Any change to the `issues`, `roles`, or `projection_meta` table shape (the
  Phase-1 `issues` table already carries every v0.1 column).

## Acceptance Criteria

- [ ] Migration `0002` up creates `epic_metadata`, `acceptance_criteria`,
  `ac_verifications`, `gate_runs`, `task_completions`, `concerns`,
  `requirements`, `agent_runs`, `labels`, `issue_labels`, `links`, `claims`,
  `activity_log`, and stub `decisions`, `escalations`, `triage_events` with the
  columns and static CHECK constraints in TRD SCHEMA. (VAL-001)
- [ ] FK column references match the schema (e.g. `concerns.issue_id` →
  `issues`, `concerns.raised_in_completion_id` → `task_completions`,
  `escalations.resolution_dec` → `decisions`). (VAL-001)
- [ ] Migration `0002` down drops only the tables it created, leaves
  `projection_meta`/`roles`/`issues` intact, and up→down→up round-trips with no
  error. (VAL-002)
- [ ] A journal event (or ordered events) whose insert effects target
  `requirements`, `labels`, `epic_metadata`, `acceptance_criteria`,
  `task_completions`, and `concerns` replays through the existing projection
  engine and the projected rows equal the effect values; population uses journal
  effects, not direct SQL INSERT. (VAL-003)
- [ ] `db.ProjectionSchemaVersion` is bumped to `2`; a stored checkpoint at the
  prior version is decided as `ProjectionRebuild`, and a current-version
  checkpoint matching the journal is decided as `ProjectionFresh`. (VAL-004)
- [ ] `fid rebuild` over the existing Phase-1-style journal still succeeds after
  the migration (no regression). (VAL-004 support)

## TDD Plan

- First failing test: migrate `0002` up on a temp DB and assert a new table
  (e.g. `concerns`) exists with its expected columns — red before the migration
  file exists.
- Additional tests: full table/column/CHECK introspection for all 16 tables;
  down-then-up round-trip; journal-effect population + read-back for the six
  representative tables (order FK-dependent inserts: issue → task_completion →
  concern); `DecideProjectionFreshness` version-mismatch + matching cases; an
  end-to-end rebuild that a prior-version checkpoint forces.
- Regression risks: down-migration FK/drop order; CHECK transcription errors;
  forgetting the version bump (old DBs would serve stale schema); breaking the
  existing `0001` round-trip.

## Documentation Updates

- `project-docs/docs/phases/phase-2/tasks/task-1-completion.md` - completion
  report citing VAL-001..004 evidence, gate outcomes, and any schema-transcription
  deviations.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios.
- Cross-surface parity impact: None (no user-facing output change in this task).
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
in a blocker report. A focused test command may be used during the red/green
loop but does not substitute for the full gate at close.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` if migration round-trip / rebuild
  tests are behind the `integration` build tag.
- `go build -o /tmp/fid ./cmd/fid` to confirm the binary still builds after the
  version bump (no new subcommand, but recommended).

## Notes For Developer

- Transcribe DDL verbatim from `03_TRD_SCHEMA.md`; column names and CHECK
  constraints are compatibility surfaces. Do not invent columns.
- Keep `0002` purely additive: do not `ALTER` or recreate Phase-1 tables.
- The generic projection engine already inserts into arbitrary tables, so VAL-003
  needs no engine change — only well-formed events. `concerns` requires a prior
  `task_completions` row (NOT NULL FK), so order effects/events accordingly.
- Stub tables (`decisions`, `escalations`, `triage_events`) ship schema-only; no
  CLI guards in this task.
- After bumping `ProjectionSchemaVersion`, confirm the Phase-1 freshness tests
  still pass and add the version-mismatch regression.
- `task_completions` static CHECKs only (e.g. `attempt_number`); its kind guard
  is a Task 3 trigger, not a CHECK.
