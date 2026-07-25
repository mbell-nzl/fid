# Phase Brief - Phase 19 (Planned, execution gated)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session.
>
> **Size target:** under 150 lines.

## Phase Goal

Phase 18 made workflow *structure* data. Phase 19 makes it **authorable**:
`fid workflow define` writes user workflows as journal events into the same
projected workflow tables, `fid workflow assign` binds a workflow to an issue
kind or to the project, `fid workflow list` reads them back in every supported
format, and the `escalation-default` workflow (PRD-WF-005a) is seeded as
system reference data. Success is observable: an agent defines a custom
workflow, assigns it to a kind, lists it, and creates an issue that picks it
up instead of the built-in default. This is v0.2a slice 2 of 3.

## Status

**Planned. Awaiting Human approval; no other precondition.** DEC-083 lifted the
DEC-074 execution gate. The deep retrospective ran 2026-07-26
(`../reviews/2026-07-26-workflow-retrospective.md`); DEC-084..086 promote this
phase's two stand-ins to standing rules and CF-023 narrows to DEC-087. **No
task spec, assertion, or scope changed.** No branch yet; next action is Human
approval, then the QA plan.

## Operating Mode

**Heavy** (DEC-075). Three heavy entry criteria apply at once: two schema
migrations with a projection-schema version bump; three new public CLI
commands plus two permissions-matrix rows that change `fid policy show`
output; and two new journal event types that become permanent replay history.
Entry criteria and artifact matrix: `../../specs/09_AI_WORKFLOW.md` §4.
Heavy mode never carries QA feasibility validation forward.

## Task Index

Full task specs and completion reports live under `../phases/phase-19/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | `workflow_assignments` schema + `escalation-default` seed | Pending | None |
| 2 | Workflow authoring core (validation + effect builder) | Pending | 1 |
| 3 | `fid workflow define` + `workflow.defined` event | Pending | 1, 2 |
| 4 | `fid workflow assign` + assignment-aware issue creation | Pending | 1-3 |
| 5 | `fid workflow list` in every output format | Pending | 1, 3, 4 |
| 6 | Closure: integrity, rebuild parity, exit scenario, spec sync | Pending | 1-5 |

## Validation Contract

Full contract, with per-assertion sources, evidence methods, and pass/fail
criteria: `../phases/phase-19/qa/validation-contract.md`. All 26 assertions
are `Draft` until the QA plan refinement round. Compact index:

| Owning window | Assertions | Validator | Theme |
|---|---|---|---|
| task-1 | VAL-001..005 | QA Lead | Assignments table, escalation seed, reversibility, reserved rule, projection version |
| task-2 | VAL-006..008 | QA Lead | Deterministic effects, rejection matrix, cancel-edge materialization |
| task-3 | VAL-009..013 | QA Lead | Define event + projection, reject-before-append, permissions, arrow parity, output shape |
| task-4 | VAL-014..017 | QA Lead | Kind/project assignment, precedence, assign rejections, creation resolution + v0.1 parity |
| task-5 | VAL-018..020 | QA Lead | Four-format list, filters, `assigned_to` shapes |
| task-6 | VAL-021..025 | QA Lead | Rebuild parity, maintenance integrity, seam absence, DEC-079 boundary, exit scenario |
| task-6 (phase-end) | VAL-026 | Human | Operator dogfood; QA Lead fans out into per-observable checklist rows |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| QA plan | `../phases/phase-19/qa/qa-plan.md` | Missing |
| Mid-phase reviews | `../phases/phase-19/qa/review-M.md` | Per QA plan (6 tasks — cadence required) |
| Final review | `../phases/phase-19/qa/final-review.md` | Pending |
| Human checklist | `../phases/phase-19/qa/human-checklist.md` | Pending |

## Active Constraints

Planning decisions DEC-075..DEC-082 are recorded in
`../phases/phase-19/decisions.md`. The binding ones during execution:

- **Assignments are journal-canonical** (DEC-076). `workflow_assignments` is
  a projected table written by `workflow.assigned` events; the seeded default
  workflows stay migration DDL as system reference data.
- **Resolution precedence:** kind assignment → project assignment → built-in
  kind default. With no assignment rows, behavior is byte-identical to v0.1.
- **Cancel edges and initial state are derived, not hard-coded** (DEC-077).
  One `cancelled` state materializes cancel edges for every non-cancelled
  state; zero materializes none; more than one is rejected. Initial state is
  the first declared `unstarted` state, else the first declared state.
- **Define accepts only the four implemented rules** (DEC-078).
  `required_gates_current`, `has_completion`, and `linked_to_decision` are
  rejected with version-scoped messages; `linked_to_decision` is seeded for
  `escalation-default` only and stays known-but-unimplemented.
- **The DEC-066 seam stays behavior-free until v0.2b.** No seeded use, no
  evaluator, no query for `required_gates_current` / `has_completion`.
- **No production Go path reads the frozen v0.1 `definitions` map** (Phase 18
  VAL-024).
- **Custom state names are not transitionable until Phase 20** (DEC-079).
  Do not auto-register permissions-matrix rows to work around it.
- **No `activity_log` rows for workflow authoring events** (DEC-081).
- **Reject before append.** Every invalid define/assign leaves journal bytes
  and projected row counts unchanged.
- **Do not weaken TRD-PERF-003 (`<30 s`) or TRD-PERF-001 (`<500 ms`)**;
  CF-024 leaves only ~8-10% rebuild headroom.
- `internal/dev` `Finish`/`Block` stay over the 60-line target (CF-017,
  v0.2e) — not this phase's refactor.
- **Evidence honesty and ID discipline are standing rules now, not phase-local
  stand-ins (DEC-084..086).** A non-terminating required command is never
  `PASS`; each result carries a verbatim terminal token. IDs are phase-scoped
  (`TRG-19-NNN` / `ESC-19-NNN`, allocated from the phase triage log) and
  TL-Lite writes the entry **before** reporting its verdict.

## Interfaces / Subsystems In Scope

`migrations/` (`0009` assignments table, `0010` escalation seed);
`internal/workflow/` (authoring core, write paths, list read, resolution);
`internal/issues/create.go` (workflow + initial-state resolution);
`internal/policy/matrix.go` (two new rows — changes `fid policy show`);
`cmd/fid/workflow.go` (new command family); `internal/format/` (markdown
renderer); `internal/maintenance/` (assignment integrity check).

## Known Risks

- **Task 4 is the heaviest.** It changes `internal/issues.Create`, the most
  exercised path in the product. The no-assignment parity half of VAL-017 is
  the regression guard; capture pre-change output bytes first.
- **Reassignment effect shape.** `applyDelete` requires exactly one affected
  row, so the event must carry `delete`+`insert` only when a prior row exists.
- **Coverage.** `internal/workflow` is a COV-001 80% package gaining a large
  amount of new code.
- **CF-023 residue (now DEC-087).** The retrospective landed the prompt and
  template fixes, but the only *mechanical* detector for a `PASS` recorded
  against a non-terminating command — a `report-lint` evidence-token check plus
  `drill-09` — is not built. QA Lead's independent rerun remains the real
  control, and heavy mode retains it.

## Open Questions

- [ ] **CF-008 disposition.** DEC-073 made Phase 19 planning the explicit
  decision point. DEC-082 re-defers it to v0.2d (governance vocabulary);
  the Human should confirm or override.

## Success / Exit Criteria

- [ ] All six tasks complete with green quality gates.
- [ ] An agent can define a custom workflow, assign it to a kind or project,
  list workflows in every format, and create an issue that picks up the
  assigned workflow instead of the built-in default (VAL-025); the
  `escalation-default` workflow is present and inspectable (VAL-003).
- [ ] VAL-001..026 all recorded with cited evidence; final review has no
  unresolved blocking finding; human checklist accepted (VAL-026 fanned out).
- [ ] Phase sign-off decision recorded.

## Next Planning Checkpoint

**Phase 20 — `role_permissions` migration + `fid role create` /
`fid policy set` + v0.2a closure.** It carries the DEC-079 unblock (custom
state names become transitionable) and the v0.2a exit smoke.
