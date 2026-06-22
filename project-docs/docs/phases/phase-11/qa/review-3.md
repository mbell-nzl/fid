# Phase 11 - Mid-Phase Review 3 (Tasks 3-3)

Reviewer: QA Reviewer
Date: 2026-06-21
Tasks covered: phase-11/task-3

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-11/tasks/task-3.md`
- Completion reports:
  - `project-docs/docs/phases/phase-11/tasks/task-3-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-11/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-11/qa/qa-plan.md`
- Grounding/current docs:
  - `AGENTS.md`
  - `project-docs/docs/current/phase_brief.md`
  - `project-docs/docs/current/decisions_active.md`
  - `project-docs/docs/current/progress_current.md`
  - `project-docs/docs/current/tl_escalations_current.md`
  - `project-docs/specs/08_QUALITY_PLAN.md`
- Design parity sections checked:
  - `project-docs/specs/05_DESIGN_PLAN.md` §8 Cross-Surface Parity Matrix, especially the claim-state row.
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md` `fid claim list`.
- Commits or changed files:
  - `47c0f8622a68913d3a918278733bc4270cb120f7` - `feat(claims): add claim list command`
  - Changed files from commit: `cmd/fid/claim.go`, `cmd/fid/claim_test.go`, `cmd/fid/root.go`, `internal/claims/list.go`, `internal/claims/list_test.go`, `internal/claims/shared.go`, `internal/format/markdown.go`, `project-docs/docs/phases/phase-11/tasks/task-3-done.md`.

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-008, VAL-009, VAL-010

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-008 | task-3 | ON TRACK | Completion report lines 26 and tests in `cmd/fid/claim_test.go` cover default JSON fields, `--by` filtering, unknown-role empty filtering, and json/toon/markdown/yaml parity. `internal/claims/list.go` exposes the required fields (`ref_type`, `ref_id`, `claimed_by`, `claimed_at`, `ttl_seconds`, `expires_at`, `is_stale`, `note`) and filters on stored `claimed_by`. Representative review check passed: `go test ./cmd/fid -run 'TestClaimListDefaultsToJSONAndListsLiveClaims|TestClaimListDispatchesFormatsAndFiltersByStoredRole|TestClaimListStaleFilterIsReadOnlyAndUsesRowIDOrder|TestClaimListOrderIsStableAcrossRebuild' -count=1` returned `ok github.com/mbell-nzl/fid/cmd/fid 2.394s`. |
| VAL-009 | task-3 | ON TRACK | Completion report line 27 cites stale membership, live exclusion, `ORDER BY rowid`, and rebuild-stable ordering tests. `internal/claims/list.go` orders by `rowid ASC`; `cmd/fid/claim_test.go` verifies stale-only output excludes the live human claim, preserves stale row order, and confirms the same list order after `fid rebuild`. `internal/claims/list_test.go` also proves live and stale filtering in the domain package. |
| VAL-010 | task-3 | ON TRACK FOR TASK-3 SLICE | Completion report line 28 correctly limits Task 3's VAL-010 claim to read-surface evidence. `cmd/fid/claim_test.go` and `internal/claims/list_test.go` prove human-held claims list as `claimed_by: human`; the non-holder transition-blocking portion remains assigned to Task 4 by the validation contract and QA plan. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task spec required gate | `gofmt -l cmd internal` | PASS | Required by task spec and QA plan; recorded in completion report line 41. |
| Task spec required gate | `golangci-lint run ./...` | PASS | Required by task spec and QA plan; recorded in completion report line 42 with `0 issues.` |
| Task spec required gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Required by task spec and QA plan; recorded in completion report line 43. |
| Task spec required gate | `go tool cover -func=coverage.out` | PASS | Required by task spec and QA plan; recorded in completion report line 44. Coverage snapshot on line 45 satisfies COV-001/COV-002 thresholds for listed domain packages and `cmd/fid`. |
| Task spec required gate | `go test -tags integration -race ./...` | PASS | Required by task spec and QA plan; recorded in completion report line 46. |
| Task spec required gate | `go build -o /tmp/fid ./cmd/fid` | PASS | Required by task spec and QA plan for Task 3; recorded in completion report line 49. |
| Conditional gate | `go test -bench . -benchmem ./...` | PASS | Task spec requires this only if benchmark fixtures or list performance paths are touched. Completion report lines 47-48 records it as run because Task 3 added a list query/read path. |
| Focused tests | `go test ./cmd/fid -run 'TestClaimListDefaultsToJSONAndListsLiveClaims|TestClaimListDispatchesFormatsAndFiltersByStoredRole|TestClaimListStaleFilterIsReadOnlyAndUsesRowIDOrder|TestClaimListOrderIsStableAcrossRebuild' -count=1` | PASS | Review-run deterministic spot check returned `ok github.com/mbell-nzl/fid/cmd/fid 2.394s`. |
| Commit hash | `git rev-parse --verify 47c0f86^{commit}` | PASS | Resolved to `47c0f8622a68913d3a918278733bc4270cb120f7`. |
| Commit subject | `git show -s --format=%H%n%s%n%b 47c0f86` | PASS | Subject is `feat(claims): add claim list command`; body includes `Refs: phase-11/task-3`. |
| Escalation baseline | `project-docs/docs/current/tl_escalations_current.md` | PASS | Re-read during review; open escalation baseline remains `None`. No new open ESC warranted. |

## Per-Task Findings

### Task 3 - `fid claim list`

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-11/tasks/task-3-done.md:41` records each close gate as `PASS`, but most full-gate entries do not include literal output excerpts or artifact pointers beyond the coverage snapshot. This is acceptable for this mid-phase review because every required command is present and the review spot check passed, but it leaves final QA with less reproducible evidence than the QA plan prefers. Fix direction: in future completion reports, include concise stdout/stderr excerpts or log/artifact paths for each full close gate.

## Cross-cutting observations

- Architecture and module boundaries are consistent with the task spec: claim query and stale classification live in `internal/claims`, command registration and flags live in `cmd/fid`, and markdown rendering support is kept in `internal/format`.
- Design parity was checked against `05_DESIGN_PLAN.md` §8 and `06_DESIGN_PLAN_CLI.md` `fid claim list`; no discrepancy found for the Task 3 read surface.
- Read-only behavior is covered by CLI tests comparing journal line counts and `activity_log` rows around `fid claim list --stale`.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-11/qa/review-3.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```
