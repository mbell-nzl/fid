# Phase 5 Validation Contract - Mutation + Workflow Transitions

> Phase: 5 - Mutation + workflow transitions (v0.1b)
> Author: Tech Lead (draft)
> Last refined: 2026-06-13 by Tech Lead (awaiting QA Lead evidence refinement)
>
> Location: `project-docs/docs/phases/phase-5/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 5 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 5, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | PRD-ISSUE-015; CLI §`fid issue update` | `fid issue update` mutates allowed scalar fields (title, priority, labels, spec_md, est_tokens, due_date) on an existing issue; updated fields read back via `fid issue get` while untouched fields are preserved; an invalid update (e.g. priority out of 0–5, blank title, unknown label) is rejected and leaves the issue unchanged. | CLI: create issue, `fid issue update` a subset of fields, `fid issue get --format json` shows new values + preserved values; separate run with invalid value asserts non-zero exit + unchanged read-back. | QA Lead | task-3 | Valid update reflected in read-back with other fields unchanged; invalid update returns actionable non-zero error and read-back is identical to pre-attempt state. | A valid update is not reflected, silently changes an unrequested field, or an invalid update mutates the issue. | Draft | `--epic`/`--inbox` reparent/promote is v0.2 (PRD-ISSUE-025); not in scope. |
| VAL-002 | TRD-ARCH-004; TRD §4.3 | One `fid issue update` produces exactly one journal event carrying the `update` effect on `issues` (plus label effects when labels change) and an `activity` effect; deleting the SQLite DB and `fid rebuild` reproduces identical projected state. | Inspect the appended JSONL line count (+1) and effects; `fid rebuild` then `fid issue get --format json` parity against pre-rebuild output. | QA Lead | task-3 | Exactly one event per update with the documented effects; post-rebuild read-back is byte-identical to pre-rebuild. | More than one event per update, missing activity effect, or rebuild diverges. | Draft | Composite commands (`fid dev finish`) are later versions. |
| VAL-003 | PRD-ISSUE-021 | Owner role and owner actor are changed via `fid issue update --owner-role/--owner-actor`; the change is journaled and projected to `activity_log` with actor and before/after detail. | CLI update of owner fields; `fid issue get` read-back; package/DB query of the resulting `activity_log` row (no public `fid activity list` until v0.1e). | QA Lead | task-3 | Owner fields read back changed and one `activity_log` row records actor + from/to; unknown owner role rejected. | Owner change not persisted, or no `activity_log` row recorded, or unknown role accepted. | Draft | `fid activity list` read command is v0.1e (PRD-OUT-017). |
| VAL-004 | TRD SCHEMA §14.6; TRD §4.3 `activity` op; CF-005 | The `activity` effect op projects exactly one `activity_log` row per audited change, and `activity_log` is populated only through the projection path (no independent CLI write surface); the §14.2 polymorphic-ref backstop still rejects an `activity_log` row whose `(ref_type, ref_id)` does not resolve. | Package test in `internal/journal` projecting an event with an `activity` effect and asserting one row; rebuild determinism test; raw-SQL test that an unresolved polymorphic ref is rejected. | QA Lead | tasks-2-3 | activity op inserts exactly one row with resolved ref; rebuild reproduces it; unresolved ref insert rejected. | activity op inserts wrong count, bypasses ref integrity, or a non-projection path writes `activity_log`. | Draft | SQL session-scoped projection guard remains deferred (DEC-006); discipline is application-layer via the activity op being the sole writer. |
| VAL-005 | CF-007; DEC-014; TRD §4.3 | A write rejected at validation or projection time (invalid create, invalid update, or invalid transition) appends **zero** lines to the canonical journal, and all subsequent reads succeed (no poisoned tail); an accepted write appends exactly one line and projects cleanly. | Regression test: journal line count unchanged after a rejected write; poisoned-tail read-recovery test; contrast test that an accepted write appends exactly one line. CLI repro of the Phase 4 invalid-create scenario followed by a successful unrelated `fid issue get`. | QA Lead | task-1 | Rejected write leaves journal line count unchanged and later reads succeed; accepted write appends exactly one line. | A rejected write appends a journal line, or a later read fails because of a rejected write. | Draft | Replaces the Phase 4 VAL-005 coverage gap noted in `phase-4/decisions.md`. |
| VAL-006 | PRD-WF-001; PRD-WF-005; PRD-WF-006; TRD §4.6 | `fid issue transition --to <state> --by <role>` advances an issue along an allowed edge of its kind's default workflow (task/epic/subtask), and `fid issue get` reads back the new status with the correct status category. | CLI: transition across an allowed edge for each kind; `fid issue get --format json` shows new `status` + `status_category`. | QA Lead | task-5 | Allowed edge transitions and read-back shows the expected status and category for each default workflow. | An allowed edge is rejected, status not updated, or wrong category mapping. | Draft | Configurable workflows (PRD-WF-011/012) are v0.2. |
| VAL-007 | PRD-WF-002; CF-007 | A transition along an edge not present in the workflow is rejected with an error naming the current state and the requested state, and appends zero journal lines. | CLI: request a disallowed edge; assert non-zero exit, error text names current + requested state; journal line count unchanged. | QA Lead | task-5 | Disallowed edge rejected with both states named and no journal append. | Disallowed edge accepted, or error omits a state, or a journal line is appended. | Draft | — |
| VAL-008 | PRD-WF-003; PRD-WF-004; PRD-WF-007; PRD-WF-008 | The `in_progress → dev_complete` task edge evaluates `all_acs_verified{developer}` and `has_completion_report` before allowing the transition; when unmet, the rejection names the specific unmet conditions including AC IDs; when satisfied (via seeded journal fixtures), the transition succeeds. | CLI/fixture: seed ACs/verifications/completion via journal fixtures; rejected run names unmet rules + AC IDs; satisfied run transitions. | QA Lead | task-5 | Both rules evaluated; unmet case names conditions + entity IDs; satisfied case transitions. | Transition allowed with unmet rules, or rejection omits the unmet condition/IDs. | Draft | v0.1c write commands for ACs/completions not yet shipped; evidence uses journal fixtures. |
| VAL-009 | PRD-WF-007 | The `qa_review → qa_pass` task edge evaluates `all_acs_verified{qa_lead}` (rule parameterised by a different role set than VAL-008); unmet rejects naming qa_lead verification gaps, satisfied transitions. | CLI/fixture both states with qa_lead-role verifications seeded via journal fixtures. | QA Lead | task-5 | Rule honours the `qa_lead` role parameter: unmet rejects, satisfied transitions. | Rule ignores the role parameter or passes/fails incorrectly. | Draft | — |
| VAL-010 | PRD-WF-010; PRD-WF-004 | The epic `in_progress → wrapped` edge evaluates `all_child_tasks_done_or_cancelled`; when an open child task remains, the rejection names blocking child IDs; when all children are done/cancelled, the transition succeeds. | CLI/fixture: epic with child tasks; rejected run names blocking child IDs; satisfied run transitions. | QA Lead | task-5 | Rule evaluated; unmet names blocking child IDs; satisfied transitions. | Transition allowed with an open child, or rejection omits child IDs. | Draft | — |
| VAL-011 | PRD-WF-013 | Any issue in a non-cancelled state can transition to `cancelled`; a transition to `cancelled` from an already-cancelled state is rejected. | CLI: cancel from at least two distinct non-cancelled states; assert reject from cancelled. | QA Lead | task-5 | Cancel allowed from non-cancelled states; rejected from cancelled. | Cancel rejected from a valid non-cancelled state, or allowed from cancelled. | Draft | — |
| VAL-012 | PRD-WF-009 | The `concerns_addressed` rule function is implemented and returns satisfied iff the issue has no open concerns; both states proven by unit test (the default v0.1 workflows do not wire this rule to an edge). | `internal/workflow` unit test: open-concern fixture → unsatisfied; no-open-concern fixture → satisfied. | QA Lead | task-4 | Function returns correct result for both the open and no-open-concern states. | Function missing, or returns wrong result for either state. | Draft | No default-workflow edge references this rule in v0.1; CLI exercise deferred to a future workflow/version. |
| VAL-013 | Phase brief; Phase 4 deep-retro lesson | Human dogfood: a real mutate/transition session (update fields, advance allowed edges, attempt an invalid transition and a rejected update) confirms rich rejection errors are actionable and the canonical journal stays clean after rejected writes. | Human checklist run against a real `.fid/` workspace; record outcomes. | Human | phase-end | Human confirms errors are actionable and no rejected write poisons the journal or blocks later reads. | Human finds an unclear error, a missing rule explanation, or a rejected write that breaks later reads. | Draft | QA Lead fans this into one human-checklist row per observable. |

## Coverage Rules

Confirmed before task execution begins:

- Phase-scope PRD/TRD requirements each map to ≥1 `VAL-NNN`: PRD-ISSUE-015→001,
  PRD-ISSUE-021→003, PRD-WF-001→006, 002→007, 003→008, 004→008/010, 005→006,
  006→006, 007→008/009, 008→008, 009→012, 010→010, 013→011; TRD-ARCH-004/§4.3→
  002/005, §14.6/`activity` op→004, CF-007→005.
- PRD-OUT-017 (`fid activity list`) has **no Phase 5 assertion by design**: the
  public activity read command ships in v0.1e. Phase 5 builds and proves only
  the `activity_log` projection write-side (VAL-004); the phase brief records
  this exclusion.
- Every `VAL-NNN` is referenced by a task spec or marked `phase-end` (VAL-013).
- Binary/conditional assertions name both states (accept/reject) — see VAL-001,
  004, 005, 007, 008, 009, 010, 011, 012.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase-review, or
  `(human)` verification and fans VAL-013 into per-observable human rows.

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-13 | VAL-001..013 | created | Phase 5 planning draft | phase brief; CF-005; CF-007 / DEC-014 |
