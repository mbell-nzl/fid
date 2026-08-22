# Phase 20 - Mid-Phase Review 5d (Task 15, Review 5c MF-02 only)

Reviewer: qa-reviewer
Date: 2026-08-16
Tasks covered: phase-20/task-15
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-20/qa/review-5c.md

## Verdict

PASS

This fresh independent re-review is limited to Review 5c MF-02. It does not
re-review or disposition Review 5 MF-01; `review-5b.md` remains its closure
record. `ESC-20-011` remains resolved by DEC-142 and is not reopened.

## Scope Reviewed

- Task specification and completion report:
  - `project-docs/docs/phases/phase-20/tasks/task-15.md`
  - `project-docs/docs/phases/phase-20/tasks/task-15-done.md`
- Phase evidence and accepted repair boundary:
  - `project-docs/docs/phases/phase-20/qa/validation-contract.md` (VAL-031)
  - `project-docs/docs/phases/phase-20/qa/qa-plan.md` (EVID-12, EVID-13,
    Review 5, and Task 15's exact six-command gate)
  - `project-docs/docs/phases/phase-20/contracts.md` §§5.2-5.3, 6
  - `project-docs/docs/phases/phase-20/decisions.md` (DEC-142)
- Prior review artifacts:
  - `project-docs/docs/phases/phase-20/qa/review-5.md`
  - `project-docs/docs/phases/phase-20/qa/review-5b.md`
  - `project-docs/docs/phases/phase-20/qa/review-5c.md`
- Repair commits, each resolved with `git rev-parse --verify <hash>^{commit}`:
  - `24f50aef13ae7a23a080a87920f396242b30f911` —
    `test(cli): strengthen inherited format guard`
    (`Refs: phase-20/task-15`)
  - `97d0aaf4c68bb86cbe568a3016ced5c1b36203f1` —
    `docs(task): record inherited format guard repair`
    (`Refs: phase-20/task-15`)

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-031

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Fresh evidence |
|---|---|---|---|
| VAL-031 | task-15 | ON TRACK | `liveUniversalFormatFlagViolations` recursively traverses every `newRootCommandAt` command and inspects `LocalNonPersistentFlags`, `PersistentFlags`, and `InheritedFlags` (`cmd/fid/contract_surface_test.go:115-135`). The clean-tree assertion passes. The repaired root-persistent injection now requires the descendant inherited violation `fid issue --format` (`:76-85`), so omission of `InheritedFlags()` would leave only root `fid --format` and fail the control. The retained local-leaf `fid ready --yaml` control remains intact (`:87-105`). |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Repair commit identity and diff hygiene | `git rev-parse --verify <hash>^{commit}`; `git show -s`; `git diff --check` for `24f50ae...` and `97d0aaf...` | PASS | Both supplied hashes resolve to the stated Conventional Commit subjects with the required Task 15 footer. `24f50ae...` changes only the inherited-control assertion; `97d0aaf...` records the bounded report repair. Neither diff reports whitespace errors. |
| Review 5c red finding and corrected report | `task-15-done.md:28-34,72-80,90-95,122-127,152-159,172-178` | PASS | The report preserves Review 5c's immutable root-only red finding, accurately states the strengthened descendant-path evidence, names VAL-031, retains the AST/rejection evidence, records every unchanged exact gate with observed termination, and records the implementation repair hash/subject. |
| Inherited/local negative controls | `go test ./cmd/fid -run '^TestLiveCobraTreeUniversalFormatGuardDetectsInjectedFlags$' -count=1` | PASS | From `/home/mark/projects/fid`: `ok  github.com/mbell-nzl/fid/cmd/fid  0.052s` (terminated `exit_code: 0`). This runs both the root-persistent descendant inherited control and the retained local-leaf control. |
| Retained guard/rejection and Task 14 regression matrix | `go test ./cmd/fid -run '^(TestLiveCobraTreeContainsNoUniversalFormatFlags|TestLiveCobraTreeUniversalFormatGuardDetectsInjectedFlags|TestLiveTestsContainNoUniversalFormatInvocation|TestIssueGetRejectsRetiredFormatFlags|TestIssueGetRejectsAllFormatValues|TestIssueListRejectsRetiredFormatFlag|TestLabelListRejectsLegacyFormatFlag|TestReadyRejectsRetiredFormatFlag|TestIssueACListCommandRejectsLegacyFormatsAndHandlesEmptyAndMissingIssues|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestV01ReadCommandsRejectRetiredFormats|TestIssueReadRejectsRetiredAndMalformedInputs|TestRootFlagsRejectBeforeEachConvertedLegacyJSONRead|TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag|TestRemainingLegacyJSONReadsDoNotDeclareFormatFlags|TestClaimListUsesDirectLegacyJSONAndFiltersByStoredRole|TestActivityListUsesDirectLegacyJSON|TestConcernListUsesDirectLegacyJSON|TestStatsTokensUsesDirectLegacyJSONAndEpicFilter|TestStatsCostUsesDirectLegacyJSONAndExcludesNullCostsFromSum|TestSQLCommandAcceptsReadStatementsAsDirectLegacyJSON|TestJournalInspectCommandUsesDirectLegacyJSONWithoutRevert|TestJournalInspectCommandDoesNotExposeRevertOrResponseEnvelope)$' -count=1` | PASS | From `/home/mark/projects/fid`: `ok  github.com/mbell-nzl/fid/cmd/fid  0.494s` (terminated `exit_code: 0`). The clean-tree guard, retained AST source guard, retired-format rejection assertions, explicit root-flag rejection, and eight-read semantic guard matrix remain green. |
| Mechanical completion-report verification | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-15-done.md project-docs/docs/phases/phase-20/tasks/task-15.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). It verifies the exact six required commands, recorded command tokens/terminations, repair hash and subject, and committed-clean Task 15 report. |
| Required independent review-window reproduction | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | From `/home/mark/projects/fid`: `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.762s` (terminated `exit_code: 0`). |

## Per-Task Findings

### Task 15 - Format shim removal and whole-tree caller convergence

#### Must-fix

- None.

#### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| Review 5c MF-02 — the root-persistent negative control accepted only `fid --format`, so it did not prove the walker inspected a descendant inherited flag set. | task-15 | RESOLVED | Commit `24f50ae...` changes the required violation to `fid issue --format` (`cmd/fid/contract_surface_test.go:82-84`). The recursive walker still calls `InheritedFlags()` and recurses through every child (`:115-135`); the retained local-leaf `fid ready --yaml` control is unchanged (`:87-105`). Fresh focused control, retained AST/rejection/Task 14 matrix, report-lint, and the required integration command all passed with observed `exit_code: 0`. |

## Cross-cutting observations

- No new finding arises from this bounded repair. The strengthened control is
  category-specific: a root persistent flag can satisfy `fid issue --format`
  only through a descendant's inherited flag set, while the separate local-leaf
  control proves local non-persistent inspection remains live.

## Blocker Candidates

- None.

## Existing Escalation References

- None. `ESC-20-011` remains resolved by DEC-142 and is outside this closure
  review.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-20/qa/review-5d.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: CONTINUE
```
