# Phase 21 - Task 18: Guarded Terminal Transition, Completion, And Resume

## Objective

Wire trusted Evidence into the shared evaluator, permit the terminal-success Issue transition only after satisfaction, reuse the evaluator for later Envelope completion, and expose bounded resume state.

## Validation Contract Coverage

- VAL-064: owning observable in the Phase 21 validation contract.
- VAL-065: owning observable in the Phase 21 validation contract.
- VAL-066: owning observable in the Phase 21 validation contract.
- VAL-067: owning observable in the Phase 21 validation contract.
- VAL-068: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Gate execution/Receipt link/evaluation precede terminal transition.
- Completion independently re-evaluates current Obligations and terminal Issue.
- Completion does not require an active Claim/Run; the journey may finish and
  release them after terminal transition and before completion.
- Never read/write `task_completions` as generic completion.

## Scope

### Files To Create Or Modify

- `internal/workflow/`
- `internal/issues/`
- `internal/evidence/`
- `internal/validation/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-18-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§6,10,16`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/docs/phases/phase-21/gate-d3-2-proposal.md`

### Out Of Scope

- Export, V1 physical removal, and recovery.

## Acceptance Criteria

- [ ] Qualifying Evidence permits terminal transition through shared evaluator. (VAL-064)
- [ ] Every nonqualifying evidence state blocks atomically. (VAL-065)
- [ ] Completion reuses evaluator and requires terminal Issue. (VAL-066)
- [ ] Failed completion has all reasons/no effects. (VAL-067)
- [ ] Resume is bounded/read-only. (VAL-068)

## First Failing Test

Add positive/negative terminal guard, evaluator-identity completion, atomic failure, and fresh-process resume tests.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-18-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/issues ./internal/evidence ./internal/validation ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-18-done.md
```

Every command is required and must terminate with session-captured evidence. The report-lint command runs before the next dispatch. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 12 and 17 close. QA Review window 6 follows tasks 16-18.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-18-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
