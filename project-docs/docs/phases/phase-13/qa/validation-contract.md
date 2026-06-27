# Phase 13 Validation Contract - Labels + v0.1d Closure

> Phase: 13 - Labels + v0.1d closure (roadmap v0.1d, slice 4 of 4)
> Author: Tech Lead (draft)
> Last refined: 2026-06-27 by Tech Lead (planning draft)
>
> Location: `project-docs/docs/phases/phase-13/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 13 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase Scoping Notes

Phase 13 is the fourth and final v0.1d slice and the v0.1 lifecycle closure. It
ships the label surface and the end-to-end lifecycle smoke that is the headline
v0.1 exit criterion:

- `fid label add` / `fid label list`: additive label CRUD (the first
  label-creation surface; before this phase a label cannot be created, so
  `requireExistingLabels` on issue create/update can never be satisfied).
- `fid issue label add` / `fid issue label remove`: per-issue assignment that
  adds or removes **one** label without replacing the full set
  (`fid issue update --labels` stays replace-all and is unchanged).
- End-to-end v0.1d closure smoke over the full gated lifecycle: create -> claim
  -> transition -> complete -> unclaim with permission enforcement, links, and
  labels.

The `labels` and `issue_labels` tables and their FKs already exist
(`migrations/0002`; `03_TRD_SCHEMA.md` section 9.6). The projector applies
insert/delete effects generically against validated table identifiers
(`internal/journal/project_effects.go`), so the new `labels` insert and
`issue_labels` insert/delete effects need **no projector change and no schema
migration**. The `label.add`, `issue.label.add`, and `issue.label.remove` rows
already appear in `03_TRD_SCHEMA.md` section 13; Tasks 1-2 append those rows to
the in-code matrix that Phase 10 intentionally kept limited to implemented
actions (DEC-026 pattern).

Three planning interpretations apply (surfaced to the human for confirmation;
to be recorded as phase decisions during execution):

- **PI-1 (`fid label add` semantics).** `labels.name` is the primary key, so a
  second `fid label add` for an existing name is rejected before append with a
  non-zero exit and zero new journal lines/rows. `--colour` and `--description`
  are optional and stored as given; any colour-format validation is deferred and
  the exact stored/validated behaviour is recorded in the Task 1 completion
  report. Label creation does **not** project to `activity_log` (it matches the
  `issue create` no-activity-row precedent the roadmap calls out; CF-011 broad
  reconciliation stays v0.1e). The write emits one `label.added` event with one
  `labels` insert effect.
- **PI-2 (`fid issue label add`/`remove` semantics).** `fid issue label add`
  requires the label to already exist (FK plus an app-layer precheck mirroring
  `requireExistingLabels`) and rejects adding a label already present on the
  issue; `fid issue label remove` rejects removing a label not present on the
  issue. Both validate before append (CF-007 seam): non-zero exit, actionable
  error, zero new journal lines on rejection. Each successful add/remove emits
  one event that mutates exactly one `issue_labels` row and projects exactly one
  `activity_log` row recording the delta (PRD-ISSUE-023), mirroring Phase 12's
  link-removal activity pattern (DEC-031). Existing labels on the issue are not
  disturbed. `fid issue update --labels` remains replace-all and is out of
  scope.
- **PI-3 (file placement / CF-014).** The new `fid issue label` subcommands live
  in a new `cmd/fid/issue_label.go`, not appended to the 984-line
  `cmd/fid/issue.go`. CF-014's full behaviour-preserving split of `issue.go` /
  `replay.go` stays deferred to v0.1e; Phase 13 must not worsen those files.

Carry-forward review at kickoff:

- CF-008 remains deferred to v0.2.
- CF-011 remains deferred to v0.1e for the broad reconciliation, but
  PRD-ISSUE-023 explicitly requires issue label add/remove to project to
  `activity_log`; Phase 13 implements that per-command projection now (see
  PI-2), exactly as Phase 12 did for link removal. `fid label add` does **not**
  project to `activity_log`.
- CF-013 remains deferred to v0.1e.
- CF-014 remains deferred to v0.1e; Phase 13 avoids growing `issue.go` (PI-3).
- Lock-ordering lesson (DEC-028): every new write command (`label add`,
  `issue label add`, `issue label remove`) refreshes the projection with the
  held-lock freshness helper inside `journal.WithLock`; no nested
  `.fid/write.lock` acquisition.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 13, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | TRD-ARCH-010; `03_TRD_SCHEMA.md` section 13; roadmap v0.1d Phase 13 | The in-code permissions matrix contains exactly one row each for `label.add`, `issue.label.add`, and `issue.label.remove` with the same role booleans as `03_TRD_SCHEMA.md` section 13 (all three: developer + tl + workflow_runner + human), and `fid policy show` displays those rows from the same source `Authorize` enforces. | Unit + CLI: matrix completeness test updated for the three actions; `Authorize` allowed/denied pairs; `fid policy show` output includes the rows and matches the matrix object. | QA Lead | tasks-1-2 | All three actions have exactly one in-code row; allowed/denied role checks match section 13; policy show matches the matrix. | A row is missing/duplicated, booleans diverge from section 13, or show/enforce drift. | Draft | Configurable role permissions remain v0.2 (PRD-ROLE-006). |
| VAL-002 | PRD-ISSUE-015; PRD-ISSUE-023; `06_DESIGN_PLAN_CLI.md` section 9 `fid label add`; PI-1 | `fid label add <name> [--colour <hex>] [--description <text>] --by <role>` creates one `labels` row (name + optional colour/description) in exactly one journal event; a second add of an existing name is rejected before append with a non-zero exit and zero new journal lines/rows. | Domain + CLI: add a label, read the row back, assert one event; re-add the same name and compare journal line count / DB state before and after. | QA Lead | task-1 | First add creates the row with the given fields and one event; duplicate-name add rejects before append leaving journal/DB unchanged. | A field is dropped, more than one event is emitted, or a duplicate-name add appends/projects or exits zero. | Draft | Colour-format validation is not promised; behaviour recorded in the completion report. |
| VAL-003 | PRD-OUT-001..004; `06_DESIGN_PLAN_CLI.md` section 9 `fid label list` | `fid label list` returns every label with `name`, `colour`, and `description`, and the json, toon, markdown, and yaml formats carry the same data. | CLI: seed several labels (with and without colour/description); run list in all four formats; compare parsed structures to the JSON reference. | QA Lead | task-1 | All four formats expose the same label set and fields, including null colour/description. | A field is missing, a label is dropped, or a format diverges from JSON. | Draft | Filtering label list is not specified for v0.1. |
| VAL-004 | TRD-ARCH-010; `03_TRD_SCHEMA.md` section 13; CF-007 | `label.add` is authorized only for roles permitted in section 13 (developer, tl, workflow_runner, human); a denied role and an unknown role are rejected before append with an error naming role + action and zero new journal lines. | CLI: add as developer (allowed) and as a denied role (e.g. qa_lead) and an unknown role; assert allow path appends one event and deny paths append nothing. | QA Lead | task-1 | Permitted roles create the label; denied/unknown roles reject before append naming role + `label.add`. | A denied/unknown role succeeds or appends, or a permitted role is rejected. | Draft | Configurable role permissions are v0.2 (PRD-ROLE-006). |
| VAL-005 | PRD-ISSUE-023; `06_DESIGN_PLAN_CLI.md` `fid issue label add`; PI-2 | `fid issue label add <issue-id> <label> --by <role>` adds exactly that one label to the issue in one journal event **without replacing the issue's other labels**; adding a label that does not exist and adding a label already present on the issue are each rejected before append with a non-zero exit and zero new journal lines. | Domain + CLI: give an issue two labels, add a third, assert all three present and the prior two untouched in one event; add a nonexistent label and an already-present label and assert zero append each. | QA Lead | task-2 | The single label is added, the existing labels remain, exactly one event is emitted; nonexistent-label and already-present add both reject before append. | An existing label is dropped/replaced, more than one `issue_labels` row changes, or a reject case appends/exits zero. | Draft | `fid issue update --labels` replace-all is unchanged and out of scope. |
| VAL-006 | PRD-ISSUE-023; `06_DESIGN_PLAN_CLI.md` `fid issue label remove`; PI-2 | `fid issue label remove <issue-id> <label> --by <role>` removes exactly that one label in one journal event **without disturbing the issue's other labels**; removing a label not present on the issue is rejected before append with a non-zero exit and zero new journal lines. | Domain + CLI: give an issue two labels, remove one, assert only the other remains in one event; remove a label not on the issue and assert zero append. | QA Lead | task-2 | The single label is removed, the remaining labels are untouched, exactly one event is emitted; removing an absent label rejects before append. | The wrong label is removed, another label is disturbed, or removing an absent label appends/exits zero. | Draft | Removing a label does not delete it from the `labels` table. |
| VAL-007 | PRD-ISSUE-023; CF-011; DEC-031 (Phase 12 activity pattern) | Each successful `fid issue label add` and `fid issue label remove` projects exactly one `activity_log` row that records the issue, the label, and whether it was added or removed; `fid label add` projects no `activity_log` row. | Domain + CLI: after an add and a remove, query `activity_log` for the issue and assert one add row and one remove row with the label recorded; after a `fid label add`, assert no new `activity_log` row. | QA Lead | task-2 | Add and remove each write one labelled `activity_log` row identifying issue + label + direction; label creation writes none. | An activity row is missing, duplicated, mislabels the direction, or `fid label add` writes one. | Draft | CF-011 broad activity-log reconciliation across all writes stays v0.1e. |
| VAL-008 | TRD-ARCH-010; `03_TRD_SCHEMA.md` section 13; CF-007 | `issue.label.add` and `issue.label.remove` are authorized only for roles permitted in section 13 (developer, tl, workflow_runner, human); a denied role and an unknown role are rejected before append naming role + action with zero new journal lines. | CLI: add/remove as developer (allowed) and as a denied role (e.g. qa_reviewer) and an unknown role; assert allow path appends one event and deny paths append nothing. | QA Lead | task-2 | Permitted roles add/remove; denied/unknown roles reject before append naming role + the label action. | A denied/unknown role succeeds or appends, or a permitted role is rejected. | Draft | Configurable role permissions are v0.2 (PRD-ROLE-006). |
| VAL-009 | Roadmap v0.1d Phase 13 exit; v0.1 headline exit criterion; PRD-NFR-008; ADR-0001; DEC-028 | The full v0.1d lifecycle runs through real commands against a temp workspace and stays rebuild-stable: create -> claim -> transition into a started-category state (requires a live caller claim) -> complete -> unclaim, with permission enforcement on every write, plus a typed link and a label assignment on the issue; each write appends exactly one event under the held lock with no nested-lock failure, and `fid rebuild` reproduces identical issue/status/claim/link/label state. | Integration test (`go test -tags integration`): scripted command sequence following the `e2e_v01d_link_test.go` pattern; assert per-write event counts, no nested-lock error, a denied-role rejection somewhere in the flow, and rebuild parity of issues/claims/links/issue_labels. | QA Lead | task-3 | The full sequence runs, each write appends one event with no nested-lock failure, a denied write is rejected before append, and rebuild reproduces identical state. | Any write fails under the lock, emits the wrong event count, a permission gate is missing, or rebuild diverges. | Draft | This is closure, not new behaviour: claim/transition/complete already shipped in Phases 5/8/11. |
| VAL-010 | `08_QUALITY_PLAN.md` sections 5.2, 5.6; phase goal | A human dogfoods a real `.fid` workspace and can tell from CLI output that a label was created, that adding a label to an issue left its other labels in place, that removing a label left the rest in place, and that the full create -> claim -> transition -> complete -> unclaim lifecycle ran end-to-end. | Human checklist with one row per observable; QA Lead generates from this VAL after deterministic QA passes. | Human | phase-end | Human confirms label creation, additive add, non-destructive remove, and the end-to-end lifecycle from CLI output without reading implementation. | Human cannot identify those outcomes from CLI output, or the checklist is not run. | Draft | QA Lead fans this into per-observable human-checklist rows. |

## Coverage Rules

Confirmed for Phase 13 execution:

- Phase-scope requirements each map to at least one `VAL-NNN`: PRD-ISSUE-015 to
  VAL-002/003; PRD-ISSUE-023 to VAL-005/006/007; TRD-ARCH-010 to
  VAL-001/004/008; roadmap v0.1d Phase 13 exit to VAL-009; phase dogfood goal to
  VAL-010.
- Every `VAL-NNN` is referenced by a task spec or explicitly marked `phase-end`
  (VAL-010).
- Binary/conditional assertions name both states: VAL-002, VAL-004, VAL-005,
  VAL-006, VAL-007, VAL-008, VAL-009.
- VAL-001 and VAL-003 carry structural/parity conformance proofs
  (matrix-vs-section-13; four-format parity-vs-JSON).
- Coverage is enforced at the gate, not as a `VAL-NNN`: COV-001 domain package
  coverage and COV-002 `cmd/fid` coverage must stay above thresholds. Each task
  runs the required coverage gate and records package percentages.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase review, or
  `(human)` verification, and fans VAL-010 into per-observable human rows.
- Standing QA-plan checks carry forward: show-vs-enforce parity for the new
  matrix rows (VAL-001), validate-before-append on every rejection
  (VAL-002/004/005/006/008), per-command `activity_log` projection on issue
  label changes (VAL-007), and no stale-projection nested-lock failure on the
  new write paths (VAL-009).

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-27 | VAL-001..010 | created | Phase 13 planning draft | roadmap v0.1d Phase 13; PRD-ISSUE-015; PRD-ISSUE-023; TRD-ARCH-010 |
