# Phase 5 - Task 3: `fid issue update` with activity-log audit - Completion Report

## Summary of what was implemented

- Added `internal/issues.Update` to mutate supported issue fields through the
  Task 1 validate-before-append seam.
- Implemented one-event `issue.updated` journaling with an `issues` update
  effect, label delta effects, and one `activity` effect projected to
  `activity_log`.
- Wired `fid issue update <id> ... --by <role>` in `cmd/fid` as a thin adapter
  over the domain service, including deterministic stderr confirmation text.
- Documented and tested the chosen no-op behavior: updates that do not change
  any field are rejected and append no journal event.

## Files created or modified

- `internal/issues/update.go`
- `internal/issues/update_test.go`
- `cmd/fid/issue.go`
- `cmd/fid/issue_update_test.go`
- `internal/journal/writer.go`
- `internal/journal/project.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-5/tasks/task-3-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - updated the `fid issue update`
  section to record no-op rejection semantics and align the stderr example with
  the shipped ASCII `->` confirmation format.

## Validation Contract Coverage

- `VAL-001` - PASS. `TestUpdateChangesTitleAndPriorityAndPreservesUntouchedFields`
  proves valid updates read back changed values while untouched fields stay
  unchanged, and `TestUpdateRejectsInvalidChangesWithoutAppending` plus
  `TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged` prove invalid
  updates reject unchanged for priority, blank title, unknown label, unknown
  owner role, and malformed due date.
- `VAL-001` - PASS. `TestIssueUpdateRequiresKnownByRole` proves `--by` is
  required and validated against known roles at the CLI surface.
- `VAL-002` - PASS.
  `TestUpdateAppendsOneEventProjectsActivityAndRebuilds` proves one successful
  update appends exactly one `issue.updated` event containing the `issues`
  update effect, label delta effects, one `activity` effect, `updated_at`
  advancement, and preconditions on the prior values of changed fields; the
  same test removes the SQLite DB, rebuilds, and proves identical read-back.
- `VAL-003` - PASS.
  `TestUpdateAppendsOneEventProjectsActivityAndRebuilds` proves owner
  role/actor changes read back correctly and project exactly one `activity_log`
  row containing the acting role plus before/after details.
- `VAL-005` - PASS (advanced). `TestUpdateRejectsInvalidChangesWithoutAppending`,
  `TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged`, and
  `TestUpdateRejectsNoOpWithoutAppending` prove rejected updates append zero
  journal lines and leave later reads working.

## Tests run

- `go test ./internal/issues -run TestUpdateChangesTitleAndPriorityAndPreservesUntouchedFields -count=1` - FAIL before implementation (`undefined: Update`, `undefined: UpdateOptions`)
- `go test ./internal/issues -run TestUpdateChangesTitleAndPriorityAndPreservesUntouchedFields -count=1` - PASS
- `go test ./internal/issues -run TestUpdateRejectsInvalidChangesWithoutAppending -count=1` - PASS
- `go test ./internal/issues -run TestUpdateAppendsOneEventProjectsActivityAndRebuilds -count=1` - FAIL before `updated_at`/date normalization fixes; PASS after fixes
- `go test ./internal/issues -run 'TestUpdate' -count=1` - PASS
- `go test ./cmd/fid -run TestIssueUpdateChangesFieldsAndPrintsConfirmation -count=1` - FAIL before CLI wiring (`unknown flag: --title`)
- `go test ./cmd/fid -run TestIssueUpdateChangesFieldsAndPrintsConfirmation -count=1` - PASS
- `go test ./cmd/fid -run 'TestIssueUpdateChangesFieldsAndPrintsConfirmation|TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged|TestIssueUpdateRequiresKnownByRole' -count=1` - PASS
- `go test ./internal/journal ./internal/issues ./cmd/fid -count=1` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS

## Deviations from spec

- Touched `internal/journal/writer.go` and `internal/journal/project.go` in
  addition to the task spec's primary files. `writer.go` gained a small
  pre-append event-preparation helper so the update event can carry the real
  `event_id`/timestamp into the `activity` effect while still dry-running
  before append. `project.go` gained extra value normalization for SQLite date
  and timestamp precondition comparisons so the new update preconditions replay
  deterministically.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `da2c21e`

## Commit subject

- `feat(issue): add audited issue updates`
