# Phase 20 - Mid-Phase Review 4 (Tasks 10-12)

Reviewer: qa-reviewer
Date: 2026-08-15
Tasks covered: phase-20/task-10, phase-20/task-11, phase-20/task-12
Review type: INITIAL
Prior review: none; `qa/review-3b.md` was a bounded Task 9 re-review and did not reset the mandatory cadence.

## Verdict

REPAIR

Two Must-fix findings remain, both bounded to an existing reviewed task and requiring no product, architecture, data-model, public-interface, dependency, security, or validation-contract decision. Review 4 does not advance the window.

## Scope Reviewed

- Task specs and reports: `tasks/task-10.md`, `tasks/task-10-done.md`, `tasks/task-11.md`, `tasks/task-11-done.md`, `tasks/task-12.md`, and `tasks/task-12-done.md`.
- Phase materials: current brief/progress/decisions/escalations; Phase 20 decisions, run ledger, validation contract, QA plan, contracts relevant to the reviewed assertions, and the mid-phase review template.
- Verified implementation commits:
  - `627ca9dea9bdef451260022718b7bfcfcae9186a` — `fix(cli): preserve missing actor authority code` (`Refs: phase-20/task-10`)
  - `4155acee3037614762d0cfc0d6710f540b3f6441` — `feat(links): enforce typed link invariants` (`Refs: phase-20/task-11`)
  - `f081929cfed7a7be737010c2fb14403c61dac7da` — `feat(workflow): migrate command contract` (`Refs: phase-20/task-12`)
- Worktree state: the only pre-existing dirty file is `project-docs/docs/phases/phase-20/run-ledger.md`, the Phase Runner append described in the request. It was not modified.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: none; `ESC-20-001..010` are resolved and do not cover these bounded repairs.
- Covered validation assertions: VAL-008, VAL-010, VAL-011, VAL-017, VAL-021, VAL-022, VAL-023, VAL-024, VAL-025, VAL-026, VAL-041, VAL-044, VAL-045, VAL-046

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-008 | tasks 10-12 | ON TRACK | Task 10's Human success tests cover AC, catalog-label, and Issue-label operations; Task 11's add/remove tests require empty stderr and stdout confirmations; Task 12's define test requires stdout-only Human success. Full retained-leaf closure remains Task 18. |
| VAL-010 | tasks 10-11 | ON TRACK | Task 10 tests successive aggregate revisions and no advance on rejected AC/Issue-label mutation. Task 11's committed link Event tests require two Issue update effects; `internal/links/add_activity_test.go` covers resulting link Issue revisions. The all-family phase outcome remains Task 18. |
| VAL-011 | tasks 10-11 | ON TRACK | Task 10's AC and Issue-label tests compare revision effect tuples and Task 6's golden replay fixture is exercised. Link add/remove tests require the five-effect Event shape with Issue update effects. Fresh link JSON mutation identity evidence is still required by MF-01. |
| VAL-017 | tasks 10-12 | NOT ON TRACK | Task 10 compares mutation envelopes to the journal tail and read envelopes to `event:null`. Task 11 has JSON list/null evidence only; its add/remove tests are Human-mode and do not assert JSON Event identity. Task 12's define test compares its Event to the tail and list checks null, but assignment asserts only non-null Event. MF-01 and MF-02 own these missing mutation cells. |
| VAL-021 | task 10 | ON TRACK | `cmd/fid/issue_ac_test.go` checks exact add keys, journal-tail identity, duplicate/missing typed JSON failures and zero append; `internal/issues/ac_test.go` covers ordinal ordering and at/over-limit states. |
| VAL-022 | task 10 | ON TRACK | `cmd/fid/label_test.go` asserts closed `colour` schemas and name order; `cmd/fid/issue_label_test.go` covers add/remove schemas, specific missing/duplicate errors, and zero append; domain tests cover limits. |
| VAL-023 | task 11 | NOT ON TRACK | `cmd/fid/link_test.go:13-114` and `:249-343` validate Human add/remove and projection effects; only `link.list` has a canonical JSON test at `:446-512`. Exact canonical JSON add/remove schemas, mutation Event identities, and sorted returned revisions are untested. MF-01 owns this. |
| VAL-024 | task 11 | NOT ON TRACK | `internal/links/cycle_test.go:8-47` proves only a three-node closing path and zero append. No test covers the required self and two-node dependency cycles, acyclic diamond success, or bounded traversal of pre-existing cyclic storage. MF-01 owns this. |
| VAL-025 | task 11 | NOT ON TRACK | Human-mode self-link and missing-remove rejection are present in `cmd/fid/link_test.go:180-186` and `:345-397`, but no canonical JSON code/detail assertions prove `invalid-argument` and `link-not-found`. MF-01 owns this. |
| VAL-026 | task 12 | NOT ON TRACK | Define rejects `--require` and invalid grammar before append; assign proves kind-over-project precedence. `internal/workflow/list.go:71-113` returns an unbounded collection, contrary to contract section 3.7, and the one-item list test cannot prove deterministic multi-item order. The assignment Event identity and root `--project` plus `--project-wide` collision cells are also absent. MF-02 owns this. |
| VAL-041 | tasks 10-12 | ON TRACK | Task 10 asserts closed permitted AC/label Effects and attribution; link tests constrain the Event to link/activity/Issue-update effects and assert actor; workflow define asserts supplied actor and no requirement effects. Full representative-leaf closure remains Task 18. |
| VAL-044 | tasks 10-12 | NOT ON TRACK | Task 10's exhaustive request matrix is strong. Task 11 lacks canonical tests for omitted endpoint-type flags and the typed endpoint error matrix. Task 12 has no test combining root `--project <path>` with `--project-wide`, despite that being an explicit task criterion. MF-01 and MF-02 own the missing cells. |
| VAL-045 | tasks 10-12 | NOT ON TRACK | Task 10 has ordering and limit evidence. Link list ordering is asserted but no at/over-limit evidence exists. Workflow list is actually unbounded (`internal/workflow/list.go:71-113`) and its only test has one workflow. MF-01 and MF-02 own these omissions. |
| VAL-046 | task 11 | NOT ON TRACK | Link operations require types in `cmd/fid/link.go:31-35`, `:97-101`, and `:163-165`, but tests provide neither the ambiguous same-ID typed-endpoint matrix nor both typed missing Issue/requirement JSON codes with zero append. MF-01 owns this. |

## Gate And Evidence Checks

Each completion report contains the six commands required by its task spec and by QA plan section 5, with exact text, a command-produced token, and observed termination. All three reports mechanically pass `report-lint`; the reports' command records are therefore adequate, although their behavior-evidence claims remain subject to the findings below.

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | Task 10/11/12 each record `exit 0 (no output)` (terminated `exit_code: 0`). |
| Lint | `golangci-lint run ./...` | PASS | Each report records `0 issues.` (terminated `exit_code: 0`). |
| Race/coverage | `go test -race -coverprofile=coverage.out ./...` | PASS | T10: `ok  github.com/mbell-nzl/fid/internal/workflow  9.974s  coverage: 80.5% of statements`; T11/T12 record the same required command with `exit_code: 0`. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | T10 records `total: ... 77.6%`; T11/T12 record `total: (statements) 77.3%` (each terminated `exit_code: 0`). |
| Integration | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | T10: `ok  github.com/mbell-nzl/fid/cmd/fid/integration  5.216s`; T11: `... 5.296s`; T12: `... 4.379s` (each terminated `exit_code: 0`). |
| Build | `go build -o /tmp/fid ./cmd/fid` | PASS | Task 10/11/12 each record `exit 0 (no output)` (terminated `exit_code: 0`). |
| Mechanical report verification | `project-docs/tools/report-lint.sh <task-done.md> <task.md>` for Tasks 10, 11, and 12 | PASS | Current-session output ends `PASS report-lint completed with no failures` for every report (terminated `exit_code: 0`). |
| Commit identity | `git rev-parse --verify <hash>^{commit}` for all three implementation hashes | PASS | Each resolved to its recorded full hash (terminated `exit_code: 0`); subject and `Refs` footer match the committed objects. |
| Required independent reproduction | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | From `/home/mark/projects/fid`: `ok  github.com/mbell-nzl/fid/cmd/fid/integration  4.543s` (terminated `exit_code: 0`). |
| Focused AC/link/workflow reproduction | `go test ./cmd/fid -run '^(TestAcceptanceCriterionAddJSONMatchesContract|TestAcceptanceCriterionAddJSONUsesSpecificAtomicFailureCodes|TestIssueLabelMutationsJSONMatchContract|TestIssueLabelJSONUsesSpecificAtomicFailureCodes|TestLinkAddCreatesDeclaredAndInverseRowsInOneJournalEvent|TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend|TestLinkRemoveRejectsInvalidInputBeforeAppend|TestLinkListReturnsTypedCanonicalJSONInContractOrder|TestWorkflowDefineJSONMatchesContract|TestWorkflowDefineRejectsRemovedAndInvalidRequestsBeforeAppend|TestWorkflowAssignJSONPrecedenceAndScopedReplacement|TestWorkflowListJSONMatchesContractAndUsesNullAssignmentScope)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid  0.249s` (terminated `exit_code: 0`). This reproduction confirms the existing assertions, not the missing matrix cells. |
| Focused cycle reproduction | `go test ./internal/links -run '^TestAddRejectsThreeNodeDependencyCycle$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/internal/links  0.070s` (terminated `exit_code: 0`). Source inspection confirms it covers only the three-node state. |

## Per-Task Findings

### Task 10 - Acceptance-criterion and label contracts

#### Must-fix

- None.

#### Improvement

- None.

### Task 11 - Link contract and dependency invariants

#### Must-fix

- [cmd/fid/link_test.go:13; internal/links/cycle_test.go:8] **MF-01 — the required typed-link evidence matrix is incomplete.** The add/remove tests exercise Human output; only list is decoded as canonical JSON. The sole dependency-cycle test is three-node. The required self/two/three-node and pre-existing-cycle negative states, acyclic diamond positive state, JSON code/detail/zero-append tests for self/missing-link, omitted-type flags, both typed missing endpoints, ambiguous same-ID endpoint types, and link-list at/over-limit behavior are absent. The report's positive-only summary cannot establish VAL-023/024/025/046 or the Task 11 portions of VAL-017/044/045. Owning task: task-11. Repair boundary: `cmd/fid/link*_test.go` and `internal/links/*_test.go`, plus only the Task 11 link implementation needed if these tests expose a contract defect. Required evidence: all named positive and negative cells decode the exact envelope/code/details, assert Event/null and sorted revisions where applicable, and assert unchanged journal/projection on every rejection; rerun Task 11's unchanged six-command gate. Re-review scope: inspect the new matrix, Task 11 report/gate record, and rerun the focused link matrix plus `go test -tags integration -race ./cmd/fid/integration/...` independently.

#### Improvement

- None.

### Task 12 - Workflow command contract

#### Must-fix

- [internal/workflow/list.go:71-113] **MF-02 — `workflow.list` has no declared collection bound, and its contract evidence is incomplete.** The query has no `LIMIT 10001` or equivalent post-query rejection, so a project with 10,001 workflows returns an unbounded Response Envelope instead of typed `result-limit-exceeded`, violating contracts section 3.7 and VAL-045. `cmd/fid/workflow_list_test.go:9-55` seeds one workflow only, so it also does not prove `name,id` order. Task 12 lacks the assignment journal-tail identity and root `--project <path>` combined with `--project-wide` no-shadow matrix required by its acceptance criteria. Owning task: task-12. Repair boundary: `internal/workflow/list.go`, `internal/workflow` tests, workflow CLI error mapping if needed, and focused `cmd/fid/workflow*_test.go`; no contract or gate change. Required evidence: exactly 10,000 workflows succeed in `name,id` order, 10,001 yields the closed `result-limit-exceeded` envelope with no truncation; define/assign/list complete their Event/null assertions; and root-path plus `--project-wide`, `--require` absence, invalid scope/grammar, and precedence have qualifying/non-qualifying command-boundary evidence. Rerun Task 12's unchanged six-command gate. Re-review scope: inspect the cap/matrix implementation and report, then rerun the focused workflow matrix and `go test -tags integration -race ./cmd/fid/integration/...` independently.

#### Improvement

- None.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-01 | task-11 | VAL-017, VAL-023, VAL-024, VAL-025, VAL-044, VAL-045, VAL-046 | Complete the bounded link contract tests; change Task 11 link code only if a new assertion exposes a real defect. | Exact JSON add/remove/list envelopes and Event/null/revision checks; typed endpoint/missing/ambiguity matrix; self/two/three/pre-existing cycle rejections and diamond success; ordering/at-over limit and zero-effect evidence; exact six gates. | Fresh QA verifies the prior finding against tests/diff/report, reruns focused links and the required integration command. |
| MF-02 | task-12 | VAL-017, VAL-026, VAL-044, VAL-045 | Enforce the contract's 10,000 workflow-list cap and add the bounded workflow command matrix. | 10,000/10,001 list behavior with typed failure and no truncation; name/id ordering; define/assign/list Event/null, project-root/project-wide, grammar/absence, invalid, and precedence evidence; exact six gates. | Fresh QA verifies the cap/matrix/diff/report, reruns focused workflow tests and the required integration command. |

## Blocker Candidates

- None. Both findings are bounded to existing task scope and preserve the accepted CLI, Event, validation, and quality-gate contracts.

## Existing Escalation References

- None.

## Final Lines

Review written: project-docs/docs/phases/phase-20/qa/review-4.md
Must-fix: 2  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
