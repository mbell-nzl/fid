# Phase 5 - Task 2: `activity` effect op + activity-log projection discipline (CF-005)

## Objective

Add the `activity` effect op to the journal projection engine and establish
`activity_log` projection discipline (CF-005). Per TRD §4.3 an effect `op` may
be `insert`, `update`, `delete`, or `activity`; `activity` effects target
`activity_log` and record who-did-what-and-when as a read-model row. Today
`applyEffect` only handles insert/update/delete, so no event can populate
`activity_log`. This task makes the projection engine understand `activity`
effects so Task 3 (`update`) and Task 5 (`transition`) can emit audited
changes.

`activity_log` is a read-model table written **only** through the projection
path (§14.6). The SQL session-scoped projection guard stays deferred (DEC-006:
unreliable under pooled connections); the discipline is enforced because the
`activity` op is the sole writer and the §14.2 polymorphic-ref backstop trigger
(shipped in Phase 2) still validates `(ref_type, ref_id)`.

## Requirements

- PRD refs: PRD-OUT-017 (activity-log is the projected source for the v0.1e
  read command; this task builds only the write-side projection)
- TRD refs: TRD §4.3 (`activity` op, effect shape), TRD SCHEMA §14.6
  (projection discipline), §14.2 (polymorphic ref backstop), TRD-ARCH-001/004
- Design refs: N/A
- Quality refs: `08_QUALITY_PLAN.md` §4.3 COV-001 (`internal/journal` ≥ 80%)

## Validation Contract Coverage

- VAL-004: `activity` op projects exactly one `activity_log` row per audited
  change; `activity_log` is populated only via projection; the §14.2
  polymorphic-ref backstop still rejects unresolved refs.

## Scope

### Files To Create Or Modify

- `internal/journal/project.go` - add `case "activity"` to `applyEffect`,
  projecting one `activity_log` row from the effect's `pk` + `values`
  (`journal_event_id`, `ref_type`, `ref_id`, `event`, `details`, `actor`,
  `occurred_at`). `details` is stored as JSON text per the schema.
- `internal/journal/project_test.go` - tests (see TDD plan).
- A helper for serialising `details` to JSON text if values arrive as a map.

### Files To Read For Context

- `internal/journal/project.go` (`applyInsert`, effect handling)
- `internal/journal/event.go` (Effect struct)
- `project-docs/specs/03_TRD_SCHEMA.md` §9.4 (activity_log DDL), §14.2, §14.6
- `project-docs/specs/02_TRD.md` §4.3 (example `activity` effect)
- `internal/db/ref_triggers_test.go` (existing polymorphic-ref backstop tests)

### Out Of Scope

- `fid activity list` read command (v0.1e, PRD-OUT-017).
- Emitting activity effects from any command (Tasks 3, 5 do that).
- Re-implementing the §14.6 SQL session guard (deferred per DEC-006).

## Acceptance Criteria

- [ ] Projecting an event with an `activity` effect inserts exactly one
  `activity_log` row with the resolved values; `details` is stored as JSON
  text. (VAL-004)
- [ ] An `activity` effect whose `(ref_type, ref_id)` does not resolve to an
  existing row is rejected by the §14.2 backstop trigger during projection.
  (VAL-004)
- [ ] Rebuild/replay of a journal containing `activity` effects reproduces the
  same `activity_log` rows deterministically. (VAL-004)
- [ ] No projection path other than the `activity` op writes `activity_log`.
  (VAL-004)
- [ ] COV-001 holds for `internal/journal` (≥ 80%).

## TDD Plan

- First failing test: project an event with an `activity` effect and assert one
  `activity_log` row with expected columns. Fails today (unsupported op).
- Additional tests: unresolved polymorphic ref rejected; rebuild determinism
  for activity rows; `details` round-trips as JSON text; multiple activity
  effects in one event each insert one row.
- Regression risks: do not break insert/update/delete dispatch; ensure the
  activity op participates in the same projection transaction so a later failed
  effect rolls back the activity row too; canonical hashing/ordering unaffected.

## Documentation Updates

- Completion report: cite VAL-004 and confirm §14.6 discipline is enforced at
  the application layer (activity op is sole writer), with the SQL guard
  remaining deferred per DEC-006.
- `project-docs/specs/03_TRD_SCHEMA.md` - only if §14.6 wording needs to record
  the application-layer enforcement; otherwise record `verified-no-change`.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: None.
- Audit or release-readiness impact: Advances CF-005 (the activity op half).
  CF-005 closes once Tasks 3 and 5 emit activity effects in real writes.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command as
`PASS` (or `BLOCKED` with reason) in the completion report.

Conditional commands:

- None.

## Notes For Developer

- Treat `activity` as an append-only insert into `activity_log` that needs no
  prior-value precondition (TRD §4.3). The PK is `{"id": "ACT-..."}`.
- ACT IDs are generated at event-build time by the emitting command (Tasks 3/5),
  not by the projector — the projector applies fully-resolved values.
- Keep `details` serialisation deterministic (sorted keys) so rebuild and
  canonical hashing stay stable.
