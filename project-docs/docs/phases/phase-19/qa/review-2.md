# Phase 19 - Mid-Phase Review 2 (Tasks 4-5)

Reviewer: QA Reviewer
Date: 2026-08-02
Tasks covered: phase-19/task-4 through phase-19/task-5

## Verdict

BLOCK

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-19/tasks/task-4.md`
  - `project-docs/docs/phases/phase-19/tasks/task-5.md`
- Completion reports:
  - `project-docs/docs/phases/phase-19/tasks/task-4-done.md`
  - `project-docs/docs/phases/phase-19/tasks/task-5-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-19/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-19/qa/qa-plan.md` §§4-6, 8
- Commits or changed files:
  - `32614bbbd50e683b61921e9ed6beabee73b46021` — `feat(workflow): add assign command and assignment-aware creation`
  - `8f79013c7e8afcd3d4bd73f9821484d9dd9d9334` — `feat(workflow): add workflow list command`
  - `cmd/fid/testdata/v01_parity/`, `internal/workflow/{assign,resolve,list}.go`, `internal/issues/create.go`, `cmd/fid/workflow.go`, and `internal/format/markdown_workflow.go`

Both reported implementation hashes resolve with `git rev-parse --verify <hash>^{commit}`. `git diff --check 50533e2..32614bb` was clean. The worktree contains a pre-existing human edit to `project-docs/docs/phases/phase-19/run-ledger.md`; it was not modified by this review.

## Summary

- Blocking findings: 2
- Improvement findings: 2
- Escalations raised: ESC-19-007
- Covered validation assertions: VAL-014..020, VAL-051..067

DEC-101's human-attestation disposition does **not** survive as a closure
control: DEC-102 already voided it after disproving its lint PASS. This review
independently re-ran all six Task 4 gates and they now pass, including lint
(`0 issues.`), so no further attested Task 4 gate is contradicted by product
execution. Task 4 nevertheless **reopens**: its completion report does not
meet required command-reporting rules and its mandated SHA-256 command fails as
written. The underlying candidate parity output passes from the only directory
where the manifest's relative paths resolve; that does not repair the required
procedure without a Tech Lead disposition (ESC-19-007).

## Validation Contract Coverage

The focused independent command below ran all named assertion tests under race
mode and terminated with `ok github.com/mbell-nzl/fid/cmd/fid 2.840s`:

```text
go test -race -count=1 -v ./internal/workflow ./internal/issues ./internal/format ./cmd/fid -run '^(TestAssignEmitsScopeAwareReplacementEffects|TestWorkflowAssignProjectsOneKindAssignment|TestWorkflowAssignRejectsWithoutMutation|TestIssueCreateUsesAssignedWorkflowAndItsInitialState|TestIssueCreateResolvesKindAssignmentBeforeProjectAssignment|TestCreateSupportsEpicTaskAndSubtaskAuthoring|TestListFiltersCountsAssignmentsAndOrdering|TestListReturnsEmptySliceForNoMatch|TestWriteWorkflowListMarkdownMatchesCommittedLiteral|TestWriteWorkflowListMarkdownRendersDocumentedJSONFieldSet|TestWorkflowListFormatsFiltersAndAssignments)$'
```

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-014 | task-4 | ON TRACK | `TestWorkflowAssignProjectsOneKindAssignment` passed; it asserts one current-schema `workflow.assigned` event and one kind row (`cmd/fid/workflow_assign_test.go:10`). |
| VAL-015 | task-4 | ON TRACK | Direct fixture first assigned `project-one`; journal trace showed its insert-only project effect, and post-rebuild SQL retained one project row. Terminal token: `workflow_id":"project-two"` after reassignment. |
| VAL-016 | task-4 | ON TRACK | `TestWorkflowAssignRejectsWithoutMutation` passed (`cmd/fid/workflow_assign_test.go:61`); direct CLI observed all six distinct rejected states, including `workflow "missing" does not exist` and `role "developer" not permitted for action "workflow.assign"`. |
| VAL-017 | task-4 | ON TRACK | `TestIssueCreateUsesAssignedWorkflowAndItsInitialState/task` and `/subtask` passed (`internal/issues/issue_test.go:171`), covering first-`unstarted` and first-declared branches. |
| VAL-051 | task-4 | ON TRACK | Direct kind/project reassignment trace used ordered delete+insert; `fid workflow list --format json` retained unassigned superseded `epic-one` and assigned `epic-two`. |
| VAL-052 | task-4 | ON TRACK | `TestAssignEmitsScopeAwareReplacementEffects/kind` passed (`internal/workflow/assign_test.go:13`); direct journal trace showed kind insert then delete+insert, and `fid rebuild --by tl` ended `rebuild complete` with one `epic-two` kind row. |
| VAL-053 | task-4 | ON TRACK | The same test's `/project` case passed; direct trace showed project insert then delete+insert with `{"scope":"project"}`, and replay left only `project-two`. |
| VAL-054 | task-4 | ON TRACK | Direct created task without a kind assignment produced `{"kind":"task","status":"queued","workflow_id":"project-two"}`. |
| VAL-055 | task-4 | ON TRACK | Direct created epic with both rows produced `{"kind":"epic","status":"queued","workflow_id":"epic-two"}`. |
| VAL-056 | task-4 | ON TRACK | `TestWorkflowAssignRejectsWithoutMutation` passed its before/after journal-byte and assignment-count checks; six direct rejection traces all exited 1 without an accepted-control write. |
| VAL-057 | task-4 | ON TRACK | `TestWorkflowAssignProjectsOneKindAssignment` passed its zero-activity assertion; direct SQL terminal token was `[{"activity_log_rows":0}]`. |
| VAL-058 | task-4 | ON TRACK | `TestCreateSupportsEpicTaskAndSubtaskAuthoring` passed (`internal/issues/issue_test.go:324`); the six independently captured v0.1 get/list outputs also match their committed default-path digests. |
| VAL-059 | task-4 | **NOT ON TRACK** | Underlying parity holds: fresh candidate capture and manifest check returned six `OK`; raw diff contained only `SHA256SUMS`, `capture.sh`, `fixture`, and `normalization.json` on the committed side. The required root command itself exited 1: `sha256sum: WARNING: 6 listed files could not be read` (ESC-19-007). |
| VAL-018 | task-5 | ON TRACK | `TestWorkflowListRendersDocumentedJSONFieldSet` passed (`cmd/fid/workflow_list_test.go:13`); companion shape assertions prove `any` null and projected counts (`:153`). |
| VAL-019 | task-5 | ON TRACK | `TestWorkflowListFormatsFiltersAndAssignments` passed and asserts qualifying `alpha-any, bravo-qa` plus excluded `delta-epic` (`cmd/fid/workflow_list_test.go:43`). |
| VAL-020 | task-5 | ON TRACK | Same test passed; `assertWorkflowListCommandShapes` checks exact `{"scope":"kind","kind":"task"}` (`cmd/fid/workflow_list_test.go:153`). |
| VAL-060 | task-5 | ON TRACK | Same test passed its JSON/TOON/Markdown/YAML decoded-structure equality loop (`cmd/fid/workflow_list_test.go:43-77`). |
| VAL-061 | task-5 | ON TRACK | `TestWriteWorkflowListMarkdownMatchesCommittedLiteral` passed (`internal/format/markdown_workflow_test.go:10`), byte-comparing the committed in-test literal. |
| VAL-062 | task-5 | ON TRACK | Same command test passed its active/unfiltered fixture ID checks, proving inactive exclusion and inclusion (`cmd/fid/workflow_list_test.go:79-85`). |
| VAL-063 | task-5 | ON TRACK | Same test passed its exact `--kind qa_review --active` intersection assertion (`cmd/fid/workflow_list_test.go:50-56`). |
| VAL-064 | task-5 | ON TRACK | `assertWorkflowListCommandShapes` passed and rejects a present project `kind` key (`cmd/fid/workflow_list_test.go:164-169`). |
| VAL-065 | task-5 | ON TRACK | The same helper passed its present, non-nil empty `assigned_to` array assertion (`cmd/fid/workflow_list_test.go:171-174`). |
| VAL-066 | task-5 | ON TRACK | `TestListFiltersCountsAssignmentsAndOrdering` passed (`internal/workflow/list_test.go:13`); the CLI test also passed its byte-identical consecutive JSON comparison (`cmd/fid/workflow_list_test.go:73-81`). |
| VAL-067 | task-5 | ON TRACK | `TestWorkflowListFormatsFiltersAndAssignments` and `TestWriteWorkflowListMarkdownRendersEmptyList` passed all four `[]\n`, exit-0 cases (`cmd/fid/workflow_list_test.go:88-95`; `internal/format/markdown_workflow_test.go:37`). |

## Gate And Evidence Checks

All results below were observed in this review session. `exit 0 (no output)`
is the literal terminal token for silent success.

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 4 gate 1 | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` |
| Task 4 gate 2 | `golangci-lint run ./...` | PASS | `0 issues.` This is the required independent control for DEC-102's former Task 4 lint failure. |
| Task 4 gate 3 | `go test -race -coverprofile=coverage.out ./...` | PASS | `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.4% of statements` |
| Task 4 gate 4 | `go tool cover -func=coverage.out` | PASS | `total: (statements) 77.0%`; required domain and CLI thresholds remain met. |
| Task 4 gate 5 | `go test -tags integration -race ./...` | PASS | `ok github.com/mbell-nzl/fid/internal/workflow 5.510s` |
| Task 4 gate 6 | `go build -o /tmp/fid ./cmd/fid` | PASS | `exit 0 (no output)` |
| Task 4 VAL-059 exact manifest command | `sha256sum -c cmd/fid/testdata/v01_parity/SHA256SUMS` | FAIL | `sha256sum: WARNING: 6 listed files could not be read`; bare manifest paths resolve against the repository root, not its artifact directory. |
| Task 4 VAL-059 underlying committed baseline | `sha256sum -c SHA256SUMS` from `cmd/fid/testdata/v01_parity/` | PASS | `issue-list-subtask.json: OK` (all six entries `OK`). |
| Task 4 VAL-059 fresh candidate | `capture.sh /tmp/fid <fresh-root>` then manifest check from `<fresh-root>` | PASS | `issue-list-subtask.json: OK` (all six entries `OK`). Raw `diff -r` exited 1 only for the four committed non-output artifacts; no JSON candidate diff. |
| Task 5 gate 1 | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` |
| Task 5 gate 2 | `golangci-lint run ./...` | PASS | `0 issues.` |
| Task 5 gate 3 | `go test -race -coverprofile=coverage.out ./...` | PASS | `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.4% of statements` |
| Task 5 gate 4 | `go tool cover -func=coverage.out` | PASS | `total: (statements) 77.0%` |
| Task 5 gate 5 | `go build -o /tmp/fid ./cmd/fid` | PASS | `exit 0 (no output)` |
| Completion-report lint, Task 4 | `project-docs/tools/report-lint.sh task-4-done.md task-4.md` | FAIL | `FAIL report-lint found 7 failure(s)`: missing `## Tests run` and six required gates not recorded `PASS`/`BLOCKED`. |
| Completion-report lint, Task 5 | `project-docs/tools/report-lint.sh task-5-done.md task-5.md` | PASS | `PASS report-lint completed with no failures` |
| Required representative command | `go test -race ./internal/issues/... ./internal/workflow/... ./internal/format/... ./cmd/fid/` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 79.674s` |
| Review artifact check | `git diff --check` | PASS | `exit 0 (no output)` |

The full coverage run reported `cmd/fid 78.1%`, `internal/issues 81.4%`,
`internal/claims 81.1%`, `internal/workflow 80.4%`, `internal/policy 86.3%`,
and `internal/journal 81.0%`; each applicable threshold holds.

## Per-Task Findings

### Task 4 - `fid workflow assign` + assignment-aware issue creation

#### Blocking

- [project-docs/docs/phases/phase-19/tasks/task-4-done.md:1] The completion report fails mechanical review: it lacks `## Tests run`, and all six required commands are `UNVERIFIED`/`FAIL` rather than `PASS` or a committed `BLOCKED` result. Review 2's green commands cannot be retroactively represented as Developer-session tokens. Fix direction: Tech Lead must disposition the reopened task and have its owning workflow produce a compliant, provenance-honest closure record.
- [project-docs/docs/phases/phase-19/tasks/task-4.md:361; project-docs/docs/phases/phase-19/qa/qa-plan.md:477-481] The required root-level `sha256sum -c` command fails despite valid artifacts, because its manifest uses bare filenames. This makes the prescribed VAL-059 proof non-reproducible as written. Fix direction: resolve ESC-19-007 with one authoritative, exact working-directory-safe procedure and reconcile the Task 4 report before resuming.

#### Improvement

- [cmd/fid/workflow_assign_test.go:61] The rejection test checks only non-empty stderr, not the required six pairwise-distinct value-naming errors. Direct Review 2 traces prove current behavior. Fix direction: add explicit expected/unique error assertions in the next scoped test change.
- [internal/workflow/assign_test.go:13] Effect-shape tests do not themselves delete/rebuild the projection and reassert the replayed row for each scope; Review 2 performed that direct CLI control. Fix direction: add permanent scope-specific rebuild assertions when Task 6 strengthens parity coverage.

### Task 5 - `fid workflow list` in every output format

#### Blocking

- None.

#### Improvement

- None.

## Escalations

- ESC-19-007 (quality-gate anomaly; cross-artifact acceptance ambiguity): the Task 4 and QA-plan SHA-256 procedure fails from the prescribed repository-root context although the artifact/candidate checks pass from the directories in which its relative filenames resolve.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-19/qa/review-2.md
Blocking: 2  Improvement: 2  Escalations: ESC-19-007
VERDICT: BLOCK
```
