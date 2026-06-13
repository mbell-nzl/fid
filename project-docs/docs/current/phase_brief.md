# Phase Brief - Phase 5 (Mutation + Workflow Transitions)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 5 is **PLANNED** (Tech Lead). Validation contract and five
> ordered tasks are drafted; QA Lead refines evidence before execution. Phase 4
> is CLOSED (QA PASS, human 5/5). Archives under `../phases/phase-4/`.

## Phase 5 Goal

v0.1b mutation + workflow transitions (roadmap §v0.1b slicing note, Phase 5
cut). Exit: an agent can mutate issue fields, observe audited changes, move an
issue through allowed workflow edges, and see deterministic rejection details
for invalid transitions. Delivered as agent-runnable CLI increments.

Scope:

- `fid issue update` (title, priority, labels, spec_md, est_tokens, due_date,
  owner role/actor) with `activity_log` projection.
- `activity` effect op + `activity_log` projection discipline (**CF-005**).
- CF-007 write-path fixup: validate-before-append seam (**High**).
- `internal/workflow` hard-coded graphs (task/epic/subtask) + rule functions
  (`all_acs_verified`, `has_completion_report`, `concerns_addressed`,
  `all_child_tasks_done_or_cancelled`).
- `fid issue transition --to <state> --by <role>` with rich rejection errors.

Explicit deferrals (do not implement this phase): `fid activity list` read
command (v0.1e, PRD-OUT-017); `--epic`/`--inbox` reparent/promote (v0.2);
claim-gated transitions (v0.1d); AC/completion/concern write commands (v0.1c);
configurable workflows + escalation workflow (v0.2). Permission matrix stays
v0.1d — `--by` only needs `RequireKnownRole` this phase.

## Validation Contract Index

`../phases/phase-5/qa/validation-contract.md` (VAL-001..VAL-013):

- VAL-001 `fid issue update` mutates fields; invalid update leaves issue
  unchanged. (task-3)
- VAL-002 one update = one event with documented effects; rebuild parity.
  (task-3)
- VAL-003 owner role/actor change journaled + projected to `activity_log`.
  (task-3)
- VAL-004 `activity` op projects one row per change; sole writer; ref backstop
  holds. (tasks-2-3)
- VAL-005 rejected write appends zero journal lines + reads still succeed; CF-007
  closure. (task-1)
- VAL-006 allowed transition edge advances status + category for all kinds.
  (task-5)
- VAL-007 invalid edge rejected naming both states; zero journal lines. (task-5)
- VAL-008 `in_progress→dev_complete` enforces `all_acs_verified{developer}` +
  `has_completion_report`; unmet names conditions + AC IDs. (task-5)
- VAL-009 `qa_review→qa_pass` enforces `all_acs_verified{qa_lead}`. (task-5)
- VAL-010 epic `in_progress→wrapped` enforces `all_child_tasks_done_or_cancelled`;
  unmet names child IDs. (task-5)
- VAL-011 cancel from any non-cancelled state; rejected from cancelled. (task-5)
- VAL-012 `concerns_addressed` rule both states by unit test (no wired edge).
  (task-4)
- VAL-013 human dogfood: rich errors actionable + journal stays clean. (phase-end)

PRD-OUT-017 has **no Phase 5 assertion by design** — the activity read command
is v0.1e; Phase 5 proves only the projection write-side (VAL-004).

## Ordered Tasks

1. **Validate-before-append seam (CF-007)** — pre-append projection-dry-run in
   `internal/journal`; route `internal/issues.Create` through it. VAL-005.
   *Foundational: Tasks 3 and 5 build update/transition on this seam.*
2. **`activity` effect op + projection discipline (CF-005 half)** — add the
   `activity` op to projection; §14.6 enforced app-layer (sole writer; §14.2
   backstop holds; SQL guard stays deferred per DEC-006). VAL-004.
3. **`fid issue update`** — mutate scalar fields + labels; emit update + activity
   effects; first real activity emitter. VAL-001/002/003 (+005).
4. **`internal/workflow` package** — hard-coded graphs + 4 rule functions, pure
   and unit-tested. VAL-012 (+ advances 006/008/009/010).
5. **`fid issue transition`** — wire workflow to projection; rich rejection
   errors; status+workflow_id preconditions + activity effect; rule evidence via
   journal fixtures. Closes CF-005. VAL-006..011 (+005).

## Planning Reminders

- DEC-007 carry-in: package-level tests for COV-001 packages
  (`internal/issues`, `internal/journal`, `internal/policy`, **new
  `internal/workflow`** ≥ 80%); hold COV-002 `cmd/fid` ≥ 70%.
- Phase 4 deep-retro lesson (now baked into VAL-005/VAL-007): every write path
  rejectable by a projection-time trigger must assert journal hygiene
  (line-count-unchanged + read-recovery), not only an actionable CLI error.
- Transition-rule evidence seeds AC/verification/completion/concern/child rows
  via journal fixtures until the v0.1c write commands ship.

## Carry-Forward Into Phase 5

| ID | Item | Disposition |
|---|---|---|
| CF-005 | `activity` effect op + `activity_log` projection discipline | Scheduled: Task 2 (op) + Tasks 3/5 (emitters); closes at Task 5. |
| CF-007 | Rejected write poisons journal tail (validate-before-append) | Scheduled: Task 1. **High — human-confirmed 2026-06-13.** |

CF-003 (role corrupt-tail tolerance) stays targeted before v0.1d.

## Next Checkpoint

QA Lead drafts `../phases/phase-5/qa/qa-plan.md` mapping every VAL-NNN to
deterministic `(qa-lead)` / mid-phase-review / `(human)` verification and fans
VAL-013 into per-observable human-checklist rows. Then Phase Runner starts the
Task 1 → Task 5 dev/QA loop.
