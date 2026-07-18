# Phase 18 - Task 1: Workflow tables migration + seed defaults + gate-obligation seam

## Objective

Add the v0.2 workflow tables (`workflows`, `workflow_states`,
`workflow_transitions`, `workflow_transition_requires`) as a new forward
migration, and seed them with the v0.1 `epic-default`, `task-default`, and
`subtask-default` definitions **verbatim** as reference data. Reserve the
DEC-066 gate-obligation seam (params_json shape + reserved rule vocabulary)
with no receipt behavior, and bump the projection schema version. This task
does not change the `issues` table or the transition engine yet.

## Validation Contract Coverage

- VAL-001: the four workflow tables exist with the specified columns incl. `workflow_transition_requires.params_json`.
- VAL-002: seeded `workflow_states` reproduce the v0.1 state→category maps for all three workflows.
- VAL-003: seeded `workflow_transitions` reproduce the v0.1 allowed edges incl. explicit `<state>→cancelled`, excluding `cancelled→cancelled` and disallowed edges.
- VAL-004: seeded `workflow_transition_requires` reproduce the v0.1 per-transition rules incl. parameterised roles in `params_json`.
- VAL-005: gate-obligation seam present with no behavior (reserved vocab + params_json shape, no seeded use, no receipt logic).
- VAL-012: `ProjectionSchemaVersion` bumped so a prior-version DB is rebuilt/caught-up.

## Binding Constraints

- **Default workflows are seeded by the migration DDL (static INSERTs), not by
  a journal event.** They are system reference data, the exact analog of the
  v0.1 `issues` CHECK constraints (which lived in schema, not the journal).
  This keeps v0.1 journal history unchanged (DEC-066) and satisfies FK ordering
  on rebuild (migrations run before journal replay). Custom-workflow
  definitions become journal events in Phase 19 — **not** in this task.
- **Seed verbatim from `02_TRD.md` §4.6 and the current
  `internal/workflow/workflow.go` definitions; the two must already agree.**
  Every state, category, edge, and rule (including parameterised roles) must
  match exactly. Add frozen fixtures capturing the v0.1 definitions so later
  tasks can assert parity against a stable reference.
- **Cancel edges are seeded explicitly.** The v0.1 engine special-cases "any
  non-cancelled state → cancelled". Seed an explicit `<state> → cancelled`
  transition row for every non-cancelled state in each workflow, and do **not**
  seed `cancelled → cancelled`. The tables become the complete source of truth
  for allowed edges; do not rely on a code special-case remaining after Phase 18.
- **`params_json` schema accommodates a future gate ID/version obligation set**
  (e.g. `{"gates":[{"id":"...","version":N}]}`) without needing a re-migration.
  Store parameterised rule params as JSON text; `NULL` for non-parameterised
  rules. Document the accepted shapes (roles set today; gate-obligation set
  reserved).
- **Reserve, do not implement, `required_gates_current` and `has_completion`.**
  Recognize these rule names as known-but-unimplemented (e.g. named constants /
  documented reservation). No seeded transition may use them and no evaluation
  path may exist for them in this phase. Any receipt/gate execution behavior is
  v0.2b and is out of scope.
- **`workflows.applies_to` mirrors the kind** (`epic`/`task`/`subtask`) for the
  three seeded defaults; `active = 1`.
- **Bump `db.ProjectionSchemaVersion`** (currently 3 → 4) so a DB projected
  under the old version is treated as stale. Provide a matching `.down.sql`.
- Migration is additive: it must not touch `issues` rows, columns, triggers,
  or the `issues_enriched` view (that is Task 4).

## Scope

### Files To Create Or Modify

- `migrations/0007_v02a_workflow_tables.up.sql` (create)
- `migrations/0007_v02a_workflow_tables.down.sql` (create)
- `internal/db/projection_meta.go` (bump `ProjectionSchemaVersion`)
- `internal/workflow/definitions_fixture.go` or equivalent (frozen v0.1
  reference definitions for parity, if not already reusable from `workflow.go`)
- New test files under `internal/db/` (schema + seed assertions) and/or
  `internal/workflow/` (fixture parity)

### Files To Read For Context

- `project-docs/specs/02_TRD.md` §4.6 (default workflows)
- `project-docs/specs/03_TRD_SCHEMA.md` §1.2 (migration placement), §3 (workflow DDL)
- `internal/workflow/workflow.go`, `internal/workflow/rules.go`
- `migrations/0006_v01b_sequence_scheduling.up.sql` (migration style)
- `internal/db/projection_meta.go`, `internal/db/migrate.go`

### Out Of Scope

- Any change to the `issues` table, its triggers, or `issues_enriched` (Task 4).
- Repointing the transition engine or any consumer to the tables (Tasks 2–5).
- `fid workflow define/assign/list`, escalation-default workflow (Phase 19).
- Any gate/receipt/evidence behavior (v0.2b).

## Acceptance Criteria

- [ ] Migration 0007 creates all four workflow tables with the specified
  columns, PKs, and the `workflow_transition_requires`→`workflow_transitions`
  FK (VAL-001)
- [ ] Seeded `workflow_states` equal the v0.1 state→category maps for all three
  workflows, no extra/missing states (VAL-002)
- [ ] Seeded `workflow_transitions` equal the v0.1 allowed edges incl. explicit
  `<state>→cancelled`; `cancelled→cancelled` and disallowed edges absent (VAL-003)
- [ ] Seeded `workflow_transition_requires` equal the v0.1 rules incl. `params_json`
  roles on the parameterised edges (VAL-004)
- [ ] No seeded row uses `required_gates_current`/`has_completion`; a
  gate-obligation-shaped `params_json` validates structurally; no receipt/gate
  code is reachable (VAL-005)
- [ ] `ProjectionSchemaVersion` bumped; a DB stored at the prior version is
  rebuilt/caught-up on the next read (VAL-012)
- [ ] `.down.sql` cleanly reverses the migration

## First Failing Test

A schema-introspection test asserting `workflow_states` contains the
`task-default` rows (`todo→unstarted`, …, `cancelled→cancelled`) — fails until
migration 0007 exists and is seeded.

## Documentation Updates

- None expected in the seven specs for this task (the tables already have draft
  DDL in `03_TRD_SCHEMA.md` §3 and §4.6). Note any deviation from the draft DDL
  in the completion report so the Tech Lead can spec-sync at wrap.

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

- `go test -tags integration -race ./...` — Phase 18 is a persistence-contract
  migration; integration/rebuild coverage is primary evidence and must be green
  at task close.

## Notes For Developer

- This is the isolated schema/reference-data task; keep the diff to migrations,
  the projection-version bump, and fixtures/tests. Do not touch `issues` or any
  consumer.
- If the v0.1 Go definitions and `02_TRD.md` §4.6 disagree in any detail, stop
  and escalate to TL-Lite before seeding — the "verbatim" guarantee depends on
  a single agreed source of truth.
- Completion report at `project-docs/docs/phases/phase-18/tasks/task-1-done.md`.
