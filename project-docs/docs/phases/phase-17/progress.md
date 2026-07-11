# Phase 17 Progress

> Phase: 17 - Maintenance + stub guard + v0.1 closure (roadmap v0.1e,
> slice 4 of 4)
> Status: APPROVED; awaiting Heavy Mode QA plan and fresh feasibility
> validation. No Developer task is active.

## Plan

Phase 17 delivers `fid maintenance check`, locked `fid maintenance vacuum`,
explicit v0.2 governance command guards, full read-format conformance, and the
v0.1 closure journey/rebuild/traceability evidence.

## Task Index

| Task | Title | Status |
|---|---|---|
| 1 | CF-019 markdown formatter split | Pending |
| 2 | Maintenance grace config contract | Pending |
| 3 | Maintenance structural errors | Pending |
| 4 | Maintenance warnings and info | Pending |
| 5 | `fid maintenance check` CLI + formats | Pending |
| 6 | `fid maintenance vacuum` | Pending |
| 7 | Governance stub-table guards | Pending |
| 8 | v0.1 read-format conformance | Pending |
| 9 | Full v0.1 closure journey | Pending |

## Validation

`qa/validation-contract.md` contains VAL-001..025 at Ready after Human approval
of Heavy Mode and DC-A/DC-B/DC-C. QA plan creation is the next hard prerequisite.

## Progress History

- Planning draft (2026-07-10): Tech Lead drafted the Heavy Mode contract and
  nine-task pack; implementation blocked on Human approval.
- Kickoff approval (2026-07-11): Human approved Heavy Mode and DC-A/DC-B/DC-C.
  DEC-049..052 record the mode, maintenance warning/config contract,
  result/exit semantics, and exact governance guard matrix. Canonical affected
  specs were synchronized. No Developer task started.

## Carry-Forward

- CF-019 scheduled in Task 1.
- CF-008 remains v0.2.
- CF-017 remains v0.2c.
- CF-018 remains at the next `internal/sqlquery` modification; Task 8 stops for
  re-planning if such a modification is required.

## Escalations

None. QA-plan feasibility validation is pending, not an escalation.
