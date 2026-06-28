# Phase 13 - Mid-Phase Review 1 (Tasks 1-2)

Reviewer: qa-reviewer
Date: 2026-06-27
Tasks covered: phase-13/task-1, phase-13/task-2

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-13/tasks/task-1.md`
  - `project-docs/docs/phases/phase-13/tasks/task-2.md`
- Completion reports:
  - `project-docs/docs/phases/phase-13/tasks/task-1-done.md`
  - `project-docs/docs/phases/phase-13/tasks/task-2-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-13/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-13/qa/qa-plan.md`
- Phase context:
  - `project-docs/docs/current/phase_brief.md`
  - `project-docs/docs/phases/phase-13/decisions.md`
  - `project-docs/docs/current/tl_escalations_current.md`
- Design sections checked for parity:
  - Task 1: `project-docs/specs/05_DESIGN_PLAN.md` section 3.1 CLI command tree; `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 9 label commands.
  - Task 2: `project-docs/specs/05_DESIGN_PLAN.md` section 3.1 CLI command tree; `project-docs/specs/06_DESIGN_PLAN_CLI.md` `fid issue label add` / `remove` command sections.
- Commits inspected:
  - `d5ad0ea2f5b4dbbaf6f468dd20523636d48a0d6a` `feat(labels): add label create and list commands`
  - `f62b6093d11e47073af8a1336cd34af6ae6da475` `feat(issue): add issue label commands`

## Summary

- Blocking findings: 0
- Improvement findings: 2
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-007, VAL-008

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | tasks 1-2 | ON TRACK | Task 1 report records `label.add` matrix/show evidence; Task 2 report completes `issue.label.add` and `issue.label.remove` matrix/show/enforce evidence. Focused policy/CLI checks passed during review. |
| VAL-002 | task 1 | ON TRACK | Task 1 report cites success, duplicate rejection, one event, and zero-append tests for `fid label add`. Focused check `go test ./cmd/fid -run 'TestLabel(AddCreatesOneRowInOneEventWithoutActivityLog\|ListFormatsMatchJSONReference\|AddPolicyShowIncludesMatrixRow)$' -count=1` passed. |
| VAL-003 | task 1 | ON TRACK | Task 1 report cites four-format list parity tests for json, toon, markdown, and yaml with null optional fields preserved. Same focused Task 1 check passed. |
| VAL-004 | task 1 | ON TRACK | Task 1 report cites denied and unknown role rejection before append naming `label.add`; quality gate and focused policy-show evidence are present. |
| VAL-005 | task 2 | ON TRACK | Task 2 report cites add-one-label preservation, nonexistent-label rejection, and duplicate-add rejection before append. Focused CLI and domain checks passed during review. |
| VAL-006 | task 2 | ON TRACK | Task 2 report cites remove-one-label preservation, absent-label rejection before append, and label-table preservation. Focused CLI and domain checks passed during review. |
| VAL-007 | tasks 1-2 | ON TRACK | Task 1 report covers no `activity_log` for `fid label add`; Task 2 report covers exactly one add and remove `activity_log` row with label + direction. Focused domain checks passed. |
| VAL-008 | task 2 | ON TRACK | Task 2 report cites matrix allowed roles plus denied and unknown role rejection before append naming `issue.label.add` / `issue.label.remove`. Focused policy/CLI checks passed. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 1 required gate | `gofmt -l cmd internal` | PASS | Recorded in Task 1 report lines 74 and 82. |
| Task 1 required gate | `golangci-lint run ./...` | PASS | Recorded in Task 1 report lines 75 and 83. |
| Task 1 required gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Recorded in Task 1 report lines 76 and 84. |
| Task 1 required gate | `go tool cover -func=coverage.out` | PASS | Recorded in Task 1 report lines 77 and 85-93; coverage package percentages included. |
| Task 1 QA-plan gate | `go build -o /tmp/fid ./cmd/fid` | PASS | Recorded in Task 1 report lines 78 and 94. |
| Task 1 conditional integration | `go test -tags integration -race ./...` | N/A | Recorded with allowed reason: no integration fixture or integration-tagged test touched, lines 95-97. |
| Task 2 required gate | `gofmt -l cmd internal` | PASS | Recorded in Task 2 report lines 84 and 92. |
| Task 2 required gate | `golangci-lint run ./...` | PASS | Recorded in Task 2 report lines 85 and 93. |
| Task 2 required gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Recorded in Task 2 report lines 86 and 94. |
| Task 2 required gate | `go tool cover -func=coverage.out` | PASS | Recorded in Task 2 report lines 87 and 95-99; coverage package percentages included. |
| Task 2 QA-plan gate | `go build -o /tmp/fid ./cmd/fid` | PASS | Recorded in Task 2 report lines 88 and 100. |
| Task 2 conditional integration | `go test -tags integration -race ./...` | N/A | Recorded with allowed reason: no integration fixture or integration-tagged test touched, lines 101-103. |
| Task 1 report hygiene | dependency note, concerns, commit hash/subject | PASS | Dependency note present lines 107-109; concerns section present lines 103-105; hash/subject present lines 111-117. `git rev-parse --verify` and `git show` verified the hash, subject, and `Refs: phase-13/task-1`. |
| Task 2 report hygiene | dependency note, concerns, commit hash/subject | PASS | Dependency note present lines 113-115; concerns section present lines 109-111; hash/subject present lines 117-123. `git rev-parse --verify` and `git show` verified the hash, subject, and `Refs: phase-13/task-2`. |
| Task 2 phase-brief constraints | PI-2 / PI-3 and no schema/projector change | PASS | DEC-033 matches PI-2/PI-3; `cmd/fid/issue.go` only adds `cmd.AddCommand(newIssueLabelCommand())` at line 32; `git diff --name-only d5ad0ea^ f62b609 -- migrations internal/journal/project_effects.go` returned no paths. |
| Open escalations | `project-docs/docs/current/tl_escalations_current.md` | PASS | Open escalations section says `None`; no new ESC is required from this review. |
| Representative Task 1 review check | `go test ./cmd/fid -run 'TestLabel(AddCreatesOneRowInOneEventWithoutActivityLog\|ListFormatsMatchJSONReference\|AddPolicyShowIncludesMatrixRow)$' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/cmd/fid 0.061s`. |
| Representative Task 2 review check | `go test ./cmd/fid ./internal/issues ./internal/policy -run 'TestIssueLabel(AddPreservesExistingLabelsAndProjectsActivity\|RemovePreservesRemainingLabelsAndProjectsActivity\|CommandsPolicyShowIncludesMatrixRows\|CommandsRejectBeforeAppendAndHonorPermissions)\|TestAuthorizeIssueLabelActionsMatchMatrixExpectations\|TestMatrixIssueLabelRowsMatchSchemaSection13' -count=1` | PASS | Output: `ok` for `cmd/fid` and `internal/policy`; `internal/issues` had no tests matching that regex. |
| Representative Task 2 domain check | `go test ./internal/issues -run 'Test(AddLabelAddsOneRowWithoutReplacingOthersAndProjectsActivity\|RemoveLabelRemovesOneRowPreservesLabelTableAndProjectsActivity\|AddLabelRejectsMissingOrExistingLabelBeforeAppend\|RemoveLabelRejectsAbsentLabelBeforeAppend)$' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/issues 0.070s`. |

## Per-Task Findings

### Task 1 - Label CRUD (`fid label add` / `fid label list`)

#### Blocking

- None.

#### Improvement

- `internal/labels/add.go:33` `Add` is about 98 lines, above the project target of 60 lines for non-test functions. This does not block the phase because the behavior is covered, but it increases review and maintenance cost around lock/write/event/projection flow. Fix direction: extract small helpers for validation/options and event construction while preserving the held-lock write sequence.

### Task 2 - Per-Issue Label Assignment (`fid issue label add` / `remove`)

#### Blocking

- None.

#### Improvement

- `internal/issues/label.go:42` and `internal/issues/label.go:141` both carry long near-duplicate write flows for add/remove, each over the 60-line target. This is deferrable because Task 2 evidence proves the required behavior, but it will make future label-write changes harder to audit consistently. Fix direction: factor the common held-lock/prepare/dry-run/append/project flow after Phase 13 closure or in the CF-014 cleanup batch.

## Cross-cutting Observations

- The Task 2 decision entry does not contradict PI-2 or PI-3: it keeps add/remove as single-row deltas, preserves `fid issue update --labels` as out of scope, records one `activity_log` row per successful issue-label change, and keeps the CLI implementation in `cmd/fid/issue_label.go`.
- The reviewed commits did not add files under `migrations/` and did not modify `internal/journal/project_effects.go`, so the phase brief's no-schema/no-projector constraint remains intact for Tasks 1-2.
- The existing review window leaves Task 3 unreviewed as requested. VAL-009 and VAL-010 remain outside this review.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-13/qa/review-1.md
Blocking: 0  Improvement: 2  Escalations: none
VERDICT: CONTINUE
```
