# Phase 18 - Mid-Phase Review 1 Remediation (Task 3 / DEC-068)

Reviewer: qa-reviewer
Date: 2026-07-19
Tasks covered: phase-18/task-3 (DEC-068 VAL-020 remediation only)

## Verdict

CONTINUE

## Scope Reviewed

- Decision: `project-docs/docs/phases/phase-18/decisions.md` (DEC-068)
- Task specification and remediation note:
  `project-docs/docs/phases/phase-18/tasks/task-3.md`
- Corrected completion report:
  `project-docs/docs/phases/phase-18/tasks/task-3-done.md`
- Prior review: `project-docs/docs/phases/phase-18/qa/review-1.md`
- Validation contract and refreshed QA plan:
  `project-docs/docs/phases/phase-18/qa/validation-contract.md`,
  `project-docs/docs/phases/phase-18/qa/qa-plan.md`
- Remediation implementation commit:
  `77400cebaa57aa3518bfad0ea7938a4fa6325cfd`
  (`test(issues): prove released claims stay absent after rebuild`)
- Report-finalization commit:
  `3e0e101a22b0ba31e8d5ac29d8a20e48ae83f26f`
  (`docs(task): finalize completion report for phase-18 task-3`)

## Summary

- Blocking findings: 0
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-020

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-020 | task-3 | PASS — remediation verified | [`internal/issues/transition_test.go:650`](../../../../../internal/issues/transition_test.go#L650) creates a fresh workspace, establishes the live `human` caller claim, enters `in_progress`, verifies the caller claim, transitions the same fixture to `done`, verifies immediate absence, rebuilds that exact workspace with `journal.RebuildWorkspace`, and verifies absence again. Independent non-cached race execution passed in 1.753s. The paired same-event release-effect test independently passed in 1.504s. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Exact released-claim rebuild characterization | `go test -race -count=1 ./internal/issues -run '^TestTransitionOutOfStartedReleasedCallerClaimRemainsAbsentAfterRebuild$'` | PASS | Independent result: `ok github.com/mbell-nzl/fid/internal/issues 1.753s`. This is the exact DEC-068 fixture sequence, not the earlier entering-started rebuild fixture. |
| Same-event release effect | `go test -race -count=1 ./internal/issues -run '^TestTransitionOutOfStartedReleasesCallerClaimInSameEvent$'` | PASS | Independent result: `ok github.com/mbell-nzl/fid/internal/issues 1.504s`; paired with the new test, this proves the event effect plus post-projection and rebuilt-projection absence required by QA-15/EVID-09. |
| Representative required task command | `gofmt -l cmd internal` | PASS | Exit 0; no files listed. This is an exact Task 3 and QA-plan required command. |
| Full Task 3 exact-command records | `task-3-done.md` compared verbatim with `task-3.md` and QA plan §5 | PASS | All five required commands are recorded as `PASS`: format, `golangci-lint run ./...`, `go test -race -coverprofile=coverage.out ./...`, `go tool cover -func=coverage.out`, and `go test -tags integration -race ./...`. The Task 3 QA mapping is QA-13..15; it adds evidence obligations, not another command. |
| Completion-report mechanical evidence | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-18/tasks/task-3-done.md project-docs/docs/phases/phase-18/tasks/task-3.md` | PASS | All headings, five verbatim required-command PASS records, reported hash, subject, and committed/clean report verified. |
| Remediation commit metadata and scope | `git rev-parse --verify 77400cebaa57aa3518bfad0ea7938a4fa6325cfd^{commit}`; `git show --check`; `git diff-tree --name-status -r` | PASS | Hash resolves; subject and `Refs: phase-18/task-3` footer match the report. Diff is limited to `internal/issues/transition_test.go` and the corrected Task 3 report; whitespace check passes. |
| No scope or contract change | Remediation commit diff against DEC-068, Task 3, QA plan, and VAL-020 | PASS | No changes to the task spec, QA plan, validation contract, decision record, run ledger, product implementation, journal shape, or migration. The report now cites the exact test and removes the previously inadequate general-rebuild claim. |

## Per-Task Findings

### Task 3 - Transition-path data-source swap (DEC-068 VAL-020 remediation)

#### Blocking

- None.

#### Improvement

- None.

#### Observation

- The remediation supplies both missing rebuilt-projection evidence and the already-required immediate-release link without weakening VAL-020: the live caller claim is present in the started state, absent immediately after leaving it, and remains absent after replaying the same workspace.

## Escalations

- None. DEC-068 correctly classifies the defect and its repair as localized to Task 3; the re-review found no new cross-task, contract, architecture, or quality-gate anomaly.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-18/qa/review-1-remediation.md
Blocking: 0  Improvement: 0  Escalations: none
VAL-020: PASS
VERDICT: CONTINUE
```
