# Phase 21 - Task 11: Generic Runs

## Objective

Replace legacy Run record behavior with start/finish/get/list bound to generic Claim, Envelope revision, Actor, repository, and immutable budgets. Keep recovery and respawn absent.

## Validation Contract Coverage

- VAL-035: owning observable in the Phase 21 validation contract.
- VAL-037: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Run IDs are never reopened or reused.
- Only explicit finish terminalizes.
- Historical `agent_runs` is frozen telemetry, not the generic Run table.

## Scope

### Files To Create Or Modify

- `internal/workflow/`
- `internal/policy/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-11-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §7`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`

### Out Of Scope

- Executor Attempt/Receipt and recovery.

## Acceptance Criteria

- [ ] Start binds every exact accepted value. (VAL-035)
- [ ] Finish is one-way and recovery/respawn stays unavailable. (VAL-037)

## First Failing Test

Add failing Run start binding and terminal/repeat/recovery-absence tests.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-11-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 10 closes.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-11-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
