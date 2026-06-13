# Phase 5 - Mid-Phase Review 4 (Tasks 5-5)

Reviewer: qa-reviewer
Date: 2026-06-13
Tasks covered: phase-5/task-5

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-5/tasks/task-5.md`
- Completion reports:
  - `project-docs/docs/phases/phase-5/tasks/task-5-done.md`
- Prior review:
  - `project-docs/docs/phases/phase-5/qa/review-3.md`
- Validation contract:
  - `project-docs/docs/phases/phase-5/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-5/qa/qa-plan.md`
- Escalation log:
  - `project-docs/docs/current/tl_escalations_current.md`
- Design sections checked:
  - `project-docs/specs/05_DESIGN_PLAN.md` §5.4 and §5.6 for rich transition errors and stderr write confirmations.
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md` §`fid issue transition` for command shape, required flags, and confirmation output.
- Commits or changed files:
  - `0173eba6d4fcff718a76e89a902f74bbfecce1da` (`feat(issue): add workflow transitions`)
  - `internal/issues/transition.go`
  - `internal/issues/transition_test.go`
  - `cmd/fid/issue.go`
  - `cmd/fid/issue_transition_test.go`
  - `project-docs/docs/phases/phase-5/tasks/task-5-done.md`

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-005, VAL-006, VAL-007, VAL-008, VAL-009, VAL-010, VAL-011

Task 5 is on track for final QA. The implementation wires the Task 4 workflow graph/rules into `internal/issues.Transition`, keeps the CLI adapter thin, validates `--by` through known-role policy, builds one status update effect plus one `activity` effect for successful transitions, and uses the dry-run-before-append seam before committing the journal event. Fixture feasibility is acceptable for this phase: tests seed acceptance criteria, AC verifications, completion reports, and child statuses via generated journal events, which matches the Task 5 allowance while v0.1c write commands remain deferred.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-005 | task-5 advances | ON TRACK | Rejected transition paths assert unchanged journal line counts and later reads: invalid edge in `internal/issues/transition_test.go:162` and CLI invalid edge in `cmd/fid/issue_transition_test.go:43`; unmet developer rule in `internal/issues/transition_test.go:232`; unmet qa_lead rule in `internal/issues/transition_test.go:331`; blocked epic wrap in `internal/issues/transition_test.go:423`; cancelled re-transition in `internal/issues/transition_test.go:530`. Successful accepted transition appends exactly one line in `internal/issues/transition_test.go:68`. |
| VAL-006 | task-5 | ON TRACK | Task allowed edge read-back is covered by package test `internal/issues/transition_test.go:13` and CLI test `cmd/fid/issue_transition_test.go:11`, including `in_progress` / `started`. Epic allowed transition to `wrapped` / `started` is covered by `internal/issues/transition_test.go:423`. Subtask allowed transition to `in_progress` / `started` is covered by `internal/issues/transition_test.go:530`. The successful transition event shape and rebuild parity are checked in `internal/issues/transition_test.go:72` and `internal/issues/transition_test.go:145`. |
| VAL-007 | task-5 | ON TRACK | `internal/issues/transition_test.go:162` rejects `todo -> done`, checks the error names `"todo"` and `"done"`, verifies unchanged read-back, and verifies unchanged journal line count. `cmd/fid/issue_transition_test.go:43` covers the same invalid edge through `fid issue transition` and later `fid issue get --format json`. |
| VAL-008 | task-5 | ON TRACK | `internal/issues/transition_test.go:232` seeds AC fixtures, rejects `in_progress -> dev_complete` while developer verifications and completion report are missing, verifies `all_acs_verified`, `has_completion_report`, and AC IDs in the error, verifies zero append on rejection, then seeds AC verifications plus a completion report and succeeds. CLI rendering of rule names and AC ID is represented by `cmd/fid/issue_transition_test.go:89`. |
| VAL-009 | task-5 | ON TRACK | `internal/issues/transition_test.go:331` drives a task to `qa_review`, rejects `qa_review -> qa_pass` until a `qa_lead` verification is seeded, checks `all_acs_verified`, `missing roles qa_lead`, and AC ID in the error, verifies zero append, then succeeds and reads back `qa_pass` / `completed`. |
| VAL-010 | task-5 | ON TRACK | `internal/issues/transition_test.go:423` creates an epic with two children, rejects `in_progress -> wrapped` while both children are blocking, verifies `all_child_tasks_done_or_cancelled` plus both child IDs in the error, verifies zero append, then seeds child statuses to `done` and `cancelled` and succeeds. |
| VAL-011 | task-5 | ON TRACK | `internal/issues/transition_test.go:530` verifies cancellation from two distinct non-cancelled states (`task todo -> cancelled` and `subtask in_progress -> cancelled`) and rejects `cancelled -> cancelled` with unchanged journal line count. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 5 completion report path | `project-docs/docs/phases/phase-5/tasks/task-5-done.md` | PASS | Path exists and records summary, files, documentation update/verified-no-change, VAL coverage, TDD evidence, exact gate results, dependency note, commit hash, and subject. |
| Task 5 commit verification | `git rev-parse --verify 0173eba^{commit}` | PASS | Resolved to `0173eba6d4fcff718a76e89a902f74bbfecce1da`. |
| Task 5 commit subject/body | `git show --no-patch --format=fuller 0173eba` | PASS | Subject is `feat(issue): add workflow transitions`; body includes `Refs: phase-5/task-5`. |
| Required Task 5 gates recorded | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | PASS | `task-5-done.md:81` records all four exact task-spec and QA-plan commands as PASS. The conditional `go build -o /tmp/fid ./cmd/fid` is recorded as not run because evidence came from automated CLI tests, not a manual CLI smoke. |
| Coverage thresholds recorded | `go test -race -coverprofile=coverage.out ./...` | PASS | `task-5-done.md:78` records `cmd/fid` 79.0%, `internal/issues` 82.8%, and `internal/workflow` 95.5%, satisfying the Task 5 COV-001/COV-002 targets. |
| Reviewer deterministic check | `go test ./internal/issues -run 'TestTransition(AllowsWorkflowEdgeProjectsActivityAndRebuilds\|RejectsInvalidEdgeWithoutAppending\|RequiresDeveloperVerificationsAndCompletionReport\|RequiresQALeadVerificationForQAPass\|RequiresDoneOrCancelledChildrenToWrapEpic\|SupportsSubtaskWorkflowAndCancellationRules)$' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/issues 0.148s`. |
| Reviewer deterministic CLI check | `go test ./cmd/fid -run 'TestIssueTransition(ChangesStateAndPrintsConfirmation\|RejectsInvalidEdgeAndLeavesStateUnchanged\|RendersRuleFailures\|RequiresTargetStateAndKnownByRole)$' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/cmd/fid 0.088s`. |
| Reviewer format spot-check | `gofmt -l cmd internal` | PASS | No files listed. |
| Fixture feasibility | Journal seed helpers in `internal/issues/transition_test.go:647`, `internal/issues/transition_test.go:678`, `internal/issues/transition_test.go:710`, and `internal/issues/transition_test.go:738` | PASS | Generated journal fixtures seed ACs, AC verifications, completion reports, and child status updates. This matches `task-5.md:15` and `qa-plan.md:185`; no v0.1c write command was added. |
| Design parity check | `05_DESIGN_PLAN.md` §5.4/§5.6; `06_DESIGN_PLAN_CLI.md` §`fid issue transition` | PASS | CLI test `cmd/fid/issue_transition_test.go:11` verifies stderr confirmation. CLI invalid and rule-error tests verify current/requested states, rule names, and entity IDs. Known-role errors reuse the existing role renderer in `cmd/fid/issue.go:355`. |
| Escalation scan | `project-docs/docs/current/tl_escalations_current.md` | PASS | No new open ESC appeared in this review window. Existing ESC-003 remains in the current log and is not newly raised by Task 5. |

## Per-Task Findings

### Task 5 - `fid issue transition` with workflow + rule enforcement

#### Blocking

- None.

#### Improvement

- `internal/issues/transition.go:97` The main transition service is doing too much in one non-test function: option validation, lock setup, DB open, state load, edge/rule evaluation, event assembly, dry-run, append, hash, and projection all sit in `Transition`, and the new file is 426 lines against the 400-line non-test file budget in `08_QUALITY_PLAN.md` §4.4. Fix direction: split event assembly and append/project orchestration into named helpers before the next transition-adjacent change so the workflow service remains reviewable.

## Cross-cutting observations

- Business behavior placement is appropriate for the window: `cmd/fid/issue.go:226` only gathers flags and delegates to `internal/issues.Transition`; workflow graph/rule semantics remain in `internal/workflow`; projection and append mechanics stay behind `internal/journal`.
- Transition rejection errors carry the required observable details for deterministic QA: invalid edges name current/requested states, rule failures name rule IDs and affected AC/child IDs or missing verifier roles.
- Security/privacy risk is low for this window: DB reads use parameterized SQL, the CLI does not shell out, and journal details include workflow state/role metadata rather than secrets.
- Documentation accuracy is acceptable: `task-5-done.md:28` records `06_DESIGN_PLAN_CLI.md` as verified-no-change, and the implemented command shape/confirmation matches the cited CLI design section.

## Escalations

- None raised.
- No new open ESC appeared in this review window. Existing ESC-003 remains listed in `project-docs/docs/current/tl_escalations_current.md`, but Task 5 does not add a new escalation condition.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-5/qa/review-4.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```
