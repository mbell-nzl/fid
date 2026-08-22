# Phase 21 - Task 6: Envelope, Obligation, And Budget Values

## Objective

Implement exact one-revision Work Envelope create/get/activate values with one selected Issue, generic Authority Grants, one complete gate-evidence Obligation, and immutable null-consumption budgets. Reject all adjacent deferred Envelope semantics.

## Validation Contract Coverage

- VAL-014: Envelope shape and one-Issue cardinality are exact.
- VAL-015: generic Obligation shape is complete with one supported type.
- VAL-016: immutable budget consumption and remaining are null.
- VAL-017: deferred Envelope/Decision/Stop/Attention inputs reject.

## Binding Constraints

- Registration must already exist before Envelope create.
- Persist references and contract values, never copied Issue lifecycle/prose.
- No revise, cancel, protected Decision, Stop, Attention, or derived-status mutation.

## Scope

### Files To Create Or Modify

- `internal/workflow/` - Work Envelope value and lifecycle files.
- `internal/cli/`
- `cmd/fid/`
- `internal/journal/`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-6-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§1,6`
- `project-docs/specs/01_PRD.md PRD-ENV-001..010`

### Out Of Scope

- Authority evaluation, Claims/Runs, completion, and other Obligation types.

## Acceptance Criteria

- [ ] Exact one-Issue Envelope round-trips without copied Issue lifecycle. (VAL-014)
- [ ] Obligation and budgets enforce exact schema/null semantics. (VAL-015, VAL-016)
- [ ] Every excluded input/leaf fails typed with no Event. (VAL-017)

## First Failing Test

Add a failing Envelope input-contract table covering the one accepted value and every excluded cardinality/type.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-6-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-6-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 5 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-6-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
