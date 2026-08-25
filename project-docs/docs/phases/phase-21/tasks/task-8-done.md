# Phase 21 Task 8: Issue Actor And Event Convergence Completion Report

## Workflow outcome
- Task outcome: COMPLETE
- Recovery TRG-21-007 completed under task-8 scope and DEC-162. The retained V1 test collateral now asserts the accepted V2 contract without reducing assertion strength.

## Summary of what was implemented
- Converged retained Issue command tests and activity fixtures on `--actor`, Event v2 attribution, and kebab-case activity values.
- Preserved durable replay identity by persisting and restoring EventRef occurrence timestamps for create, update, and transition idempotency results.
- Derived issue aggregate sequences from the existing journal before update and transition commits so mixed retained V1/V2 histories remain appendable.
- Removed obsolete V1 policy-action helpers and resolved the staticcheck issue in actor parsing.

## Files created or modified
- `cmd/fid/activity_test.go`
- `cmd/fid/e2e_v01d_closure_test.go`
- `cmd/fid/e2e_v01e_activity_test.go`
- `cmd/fid/issue_test.go`
- `cmd/fid/issue_transition_test.go`
- `cmd/fid/issue_update_test.go`
- `cmd/fid/policy_test.go`
- `internal/activity/activity_test.go`
- `internal/issues/actor.go`
- `internal/issues/create_commit.go`
- `internal/issues/idempotency_results.go`
- `internal/issues/issue_test.go`
- `internal/issues/policy_actions.go`
- `internal/issues/task8_test.go`
- `internal/issues/transition.go`
- `internal/issues/transition_test.go`
- `internal/issues/update.go`
- `internal/issues/update_test.go`
- `internal/journal/aggregate.go`
- `project-docs/docs/phases/phase-21/tasks/task-8-done.md`

## Documentation updated
- Completion report only. No task, QA, specification, decision, or run-ledger documents were changed.

## Validation Contract Coverage
- `VAL-025`: `cmd/fid/task8_contract_test.go:9`, `:43`, `:86`, and `:132` prove `--actor` is required, Event v2 metadata is emitted, and create/update/transition idempotency replays preserve the EventRef identity. The real `fid` entrypoint replay recorded the same `EVT-20260825-000002` event ID and occurrence timestamp with `replayed: true`.
- `VAL-026`: `internal/issues/task8_test.go:16` verifies V2 create/update/transition payload and fixed same-second event ordering; `cmd/fid/issue_transition_test.go` verifies the persisted V2 attribution/provenance and allowed idempotency/revision effects; `internal/issues/idempotency_results.go` restores the stored EventRef occurrence timestamp.
- `VAL-027`: `internal/issues/task8_test.go:111` verifies source file budgets. `cmd/fid/task8_contract_test.go:178` rejects legacy `--by`; `cmd/fid/issue_test.go`, `issue_update_test.go`, and `issue_transition_test.go` prove standalone actors remain attribution only rather than role authority.

## Tests run
- Historical red: `go test -race ./internal/issues -run '^TestTask8Legacy'` - `FAIL` (terminated `exit_code: 1`); historical green after contract implementation - `ok\tgithub.com/mbell-nzl/fid/internal/issues\t1.999s` (terminated `exit_code: 0`).
- Historical red: focused Issue create command-contract test - `FAIL` (terminated `exit_code: 1`) while `--actor` was undeclared; historical green - `ok\tgithub.com/mbell-nzl/fid/cmd/fid\t1.813s` (terminated `exit_code: 0`).
- Historical red: focused idempotency replay test - `FAIL` (terminated `exit_code: 1`) while replay appended another event; historical green - `ok\tgithub.com/mbell-nzl/fid/cmd/fid\t1.876s` (terminated `exit_code: 0`).
- Recovery red: `go test -race ./internal/issues ./internal/activity` - `FAIL` (terminated `exit_code: 1`) because the retained test incorrectly required same-clock `UpdatedAt` to strictly advance; green after preserving the accepted nondecreasing timestamp contract - `ok\tgithub.com/mbell-nzl/fid/internal/journal\t(cached)` (terminated `exit_code: 0`).
- Focused recovery green: `go test -race ./cmd/fid --run '^(TestIssueTransitionCanonicalOutcomeTable|TestIssueMutationCanonicalAttributionAvoidsProtectedEffects|TestIssueUpdateCanonicalChangesExcludeOwnerAndPreserveProjection|TestV01DClosureLifecycleSmokeAndRebuildParity)$'` - `ok\tgithub.com/mbell-nzl/fid/cmd/fid\t8.753s` (terminated `exit_code: 0`).
- Real entrypoint: `go build -o /tmp/fid-task8 ./cmd/fid` - `exit 0 (no output)` (terminated `exit_code: 0`). `/tmp/fid-task8 --project /tmp/fid-task8-entrypoint.bh20iY --json init --prefix FID` and two identical `issue create --actor agent:task8 --idempotency-key task8-entrypoint-key` calls both succeeded (terminated `exit_code: 0`); the second response retained EventRef `EVT-20260825-000002` and reported `replayed: true`.
- SQLite contention diagnosis: `go test -race ./internal/journal --run '^TestCommitIdempotentSerializesConcurrentSameKeyCalls$' -count=10` - `ok\tgithub.com/mbell-nzl/fid/internal/journal\t5.205s` (terminated `exit_code: 0`); `go test -race -cover ./internal/journal --run '^TestCommitIdempotentSerializesConcurrentSameKeyCalls$' -count=20` - `ok\tgithub.com/mbell-nzl/fid/internal/journal\t9.286s\tcoverage: 30.8% of statements` (terminated `exit_code: 0`).

## Exact quality-gate results
- `go test -race ./internal/issues ./internal/activity ./internal/journal ./internal/cli ./cmd/fid/...` - PASS - `ok\tgithub.com/mbell-nzl/fid/cmd/fid\t154.622s` (terminated `exit_code: 0`).
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok\tgithub.com/mbell-nzl/fid/internal/workflow\t38.4% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` - PASS - `total:\t(statements)\t70.6%` (terminated `exit_code: 0`).

## Deviations from spec
- No accepted-contract deviations.
- Retained evidence: the pre-recovery task command failed twice with final token `FAIL` (terminated `exit_code: 1`): first due to stale V1 Role/Command, snake-case activity, and fixed-date test expectations; second due to retained standalone-role-authority expectations. Tests were adapted under DEC-162, and the final byte-identical command passed.
- Retained evidence: the historical exact lint attempt failed with token `4 issues:` (terminated `exit_code: 1`) for the actor parsing simplification and obsolete V1 helpers. The minimal in-scope cleanup produced the final `0 issues.` result.
- Retained evidence: the historical exact coverage attempt failed with token `FAIL` (terminated `exit_code: 1`) while retained V1 expectations remained. A later exact coverage attempt also failed with token `FAIL` (terminated `exit_code: 1`) because `TestCommitIdempotentSerializesConcurrentSameKeyCalls` observed `database is locked (5) (SQLITE_BUSY)`. The focused repeated race and race-plus-coverage diagnostics passed without code, test, command, or threshold changes; the later byte-identical repository coverage command passed.
- Retained entrypoint evidence: the original temporary project initialization invocation ended `error: invalid project initialization argument` (terminated `exit_code: 1`) because the generated directory name was used as an invalid implicit prefix. The real entrypoint validation supplied the explicit valid `--prefix FID` and passed.

## Concerns, technical debt, or notes for the Tech Lead
- None.

## Dependency-budget note
- No dependency changes.

## Commit hash
- Pending commit

## Commit subject
- Pending commit
