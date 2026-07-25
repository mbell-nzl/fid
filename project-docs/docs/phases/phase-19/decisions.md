# Phase 19 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes,
> and wrap-up decisions. DEC numbering is global and never resets.

Phase: 19 — `fid workflow define/assign/list` + escalation default workflow
(v0.2a slice 2 of 3).

## Planning Decisions (2026-07-26, pending Human approval with the task list)

### DEC-075: Phase 19 operating mode is **Heavy**

- **Date:** 2026-07-26
- **Context:** `09_AI_WORKFLOW.md` §4.3 makes heavy mode mandatory for any
  phase carrying a data migration, a public-interface change, or a
  persistence-contract change.
- **Decision:** Phase 19 runs in **heavy mode**. It hits three heavy entry
  criteria at once: two new schema migrations (`0009` assignments table,
  `0010` escalation-default seed) with a projection-schema version bump;
  three new public CLI commands (`fid workflow define/assign/list`) plus two
  new permissions-matrix rows that change `fid policy show` output; and two
  new journal event types (`workflow.defined`, `workflow.assigned`) that
  become permanent replay history.
- **Consequences:** Full artifact set — validation contract, QA plan with
  **fresh** command and feasibility validation (heavy mode never carries
  feasibility forward), mid-phase review cadence per the QA plan, final
  review, human checklist, spec sync across all canonical specs at wrap, and
  Phase Runner orchestration.

### DEC-076: Workflow assignments are a new projected table written by journal events

- **Date:** 2026-07-26
- **Context:** `03_TRD_SCHEMA.md` §3 ships four workflow tables but **no**
  assignment storage, while `06_DESIGN_PLAN_CLI.md` specifies
  `fid workflow assign --name <slug> [--kind <kind>] [--project]` and a
  `fid workflow list` payload containing
  `"assigned_to": [{"scope": "kind", "kind": "task"}]`. The spec is
  incomplete, not wrong: assignment is required by PRD-WF-012 and has no home.
- **Decision:** Add `workflow_assignments` in migration `0009` with exactly
  this shape:

  ```sql
  CREATE TABLE workflow_assignments (
    scope       TEXT NOT NULL,
    kind        TEXT,
    workflow_id TEXT NOT NULL REFERENCES workflows(id),
    assigned_at DATETIME,
    CHECK (scope IN ('kind', 'project')),
    CHECK ((scope = 'kind' AND kind IS NOT NULL)
        OR (scope = 'project' AND kind IS NULL))
  );

  CREATE UNIQUE INDEX idx_workflow_assignments_kind
    ON workflow_assignments(kind) WHERE scope = 'kind';

  CREATE UNIQUE INDEX idx_workflow_assignments_project
    ON workflow_assignments(scope) WHERE scope = 'project';
  ```

  Rules that follow:
  1. **Assignments are user data, so they are journal-canonical** (ADR-0001).
     `fid workflow assign` appends a `workflow.assigned` event whose effects
     write the row; the table is a derived projection like every other. The
     seeded default workflows stay migration DDL because they are system
     reference data (unchanged from Phase 18).
  2. **Resolution precedence at issue creation:** kind-scope assignment →
     project-scope assignment → the built-in kind default
     (`epic-default` / `task-default` / `subtask-default`). With no
     assignment rows present, behavior is byte-identical to v0.1.
  3. **Reassignment is last-write-wins,** expressed as an explicit ordered
     `delete` + `insert` effect pair in one event when a prior row exists,
     and a bare `insert` when it does not. Both forms replay deterministically
     because the effects are recorded literally.
  4. **Two partial unique indexes, not a composite primary key.** SQLite
     permits multiple NULLs in a primary key, so `PRIMARY KEY (scope, kind)`
     would silently allow several project-scope rows.
  5. **`--name` sets `workflows.id` (the stable slug) and `workflows.name`**
     to the same value. The documented flag set has no display-name flag, and
     inventing one is an unnecessary CLI surface change; `workflow list`
     still reports both fields.
  6. **Slug and state-name validation:** workflow slug matches
     `^[a-z0-9][a-z0-9-]{0,63}$`; state names match `^[a-z][a-z0-9_]{0,63}$`;
     categories are exactly `unstarted`, `started`, `completed`, `cancelled`.
- **Spec sync obligation:** `03_TRD_SCHEMA.md` §3 (add the table, indexes,
  and CHECKs) and `02_TRD.md` §4.6 (assignment resolution precedence) are
  updated **in Task 1**, with the change itself.

### DEC-077: Cancel-edge materialization and initial-state resolution for custom workflows

- **Date:** 2026-07-26
- **Context:** PRD-WF-013 ("any issue may transition to cancelled from any
  non-cancelled state") is a v0.1 **Must**, and the seeded defaults satisfy it
  by seeding every `<state> → cancelled` edge explicitly rather than
  special-casing cancel in code (`02_TRD.md` §4.6). A user-authored
  `--transitions` list would not do that by hand. Separately,
  `internal/issues.Create` hard-codes the initial status `"todo"`, which does
  not exist in an arbitrary custom workflow and would fail the
  `(workflow_id, status)` composite FK.
- **Decision:**
  1. **Cancel edges.** If a definition declares **exactly one**
     `cancelled`-category state, `fid workflow define` materializes
     `<state> → <cancelled-state>` for every non-cancelled declared state,
     deduplicated against edges the user declared explicitly. If it declares
     **zero** cancelled-category states, no cancel edge is materialized —
     this is the `escalation-default` precedent (`open` → `resolved` only,
     `02_TRD.md` §4.6). If it declares **more than one**, define is rejected
     before any append because the cancel target is ambiguous.
  2. **Initial state.** An issue created under a workflow starts in the
     **first declared state whose category is `unstarted`**; if the workflow
     declares no `unstarted` state, it starts in the **first declared state**.
     `--states` is documented as "Ordered state names", so declaration order
     is the `workflow_states.ordinal` and the rule is deterministic. For all
     three seeded defaults this resolves to `todo`, preserving v0.1 behavior
     exactly.
- **Spec sync obligation:** PRD-WF-013 gains a scoping note (cancel edges
  exist for workflows that declare a cancelled-category state; the seeded
  escalation default deliberately declares none) and
  `06_DESIGN_PLAN_CLI.md` documents both rules on `fid workflow define`.
  Applied in Task 3.

### DEC-078: `fid workflow define` accepts only the implemented rule vocabulary

- **Date:** 2026-07-26
- **Context:** `decisions_active.md` carries a standing commitment that the
  DEC-066 gate-obligation seam stays **behavior-free until v0.2b**: no seeded
  use and no evaluation of `required_gates_current` / `has_completion`.
  `internal/issues.evaluateTransitionRule` returns
  `unsupported workflow rule %q` for any name it cannot evaluate, so a
  workflow authored with an unimplemented rule would define successfully and
  then fail at transition time — a latent defect, not a seam.
- **Decision:** `fid workflow define --require` accepts exactly the four
  implemented rules: `all_acs_verified`, `has_completion_report`,
  `concerns_addressed`, `all_child_tasks_done_or_cancelled`. It **rejects**
  before any append:
  - `required_gates_current` — error names v0.2b;
  - `has_completion` — error names v0.2c;
  - `linked_to_decision` — error names v0.2d;
  - any other name — unknown-rule error naming the offending value.

  `linked_to_decision` is nevertheless **seeded** for `escalation-default` by
  migration `0010` (system reference data, `02_TRD.md` §4.6) and is added to
  the registry loader's known-but-unimplemented set alongside the two
  reserved seam names, so `workflow.Load`, `fid workflow list`, and
  `fid maintenance check` surface it structurally without evaluating it. No
  issue can reach it in v0.2a: `issues.kind` is still
  `CHECK (kind IN ('epic','task','subtask'))`, and the escalation entity CLI
  is v0.2d.
- **Spec sync obligation:** `06_DESIGN_PLAN_CLI.md` documents the accepted
  rule vocabulary and the three version-scoped rejection messages (Task 3);
  `03_TRD_SCHEMA.md` §3 notes `linked_to_decision` joins the
  known-but-unimplemented set (Task 1).

### DEC-079: Custom **state names** are not transitionable until Phase 20

- **Date:** 2026-07-26
- **Context:** `policy.Authorize` resolves `issue.transition.<state>` against
  the in-code matrix and returns
  `authorize action %q: no permissions matrix entry` when the action is
  absent. The matrix has rows only for the v0.1 state names. A custom
  workflow with a state named `backlog` therefore cannot be transitioned into
  by any role, including `human`. The fix is exactly Phase 20's deliverable:
  the `role_permissions` migration plus `fid role create` / `fid policy set`.
- **Decision:** Phase 19 ships **authoring, assignment, listing, and issue
  creation** under custom workflows. Transitioning into a custom-named state
  is an accepted, documented boundary that fails with the existing
  no-matrix-entry error; it is **not** silently worked around by
  auto-registering matrix rows at define time (that would be an unreviewed
  permissions grant). Custom workflows that reuse v0.1 state names are fully
  usable end to end today, which is what the Phase 19 exit demonstration
  uses.
- **Consequences:** VAL-024 asserts both states (custom-named state rejected;
  v0.1-named state under a custom workflow accepted).
- **Spec sync obligation:** `06_DESIGN_PLAN_CLI.md` records the limitation
  and points at `fid policy set` (v0.2a Phase 20). Applied in Task 6.

### DEC-080: `--applies-to` vocabulary and `--project` compatibility

- **Date:** 2026-07-26
- **Context:** `06_DESIGN_PLAN_CLI.md` lists `--applies-to` values as
  `epic`, `task`, `subtask`, `qa_review`, `any`, but `02_TRD.md` §4.6 seeds
  an `escalation-default` workflow, and `03_TRD_SCHEMA.md` §3 comments the
  column as `'epic' | 'task' | 'subtask' | 'qa_review' | NULL (any)`. The
  seeded value has no accepted CLI spelling.
- **Decision:** `--applies-to` accepts `epic`, `task`, `subtask`,
  `qa_review`, `escalation`, and `any`. `any` projects as
  `workflows.applies_to IS NULL`. `fid workflow assign --kind <k>` requires
  the workflow's `applies_to` to be `<k>` or `any`;
  `fid workflow assign --project` requires `applies_to` to be `any`, because
  a project default must be legal for every kind.
- **Spec sync obligation:** `06_DESIGN_PLAN_CLI.md` and `03_TRD_SCHEMA.md` §3
  column comment updated in Task 3 and Task 4 respectively.

### DEC-081: No `activity_log` rows for workflow authoring events in Phase 19

- **Date:** 2026-07-26
- **Context:** `activity_log` is polymorphic over
  `('issue','requirement','decision','escalation')` and is guarded by the
  migration `0005` `activity_log_ref_integrity_guard_insert` trigger. Adding
  a `workflow` ref type means a trigger rebuild and a new ref family.
- **Decision:** `workflow.defined` and `workflow.assigned` write **no**
  `activity_log` rows in Phase 19. The journal is the audit record for them
  and is inspectable via `fid journal`. Extending the `activity_log` ref-type
  whitelist to non-issue entities is v0.2d work, where the decision,
  escalation, and comment entity CLIs land together.
- **Consequences:** the validation contract records this as an explicit
  out-of-scope note on VAL-009/VAL-014 so a later reviewer does not read the
  absence as an omission.

### DEC-082: Carry-forward disposition at Phase 19 planning

- **Date:** 2026-07-26
- **Decision:** each open carry-forward item is dispositioned as follows.

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | **Re-deferred to v0.2d** (was: "explicit decision point at Phase 19 planning", DEC-073) | The resolution-reason / `cancel --reason` enhancement is an `issues` data-model plus `activity_log` audit change, orthogonal to workflow authoring. It belongs with the governance vocabulary in v0.2d, where decision/escalation/triage entities and their reason semantics land together, and where DEC-081's `activity_log` ref-type extension is already scheduled. Folding it into a heavy migration phase would mix two unrelated schema changes against the task-slicing rule. **Presented to the Human for confirmation with this plan.** |
| CF-017 | Re-deferred to v0.2e | `internal/dev` is not touched by Phase 19; trigger has not fired. |
| CF-018 | Re-deferred | `internal/sqlquery` is consumed only by `fid sql`; Phase 19 does not reopen it. |
| CF-020 | Re-deferred to v0.2 | `occurred_at` monotonicity is untouched. |
| CF-021 | **Scheduled — Task 6** | Trigger fired in Phase 18 and was missed. Task 6 already touches `internal/maintenance` for the assignment-integrity check and runs the full Go gate, so the doc comments on `maintenance.VacuumResult` / `Vacuum` ride along. |
| CF-022 | **Partially taken — Task 5** | Trigger fires: Task 5 adds a markdown renderer in `internal/format`. Task 5 must ship exact-output markdown characterization for the **new** workflow-list renderer (VAL-018). The broader backfill for pre-existing renderers stays deferred to the next `internal/format` markdown reopen. |
| CF-023 | **Re-deferred to the workflow overhaul — see DEC-083** (supersedes the DEC-074 blocking precondition) | Human direction 2026-07-26: a broader agentic-workflow overhaul is close, and these tweaks would be discarded by it. |
| CF-024 | Re-deferred to v0.2b, **with a Phase 19 measurement obligation** | Phase 19 adds no per-event projection work, but it adds a fourth workflow to every `workflow.Load` and an assignment lookup to `issue create`. Task 6 runs the benchmark gate and records the 100k-event rebuild (TRD-PERF-003, `<30 s`) and 10k-issue read (TRD-PERF-001, `<500 ms`) numbers against the controls. Do not weaken either target. |

### DEC-083: CF-023 is re-deferred to the workflow overhaul; Phase 19 execution is no longer gated on it

- **Date:** 2026-07-26
- **Supersedes:** DEC-074 (which made the CF-023 workflow-tooling pass a
  required precondition to Phase 19 execution).
- **Context:** The Human is mid-way through an overhaul of the agentic
  engineering workflow that makes changes substantially larger than CF-023's
  scope, and expects to bring it into fid before long. Landing CF-023's
  prompt, template, and lint tweaks now means discarding them at the
  overhaul. The Human asked for a residual-risk read before deferring.
- **Risk read (Tech Lead):** **No blocking quality risk.** Seven of the eight
  CF-023 items — (a) commit-metadata template affordance, (b) the
  report-is-its-own-commit prompt rule, (c) report-lint dirty-vs-parse
  messaging, (d) `ESC-NNN`/`TRG-NNN` global uniqueness, (e) minting an
  `ESC-NNN` plus triage-log row on every `BLOCK`, (f) run-ledger
  chronological append, (g) the wrap-time reopen audit — are **audit-trail
  and orchestration-friction** defects. They raise reconstruction cost and
  burn STOP/resume turns; none of them can let a defect reach the product,
  and all seven are the exact class of artifact an overhaul replaces
  wholesale.

  Item (h), **evidence-honesty hardening** (a non-terminating required
  command must not be recordable as `PASS`), is the one genuine
  product-quality risk: it is a false-PASS on the deterministic gate, and in
  Phase 18 it cost a full final-review rerun plus a Developer remediation
  cycle. It is deferrable here because it is a **rule, not tooling**, and the
  rule is already in force for this phase (see mitigations), while the
  control that actually detected it — QA Lead's independent rerun — is
  retained by heavy mode regardless.
- **Decision:** CF-023 is **re-deferred to the workflow overhaul**. Phase 19
  execution is **not** gated on it. CF-023 stays open in the carry-forward
  register with the overhaul as its target, so it is not lost if the overhaul
  slips or lands with a narrower scope.
- **Zero-cost mitigations carried by Phase 19 instead of the tooling pass:**
  1. **Evidence honesty (h).** Every Phase 19 task spec's
     `Exact Commands To Run` section already states that a non-terminating
     command is never `PASS`; Task 6 additionally requires recording that
     every command terminated before a result is claimed. This is repeated in
     the phase brief's active constraints.
  2. **ID collisions (d) and missing mints (e).** Phase 19 uses
     **phase-scoped identifiers** — `ESC-19-NNN` and `TRG-19-NNN` — which
     removes the cross-phase collision class for this phase at no cost, and
     the phase brief requires the ID **and** the triage-log row to be written
     on every `BLOCK`. Phase-scoped IDs are a phase-local convention; the
     overhaul decides the durable format.
  3. **Reopen audit (g).** Already part of the Tech Lead wrap-up procedure;
     no tooling needed.
- **Consequences:** Phase 19 may start as soon as the plan is approved. The
  deep retrospective (DEC-073) is decoupled from the tooling pass and becomes
  an **input to the overhaul** rather than a driver of a CF-023 pass; the
  Human decides its timing.
- **Refs:** DEC-073, DEC-074, CF-023; `../phases/phase-18/decisions.md`
  Triage-Log Audit and Workflow Notes.

## Spec Sync Verdicts

Recorded at wrap-up for the full canonical spec suite. Known obligations from
the planning decisions above:

| Doc | Expected verdict | Owning task |
|---|---|---|
| `00_PROJECT_CHARTER.md` | TBD at wrap | — |
| `GLOSSARY.md` | touched (workflow assignment, custom workflow, initial state) | Task 6 |
| `01_PRD.md` | touched — PRD-WF-013 scoping note (DEC-077) | Task 3 |
| `02_TRD.md` | touched — §4.6 assignment precedence + escalation seed (DEC-076) | Task 1 |
| `03_TRD_SCHEMA.md` | touched — §3 `workflow_assignments`, `applies_to` vocabulary, `linked_to_decision` (DEC-076/078/080) | Tasks 1, 4 |
| `04_TRD_CONFIG.md` | TBD at wrap | — |
| `05_DESIGN_PLAN.md` | TBD at wrap | — |
| `06_DESIGN_PLAN_CLI.md` | touched — define/assign/list flags, rules, errors, limits (DEC-077/078/079/080) | Tasks 3, 4, 5, 6 |
| `07_DESIGN_PLAN_FLOWS.md` | TBD at wrap | — |
| `08_QUALITY_PLAN.md` | TBD at wrap | — |
| `09_AI_WORKFLOW.md` | TBD at wrap | — |
| `10_IMPLEMENTATION_ROADMAP.md` | touched — Phase 19 closure note | Wrap |

## Workflow Notes

Appended at wrap-up.
