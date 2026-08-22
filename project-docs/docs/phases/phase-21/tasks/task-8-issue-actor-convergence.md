# Phase 21 - Task 8: Issue Actor And Idempotency Convergence

## Objective

Convert retained standalone Issue mutations and execution-aware Issue transition request plumbing to Event v2, `--actor`, and durable idempotency without granting kit trust. Preserve the accepted Phase 20 Issue contract outside the explicit Phase 21 changes.

## Validation Contract Coverage

- VAL-001: converted Issue Events use Event v2.
- VAL-003: retryable Issue mutations replay durably.
- VAL-020: standalone Issue attribution remains available but untrusted.
- VAL-062: burst Issue Events use observed time without synthetic future dates.
- VAL-063: reopened Issue creation code meets structure budgets with parity.

## Binding Constraints

- Standalone attribution is a project-local fact and cannot resolve an Activation profile.
- Preserve Phase 20 JSON shapes and revision rules except the exact Actor/idempotency additions.
- Do not implement the execution-bound transition decision until task 12.
- CF-020 resolves ordering through Event/aggregate sequence, not future-dated
  occurrence timestamps; CF-027 requires behavior-preserving `create.go`
  decomposition in this task.

## Scope

### Files To Create Or Modify

- `internal/issues/`
- `internal/activity/`
- `internal/maintenance/`
- `internal/cli/`
- `cmd/fid/issue*.go`
- `internal/journal/`
- `project-docs/docs/phases/phase-21/tasks/task-8-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§1,3,12,14`
- `project-docs/docs/phases/phase-20/contracts.md`

### Out Of Scope

- Non-Issue retained mutations, V1 policy deletion, and transition authorization.

## Acceptance Criteria

- [ ] Issue mutations emit v2 Actor/idempotency facts and replay once. (VAL-001, VAL-003)
- [ ] Standalone `--actor` persists attribution but cannot authorize kit operations. (VAL-020)
- [ ] Burst writes remain current without synthetic future timestamps. (VAL-062)
- [ ] Reopened Issue creation code meets the recorded structure targets with unchanged conformance behavior. (VAL-063)
- [ ] Unchanged Phase 20 Issue behavior remains conformant. (VAL-020)

## First Failing Test

First add failing fixed-clock burst/recent-activity and creation-structure
guards, then the Issue-create retry test asserting Event-v2 Actor provenance
and no kit authority.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-8-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/issues ./internal/activity ./internal/maintenance ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-8-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 7 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-8-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
