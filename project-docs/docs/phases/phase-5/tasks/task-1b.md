# Phase 5 - Task 1b: Close VAL-005 evidence gap (rebuild recovery + accepted one-line append)

## Objective

Close ESC-003 / TRG-001 (review-1 BLOCK on Task 1). The validate-before-append
seam from Task 1 is implemented correctly, but the Task 1 evidence is
incomplete against VAL-005, AC #2/#3 of `task-1.md`, and QA-plan EVID-01:

1. The rejected-create regression
   (`TestCreateRejectsProjectionFailingWriteBeforeAppend`) stops after a later
   `Get` and never rebuilds, so it does not prove the canonical journal still
   rebuilds cleanly after a rejected write.
2. No focused test asserts that an **accepted** create appends **exactly one**
   journal line, and the Task 1 report does not cite such an assertion.

This task adds the missing evidence. It is **test-and-report only**: the seam
implementation (`internal/journal/project.go`, `internal/issues/issue.go`) is
expected to be unchanged. If a new assertion fails, that exposes a real seam
defect to fix — record it as a deviation and fix the seam, do not weaken the
assertion.

## Requirements

- PRD refs: PRD-OPS-001, PRD-OPS-002 (journal append is the commit point)
- TRD refs: TRD §4.3 (validate-before-append write path), TRD-ARCH-001
- Design refs: N/A
- Quality refs: `08_QUALITY_PLAN.md` §4.3 COV-001 (`internal/issues`,
  `internal/journal` ≥ 80%); Phase 5 QA plan EVID-01, QA-03

## Validation Contract Coverage

- VAL-005: A rejected write appends zero journal lines **and all later reads
  succeed including after a full rebuild** (no poisoned tail); an accepted write
  appends **exactly one** line and projects cleanly.

## Scope

### Files To Create Or Modify

- `internal/issues/issue_test.go` - add the two missing VAL-005 assertions:
  - rebuild recovery after a rejected create (delete the SQLite DB and
    `journal.RebuildWorkspace`, then `Get` the previously-valid issue and assert
    success);
  - accepted-create exactly-one-line append (capture journal line count
    before/after an accepted `Create`, assert delta == 1).
- `project-docs/docs/phases/phase-5/tasks/task-1b-done.md` - completion report
  citing the new tests and the existing dry-run/`projection_meta` evidence
  against VAL-005 and EVID-01.

### Files To Read For Context

- `internal/issues/issue_test.go`
  (`TestCreateRejectsProjectionFailingWriteBeforeAppend` at :475;
  `TestCreatePreservesAuthoringFieldsAcrossRebuild` at :540 for the
  delete-DB + `journal.RebuildWorkspace` pattern;
  `TestCreateAndGetRoundTrip` at :627)
- `internal/journal/project_test.go` (existing dry-run rollback tests)
- `project-docs/docs/phases/phase-5/qa/review-1.md` (the BLOCK finding)
- `project-docs/docs/phases/phase-5/qa/qa-plan.md` EVID-01, QA-03

### Out Of Scope

- Any change to the seam implementation in `internal/journal/project.go` or
  `internal/issues/issue.go`, unless a newly added assertion fails (then fix the
  defect and record it as a deviation).
- `fid issue update` / `transition` / `activity` op (Tasks 2, 3, 5).
- Auto-repairing a workspace that already has a poisoned tail (out of scope, per
  task-1.md).

## Acceptance Criteria

- [ ] After a rejected create, deleting the SQLite DB and running
  `journal.RebuildWorkspace` succeeds, and a `Get` of the previously-valid issue
  succeeds post-rebuild (no poisoned tail survives a rebuild). (VAL-005)
- [ ] An accepted create appends **exactly one** journal line, asserted by a
  before/after line-count delta of 1. (VAL-005)
- [ ] The Task 1b completion report cites these two assertions plus the existing
  rejected-create zero-append, later-`Get`, dry-run rollback, and unchanged
  `projection_meta` evidence as the full VAL-005 / EVID-01 picture. (VAL-005)
- [ ] COV-001 still holds for `internal/issues` and `internal/journal` (≥ 80%).

## TDD Plan

- First failing test: extend `TestCreateRejectsProjectionFailingWriteBeforeAppend`
  (or add a sibling, e.g. `TestRejectedCreateRebuildsCleanly`) so that after the
  rejected create it removes `.fid/fid.sqlite`, calls
  `journal.RebuildWorkspace(workspace, ...)`, and `Get`s the valid task. This is
  the gap the BLOCK identified; it should pass if the seam is correct.
- Additional test: accepted-create one-line append — capture
  `countJournalLines(journalPath)` before and after a valid `Create` and assert
  the delta is exactly 1.
- Regression risks: reuse the existing `journal.RebuildWorkspace` /
  `countJournalLines` helpers (see :540 and :505); do not duplicate seam logic;
  keep tests `t.Parallel()`-safe with isolated workspaces.

## Documentation Updates

- `project-docs/docs/phases/phase-5/tasks/task-1b-done.md` - completion report
  with VAL-005 / EVID-01 evidence and the four exact gate results.
- No spec change expected; if the seam needed a fix, record `07_DESIGN_PLAN_FLOWS.md`
  impact (or `verified-no-change`) in the report.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: None (no CLI surface change).
- Audit or release-readiness impact: Completes CF-007 / DEC-014 closure evidence;
  resolves ESC-003. Note the closure in the completion report.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. A focused test command
(e.g. `go test ./internal/issues -run TestRejectedCreateRebuildsCleanly -count=1`)
may be used during the red/green loop but does not substitute for the full gate.
Record each command as `PASS` (or `BLOCKED` with reason) in the completion
report.

Conditional commands:

- None.

## Notes For Developer

- This is evidence remediation: the seam should already satisfy both new
  assertions. Add the tests first, watch them pass against the current seam, and
  only touch `project.go` / `issue.go` if an assertion genuinely fails.
- Mirror the delete-DB + `journal.RebuildWorkspace` pattern in
  `TestCreatePreservesAuthoringFieldsAcrossRebuild` (issue_test.go:587-594).
- Housekeeping: `task-1-done.md` was reported modified relative to commit
  `263d523` (the committed copy lacked the final hash/subject). Commit the
  current `task-1-done.md` state alongside this task so the working tree is clean.
- Commit with subject referencing `phase-5/task-1b` (e.g.
  `test(journal): prove rebuild recovery + one-line append for VAL-005`).
