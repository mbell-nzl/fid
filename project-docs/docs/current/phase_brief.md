# Phase Brief - Phase 6 (Ready Queue + v0.1b Closure)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 6 is **NOT YET PLANNED**. Phase 5 (Mutation + Workflow
> Transitions) is **CLOSED** (QA PASS WITH IMPROVEMENTS; human dogfood 6/6 PASS;
> VAL-001..VAL-013 all PASS). Archives under `../phases/phase-5/`. Tech Lead
> drafts the Phase 6 validation contract and task list next.

## Phase 6 Goal (from roadmap §v0.1b, Phase 6 cut)

`fid ready` + v0.1b closure. **Exit:** the workflow runner can ask for the next
actionable task and receive the same scheduling answer a human expects.

Scope (roadmap §10):

- `fid ready` default mode: sequential epics, prioritised epics, unattached
  tasks, cross-project selection, empty-result explanations, format dispatch.
- End-to-end v0.1b smoke over create / list / get / update / transition / ready.

Explicit deferrals (do not implement this phase): `fid ready --explain`
(v0.2, PRD-ISSUE-024); task reparenting / inbox promotion (v0.2, PRD-ISSUE-025);
label CRUD (later label/links slice — use existing `labels` rows). v0.1c
AC/completion/concern write commands and v0.1d claims/permissions/links remain
out of scope.

## Carry-Forward Into Phase 6

Tech Lead reviews the register in `decisions_active.md` at phase start. Open
items to schedule, re-defer, or cancel:

| ID | Item | Note |
|---|---|---|
| CF-003 | Role-validation corrupt-tail tolerance | Target before v0.1d; likely not Phase 6. |
| CF-008 | Resolution type / cancellation reason (DEC-015) | Product decision, v0.2; do not widen v0.1 graphs. |
| CF-009 | `fid issue update` one-field-per-line output (DEC-016) | Low CLI polish; fold into next update-adjacent change. |
| CF-010 | Non-test file-size split (transition/update/project) | Fold into next adjacent change to those files. |

CF-004/CF-005/CF-006/CF-007 are closed (see register footnotes).

## Planning Reminders

- DEC-007 carry-in: package-level COV-001 (`internal/issues`,
  `internal/journal`, `internal/policy`, `internal/workflow` ≥ 80%); hold
  COV-002 `cmd/fid` ≥ 70%. Phase 5 ended at `cmd/fid` 79.0%, `internal/issues`
  82.8%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow`
  95.5%.
- Phase 4/5 deep-retro lesson (write-command slicing checklist): when an
  acceptance criterion asserts *recovery* or *journal hygiene* for a rejected
  write, the task's **original** evidence must demonstrate that state — "no bad
  append" and "system still rebuilds/reads" are distinct observables.
- `fid ready` is read-only scheduling; the human dogfood remains the place where
  product-feel questions surface (DEC-015 came from the Phase 5 dogfood).
- This phase closes v0.1b: plan the cross-command end-to-end smoke as an explicit
  task and decide whether it becomes the first held-out scenario harness input.

## Next Checkpoint

Tech Lead: read the roadmap §v0.1b Phase 6 entry and the `fid ready` /
scheduling PRD-ISSUE/PRD-OUT requirements, create
`../phases/phase-6/qa/validation-contract.md`, break the phase into ordered
one-session tasks, refresh this brief with the VAL index, and present the
contract + task list to the human.
