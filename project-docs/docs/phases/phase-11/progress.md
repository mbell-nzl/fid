# Phase 11 Progress

> Phase: 11 - Claims layer (roadmap v0.1d, slice 2 of 4)
> Status: PLANNED, awaiting human review of validation contract and task list.

## Plan

Phase 11 delivers primitive issue claims and transition integration:

- `fid issue claim` with normal acquire, stale takeover, force-acquire, renew,
  TTL, and note.
- `fid issue unclaim` holder release.
- `fid claim list` with stale/by filters and all four output formats.
- Transition integration requiring live caller claims for started-category entry
  and auto-release on exit.
- Claim action rows added to the Phase 10 in-code permissions matrix.

No schema migration is planned; the `claims` table/triggers already exist.

## Task Index

| Task | Title | Status |
|---|---|---|
| 1 | Claim matrix rows and domain claim semantics | Pending |
| 2 | `fid issue claim` and `fid issue unclaim` | Pending |
| 3 | `fid claim list` | Pending |
| 4 | Claim and transition integration | Pending |

## Validation

Validation contract drafted at `qa/validation-contract.md` with VAL-001..014.
VAL-014 is phase-end human dogfood; all other VALs are deterministic QA Lead or
QA Reviewer evidence.

## Carry-Forward

- CF-008 re-deferred to v0.2.
- CF-011 broad reconciliation re-deferred to v0.1e; Phase 11 claim renewal still
  projects to `activity_log` per PRD-CLAIM-008a.
- CF-013 re-deferred to v0.1e.

## Escalations

None open at kickoff.
