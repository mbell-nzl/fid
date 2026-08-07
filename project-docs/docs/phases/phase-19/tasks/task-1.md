# Phase 19 - Task 1: `workflow_assignments` schema + `escalation-default` seed

## Objective

Add the two schema migrations Phase 19 needs before any command work:
migration `0009` creates the `workflow_assignments` table that
`fid workflow assign` will write to, and migration `0010` seeds the
`escalation-default` workflow (PRD-WF-005a) as system reference data. Teach
the registry loader that `linked_to_decision` is a known-but-unimplemented
rule, and bump the projection schema version so existing projections rebuild.
No CLI change in this task.

## Validation Contract Coverage

- VAL-001: `workflow_assignments` has exactly the DEC-076 **structure**
  (columns, types, CHECKs, FK, both partial unique indexes).
- VAL-027: kind-scope uniqueness both ways — duplicate kind rejected, two
  distinct kinds accepted.
- VAL-028: project-scope uniqueness both ways — first project row accepted,
  second rejected.
- VAL-029: the scope/kind pairing CHECK rejects both malformed shapes and
  accepts both well-formed ones.
- VAL-002: migrations `0009`/`0010` are fully reversible (`up → down → up`
  leaves no residue and no fingerprint drift).
- VAL-003: `escalation-default` is seeded exactly per TRD §4.6, tuple-for-tuple,
  with nothing extra.
- VAL-030: the three Phase 18 defaults are tuple-identical to the frozen
  reference before and after `0010` (counts are not evidence).
- VAL-004: the registry loader surfaces `linked_to_decision` as
  known-but-unimplemented, with no panic.
- VAL-031: `workflow.Load` over the seeded DB writes nothing and no production
  evaluator exists for `linked_to_decision`.
- VAL-005: the projection schema version bump rebuilds a stale projection and
  does not rebuild a fresh current-version one.

Assertions live in
`project-docs/docs/phases/phase-19/qa/validation-contract.md`.

## Binding Constraints

- **Use exactly the DEC-076 DDL** for migration
  `0009_v02a_workflow_assignments.up.sql`:

  ```sql
  CREATE TABLE workflow_assignments (
    scope       TEXT NOT NULL,
    kind        TEXT,
    workflow_id TEXT NOT NULL REFERENCES workflows(id),
    assigned_at DATETIME,
    CHECK (scope IN ('kind', 'project')),
    CHECK ((scope = 'kind' AND kind IS NOT NULL)
        OR (scope = 'project' AND kind IS NULL))
  );

  CREATE UNIQUE INDEX idx_workflow_assignments_kind
    ON workflow_assignments(kind) WHERE scope = 'kind';

  CREATE UNIQUE INDEX idx_workflow_assignments_project
    ON workflow_assignments(scope) WHERE scope = 'project';
  ```

  **Do not** substitute `PRIMARY KEY (scope, kind)`. SQLite allows multiple
  NULLs in a primary key, which would silently permit several project-scope
  rows. The two partial unique indexes are the constraint.
- **Migration `0010_v02a_escalation_default_workflow.up.sql` seeds exactly
  the TRD §4.6 escalation table and nothing more**: workflow
  `('escalation-default', 'Escalation default', 'escalation', 1)`; states
  `open`/`started`/ordinal 1 and `resolved`/`completed`/ordinal 2; the single
  transition `open → resolved`; one requirement row
  `('escalation-default','open','resolved','linked_to_decision', NULL)`.
  **No cancel edge** — the escalation default deliberately declares no
  `cancelled`-category state (DEC-077). Do not touch the three v0.1 default
  workflows.
- **Two separate migration files, not one.** The assignments table and the
  escalation seed are independently reversible concerns; keep the down
  migrations exact mirrors (drop both indexes then the table for `0009`;
  delete the seeded rows in reverse FK order for `0010`).
- **`linked_to_decision` joins the known-but-unimplemented rule set.**
  `internal/workflow/rules.go` already declares `RuleRequiredGatesCurrent`
  and `RuleHasCompletion` and `registry.go` flags them via
  `isReservedUnimplementedRule`. Add `RuleLinkedToDecision` the same way. The
  loader must surface it structurally; it must not gain an evaluator in
  `internal/issues/transition_rules.go`.
- **Bump `db.ProjectionSchemaVersion` 4 → 5.** A new projected table means a
  DB stamped with version 4 must rebuild before serving a read. Update the
  test that pins the constant.
- **This task adds no journal event and no command.** `workflow_assignments`
  stays empty after migration; nothing writes to it yet.

## Scope

### Files To Create Or Modify

- `migrations/0009_v02a_workflow_assignments.up.sql` /
  `migrations/0009_v02a_workflow_assignments.down.sql`
- `migrations/0010_v02a_escalation_default_workflow.up.sql` /
  `migrations/0010_v02a_escalation_default_workflow.down.sql`
- `internal/db/projection_meta.go` — `ProjectionSchemaVersion` 4 → 5
- `internal/db/workflow_assignments_test.go` (new) — schema introspection,
  both-direction uniqueness, down/up reversibility
- `internal/db/escalation_workflow_test.go` (new) — seed parity vs TRD §4.6
- `internal/db/migrate_v01_tables_test.go` — projection-version assertions
- `internal/workflow/rules.go`, `internal/workflow/registry.go`,
  `internal/workflow/registry_test.go` — `linked_to_decision` recognition
- `project-docs/specs/03_TRD_SCHEMA.md` §3 — add `workflow_assignments`,
  extend the `applies_to` comment with `escalation`, note
  `linked_to_decision` in the known-but-unimplemented set
- `project-docs/specs/02_TRD.md` §4.6 — record that `escalation-default` is
  now seeded, and the DEC-076 assignment resolution precedence

### Files To Read For Context

- `migrations/0007_v02a_workflow_tables.up.sql`,
  `migrations/0008_v02a_issues_workflow_fk.up.sql`
- `internal/workflow/registry.go`, `internal/workflow/rules.go`
- `internal/db/projection_meta.go`, `internal/db/migrate.go`
- `project-docs/specs/03_TRD_SCHEMA.md` §3, `project-docs/specs/02_TRD.md` §4.6
- `project-docs/docs/phases/phase-19/decisions.md` (DEC-076, DEC-077, DEC-078)

### Out Of Scope

- Any `fid workflow` command (Tasks 3, 4, 5).
- Any journal event type or effect builder (Task 2 onward).
- Any change to `internal/issues.Create` workflow resolution (Task 4).
- An evaluator for `linked_to_decision` (v0.2d).
- Any change to the three v0.1 default workflows.

## Acceptance Criteria

- [ ] `workflow_assignments` exists with the exact DEC-076 columns, types,
  CHECKs, FK, and two partial unique indexes (VAL-001)
- [ ] A second `kind`-scope row for the same kind is rejected and rows for two
  different kinds are both accepted (VAL-027)
- [ ] The first `project`-scope row is accepted and a second is rejected
  (VAL-028)
- [ ] `scope='kind'` with NULL `kind` and `scope='project'` with a non-NULL
  `kind` are both rejected; both well-formed shapes are accepted (VAL-029)
- [ ] `up → down → up` over migrations `0009`/`0010` leaves an identical
  `sqlite_schema` fingerprint and row set, with no residual index, table, or
  seeded row after `.down` (VAL-002)
- [ ] `escalation-default` is seeded tuple-for-tuple per TRD §4.6, with no
  cancel edge and no extra row in any of the four tables (VAL-003)
- [ ] The three v0.1 default workflows are tuple-identical to the frozen
  `internal/workflow/definitions_fixture.go` reference before and after `0010`
  — assert tuples, not counts (VAL-030)
- [ ] `workflow.Load` returns `escalation-default` with `linked_to_decision`
  surfaced as known-but-unimplemented and no panic (VAL-004)
- [ ] `workflow.Load` over the seeded DB leaves journal bytes, `sqlite_schema`,
  and all four workflow-table row sets unchanged, and no production evaluator
  exists for `linked_to_decision` (VAL-031)
- [ ] `ProjectionSchemaVersion` is 5; a DB stamped 4 rebuilds before serving a
  read, and a fresh version-5 projection is not rebuilt for schema-version
  reasons (VAL-005)
- [ ] `03_TRD_SCHEMA.md` §3 and `02_TRD.md` §4.6 updated in this task

## First Failing Test

`TestMigration0009CreatesWorkflowAssignmentsWithScopedUniqueness`: assert the
partial unique index rejects a second `kind`-scope row for the same kind —
fails until migration `0009` exists.

## Documentation Updates

- `project-docs/specs/03_TRD_SCHEMA.md` §3 — `workflow_assignments` DDL and
  indexes; `applies_to` vocabulary gains `escalation`; `linked_to_decision`
  noted as known-but-unimplemented (DEC-076, DEC-078, DEC-080).
- `project-docs/specs/02_TRD.md` §4.6 — `escalation-default` is seeded by
  migration `0010`; assignment resolution precedence per DEC-076.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. A non-terminating command is
never `PASS`.

Conditional commands required by this task:

- `go test -tags integration -race ./...` — migration and rebuild paths are
  integration-tagged; a schema change must prove them green.
- `go build -o /tmp/fid ./cmd/fid` — the projection-version bump touches the
  command startup path.

## Notes For Developer

- This is the foundation task; Tasks 2-6 all depend on it.
- The DEC-076 DDL is binding, not a suggestion. If SQLite rejects any part of
  it, stop and escalate rather than substituting a different constraint
  strategy — the uniqueness semantics are the point.
- `foreign_keys` is ON for every connection; the `0010` down migration must
  delete requirements before transitions before states before the workflow
  row.
- Do not add an evaluator for `linked_to_decision`. Reaching it from a
  production transition path is a defect this phase asserts against
  (VAL-004, VAL-023).
- Completion report at
  `project-docs/docs/phases/phase-19/tasks/task-1-done.md`.
