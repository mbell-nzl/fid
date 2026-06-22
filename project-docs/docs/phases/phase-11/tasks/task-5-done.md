## Summary of what was implemented

- Added a lock-aware journal freshness path in `internal/journal` so claim
  writes can catch up or rebuild a stale SQLite projection while already
  holding the workspace write lock.
- Switched `internal/claims.Acquire`, `Renew`, and `Release` to use the
  held-lock freshness helper before validation, preserving the existing claim
  write flow and avoiding nested `.fid/write.lock` acquisition.
- Added stale-projection claim tests for acquire conflict, renew success, and
  release success, while keeping the concurrent acquire atomicity proof green.

## Files created or modified

- `internal/journal/replay.go`
- `internal/claims/acquire.go`
- `internal/claims/renew.go`
- `internal/claims/release.go`
- `internal/claims/claims_test.go`
- `project-docs/docs/phases/phase-11/tasks/task-5-done.md`

## Documentation updated

- Added this completion report:
  `project-docs/docs/phases/phase-11/tasks/task-5-done.md`.
- `project-docs/docs/phases/phase-11/decisions.md`: no update required because
  the implementation stayed within DEC-028.
- `project-docs/specs/02_TRD.md`, `03_TRD_SCHEMA.md`, and
  `06_DESIGN_PLAN_CLI.md`: verified no change required because the fix
  preserved the documented write path, schema contract, and public claim CLI
  behavior.

## Validation Contract Coverage

- `VAL-015` PASS: stale-projection acquire now catches up before validation and
  rejects with the existing live-claim holder/expiry error without appending;
  stale-projection renew and release catch up before validation and append
  exactly one new event while projecting the expected current-row state.
- `VAL-013` PASS: concurrent acquire atomicity remained green in focused
  `-race` runs and in the full task gate.

## Tests run

- `go test ./internal/claims -run TestAcquireRejectsLiveConflictWhenJournalIsAheadOfProjection -count=1` (`red`)
- `go test ./internal/claims -run TestAcquireRejectsLiveConflictWhenJournalIsAheadOfProjection -count=1` (`green`)
- `go test ./internal/claims -run 'TestAcquireRejectsLiveConflictWhenJournalIsAheadOfProjection|TestRenewCatchesUpJournalAheadClaimBeforeValidation|TestReleaseCatchesUpJournalAheadClaimBeforeValidation' -count=1`
- `go test ./internal/claims -count=1`
- `go test ./internal/journal -run TestEnsureProjectionFreshUsesProjectionMetaMatrix -count=1`
- `go test -race ./internal/claims -run TestAcquireIsAtomicAcrossConcurrentAttempts -count=10`

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS
  Coverage excerpts: `cmd/fid` 77.8%, `internal/issues` 81.8%,
  `internal/claims` 81.4%, `internal/workflow` 95.5%,
  `internal/policy` 81.2%, `internal/journal` 80.0%.
- `go test -tags integration -race ./...` — PASS
- `go test -bench . -benchmem ./...` — PASS
- `go build -o /tmp/fid ./cmd/fid` — PASS

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `de6a73c90c8051624ae8f69ad1d7f88268a34369` (`de6a73c`)

## Commit subject

- `fix(claims): refresh stale projections under lock`
