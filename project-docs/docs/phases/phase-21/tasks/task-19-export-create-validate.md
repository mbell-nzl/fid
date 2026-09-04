# Phase 21 - Task 19: Deterministic Export Create And Validate

## Objective

Create the deterministic v2 export archive and validate it without mutating project state. Keep import absent and portable capability unavailable.

## Validation Contract Coverage

- VAL-069: owning observable in the Phase 21 validation contract.
- VAL-070: owning observable in the Phase 21 validation contract.
- VAL-071: owning observable in the Phase 21 validation contract.

## Binding Constraints

- No wall-clock archive field or unstable ZIP metadata.
- Validate never repairs or imports.
- Do not advertise `export.portable.v1`.

## Scope

### Files To Create Or Modify

- `internal/export/`
- `internal/evidence/`
- `internal/journal/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-19-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §11`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md §10`

### Out Of Scope

- Import, identity/history reconstruction, migration 0014, and final parity.

## Acceptance Criteria

- [ ] Repeated create from identical facts is byte-identical and complete. (VAL-069)
- [ ] Every archive defect rejects without mutation. (VAL-070)
- [ ] Import is absent and capability false. (VAL-071)

## First Failing Test

Add failing repeated-digest/inventory test, corruption corpus, and import/capability absence guard.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-19-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/export ./internal/evidence ./internal/journal ./internal/cli ./cmd/fid/...
rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"import|newExportImport|export\.import' cmd/fid internal
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Expected Non-Zero Commands

- `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"import|newExportImport|export\.import' cmd/fid internal` - expected `exit_code: 1` - no export import command/path

## Notes For Developer

- Start after task 18 closes.
- This task creates `internal/export/`; its focused package command is valid only from this task onward.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-19-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
