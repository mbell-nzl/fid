# Phase 11 Progress

> Phase: 11 - Claims layer (roadmap v0.1d, slice 2 of 4)
> Status: IN PROGRESS, Task 5 pending after ESC-008 resolution.

## Plan

Phase 11 delivers primitive issue claims and transition integration:

- `fid issue claim` with normal acquire, stale takeover, force-acquire, renew,
  TTL, and note.
- `fid issue unclaim` holder release.
- `fid claim list` with stale/by filters and all four output formats.
- Transition integration requiring live caller claims for started-category entry
  and auto-release on exit.
- Claim action rows added to the Phase 10 in-code permissions matrix.
- ESC-008 fixup: claim writes refresh stale projections without nested write-lock
  acquisition before final QA.

No schema migration is planned; the `claims` table/triggers already exist.

## Task Index

| Task | Title | Status |
|---|---|---|
| 1 | Claim matrix rows and domain claim semantics | Complete; Review 1 CONTINUE |
| 2 | `fid issue claim` and `fid issue unclaim` | Complete; Review 2 CONTINUE |
| 3 | `fid claim list` | Complete; Review 3 CONTINUE |
| 4 | Claim and transition integration | Complete; review pending |
| 5 | Claim write projection freshness lock ordering | Pending |

## Validation

Validation contract at `qa/validation-contract.md` now includes VAL-001..015.
VAL-014 is phase-end human dogfood; VAL-015 is the ESC-008 stale-projection
lock-ordering proof. QA plan at `qa/qa-plan.md` requires Review 5 before final
QA.

## Progress History

- Kickoff: drafted VAL-001..014 and Tasks 1-4; re-deferred CF-008/011/013.
- Task 1 complete: claim matrix rows and domain acquire/renew semantics.
- Review 1: CONTINUE with one improvement note about explicit policy-show
  evidence.
- Task 2 complete: `fid issue claim` and `fid issue unclaim` CLI.
- Review 2: CONTINUE with one improvement note about timing-sensitive stale
  takeover CLI fixture.
- Task 3 complete: `fid claim list`.
- Review 3: CONTINUE with one improvement note about fuller gate excerpts.
- Task 4 complete: transition claim gate and started-exit auto-release.
- TRG-001/ESC-008: Task 4 report surfaced claim-write projection refresh lock
  ordering risk.
- DEC-028: Tech Lead resolved ESC-008 as do-now Task 5 with VAL-015.

## Carry-Forward

- CF-008 re-deferred to v0.2.
- CF-011 broad reconciliation re-deferred to v0.1e; Phase 11 claim renewal still
  projects to `activity_log` per PRD-CLAIM-008a.
- CF-013 re-deferred to v0.1e.

## Escalations

- ESC-008 resolved by DEC-028; Task 5 is pending before final QA.
