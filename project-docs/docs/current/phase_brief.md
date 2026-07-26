# Phase Brief - Phase 19 (Approved, awaiting QA plan unblock)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session.
>
> **Size target:** under 150 lines.

## Phase Goal

Phase 18 made workflow *structure* data. Phase 19 makes it **authorable**:
`fid workflow define` writes user workflows as journal events into the same
projected workflow tables, `fid workflow assign` binds a workflow to an issue
kind or to the project, `fid workflow list` reads them back in every supported
format, and `escalation-default` (PRD-WF-005a) is seeded as system reference
data. Success is observable: an agent defines a custom workflow, assigns it,
lists it, and creates an issue that picks it up instead of the built-in
default. v0.2a slice 2 of 3.

## Status

**Approved (DEC-093, 2026-07-26)** — mode, six-task list, refined contract
(VAL-001..070), CF-008 re-defer. The QA plan opened **BLOCKED** on
QA-BLK-01..07; all seven are dispositioned by **DEC-089..DEC-093**
(`../phases/phase-19/decisions.md`): `internal/workflow` boundary resolved in
the task pack's favor (`02_TRD.md` §3 synchronized); PRD-WF-013 preserved via
DEC-090; contract unbundled to 70 `Ready` assertions; `05_DESIGN_PLAN.md` §3.1
gains the `workflow` family. **Remaining precondition to Task 1: the QA Lead
re-verifies the plan against the refined contract and moves it to Ready.**

## Operating Mode

**Heavy** (DEC-075): two schema migrations with a projection-version bump;
three new public CLI commands plus two permissions-matrix rows; two new journal
event types that become permanent replay history. Criteria and artifact matrix:
`../../specs/09_AI_WORKFLOW.md` §4. Heavy never carries feasibility forward.

## Task Index

Specs and completion reports: `../phases/phase-19/tasks/`.

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | `workflow_assignments` schema + `escalation-default` seed | Pending | None |
| 2 | Workflow authoring core (validation + effect builder) | Pending | 1 |
| 3 | `fid workflow define` + `workflow.defined` event | Pending | 1, 2 |
| 4 | `fid workflow assign` + assignment-aware issue creation | Pending | 1-3 |
| 5 | `fid workflow list` in every output format | Pending | 1, 3, 4 |
| 6 | Closure: integrity, rebuild parity, exit scenario, spec sync | Pending | 1-5 |

## Validation Contract

Sources, evidence, and pass/fail criteria:
`../phases/phase-19/qa/validation-contract.md`. **70 assertions, all `Ready`**
after DEC-091 (14 originals narrowed; VAL-027..070 added):

| Window | Assertions | # | Theme |
|---|---|---|---|
| task-1 | VAL-001..005, 027..031 | 10 | Assignments table + constraint enforcement, escalation seed, default parity, reversibility, projection version |
| task-2 | VAL-006..008, 032..042 | 14 | Deterministic effects, purity, rejection families, rule vocabulary, cancel-edge policy, initial state, assignment validation |
| task-3 | VAL-009..013, 043..050 | 13 | Define event/projection/load, no `activity_log`, reject-before-append, permissions, arrow aliases, output shape, `--applies-to` |
| task-4 | VAL-014..017, 051..059 | 13 | Assignment + effect shapes, precedence, assign rejections, creation resolution, v0.1 parity vs an immutable baseline |
| task-5 | VAL-018..020, 060..067 | 11 | Four-format list, markdown golden, filters, `assigned_to` shapes, ordering, empty result |
| task-6 (incl. phase-end) | VAL-021..025, 068..070 | 8 | Rebuild parity, maintenance integrity, seam + frozen-map inventories, DEC-079 boundary, exit scenario |
| phase-end (Human) | VAL-026 | 1 | Operator dogfood; QA Lead fans out into per-observable checklist rows |

## QA Artifacts

All under `../phases/phase-19/qa/`. `qa-plan.md` is **BLOCKED — awaiting QA
Lead re-verification** (QA-BLK-01..07 dispositioned by DEC-089..093); it sets
the mid-phase cadence (`review-1` after Tasks 1-3, `review-2` after Tasks 4-5).
`final-review.md` / `human-checklist.md` are pending at phase end.

## Active Constraints

DEC-075..DEC-082 (planning) and DEC-089..DEC-093 (QA-blocker dispositions) are
recorded in `../phases/phase-19/decisions.md`. Binding during execution:

- **`internal/workflow` owns authoring, assignment, and resolution** (DEC-089;
  `02_TRD.md` §3 synchronized). It still must not create issues, verify ACs,
  own migrations, or format output.
- **Issue-applicable workflows must declare exactly one `cancelled` state**
  (DEC-090, amending DEC-077): `--applies-to epic|task|subtask|any` with zero
  is rejected at define time, so PRD-WF-013 (a v0.1 Must) can never be violated
  for an issue. Zero stays legal for `escalation` / `qa_review` only.
- **Assignments are journal-canonical** (DEC-076): `workflow_assignments` is a
  projected table written by `workflow.assigned` events; seeded defaults stay
  migration DDL. **Precedence:** kind → project → built-in default; with no
  assignment rows, behavior is byte-identical to v0.1.
- **Cancel edges and initial state are derived, not hard-coded** (DEC-077).
  One `cancelled` state materializes edges for every non-cancelled state; zero
  materializes none (non-issue only, DEC-090); more than one is rejected.
  Initial state is the first `unstarted` state, else the first declared state.
- **Define accepts only the four implemented rules** (DEC-078);
  `required_gates_current`/`has_completion`/`linked_to_decision` are rejected
  with version-scoped messages. `linked_to_decision` is seeded for
  `escalation-default` only and stays known-but-unimplemented.
- **The DEC-066 seam stays behavior-free until v0.2b** (no seeded use, no
  evaluator, no query), and **no production Go path reads the frozen v0.1
  `definitions` map** (Phase 18 VAL-024; re-proven by VAL-069/070).
- **Custom state names are not transitionable until Phase 20** (DEC-079); do
  not auto-register matrix rows to work around it. **No `activity_log` rows for
  workflow authoring** (DEC-081). **Reject before append**: every invalid
  define/assign leaves journal bytes and projected rows unchanged.
- **Do not weaken TRD-PERF-003 (`<30 s`) or TRD-PERF-001 (`<500 ms`)**; CF-024
  leaves ~8-10% rebuild headroom. `internal/dev` `Finish`/`Block` stay over the
  60-line target (CF-017, v0.2e).
- **Evidence honesty and ID discipline are standing rules (DEC-084..086).** A
  non-terminating required command is never `PASS`; each result carries a
  verbatim terminal token. IDs are phase-scoped (`TRG-19-NNN` / `ESC-19-NNN`),
  written **before** the verdict.

## Interfaces / Subsystems In Scope

`migrations/` (`0009` assignments, `0010` escalation seed);
`internal/workflow/` (authoring core, write paths, list read, resolution);
`internal/issues/create.go`; `internal/policy/matrix.go` (two rows — changes
`fid policy show`); `cmd/fid/workflow.go`; `internal/format/` (markdown);
`internal/maintenance/` (assignment integrity).

## Known Risks

- **Task 4 is the heaviest.** It changes `internal/issues.Create`, the most
  exercised path in the product. VAL-058/059 are the regression guard; capture
  the immutable baseline bytes **before** editing `create.go`.
- **Reassignment effect shape.** `applyDelete` requires exactly one affected
  row, so the event carries `delete`+`insert` only when a prior row exists.
- **Coverage.** `internal/workflow` is a COV-001 80% package gaining a lot of
  new code. **CF-023 residue (DEC-087):** no mechanical detector exists for a
  `PASS` against a non-terminating command; QA Lead's rerun is the control.

## Open Questions

None — CF-008's re-defer to v0.2d (DEC-082) was confirmed at approval (DEC-093).

## Success / Exit Criteria

- [ ] All six tasks complete with green quality gates.
- [ ] An agent can define a custom workflow, assign it, list it in every
  format, and create an issue that picks it up instead of the built-in default
  (VAL-025); `escalation-default` is present and inspectable (VAL-003).
- [ ] VAL-001..070 all recorded with cited evidence; final review has no
  unresolved blocking finding; human checklist accepted (VAL-026 fanned out).
- [ ] Phase sign-off decision recorded.

## Next Planning Checkpoint

**Phase 20 — `role_permissions` migration + `fid role create` /
`fid policy set` + v0.2a closure**, carrying the DEC-079 unblock and the v0.2a
exit smoke.
