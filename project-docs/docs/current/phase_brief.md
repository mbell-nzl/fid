# Phase Brief - Phase 11 PLANNED (v0.1d slice 2/4)

> **Purpose.** Default execution brief for the current phase. Loaded by every
> agent session. Keep under 150 lines.

## Phase Goal

Phase 11 delivers the **claims layer**: issue claim acquire/renew/force/release,
claim listing, and transition integration. Success means an agent can acquire an
exclusive issue claim with correct TTL behavior, reject conflicts atomically,
list live/stale claims, enter started states only while holding a live claim,
and see the claim auto-released when leaving started work.

Scope refs: PRD-CLAIM-001..008a, PRD-ERR-003, TRD-ARCH-010..011,
`03_TRD_SCHEMA.md` section 11.1, `06_DESIGN_PLAN_CLI.md` claim commands.
The `claims` table/triggers already exist; no schema migration is planned.

## Status

**Planned, awaiting human review.** Current branch: `main`. Phase 11 validation
contract and task specs are drafted under
`project-docs/docs/phases/phase-11/`. Next action: human reviews the ordered task
list and contract, then QA Lead creates the QA plan before Developer work starts.

## Task Index

| Task | Title | Status | Dependencies |
|---|---|---|---|
| 1 | Claim matrix rows and domain claim semantics | Pending | None |
| 2 | `fid issue claim` and `fid issue unclaim` | Pending | Task 1 |
| 3 | `fid claim list` | Pending | Tasks 1-2 |
| 4 | Claim and transition integration | Pending | Tasks 1-2 |

## Validation Contract

Full contract:
`../phases/phase-11/qa/validation-contract.md`.

| VAL ID | Source | Owning Task/Window | Validator | Status |
|---|---|---|---|---|
| VAL-001 | TRD-ARCH-010; schema section 13 | task-1 | QA Lead | Draft |
| VAL-002 | PRD-CLAIM-001 | tasks-1-2 | QA Lead | Draft |
| VAL-003 | PRD-CLAIM-002; PRD-ERR-003 | tasks-1-2 | QA Lead | Draft |
| VAL-004 | PRD-CLAIM-003 | tasks-1-2 | QA Lead | Draft |
| VAL-005 | PRD-CLAIM-004 | tasks-1-2 | QA Lead | Draft |
| VAL-006 | PRD-CLAIM-008a; CF-011 | tasks-1-2 | QA Lead | Draft |
| VAL-007 | roadmap issue unclaim | task-2 | QA Lead | Draft |
| VAL-008 | PRD-CLAIM-007 | task-3 | QA Lead | Draft |
| VAL-009 | PRD-CLAIM-007; read-ordering principle | task-3 | QA Lead | Draft |
| VAL-010 | PRD-CLAIM-008 | tasks-2-4 | QA Lead | Draft |
| VAL-011 | PRD-CLAIM-005 | task-4 | QA Lead | Draft |
| VAL-012 | PRD-CLAIM-006 | task-4 | QA Lead | Draft |
| VAL-013 | TRD-ARCH-011 | task-1 | QA Lead | Draft |
| VAL-014 | Quality plan human dogfood | phase-end | Human | Draft |

## QA Artifacts

| Artifact | Path | Status |
|---|---|---|
| Validation contract | `../phases/phase-11/qa/validation-contract.md` | Draft |
| QA plan | `../phases/phase-11/qa/qa-plan.md` | Missing |
| Mid-phase reviews | `../phases/phase-11/qa/review-M.md` | Pending; recommend review after each task |
| Final review | `../phases/phase-11/qa/final-review.md` | Pending |
| Human checklist | `../phases/phase-11/qa/human-checklist.md` | Pending |

## Active Constraints

- All claim writes are gated by the Phase 10 matrix; Task 1 appends only the
  implemented claim action rows to the in-code matrix.
- `fid claim list --by` is a read filter, not an actor flag (DEC-027).
- Claim renewal must project to `activity_log` now; broader CF-011 remains v0.1e.
- `fid claim list` default ordering uses `claims.rowid` unless a spec changes.
- Each task must reconcile completion-report commit hash/subject in the task
  commit or an immediate `docs(task): reconcile ...` commit.

## Interfaces / Subsystems In Scope

- `internal/claims/` - new domain package for claim semantics.
- `internal/policy/` - claim action rows in the in-code matrix.
- `cmd/fid/issue.go` - `issue claim` and `issue unclaim`.
- `cmd/fid/claim.go` - root `claim list`.
- `internal/issues/transition.go` - live-claim gate and auto-release.

## Known Risks

- **Error precedence.** If claim-gate and workflow-rule failures conflict in
  transition tests, Developer must escalate before freezing public behavior.
- **Validate-before-append.** Every reject path must prove zero new journal
  lines/rows; QA plan should make this explicit.
- **Activity-log scope.** Renewal is in scope per PRD-CLAIM-008a; broad CF-011
  reconciliation stays deferred.

## Success / Exit Criteria

- [ ] All four tasks complete with required gates recorded.
- [ ] VAL-001..014 have PASS or Tech Lead-dispositioned outcomes.
- [ ] Final QA review has no unresolved blocking findings.
- [ ] Human dogfood checklist is complete.
- [ ] Phase sign-off and spec sync verdicts are recorded.

## Next Planning Checkpoint

Phase 12 plans requirements sync plus links after Phase 11 wrap.
