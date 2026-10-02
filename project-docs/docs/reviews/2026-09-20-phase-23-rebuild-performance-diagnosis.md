# Phase 23 rebuild performance diagnosis

User-requested investigation, 2026-09-20. Investigated HEAD
`150569adb859cbd28b177ea1950f6dfdde6cd51d`. This is diagnostic evidence and a
repair recommendation, not task completion, a gate-definition change, or
authority to resume the stopped phase. Existing dirty triage/escalation/ledger
artifacts were preserved. Product source and accepted assertions are unchanged.

## Finding

Keep the reference-hardware 100k rebuild target below 30 seconds. The current
failure reproduces, and there is a measured optimization path with useful
headroom. It does not require weakening constraints, batching transactions, or
changing locking/durability settings.

The same benchmark at the immediate pre-implementation commit
`7a865d8ccbb073b51f485f393c08728189117b98` takes 27.687–28.002 seconds;
current HEAD takes 31.608–32.977 seconds across diagnostic variants. Migration
execution itself is only 58 ms in the timed current run; replay takes 32.918 s.
The failure is in ongoing per-Event work, not the one-time migration.

CPU profiles show repeated SQL compilation as a large avoidable cost. Every
Event begins a transaction, applies Effects through uncached SQL execution,
upserts projection metadata, and commits. The current sqlite driver's
`conn.exec` prepares and closes a statement on every call. Its explicit
prepared-statement implementation does retain the SQLite statement handle.

| Cumulative sampled CPU, not elapsed wall time | Pre-task profile | Current profile | Statement-reuse prototype |
|---|---:|---:|---:|
| SQLite statement preparation | 8.83 s | 11.20 s | 0.71 s |
| SQLite foreign-key check code generation | 1.63 s | 3.26 s | Below displayed profile cutoff |
| Child-FK scan code generation | 1.38 s | 2.93 s | Below displayed profile cutoff |

These rows overlap; do not add them. The new schema introduces several more
references to `issues`. The benchmark inserts Issues even though it does not
yet exercise the new Review/Judgment Event families. The before/after profiles
strongly implicate the increased cost of compiling FK-aware writes as a
contributor to the regression. They do not assign every millisecond of the
wall-time difference to one schema object. Per-Event commit and metadata costs
were already substantial before Phase 23.

An external diagnostic overlay reuses seven SQL statement shapes for Effect
writes and projection metadata. It retains per-Event transactions, metadata
updates, ordering, FK checks, validation, and the existing rebuild pragmas.
The profiled prototype takes 24.129 s; a subsequent fresh unprofiled process
takes 23.367 s. Both the retained raw-file parity test
and the scale replay determinism test pass with this overlay.

The prototype is deliberately not production code: its serial-benchmark
globals must be replaced by explicit ownership scoped to one rebuild and its
connection/transaction lifecycle. It is evidence for the optimization, not a
patch to dispatch or proof of full correctness. The effect-write-only prototype
was insufficient and its failed result is retained below.

## Recommendation and boundaries

1. Assign a bounded rebuild-performance repair covering `internal/journal`
   and the projection-metadata seam in `internal/db`, with the 30-second target
   and existing semantic/parity assertions unchanged. Reuse prepared statements
   with explicit lifetime management; avoid global caches and cross-workspace
   state. Preserve per-Event rollback, conflict location, checkpoint metadata,
   and staged replacement on failure. Close statements on every exit path.
2. Require focused regression evidence for those behaviors, the existing
   parity checks, the exact task gates, and the accepted three fresh phase-end
   samples. The diagnostic prototype is not any of those task-gate receipts.
   Aim for operating headroom around the measured 24 s or better, rather than
   treating 29.9 s as a durable solution.
3. Keep the 30-second product target tied to the recorded reference instance.
   It corresponds to roughly 3,333 Events/second and is an achievable recovery
   budget. No user research establishes 30 rather than 35 seconds as an exact
   usability boundary; nevertheless, raising it now would absorb a reproducible,
   avoidable regression. Shared-host load and storage still affect measurements;
   the current three-fresh-sample rule is the accepted policy, not permission
   to retry until green. Any future measurement-policy change is separate.
4. Assign separate fixture repair for the execution-transition assertion.
   `executionTransitionFixture` clones a cached workspace containing a claim
   acquired with TTL 60 seconds; `loadExecutionClaim` uses `time.Now()`. The
   cached claim ages throughout the suite. An expired qualifying fixture explains
   both extra reasons (`claim-expired`, then `run-inactive`) without changing
   transition semantics. This is a source-supported diagnosis, not a fresh
   reproduction of the entire race-suite failure. Prefer deterministic clock
   control or fresh ephemeral claim/run setup per case; preserve explicit expiry
   tests. Assertion provenance alone misses the later cached-fixture change:
   `d17f524ce3d7bdf7de4e715a03514cba2fe8d7a8` introduced that reuse.
5. Investigate the separate 600-second package budget with suite profiling.
   This session did not run or clear that gate; a faster rebuild does not prove
   the ordinary suite meets its budget.

The earlier report's `27.170 ms/rebuild` exclusive-lock experiment has no raw
output here and is not relied upon. Its stated unit should be audited before
reuse. Raw SQLite-file parity is also stronger than the schema spec's
byte-identical *canonical reads*: SQLite file headers contain change counters
affected by transaction/unlock history. That distinction merits separate review
if locking/batching is ever proposed, but the measured statement-reuse approach
passes the existing raw-byte assertion and needs no such change.

The benchmark itself checks replay count and time. Full projection parity is
covered by separate tests, not a full-state assertion inside the 100k benchmark.

## Session evidence

Reference identity observed in this session: Linux/amd64,
`6.6.114.1-microsoft-standard-WSL2`, Go `go1.26.2`, six exposed logical CPUs,
AMD Ryzen 7 5700G, `MemTotal: 8133236 kB`. Initial load average was
`0.17, 0.16, 0.37`. Benchmarks were run sequentially, not alongside another
benchmark or test suite. Profiling changes measurement conditions and is
identified explicitly below.

Scratch root: `/tmp/fid-rebuild-diagnosis-kLezCb`. Baseline is a `git archive`
extraction of the exact pre-implementation commit's `go.mod`, `go.sum`,
`internal`, and `migrations`. Scratch profiles, binaries, overlay JSON/source,
and readable pprof summaries remain there for inspection; scratch is not
durable project storage.

For the table, `R` is `/home/mark/projects/fid`, `S` is the scratch root, and
`B` is the exact command
`go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/`.
Extra options listed in the table were inserted before the package argument.
These are diagnostic outcomes recorded from this session's own command results.

| Working directory / command variant | Outcome | Measured rebuild | Verbatim output token | Observed termination |
|---|---|---:|---|---|
| R / B unchanged | FAIL | 31.607555266 s | `100k-event rebuild = 31.607555266s, want <30s` | `exit_code: 1` |
| R / B + `-cpuprofile S/current.cpu -o S/journal.test` | FAIL | 32.440746598 s | `100k-event rebuild = 32.440746598s, want <30s` | `exit_code: 1` |
| S/baseline / B unchanged | PASS | 28.002 s | `28002 ms/rebuild` | `exit_code: 0` |
| S/baseline / B + `-cpuprofile S/baseline.cpu -o S/baseline.test` | PASS | 27.687 s | `27687 ms/rebuild` | `exit_code: 0` |
| R / B + `-overlay S/timing.json` | FAIL | 32.976916641 s | `100k-event rebuild = 32.976916641s, want <30s` | `exit_code: 1` |
| R / B + `-overlay S/cache.json` (Effect writes only) | FAIL | 31.680707425 s | `100k-event rebuild = 31.680707425s, want <30s` | `exit_code: 1` |
| R / B + `-overlay S/cache-full.json -cpuprofile S/cache.cpu -o S/cache.test` | PASS | 24.129 s | `24129 ms/rebuild` | `exit_code: 0` |
| R / B + `-overlay S/cache-full.json` (fresh unprofiled confirmation) | PASS | 23.367 s | `23367 ms/rebuild` | `exit_code: 0` |

Focused diagnostic verification, working directory R:

- `go test -overlay /tmp/fid-rebuild-diagnosis-kLezCb/cache-full.json -race ./cmd/fid -run '^TestPristineWorkspaceClonePreservesRebuildParityAndIsolation$' -count=1`
  — PASS — `4.791s`
  (terminated `exit_code: 0`).
- `go test -overlay /tmp/fid-rebuild-diagnosis-kLezCb/cache-full.json ./internal/journal -run '^TestRebuildWorkspaceIsDeterministicAtSingleProjectScale$' -count=1`
  — PASS — `1.331s`
  (terminated `exit_code: 0`).

The existing failed Task 1 attempts remain unchanged. No threshold, assertion,
phase decision, recovery authority, or ledger state was edited by this diagnosis.

## References

- `project-docs/specs/02_TRD.md` section 9; Phase 18 DEC-071 and Phase 21 DEC-158.
- `internal/journal/rebuild_bench_test.go`, `projection.go`, `project.go`,
  `project_effects.go`; `internal/db/projection_meta.go` and `db.go`.
- Driver v1.51.0 `conn.go` and `stmt.go` in the local module cache.
- [SQLite locking-mode behavior](https://www.sqlite.org/pragma.html#pragma_locking_mode).
- [SQLite file change counter](https://www.sqlite.org/fileformat.html#file_change_counter).

## Implemented resolution

ESC-23-002 is **RESOLVED 2026-09-20 under DEC-238**. The Human delegated
bounded repair authority above the fixed agent roles. Rebuild statement reuse
and fresh/isolated test setup repair the causes without changing any gate or
assertion. Fresh exact gates pass: CLI **527.099s**, workflow **231.010s**,
tagged integration **669.586s**; rebuild samples **22.728s / 22.623s / 22.489s**.
Integration binary directories stay **0 → 0**. Task 1 is complete in its report;
the Runner must reconcile the retained STOP before advancing. Tasks 9-11 stay
held. See [repair evidence](../phases/phase-23/qa/esc-23-002-repair.md) and
[resume handoff](../phases/phase-23/qa/resume-after-esc-23-002.md).

The Human authorized implementation and cross-role closure under DEC-238.
The proposal above is now implemented; no performance threshold was raised.
