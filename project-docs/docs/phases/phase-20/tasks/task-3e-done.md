# Phase 20 - Task 3e: Production initialization benchmark - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: The MF-01 report-only repair records fresh terminating evidence for every exact Task 3e command, and the repaired report passes the canonical checker.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

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
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-3e-done.md project-docs/docs/phases/phase-20/tasks/task-3e.md` - FAIL - `FAIL report-lint found 25 failure(s)` (terminated `exit_code: 25`)
- `go test -run '^$' -bench '^BenchmarkInitializeWorkspace$' -benchmem -count=3 ./internal/bootstrap` - PASS - `ok  	github.com/mbell-nzl/fid/internal/bootstrap	3.588s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -run '^$' -bench '^BenchmarkInitializeWorkspace$' -benchmem -count=3 ./internal/bootstrap` - PASS - `ok  	github.com/mbell-nzl/fid/internal/bootstrap	3.588s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- The first report-lint run failed with `FAIL report-lint found 4 failure(s)`.
  Root cause: four valid `go test` output tokens were wrapped onto continuation
  lines, while the validator deliberately requires each token on the command's
  result bullet. The report formatting was corrected without changing any
  command, result, test, assertion, or validator; the byte-identical lint
  command is rerun after this correction is committed.
- The MF-01 pre-repair report-lint run failed with `FAIL report-lint found 25
  failure(s)` (terminated `exit_code: 25`). Root cause: this historical report
  lacked the canonical typed workflow outcome and current evidence-row syntax.
  This report-only continuation adds current-session terminating evidence and
  preserves the prior expected-red evidence and deviations without changing any
  task command, product behavior, test, assertion, or validator.

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
