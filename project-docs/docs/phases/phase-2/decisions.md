# Phase 2 Decisions

> Purpose: phase-local decision history, spec sync verdicts, triage audit,
> workflow notes, and wrap-up decisions. DEC numbering is global and never
> resets (Phase 1 ended at DEC-009).

## Decisions

### DEC-010: ProjectionSchemaVersion bumps for projected-shape changes only

- **Date:** 2026-06-07
- **Context:** Phase 2 added four migrations. Tasks 1 and 2 bumped
  `internal/db.ProjectionSchemaVersion` (1→2 for the `0002` base/stub tables,
  2→3 for the `0003` `issues_enriched` view + indexes) so existing derived DBs
  rebuild on upgrade. Tasks 3 and 4 added the `0004` kind backstop triggers and
  the `0005` polymorphic-ref/scope backstop triggers but **kept the version at
  3**. The Developers reasoned that BEFORE INSERT/UPDATE backstop triggers do
  not change the projected row shape or view/index generation, so existing
  projected rows remain structurally valid and need no fresh generation; the
  migration applies the new backstops in place and they only constrain future
  writes.
- **Decision:** Adopt as the standing rule: bump `ProjectionSchemaVersion` when
  a migration changes the **projected structure** an existing DB would have to
  reproduce (tables, columns, views, indexes that feed reads). Do **not** bump
  for migrations that add only backstop triggers, because triggers gate future
  writes and do not invalidate already-projected rows. `ProjectionSchemaVersion`
  is an internal constant, not a public interface (VAL-004 note).
- **Consequences:** No spec edit required; this is internal projection-freshness
  discipline under TRD §4.3. Validated end-to-end: VAL-004 proves stale/fresh
  freshness behavior at the bumped version, and VAL-016 proves a deterministic,
  non-destructive full-schema rebuild with all triggers active (delete-DB +
  `fid rebuild` parity), so the no-bump trigger migrations replay correctly.
- **Refs:** VAL-004, VAL-016 (`qa/validation-contract.md`);
  `tasks/task-1-completion.md`, `tasks/task-2-completion.md`,
  `tasks/task-3-completion.md`, `tasks/task-4-completion.md`;
  `internal/db/projection_meta.go`; TRD §4.3.

### DEC-011: `task_completions` kind guard is task-or-subtask; reconcile §4.1

- **Date:** 2026-06-07
- **Context:** The Phase 2 plan flagged a spec ambiguity for the
  `task_completions` kind-guard trigger (VAL-010). `03_TRD_SCHEMA.md` §14.1
  defines the `task_completions` guard as "same" as `acceptance_criteria`,
  i.e. `issue_id` must reference a `task` **or** `subtask`. The `02_TRD.md` §4.1
  entity-inventory note read `(task only)`, contradicting §14.1. The planning
  brief directed the Developer to follow §14.1 and record the resolution in the
  Task 3 completion report, then revisit at spec sync.
- **Decision:** The shipped `0004` trigger guards `task_completions.issue_id` as
  **task or subtask**, following the authoritative trigger spec
  (`03_TRD_SCHEMA.md` §14.1). At wrap-up, reconcile the durable spec to the
  validated behavior rather than weaken the assertion: `02_TRD.md` §4.1 line for
  `task_completions` updated from `(task only)` to `(task/subtask only)`,
  matching the adjacent `acceptance_criteria` row and §14.1. (Same reconcile-to-
  validated-behavior pattern as Phase 1 DEC-009.)
- **Consequences:** `02_TRD.md` §4.1 edited (see Spec Sync Verdict). VAL-010
  behavior is unchanged and now matches the durable spec. No code change.
- **Refs:** VAL-010, QA-11/QA-13 (`qa/validation-contract.md`, `qa/qa-plan.md`);
  `tasks/task-3-completion.md`; `03_TRD_SCHEMA.md` §14.1; `02_TRD.md` §4.1.

## Accepted Developer Deviations (test-only, no spec impact)

Task 3 recorded two test-only deviations from its file list; both were reviewed
and accepted. Neither weakens a trigger or a validation assertion.

- Edited `internal/issues/issue_test.go` (outside the Task 3 file list): the
  pre-existing fixture inserted a task under a task, which the new VAL-009
  parent-kind backstop rejects. Only the parent fixture row was changed to an
  epic so the existing `Get()` behavior test stays valid under the new trigger.
- Hardened `internal/db/migrate_v01_tables_test.go` from `Steps(-2)` to
  `Migrate(1)` so the Task 1 Phase-1-schema regression test keeps targeting the
  Phase-1 schema boundary after migration `0004` was added.

## Validation Contract Audit

- All VAL-001..VAL-017 final outcomes are **PASS** with cited deterministic
  evidence (`qa/final-review.md` Validation Contract Outcomes; QA `VAL outcomes:
  17 PASS, 0 FAIL, 0 BLOCKED, 0 CARRY-FORWARD`).
- No FAIL, BLOCKED, or CARRY-FORWARD VAL to disposition.
- The `final-review.md` Requirement Compliance "Deferred Phase 2 scope" row is
  CARRY-FORWARD and maps to existing register items CF-004 (sequence-vs-
  scheduling trigger) and CF-005 (`activity_log` projection-only guard /
  `activity` op); both remain open in `decisions_active.md`.
- No validation-contract deviation was required; every assertion was proven as
  written, including both-state binary/conditional cases.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not a default skip.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 2 widened existing v0.1 schema scope; no identity/scope/constraint change. |
| `01_PRD.md` | verified-no-change | PRD-ISSUE-006, PRD-DEV-007, PRD-OUT-001/004, PRD-OPS-003/004 MET as written (final review Requirement Compliance); no requirement reworded. |
| `02_TRD.md` | updated | §4.1 `task_completions` relationship note reconciled `(task only)` → `(task/subtask only)` to match §14.1 and the shipped `0004` trigger (DEC-011). |
| `03_TRD_SCHEMA.md` | verified-no-change | Migrations `0002`-`0005` implement §1.2/§2/§5/§7/§9/§10/§10.1/§11 tables, the §2.2 `issues_enriched` view, and the §14.1/§14.2/§14.3/§14.4/§14.7 backstops as written; §14.1 was already correct (DEC-011 fixed the §4.1 echo, not §14.1). |
| `04_TRD_CONFIG.md` | N/A | Phase 2 added no `.fid/` layout or config surface. |
| `05_DESIGN_PLAN.md` | verified-no-change | Enriched `fid issue get` realized per §5.2/§6.1 JSON shape; no design-language change. |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | `fid issue get --format json` enriched fields match the documented shape; TOON/markdown/YAML and `--section` remain Phase 3 / v0.1b. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | No flow change; widened schema follows the documented journal→projection→rebuild flow. |
| `08_QUALITY_PLAN.md` | verified-no-change | COV-001 (`internal/journal` 80.6%, `internal/issues` 85.1%) and COV-002 (`cmd/fid` 70.5%) held; no gate/process change this phase. |
| `09_AI_WORKFLOW.md` | verified-no-change | Workflow ran as specified (validation contract → QA plan → dev/QA loop → mid-phase review → final review → human checklist); no workflow change. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 2 delivered the v0.1a "widen the schema vertically" slice per the slicing note; deferrals (CF-004/CF-005) recorded, no re-slice required. |

## Triage-Log Audit

- [x] Walked `triage-log.md`: **no TRG-NNN entries were recorded in Phase 2.**
  The Phase Runner / TL-Lite triage path was not exercised — no `VERDICT:
  RESUME` or `VERDICT: BLOCK` events occurred. The mid-phase review (review-3)
  returned CONTINUE with zero findings and raised no triage.
- [x] No RESUME verdicts to confirm in hindsight; none to flag as
  should-have-escalated.
- [x] Autonomous TL-Lite resume remains permanently ENABLED (DEC-002); there is
  no warm-up tally to maintain.

## Workflow Notes

Answer in 1-3 sentences each. Completed at phase wrap-up.

- **What worked well?** The Phase-1 DEC-007 carry-in was applied at contract-
  authoring time: VAL-003 named package-level `internal/journal` tests and
  VAL-007 named `internal/issues` tests, so coverage held (journal 80.6%,
  issues 85.1%, cmd/fid 70.5%) with no end-of-phase coverage remediation cycle —
  exactly the planning seam Phase 1 paid for. The "exercise via journal effects,
  not direct SQL" rule gave every new table real end-to-end signal.
- **What was inefficient?** Very little. The one pre-identified spec ambiguity
  (`task_completions` kind guard) was resolved cleanly via the planning note and
  closed at spec sync (DEC-011) without a rework loop.
- **Were tasks the right size?** Yes. Four migration-aligned vertical slices
  (0002 tables, 0003 enriched view + CLI, 0004 kind triggers, 0005 ref/scope
  triggers + phase-end) were each one focused session with clean red/green TDD
  and no blockers.
- **Did the QA review cadence feel right?** Yes. A single mid-phase review after
  Task 3 (covering VAL-001..012) plus the phase-end final review was sufficient
  for a four-task phase with zero blocking findings; an earlier review was not
  warranted.
- **Were there token-efficiency opportunities missed?** Minor. The final review
  was a consolidation of existing task evidence (no new command runs), which was
  efficient; no material waste observed.

### Deep Retrospective Recommendation

**No retro warranted** beyond this note. Phase 2 ran clean: zero blocking,
major, or improvement findings; zero escalations; zero triage events; all 17
VALs PASS first time; human checklist 5/5 PASS. The only durable lesson —
apply COV-001 package-level evidence at contract time — was carried in from
Phase 1 DEC-007 and demonstrably prevented a repeat coverage cycle. Continue the
practice; nothing new to extract.
