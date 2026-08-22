# Phase 21 - Task 9: Retained Inventory Attribution And Judgment Freeze

## Objective

Convert every retained non-Issue inventory mutation to project-local `--actor` plus Event-v2 idempotency, remove the public V1 policy command after generic authority passes, isolate the legacy action registry to enumerated dev/completion callers until task 20, and make Phase 23 judgment mutations fail closed. Rename stats grouping without treating historical attribution as authority.

## Validation Contract Coverage

- VAL-028: owning observable in the Phase 21 validation contract.
- VAL-029: owning observable in the Phase 21 validation contract.
- VAL-030: owning observable in the Phase 21 validation contract.
- VAL-031: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Cover workflow define/assign, AC add, labels, links, requirements, activity
  normalization, concerns, maintenance vacuum, and stats.
- No public role/judgment requirement authoring survives.
- Do not delete the legacy policy implementation while the enumerated
  dev/completion callers remain; prove new generic paths cannot reach it.
- Do not delete replacement-bound dev/completion packages or tables yet.

## Scope

### Files To Create Or Modify

- `internal/policy/`
- `internal/workflow/`
- `internal/issues/`
- `internal/labels/`
- `internal/links/`
- `internal/requirements/`
- `internal/activity/`
- `internal/stats/`
- `internal/maintenance/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-9-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/public-interface.md §§4-5`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`

### Out Of Scope

- Generic Claims/Runs and final V1 removal.

## Acceptance Criteria

- [ ] Every retained family has exact attributed/idempotent behavior. (VAL-028)
- [ ] Role/judgment workflow inputs cannot be authored. (VAL-029)
- [ ] Judgment mutations are unavailable/no-writer. (VAL-030)
- [ ] Stats uses exact grouping/attribution contract. (VAL-031)

## First Failing Test

Add one failing conformance row per retained/fail-closed family before shared conversion.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-9-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/policy ./internal/workflow ./internal/issues ./internal/labels ./internal/links ./internal/requirements ./internal/activity ./internal/stats ./internal/maintenance ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-9-done.md
```

Every command is required and must terminate with session-captured evidence. The report-lint command runs before the next dispatch. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 7-8 close. QA Review window 3 follows tasks 7-9.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-9-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
