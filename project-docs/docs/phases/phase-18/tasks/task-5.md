# Phase 18 - Task 5: Remaining consumer swap + maintenance parity + Phase 18 closure

## Objective

Repoint the remaining consumers of the package-global workflow Go map to the
DB-backed registry — `fid maintenance check` (§14.11 `(workflow_id, status)`
audit), issue get/list status-category display, the `fid dev *` composites, and
`fid ready` — then prove the fully-migrated system end-to-end: v0.1 behavior
parity, deterministic rebuild, non-destructive delete-and-rebuild, and
unchanged v0.1 journal history.

## Validation Contract Coverage

- VAL-010: `fid maintenance check` validates `(workflow_id, status)` against the workflow tables — clean on healthy, flags a seeded invalid pair, non-zero under `--warn-as-error`.
- VAL-011: `fid rebuild` reproduces the migrated projection deterministically; delete-and-rebuild is non-destructive with normalized parity of **every** projected table (not only workflow tables + issues) across a first and a second rebuild, plus equal `sqlite_schema`/FK integrity.
- VAL-013: a golden v0.1 journal (`issue.created`/`issue.transitioned`) replays into the migrated schema producing the expected migrated projection.
- VAL-021: every journal file on disk is byte-identical across migration, stale-version refresh, delete-and-rebuild, and a second rebuild (history never rewritten).
- VAL-017: phase-end proof that the gate-obligation seam stays behavior-free — no production path can satisfy or query receipt/gate state for `required_gates_current`/`has_completion` (runtime injected-rule outcome + static production-source inventory).
- VAL-022: post-swap `fid issue get`/`list` status-category presentation and `fid ready` selection/order match frozen v0.1 expectations.
- VAL-023: post-swap `fid dev start`/`finish`/`block` outcomes and effects match frozen v0.1 expectations, using DB-resolved categories.
- VAL-024: no production (non-test) Go file reads the frozen `definitions` map after the swap; it remains only as test/parity fixtures.
- Performance controls (no VAL): a 100k-event rebuild benchmark (<30s, TRD-PERF-003) and a 10k-issue registry-backed read benchmark (<500ms, TRD-PERF-001) proven by the Task 5 benchmark gate (QA-BLK-18-03).

## Binding Constraints

- **Swap every remaining `workflow.*` package-function consumer to the registry.**
  From the current call sites: `internal/maintenance/check.go` +
  `internal/maintenance/warnings.go` (§14.11 workflow-pair audit and any
  category-derived warnings), `internal/issues/get.go` / `issue.go` /
  `detail.go` / `list.go` / `ready.go` (status-category display and ordering),
  and `internal/dev/{start,finish,block,dev}.go` (category checks in the
  composites). After this task, **no production code path reads the Go
  `definitions` map**; it remains only as frozen parity fixtures for tests.
- **Maintenance §14.11 now validates against the workflow tables.** The audit
  must resolve each issue's `(workflow_id, status)` against `workflow_states`
  (the migrated source of truth) rather than the Go map. Behavior is otherwise
  unchanged: deterministic findings ordered by check, entity ID, then message;
  `--warn-as-error` exits non-zero on a violation.
- **Behavior parity everywhere.** Category display, ready-queue ordering, and
  dev-composite category gating must be identical to v0.1. Because the seed and
  the Go map are behavior-identical, this is a source swap, not a semantic change.
- **Closure evidence.** Provide a delete-and-rebuild parity check over the
  **complete** projection — normalized snapshots of every projected table (all
  four workflow tables and every pre-existing table), compared after a first
  rebuild and again after a second rebuild, with equal `sqlite_schema` and clean
  `foreign_key_check`; document any intentionally normalized volatile field
  (VAL-011). Provide a golden replay fixture proving a v0.1 journal reprojects
  into the expected migrated projection (VAL-013). Assert every journal file is
  byte-identical across four boundaries — migration, stale-version refresh,
  delete-and-rebuild, and a second rebuild (VAL-021; DEC-066: history never
  rewritten).
- **Performance controls (QA-BLK-18-03).** Add a 100k-event rebuild benchmark
  asserting <30s (TRD-PERF-003) and a 10k-issue registry-backed single-read
  benchmark asserting <500ms (TRD-PERF-001), reusing the configurable
  `seedScaleReplayWorkspace` helper as the fixture starting point. Both must run
  under the `go test -bench . -benchmem ./...` gate in this task. The existing
  1k/10k rebuild benchmarks remain (10k rebuild <5s, TRD-PERF-002).
- **Seam behavior-absence (phase-end).** Re-confirm the gate-obligation seam is
  behavior-free after every swap: inject a reserved-rule requirement into a
  transition fixture and prove no receipt/gate evaluator runs, and run the static
  production-source inventory showing no receipt/gate evaluator or query for
  `required_gates_current`/`has_completion` (VAL-017; VAL-005 also stays true).

## Scope

### Files To Create Or Modify

- `internal/maintenance/check.go`, `internal/maintenance/warnings.go`
- `internal/issues/get.go`, `internal/issues/issue.go`, `internal/issues/detail.go`,
  `internal/issues/list.go`, `internal/issues/ready.go`
- `internal/dev/start.go`, `internal/dev/finish.go`, `internal/dev/block.go`, `internal/dev/dev.go`
- Their tests, plus a phase-closure test/fixture (full-projection
  delete-and-rebuild parity; golden v0.1 journal replay; journal-byte
  immutability; seam behavior-absence). Existing e2e closure tests under
  `cmd/fid/` updated/extended as needed.
- Benchmark additions under `internal/journal/` (100k-event rebuild) and
  `internal/issues/` (10k-issue registry-backed read), reusing
  `seedScaleReplayWorkspace` (QA-BLK-18-03).

### Files To Read For Context

- `internal/workflow/registry.go` (loader)
- `project-docs/specs/03_TRD_SCHEMA.md` §14.11 (maintenance parity)
- Current consumers listed above (call sites of `workflow.CategoryForState` /
  `RulesForEdge` / `Workflow`)
- `internal/journal/` rebuild + golden-fixture tests

### Out Of Scope

- Deleting the Go `definitions` map (kept as parity fixtures).
- `fid workflow define/assign/list`, escalation-default (Phase 19).
- `role_permissions` migration (Phase 20); any gate/receipt behavior (v0.2b).

## Acceptance Criteria

- [ ] No production (non-test) Go file reads the frozen `definitions` map; all
  consumers use the DB registry, proven by the two `rg` inventory commands
  returning only test/fixture references, paired with runtime parity (VAL-024)
- [ ] `fid maintenance check` resolves `(workflow_id, status)` from the workflow
  tables: clean on healthy, flags a seeded invalid pair, non-zero under
  `--warn-as-error` (VAL-010)
- [ ] `fid issue get`/`list` status-category presentation and `fid ready`
  selection/order match frozen v0.1 expectations across the parity matrix (VAL-022)
- [ ] `fid dev start`/`finish`/`block` outcomes and effects match frozen v0.1
  expectations using DB-resolved categories (VAL-023)
- [ ] Delete-and-rebuild is non-destructive with normalized parity of every
  projected table across a first and second rebuild, plus equal `sqlite_schema`
  and clean `foreign_key_check` (VAL-011)
- [ ] A golden v0.1 journal replays into the expected migrated projection (VAL-013)
- [ ] Every journal file is byte-identical across migration, stale-version
  refresh, delete-and-rebuild, and a second rebuild (VAL-021)
- [ ] Phase-end proof that no production path evaluates or queries receipt/gate
  state for `required_gates_current`/`has_completion` (runtime + static
  inventory) (VAL-017)
- [ ] Benchmark gate: 100k-event rebuild <30s (TRD-PERF-003) and 10k-issue
  registry-backed read <500ms (TRD-PERF-001), alongside the retained 10k-event
  rebuild <5s (TRD-PERF-002) (QA-BLK-18-03)
- [ ] Full phase-level gate green

## First Failing Test

A maintenance test that seeds an issue whose `(workflow_id, status)` pair is not
in `workflow_states` and asserts `fid maintenance check` reports a structural
finding sourced from the tables — fails until the audit reads the registry.

## Documentation Updates

- None authored by the Developer. The Tech Lead syncs `03_TRD_SCHEMA.md`
  §14.11 / §4.6 and `02_TRD.md` wording (maintenance now audits against the
  workflow tables; defaults are migrated) at phase wrap. Flag any behavior
  nuance discovered here in the completion report.

## Exact Commands To Run

Run exactly the full phase-level gate:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates.

## Notes For Developer

- Closure task: depends on Tasks 1–4. Run the full phase-level gate (this task
  uses the phase-level gate, not the shorter task gate) because it proves the
  v0.2a slice-1 exit criteria.
- Record the 1k/10k/100k rebuild and 10k-issue read benchmark durations
  (ns/op) verbatim in the completion report so QA can check them against
  TRD-PERF-001/002/003; a benchmark that regresses past its target is a blocker,
  not a silent pass.
- If any consumer's behavior would change under the registry, stop and escalate
  rather than adjusting the assertion.
- Completion report at `project-docs/docs/phases/phase-18/tasks/task-5-done.md`.

## Remediation note (DEC-071) — TRD-PERF-003 100k-rebuild failure

The consumer swap and closure evidence are sound; the only open item is the
performance control. On the DEC-070 completion run the full gate ran to
completion and the 100k-event rebuild self-asserted a genuine failure
(`33.3846063s`, want `<30s`; TRD-PERF-003). The benchmark rows are linear at
~0.34 ms/event (1k `361.9ms`, 10k `3447ms`, 100k `33.38s`), so the fix must shave
the per-event constant, not chase a super-linear hotspot. Root cause: the rebuild
replays each event in its own transaction (`ProjectEvent` → begin/commit +
per-event `WriteProjectionMeta`), and the rebuild connection uses SQLite defaults
(`synchronous=FULL` + rollback journal), so a 100k-event rebuild pays ~100k
fsyncs. Retain the existing working tree; **do not weaken the `<30s` target.**

Do exactly this remediation (no assertion/contract/scope change; DEC-071):

- **Relax durability on the rebuild-only connection.** Open the connection used
  by `rebuildWorkspaceLocked` (`internal/journal/projection.go`) with
  `synchronous=OFF` and `journal_mode=MEMORY`, added via the DSN `_pragma`
  mechanism in `internal/db` (mirror the existing `requiredPragmas` in
  `sqliteDSN`; add a rebuild-scoped open variant such as `OpenForRebuild`). Do
  **not** set these with `conn.Exec` on the pooled `*sql.DB` (a pragma set that
  way applies to only one pooled connection); the DSN `_pragma` list is applied
  to every connection the driver opens.
- **Keep `foreign_keys=ON` and `recursive_triggers=ON`** on the rebuild
  connection. Use `journal_mode=MEMORY` (not `OFF`) so transactional rollback
  still works and per-event error surfacing (e.g. `PreconditionConflictError`)
  is unchanged.
- **Do not touch the global `db.Open` / normal command-path durability.** Only
  the rebuild connection relaxes durability; normal writes keep full durability.
- **Add a focused test** proving the rebuild connection carries the relaxed
  pragmas (`synchronous`/`journal_mode`) while the normal `Open` connection does
  not, and confirm every rebuild-parity test (VAL-011/013/021) stays green.

Why this is behavior-preserving and safe: pragmas change durability/journaling,
not row data, so the projection stays byte-identical (VAL-011/013/021 hold), and
rebuild reconstructs a *derived* store from the canonical journal — from scratch
(the DB is removed first) and under the write lock — so a crash mid-rebuild just
discards the partial DB and re-rebuilds (ADR-0001; TRD-ARCH-001). Durability
during a rebuild has no contract value.

Then resume the DEC-070 completion steps unchanged: run the full phase-level gate
(the seven commands above) **to completion in order**, capturing verbatim
`ns/op`/`ms` for the 1k/10k/100k rebuild and the 10k-issue registry read; run
`go build`; re-run the VAL-024 two-`rg` production inventory + runtime parity;
replace this blocker report with a normal completion report citing
VAL-010/011/013/017/021/022/023/024 with the benchmark durations; produce the
two-commit sequence.

**Guardrail (DEC-071).** If the 100k-event rebuild is still `≥30s` after this
fix, do **not** weaken the target — stop and escalate to the Tech Lead (the next
lever, batched commits, is a higher-risk change the Tech Lead owns). If any
rebuild-parity assertion breaks, that is a genuine failure — stop and escalate
rather than adjusting the assertion.

## Remediation note (DEC-072) — registry-loader connection deadlock

DEC-071 is done and committed (`37b3fb5`); the 100k rebuild passes and every
VAL-001..024 assertion has evidenced PASS. The QA Lead final review nonetheless
returned **BLOCK**: `workflow.Load` deadlocks on a one-connection pool, so four
of the seven required gate commands (coverage, unit, integration, benchmark)
**never terminate** on HEAD. This is the only product-code work left in Task 5.

Reproduce the failure first (it must fail before your fix):

```text
go test -race -count=1 -v -timeout 5s ./internal/maintenance -run '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'
```

Expected today: `panic: test timed out after 5s`, with
`database/sql.(*DB).conn` ← `workflow.loadStates`
(`internal/workflow/registry.go:112`) ← `workflow.Load` (line 43) ←
`maintenance.checkWorkflowStates` (`internal/maintenance/check.go:208`).

**Root cause (already diagnosed — do not re-litigate).** `workflow.Load`
(`internal/workflow/registry.go:23`) opens `SELECT id FROM workflows ORDER BY id`
and then, **while still iterating that `*sql.Rows`** (line 38), calls
`loadDefinition` → `loadStates` / `loadTransitions` /
`loadTransitionRequirements`, each of which issues a fresh `conn.Query(...)`
(lines 112, 137, 173). `database/sql` keeps the connection checked out until the
outer rows are drained or closed, so with `SetMaxOpenConns(1)` the first nested
query waits forever for a connection the outer rows still own. Production
`db.Open` sets no pool cap, which is why the real `fid maintenance check`
entrypoint works and the fixture hangs — the loader is still wrong.

Do exactly this remediation (no assertion, contract, scope, or gate change; DEC-072):

- **Fix `workflow.Load` to release the outer rows before any nested query.**
  Materialize the workflow IDs from the first query into a slice, close/drain
  that result set, **then** loop the slice calling `loadDefinition`. Preserve
  `ORDER BY id` (slice order), read-only behavior, error-wrapping strings, and
  registry contents exactly — this is a connection-lifetime fix, not a semantic
  one. `maintenance.checkActivityReferences` (`internal/maintenance/check.go:134`)
  already uses this drain-then-lookup pattern; match it.
- **Add a bounded-timeout regression test** in `internal/workflow` that calls
  `Load` against a pool with `SetMaxOpenConns(1)` and **fails fast instead of
  hanging** — run `Load` in a goroutine and `select` against a short
  `time.After`, or drive it with a context deadline. It must go red (timeout)
  before the fix and green after. A future regression has to surface as a test
  **failure**, never as a suite that never returns.
- **Do not mask the defect.** Do not raise or remove `SetMaxOpenConns(1)` in any
  fixture (`internal/maintenance/test_helpers_test.go:54`,
  `cmd/fid/maintenance_test.go:149`, `internal/db/db_test.go:68` — these are
  deliberate determinism controls that caught a real bug), do not add `-timeout`
  flags to the gate commands, do not skip or delete
  `TestCheckWorkflowStatesReadsProjectedWorkflowTables`, and do not change
  `db.Open` pooling. Run the seven commands exactly as listed above.
- **No other `Load` call site needs changing.** Every caller
  (`internal/maintenance/{check,warnings}.go`,
  `internal/issues/{get,list,ready,transition}.go`,
  `internal/dev/{start,finish,block,gate}.go`) calls `Load` *before* opening its
  own rows, so fixing `Load` fixes all of them.

Then finish the task:

- Run the **full seven-command phase-level gate to completion in order** (they
  will now terminate), capturing verbatim `ns/op` for the 1k/10k/100k rebuild and
  the 10k-issue registry read, plus the fresh `go tool cover -func` total.
- Re-run the focused maintenance evidence (healthy, seeded-invalid pair, and
  `--warn-as-error` non-zero exit) and the VAL-024 two-`rg` production inventory
  plus runtime parity.
- **Correct `task-5-done.md`.** Its current records for commands 3, 5, and 6 are
  over-claimed PASS (they did not terminate), and command 4's total came from a
  retained profile, not a fresh one. Replace all four with the fresh verbatim
  results and cite VAL-010 with both the focused maintenance evidence and the
  real-entrypoint check. Then produce the standard commit sequence.

**Guardrail (DEC-072).** If the fix changes registry contents or ordering, if any
VAL parity assertion breaks (VAL-002/003/004/014 parity set or
VAL-011/013/021 rebuild set), if the 100k rebuild returns to `≥30s`, or if any
gate command still fails to terminate — **stop and escalate to the Tech Lead**
rather than adjusting the assertion, the fixture, or the command.
