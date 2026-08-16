# Phase 20 - Task 3k: Deterministic claims time - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: gate-repair
- Condition: The current exact repeated claims/CLI race gate timed out after Go's 10-minute limit in `cmd/fid`; this report-only continuation cannot repair product tests or their harness.
- Recovery attempts: 1
- Next action: Dispatch a same-task continuation with ordinary Task 3k source/test scope to diagnose the timeout, preserve this failed attempt, and rerun the byte-identical gate.
- Human attention required: no

## Summary of what was implemented

- Added an optional operation timestamp to acquire, renew, release, and list
  options. A zero value resolves to current UTC time.
- Reused one resolved timestamp for claim liveness validation, takeover reason,
  renewal/release event time, and list stale classification.
- Added fixed just-before, exact-boundary, and just-after expiry tests plus
  focused coverage for renew, release, list, shared live-claim validation, UTC
  normalization, and the zero-value production path.
- Replaced all three claim command sleeps with deterministic temporary
  projection state while preserving ownership, stale-filter, read-only, and
  takeover assertions.
- Reclassified six isolated validation/list tests onto the existing real
  initialized workspace clone and bounded parallel scheduler so the unchanged
  ten-count required gate completes before Go's default ten-minute timeout.
  Rebuild, durability, and compatibility owners remain real-init and serial.
- This bounded continuation changes only completion-report evidence syntax and
  records the current gate attempt; it does not reopen product implementation.

## Files created or modified

- `internal/claims/clock_test.go`
- `internal/claims/claims.go`
- `internal/claims/acquire.go`
- `internal/claims/renew.go`
- `internal/claims/release.go`
- `internal/claims/list.go`
- `internal/claims/shared.go`
- `internal/claims/claims_test.go`
- `cmd/fid/claim_test.go`
- `cmd/fid/issue_claim_test.go`
- `cmd/fid/issue_test.go`
- `cmd/fid/issue_ac_test.go`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/run-ledger.md`
- `project-docs/docs/phases/phase-20/tasks/task-3k-done.md`

## Documentation updated

- `infra-improvement-pack.md` records the clock seam, eliminated sleep
  inventory, initial timeout, scoped fixture repair, and repeated-race result.
- `run-ledger.md` preserves the first exact-gate failure and byte-identical
  repaired pass in wall-clock order.
- This report preserves the historical test/deviation evidence, marks records
  without captured termination as `UNVERIFIED`, and records this continuation's
  observed results in canonical syntax.

## Validation Contract Coverage

- No VAL-NNN applies. Existing claim lifecycle and CLI assertions are
  unchanged. Fixed-time tests prove the existing strict expiry boundary.

## Tests run

- `go test -run '^TestAcquireUsesSuppliedTimeAtExpiryBoundary$' ./internal/claims` - UNVERIFIED - `FAIL\tgithub.com/mbell-nzl/fid/internal/claims [build failed]` (historical expected-red record; the original report did not retain an observed termination result).
- `go test -race -run '^(TestOperationTime|TestAcquireUsesSuppliedTimeAtExpiryBoundary|TestRenewUsesSuppliedTime|TestReleaseUsesSuppliedTime|TestListAndLiveValidation)' ./internal/claims` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/internal/claims\t2.862s` (historical focused-pass record; the original report did not retain an observed termination result).
- `go test -race -run '^(TestClaimListDispatchesFormatsAndFiltersByStoredRole|TestClaimListStaleFilterIsReadOnlyAndUsesRowIDOrder|TestIssueClaimCommandTakesOverStaleClaimWithoutForce)$' ./cmd/fid` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t3.499s` (historical focused-pass record; the original report did not retain an observed termination result).
- Focused six-test fixture/parallel repair - UNVERIFIED historical focused-pass record - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t2.771s` (the original report did not retain an observed termination result).
- Repaired one-count diagnostic - UNVERIFIED historical focused-pass record - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t57.903s` (the original report did not retain an observed termination result).
- Claim sleep inventory - UNVERIFIED historical no-match record - `exit 1 (no output)` (the original report did not retain an observed termination result).
- `/tmp/fid --help` - PASS - `Use "fid [command] --help" for more information about a command.` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race -count=10 ./internal/claims ./cmd/fid` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t578.869s` (historical repaired-pass record; the original report did not retain an observed termination result).
- `go test -race -count=10 ./internal/claims ./cmd/fid` - BLOCKED - `FAIL` (terminated `exit_code: 1`; `cmd/fid` timed out after 10m0s while `internal/claims` completed at `ok  \tgithub.com/mbell-nzl/fid/internal/claims\t93.647s`).
- `gofmt -l cmd internal` - UNVERIFIED - `exit 0 (no output)` (historical pass record; the original report did not retain an observed termination result).
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - UNVERIFIED - `0 issues.` (historical pass record; the original report did not retain an observed termination result).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.5% of statements` (historical pass record; the original report did not retain an observed termination result).
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - UNVERIFIED - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t77.7%` (historical pass record; the original report did not retain an observed termination result).
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - UNVERIFIED - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t8.440s` (historical pass record; the original report did not retain an observed termination result).
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.567s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - UNVERIFIED - `exit 0 (no output)` (historical pass record; the original report did not retain an observed termination result).
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- The first affected CLI fixture run failed with `FAIL\tgithub.com/mbell-nzl/fid/cmd/fid\t2.421s`. Aging `claimed_at` to midnight caused the SQLite `DATETIME` driver to normalize the actual precondition value to a date-only value. The fixture now uses noon on the same fixed date; no product code or assertion changed, and the byte-identical focused command passed at 3.499s.
- The original first required ten-count gate terminated with `panic: test timed out after 10m0s` and `FAIL\tgithub.com/mbell-nzl/fid/cmd/fid\t600.458s`; `internal/claims` had already passed in 70.440s. Per-test JSON timing showed 68.797s for one command-package count and identified 13.10s of validation/list work left serial by earlier file-level classification. Six isolated tests were moved to the existing real initialized clone and bounded parallel scheduler without changing assertions. The repaired one-count run took 57.903s and the historical byte-identical required command recorded `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t578.869s`; because that prior report omitted its observed termination, this continuation retains it as `UNVERIFIED` rather than PASS.
- This continuation reran the required `go test -race -count=10 ./internal/claims ./cmd/fid` command. Its first current-session attempt ended `FAIL` (terminated `exit_code: 1`) after Go reported `panic: test timed out after 10m0s`; `internal/claims` had completed at `ok  \tgithub.com/mbell-nzl/fid/internal/claims\t93.647s`. The failure is in product/test-harness scope, which this report-only continuation is forbidden to repair. No narrower or substituted rerun was used.
- The timeout repair added `cmd/fid/issue_test.go` and
  `cmd/fid/issue_ac_test.go` to the task's planned file list. Only
  validation-only rejection matrices changed harness setup; rebuild,
  durability, and compatibility owners remain unchanged.

## Concerns, technical debt, or notes for the Tech Lead

- The current repeated claims/CLI race gate again exceeded Go's ten-minute
  timeout. A same-task source/test continuation must determine whether this is
  a reproducible harness regression before Task 3k can return to COMPLETE.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `bec40c86d49932c862079aa67163003cab6a62ee`

## Commit subject

- `test(claims): make expiry time deterministic`
