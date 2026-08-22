# Phase 21 - Task 1: Operational Schema And Projection

## Objective

Implement only migration 0013 and its schema/constraint integration coverage for the additive Phase 21 operational kernel. Prove the exact reversible inventory before any domain family writes the new tables; Event application and domain row adapters begin in task 2 or their owning domain task.

## Validation Contract Coverage

- VAL-001: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Do not create or apply migration 0014 in this task.
- Use `work_claims` for the collision-free generic table.
- Task 1 remains unauthorized pending Human plan approval, a Ready QA plan, and Tech Lead refinements.

## Scope

### Files To Create Or Modify

- `migrations/0013_p21_operational_kernel.up.sql`
- `migrations/0013_p21_operational_kernel.down.sql`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-1-done.md`

### Files To Read For Context

- `project-docs/specs/03_TRD_SCHEMA.md §8.2`
- `project-docs/docs/phases/phase-21/contracts.md`
- `project-docs/docs/phases/phase-21/public-interface.md`

### Out Of Scope

- Event application, domain row adapters/commands, migration 0014, and V1 deletion.

## Acceptance Criteria

- [ ] 0013 up/down/up matches the accepted tables, indexes, checks, and foreign keys exactly. (VAL-001)
- [ ] Representative row adapters accept valid rows and reject invalid constraints. (VAL-001)

## First Failing Test

Add the failing 0013 schema inventory and up/down/up migration integration test.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-1-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/db/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- None; this is the first ordered task.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-1-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
