# Phase 21 - Task 7: Envelope Authority Evaluator

## Objective

Implement the generic authority evaluator that intersects an activation-resolved Actor with exact Envelope grants. Make absence and every dimension mismatch deny atomically, without consulting V1 roles or recalculating kit policy.

## Validation Contract Coverage

- VAL-018: exact Actor/grant matches authorize only their recorded scope.
- VAL-019: missing, forged, inactive, expired, or mismatched authority denies without effects.

## Binding Constraints

- The evaluator accepts resolved Actor facts, never caller capabilities.
- Absence is deny and all failures use stable typed reasons.
- Do not remove V1 policy callers until their retained mutations are converted in tasks 8-9.

## Scope

### Files To Create Or Modify

- `internal/policy/` - replace V1 roles with the generic authority evaluator.
- `internal/workflow/`
- `internal/cli/`
- `project-docs/docs/phases/phase-21/tasks/task-7-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §6`
- `project-docs/specs/01_PRD.md PRD-AUTH-001..007`

### Out Of Scope

- Protected Human authority, kit risk evaluation, and caller migration.

## Acceptance Criteria

- [ ] The complete positive decision table authorizes exact matches. (VAL-018)
- [ ] Every independent mismatch denies with no Event/effect. (VAL-019)
- [ ] No V1 role lookup or caller role string participates. (VAL-018, VAL-019)

## First Failing Test

Add a failing table-driven authority test varying one grant/Actor dimension per row.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-7-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/policy ./internal/workflow
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-7-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after tasks 4-6 and QA Review 2 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-7-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
