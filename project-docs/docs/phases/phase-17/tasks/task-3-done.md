# Phase 17 - Task 3: Implement Maintenance Structural Error Checks - Completion Report

## Summary of what was implemented

- Created the `internal/maintenance` domain boundary with a read-only structural audit over the SQLite projection and journal.
- Added deterministic error findings for v0.1 hierarchy, dependent-row kind, sequence/scheduling, link, claim, activity-reference, activity journal-event, foreign-key, and workflow/state violations.
- The audit returns every finding ordered by check, entity ID, then message; it performs no repair, projection write, or journal append.

## Files created or modified

- `internal/maintenance/check.go`
- `internal/maintenance/errors.go`
- `internal/maintenance/check_test.go`
- `internal/maintenance/test_helpers_test.go`
- `AGENTS.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/docs/phases/phase-17/tasks/task-3-done.md`

## Documentation updated

- `AGENTS.md` — records `internal/maintenance` as a core/domain boundary.
- `project-docs/specs/02_TRD.md` — records the maintenance package ownership and write prohibitions.
- `project-docs/specs/03_TRD_SCHEMA.md` — documents v0.1 maintenance structural-audit/trigger parity and activity journal-event validation.

## Validation Contract Coverage

- VAL-003: PASS — tests seed every active v0.1 hierarchy, dependent-row kind, and sequence/scheduling error family plus valid counterparts.
- VAL-004: PASS — tests cover orphaned and resolved links, claims, all supported activity ref types, and present/missing activity journal-event IDs.
- VAL-005: PASS — test-only FK bypass seeds two `foreign_key_check` rows; the audit preserves both with table/row/parent identity, while the valid fixture yields none.
- VAL-006: PASS — the valid fixture exercises every state in the three default workflows; unknown workflow and state pairs are reported through `internal/workflow.CategoryForState`.

## Tests run

- `go test ./internal/maintenance` - PASS (initial red: package/types did not exist; green after implementation)
- `go test -count=1 ./internal/maintenance` - PASS (per-row missing journal-event regression)
- `go test -v -count=1 ./internal/maintenance` - PASS (focused structural families, valid fixture, ordering, and read-only checks)
- Runtime entrypoint - N/A: the `fid maintenance check` CLI is explicitly Task 5 scope; Task 3 exposes and tests the domain API only.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS (total: 76.4%; `internal/maintenance`: 77.7%)
- Conditional commands - N/A (the task spec lists none).

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
