# Progress - Current

> **Purpose.** In-flight phase status only. Loaded by default in every agent
> session. Historical phase progress lives under `../phases/phase-N/progress.md`.
>
> **Size target:** under 100 lines.

## Project State

**Status:** Phase 4 (v0.1b Issue authoring and discovery) is **PLANNED -
PENDING REVIEW** on branch `main`. Phase 3 is closed and archived under
`../phases/phase-3/`.

Planning completed in this pass:

- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` now records the accepted
  v0.1b three-phase split: Phase 4 authoring/discovery, Phase 5
  mutation/workflow transitions, Phase 6 ready queue/v0.1b closure.
- Phase 4 directories exist under `../phases/phase-4/`.
- Draft validation contract exists with VAL-001..VAL-010:
  `../phases/phase-4/qa/validation-contract.md`.
- Five Developer task specs exist:
  task 1 sequence scheduling backstop; task 2 issue authoring domain service;
  task 3 `fid issue create` CLI; task 4 centralised issue detail + `--section`;
  task 5 `fid issue list` filters/order/formats.

Carry-forward disposition for this plan:

- CF-004 scheduled in Phase 4 task 1 (and CLI evidence in task 3).
- CF-006 scheduled in Phase 4 task 4.
- CF-005 re-deferred to Phase 5, where `activity_log` population and
  `activity` effects belong with `fid issue update`.
- CF-003 remains later, before v0.1d role permission matrix / mutable role
  events.

Next action:

- Human reviews the Phase 4 validation contract and ordered task list. If
  accepted, QA Lead writes `qa/qa-plan.md`; implementation starts only after the
  QA plan is ready.
