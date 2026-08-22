# Phase 21 - Task 20: Maintenance And Full Rebuild Parity

## Objective

Extend maintenance to all Phase 21 integrity classes and prove full live/rebuild equivalence across capabilities, Envelope, evidence, completion, resume, and export. Maintenance reports facts only and never repairs state.

## Validation Contract Coverage

- VAL-054: maintenance detects every required integrity finding class deterministically.
- VAL-055: rebuild reproduces all named Phase 21 semantics and export bytes.
- VAL-064: a 100k-Event v2 rebuild meets the accepted performance target.

## Binding Constraints

- An Attempt without Receipt is non-passing, but only becomes abandoned under the accepted grace observation.
- Frozen-table writes and removed V1 state are integrity findings.
- Do not add automatic repair, import, or recovery.

## Scope

### Files To Create Or Modify

- `internal/maintenance/`
- `internal/db/`
- `internal/journal/`
- `internal/export/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-20-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §16`
- `project-docs/specs/01_PRD.md PRD-FOUND-003 and PRD-OPS-003`

### Out Of Scope

- Repair actions, baseline consolidation, and import parity.

## Acceptance Criteria

- [ ] Each corrupt fixture yields the exact typed/order-stable finding. (VAL-054)
- [ ] Delete projection/rebuild reproduces all named reads and archive bytes. (VAL-055)
- [ ] Recorded reference hardware rebuilds 100k v2 Events under 30 seconds with parity. (VAL-064)

## First Failing Test

Add one failing corruption fixture per maintenance class, then the full state snapshot rebuild comparison.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-20-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/maintenance ./internal/db ./internal/journal ./internal/export ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
go test -bench . -benchmem ./...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-20-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 19 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-20-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
