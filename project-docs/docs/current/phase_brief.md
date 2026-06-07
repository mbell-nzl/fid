# Phase Brief - Phase 4 Planning

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.
>
> **State:** Phase 4 (v0.1b Issue authoring and discovery) is **PLANNED -
> PENDING REVIEW**. Implementation must not begin until the human approves this
> validation contract/task list and QA Lead writes `qa/qa-plan.md`.

## Phase Goal

Phase 4 is the first v0.1b slice. It makes issue authoring and discovery usable
end to end: create epics/tasks/subtasks with scheduling and owner metadata,
enforce the sequence scheduling backstop, centralise enriched issue reads, add
structured `fid issue get --section`, and add `fid issue list` filters/order
with all read formats.

## Task Index

Task specs live under `../phases/phase-4/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Sequence Scheduling Backstop | Pending | None |
| 2 | Issue Authoring Domain Service | Pending | Task 1 |
| 3 | `fid issue create` CLI Surface | Pending | Tasks 1-2 |
| 4 | Centralised Issue Detail Read Model + Sections | Pending | Tasks 1-3 |
| 5 | `fid issue list` Filters, Ordering, And Formats | Pending | Tasks 1-4 |

## Validation Contract

Full contract: `../phases/phase-4/qa/validation-contract.md`.

| VAL ID | Source | Owning Task/Window | Validator | Status |
|---|---|---|---|---|
| VAL-001 | CF-004; PRD-ISSUE-013/014; schema 14.1 | task-1, task-3 | QA Lead | Draft |
| VAL-002 | PRD-ISSUE-001..004; PRD-WF-005/006 | tasks 2-3 | QA Lead | Draft |
| VAL-003 | PRD-ISSUE-015/016/019/020 | tasks 2-3 | QA Lead | Draft |
| VAL-004 | Epic metadata; schema 2.1 | tasks 2-3 | QA Lead | Draft |
| VAL-005 | PRD-OPS-001/003/004; TRD-ARCH-001/002 | tasks 2-3 | QA Lead | Draft |
| VAL-006 | CF-006; PRD-ISSUE-006; TRD-ARCH-005 | task-4 | QA Lead | Draft |
| VAL-007 | PRD-ISSUE-018; CLI `--section` | task-4 | QA Lead | Draft |
| VAL-008 | PRD-ISSUE-007 list filters | task-5 | QA Lead | Draft |
| VAL-009 | PRD-ISSUE-007; PRD-OUT-004 list formats/order | task-5 | QA Lead | Draft |
| VAL-010 | Human dogfood authoring/discovery loop | phase-end | Human | Draft |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-4/qa/validation-contract.md` | Draft |
| QA plan | `../phases/phase-4/qa/qa-plan.md` | Missing |
| Mid-phase review | `../phases/phase-4/qa/review-3.md` | Recommended after Task 3 |
| Final review | `../phases/phase-4/qa/final-review.md` | Pending |
| Human checklist | `../phases/phase-4/qa/human-checklist.md` | Pending |

## Active Constraints

- Do not implement `fid issue update`, `fid issue transition`, `fid ready`,
  activity-log population, workflow transition rules, `--context-pack`, label
  CRUD, `fid ready --explain`, or task reparenting/inbox promotion in Phase 4.
- Label assignment/filtering uses existing `labels` rows only. Do not silently
  add label auto-create semantics.
- Epic metadata is scheduling-only in Phase 4; brief path and target dates need
  a recorded CLI/update surface before implementation.
- A task under an epic missing `epic_metadata` is invalid in Phase 4; create
  epics with scheduling metadata first.
- `fid issue get --section` means named structured sections for this phase:
  `acceptance`, `concerns`, `completions`, `links`, `claims`.
- Full markdown-body heading extraction from `spec_md_path` is out of scope
  unless the human requests a spec update before execution.
- Apply DEC-007 carry-in: package-level tests for COV-001 packages
  (`internal/issues`, `internal/journal`, `internal/policy`); hold COV-002
  `cmd/fid` >= 70%.

## Interfaces / Subsystems In Scope

- `migrations/` and `internal/db/` - sequence scheduling backstop.
- `internal/issues/` - issue authoring, detail read model, list query.
- `internal/workflow/` - default workflow/status/category constants only if
  useful; transition validation waits for Phase 5.
- `internal/format/` - markdown list renderer if needed; reuse generic
  JSON/TOON/YAML format path.
- `cmd/fid/issue.go` - thin Cobra adapter for create/get/list.

## Known Risks

- **`--section` ambiguity.** PRD wording mentions markdown parsing while CLI
  names structured sections. Phase 4 implements structured sections only; if
  that is rejected, Task 4 must be re-scoped before execution.
- **Label semantics.** Label CRUD is later roadmap scope. Phase 4 may test
  label assignment/filtering with seeded existing labels but must not invent
  auto-create behavior.
- **Task 2 size.** Domain authoring spans several fields. Keep transition,
  update, ready, and label CRUD out to preserve one-session sizing.

## Success / Exit Criteria

- [ ] All five tasks complete with green required task gates.
- [ ] Phase-level QA final review has no unresolved blocking findings.
- [ ] VAL-001..VAL-010 have evidenced outcomes or accepted dispositions.
- [ ] Human checklist accepts the create/get/section/list dogfood loop or records
      follow-up work.

## Next Planning Checkpoint

Phase 5 should plan v0.1b mutation and workflow transitions: `fid issue update`,
activity-log projection discipline / activity effects (CF-005), hard-coded
workflow graphs/rules, and `fid issue transition`.
