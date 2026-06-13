# Phase 5 - Mid-Phase Review 2 (tasks 2-3)

Reviewer: qa-reviewer
Date: 2026-06-13
Tasks covered: phase-5/task-2 and phase-5/task-3

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-5/tasks/task-2.md`
  - `project-docs/docs/phases/phase-5/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-5/tasks/task-2-done.md`
  - `project-docs/docs/phases/phase-5/tasks/task-3-done.md`
- Prior review addendum:
  - `project-docs/docs/phases/phase-5/qa/review-1-addendum.md`
- Validation contract:
  - `project-docs/docs/phases/phase-5/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-5/qa/qa-plan.md`
- Escalation log:
  - `project-docs/docs/current/tl_escalations_current.md`
- Design sections checked:
  - Task 2: `project-docs/specs/05_DESIGN_PLAN.md` §1-2 and `project-docs/specs/06_DESIGN_PLAN_CLI.md` §`fid activity list`; no public activity read surface is in scope.
  - Task 3: `project-docs/specs/05_DESIGN_PLAN.md` §1, §5.3 and `project-docs/specs/06_DESIGN_PLAN_CLI.md` §`fid issue update`.
- Commits or changed files:
  - `3db4c0df2843649e91c52fb25f7fc7b916ce859a` (`feat(journal): add activity projection discipline`)
  - `da2c21e48a02b8391a05e3e25474188da0f7bc9f` (`feat(issue): add audited issue updates`)
  - `internal/journal/project.go`
  - `internal/journal/project_test.go`
  - `internal/issues/update.go`
  - `internal/issues/update_test.go`
  - `cmd/fid/issue.go`
  - `cmd/fid/issue_update_test.go`
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md`

## Summary

- Blocking findings: 0
- Improvement findings: 2
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005

Review 2 is on track. Task 2 establishes the `activity` projection operation and rejects generic `activity_log` writes in the projection engine. Task 3 uses the Task 1 validate-before-append seam, emits one `issue.updated` event with update/label/activity effects, proves rebuild parity, and covers invalid-update journal hygiene. No new open ESC was found. Existing ESC-003 is still present in `tl_escalations_current.md`, but `review-1-addendum.md` records that its closure condition is met.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-3 | ON TRACK | `internal/issues/update_test.go:88` covers invalid priority, blank title, unknown label, unknown owner role, and malformed due date with unchanged read-back and unchanged journal line count. `cmd/fid/issue_update_test.go:9` proves CLI update confirmation and JSON read-back for title/priority. `cmd/fid/issue_update_test.go:94` proves missing and unknown `--by` reject. Reviewer reran `go test ./internal/issues -run 'TestUpdateRejectsInvalidChangesWithoutAppending|TestUpdateAppendsOneEventProjectsActivityAndRebuilds|TestUpdateRejectsNoOpWithoutAppending' -count=1` PASS and `go test ./cmd/fid -run 'TestIssueUpdateChangesFieldsAndPrintsConfirmation|TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged|TestIssueUpdateRequiresKnownByRole' -count=1` PASS. |
| VAL-002 | task-3 | ON TRACK | `internal/issues/update.go:120` prepares a single `issue.updated` event before dry-run; `internal/issues/update.go:142` adds one `issues` update effect, label delta effects, and one `activity` effect. `internal/issues/update_test.go:307` asserts one appended journal line; `internal/issues/update_test.go:311` inspects the stored event shape; `internal/issues/update_test.go:399` removes the DB, rebuilds, and asserts post-rebuild read-back parity. |
| VAL-003 | task-3 | ON TRACK | `internal/issues/update_test.go:294` proves owner role/actor read-back; `internal/issues/update_test.go:370` queries exactly one `activity_log` row for the update event, with actor and before/after owner details. Unknown owner role rejection is covered at `internal/issues/update_test.go:130`. |
| VAL-004 | task-2, task-3 | ON TRACK | `internal/journal/project.go:153` rejects non-`activity` effects targeting `activity_log`; `internal/journal/project.go:177` constrains `activity` effects to `activity_log`. `internal/journal/project_test.go:267` proves one row and JSON details; `internal/journal/project_test.go:376` proves one row per activity effect; `internal/journal/project_test.go:428` proves unresolved refs reject and leave `activity_log` empty; `internal/journal/project_test.go:490` proves generic insert/update/delete effects are rejected; `internal/journal/project_test.go:590` proves rebuild determinism. Reviewer reran the focused activity tests PASS. |
| VAL-005 | task-3, with task-1 addendum context | ON TRACK | Task 3 uses `journal.DryRunProjectEvent` before `journal.AppendEvent` at `internal/issues/update.go:171`. Invalid update line-count and later-read recovery are covered in `internal/issues/update_test.go:182` and `internal/issues/update_test.go:197`; CLI priority rejection also compares before/after JSON and journal line counts at `cmd/fid/issue_update_test.go:61` and `cmd/fid/issue_update_test.go:88`. Task 1 accepted-create and rejected-create evidence was cleared by `review-1-addendum.md`. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 2 completion report path | `project-docs/docs/phases/phase-5/tasks/task-2-done.md` | PASS | Path exists and records summary, files, VAL-004 evidence, exact gates, dependency note, commit hash, and subject. |
| Task 3 completion report path | `project-docs/docs/phases/phase-5/tasks/task-3-done.md` | PASS | Path exists and records summary, files, VAL-001/002/003/005 evidence, exact gates, dependency note, commit hash, and subject. |
| Task 2 commit verification | `git rev-parse --verify 3db4c0d^{commit}` | PASS | Resolved to `3db4c0df2843649e91c52fb25f7fc7b916ce859a`; subject/body include `feat(journal): add activity projection discipline` and `Refs: phase-5/task-2`. |
| Task 3 commit verification | `git rev-parse --verify da2c21e^{commit}` | PASS | Resolved to `da2c21e48a02b8391a05e3e25474188da0f7bc9f`; subject/body include `feat(issue): add audited issue updates` and `Refs: phase-5/task-3`. |
| Required Task 2 gates recorded | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | PASS | Task 2 completion report records all four exact commands as PASS; no conditional command applied. |
| Required Task 3 gates recorded | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | PASS | Task 3 completion report records all four exact commands as PASS. Conditional `go build` was not required because no manual CLI smoke was used as evidence. |
| Reviewer format spot-check | `gofmt -l cmd internal` | PASS | No files listed. |
| Reviewer Task 2 deterministic check | `go test ./internal/journal -run 'TestProjectEventProjectsActivityEffectIntoActivityLog|TestProjectEventRejectsActivityEffectWithUnresolvedRef|TestProjectEventRejectsGenericWritesToActivityLog|TestRebuildWorkspaceReplaysActivityLogDeterministically' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/journal 0.077s`. |
| Reviewer Task 3 domain deterministic check | `go test ./internal/issues -run 'TestUpdateRejectsInvalidChangesWithoutAppending|TestUpdateAppendsOneEventProjectsActivityAndRebuilds|TestUpdateRejectsNoOpWithoutAppending' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/issues 0.172s`. |
| Reviewer Task 3 CLI deterministic check | `go test ./cmd/fid -run 'TestIssueUpdateChangesFieldsAndPrintsConfirmation|TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged|TestIssueUpdateRequiresKnownByRole' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/cmd/fid 0.068s`. |
| Design parity check | `05_DESIGN_PLAN.md` §1/§5.3; `06_DESIGN_PLAN_CLI.md` §`fid issue update` | PASS | CLI adapter delegates to `internal/issues.Update` (`cmd/fid/issue.go:188`), writes confirmation to stderr (`cmd/fid/issue.go:209`), and implemented flags match the v0.1 update field set while deferring `--epic`/`--inbox`. |
| Existing escalation context | `review-1-addendum.md`; `tl_escalations_current.md` | PASS | No new open ESC. ESC-003 remains listed open, but Review 1 addendum states its closure condition is met and no product-code follow-up is needed in this window. |

## Per-Task Findings

### Task 2 - `activity` effect op + activity-log projection discipline (CF-005)

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-5/qa/qa-plan.md:89` VAL-004 is split across Task 2 and the later transition emitter; Task 2 proves projector discipline and Task 3 proves the first emitter, but CF-005 cannot be treated as fully closed until Review 4 verifies Task 5 reuses the `activity` op for transitions. Fix direction: carry the same one-row-per-change, unresolved-ref, and rebuild-parity checks into Review 4 for the transition emitter.

### Task 3 - `fid issue update` with activity-log audit

#### Blocking

- None.

#### Improvement

- `cmd/fid/issue_update_test.go:94` proves missing and unknown `--by` reject, but unlike the invalid-priority CLI path at `cmd/fid/issue_update_test.go:61`, it does not also capture journal line-count/read-back evidence. This is not blocking because package and CLI tests already prove invalid update journal hygiene for representative validation failures, and actor-role validation happens before the write lock/append path. Fix direction: final QA QA-06 should include line-count and later-read evidence for missing or unknown `--by`, or a follow-up CLI test can add the same evidence pattern.

## Cross-cutting observations

- Business behavior is placed in `internal/issues.Update`; `cmd/fid` remains a thin flag/format adapter. Journal projection details stay in `internal/journal`.
- The Task 3 deviation touching `internal/journal/writer.go` is justified by the need to know the event ID before building the `activity` effect, and the implementation still dry-runs before append under the write lock.
- Documentation is current for the reviewed surface: `06_DESIGN_PLAN_CLI.md` records `fid issue update` flags, no-op rejection semantics, and the ASCII `->` confirmation format.
- No dependency changes, security/privacy concern, or data-model change was introduced in the reviewed commits.

## Escalations

- None raised.
- No new open ESC appeared in this review window. Existing ESC-003 remains in `tl_escalations_current.md`; Review 1 addendum records that its closure condition has been satisfied by Task 1b evidence.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-5/qa/review-2.md
Blocking: 0  Improvement: 2  Escalations: none
VERDICT: CONTINUE
```
