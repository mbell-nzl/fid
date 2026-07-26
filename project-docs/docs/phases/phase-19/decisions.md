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

## QA Plan Blocker Dispositions (2026-07-26)

The Phase 19 QA plan opened **BLOCKED** on QA-BLK-01..07. DEC-089..DEC-093
disposition all seven. Nothing in the task pack's scope was widened or
weakened; the only behavioral addition is DEC-090's define-time rule.

| Blocker | Disposition | Decision |
|---|---|---|
| QA-BLK-01 | Boundary resolved in favor of the task pack; `02_TRD.md` §3 synchronized now | DEC-089 |
| QA-BLK-02 | PRD-WF-013 preserved unamended; issue-applicable workflows must declare exactly one `cancelled` state | DEC-090 |
| QA-BLK-03 | 14 bundled assertions narrowed; split-out observables added one-for-one in the same task windows | DEC-091 |
| QA-BLK-04 | 20 assertions added for uncovered binding behavior | DEC-091 |
| QA-BLK-05 | Sources corrected; VAL-003/012/059 evidence made deterministic and reproducible | DEC-091 |
| QA-BLK-06 | `05_DESIGN_PLAN.md` §3.1 command tree updated now; per-task verification ownership assigned | DEC-092 |
| QA-BLK-07 | Human approval recorded; contract marked `Ready`; CF-008 disposition confirmed | DEC-093 |

### DEC-089: `internal/workflow` owns workflow authoring, assignment, and resolution

- **Date:** 2026-07-26
- **Resolves:** QA-BLK-01.
- **Context:** `02_TRD.md` §3 described `internal/workflow/` as a *read-only
  registry* and listed "workflow authoring/mutation" in its must-not column.
  That wording was written for v0.1, when workflow structure was a hard-coded
  Go map and there was no authoring surface at all. PRD-WF-011/012 (v0.2)
  require authoring and assignment, and Tasks 2-4 place them in
  `internal/workflow/`. QA correctly refused to choose between changing the
  architecture and changing the spec.
- **Decision:** the **task pack is correct and the TRD sentence is stale**.
  Workflow authoring, assignment validation, effect construction, the
  authoring/assignment write paths, and assignment resolution belong in
  `internal/workflow/`, for three reasons:
  1. It is the package that owns the workflow vocabulary (slugs, states,
     categories, rules, cancel-edge policy). Splitting authoring away from the
     registry would duplicate that vocabulary in a second package.
  2. The write-path shape is the **established domain-package pattern**, not a
     new one: `internal/issues` and `internal/dev` already orchestrate
     `journal.WithLock` → `PrepareEvent` → `AppendEvent` → `ProjectEvent`.
     `internal/workflow` doing the same for its own entity is symmetric.
  3. The alternatives are worse: authoring in `cmd/fid/` violates TRD-ARCH-005
     (no business logic in the interface adapter), and a new
     `internal/workflowauthoring` package would split one entity's rules
     across two packages for no benefit.
- **Boundary that stays:** `internal/workflow` still must not create issues,
  verify ACs, own schema/migrations, or format output. Task 4's
  `internal/workflow/resolve.go` is consumed **by** `internal/issues`; issue
  creation stays in `internal/issues`.
- **Spec sync:** `02_TRD.md` §3 (the layer-responsibility table row and the
  package tree annotation) is updated **now**, by the Tech Lead, as part of
  this disposition — the contradiction must be gone before Task 2 starts, not
  after. Task 6 re-verifies the row against what shipped.

### DEC-090: PRD-WF-013 stays unamended; issue-applicable workflows must declare exactly one `cancelled` state

- **Date:** 2026-07-26
- **Resolves:** QA-BLK-02. **Amends:** DEC-077 clause 1 (zero-cancelled branch).
- **Context:** PRD-WF-013 ("any issue may transition to cancelled from any
  non-cancelled state") is a v0.1 **Must**. DEC-077 allowed a definition to
  declare zero `cancelled`-category states and materialize no cancel edge —
  which is right for `escalation-default`, but nothing stopped such a workflow
  from being assigned to an issue kind, producing an issue that can never be
  cancelled. That is a silent violation of a Must.
- **Options considered:** (a) require one cancelled state for issue-applicable
  workflows; (b) restrict *assignment* of zero-cancelled workflows;
  (c) formally revise PRD-WF-013 to be conditional.
- **Decision: (a), with (b) as a free by-product.** `fid workflow define`
  rejects, before any append, a definition that declares **zero**
  `cancelled`-category states when `--applies-to` is `epic`, `task`,
  `subtask`, or `any` (the issue-applicable vocabulary). The error names the
  `applies_to` value and the PRD-WF-013 obligation. Zero cancelled states
  stays legal only for non-issue applicability (`escalation`, `qa_review`),
  which is exactly the `escalation-default` precedent. Exactly one cancelled
  state behaves as DEC-077 specifies; more than one is still rejected.
- **Why not (c):** weakening a v0.1 Must to accommodate an authoring
  convenience is the wrong trade. The requirement is cheap to keep: one
  validation rule in the Task 2 core, which already owns the cancelled-state
  count.
- **Defense in depth:** DEC-080's assignment rules already require
  `applies_to` ∈ {`<kind>`, `any`} for `--kind` and `any` for `--project`, so
  after this rule no assignable workflow can lack a cancel target. Both guards
  are asserted (VAL-038 at define time, VAL-016/VAL-041 at assign time).
- **Consequences:** Task 2 gains one rejection rule; VAL-038 replaces the
  ambiguous zero-cancelled pass state with an explicit conditional. The Phase
  19 exit scenario's `task-strict` workflow declares a `cancelled` state, as
  it already did.
- **Spec sync:** the PRD-WF-013 scoping note (already owned by Task 3) now
  records the DEC-090 rule rather than an open-ended exemption;
  `06_DESIGN_PLAN_CLI.md` documents the define-time rejection (Task 3).

### DEC-091: Validation contract refinement — unbundling, coverage, and evidence

- **Date:** 2026-07-26
- **Resolves:** QA-BLK-03, QA-BLK-04, QA-BLK-05.
- **Decision:** `qa/validation-contract.md` is refined in place and marked
  `Ready`. Three changes, all additive to rigor:
  1. **Unbundling (QA-BLK-03).** VAL-001, 003, 004, 007, 009, 010, 011, 014,
     015, 016, 017, 018, 019, and 023 are narrowed to one observable each;
     every behavior they carried is now owned by exactly one assertion, in the
     **same task window**. No ID is reused or retired.
  2. **Coverage (QA-BLK-04).** New assertions cover the accepted rule
     vocabulary, `--require` params and edge validation, duplicate state and
     duplicate transition rejection, the DEC-080 `--applies-to` vocabulary and
     `any`→NULL, DEC-081 no-`activity_log` for both event types, DEC-076
     insert-versus-delete+insert effect shapes for both scopes, deterministic
     list ordering, the all-format empty result, and the standing prohibition
     on production reads of the frozen definitions fixture.
  3. **Evidence (QA-BLK-05).** VAL-002's, VAL-009's, and VAL-021's sources are
     corrected; VAL-003 gains a tuple-for-tuple frozen reference (with VAL-030
     for the three Phase 18 defaults); VAL-012 names all four documented arrow
     aliases; VAL-059 specifies a reproducible immutable baseline — the Task 3
     completion commit, a committed capture script, committed artifacts under
     `cmd/fid/testdata/v01_parity/`, and a committed, enumerated
     volatile-field normalization allow-list.
- **Result:** 26 → 70 assertions (10/14/13/13/11/9 across Tasks 1-6). Task
  ownership, scope, and the six-task slicing are unchanged; the only new
  implementation obligation in the whole disposition is DEC-090's rule.
- **Obligation on the QA Lead:** the QA plan's Section 4 evidence map and
  Section 7 deterministic checklist must gain one step per new assertion
  before the plan leaves BLOCKED. No new assertion may inherit another's
  outcome.

### DEC-092: `05_DESIGN_PLAN.md` §3.1 command-tree ownership

- **Date:** 2026-07-26
- **Resolves:** QA-BLK-06.
- **Context:** §3.1 is the canonical CLI command tree and had no `workflow`
  family, and no task owned adding it. The tree already documents
  not-yet-implemented commands with version annotations (`gate` v0.2b,
  `envelope` v0.3a), so it is a design index, not a shipped-surface mirror.
- **Decision:** the Tech Lead adds the `workflow` family
  (`define` / `assign` / `list`, annotated v0.2a) to §3.1 **now**, so the
  canonical tree is not silently stale during execution. Per-task
  verification ownership: Task 3 (`define`), Task 4 (`assign`), and Task 5
  (`list`) each confirm the §3.1 entry matches the flag surface they shipped
  and correct it in the same task if it diverges; Task 6 records the final
  §3.1 and `02_TRD.md` §3 verdicts in the spec-sync sweep.
- **Consequence:** `05_DESIGN_PLAN.md` moves from `TBD at wrap` to `touched`
  in the spec-sync table below.

### DEC-093: Human approval of the Phase 19 plan; CF-008 confirmed

- **Date:** 2026-07-26
- **Resolves:** QA-BLK-07.
- **Approval recorded (Human, 2026-07-26):** operating mode **heavy**
  (DEC-075); the **six-task** list and its ordering; the refined validation
  contract **VAL-001..VAL-070**, now `Ready`; the **CF-008** re-defer to v0.2d
  (DEC-082), which is hereby confirmed rather than pending; and the DEC-089,
  DEC-090, and DEC-092 dispositions above.
- **Consequences:** the phase brief and `progress_current.md` no longer carry
  an approval precondition. The remaining precondition to Task 1 is
  procedural: the QA Lead re-verifies the QA plan against the refined contract
  (DEC-091) and moves it from BLOCKED to Ready, with fresh heavy-mode command
  and feasibility validation retained from the 2026-07-26 run only where the
  underlying artifact did not change.
- **Standing:** DEC-084..DEC-086 (phase-scoped IDs, triage row before verdict,
  verbatim terminal token per required command) remain binding for execution.

## Spec Sync Verdicts

Recorded at wrap-up for the full canonical spec suite. Known obligations from
the planning decisions above:

| Doc | Expected verdict | Owning task |
|---|---|---|
| `00_PROJECT_CHARTER.md` | TBD at wrap | — |
| `GLOSSARY.md` | touched (workflow assignment, custom workflow, initial state) | Task 6 |
| `01_PRD.md` | touched — PRD-WF-013 scoping note (DEC-077, **as amended by DEC-090**) | Task 3 |
| `02_TRD.md` | touched — §3 `internal/workflow` boundary (DEC-089, **applied at planning**); §4.6 assignment precedence + escalation seed (DEC-076) | Planning (§3), Task 1 (§4.6); Task 6 verifies |
| `03_TRD_SCHEMA.md` | touched — §3 `workflow_assignments`, `applies_to` vocabulary, `linked_to_decision` (DEC-076/078/080); §2.1 resolution source | Tasks 1, 4 |
| `04_TRD_CONFIG.md` | TBD at wrap | — |
| `05_DESIGN_PLAN.md` | touched — §3.1 `workflow` command family (DEC-092, **applied at planning**) | Planning; Tasks 3, 4, 5 verify their subcommand; Task 6 records the verdict |
| `06_DESIGN_PLAN_CLI.md` | touched — define/assign/list flags, rules, errors, limits (DEC-077/078/079/080/090) | Tasks 3, 4, 5, 6 |
| `07_DESIGN_PLAN_FLOWS.md` | TBD at wrap | — |
| `08_QUALITY_PLAN.md` | TBD at wrap | — |
| `09_AI_WORKFLOW.md` | TBD at wrap | — |
| `10_IMPLEMENTATION_ROADMAP.md` | touched — Phase 19 closure note | Wrap |

## Workflow Notes

Appended at wrap-up.
