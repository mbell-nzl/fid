# Phase 21 - Task 10: Generic Mutation Claims - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added generic Claim `acquire`, `renew`, `release`, and `list` operations backed by `work_claims`, journal Event v2 commits, and idempotent replay.
- Bound mutation Claims to an active Envelope, Actor, Issue revision, worktree, mutation scope, exact version, and TTL; conflicts, expiry, and wrong-holder mutations fail closed.
- Replaced public `issue claim` and `issue unclaim` leaves with the generic Claim surface. The old spellings are not registered aliases.
- Continuation attempt 1 did not change product code: the retained real-process closure fixture cannot be converted to the accepted generic transport without also changing its out-of-scope V1 transition dependency.
- Continuation attempt 2 is report-only: TRG-21-011 confirms that the retained transition dependency belongs to the existing Task 12 recovery, so Task 10 closes on its unchanged five-command contract.

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
- `project-docs/docs/phases/phase-21/tasks/task-10-done.md`

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
- `go test -tags integration -race ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`): `TestV01ClosureJourneyAndRebuildParity` at `cmd/fid/integration/e2e_v01_closure_test.go:33` received `issue "FID-8e5336" has no live claim held by "qa_reviewer"`, while the retained V1 assertion requires `role "qa_reviewer" not permitted`. The failure occurs before the old `issue claim` call at line 40.

## Exact quality-gate results

- Continuation attempt 1 exact gate:
- `go test -race ./internal/claims ./internal/policy ./internal/workflow ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	(cached)` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 33.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)							65.3%` (terminated `exit_code: 0`)

- `go test -race ./internal/claims ./internal/policy ./internal/workflow ./internal/cli ./cmd/fid/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	156.550s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 39.7% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					66.9%` (terminated `exit_code: 0`)

## Deviations from spec

- Required command repair: the first `go test -race ./internal/claims ./internal/policy ./internal/workflow ./internal/cli ./cmd/fid/...` attempt ended `FAIL` (terminated `exit_code: 1`) because the exact registered-code inventory still listed the pre-Claim count. Updated the in-scope `internal/cli/codes_test.go` registry expectation; the byte-identical rerun ended `ok  	github.com/mbell-nzl/fid/cmd/fid	156.550s` (terminated `exit_code: 0`).
- Required command repair: the first `golangci-lint run ./...` attempt ended `3 issues:` (terminated `exit_code: 1`) because unregistered legacy Issue Claim command code and its TTL helper were dead. Removed that obsolete public-alias implementation and unused helper; the byte-identical rerun ended `0 issues.` (terminated `exit_code: 0`).
- Supplementary tagged-gate boundary: `go test -tags integration -race ./cmd/fid/integration/...` ended `FAIL` (terminated `exit_code: 1`) at the retained closure Claim assertion. This is not a Task 10 defect: the retained transition calls `claims.RequireLiveIssueClaimHeld` at `internal/issues/transition.go:221` before the removed legacy `issue claim` command, while the accepted generic Claim transport writes only `work_claims`. TRG-21-011 routes the required generic Claim-observation integration to the existing Task 12 recovery. The tagged FAIL remains supplementary evidence; it is not a PASS and Task 10 makes no product or test change for it.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `0b67c8304a84e4e5bcb8df44d79f899e40a16acb`

## Commit subject

- `feat(claims): add generic mutation claim lifecycle`
