# Phase 17 - Task 0 Completion Report

## Summary of what was implemented

- Added behavioral journal tests that restore the existing COV-001 package coverage gate without changing production code, compatibility surfaces, policies, or thresholds.
- `PrepareEvent` now has coverage proving it assigns the next event ID while leaving journal bytes unchanged.
- `EnsureProjectionFreshWithLock` now has temporary-SQLite-workspace coverage proving a caller-held lock rebuilds a missing projection, restores projected data, and leaves journal bytes unchanged.
- Added a table-driven observable verdict matrix for all remaining supported journal revert strategies, including guarded inverses, retractions, link row preconditions, and explicit non-revertible guidance.

## Files created or modified

- `internal/journal/inspect_test.go` — supported event-type reversal strategy coverage.
- `internal/journal/replay_test.go` — held-lock SQLite projection rebuild coverage.
- `internal/journal/writer_test.go` — prepare-without-append coverage.
- `project-docs/docs/phases/phase-17/tasks/task-0-done.md` — this completion report.

## Documentation updated

- Added this required Task 0 completion report only. No product, policy, threshold, compatibility, or specification documentation changed.

## Validation Contract Coverage

- **VAL-024 — COV-001 entry remediation:** focused `internal/journal` coverage rose from the QA-BLK-01 baseline of 78.9% to **81.0%**, exceeding the unchanged 80.0% threshold. The unchanged full-suite coverage run independently reports `internal/journal` at **81.0%**.
- The new tests assert observable outcomes and can fail if journal preparation writes data, lock-aware projection refresh does not rebuild/preserve state, or supported reversal strategies regress.
- `git diff --exit-code` excluding the three allowed test files passed before this report was written; no production files or compatibility surfaces changed.

## Tests run

- Red baseline: `go test -race -coverprofile=/tmp/fid-journal-coverage.out ./internal/journal/...` — PASS as a test run; coverage was **78.9%**, reproducing QA-BLK-01's COV-001 failure before this task's tests.
- Focused behavioral checks: `go test -race ./internal/journal -run '^(TestPrepareEventAssignsNextIDWithoutAppendingJournalBytes|TestEnsureProjectionFreshWithLockRebuildsMissingSQLiteProjection)$' -count=1 -v` — PASS.
- Focused reversal matrix: `go test -race ./internal/journal -run '^TestRevertVerdictKeepsSpecificStrategiesForSupportedEventTypes$' -count=1 -v` — PASS (14 subcases).
- No CLI-entrypoint interaction was required: this task changes test coverage only and does not change runtime behavior.

## Exact quality-gate results

- PASS — `gofmt -l cmd internal` — no output.
- PASS — `golangci-lint run ./...` — `0 issues.`
- PASS — `go test -race -coverprofile=/tmp/fid-journal-coverage.out ./internal/journal/...` — `ok github.com/mbell-nzl/fid/internal/journal 17.639s coverage: 81.0% of statements`.
- PASS — `go tool cover -func=/tmp/fid-journal-coverage.out` — `total: (statements) 81.0%`; `RevertVerdict` and `inverseUpdateVerdict` report 100.0%, and `PrepareEvent` reports 100.0%.
- PASS — `go test -race -coverprofile=coverage.out ./...` — `ok github.com/mbell-nzl/fid/internal/journal (cached) coverage: 81.0% of statements` (all packages passed).
- PASS — `go tool cover -func=coverage.out` — `total: (statements) 75.9%`; the full-profile `internal/journal` function entries include `RevertVerdict 100.0%`, `PrepareEvent 100.0%`, and `EnsureProjectionFreshWithLock 66.7%`.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- 26df9d52043d998a631aed4288613ee9c1aa176f

## Commit subject

- test(journal): restore journal coverage gate
