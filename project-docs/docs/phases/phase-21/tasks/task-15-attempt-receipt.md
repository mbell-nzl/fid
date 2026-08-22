# Phase 21 - Task 15: Execution Attempts And Receipts

## Objective

Persist an Attempt before spawn and at most one immutable terminal Receipt after observed termination, including exact process non-success outcomes, the frame-only handled `fid-error` case, and the streaming terminal response.

## Validation Contract Coverage

- VAL-051: owning observable in the Phase 21 validation contract.
- VAL-052: owning observable in the Phase 21 validation contract.
- VAL-053: owning observable in the Phase 21 validation contract.
- VAL-054: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Wrapper completion or missing exit observation is not termination.
- A non-success executable still records its Receipt and exits with class 8.
- An Attempt alone never satisfies Evidence.

## Scope

### Files To Create Or Modify

- `internal/evidence/`
- `internal/journal/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-15-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§8,14`
- `project-docs/docs/phases/phase-21/public-interface.md §§4,7,8`

### Out Of Scope

- Artifact retention and Evidence linking.

## Acceptance Criteria

- [ ] Attempt precedes every spawn. (VAL-051)
- [ ] Success Receipt is exact and immutable. (VAL-052)
- [ ] All non-success outcomes/exit classes are exact. (VAL-053)
- [ ] Attempt-only/caller proof remains unsatisfied. (VAL-054)

## First Failing Test

Add injected append/spawn ordering, success, four-outcome, and missing-Receipt tests.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-15-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 13-14 close. QA Review window 5 follows tasks 13-15.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-15-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
