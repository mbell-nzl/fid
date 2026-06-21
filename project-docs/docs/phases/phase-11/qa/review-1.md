# Phase 11 - Mid-Phase Review 1 (Tasks 1-1)

Reviewer: QA Reviewer
Date: 2026-06-21
Tasks covered: phase-11/task-1

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-11/tasks/task-1.md`
- Completion reports:
  - `project-docs/docs/phases/phase-11/tasks/task-1-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-11/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-11/qa/qa-plan.md`
- Grounding:
  - `AGENTS.md`
  - `project-docs/docs/current/phase_brief.md`
  - `project-docs/docs/current/decisions_active.md`
  - `project-docs/docs/current/progress_current.md`
  - `project-docs/docs/current/tl_escalations_current.md`
  - `project-docs/specs/08_QUALITY_PLAN.md`
- Design sections checked for parity mode:
  - `project-docs/specs/05_DESIGN_PLAN.md` section 3.1, CLI command tree:
    claim and policy surfaces remain command-surface concerns; Task 1 stayed in
    `internal/claims` and `internal/policy`.
  - `project-docs/specs/05_DESIGN_PLAN.md` section 4, primitive command
    behavior: domain layer supports permission enforcement and one-event write
    semantics for later CLI wiring.
  - `project-docs/specs/05_DESIGN_PLAN.md` sections 5.3-5.4, output/error
    standards: Task 1 domain errors carry issue, holder, and expiry details for
    Task 2 rendering.
- Commits or changed files:
  - `2a8a3aacb6a3a845df890539ab4d9173d8623119`
  - `internal/claims/acquire.go`
  - `internal/claims/claims.go`
  - `internal/claims/renew.go`
  - `internal/claims/shared.go`
  - `internal/claims/claims_test.go`
  - `internal/policy/matrix.go`
  - `internal/policy/matrix_test.go`

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-006, VAL-013

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-001 | task-1 | ON TRACK | `internal/policy/matrix.go` rows 54-57 contain exactly the four implemented claim actions with schema-matching booleans. `internal/policy/matrix_test.go` rows 5-68 checks implemented action uniqueness; rows 70-132 checks row shape. `cmd/fid/policy.go` rows 26-27 renders `internalpolicy.MatrixRows()`, and `cmd/fid/policy_test.go` rows 16-74 verifies `fid policy show` output equals `MatrixRows()` across formats. Focused QA spot-check passed: `go test ./internal/claims ./internal/policy ./cmd/fid -run 'Test(Acquire|Renew|Matrix|PolicyShow)'`. |
| VAL-002 | task-1 | ON TRACK | `internal/claims/claims_test.go` rows 22-105 covers default TTL 1800, holder, expiry, event append, and one current DB row; rows 107-152 covers custom TTL and note read-back. Completion report lines 45-50 cites the same evidence and required full gates passed at lines 110-121. |
| VAL-003 | task-1 | ON TRACK | `internal/claims/claims_test.go` rows 154-220 covers live conflict rejection, `LiveClaimConflictError` holder/expiry details, unchanged journal line count, and unchanged one-row DB state. |
| VAL-004 | task-1 | ON TRACK | `internal/claims/claims_test.go` rows 222-331 covers stale takeover under normal acquire, claims update effect, `claim_superseded` activity effect, prior holder/expiry evidence, and new current holder. |
| VAL-005 | task-1 | ON TRACK | `internal/claims/claims_test.go` rows 333-392 covers permitted `tl` force-acquire with takeover activity evidence; rows 394-438 covers denied developer force-acquire as `PermissionDeniedError` for `issue.claim.force_acquire` with unchanged journal line count. |
| VAL-006 | task-1 | ON TRACK | `internal/claims/claims_test.go` rows 440-551 covers holder renewal preserving TTL, overriding supplied TTL, advancing `claimed_at`, preserving note, and projecting `claim_renewed` activity. Rows 553-596 covers non-holder renewal rejection before append. |
| VAL-013 | task-1 | ON TRACK | `internal/claims/claims_test.go` rows 598-667 covers two concurrent acquire attempts with exactly one success, one `LiveClaimConflictError` matching winner holder/expiry, and exactly one current DB row. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task spec required gate | `gofmt -l cmd internal` | PASS | Completion report line 110 records PASS. |
| Task spec required gate | `golangci-lint run ./...` | PASS | Completion report line 111 records PASS with `0 issues.` |
| Task spec required gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report lines 112-116 record PASS and package coverage: `cmd/fid` 77.8%, `internal/claims` 78.2%, `internal/issues` 81.9%, `internal/journal` 80.3%, `internal/policy` 81.2%, `internal/workflow` 95.5%. |
| Task spec required gate | `go tool cover -func=coverage.out` | PASS | Completion report lines 117-120 record PASS and total statement coverage 76.4%. |
| Task spec required gate | `go test -tags integration -race ./...` | PASS | Completion report line 121 records PASS. |
| Conditional gate | `go build -o /tmp/fid ./cmd/fid` | N/A | Completion report lines 122-123 records N/A because no command registration or CLI adapter files were touched. `git show --name-only` confirms touched files are `internal/claims`, `internal/policy`, and the completion report only. |
| Conditional gate | `go test -bench . -benchmem ./...` | N/A | Completion report lines 124-125 records N/A because no benchmark fixtures or journal scale helpers were touched. Touched-file list confirms no benchmark/helper files. |
| QA spot-check | `go test ./internal/claims ./internal/policy ./cmd/fid -run 'Test(Acquire|Renew|Matrix|PolicyShow)'` | PASS | QA Reviewer ran this focused deterministic check. Output: `ok github.com/mbell-nzl/fid/internal/claims`; `ok github.com/mbell-nzl/fid/internal/policy`; `ok github.com/mbell-nzl/fid/cmd/fid`. |
| Commit verification | `git rev-parse --verify 2a8a3aacb6a3a845df890539ab4d9173d8623119^{commit}` | PASS | Returned `2a8a3aacb6a3a845df890539ab4d9173d8623119`. |
| Commit subject/body | `git show --no-patch --format='%H%n%s%n%b' 2a8a3aacb6a3a845df890539ab4d9173d8623119` | PASS | Subject matches report: `feat(claims): implement issue claim domain semantics`; body includes `Refs: phase-11/task-1`. |
| Escalation baseline | `project-docs/docs/current/tl_escalations_current.md` | PASS | Open escalations remain `None.` No new Blocking finding or qualifying open ESC was found. |

## Per-Task Findings

### Task 1 - Claim Matrix Rows And Domain Claim Semantics

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-11/tasks/task-1-done.md:38` VAL-001 evidence names the matrix tests but omits the policy-show side of the QA-plan evidence, even though current code renders `internalpolicy.MatrixRows()` and existing `cmd/fid` tests verify `fid policy show` against that source. Fix direction: in the next reconciliation or final QA evidence roll-up, cite the policy-show test or include a small `fid policy show` excerpt for the four `issue.claim.*` rows so the show-vs-enforce parity proof is explicit.

## Cross-Cutting Observations

- Architecture boundary is appropriate for Task 1: claim business behavior is in `internal/claims`, policy rows are in `internal/policy`, and no `cmd/fid` adapter files were touched.
- Conditional/binary VAL evidence is credible for this window: conflict, denied force, non-holder renewal, and concurrent loser states all have negative-state tests with journal line-count or row-count checks.
- The working tree currently shows `project-docs/docs/phases/phase-11/qa/qa-plan.md` as untracked. This review treats it as a provided grounding artifact, but phase wrap should ensure the QA plan artifact is committed with the phase documentation set.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-11/qa/review-1.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```
