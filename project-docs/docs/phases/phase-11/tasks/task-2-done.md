## Summary of what was implemented

- Wired `fid issue claim` into [cmd/fid/issue.go](/home/mark/projects/fid/cmd/fid/issue.go) with `--by`, `--ttl`, `--force`, `--renew`, and `--note`, delegating to `internal/claims` for acquire and renew flows.
- Wired `fid issue unclaim` into [cmd/fid/issue.go](/home/mark/projects/fid/cmd/fid/issue.go) and added `internal/claims.Release` for holder-only release.
- Added CLI coverage in [cmd/fid/issue_claim_test.go](/home/mark/projects/fid/cmd/fid/issue_claim_test.go) for default/custom TTL, note persistence, conflict rejection, stale takeover, force-acquire allow/deny, renew allow/deny, force+renew exclusivity, human claims, unclaim success, and unclaim rejection paths.
- Extended `internal/claims` result types so CLI confirmations can render force-acquire previous-holder context, and moved claim `EnsureProjectionFresh` calls under the claim write lock to keep concurrent acquire losers from surfacing projection-time unique-constraint errors.

## Files created or modified

- Modified [cmd/fid/issue.go](/home/mark/projects/fid/cmd/fid/issue.go)
- Added [cmd/fid/issue_claim_test.go](/home/mark/projects/fid/cmd/fid/issue_claim_test.go)
- Modified [internal/claims/acquire.go](/home/mark/projects/fid/internal/claims/acquire.go)
- Modified [internal/claims/claims.go](/home/mark/projects/fid/internal/claims/claims.go)
- Modified [internal/claims/claims_test.go](/home/mark/projects/fid/internal/claims/claims_test.go)
- Added [internal/claims/release.go](/home/mark/projects/fid/internal/claims/release.go)
- Modified [internal/claims/renew.go](/home/mark/projects/fid/internal/claims/renew.go)
- Added [project-docs/docs/phases/phase-11/tasks/task-2-done.md](/home/mark/projects/fid/project-docs/docs/phases/phase-11/tasks/task-2-done.md)

## Documentation updated

- Added [project-docs/docs/phases/phase-11/tasks/task-2-done.md](/home/mark/projects/fid/project-docs/docs/phases/phase-11/tasks/task-2-done.md).
- Verified no change needed in `project-docs/specs/06_DESIGN_PLAN_CLI.md`; implemented output matched the documented claim, renew, force-acquire, and unclaim shapes.
- No Phase 11 decision update was needed because no CLI/flag behavior deviated from the task spec or design doc.

## Validation Contract Coverage

- `VAL-002` PASS: `TestIssueClaimCommandCreatesLiveClaimWithDefaultTTL` and `TestIssueClaimCommandSupportsCustomTTLNoteAndHumanReadBack/custom ttl and note` verified default TTL `1800`, custom TTL parsing, note persistence, and single current-row readback through `issue get --section claims`.
- `VAL-003` PASS: `TestIssueClaimCommandRejectsLiveConflictBeforeAppend` verified non-holder live-claim rejection, holder/expiry error text, non-zero exit, and unchanged journal line count.
- `VAL-004` PASS: `TestIssueClaimCommandTakesOverStaleClaimWithoutForce` verified CLI stale takeover without `--force` and current-holder readback; Task 1 domain evidence in `internal/claims` continues to cover superseded-holder event/activity details.
- `VAL-005` PASS: `TestIssueClaimCommandForceAcquireReportsPreviousHolder`, `TestIssueClaimCommandRejectsDeniedForceBeforeAppend`, and `TestIssueClaimCommandMatchesPolicyShowForForceAcquire` verified permitted/denied force-acquire behavior, previous-holder confirmation text, permission denial text, and show-vs-enforce parity.
- `VAL-006` PASS: `TestIssueClaimCommandRenewsHolderClaimAndPreservesTTL`, `TestIssueClaimCommandRenewsWithOverrideTTLAndRejectsNonHolder/override ttl`, and `TestIssueClaimCommandRenewsWithOverrideTTLAndRejectsNonHolder/non-holder` verified holder-only renew, TTL preservation, TTL override, note preservation, `activity_log` projection, and zero-append non-holder rejection.
- `VAL-007` PASS: `TestIssueUnclaimCommandReleasesHolderClaim` and `TestIssueUnclaimCommandRejectsWrongHolderRepeatAndRoleErrorsBeforeAppend` verified holder release and wrong-holder, repeat, missing-role, and unknown-role rejection paths with unchanged journal evidence on rejects.
- `VAL-010` PARTIAL PASS FOR TASK-2 WINDOW: `TestIssueClaimCommandSupportsCustomTTLNoteAndHumanReadBack/human role` verified `--by human` persists and reads back as `human`. Transition blocking for non-holder agents remains Task 4 evidence.

## Tests run

- `go test ./cmd/fid -run TestIssueClaimCommandCreatesLiveClaimWithDefaultTTL -count=1`
- `go test ./cmd/fid -run TestIssueClaimCommandRenewsHolderClaimAndPreservesTTL -count=1`
- `go test ./cmd/fid -run TestIssueUnclaimCommandReleasesHolderClaim -count=1`
- `go test ./cmd/fid -run TestIssueClaimCommandForceAcquireReportsPreviousHolder -count=1`
- `go test ./cmd/fid -run 'TestIssueClaimCommand|TestIssueUnclaimCommand' -count=1`
- `go test ./internal/claims -run 'TestReleaseDeletesHolderClaim|TestReleaseRejectsNonHolderBeforeAppend' -count=1`
- `go test ./internal/claims -run TestAcquireIsAtomicAcrossConcurrentAttempts -count=100`
- `go test ./cmd/fid ./internal/claims -count=1`
- `gofmt -l cmd internal`
- `golangci-lint run ./...`
- `go test -race -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`
- `go test -tags integration -race ./...`
- `go build -o /tmp/fid ./cmd/fid`

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS; `cmd/fid` 78.1%, `internal/issues` 81.9%, `internal/claims` 83.4%, `internal/workflow` 95.5%, `internal/policy` 81.2%, `internal/journal` 80.3%
- `go test -tags integration -race ./...` — PASS
- `go build -o /tmp/fid ./cmd/fid` — PASS
- `go test -bench . -benchmem ./...` — N/A; benchmark fixtures and journal scale helpers were not touched in Task 2 scope

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `fbb8c2f`

## Commit subject

- `feat(issue): add claim and unclaim commands`
