# Phase 21 - Task 11: Generic Runs And Budget Bindings

## Objective

Implement generic Run start, finish, get, and list bound to the active Claim, Envelope revision, Actor, repository identity, and immutable budgets. Replace legacy run recording while proving consumption, Stops, recovery, and continuation stay absent.

## Validation Contract Coverage

- VAL-016: Run budget values retain null consumption and remaining.
- VAL-024: Run identity binds all accepted facts immutably.
- VAL-025: deferred Run/budget/Stop/recovery semantics remain absent.

## Binding Constraints

- Only explicit finish terminalizes; a Run identity is never reused.
- Do not add predecessor, recovery, telemetry, counter, warning, Stop, or Attention state.
- Remove legacy `run record`/`internal/runs` only after generic tests pass; `agent_runs` remains frozen.

## Scope

### Files To Create Or Modify

- `internal/runs/`
- `internal/workflow/` - generic Run aggregate files.
- `internal/cli/`
- `cmd/fid/`
- `internal/journal/`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-11-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§1,7,15`
- `project-docs/specs/01_PRD.md PRD-RUN-005..014`

### Out Of Scope

- Recovery, telemetry, budget consumption, Stops, Attention, and automatic continuation.

## Acceptance Criteria

- [ ] Start/read/finish preserve exact immutable bindings and terminal fields. (VAL-024)
- [ ] Envelope budgets bind to Run with null consumption/remaining throughout. (VAL-016)
- [ ] Deferred state/leaves/events are absent and legacy run recording is replaced. (VAL-025)

## First Failing Test

Add a failing Run-start contract test followed by a finish immutability and absent-recovery matrix.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-11-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/runs ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-11-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 10 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-11-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
