# Phase 11 Task 1 Completion Report

## Summary of what was implemented

- Added a new `internal/claims` domain package for issue-claim acquire and
  renew semantics, keeping claim business rules out of `cmd/fid`.
- Implemented journaled claim acquire with default/custom TTL, note
  persistence, live-conflict rejection, stale takeover, permitted force-acquire,
  and holder-only renewal with `activity_log` projection.
- Added an atomic concurrent-acquire proof using two simultaneous claim attempts
  against the same SQLite-backed workspace.
- Appended the implemented `issue.claim.*` rows to the in-code permissions
  matrix and added exact row-shape tests for the Phase 11 claim actions.

## Files created or modified

- Created `internal/claims/acquire.go`
- Created `internal/claims/claims.go`
- Created `internal/claims/renew.go`
- Created `internal/claims/shared.go`
- Created `internal/claims/claims_test.go`
- Modified `internal/policy/matrix.go`
- Modified `internal/policy/matrix_test.go`
- Created `project-docs/docs/phases/phase-11/tasks/task-1-done.md`

## Documentation updated

- Created `project-docs/docs/phases/phase-11/tasks/task-1-done.md`.
- `verified-no-change`: `project-docs/docs/phases/phase-11/decisions.md`
  remained correct because implementation matched DEC-026 and found no mismatch
  with `03_TRD_SCHEMA.md` section 13.
- `verified-no-change`: `project-docs/specs/03_TRD_SCHEMA.md` already matched
  the implemented claim rows and claim-table semantics; no spec edit was
  required.

## Validation Contract Coverage

- `VAL-001` PASS.
  Evidence:
  `TestMatrixContainsExactlyTheImplementedV01WriteActions` and
  `TestMatrixClaimRowsMatchSchemaSection13` prove the in-code matrix contains
  exactly one row each for `issue.claim.acquire`, `issue.claim.renew`,
  `issue.claim.force_acquire`, and `issue.claim.release`, with role booleans
  matching `03_TRD_SCHEMA.md` section 13.
- `VAL-002` PASS.
  Evidence:
  `TestAcquireCreatesLiveIssueClaimWithDefaultTTL` failed red on the missing
  matrix row and then passed green once journaled acquire was implemented.
  `TestAcquireSupportsCustomTTLAndNote` proves custom TTL override and note
  persistence.
- `VAL-003` PASS.
  Evidence:
  `TestAcquireRejectsLiveConflictBeforeAppend` failed red on the missing
  `LiveClaimConflictError`, then passed green once live conflicts rejected
  before append with holder and expiry details while leaving journal row count
  unchanged.
- `VAL-004` PASS.
  Evidence:
  `TestAcquireTakesOverStaleClaimWithoutForcePermission` failed red when stale
  takeover still attempted an `INSERT`, then passed green once stale claims were
  superseded under normal acquire permission with `claim_superseded`
  activity evidence naming the prior holder and expiry.
- `VAL-005` PASS.
  Evidence:
  `TestAcquireForceSupersedesLiveClaimForPermittedRole` and
  `TestAcquireForceRejectsDeniedRoleBeforeAppend` prove permitted force-acquire
  succeeds and denied force-acquire fails before append with
  `issue.claim.force_acquire` permission enforcement.
- `VAL-006` PASS.
  Evidence:
  `TestRenewUpdatesClaimTimePreservesOrOverridesTTLAndProjectsActivity` failed
  red on the missing renew row/API, then passed green once holder renewal
  preserved TTL when omitted, overrode TTL when supplied, advanced `claimed_at`,
  and projected a `claim_renewed` activity row.
  `TestRenewRejectsNonHolderBeforeAppend` proves non-holder renewal rejects with
  zero new journal lines.
- `VAL-013` PASS.
  Evidence:
  `TestAcquireIsAtomicAcrossConcurrentAttempts` proves two concurrent acquires
  on one issue produce exactly one success, one `LiveClaimConflictError`, and
  one current `claims` row in SQLite.

## Tests run

- `go test ./internal/claims -run TestAcquireCreatesLiveIssueClaimWithDefaultTTL`
  RED: failed with `authorize action "issue.claim.acquire": no permissions
  matrix entry`, then GREEN after the initial acquire path and matrix row were
  added.
- `go test ./internal/claims -run TestAcquireRejectsLiveConflictBeforeAppend`
  RED: failed on missing `LiveClaimConflictError`, then GREEN.
- `go test ./internal/claims -run TestAcquireTakesOverStaleClaimWithoutForcePermission`
  RED: failed with `UNIQUE constraint failed: claims.ref_type, claims.ref_id`,
  then GREEN after stale takeover switched to an update with activity evidence.
- `go test ./internal/claims -run 'TestAcquireForce(SupersedesLiveClaimForPermittedRole|RejectsDeniedRoleBeforeAppend)'`
  RED: failed on missing `Force` in `AcquireOptions`, then GREEN.
- `go test ./internal/claims -run 'TestRenew(UpdatesClaimTimePreservesOrOverridesTTLAndProjectsActivity|RejectsNonHolderBeforeAppend)'`
  RED: failed on missing renew API/helpers, then on missing renew matrix row,
  then on same-second timestamp semantics; GREEN after implementing renewal and
  timestamp advancement.
- `go test ./internal/policy -run TestMatrixContainsExactlyTheImplementedV01WriteActions`
  RED: failed with `len(MatrixRows()) = 30, want 26`, then GREEN once the claim
  rows were added to the expected implemented-action set.
- `go test ./internal/claims`
  PASS.
- `go test ./internal/policy`
  PASS.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS.
- `golangci-lint run ./...` - PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS.
  Key package coverage from command output: `cmd/fid` 77.8%,
  `internal/claims` 78.2%, `internal/issues` 81.9%,
  `internal/journal` 80.3%, `internal/policy` 81.2%,
  `internal/workflow` 95.5%.
- `go tool cover -func=coverage.out` - PASS.
  Total statement coverage: 76.4%. Confirmed key package coverage stayed above
  the Phase 11 thresholds, including `internal/claims` 78.2% and `cmd/fid`
  77.8%.
- `go test -tags integration -race ./...` - PASS.
- `go build -o /tmp/fid ./cmd/fid` - N/A; no command registration or CLI
  adapter files were touched.
- `go test -bench . -benchmem ./...` - N/A; no benchmark fixtures or journal
  scale helpers were touched.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `PENDING - reconcile after task commit`

## Commit subject

- `PENDING - reconcile after task commit`
