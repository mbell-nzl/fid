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

### DEC-094: VAL-031 evidence correction and the no-match inventory token convention

- **Date:** 2026-07-29
- **Resolves:** ESC-19-001 / TRG-19-001 (Task 1 completion report records the
  no-match `rg` inventory as `PASS - exit 0 (no output)`; the exact command
  terminates with no output and **exit 1**).
- **Finding:** an evidence-recording defect, not an implementation defect.
  `rg` exits 1 when it finds no match, so `exit 0 (no output)` — the workflow's
  phrase for a *silent success* — cannot have been copied from that command's
  output. DEC-086 makes the verbatim token mandatory and mode-invariant, so the
  recorded result is invalid on its face. The VAL-031 observable itself (no
  production `linked_to_decision` evaluator under `internal/issues` /
  `cmd/fid`) holds: it was confirmed by the QA plan feasibility pass and again
  by the Phase Runner's independent rerun. The Task 1 implementation commit
  `a509d59` needs no change; the six required gate commands are unaffected, and
  the `rg` inventory is not one of them.
- **Decision — correction, not re-implementation.** Before Task 2 starts:
  1. Re-execute the exact command
     `rg -n --glob '*.go' --glob '!**/*_test.go' 'RuleLinkedToDecision|linked_to_decision' internal/issues cmd/fid`
     and copy the observed token. Do not transcribe the token from this
     decision, the ledger, or the QA plan — the point of the rule is that the
     token comes from a real terminal.
  2. Correct **both** occurrences in `tasks/task-1-done.md` (the VAL-031
     coverage line and the `## Tests run` line) to the observed
     `exit 1 (no output)`, stating inline that exit 1 **is** the no-match pass
     condition for this inventory.
  3. Add a correction note under `## Concerns, technical debt, or notes for the
     Tech Lead` naming TRG-19-001, so the audit trail survives in the report
     itself rather than only in the ledger.
  4. Commit the replacement report as its own `docs(task): ...` commit
     (`09_AI_WORKFLOW.md` step 13 — the `## Commit hash` field still names task
     commit `a509d59` and does not change), rerun
     `project-docs/tools/report-lint.sh` against the newly committed file, then
     resume the Phase Runner at Task 2.
  Task 1 is recorded complete only after that rerun passes.
- **Convention (binding for the rest of the phase).** For inventory-style
  assertions proven by the *absence* of matches — VAL-031, VAL-042, VAL-069,
  VAL-070 — the expected terminal result is **exit 1 with no output**, and that
  is what the report records: `PASS (no match) - exit 1 (no output)`. Exit 0
  with matches printed is the *failing* result for these assertions. Recording
  `exit 0 (no output)` for a search command is a DEC-086 violation and a hard
  stop. The QA plan already wrote this convention down for VAL-069/070; it is
  now stated in the validation contract for all four, and repeated in the Task
  2 and Task 6 specs that own the remaining three.
- **Canonical spec edits (justified by this decision):** `AGENTS.md`
  (Commands) and `09_AI_WORKFLOW.md` §11 gain one sentence each — a search or
  inventory command that finds nothing terminates non-zero, and the recorded
  token is the real exit status, never the silent-success phrase.
- **No contract weakening.** VAL-031 keeps its pass criterion, its evidence
  method, and its QA Lead validator; the QA Lead still re-runs the inventory at
  phase end. Nothing about this correction pre-clears VAL-031.

### DEC-095: The phase triage log is seeded at phase start; a verdict without a written record is a defective triage

- **Date:** 2026-07-29
- **Resolves:** the second defect in TRG-19-001 — TL-Lite returned
  `VERDICT: BLOCK (triage record could not be written)` and created neither
  `phases/phase-19/triage-log.md` nor an `ESC-19-NNN`.
- **Finding:** not a stale-prompt problem. `.claude/agents/tl-lite.md` is
  byte-current with `project-docs/agents/tl-lite.md`, grants `Write`, names the
  template, and explicitly says to *start at `001` when the file does not
  exist*. TL-Lite had everything it needed and still treated an absent file as
  an allocation blocker. This is the same failure mode DEC-085 was written to
  stop, and it is now the second phase running (Phase 18: five of six triages
  reported an escalation that was never minted). Cost: a full Tech Lead
  reconstruction session for a two-line report fix.
- **Decision:**
  1. **The Tech Lead seeds `triage-log.md` at phase start**, template header
     and empty table, as part of the phase directory. An empty allocation
     source removes the create-a-file step from TL-Lite's critical path
     entirely, which is cheaper than trusting prompt compliance a third time.
     Phase 19's is created now (TRG-19-001 backfilled from the ledger); this
     becomes a standing planning obligation from Phase 20.
  2. **Phase Runner precondition:** if `triage-log.md` is absent when a triage
     is needed, that is itself a phase-setup defect — record it in the ledger
     and stop; do not ask TL-Lite to invent an ID.
  3. **A verdict without a durable record is defective triage** and never
     satisfies a stop. The Phase Runner correctly recorded the defect instead
     of resuming; that handling is confirmed as the standard.
  4. **Escalation of last resort:** when TL-Lite genuinely cannot write, the
     Phase Runner records `TRIAGE (identifier unavailable)` — as it did here —
     and the full Tech Lead backfills both the TRG and the ESC. Backfilled
     entries are labelled as Tech Lead reconstructions so the audit trail never
     implies TL-Lite wrote them.
- **Carry-forward impact:** this is the recurrence CF-023/DEC-087 predicted.
  Its remaining scope (drill-09 for the non-terminating / mis-recorded-token
  defect class, plus the `report-lint.sh` evidence-token check) stays in the
  workflow-overhaul bucket, but is **raised from Low to Medium** and gains a
  concrete first check: flag any recorded result whose command begins with
  `rg`/`grep` and whose token is `exit 0 (no output)`. A mechanical check would
  have caught this defect at Task 1 close without a human in the loop; a drill
  for the missing-triage-record class belongs in the same pass.
- **Not decided here:** whether TL-Lite stays a separate role at all. That is
  deep-retrospective scope (DEC-073), not a mid-phase change.

#### Correction (2026-07-29, per DEC-097)

The **Finding** above is wrong on its central claim and is corrected here
rather than rewritten, so the reasoning trail stays auditable.

It asserted this was "not a stale-prompt problem" on the evidence that
`.claude/agents/tl-lite.md` is byte-current, grants `Write`, and says to start
at `001` when the file does not exist — concluding that "TL-Lite had everything
it needed and still treated an absent file as an allocation blocker."

**The failing agent was the Codex `tl_lite`, not the Claude `tl-lite`.** The
Claude definition was never the one executing. `.codex/agents/tl_lite.toml`
declares no tools, and the shared prompt body told it three times that it could
not run commands — the surface through which Codex writes files. So the
prompt *was* the problem, in the provider that was actually running, and the
exoneration above was reached by reading the wrong file.

What survives: the remedy. Seeding `triage-log.md` at phase start is worth
keeping on its own merits and made both backfills cheap. What does not survive:
the root-cause attribution, and the inference that prompt compliance had been
"trusted and failed." See DEC-097 for the corrected analysis and the fix.

### DEC-096: `internal/workflow` may depend on `internal/journal` for the effect data contract; the Task 2 cycle is repaired by test packaging

- **Date:** 2026-07-29
- **Resolves:** ESC-19-002 / TRG-19-002 — Task 2's required
  `golangci-lint run ./...` gate fails `typecheck` with
  `import cycle not allowed in test`.
- **Context:** Task 2 binds the authoring core's return type to
  `[]journal.Effect`. `internal/journal` imports `internal/db`, and
  `internal/db/workflow_tables_test.go` is an **in-package** (`package db`) test
  file that imports `internal/workflow`. Go compiles in-package test files into
  package `db` itself, so the `db` test binary closes
  `db [test] → workflow → journal → db`. The Developer correctly refused to
  self-authorize either an out-of-scope file edit or a dependency-boundary
  change.
- **Finding — this is a test-packaging artifact, not a design defect:**
  1. **There is no production cycle.** `internal/workflow` production code uses
     exactly one symbol from `internal/journal`: the `Effect` struct
     (`internal/journal/event.go:31`), referenced only in
     `internal/workflow/authoring_effects.go`. `internal/journal` does not
     import `internal/workflow`.
  2. **The dependency direction is the codebase norm.** `internal/issues`,
     `internal/dev`, `internal/links`, `internal/claims`, `internal/labels`,
     `internal/runs`, `internal/requirements`, `internal/policy`,
     `internal/activity`, and `internal/bootstrap` all build `journal.Effect`
     values. Under DEC-089 `internal/workflow` owns authoring, and the effect
     set *is* the authoring output. `workflow` joining that list is consistent
     layering, not an inversion.
  3. **`journal.Effect` is schema-bearing.** It is the `effects` array of the
     event JSON that becomes permanent replay history. Re-typing or duplicating
     it inside `internal/workflow` would create exactly the drift risk this
     phase's Known Risks call out.
  4. **The blocked file is the only one affected.**
     `internal/db/workflow_tables_test.go` is the sole `package db` file
     importing `internal/workflow`. `internal/journal` has no test importing
     `internal/workflow`. The blast radius is one file.
  5. **The external test package already exists and already does this.**
     `internal/db` has four `package db_test` files
     (`issues_enriched_test.go`, `kind_triggers_test.go`, `ref_triggers_test.go`,
     `issues_workflow_fk_cli_test.go`); three of them already import
     `internal/journal`, which is legal precisely because an external test
     package is not compiled into `db`.
- **Decision:**
  1. **`internal/workflow` → `internal/journal` is confirmed as intended and is
     bounded to the `Effect` data contract.** The Task 2 API stays
     `([]journal.Effect, error)`; it is not renegotiated. `internal/workflow`
     must not grow a dependency on any other `internal/journal` surface
     (no append, no projection, no replay) — that would breach DEC-089's
     "must not own migrations / must not write" boundary.
  2. **The narrow repair is authorized** and added to Task 2's scope: split
     `internal/db/workflow_tables_test.go` so that the two tests needing
     `internal/workflow` — `TestMigration0007SeedsFrozenV01WorkflowDefinitions`
     and `TestMigration0007ReservesGateObligationShapeWithoutSeedingReservedRules`
     — move, with the helpers only they use (`assertSeededWorkflows`,
     `assertSeededStates`, `assertSeededTransitions`, `assertSeededRequirements`,
     the three tuple types and their key functions), into a new
     `package db_test` file. The remaining three tests stay in `package db`.
  3. **Test-only change, behavior-preserving.** No production file may change
     under this authorization. Assertions move verbatim; nothing is weakened,
     renamed, or dropped. The moved file re-declares the few small helpers it
     needs (`testMigrationsDir`, `closeDB`, `closeRows`,
     `migratedWorkflowTablesDB`) locally under workflow-scoped names, matching
     the existing `db_test` convention (`closeIssuesEnrichedDB` et al.).
  4. **No new exported production API.** `driverName` stays unexported;
     `newTestMigrator`/`closeMigrator` stay in `package db` with
     `TestMigration0007DownAndUpRestoresWorkflowTables`, which does not import
     `internal/workflow` and therefore does not move.
- **Rejected alternatives:**
  - *Extract `Effect` into a leaf package with a `type Effect = ...` alias in
    `journal`.* Architecturally tidier, but 154 references across 51 files, a
    new package, and an edit to the event schema's home file — mid-phase, for a
    change a type alias makes invisible. Disproportionate.
  - *Give `internal/workflow` its own effect type and map it in Task 3.*
    Contradicts the Task 2 spec, duplicates a schema-bearing type, and adds a
    mapping layer whose only job is to drift.
  - *Break `internal/journal` → `internal/db`.* Far outside phase scope.
- **Standing constraint (binding beyond Phase 19):** **no in-package
  (`package db`) test file in `internal/db` may import `internal/workflow`**,
  because `internal/workflow` now sits above `internal/journal`, which imports
  `internal/db`. Tests that need both belong in `package db_test`. Same rule
  applies to `internal/journal`'s in-package tests.
- **Spec impact:** `02_TRD.md` §3 gains the `internal/workflow` →
  `internal/journal` edge and the `Effect`-only bound; Task 6 verifies at wrap.

### DEC-097: TL-Lite's missing triage records are a provider-prompt defect on Codex, not agent unreliability

- **Date:** 2026-07-29
- **Resolves:** the second defect in TRG-19-002 — TL-Lite returned
  `VERDICT: BLOCK (required triage/escalation records could not be written; no
  identifier allocated)` and wrote neither `TRG-19-002` nor `ESC-19-002`.
- **Supersedes an incorrect first finding.** This decision was first drafted
  holding that TL-Lite's report of having "no filesystem read/write mechanism
  available" was a hallucinated capability claim, on the evidence that
  `.claude/agents/tl-lite.md:4` grants `Read, Glob, Grep, Edit, Write`. **That
  was wrong, and wrong because it inspected the wrong provider.** The failing
  agent was the **Codex** `tl_lite`, not the Claude `tl-lite`. The original
  finding is recorded here rather than deleted, because the same mistake sits
  inside DEC-095 (see the correction note appended there).
- **Actual root cause — the Codex prompt tells TL-Lite it cannot write:**
  1. `.codex/agents/tl_lite.toml` declares **no tools at all** — only `name`,
     `description`, `model`, `model_reasoning_effort`, and
     `developer_instructions`. Tool access comes solely from the global
     `sandbox_mode = "workspace-write"` in `.codex/config.toml`. The
     `tools: ... Write` / `disallowedTools: Bash` grant lives exclusively in the
     `claude:` frontmatter block of `project-docs/agents/tl-lite.md` and is
     emitted **only** into `.claude/agents/tl-lite.md`.
  2. The **shared** prompt body told TL-Lite three separate times that it could
     not run commands: "You do not write code, **run commands**…" (opening),
     "never **run commands**" (Write Boundaries), and "**You cannot run
     commands**, so you never verify commits yourself" (Triage Rules).
  3. On Claude those two facts are coherent: `Write`/`Edit` are distinct tools
     from `Bash`, so "never run commands" plainly means "don't use Bash."
     **On Codex they are contradictory**, because Codex exposes file editing
     through the same command/`apply_patch` surface as execution. An agent told
     it cannot run commands, given no declared write tool, has a defensible
     basis to conclude it has no filesystem write mechanism — which is
     precisely, and accurately, what it reported.
  4. This is therefore a **prompt-authoring defect in a provider-shared body**,
     not agent unreliability. The agent behaved correctly given its
     instructions, and reported the limitation explicitly rather than inventing
     a placeholder ID — which is what its prompt asks for in that situation.
- **Decision:**
  1. **Fix the prompt at the canonical source.**
     `project-docs/agents/tl-lite.md` is amended: the blanket "run commands"
     prohibition is narrowed to its real intent (no quality gates, tests,
     builds, git, or verification commands); the opening states that writing
     the triage record is a file write and never optional; and a
     `<!-- body: codex -->` segment names `apply_patch` under
     `workspace-write` as the write affordance, with a matching
     `<!-- body: claude -->` segment naming `Write`/`Edit`. A new triage rule
     requires an *attempted* write before any claim of inability, quoting the
     concrete error.
  2. **`project-docs/tools/sync-agents.sh` must be run** to regenerate
     `.claude/agents/tl-lite.md` and `.codex/agents/tl_lite.toml`; this session
     has no shell. Until it runs, the generated prompts are stale and the defect
     is still live. `sync-agents.sh --check` verifies.
  3. **The DEC-095 mitigations stand** — seeding the triage log is independently
     worthwhile and made both backfills cheap — but its *finding* is corrected,
     not its remedy.
  4. **The Phase Runner's handling is confirmed correct again**: it recorded
     `TRIAGE (identifier unavailable)`, refused the placeholder, and stopped.
  5. **No role change, and the deep-retrospective question is withdrawn.** The
     draft escalated "should TL-Lite exist at all" to the retrospective on the
     strength of the incorrect finding. With a concrete, cheap, provider-local
     cause identified, that escalation is not justified. What *is* referred to
     the retrospective is narrower and better founded: **provider-shared agent
     bodies can encode capability claims that are true for one provider and
     false for another**, and nothing checks for it.
- **Open verification item (not decided here):** whether the Codex agent config
  schema supports a per-agent sandbox/tool declaration that would make the write
  affordance explicit rather than inherited from the global config. No such key
  was invented here — an unverified key risks breaking agent loading. Confirm
  against the Codex config schema before adding one.
- **Carry-forward impact:** **CF-023 gains a drill case of a different and more
  useful shape** than the draft's — a provider-shared prompt segment that
  contradicts one provider's tool model — alongside the DEC-095 checks. A
  cheap mechanical check falls out of it: flag any shared body asserting a
  capability limit ("cannot run commands", "no access to") that is enforced by
  only one provider's frontmatter. CF-023 stays **Medium**.

### DEC-098: The VAL-059 parity baseline is pinned to the Task 3 implementation commit; ambient `HEAD` is prohibited

- **Date:** 2026-07-30
- **Resolves:** ESC-19-003 / TRG-19-003, raised by QA Review 1.
- **Context.** Three artifacts described the same VAL-059 baseline differently.
  Task 4's binding constraint said "record the **Task 3 completion commit**
  hash with `git rev-parse HEAD`". The validation contract said "the Task 3
  completion commit". QA-plan §6 said the Task 3 **implementation** commit,
  "not the following report commit and not ambient `HEAD`". At Review 1 those
  selectors demonstrably resolved to different commits: ambient `HEAD` was the
  Task 3 *report* commit `697c673`, while the implementation commit is
  `50533e2`.
- **The QA plan is right, and the task spec is the artifact that drifted.** A
  baseline is only immutable if its selector is immutable. `HEAD` moves with
  every commit — including this escalation's own doc commits — so by the time a
  Developer reaches Task 4 the ambient-`HEAD` capture would be taken from a
  *later* tree than the one VAL-059 is supposed to freeze. That defeats the
  purpose of the assertion: the parity evidence would be unreproducible by the
  QA Lead and could silently include post-Task-3 artifacts. "Completion commit"
  was itself the root ambiguity — Task N's work lands in two commits (the
  implementation commit and the `docs(task)` report commit), and only the first
  is a meaningful build target.
- **Decision.**
  1. **The baseline is pinned, not selected.** Task 4 and the validation
     contract now carry the literal verified hash
     `50533e2b813bf073cd75f75eb70dd09cd15e9040`. The Developer re-verifies it
     with `git rev-parse --verify <hash>^{commit}` **and** cross-checks it
     against the `## Commit hash` in the committed `task-3-done.md` (`50533e2`)
     before use, so a wrong pin cannot propagate silently. `git rev-parse HEAD`,
     branch names, and any other moving selector are prohibited.
  2. **Task 4's procedure is synchronized to QA-plan §6 verbatim** — the
     detached worktree at `/tmp/fid-phase19-v01-baseline`, the
     `/tmp/fid-phase19-v01-baseline.bin` binary plus its recorded `sha256sum`,
     the two-positional-argument `capture.sh`, the fixed
     `cmd/fid/testdata/v01_parity/fixture/` with fixed IDs and journal
     timestamps, the exact six output filenames, the candidate output root
     `/tmp/fid-phase19-v01-parity-candidate/`, and the committed `SHA256SUMS`
     manifest. The earlier improvised paths (`/tmp/fid-baseline`,
     `<command>-<kind>.json`) are withdrawn.
  3. **Normalization is the committed no-op** `{"json_pointers":[],"replacement":null}`;
     all six files are raw-byte compared. The task spec's "allow-list of
     creation/update timestamps" and the contract's matching phrase are
     withdrawn as looser than QA-plan §6. If the fixed fixture cannot hold
     timestamps stable, **Task 4 is `BLOCKED` and escalates**; the allow-list is
     not widened and no expectation is hand-written. This tightens VAL-059; it
     does not weaken it.
  4. **The Developer does not read the QA plan**, so every QA-plan §6 command
     is now written into Task 4's `## Exact Commands To Run`, including the
     `sha256sum -c` and `diff -r` verification pair and the worktree cleanup.
- **Consequences.** Task 4 is unblocked and may start. VAL-059's evidence is
  reproducible by the QA Lead from committed artifacts alone. No implementation
  behavior, API, or other assertion changes. `qa-plan.md` needs **no** edit — it
  was correct throughout; the fix is in `task-4.md` and `validation-contract.md`.
- **Workflow lesson (for the wrap-up notes).** The QA plan refined VAL-059's
  procedure under DEC-091, but the already-written Task 4 spec was not
  re-synchronized, and nothing mechanically compares a task spec's evidence
  procedure against the QA plan's. Two independent prose copies of one procedure
  drifted, and only a human-grade QA read caught it. Candidate rule: when a QA
  plan fixes a procedure for a `VAL-NNN`, the task spec should **cite** that
  section rather than paraphrase it. Referred to the deep retrospective.

### DEC-099: Task 4 recovery — preserve the partial work; re-verify the baseline by reproduction, not by restart

- **Date:** 2026-08-01
- **Resolves:** ESC-19-004 / TRG-19-004, raised from the Task 4 blocker report.
- **Context.** Task 4's first required VAL-059 command,
  `capture.sh /tmp/fid-phase19-v01-baseline.bin cmd/fid/testdata/v01_parity`,
  failed with `cannot create cmd/fid/testdata/v01_parity/issue-get-epic.json:
  Directory nonexistent`. The Developer repaired the script, retried
  successfully, continued through three red/green cycles, then stopped without
  running the six required gates and without committing. TL-Lite recorded BLOCK
  on the ground that "the retry cannot substitute for the failed capture" and
  recommended discarding or reconciling the partial work and rerunning the whole
  VAL-059 procedure and every gate **from the beginning**.
- **Finding 1 — the failed command was a defect in a Task 4 deliverable, not a
  gate failure.** `capture.sh` is created *by this task* (Task 4 scope list,
  `cmd/fid/testdata/v01_parity/`). Root cause is now visible in the repaired
  script: it `cd`s into a `mktemp -d` workspace and then wrote to the
  **relative** `output_root`, so `cmd/fid/testdata/v01_parity/...` resolved
  inside the temp dir and did not exist. The fix absolutizes `output_root`
  (lines 12-15) and adds `mkdir -p` (line 19). That is precisely the authoring
  the task asked for. The same report records three `RED expected failure` →
  `PASS` test cycles that nobody treats as failed gates; a harness reaching
  green is the identical loop.
- **Finding 2 — the failed run wrote nothing.** It aborted on the first
  redirection. It cannot have produced, contaminated, or partially overwritten
  any artifact; with respect to evidence it is a no-op.
- **Finding 3 — the baseline's authenticity never depended on that first
  capture, because DEC-098 already replaced the temporal guard with a
  structural one.** The baseline binary is built in a **detached worktree at
  the pinned commit** `50533e2…`, and `capture.sh` consumes only that binary
  plus the committed, fully-fixed fixture (verified: fixed IDs, fixed event
  IDs, fixed `occurred_at`, and outputs whose `created_at`/`updated_at` derive
  solely from the journal). Baseline output is therefore a pure function of
  *(pinned commit, committed fixture)*. Uncommitted working-tree edits to
  `internal/issues/create.go` are structurally incapable of affecting it. Task 4
  says so itself: "The QA Lead must be able to regenerate the baseline from the
  recorded hash and committed script alone."
- **The BLOCK is sustained in outcome and corrected in reasoning.** Stopping was
  right and the escalation was correctly routed (trigger 4: partial work
  spanning multiple sessions and a public creation path). Task 4 is genuinely
  incomplete — no implementation or report commit, six required gates never run,
  and the three VAL-059 verification commands (candidate capture, `sha256sum
  -c`, `diff -r`) never run. But "restart the whole procedure from the
  beginning" does not follow, and no rule requires it. What was in doubt is
  **provenance of the six committed baseline outputs** — the report's prose
  leaves the retry's timing relative to the `create.go` edit ambiguous. That
  doubt is cheaply and mechanically decidable, and a narrative cannot settle it
  either way.
- **Decision.**
  1. **Preserve the uncommitted work; do not discard it.** `09_AI_WORKFLOW.md`
     §2.3 responsibility 2 directs a Developer to inspect in-scope partial work
     and continue from that state, and this phase already set the precedent at
     ESC-19-002/DEC-096, where Task 2's uncommitted work was preserved intact.
     Discarding would destroy correct, evidenced work and re-run the phase's
     highest-risk change for no evidentiary gain.
  2. **Settle provenance by reproduction, not by narrative.** Before any other
     Task 4 work, the recovery session removes the stale scratch worktree,
     rebuilds the baseline binary from the pinned commit, re-captures into a
     **fresh scratch root**, and checks the result against the committed
     `SHA256SUMS`. A match proves the committed baselines authentic regardless
     of when they were written; a mismatch catches contamination and the
     committed baselines are replaced from the reproduction. Either way the
     ambiguity is closed by a deterministic test. Exact commands are written
     into `task-4.md`.
  3. **Remove the scratch worktree through the authorized path** —
     `git worktree remove /tmp/fid-phase19-v01-baseline` then
     `git worktree prune` — before re-adding it, since the stopped session left
     it registered and `git worktree add` would otherwise fail.
  4. **Then complete Task 4 normally:** the candidate build, the candidate
     capture, `sha256sum -c`, `diff -r`, the six required gates, the three named
     spec updates, the implementation commit, and the `docs(task)` report
     commit.
  5. **Report the first failure honestly as a deviation** — its verbatim token,
     the root cause, the fix, and the final `PASS` token — under DEC-100 below.
     It is not recorded as a passed gate and it is not concealed.
- **Spec correction.** `task-4.md` still carried the pre-DEC-098 *temporal*
  guard ("the first three run **before** any `internal/issues/create.go` edit";
  "once `create.go` changes, the baseline can no longer be captured from the
  working tree"). Under DEC-098 that is obsolete and actively misleading — the
  second clause refutes itself by conceding the pinned commit can always
  reproduce it, and reading it as a binding precondition is what turned a
  fixable harness bug into a full-task hard stop. It is demoted from a binding
  precondition to sequencing guidance, with the structural guarantee stated
  explicitly. VAL-059 is unchanged and not weakened: the baseline is still the
  pinned commit, normalization is still the committed no-op, and comparison is
  still raw-byte.
- **Consequences.** Task 4 resumes in place from preserved work. No change to
  any assertion, API, or behavior. `validation-contract.md` and `qa-plan.md`
  need no edit.

### DEC-100: A task-authored evidence harness that fails while being authored is a deviation, not a failed gate

- **Date:** 2026-08-01
- **Status:** binding beyond Phase 19.
- **Context.** DEC-086 and `AGENTS.md` (Commands) make required-command evidence
  strict: `PASS` or `BLOCKED`, verbatim tokens, no missing/skipped/narrower
  substitutes, and a non-terminating command is never `PASS`. Those rules were
  written for gates run *against finished work*. They say nothing about a
  required command whose **harness is itself a deliverable of the same task**,
  and Phase 19 Task 4 shows how that gap reads in practice: a one-line path bug
  in a script the task existed to write was treated as a permanently failed gate
  requiring a from-scratch restart.
- **Decision.** A required command that fails, is repaired, and is rerun
  **byte-identically** to a deterministic pass is recorded as `PASS` with a
  **mandatory deviation note**, provided **all four** hold:
  1. The failing artifact is a **deliverable of the task being authored** — a
     fixture, capture/measurement script, or golden generator. **Not** the
     product under test, **not** a test asserting product behavior, and **not** a
     pre-existing committed gate, tool, or script.
  2. The rerun command is **byte-identical**. Nothing is narrowed, substituted,
     re-scoped, or re-pointed.
  3. The final result is **deterministic and reproducible by a third party from
     committed artifacts alone** — not a flaky or environment-dependent pass.
  4. The report records the first failure's **verbatim token, the root cause,
     and the fix**, plus the final result's own token.
- **Otherwise the strict rule stands.** A required command that fails against
  the product under test, or that passes only after the command, the test, or
  the assertion is weakened, is a **failed gate** and escalates. Reruns that
  mine for a flaky green are excluded by (3); "fix the test until it passes" is
  excluded by (1) and (2).
- **Rationale.** Red→green is the workflow's core method, not an exception to
  it. A required-command list is a statement about the **final, reproducible
  state** of a task's evidence, not a prohibition on iteration. The strict
  reading would require any task that ships its own evidence tooling to write
  that tooling correctly on the first attempt — a standard no rule imposes on
  any other deliverable, and one that would have made Task 4 unrunnable.
- **Consequences.** `AGENTS.md` (Commands) and `09_AI_WORKFLOW.md` §11 carry the
  rule. TL-Lite gains an explicit distinction it lacked: **"the task is
  incomplete, stop"** and **"the evidence is tainted, restart"** are different
  verdicts with different remedies, and a harness bug is normally the first.
- **Carry-forward impact.** CF-023 gains a drill case: a required command that
  fails on a task-authored harness, repaired and rerun byte-identically —
  correct handling is `PASS` **plus** a deviation note, not a blocker report and
  not a silent rerun.

## Spec Sync Verdicts

Recorded at wrap-up for the full canonical spec suite. Known obligations from
the planning decisions above:

| Doc | Expected verdict | Owning task |
|---|---|---|
| `00_PROJECT_CHARTER.md` | TBD at wrap | — |
| `GLOSSARY.md` | touched (workflow assignment, custom workflow, initial state) | Task 6 |
| `01_PRD.md` | touched — PRD-WF-013 scoping note (DEC-077, **as amended by DEC-090**) | Task 3 |
| `02_TRD.md` | touched — §3 `internal/workflow` boundary (DEC-089, **applied at planning**); §3 `internal/workflow` → `internal/journal` `Effect`-only edge + the `package db` test-import constraint (**DEC-096**); §4.6 assignment precedence + escalation seed (DEC-076) | Planning (§3), Task 1 (§4.6), Task 2 (DEC-096 edge); Task 6 verifies |
| `03_TRD_SCHEMA.md` | touched — §3 `workflow_assignments`, `applies_to` vocabulary, `linked_to_decision` (DEC-076/078/080); §2.1 resolution source | Tasks 1, 4 |
| `04_TRD_CONFIG.md` | TBD at wrap | — |
| `05_DESIGN_PLAN.md` | touched — §3.1 `workflow` command family (DEC-092, **applied at planning**) | Planning; Tasks 3, 4, 5 verify their subcommand; Task 6 records the verdict |
| `06_DESIGN_PLAN_CLI.md` | touched — define/assign/list flags, rules, errors, limits (DEC-077/078/079/080/090) | Tasks 3, 4, 5, 6 |
| `07_DESIGN_PLAN_FLOWS.md` | TBD at wrap | — |
| `08_QUALITY_PLAN.md` | TBD at wrap | — |
| `09_AI_WORKFLOW.md` | touched — §11 no-match search-command token rule (DEC-094, **applied mid-phase**); §11 task-authored harness deviation rule (**DEC-100**, applied mid-phase); `AGENTS.md` Commands carries both | Mid-phase (DEC-094, DEC-100); Task 6 verifies |
| `10_IMPLEMENTATION_ROADMAP.md` | touched — Phase 19 closure note | Wrap |

## Workflow Notes

Appended at wrap-up.
