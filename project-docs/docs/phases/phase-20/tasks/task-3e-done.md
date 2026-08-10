# Phase 20 - Task 3e: Production initialization benchmark - Completion Report

## Summary of what was implemented

- Added `BenchmarkInitializeWorkspace`, which exercises the complete production
  initialization path in a fresh workspace for every operation.
- Used the existing test-visible init checkpoints to report stage intervals
  without adding production branches or altering durability behavior.
- Kept workspace creation, successful-result verification, and cleanup outside
  the timed region while retaining migrations, journal append, rebuild parity,
  filesystem sync, and activation inside it.
- Recorded three same-host samples: 34.063, 34.544, and 35.501 ms/op (median
  34.544 ms/op), approximately 3.84 MB and 5,243 allocations per operation.

## Files created or modified

- `internal/bootstrap/init_benchmark_test.go`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3e-done.md`

## Documentation updated

- `infra-improvement-pack.md` records Task 3e's benchmark samples, median,
  allocations, dominant intervals, and unchanged production boundary.

## Validation Contract Coverage

- No VAL-NNN applies. This is test/performance instrumentation governed by
  DEC-133. The benchmark invokes `initializeWorkspaceWithHooks` with the same
  empty hooks semantics as `InitializeWorkspace`; product source is unchanged.

## Tests run

- `rg -n '^func BenchmarkInitializeWorkspace' internal/bootstrap` - EXPECTED
  RED - `exit 1 (no output)`.
- `go test -run '^$' -bench '^BenchmarkInitializeWorkspace$' -benchmem -count=3 ./internal/bootstrap` - PASS - `ok  	github.com/mbell-nzl/fid/internal/bootstrap	3.539s`.

## Exact quality-gate results

- `go test -run '^$' -bench '^BenchmarkInitializeWorkspace$' -benchmem -count=3 ./internal/bootstrap` - PASS - `ok  	github.com/mbell-nzl/fid/internal/bootstrap	3.539s`.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`.
- `golangci-lint run ./...` - PASS - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements`.
- `go tool cover -func=coverage.out` - PASS -
  `total:										(statements)					77.7%`.
- `go test -tags integration -race ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`.

## Deviations from spec

- The first report-lint run failed with `FAIL report-lint found 4 failure(s)`.
  Root cause: four valid `go test` output tokens were wrapped onto continuation
  lines, while the validator deliberately requires each token on the command's
  result bullet. The report formatting was corrected without changing any
  command, result, test, assertion, or validator; the byte-identical lint
  command is rerun after this correction is committed.

## Concerns, technical debt, or notes for the Tech Lead

- Rebuild verification and schema initialization dominate the approximately
  34.5 ms production path. Per DEC-133, this task does not optimize them; Task
  3h avoids repeating that correct production work where initialization is not
  the test subject.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `af615cbb18abbf54130547f26e3d73c2625d73e5`

## Commit subject

- `test(bootstrap): benchmark durable initialization`
