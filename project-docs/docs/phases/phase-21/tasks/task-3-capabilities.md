# Phase 21 - Task 3: Capability Discovery Matrix

## Objective

Implement capability discovery as a deterministic binary/migration-derived read. Return the exact 12-row availability, limit, and trust-source matrix and reject any caller attempt to upgrade or weaken it.

## Validation Contract Coverage

- VAL-007: discovery returns the complete exact capability matrix.
- VAL-008: incomplete capabilities remain unavailable without caller override.

## Binding Constraints

- Walking-skeleton command existence does not imply semantic capability availability.
- Rows sort by capability ID and use canonical `fid-cli.v1` output.
- No environment/config toggle changes semantic availability.

## Scope

### Files To Create Or Modify

- `internal/project/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-3-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §2`
- `project-docs/specs/02_TRD.md §3`

### Out Of Scope

- Activation and any unavailable capability implementation.

## Acceptance Criteria

- [ ] Clean, restarted, and rebuilt projects return the exact matrix. (VAL-007)
- [ ] Every override/unknown schema request fails typed and cannot enable a row. (VAL-008)

## First Failing Test

Add a failing command golden test for all 12 rows and their deterministic order.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-3-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/project ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-3-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after tasks 1-2 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-3-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
