# Phase 5 - Task 1b: Close VAL-005 evidence gap (rebuild recovery + accepted one-line append) - Completion Report

## Summary of what was implemented

- Added the missing `VAL-005` evidence in `internal/issues/issue_test.go` only.
- Extended `TestCreateRejectsProjectionFailingWriteBeforeAppend` so the rejected-create fixture now deletes `.fid/fid.sqlite`, runs `journal.RebuildWorkspace`, and proves a later `Get` still succeeds after the rebuild.
- Strengthened `TestCreateAndGetRoundTrip` so an accepted create now asserts a journal line-count delta of exactly `+1`.
- Product code was unchanged because both new assertions passed against the existing Task 1 seam.

## Files created or modified

- `internal/issues/issue_test.go`
- `project-docs/docs/phases/phase-5/tasks/task-1b-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-5/tasks/task-1b-done.md` - added the Task 1b completion report.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` - verified-no-change. This task only adds evidence for the already-delivered validate-before-append seam; no flow narrative changed.

## Validation Contract Coverage

- `VAL-005` - PASS. `TestCreateRejectsProjectionFailingWriteBeforeAppend` now proves all three rejected-write checks in one fixture: projection-time rejection, zero journal-line append, and successful later reads both before and after a full rebuild (`internal/issues/issue_test.go:475`).
- `VAL-005` - PASS. `TestCreateAndGetRoundTrip` now asserts that an accepted create appends exactly one journal line via a before/after delta of `1` (`internal/issues/issue_test.go:648`).
- `VAL-005` - PASS. Existing seam evidence remains in place: `TestDryRunProjectEventRollsBackEffectsAndLeavesProjectionMetaUnchanged` and `TestDryRunProjectEventRollsBackOnProjectionFailure` prove the dry-run projection commits no rows and leaves `projection_meta` unchanged on both success and failure (`internal/journal/project_test.go:418`, `internal/journal/project_test.go:508`).
- `VAL-005` / `EVID-01` - PASS. Together, the new `internal/issues` assertions plus the existing `internal/journal` dry-run rollback tests close the `review-1` evidence gap called out by `ESC-003`.

## Tests run

- `go test ./internal/issues -run 'TestCreateRejectsProjectionFailingWriteBeforeAppend|TestCreateAndGetRoundTrip' -count=1` - PASS
- `go test ./internal/journal -run TestDryRunProjectEvent -count=1` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `c082152`

## Commit subject

- `test(journal): close VAL-005 evidence gap`
