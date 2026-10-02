# ESC-23-005 test-parallelism repair evidence

Authority: Human delegation recorded in DEC-258, 2026-09-27. Working directory
for every command: `/home/mark/projects/fid`. Repair commit:
`3dff4574a71438900100449e36d7044969ba9162`
(`perf(test): parallelize serial test suites for ESC-23-005`). No command,
budget, assertion, dependency, Event, schema, or product source was changed.
Earlier failed attempts stay in `tasks/task-11-done.md` and TRG-23-005.

## Diagnosis

The pre-repair measurements were taken at `c1588fc` with
`go test -race -cover -count=1 -timeout 30m -json ./cmd/fid`, with the package
run alone:

- `cmd/fid` took 535.092 s wall (`ok ... 535.092s coverage: 69.3%`,
  terminated `exit_code: 0`). `/usr/bin/time` reported 172% CPU, so about 1.7
  of 6 cores were used.
- Per-test JSON: 103 of 305 top-level tests had no `t.Parallel()`. Go runs
  those one after another before any parallel test resumes, and they ended at
  370.6 s of the 531.7 s test wall, on about one core. The largest was
  `TestStopPredicatesAcrossIssueAndEnvelopeCommands` at 79.3 s: 10 serial
  subtests, each building its own workspace.
- All 85 tagged-suite top-level tests were serial.
- Recorded same-code spread under `./...`: 468.708/564.337 s (task 10) and
  593.812/614.601 s (task 11's parent). Task 11's attempts were 686.386 s and
  630.927 s. The margin was under 10%; the spread was 15-20%.

## Change

- Existing top-level tests call `t.Parallel()`: 102 in `cmd/fid`, 84 in
  `cmd/fid/integration`, and 337 across `internal/*`. Subtests are parallel in
  the eight longest `cmd/fid` matrices. Every subtest builds its own workspace
  (`executionTransitionFixture`, `task18TerminalGuardFixture`,
  `gateReadyEvidenceEnvelopeForTest`, `runEnvelopeForTest`,
  `stopOperationEnvelopeForTest`).
- Audit before the change. The following were searched for in tests and
  production: `Setenv`/`Chdir`/`Umask`/stdio reassignment, package-level
  mutable variables, fixed temp paths, and elapsed-time assertions.
  - `rootFlagCapabilities` is filled only in `init()`.
  - The only swappable hook is `journal.hashSuffixGenerator`.
- Tests left serial, with a `//nolint:paralleltest` reason or automatic
  `paralleltest` exemption:
  - `t.Setenv` users: `TestLiveOutputGoldenFixtures` (through its helper),
    `TestExecuteAppliesAcceptedGatePolicy`,
    `TestRunLiveFramesInterleavedOutputControlsAndFailure`, and
    `TestObserveRepositoryRejectsChangeDuringCapture`.
  - The hook swapper `TestGenerateIDRetriesOnCollision`.
  - `TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat`, below.
- `.golangci.yml` enables `paralleltest` for top-level tests
  (`ignore-missing-subtests: true`) through the already-required linter.

## Red observations during the repair

1. The first full measurement after the change, before the maintenance
   exception, was `go test -race -cover -count=1 -timeout 30m -json ./cmd/fid`
   — FAIL (`exit_code: 1`) in
   `TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat`:
   `maintenance warnings differs across fixed-grace branches`, with FID-299
   warned in the second run only.
   - The test was **already parallel before this repair**. It seeds FID-299 at
     exactly 299 s old and runs two CLI commands against the wall clock, which
     leaves about one second of margin. Heavier contention used that margin up.
   - Repair: it runs serial, with an explanatory comment, and its assertions
     are unchanged. `internal/maintenance` proves the exact 299/300 boundary
     deterministically with an injected `Now`
     (`internal/maintenance/test_helpers_test.go:381-384`).
2. The first exact gate (attempt 1 below) failed for a reason unrelated to
   this repair.

## Flake guards (DEC-189 / DEC-254, widened)

- `go test -race -count=2 -shuffle=on -timeout 60m ./...` — PASS, every package
  `ok` (terminated `exit_code: 0`). The largest package lines were `cmd/fid`
  903.214 s, `internal/issues` 576.472 s, `internal/journal` 555.994 s,
  `internal/export` 524.168 s, and `internal/workflow` 316.838 s. Both counts
  ran concurrently across packages, so these are not budget figures.
- `go test -tags integration -race -count=2 -shuffle=on -timeout 60m ./cmd/fid/integration/...`
  — PASS — `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	615.509s`
  (terminated `exit_code: 0`), about 308 s per pass.
- These two logs were written to the session scratch directory, which was
  removed when the session environment changed. The tokens above were copied
  from their terminal output before that happened.

## Exact task-11 gates

Full logs: `esc-23-005-evidence/exact-gates-attempt-1.log` and
`exact-gates-attempt-2.log`. Each ran the five task-spec commands in order on
`3dff457`, after `go clean -testcache`.

### Attempt 1 — FAIL (retained)

- `go test -race -timeout 30m -coverprofile=coverage.out ./...` — FAIL —
  `FAIL	github.com/mbell-nzl/fid/cmd/fid	525.379s` (terminated
  `exit_code: 1`).
  - Every failure was `journal files = 2, want 1`, raised by `mustJournalPath`
    (`cmd/fid/init_test.go:698-709`). It hit 14 subtests of
    `TestExecutionBoundTransitionPredicateMatrix`,
    `TestTerminalSuccessTransitionRejectsEveryOtherNonqualifyingEvidenceState`,
    and `TestResumeProjectsExactActiveAndCompletedDurableStateWithoutRecovery`.
  - The run crossed 00:00 UTC (`date -u` shortly after:
    `Sun Sep 27 00:00:54 UTC 2026`). Journal files are named by UTC day
    (`internal/journal/writer.go:39`). Fixture templates cached before midnight
    held `2026-09-26.jsonl`, and clones appending after midnight created
    `2026-09-27.jsonl`.
  - This is pre-existing and independent of parallelism: a serial run crossing
    midnight fails the same way. It is carried forward, not changed here.
- The other four commands passed; the tagged gate was
  `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	335.029s`
  (terminated `exit_code: 0`).

### Attempt 2 — PASS

- `gofmt -l cmd internal` — PASS — `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` — PASS — `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` — PASS —
  `ok  	github.com/mbell-nzl/fid/internal/workflow	165.920s	coverage: 64.6% of statements`
  (terminated `exit_code: 0`), with 536 s of command wall time. `cmd/fid` was
  531.533 s.
- `go tool cover -func=coverage.out` — PASS —
  `total:											(statements)					62.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` —
  PASS — `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	326.421s`
  (terminated `exit_code: 0`)

| Package (DEC-227 times) | Task 11 before (attempts 1/2) | Attempt 1 | Attempt 2 | Budget |
|---|---|---|---|---|
| `cmd/fid/integration` | 778.278 / 766.070 s | 335.029 s | 326.421 s | 900 s |
| `cmd/fid` | 686.386 / 630.927 s | 525.379 s (FAIL, midnight) | 531.533 s | 600 s |
| `internal/workflow` | 277.603 / 252.947 s | 172.289 s | 165.920 s | 600 s |

Wall time for the ordinary gate command fell from roughly 690 s to 536 s, and
for the tagged gate from about 770 s to 336 s. `cmd/fid` passes with 68 s of
headroom. That is 88.6% of its budget, above DEC-258's 480 s headroom-warning
line, so the warning is raised here for the Tech Lead. Attempt 1 and attempt 2
differ by 6 s, against the earlier 15-20% spread.

## Remaining cost profile (for the headroom follow-up)

The profile was taken with
`go test -race -count=1 -timeout 30m -cpuprofile .tmp/prof/cmdfid.cpu ./cmd/fid`
— PASS — `ok  	github.com/mbell-nzl/fid/cmd/fid	286.641s` (terminated
`exit_code: 0`), 943 s of CPU samples.

- 33% is SQLite reading and parsing the full schema on each new connection
  (`sqlite3ReadSchema`/`InitOne`). Of the `internal/db.open` time, 47% comes
  from `journal.ensureProjectionFresh`, which opens and closes the projection
  before commands, and 33% from `workflow.openWorkflowDB`, which has 32 call
  sites. One command opens the database several times.
- 32% is `checkptr` instrumentation, which `-race` enables, running inside
  modernc's transpiled SQLite. fid production code imports no `unsafe`.
- With all packages parallel, the ordinary gate is CPU-throughput bound. Every
  large package finishes near the gate's wall time (about 535 s), so package
  times now track total suite CPU on this 6-core machine rather than one
  package's serial chain. Further headroom needs less CPU work: fewer
  connection opens per command (a product change) or dependency-scoped
  `-d=checkptr=0` (a gate-definition change). Both need separate authority.
