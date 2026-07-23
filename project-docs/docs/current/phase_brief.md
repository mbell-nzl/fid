# Phase Brief - Phase 18 (In execution — final-review remediation)

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

**In execution (Heavy Mode) — final-review remediation, 2026-07-23.** All five
tasks are **implemented and committed** (Task 5 at `37b3fb5`, reports finalized
at `d7f1199`). Five mid-phase blocks were resolved by the Tech Lead as do-now,
in-window remediations with no contract change: **DEC-067** (QA-plan block /
TRG-001), **DEC-068** (Review 1 / TRG-002, VAL-020 rebuild evidence), **DEC-069**
(Review 2 / TRG-003, VAL-006 valid-pair UPDATE evidence), **DEC-070** (Task 5 /
TRG-004, interrupted benchmark+build gate), and **DEC-071** (Task 5 / TRG-005,
TRD-PERF-003 100k-rebuild failure, fixed by rebuild-scoped durability
relaxation — target unchanged).

**Current blocker: QA Lead final review returned `BLOCK`** (`qa/final-review.md`)
with **24 PASS / 0 FAIL** on the validation contract but two Blocking findings,
triaged as **TRG-006 → ESC-001** and resolved by **DEC-072**:

1. **QA-BLK-18-FINAL-01 — `workflow.Load` deadlocks on a one-connection pool.**
   It issues nested per-workflow queries while the outer `workflows` rows still
   hold the connection (`internal/workflow/registry.go:23/38/112`), so four of
   the seven required gate commands (coverage, unit, integration, benchmark)
   **never terminate** on HEAD. Production `db.Open` sets no pool cap, so the
   real `fid maintenance check` entrypoint is unaffected (VAL-010 PASS via the
   real-entrypoint check) — but the heavy-mode gate is not green. Disposition:
   **fix the loader, not the fixture**; do-now within Task 5.
2. **QA-BLK-18-FINAL-02 — required Tech Lead spec sync absent (QA-30).**
   **Resolved by DEC-072:** `02_TRD.md` §§4.4/4.6, `03_TRD_SCHEMA.md` §§3/14.11/
   §15, the canonical-suite verdict table, and these current docs are now synced.

**Next concrete action:** resume the Phase Runner from the ESC-001 / TRG-006
STOP — re-spawn a Developer to apply the DEC-072 `workflow.Load` fix plus a
bounded-timeout single-connection regression test, rerun the full seven-command
phase gate **to completion in order** (capturing 1k/10k/100k rebuild + 10k-read
`ns/op` and the fresh coverage total), re-run the focused maintenance evidence
and the VAL-024 `rg` inventory, correct the over-claimed command records in
`task-5-done.md`, and commit; then the QA Lead re-runs the **final review**
against fresh evidence (Task 5 is the last task — no further mid-phase review),
followed by the human checklist. A still-hanging command, a changed registry
ordering, a broken parity assertion, or a 100k rebuild back at/over `30s`
escalates rather than passing. See DEC-067..DEC-072 in
`../phases/phase-18/decisions.md`.

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
| 1 | Workflow tables migration + seed defaults + gate seam | Complete | None |
| 2 | DB-backed workflow registry loader + v0.1 parity fixtures | Complete | 1 |
| 3 | Transition-path data-source swap (behavior-preserving) | Complete (VAL-020 remediated, DEC-068) | 1, 2 |
| 4 | issues.workflow_id/status FK migration (table rebuild) | Complete (VAL-006 remediated, DEC-069, re-review PASS) | 1 |
| 5 | Remaining consumer swap + maintenance parity + closure | Committed (`37b3fb5`, DEC-071 applied) — reopened for the DEC-072 loader-deadlock fix + gate rerun | 1–4 |

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
| QA plan | `../phases/phase-18/qa/qa-plan.md` | Active (QA-01..QA-30); blocks resolved by DEC-067 |
| Mid-phase reviews | `../phases/phase-18/qa/review-M.md` | Complete — `review-1`/`review-2` plus both remediation re-reviews `CONTINUE` |
| Final review | `../phases/phase-18/qa/final-review.md` | **BLOCK** (2026-07-20); re-run required after the DEC-072 fix |
| Human checklist | `../phases/phase-18/qa/human-checklist.md` | Generated, 4 items, all pending — execute after the gate is green |

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

- [x] QA Lead mid-phase review cadence for a 5-task heavy phase — settled: a
  review after Task 3 and after Task 4, each with a remediation re-review.

## Spec-Sync Note (Tech Lead)

**Applied 2026-07-23 (DEC-072).** `02_TRD.md` §§4.4/4.6 (+ the module map and
entity table) and `03_TRD_SCHEMA.md` §§3/14.11/15 now record that the workflow
tables are the runtime source of truth, that maintenance audits against the
projected tables, and that rebuild uses a relaxed-durability connection. The
canonical-suite verdict table is in `../phases/phase-18/decisions.md`. QA-30 is
verification at re-review, not authorship.

## Success / Exit Criteria

- [x] All five tasks implemented and committed.
- [ ] All five tasks complete **with green quality gates** — blocked by the
  DEC-072 loader deadlock; four required commands must terminate.
- [x] VAL-001..024 PASS with cited evidence (behavior parity + rebuild
  determinism + seam-present-no-behavior) — 24 PASS / 0 FAIL at final review,
  to be re-confirmed against fresh evidence.
- [ ] Final QA review has no unresolved blocking findings.
- [ ] Human checklist accepted (heavy mode).
- [x] TRD spec-sync applied (DEC-072).
- [ ] Phase sign-off recorded.

## Next Planning Checkpoint

After Phase 18 wrap: plan **Phase 19 — `fid workflow define/assign/list` +
escalation default workflow**, then Phase 20 (role_permissions migration +
closure) to finish v0.2a.
