# Phase 21 - Task 5: Activated Gate Registration

## Objective

Create immutable Gate Registrations only from finalized activated Gate Definition values. Bind definition and execution policy to project, repository, and scope before any Envelope can reference the Registration.

## Validation Contract Coverage

- VAL-012: Gate Registration resolves activated values and precedes Envelope creation.
- VAL-013: Registration content and binding are immutable and versioned.

## Binding Constraints

- Do not execute processes in this task.
- Registration accepts no argv, cwd, environment, limits, or policy override from the caller.
- Policy validation is structural here; portable enforcement lands in task 13.

## Scope

### Files To Create Or Modify

- `internal/evidence/`
- `internal/cli/`
- `cmd/fid/`
- `internal/journal/`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-5-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §5`
- `project-docs/specs/03_TRD_SCHEMA.md §8.2`

### Out Of Scope

- Executor implementation, Attempts, Receipts, and Evidence linking.

## Acceptance Criteria

- [ ] Only a finalized active definition registers and inline/unactivated input rejects. (VAL-012)
- [ ] Same content/binding reuses identity; changed digest creates a new version without mutation. (VAL-013)

## First Failing Test

Add a failing registration-order test proving inline and unfinalized definitions cannot register.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-5-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-5-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 4 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-5-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
