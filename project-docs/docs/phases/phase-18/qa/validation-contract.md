# Phase 18 Validation Contract - Workflow Tables Migration + Engine Data-Source Swap + Gate Seams

> Phase: 18 - Workflow Engine Migration (v0.2a, slice 1 of 3)
> Author: Tech Lead (draft)
> Last refined: 2026-07-18 by Tech Lead

## Purpose

This file is the finite set of validation assertions that must be true when
Phase 18 is complete. The validation contract specifies **what must be proven**;
the QA plan (QA Lead) specifies **how to prove it**.

Phase 18 converts the v0.1 hard-coded workflow state machines into
data-configurable tables and repoints the transition engine to resolve from
those tables, while preserving v0.1 behavior exactly. It also reserves the
gate-obligation seam (DEC-066) without any receipt behavior. No new user-facing
command ships in this phase; behavior is observed through existing commands
(`fid issue transition`, `fid issue get/list`, `fid maintenance check`,
`fid rebuild`) plus migration/replay tests.

If a required validation assertion cannot be satisfied, stop and escalate
instead of weakening the assertion or silently adapting the spec.

## Numbering

- IDs use `VAL-NNN`, scoped to Phase 18, sequential from `VAL-001`.
- Retired or superseded assertions stay in the table with status and rationale;
  IDs are not reused.

## Assertions

| ID | Source requirement | Behavior or invariant | Evidence method | Validator | Owning task/window | Pass criterion | Fail criterion | Status | Out-of-scope notes |
|---|---|---|---|---|---|---|---|---|---|
| VAL-001 | TRD-ARCH-009; TRD_SCHEMA §3 | Migration creates `workflows`, `workflow_states`, `workflow_transitions`, and `workflow_transition_requires` with the specified columns, including `workflow_transition_requires.params_json`. | Schema-introspection test (`PRAGMA table_info` / DDL assertions) on a migrated DB. | QA Lead | task-1 | All four tables and their required columns (incl. `params_json`, PKs, and the `workflow_transition_requires`→`workflow_transitions` FK) are present. | Any table or required column is missing or misnamed. | Draft | Custom-workflow authoring CLI is Phase 19; not asserted here. |
| VAL-002 | TRD §4.6; roadmap v0.2a "seed defaults verbatim" | Seeded `workflow_states` reproduce the v0.1 state→category maps for `epic-default`, `task-default`, and `subtask-default` exactly. | Query `workflow_states`; compare against frozen v0.1 fixtures. | QA Lead | task-1 | Every v0.1 state is present with the identical category; no extra or missing states for any of the three workflows. | Any state missing, added, or mapped to a different category. | Draft | Escalation-default workflow is Phase 19. |
| VAL-003 | TRD §4.6 | Seeded `workflow_transitions` reproduce the v0.1 allowed edges for all three workflows, including an explicit `<state> → cancelled` edge for every non-cancelled state, and excluding `cancelled → cancelled` and all v0.1-disallowed edges. | Query `workflow_transitions`; compare against frozen v0.1 fixtures, asserting both present edges and representative absent edges (e.g. `todo→done`, `cancelled→cancelled`). | QA Lead | task-1 | Present edges match v0.1 exactly; each non-cancelled state has an explicit `→cancelled` edge; disallowed and `cancelled→cancelled` edges are absent. | Any allowed edge missing, any disallowed/`cancelled→cancelled` edge present. | Draft | — |
| VAL-004 | TRD §4.6; PRD-WF-007–010 | Seeded `workflow_transition_requires` reproduce the v0.1 per-transition rules, including parameterised roles in `params_json` (`all_acs_verified {"roles":["developer"]}` on `task-default` `in_progress→dev_complete`; `{"roles":["qa_lead"]}` on `qa_review→qa_pass`; `has_completion_report` on `in_progress→dev_complete`; `all_child_tasks_done_or_cancelled` on `epic-default` `in_progress→wrapped`). | Query `workflow_transition_requires`; compare rule names and parsed `params_json` against frozen v0.1 fixtures. | QA Lead | task-1 | Every v0.1 rule (name + params) is present on the correct edge; no extra or missing rules. | Any rule missing, added, on the wrong edge, or with different params. | Draft | New v0.2b/c rules are seam-only (VAL-005). |
| VAL-005 | PRD-WF-014/015 (seam); DEC-066 | The gate-obligation seam is present with no behavior: `params_json` structurally accommodates a gate ID/version obligation shape; rule names `required_gates_current` and `has_completion` are reserved (recognized as known-but-unimplemented); no seeded transition uses them; and no transition path evaluates any receipt/gate logic. | Test: no seeded row references the reserved names; a `workflow_transition_requires` row carrying a gate-obligation-shaped `params_json` is structurally accepted; a transition never triggers receipt/gate evaluation code. | QA Lead | task-1 | Reserved names recognized but unseeded; gate-obligation `params_json` shape validates structurally; zero receipt/gate behavior is reachable from any transition. | A reserved rule is seeded/enforced, the shape is rejected, or any receipt/gate logic executes. | Draft | Actual receipt/execution behavior is v0.2b; explicitly excluded. |
| VAL-006 | TRD-ARCH-009; TRD_SCHEMA §2.1 | After migration `issues.workflow_id` is a real FK to `workflows(id)` and `(issues.workflow_id, issues.status)` is a composite FK to `workflow_states(workflow_id, state)`: an insert/update with an unknown `workflow_id` or an invalid `(workflow_id, status)` pair is rejected, while a valid pair is accepted. | FK-enforcement test asserting both the rejected (invalid) and accepted (valid) paths; `PRAGMA foreign_key_check` returns no rows post-migration. | QA Lead | task-4 | Invalid workflow_id / (workflow_id,status) rejected by FK; valid pair accepted; `foreign_key_check` clean. | An invalid pair is accepted, a valid pair rejected, or `foreign_key_check` reports rows. | Draft | Enum CHECK constraints are replaced, not supplemented. |
| VAL-007 | TRD-ARCH-008; TRD_SCHEMA §2, §14 | The issues-table rebuild preserves all existing issue rows and re-creates the `issues_enriched` view, the kind/sequence-scheduling/kind-immutability triggers, and the issues indexes; a migrated v0.1 workspace returns identical `fid issue get`/`fid issue list` output for pre-existing issues. | Migrate a seeded v0.1 DB; diff `fid issue get`/`list` output pre/post migration; assert view + trigger existence. | QA Lead | task-4 | Row count and issue reads are identical pre/post; view and all issues triggers exist post-migration. | Any row lost/changed, missing view/trigger, or differing read output. | Draft | New behavior of the workflow tables is covered by other VALs. |
| VAL-008 | v0.2 exit criterion ("v0.1 defaults migrate cleanly"); PRD-WF-001–010 | Every v0.1 transition behaves identically once the engine resolves from the DB tables: allowed edges succeed; disallowed edges are rejected with the same `InvalidTransitionError` message; blocked rules produce the same `BlockedTransitionError` naming the same unmet conditions, entity IDs, and missing roles; cancel-from-any-non-cancelled-state succeeds. | Characterization tests plus `fid issue transition` e2e across representative allowed, invalid-edge, blocked-rule, and cancel edges, comparing outcomes and error text to v0.1. | QA Lead | task-3 | All four categories (allowed, invalid-edge, blocked-rule, cancel) produce v0.1-identical outcomes and messages. | Any edge outcome or message differs from v0.1. | Draft | Custom-workflow edges are Phase 19. |
| VAL-009 | PRD-CLAIM-005/006 | Started-category claim enforcement is unchanged using DB-resolved categories: a transition into a started-category state still requires a live caller claim, and a transition out of the started category still auto-releases the claim. | Transition claim tests covering both the required-claim (into started) and auto-release (out of started) paths. | QA Lead | task-3 | Into-started transition rejected without a live claim and accepted with one; out-of-started transition releases the claim. | Either the claim requirement or the auto-release regresses. | Draft | — |
| VAL-010 | PRD-OPS-005–008; TRD_SCHEMA §14.11 | `fid maintenance check` validates each issue's `(workflow_id, status)` pair against the workflow tables (not the Go map): it stays clean on a healthy migrated workspace and reports a structural error for a seeded invalid pair. | Maintenance test on a healthy DB (no finding) and on a DB with a seeded invalid `(workflow_id, status)` (finding reported); `--warn-as-error` exits non-zero on the seeded violation. | QA Lead | task-5 | Healthy workspace produces no workflow-pair finding; seeded invalid pair produces a deterministic finding and non-zero `--warn-as-error` exit. | Healthy workspace flags a false positive, or a seeded invalid pair is not detected. | Draft | Other maintenance checks unchanged. |
| VAL-011 | PRD-OPS-003–004; DEC-066 | `fid rebuild` reproduces the fully-migrated projection deterministically, and delete-and-rebuild is non-destructive: after migrate + replay the workflow tables, issues FKs, and issue states match, with per-table parity between the original and the rebuilt DB. | Delete-and-rebuild per-table parity test over the migrated schema (workflow tables + issues). | QA Lead | task-5 (phase-end) | Rebuilt DB is per-table identical to the pre-delete DB across workflow tables and issues; rebuild is deterministic across runs. | Any per-table divergence or non-deterministic rebuild output. | Draft | — |
| VAL-012 | TRD §4.3; `projection_meta` | `ProjectionSchemaVersion` is bumped so a DB projected under the prior projection-schema version is treated as stale and rebuilt/caught-up rather than served as-is. | Test that stores the prior `projection_schema_version` and asserts the next read forces a rebuild/catch-up to the current version. | QA Lead | task-1 | Stored prior version triggers a rebuild/catch-up; post-read version equals the current `ProjectionSchemaVersion`. | A stale-version DB is served without rebuild, or the version is not bumped. | Draft | — |
| VAL-013 | DEC-066; TRD-ARCH-002–003 | v0.1 journal history is never rewritten by this migration: existing `issue.created` / `issue.transitioned` events replay unchanged into the migrated schema, and journal files on disk are byte-identical before and after migrate + rebuild. | Golden replay fixture: replay a v0.1 journal into the migrated schema and assert the expected projection; assert journal file bytes are unchanged by migrate/rebuild. | QA Lead | task-5 | Golden v0.1 journal replays into the migrated projection with the expected result; journal bytes unchanged. | Replay fails/diverges, or any journal file is modified. | Draft | New event types for custom workflows are Phase 19. |
| VAL-014 | Roadmap v0.2a; TRD §4.6 | The DB-backed workflow registry loaded from the seeded tables is behavior-equivalent to the frozen v0.1 Go definitions (states, categories, allowed edges, and per-edge rules) for all three workflows. | Characterization test comparing the loaded registry's edges/categories/rules against the frozen v0.1 Go definitions for `epic-default`, `task-default`, and `subtask-default`. | QA Lead | task-2 | Registry output equals the frozen v0.1 definitions for every state, category, edge, and rule across all three workflows. | Any divergence between the loaded registry and the v0.1 definitions. | Draft | Consumers still reading the Go map before their swap are covered by their own tasks. |

## Coverage Rules

- Every in-scope requirement (TRD-ARCH-009, TRD §4.6, PRD-WF-014/015 seam,
  DEC-066 compatibility, and the "v0.1 defaults migrate cleanly" v0.2 exit
  criterion) is referenced by at least one `VAL-NNN` above.
- Every `VAL-NNN` is owned by exactly one task or `phase-end` window and has
  exactly one validator role (all QA Lead — this phase is fully
  machine-checkable; no Human perception assertion applies).
- Out-of-scope for Phase 18 (own VALs in later phases): `fid workflow
  define/assign/list` and escalation-default workflow (Phase 19);
  `role_permissions` migration and `fid role create`/`fid policy set`
  (Phase 20); any gate execution/receipt/evidence behavior (v0.2b).

## Outcomes

Phase-end final review records every `VAL-NNN` as `PASS`, `FAIL`, `BLOCKED`, or
`CARRY-FORWARD` with cited evidence. `PASS` without evidence is invalid.

## Revision Log

| Date | Assertion | Change | Reason | Decision reference |
|---|---|---|---|---|
| 2026-07-18 | VAL-001..014 | created | Phase 18 planning | DEC-066 |
