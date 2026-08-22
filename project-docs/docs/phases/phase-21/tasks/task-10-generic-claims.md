# Phase 21 - Task 10: Generic Mutation Claims

## Objective

Replace the Claim package behavior and CLI with exact generic acquire/renew/release/list operations under active Envelope authority. Use `work_claims` until migration 0014.

## Validation Contract Coverage

- VAL-032: owning observable in the Phase 21 validation contract.
- VAL-033: owning observable in the Phase 21 validation contract.
- VAL-034: owning observable in the Phase 21 validation contract.

## Binding Constraints

- One active mutation Claim per Issue; no force/supersession/recovery.
- Exact Actor/worktree/version/TTL bindings are mandatory.
- Old Issue claim/unclaim spelling is not retained as an alias.

## Scope

### Files To Create Or Modify

- `internal/claims/`
- `internal/policy/`
- `internal/workflow/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-10-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §7`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`

### Out Of Scope

- Final table rename/removal and recovery.

## Acceptance Criteria

- [ ] Acquire binds every accepted dimension and one-active invariant. (VAL-032)
- [ ] Conflict/expiry/wrong-holder paths fail closed. (VAL-033)
- [ ] Renew/release versioning and replay are exact. (VAL-034)

## First Failing Test

Add failing acquire binding, conflict/expiry, and version lifecycle tables.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-10-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/claims ./internal/policy ./internal/workflow ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 6-7 close.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-10-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
