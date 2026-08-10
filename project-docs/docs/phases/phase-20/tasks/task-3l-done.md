# Phase 20 - Task 3l: Deterministic maintenance time - Completion Report

## Summary of what was implemented

- Added an explicit fixed-time table for started-without-claim eligibility at
  just before, exactly at, and just after the 300-second grace boundary.
- Added an AST guard that rejects `time.Sleep` calls in the closure journey.
- Replaced the closure journey's four-second no-rebuild control with one fixed
  maintenance evaluation instant, derived as 24 hours after the latest
  journal event and supplied through existing `maintenance.CheckOptions.Now`.
- Required exact equality for all maintenance findings and info counters before
  and after delete-and-rebuild, alongside the existing byte-equivalent
  projection and journal checks.
- Preserved a separate real-binary maintenance command assertion, so the
  production CLI's zero-value real UTC path remains exercised.
- Reduced the dedicated integration gate from 8.440s to 4.331s (48.7%).

## Files created or modified

- `internal/maintenance/warnings_test.go`
- `cmd/fid/integration/classification_test.go`
- `cmd/fid/integration/e2e_v01_closure_test.go`
- `cmd/fid/integration/e2e_v01_closure_helpers_test.go`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3l-done.md`

## Documentation updated

- `infra-improvement-pack.md` records the deterministic evaluation instant,
  exact rebuild result, eliminated sleep, remaining sleep classifications, and
  dedicated integration timing.

## Validation Contract Coverage

- No VAL-NNN applies. Existing maintenance and closure assertions are
  unchanged or strengthened; the task adds deterministic timing evidence only.

## Tests run

- `go test -race -run '^TestCheckStartedWithoutClaimGraceBoundary$' ./internal/maintenance` - PASS - `ok  	github.com/mbell-nzl/fid/internal/maintenance	1.298s`.
- `go test -tags integration -run '^TestClosureJourneyHasNoWallClockSleep$' ./cmd/fid/integration/...` - EXPECTED RED - `FAIL	github.com/mbell-nzl/fid/cmd/fid/integration	0.562s`.
- `go test -tags integration -race -run '^(TestClosureJourneyHasNoWallClockSleep|TestV01ClosureJourneyAndRebuildParity)$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.188s`.
- Repository `time.Sleep` inventory - PASS with all matches classified - `./cmd/fid/integration/classification_test.go:60:			t.Errorf("closure journey contains time.Sleep at %s", path)`.

## Exact quality-gate results

- `go test -tags integration -race -count=10 ./internal/maintenance ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	32.305s`.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`.
- `golangci-lint run ./...` - PASS - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements`.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.7%`.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.331s`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- Remaining `time.Sleep` calls are limited to the journal lock backoff seam and
  its intentional contention/cancellation test. The integration source guard's
  error message also contains the text `time.Sleep` but is not a call.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `16978f597eb8c31a727f02670165013805efea90`

## Commit subject

- `test(maintenance): make closure timing deterministic`
