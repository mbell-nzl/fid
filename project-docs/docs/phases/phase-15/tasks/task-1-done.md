## Summary of what was implemented

Implemented the Phase 15 Task 1 behavior-preserving cleanup for CF-014 and CF-015 only.

- Split `cmd/fid/issue.go` into smaller sibling files in `cmd/fid` without changing command wiring, flags, or output behavior.
- Split `internal/journal/replay.go` into `replay.go`, `projection.go`, and `inspect.go` without changing replay, rebuild, or freshness behavior.
- Extracted private helpers in `internal/labels/add.go` and `internal/issues/label.go` so the targeted non-test functions are each <=60 lines while keeping journal event types, effect shapes, and write flow behavior unchanged.

## Files created or modified

- `cmd/fid/issue.go`
- `cmd/fid/issue_ac.go`
- `cmd/fid/issue_claim.go`
- `cmd/fid/issue_complete.go`
- `cmd/fid/issue_create.go`
- `cmd/fid/issue_get_list.go`
- `cmd/fid/issue_helpers.go`
- `cmd/fid/issue_transition.go`
- `cmd/fid/issue_update.go`
- `internal/journal/replay.go`
- `internal/journal/projection.go`
- `internal/journal/inspect.go`
- `internal/labels/add.go`
- `internal/issues/label.go`
- `project-docs/docs/phases/phase-15/tasks/task-1-done.md`

## Documentation updated

- Added `project-docs/docs/phases/phase-15/tasks/task-1-done.md` with Task 1 implementation and gate evidence.

## Validation Contract Coverage

- `VAL-001` PASS.
  `wc -l` proof for split `cmd/fid` files:
  `22 cmd/fid/issue.go`
  `215 cmd/fid/issue_ac.go`
  `149 cmd/fid/issue_claim.go`
  `90 cmd/fid/issue_create.go`
  `87 cmd/fid/issue_complete.go`
  `144 cmd/fid/issue_get_list.go`
  `228 cmd/fid/issue_helpers.go`
  `48 cmd/fid/issue_transition.go`
  `69 cmd/fid/issue_update.go`
  `wc -l` proof for split `internal/journal` files:
  `205 internal/journal/replay.go`
  `192 internal/journal/projection.go`
  `128 internal/journal/inspect.go`
  All resulting files are <=400 lines.
- `VAL-001` behavior-preservation evidence.
  Focused package safety net before refactor: `go test ./cmd/fid ./internal/journal ./internal/labels ./internal/issues` PASS.
  Focused package safety net after refactor: `go test ./cmd/fid ./internal/journal ./internal/labels ./internal/issues` PASS.
  Full task gate PASS with no test changes.
- `VAL-002` PASS.
  Function line-span evidence from `nl -ba` inspection:
  `internal/labels/add.go`: `Add` lines 33-50, `prepareLabelAdd` lines 52-71, `addLabelLocked` lines 73-98, `prepareLabelAddEvent` lines 104-125, `appendAndProjectLabelEvent` lines 138-148.
  `internal/issues/label.go`: `AddLabel` lines 42-69, `RemoveLabel` lines 71-98, `prepareIssueLabelWrite` lines 123-157, `issueLabelWrite` lines 159-171, `issueLabelWriteLocked` lines 173-198, `prepareIssueLabelEvent` lines 204-231, `issueLabelEffects` lines 233-258, `appendAndProjectIssueLabelEvent` lines 260-270.
  Each targeted non-test function is <=60 lines.
- Coverage percentages from `go test -race -coverprofile=coverage.out ./...`:
  `cmd/fid` 77.4%
  `internal/issues` 81.7%
  `internal/journal` 80.5%
  `internal/labels` 82.4%

## Tests run

- `go test ./cmd/fid ./internal/journal ./internal/labels ./internal/issues` before edits
- `go test ./cmd/fid ./internal/journal ./internal/labels ./internal/issues` after edits
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`
- `go test -tags integration -race ./...`
- `go build -o /tmp/fid ./cmd/fid`

## Exact quality-gate results

- PASS `gofmt -l cmd internal`
  Output: none.
- PASS `golangci-lint run ./...`
  Output excerpt: `0 issues.`
- PASS `go test -race -coverprofile=coverage.out ./...`
  Output excerpt: `ok github.com/mbell-nzl/fid/cmd/fid ... coverage: 77.4%`, `ok github.com/mbell-nzl/fid/internal/issues ... coverage: 81.7%`, `ok github.com/mbell-nzl/fid/internal/journal ... coverage: 80.5%`, `ok github.com/mbell-nzl/fid/internal/labels ... coverage: 82.4%`.
- PASS `go tool cover -func=coverage.out`
  Output excerpt: `total: (statements) 75.8%`.
- PASS `go test -tags integration -race ./...`
  Output excerpt: `ok github.com/mbell-nzl/fid/cmd/fid 60.024s`, `ok github.com/mbell-nzl/fid/internal/issues 21.315s`, `ok github.com/mbell-nzl/fid/internal/journal 21.089s`, `ok github.com/mbell-nzl/fid/internal/labels 2.055s`.
- PASS `go build -o /tmp/fid ./cmd/fid`
  Output: none.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0906428`
- Full hash: `0906428adf9bb4d89b41dc94ad5b147b1b7c538f`

## Commit subject

- `refactor: split issue and replay cleanup paths`
