# Phase 23 V1 Convergence Inventory

> DRAFT, pending Human acceptance with the slice contract.
> Inventory base `0195897a347dfe1872db0f7be3800c5c3ccdefcf`.
> Git history is the archive; no legacy compatibility mode is promised.

## 1. Exact Live Dispositions

| Surface / exact paths | Current fact | Disposition / first owner / executing gate |
|---|---|---|
| `internal/workflow/envelope.go` single hard-coded `gate-evidence` Obligation and `Obligations[0]` readers; `cmd/fid/envelope.go`, `envelope_test.go`, `envelope_fixture_test.go` | Exactly one Obligation of one type; evaluation returns one row | Replace with 1..16 typed Obligations and per-Obligation evaluation, task 2; ordinary race/coverage plus tagged integration |
| `migrations/0013_p21_operational_kernel.up.sql` `envelope_obligations`, `work_claims.mode`; `migrations/0015_p22_recovery.up.sql` `stops.code`, `attention_items.code` | Baseline CHECK domains exclude Phase 23 values | Exact textual baseline changes in `schema.md` section 1, task 1; both gates |
| `internal/claims/generic.go`, `cmd/fid/claim.go`, `cmd/fid/generic_claim_test.go` | Claim mode is `mutation` only | Add non-mutating `review` mode and its adoption fence, task 3; both gates |
| `internal/evidence/evidence.go`, `execution.go`, `policy.go`; `cmd/fid/evidence.go`, `evidence_test.go` | Evidence is receipt-only with `available`/`missing` artifact dimension | Widen to three evidence types and `not-applicable`, tasks 6 and 10; both gates |
| `internal/workflow/budget.go`, `budget_policy_test.go`, `run_budget.go`, `recovery_budget.go`; `cmd/fid/live_budget_test.go` | Units exclude `review-cycles`; exhaustion action is `budget-exhausted` only | Add `review-cycles` and `review-rework-exhausted`, task 6; both gates |
| `internal/workflow/evaluation.go`, `resolve.go`, `rules.go`; `internal/workflow/attention_test.go` | Seven Stop codes, six Condition types, no judgment facts | Add three Stop codes and four Condition types, tasks 6, 9, 10, 11; both gates |
| `internal/issues/ac_verify.go`, `ac.go`, `read.go`; `cmd/fid/issue_ac.go` verify leaf; `cmd/fid/issue_ac_test.go`; `internal/issues/ac_test.go` | Fail-closed `issue ac verify` with role-keyed `latest_verifications` projection | Remove leaf, writer, and role projection, task 15; both gates. `acceptance_criteria` add/list are retained |
| `internal/issues/concern_raise.go`, `concern_address.go`, `concern_list.go`; `cmd/fid/concern.go` raise/address; `cmd/fid/concern_test.go` | Fail-closed completion-linked V1 Concern surface | Replace with generic Concerns in task 8, remove V1 leaves, writers and readers in task 15; both gates |
| `internal/issues/detail.go`, `get.go`, `list.go`, `issue.go`, `transition_rules.go`; `internal/issues/detail_test.go`, `list_test.go`, `issue_test.go` | Issue detail and list read `task_completions`, `ac_verifications`, and V1 `concerns` | Re-point to generic Concerns, drop completion and verification projections, task 15; both gates |
| `internal/workflow/rules.go` `all_acs_verified`, `has_completion_report`, `concerns_addressed`, `linked_to_decision`; `rules_test.go`; `internal/workflow/registry.go` `RequiredRoles`; `internal/workflow/workflow.go` | Role-parameterised and completion-bound transition rules; `linked_to_decision` reserved | Delete the first three rules and the roles parameter, implement `linked_to_decision` against `decision_refs`, tasks 9 and 15; both gates |
| `internal/requirements/` all files; `cmd/fid/requirement.go`; `cmd/fid/requirement_test.go`; `internal/journal/current_event_catalog.go` `requirement.synced` | Markdown-parsing requirement authoring with `last_synced_hash` | Replace with `requirement_refs`/`requirement_issue_links` and remove the package, leaf, Event, and rule, task 12; both gates |
| `internal/links/registry.go`, `add.go`, `list.go`, `errors.go`, `types.go`; `cmd/fid/link.go`; `cmd/fid/link_test.go`, `link_limit_test.go`; `cmd/fid/e2e_v01d_link_test.go` | `links` accepts `requirement` endpoints and `implements`/`implemented_by` | Restrict `links` to Issue endpoints; requirement traceability moves to `requirement_issue_links`, task 12; both gates. This removes the second authoritative model for one fact |
| `internal/maintenance/errors.go`, `warnings.go`, `operational.go`, `check.go`; `internal/maintenance/*_test.go` | Integrity, `stale_verifications`, `telemetry_gaps`, and `frozen-state-present` checks read the removed tables | Re-point to the Phase 23 tables and delete the checks whose subject is removed, task 15; both gates |
| `internal/stats/stats.go`, `stats_test.go`; `cmd/fid/stats.go`, `stats_test.go` | Historical stats join `agent_runs.completion_id` | Converge only the dropped column; `agent_runs` and historical stats stay Phase 25, task 15; both gates |
| `migrations/0003_v01_enriched_view.up.sql` `issues_enriched`, `idx_concerns_issue_status`; `migrations/0005_v01_ref_triggers.up.sql` link and activity triggers; `migrations/0016` activity trigger | View and triggers reference V1 `concerns`, `requirements`, `decisions`, and requirement link endpoints | Rebuilt by migration 0018 in task 15; both gates. Historical migration tests keep their explicit earlier scope |
| `internal/db/migrate_v01_tables_test.go`, `issues_enriched_test.go`, `ref_triggers_test.go`, `kind_triggers_test.go`, `issues_workflow_fk_test.go`, `operational_kernel_migration_test.go` | Latest-schema expectations include the removed tables | Latest-schema expectations change in task 15; explicitly earlier-migration-scoped tests keep their historical assertions; both gates |
| `internal/journal/v01_tables_test.go`, `phase2_rebuild_test.go`, `replay_test.go`, `rebuild_scale_test.go`, `escalation_removal_test.go`, `project_effects.go` | Effect rules, fixtures, and scale corpora carry V1 judgment tables and the `requirement.synced` Event | Progressive fence in task 1, per-Event deletion by its owner, corpus convergence in tasks 12 and 15; both gates |
| `cmd/fid/integration/testdata/p21_conformance/` `concern.raise.json`, `concern.address.json`, `issue.ac.verify.json`, `requirement.sync.json`, `manifest.json`; `p21_conformance_test.go`; `exit_contract_test.go`; `cmd/fid/contract_surface_test.go`; `cmd/fid/task9_contract_test.go` | Pinned fail-closed capability-unavailable fixtures for the four judgment leaves | Replaced by the Phase 23 available-operation fixtures in tasks 13 and 17; the negative rows are deleted only when their replacement passes; tagged integration gate |
| `internal/issues/judgment_unavailable_test.go`; `internal/issues/capability_unavailable.go` | Fail-closed judgment guard and its tests | Removed when tasks 4, 8, and 10 deliver real writers; final absence proved in task 17; both gates |
| `internal/workflow/capabilities.go`; `internal/export/archive_profiles.go`; `cmd/fid/capabilities_test.go`; `cmd/fid/integration/p22_capability_matrix_test.go` | Four Phase-23/24/25-unavailable rows and `human_outcome_resolution:false` | Final binary matrix in task 17 after every replacement and removal proof; both gates |
| `internal/workflow/resume.go` `ResumeFindings{Available:false}`; `cmd/fid/resume.go`, `resume_human.go`, `resume_human_guidance.go`; `internal/workflow/resume_test.go`; `cmd/fid/integration/p22_resume_*_test.go` | Findings placeholder and Phase 22 resume sections | Real Findings and judgment projections, task 14; both gates |
| `internal/export/export.go`, `archive_profiles.go`; `cmd/fid/export.go`; `cmd/fid/integration/p22_portable_*_test.go` | Manifest v2 pinned to projection schema 16 | Schema 17 during implementation and 18 at phase end, tasks 1, 15, 16; both gates |
| `internal/activity/`; `internal/journal/project_effects.go` activity writers | Closed ActivityEvent set includes `ac-verified`, `concern-raised`, `concern-addressed`, `task-completed` | Replace with the Phase 23 event names in `public-interface.md` section 2, tasks 8 and 15; both gates |

No product escalation, QA-plan, validation-contract, template, retrospective, or
carry-forward table exists at the planning base; task 17 proves they remain
absent rather than creating an alias. Searches must distinguish product state
from `project-docs` delivery artifacts, which are outside product removal scope.

## 2. Provenance Anchors

Introducing commits are established by the owning Developer with
`git log --diff-filter=A`, `git log -S`, and `git blame` before any row is
called inherited or displaced. Labels, filenames, and embedded task numbers
never assign ownership (DEC-162/165/166/194). The convergence rules in force:

- DEC-167/173: a task's required gates execute every owned scope, and its scope
  covers every path its gates walk. Amend an incomplete assignment before
  dispatch, never after the failure.
- DEC-175: a latent implementation defect first observable under a new gate
  belongs to the observing task with minimum causal source and regression scope.
- DEC-189/193: a nondeterministic or calendar-dependent retained test is
  neither convergence nor DEC-175; it is repaired at its source by the observing
  task with a bounded scope amendment.
- DEC-192: a progressive availability fence is removed key by key by the task
  that owns the fenced Event; the last owner removes the empty map and guard.
- DEC-195: an accepted clause that adds a closed-domain member requires a
  matching schema row; amend the schema, never the assertion.

## 3. Retained With Expiry

- `triage_events`, the triage-status projection in `issues_enriched`, broader
  Claim and Run shapes, and multi-Issue revision: Phase 24.
- `agent_runs` minus its dropped `completion_id`, `epic_metadata`, historical
  `internal/stats` behavior, old migrations, migration-scoped Event fixtures,
  `v01_parity_test.go`, and general historical residue: Phase 25.
- `internal/sqlquery/sqlquery.go`: diagnostic only. No planned production edit;
  CF-018's split fires only if a causal change reopens it.
- `acceptance_criteria` add and list: retained target behavior, not residue.

Replacement-bound removal cannot be deferred to task 17 or Phase 25 merely
because a test breaks. Any newly found surface gets a recorded exact owner,
expiry, and an executing gate before dispatch.
