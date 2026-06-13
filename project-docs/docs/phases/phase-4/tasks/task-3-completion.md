## Summary of what was implemented

- Expanded `fid issue create` to expose the Phase 4 Task 3 CLI flags for hierarchy, authoring fields, labels, ownership, inbox, and epic scheduling while keeping `cmd/fid` limited to flag parsing, normalization, and error rendering.
- Added CLI coverage for epic/task/subtask creation, unattached tasks, invalid parent-flag usage, scheduling/sequence backstop errors, unknown-label rejection, one-event journal behavior, and rebuild parity for a representative hierarchy.

## Files created or modified

- `cmd/fid/issue.go`
- `cmd/fid/issue_test.go`
- `project-docs/docs/phases/phase-4/tasks/task-3-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-4/tasks/task-3-completion.md`

## Validation Contract Coverage

- `VAL-001` PASS: `TestIssueCreateRejectsInvalidFlagsAndConstraintViolations` proves actionable create errors for `issues: task under sequential epic requires sequence`, `issues: task under prioritised epic must not have sequence`, `issues: unattached task must not have sequence`, and `issues: task parent epic must have epic_metadata`; `go test -race -coverprofile=coverage.out ./...` passed with this test included.
- `VAL-002` PASS: `TestIssueCreateEpicAcceptsSchedulingAndReturnsDefaultReadShape` proves epic create/read defaults; `TestIssueCreateSupportsHierarchyAndAuthoringFlags` proves task under epic, unattached task, and subtask create/read defaults and hierarchy via `fid issue get --format json`.
- `VAL-003` PASS: `TestIssueCreateSupportsHierarchyAndAuthoringFlags` proves persisted CLI authoring fields for `priority`, `sequence`, `spec_md_path`, `est_tokens`, labels, `in_inbox`, `owner_role`, and `owner_actor`; the same test verifies `due_date` via direct SQLite query because the current `fid issue get` JSON shape does not expose `due_date`; `TestIssueCreateRejectsInvalidFlagsAndConstraintViolations` proves unknown labels are rejected without auto-create.
- `VAL-004` PASS: `TestIssueCreateEpicAcceptsSchedulingAndReturnsDefaultReadShape` verifies an epic create writes one matching `epic_metadata.scheduling=sequential` row; `TestIssueCreateRejectsInvalidFlagsAndConstraintViolations` proves non-epic parent-flag misuse is rejected at the CLI boundary and missing epic metadata is surfaced as an actionable create failure.
- `VAL-005` PASS: `TestIssueCreateUsesOneJournalEventPerCommandAndRebuildsRepresentativeHierarchy` proves each CLI create appends exactly one journal line and that epic/task/subtask JSON read-back remains identical after deleting `.fid/fid.sqlite` and running `fid rebuild --by developer`.

## Tests run

- `go test ./cmd/fid -run TestIssueCreateEpicAcceptsSchedulingAndReturnsDefaultReadShape -count=1`
- `go test ./cmd/fid -run TestIssueCreateSupportsHierarchyAndAuthoringFlags -count=1`
- `go test ./cmd/fid -run 'TestIssueCreateRejectsInvalidFlagsAndConstraintViolations|TestIssueCreateUsesOneJournalEventPerCommandAndRebuildsRepresentativeHierarchy' -count=1`
- `go test ./cmd/fid`
- `go test -race -coverprofile=coverage.out ./...`

## Exact quality-gate results

- `PASS` `gofmt -l cmd internal` - no output.
- `PASS` `golangci-lint run ./...` - `0 issues.`
- `PASS` `go test -race -coverprofile=coverage.out ./...` - `cmd/fid` `78.4%`; `internal/issues` `86.5%`; `internal/journal` `80.7%`; `internal/policy` `85.0%`.
- `PASS` `go tool cover -func=coverage.out` - `cmd/fid/issue.go newIssueCreateCommand 93.5%`; total statements `77.0%`.
- `PASS` `go build -o /tmp/fid ./cmd/fid` - exited 0.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `2bfe94f`

## Commit subject

- `feat(cli): wire issue create authoring flags`
