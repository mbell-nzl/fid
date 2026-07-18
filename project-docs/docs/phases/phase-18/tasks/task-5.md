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
