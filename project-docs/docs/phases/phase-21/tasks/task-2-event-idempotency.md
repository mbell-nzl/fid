# Phase 21 - Task 2: Event V2, Idempotency, And Role-Free Bootstrap

## Objective

Implement the shared Event-v2 and durable idempotency foundation, then convert `project.initialized` and init/rebuild adapters to the exact role-free bootstrap fixture. This task owns `internal/bootstrap` explicitly because genesis is the first Event-v2 aggregate.

## Validation Contract Coverage

- VAL-003: owning observable in the Phase 21 validation contract.
- VAL-004: owning observable in the Phase 21 validation contract.
- VAL-005: owning observable in the Phase 21 validation contract.
- VAL-006: owning observable in the Phase 21 validation contract.
- VAL-007: owning observable in the Phase 21 validation contract.
- VAL-008: owning observable in the Phase 21 validation contract.

## Binding Constraints

- The genesis Event has the `fid-bootstrap` service Actor, project sequence 1, null repository/idempotency, and empty Effects.
- Initialization remains naturally exclusive and has no idempotency key.
- This task removes role seeding from the new fixture but does not drop the live `roles` table before task 20.

## Scope

### Files To Create Or Modify

- `internal/journal/`
- `internal/bootstrap/`
- `internal/cli/`
- `cmd/fid/init.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/init_test.go`
- `cmd/fid/rebuild_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-2-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§3,14`
- `project-docs/docs/phases/phase-21/public-interface.md §§1,4,6`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`

### Out Of Scope

- Other mutation-family conversions and migration 0014.

## Acceptance Criteria

- [ ] Converted fixtures use exact closed Event v2 and aggregate sequence rules. (VAL-003, VAL-004)
- [ ] Same-key replay, digest conflict, and concurrent serialization match the accepted contract. (VAL-005, VAL-006, VAL-007)
- [ ] Bootstrap emits the exact role-free `project.initialized` fixture. (VAL-008)

## First Failing Test

First add the failing golden `project.initialized` Event-v2 fixture, then idempotency replay/conflict/concurrency tests.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-2-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/journal ./internal/bootstrap ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-2-done.md
```

Every command is required and must terminate with session-captured evidence. The report-lint command runs before the next dispatch. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 1 closes.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-2-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
