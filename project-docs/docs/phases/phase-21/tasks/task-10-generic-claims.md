# Phase 21 - Task 10: Generic Mutation Claims

## Objective

Implement generic Claim acquire, renew, expiry, release, list, conflict,
version, and idempotency semantics against `work_claims`. Replace the
Issue-specific `internal/claims` implementation in place after generic
behavior passes while retaining the table rename for task 21.

## Validation Contract Coverage

- VAL-022: Claim lifecycle binds exact identities and versions idempotently.
- VAL-023: one active mutation Claim per Issue and stale/wrong requests fail closed.

## Binding Constraints

- One-Issue mutation only; no supersession, parallel-scope proof, or recovery.
- Use an injected clock for deterministic expiry tests.
- Remove legacy Claim commands/actions after replacement tests pass, but retain
  `internal/claims` as the generic owner and do not drop legacy `claims` yet.

## Scope

### Files To Create Or Modify

- `internal/claims/`
- `internal/cli/`
- `cmd/fid/`
- `internal/journal/`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-10-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §7`
- `project-docs/specs/01_PRD.md PRD-RUN-001..004`

### Out Of Scope

- Claim supersession/recovery, parallel mutation, and migration 0014.

## Acceptance Criteria

- [ ] Acquire/renew/expire/release/list round-trip exact bindings, versions, and replay. (VAL-022)
- [ ] Concurrent/stale/wrong Actor/worktree cases reject with zero authorization/effects. (VAL-023)
- [ ] Issue-specific Claim leaves/package are removed after generic parity. (VAL-022)

## First Failing Test

Add a failing fake-clock Claim lifecycle table and concurrent same-Issue acquire test.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-10-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/claims ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-10-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after tasks 6-9 and QA Review 3 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-10-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
