# Phase 21 - Task 7: Generic Authority Evaluator

## Objective

Evaluate exact Envelope Grants against the active Activation-resolved Actor and deny by default. Prove role-shaped project attribution cannot satisfy any protected operation.

## Validation Contract Coverage

- VAL-023: owning observable in the Phase 21 validation contract.
- VAL-024: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Absence is deny and every binding dimension must match.
- Do not consult `roles`, legacy role action maps, or Actor-looking historical strings.
- Retain `internal/policy` as the generic evaluator boundary.

## Scope

### Files To Create Or Modify

- `internal/policy/`
- `internal/workflow/`
- `internal/validation/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-7-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §6`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`
- `project-docs/specs/01_PRD.md PRD-AUTH-001..004`

### Out Of Scope

- Retained standalone mutation conversion and V1 policy deletion.

## Acceptance Criteria

- [ ] Only the all-match Grant row authorizes. (VAL-023)
- [ ] Every role-looking attribution remains non-authoritative. (VAL-024)

## First Failing Test

Add the failing one-dimension-at-a-time Grant matrix and seeded-role adversarial test.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-7-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/policy ./internal/workflow ./internal/validation ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-7-done.md
```

Every command is required and must terminate with session-captured evidence. The report-lint command runs before the next dispatch. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 4 and 6 close.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-7-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
