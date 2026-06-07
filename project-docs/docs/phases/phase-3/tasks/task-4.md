# Phase 3 - Task 4: Rebuild Performance Verification At Single-Project Scale

## Objective

Verify and lock in `fid rebuild` performance and determinism at scale. Add a
seeded-journal benchmark (`go test -bench`) that records rebuild elapsed at
single-project scale (~1k events) and at the TRD-PERF-002 tier (10k events), and
a deterministic package-level test proving rebuild over the scale fixture is
non-destructive and reproducible. No new rebuild *features* — this is the
performance/foundation-hardening proof the roadmap requires.

## Requirements

- PRD refs: PRD-OPS-003 (auto-rebuild before serving stale reads), PRD-OPS-004
  (`fid rebuild`)
- TRD refs: TRD-PERF-002 (full rebuild at 10k events < 5s), TRD-ARCH-001
  (DB is a deterministic projection of the journal)
- Design refs: `06_DESIGN_PLAN_CLI.md` §18 (`fid rebuild`)
- Quality refs: Quality Plan §4.3 (COV-001 `internal/journal` ≥ 80%); phase-level
  gate includes `go test -bench . -benchmem ./...`

## Validation Contract Coverage

Assertions live in `project-docs/docs/phases/phase-3/qa/validation-contract.md`.

- VAL-006: `fid rebuild` is within the TRD-PERF-002 budget at 10k events and
  sub-second at single-project scale (~1k) — benchmark evidence.
- VAL-007: rebuild over the scale journal is deterministic and non-destructive
  (event count matches; delete-DB + rebuild identical; rebuild-vs-rebuild
  identical).

## Scope

### Files To Create Or Modify

- `internal/journal/rebuild_bench_test.go` (new: a journal-seeding helper +
  `Benchmark*` for ~1k and 10k events; records elapsed)
- `internal/journal/replay_test.go` or a new `*_test.go` (deterministic
  scale-rebuild test: event count, pre/post delete-DB equality, rebuild×2
  equality)
- A reusable seed helper for valid `issue.created` (and a small representative
  mix) events at a parameterised count
- `project-docs/docs/phases/phase-3/tasks/task-4-completion.md`

### Files To Read For Context

- `internal/journal/replay.go` (`RebuildWorkspace`, `ReplayJournal`,
  `EnsureProjectionFresh`)
- `internal/journal/writer.go` (event append) and `internal/journal/project.go`
  (projection) for building valid journals
- `internal/journal/phase2_rebuild_test.go` (existing rebuild determinism test to
  extend the scale pattern from)
- `02_TRD.md` §7 (TRD-PERF-002), `06_DESIGN_PLAN_CLI.md` §18

### Out Of Scope

- The 100k-event tier (TRD-PERF-003) — out of scope this phase.
- Read-latency benchmarks (TRD-PERF-001) and write-latency (TRD-PERF-005).
- Checkpoint/partial replay (TRD §11 open question) — full replay only.
- Any change to rebuild behaviour; this task measures and proves, it does not
  refactor the replay path.

## Acceptance Criteria

- [ ] A benchmark seeds a valid journal of N events and rebuilds it; running it
  records elapsed for ~1k and 10k events. (VAL-006)
- [ ] Recorded benchmark evidence shows the 10k-event rebuild < 5s and the
  ~1k-event rebuild < 1s. (VAL-006)
- [ ] A deterministic test proves: replayed event count equals the seeded count;
  the projection is identical before vs after delete-DB + `fid rebuild`; two
  successive rebuilds produce identical projection state. (VAL-007)
- [ ] No trigger false-positive rejects a valid journal-applied row during the
  scale rebuild. (VAL-007)

## TDD Plan

- First failing test (`internal/journal`): assert the scale-rebuild reconstructs
  N events and matches pre-delete state — red until the seed helper + scale test
  exist.
- Additional tests/benchmarks: rebuild×2 equality; the `Benchmark*` for ~1k and
  10k; (optional) `b.ReportMetric` for elapsed.
- Regression risks: seeding must use valid journal events (preconditions/kind
  guards satisfied) so triggers do not falsely fire; the benchmark must not be
  flaky — keep the threshold margin large and document the measured numbers; hold
  COV-001 `internal/journal` ≥ 80%.

## Documentation Updates

- `project-docs/docs/phases/phase-3/tasks/task-4-completion.md` - completion
  report citing VAL-006/007 evidence, including the captured benchmark output
  (elapsed for ~1k and 10k events) as deterministic evidence.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None. Do not read held-out scenarios. (The journal
  corruption-recovery scenario theme is exercised structurally in task 5.)
- Cross-surface parity impact: None.
- Audit or release-readiness impact: establishes the first performance benchmark;
  note it as the seed for future perf-tier checks.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Every required command
must be recorded in the completion report as `PASS`, or as `BLOCKED` with reason
in a blocker report. Missing, silently skipped, or narrower substituted commands
are failed gates.

Conditional commands (REQUIRED for this task):

```text
go test -bench . -benchmem ./...
```

List any conditional command required by this task and why:

- `go test -bench . -benchmem ./...` is **required**: VAL-006 evidence is the
  recorded benchmark output. Capture the elapsed numbers for ~1k and 10k events
  in the completion report. (Run the race-enabled suite separately; do not run
  benchmarks under `-race` for the timing numbers.)

## Notes For Developer

- VAL-006 is environment-sensitive; the deterministic anchor is VAL-007. Use a
  generous threshold margin and record the actual measured elapsed as evidence so
  QA Lead can confirm sub-second at single-project scale.
- Build the seed journal from valid events that satisfy kind/precondition guards
  (Phase 2 triggers are active) so replay does not falsely reject rows.
- Reuse the existing rebuild path (`RebuildWorkspace` / `ReplayJournal`); do not
  add a faster bespoke replay just for the benchmark — measure the real path.
- Keep all logic in `internal/journal` (COV-001); no `cmd/fid` change needed.
