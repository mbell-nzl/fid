# Phase 20 - Mid-Phase Review 5 (Tasks 13-15)

Reviewer: qa-reviewer
Date: 2026-08-16
Tasks covered: phase-20/task-13, phase-20/task-14, phase-20/task-15
Review type: INITIAL
Prior review: none

## Scope Reviewed

- Task specs and completion reports:
  - `project-docs/docs/phases/phase-20/tasks/task-13.md`
  - `project-docs/docs/phases/phase-20/tasks/task-13-done.md`
  - `project-docs/docs/phases/phase-20/tasks/task-14.md`
  - `project-docs/docs/phases/phase-20/tasks/task-14-done.md`
  - `project-docs/docs/phases/phase-20/tasks/task-15.md`
  - `project-docs/docs/phases/phase-20/tasks/task-15-done.md`
- Phase contract and QA evidence plan:
  - `project-docs/docs/phases/phase-20/contracts.md` §§5.2-5.3, 6
  - `project-docs/docs/phases/phase-20/decisions.md` DEC-142
  - `project-docs/docs/phases/phase-20/qa/validation-contract.md`
  - `project-docs/docs/phases/phase-20/qa/qa-plan.md` §§4-5, Review 5
- Reconciled implementation commits:
  - `161f89ca18256cdc2a0ed3c9da6df094cdb8f427` — `feat(maintenance): migrate inspection response contract` (`Refs: phase-20/task-13`)
  - `f712c277b4debea077f0ae4173fad9a3590443c2` — `feat(cli): convert legacy JSON reads directly` (`Refs: phase-20/task-14`)
  - `f95db81f938e15d59357a35936e43a24c708b653` — `test(cli): converge format shim callers` (`Refs: phase-20/task-15`)

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-002, VAL-008, VAL-017, VAL-027, VAL-028, VAL-029, VAL-031, VAL-044, VAL-045

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-002 | task-13 | ON TRACK | `TestMaintenanceCheckJSONMatchesContract` decodes exact successful and `projection-divergence` failure envelope partitions and channels; the focused reproduction passed. |
| VAL-008 | task-13 | ON TRACK | `TestMaintenanceCheckHumanSuccessAndVacuumJSONRejection` asserts stdout-only Human success, and the failure branch asserts stderr-only output; the focused reproduction passed. |
| VAL-017 | task-13 | ON TRACK | The maintenance JSON test asserts `event:null`; the report and committed integration test retain read-response coverage. |
| VAL-027 | task-13 | NOT ON TRACK | The warning/promotion test compares `errors`, `warnings`, and `info` from the same fixture and the internal checker tests the mathematical pre-boundary/exact-boundary states, but the command-level config-independence test claims the required 299/300-second pair while actually seeds 298/300 seconds. QA-plan EVID-11 and Review 5 require the exact 299/300 evidence pair. MF-01 owns this gap. |
| VAL-028 | task-14 | ON TRACK | The reproduced direct-JSON/journal matrix passed. `journal_test.go` decodes retained event/location fields, excludes `revert`, `contract_version`, and envelope `event`; the root-flag table rejects both flags before the handler. |
| VAL-029 | task-14/task-15 | ON TRACK | The reproduced eight-read matrix decodes semantic direct JSON for policy, claim, activity, concern, both stats reads, SQL, and journal inspection. It passed, as did the mandatory integration gate. |
| VAL-031 | task-14/task-15 | NOT ON TRACK | Task 14's eight converted constructors have no `format` flag and the repository source guard/rejection tests pass, but the Task 15-required dynamic root command-tree guard has not been added. MF-02 owns the missing proof for zero universal runtime flags. |
| VAL-044 | task-13/task-14 | ON TRACK | Maintenance explicitly accepts a selected project root and rejects removed `--format`; the root-flag table exercises both `--json` and `--project` against all eight converted legacy reads before their handlers. |
| VAL-045 | task-13 | ON TRACK | Ordered structural/warning fixtures and the 1,001-warning typed over-limit/no-truncation test remain committed and passed in the focused maintenance reproduction. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task-spec and QA-plan command coverage | Six exact commands in each Task 13-15 report | PASS | All reports record the exact `gofmt`, lint, race/coverage, cover, dedicated integration, and build commands as PASS with command-produced tokens and observed `exit_code: 0`; no QA-plan command adds a different Task 13-15 gate. Task 14 retains its immutable first failed gate attempts and byte-identical closing passes under DEC-142. |
| Mechanical report verification | `project-docs/tools/report-lint.sh <task-done.md> <task.md>` for Tasks 13, 14, and 15 | PASS | Each current-session command ended `PASS report-lint completed with no failures` (terminated `exit_code: 0`), including committed-and-clean report verification. |
| Commit identity and diff hygiene | `git rev-parse --verify <hash>^{commit}`; `git show -s`; `git diff --check` for the three implementation commits | PASS | Each recorded hash resolved to the stated Conventional Commit subject with the required task `Refs` footer. All three committed diffs pass whitespace checking. |
| Required independent reproduction | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | From `/home/mark/projects/fid`: `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.992s` (terminated `exit_code: 0`). |
| Maintenance branch/bound/limit reproduction | `go test ./cmd/fid -run '^(TestMaintenanceCheckJSONMatchesContract|TestMaintenanceCheckWarningPromotionPreservesFindingsAndInfo|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestMaintenanceCheckHumanSuccessAndVacuumJSONRejection|TestMaintenanceCheckUsesExplicitProjectRoot|TestMaintenanceCheckRejectsFindingLimitWithoutTruncation)$' -count=1` | PASS, evidence gap found | `ok  github.com/mbell-nzl/fid/cmd/fid  0.356s` (terminated `exit_code: 0`). Inspection shows the test seeds `FID-299` at 298 seconds (`cmd/fid/maintenance_test.go:130-136`), not the QA-required 299-second state. |
| Legacy semantic JSON and root-flag reproduction | `go test ./cmd/fid -run '^(TestPolicyShowUsesDirectLegacyJSONWithoutFormatFlag|TestRemainingLegacyJSONReadsDoNotDeclareFormatFlags|TestClaimListUsesDirectLegacyJSONAndFiltersByStoredRole|TestActivityListUsesDirectLegacyJSON|TestConcernListUsesDirectLegacyJSON|TestStatsTokensUsesDirectLegacyJSONAndEpicFilter|TestStatsCostUsesDirectLegacyJSONAndExcludesNullCostsFromSum|TestSQLCommandAcceptsReadStatementsAsDirectLegacyJSON|TestJournalInspectCommandUsesDirectLegacyJSONWithoutRevert|TestJournalInspectCommandDoesNotExposeRevertOrResponseEnvelope|TestRootFlagsRejectBeforeEachConvertedLegacyJSONRead)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid  0.311s` (terminated `exit_code: 0`). The test matrix covers all eight named legacy reads and both explicit root flags. |
| Journal implementation reproduction | `go test ./internal/journal -count=1` | PASS | `ok  github.com/mbell-nzl/fid/internal/journal  1.930s` (terminated `exit_code: 0`). |
| Shim/rejection/source-guard reproduction | `go test ./cmd/fid -run '^(TestLiveTestsContainNoUniversalFormatInvocation|TestIssueGetRejectsRetiredFormatFlags|TestIssueGetRejectsAllFormatValues|TestIssueListRejectsRetiredFormatFlag|TestLabelListRejectsLegacyFormatFlag|TestReadyRejectsRetiredFormatFlag|TestIssueACListCommandRejectsLegacyFormatsAndHandlesEmptyAndMissingIssues|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestV01ReadCommandsRejectRetiredFormats|TestIssueReadRejectsRetiredAndMalformedInputs)$' -count=1` | PASS, evidence gap found | `ok  github.com/mbell-nzl/fid/cmd/fid  0.342s` (terminated `exit_code: 0`). The retained negative assertions and AST source guard pass; `contract_surface_test.go:14-61` does not construct/walk the Cobra tree. |
| Live non-test literal inventory | `rg -n --glob '*.go' '"(format|toon|yaml|markdown)"' cmd/fid -g '!**/*_test.go'` | PASS, insufficient alone | Expected no-match: `exit 1 (no output)` (terminated `exit_code: 1`); no-match is the pass condition. This static inventory cannot replace the missing dynamic command-tree guard. |

## Per-Task Findings

### Task 13 - Maintenance inspection contract

#### Must-fix

- [cmd/fid/maintenance_test.go:130-136; task-13-done.md:38] **MF-01 — the claimed 299/300-second command-level boundary evidence is weaker than the accepted evidence contract.** The command fixture labels the pre-boundary issue `FID-299`, but seeds it at 298 seconds. `validation-contract.md:50`, QA-plan EVID-11, QA-12, and the Review 5 focus expressly require 299/300-second cases; the report's claim that this test proves that exact pair is therefore false. Owning task: task-13. Repair boundary: add deterministic command-level 299-second non-qualifying and 300-second qualifying fixture evidence while preserving config-independence and the same-fixture warn/promoted comparison; change production code only if the exact test exposes a defect. Required evidence: decoded canonical command results show 299 absent and 300 present after an Attempt-config grace override, the existing warn and `--warn-as-error` branches retain identical findings/info, the Task 13 report accurately records the repair, and the unchanged six-command Task 13 gate passes. Re-review scope: inspect the exact 299/300 fixture, same-fixture branch assertions, report/diff, and rerun the focused maintenance matrix plus the mandatory integration command independently.

#### Improvement

- None.

### Task 14 - Legacy JSON conversion and caller convergence

#### Must-fix

- None.

#### Improvement

- None.

### Task 15 - Format shim removal and whole-tree caller convergence

#### Must-fix

- [cmd/fid/contract_surface_test.go:14-61; task-15.md:52-53; task-15-done.md:52-54] **MF-02 — the required dynamic root command-tree zero-format guard is absent.** `TestLiveTestsContainNoUniversalFormatInvocation` parses only test files for live `--format` arguments; `TestRemainingLegacyJSONReadsDoNotDeclareFormatFlags` checks the eight converted constructors only. Neither constructs and recursively walks the complete live Cobra tree to assert zero `format`, `toon`, `yaml`, or `markdown` flags, although Task 15 explicitly requires that guard and its report says it is included. Owning task: task-15. Repair boundary: add a focused command-tree guard in `cmd/fid` (and a narrow helper/test-only adjustment only if necessary) that traverses every live command and inherited/local flag set; retain the existing source guard and all rejection tests unchanged. Required evidence: the new guard fails against an injected universal flag or otherwise demonstrates a complete live-tree traversal, passes on the current tree, the retained rejected-format/root-flag tests remain green, the Task 15 report corrects the unsupported claim, and the unchanged six-command Task 15 gate passes. Re-review scope: inspect the guard's complete traversal and negative control, all retained rejection tests and source guard, then rerun the focused guard/rejection matrix plus the mandatory integration command independently.

#### Improvement

- None.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-01 | task-13 | VAL-027 | Deterministic command-level 299/300 fixture and its Task 13 report; production maintenance code only if that exact boundary exposes a defect. | Config-independent 299 absent/300 present decoded result; same-fixture warning/promotion equality; unchanged six gates with observed termination. | New reviewer inspects the fixture/report/diff, reruns the focused maintenance matrix and `go test -tags integration -race ./cmd/fid/integration/...`. |
| MF-02 | task-15 | VAL-031 | Add the complete Cobra command-tree universal-flag guard and correct Task 15 evidence; preserve existing source and rejection guards. | Traversal over all live commands and inherited/local flags, a demonstrated negative control, retained rejection/root-flag tests, and unchanged six gates with observed termination. | New reviewer inspects the guard/report/diff, reruns the focused command-tree/rejection suite and `go test -tags integration -race ./cmd/fid/integration/...`. |

## Cross-cutting observations

- DEC-142's caller-ownership repair is consistent with the accepted amended boundary: Task 14's recovered lint, full-race, and integration caller repairs are in its commit and report; Task 15 does not reopen those callers. The three Task 14 first failed attempts are preserved and followed by byte-identical passes. No basis was found to reopen resolved ESC-20-011.
- The Task 14 tests assert observable legacy JSON semantics and the journal boundary rather than byte-compatibility: `journal inspect` has retained event/location fields, is unenveloped, excludes revert semantics, and rejects both root flags before its handler.
- The source inventory is clean apart from intentional test/rejection references, and the committed Task 15 diff removes the migrated argument/output shim names. That static state is useful supporting evidence, but not a substitute for the required dynamic command-tree guard.

## Blocker Candidates

- None. Both findings are bounded evidence/test and completion-report corrections inside reviewed existing tasks. Neither changes the accepted product, architecture, public interface, dependency budget, security posture, validation contract, or quality-gate definition.

## Existing Escalation References

- None. `ESC-20-011` is resolved by DEC-142 and is not reopened by either finding.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-20/qa/review-5.md
Must-fix: 2  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
```
