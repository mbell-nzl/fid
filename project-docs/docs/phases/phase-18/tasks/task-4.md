# Phase 18 - Task 4: Migrate issues.workflow_id / status to real FKs (table rebuild)

## Objective

Migrate `issues.workflow_id` to a real FK to `workflows(id)` and
`(issues.workflow_id, issues.status)` to a composite FK to
`workflow_states(workflow_id, state)`, replacing the enum CHECK constraints.
Because SQLite cannot alter constraints in place, rebuild the `issues` table
using the SQLite 12-step procedure, preserving all data and recreating the
`issues_enriched` view, issues triggers, and indexes.

## Validation Contract Coverage

- VAL-006: `issues.workflow_id` FK and `(workflow_id, status)` composite FK enforced — invalid rejected, valid accepted, `PRAGMA foreign_key_check` clean.
- VAL-007: the rebuild preserves all issue rows and recreates the view, triggers, and indexes; migrated v0.1 reads are identical.

## Binding Constraints

- **Use the SQLite table-rebuild procedure** (per the SQLite "Making Other
  Kinds Of Table Schema Changes" guidance): under `PRAGMA foreign_keys=OFF`,
  create `issues_new` with the new FK constraints, `INSERT ... SELECT` all
  columns from `issues`, drop the old table, rename, then recreate every object
  that depended on `issues`. Re-enable `PRAGMA foreign_keys=ON` and run
  `PRAGMA foreign_key_check` — it must return no rows.
- **Replace, do not supplement, the enum CHECKs.** The v0.1 kind→workflow_id
  CHECK and the per-workflow status CHECK are removed; their guarantees now come
  from the FKs (workflow_id must exist in `workflows`; `(workflow_id, status)`
  must exist in `workflow_states`). Preserve the other CHECKs (kind enum,
  priority range, epic est_tokens nullability, positive sequence, `UNIQUE
  (parent_id, sequence)`).
- **Recreate everything attached to `issues`:** the `issues_enriched` view; the
  kind-hierarchy triggers, sequence-scheduling triggers, kind-immutability
  trigger, and polymorphic/ref triggers that reference `issues`; and all issues
  indexes. Verify each exists post-migration. Cross-check against migrations
  0001–0006 so nothing is dropped silently.
- **Data preservation is mandatory.** Row count and every column value must be
  identical pre/post migration; seeded workflow rows (Task 1) guarantee FK
  targets exist for all current issue `(workflow_id, status)` pairs.
- **Deterministic rebuild.** After migrate + `fid rebuild`, the projection must
  be reproducible; the migration must not depend on wall-clock or ordering that
  breaks replay.
- Provide a `.down.sql` that restores the enum-CHECK form (best-effort reverse
  table rebuild) so the migration is reversible in tests.
- Do not change the transition engine or other consumers (done in Tasks 3/5).

## Scope

### Files To Create Or Modify

- `migrations/0008_v02a_issues_workflow_fk.up.sql` (create)
- `migrations/0008_v02a_issues_workflow_fk.down.sql` (create)
- New tests under `internal/db/` (FK enforcement, `foreign_key_check`, view/
  trigger/index existence, data preservation) and/or `internal/journal/`
  (rebuild determinism over the rebuilt schema)

### Files To Read For Context

- `migrations/0001_v01a_base.up.sql` (issues DDL), `0002`–`0006` (triggers/view/indexes on issues)
- `project-docs/specs/03_TRD_SCHEMA.md` §2.1, §14 (issues DDL + triggers)
- `migrations/0007_v02a_workflow_tables.up.sql` (Task 1 — FK targets)
- `internal/db/migrate.go`, `internal/journal/` rebuild tests for the parity harness

### Out Of Scope

- Consumer swaps (Task 5) and the transition-path swap (Task 3).
- `fid maintenance check` changes (Task 5).
- Any new workflow or role tables (Phase 19/20).

## Acceptance Criteria

- [ ] Migration 0008 rebuilds `issues` with `workflow_id` FK → `workflows(id)`
  and `(workflow_id, status)` FK → `workflow_states(workflow_id, state)`; enum
  CHECKs removed, other CHECKs preserved
- [ ] Invalid `workflow_id` or invalid `(workflow_id, status)` insert/update is
  rejected; a valid pair is accepted (VAL-006)
- [ ] `PRAGMA foreign_key_check` returns no rows after migration (VAL-006)
- [ ] All existing issue rows preserved unchanged; `issues_enriched` view,
  issues triggers, and indexes recreated and present (VAL-007)
- [ ] `fid issue get`/`list` output identical pre/post migration for existing
  issues (VAL-007)
- [ ] `.down.sql` reverses the migration; full suite (incl. integration) green

## First Failing Test

An FK-enforcement test attempting to insert an issue with a `(workflow_id,
status)` pair absent from `workflow_states` and asserting rejection — fails
until migration 0008 replaces the enum CHECK with the composite FK.

## Documentation Updates

- None expected in specs (issues FK form is already described in
  `03_TRD_SCHEMA.md` §2.1). Flag any DDL deviation in the completion report.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates.

Conditional command required by this task:

- `go test -tags integration -race ./...` — the migrate + rebuild + `fid issue`
  read parity is primary evidence for VAL-006/007.

## Notes For Developer

- Highest-risk task in the phase. The table rebuild must recreate **every**
  object attached to `issues`; enumerate them from migrations 0001–0006 before
  writing 0008 and assert their existence post-migration. If any object cannot
  be safely recreated, stop and escalate.
- Depends on Task 1 (FK targets). Independent of Task 3, but plan for it to land
  after Task 3 per the phase order.
- Completion report at `project-docs/docs/phases/phase-18/tasks/task-4-done.md`.
