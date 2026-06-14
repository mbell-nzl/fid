# Phase 6 Decisions

> Per-phase decision log. Closed phase decision history lives here; live
> carry-forward constraints live in `../../current/decisions_active.md`.
>
> DEC-017 (readiness semantics) and DEC-018 (empty-state contract) were
> human-confirmed at planning (2026-06-14) and are formally recorded at phase
> wrap per the phase brief.

## DEC-019 — `fid ready` deterministic selection order (resolves ESC-004 / TRG-001)

- **Status:** Confirmed by human 2026-06-14 (option A).
- **Context:** Task 3 implemented the project-wide / prioritised-epic tie-break
  as `priority ASC` (nil last) → `created_at ASC` → `id ASC`. The PRD pins the
  behavior (PRD-ISSUE-008/012 deterministic; PRD-ISSUE-011 highest-priority) but
  did not spell out the full ordering tuple or nil-priority handling. TRG-001
  BLOCKed and raised ESC-004 because Task 4 makes this a user-visible scheduling
  contract.
- **Decision:** The bare `fid ready` selection order is, in precedence:
  1. `priority` ascending (PRD-ISSUE-015: `0` critical → lowest wins); a task
     with **no priority ranks last** (human-confirmed option A).
  2. `created_at` ascending.
  3. `id` ascending — final deterministic tie-break only, no product meaning.
  Sequential epics continue to order by `sequence` (PRD-ISSUE-009) and
  contribute only their current actionable task project-wide (PRD-ISSUE-012).
- **Rationale:** This matches the selection factors the design plan already
  documented inside the v0.2 `--explain` example (`priority`, `created_at`);
  `id` only guarantees a single deterministic winner. Nil-last treats an
  unprioritised task as least urgent ("if you didn't prioritise it, it waits").
- **Scope of change:** Documentation only. Task 3's implementation stands — no
  code change. Task 4 may proceed and wire the CLI surface against this order.
- **Spec sync:** `06_DESIGN_PLAN_CLI.md` §`fid ready` now states the bare
  selection order explicitly (lifted out of the `--explain`-only block).
  `05_DESIGN_PLAN.md` §6.3 documents the output *shape* only (not ordering):
  verified-no-change. No other spec touched.
- **References:** ESC-004, TRG-001, `tasks/task-3-done.md`, VAL-006; PRD-ISSUE-008,
  011, 012, 015.
