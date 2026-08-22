# Phase 21 - Task 15: Trusted Attempt And Receipt Execution

## Objective

Wire `gate run` to append an Attempt before standard-library process spawn and commit at most one immutable terminal Receipt after observed termination. Cover success, failure, timeout, signal, spawn failure, output metadata, and durable idempotent replay.

## Validation Contract Coverage

- VAL-037: Attempt durability precedes process spawn.
- VAL-038: one immutable Receipt records each exact observed terminal outcome.
- VAL-039: absent or caller-authored Receipts never pass evidence.

## Binding Constraints

- The caller supplies only Registration, Run, Actor, repository observation, and idempotency key.
- A timeout waits for observed termination; wrapper/session completion is not process termination.
- Receipt success is determined only from registered success exit codes and observed outcome.

## Scope

### Files To Create Or Modify

- `internal/evidence/` - executor, Attempt, and Receipt files.
- `internal/cli/`
- `cmd/fid/`
- `internal/journal/`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-15-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§5,8,14`
- `project-docs/specs/01_PRD.md PRD-GATE-003..009`

### Out Of Scope

- Artifact retention and Evidence linking.

## Acceptance Criteria

- [ ] Injected append failure spawns nothing and every observed spawn has a prior Attempt. (VAL-037)
- [ ] Terminal outcome matrix creates one exact immutable Receipt. (VAL-038)
- [ ] Missing/caller-authored proof remains non-passing and execution replay does not respawn. (VAL-039)

## First Failing Test

Add a failing process-fixture trace proving Attempt commit occurs before the fixture observes spawn.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-15-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-15-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after tasks 13-14 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-15-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
