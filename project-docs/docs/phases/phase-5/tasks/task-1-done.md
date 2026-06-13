# Phase 5 - Task 1: Validate-before-append write-path seam (CF-007) - Completion Report

## Summary of what was implemented

- Added `journal.DryRunProjectEvent`, a pre-append projection seam that applies an event inside a transaction and always rolls it back, surfacing the same projection-time trigger or constraint failures without committing projected rows or advancing `projection_meta`.
- Routed `internal/issues.Create` through the dry-run seam before `AppendEvent`, so projection-rejected creates now fail before the canonical journal append while the existing post-append `ProjectEvent` path remains unchanged for accepted writes.
- Added regressions for VAL-005 at both layers: `internal/issues` now proves a rejected create leaves the journal line count unchanged and does not poison later reads, and `internal/journal` proves dry runs roll back on both success and failure.

## Files created or modified

- `internal/journal/project.go`
- `internal/journal/project_test.go`
- `internal/issues/issue.go`
- `internal/issues/issue_test.go`
- `project-docs/docs/phases/phase-5/tasks/task-1-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-5/tasks/task-1-done.md` - added the Task 1 completion report and VAL-005 evidence.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - verified-no-change. The existing flow already matches TRD §4.3's validate/build-before-append contract; Task 1 fixed implementation conformance rather than changing the user flow narrative.

## Validation Contract Coverage

- `VAL-005` - PASS. `TestCreateRejectsProjectionFailingWriteBeforeAppend` proves a projection-rejected create appends zero journal lines and a later `Get` of an unrelated valid issue still succeeds. `TestDryRunProjectEventRollsBackEffectsAndLeavesProjectionMetaUnchanged` and `TestDryRunProjectEventRollsBackOnProjectionFailure` prove the pre-append seam never commits dry-run effects and never advances `projection_meta`.
- `VAL-005` - PASS. Accepted writes still use the real append-then-project path, and the full gate preserved `internal/issues` coverage at 85.6% and `internal/journal` coverage at 80.9%, satisfying the task's COV-001 requirement.

## Tests run

- `go test ./internal/issues -run TestCreateRejectsProjectionFailingWriteBeforeAppend -count=1` - PASS after initial RED (`journal lines after rejected create = 4, want unchanged 3`).
- `go test ./internal/journal -run TestDryRunProjectEvent -count=1` - PASS after initial RED (undefined `DryRunProjectEvent` seam).
- `go test ./internal/journal ./internal/issues` - PASS.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency Budget / Dependency-budget note

- No dependency changes.

## Commit hash

- `263d523`

## Commit subject

- `fix(journal): validate issue creates before append`
