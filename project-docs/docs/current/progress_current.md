# Progress - Current

> **Purpose.** In-flight phase status only. Loaded by default in every agent
> session. Historical phase progress lives under `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 5 (v0.1b Mutation + Workflow Transitions) is **CLOSED** on
branch `phase-5`. QA final review VERDICT: **PASS WITH IMPROVEMENTS** (0 blocking
/ 0 major / 1 improvement); human dogfood checklist 6/6 PASS. VAL-001..VAL-013
all PASS. Phase 4 and earlier are archived under `../phases/phase-N/`.

Phase 5 delivered:

- Validate-before-append seam in `internal/journal`; `issues.Create` routed
  through it (Tasks 1/1b, **CF-007 closed**, VAL-005).
- `activity` effect op + sole-writer `activity_log` projection discipline (Task
  2, **CF-005 closed**, VAL-004; SQL session guard stays deferred per DEC-006).
- `fid issue update` (scalars + labels + owner role/actor) with one-event /
  activity-effect audit and rebuild parity (Task 3, VAL-001/002/003).
- `internal/workflow` hard-coded graphs + 4 rule functions (Task 4, VAL-012,
  95.5% coverage).
- `fid issue transition` with rich rejection errors and rule-evidence fixtures
  (Task 5, VAL-006..011).

Five tasks + one remediation (Task 1b), one BLOCK→resolved (ESC-003), zero spec
deviations, zero dependency changes.

**Human-raised at wrap (non-blocking):**

- DEC-015 / **CF-008**: `todo → done` rejection is correct; "close without
  completing" is already `todo → cancelled` (PRD-WF-013). A *resolution
  type/reason* (won't-complete, duplicate, …) is a real but separate v0.2
  product enhancement — captured, not actioned. Suggest dogfooding
  `todo → cancelled` to confirm it serves the redundant-backlog case.
- DEC-016 / **CF-009**: `fid issue update` should render one changed field per
  line instead of one comma list. Deferred CLI polish.
- **CF-010**: split oversized non-test files (`transition.go` 426,
  `update.go` 551, `project.go` 524) on next adjacent change.

Open carry-forward register: CF-003 (role corrupt-tail, before v0.1d), CF-008,
CF-009, CF-010.

## Next Action

- Plan **Phase 6** (v0.1b closure): `fid ready` (sequential/prioritised epics,
  unattached tasks, cross-project, empty-result explanations, format dispatch)
  plus an end-to-end v0.1b smoke over create/list/get/update/transition/ready.
  Draft `../phases/phase-6/qa/validation-contract.md` and the ordered task list.
