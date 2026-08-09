# Phase 20 - Mid-Phase Review 1b (DEC-130 recovery: tasks 3a-3b)

Reviewer: qa-reviewer
Date: 2026-08-09
Tasks covered: phase-20/task-3a, phase-20/task-3b

## Verdict

BLOCK

## Scope Reviewed

- Sustained Review 1 Blocking finding A: exhaustive VAL-002 reason-code detail value/null/nested-schema, subject, and redaction evidence.
- Sustained Review 1 Blocking finding B: exhaustive VAL-043 strict-config and init boundary evidence.
- Review 1 Medium improvement only: committed dangling `.fid` symlink-collision regression.
- Recovery commits:
  - `51763124db2cd6f6c6d00d15548de37532bbc273` — `fix(cli): enforce closed reason-code details`
  - `4f308c864cd9dd526886bf2469a641f7f2250efa` — `test(config): close strict init boundary matrix`
- Completion reports: `tasks/task-3a-done.md`, `tasks/task-3b-done.md`.
- Validation contract and QA-plan recovery evidence: VAL-002, VAL-039, VAL-042, VAL-043; EVID-01 and EVID-16.

`git rev-parse --verify` resolved both recovery hashes. Their subjects and
`Refs: phase-20/task-3a` / `Refs: phase-20/task-3b` footers match. Independent
`report-lint.sh` passed for both committed completion reports, including all six
required task/QA-plan commands. DEC-130 keeps VAL-002 and VAL-043 unchanged,
neither deferred nor closable from source inspection; this review applies that
decision. The current open escalation register contains no open ESC IDs.

## Summary

- Blocking findings: 2
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-002, VAL-039, VAL-042, VAL-043

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-002 | task-3a; task-3b supporting collision capture | ON TRACK for Review 1 finding A | `internal/cli/codes_test.go:37-130` declares all 23 codes and checks registry order/category/retryability/subjects. Lines 268-315 exercise valid and invalid scalar/null/nested/object-list details, per-code subject rejection, equal-length duplicate identity, and closed-detail redaction; lines 355-433 decode the exact serialized failure shape. The targeted current-session test passed. Commit `51763124` changes product code only to reject unequal `duplicate-entity.identity.fields`/`values`, the DEC-130-authorized repair. |
| VAL-039 | task-3b supporting | ON TRACK | `internal/config/config_test.go:189-213` covers normal and `math.MaxInt64` byte-stable writer/parser round trips; successful init rows use the strict parser against the activated config at `cmd/fid/init_test.go:331-341`. |
| VAL-042 | task-3b supporting | ON TRACK for Review 1 Medium improvement | `cmd/fid/init_test.go:358-380` creates a disposable dangling `.fid` symlink, asserts non-zero typed `project-already-initialized`, verifies the identical link target with `Readlink` and `Lstat`, and proves no staging residue. The targeted current-session test passed. |
| VAL-043 | task-3b | NOT ON TRACK | The parser/init matrix covers the structural cases and most domain pairs, but it does not independently exercise rejected lower-bound, type, or overflow values for `executor_limits.max_excerpt_bytes`; see the Task 3b Blocking finding. The required current-session integration gate also fails, so the report's passing gate record cannot be independently reproduced. |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Recovery report command coverage | `project-docs/tools/report-lint.sh task-3a-done.md task-3a.md`; corresponding task-3b invocation | PASS | Each report records all six exact commands required by its spec and the QA plan; lint also verified committed/clean reports, hashes, and subjects. |
| Scoped test reproduction | `go test -race ./internal/cli ./internal/config ./cmd/fid -run 'TestReasonCodeDetailSchemaIsExactForEveryCode|TestDuplicateEntityRequiresEqualLengthIdentityArrays|TestParseProjectConfigBoundaryDomains|TestProjectConfigWriterRoundTripsByteStableAtInt64Boundaries|TestReadProjectConfigRejectsArtifactRootSymlinkComponentEscape|TestInitCommandPrefixBoundaryMatrix|TestInitCommandPreservesDanglingWorkspaceSymlinkCollision' -count=1` | PASS | Current session: `ok github.com/mbell-nzl/fid/internal/cli`, `internal/config`, and `cmd/fid`; the changed tests execute behavioral registry/parser/CLI boundaries. |
| Required independent reproduction (repository root) | `go test -tags integration -race ./...` | FAIL | Current session exited 1. Token: `FAIL`. `cmd/fid` fails `TestV01ClosureJourneyAndRebuildParity`, `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity`, and `TestMaintenanceCheckReportsIncompatibleWorkflowAssignment` because their `fid init` call exits 1 with `error: invalid project initialization argument`. This exact recovery-spec/QA-plan command cannot support the reports' recorded PASS. |
| Authorized product-change and move scope | recovery commit diffs; `wc -l internal/config/config.go internal/config/project_config.go internal/bootstrap/init.go` | PARTIAL | `51763124` adds only the DEC-130 equal-length rejection in product code. `4f308c8` relocates the strict project-config implementation from `config.go` to `project_config.go`; exported declarations and call sites remain in package `config`, with no implementation change found beyond the move. The files are 230 and 295 lines, respectively, within the 400-line budget. `internal/bootstrap/init.go` remains 427 lines and is not in the recovery diff, preserving CF-028. |

## Per-Task Findings

### Task 3a - Reason-code schema and redaction evidence

#### Blocking

- None. The table-driven checks are contract-aligned and behavioral: they prove both accepted and rejected details/subjects, nullability, nested members, serialized output, and redaction for the closed registry. The report records a red focused run before the authorized equal-length repair; the recovery commit contains the test-led repair only.

#### Improvement

- None.

### Task 3b - Strict-config and init boundary evidence

#### Blocking

- [High][internal/config/config_test.go:151-158] `TestParseProjectConfigBoundaryDomains` proves accepted minimum/maximum values for both executor fields, but every malformed/lower-bound/overflow/type row changes only `max_output_bytes`. VAL-043 and task-3b require both qualifying and non-qualifying evidence for every numeric boundary; no row proves that `max_excerpt_bytes` itself rejects `0`, a negative integer, an over-`MaxInt64` literal, a float, or a quoted value. Source symmetry in `project_config.go` cannot close this evidence gap. Fix direction: add explicit rejected `max_excerpt_bytes` rows for each required invalid form while retaining the existing equal-to-output and valid-boundary rows.

- [High][project-docs/docs/phases/phase-20/tasks/task-3b-done.md:75-80] The report records the exact required integration command PASS, but independent execution from repository root now exits 1 with token `FAIL` in three `cmd/fid` integration scenarios. A passing completion-report gate needs current reproducible evidence; the failed mandated gate keeps the recovery incomplete regardless of the focused recovery-test pass. Fix direction: restore the failing integration scenarios without weakening their behavioral assertions, then rerun and record the byte-identical required command with its real terminal token.

#### Improvement

- None. Review 1's dangling-symlink collision improvement is resolved by committed behavioral evidence at `cmd/fid/init_test.go:358-380`. The pure config split is within the 400-line budget, and CF-028 remains untouched as DEC-130 requires.

## Escalations

- None. No newly appearing open ESC ID is present. This bounded review does not reopen unrelated Task 1-3 findings; the two Blocking findings above are confined to the mandated DEC-130 recovery evidence and gate.

## Final Lines

Blocking: 2  Improvement: 0  Escalations: none
VERDICT: BLOCK
