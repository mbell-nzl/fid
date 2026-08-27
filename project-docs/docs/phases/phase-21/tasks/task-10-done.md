# Phase 21 - Task 10: Generic Mutation Claims - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added generic Claim `acquire`, `renew`, `release`, and `list` operations backed by `work_claims`, journal Event v2 commits, and idempotent replay.
- Bound mutation Claims to an active Envelope, Actor, Issue revision, worktree, mutation scope, exact version, and TTL; conflicts, expiry, and wrong-holder mutations fail closed.
- Replaced public `issue claim` and `issue unclaim` leaves with the generic Claim surface. The old spellings are not registered aliases.

## Files created or modified

- `internal/claims/generic.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `cmd/fid/claim.go`
- `cmd/fid/issue.go`
- `cmd/fid/issue_claim.go` (removed)
- `cmd/fid/issue_helpers.go`
- `cmd/fid/generic_claim_test.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/capabilities_test.go`
- `cmd/fid/issue_test.go`
- `cmd/fid/issue_update_test.go`
- `cmd/fid/root_flags_test.go`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-10-done.md`

## Validation Contract Coverage

- VAL-032: PASS - `TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL` at `cmd/fid/generic_claim_test.go:12` proves acquire persists the active Envelope, Actor, Issue revision, worktree, scope, TTL, version, and active status; the same test proves a second active Claim for the Issue returns the typed conflict without a second mutation.
- VAL-033: PASS - `TestClaimConflictsExpiryAndWrongHolderFailWithoutMutation` at `cmd/fid/generic_claim_test.go:72` exercises duplicate active Claim conflict, expired Claim renew, and a different Actor release; each rejects without appending a mutation.
- VAL-034: PASS - `TestClaimRenewReleaseVersionAndReplayAreExact` at `cmd/fid/generic_claim_test.go:113` proves renew and release advance exact versions, stale versions fail, and matching idempotency keys replay the original events.

## Tests run

- `go test -race ./cmd/fid -run '^TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL$' -count=1` - FAIL (expected red) - `FAIL` (terminated `exit_code: 1`)
- `go test -race ./cmd/fid -run '^TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.147s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -run '^(TestClaimAcquireBindsActiveEnvelopeActorWorktreeScopeAndTTL|TestClaimConflictsExpiryAndWrongHolderFailWithoutMutation|TestClaimRenewReleaseVersionAndReplayAreExact|TestLegacyIssueClaimSpellingsAreNotAliases)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	7.701s` (terminated `exit_code: 0`)
- `go test -race ./cmd/fid -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	156.628s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-task10-D5yAq3 init --prefix FID` - PASS - `initialized fid project in .fid/` (terminated `exit_code: 0`)
- `/tmp/fid --project /tmp/fid-task10-D5yAq3 --json claim list` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"claim.list","result":{"items":[]}}` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/claims ./internal/policy ./internal/workflow ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	156.550s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 39.7% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					66.9%` (terminated `exit_code: 0`)

## Deviations from spec

- Required command repair: the first `go test -race ./internal/claims ./internal/policy ./internal/workflow ./internal/cli ./cmd/fid/...` attempt ended `FAIL` (terminated `exit_code: 1`) because the exact registered-code inventory still listed the pre-Claim count. Updated the in-scope `internal/cli/codes_test.go` registry expectation; the byte-identical rerun ended `ok  	github.com/mbell-nzl/fid/cmd/fid	156.550s` (terminated `exit_code: 0`).
- Required command repair: the first `golangci-lint run ./...` attempt ended `3 issues:` (terminated `exit_code: 1`) because unregistered legacy Issue Claim command code and its TTL helper were dead. Removed that obsolete public-alias implementation and unused helper; the byte-identical rerun ended `0 issues.` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0b67c8304a84e4e5bcb8df44d79f899e40a16acb`

## Commit subject

- `feat(claims): add generic mutation claim lifecycle`
