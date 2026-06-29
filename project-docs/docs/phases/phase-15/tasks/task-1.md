# Phase 15 - Task 1: Behavior-preserving cleanup (CF-014 + CF-015)

## Objective

Discharge two scheduled carry-forward cleanups before adding new composite code,
so later tasks build on smaller, clearer files. **No behavior change, no
exported-surface change, no CLI/JSON/journal contract change.**

- **CF-014:** split `cmd/fid/issue.go` (984 lines) and
  `internal/journal/replay.go` (485 lines) under the 400-line non-test file
  target.
- **CF-015:** extract shared helpers in `internal/labels/add.go` (`Add`, ~98
  lines) and `internal/issues/label.go` (`AddLabel`/`RemoveLabel`, long and
  partly duplicated) under the 60-line non-test function target.

CF-013 (`internal/issues/complete.go`) is **not** in this task; it ships in
Task 3 alongside `fid dev finish`/`block`, where the extracted completion-effect
builder is actually consumed.

## Requirements

- PRD refs: N/A (refactor only)
- TRD refs: `02_TRD.md` §3 (module boundaries — unchanged)
- Design refs: N/A
- Quality refs: `08_QUALITY_PLAN.md` size budgets (≤400-line files, ≤60-line
  non-test functions); AGENTS.md "Keep changes small"

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-15/qa/validation-contract.md`.

- VAL-001: `cmd/fid/issue.go` and `internal/journal/replay.go` each ≤400 lines
  after the split, with no exported-surface or behavior change.
- VAL-002: `internal/labels/add.go` and the `internal/issues/label.go`
  add/remove write flows refactored so each non-test function is ≤60 lines, with
  label add/remove behavior unchanged.

## Scope

### Files To Create Or Modify

CF-014 (`cmd/fid`, `package main` — pure file moves between files in one
package):

- `cmd/fid/issue.go` - keep `newIssueCommand()` wiring; move subcommand
  constructors and shared helpers out into sibling files. Target: each resulting
  file ≤400 lines and `issue.go` itself small.
- New `cmd/fid/issue_*.go` files grouping the existing constructors, e.g.
  `issue_ac.go` (`newIssueACCommand`/`add`/`verify`/`list`),
  `issue_create.go`, `issue_complete.go`, `issue_transition.go`,
  `issue_update.go`, `issue_get_list.go`, and `issue_helpers.go`
  (`renderIssueError`, `writeIssueLine`, `issueMigrationsDir`,
  `joinCommaSeparated`, the `issueCreate*`/`issueClaim*`/`issueUpdate*` helpers).
  Claim/unclaim/label constructors may stay or move; `issue_label.go` already
  exists. Choose a grouping that keeps each file ≤400 lines without splitting a
  single constructor across files.

CF-014 (`internal/journal`, one package — file moves):

- `internal/journal/replay.go` - keep the core replay path
  (`ReplayJournal`, `replayJournalFile`, `orderedJournalFiles`,
  `CorruptJournalError`, `ReplayStats`/`ReplayOptions`).
- New `internal/journal/projection.go` - move the freshness/rebuild/catch-up
  functions (`EnsureProjectionFresh*`, `ensureProjectionFresh`,
  `RebuildWorkspace*`, `catchUp*`, `EnsureProjectionResult`, `journalCheckpoint`).
- New `internal/journal/inspect.go` - move `inspectJournal`, `journalInspection`,
  and their helpers.
  (Pick the split that lands each file ≤400 lines; the three-way cut above is a
  recommendation, not a mandate.)

CF-015 (function extraction; behavior identical):

- `internal/labels/add.go` - extract the held-lock build/append/project body of
  `Add` into small helpers (validate/options, event build, append+project) so
  each non-test function is ≤60 lines.
- `internal/issues/label.go` - extract the duplicated
  validate→prepare→effects→DryRun→Append→Project flow shared by `AddLabel` and
  `RemoveLabel` into a small internal helper (e.g. an
  `issueLabelWrite(...)` seam parameterised by op/event-type/effects), keeping
  each non-test function ≤60 lines. Do not change effect shapes or event types.

### Files To Read For Context

- `cmd/fid/issue.go`, `cmd/fid/root.go` (command wiring)
- `internal/journal/replay.go`
- `internal/labels/add.go`, `internal/issues/label.go`

### Out Of Scope

- CF-013 (`internal/issues/complete.go`) — Task 3.
- Any new command, flag, event type, effect, or behavior change.
- Renaming exported symbols or changing CLI/JSON/journal output.

## Acceptance Criteria

- [ ] `cmd/fid/issue.go` and `internal/journal/replay.go` are each ≤400 lines;
  no new resulting file exceeds 400 lines. (VAL-001)
- [ ] All pre-existing `cmd/fid` and `internal/journal` tests pass unmodified
  except for mechanical moves; no exported symbol, CLI flag, JSON shape, or
  journal event type changes. (VAL-001)
- [ ] Every non-test function in `internal/labels/add.go` and the
  `internal/issues/label.go` `AddLabel`/`RemoveLabel` flows is ≤60 lines. (VAL-002)
- [ ] Label add/remove behavior is unchanged; pre-existing `internal/labels` and
  issue-label tests pass unmodified except for moves. (VAL-002)

## TDD Plan

- This is a behavior-preserving refactor: the existing test suites are the
  red/green safety net. Run the full gate before and after; no test should need
  a behavior change. If a test must change beyond a mechanical move, stop — that
  signals a behavior change and breaks VAL-001/002.
- Regression risks: accidental exported-surface change during file moves;
  splitting a function across files; altering an effect shape in the label
  helper extraction. Keep effect maps and event types byte-identical.

## Documentation Updates

- None required (no contract change). If file moves are large enough to warrant
  it, note the new file layout in the completion report only.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: None (no surface changes).
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command in
the completion report as `PASS`, or `BLOCKED` with reason. Record the
`internal/labels`, `internal/issues`, `internal/journal`, and `cmd/fid` coverage
percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` — required: the package set under
  refactor (`cmd/fid`) contains integration-tagged tests; record PASS/BLOCKED.
- `go build -o /tmp/fid ./cmd/fid` — required: prove the CLI still builds after
  the `cmd/fid` file split.

## Notes For Developer

- This is the lowest-risk task in the phase; do it first to shrink `issue.go`
  before any composite registration lands.
- Behavior preservation is the whole point: prefer pure cut/paste of function
  bodies into new files in the same package over any rewrite.
- For CF-015, the goal is the 60-line target via extraction, not a new
  abstraction layer — keep helpers private to the package and effect shapes
  identical.
- Confirm `gofmt`/`goimports` on every moved file; file splits commonly leave
  unused imports.
