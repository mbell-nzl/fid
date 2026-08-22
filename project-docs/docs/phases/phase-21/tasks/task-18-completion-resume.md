# Phase 21 - Task 18: Envelope Completion And Resume

## Objective

Implement atomic Envelope completion through the shared evaluator and a bounded durable resume projection. Completion must require terminal-success Issue plus satisfied current Obligation; resume must reconstruct facts and next actions without chat, mutation, recovery, or respawn.

## Validation Contract Coverage

- VAL-048: completion succeeds only when Issue and Obligation qualify.
- VAL-049: failed completion returns all reasons with no partial or legacy state.
- VAL-050: resume is bounded, durable, read-only, and rebuild-stable.

## Binding Constraints

- Completion never writes or reads `task_completions` as authority.
- Use the task-17 evaluator; do not duplicate readiness rules.
- Resume lists only supported explicit next actions and makes no liveness guess.

## Scope

### Files To Create Or Modify

- `internal/workflow/`
- `internal/validation/`
- `internal/report/`
- `internal/cli/`
- `cmd/fid/`
- `internal/journal/`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-18-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§6,10`
- `project-docs/specs/01_PRD.md PRD-ENV-008 and PRD-RUN-014`

### Out Of Scope

- Recovery, successor Runs, Attention/Stops, and Human validation brief.

## Acceptance Criteria

- [ ] Only terminal-success plus satisfied applicable Obligation completes once. (VAL-048)
- [ ] Every failed completion is atomic and complete in reasons. (VAL-049)
- [ ] Fresh-process resume is byte-stable across rebuild and appends nothing. (VAL-050)

## First Failing Test

Add a failing completion condition matrix, then a fresh-process resume before/after rebuild test.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-18-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/validation ./internal/report ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-18-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 17 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-18-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
