# Phase Brief - Phase 18 (Planning — pending Human approval)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session.
>
> **Size target:** under 150 lines.

## Phase Goal

Convert the v0.1 hard-coded workflow state machines into data-configurable
tables and repoint the transition engine to resolve from those tables —
**preserving v0.1 behavior exactly** — and reserve the DEC-066 gate-obligation
seam with no receipt behavior. This is v0.2a slice 1 of 3 (Phases 18–20) and
the foundational dependency for all later evidence obligations, which hang off
data-backed workflow rules. Success: after migration every issue carries
FK-backed workflow/status, every v0.1 transition behaves identically, a clean
`fid rebuild` reproduces the migrated state deterministically, and the seam is
present in schema with zero behavioral change.

## Status

**Approved (Heavy Mode) — 2026-07-19.** The Human approved the operating mode,
validation contract (VAL-001..024, Ready), and five-task list. The QA-plan
block (QA-BLK-18-01..03 / TRG-001) was resolved by the Tech Lead: contract
refined to **Ready** with VAL-015..024; Task 4 dependent-object inventory made
exact; Task 5 performance benchmarks added — see DEC-067 in
`../phases/phase-18/decisions.md`. Next concrete action: **resume the Phase
Runner** from the existing RUN-START ledger entry. On resume it re-spawns the
QA Lead to refresh the QA plan's IDs/mappings (activating VAL-015..024) and
clear its `BLOCKED BEFORE TASK 1` status, then drives the dev/QA loop from
Task 1.

## Operating Mode

**Heavy.** Entry criteria met: this is a data-migration / persistence-contract /
public-interface change (`09_AI_WORKFLOW.md` §4.3). Required artifacts:
validation contract (Ready, VAL-001..024), QA plan with fresh command/feasibility
validation, mid-phase cadence per QA plan, final review, human checklist, and
spec sync across all seven specs at wrap. Orchestration via Phase Runner.

## Task Index

Full task specs and completion reports live under `../phases/phase-18/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Workflow tables migration + seed defaults + gate seam | Pending | None |
| 2 | DB-backed workflow registry loader + v0.1 parity fixtures | Pending | 1 |
| 3 | Transition-path data-source swap (behavior-preserving) | Pending | 1, 2 |
| 4 | issues.workflow_id/status FK migration (table rebuild) | Pending | 1 |
| 5 | Remaining consumer swap + maintenance parity + closure | Pending | 1–4 |

## Validation Contract

Full contract (24 assertions, **Ready**):
`../phases/phase-18/qa/validation-contract.md`. All assertions are
machine-checkable; validator is QA Lead throughout (no Human perception
assertion applies). Refined 2026-07-19 (DEC-067) to unbundle independent
observables and close coverage gaps (QA-BLK-18-01/02); VAL-015..024 added.

| Owning window | Assertions |
|---|---|
| task-1 | VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-016, VAL-012 |
| task-2 | VAL-014, VAL-015 (advances VAL-002/003/004) |
| task-3 | VAL-008, VAL-009, VAL-020 |
| task-4 | VAL-006, VAL-007, VAL-018, VAL-019 |
| task-5 (incl. phase-end) | VAL-010, VAL-011, VAL-013, VAL-017, VAL-021, VAL-022, VAL-023, VAL-024 |

Performance controls (TRD-PERF-001/003 read + 100k-rebuild spot-checks) are a
Task 5 benchmark gate, not a `VAL-NNN` (QA-BLK-18-03).

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| QA plan | `../phases/phase-18/qa/qa-plan.md` | Drafted; blocked pending QA-plan mapping refresh on resume (blocks resolved by DEC-067) |
| Mid-phase reviews | `../phases/phase-18/qa/review-M.md` | Per QA plan (heavy; 5 tasks) |
| Final review | `../phases/phase-18/qa/final-review.md` | Pending |
| Human checklist | `../phases/phase-18/qa/human-checklist.md` | Pending |

## Active Constraints

- **Behavior-preserving migration.** The seeded tables and the v0.1 Go
  `definitions` map must be behavior-identical; the Go map is retained as frozen
  parity fixtures until every consumer is swapped (Task 5), then kept for tests.
- **Default workflows seed via migration DDL, not journal events** (system
  reference data; keeps v0.1 history unchanged, satisfies FK ordering on
  rebuild). Custom-workflow definitions become journal events in Phase 19.
- **Seam only.** Reserve `required_gates_current` / `has_completion` and a
  gate-obligation `params_json` shape; no seeded use, no receipt/gate code.
- **Cancel edges seeded explicitly** (`<state>→cancelled` for every non-cancelled
  state; no `cancelled→cancelled`); do not rely on a residual code special-case.
- No new user-facing command ships; behavior is observed through existing
  commands (`fid issue transition/get/list`, `fid maintenance check`, `fid rebuild`).

## Interfaces / Subsystems In Scope

- `migrations/0007_*`, `migrations/0008_*` — workflow tables + issues FK rebuild.
- `internal/workflow/` — DB-backed registry loader; Go map → frozen fixtures.
- `internal/issues/` — transition path, category display, ready; `internal/dev/*`.
- `internal/maintenance/` — §14.11 `(workflow_id, status)` audit source swap.
- `internal/db/`, `internal/journal/` — projection version bump, rebuild parity.

## Known Risks

- **issues-table rebuild (Task 4).** SQLite 12-step rebuild must recreate every
  object attached to `issues` (view, all triggers, indexes). Mitigation:
  enumerate from migrations 0001–0006; assert existence post-migration;
  `PRAGMA foreign_key_check` must be clean.
- **Consumer blast radius.** ~10 files read the Go map. Mitigation: registry
  proven equal to v0.1 defs in Task 2 before any swap; every intermediate merge
  stays behavior-consistent because seed == Go map.
- **Silent behavior drift.** Mitigation: characterization tests lock error
  messages and category/edge decisions before source swaps (Tasks 3, 5).

## Open Questions

- [ ] QA Lead: confirm the mid-phase review cadence for a 5-task heavy phase
  (recommend a review after Task 3 and after Task 4, the two highest-risk swaps).

## Spec-Sync Note (Tech Lead, at wrap)

TRD §14.11 and §4.6 wording ("validate against the canonical v0.1 Go
definitions" / "v0.1 hard-coded") will need updating to reflect that maintenance
now audits against the workflow tables and defaults are migrated. Track as a
phase-decision spec-sync item; not a Developer doc task.

## Success / Exit Criteria

- [ ] All five tasks complete with green quality gates.
- [ ] VAL-001..014 PASS with cited evidence (behavior parity + rebuild
  determinism + seam-present-no-behavior).
- [ ] Final QA review has no unresolved blocking findings.
- [ ] Human checklist accepted (heavy mode).
- [ ] Phase sign-off recorded; TRD spec-sync applied.

## Next Planning Checkpoint

After Phase 18 wrap: plan **Phase 19 — `fid workflow define/assign/list` +
escalation default workflow**, then Phase 20 (role_permissions migration +
closure) to finish v0.2a.
