# Phase Brief - Phase 23 (Approved)

## Phase Goal

Deliver the Review and Human judgment slice: an independent Reviewer records
coverage and Findings against an exact Envelope revision and repository
identity; Builder rework preserves review independence; protected or subjective
Obligations receive provenanced Human Outcomes; completion is re-evaluated from
those facts. It also replaces V1 requirement authoring with references and
removes the role-specific verification columns, QA and TL composites, and
role-bound completion behavior this slice owns. Success is a Human who can tell,
from fid output alone, what a Review resolved, what is merely addressed, what
still blocks, and which judgment unblocked it.

## Status

**Approved, not started; two contract amendments still open.** Branch
`phase-23` at planning base `0195897a347dfe1872db0f7be3800c5c3ccdefcf`. The QA
Lead's Heavy-mode plan returned thirteen Must-fix revisions and sixteen
improvements; all are accepted and applied by the Tech Lead on 2026-09-15, with
the disposition table and reasoning in
[phase decisions](../phases/phase-23/decisions.md). On 2026-09-19 the Human
approved DEC-223 and DEC-227 and confirmed DEC-225. Next action: the **Human**
decides DEC-224 and DEC-226, after the requested detail.

| DEC | Amendment | State |
|---|---|---|
| DEC-223 | Link-preserving `evidence_records` rebuild in 0017/0017-down | **Approved 2026-09-19** |
| DEC-224 | Four named Events raise `human-outcome-required`, two gain effect rows | **Open** - detail supplied 2026-09-19 |
| DEC-225 | Independence conflicts with the mutating side only | **Confirmed 2026-09-19** |
| DEC-226 | Code producers, precedence, denials; one code demoted | **Open** - four separable parts, detail 2026-09-19 |
| DEC-227 | Explicit `-timeout 30m` on the two long gates plus measured wall-time budgets | **Approved 2026-09-19** |

Once DEC-224 and DEC-226 are decided, the QA Lead runs the `qa-plan.md` section
11 readiness refresh, then the Runner's readiness check. No task dispatches
before both. A rejected amendment reverts its dependent validation rows to
`Draft` and needs replanning.

## Operating Mode

**Heavy** (`../../specs/09_AI_WORKFLOW.md` section 4). Four entry criteria met:
persistence and migration change; public-interface change (eighteen leaves,
seven grant actions, a new closed failure-code set); an authority and
trust surface; replacement-bound removal of fail-closed V1 behavior.

## Accepted Contract Set

Accepted under DEC-222 and the only implementation authority, all under
`../phases/phase-23/`: `contracts.md`, `public-interface.md`, `schema.md`,
`execution-scope.md`, `v1-convergence-inventory.md`, and `contract-coverage.md`.
Each carries the 2026-09-15 amendments inline, marked with their DEC and with
their current approval state.

## Task Index

Full specs and reports: `../phases/phase-23/tasks/`. All 18 are **Pending** and
strictly sequential: each depends on the one before, task 1 also on
DEC-222/DEC-223, task 10 also on DEC-224, task 12 also on the checkpoint after
task 11.

1 schema baseline; 2 typed multi-Obligation Envelope; 3 review Claim mode and
Run; 4 Review Records and computed independence; 5 Finding dispositions; 6
review evidence, cycles, exhaustion; 7 freshness and supersession; 8 Concerns
and Deviations; 9 Decision References and protected Decisions; 10 Human
Outcomes; 11 Attention resolution; 12 requirement references and one
traceability model; 13 CLI, codes, idempotency; 14 Human projections; 15
migration 0018 and V1 removal; 16 portable recovery; 17 capability matrix; 18
conformance and wire audit.

## Validation Contract

`../phases/phase-23/qa/validation-contract.md` holds **79** assertions (75
deterministic, 4 Human), all `Ready` after the 2026-09-15 refinement. Ownership
by range: VAL-001..005/075 schema baseline, fence, parity (tasks 1, 15, 16);
VAL-006..013 Obligations, evaluation, review Claims (2, 3); VAL-014..023 Review
Records, independence, Findings (4, 5); VAL-024..033/077 review evidence,
cycles, exhaustion, freshness (6, 7); VAL-034..046/074 Concerns, Deviations,
Decisions, Outcomes (8, 9, 10); VAL-047..054/078 Attention, Conditions,
requirements, links (11, 12, 17); VAL-055..065/076 CLI, Human projections,
removal, portability (13-16); VAL-066..069/079 availability, conformance audit,
gate budgets (17, 18); VAL-070..073 Human comprehension and product judgment at
phase end, with runnable fixtures from task 14. Validators are QA Lead or QA
Reviewer except those four Human rows.

## QA Artifacts

All under `../phases/phase-23/qa/`. Validation contract refined 2026-09-15, all
79 rows `Ready`. QA plan authored but **NOT READY** until its section 11
readiness refresh, which Heavy mode requires to be freshly measured. Mid-phase
`review-M.md` for windows 1-3, 4-6, 7-9, 10-12, 13-15, 16-18;
`checkpoint-after-task-11.md` must PASS before task 12; `human-fixtures.md` and
the runnable VAL-070..072 fixtures come from task 14 and are regenerated at
task 18; final review and human checklist pending.

## Active Constraints

- Nothing dispatches before DEC-224 and DEC-226 are decided, the QA readiness
  refresh passes, and the Runner's readiness check passes.
- Every task runs the tagged integration gate from task 1 with an explicit
  `-timeout 30m`, and records the `cmd/fid/integration`, `cmd/fid`, and
  `internal/workflow` package times (DEC-227). Hosted CI does not run that gate
  (CF-041), so a green PR is never phase-level evidence.
- Removal is replacement-bound: task 15 removes nothing whose replacement is not
  already green, and defers nothing because a test broke.
- Independence is computed from Claim and Run history, over mutating Claims and
  their Runs only (DEC-225); no role string may decide it. A Human Outcome is
  never stored, evaluated, or rendered as deterministic proof.
- Phase 23 adds no denial Event: a Stop-blocked request commits nothing and
  returns the existing Stop and Attention IDs with a null EventRef (DEC-226).
- Envelope revision stays 1 and the primary Issue single; parallel Claims,
  multi-Issue Envelopes, and `retry-exhausted` remain Phase 24.

## Interfaces / Subsystems In Scope

`internal/review/` (new) owns Review, Finding, Concern, Deviation, Decision
Reference, and Human Outcome rules. `internal/workflow/`, `internal/evidence/`,
`internal/claims/` carry Obligations, protected Decisions, budgets, Stops,
Attention, freshness, and review Claim mode; the closed grant-action allowlist
is `workflow/activation.go`, now in the scope of tasks 4, 5, 8, 9, 10
(DEC-167/173). `internal/journal/`, `internal/db/`, `migrations/` own 0017,
0018, the fence, effect rules, and parity. `internal/links/` is restricted and
`internal/requirements/` is removed, including its archive semantic rule in
`internal/export`. `internal/export/`, `internal/maintenance/`,
`internal/activity/`, and `cmd/fid/` carry adapters, projections, and removal
convergence.

## Known Risks

Wide-reach change: Obligation generalization (task 2) and Evidence widening
(task 6), each isolated in its own task. Clause without a schema row (DEC-195):
the coverage matrix plus the QA Lead's independent re-derivation caught two
before dispatch (R-02, R-12), both now assertion-backed. Late Human-readability
discovery: mitigated by the deterministic Human-text check in window 13-15.
Removal ordering: task 15 depends on fourteen replacements, and a slip amends
the task rather than deferring a removal. Test-time headroom: the tagged suite
ran 560 s of an implicit 600 s default before this phase adds journeys, now
bounded and budgeted by DEC-227.

## Success / Exit Criteria

- [ ] All 18 tasks complete with green quality gates.
- [ ] All 79 `VAL-NNN` outcomes recorded with cited evidence.
- [ ] Final QA review has no unresolved blocking finding.
- [ ] Human checklist accepted, including VAL-070 to VAL-073.
- [ ] Sign-off recorded and all twelve canonical specs synced.

## Next Planning Checkpoint

Phase 24: multi-Issue Envelope revision, scope conflict, parallel Claims,
partial progress, generic `retry-exhausted`, `triage_events` and V1 Claim/Run
removals.
