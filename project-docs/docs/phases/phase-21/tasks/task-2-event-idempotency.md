# Phase 21 - Task 2: Event V2 And Durable Idempotency

## Objective

Introduce the canonical Event v2 envelope and operation/project-scoped idempotency foundation used by every later Phase 21 mutation. Prove serialization, aggregate sequencing, replay, digest conflict, and concurrent single-commit behavior before any new aggregate depends on them.
Use the task-1 schema to prove live/rebuild equivalence for the shared fixture.

## Validation Contract Coverage

- VAL-001: converted mutations use the exact Event v2 envelope.
- VAL-002: aggregate sequences remain contiguous and rejected mutations append nothing.
- VAL-003: same-key/same-request retries replay one stored result.
- VAL-004: changed or concurrent idempotency digests fail closed.
- VAL-006: live and rebuilt shared projections are identical.

## Binding Constraints

- Do not convert product families outside the focused fixture used to prove the shared primitive; owning tasks convert their families.
- Store the canonical response and idempotency row atomically with the Event/effects transaction.
- Event v1 development journals require no upgrader under ADR-0011.

## Scope

### Files To Create Or Modify

- `internal/journal/`
- `internal/db/`
- `internal/project/`
- `project-docs/docs/phases/phase-21/tasks/task-2-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§3,14`
- `project-docs/specs/03_TRD_SCHEMA.md §8.2`
- `project-docs/docs/phases/phase-21/qa/validation-contract.md`

### Out Of Scope

- Activation, capability discovery, and wholesale family conversion.

## Acceptance Criteria

- [ ] Canonical Event v2 round-trips and rejects missing/unknown/invalid aggregate or idempotency fields. (VAL-001)
- [ ] Committed aggregate sequences are contiguous under success/conflict concurrency. (VAL-002)
- [ ] Idempotent replay and changed-digest conflict have the exact durable semantics. (VAL-003, VAL-004)
- [ ] The shared v2 fixture has byte-identical live/rebuild reads. (VAL-006)

## First Failing Test

Add a failing golden Event-v2 decoder/encoder test, followed by same-key concurrency tests.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-2-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/journal ./internal/db/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-2-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 1 closes with the exact schema inventory.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-2-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
