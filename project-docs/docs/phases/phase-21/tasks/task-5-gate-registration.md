# Phase 21 - Task 5: Activated Gate Registration

## Objective

Register and resolve one immutable Gate Definition only from finalized Activation values. Validate the full portable policy at registration and version deterministic bindings.

## Validation Contract Coverage

- VAL-016: owning observable in the Phase 21 validation contract.
- VAL-017: owning observable in the Phase 21 validation contract.
- VAL-018: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Registration precedes Envelope creation.
- No runtime argv/policy override is accepted.
- Stronger unsupported restrictions reject; no OS-containment claim.

## Scope

### Files To Create Or Modify

- `internal/evidence/`
- `internal/workflow/`
- `internal/policy/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-5-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §5`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/specs/04_TRD_CONFIG.md`

### Out Of Scope

- Process execution and live output.

## Acceptance Criteria

- [ ] Only finalized activated Definitions register. (VAL-016)
- [ ] Each unsupported policy field rejects without persistence. (VAL-017)
- [ ] Same binding reuses identity; changed binding versions immutably. (VAL-018)

## First Failing Test

Add the failing activated-source/order test, then policy and version tables.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-5-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 4 closes.
- This task creates `internal/evidence/`; the focused package command is valid only from this task onward.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-5-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
