# Phase 11 - Mid-Phase Review 5 (Tasks 5-5)

Reviewer: QA Reviewer
Date: 2026-06-21
Tasks covered: phase-11/task-5

## Verdict

CONTINUE

## Scope Reviewed

- Task specs:
  - `project-docs/docs/phases/phase-11/tasks/task-5.md`
- Completion reports:
  - `project-docs/docs/phases/phase-11/tasks/task-5-done.md`
- Validation contract:
  - `project-docs/docs/phases/phase-11/qa/validation-contract.md`
- QA plan:
  - `project-docs/docs/phases/phase-11/qa/qa-plan.md`
- Decisions and escalations:
  - `project-docs/docs/phases/phase-11/decisions.md` DEC-028
  - `project-docs/docs/phases/phase-11/tl-escalations.md` ESC-008
  - `project-docs/docs/current/tl_escalations_current.md`
- Commits or changed files:
  - `de6a73c90c8051624ae8f69ad1d7f88268a34369`
  - `internal/journal/replay.go`
  - `internal/claims/acquire.go`
  - `internal/claims/renew.go`
  - `internal/claims/release.go`
  - `internal/claims/claims_test.go`

## Summary

- Blocking findings: 0
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-013, VAL-015

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-013 | task-5 | ON TRACK | `internal/claims/claims_test.go:1037` keeps the concurrent acquire proof: exactly one success, one loser conflict with holder/expiry, and one current claim row. Completion report records `go test -race ./internal/claims -run TestAcquireIsAtomicAcrossConcurrentAttempts -count=10` plus the full `go test -race -coverprofile=coverage.out ./...` gate as PASS. Reviewer spot-check `go test ./internal/claims -run 'TestAcquireRejectsLiveConflictWhenJournalIsAheadOfProjection\|TestRenewCatchesUpJournalAheadClaimBeforeValidation\|TestReleaseCatchesUpJournalAheadClaimBeforeValidation\|TestAcquireIsAtomicAcrossConcurrentAttempts' -count=1` passed. |
| VAL-015 | task-5 | ON TRACK | Claim writes now call `journal.EnsureProjectionFreshWithLock` while already inside `journal.WithLock` in `internal/claims/acquire.go:35`, `internal/claims/renew.go:35`, and `internal/claims/release.go:35`; the helper avoids recursive lock acquisition by dispatching catch-up/rebuild through held-lock paths in `internal/journal/replay.go:94`. Stale-projection tests cover acquire conflict with zero append and caught-up holder (`internal/claims/claims_test.go:343`), renew success with one appended renew event and one current row (`internal/claims/claims_test.go:737`), and release success with one appended release event and no current row (`internal/claims/claims_test.go:982`). Fresh-behavior regressions remain covered by existing acquire, renew, release, stale takeover, force, and rejection tests in the same file and by the full reported gate. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Task spec required command | `gofmt -l cmd internal` | PASS | Completion report records PASS at `project-docs/docs/phases/phase-11/tasks/task-5-done.md:52`; reviewer spot-check produced no output. |
| Task spec required command | `golangci-lint run ./...` | PASS | Completion report records PASS at `project-docs/docs/phases/phase-11/tasks/task-5-done.md:53`. |
| Task spec required command | `go test -race -coverprofile=coverage.out ./...` | PASS | Completion report records PASS at `project-docs/docs/phases/phase-11/tasks/task-5-done.md:54`. |
| Task spec required command | `go tool cover -func=coverage.out` | PASS | Completion report records PASS and the required package excerpts at `project-docs/docs/phases/phase-11/tasks/task-5-done.md:55`. Reported package coverage meets COV-001/COV-002: `cmd/fid` 77.8%, `internal/issues` 81.8%, `internal/claims` 81.4%, `internal/workflow` 95.5%, `internal/policy` 81.2%, `internal/journal` 80.0%. |
| Task spec required command | `go test -tags integration -race ./...` | PASS | Completion report records PASS at `project-docs/docs/phases/phase-11/tasks/task-5-done.md:59`. |
| Task spec required command | `go test -bench . -benchmem ./...` | PASS | Completion report records PASS at `project-docs/docs/phases/phase-11/tasks/task-5-done.md:60`; this matches the Task 5 and QA-plan requirement. |
| Task spec required command | `go build -o /tmp/fid ./cmd/fid` | PASS | Completion report records PASS at `project-docs/docs/phases/phase-11/tasks/task-5-done.md:61`; this matches the Task 5 and QA-plan requirement. |
| Focused reviewer check | `go test ./internal/claims -run 'TestAcquireRejectsLiveConflictWhenJournalIsAheadOfProjection\|TestRenewCatchesUpJournalAheadClaimBeforeValidation\|TestReleaseCatchesUpJournalAheadClaimBeforeValidation\|TestAcquireIsAtomicAcrossConcurrentAttempts' -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/internal/claims 0.049s`. |
| Commit verification | `git rev-parse --verify de6a73c90c8051624ae8f69ad1d7f88268a34369^{commit}` and `git show --no-patch --format=%H%n%s%n%b ...` | PASS | Hash verified. Subject is `fix(claims): refresh stale projections under lock`; body includes `Refs: phase-11/task-5`, matching the completion report. |
| ESC-008 / DEC-028 alignment | Archived ESC-008, DEC-028, Task 5 implementation and evidence | PASS | ESC-008 is archived as resolved to do-now Task 5. DEC-028 requires fresh claim validation without nested `.fid/write.lock` acquisition and no public claim CLI/semantic change. The narrow helper plus acquire/renew/release call sites align with that decision; current escalations list has no open items, so no new ESC was raised. |
| Public/scope drift check | `git show --name-only --format= de6a73c90c8051624ae8f69ad1d7f88268a34369` | PASS | Changed files are limited to `internal/journal/replay.go`, claim write call sites/tests, and the completion report. No schema migration or CLI output contract file changed. |

## Per-Task Findings

### Task 5 - Claim Write Projection Freshness Lock Ordering

#### Blocking

- None.

#### Improvement

- None.

## Cross-Cutting Observations

- ESC-008 is treated as resolved, not carried forward: the archived escalation requires Task 5 plus Review 5 verification, and this review found VAL-015 evidence in the stale-projection acquire, renew, and release fixtures.
- The completion report's required gate list exactly matches both Task 5's `Exact Commands To Run` and the QA plan's Task 5 required commands, including benchmark and build gates.

## Escalations

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-11/qa/review-5.md
Blocking: 0  Improvement: 0  Escalations: none
VERDICT: CONTINUE
```
