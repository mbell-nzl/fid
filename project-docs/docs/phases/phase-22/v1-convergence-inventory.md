# Phase 22 V1 Convergence Inventory

> ACCEPTED under DEC-188; inventory base 89d964a84dbdd297fdf472d6f78b001e2efac702.
> Git history is the archive; no legacy compatibility mode is promised.

## Exact Live Dispositions

| Surface / exact paths | Current fact | Disposition / first owner / executing gate |
|---|---|---|
| `internal/claims/generic.go`; `cmd/fid/claim.go`; `cmd/fid/generic_claim_test.go` | Expired acquisition can replace a lease without executor recovery | Replace with explicit expiry and unresolved-Run fence, task 3; ordinary + tagged |
| `internal/workflow/run.go`, `run_read.go`; `cmd/fid/run.go`, `generic_run_test.go`, `run_test.go` | Caller-selected finish outcome; no predecessor/executor/version | Binding change task 3; derived close/abandon task 8; remove finish leaf and current run.finished acceptance task 9; both gates each |
| `internal/workflow/activation.go`; `cmd/fid/activation.go`; `internal/cli/codes.go`; activation tests | Single happy activation; post-commit recovery-unavailable | Resource/predecessor/domain recovery task 10; new leaves/errors task 11; both gates |
| `internal/workflow/envelope.go`, `run.go`, `run_read.go`; `internal/evidence/execution.go` | Immutable budgets with null consumption/remaining and unavailable action | Preserve immutable policy, replace with observed sample projection and atomic limits, task 7; both gates |
| `internal/workflow/resume.go`; `cmd/fid/resume.go`; `task18_completion_resume_test.go` | Stored target and string actions without recovery facts | Current target, predicate-derived actions, Stops/Attention/budgets/lineage, task 12; both gates |
| `internal/export/export.go`; `cmd/fid/export.go`; export tests | v1 archive create/validate; hard-coded schema 14 | Schema reporting task 1, v2 format task 13, import task 14/15; retain only non-mutating v1 validator; both gates |
| `internal/workflow/capabilities.go`; `internal/export/export.go`; capability fixtures | Four Phase-22-unavailable rows | Task 17 final binary matrix; three true, Claims false with Phase 24 owner; no obscured Attention identifier workaround |
| `migrations/0002_v01_tables.up.sql` escalations table; `0005_v01_ref_triggers.up.sql` activity trigger; `0010_v02a_escalation_default_workflow.up.sql` seed | Historical schema producers | Latest schema removes through exact 0016, task 16; old files only migration-scoped fixtures until Phase 25 consolidation |
| `triage_events.esc_id`, `concerns.escalated_to_esc_id` | Foreign keys prevent dropping escalations | Remove only those columns in task 16; preserve triage/Concern lifecycle for Phase 24/23 respectively |
| `internal/maintenance/check.go`, `errors.go`, `test_helpers_test.go`; `internal/db/escalation_workflow_test.go`, `ref_triggers_test.go`, `migrate_v01_tables_test.go`; matching journal fixtures | Latest versus historical escalation reference assumptions | Task 16 changes latest-schema expectations and maintenance refs; historical tests retain explicit earlier migration scope |
| `internal/workflow/rules.go` RuleLinkedToDecision escalation comment; registry/rules tests | Retired escalation seed rationale | Task 16 removes escalation-only authority; retain generic project Decision scope only where still explicitly owned by Phase 23 |
| `internal/activity/`; `internal/journal/project_effects.go`; maintenance activity readers | Event-derived activity rows, some frozen V1 attribution | Retain read-only projection; prohibit Stop/Run takeover truth in activity; task 16 classifies current writers/readers, task 17 audits integrity |
| `internal/issues/get.go`, `issue.go`; `issues_enriched`; `triage_events`; `idx_triage_issue_created` | Legacy triage-status inspection still live | Frozen, no recovery authority. Phase 24 removal; task 16 only converges removed escalation refs, not triage semantics |
| `cmd/fid/integration/` and `testdata/p21_conformance/`; `cmd/fid/task*` tests | Historical phase/task names and exact schema/command/availability fixtures | Subject-preserving convergence belongs to each invalidating task above; full tagged gate runs from task 1 |
| `cmd/fid/gate_run_test.go` TestGateRunRecordsExactNonSuccessTerminalOutcomes and TestGateRunHelperProcess | Signal helper can exit success; 1 ms timeout can precede spawn; both predate Phase 22 | Task 8 owns synchronized signal/timeout fixtures alongside cancellation; ordinary + focused repeated + two full benchmark gates. Earlier failure gets bounded same-task scope amendment per DEC-186, never a waiver |
| `internal/workflow/activation.go` accepted actions and `internal/workflow/envelope.go` grants | Phase 21 13-name set includes run.finish | Exact staged additions under public-interface section 2.1; task 9 removes run.finish and task 18 audits the final 23-member set through both gates |

No escalation CLI leaf exists at the planning base. Task 16 proves it remains
absent rather than creating an alias to Stop/Attention. Searches must distinguish
product escalation state from `project-docs` delivery TRG/ESC artifacts, which
are outside product removal scope.

## Provenance Anchors

Introducing commits observed with `git log --diff-filter=A`, all before planning:

| Artifact | Introducing commit / phase task |
|---|---|
| generic Claim implementation | 0b67c83 / phase-21/task-10 |
| generic Run implementation | 5ca732f / phase-21/task-11 |
| Activation implementation | ad9098f / phase-21/task-4 |
| Attempt/Receipt implementation | b694991 / phase-21/task-15 |
| Signal helper and timeout_ms=1 fixture | b69499116e70d4c88423a30daa2010a4ff221f14 / phase-21/task-15; git log -S and blame observed by TL |
| Signal-case timeout extension only | 76ad781727b9e8545e069c4ac0e9f926012e6af9 / phase-21/task-21; does not change timeout_ms=1 (corrects QA R-14 prose) |
| resume implementation | 2313235 / phase-21/task-18 |
| export implementation | 35a75ba / phase-21/task-19 |
| escalation workflow test/seed | a509d59 / phase-19/task-1 |
| V1 tables migration | 5b1467f / phase-2/task-1 |
| real-process TestMain | 1edb689 / phase-20/task-3g |

These establish artifact origin, not the provenance of every retained assertion.
Developer reports establish per-row provenance before claiming a particular row
is inherited or displaced. New failures are not classified by filename.

## Retained With Expiry

- `task_completions`, `ac_verifications`, `concerns`, `decisions`, role-specific
  judgment columns: frozen readers only; Phase 23.
- `triage_events` apart from removed escalation FK, triage-status projection,
  broader Claim/Run shapes and multi-Issue revision: Phase 24.
- `agent_runs`, historical stats, old migrations and strictly migration-scoped
  Event fixtures, general historical residue: Phase 25 qualification.
- `internal/sqlquery/sqlquery.go`: diagnostic only. No planned production edit;
  CF-018's split fires if a causal change reopens it. Its tests may converge
  without firing the production-file trigger.

Replacement-bound removal cannot be deferred to task 18 or Phase 25 merely
because a test breaks. Any newly found surface gets a recorded exact owner and
expiry before dispatch, with a gate that actually exercises it.
