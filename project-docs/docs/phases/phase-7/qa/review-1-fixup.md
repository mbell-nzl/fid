# Phase 7 - Mid-Phase Review 1 Fixup (Task 2b)

Reviewer: qa-reviewer
Date: 2026-06-14
Tasks covered: phase-7/task-2b

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-7/tasks/task-2b.md`
- Completion reports:
  - `project-docs/docs/phases/phase-7/tasks/task-2b-done.md`
  - `project-docs/docs/phases/phase-7/tasks/task-1-done.md` committed-state proof only
- Prior review and governing docs:
  - `project-docs/docs/phases/phase-7/qa/review-1.md`
  - `project-docs/docs/phases/phase-7/qa/qa-plan.md`
  - `project-docs/docs/phases/phase-7/decisions.md`
  - `project-docs/docs/current/tl_escalations_current.md`
  - `project-docs/docs/current/decisions_active.md`
- Design parity sections checked:
  - `project-docs/specs/05_DESIGN_PLAN.md` 3.1 CLI Command Tree
  - `project-docs/specs/05_DESIGN_PLAN.md` 6.1 `fid issue get` JSON
  - `project-docs/specs/05_DESIGN_PLAN.md` 8 Cross-Surface Parity Matrix
  - `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 3 per Task 2b design reference
- Commits or changed files:
  - `69e761d8f01e8708b32e52bda861b6f79b3f4e9f` - `fix(issue): make AC verification latest deterministic`
  - `77a50def318866655c1bc6e8a93bd91dbf7957ab` - `docs(phase-7): finalize task-2b report`
  - `internal/issues/detail.go`
  - `internal/issues/transition_rules.go`
  - `internal/issues/ac_test.go`
  - `cmd/fid/issue_ac_test.go`
  - `internal/issues/ac.go`
  - `internal/issues/ac_verify.go`

## Summary

- Blocking findings: 0
- Improvement findings: 1
- Escalations raised: none
- Covered validation assertions: VAL-004, VAL-009

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-004 | task-2b | ON TRACK | `internal/issues/detail.go:156` selects latest per role by `ac_verifications.rowid DESC`, matching DEC-021 append-order semantics. `internal/issues/ac_test.go:586` inserts same-timestamp same-role fail/pass rows with IDs that would fail the old `id DESC` tie-break and expects pass. `internal/issues/ac_test.go:681` and `cmd/fid/issue_ac_test.go:295` prove fail-only latest is fail, immediate fail->pass latest is pass, and two developer rows remain. `rg -n "time\\.Sleep|verified_at DESC|id DESC|rowid DESC|latestVerificationsByRole|loadAcceptanceCriterionVerifications" internal/issues cmd/fid` found no `time.Sleep`; only non-AC ordering uses `id DESC`. Focused reviewer check passed: `go test ./internal/issues ./cmd/fid -run 'TestLatestVerificationSelectionPrefersAppendOrderWithinTimestampTie|TestVerifyAcceptanceCriterionPreservesHistoryAndLatestPerRole|TestVerifyAcceptanceCriterionLatestPerRoleSurvivesRebuild|TestIssueACVerifyCommandPreservesFailToPassHistory|TestIssueACVerifyCommandLatestPerRoleSurvivesRebuild'`. |
| VAL-009 | task-2b targeted fixup | ON TRACK | `internal/issues/ac_test.go:811` and `cmd/fid/issue_ac_test.go:344` delete the SQLite DB, rebuild from journal, and compare latest-per-role before/after for developer fail->pass plus independent `qa_lead` fail. The same focused reviewer check passed. `internal/issues/transition_rules.go:81` loads AC verification state through `issueAcceptanceCriteria`, so targeted rebuild/latest semantics are shared with the transition-rule loader after DEC-021. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task 2b required headings | Completion report sections | PASS | Report includes summary, files, docs, VAL coverage, tests, exact gates, deviations, concerns, dependency note, commit hash, and commit subject. |
| Exact task commands | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | PASS | All required commands are recorded as `PASS` in `task-2b-done.md:81`. |
| Required build command | `go build -o /tmp/fid ./cmd/fid` | PASS | Recorded as `PASS` in `task-2b-done.md:94`, with rationale that it generated deterministic fail->pass and rebuild transcripts. |
| Coverage thresholds | `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | PASS | Report records `cmd/fid` 79.2% and `internal/issues` 82.6%, above the Task 2b thresholds; total coverage is 77.2%. |
| Reviewer deterministic spot-check | `go test ./internal/issues ./cmd/fid -run 'TestLatestVerificationSelectionPrefersAppendOrderWithinTimestampTie|TestVerifyAcceptanceCriterionPreservesHistoryAndLatestPerRole|TestVerifyAcceptanceCriterionLatestPerRoleSurvivesRebuild|TestIssueACVerifyCommandPreservesFailToPassHistory|TestIssueACVerifyCommandLatestPerRoleSurvivesRebuild'` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/issues (cached)` and `ok github.com/mbell-nzl/fid/cmd/fid (cached)`. |
| No sleep/timestamp-boundary ordering proof | `rg -n "time\\.Sleep|verified_at DESC|id DESC|rowid DESC|latestVerificationsByRole|loadAcceptanceCriterionVerifications" internal/issues cmd/fid` | PASS | No `time.Sleep` remains in `internal/issues` or `cmd/fid`; AC latest query uses `rowid DESC` at `detail.go:165`. |
| Task 1 committed report proof | `git show HEAD:project-docs/docs/phases/phase-7/tasks/task-1-done.md` | PASS | Committed report contains `59f02df` and `feat(issue): add acceptance criteria command`; it does not contain `Pending git commit`. |
| Scoped file-size check | `wc -l internal/issues/ac.go internal/issues/ac_verify.go internal/issues/issue.go` | PASS | `ac.go` is 190 lines and `ac_verify.go` is 255 lines. `issue.go` is 677 lines but is explicitly out of Task 2b scope under DEC-022 / CF-012 and is not a blocker here. |
| Commit verification | `git rev-parse --verify 69e761d^{commit}`; `git rev-parse --verify 77a50de^{commit}`; `git show -s --format=%h%x09%s 69e761d 77a50de` | PASS | Hashes resolve to `69e761d8f01e8708b32e52bda861b6f79b3f4e9f` and `77a50def318866655c1bc6e8a93bd91dbf7957ab`; subjects match the requested code-behavior and report-finalization commits. |
| Open escalation state | `project-docs/docs/current/tl_escalations_current.md` | PASS | ESC-005 remains listed as decided by DEC-021 and says it closes when Task 2b lands and re-review clears VAL-004. This review clears the re-review condition; no new ESC is raised. |

## Per-Task Findings

### Task 2b - Deterministic AC verification latest-per-role

#### Blocking

- None.

#### Improvement

- `project-docs/docs/phases/phase-7/tasks/task-2b-done.md:109` records the code-behavior commit `69e761d` and subject, while the separate report-finalization commit `77a50de` is verified from git history rather than self-recorded in the completion report. This is not blocking because `77a50de` is HEAD, resolves, has subject `docs(phase-7): finalize task-2b report`, and includes `Refs: phase-7/task-2b`. Fix direction: carry this review's finalization-commit proof into final QA rather than requiring a self-referential report edit.

## Cross-cutting observations

- The Task 2b fixup is consistent with `05_DESIGN_PLAN.md` 3.1 and 6.1: `issue ac verify` and `issue get --section acceptance` retain the same command/output shape while the per-role `latest_verifications` object now has deterministic ordering semantics.
- The targeted VAL-009 check matches `05_DESIGN_PLAN.md` 8 parity expectations for AC verification state by preserving the same latest-per-role state before and after delete-DB rebuild.
- No security/privacy issue was found in this fixup scope. The change narrows ordering to a projection-local append-order key and does not add new input handling, dependency, persistence schema, or public CLI surface.

## Escalations

- None raised. Existing ESC-005 is decided by DEC-021 and has Task 2b evidence sufficient for this re-review window; closure can be handled by the normal phase escalation/triage flow.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-7/qa/review-1-fixup.md
Blocking: 0  Improvement: 1  Escalations: none
VERDICT: CONTINUE
```
