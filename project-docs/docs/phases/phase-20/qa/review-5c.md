# Phase 20 - Mid-Phase Review 5c (Task 15, Review 5 MF-02 only)

Reviewer: qa-reviewer
Date: 2026-08-16
Tasks covered: phase-20/task-15
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-20/qa/review-5.md

## Scope Reviewed

- This fresh independent re-review is limited to Review 5 MF-02 and the
  bounded Task 15 repair. It does not re-review or disposition Review 5 MF-01;
  `review-5b.md` remains its closed re-review record.
- Task 15 specification and completion report:
  - `project-docs/docs/phases/phase-20/tasks/task-15.md`
  - `project-docs/docs/phases/phase-20/tasks/task-15-done.md`
- Phase evidence/contract context:
  - `project-docs/docs/phases/phase-20/qa/validation-contract.md` (VAL-031)
  - `project-docs/docs/phases/phase-20/qa/qa-plan.md` (EVID-12/EVID-13 and
    Task 15's six-command gate)
  - `project-docs/docs/phases/phase-20/contracts.md` §§5.2-5.3, 6
  - `project-docs/docs/phases/phase-20/decisions.md` (DEC-142)
- Repair commits reconciled with `git rev-parse --verify <hash>^{commit}`:
  - `b720fd174fae806449be6ff28cf3ec8a9d237f07` —
    `test(cli): add dynamic Cobra format guard` (`Refs: phase-20/task-15`)
  - `d7db9fce9e627df823f8f5687f2d8c6d79ae1f78` —
    `docs(task): finalize MF-02 repair report` (`Refs: phase-20/task-15`)

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-031

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Fresh evidence |
|---|---|---|---|
| VAL-031 | task-15 | NOT ON TRACK | `liveUniversalFormatFlagViolations` recursively walks `newRootCommandAt` and directly calls `LocalNonPersistentFlags`, `PersistentFlags`, and `InheritedFlags` (`cmd/fid/contract_surface_test.go:115-135`). The clean-tree and local-leaf controls pass, but the root-persistent injection is accepted as `fid --format` (`:76-85`), which comes from the root's own persistent set. It does not require a descendant inherited-flag violation, so it would still pass if the `InheritedFlags` inspection were removed. |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Repair commit identity and hygiene | `git rev-parse --verify <hash>^{commit}`; `git show -s`; `git diff --check` for both repair commits | PASS | Both supplied hashes resolve to the stated Conventional Commit subjects and Task 15 `Refs` footer. The implementation commit changes only the dynamic guard and Task 15 report; the follow-up only corrects the report commit identity. No whitespace error was reported. |
| Dynamic live-tree guard | `go test ./cmd/fid -run '^TestLiveCobraTree' -count=1` | PASS, insufficient negative evidence | `ok  github.com/mbell-nzl/fid/cmd/fid  0.053s` (terminated `exit_code: 0`). The clean tree has no named flag and both current controls return a violation. The inherited-root assertion is not category-specific: its expected `fid --format` is found on the root's persistent flags before any child inherited set is needed. |
| Retained source/rejection/root-flag/Task 14 matrix | `go test ./cmd/fid -run '^(TestLiveTestsContainNoUniversalFormatInvocation|TestIssueGetRejectsRetiredFormatFlags|TestIssueGetRejectsAllFormatValues|TestIssueListRejectsRetiredFormatFlag|TestLabelListRejectsLegacyFormatFlag|TestReadyRejectsRetiredFormatFlag|TestIssueACListCommandRejectsLegacyFormatsAndHandlesEmptyAndMissingIssues|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestV01ReadCommandsRejectRetiredFormats|TestIssueReadRejectsRetiredAndMalformedInputs|TestRootFlagsRejectBeforeEachConvertedLegacyJSONRead|TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag|TestRemainingLegacyJSONReadsDoNotDeclareFormatFlags|TestClaimListUsesDirectLegacyJSONAndFiltersByStoredRole|TestActivityListUsesDirectLegacyJSON|TestConcernListUsesDirectLegacyJSON|TestStatsTokensUsesDirectLegacyJSONAndEpicFilter|TestStatsCostUsesDirectLegacyJSONAndExcludesNullCostsFromSum|TestSQLCommandAcceptsReadStatementsAsDirectLegacyJSON|TestJournalInspectCommandUsesDirectLegacyJSONWithoutRevert|TestJournalInspectCommandDoesNotExposeRevertOrResponseEnvelope)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid  0.554s` (terminated `exit_code: 0`). The AST source guard, all named retired-format rejection tests, explicit root-flag rejection, and every Task 14 eight-read semantic guard remain present and pass. |
| Required independent review-window reproduction | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | From `/home/mark/projects/fid`: `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.822s` (terminated `exit_code: 0`). |
| Task 15 report and exact six-command coverage | `task-15-done.md:64-155` against `task-15.md:85-94` | PASS except finding below | The report accurately names VAL-031's walker and retains prior red-first evidence (`:84`), the six byte-identical commands with observed termination (`:99-112`), and repair commit identity (`:149-155`). Its claim that the inherited control proves that category (`:73-74`, `:131-135`) overstates the current assertion. |
| Mechanical completion-report verification | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-15-done.md project-docs/docs/phases/phase-20/tasks/task-15.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). It confirms the exact six commands, their recorded termination, and the reported `b720fd...` identity on the committed-clean report. |

## Per-Task Findings

### Task 15 - Format shim removal and whole-tree caller convergence

#### Must-fix

- [cmd/fid/contract_surface_test.go:76-85; task-15-done.md:73-74,131-135] **MF-02 remains open — the purported inherited-root negative control does not prove traversal of an inherited flag set.** It injects `format` into `root.PersistentFlags()` and accepts `fid --format`; the root's own `PersistentFlags()` inspection produces that result, so removing `command.InheritedFlags()` from the walker would leave this control green. Owning task: task-15. Repair boundary: strengthen only the inherited-root negative test to require the injected root persistent flag on a descendant command's inherited flag set (for example, `fid issue --format`), retain the local-leaf control, and correct the Task 15 report's evidence claim. Required evidence: the enhanced control fails if inherited inspection is absent and passes on the complete current tree; the AST/rejection/root-flag/Task 14 matrix remains green; all unchanged exact six Task 15 gate commands and report-lint pass with observed termination. Re-review scope: a fresh QA Reviewer inspects the category-specific inherited and local controls, Task 15 report/diff and repair commits, then reruns the focused guard, retained matrix, and `go test -tags integration -race ./cmd/fid/integration/...` independently.

#### Improvement

- None.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-02 — no dynamic root command-tree guard proved the complete live Cobra tree has zero universal format flags. | task-15 | STILL OPEN | Commit `b720fd...` adds a recursive walker over every `newRootCommandAt` command and all three required flag sets; the fresh guard and retained matrix pass. Its local-leaf injection is category-specific, but the root-persistent injection asserts only the root's `PersistentFlags()` result. The absence condition is therefore covered, while the required inherited-set negative evidence remains incomplete. |

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-02 | task-15 | VAL-031 | In `cmd/fid/contract_surface_test.go`, make the root persistent injection assert a descendant inherited violation; update only the corresponding Task 15 report wording/evidence. Preserve the walker, local-leaf control, AST source guard, rejection tests, and Task 14 guards. | Category-specific inherited-root and local-leaf negative controls; clean-tree pass; retained matrix; six byte-identical Task 15 gates; report-lint with real termination. | Inspect the changed assertions/report and commit identity; rerun the focused guard, retained matrix, and required integration command. |

## Cross-cutting observations

- `liveUniversalFormatFlagViolations` itself currently invokes all three required
  Cobra accessors and recurses through `command.Commands()`. The finding is
  limited to the strength of the inherited-root control, not a claim that the
  implementation currently omits that accessor.
- Review 5 MF-01 is not reconsidered here. `review-5b.md` remains the fresh
  independent closure record for Task 13.

## Blocker Candidates

- None. The stronger control and report correction are bounded, decision-free
  Task 15 test/documentation work. They require no product, architecture,
  data-model, public-interface, dependency, security, validation-contract, or
  gate-definition authority.

## Existing Escalation References

- None. `ESC-20-011` is resolved by DEC-142 and is not reopened by this
  in-scope evidence repair.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-20/qa/review-5c.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
