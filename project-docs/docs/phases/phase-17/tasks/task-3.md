# Phase 17 - Task 3: Implement Maintenance Structural Error Checks

## Objective

Create the `internal/maintenance` domain boundary and implement the four
structural error families required by PRD-OPS-005/006. The checker reports all
violations from one audit instead of stopping at the first row, and performs no
repair or journal write.

## Validation Contract Coverage

- VAL-003: active v0.1 kind hierarchy and sequence/scheduling violations.
- VAL-004: orphaned polymorphic and activity journal-event references.
- VAL-005: SQLite foreign-key integrity violations.
- VAL-006: workflow/state pairs validated against `internal/workflow`.

## Binding Constraints

- Business rules live in `internal/maintenance`; `cmd/fid` remains untouched.
- Mirror every active v0.1 trigger-backed kind/ref rule, including schema-only
  decision/triage rows, but do not query tables that do not exist in v0.1.
- Workflow validity must use `internal/workflow` as the source of truth, not a
  second status map.
- Return all findings in deterministic order and wrap query errors with context.
- Test-only trigger/FK bypass is allowed only to seed impossible states; normal
  production connections retain required pragmas.

## Scope

### Files To Create Or Modify

- `internal/maintenance/check.go` - audit orchestration and result types.
- `internal/maintenance/errors.go` - structural check implementations.
- `internal/maintenance/check_test.go` - valid and invalid rule-family tests.
- `internal/maintenance/test_helpers_test.go` - narrow corrupt-state fixtures.
- `internal/workflow/workflow.go` and tests - only if a small read-only helper is
  required; preserve existing definitions and transition behavior.
- `AGENTS.md` - record `internal/maintenance` in the core/domain layer.
- `project-docs/specs/02_TRD.md` - record the maintenance boundary.
- `project-docs/specs/03_TRD_SCHEMA.md` - document the active audit/trigger
  parity and activity journal-reference check.

### Files To Read For Context

- `migrations/0001_v01a_base.up.sql`
- `migrations/0004_v01_kind_triggers.up.sql`
- `migrations/0005_v01_ref_triggers.up.sql`
- `migrations/0006_v01b_sequence_scheduling.up.sql`
- `internal/workflow/workflow.go`
- `internal/journal/inspect.go`
- `project-docs/specs/03_TRD_SCHEMA.md` sections 10 and 14

### Out Of Scope

- Warning/info checks (task-4) and CLI rendering/exit behavior (task-5).
- Repair, deletion, scratch rebuild/diff, or schema migration.
- v0.2+ comments, QA, validation, and governance command behavior.

## Acceptance Criteria

- [ ] Every active v0.1 kind/sequence invariant has an invalid and valid fixture;
  invalid rows produce one actionable error each. (VAL-003)
- [ ] Links, claims, activity refs, and activity journal-event IDs prove both
  orphan and resolved states. (VAL-004)
- [ ] All `foreign_key_check` rows are preserved as findings with enough identity
  to locate the violation. (VAL-005)
- [ ] All default workflow states pass; unknown workflow/state pairs fail using
  the canonical workflow definitions. (VAL-006)
- [ ] Repeated audit runs return deeply equal ordered results and append no
  journal event. (VAL-003..006)

## First Failing Test

`TestCheckErrorsReportsAllSeededStructuralViolationFamilies` fails because the
`internal/maintenance` package does not exist.

## Documentation Updates

- `AGENTS.md` - add the package to the architecture map.
- `project-docs/specs/02_TRD.md` - add package ownership/dependencies.
- `project-docs/specs/03_TRD_SCHEMA.md` - define v0.1 audit parity.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Required-command rules live in `AGENTS.md`. Every command is reported `PASS`,
or `BLOCKED` with reason; the full gate is required before task close.

List any conditional command required by this task and why:

- None; corrupt-state coverage belongs in package tests using temp SQLite.

## Notes For Developer

- Heavy Mode is approved by DEC-049. This task does not depend on task-2
  implementation.
- Keep each check independently testable and files/functions within quality
  budgets; do not build one monolithic SQL string.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-3-done.md`.
