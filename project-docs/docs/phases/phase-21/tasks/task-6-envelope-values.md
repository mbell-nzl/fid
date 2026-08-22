# Phase 21 - Task 6: One-Revision Envelope Values

## Objective

Create/get/activate the exact one-Issue, one-revision Envelope with prior Registration, generic gate-evidence Obligation, Authority Grants, sources, and immutable budgets. Fail closed on every deferred widening input.

## Validation Contract Coverage

- VAL-019: owning observable in the Phase 21 validation contract.
- VAL-020: owning observable in the Phase 21 validation contract.
- VAL-021: owning observable in the Phase 21 validation contract.
- VAL-022: owning observable in the Phase 21 validation contract.

## Binding Constraints

- One primary Issue, no related Issues, revision 1, and one gate-evidence Obligation only.
- The Obligation guards both terminal-success transition and completion.
- Routine is external input lane data, not fid lifecycle state.

## Scope

### Files To Create Or Modify

- `internal/workflow/`
- `internal/validation/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-6-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§1,6`
- `project-docs/docs/phases/phase-21/public-interface.md §§2-4`
- `project-docs/specs/01_PRD.md PRD-ENV-001..008`

### Out Of Scope

- Completion, revision, cancellation, Stops, Attention, and budget consumption.

## Acceptance Criteria

- [ ] Exact one-revision boundary and prior Registration hold. (VAL-019)
- [ ] Obligation shape/dual guards are exact. (VAL-020)
- [ ] Budget values remain null/immutable. (VAL-021)
- [ ] Deferred values reject without Event. (VAL-022)

## First Failing Test

Add the failing closed Envelope input round-trip and widening rejection table.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-6-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/validation ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-6-done.md
```

Every command is required and must terminate with session-captured evidence. The report-lint command runs before the next dispatch. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 5 closes. QA Review window 2 follows tasks 4-6.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-6-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
