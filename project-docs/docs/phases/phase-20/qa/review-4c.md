# Phase 20 - Mid-Phase Review 4c (task 12)

Reviewer: qa-reviewer
Date: 2026-08-15
Tasks covered: phase-20/task-12
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-20/qa/review-4.md

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-017, VAL-026, VAL-044, VAL-045

## Scope Reviewed

- This fresh review is limited to Review 4 MF-02 and its Task 12 repair.
- `project-docs/docs/phases/phase-20/qa/review-4b.md` was the Task 11-only
  MF-01 re-review; it did not disposition MF-02 and is not otherwise reopened
  here.
- Reconciled Task 12's repair and report commits:
  - `49c179add40724871bae582fc3812ffc0aff04fc` — `fix(workflow): enforce workflow list limit` (`Refs: phase-20/task-12`)
  - `899c445c8e5fc7d71dba23a0fd871128a2fe3acc` — `docs(task): finalize completion report for phase-20 task-12`
  - `bac3a5045c8d3d75f8ed55ec9f0fe4616228eb84` — `docs(task): correct phase-20 task-12 completion report`
- Each required hash resolves with `git rev-parse --verify <hash>^{commit}`.
  The repair commit has the required footer; the final report is committed and
  has no worktree diff. The only pre-existing dirty file is the out-of-scope
  `project-docs/docs/phases/phase-20/run-ledger.md`, which this review did not
  modify.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Fresh evidence |
|---|---|---|---|
| VAL-017 | task-12 | ON TRACK | Define and assign compare both returned Event ID and `occurred_at` to the journal tail (`cmd/fid/workflow_define_test.go:44-47`, `cmd/fid/workflow_assign_test.go:42-45`). Successful list reads assert `event:null`, including the 10,000-item state (`cmd/fid/workflow_list_test.go:79-95`). |
| VAL-026 | task-12 | ON TRACK | The committed matrix verifies closed define/assign/list schemas, `--require` rejection, invalid grammar/scope before append, kind-over-project resolution, a root `--project <path>` combined with `--project-wide`, and absence of local `workflow assign --project` (`cmd/fid/workflow_define_test.go:50-78`, `cmd/fid/workflow_assign_test.go:14-119`, `cmd/fid/workflow.go:70-120`). |
| VAL-044 | task-12 | ON TRACK | Removed authoring, legacy applicability, malformed CSV/arrow, over-limit state input, absent/both/invalid assignment scope, unknown workflow, and non-Issue list kind have typed command-boundary rejection evidence. The define and assignment rejection cases compare journal bytes before and after (`cmd/fid/workflow_define_test.go:50-76`, `cmd/fid/workflow_assign_test.go:61-90`, `cmd/fid/workflow_list_test.go:59-67`). The domain duplicate path also preserves journal bytes (`internal/workflow/define_test.go:80-84`). |
| VAL-045 | task-12 | ON TRACK | `workflow.list` orders SQL results by `name, id` and requests at most 10,001 rows (`internal/workflow/list.go:104-125`). The domain test proves the equal-name ID tie-break and 10,000/10,001 boundary (`internal/workflow/list_test.go:61-116`); the command test proves exactly 10,000 complete rows succeed and 10,001 returns the closed `result-limit-exceeded` envelope, leaves the journal unchanged, and leaves all 10,001 projection rows intact (`cmd/fid/workflow_list_test.go:69-127`). |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Task-spec/QA-plan command coverage | Task 12 exact six-command gate | PASS | `task-12-done.md:58-63` records each byte-identical required command as PASS with a command-produced token and observed `exit_code: 0`; no QA-plan command adds a different Task 12 gate. |
| Report outcome, deviations, and provenance | `task-12-done.md` and its final commits | PASS | Outcome is COMPLETE; the report retains the focused red/green cap and fixture-repair attempts and its report-lint correction (`task-12-done.md:42-69`). No specification or gate-contract deviation is claimed. |
| Mechanical report verification | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-12-done.md project-docs/docs/phases/phase-20/tasks/task-12.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). It verifies all six exact required commands, the repair hash/subject, and that the report file is committed and clean. |
| Reproduced focused CLI matrix | `go test ./cmd/fid -run '^(TestWorkflowDefineJSONMatchesContract|TestWorkflowDefineRejectsRemovedAndInvalidRequestsBeforeAppend|TestWorkflowAssignJSONPrecedenceAndScopedReplacement|TestWorkflowAssignRequiresExactlyOneScopeAndMapsUnknownWorkflow|TestWorkflowAssignCombinesRootProjectSelectorAndProjectWideScope|TestWorkflowListJSONMatchesContractAndUsesNullAssignmentScope|TestWorkflowListRejectsNonIssueKind|TestWorkflowListEnforcesCollectionLimitWithoutTruncation)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid  1.431s` (terminated `exit_code: 0`). |
| Reproduced workflow-list domain matrix | `go test ./internal/workflow --run '^(TestListFiltersCountsAssignmentsAndOrdering|TestListReturnsEmptySliceForNoMatch|TestListOrdersWorkflowsByNameThenID|TestListRejectsCollectionAboveLimit)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/internal/workflow  1.459s` (terminated `exit_code: 0`). |
| Required independent review-window reproduction | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.646s` (terminated `exit_code: 0`). |

## Task 12 - Workflow command contract

### Must-fix

- None.

### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-02 — `workflow.list` lacked the 10,000-item contract bound and Task 12 lacked the required list/order, Event/null, root-project/project-wide, and rejection evidence matrix. | task-12 | RESOLVED | Commit `49c179a` adds the 10,001-row query cap and typed CLI mapping (`internal/workflow/list.go:12`, `:72-125`; `cmd/fid/workflow.go:232-244`), without changing the accepted contract or gate. The independently rerun CLI and domain matrices passed, including qualifying 10,000 and non-qualifying 10,001 states. Their committed assertions also prove `name,id` ordering, no truncation, unchanged journal/projection on over-limit read, define/assign tail Event identity, list null Event, no local flag shadow, forbidden `--require`, grammar/scope rejection without append, and kind precedence. The required integration command passed independently. |

## Cross-cutting observations

- The error is a bounded typed domain error (`ResultLimitExceededError`) mapped
  only by the workflow adapter, so the repair remains entirely inside Task 12.
- Tests assert observable command envelopes, journal-tail identity, and
  projection/journal invariance rather than only implementation details.

## Blocker Candidates

- None.

## Existing Escalation References

- None.

## Final Lines

Review written: project-docs/docs/phases/phase-20/qa/review-4c.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
