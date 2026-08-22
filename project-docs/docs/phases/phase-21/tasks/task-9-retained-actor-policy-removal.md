# Phase 21 - Task 9: Retained Actor Convergence And Policy Removal

## Objective

Convert every remaining retained mutation to `--actor`, Event v2, and required
idempotency, rename stats grouping to `--group-by`, then delete the V1 policy
command/implementation after generic authority conformance while retaining
`internal/policy` as the generic owner. Remove the mixed Phase 20 global-flag
restriction.

## Validation Contract Coverage

- VAL-001: every retained mutation family emits Event v2.
- VAL-003: retained retryable mutations use durable replay.
- VAL-021: no live `--by` spelling remains.

## Binding Constraints

- Convert callers before deleting V1 role-policy files; retain the generic
  `internal/policy` package and leave `roles` for migration 0014 in task 21.
- Help, errors, fixtures, and integration commands must use the replacement spellings consistently.
- Do not treat stats grouping as Actor attribution.

## Scope

### Files To Create Or Modify

- `cmd/fid/`
- `internal/cli/`
- `internal/workflow/`
- `internal/labels/`
- `internal/links/`
- `internal/requirements/`
- `internal/maintenance/`
- `internal/policy/`
- `project-docs/docs/phases/phase-21/tasks/task-9-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§12,14,15`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md §§1,8`

### Out Of Scope

- Legacy Claim/Run/dev/completion removals and dropping `roles`.

## Acceptance Criteria

- [ ] All retained mutation families use resolved/attributed Actor and Event v2/idempotency as contracted. (VAL-001, VAL-003)
- [ ] Mutation `--actor` and stats `--group-by` work; `--by` rejects everywhere. (VAL-021)
- [ ] V1 policy command/role implementation is absent only after generic authority tests pass. (VAL-021)

## First Failing Test

Add a failing command-tree matrix proving replacement spellings and rejection of every former `--by` use.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-9-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/labels ./internal/links ./internal/requirements ./internal/maintenance ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-9-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after tasks 7-8 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-9-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
