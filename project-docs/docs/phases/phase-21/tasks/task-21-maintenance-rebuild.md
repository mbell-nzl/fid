# Phase 21 - Task 21: Final-Schema Maintenance, Init, Rebuild, And Export Parity

## Objective

Against the post-0014 binary, prove fresh init and delete-projection/rebuild are role-free, extend deterministic maintenance findings, compare all named live/rebuilt reads and export bytes, and run the 100k Event-v2 check.

## Validation Contract Coverage

- VAL-076: owning observable in the Phase 21 validation contract.
- VAL-077: owning observable in the Phase 21 validation contract.
- VAL-078: owning observable in the Phase 21 validation contract.
- VAL-079: owning observable in the Phase 21 validation contract.
- VAL-080: owning observable in the Phase 21 validation contract.
- VAL-081: owning observable in the Phase 21 validation contract.
- VAL-082: owning observable in the Phase 21 validation contract.
- VAL-083: owning observable in the Phase 21 validation contract.

## Binding Constraints

- This task starts only after destructive removal; it cannot use a pre-0014 bootstrap fixture.
- Maintenance reports but never repairs.
- Record hardware/toolchain and observed termination for the 100k threshold.

## Scope

### Files To Create Or Modify

- `internal/bootstrap/`
- `internal/maintenance/`
- `internal/journal/`
- `internal/db/`
- `internal/export/`
- `internal/cli/`
- `cmd/fid/`
- `cmd/fid/integration/`
- `project-docs/docs/phases/phase-21/tasks/task-21-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§14-16`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`
- `project-docs/specs/03_TRD_SCHEMA.md §8.2`

### Out Of Scope

- Import, repair, held-out scenarios, and final conformance journey.

## Acceptance Criteria

- [ ] Final-binary fresh init succeeds without roles. (VAL-076)
- [ ] Delete/rebuild reproduces role-free bootstrap/capabilities. (VAL-077)
- [ ] All three atomic maintenance classes are exact. (VAL-078, VAL-079, VAL-080)
- [ ] Named reads and export are live/rebuild identical. (VAL-081, VAL-082)
- [ ] 100k v2 rebuild terminates under threshold with parity. (VAL-083)

## First Failing Test

Add the final-binary temp-project init/rebuild test first, then corrupt fixtures, parity snapshot, and bounded benchmark.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-21-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/bootstrap ./internal/maintenance ./internal/journal ./internal/db ./internal/export ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
go test -bench . -benchmem ./...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-21-done.md
```

Every command is required and must terminate with session-captured evidence. The report-lint command runs before the next dispatch. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 20 closes. QA Review window 7 follows tasks 19-21.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-21-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
