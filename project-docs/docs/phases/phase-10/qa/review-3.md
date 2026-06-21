# Phase 10 - Mid-Phase Review 3 (tasks 3-3)

Reviewer: qa-reviewer
Date: 2026-06-21
Tasks covered: phase-10/task-3

## Verdict

BLOCK

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-10/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-10/tasks/task-3-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-10/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-10/qa/qa-plan.md`
- Decisions and brief:
  - `project-docs/docs/phases/phase-10/decisions.md`
  - `project-docs/docs/current/phase_brief.md`
- Design sections checked for polish/parity:
  - `project-docs/specs/05_DESIGN_PLAN.md` §1, especially "Errors that teach"
  - `project-docs/specs/05_DESIGN_PLAN.md` §3.1 CLI command tree for in-scope issue/AC commands
  - `project-docs/specs/05_DESIGN_PLAN.md` §4 primitive command behavior
  - `project-docs/specs/05_DESIGN_PLAN.md` scenario row "Agent with wrong role attempts forbidden action"
- Commits or changed files:
  - `5293fd2674b0afdc902feb2f7079d12b45e82551` (`feat(policy): enforce write permissions on commands`)
  - `git status --short` shows `project-docs/docs/phases/phase-10/tasks/task-3-done.md` modified after the commit.

## Summary

- Blocking findings: 1
- Improvement findings: 0
- Escalations raised: ESC-007
- Covered validation assertions: VAL-004, VAL-005, VAL-006, VAL-007, VAL-008

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-004 | task-3 | ON TRACK, blocked only by report commit-state evidence | Task report lines 29 and 49-55 claim denied epic create, permitted counterpart, and full gates. Source spot-check shows `issue create` calls `policy.Authorize` before write setup at `internal/issues/create.go:33`; focused `go test` below passed. |
| VAL-005 | task-3 | ON TRACK, blocked only by report commit-state evidence | Task report line 30 cites missing/unknown-role tests. `rg "RequireKnownRole\|Authorize\\(" internal/issues cmd/fid -n` found `Authorize` on scoped writes and only owner-role `RequireKnownRole` exceptions in `internal/issues/update_request.go:69` and `internal/issues/create.go:291`, matching the task spec exception. |
| VAL-006 | task-3 | ON TRACK, blocked only by report commit-state evidence | `cmd/fid/issue_ac_test.go:533` verifies `--by human`; `cmd/fid/issue_ac_test.go:547` reads back latest verification under `human`, while workflow_runner denial is checked at `cmd/fid/issue_ac_test.go:494`. |
| VAL-007 | task-3 | ON TRACK, blocked only by report commit-state evidence | Action derivation helpers live in `internal/issues/policy_actions.go:13`; AC verifier denial/permit evidence is in `cmd/fid/issue_ac_test.go:494` and `cmd/fid/issue_ac_test.go:519`; transition denial/permit evidence is in `cmd/fid/issue_transition_test.go:215` and `cmd/fid/issue_transition_test.go:234`. |
| VAL-008 | task-3 | ON TRACK, blocked only by report commit-state evidence | Show-vs-enforce parity test checks displayed denied and permitted cells at `cmd/fid/policy_test.go:136`, then denied and permitted writes at `cmd/fid/policy_test.go:149` and `cmd/fid/policy_test.go:168`. Focused test command below passed. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Commit hash verification | `git rev-parse --verify 5293fd2674b0afdc902feb2f7079d12b45e82551^{commit}` | PASS | Returned the same hash. |
| Committed report state | `git show 5293fd2674b0afdc902feb2f7079d12b45e82551:project-docs/docs/phases/phase-10/tasks/task-3-done.md \| tail -n 12` | FAIL | The committed report still says `Pending commit; update after ...`; current working tree has the filled hash/subject. Blocking finding below. |
| Exact task commands in report | Task report lines 49-55 | PASS | All exact required commands are recorded as PASS: `gofmt`, `golangci-lint`, race coverage, `go tool cover`, integration, build. Conditional benchmark command is also recorded PASS with a reason. |
| Coverage thresholds | Task report line 52 | PASS | Report records COV-001 package thresholds above 80% and COV-002 `cmd/fid` at 77.8%. |
| Representative deterministic check | `go test ./cmd/fid ./internal/issues ./internal/policy -run 'TestIssueCreateAndRebuildValidateRoleExistence\|TestIssueACVerifyCommandUsesRoleSpecificPermissionsAndHumanAttribution\|TestIssueTransitionUsesStateSpecificPermissions\|TestPolicyShowMatchesWriteEnforcementForPermittedAndDeniedCells\|TestMatrixContainsExactlyTheImplementedV01WriteActions\|TestAuthorizePermitsAllowedCellAndDeniesKnownRoleWithoutPermission'` | PASS | Passed for `cmd/fid`, `internal/issues` (no tests to run for selected names), and `internal/policy`. |
| Format spot-check | `gofmt -l cmd internal` | PASS | No files listed. |
| Write-surface source check | `rg "RequireKnownRole\|Authorize\\(" internal/issues cmd/fid -n` | PASS | Scoped writes use `Authorize`; remaining `RequireKnownRole` instances are owner-role validation exceptions, as allowed by task spec lines 48-50. |
| Fixture-role-change note | `project-docs/docs/phases/phase-10/decisions.md:6` | PASS | DEC-E execution note lists rebuild, raw journal fixture, transition, and setup role moves; task report lines 22-23 match the docs/decision note. |
| Documentation notes | Task report lines 22-25 and 57-59 | PASS | Decisions and `03_TRD_SCHEMA.md` updates are recorded; `06_DESIGN_PLAN_CLI.md` and `01_PRD.md` are recorded as verified-no-change. |
| Concerns/dependency-budget fields | Task report lines 61-67 | PASS | Concerns are recorded as none; dependency-budget note says no dependency changes. |
| Commit fields | Task report lines 69-75 | FAIL | Fields are present only in the working tree; the task commit still contains pending placeholders. |
| QA-plan evidence requirements | QA plan lines 96-100 and 121-130 | PASS with final-review note | Tests and source order cover the Task 3 VALs. The final QA lead should still collect literal CLI/DB row-count excerpts for QA-06..QA-11 because the task report primarily cites tests and journal line counts rather than pasted DB row-count output. |

## Per-Task Findings

### Task 3 - Retrofit permission enforcement onto every write command

#### Blocking

- `project-docs/docs/phases/phase-10/tasks/task-3-done.md:71` - The reviewed report's commit hash/subject fields are only present as an uncommitted working-tree edit. The committed Task 3 report in `5293fd2674b0afdc902feb2f7079d12b45e82551` still says "Pending commit", and `08_QUALITY_PLAN.md` §10 requires the completion report to be committed by the task author or an immediate docs commit before the task is done. Fix direction: commit or otherwise reconcile the Task 3 completion-report commit-field update before phase-end QA treats Task 3 evidence as closed. Reference: ESC-007.

#### Improvement

- None.

## Cross-cutting observations

- The implementation evidence tracks toward the behavior criteria: action keys are derived in the domain layer (`internal/issues/policy_actions.go:13`), scoped writes call `policy.Authorize` before projection/journal work, and the parity tests exercise both displayed permission booleans and actual write outcomes.
- The direct-prerequisite rows added in `internal/policy/matrix.go:40` and `project-docs/specs/03_TRD_SCHEMA.md:886` are documented in decisions lines 30-36 and avoid missing-entry errors for already-derived actions.
- The same commit-field drift is visible on Task 1 and Task 2 completion reports, but this review only blocks on Task 3 because this artifact covers Task 3 only. ESC-007 records the wider phase-runner cleanup need.

## Escalations

- ESC-007 (QA Reviewer trigger 7): Phase 10 task completion reports have uncommitted commit-field evidence; Task 3's committed report still contains pending placeholders while the working-tree report contains the reviewed hash/subject.

## Final Lines

Review written: project-docs/docs/phases/phase-10/qa/review-3.md
Blocking: 1  Improvement: 0  Escalations: ESC-007
VERDICT: BLOCK

BLOCK
