# Phase 5 - Mid-Phase Review 1 Addendum (tasks 1-1b)

Reviewer: qa-reviewer
Date: 2026-06-13
Tasks covered: phase-5/task-1 and phase-5/task-1b

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-5/tasks/task-1.md`
  - `project-docs/docs/phases/phase-5/tasks/task-1b.md`
- Completion reports:
  - `project-docs/docs/phases/phase-5/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-5/tasks/task-1b-done.md`
- Prior review:
  - `project-docs/docs/phases/phase-5/qa/review-1.md`
- Validation contract:
  - `project-docs/docs/phases/phase-5/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-5/qa/qa-plan.md`
- Escalation and triage context:
  - `project-docs/docs/current/tl_escalations_current.md`
  - `project-docs/docs/phases/phase-5/triage-log.md`
- Commits or changed files:
  - `263d52339f39a3426e429e64c3bf991e976233e7` (`fix(journal): validate issue creates before append`)
  - `c082152a220494478935d89113560207f15b0e52` (`test(journal): close VAL-005 evidence gap`)
  - `internal/issues/issue_test.go`
  - `internal/journal/project_test.go`

## Summary

- Blocking findings: 0
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-005

The prior Review 1 BLOCK is cleared. Task 1b supplies the missing VAL-005
evidence called out in `review-1.md` and ESC-003: the rejected-create fixture now
proves rebuild recovery after the failed write, and the accepted-create path now
asserts exactly one journal line is appended.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-005 | task-1, task-1b | ON TRACK | `TestCreateRejectsProjectionFailingWriteBeforeAppend` now proves the Phase 4 invalid-create repro returns the expected projection rejection, leaves the journal line count unchanged, allows a later `Get`, deletes `.fid/fid.sqlite`, runs `journal.RebuildWorkspace`, and successfully `Get`s the previously valid task after rebuild (`internal/issues/issue_test.go:475`, rebuild at `internal/issues/issue_test.go:539`). `TestCreateAndGetRoundTrip` now captures journal line count before an accepted create and asserts the post-create delta is exactly `1` (`internal/issues/issue_test.go:648`, assertion at `internal/issues/issue_test.go:718`). Existing dry-run seam tests still prove rollback and unchanged `projection_meta` for success and failure (`internal/journal/project_test.go:418`, `internal/journal/project_test.go:508`). Reviewer reran the focused Task 1b tests: `go test ./internal/issues -run 'TestCreateRejectsProjectionFailingWriteBeforeAppend|TestCreateAndGetRoundTrip' -count=1` PASS; `go test ./internal/journal -run TestDryRunProjectEvent -count=1` PASS. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1b completion report path | `project-docs/docs/phases/phase-5/tasks/task-1b-done.md` | PASS | Path exists and records the Task 1b summary, VAL-005 evidence, exact gates, dependency note, commit hash, and commit subject. |
| Task 1b commit verification | `git rev-parse --verify c082152^{commit}` | PASS | Resolved to `c082152a220494478935d89113560207f15b0e52`. |
| Task 1b commit subject/body | `git show -s --format=fuller c082152` | PASS | Subject is `test(journal): close VAL-005 evidence gap`; body includes `Refs: phase-5/task-1b`. |
| Task 1b changed files | `git show --name-only --format=fuller c082152` | PASS | Commit changes only `internal/issues/issue_test.go` and `project-docs/docs/phases/phase-5/tasks/task-1b-done.md`, matching the test-and-report scope. |
| Format gate recorded | `gofmt -l cmd internal` | PASS | Task 1b completion report records PASS; reviewer spot-check produced no listed files. |
| Lint gate recorded | `golangci-lint run ./...` | PASS | Task 1b completion report records PASS. Not rerun in this addendum. |
| Race coverage gate recorded | `go test -race -coverprofile=coverage.out ./...` | PASS | Task 1b completion report records PASS. Not rerun in this addendum. |
| Coverage report recorded | `go tool cover -func=coverage.out` | PASS | Task 1b completion report records PASS. |
| Rejected-create rebuild recovery | `go test ./internal/issues -run 'TestCreateRejectsProjectionFailingWriteBeforeAppend|TestCreateAndGetRoundTrip' -count=1` | PASS | Reviewer rerun passed: `ok github.com/mbell-nzl/fid/internal/issues 0.124s`. Covers rejected-write zero append, later read, rebuild recovery, and accepted-create one-line append. |
| Dry-run rollback evidence | `go test ./internal/journal -run TestDryRunProjectEvent -count=1` | PASS | Reviewer rerun passed: `ok github.com/mbell-nzl/fid/internal/journal 0.048s`. Covers dry-run rollback and unchanged `projection_meta`. |

## Per-Task Findings

### Task 1 - Validate-before-append write-path seam (CF-007)

#### Blocking

- None. The prior blocking finding from `review-1.md` is cleared by Task 1b evidence.

#### Improvement

- None.

### Task 1b - Close VAL-005 evidence gap (rebuild recovery + accepted one-line append)

#### Blocking

- None.

#### Improvement

- None.

## Cross-cutting observations

- ESC-003 remains listed as open in `project-docs/docs/current/tl_escalations_current.md`, but its stated closure condition is now met by Task 1b plus this addendum. No new open ESC was discovered, and no new escalation was opened.
- Task 1b stayed within the requested test-and-report scope: product code was not changed in commit `c082152`.

## Escalations

- None raised.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-5/qa/review-1-addendum.md
Blocking: 0  Improvement: 0  Escalations: none
VERDICT: CONTINUE
```
