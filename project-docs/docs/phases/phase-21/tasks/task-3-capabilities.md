# Phase 21 - Task 3: Capability Discovery

## Objective

Implement the migration-derived capability projection and `fid capabilities` transport. Keep every incomplete semantic capability unavailable with exact limits and trust sources.

## Validation Contract Coverage

- VAL-009: owning observable in the Phase 21 validation contract.
- VAL-010: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Availability is binary; there is no partial or experimental member.
- Caller input cannot override binary/migration-derived rows.
- Use exact operation/result/error contracts from `public-interface.md`.

## Scope

### Files To Create Or Modify

- `internal/workflow/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-3-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §2`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/specs/02_TRD.md §3.1`

### Out Of Scope

- Activation lifecycle and capability qualification beyond the matrix.

## Acceptance Criteria

- [ ] Discovery emits every exact sorted row. (VAL-009)
- [ ] Unavailable rows remain false despite internal operations. (VAL-010)

## First Failing Test

Add a failing full-matrix JSON golden and a false-capability conformance table.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-3-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-3-done.md
```

Every command is required and must terminate with session-captured evidence. The report-lint command runs before the next dispatch. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 1-2 close. QA Review window 1 follows tasks 1-3.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-3-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
