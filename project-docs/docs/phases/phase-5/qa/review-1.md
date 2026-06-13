# Phase 5 - Mid-Phase Review 1 (tasks 1-1)

Reviewer: qa-reviewer
Date: 2026-06-13
Tasks covered: phase-5/task-1

## Verdict

BLOCK

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-5/tasks/task-1.md`
- Completion reports:
  - `project-docs/docs/phases/phase-5/tasks/task-1-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-5/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-5/qa/qa-plan.md`
- Commits or changed files:
  - `263d52339f39a3426e429e64c3bf991e976233e7` (`fix(journal): validate issue creates before append`)
  - `internal/issues/issue.go`
  - `internal/issues/issue_test.go`
  - `internal/journal/project.go`
  - `internal/journal/project_test.go`

## Summary

- Blocking findings: 1
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-005

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-005 | task-1 | NOT ON TRACK | `internal/issues/issue.go:141` refreshes projection before the locked write path and `internal/issues/issue.go:265` calls `journal.DryRunProjectEvent` before `AppendEvent`, which is the right implementation shape. `TestCreateRejectsProjectionFailingWriteBeforeAppend` proves the Phase 4 invalid-create repro returns the trigger error, leaves the journal line count unchanged, and allows a later `Get` (`internal/issues/issue_test.go:475`). `TestDryRunProjectEventRollsBackEffectsAndLeavesProjectionMetaUnchanged` and `TestDryRunProjectEventRollsBackOnProjectionFailure` prove dry-run rollback and unchanged `projection_meta` (`internal/journal/project_test.go:418`, `internal/journal/project_test.go:508`). However, the Task 1 evidence does not prove rebuild recovery after the rejected create, which is explicitly required by the task spec and QA plan. Accepted-create one-line append is present in older coverage (`internal/issues/issue_test.go:627`), but the completion report does not cite that test as Task 1 evidence. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Completion report path | `project-docs/docs/phases/phase-5/tasks/task-1-done.md` | PASS | Path exists and current working-tree content records commit `263d523` and subject `fix(journal): validate issue creates before append` (`task-1-done.md:52`). |
| Commit verification | `git rev-parse --verify 263d523^{commit}` | PASS | Resolved to `263d52339f39a3426e429e64c3bf991e976233e7`; `git show -s --format=fuller 263d523` confirms subject and `Refs: phase-5/task-1`. |
| Format gate recorded | `gofmt -l cmd internal` | PASS | Completion report records PASS (`task-1-done.md:35`). Reviewer spot-check also produced no listed files. |
| Lint gate recorded | `golangci-lint run ./...` | PASS | Completion report records PASS (`task-1-done.md:36`). Not rerun in this review. |
| Race coverage gate recorded | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records PASS (`task-1-done.md:37`). Not rerun in this review. |
| Coverage report recorded | `go tool cover -func=coverage.out` | PASS | Completion report records PASS (`task-1-done.md:38`). Reviewer spot-check of existing `coverage.out` showed `internal/issues/issue.go` at 85.6% and `internal/journal/project.go` at 80.9%, matching the COV-001 claim. |
| Focused rejected-create regression | `go test ./internal/issues -run TestCreateRejectsProjectionFailingWriteBeforeAppend -count=1` | PASS | Reviewer rerun passed: `ok github.com/mbell-nzl/fid/internal/issues 0.032s`. Covers zero append and later `Get`, but not rebuild. |
| Focused dry-run seam regression | `go test ./internal/journal -run TestDryRunProjectEvent -count=1` | PASS | Reviewer rerun passed: `ok github.com/mbell-nzl/fid/internal/journal 0.021s`. Covers success rollback, failure rollback, and unchanged `projection_meta`. |
| Focused accepted/rebuild existing coverage | `go test ./internal/issues -run 'TestCreateRejectsProjectionFailingWriteBeforeAppend|TestCreateAndGetRoundTrip|TestCreatePreservesAuthoringFieldsAcrossRebuild' -count=1` | PASS | Reviewer rerun passed: `ok github.com/mbell-nzl/fid/internal/issues 0.060s`. Existing tests cover accepted create and generic rebuild, but not rebuild after the rejected-create repro. |

## Per-Task Findings

### Task 1 - Validate-before-append write-path seam (CF-007)

#### Blocking

- [project-docs/docs/phases/phase-5/qa/qa-plan.md:114] EVID-01 requires Task 1 evidence for rejected-create line count unchanged, later `issue get`/rebuild succeeds, accepted create appends one line, and dry-run rollback. The completion report cites the new rejected-create test only for line count and later `Get` (`project-docs/docs/phases/phase-5/tasks/task-1-done.md:24`), while the test itself stops after `Get` and never runs `journal.RebuildWorkspace` in the rejected-write fixture (`internal/issues/issue_test.go:527`). The accepted one-line append requirement is only indirectly covered by existing create round-trip tests and is not cited in the Task 1 report as the QA plan requires. Fix direction: add or extend Task 1 evidence so the rejected-create fixture also proves rebuild succeeds after the failed write, and explicitly cite an accepted-create one-line append assertion in the completion evidence.

#### Improvement

- None.

## Cross-cutting observations

- The implementation placement matches the intended module boundary: `internal/issues.Create` keeps write-path orchestration in the domain service, while `internal/journal.DryRunProjectEvent` owns projection simulation.
- The commit is valid and correctly references `phase-5/task-1`, but `git status --short` showed `project-docs/docs/phases/phase-5/tasks/task-1-done.md` modified relative to `263d523`; the current working-tree report has the hash/subject that the committed copy lacked.
- `project-docs/docs/phases/phase-5/qa/qa-plan.md` is currently untracked in this workspace. This review treated it as authoritative because it was in the invoking scope.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-5/qa/review-1.md
Blocking: 1  Improvement: 0  Escalations: none
VERDICT: BLOCK
```
