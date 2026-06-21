# Phase Brief - Phase 11 IN PROGRESS (v0.1d slice 2/4)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.

## Phase Goal

Phase 11 delivers the **claims layer**: issue claim acquire/renew/force/release,
claim listing, and transition integration. Success means an agent can acquire an
exclusive issue claim with correct TTL behavior, reject conflicts atomically,
list live/stale claims, enter started states only while holding a live claim,
and see the claim auto-released when leaving started work.

Scope refs: PRD-CLAIM-001..008a, PRD-ERR-003, PRD-OPS-003, PRD-NFR-008,
TRD-ARCH-010..011, `02_TRD.md` §4.3/§4.9, `03_TRD_SCHEMA.md` section 11.1,
and `06_DESIGN_PLAN_CLI.md` claim commands. The `claims` table/triggers already
exist; no schema migration is planned.

## Status

**In progress on branch `phase-11`.** Tasks 1-4 have completion reports and
Reviews 1-3 are CONTINUE. ESC-008 is resolved by DEC-028 with do-now Task 5
before final QA. Next action: Developer completes Task 5, then QA Reviewer writes
Review 4 for Task 4 and Review 5 for Task 5 before final QA.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Claim matrix rows and domain claim semantics | Complete; Review 1 CONTINUE | None |
| 2 | `fid issue claim` and `fid issue unclaim` | Complete; Review 2 CONTINUE | Task 1 |
| 3 | `fid claim list` | Complete; Review 3 CONTINUE | Tasks 1-2 |
| 4 | Claim and transition integration | Complete; review pending | Tasks 1-2 |
| 5 | Claim write projection freshness lock ordering | Pending | Tasks 1-4; DEC-028 |

## Validation Contract

Full contract:
`../phases/phase-11/qa/validation-contract.md`.

| VAL ID | Source | Owning Task/Window | Validator | Status |
|---|---|---|---|---|
| VAL-001 | TRD-ARCH-010; schema section 13 | task-1 | QA Lead | On track |
| VAL-002 | PRD-CLAIM-001 | tasks-1-2 | QA Lead | On track |
| VAL-003 | PRD-CLAIM-002; PRD-ERR-003 | tasks-1-2 | QA Lead | On track |
| VAL-004 | PRD-CLAIM-003 | tasks-1-2 | QA Lead | On track |
| VAL-005 | PRD-CLAIM-004 | tasks-1-2 | QA Lead | On track |
| VAL-006 | PRD-CLAIM-008a; CF-011 | tasks-1-2 | QA Lead | On track |
| VAL-007 | roadmap issue unclaim | task-2 | QA Lead | On track |
| VAL-008 | PRD-CLAIM-007 | task-3 | QA Lead | On track |
| VAL-009 | PRD-CLAIM-007; read-ordering principle | task-3 | QA Lead | On track |
| VAL-010 | PRD-CLAIM-008 | tasks-2-4 | QA Lead | Pending Review 4 |
| VAL-011 | PRD-CLAIM-005 | task-4 | QA Lead | Pending Review 4 |
| VAL-012 | PRD-CLAIM-006 | task-4 | QA Lead | Pending Review 4 |
| VAL-013 | TRD-ARCH-011 | tasks-1,5 | QA Lead | On track; Task 5 regression pending |
| VAL-014 | Quality plan human dogfood | phase-end | Human | Pending |
| VAL-015 | ESC-008; PRD-OPS-003; PRD-NFR-008 | task-5 | QA Lead | Pending |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-11/qa/validation-contract.md` | Updated for VAL-015 |
| QA plan | `../phases/phase-11/qa/qa-plan.md` | Updated for Task 5 |
| Reviews 1-3 | `../phases/phase-11/qa/review-1.md`..`review-3.md` | CONTINUE |
| Review 4 | `../phases/phase-11/qa/review-4.md` | Pending |
| Review 5 | `../phases/phase-11/qa/review-5.md` | Pending |
| Final review | `../phases/phase-11/qa/final-review.md` | Pending |
| Human checklist | `../phases/phase-11/qa/human-checklist.md` | Pending |

## Active Constraints

- All claim writes are gated by the Phase 10 matrix; Task 1 added only the
  implemented claim action rows to the in-code matrix.
- `fid claim list --by` is a read filter, not an actor flag (DEC-027).
- Claim renewal must project to `activity_log` now; broader CF-011 remains v0.1e.
- `fid claim list` default ordering uses `claims.rowid` unless a spec changes.
- DEC-028: claim writes must refresh stale projections without nested
  `.fid/write.lock` acquisition before final QA continues.
- Each task must reconcile completion-report commit hash/subject in the task
  commit or an immediate `docs(task): reconcile ...` commit.

## Interfaces / Subsystems In Scope

- `internal/claims/` - claim semantics and Task 5 lock-ordering fix.
- `internal/journal/` - lock-aware projection freshness if needed for Task 5.
- `internal/policy/` - claim action rows in the in-code matrix.
- `cmd/fid/issue.go` - `issue claim`, `issue unclaim`, and `issue transition`.
- `cmd/fid/claim.go` - root `claim list`.
- `internal/issues/transition.go` - live-claim gate and auto-release.

## Known Risks

- **ESC-008 regression.** Task 5 must prove stale-projection claim writes do not
  fail through nested lock contention or validate against stale SQLite.
- **Validate-before-append.** Every reject path must prove zero new journal
  lines/rows, including Task 5 stale-projection conflicts.
- **Activity-log scope.** Renewal is in scope per PRD-CLAIM-008a; broad CF-011
  reconciliation stays deferred.

## Success / Exit Criteria

- [ ] All five tasks complete with required gates recorded.
- [ ] VAL-001..015 have PASS or Tech Lead-dispositioned outcomes.
- [ ] Reviews 4-5 and final QA review have no unresolved blocking findings.
- [ ] Human dogfood checklist is complete.
- [ ] Phase sign-off and spec sync verdicts are recorded.

## Next Planning Checkpoint

Phase 12 plans requirements sync plus links after Phase 11 wrap.
