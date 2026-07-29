# Phase 19 - Mid-Phase Review 1 (Tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-07-30
Tasks covered: phase-19/task-1 through phase-19/task-3

## Verdict

BLOCK

## Scope Reviewed

- Task specs and committed completion reports:
  - `project-docs/docs/phases/phase-19/tasks/task-1.md`
  - `project-docs/docs/phases/phase-19/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-19/tasks/task-2.md`
  - `project-docs/docs/phases/phase-19/tasks/task-2-done.md`
  - `project-docs/docs/phases/phase-19/tasks/task-3.md`
  - `project-docs/docs/phases/phase-19/tasks/task-3-done.md`
- Validation contract and QA plan:
  - `project-docs/docs/phases/phase-19/qa/validation-contract.md`
  - `project-docs/docs/phases/phase-19/qa/qa-plan.md` §5, §6, and §8
- Decisions, ledger, current status, implementation diffs, and worktree state.
- Implementation commits: `a509d59efcc826df417f64ea72a43d743936e669`,
  `10d4d242403851b5c676a4308cde598d251516ea`, and
  `50533e2b813bf073cd75f75eb70dd09cd15e9040`.

## Summary

- Blocking findings: 1
- Improvement findings: 0
- Escalations raised: ESC-19-003
- Covered validation assertions: VAL-001..013 and VAL-027..050

All 37 in-window assertions have specific committed test/evidence citations.
All required Task 1-3 commands are present with matching terminal tokens, and
the independent representative command passed. The block is a Task 4
precondition: its parity-baseline instruction conflicts with the QA plan's
immutable implementation-commit procedure.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | `TestMigration0009WorkflowAssignmentsSchemaAndConstraints` introspects columns, checks, FK, and exact partial indexes; reviewed migration `0009` matches DEC-076. |
| VAL-002 | task-1 | ON TRACK | `TestMigrations0009And0010RoundTripWithoutResidue` snapshots schema and rows across `up → down → up`. |
| VAL-003 | task-1 | ON TRACK | `TestMigration0010SeedsEscalationDefaultExactlyAndPreservesV01Tuples` deep-compares all four seeded-table tuple sets. |
| VAL-004 | task-1 | ON TRACK | `TestLoadSurfacesEscalationRuleAsUnimplementedWithoutWriting` uses `PRAGMA query_only` and asserts the structural unimplemented rule. |
| VAL-005 | task-1 | ON TRACK | `TestDecideProjectionFreshnessRequiresPhase19SchemaVersion` pins version 5 and exercises stale-rebuild/current-fresh branches. |
| VAL-027 | task-1 | ON TRACK | Assignment test rejects duplicate kind and accepts distinct kind rows. |
| VAL-028 | task-1 | ON TRACK | Schema test accepts the first project row, rejects the second, and retains one row. |
| VAL-029 | task-1 | ON TRACK | Schema test exercises both malformed and both well-formed scope/kind shapes. |
| VAL-030 | task-1 | ON TRACK | Ordered default tuples are captured before/after `0010`; external `db_test` parity test checks frozen definitions. |
| VAL-031 | task-1 | ON TRACK | Query-only loader test plus independent exact inventory no-match: `rg` exited 1 with no output, the DEC-094 pass condition. |
| VAL-006 | task-2 | ON TRACK | `TestBuildDefinitionEffectsIsOrderedAndDeterministic` compares the ordered effect fixture and a second build. |
| VAL-007 | task-2 | ON TRACK | Identifier matrix covers invalid/duplicate workflow and state identifiers with accepted control. |
| VAL-008 | task-2 | ON TRACK | Cancelled-state policy test proves materialized edges and declared-edge deduplication. |
| VAL-032 | task-2 | ON TRACK | States/categories matrix covers empty, length mismatch, and unknown category. |
| VAL-033 | task-2 | ON TRACK | Transition test rejects undeclared states and duplicate edges with edge-specific errors. |
| VAL-034 | task-2 | ON TRACK | Rule vocabulary test proves all four implemented rules and their effect rows/params. |
| VAL-035 | task-2 | ON TRACK | Reserved-rule matrix records v0.2b/v0.2c/v0.2d messages and unknown-rule rejection. |
| VAL-036 | task-2 | ON TRACK | Params test rejects malformed JSON/empty roles and preserves valid params verbatim. |
| VAL-037 | task-2 | ON TRACK | Requirement-edge tests prove rejected unknown and accepted declared/materialized edge cases. |
| VAL-038 | task-2 | ON TRACK | Zero-cancel tests prove both non-issue acceptance/no-edge and issue-applicable rejection branches. |
| VAL-039 | task-2 | ON TRACK | Multiple-cancelled test names both states and returns no effect slice. |
| VAL-040 | task-2 | ON TRACK | `TestInitialState` proves first-`unstarted` and first-declared branches. |
| VAL-041 | task-2 | ON TRACK | Assignment validator matrix covers five rejections and compatible kind/project controls. |
| VAL-042 | task-2 | ON TRACK | Query-only/snapshot tests plus independent mutating-SQL inventory no-match (exit 1, no output). |
| VAL-009 | task-3 | ON TRACK | Built-command test asserts one current-schema `workflow.defined` event containing all effects. |
| VAL-010 | task-3 | ON TRACK | Rejection matrix byte-compares journal data with an accepted one-event control. |
| VAL-011 | task-3 | ON TRACK | Matrix test asserts exactly `tl` and `human` for both workflow actions. |
| VAL-012 | task-3 | ON TRACK | Four arrow spellings project the same transition tuples. |
| VAL-013 | task-3 | ON TRACK | Command test asserts the exact stderr confirmation and SQL-backed counts. |
| VAL-043 | task-3 | ON TRACK | Command test asserts ordered workflow, state, transition, and requirement tuples. |
| VAL-044 | task-3 | ON TRACK | Defined workflow resolves through `workflow.Load` in the built-command test. |
| VAL-045 | task-3 | ON TRACK | Successful define test compares `activity_log` count before and after. |
| VAL-046 | task-3 | ON TRACK | Rejection matrix snapshots all four workflow tables before/after each case. |
| VAL-047 | task-3 | ON TRACK | Rejection matrix asserts non-zero, value-naming, pairwise-distinct stderr. |
| VAL-048 | task-3 | ON TRACK | Policy-show test decodes JSON, TOON, Markdown, and YAML and finds both rows. |
| VAL-049 | task-3 | ON TRACK | Denied developer path preserves journal bytes; `tl` success is the paired control. |
| VAL-050 | task-3 | ON TRACK | Six accepted applicability projections, `any` SQL NULL, and invalid-value no-mutation rejection are asserted. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 exact gates | Six Task 1/QA-plan commands | PASS | Report records format `exit 0 (no output)`, lint `0 issues.`, coverage test `exit 0 (no output)`, coverage `total: (statements) 76.6%`, integration `exit 0 (no output)`, and build `exit 0 (no output)`. |
| Task 2 exact gates | Four Task 2/QA-plan commands | PASS | Report records format `exit 0 (no output)`, lint `0 issues.`, coverage test `exit 0 (no output)`, and coverage `total: (statements) 77.1%`. No required command was substituted or omitted. |
| Task 3 exact gates | Six Task 3/QA-plan commands | PASS | Report records format `exit 0 (no output)`, lint `0 issues.`, coverage test `exit 0 (no output)`, coverage `total: (statements) 77.1%`, integration `exit 0 (no output)`, and build `exit 0 (no output)`. |
| Completion-report mechanical check | `project-docs/tools/report-lint.sh <report> <task-spec>` for Tasks 1-3 | PASS | All three runs ended `PASS report-lint completed with no failures`, including committed-and-clean report checks. |
| Commit metadata | Six `git rev-parse --verify <hash>^{commit}` checks and commit-body inspection | PASS | All implementation/report hashes resolve; each has its reported Conventional subject and `Refs: phase-19/task-N`. |
| Worktree state | `git status --porcelain=v1`, report path checks, `git diff --check` | PASS | Only the pre-existing Phase Runner modification to `run-ledger.md` was present; all three reports are committed and clean; no whitespace error. |
| Migration and seed integrity | Reviewed SQL plus Task 1 tests | PASS | `0009` supplies the exact two partial unique indexes; `0010` seeds and reverses the required escalation workflow in FK-safe order. |
| DEC-096 architecture/test packaging | Source inventory and Task 2 diff | PASS | Production `internal/workflow` imports only `journal.Effect`; no production journal write/projection/replay dependency exists. The only `internal/db` workflow import is `workflow_seed_parity_test.go`, declared `package db_test`; no in-package `db` test imports workflow. |
| Independent reproduction | `go test -race ./internal/db/... ./internal/workflow/... ./internal/policy/... ./cmd/fid/...` | PASS | Exit 0. Terminating output token: `ok github.com/mbell-nzl/fid/cmd/fid 80.225s`. This is additional QA evidence, not a replacement for Developer gates. |
| VAL-059 implementation baseline | `git rev-parse --verify 50533e2b813bf073cd75f75eb70dd09cd15e9040^{commit}` | PASS | Exact output: `50533e2b813bf073cd75f75eb70dd09cd15e9040`; subject: `feat(workflow): add define command and journal event`. No `internal/issues/create.go` or `cmd/fid/testdata/v01_parity` change exists from this commit through current HEAD. |
| Task 4 baseline instruction | task-4 lines 82-85 versus QA-plan lines 417-437 | FAIL | Task spec selects ambient `HEAD`; QA plan requires the verified Task 3 implementation commit and explicitly excludes ambient `HEAD`. Current HEAD is the Task 3 report commit `697c673620b336d81e2d5370c2dd373b46f60202`, proving the selectors differ. |

## Per-Task Findings

### Task 1 - `workflow_assignments` schema + `escalation-default` seed

#### Blocking

- None.

#### Improvement

- None.

### Task 2 - Workflow authoring core (validation + effect builder)

#### Blocking

- None.

#### Improvement

- None.

### Task 3 - `fid workflow define` command + `workflow.defined` event

#### Blocking

- None.

#### Improvement

- None.

## Cross-cutting observations

- The DEC-096 repair is faithful to its authorized test-only scope. The moved
  parity coverage is external-package `db_test`; the verified production
  boundary is `internal/workflow → internal/journal.Effect` only.
- The review has copied and independently verified the exact Task 3
  implementation commit for VAL-059:
  `PHASE19_TASK3_IMPLEMENTATION_COMMIT=50533e2b813bf073cd75f75eb70dd09cd15e9040`.
  This is not the current report commit or ambient `HEAD`.

### Task 4 baseline readiness — Blocking

- [project-docs/docs/phases/phase-19/tasks/task-4.md:84] The baseline
  instruction selects `git rev-parse HEAD`, but QA-plan §6 requires the
  verified Task 3 implementation commit and explicitly rejects ambient HEAD.
  At this review HEAD is Task 3's report commit, not its implementation
  commit, so the selectors are demonstrably different. Fix direction:
  synchronize Task 4's procedure to QA-plan §6; set and verify
  `PHASE19_TASK3_IMPLEMENTATION_COMMIT=50533e2b813bf073cd75f75eb70dd09cd15e9040`
  and create the detached baseline worktree from that value before any
  `internal/issues/create.go` edit. ESC-19-003.

## Escalations

- ESC-19-003 (triggers 5 and 6): Task 4's ambient-HEAD baseline selector is a
  cross-task, material deviation from QA-plan §6's immutable Task 3
  implementation-commit procedure and can invalidate VAL-059 evidence.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-19/qa/review-1.md
Blocking: 1  Improvement: 0  Escalations: ESC-19-003
VERDICT: BLOCK
```
