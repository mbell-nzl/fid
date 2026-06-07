# Phase Brief - Between Phases (Phase 3 CLOSED → Phase 4 planning)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 3 (v0.1a Output formats + foundation hardening) is **CLOSED**
> and wrapped (2026-06-07, branch `phase-3/output-formats-hardening`). Phase 4
> (v0.1b Issue Lifecycle) is **NOT YET PLANNED**. There is no active execution
> phase until the Tech Lead drafts the Phase 4 validation contract and task
> specs and the human approves them.

## What Just Shipped (Phase 3)

The `--format` engine is complete for `fid issue get`: JSON, TOON, markdown, and
YAML from one code path. Persistence hardening landed: rebuild is within budget
at scale and deterministic; corruption (malformed tail) and precondition
conflicts render structurally on both `fid rebuild` and the read catch-up path.

- All VAL-001..VAL-010 **PASS**. QA final **PASS** (0 blocking/major/improvement).
  Human checklist **PASS** (2/2). Archive: `../phases/phase-3/`.
- Decisions: DEC-012 (`gopkg.in/yaml.v3`, MIT, for `--format yaml`), DEC-013
  (generic render pipeline + typed markdown seam). `02_TRD.md` §9 gained a
  `gopkg.in/yaml.v3` dependency row at wrap.

## Next: Plan Phase 4 (v0.1b Issue Lifecycle)

Roadmap `10_IMPLEMENTATION_ROADMAP.md` → "v0.1b — Issue Lifecycle". Candidate
scope (Tech Lead owns the final slice):

- `fid issue create` (epic/task/subtask; kind hierarchy; default workflow per
  kind), `fid issue get` (enrich + `--section`), `fid issue list`,
  `fid issue update`, `fid issue transition` (workflow state machine).
- `epic_metadata` scheduling (sequential/prioritised), `fid ready`.
- Hard-coded workflows + transition rules (`02_TRD.md` §4.6, TRD-ARCH-009).

v0.1b is large; consider slicing into more than one phase, each exiting on an
agent-runnable command (roadmap Guiding Principle #4, tracer-bullet rule).

## Phase-Start Checklist (Tech Lead)

1. Read the roadmap v0.1b entry and relevant PRD/TRD/design sections.
2. Review the Carry-Forward Register in `decisions_active.md` — several items
   target v0.1b (see below); schedule, re-defer, or fold each into a task.
3. Create `../phases/phase-4/` dirs; draft `qa/validation-contract.md`.
4. Break into ordered one-session tasks; cite VAL IDs; list the task gate and
   any QA-plan command per task; list doc updates.
5. Refresh this brief with the Phase 4 validation contract index and present to
   the human for approval before execution.

## Carry-In Register (open items to weigh at Phase 4 planning)

- **CF-004 (v0.1b):** sequence-vs-scheduling backstop trigger — add the
  cross-table rule when `epic_metadata` population, tasks-under-epics, and
  `fid ready` land. Exactly this phase's scope.
- **CF-005 (v0.1b):** `activity_log` projection-only guard + `activity` effect op
  — add when `activity_log` population ships (issue lifecycle writes).
- **CF-006 (v0.1b):** centralise the enriched issue read-model in
  `internal/issues` before `fid issue list` / `fid issue get --section` reuse the
  markdown-only read shape now living in `cmd/fid/issue.go` (TRD-ARCH-005).
- **CF-003 (later, v0.1d):** policy bootstrap corrupt-tail tolerance — not v0.1b
  scope unless mutable role events appear; otherwise re-defer.
- **COV carry-in (DEC-007):** name package-level tests in COV-001 packages
  (`internal/issues`, `internal/journal`, `internal/policy`); hold COV-002
  `cmd/fid` ≥ 70%.

## Current Code State (planning facts)

- `internal/format/` ships JSON/TOON/markdown/YAML; `fid issue get` is still the
  only read command with `--format`. Adding a read command means wiring it
  through `fidformat.Write` (generic) or a typed renderer (markdown-style).
- `internal/journal` replayer renders `CorruptJournalError` and
  `PreconditionConflictError` structurally on rebuild and read catch-up.
- Full v0.1 base schema + triggers exist (Phase 2); `issues_enriched` view is
  live. No `fid issue create`/lifecycle commands yet — issues are created via
  raw journal events in tests only.
- No `integration` build tag in use; CLI tests run under `go test -race ./...`.
