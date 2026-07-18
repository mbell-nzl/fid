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

**Planning — pending Human approval.** No branch, no implementation. On
approval, plan hands to Phase Runner for the dev/QA loop. Next concrete action:
Human approves operating mode, validation contract, and task list; QA Lead then
drafts the QA plan (refining the contract) before task execution.

## Operating Mode

**Heavy.** Entry criteria met: this is a data-migration / persistence-contract /
public-interface change (`09_AI_WORKFLOW.md` §4.3). Required artifacts:
validation contract (done, draft), QA plan with fresh command/feasibility
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

Full contract: `../phases/phase-18/qa/validation-contract.md`. All assertions
are machine-checkable; validator is QA Lead throughout (no Human perception
assertion applies).

| VAL ID | Source | Owning Task/Window | Validator | Status |
|---|---|---|---|---|
| VAL-001 | TRD-ARCH-009; TRD_SCHEMA §3 | task-1 | QA Lead | Draft |
| VAL-002 | TRD §4.6 (seed verbatim) | task-1 | QA Lead | Draft |
| VAL-003 | TRD §4.6 | task-1 | QA Lead | Draft |
| VAL-004 | TRD §4.6; PRD-WF-007–010 | task-1 | QA Lead | Draft |
| VAL-005 | PRD-WF-014/015 seam; DEC-066 | task-1 | QA Lead | Draft |
| VAL-006 | TRD-ARCH-009; TRD_SCHEMA §2.1 | task-4 | QA Lead | Draft |
| VAL-007 | TRD-ARCH-008; TRD_SCHEMA §2, §14 | task-4 | QA Lead | Draft |
| VAL-008 | v0.2 exit; PRD-WF-001–010 | task-3 | QA Lead | Draft |
| VAL-009 | PRD-CLAIM-005/006 | task-3 | QA Lead | Draft |
| VAL-010 | PRD-OPS-005–008; §14.11 | task-5 | QA Lead | Draft |
| VAL-011 | PRD-OPS-003–004; DEC-066 | task-5 (phase-end) | QA Lead | Draft |
| VAL-012 | TRD §4.3; projection_meta | task-1 | QA Lead | Draft |
| VAL-013 | DEC-066; TRD-ARCH-002–003 | task-5 | QA Lead | Draft |
| VAL-014 | Roadmap v0.2a; TRD §4.6 | task-2 | QA Lead | Draft |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| QA plan | `../phases/phase-18/qa/qa-plan.md` | Missing (QA Lead to draft) |
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
