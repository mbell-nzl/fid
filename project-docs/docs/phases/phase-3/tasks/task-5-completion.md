# Phase 3 - Task 5: Conflict / Corruption Error Fidelity (Rebuild + Read Catch-Up) - Completion Report

## Summary of what was implemented

- Added shared `cmd/fid` journal replay error rendering so the CLI formats typed corruption and precondition-conflict failures consistently without moving detection logic out of `internal/journal`.
- Extended `fid rebuild` to render `PreconditionConflictError` structurally with journal file/line, event id/type, entity primary key, expected values, and actual values.
- Extended `fid issue get` to stop automatic projection catch-up on malformed trailing lines and replay precondition conflicts, emit structured stderr, and avoid serving stdout on failure.
- Added `cmd/fid` tests covering rebuild conflict rendering and read-path corruption/conflict failures, with clean read success asserted before injecting each bad tail/event.

## Files created or modified

- `cmd/fid/journal_errors.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/issue.go`
- `cmd/fid/rebuild_test.go`
- `cmd/fid/issue_test.go`
- `project-docs/docs/phases/phase-3/tasks/task-5-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-3/tasks/task-5-completion.md`

## Validation Contract Coverage

- `VAL-008` - PASS. `cmd/fid/issue_test.go` adds `TestIssueGetStopsOnJournalCatchUpErrors/malformed_trailing_line`, which proves a clean `fid issue get` succeeds before catch-up, then a malformed trailing line appended after the checkpoint causes non-zero exit, empty stdout, and structured stderr with journal file/line, parse cause, and manual-repair guidance. `cmd/fid/rebuild_test.go` keeps corruption coverage on `fid rebuild` via `TestRebuildCommandReportsCorruptJournal`, while `TestRebuildCommandReplaysJournalAndReportsCounts` covers the clean rebuild state.
- `VAL-009` - PASS. `cmd/fid/rebuild_test.go` adds `TestRebuildCommandReportsPreconditionConflict`, proving `fid rebuild` renders the replay conflict structurally with event id/type, entity pk, expected value, and actual value. `cmd/fid/issue_test.go` adds `TestIssueGetStopsOnJournalCatchUpErrors/precondition_conflict`, proving `fid issue get` succeeds on the clean checkpointed state, then rejects a conflicting catch-up event before serving stdout and emits the same structured conflict fields.

## Tests run

- `go test -race ./cmd/fid -run TestRebuildCommandReportsPreconditionConflict` - FAIL (expected red: rebuild fell through to the generic `PreconditionConflictError.Error()` string instead of the required structured CLI lines and `.fid/journal/...` path)
- `go test -race ./cmd/fid -run TestRebuildCommandReportsPreconditionConflict` - PASS
- `go test -race ./cmd/fid -run TestIssueGetStopsOnJournalCatchUpErrors` - PASS
- `go test -race ./cmd/fid -run 'Test(RebuildCommand|IssueGetStopsOnJournalCatchUpErrors)'` - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS (`internal/journal` 80.6%, `cmd/fid` 74.3%, total 75.2%)
- `go build -o /tmp/fid ./cmd/fid` - PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `f1b18fb`

## Commit subject

- `fix(fid): render replay errors structurally`
