# Phase 17 Progress

> Phase: 17 - Maintenance + stub guard + v0.1 closure (roadmap v0.1e,
> slice 4 of 4)
> Status: APPROVED / BLOCKED AT ENTRY. QA-BLK-01 now has Task 0 remediation;
> ESC-009 model qualification and a revised Ready QA plan remain prerequisites.
> No Developer task is active.

## Plan

Phase 17 delivers `fid maintenance check`, locked `fid maintenance vacuum`,
explicit v0.2 governance command guards, full read-format conformance, and the
v0.1 closure journey/rebuild/traceability evidence.

## Task Index

| Task | Title | Status |
|---|---|---|
| 0 | Restore `internal/journal` COV-001 coverage | Pending; entry prerequisite |
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

`qa/validation-contract.md` contains VAL-001..025 at Ready. DEC-054 adds Task 0
evidence to VAL-024 without changing its threshold or QA Lead ownership. The QA
plan remains BLOCKED until ESC-009 is resolved and QA Lead revises it.

## Progress History

- Planning draft (2026-07-10): Tech Lead drafted the Heavy Mode contract and
  nine-task pack; implementation blocked on Human approval.
- Kickoff approval (2026-07-11): Human approved Heavy Mode and DC-A/DC-B/DC-C.
  DEC-049..052 record the mode, maintenance warning/config contract,
  result/exit semantics, and exact governance guard matrix. Canonical affected
  specs were synchronized. No Developer task started.
- Entry re-plan (2026-07-11): QA-BLK-01 accepted. DEC-054 adds tests-only Task
  0 plus an entry review; no coverage waiver. QA-BLK-02 is tracked as ESC-009.

## Carry-Forward

- CF-019 scheduled in Task 1.
- CF-008 remains v0.2.
- CF-017 remains v0.2c.
- CF-018 remains at the next `internal/sqlquery` modification; Task 8 stops for
  re-planning if such a modification is required.

## Escalations

- ESC-009: qualify the DEC-053 QA routes with durable Human-scored
  seeded-defect evidence. DEC-055 rejects any routing fallback; Developer
  execution is blocked.
