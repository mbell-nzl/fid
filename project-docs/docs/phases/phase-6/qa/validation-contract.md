# Phase 6 Validation Contract - Ready Queue + v0.1b Closure

> Phase: 6 - Ready queue + v0.1b closure (v0.1b)
> Author: Tech Lead (draft)
> Last refined: 2026-06-14 by Tech Lead (awaiting QA Lead evidence refinement)
>
> Location: `project-docs/docs/phases/phase-6/qa/validation-contract.md`.

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 6 is complete. The contract specifies **what must be proven**; the QA
plan specifies **how**.

If a required assertion cannot be satisfied, stop and escalate rather than
weakening it.

## Phase scoping note — "blocked / actionable" and empty-state contract

v0.1b ships **no dependency (`blocked_by`) links and no claims** — both are
v0.1d (GLOSSARY "Sequential epic"; roadmap §v0.1d). Phase 6 therefore defines
readiness from workflow **status category** alone (proposed DEC-017, to confirm
with human at planning):

- **Actionable / ready** = a task whose status category is `unstarted` (`todo`).
- **Unfinished** = status category is not `completed` and not `cancelled`.
- A sequential epic's lowest-sequence unfinished task that sits in a `started`
  category (e.g. `in_progress`) is **not actionable**; the epic then yields no
  ready task and the empty result explains it (current task already in flight).
- Dependency-link blocking and claim blocking are explicitly out of scope until
  v0.1d; `fid ready --explain` is v0.2 (PRD-ISSUE-024).

Empty-state output contract (proposed DEC-018, to confirm with human): bare
`fid ready` that finds nothing ready exits **0**, emits `{"task": null}` on
stdout (frozen shape; Workflow Runner depends on it), and writes a `note:`
explanation to stderr. Empty is a valid answer, not an error.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 6, sequential from `VAL-001`.
- Carry-forward assertions are re-issued with a new `VAL-NNN` and a source note.
- Retired/superseded assertions stay in the table; IDs are never reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | PRD-ISSUE-008; design §6.3 | Bare `fid ready` returns the next actionable task as a single object in the frozen §6.3 JSON shape (`task` with `id`, `title`, `status`, `epic_id`, `epic_title`, `sequence`, `priority`, `est_tokens`). | CLI: seed a workspace via journal fixtures / create commands, run `fid ready --format json`, assert the returned object matches the §6.3 shape and names the expected task. | QA Lead | task-4 | A ready task exists and `fid ready` returns exactly that task in the §6.3 shape with all documented keys. | Output omits/renames a §6.3 key, returns more than one task, or returns the wrong task. | Draft | `--explain` diagnostic shape is v0.2 (PRD-ISSUE-024). |
| VAL-002 | PRD-ISSUE-009; GLOSSARY "Sequential epic" | For a sequential epic, `fid ready --epic <id>` returns the lowest-`sequence` unfinished task when that task is actionable (`todo`), never a higher-sequence task. | Unit tests in the selection engine over a sequential fixture (lowest-seq todo selected; higher-seq todo skipped); CLI confirmation with `fid ready --epic --format json`. | QA Lead | tasks-3-4 | The lowest-sequence unfinished actionable task is returned; a higher-sequence task is never chosen ahead of it. | A higher-sequence task is returned, or the lowest actionable task is skipped. | Draft | Dependency-link ordering is v0.1d; v0.1b uses `sequence` only. |
| VAL-003 | PRD-ISSUE-010; DEC-017 | For a sequential epic whose lowest-`sequence` unfinished task is not actionable (already `started`), `fid ready --epic <id>` returns no task and the result explains the blocker; when that same task is actionable it is returned. | Unit + CLI over two fixtures: (a) lowest-seq task `in_progress` → no task + explanatory `note`; (b) lowest-seq task `todo` → task returned. | QA Lead | tasks-3-4 | Blocked case returns no task with a `note` naming the blocking current task; actionable case returns the task. | Blocked case returns a (wrong) task, omits the explanation, or actionable case returns nothing. | Draft | "Blocked" = status-category in v0.1b; `blocked_by` link blocking is v0.1d. |
| VAL-004 | PRD-ISSUE-011; GLOSSARY "Prioritised epic" | For a prioritised epic, `fid ready --epic <id>` returns the highest-priority actionable (`todo`) task within the epic; a started or lower-priority task is not chosen ahead of it. | Unit tests over a prioritised fixture (mixed priorities + one started); CLI confirmation. | QA Lead | tasks-3-4 | The highest-priority actionable task is returned; started and lower-priority tasks are not returned ahead of it. | A lower-priority or started task is returned, or the highest-priority actionable task is skipped. | Draft | Priority tie-break ordering recorded in task spec; no new ranking semantics beyond PRD. |
| VAL-005 | PRD-ISSUE-012 | An unattached (no-epic) actionable task is a ready candidate in project-wide `fid ready`; an unattached task that is started/done/cancelled is excluded. | Unit + CLI: project-wide `fid ready` over a fixture with an unattached `todo` task and an unattached `done` task; assert the todo is eligible and the done is excluded. | QA Lead | tasks-3-4 | Unattached `todo` task is eligible for selection; unattached non-actionable task is excluded. | An unattached non-actionable task is selected, or an unattached actionable task is never eligible. | Draft | Reparenting / inbox promotion is v0.2 (PRD-ISSUE-025). |
| VAL-006 | PRD-ISSUE-012 | Project-wide `fid ready` (no `--epic`) returns the single top task across prioritised epics and unattached tasks, with each sequential epic contributing only its current actionable task. | Unit + CLI over a multi-epic fixture (≥1 sequential, ≥1 prioritised, ≥1 unattached) asserting the single deterministic winner and that sequential epics contribute only their current task. | QA Lead | tasks-3-4 | One deterministic top task returned across the project; sequential epics contribute at most their current actionable task. | More than one task returned, a non-current sequential task contributes, or selection is non-deterministic. | Draft | Cross-database/multi-project federation is later (roadmap v0.6+). |
| VAL-007 | PRD-ISSUE-008/010; design §`fid ready` empty state; DEC-018 | When nothing is ready, `fid ready` exits 0, emits `{"task": null}` on stdout, and writes a `note:` explanation to stderr; when something is ready it returns the task with no `note`. | CLI over two fixtures: (a) all tasks completed/cancelled/in-flight → exit 0, `{"task": null}`, stderr `note:`; (b) one ready task → task object, no stderr `note`. | QA Lead | task-4 | Empty case: exit 0, `{"task": null}`, explanatory stderr note. Non-empty case: task returned, no note. | Empty case errors out, omits the null/`note`, or non-empty case prints a spurious note. | Draft | Empty shape frozen for the Workflow Runner; do not change without a recorded decision. |
| VAL-008 | PRD-OUT-001–004; design §`fid ready` | `fid ready --format json\|toon\|markdown\|yaml` renders both the selected-task and empty states without error; JSON is the default when `--format` is omitted. | CLI: run all four formats for a ready fixture and for an empty fixture; assert each renders successfully and default (no flag) equals JSON. | QA Lead | task-4 | All four formats render selected and empty states; omitting `--format` yields JSON. | A format errors, produces malformed output, or default is not JSON. | Draft | Format engine itself is Phase 3; this proves `ready` dispatches through it. |
| VAL-009 | CF-009; DEC-016; design §`fid issue update` | `fid issue update` confirmation renders one changed field per line; a multi-field update prints multiple lines and a single-field update prints exactly one field line. | CLI: a multi-field update prints one line per changed field; a single-field update prints exactly one field line. Updated design example matches. | QA Lead | task-1 | Multi-field update → one line per changed field; single-field update → one field line; design doc example updated to match. | Changed fields rendered as one comma-joined line, or line count mismatches changed-field count. | Draft | Does not change which fields are updatable or the journal/activity audit. |
| VAL-010 | CF-010; QP STRUCT-002 (§ file-size) | `internal/issues/transition.go`, `internal/issues/update.go`, and `internal/journal/project.go` are each ≤400 non-test lines after a behavior-preserving split, and behavior is unchanged (full suite green; `fid rebuild` parity holds). | `wc -l` (or equivalent) on the three files post-split; full task gate green; a rebuild-parity test/CLI run before vs after shows identical projected state. | QA Lead | task-2 | Each of the three files ≤400 lines and the full suite + rebuild parity pass unchanged. | Any of the three files >400 lines, or any behavior/parity regression. | Draft | Refactor only — no new behavior, no public-interface or schema change. |
| VAL-011 | Phase brief; roadmap §v0.1b exit | An end-to-end v0.1b smoke over create → list → get → update → transition → ready against a real `.fid/` workspace produces consistent read-back at each step, and deleting the DB then `fid rebuild` reproduces identical projected state. | Integration test (and/or scripted CLI run) exercising the six commands in sequence with assertions at each step plus a post-rebuild parity assertion. | QA Lead | task-5 | The full sequence runs green with consistent read-back at each step and byte-identical post-rebuild state. | Any step diverges from expected read-back or rebuild parity fails. | Draft | Held-out scenario harness location is a separate scheduled commitment; this task only recommends. |
| VAL-012 | Phase brief; roadmap §v0.1b exit | Human dogfood confirms `fid ready` returns the scheduling answer a human expects across a sequential epic, a prioritised epic, and an empty/all-blocked case in a real workspace. | Human checklist run against a real `.fid/` workspace; record outcomes per observable. | Human | phase-end | Human confirms the selected task (and empty explanations) match expectation for sequential, prioritised, and empty cases. | Human finds a selection that contradicts expectation or an unclear/incorrect empty explanation. | Draft | QA Lead fans this into one human-checklist row per observable (sequential / prioritised / empty). |

## Coverage Rules

Confirmed before task execution begins:

- Phase-scope requirements each map to ≥1 `VAL-NNN`: PRD-ISSUE-008→001/007,
  009→002, 010→003/007, 011→004, 012→005/006; PRD-OUT-001–004 (formats on
  `ready`)→008; CF-009/DEC-016→009; CF-010/STRUCT-002→010; v0.1b cross-command
  exit→011; human scheduling-feel→012.
- PRD-ISSUE-024 (`fid ready --explain`) and PRD-ISSUE-025 (reparent / inbox
  promote) have **no Phase 6 assertion by design** — both are v0.2; the phase
  brief records the exclusion.
- Every `VAL-NNN` is referenced by a task spec or marked `phase-end` (VAL-012).
- Binary/conditional assertions name both states — VAL-002, 003, 004, 005, 007,
  008, 009.
- DEC-017 (status-category readiness) and DEC-018 (empty-state contract) are
  planning decisions to confirm with the human; if either is rejected, the
  affected assertions (VAL-003, VAL-007) must be re-drafted before execution.
- QA Lead maps each `VAL-NNN` to deterministic `(qa-lead)`, mid-phase-review, or
  `(human)` verification and fans VAL-012 into per-observable human rows.

## Outcomes

Phase-end final review records every `VAL-NNN` outcome: `PASS`, `FAIL`,
`BLOCKED`, or `CARRY-FORWARD`. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-06-14 | VAL-001..012 | created | Phase 6 planning draft | phase brief; CF-009/DEC-016; CF-010; PRD-ISSUE-008..012 |
