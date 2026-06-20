# Phase 9 - Mid-Phase Review 1 (Tasks 1-3)

Reviewer: QA Reviewer
Date: 2026-06-20
Tasks covered: phase-9/task-1 through phase-9/task-3

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-9/tasks/task-1.md`
  - `project-docs/docs/phases/phase-9/tasks/task-2.md`
  - `project-docs/docs/phases/phase-9/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-9/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-9/tasks/task-2-done.md`
  - `project-docs/docs/phases/phase-9/tasks/task-3-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-9/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-9/qa/qa-plan.md`
- Current context and triage:
  - `AGENTS.md`
  - `project-docs/docs/current/decisions_active.md`
  - `project-docs/docs/current/tl_escalations_current.md`
  - `project-docs/docs/phases/phase-9/triage-log.md`
- Design sections checked for parity/polish mode:
  - `project-docs/specs/05_DESIGN_PLAN.md` lines 83-86: concern command tree.
  - `project-docs/specs/05_DESIGN_PLAN.md` lines 280-312: issue get/list concern fields.
  - `project-docs/specs/05_DESIGN_PLAN.md` lines 406-413: cross-surface parity matrix for open concerns.
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md` lines 696-776: concern raise/address/list command specs.
- Commits or changed files:
  - `f35b3febd7a9976c21be0370d3d701c6aec8e008`
  - `3ea2d4f89002dfe14620d2ddedfe2eb108fef7ee`
  - `a172634ceb47ed55ddcaf8f01d40e88ad0046129`

## Summary

- Blocking findings: 0
- Improvement findings: 3
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-008, VAL-009, VAL-010

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | Report cites `TestRaiseConcernCreatesOpenConcernAndFlipsDerivedFlags` and `TestConcernRaiseCreatesConcernAndPrintsConfirmation`; spot-check `go test -race ./internal/issues -run TestRaiseConcernCreatesOpenConcernAndFlipsDerivedFlags -count=1` passed. |
| VAL-002 | task-1 | ON TRACK | Report cites `TestRaiseConcernRejectsBeforeAppend` and unchanged journal/row counts for reject paths; required gates are recorded. |
| VAL-003 | task-1 | ON TRACK | Report cites before/after issue get/list false/0 to true/1 evidence; `05_DESIGN_PLAN.md` issue get/list fields checked. |
| VAL-004 | task-2 | ON TRACK | Report cites in-place address tests proving retained row, `status="addressed"`, and `addressed_in_completion_id`. |
| VAL-005 | task-2 | ON TRACK | Report cites unknown role, missing concern, already-addressed, missing completion, and cross-issue completion rejections before append. |
| VAL-006 | task-2 | ON TRACK | Report cites partial state true/1 and final clear false/0 on both Get and List. |
| VAL-007 | task-2 | ON TRACK | Report cites `TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites`; spot-check `go test -tags integration -race ./internal/issues -run TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites -count=1` passed. |
| VAL-008 | task-3 | ON TRACK | Report cites domain and CLI list/filter/format tests across `json`, `toon`, `markdown`, and `yaml`; design §6 row shape includes `issue_id`. |
| VAL-009 | task-3 | ON TRACK | Report cites `TestV01CConcernRebuildParity`, snapshot/delete DB/rebuild/exact parity over concern rows and derived flags. |
| VAL-010 | task-3 | ON TRACK | Report cites `TestV01CClosureSmoke`; spot-check `go test -tags integration -race ./cmd/fid -run 'TestV01C(ConcernRebuildParity|ClosureSmoke)' -count=1` passed. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 required headings | `task-1-done.md` lines 1-91 | PASS with improvement | All required `##` sections are present, including non-empty concerns, dependency-budget, commit hash, and subject sections. Top-level template title is missing; see improvement. |
| Task 1 exact command coverage | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go build -o /tmp/fid ./cmd/fid`; `go test -tags integration -race ./...` N/A | PASS | Matches Task 1 spec and QA plan. Integration N/A is explicitly reasoned because no integration-tagged test was added. |
| Task 1 DEC-007 coverage | `task-1-done.md` lines 60-65 | PASS | `cmd/fid` 78.9%; `internal/issues` 81.8%; `internal/journal` 80.3%; `internal/policy` 85.0%; `internal/workflow` 95.5%. |
| Task 1 commit verification | `git rev-parse --verify f35b3fe^{commit}`; `git log -1 --format=%s f35b3fe` | PASS | Full hash `f35b3febd7a9976c21be0370d3d701c6aec8e008`; subject `feat(concern): add concern raise command`. |
| Task 1 representative check | `go test -race ./internal/issues -run TestRaiseConcernCreatesOpenConcernAndFlipsDerivedFlags -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/issues 1.399s`. |
| Task 2 required headings | `task-2-done.md` lines 1-119 | PASS with improvement | All required `##` sections are present, including non-empty concerns, dependency-budget, commit hash, and subject sections. Top-level template title is missing; see improvement. |
| Task 2 exact command coverage | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go build -o /tmp/fid ./cmd/fid`; `go test -tags integration -race ./...` | PASS | Matches Task 2 spec and QA plan. Integration suite was required because an integration-tagged rule-seam test was added, and the report records PASS. |
| Task 2 DEC-007 coverage | `task-2-done.md` lines 85-94 | PASS | `cmd/fid` 78.2%; `internal/issues` 81.8%; `internal/journal` 80.3%; `internal/policy` 85.0%; `internal/workflow` 95.5%. |
| Task 2 commit verification | `git rev-parse --verify 3ea2d4f^{commit}`; `git log -1 --format=%s 3ea2d4f` | PASS | Full hash `3ea2d4f89002dfe14620d2ddedfe2eb108fef7ee`; subject `feat(concern): add concern address command`. |
| Task 2 representative check | `go test -tags integration -race ./internal/issues -run TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/issues 1.595s`. |
| Task 3 required headings | `task-3-done.md` lines 1-125 | PASS with improvement | All required sections are present, including non-empty concerns, dependency-budget, commit hash, and subject sections. VAL rows should carry explicit PASS labels for consistency; see improvement. |
| Task 3 exact command coverage | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid`; `go test -bench . -benchmem ./...` N/A | PASS | Matches Task 3 spec and QA plan. Benchmark N/A is explicitly reasoned because no new hot concern-read query path was added. |
| Task 3 DEC-007 coverage | `task-3-done.md` lines 88-94 | PASS | `cmd/fid` 77.8%; `internal/issues` 81.8%; `internal/journal` 80.3%; `internal/policy` 85.0%; `internal/workflow` 95.5%. |
| Task 3 commit verification | `git rev-parse --verify a172634^{commit}`; `git log -1 --format=%s a172634` | PASS | Full hash `a172634ceb47ed55ddcaf8f01d40e88ad0046129`; subject `feat: add concern list and v0.1c closure tests`. |
| Task 3 representative check | `go test -tags integration -race ./cmd/fid -run 'TestV01C(ConcernRebuildParity|ClosureSmoke)' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/cmd/fid 2.826s`. |
| DEC-A no default edge | `internal/workflow/workflow.go` lines 28-57, 68-80, 90-96; `rg RuleConcernsAddressed` | PASS | Default workflow graph contains no `concerns_addressed` edge. Rule remains at evaluation seam in `internal/issues/transition_rules.go`. |
| DEC-B concern boundary | `internal/issues/concern_raise.go`, `concern_address.go`, `concern_list.go`; `cmd/fid/concern.go` | PASS | Business behavior remains in `internal/issues`; CLI namespace wiring remains in `cmd/fid`. |
| DEC-C/DEC-D validate-before-append and in-place address | `internal/issues/concern_raise.go` lines 96-118; `internal/issues/concern_address.go` lines 88-105 and 173-190 | PASS | Raise validates issue/completion ownership before append; address rejects non-open before append and emits an update effect on the same concern row. |
| DEC-E no activity log expectation | `rg activity_log` over concern files and Phase 9 reports | PASS | No concern implementation or task evidence introduces concern `activity_log` assertions. |
| Task 3 deviation disposition | `task-3-done.md` lines 102-109; `triage-log.md` TRG-001 | PASS | Deviation matches TRG-001 RESUME: markdown renderer was required for VAL-008; transition rule test change was test-only stability. No new escalation needed. |

## Per-Task Findings

### Task 1 - `fid concern raise`

Observation: Task 1 evidence tracks toward VAL-001, VAL-002, and VAL-003. The report records required command coverage, coverage thresholds, dependency-budget note, commit hash, and subject. Design checks covered `05_DESIGN_PLAN.md` command tree lines 83-86, issue concern fields lines 280-312, cross-surface parity matrix line 412, and `06_DESIGN_PLAN_CLI.md` concern raise lines 698-720.

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-9/tasks/task-1-done.md:1` The completion report starts at `## Summary of what was implemented` and omits the template's top-level report title. Fix direction: add a `# Phase 9 - Task 1: fid concern raise - Completion Report` title if the report is touched again.

### Task 2 - `fid concern address` + `concerns_addressed` rule proof

Observation: Task 2 evidence tracks toward VAL-004, VAL-005, VAL-006, and VAL-007. The integration-tagged rule-seam proof was reported and spot-checked, and DEC-A is preserved because no default workflow edge was added. Design checks covered `05_DESIGN_PLAN.md` command tree lines 83-86, issue concern fields lines 280-312, cross-surface parity matrix line 412, and `06_DESIGN_PLAN_CLI.md` concern address lines 722-742.

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-9/tasks/task-2-done.md:1` The completion report starts at `## Summary of what was implemented` and omits the template's top-level report title. Fix direction: add a `# Phase 9 - Task 2: fid concern address + concerns_addressed rule proof - Completion Report` title if the report is touched again.

### Task 3 - `fid concern list`, rebuild parity, v0.1c closure smoke

Observation: Task 3 evidence tracks toward VAL-008, VAL-009, and VAL-010. The reported deviation is consistent with DEC-A through DEC-E and with `triage-log.md` TRG-001 RESUME; no standing decision or carry-forward constraint is contradicted. Design checks covered `05_DESIGN_PLAN.md` command tree lines 83-86, issue concern fields lines 280-312, cross-surface parity matrix line 412, and `06_DESIGN_PLAN_CLI.md` concern list lines 744-776.

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-9/tasks/task-3-done.md:34` The VAL-008 through VAL-010 entries cite evidence but do not explicitly mark each assertion `PASS`, unlike Tasks 1 and 2. Fix direction: add explicit PASS statuses to each VAL entry for consistency with the completion-report template.

## Cross-cutting observations

- Required command coverage is plausible for all three tasks: no missing, silently skipped, or narrowed required gates were found.
- DEC-007 thresholds are reported above the required floors across `cmd/fid`, `internal/issues`, `internal/journal`, `internal/policy`, and `internal/workflow` for every task.
- The Phase 9 concern work preserves the planned architecture: concern domain operations are in `internal/issues`, CLI plumbing is in `cmd/fid`, persistence remains journal-canonical with rebuild parity evidence, and concern writes do not add `activity_log` scope.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-9/qa/review-1.md
Blocking: 0  Improvement: 3  Escalations: none
VERDICT: CONTINUE
```
