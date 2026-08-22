# Phase 21 - Task 21: V1 Authority And Execution Removal

## Objective

After every generic replacement passes, remove the remaining V1 dev/gate/completion surfaces, apply migration 0014, and prove frozen legacy tables have neither writer nor authoritative reader. Finish the exact source, command, Event, schema, and writer removal inventory.

## Validation Contract Coverage

- VAL-005: migration 0014 performs the exact accepted drop/rename sequence.
- VAL-056: all owned V1 authority/execution surfaces are absent.
- VAL-057: frozen completion/telemetry tables have no Phase 21 writer or authority read.

## Binding Constraints

- Do not drop `task_completions` or `agent_runs`; document and test their frozen status.
- Remove invalidated callers/tests with each implementation; do not retain compatibility aliases.
- Do not modify the protected removal inventory or migration semantics to accommodate failures.

## Scope

### Files To Create Or Modify

- `cmd/fid/`
- `internal/dev/`
- `internal/runs/`
- `internal/issues/complete.go`
- `internal/issues/complete_support.go`
- `internal/issues/complete_test.go`
- `migrations/0014_p21_retire_v1_authority_execution.up.sql`
- `migrations/0014_p21_retire_v1_authority_execution.down.sql`
- `internal/db/`
- `internal/maintenance/`
- `project-docs/docs/phases/phase-21/tasks/task-21-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §15`
- `project-docs/specs/03_TRD_SCHEMA.md §§6,8.2`
- `internal/policy/` - verify only the generic evaluator remains.
- `internal/claims/` - verify only generic Claim operations remain.

### Out Of Scope

- Phase 23/25 deletion of frozen tables and any legacy compatibility mode.

## Acceptance Criteria

- [ ] Migration 0014 drops V1 collision tables and renames generic Claims exactly. (VAL-005)
- [ ] V1 implementations/leaves/Events/readers/writers and obsolete
  `internal/dev`/`internal/runs` packages are absent after generic parity. (VAL-056)
- [ ] Frozen tables remain structurally only, with no Phase 21 write/authority read. (VAL-057)

## First Failing Test

Add failing removal/conformance guards before deleting each remaining V1 surface and enabling migration 0014.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-21-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/db ./internal/maintenance ./internal/cli ./cmd/fid/...
ls internal/dev internal/runs
ls cmd/fid/dev.go cmd/fid/policy.go cmd/fid/issue_claim.go cmd/fid/issue_complete.go
ls internal/issues/complete.go internal/issues/complete_support.go
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-21-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Expected Non-Zero Commands

- `ls internal/dev internal/runs` - expected `exit_code: 2` - obsolete packages must not exist
- `ls cmd/fid/dev.go cmd/fid/policy.go cmd/fid/issue_claim.go cmd/fid/issue_complete.go` - expected `exit_code: 2` - removed paths must not exist
- `ls internal/issues/complete.go internal/issues/complete_support.go` - expected `exit_code: 2` - removed paths must not exist

## Notes For Developer

- Start only after tasks 7-20 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-21-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
