# Phase 21 - Task 1: Operational Schema And Projection

## Objective

Land the accepted 0013 operational-kernel schema and projection row adapters
for all Phase 21 aggregates. Prove exact up/down/up inventory before Event and
domain operations begin populating the tables.

## Validation Contract Coverage

- VAL-005: migration 0013 implements the accepted additive schema sequence.

## Binding Constraints

- Implement only 0013 now. Do not create migration 0014 before task 21 because
  the normal migration runner would apply its destructive V1 removal.
- Use `work_claims` until task 21 removes legacy `claims`.
- Do not add tables for any deferred capability.

## Scope

### Files To Create Or Modify

- `migrations/0013_p21_operational_kernel.up.sql`
- `migrations/0013_p21_operational_kernel.down.sql`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-1-done.md`

### Files To Read For Context

- `project-docs/specs/03_TRD_SCHEMA.md §8.2`
- `project-docs/docs/phases/phase-21/contracts.md §§14,15`

### Out Of Scope

- Domain validation, CLI commands, and V1 removal.

## Acceptance Criteria

- [ ] Fresh 0013 up/down/up schema inventory matches §8.2 exactly. (VAL-005)
- [ ] Projection adapters apply/read representative valid and invalid rows
  against the accepted constraints. (VAL-005)
- [ ] No deferred table or second authority model is introduced. (VAL-005)

## First Failing Test

Add a failing migration inventory and up/down/up integration test for migration
0013, including proof that 0014 is not yet present or applied.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-1-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/db/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-1-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- No other Developer task may start before this schema task closes.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-1-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
