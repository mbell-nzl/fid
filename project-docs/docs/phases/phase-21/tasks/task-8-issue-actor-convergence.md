# Phase 21 - Task 8: Issue Actor And Event Convergence

## Objective

Convert retained Issue create/update/transition transport and Events to `--actor`, durable idempotency, and Event v2. Apply the same-clock ordering and reopened creation-size carry-forwards without changing the standalone Issue contract.

## Validation Contract Coverage

- VAL-025: owning observable in the Phase 21 validation contract.
- VAL-026: owning observable in the Phase 21 validation contract.
- VAL-027: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Standalone Actor is attribution only; Envelope-bound protected authority still uses task 7.
- Do not implement positive terminal-success Evidence in this task.
- CF-020 and CF-027 are triggered.

## Scope

### Files To Create Or Modify

- `internal/issues/`
- `internal/activity/`
- `internal/journal/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-8-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/public-interface.md §4.1`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`
- `project-docs/docs/current/decisions_active.md`

### Out Of Scope

- Other retained mutations, generic Claims, and terminal Evidence guard.

## Acceptance Criteria

- [ ] Issue mutations match exact Actor/key/Event/replay surface and reject `--by`. (VAL-025)
- [ ] Same-clock Events are nonfuture and sequence-ordered. (VAL-026)
- [ ] Creation file/function targets are met without behavior drift. (VAL-027)

## First Failing Test

Add the shared Issue v2/replay corpus, fixed-clock test, and size assertion before conversion.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-8-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/issues ./internal/activity ./internal/journal ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 2 and 7 close.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-8-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
