# Phase 21 - Task 12: Shared Evaluator And Nonterminal Transition

## Objective

Implement the execution-bound nonterminal transition predicate and one shared Obligation evaluator used fail-closed by both terminal-success transition and completion. Positive trusted Evidence remains Group 2 work.

## Validation Contract Coverage

- VAL-038: owning observable in the Phase 21 validation contract.
- VAL-039: owning observable in the Phase 21 validation contract.
- VAL-040: owning observable in the Phase 21 validation contract.
- VAL-041: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Terminal-success cannot succeed in this task because qualifying Evidence does not exist yet.
- Return all reasons in fixed kebab-case order with atomic failure.
- Apply CF-032 invalid revision-effect tests first.

## Scope

### Files To Create Or Modify

- `internal/workflow/`
- `internal/issues/`
- `internal/policy/`
- `internal/validation/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-12-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§6,7,10`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`

### Out Of Scope

- Evidence linking, positive terminal transition, and completion mutation.

## Acceptance Criteria

- [ ] Only all-qualified nonterminal transition commits. (VAL-038)
- [ ] Failures return all reasons with no effect. (VAL-039)
- [ ] Both guarded actions share the same unsatisfied evaluator result. (VAL-040)
- [ ] Invalid revision effects reject. (VAL-041)

## First Failing Test

Add CF-032 tests first, then the nonterminal predicate and dual-call fail-closed evaluator goldens.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-12-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 7 and 10-11 close. QA Review window 4 and the independent Group 1 checkpoint follow; task 13 cannot start without PASS.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-12-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
