# Phase 20 - Task 3k: Deterministic claims time - Completion Report

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

## Validation Contract Coverage

- No VAL-NNN applies. Existing claim lifecycle and CLI assertions are
  unchanged. Fixed-time tests prove the existing strict expiry boundary.

## Tests run

- `go test -run '^TestAcquireUsesSuppliedTimeAtExpiryBoundary$' ./internal/claims` - EXPECTED RED - `FAIL	github.com/mbell-nzl/fid/internal/claims [build failed]`.
- `go test -race -run '^(TestOperationTime|TestAcquireUsesSuppliedTimeAtExpiryBoundary|TestRenewUsesSuppliedTime|TestReleaseUsesSuppliedTime|TestListAndLiveValidation)' ./internal/claims` - PASS - `ok  	github.com/mbell-nzl/fid/internal/claims	2.862s`.
- `go test -race -run '^(TestClaimListDispatchesFormatsAndFiltersByStoredRole|TestClaimListStaleFilterIsReadOnlyAndUsesRowIDOrder|TestIssueClaimCommandTakesOverStaleClaimWithoutForce)$' ./cmd/fid` - PASS on byte-identical rerun after the fixture correction below - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.499s`.
- Focused six-test fixture/parallel repair - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.771s`.
- Repaired one-count diagnostic - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	57.903s`.
- Claim sleep inventory - PASS no-match condition - `exit 1 (no output)`.

## Exact quality-gate results

- `go test -race -count=10 ./internal/claims ./cmd/fid` - PASS on byte-identical rerun after the harness repair below - `ok  	github.com/mbell-nzl/fid/cmd/fid	578.869s`.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`.
- `golangci-lint run ./...` - PASS - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements`.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.7%`.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.440s`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`.

## Deviations from spec

- The first affected CLI fixture run failed with `FAIL	github.com/mbell-nzl/fid/cmd/fid	2.421s`. Aging `claimed_at` to midnight caused the SQLite `DATETIME` driver to normalize the actual precondition value to a date-only value. The fixture now uses noon on the same fixed date; no product code or assertion changed, and the byte-identical focused command passed at 3.499s.
- The first required ten-count gate terminated with `panic: test timed out after 10m0s` and `FAIL	github.com/mbell-nzl/fid/cmd/fid	600.458s`; `internal/claims` had already passed in 70.440s. Per-test JSON timing showed 68.797s for one command-package count and identified 13.10s of validation/list work left serial by earlier file-level classification. Six isolated tests were moved to the existing real initialized clone and bounded parallel scheduler without changing assertions. The repaired one-count run took 57.903s and the byte-identical required command passed at 578.869s.
- The timeout repair added `cmd/fid/issue_test.go` and
  `cmd/fid/issue_ac_test.go` to the task's planned file list. Only
  validation-only rejection matrices changed harness setup; rebuild,
  durability, and compatibility owners remain unchanged.

## Concerns, technical debt, or notes for the Tech Lead

- The one-off ten-count gate passed only 21.131s before Go's default timeout.
  It is valid evidence for this task but is too close to the alarm for a
  canonical recurring gate. The canonical ordinary coverage and dedicated
  integration gates do not use `-count=10`; future repeated stability commands
  should declare an explicit timeout budget when their expected runtime exceeds
  roughly eight minutes.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `PENDING`

## Commit subject

- `test(claims): make expiry time deterministic`
