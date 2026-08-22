# Phase 21 - Task 19: Deterministic Export Create And Validate

## Objective

Implement deterministic archive creation and non-mutating validation for the accepted manifest, journal prefix, source metadata, capability snapshot, and available Artifact blobs. Keep import absent and `export.portable.v1` unavailable.

## Validation Contract Coverage

- VAL-051: export creation is byte-deterministic and integrity-bound.
- VAL-052: validation detects safety/integrity defects without mutation.
- VAL-053: import is absent and portable capability remains unavailable.

## Binding Constraints

- Fixed ZIP ordering, timestamp, mode, extras, and compression are contract data.
- Validation must not initialize, append, rebuild, repair, or alter the target project.
- Do not claim PRD-REC-001/004 reconstruction parity before Phase 22 import.

## Scope

### Files To Create Or Modify

- `internal/export/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-19-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§2,11`
- `project-docs/specs/01_PRD.md PRD-REC-001..004`

### Out Of Scope

- Export import, reconstruction, identity/history preservation, and portable capability availability.

## Acceptance Criteria

- [ ] Repeated identical export bytes/digests and exact manifest inventory match. (VAL-051)
- [ ] Every corruption/path/inventory defect rejects with unchanged project state. (VAL-052)
- [ ] No import path exists and capability stays false. (VAL-053)

## First Failing Test

Add a failing deterministic two-export byte comparison, then one corruption row per validation rule.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-19-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/export ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-19-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after tasks 16-18 and QA Review 6 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-19-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
