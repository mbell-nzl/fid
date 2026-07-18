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
- VAL-007: the rebuild preserves every `issues` row and column value exactly (row count and per-column values identical pre/post migration 0008).
- VAL-018: exact dependent-schema-object preservation — the `sqlite_schema`/PRAGMA fingerprint after migration (and after down/up) matches the enumerated inventory in the binding constraints below.
- VAL-019: `fid issue get`/`list` output is byte-identical pre/post migration 0008 for pre-existing issues.

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
- **Recreate exactly the objects dropped with `issues`, and verify (do not
  drop) the objects that merely reference it.** Migrations 0001–0006 define **no
  named index directly on `issues`**; the earlier "recreate all issues indexes"
  wording was inexact and is replaced by this explicit `sqlite_schema` inventory
  (QA-BLK-18-02). Dropping and renaming `issues` removes the objects defined
  `ON issues`; the view and other-table triggers referencing `issues` survive in
  `sqlite_schema` but must be validated as still correct. The exact inventory:
  - **Five triggers defined `ON issues` (dropped with the table — MUST be
    recreated):** `issues_parent_kind_guard_insert`,
    `issues_parent_kind_guard_update`, `issues_kind_immutable` (migration 0004),
    `issues_sequence_scheduling_guard_insert`, and
    `issues_sequence_scheduling_guard_update` (migration 0006).
  - **`issues_enriched` view (migration 0003):** it selects `i.*` from `issues`,
    so drop and recreate it as part of the 12-step procedure and verify identical
    definition and output.
  - **Two named indexes `idx_concerns_issue_status` and `idx_triage_issue_created`
    (migration 0003):** these are on `concerns`/`triage_events`, NOT on `issues`.
    They MUST survive the rebuild unchanged — do not drop or recreate them; assert
    they are still present and unchanged.
  - **Other-table triggers whose body references `issues` (survive; verify
    present and functional):** the 0004 kind-guard triggers
    `acceptance_criteria_issue_kind_guard_insert`,
    `task_completions_issue_kind_guard_insert`,
    `epic_metadata_issue_kind_guard_insert`,
    `triage_events_epic_kind_guard_insert`, `decisions_epic_kind_guard_insert`,
    plus any 0005 ref-integrity triggers that reference `issues`.
  - **PK and UNIQUE autoindexes:** the `issues(id)` PK autoindex and the
    `UNIQUE(parent_id, sequence)` autoindex are recreated by the new table DDL;
    verify them structurally via `PRAGMA index_list`/`index_info`.
  Capture a named `sqlite_schema` + PRAGMA fingerprint before and after the
  migration (and after `.down`/`.up`) and diff against this inventory.
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
- [ ] Every existing `issues` row and column value is preserved exactly
  (row count and per-column values identical pre/post migration 0008) (VAL-007)
- [ ] The post-migration `sqlite_schema`/PRAGMA fingerprint matches the exact
  dependent-object inventory in the binding constraints: the five `ON issues`
  triggers and the `issues_enriched` view are recreated; the two named indexes on
  `concerns`/`triage_events` survive unchanged; other-table triggers referencing
  `issues` remain present; the `issues(id)` PK and `UNIQUE(parent_id, sequence)`
  autoindexes are present per PRAGMA; the fingerprint is stable across `.down`/`.up` (VAL-018)
- [ ] `fid issue get`/`list` output is byte-identical pre/post migration for
  pre-existing issues (VAL-019)
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

- Highest-risk task in the phase. Recreate exactly the objects dropped with
  `issues` (the five `ON issues` triggers and the `issues_enriched` view) and
  verify — do not drop — the objects that only reference `issues` (the two named
  indexes on `concerns`/`triage_events` and the other-table kind/ref triggers),
  per the explicit inventory in the binding constraints. Assert the full
  `sqlite_schema`/PRAGMA fingerprint post-migration. If any object cannot be
  safely recreated or verified, stop and escalate.
- Depends on Task 1 (FK targets). Independent of Task 3, but plan for it to land
  after Task 3 per the phase order.
- Completion report at `project-docs/docs/phases/phase-18/tasks/task-4-done.md`.
